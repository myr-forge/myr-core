---
categorie: Développement autour de MYR
titre: "Utilisation de l'API"
probabilite: 1
impact: 0
importance: 0
etat: analyse
tags:
  - couche/analyse
  - type/use-case
  - famille/UCDEV
  - uc/UCDEV01
  - rm/RM01
  - rm/RM07
  - rm/RM08
  - enf/ENF12
---

# Utilisation de l'API

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Développeur" as Dev

rectangle "API REST MYR\n(adapters/in/rest/)" {
    usecase "S'authentifier à l'API" as UC1
    usecase "Lire des données (GET)" as UC2
    usecase "Écrire des données (POST/PUT/DELETE)" as UC3
    usecase "Consulter la documentation API" as UC4
}

Dev --> UC1
Dev --> UC4
UC1 .> UC2 : <<include>>
UC1 .> UC3 : <<include>>

@enduml
```

## Contexte

Les développeurs tiers peuvent utiliser l'API REST MYR pour intégrer ses fonctionnalités dans leurs applications (boutiques, éditeurs 3D, plugins CAO, systèmes de gestion). L'API REST est exposée sous le préfixe `/api/` par `myr-api`.

L'API (`adapters/in/rest/`) couvre les opérations sur les composants, modules, liaisons, canaux, réseaux, licences et l'identité (enrôlement CA, sessions par token opaque — pas de JWT, voir `specs/3-Conception/DC_D1_Auth_Identity.md`).

L'acteur "Développeur" est distinct de l'Administrateur : il accède à MYR exclusivement via l'API REST depuis son environnement (jamais via le CLI qui est réservé au serveur). Il possède un compte avec un rôle approprié (minimum Lecteur pour les lectures, Concepteur pour les écritures).

## Pré-conditions

- `myr-api` est en cours d'exécution et accessible
- Le développeur dispose d'un compte MYR avec les droits nécessaires (rôle Lecteur minimum)
- Le développeur dispose d'un token de session opaque valide (obtenu via `POST /api/identity/session` ou `POST /api/identity/guest`)
- La documentation de l'API REST (spec OpenAPI dans `api/`) est disponible

## Scénario

**Étape initiale :** Le développeur consulte la documentation API MYR

### Flux nominal — Lecture de données (query)

1. Le développeur obtient un token de session via `POST /api/identity/session` (secret CA) ou `POST /api/identity/guest` (accès invité)
2. Il effectue un appel GET avec l'en-tête `X-Myr-Token: <token>`
3. Exemple : `GET /api/components?channel=greenchannel`
4. Le serveur vérifie le token (`requireAuth`)
5. Le handler retourne les données en JSON
6. Le développeur intègre les données dans son application tierce

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 79 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : teste route `/api/components/`
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : teste route `/api/components/`
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : teste route `/api/components/`
> - 🟡 [TestChannels_PUT_SwitchesChannel](../../../docs/tests/adapters-in-rest/TestChannels_PUT_SwitchesChannel.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestComponentInterfaces_GET](../../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponentInterfaces_POST_Created](../../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponentTree_GET](../../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponentTree_MethodNotAllowed](../../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_AddAssembly](../../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_AddInstance](../../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_AddTwoInstancesSequentially](../../../docs/tests/adapters-in-rest/TestComponent_AddTwoInstancesSequentially.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_DELETE_NoContent](../../../docs/tests/adapters-in-rest/TestComponent_DELETE_NoContent.md) — déduit : teste route `/api/components/`
> - … et 64 autre(s) : voir la [matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux nominal — Écriture de données (invoke)

1. Le développeur dispose d'un token de session dont le rôle porte la permission `write` (rôle `contributor` par défaut — voir `specs/3-Conception/DC_D1_Auth_Identity.md` et `specs/roadmap_dev.md` § Écarts Identité & Session pour l'état de la synchronisation avec le rôle CA)
2. Il effectue un appel POST/PUT/PATCH/DELETE avec l'en-tête `X-Myr-Token`
3. Exemple : `POST /api/components` avec body JSON
4. Le serveur vérifie le token ET la permission (`requireRole(rbac.PermWrite, ...)`)
5. La validation métier est effectuée côté serveur
6. Si valide : la transaction est soumise à Fabric ; la réponse JSON confirme le succès

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 69 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : teste route `/api/components/`
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : teste route `/api/components/`
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponentInterfaces_GET](../../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponentInterfaces_POST_Created](../../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponentTree_GET](../../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponentTree_MethodNotAllowed](../../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_AddAssembly](../../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_AddInstance](../../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_AddTwoInstancesSequentially](../../../docs/tests/adapters-in-rest/TestComponent_AddTwoInstancesSequentially.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_DELETE_NoContent](../../../docs/tests/adapters-in-rest/TestComponent_DELETE_NoContent.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_DELETE_ServiceError](../../../docs/tests/adapters-in-rest/TestComponent_DELETE_ServiceError.md) — déduit : teste route `/api/components/`
> - … et 54 autre(s) : voir la [matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux alternatif — Token expiré ou invalide

1. Le développeur effectue un appel avec un token expiré ou révoqué
2. Le serveur retourne `401 Unauthorized`
3. **Il n'existe pas de mécanisme de rafraîchissement** — le développeur doit ré-obtenir un nouveau token via `POST /api/identity/session` (nouvel enrôlement CA) ou `POST /api/identity/guest`
4. Le développeur réessaie l'appel avec le nouveau token

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 10 déduit(s)
> - 🟡 [TestChannels_PUT_SwitchesChannel](../../../docs/tests/adapters-in-rest/TestChannels_PUT_SwitchesChannel.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentityGuest_AllowedDeliversToken](../../../docs/tests/adapters-in-rest/TestHandleIdentityGuest_AllowedDeliversToken.md) — déduit : teste route `/api/identity/guest`
> - 🟡 [TestHandleIdentityGuest_DisallowedReturnsForbidden](../../../docs/tests/adapters-in-rest/TestHandleIdentityGuest_DisallowedReturnsForbidden.md) — déduit : teste route `/api/identity/guest`
> - 🟡 [TestHandleIdentitySession_EnrollError](../../../docs/tests/adapters-in-rest/TestHandleIdentitySession_EnrollError.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_InvalidJSON](../../../docs/tests/adapters-in-rest/TestHandleIdentitySession_InvalidJSON.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_MethodNotAllowed](../../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MethodNotAllowed.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_MissingFields](../../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MissingFields.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_NoIdentityService](../../../docs/tests/adapters-in-rest/TestHandleIdentitySession_NoIdentityService.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_Success](../../../docs/tests/adapters-in-rest/TestHandleIdentitySession_Success.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestRequireAuth_BlockchainWithValidToken](../../../docs/tests/adapters-in-rest/TestRequireAuth_BlockchainWithValidToken.md) — déduit : teste route `/api/identity/session`
<!-- tests-obsidian:end -->

### Flux erreur — Droits insuffisants

1. Le développeur tente une écriture avec un rôle Lecteur
2. Le serveur retourne `403 Forbidden` : "Droits insuffisants"
3. Aucune modification n'est effectuée

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux erreur — Ressource introuvable

1. Le développeur appelle `GET /api/components/{id}` avec un ID inexistant
2. Le serveur retourne `404 Not Found`
3. Un corps JSON structuré décrit l'erreur

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 69 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : teste route `/api/components/`
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : teste route `/api/components/`
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponentInterfaces_GET](../../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponentInterfaces_POST_Created](../../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponentTree_GET](../../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponentTree_MethodNotAllowed](../../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_AddAssembly](../../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_AddInstance](../../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_AddTwoInstancesSequentially](../../../docs/tests/adapters-in-rest/TestComponent_AddTwoInstancesSequentially.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_DELETE_NoContent](../../../docs/tests/adapters-in-rest/TestComponent_DELETE_NoContent.md) — déduit : teste route `/api/components/`
> - 🟡 [TestComponent_DELETE_ServiceError](../../../docs/tests/adapters-in-rest/TestComponent_DELETE_ServiceError.md) — déduit : teste route `/api/components/`
> - … et 54 autre(s) : voir la [matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

## Post-conditions

- Pour une lecture : les données sont retournées en JSON, prêtes à être intégrées
- Pour une écriture : la ressource est créée/modifiée sur la blockchain Fabric ; l'ID est retourné
- Le token de session reste valide pour les appels suivants jusqu'à son expiration (7 jours) ou sa révocation par un administrateur

## Diagramme de séquence

```plantuml
@startuml
title UCDEV01 — Utilisation de l'API REST

participant "Développeur\n(HTTP client / terminal)" as Dev
participant "REST API\n(adapters/in/rest/)" as REST
participant "Service Domaine\n(domain/*/)" as Service

Dev -> REST : POST /api/identity/session {name, secret, org_id}
REST --> Dev : {token, role, pseudo, channel}

Dev -> REST : GET /api/components?channel=X\nX-Myr-Token: <token>
REST -> REST : Vérifier token (requireAuth)
REST -> Service : List(channelID)
Service --> REST : []Component
REST --> Dev : 200 OK — [{id, name, ...}]

Dev -> REST : POST /api/components {name, channelID, tags…}\nX-Myr-Token: <token>
REST -> REST : Vérifier token + permission write
REST -> Service : Add(params)
Service --> REST : Component créé

alt Succès
    REST --> Dev : 201 Created — {id, name, ...}
else Validation KO (ex: anti-plagiat)
    REST --> Dev : 422 Unprocessable — {error: "duplicate_hash", message: "..."}
else Token expiré ou révoqué
    REST --> Dev : 401 Unauthorized
    Dev -> REST : POST /api/identity/session {name, secret, org_id}\n(pas de rafraîchissement — nouvel enrôlement requis)
    REST --> Dev : {token, role, pseudo, channel}
    Dev -> REST : Réessayer l'appel avec le nouveau token
end
@enduml
```

## Règles métier déclenchées

- **EF55** (ENF55) : API REST disponible pour les intégrations tierces
- **RM07** : Validation côté serveur avant toute soumission blockchain
- **ENF12** : Vérification du rôle côté serveur (jamais côté client uniquement)
- **RM01** : Anti-plagiat SHA-256 + similarité SCM > 50 % avant ajout d'asset `base`
- **RM08** : Vérification compatibilité de licence si `ParentID != ""` et `LicenseID != ""`

## Exigences non-fonctionnelles

- L'API doit répondre en moins de 500 ms pour les opérations de lecture (hors Fabric)
- Le format JSON des erreurs doit être standardisé (code machine + message lisible)
- L'API est versionnée implicitement — les changements cassants doivent être anticipés
- La documentation de l'API doit couvrir toutes les routes exposées dans `server.go`

## Notes d'implémentation

**Routes exposées** dans `adapters/in/rest/server.go` :
- Identité : `/api/identity/policy`, `/api/identity/enroll`, `/api/identity/session`, `/api/identity/guest`, `/api/identity/request`, `/api/identity/requests`, `/api/identity/wallets`, `/api/identity/status`
- Données : `/api/components`, `/api/components/`, `/api/connections`, `/api/connections/`, `/api/assembly-links`, `/api/virtual-connect`, `/api/modules`, `/api/modules/`, `/api/interfaces/`
- Référentiel : `/api/refs`, `/api/refs/categories`, `/api/refs/types`, `/api/refs/units`
- Réseau : `/api/channels`, `/api/networks`, `/api/networks/active`
- Utilitaire : `/api/ping`, `/api/status`, `/api/health`, `/metrics`
- Licences : `/api/licenses`, `/api/licenses/`
- Admin : `/api/admin/sessions`, `/api/admin/sessions/`

**Authentification :** token opaque de session (`X-Myr-Token`) via `requireAuth()` — pas de JWT. Accès en lecture : permission `read` (rôle Lecteur minimum). Accès en écriture (POST/PUT/PATCH/DELETE) : permission `write` (rôle Concepteur `contributor` par défaut) via le middleware `contrib`, vérifiée par `domain/role.RoleService.HasPermission`.

**Documentation API :** `api/` contient les specs OpenAPI — la manière de la publier (portail statique, Swagger UI, etc.) relève du dépôt GUI externe ou de tout autre outil de documentation, hors périmètre de `myr`.

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Étape suivante — conception**
- _Aucun document de conception ne traite ce use case : la chaîne s'arrête à l'analyse._

**Code cité par cette analyse (sans conception : non relié)**
- REST `/api/identity/session`
- REST `/api/identity/guest`
- REST `/api/components/`
- REST `/api/identity/policy`
- REST `/api/identity/enroll`
- REST `/api/identity/request`
- REST `/api/identity/requests`
- REST `/api/identity/wallets`
- REST `/api/identity/status`
- REST `/api/connections/`
- REST `/api/assembly-links`
- REST `/api/virtual-connect`
- REST `/api/modules/`
- REST `/api/interfaces/`
- REST `/api/refs`
- REST `/api/refs/categories`
- REST `/api/refs/types`
- REST `/api/refs/units`
- REST `/api/channels`
- REST `/api/networks`
- REST `/api/networks/active`
- REST `/api/ping`
- REST `/api/status`
- REST `/api/health`
- REST `/api/licenses/`
- REST `/api/admin/sessions/`
- `RoleService.HasPermission`

<!-- liens-obsidian:end -->
