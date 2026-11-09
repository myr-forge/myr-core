// adapters/in/rest/limits.go — adresse du client derrière un proxy de
// confiance, et bornes de taille des corps de requête.
package rest

import (
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strings"
)

const (
	// maxJSONBody borne un corps JSON ordinaire.
	maxJSONBody = 1 << 20
	// maxThumbnailBody borne l'envoi d'une miniature : une image de 2 Mo
	// (model.MaxThumbnailBytes) encodée en base64 dans une data URL.
	maxThumbnailBody = 3 << 20
	// DefaultMaxUpload borne un envoi de fichier (multipart) : le fichier CAO
	// n'est reçu que pour en calculer l'empreinte (RM42).
	DefaultMaxUpload = 256 << 20
)

// loopbackProxies sont les proxys de confiance par défaut : un proxy inverse
// (nginx) installé sur le même serveur que l'API.
var loopbackProxies = []netip.Prefix{
	netip.MustParsePrefix("127.0.0.0/8"),
	netip.MustParsePrefix("::1/128"),
}

// WithTrustedProxies remplace les proxys de confiance (adresses IP ou plages
// CIDR) dont l'en-tête X-Forwarded-For désigne le client : ceux de la boucle
// locale par défaut, l'adresse d'un proxy distant quand myr-web et l'API sont
// sur des serveurs différents. Une liste vide n'accorde la confiance à aucun.
func (h *Handler) WithTrustedProxies(entries []string) (*Handler, error) {
	prefixes := make([]netip.Prefix, 0, len(entries))
	for _, e := range entries {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}
		if p, err := netip.ParsePrefix(e); err == nil {
			prefixes = append(prefixes, p.Masked())
			continue
		}
		a, err := netip.ParseAddr(e)
		if err != nil {
			return h, fmt.Errorf("proxy de confiance invalide %q : adresse IP ou plage CIDR attendue", e)
		}
		prefixes = append(prefixes, netip.PrefixFrom(a.Unmap(), a.Unmap().BitLen()))
	}
	h.trustedProxies = prefixes
	return h, nil
}

// WithMaxUpload borne la taille d'un envoi de fichier (octets ; DefaultMaxUpload
// si n <= 0).
func (h *Handler) WithMaxUpload(n int64) *Handler {
	if n <= 0 {
		n = DefaultMaxUpload
	}
	h.maxUpload = n
	return h
}

// trusted indique si ip appartient à un proxy de confiance.
func (h *Handler) trusted(ip string) bool {
	a, err := netip.ParseAddr(ip)
	if err != nil {
		return false
	}
	a = a.Unmap()
	for _, p := range h.trustedProxies {
		if p.Contains(a) {
			return true
		}
	}
	return false
}

// clientIP retourne l'adresse du client, clé de la limite de débit. Une
// connexion venant d'un proxy de confiance est attribuée à la dernière adresse
// de X-Forwarded-For qui n'est pas elle-même un proxy de confiance — les
// adresses plus à gauche sont fournies par le client et ne sont jamais crues.
// Toute autre connexion est attribuée à son adresse réseau.
func (h *Handler) clientIP(r *http.Request) string {
	host := r.RemoteAddr
	if hst, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		host = hst
	}
	if !h.trusted(host) {
		return host
	}
	var hops []string
	for _, v := range r.Header.Values("X-Forwarded-For") {
		hops = append(hops, strings.Split(v, ",")...)
	}
	for i := len(hops) - 1; i >= 0; i-- {
		if ip := strings.TrimSpace(hops[i]); ip != "" && !h.trusted(ip) {
			return ip
		}
	}
	return host
}

// limitBodies borne le corps de chaque requête : maxUpload pour un envoi de
// fichier (multipart), maxThumbnailBody pour une miniature, maxJSONBody sinon.
// Une longueur déclarée au-delà est refusée d'emblée (413) ; un corps de
// longueur non déclarée est interrompu à la borne.
func (h *Handler) limitBodies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limit := int64(maxJSONBody)
		switch {
		case strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/"):
			limit = h.maxUpload
		case strings.HasSuffix(r.URL.Path, "/thumbnail"):
			limit = maxThumbnailBody
		}
		if r.ContentLength > limit {
			jsonError(w, fmt.Sprintf("corps de requête trop volumineux (%d octets au plus)", limit), http.StatusRequestEntityTooLarge)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, limit)
		next.ServeHTTP(w, r)
	})
}
