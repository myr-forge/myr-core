---
tags:
  - couche/conception
  - type/conception
  - domaine/identity
  - domaine/role
  - uc/UCA01
  - uc/UCA02
  - uc/UCA04
  - uc/UCA07
  - uc/UCA08
---
# DC — CLI Identité & Session : Référence des commandes `myr identity` / `myr session`

> Phase 3 — Arrington | Use cases : UCA01–02, UCA04, UCA07–08 | Outil : `myr` (`bin/myr-cli`)

---

## 1. Objectif

Ce document définit le **contrat d'interface CLI** couvrant les capacités du domaine `identity` (compte CA, wallets locaux, demandes d'accès) et du domaine `session` (compte local de la machine exécutant le CLI). Il joue pour ces use cases le rôle que `DC_CLI_Admin.md` joue pour l'administration réseau et `DC_CLI_Model.md` pour les composants/modules : référence unique du nommage des commandes.

**Parité CLI/REST — nuance propre à ce domaine :** contrairement à `model`, `network` ou `channel`, l'identité n'a pas un contrat CLI/REST parfaitement symétrique par construction :
- `myr identity wallets/status/enroll/request/requests/set-role` ont un équivalent REST direct (`GET /api/identity/wallets`, `GET /api/identity/status`, `POST /api/identity/enroll`, `POST /api/identity/request`, `GET /api/identity/requests`) — même service domaine des deux côtés.
- `myr identity re-enroll` n'a pas de route REST dédiée : côté REST, le ré-enrôlement se produit implicitement à chaque `POST /api/identity/session` (qui appelle `Enroll`, pas `ReEnroll`) — voir `DC_D1_Auth_Identity.md` § Session (`domain/identity`).
- `myr session` (`domain/session`, compte local **machine**) n'a **aucun** équivalent REST par conception : c'est un concept propre à la machine qui exécute le CLI, distinct des sessions REST (jeton `X-Myr-Token`) — voir `DC_D1_Auth_Identity.md` § 1.

**Ce que ce document ne couvre pas :** les sessions REST elles-mêmes (`POST /api/identity/session`, `POST /api/identity/guest`) sont un mécanisme propre à l'adapter REST, sans commande CLI équivalente (le CLI, exécuté par l'administrateur sur le serveur, n'a pas besoin d'un jeton opaque pour s'authentifier auprès de lui-même) — voir UCA02.

---

## 2. Arbre de commandes

```
myr identity
├── wallets                              — lister les wallets locaux                         —
├── status --handle <pseudo@org>         — interroger le statut CA d'un wallet                UCA04, UCA07 (voir note ci-dessous)
├── enroll                               — enrôler avec un secret, sauvegarder le wallet local UCA02 (variante sans session REST)
│     --name <pseudo> --secret <secret> --org-id <id>
├── re-enroll <pseudo@org>               — renouveler le certificat d'un wallet local          UCA08 (suite d'un changement de rôle)
├── request                              — soumettre une demande d'accès pour un tiers          UCA01
│     --pseudo <p> --email <e> --org-id <id> [--display-name <n>] [--message <m>]
├── requests                             — lister les demandes d'accès en attente               UCA01
└── set-role --id <pseudo@org> --role <nom>  — changer le rôle CA d'une identité                UCA08

myr session                              — compte local de la machine exécutant le CLI (bootstrap, un seul enregistrement)
├── create --name <n> [--org-id <id>]    — créer le compte local (premier lancement)
├── show                                 — afficher le compte local, s'il existe
└── logout                               — effacer le compte local
```

> **Note UCA04/UCA07 :** `myr identity status` interroge la Fabric CA pour le statut d'enrôlement d'un wallet (`pending` / `active` / `suspended`) — c'est l'équivalent direct de `GET /api/identity/status`. Les use cases UCA04 (« Vérification de connexion ») et UCA07 (« Vérification du rôle attribué ») citent cette commande, mais décrivent des besoins plus larges (détails de connexion réseau/organisation/peer pour UCA04, rôle RBAC pour UCA07) que `GetStatus` seul ne couvre pas — voir l'annotation `#remarque` posée dans ces deux fichiers.

---

## 3. Détail des commandes

### 3.1 `myr identity wallets`

```
myr identity wallets
```

**Comportement :** Liste les wallets présents localement (`~/.Myr/wallets/`).

**Sortie :**
```
HANDLE          PSEUDO   ORG          STATUT
alice@Org1MSP   alice    Org1MSP      active
```

**Service :** `identitySvc.ListLocalWallets()`

---

### 3.2 `myr identity status`

```
myr identity status --handle <pseudo@org>
```

**Comportement :** Découpe le handle en `pseudo` + `org` (suffixe `MSP` ajouté automatiquement), puis interroge la CA.

**Sortie :** `alice@Org1 : active`

**Erreurs :** `handle invalide — format attendu : pseudo@org`

**Service :** `identitySvc.GetStatus(ctx, WalletEntry{Handle, Name, OrgID})`

---

### 3.3 `myr identity enroll`

```
myr identity enroll --name <pseudo> --secret <secret> --org-id <id>
```

**Comportement :** Enrôle l'identité (déjà connue de la CA) avec le secret d'enrôlement fourni et sauvegarde le wallet localement — sans créer de session REST. Équivalent CLI de `POST /api/identity/enroll`.

**Sortie :** `Identité enrôlée : alice@Org1MSP  statut=pending`

**Service :** `identitySvc.Enroll(ctx, name, secret, orgID)`

---

### 3.4 `myr identity re-enroll`

```
myr identity re-enroll <pseudo@org>
```

**Comportement :** Recherche le wallet local correspondant au handle (`ListLocalWallets`), puis demande à la CA un nouveau certificat. Nécessaire après un changement de rôle (`myr identity set-role`, UCA08) pour que le nouveau rôle prenne effet dans le certificat actif.

**Sortie :** `Identité alice@Org1MSP ré-enrôlée : statut=active`

**Erreurs :** `aucun wallet local pour "<handle>" — utilisez d'abord 'myr identity enroll'`

**Service :** `identitySvc.ListLocalWallets()` puis `identitySvc.ReEnroll(ctx, wallet)`

---

### 3.5 `myr identity request`

```
myr identity request --pseudo <p> --email <e> --org-id <id> [--display-name <n>] [--message <m>]
```

**Comportement :** Équivalent CLI exact de `POST /api/identity/request` (UCA01) — un administrateur soumet une demande d'accès pour le compte d'un tiers. Si le réseau actif a `AllowAutoRegister=true`, la demande est immédiatement transformée en identité active et le secret d'enrôlement est affiché ; sinon la demande reste `pending`.

**Sortie (auto-enregistrement) :**
```
Identité créée automatiquement : alice
Secret d'enrôlement : <secret>
```

**Sortie (en attente) :**
```
Demande enregistrée (id=req-123, statut=pending) — un administrateur doit créer l'identité manuellement.
```

**Service :** `identitySvc.SubmitRequest(req)` puis, si le réseau l'autorise, `identitySvc.AutoRegister(ctx, req, profile.AutoRegisterRole)` (`networkSvc.GetActive()` pour lire la politique du réseau actif).

---

### 3.6 `myr identity requests`

```
myr identity requests
```

**Comportement :** Liste les demandes d'accès en attente. Équivalent CLI de `GET /api/identity/requests`.

**Service :** `identitySvc.ListRequests()`

---

### 3.7 `myr identity set-role`

Voir UCA08 (`specs/2-Analyse/UCA-Compte_et_Acces/UCA08.md`) pour le détail complet — déjà implémenté avant ce document.

```
myr identity set-role --id <pseudo@org> --role <nom>
```

**Service :** `identitySvc.SetRole(ctx, id, role)`

---

### 3.8 `myr session create` / `show` / `logout`

```
myr session create --name <n> [--org-id <id>]
myr session show
myr session logout
```

**Comportement :** Gère le compte local **de la machine** qui exécute le CLI (`{Name, OrgID, CreatedAt}`, un seul enregistrement JSON) — sans rapport avec l'identité CA ni les sessions REST (voir `DC_D1_Auth_Identity.md` § 1). `create` échoue si un compte existe déjà (utiliser `logout` avant de le remplacer).

**Sortie (`show`, aucun compte) :** `Aucun compte local — utilisez 'myr session create'.`

**Service :** `sessionSvc.Create(name, orgID)` / `sessionSvc.Current()` / `sessionSvc.Logout()`

---

## 4. Points ouverts

| # | Écart / question | Impact |
|---|---|---|
| 1 | Aucun mécanisme (CLI ou REST) ne permet d'approuver *a posteriori* une `AccountRequest` restée `pending` (réseau sans auto-enregistrement) — voir UCA01 § Notes d'implémentation. `myr identity request` ne fait qu'imiter la logique d'auto-approbation au moment de la soumission, comme l'endpoint REST. | L'administrateur doit toujours créer l'identité manuellement via l'outillage Fabric CA hors `myr` pour traiter une demande déjà en attente. |
| 2 | `IdentityService.Register(ctx, RegisterRequest)` (enregistrement direct en CA, sans passer par une `AccountRequest`) existe côté domaine mais n'est exposé ni en CLI ni en REST — `UCA01.md` (1-Expression) qualifie l'enregistrement manuel CA d'« hors périmètre applicatif ». Tension entre ce que permettrait déjà le code et ce que la spec UC décrit comme hors périmètre — à trancher avant d'exposer une éventuelle commande `myr identity register`. | Pas d'implémentation CLI/REST tant que ce point n'est pas tranché par le product owner. |
| 3 | `IdentityService.LoadGuestWallet(certPath, keyPath, caCertPath, orgID)` (chargement d'un wallet invité pré-enrôlé par l'admin) n'est exposé ni en CLI ni en REST. | Le provisionnement initial du wallet `guest` partagé (voir `GuestHandle` dans `adapters/in/rest/handlers_identity.go`) reste une opération manuelle hors `myr`. |
| 4 | `IdentityService.WalletDir()` (chemin du répertoire de wallets) n'est exposé nulle part — utilité limitée en dehors d'un contexte de diagnostic/support. | Mineur — à exposer seulement si un besoin de diagnostic CLI se présente. |

---

## 5. Décisions de conception

| ID | Décision | Raison |
|----|---------|--------|
| DC-CLII-01 | `myr identity request` reproduit l'orchestration (soumission + auto-enregistrement conditionnel) déjà présente dans `handleIdentityRequest` (REST), plutôt que de se limiter à `SubmitRequest` | Nécessaire pour une parité fonctionnelle réelle avec `POST /api/identity/request` (UCA01) — un administrateur agissant pour le compte d'un tiers doit obtenir le même résultat, quel que soit le canal. Cette orchestration (lire la politique réseau puis décider d'auto-enregistrer) vit aujourd'hui dans les deux adapters `in/` plutôt que dans le service domaine — dupplication à surveiller, cf. règle de centralisation (28), non résolue par ce document. |
| DC-CLII-02 | `myr session` n'a aucune route REST | `domain/session` modélise le compte local de la machine qui exécute le CLI — un client REST n'a pas de « machine locale » à ce sens, cf. `DC_D1_Auth_Identity.md` § 1. |
| DC-CLII-03 | `myr identity register` et `myr identity guest-wallet load` ne sont pas exposés dans cette itération | Points ouverts (§4, #2 et #3) nécessitant une décision produit avant exposition — pas un oubli d'adapter. |

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Étape suivante — code, par section**
- [§ 1. Objectif](DC_CLI_Identity.md#1.%20Objectif) → [REST /api/identity/wallets](../../docs/code/routes/REST%20api-identity-wallets.md) · [REST /api/identity/status](../../docs/code/routes/REST%20api-identity-status.md) · [REST /api/identity/enroll](../../docs/code/routes/REST%20api-identity-enroll.md) · [REST /api/identity/request](../../docs/code/routes/REST%20api-identity-request.md) · [REST /api/identity/requests](../../docs/code/routes/REST%20api-identity-requests.md) · [REST /api/identity/session](../../docs/code/routes/REST%20api-identity-session.md) · [REST /api/identity/guest](../../docs/code/routes/REST%20api-identity-guest.md) · [CLI myr identity wallets](../../docs/code/commandes/CLI%20myr-identity-wallets.md) · [CLI myr identity re-enroll](../../docs/code/commandes/CLI%20myr-identity-re-enroll.md) · [CLI myr session](../../docs/code/commandes/CLI%20myr-session.md) · [CAPort.Enroll](../../docs/code/fonctions/identity.CAPort.Enroll.md) · [CAPort.ReEnroll](../../docs/code/fonctions/identity.CAPort.ReEnroll.md) · [IdentityService.Enroll](../../docs/code/fonctions/identity.IdentityService.Enroll.md) · [IdentityService.ReEnroll](../../docs/code/fonctions/identity.IdentityService.ReEnroll.md)
- [§ 2. Arbre de commandes](DC_CLI_Identity.md#2.%20Arbre%20de%20commandes) → [REST /api/identity/status](../../docs/code/routes/REST%20api-identity-status.md) · [CLI myr identity](../../docs/code/commandes/CLI%20myr-identity.md) · [CLI myr session](../../docs/code/commandes/CLI%20myr-session.md) · [CLI myr identity status](../../docs/code/commandes/CLI%20myr-identity-status.md) · [CAPort.GetStatus](../../docs/code/fonctions/identity.CAPort.GetStatus.md) · [IdentityService.GetStatus](../../docs/code/fonctions/identity.IdentityService.GetStatus.md)
- [§ 3.1 `myr identity wallets`](DC_CLI_Identity.md#3.1%20`myr%20identity%20wallets`) → [CLI myr identity wallets](../../docs/code/commandes/CLI%20myr-identity-wallets.md) · [IdentityService.ListLocalWallets](../../docs/code/fonctions/identity.IdentityService.ListLocalWallets.md)
- [§ 3.2 `myr identity status`](DC_CLI_Identity.md#3.2%20`myr%20identity%20status`) → [CLI myr identity status](../../docs/code/commandes/CLI%20myr-identity-status.md) · [CAPort.GetStatus](../../docs/code/fonctions/identity.CAPort.GetStatus.md) · [IdentityService.GetStatus](../../docs/code/fonctions/identity.IdentityService.GetStatus.md)
- [§ 3.3 `myr identity enroll`](DC_CLI_Identity.md#3.3%20`myr%20identity%20enroll`) → [REST /api/identity/enroll](../../docs/code/routes/REST%20api-identity-enroll.md) · [CLI myr identity enroll](../../docs/code/commandes/CLI%20myr-identity-enroll.md) · [CAPort.Enroll](../../docs/code/fonctions/identity.CAPort.Enroll.md) · [IdentityService.Enroll](../../docs/code/fonctions/identity.IdentityService.Enroll.md)
- [§ 3.4 `myr identity re-enroll`](DC_CLI_Identity.md#3.4%20`myr%20identity%20re-enroll`) → [CLI myr identity re-enroll](../../docs/code/commandes/CLI%20myr-identity-re-enroll.md) · [CLI myr identity set-role](../../docs/code/commandes/CLI%20myr-identity-set-role.md) · [CLI myr identity enroll](../../docs/code/commandes/CLI%20myr-identity-enroll.md) · [CAPort.ReEnroll](../../docs/code/fonctions/identity.CAPort.ReEnroll.md) · [IdentityService.ListLocalWallets](../../docs/code/fonctions/identity.IdentityService.ListLocalWallets.md) · [IdentityService.ReEnroll](../../docs/code/fonctions/identity.IdentityService.ReEnroll.md)
- [§ 3.5 `myr identity request`](DC_CLI_Identity.md#3.5%20`myr%20identity%20request`) → [REST /api/identity/request](../../docs/code/routes/REST%20api-identity-request.md) · [CLI myr identity request](../../docs/code/commandes/CLI%20myr-identity-request.md) · [IdentityService.AutoRegister](../../docs/code/fonctions/identity.IdentityService.AutoRegister.md) · [IdentityService.SubmitRequest](../../docs/code/fonctions/identity.IdentityService.SubmitRequest.md) · [NetworkService.GetActive](../../docs/code/fonctions/network.NetworkService.GetActive.md)
- [§ 3.6 `myr identity requests`](DC_CLI_Identity.md#3.6%20`myr%20identity%20requests`) → [REST /api/identity/requests](../../docs/code/routes/REST%20api-identity-requests.md) · [CLI myr identity requests](../../docs/code/commandes/CLI%20myr-identity-requests.md) · [IdentityService.ListRequests](../../docs/code/fonctions/identity.IdentityService.ListRequests.md)
- [§ 3.7 `myr identity set-role`](DC_CLI_Identity.md#3.7%20`myr%20identity%20set-role`) → [CLI myr identity set-role](../../docs/code/commandes/CLI%20myr-identity-set-role.md) · [IdentityService.SetRole](../../docs/code/fonctions/identity.IdentityService.SetRole.md)
- [§ 3.8 `myr session create` / `show` / `logout`](DC_CLI_Identity.md#3.8%20`myr%20session%20create`%20/%20`show`%20/%20`logout`) → [CLI myr session create](../../docs/code/commandes/CLI%20myr-session-create.md) · [CLI myr session show](../../docs/code/commandes/CLI%20myr-session-show.md) · [CLI myr session logout](../../docs/code/commandes/CLI%20myr-session-logout.md) · [RoleService.Create](../../docs/code/fonctions/role.RoleService.Create.md) · [SessionService.Current](../../docs/code/fonctions/session.SessionService.Current.md) · [SessionService.Logout](../../docs/code/fonctions/session.SessionService.Logout.md)
- [§ 4. Points ouverts](DC_CLI_Identity.md#4.%20Points%20ouverts) → [CLI myr identity request](../../docs/code/commandes/CLI%20myr-identity-request.md) · [CLI myr identity](../../docs/code/commandes/CLI%20myr-identity.md) · [CAPort.Register](../../docs/code/fonctions/identity.CAPort.Register.md) · [IdentityService.LoadGuestWallet](../../docs/code/fonctions/identity.IdentityService.LoadGuestWallet.md) · [IdentityService.Register](../../docs/code/fonctions/identity.IdentityService.Register.md) · [IdentityService.WalletDir](../../docs/code/fonctions/identity.IdentityService.WalletDir.md)
- [§ 5. Décisions de conception](DC_CLI_Identity.md#5.%20Décisions%20de%20conception) → [REST /api/identity/request](../../docs/code/routes/REST%20api-identity-request.md) · [CLI myr identity request](../../docs/code/commandes/CLI%20myr-identity-request.md) · [CLI myr session](../../docs/code/commandes/CLI%20myr-session.md) · [CLI myr identity](../../docs/code/commandes/CLI%20myr-identity.md) · [IdentityService.SubmitRequest](../../docs/code/fonctions/identity.IdentityService.SubmitRequest.md)

<!-- liens-obsidian:end -->
