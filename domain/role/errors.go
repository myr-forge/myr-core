package role

import "errors"

var (
	ErrNotFound          = errors.New("role: introuvable")
	ErrAlreadyExists     = errors.New("role: existe déjà")
	ErrAdminImmutable    = errors.New("role: le rôle admin est immuable — ni modifiable ni supprimable")
	ErrInvalidPermission = errors.New("role: permission inconnue")
	ErrInvalidName       = errors.New("role: nom invalide")
	ErrAdminAutomatic    = errors.New("role: un rôle d'administration ne s'attribue jamais automatiquement")
	// ErrNoTable est retournée par toute modification de la table quand aucun
	// réseau ne la porte (Repo absent) : seule la table initiale s'applique.
	ErrNoTable = errors.New("role: aucune table des rôles publiable — réseau non configuré")
	// ErrSignerRequired est retournée par Repo quand la modification n'est pas
	// signée par l'identité de l'administrateur qui l'initie (ADR-12).
	ErrSignerRequired = errors.New("role: publication de la table des rôles impossible sans identité signataire de l'administrateur (wallet)")
)
