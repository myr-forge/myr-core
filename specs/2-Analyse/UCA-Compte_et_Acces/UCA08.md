---
categorie: Compte et Accès
titre: "Demander un rôle"
probabilite: 4
impact: 4
importance: 16
etat: analyse
tags:
  - couche/analyse
  - type/use-case
  - famille/UCA
  - domaine/identity
  - domaine/role
  - uc/UCA08
  - rm/RM22
  - enf/ENF12
---

# Demander un rôle

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Utilisateur\n(tout rôle authentifié)" as U
actor "Administrateur" as ADM

rectangle "API myr" {
    usecase "Soumettre une demande\n(message libre)" as UC1
    usecase "Changer le rôle d'une identité\n(myr identity set-role)" as UC2
}

U --> UC1
ADM --> UC2

@enduml
```

## Contexte

**Il n'existe aucun mécanisme en libre-service pour qu'un utilisateur change son propre rôle.** Le changement de rôle est une opération strictement administrateur :

```
myr identity set-role --id <pseudo@org> --role <nom-du-rôle>
```

Cette commande CLI modifie l'attribut `Myr.role` de l'identité auprès de la Fabric CA (`identitySvc.SetRole` → `CAPort.UpdateAttributes`). **Le nouveau rôle ne s'applique qu'au prochain ré-enrôlement de l'identité** (propriété de la Fabric CA, pas une limitation de `myr`) — l'utilisateur doit donc se reconnecter (UCA02) après le changement.

Le rôle demandé doit exister dans le catalogue RBAC (`domain/role`) — soit l'un des 4 rôles intégrés (`reader`, `contributor`, `auditor`, `admin`), soit un rôle personnalisé déjà créé par un administrateur (`myr role create`, voir UCA05).

Un utilisateur ne peut transmettre une requête d'accès qu'en texte libre : le champ `message` de `POST /api/identity/request` (UCA01) ou de `POST /api/identity/guest` (UCA02) permet d'écrire une demande, mais **il n'existe pas de champ structuré « rôle souhaité » dans `AccountRequest`** — l'administrateur doit lire le message et agir manuellement via `myr identity set-role`.

> ⚠️ En complément, le rôle de session REST étant actuellement figé à `"contributor"` à chaque connexion (voir écart documenté dans UCA02), un changement de rôle effectué via `myr identity set-role` **n'a aujourd'hui aucun effet visible sur la session REST** obtenue via `POST /api/identity/session` — seul un usage CLI direct de l'identité (hors REST) refléterait le nouveau rôle CA.

## Pré-conditions

- L'identité existe déjà auprès de la CA (UCA01 complété)
- L'administrateur dispose des droits pour exécuter `myr identity set-role` (accès CLI serveur)

## Scénario

### Flux nominal — Demande informelle puis changement manuel

1. L'utilisateur transmet sa demande de rôle par le champ `message` libre d'une requête `POST /api/identity/request` ou `POST /api/identity/guest`, ou par un canal hors `myr` (email, ticket, etc.)
2. L'administrateur consulte les demandes en attente via `GET /api/identity/requests`
3. L'administrateur exécute `myr identity set-role --id <pseudo@org> --role <nouveau-rôle>`
4. Le service appelle `CAPort.UpdateAttributes(ctx, name, {"Myr.role": newRole})`
5. Le CLI affiche : « Rôle de "pseudo@org" mis à jour : nouveau-rôle. Le nouveau rôle s'applique au prochain ré-enrôlement de l'identité. »
6. L'identité doit être ré-enrôlée pour qu'un nouveau certificat portant l'attribut à jour soit émis — soit côté CLI via `myr identity re-enroll <pseudo@org>` (renouvelle le wallet local sans passer par une session REST), soit côté REST via une nouvelle connexion (`POST /api/identity/session`). Dans ce second cas, **cela ne met pas à jour le rôle de la session REST**, actuellement toujours fixé à `"contributor"` (voir écart UCA02)

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 17 déduit(s)
> - 🟡 [TestChannels_PUT_SwitchesChannel](../../../docs/tests/adapters-in-rest/TestChannels_PUT_SwitchesChannel.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentityGuest_AllowedDeliversToken](../../../docs/tests/adapters-in-rest/TestHandleIdentityGuest_AllowedDeliversToken.md) — déduit : teste route `/api/identity/guest`
> - 🟡 [TestHandleIdentityGuest_DisallowedReturnsForbidden](../../../docs/tests/adapters-in-rest/TestHandleIdentityGuest_DisallowedReturnsForbidden.md) — déduit : teste route `/api/identity/guest`
> - 🟡 [TestHandleIdentityRequest_InvalidJSON](../../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_InvalidJSON.md) — déduit : teste route `/api/identity/request`
> - 🟡 [TestHandleIdentityRequest_MissingFields](../../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_MissingFields.md) — déduit : teste route `/api/identity/request`
> - 🟡 [TestHandleIdentityRequest_ServiceError](../../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_ServiceError.md) — déduit : teste route `/api/identity/request`
> - 🟡 [TestHandleIdentityRequest_Success](../../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_Success.md) — déduit : teste route `/api/identity/request`
> - 🟡 [TestHandleIdentityRequests_Admin_EmptyList](../../../docs/tests/adapters-in-rest/TestHandleIdentityRequests_Admin_EmptyList.md) — déduit : teste route `/api/identity/requests`
> - 🟡 [TestHandleIdentityRequests_Admin_ReturnsList](../../../docs/tests/adapters-in-rest/TestHandleIdentityRequests_Admin_ReturnsList.md) — déduit : teste route `/api/identity/requests`
> - 🟡 [TestHandleIdentityRequests_NonAdmin_Forbidden](../../../docs/tests/adapters-in-rest/TestHandleIdentityRequests_NonAdmin_Forbidden.md) — déduit : teste route `/api/identity/requests`
> - 🟡 [TestHandleIdentitySession_EnrollError](../../../docs/tests/adapters-in-rest/TestHandleIdentitySession_EnrollError.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_InvalidJSON](../../../docs/tests/adapters-in-rest/TestHandleIdentitySession_InvalidJSON.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_MethodNotAllowed](../../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MethodNotAllowed.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_MissingFields](../../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MissingFields.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_NoIdentityService](../../../docs/tests/adapters-in-rest/TestHandleIdentitySession_NoIdentityService.md) — déduit : teste route `/api/identity/session`
> - … et 2 autre(s) : voir la [matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux erreur — Identité ou rôle invalide

1. `--id` ou `--role` manquant : cobra refuse la commande (`MarkFlagRequired`)
2. Rôle inexistant dans le catalogue RBAC : l'appel CA échoue ou le rôle reste sans effet RBAC tant qu'il n'est pas créé (`myr role create`)

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

## Post-conditions

- L'attribut `Myr.role` de l'identité CA est mis à jour
- Le nouveau rôle ne prend effet qu'au prochain enrôlement CA — pas immédiatement
- La session REST en cours (et toute nouvelle session créée via `POST /api/identity/session` tant que l'écart de rôle figé n'est pas corrigé) n'est pas affectée par ce changement

## Diagramme de séquence

```plantuml
@startuml
participant "Administrateur" as ADM
participant "CLI\n(adapters/in/cli/identity.go)" as CLI
participant "Identity Service\n(domain/identity/service.go)" as Service
participant "Fabric CA" as CA

ADM -> CLI : myr identity set-role --id alice@org1 --role contributor
CLI -> Service : SetRole(ctx, "alice@org1", "contributor")
Service -> CA : UpdateAttributes(ctx, "alice@org1", {"Myr.role": "contributor"})
CA --> Service : ok
Service --> CLI : nil
CLI --> ADM : "Rôle de alice@org1 mis à jour : contributor.\nLe nouveau rôle s'applique au prochain ré-enrôlement."
@enduml
```

## Règles métier déclenchées

- **RM22** — Le changement de rôle est réservé à l'administrateur ; aucun utilisateur ne peut se l'auto-attribuer.

## Exigences non-fonctionnelles

- **ENF12** — L'attribution de rôle est strictement côté serveur/CA — le client ne peut pas s'auto-attribuer un rôle (le token de session est opaque, non falsifiable).

## Notes d'implémentation

**Commande réelle :** `myr identity set-role --id <pseudo@org> --role <nom>` (`adapters/in/cli/identity.go`) → `IdentityService.SetRole` (`domain/identity/service.go`) → `CAPort.UpdateAttributes`.

**Ré-enrôlement CLI :** `myr identity re-enroll <pseudo@org>` (`adapters/in/cli/identity.go`) → `IdentityService.ReEnroll` (retrouve le wallet local correspondant au handle via `ListLocalWallets`, puis émet un nouveau certificat). Ferme l'écart qui obligeait auparavant à passer par une session REST (`POST /api/identity/session`) pour tout ré-enrôlement — utile en particulier quand l'identité n'a pas vocation à obtenir de session REST.

**⚠️ Écart — pas de canal structuré pour la demande :** `AccountRequest` (`domain/identity/entity.go`) n'a pas de champ « rôle souhaité » — seul un champ `message` libre existe. Une future itération pourrait ajouter un champ `requested_role` et un endpoint/commande d'approbation dédiés (voir écart similaire documenté dans UCA01, absence de flux d'approbation).

**⚠️ Écart — déconnexion entre rôle CA et rôle de session REST :** tant que `handleIdentitySession` fixera le rôle de session à `"contributor"` en dur (écart documenté dans UCA02), ce use case n'aura aucun effet observable via l'API REST — seul un usage direct de l'identité CA (CLI Fabric, ou un futur code REST corrigé) en bénéficierait.

**Pas d'auto-distribution configurée par réseau :** contrairement à un système de règles par réseau (auto vs validation), il n'existe qu'un seul mécanisme aujourd'hui : l'action manuelle de l'administrateur via `myr identity set-role`.

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Étape suivante — conception**
- **DC_CLI_Identity** : [§ DC — CLI Identité & Session : Référence des commandes `myr identity` / `myr session`](../../3-Conception/DC_CLI_Identity.md#DC%20—%20CLI%20Identité%20&%20Session%20:%20Référence%20des%20commandes%20`myr%20identity`%20/%20`myr%20session`) · [§ 2. Arbre de commandes](../../3-Conception/DC_CLI_Identity.md#2.%20Arbre%20de%20commandes) · [§ 3.4 `myr identity re-enroll`](../../3-Conception/DC_CLI_Identity.md#3.4%20`myr%20identity%20re-enroll`) · [§ 3.7 `myr identity set-role`](../../3-Conception/DC_CLI_Identity.md#3.7%20`myr%20identity%20set-role`)

<!-- liens-obsidian:end -->
