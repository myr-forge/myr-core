// adapters/out/webimage/fetcher_test.go
// Test boîte blanche (package webimage) : le contrôle des adresses contactées
// (Fetcher.checkAddr, checkPublicAddr) n'est pas exporté — même raison que
// l'exception des tests de adapters/in/cli/.
package webimage

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"sync/atomic"
	"testing"
)

// pngHeader suffit à http.DetectContentType pour reconnaître une image PNG.
var pngHeader = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR")

// countingServer démarre un serveur de test et compte les requêtes reçues.
func countingServer(t *testing.T, h http.HandlerFunc) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		h(w, r)
	}))
	t.Cleanup(srv.Close)
	return srv, &hits
}

// pageWithOGImage sert une page HTML dont la balise og:image vaut imgURL.
func pageWithOGImage(imgURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<html><head><meta property="og:image" content="%s"></head></html>`, imgURL)
	}
}

// allowOnly autorise uniquement les serveurs de test donnés — les autres
// adresses passent par le contrôle de production (checkPublicAddr).
func allowOnly(servers ...*httptest.Server) func(string) error {
	return func(address string) error {
		for _, s := range servers {
			if s.Listener.Addr().String() == address {
				return nil
			}
		}
		return checkPublicAddr(address)
	}
}

// / @brief  isPublicAddr n'accepte que les adresses publiques d'Internet
// / @input  adresses loopback, privées, lien local (métadonnées cloud), CGNAT, NAT64, 6to4, IPv4 mappée, multicast, publiques
// / @expect false pour toute adresse non publique, true pour 8.8.8.8 et 2001:4860:4860::8888
func TestIsPublicAddr_Ranges(t *testing.T) {
	cases := map[string]bool{
		"127.0.0.1":            false,
		"10.1.2.3":             false,
		"172.16.0.1":           false,
		"192.168.16.5":         false,
		"169.254.169.254":      false,
		"100.100.100.200":      false,
		"0.0.0.0":              false,
		"0.1.2.3":              false,
		"255.255.255.255":      false,
		"224.0.0.1":            false,
		"::1":                  false,
		"::":                   false,
		"fc00::1":              false,
		"fe80::1":              false,
		"::ffff:127.0.0.1":     false,
		"64:ff9b::7f00:1":      false,
		"2002:7f00:1::":        false,
		"8.8.8.8":              true,
		"2001:4860:4860::8888": true,
	}
	for s, want := range cases {
		if got := isPublicAddr(netip.MustParseAddr(s)); got != want {
			t.Errorf("isPublicAddr(%s) = %v, want %v", s, got, want)
		}
	}
}

// / @brief  FetchOGImage refuse une page hébergée sur une adresse interne, sans la contacter
// / @input  page de test sur 127.0.0.1, Fetcher de production (New)
// / @expect erreur errForbiddenDestination, aucune requête reçue par le serveur
func TestFetchOGImage_LoopbackPage_Rejected(t *testing.T) {
	page, hits := countingServer(t, pageWithOGImage("https://cdn.example/img.png"))

	_, err := New().FetchOGImage(page.URL)
	if !errors.Is(err, errForbiddenDestination) {
		t.Fatalf("err = %v, want errForbiddenDestination", err)
	}
	if n := hits.Load(); n != 0 {
		t.Errorf("le serveur interne a reçu %d requête(s), want 0", n)
	}
}

// / @brief  FetchOGImage refuse une og:image pointant vers un service interne
// / @input  page autorisée dont l'og:image désigne un second serveur sur 127.0.0.1 (non autorisé)
// / @expect erreur errForbiddenDestination, aucune requête reçue par le service interne
func TestFetchOGImage_ImageOnInternalHost_Rejected(t *testing.T) {
	internal, hits := countingServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(pngHeader)
	})
	page, _ := countingServer(t, pageWithOGImage(internal.URL+"/_all_dbs"))

	f := &Fetcher{checkAddr: allowOnly(page)}
	_, err := f.FetchOGImage(page.URL)
	if !errors.Is(err, errForbiddenDestination) {
		t.Fatalf("err = %v, want errForbiddenDestination", err)
	}
	if n := hits.Load(); n != 0 {
		t.Errorf("le service interne a reçu %d requête(s), want 0", n)
	}
}

// / @brief  FetchOGImage refuse une redirection vers un service interne
// / @input  page autorisée répondant 302 vers un second serveur sur 127.0.0.1 (non autorisé)
// / @expect erreur errForbiddenDestination, aucune requête reçue par le service interne
func TestFetchOGImage_RedirectToInternal_Rejected(t *testing.T) {
	internal, hits := countingServer(t, pageWithOGImage("https://cdn.example/img.png"))
	page, _ := countingServer(t, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, internal.URL, http.StatusFound)
	})

	f := &Fetcher{checkAddr: allowOnly(page)}
	_, err := f.FetchOGImage(page.URL)
	if !errors.Is(err, errForbiddenDestination) {
		t.Fatalf("err = %v, want errForbiddenDestination", err)
	}
	if n := hits.Load(); n != 0 {
		t.Errorf("le service interne a reçu %d requête(s), want 0", n)
	}
}

// / @brief  FetchOGImage n'accepte que http et https, pour la page comme pour l'og:image
// / @input  liens file://, ftp://, gopher:// ; page autorisée dont l'og:image vaut file:///etc/passwd
// / @expect erreur errForbiddenDestination dans chaque cas
func TestFetchOGImage_NonHTTPScheme_Rejected(t *testing.T) {
	for _, link := range []string{"file:///etc/passwd", "ftp://example.com/x", "gopher://example.com/x", "/relatif"} {
		if _, err := New().FetchOGImage(link); !errors.Is(err, errForbiddenDestination) {
			t.Errorf("FetchOGImage(%q) err = %v, want errForbiddenDestination", link, err)
		}
	}

	page, _ := countingServer(t, pageWithOGImage("file:///etc/passwd"))
	f := &Fetcher{checkAddr: allowOnly(page)}
	if _, err := f.FetchOGImage(page.URL); !errors.Is(err, errForbiddenDestination) {
		t.Errorf("og:image file:// : err = %v, want errForbiddenDestination", err)
	}
}

// / @brief  FetchOGImage refuse un contenu qui n'est pas une image matricielle
// / @input  page et serveur d'image autorisés, l'image répondant en JSON, en SVG, puis sans Content-Type avec un corps texte
// / @expect une erreur dans chaque cas, jamais de data URL
func TestFetchOGImage_NonImageContent_Rejected(t *testing.T) {
	for _, tc := range []struct{ contentType, body string }{
		{"application/json", `{"secret":"x"}`},
		{"image/svg+xml", `<svg xmlns="http://www.w3.org/2000/svg"><script>alert(1)</script></svg>`},
		{"", "texte brut"},
	} {
		img, _ := countingServer(t, func(w http.ResponseWriter, _ *http.Request) {
			w.Header()["Content-Type"] = []string{tc.contentType}
			w.Write([]byte(tc.body))
		})
		page, _ := countingServer(t, pageWithOGImage(img.URL+"/img"))

		f := &Fetcher{checkAddr: allowOnly(page, img)}
		if dataURL, err := f.FetchOGImage(page.URL); err == nil {
			t.Errorf("Content-Type %q : data URL %q produite, want une erreur", tc.contentType, dataURL)
		}
	}
}

// / @brief  FetchOGImage retourne l'og:image d'une page autorisée en data URL, URL relative résolue
// / @input  page autorisée dont l'og:image est relative (/img.png), servie en image/png
// / @expect data URL "data:image/png;base64,…", sans erreur
func TestFetchOGImage_AllowedImage_OK(t *testing.T) {
	srv, _ := countingServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/img.png" {
			w.Header().Set("Content-Type", "image/png")
			w.Write(pngHeader)
			return
		}
		pageWithOGImage("/img.png")(w, r)
	})

	f := &Fetcher{checkAddr: allowOnly(srv)}
	dataURL, err := f.FetchOGImage(srv.URL + "/produit")
	if err != nil {
		t.Fatalf("FetchOGImage: %v", err)
	}
	if !strings.HasPrefix(dataURL, "data:image/png;base64,") {
		t.Errorf("data URL = %q, want préfixe data:image/png;base64,", dataURL)
	}
}
