// adapters/in/rest/tests/publicread_test.go — lecture publique du catalogue
// sans compte, quand le réseau l'autorise (RM49).
package rest_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"myr-core/adapters/in/rest"
	"myr-core/domain/model"
	"myr-core/domain/network"
)

// catalogueMux construit le mux REST d'un réseau dont la lecture publique est
// ouverte (open) ou fermée, avec un catalogue d'un asset soumis « pub1 » et
// d'un brouillon « draft1 ».
func catalogueMux(t *testing.T, open bool) http.Handler {
	t.Helper()
	assets := map[string]*model.Model3D{
		"pub1":   {ID: "pub1", Name: "vis publiée", Status: model.ModuleSubmitted},
		"draft1": {ID: "draft1", Name: "brouillon privé", Status: model.ModuleDraft},
	}
	svc := &mockSvc{
		list: func(string) ([]*model.Model3D, error) {
			return []*model.Model3D{assets["pub1"], assets["draft1"]}, nil
		},
		get: func(id string) (*model.Model3D, error) {
			if m, ok := assets[id]; ok {
				return m, nil
			}
			return nil, model.ErrNotFound
		},
	}
	netSvc := &mockNetworkSvc{getActive: func() (*network.NetworkProfile, error) {
		return &network.NetworkProfile{Name: "net", AllowAutoGuest: open}, nil
	}}
	store := &mockStore{}
	h := rest.NewHandler(svc, store).WithIdentityService(&mockIdentitySvc{}).WithNetworkService(netSvc)
	mux, err := rest.NewServer(h, "localhost:0").Handler()
	if err != nil {
		t.Fatalf("Handler(): %v", err)
	}
	return mux
}

func visit(mux http.Handler, method, path, token string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, nil)
	if token != "" {
		r.Header.Set("X-Myr-Token", token)
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	return w
}

// / @brief  Réseau ouvert : un visiteur sans jeton lit le catalogue publié, jamais les brouillons (RM49)
// / @input  GET /api/components, /api/components/pub1, /api/components/draft1, /api/licenses sans jeton
// / @expect liste avec pub1 sans draft1 ; pub1 200 ; draft1 404 ; licences 200
func TestPublicRead_RM49_Open_PublishedOnly(t *testing.T) {
	mux := catalogueMux(t, true)
	w := visit(mux, http.MethodGet, "/api/components", "")
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "pub1") || strings.Contains(w.Body.String(), "draft1") {
		t.Errorf("liste visiteur : %d %s", w.Code, w.Body.String())
	}
	if c := visit(mux, http.MethodGet, "/api/components/pub1", "").Code; c != http.StatusOK {
		t.Errorf("asset publié : 200 attendu, obtenu %d", c)
	}
	if c := visit(mux, http.MethodGet, "/api/components/draft1", "").Code; c != http.StatusNotFound {
		t.Errorf("brouillon : 404 attendu, obtenu %d", c)
	}
	if c := visit(mux, http.MethodGet, "/api/components/draft1/interfaces", "").Code; c != http.StatusNotFound {
		t.Errorf("sous-ressource d'un brouillon : 404 attendu, obtenu %d", c)
	}
	if c := visit(mux, http.MethodGet, "/api/licenses", "").Code; c != http.StatusOK {
		t.Errorf("licences : 200 attendu, obtenu %d", c)
	}
}

// / @brief  Un visiteur n'écrit jamais, ne lit pas les routes hors catalogue, et un jeton invalide reste refusé (RM49)
// / @input  réseau ouvert : POST /api/components sans jeton ; GET /api/connections sans jeton ; GET /api/components avec un jeton inconnu
// / @expect 401 dans les trois cas
func TestPublicRead_RM49_WritesAndOtherRoutes_Refused(t *testing.T) {
	mux := catalogueMux(t, true)
	if c := visit(mux, http.MethodPost, "/api/components", "").Code; c != http.StatusUnauthorized {
		t.Errorf("écriture visiteur : 401 attendu, obtenu %d", c)
	}
	if c := visit(mux, http.MethodGet, "/api/connections", "").Code; c != http.StatusUnauthorized {
		t.Errorf("route hors catalogue : 401 attendu, obtenu %d", c)
	}
	if c := visit(mux, http.MethodGet, "/api/components", "jeton-inconnu").Code; c != http.StatusUnauthorized {
		t.Errorf("jeton invalide : 401 attendu, obtenu %d", c)
	}
}

// / @brief  Réseau fermé à la lecture publique : le catalogue exige une session (RM49)
// / @input  réseau AllowAutoGuest=false, GET /api/components sans jeton, puis avec une session
// / @expect 401 sans jeton ; avec une session, la liste inclut aussi les brouillons
func TestPublicRead_RM49_Closed_SessionRequired(t *testing.T) {
	mux := catalogueMux(t, false)
	if c := visit(mux, http.MethodGet, "/api/components", "").Code; c != http.StatusUnauthorized {
		t.Errorf("sans jeton : 401 attendu, obtenu %d", c)
	}
	w := visit(mux, http.MethodGet, "/api/components", loginToken(t, mux))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "draft1") {
		t.Errorf("avec session : %d %s", w.Code, w.Body.String())
	}
}
