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
- `myr identity wallets/status/enroll/request/requests/approve/reject/register` ont un équivalent REST direct (`GET /api/identity/wallets`, `GET /api/identity/status`, `POST /api/identity/enroll`, `POST /api/identity/request`, `GET /api/identity/requests`, `POST /api/identity/requests/{id}/approve`, `POST /api/identity/requests/{id}/reject`, `POST /api/identity/accounts`) — même service domaine des deux côtés.
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
├── requests                             — lister les demandes d'accès et leur statut            UCA01
├── approve <request-id>                 — approuver une demande en attente (identité active)   UCA01
├── reject <request-id>                  — refuser une demande en attente (identité suspendue)  UCA01
├── register                             — créer un compte actif (administrateur)               UCA01
│     --pseudo <p> --org-id <id> [--role <r>] [--email <e>]
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

**Comportement :** Équivalent CLI exact de `POST /api/identity/request` (UCA01) — un administrateur soumet une demande d'accès pour le compte d'un tiers, selon les règles d'inscription du réseau actif. Réseau fermé (`AllowAutoRegister=false`, ou aucun réseau actif) : la commande échoue, rien n'est enregistré — utiliser `myr identity register`. Réseau ouvert : l'identité est créée avec le rôle par défaut du réseau, active ou en attente d'approbation (`RequireApproval`), et le secret d'enrôlement est affiché.

**Sortie (sans approbation) :**
```
Identité créée : alice@Org1 (demande req-123, statut=approved)
Secret d'enrôlement : <secret>
```

**Sortie (approbation exigée) :**
```
Identité créée : alice@Org1 (demande req-123, statut=pending)
Secret d'enrôlement : <secret>
L'identité reste inactive jusqu'à son approbation : myr identity approve req-123
```

**Service :** `identitySvc.RequestAccess(ctx, req, identity.PolicyOf(networkSvc.GetActive()))`.

---

### 3.5b `myr identity approve` / `reject`

```
myr identity approve <request-id>
myr identity reject <request-id>
```

**Comportement :** Décide une demande en attente (UCA01). `approve` passe l'identité à `Myr.status=active` ; `reject` la passe à `Myr.status=suspended` — rien n'est supprimé. La décision s'applique au prochain enrôlement de l'identité. Échoue si la demande est inconnue ou n'est plus en attente. Équivalents REST : `POST /api/identity/requests/{id}/approve|reject`.

**Sortie :**
```
Demande req-123 : approved (alice@Org1).
La décision s'applique au prochain enrôlement de l'identité.
```

**Service :** `identitySvc.ApproveRequest(ctx, id)` / `identitySvc.RejectRequest(ctx, id)`

---

### 3.5c `myr identity register`

```
myr identity register --pseudo <p> --org-id <id> [--role <r>] [--email <e>]
```

**Comportement :** Crée directement, à l'initiative de l'administrateur, une identité active avec le rôle donné (`reader` par défaut) — seule voie d'inscription sur un réseau fermé à l'auto-enregistrement (ADR-12). Le secret d'enrôlement est affiché une fois, à transmettre à l'utilisateur. Équivalent REST : `POST /api/identity/accounts`.

**Sortie :**
```
Identité créée : alice@Org1 (active)
Secret d'enrôlement : <secret>
```

**Service :** `identitySvc.CreateAccount(ctx, pseudo, orgID, role, email)`

---

### 3.6 `myr identity requests`

```
myr identity requests
```

**Comportement :** Liste les demandes d'accès et leur statut (`pending`, `approved`, `rejected`, `failed`). Équivalent CLI de `GET /api/identity/requests`.

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
| 1 | `IdentityService.LoadGuestWallet(certPath, keyPath, caCertPath, orgID)` (chargement d'un wallet invité pré-enrôlé par l'admin) n'est exposé ni en CLI ni en REST. | Le provisionnement initial du wallet `guest` partagé (voir `GuestHandle` dans `adapters/in/rest/handlers_identity.go`) reste une opération manuelle hors `myr`. |
| 2 | `IdentityService.WalletDir()` (chemin du répertoire de wallets) n'est exposé nulle part — utilité limitée en dehors d'un contexte de diagnostic/support. | Mineur — à exposer seulement si un besoin de diagnostic CLI se présente. |

---

## 5. Décisions de conception

| ID | Décision | Raison |
|----|---------|--------|
| DC-CLII-01 | `myr identity request` et `POST /api/identity/request` appellent le même `RequestAccess`, avec la même politique `identity.PolicyOf(réseau actif)` | Parité fonctionnelle (UCA01) — un administrateur agissant pour le compte d'un tiers obtient le même résultat, quel que soit le canal. Toute la décision (refus si fermé, statut actif ou en attente, rôle par défaut) vit dans le service domaine ; les adapters `in/` ne font que lire le réseau actif. |
| DC-CLII-02 | `myr session` n'a aucune route REST | `domain/session` modélise le compte local de la machine qui exécute le CLI — un client REST n'a pas de « machine locale » à ce sens, cf. `DC_D1_Auth_Identity.md` § 1. |
| DC-CLII-03 | `myr identity guest-wallet load` n'est pas exposé | Point ouvert (§4, #1) nécessitant une décision produit avant exposition — pas un oubli d'adapter. `myr identity register` est exposé : sur un réseau fermé à l'auto-enregistrement, c'est la seule voie d'inscription (ADR-12). |
| DC-CLII-04 | Le CLI et l'API partagent le stockage des demandes (`account-requests.json` du répertoire de données serveur) | Une demande reçue par l'API doit pouvoir être décidée par `myr identity approve/reject`, exécuté par l'administrateur sur le même serveur. |

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
