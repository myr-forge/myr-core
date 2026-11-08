// adapters/out/fabric/role_store.go
package fabric

import (
	"encoding/json"
	"errors"
	"fmt"

	"myr-core/domain/role"
)

var _ role.Repo = (*FabricRoleStore)(nil)

// FabricRoleStore implémente role.Repo via le chaincode myrcc (PutRole,
// DeleteRole, GetRoles) : chaque canal du réseau porte un exemplaire de la
// table des rôles (ADR-12). Une modification est publiée sur tous les canaux ;
// la table est lue sur le premier.
type FabricRoleStore struct {
	gc       GatewayProvider
	channels []string // vide : canal par défaut de la passerelle
	// userSigned indique que gc signe avec l'identité de l'administrateur qui
	// initie la modification : seule une telle passerelle publie la table.
	userSigned bool
}

// NewFabricRoleStore retourne un accès en lecture seule à la table des rôles :
// gc signe avec une identité d'organisation partagée, qui ne publie jamais —
// Save et Delete retournent role.ErrSignerRequired.
func NewFabricRoleStore(gc GatewayProvider, channels []string) *FabricRoleStore {
	return &FabricRoleStore{gc: gc, channels: channels}
}

// NewUserFabricRoleStore retourne un accès en lecture et en publication à la
// table des rôles : gc signe avec l'identité de l'administrateur.
func NewUserFabricRoleStore(gc GatewayProvider, channels []string) *FabricRoleStore {
	return &FabricRoleStore{gc: gc, channels: channels, userSigned: true}
}

func (f *FabricRoleStore) contract(channel string) ContractCaller {
	if channel == "" {
		return f.gc.Contract()
	}
	return f.gc.ContractFor(channel)
}

// FindAll lit la table publiée sur le premier canal du réseau.
func (f *FabricRoleStore) FindAll() ([]*role.Role, error) {
	channel := ""
	if len(f.channels) > 0 {
		channel = f.channels[0]
	}
	out, err := f.contract(channel).EvaluateTransaction("GetRoles")
	if err != nil {
		return nil, wrapGatewayErr("GetRoles", err)
	}
	var roles []*role.Role
	if err := json.Unmarshal(out, &roles); err != nil {
		return nil, fmt.Errorf("fabric GetRoles : table illisible : %w", err)
	}
	return roles, nil
}

// Save publie le rôle r sur chaque canal du réseau.
func (f *FabricRoleStore) Save(r *role.Role) error {
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return f.publish("PutRole", string(b))
}

// Delete retire le rôle name de la table de chaque canal du réseau.
func (f *FabricRoleStore) Delete(name string) error {
	return f.publish("DeleteRole", name)
}

// publish soumet la transaction fn(arg) sur chaque canal. Un échec sur un
// canal n'interrompt pas les autres : l'erreur nomme les canaux en échec, et
// relancer la même modification est sans effet sur ceux déjà à jour.
func (f *FabricRoleStore) publish(fn, arg string) error {
	if !f.userSigned {
		return role.ErrSignerRequired
	}
	channels := f.channels
	if len(channels) == 0 {
		channels = []string{""}
	}
	var errs []error
	for _, ch := range channels {
		if _, err := f.contract(ch).SubmitTransaction(fn, arg); err != nil {
			label := ch
			if label == "" {
				label = "par défaut"
			}
			errs = append(errs, fmt.Errorf("canal %s : %w", label, wrapGatewayErr(fn, err)))
		}
	}
	return errors.Join(errs...)
}
