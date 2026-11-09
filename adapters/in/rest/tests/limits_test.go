// adapters/in/rest/tests/limits_test.go — limite de débit derrière un proxy
// (R-05), routes publiques bornées (R-06), corps de requête bornés (R-07).
package rest_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"myr-core/adapters/in/rest"
	"myr-core/domain/model"
)

// newLimitsMux construit le mux REST complet autour de h.
func newLimitsMux(t *testing.T, h *rest.Handler) http.Handler {
	t.Helper()
	mux, err := rest.NewServer(h, "localhost:0").Handler()
	if err != nil {
		t.Fatalf("Handler(): %v", err)
	}
	return mux
}

// sessionAttempt envoie une tentative de connexion (corps invalide : seule la
// limite de débit, appliquée d'abord, est observée) et retourne le code HTTP.
func sessionAttempt(mux http.Handler, remote string, xff ...string) int {
	r := httptest.NewRequest(http.MethodPost, "/api/identity/session", strings.NewReader("{"))
	r.RemoteAddr = remote
	for _, v := range xff {
		r.Header.Add("X-Forwarded-For", v)
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w.Code
}

func newIdentityHandler() *rest.Handler {
	store := &mockStore{}
	return rest.NewHandler(&mockSvc{}, store, store).WithIdentityService(&mockIdentitySvc{})
}

// / @brief  Un X-Forwarded-For forgé par un client externe ne contourne pas la limite de débit (R-05)
// / @input  11 tentatives depuis 203.0.113.7, chacune avec un X-Forwarded-For différent
// / @expect la 11e est refusée (429)
func TestRateLimit_R05_SpoofedForwardedFor_Ignored(t *testing.T) {
	mux := newLimitsMux(t, newIdentityHandler())
	for i := 0; i < 10; i++ {
		if c := sessionAttempt(mux, "203.0.113.7:4000", "198.51.100."+string(rune('1'+i))); c == http.StatusTooManyRequests {
			t.Fatalf("tentative %d refusée trop tôt", i+1)
		}
	}
	if c := sessionAttempt(mux, "203.0.113.7:4000", "198.51.100.200"); c != http.StatusTooManyRequests {
		t.Errorf("11e tentative : 429 attendu, obtenu %d", c)
	}
}

// / @brief  Derrière un proxy de confiance (boucle locale), le client est la dernière adresse non fiable de X-Forwarded-For (R-05)
// / @input  connexions de 127.0.0.1 ; client A (adresses de gauche forgées) 10 fois, puis client B, puis A
// / @expect B passe ; la 11e tentative de A est refusée (429) malgré les adresses de gauche changeantes
func TestRateLimit_R05_TrustedProxy_UsesForwardedClient(t *testing.T) {
	mux := newLimitsMux(t, newIdentityHandler())
	for i := 0; i < 10; i++ {
		sessionAttempt(mux, "127.0.0.1:5000", "10.9.9."+string(rune('1'+i))+", 203.0.113.10")
	}
	if c := sessionAttempt(mux, "127.0.0.1:5000", "203.0.113.20"); c == http.StatusTooManyRequests {
		t.Error("client B : ne doit pas partager la limite de A")
	}
	if c := sessionAttempt(mux, "127.0.0.1:5000", "1.2.3.4", "203.0.113.10"); c != http.StatusTooManyRequests {
		t.Errorf("client A, 11e tentative : 429 attendu, obtenu %d", c)
	}
}

// / @brief  Les adresses IPv6 sont distinguées entre elles (R-05)
// / @input  11 tentatives depuis [2001:db8::1], puis une depuis [2001:db8::2]
// / @expect la 11e de ::1 est refusée, celle de ::2 passe
func TestRateLimit_R05_IPv6Clients_Distinct(t *testing.T) {
	mux := newLimitsMux(t, newIdentityHandler())
	for i := 0; i < 10; i++ {
		sessionAttempt(mux, "[2001:db8::1]:443")
	}
	if c := sessionAttempt(mux, "[2001:db8::1]:443"); c != http.StatusTooManyRequests {
		t.Errorf("11e tentative IPv6 : 429 attendu, obtenu %d", c)
	}
	if c := sessionAttempt(mux, "[2001:db8::2]:443"); c == http.StatusTooManyRequests {
		t.Error("une autre adresse IPv6 ne doit pas partager la limite")
	}
}

// / @brief  WithTrustedProxies accepte adresses et plages CIDR, et refuse une entrée invalide (R-05)
// / @input  proxy de confiance 10.0.0.0/8 ; client 203.0.113.30 relayé par 10.1.2.3 ; puis l'entrée « proxy »
// / @expect les tentatives relayées comptent pour 203.0.113.30 ; l'entrée invalide est refusée
func TestWithTrustedProxies_R05_CIDR(t *testing.T) {
	h, err := newIdentityHandler().WithTrustedProxies([]string{"10.0.0.0/8", " "})
	if err != nil {
		t.Fatalf("WithTrustedProxies : %v", err)
	}
	mux := newLimitsMux(t, h)
	for i := 0; i < 10; i++ {
		sessionAttempt(mux, "10.1.2.3:80", "203.0.113.30")
	}
	if c := sessionAttempt(mux, "10.4.5.6:80", "203.0.113.30"); c != http.StatusTooManyRequests {
		t.Errorf("même client par un autre proxy : 429 attendu, obtenu %d", c)
	}
	if c := sessionAttempt(mux, "127.0.0.1:80", "203.0.113.40"); c == http.StatusTooManyRequests {
		t.Error("127.0.0.1 n'est plus un proxy de confiance : sa propre limite s'applique, pas celle du client")
	}
	if _, err := newIdentityHandler().WithTrustedProxies([]string{"proxy"}); err == nil {
		t.Error("entrée invalide : erreur attendue")
	}
}

// / @brief  POST /api/identity/enroll est limité en débit (R-06)
// / @input  11 tentatives d'enrôlement depuis la même adresse
// / @expect la 11e est refusée (429)
func TestEnroll_R06_RateLimited(t *testing.T) {
	mux := newLimitsMux(t, newIdentityHandler())
	var c int
	for i := 0; i < 11; i++ {
		r := httptest.NewRequest(http.MethodPost, "/api/identity/enroll", strings.NewReader(`{"name":"a","secret":"x","org_id":"Org1MSP"}`))
		r.RemoteAddr = "203.0.113.50:1"
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		c = w.Code
	}
	if c != http.StatusTooManyRequests {
		t.Errorf("11e tentative : 429 attendu, obtenu %d", c)
	}
}

// / @brief  GET /api/status ne parcourt le ledger qu'une fois par délai de cache (R-06)
// / @input  3 appels successifs à /api/status
// / @expect List appelé une seule fois ; le nombre d'assets est le même dans les 3 réponses
func TestStatus_R06_AssetCountCached(t *testing.T) {
	calls := 0
	svc := &mockSvc{list: func(string) ([]*model.Model3D, error) {
		calls++
		return []*model.Model3D{{ID: "a"}, {ID: "b"}}, nil
	}}
	store := &mockStore{}
	mux := newLimitsMux(t, rest.NewHandler(svc, store, store))
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/status", nil))
		if !strings.Contains(w.Body.String(), `"assets":2`) {
			t.Errorf("appel %d : %s", i+1, w.Body.String())
		}
	}
	if calls != 1 {
		t.Errorf("List appelé %d fois, attendu 1", calls)
	}
}

// / @brief  Les corps de requête sont bornés : JSON, miniature, envoi de fichier (R-07)
// / @input  JSON de 2 Mo ; miniature de 4 Mo ; fichier multipart de 2 Ko avec une limite d'envoi de 1 Ko
// / @expect 413 pour les trois
func TestLimitBodies_R07_TooLarge_Rejected(t *testing.T) {
	store := &mockStore{}
	mux := newLimitsMux(t, rest.NewHandler(&mockSvc{}, store, store).WithMaxUpload(1024))

	big := `{"name":"` + strings.Repeat("x", 2<<20) + `"}`
	thumb := `{"thumbnail":"data:image/png;base64,` + strings.Repeat("A", 4<<20) + `"}`
	var form bytes.Buffer
	mw := multipart.NewWriter(&form)
	fw, _ := mw.CreateFormFile("file", "piece.stl")
	_, _ = fw.Write(bytes.Repeat([]byte("s"), 2048))
	_ = mw.Close()

	cases := []struct {
		name, path, ctype string
		body              []byte
	}{
		{"JSON", "/api/connections", "application/json", []byte(big)},
		{"miniature", "/api/components/a1/thumbnail", "application/json", []byte(thumb)},
		{"fichier", "/api/components", mw.FormDataContentType(), form.Bytes()},
	}
	for _, c := range cases {
		r := httptest.NewRequest(http.MethodPost, c.path, bytes.NewReader(c.body))
		r.Header.Set("Content-Type", c.ctype)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("%s : 413 attendu, obtenu %d", c.name, w.Code)
		}
	}
}
