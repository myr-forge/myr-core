---
tags:
  - couche/conception
  - type/conception
  - rm/RM46
  - rm/RM47
  - rm/RM48
  - relecture/ecart
---
# DC — D1 : Identité, RBAC et Session

## 1. Objectif

Ce document couvre les domaines qui composent la couche accès du système Myr. **Il n'existe pas de domaine `auth` au sens email/mot de passe/JWT** — l'architecture repose sur quatre pièces distinctes :

- **`domain/identity`** : identité cryptographique Fabric CA — enrôlement X.509, wallets locaux MSP, demandes d'accès (`AccountRequest`). C'est la seule source de vérité de « qui est quelqu'un ».
- **`domain/role`** : RBAC dynamique — un catalogue de permissions (`read`, `write`, `network.admin`, `role.admin`, `identity.admin`, `admin`) et des rôles (4 intégrés + rôles personnalisés) associant un sous-ensemble de ces permissions.
- **Sessions REST** (`adapters/in/rest/session.go`, `session_redis.go`) : couche adaptateur (pas un domaine) qui émet un **token opaque** (32 octets aléatoires hex, pas un JWT) après un enrôlement CA réussi ou un accès invité, et le fait correspondre à `{role, pseudo, channel}` pendant 7 jours.
- **`domain/session`** : concept **sans rapport** avec les sessions REST ci-dessus — c'est le compte local de la machine exécutant le CLI `myr` (`{Name, OrgID}`, un seul enregistrement, bootstrap au premier lancement). Nommage historique malheureux à ne pas confondre.

Il n'y a ni bcrypt, ni JWT, ni table SQL utilisateur nulle part dans le code — aucun adapter `sqlite` n'existe dans `adapters/out/`.

---

## 2. Diagramme de classes

```plantuml
@startuml
skinparam classAttributeIconSize 0
skinparam class {
  BackgroundColor #FEFECE
  BorderColor #A80036
}

package "domain/identity" {

  class MyrIdentity {
    + ID : string
    + DisplayName : string
    + LegalName : string
    + Email : string
    + Country : string
    + Organization : string
    + LicenseDefault : string
    + IPAgreement : bool
    + IPAgreedAt : time.Time
    + CreatedAt : time.Time
    + Status : IdentityStatus
    + Role : string
  }

  enum IdentityStatus {
    pending
    active
    suspended
  }

  class WalletEntry {
    + Handle : string <<pseudo@org>>
    + Name : string
    + OrgID : string
    + Status : string <<Myr.status du certificat>>
    + Role : string <<Myr.role du certificat>>
    + MSPDir : string <<chemin absolu MSP>>
    + CanSignIn() : error
  }

  class RegisterRequest {
    + Name : string <<pseudo@org>>
    + OrgID : string
    + Password : string <<secret enrollment, jamais stocké>>
    + Role : string
    + Status : string <<vide = pending>>
    + DisplayName : string
    + LegalName : string
    + Email : string
    + Country : string
  }

  class AccountRequest {
    + ID : string
    + Pseudo : string
    + DisplayName : string
    + Email : string
    + OrgID : string
    + Message : string <<texte libre — pas de champ "rôle souhaité">>
    + Status : RequestStatus
    + CreatedAt : string <<ISO 8601>>
  }

  enum RequestStatus {
    pending
    approved
    rejected
    failed
  }

  class RegistrationPolicy {
    + AllowAutoRegister : bool
    + RequireApproval : bool
    + DefaultRole : string <<vide = reader>>
  }
  note bottom of RegistrationPolicy
    Règles d'inscription du réseau,
    PolicyOf(NetworkProfile)
  end note

  MyrIdentity ..> IdentityStatus : uses
  WalletEntry ..> IdentityStatus : uses (statut lu depuis le certificat)
  AccountRequest ..> RequestStatus : uses
  AccountRequest ..> RegistrationPolicy : traitée selon
}

package "domain/role" {

  class Role {
    + ID : string
    + Name : string
    + Permissions : []Permission
    + BuiltIn : bool
    + CreatedAt : time.Time
  }

  enum Permission {
    read
    write
    network.admin
    role.admin
    identity.admin
    admin
  }

  Role "1" --> "0..*" Permission : grants
}

package "adapters/in/rest (pas un domaine)" {

  class myrSession {
    + Token : string <<32 octets aléatoires, hex>>
    + Role : string
    + Pseudo : string
    + Handle : string <<pseudo@org signataire, vide pour un invité>>
    + Channel : string
    + NetworkID : string
    + ExpiresAt : time.Time
  }
}

package "domain/session" {

  class Session {
    + Name : string <<pseudo@org>>
    + OrgID : string
    + CreatedAt : time.Time
  }
  note right of Session
    Compte local CLI (bootstrap machine),
    SANS RAPPORT avec myrSession ci-dessus
    malgré le nom partagé.
  end note
}

' Relations inter-domaines
WalletEntry ..> MyrIdentity : vue locale (fichiers MSP)\nd'une identité enregistrée en CA
myrSession ..> MyrIdentity : créée après un Enroll() réussi\n(rôle lu dans le certificat, figé pour la durée de la session)
myrSession ..> Role : Role vérifié via RoleService.HasPermission

@enduml
```

---

## 3. Description des classes

### MyrIdentity (`domain/identity`)

- **Rôle :** Identité blockchain d'un auteur dans le réseau Myr. Les attributs `Status` et `Role` sont écrits dans le certificat X.509 par la Fabric CA (attributs `Myr.status` et `Myr.role`).
- **Invariants :**
  - `IPAgreement` doit être `true` avant tout acte de publication sur la blockchain.
  - `LicenseDefault` définit la licence par défaut appliquée aux créations de cet auteur.
  - `Role` est une chaîne libre — sa validité vis-à-vis du catalogue RBAC (`domain/role`) n'est pas vérifiée à ce niveau.
  - `Status` reflète l'état de l'identité dans la CA Fabric.
- **Cycle de vie :** créée `active` (création par l'administrateur, ou auto-enregistrement sans approbation) ou `pending` (auto-enregistrement avec approbation) ; `pending` → `active` à l'approbation, `pending` → `suspended` au refus (`CAPort.UpdateAttributes`). Le changement apparaît dans le certificat au prochain enrôlement.

### WalletEntry (`domain/identity`)

- **Rôle :** Vue locale d'un wallet Fabric (fichiers MSP sur disque, `~/.Myr/wallets/<handle>/msp/`). C'est la **seule** représentation du wallet — aucun miroir chiffré en base ne coexiste (pas de BDD relationnelle, cf. principe de décentralisation).
- **Invariants :**
  - `Handle` suit le format `pseudo@org` (ex : `alice@Org1`).
  - `MSPDir` est un chemin absolu vers le répertoire MSP local.
  - `Status` et `Role` sont lus dans l'extension d'attributs du certificat X.509 local (`CAPort.CertAttributes`), pas depuis une base de données ; ils restent vides si le certificat ne les porte pas.
  - `CanSignIn` n'autorise l'ouverture d'une session qu'à une identité `active` dotée d'un rôle (`ErrIdentityPending`, `ErrIdentitySuspended`, `ErrIdentityInactive`, `ErrIdentityWithoutRole` sinon).

### RegisterRequest (`domain/identity`)

- **Rôle :** DTO de création d'identité auprès de la Fabric CA (`CAPort.Register`). Valeur éphémère, non persistée.
- **Invariants :** `Password` (secret d'enrôlement) n'est jamais stocké par `myr` — seul le certificat/clé résultant de l'enrôlement l'est (dans le wallet local). `Status` fixe l'attribut `Myr.status` à l'enregistrement : `active` ou `suspended` explicitement, `pending` sinon.

### RegistrationPolicy (`domain/identity`)

- **Rôle :** Règles d'inscription d'un réseau, fixées par son administrateur sur le profil réseau (`AllowAutoRegister`, `RequireApproval`, `AutoRegisterRole`) et converties par `identity.PolicyOf(*NetworkProfile)` — sans réseau, la politique est fermée. Les adapters `in/` (CLI, REST) la transmettent à `RequestAccess` sans l'interpréter.

### AccountRequest (`domain/identity`)

- **Rôle :** Demande d'accès soumise par un utilisateur sans compte, acceptée seulement si le réseau autorise l'auto-enregistrement (`RegistrationPolicy`) — ou demande de rôle d'une identité déjà inscrite (`Kind = role`, rôle demandé dans `Role`, UCA08, RM47). Stockée localement (`adapters/out/localstorage/request_store.go`, fichier partagé par le CLI et l'API), consultable et décidée par un administrateur.
- **Invariants :**
  - Aucune demande n'est conservée sur un réseau fermé (`ErrRegistrationClosed`).
  - L'identité CA `pseudo@org` est créée dès la demande ; la demande ne fait que tracer la décision.
  - `CreatedAt` est une chaîne ISO 8601 (`string`, pas `time.Time`).
  - Une identité n'a jamais plus d'une demande de rôle en attente (`ErrRequestAlreadyPending`).
- **Cycle de vie (inscription) :** `pending` à l'enregistrement → `approved` (identité créée active) | reste `pending` (identité créée en attente, approbation exigée) → `approved` (`ApproveRequest`) ou `rejected` (`RejectRequest`, identité suspendue) | `failed` si la CA refuse l'enregistrement. Seule une demande `pending` peut être décidée (`ErrRequestNotPending`).
- **Cycle de vie (demande de rôle) :** `approved` dès la demande si le rôle est en attribution automatique (rôle attribué) | `pending` sinon → `approved` (`ApproveRequest`, rôle attribué s'il figure toujours dans la table, RM46) ou `rejected` (`RejectRequest`, identité inchangée) | `failed` si la CA refuse l'attribution automatique.


**⚠️ Écart E4 (critique) — l'enregistrement dépend d'une configuration invisible depuis le `NetworkProfile` :** le client CA utilisé par `identitySvc.RequestAccess` (`cmd/api/main.go`) est construit une fois au démarrage depuis les variables d'environnement (`FABRIC_CA_ENDPOINT`, `FABRIC_CA_ADMIN_CERT_PATH`...) ou un `fabric.env` optionnel — jamais depuis les champs `CAEndpoint`/`CAAdminCertPath`/`CAAdminKeyPath` du `NetworkProfile` actif, alors que ces champs existent précisément pour ça et sont déjà utilisés ainsi par le CLI et par le pool blockchain multi-réseau du serveur (`ConfigFromProfile`, voir `specs/roadmap_dev.md` § Écarts Identité & Session pour le détail vérifié). Conséquence : `myr network update --auto-register` peut afficher une configuration entièrement correcte (CLI et `GET /api/identity/policy` inclus) sans que l'auto-enregistrement ne fonctionne réellement — la demande échoue faute de CA.

#ecart "a trancher"
 

### Role / Permission (`domain/role`)

- **Rôle :** RBAC dynamique. `Permission` est un catalogue plat de 6 valeurs. `Role` associe un nom à un sous-ensemble de permissions et à un réglage d'attribution sur demande (`AutoGrant`, UCA08, RM47).
- **Stockage :** la table des rôles est publiée sur le ledger, un exemplaire par canal du réseau (port `role.Repo`, adapter `adapters/out/fabric/role_store.go`, chaincode `PutRole`/`DeleteRole`/`GetRoles` — `Chaincode.md` §4.3, ADR-12). Le CLI la publie avec le wallet du compte local ; l'API la lit seulement, et la relit au plus toutes les 30 secondes (`WithRefresh`).
- **Invariants (RM48) :**
  - Tant qu'aucune table n'est publiée, la table initiale s'applique (`DefaultRoles` : `reader`, `contributor`, `auditor`, `admin`, marqués `BuiltIn`) ; la première publication part de cette table.
  - Tout rôle se modifie ou se supprime, sauf `admin`, qui porte toujours toutes les permissions (`ErrAdminImmutable`) — le domaine le rétablit à chaque lecture.
  - Le nom d'un rôle doit respecter `^[a-z][a-z0-9_-]{1,63}$`.
  - Un rôle en attribution automatique ne porte aucune permission d'administration (`ErrAdminAutomatic`).
  - `HasPermission` échoue fermé : un rôle inconnu, ou supprimé, ne porte aucune permission.
  - Table illisible (réseau injoignable) : la dernière table lue s'applique, à défaut la table initiale ; sans réseau, la table ne se modifie pas (`ErrNoTable`).
- **Cycle de vie :** `Create`/`Update`/`Delete` via `myr role ...` (CLI uniquement — pas de CRUD REST, seulement une consultation implicite par le middleware `requireRole`).

### myrSession (adaptateur REST, `adapters/in/rest/session.go`)

- **Rôle :** Jeton opaque de session HTTP. **Ce n'est pas une entité de domaine** — c'est un détail d'implémentation de l'adaptateur `in/rest`, documenté ici car il joue le rôle qu'un domaine `auth` aurait joué dans une architecture plus classique.
- **Invariants :**
  - `Token` est généré par `crypto/rand`, jamais dérivé de données utilisateur.
  - `ExpiresAt` = création + 7 jours (`sessionTTL`), non renouvelable automatiquement.
  - `Role` est le `Myr.role` du certificat émis à la connexion, fixé **une fois pour toutes à la création** de la session, sans resynchronisation ultérieure avec l'identité CA : un changement de rôle s'applique à la session suivante.
  - `Handle` est l'identité `pseudo@org` dont le wallet signe les écritures de l'utilisateur sur le ledger (ADR-12) ; vide pour un invité, qui ne peut pas écrire.
- **Cycle de vie :** `create` (connexion ou invité) → `get` (vérifié à chaque requête protégée) → `delete` (révocation admin uniquement, pas de self-service) ou expiration naturelle.

### Session (`domain/session`)

- **Rôle :** Compte local **de la machine** exécutant le CLI `myr` — pas un compte utilisateur au sens web. Un seul enregistrement possible, persisté en JSON.
- **Invariants :**
  - `Name` suit le format `pseudo@org`.
  - Un seul fichier session à la fois.
- **Cycle de vie :** `Create` au premier lancement du CLI ; `Logout`/`Clear` réinitialise le fichier.
- **Exposition :** `myr session create --name <n> [--org-id <id>]` / `myr session show` / `myr session logout` (`adapters/in/cli/session.go`), câblés dans `cmd/cli/main.go` (`localstorage.NewJSONSessionStore` + `session.NewService`). N'est volontairement pas exposé en REST : c'est un concept propre à la machine qui exécute le CLI (voir § « Rôle » ci-dessus), sans équivalent côté session HTTP.

---

## 4. Relations et dépendances

- **WalletEntry → MyrIdentity** : `WalletEntry` est la vue fichier locale d'une identité enregistrée auprès de la CA ; il n'existe pas de seconde vue chiffrée en base (pas de BDD relationnelle, cf. principe de décentralisation).
- **myrSession → MyrIdentity** : une session REST est créée après un enrôlement CA réussi d'une identité active dotée d'un rôle (`POST /api/identity/session`) ; elle retient son identifiant (`Handle`, signataire des écritures) et une copie de son rôle, figée au moment de la création et jamais resynchronisée.
- **myrSession → Role** : à chaque requête protégée, `requireRole` interroge `RoleService.HasPermission(session.Role, permission)` — c'est le seul point de contact entre la session REST et le RBAC.
- **Session (domain/session) ↔ WalletEntry** : lien conceptuel seulement (`Session.Name` pourrait correspondre à un `WalletEntry.Handle`) — aucune vérification ni jointure automatique entre les deux n'est requise par la conception.

---

## 5. Ports (interfaces Go)

### Port entrant — IdentityService

```plantuml
@startuml
skinparam classAttributeIconSize 0
interface IdentityService {
  + ListLocalWallets() : ([]WalletEntry, error)
  + Wallet(handle) : (WalletEntry, error)
  + CreateAccount(ctx, pseudo, orgID, role, email) : (handle, secret string, err)
  + Enroll(ctx, name, secret, orgID) : (WalletEntry, error)
  + GetStatus(ctx, wallet WalletEntry) : (string, error)
  + ReEnroll(ctx, wallet WalletEntry) : (WalletEntry, error)
  + LoadGuestWallet(certPath, keyPath, caCertPath, orgID) : (WalletEntry, error)
  + RequestAccess(ctx, req AccountRequest, policy RegistrationPolicy) : (*AccountRequest, secret string, err)
  + ApproveRequest(ctx, id) : (*AccountRequest, error)
  + RejectRequest(ctx, id) : (*AccountRequest, error)
  + ListRequests() : ([]*AccountRequest, error)
  + WalletDir() : string
  + SetRole(ctx, name, newRole string) : error
}
@enduml
```

### Port sortant — CAPort / RequestStore (identity)

```plantuml
@startuml
skinparam classAttributeIconSize 0
interface CAPort {
  + Register(ctx, req RegisterRequest) : (secret string, err)
  + Enroll(ctx, name, secret) : (certPEM, keyPEM, caCertPEM string, err)
  + GetStatus(ctx, name) : (status string, err)
  + ReEnroll(ctx, name, certPEM, keyPEM) : (newCertPEM string, err)
  + UpdateAttributes(ctx, name, attrs map[string]string) : error
  + CertAttributes(certPEM) : (map[string]string, error)
}
note right of CAPort
  CertAttributes lit localement l'extension
  d'attributs du certificat (format propre
  à la CA), sans appel réseau
end note

interface RequestStore {
  + Save(req *AccountRequest) : error
  + FindAll() : ([]*AccountRequest, error)
}
@enduml
```

### Port entrant/sortant — RoleService / Repo (role)

```plantuml
@startuml
skinparam classAttributeIconSize 0
interface RoleService {
  + Create(name, permissions) : (*Role, error)
  + Update(name, permissions) : (*Role, error)
  + Delete(name) : error
  + Get(name) : (*Role, error)
  + List() : ([]*Role, error)
  + HasPermission(roleName, p Permission) : bool
}

interface Repo {
  + Save(r *Role) : error
  + FindAll() : ([]*Role, error)
  + FindByName(name) : (*Role, error)
  + Delete(name) : error
}
@enduml
```

### Port entrant/sortant — SessionService / Store (domain/session — compte local CLI)

```plantuml
@startuml
skinparam classAttributeIconSize 0
interface SessionService {
  + Create(name, orgID) : (*Session, error)
  + Current() : (*Session, error)
  + Logout() : error
}

interface Store {
  + Save(s *Session) : error
  + Load() : (*Session, error)
  + Clear() : error
}
@enduml
```

### Interface interne (non-domaine) — sessionBackend (REST)

```plantuml
@startuml
skinparam classAttributeIconSize 0
interface sessionBackend {
  + create(pseudo, handle, role, channel) : (*myrSession, error)
  + get(token) : *myrSession
  + setNetwork(token, networkID) : bool
  + setChannel(token, channel) : bool
  + delete(token) : void
  + list() : []*myrSession
}
note right : Implémentée par sessionStore (mémoire + JSON optionnel)\net redisSessionStore (multi-instances)
@enduml
```

---

## 6. Décisions de conception

| ID  | Décision | Raison |
|-----|----------|--------|
| DC-D1-01 | Pas de domaine `auth` séparé — l'identité CA est l'unique source de vérité | Une architecture à deux systèmes d'identité (compte web + identité blockchain) a été abandonnée au profit d'un système unique : l'identité CA. La session REST n'est qu'une couche d'adaptateur au-dessus, pas un second domaine d'identité. |
| DC-D1-02 | Le wallet local n'est représenté que par des fichiers MSP sur disque, sans miroir chiffré en base | Source unique pour éviter la divergence entre deux vues du même wallet — voir `Conception_intro.md` ADR-03 pour l'écart de sécurité qui en résulte (absence de chiffrement au repos). |
| DC-D1-03 | Le certificat nomme le rôle, `domain/role` en définit les permissions | Le rôle d'une session est le `Myr.role` du certificat de l'identité, jamais un rôle choisi par le client ni figé par le serveur ; `RoleService.HasPermission(rôle, permission)` traduit ce nom en permissions (table des rôles, éditable par l'administrateur — ADR-12). |
| DC-D1-04 | Le token de session REST est opaque, pas un JWT | Évite toute dépendance à une bibliothèque JWT et tout risque de claims falsifiables côté client ; en contrepartie, aucune information (rôle, pseudo) n'est extractible du token lui-même — elle doit être recherchée côté serveur (`sessionBackend.get`). |
| DC-D1-05 | `domain/session` (compte local CLI) reste un domaine à part, non couplé conceptuellement aux autres domaines (identité CA, RBAC) | Exposé en CLI (`myr session create/show/logout`) comme bootstrap machine, mais `Session.Name` n'est ni validé ni synchronisé automatiquement avec `WalletEntry.Handle` (`domain/identity`) ou un rôle RBAC (`domain/role`) — un lien purement conceptuel, pas une jointure applicative (voir § 4 « Relations et dépendances »). |
| DC-D1-06 | Lecture des composants (`GET /api/components`) ouverte aux visiteurs non authentifiés, inconditionnellement — écriture (POST/PUT/PATCH/DELETE) reste gated par rôle | EF17 (D4, `specs/2-Analyse/UCCL-Composant_Lecture/UCCL01.md`) exige un accès public en lecture, sans condition. Le RBAC dynamique (`domain/role`) et la session REST (`myrSession`) ne changent pas de mécanisme : la route GET est dispensée de `requireAuth` par distinction de méthode HTTP dans le routeur (`server.go`, pattern déjà décrit dans `API_REST.md` §1 « routes mixtes ») — pas par un nouveau système d'identité « visiteur ». |

**⚠️ Incohérence relevée entre deux documents d'analyse (ni l'un ni l'autre n'est corrigé ici — signalée pour arbitrage PO) :** `UCCL01.md` (composants) décrit l'accès visiteur comme **toujours ouvert**, sans condition. `specs/2-Analyse/UCMOD-Module/UCMOD04.md` (ligne 134, lecture d'un module) décrit au contraire un accès **conditionné par la configuration réseau** (« sauf si la configuration réseau autorise les Visiteurs »). DC-D1-06 ci-dessus ne couvre donc que `GET /api/components`, conformément à `UCCL01.md` — il ne s'étend pas à `GET /api/modules`, dont la lecture reste `Auth` par défaut (voir `API_REST.md` §6) tant que le PO n'a pas tranché si D4 doit avoir un comportement uniforme entre composants et modules, ou si cette différence (composant toujours public, module conditionnel) est intentionnelle. C'est aussi la question déjà notée `specs/2-Analyse/todo.md` (accès visiteur conditionné par un flag réseau) — elle ne semble concerner que les modules, pas les composants, au vu du texte actuel des deux UC.

---

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Étape suivante — code, par section**
- [§ 2. Diagramme de classes](DC_D1_Auth_Identity.md#2.%20Diagramme%20de%20classes) → [RoleService.HasPermission](../../docs/code/fonctions/role.RoleService.HasPermission.md)
- [§ MyrIdentity (`domain/identity`)](DC_D1_Auth_Identity.md#MyrIdentity%20%28`domain/identity`%29) → [CAPort.UpdateAttributes](../../docs/code/fonctions/identity.CAPort.UpdateAttributes.md)
- [§ WalletEntry (`domain/identity`)](DC_D1_Auth_Identity.md#WalletEntry%20%28`domain/identity`%29) → [CAPort.CertAttributes](../../docs/code/fonctions/identity.CAPort.CertAttributes.md)
- [§ RegisterRequest (`domain/identity`)](DC_D1_Auth_Identity.md#RegisterRequest%20%28`domain/identity`%29) → [CAPort.Register](../../docs/code/fonctions/identity.CAPort.Register.md)
- [§ RegistrationPolicy (`domain/identity`)](DC_D1_Auth_Identity.md#RegistrationPolicy%20%28`domain/identity`%29) → [IdentityService.RequestAccess](../../docs/code/fonctions/identity.IdentityService.RequestAccess.md)
- [§ AccountRequest (`domain/identity`)](DC_D1_Auth_Identity.md#AccountRequest%20%28`domain/identity`%29) → [REST /api/identity/policy](../../docs/code/routes/REST%20api-identity-policy.md) · [CLI myr network update](../../docs/code/commandes/CLI%20myr-network-update.md) · [IdentityService.ApproveRequest](../../docs/code/fonctions/identity.IdentityService.ApproveRequest.md) · [IdentityService.RejectRequest](../../docs/code/fonctions/identity.IdentityService.RejectRequest.md) · [IdentityService.RequestAccess](../../docs/code/fonctions/identity.IdentityService.RequestAccess.md)
- [§ Role / Permission (`domain/role`)](DC_D1_Auth_Identity.md#Role%20/%20Permission%20%28`domain/role`%29) → [CLI myr role](../../docs/code/commandes/CLI%20myr-role.md) · [RoleService.HasPermission](../../docs/code/fonctions/role.RoleService.HasPermission.md)
- [§ Session (`domain/session`)](DC_D1_Auth_Identity.md#Session%20%28`domain/session`%29) → [CLI myr session create](../../docs/code/commandes/CLI%20myr-session-create.md) · [CLI myr session logout](../../docs/code/commandes/CLI%20myr-session-logout.md) · [CLI myr session show](../../docs/code/commandes/CLI%20myr-session-show.md)
- [§ 4. Relations et dépendances](DC_D1_Auth_Identity.md#4.%20Relations%20et%20dépendances) → [REST /api/identity/session](../../docs/code/routes/REST%20api-identity-session.md) · [RoleService.HasPermission](../../docs/code/fonctions/role.RoleService.HasPermission.md)
- [§ Port entrant — IdentityService](DC_D1_Auth_Identity.md#Port%20entrant%20—%20IdentityService) → [IdentityService.ApproveRequest](../../docs/code/fonctions/identity.IdentityService.ApproveRequest.md) · [IdentityService.CreateAccount](../../docs/code/fonctions/identity.IdentityService.CreateAccount.md) · [IdentityService.ListLocalWallets](../../docs/code/fonctions/identity.IdentityService.ListLocalWallets.md) · [IdentityService.ListRequests](../../docs/code/fonctions/identity.IdentityService.ListRequests.md) · [IdentityService.LoadGuestWallet](../../docs/code/fonctions/identity.IdentityService.LoadGuestWallet.md) · [IdentityService.RejectRequest](../../docs/code/fonctions/identity.IdentityService.RejectRequest.md) · [IdentityService.RequestAccess](../../docs/code/fonctions/identity.IdentityService.RequestAccess.md) · [IdentityService.SetRole](../../docs/code/fonctions/identity.IdentityService.SetRole.md) · [IdentityService.Wallet](../../docs/code/fonctions/identity.IdentityService.Wallet.md) · [IdentityService.WalletDir](../../docs/code/fonctions/identity.IdentityService.WalletDir.md)
- [§ Port sortant — CAPort / RequestStore (identity)](DC_D1_Auth_Identity.md#Port%20sortant%20—%20CAPort%20/%20RequestStore%20%28identity%29) → [CAPort.CertAttributes](../../docs/code/fonctions/identity.CAPort.CertAttributes.md) · [CAPort.UpdateAttributes](../../docs/code/fonctions/identity.CAPort.UpdateAttributes.md)
- [§ Port entrant/sortant — RoleService / Repo (role)](DC_D1_Auth_Identity.md#Port%20entrant/sortant%20—%20RoleService%20/%20Repo%20%28role%29) → [RoleService.HasPermission](../../docs/code/fonctions/role.RoleService.HasPermission.md)
- [§ 6. Décisions de conception](DC_D1_Auth_Identity.md#6.%20Décisions%20de%20conception) → [REST /api/components](../../docs/code/routes/REST%20api-components.md) · [REST /api/modules](../../docs/code/routes/REST%20api-modules.md) · [CLI myr session create](../../docs/code/commandes/CLI%20myr-session-create.md) · [RoleService.HasPermission](../../docs/code/fonctions/role.RoleService.HasPermission.md)

<!-- liens-obsidian:end -->
