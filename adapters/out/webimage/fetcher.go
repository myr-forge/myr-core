// adapters/out/webimage/fetcher.go — adapter de récupération d'image représentative
// (balise og:image) depuis une page web, pour le port model.OGImageFetcher.
package webimage

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"regexp"
	"strings"
	"syscall"
	"time"
)

var reOGImage = regexp.MustCompile(
	`(?i)<meta[^>]+property=["']og:image["'][^>]+content=["']([^"']+)["']|` +
		`<meta[^>]+content=["']([^"']+)["'][^>]+property=["']og:image["']`,
)

const (
	fetchTimeout  = 10 * time.Second
	pageMaxBytes  = 2 << 20 // 2 Mo
	imageMaxBytes = 5 << 20 // 5 Mo
	maxRedirects  = 3
)

// errForbiddenDestination signale une URL ou une adresse que le Fetcher refuse
// de contacter : la page et l'image sont désignées par l'utilisateur
// (Model3D.Links) puis par la page elle-même, et ne doivent jamais permettre
// d'atteindre un service interne au serveur (SSRF).
var errForbiddenDestination = errors.New("destination interdite : seules les adresses publiques en http/https sont autorisées")

// nonPublicPrefixes complète les catégories de netip.Addr (privée, loopback,
// lien local…) avec les plages qui ne désignent pas une machine publique
// d'Internet ou qui encapsulent une IPv4 arbitraire, donc potentiellement interne.
var nonPublicPrefixes = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"),      // « ce réseau »
	netip.MustParsePrefix("100.64.0.0/10"),  // NAT opérateur, métadonnées de certains clouds
	netip.MustParsePrefix("192.0.0.0/24"),   // affectations IETF
	netip.MustParsePrefix("198.18.0.0/15"),  // tests d'interconnexion
	netip.MustParsePrefix("240.0.0.0/4"),    // réservé, diffusion
	netip.MustParsePrefix("64:ff9b::/96"),   // NAT64
	netip.MustParsePrefix("64:ff9b:1::/48"), // NAT64 local
	netip.MustParsePrefix("2001::/32"),      // Teredo
	netip.MustParsePrefix("2002::/16"),      // 6to4
}

// Fetcher implémente model.OGImageFetcher via des requêtes HTTP directes,
// limitées aux adresses publiques d'Internet (voir errForbiddenDestination).
type Fetcher struct {
	// checkAddr valide chaque adresse "ip:port" effectivement contactée,
	// redirections comprises — nil = checkPublicAddr. Remplacé uniquement en test.
	checkAddr func(address string) error
}

func New() *Fetcher { return &Fetcher{} }

// FetchOGImage récupère l'image og:image d'une page web et la retourne en data URL base64.
func (f *Fetcher) FetchOGImage(pageURL string) (string, error) {
	base, err := parseFetchURL(pageURL)
	if err != nil {
		return "", err
	}
	client := f.client()

	body, _, err := get(client, base.String(), pageMaxBytes)
	if err != nil {
		return "", fmt.Errorf("chargement page: %w", err)
	}

	matches := reOGImage.FindSubmatch(body)
	imgRef := ""
	if len(matches) >= 3 {
		if len(matches[1]) > 0 {
			imgRef = string(matches[1])
		} else {
			imgRef = string(matches[2])
		}
	}
	if imgRef == "" {
		return "", fmt.Errorf("og:image introuvable sur %s", pageURL)
	}
	ref, err := url.Parse(imgRef)
	if err != nil {
		return "", fmt.Errorf("URL og:image invalide: %w", err)
	}
	// ResolveReference couvre les URL absolues, relatives et sans schéma (//hôte/…).
	imgURL, err := parseFetchURL(base.ResolveReference(ref).String())
	if err != nil {
		return "", err
	}

	imgData, contentType, err := get(client, imgURL.String(), imageMaxBytes)
	if err != nil {
		return "", fmt.Errorf("chargement image: %w", err)
	}
	ct, err := imageMediaType(contentType, imgData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("data:%s;base64,%s", ct, base64.StdEncoding.EncodeToString(imgData)), nil
}

// client construit un client HTTP qui n'atteint que des adresses autorisées :
// le contrôle a lieu à l'établissement de chaque connexion, après résolution
// DNS — il couvre donc aussi les redirections et le DNS rebinding.
func (f *Fetcher) client() *http.Client {
	check := f.checkAddr
	if check == nil {
		check = checkPublicAddr
	}
	dialer := &net.Dialer{
		Timeout: fetchTimeout,
		Control: func(_, address string, _ syscall.RawConn) error { return check(address) },
	}
	return &http.Client{
		Timeout: fetchTimeout,
		Transport: &http.Transport{
			Proxy:               nil, // un proxy contournerait le contrôle des adresses
			DialContext:         dialer.DialContext,
			TLSHandshakeTimeout: fetchTimeout,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return fmt.Errorf("plus de %d redirections", maxRedirects)
			}
			_, err := parseFetchURL(req.URL.String())
			return err
		},
	}
}

// get télécharge rawURL (au plus maxBytes) et retourne le corps et le Content-Type.
func get(client *http.Client, rawURL string, maxBytes int64) ([]byte, string, error) {
	resp, err := client.Get(rawURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, "", fmt.Errorf("%s : HTTP %d", rawURL, resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxBytes))
	if err != nil {
		return nil, "", err
	}
	return body, resp.Header.Get("Content-Type"), nil
}

// parseFetchURL n'accepte qu'une URL http(s) absolue avec un hôte.
func parseFetchURL(rawURL string) (*url.URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("URL invalide: %w", err)
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" {
		return nil, fmt.Errorf("%w (%s)", errForbiddenDestination, rawURL)
	}
	return u, nil
}

// checkPublicAddr refuse toute adresse qui n'est pas une adresse publique
// d'Internet, ainsi que les adresses du serveur lui-même (son IP publique
// donnerait accès à ses propres services exposés).
func checkPublicAddr(address string) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}
	ip, err := netip.ParseAddr(host)
	if err != nil || !isPublicAddr(ip) || isOwnAddr(ip) {
		return fmt.Errorf("%w (%s)", errForbiddenDestination, host)
	}
	return nil
}

// isPublicAddr indique si ip désigne une machine publique d'Internet.
func isPublicAddr(ip netip.Addr) bool {
	ip = ip.Unmap()
	if !ip.IsGlobalUnicast() || ip.IsPrivate() {
		return false
	}
	for _, p := range nonPublicPrefixes {
		if p.Contains(ip) {
			return false
		}
	}
	return true
}

// isOwnAddr indique si ip est attribuée à une interface du serveur. En cas
// d'erreur de lecture des interfaces, l'adresse est refusée (sûreté par défaut).
func isOwnAddr(ip netip.Addr) bool {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return true
	}
	ip = ip.Unmap()
	for _, a := range addrs {
		if n, ok := a.(*net.IPNet); ok {
			if own, ok := netip.AddrFromSlice(n.IP); ok && own.Unmap() == ip {
				return true
			}
		}
	}
	return false
}

// imageMediaType retourne le type MIME de l'image téléchargée (en-tête, sinon
// détection sur le contenu) et refuse tout contenu qui n'est pas une image.
// SVG est exclu : c'est un document actif (scripts) — la miniature doit rester
// une image matricielle inerte une fois servie en data URL au client.
func imageMediaType(contentType string, data []byte) (string, error) {
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	mt, _, err := mime.ParseMediaType(contentType)
	if err != nil || !strings.HasPrefix(mt, "image/") || mt == "image/svg+xml" {
		return "", fmt.Errorf("le contenu récupéré n'est pas une image (%s)", contentType)
	}
	return mt, nil
}
