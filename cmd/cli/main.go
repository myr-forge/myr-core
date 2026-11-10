// cmd/cli/main.go — assemblage final
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"myr-core/adapters/in/cli"
	fabricadapter "myr-core/adapters/out/fabric"
	"myr-core/adapters/out/localstorage"
	"myr-core/adapters/out/webimage"
	"myr-core/domain/channel"
	"myr-core/domain/identity"
	"myr-core/domain/model"
	"myr-core/domain/network"
	"myr-core/domain/payment"
	"myr-core/domain/role"
	"myr-core/domain/session"
)

// noopPaymentGateway satisfait payment.PaymentGatewayPort tant qu'aucun adapter
// de paiement n'est câblé — évite un service nil (panique) côté CLI.
type noopPaymentGateway struct{}

var errNoPaymentGateway = fmt.Errorf("aucun gateway de paiement configuré pour le CLI")

// version est injectée via -ldflags "-X main.version=..." lors de la compilation
// (voir Makefile : VERSION dérivé de `git describe --tags`).
var version = "dev"

func (noopPaymentGateway) Transfer(*payment.Payment) error { return errNoPaymentGateway }
func (noopPaymentGateway) GetHistory(string) ([]*payment.Payment, error) {
	return nil, errNoPaymentGateway
}

// dataDir retourne le répertoire de données partagé avec myr-api.
// Priorité : MYR_DATA_DIR > ~/.Myr/server-data
func dataDir() string {
	if d := os.Getenv("MYR_DATA_DIR"); d != "" {
		return d
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".Myr", "server-data")
}

func main() {
	data := dataDir()
	_ = os.MkdirAll(data, 0700)

	networkStore := localstorage.NewJSONNetworkStore(filepath.Join(data, "networks.json"))
	networkSvc := network.NewService(networkStore, fabricadapter.GatewayTester{}).
		WithProvisioner(fabricadapter.NewFabricPeerProvisioner()).
		WithBootstrapper(fabricadapter.NewFabricNetworkBootstrapper())

	active, _ := networkSvc.GetActive()

	var caPort identity.CAPort
	if active != nil && active.CAEndpoint != "" {
		cfg := fabricadapter.ConfigFromProfile(active)
		if ca, err := fabricadapter.NewCAClient(cfg); err == nil {
			caPort = ca
		}
	}
	walletDir := filepath.Join(data, "wallets")
	_ = os.MkdirAll(walletDir, 0700)
	// Même fichier que myr-api : une demande reçue par l'API se décide en CLI.
	requestStore := localstorage.NewJSONRequestStore(filepath.Join(data, "account-requests.json"))
	identitySvc := identity.NewService(walletDir, caPort).WithRequestStore(requestStore)

	sessionStore := localstorage.NewJSONSessionStore(filepath.Join(data, "session.json"))
	sessionSvc := session.NewService(sessionStore)

	// ── Connexion Gateway Fabric (non bloquante) ──────────────────────────────
	// bc reste nil si Fabric est indisponible — pas de repli silencieux vers un
	// autre stockage : les opérations sur les modèles retournent alors
	// model.ErrBlockchainUnavailable (voir domain/model/service.go).
	// L'identité du profil réseau ne fait que lire : les écritures sont signées
	// par le wallet de l'identité du compte local CLI (myr session, ADR-12).
	// La table des rôles est celle publiée sur le ledger du réseau actif, sur
	// chacun de ses canaux ; sans réseau joignable, la table initiale
	// s'applique et ne peut pas être modifiée (role.ErrNoTable).
	var bc model.BlockchainPort
	var roleRepo role.Repo
	if active != nil {
		cfg := fabricadapter.ConfigFromProfile(active)
		if gw, err := fabricadapter.NewGatewayClient(cfg); err == nil {
			defer gw.Close()
			channels := active.AllChannels()
			bc = fabricadapter.NewFabricBlockchain(gw)
			roleRepo = fabricadapter.NewFabricRoleStore(gw, channels)
			if handle, ok := cli.SessionIdentity(sessionSvc); ok {
				if w, err := identitySvc.Wallet(handle); err == nil {
					if ugw, err := gw.ForIdentity(w.OrgID, w.CertPath(), w.KeyPath()); err == nil {
						defer ugw.Close()
						bc = fabricadapter.NewUserFabricBlockchain(ugw)
						roleRepo = fabricadapter.NewUserFabricRoleStore(ugw, channels)
					} else {
						fmt.Fprintf(os.Stderr, "myr: wallet %s inutilisable, accès au ledger en lecture seule : %v\n", handle, err)
					}
				}
			}
		}
	}
	// Table des rôles : référence des rôles attribuables, partagée par le
	// réseau (rôle d'inscription par défaut) et l'identité (RM46).
	roleSvc := role.NewService(roleRepo)
	networkSvc.WithRoles(roleSvc)
	identitySvc.WithRoles(roleSvc)

	// store partage le même fichier assets.json que myr-api (voir dataDir ci-dessus) —
	// donne au CLI accès aux connexions/miniatures/interfaces hors blockchain gérées côté GUI.
	store := localstorage.NewJSONBlockchain(filepath.Join(data, "assets.json"))
	// myr-core ne conserve jamais le fichier 3D transmis : seule son empreinte
	// SHA-256 est retenue (traçabilité RM01) — voir RM42 pour la vérification
	// des emplacements externes déclarés.
	modelSvc := model.NewService(bc).
		WithConnStore(store).
		WithThumbStore(store).
		WithIfaceStore(store).
		WithDraftStore(store).
		WithRemovedStore(store).
		WithOGImageFetcher(webimage.New())

	nodeStore := localstorage.NewJSONNodeStore(filepath.Join(data, "nodes.json"))
	channelStore := localstorage.NewJSONChannelStore(filepath.Join(data, "channels.json"))
	channelSvc := channel.NewService(channelStore).
		WithFabricConfig(fabricadapter.NewFabricChannelConfig(networkSvc, nodeStore))

	paymentSvc := payment.NewService(noopPaymentGateway{})

	cli.Version = version
	// La CLI agit pour l'identité du compte local (myr session) : seul le
	// propriétaire d'un asset le modifie, administrateur compris (FT-001).
	actor, _ := cli.SessionIdentity(sessionSvc)
	cli.Execute(modelSvc.As(actor), channelSvc, paymentSvc, networkSvc, roleSvc, identitySvc, sessionSvc)
}
