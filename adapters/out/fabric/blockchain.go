// adapters/out/fabric/blockchain.go
package fabric

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"myr-core/domain/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validID = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]{1,128}$`)

func validateID(v, label string) error {
	if !validID.MatchString(v) {
		return fmt.Errorf("%s invalide : caractères non autorisés ou longueur hors limites", label)
	}
	return nil
}

// wrapGatewayErr classe l'erreur d'un appel Evaluate/Submit au gateway Fabric : si le code
// gRPC indique que l'infrastructure est injoignable (transport indisponible, ou aucun pair
// disponible pour évaluer/endosser le chaincode), l'erreur est chaînée à
// model.ErrBlockchainUnreachable (503 côté REST via internalErr) plutôt que remontée telle
// quelle (500 générique) — voir spécification Architecture_Hexagonale.md §1.1.
func wrapGatewayErr(op string, err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	unreachable := ok && (st.Code() == codes.Unavailable ||
		(st.Code() == codes.FailedPrecondition && strings.Contains(st.Message(), "no peers available")))
	if unreachable {
		return fmt.Errorf("fabric %s : %w: %w", op, model.ErrBlockchainUnreachable, err)
	}
	return fmt.Errorf("fabric %s : %w", op, err)
}

// ledgerLookupErr traduit le refus du chaincode pour un asset absent du
// ledger (« GetModel : asset … introuvable ») en model.ErrNotFound.
func ledgerLookupErr(id string, err error) error {
	if err != nil && strings.Contains(err.Error(), "introuvable") {
		return fmt.Errorf("%w : asset %s sur le ledger : %w", model.ErrNotFound, id, err)
	}
	return err
}

var _ model.BlockchainPort = (*FabricBlockchain)(nil)

// FabricBlockchain implémente model.BlockchainPort via le chaincode myrcc.
type FabricBlockchain struct {
	gc GatewayProvider
	// userSigned indique que gc signe avec l'identité de l'utilisateur qui
	// initie l'opération : seule une telle passerelle écrit sur le ledger.
	userSigned bool
}

// NewFabricBlockchain retourne un accès au ledger en lecture seule : gc signe
// avec une identité d'organisation partagée, qui n'écrit jamais sur le ledger
// (ADR-12) — StoreModelRecord retourne model.ErrSignerRequired.
func NewFabricBlockchain(gc GatewayProvider) *FabricBlockchain {
	return &FabricBlockchain{gc: gc}
}

// NewUserFabricBlockchain retourne un accès au ledger en lecture et en
// écriture : gc signe avec l'identité de l'utilisateur qui initie les
// écritures (ADR-12, voir GatewayClient.ForIdentity).
func NewUserFabricBlockchain(gc GatewayProvider) *FabricBlockchain {
	return &FabricBlockchain{gc: gc, userSigned: true}
}

// StoreModelRecord soumet une transaction sur le canal de l'asset (m.ChannelID).
func (f *FabricBlockchain) StoreModelRecord(m *model.Model3D) error {
	if !f.userSigned {
		return model.ErrSignerRequired
	}
	if err := validateID(m.ChannelID, "channelID"); err != nil {
		return err
	}
	if err := validateID(m.ID, "model ID"); err != nil {
		return err
	}
	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("fabric StoreModelRecord marshal : %w", err)
	}
	if _, err := f.gc.ContractFor(m.ChannelID).SubmitTransaction("StoreModel", string(data)); err != nil {
		return wrapGatewayErr("StoreModelRecord submit", err)
	}
	return nil
}

// UpdateModelMetadata soumet le remplacement des métadonnées descriptives d'un
// asset enregistré — signé par l'utilisateur, que le chaincode vérifie être le
// propriétaire (ADR-12).
func (f *FabricBlockchain) UpdateModelMetadata(id, channelID string, meta model.AssetMetadata) error {
	if !f.userSigned {
		return model.ErrSignerRequired
	}
	if err := validateID(channelID, "channelID"); err != nil {
		return err
	}
	if err := validateID(id, "model ID"); err != nil {
		return err
	}
	data, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("fabric UpdateModelMetadata marshal : %w", err)
	}
	if _, err := f.gc.ContractFor(channelID).SubmitTransaction("UpdateMetadata", id, string(data)); err != nil {
		return wrapGatewayErr("UpdateModelMetadata submit", err)
	}
	return nil
}

// GetModelRecord évalue une query sur le canal spécifié.
// channelID="" utilise le canal par défaut de la configuration.
func (f *FabricBlockchain) GetModelRecord(id, channelID string) (*model.Model3D, error) {
	if err := validateID(id, "model ID"); err != nil {
		return nil, err
	}
	result, err := f.gc.ContractFor(channelID).EvaluateTransaction("GetModel", id)
	if err != nil {
		return nil, ledgerLookupErr(id, wrapGatewayErr("GetModelRecord evaluate", err))
	}
	var m model.Model3D
	if err := json.Unmarshal(result, &m); err != nil {
		return nil, fmt.Errorf("fabric GetModelRecord unmarshal : %w", err)
	}
	return &m, nil
}

// ListModelRecords évalue une query sur le canal spécifié.
func (f *FabricBlockchain) ListModelRecords(channelID string) ([]*model.Model3D, error) {
	if err := validateID(channelID, "channelID"); err != nil {
		return nil, err
	}
	result, err := f.gc.ContractFor(channelID).EvaluateTransaction("ListModels", channelID)
	if err != nil {
		return nil, wrapGatewayErr("ListModelRecords evaluate", err)
	}
	var list []*model.Model3D
	if err := json.Unmarshal(result, &list); err != nil {
		return nil, fmt.Errorf("fabric ListModelRecords unmarshal : %w", err)
	}
	return list, nil
}
