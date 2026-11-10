// adapters/in/rest/server.go — adaptateur entrant HTTP (API REST uniquement)
package rest

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	rbac "myr-core/domain/role"
)

// Server expose l'API REST sur un port local. Aucune interface graphique
// n'est servie ici — le GUI est un dépôt externe, client de cette API.
type Server struct {
	handler *Handler
	addr    string
}

func NewServer(h *Handler, addr string) *Server {
	return &Server{handler: h, addr: addr}
}

// legacyModulesAlias adapte les anciennes routes /api/modules/* vers les
// handlers unifiés /api/components/* (ADR-11, specs/3-Conception/Conception_intro.md
// — fenêtre de transition pour myr-web, durée non tranchée). Le corps de
// réponse suit désormais la forme componentDTO, plus moduleDTO — attendu, la
// fusion porte justement sur l'unification du contrat, pas sur la préservation
// de l'ancienne forme derrière l'ancienne URL.
func legacyModulesAlias(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Deprecation", "true")
		w.Header().Set("Link", `</api/components`+strings.TrimPrefix(r.URL.Path, "/api/modules")+`>; rel="successor-version"`)
		r.URL.Path = "/api/components" + strings.TrimPrefix(r.URL.Path, "/api/modules")
		next(w, r)
	}
}

// secureHeaders ajoute les en-têtes de sécurité HTTP à chaque réponse.
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'none'")
		next.ServeHTTP(w, r)
	})
}

// Handler construit et retourne le http.Handler complet de l'API REST.
func (s *Server) Handler() (http.Handler, error) {
	return s.buildMux()
}

// Start démarre le serveur et bloque jusqu'à une erreur (usage dev/simple).
func (s *Server) Start() error {
	mux, err := s.buildMux()
	if err != nil {
		return err
	}
	log.Printf("Serveur REST démarré sur http://%s\n", s.addr)
	return newHTTPServer(s.addr, mux).ListenAndServe()
}

// Délais du serveur HTTP (R-07) : un client qui envoie ses en-têtes ou son
// corps au compte-gouttes ne bloque pas une connexion indéfiniment. Les délais
// de lecture et d'écriture laissent passer un envoi de fichier de taille
// maximale sur un lien lent, et une transaction Fabric (ENF02).
const (
	readHeaderTimeout = 10 * time.Second
	readTimeout       = 10 * time.Minute
	writeTimeout      = 15 * time.Minute
	idleTimeout       = 2 * time.Minute
)

// newHTTPServer retourne le serveur HTTP de l'API, avec ses délais.
func newHTTPServer(addr string, h http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
}

// StartWithShutdown démarre le serveur et se termine proprement quand ctx est annulé.
// À utiliser en production avec signal.NotifyContext.
func (s *Server) StartWithShutdown(ctx context.Context) error {
	mux, err := s.buildMux()
	if err != nil {
		return err
	}
	srv := newHTTPServer(s.addr, mux)
	log.Printf("Serveur REST démarré sur http://%s\n", s.addr)

	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()

	select {
	case <-ctx.Done():
		log.Println("Arrêt propre du serveur HTTP…")
		shutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if shutErr := srv.Shutdown(shutCtx); shutErr != nil {
			log.Printf("Erreur shutdown : %v", shutErr)
		}
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func (s *Server) buildMux() (http.Handler, error) {
	mux := http.NewServeMux()

	// Routes d'identité — ouvertes (point d'entrée, pas d'auth requise)
	mux.HandleFunc("/api/identity/policy", s.handler.handleIdentityPolicy)
	mux.HandleFunc("/api/identity/enroll", s.handler.handleIdentityEnroll)
	mux.HandleFunc("/api/identity/session", s.handler.handleIdentitySession)
	mux.HandleFunc("/api/identity/request", s.handler.handleIdentityRequest)

	// Routes d'identité — protégées
	mux.HandleFunc("/api/identity/requests", s.handler.requireRole(rbac.PermIdentityAdmin, s.handler.handleIdentityRequests))
	mux.HandleFunc("/api/identity/requests/", s.handler.requireRole(rbac.PermIdentityAdmin, s.handler.handleIdentityRequestDecision))
	mux.HandleFunc("/api/identity/accounts", s.handler.requireRole(rbac.PermIdentityAdmin, s.handler.handleIdentityAccounts))
	mux.HandleFunc("/api/identity/wallets", s.handler.requireAuth(s.handler.handleIdentityWallets))
	mux.HandleFunc("/api/identity/status", s.handler.requireAuth(s.handler.handleIdentityStatus))
	mux.HandleFunc("/api/identity/roles/request", s.handler.requireAuth(s.handler.handleIdentityRoleRequest))

	// Routes de données — lecture : reader+, écriture (POST/PUT/PATCH/DELETE) : contributor+.
	// catalogue : lecture du catalogue, ouverte aux visiteurs sans compte
	// quand le réseau l'autorise (RM49) ; les autres lectures exigent une session.
	auth := s.handler.requireAuth
	writeOr := func(read func(http.HandlerFunc) http.HandlerFunc) func(http.HandlerFunc) http.HandlerFunc {
		return func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
					s.handler.requireRole(rbac.PermWrite, next)(w, r)
				default:
					read(next)(w, r)
				}
			}
		}
	}
	contrib := writeOr(auth)
	catalogue := writeOr(s.handler.requireReader)
	mux.HandleFunc("/api/components", catalogue(s.handler.handleComponents))
	mux.HandleFunc("/api/components/", catalogue(s.handler.handleComponent))
	mux.HandleFunc("/api/connections", contrib(s.handler.handleConnections))
	mux.HandleFunc("/api/connections/", contrib(s.handler.handleConnection))
	mux.HandleFunc("/api/assembly-links", contrib(s.handler.handleAssemblyLinks))
	// /api/modules/* — alias déprécié, ADR-11 : Model3D est désormais une
	// ressource REST unique (component), retirer ces deux lignes clôt la
	// fenêtre de transition une fois myr-web basculé.
	mux.HandleFunc("/api/modules", catalogue(legacyModulesAlias(s.handler.handleComponents)))
	mux.HandleFunc("/api/modules/", catalogue(legacyModulesAlias(s.handler.handleComponent)))
	mux.HandleFunc("/api/interfaces/", contrib(s.handler.handleInterface))
	mux.HandleFunc("/api/refs", catalogue(s.handler.handleRefs))
	mux.HandleFunc("/api/refs/categories", catalogue(s.handler.handleRefCategories))
	mux.HandleFunc("/api/refs/types", catalogue(s.handler.handleRefTypes))
	mux.HandleFunc("/api/refs/units", catalogue(s.handler.handleRefUnits))
	mux.HandleFunc("/api/channels", s.handler.handleChannels)             // GET public, PUT protégé
	mux.HandleFunc("/api/networks", s.handler.handleNetworks)             // GET liste des réseaux
	mux.HandleFunc("/api/networks/active", s.handler.handleNetworkActive) // GET/PUT réseau actif session
	mux.HandleFunc("/api/ping", s.handler.handlePing)                     // liveness — pas d'appel Fabric
	mux.HandleFunc("/api/status", s.handler.handleStatus)                 // status toujours accessible
	mux.HandleFunc("/api/licenses", catalogue(s.handler.handleLicenses))
	mux.HandleFunc("/api/licenses/", catalogue(s.handler.handleLicense))

	// Health + métriques — non authentifiés (appelés par infra/load-balancer)
	mux.HandleFunc("/api/health", s.handler.handleHealth)
	mux.Handle("/metrics", s.handler.metricsHandler())

	// Admin — réservé au rôle "admin"
	adminOnly := func(next http.HandlerFunc) http.HandlerFunc {
		return s.handler.requireRole(rbac.PermAdmin, next)
	}
	mux.HandleFunc("/api/admin/sessions", adminOnly(s.handler.handleAdminSessions))
	mux.HandleFunc("/api/admin/sessions/", adminOnly(s.handler.handleAdminSession))

	return secureHeaders(s.handler.limitBodies(mux)), nil
}
