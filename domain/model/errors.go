package model

import "errors"

var ErrNotFound = errors.New("model: not found")
var ErrIntegrityFailed = errors.New("model: integrity check failed")

// ErrBlockchainUnavailable est retournée par toute méthode du service qui
// nécessite BlockchainPort quand celui-ci n'est pas configuré (aucun adapter
// blockchain injecté). Le service ne substitue jamais silencieusement un
// autre stockage — voir domain/channel.ErrFabricUnavailable pour le même
// principe appliqué au domaine channel.
var ErrBlockchainUnavailable = errors.New("model: adapter blockchain non configuré")

// ErrBlockchainUnreachable est retournée par BlockchainPort quand l'adapter
// EST configuré mais que l'infrastructure sous-jacente échoue à répondre
// (nœuds injoignables, connectivité réseau) — panne transitoire, à distinguer
// de ErrBlockchainUnavailable (absence de configuration). Les deux méritent
// le même traitement côté appelant (retry, message "réessayez plus tard").
var ErrBlockchainUnreachable = errors.New("model: infrastructure blockchain injoignable")

// ErrNoThumbnailSource est retournée par Service.RegenerateThumbnail quand l'asset
// n'a aucun lien externe (Links) — seule source de miniature régénérable côté serveur.
var ErrNoThumbnailSource = errors.New("model: aucun lien externe enregistré pour régénérer la miniature")

// IsDegradedListErr indique si err signale que List/ListModules a échoué à joindre la
// blockchain (ErrBlockchainUnavailable ou ErrBlockchainUnreachable) — dans ce cas la méthode
// retourne quand même les brouillons locaux disponibles plutôt que rien. L'appelant (adapter
// in) reste responsable de vérifier explicitement cette fonction avant d'exploiter des données
// accompagnées d'une erreur non-nil, pour ne jamais présenter un résultat partiel comme complet
// sans le signaler à l'utilisateur.
func IsDegradedListErr(err error) bool {
	return errors.Is(err, ErrBlockchainUnavailable) || errors.Is(err, ErrBlockchainUnreachable)
}

// DegradedListWarning est le message que les adapters in (REST, CLI) présentent à
// l'utilisateur quand List/ListModules retourne un résultat partiel — voir IsDegradedListErr.
// Centralisé ici pour que CLI et REST affichent exactement le même avertissement (règle de
// parité fonctionnelle des adapters d'entrée).
const DegradedListWarning = "blockchain indisponible — liste limitée aux brouillons locaux non soumis"
