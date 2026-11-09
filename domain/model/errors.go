package model

import "errors"

var ErrNotFound = errors.New("model: not found")

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

// ErrSignerRequired est retournée par BlockchainPort.StoreModelRecord quand
// aucune identité d'utilisateur n'est disponible pour signer l'écriture :
// toute écriture sur le ledger est signée par le certificat de l'utilisateur
// qui l'initie, jamais par une identité d'organisation partagée (ADR-12).
var ErrSignerRequired = errors.New("model: écriture sur le ledger impossible sans identité signataire de l'utilisateur (wallet)")

// ErrAssetSubmitted est retournée par toute modification d'un asset déjà
// soumis — licence, interfaces, composition, re-soumission : il est immuable
// sur le ledger (RM19, ADR-12). Toute évolution passe par un nouvel asset
// dérivé (fork, ParentID vers l'asset d'origine) ; seules ses métadonnées
// descriptives restent modifiables par son propriétaire (UpdateAsset).
var ErrAssetSubmitted = errors.New("model: asset soumis, immuable — créez un asset dérivé (fork) pour le faire évoluer")

// ErrUnknownLicense est retournée quand la licence proposée pour un asset
// n'existe pas dans le catalogue.
var ErrUnknownLicense = errors.New("model: licence absente du catalogue")

// ErrThumbnailTooLarge est retournée par SaveThumbnail pour une image
// dépassant MaxThumbnailBytes.
var ErrThumbnailTooLarge = errors.New("model: miniature trop volumineuse")

// ErrLicenseIncompatible est retournée quand la licence proposée pour un
// asset dérivé n'est pas compatible avec celle de son parent (RM03).
var ErrLicenseIncompatible = errors.New("model: incompatibilité de licence avec l'asset parent")

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

// Règles d'assemblage (UCAM01 : RM09, RM10, RM50, RM51, RM52).
var (
	// ErrInterfaceInUse : une interface (d'une instance donnée) est déjà engagée
	// dans une liaison (RM09).
	ErrInterfaceInUse = errors.New("model: interface déjà engagée dans une liaison")
	// ErrIncompatibleInterfaces : la paire d'interfaces n'est pas compatible
	// (RM10, critères RM11).
	ErrIncompatibleInterfaces = errors.New("model: interfaces incompatibles")
	// ErrSelfLink : une interface ne peut pas être reliée à elle-même.
	ErrSelfLink = errors.New("model: une interface ne peut pas être reliée à elle-même")
	// ErrConnectionExists : la liaison demandée existe déjà ; le client la
	// retrouve dans la liste des liaisons, le serveur ne la devine pas.
	ErrConnectionExists = errors.New("model: liaison déjà existante")
	// ErrConnectionNotInModule : la liaison rattachée à un module n'existe pas,
	// ou ne relie pas deux de ses instances.
	ErrConnectionNotInModule = errors.New("model: liaison inexistante ou étrangère aux instances du module")
	// ErrCyclicComposition : un asset ne peut pas contenir, même indirectement,
	// le module auquel on l'ajoute.
	ErrCyclicComposition = errors.New("model: composition cyclique — un module ne peut pas se contenir lui-même")
	// ErrInterfaceCarrierChange : une interface ne change pas d'asset porteur.
	ErrInterfaceCarrierChange = errors.New("model: une interface ne change pas d'asset porteur")
)
