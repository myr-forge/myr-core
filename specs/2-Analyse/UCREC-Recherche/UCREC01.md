---
categorie: Recherche
titre: "Rechercher une référence existante"
probabilite: 4
impact: 5
importance: 20
etat: analyse
tags:
  - couche/analyse
  - type/use-case
  - famille/UCREC
  - domaine/model
  - uc/UCREC01
  - rm/RM22
  - enf/ENF12
---

# Rechercher une référence existante

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Client\n(dépôt GUI externe)" as U

rectangle "API myr" {
    usecase "Rechercher une référence\n(UUID ou nom)" as UC1
}

U --> UC1

@enduml
```

## Contexte

La recherche par référence est le point d'entrée principal pour trouver un Composant ou un Module sur le réseau. Elle interroge la blockchain via l'API REST et retourne les assets correspondants au client (le dépôt GUI externe, typiquement, pour affichage et exploration ultérieure).

Cette fonctionnalité est accessible à tout utilisateur authentifié (`Lecteur`, `Concepteur`, `Consommateur`…). Elle est le use case de recherche le plus fréquent (importance 20 — la plus haute du domaine D8).

> La présentation des résultats (Explorer, Asset UI, MenuBar, etc.) relève du dépôt GUI externe — hors périmètre de ce document. Seuls les contrats REST et CLI ci-dessous font partie de `myr`.

## Pré-conditions

- Utilisateur authentifié (rôle `Lecteur` minimum)
- Connexion réseau active
- Au moins un asset publié sur la blockchain du réseau courant

## Scénario

### Flux nominal — Référence trouvée (UUID exact)

1. Le client envoie l'UUID ou la référence exacte du Composant/Module recherché
2. Le système appelle `GET /api/components?name=<ref>` ou `GET /api/modules?name=<ref>`
3. Service : `List(channelID)` puis filtrage côté handler sur `ID == ref` ou `Name == ref`
4. Le Composant ou Module correspondant est retourné avec : nom, catégorie, statut, description courte

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 75 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : teste route `/api/components/`
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : teste route `/api/components/`
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : teste route `/api/components/`
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : teste route `/api/modules/`
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
> - … et 60 autre(s) : voir la [matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux nominal — Recherche par nom (correspondance partielle)

1. Le client fournit un nom partiel (ex : "moteur")
2. Le système filtre les assets dont `Name` contient la chaîne (insensible à la casse)
3. La liste des résultats correspondants est retournée

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux nominal — Référence introuvable

1. Aucun asset ne correspond à la référence fournie
2. Le système retourne une liste vide

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux alternatif — Résultats mixtes (Composants ET Modules)

1. La recherche retourne à la fois des Composants et des Modules
2. Chaque résultat porte une indication de son type (Composant / Module)

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

## Post-conditions

- Le client reçoit la liste des assets correspondants
- Aucune modification de la blockchain

## Diagramme de séquence

```plantuml
@startuml
participant "Client\n(dépôt GUI externe)" as Client
participant "REST Handler\n(adapters/in/rest/)" as REST
participant "Model Service\n(domain/model/)" as Service
database "Fabric\n(adapters/out/fabric/)" as Fabric

Client -> REST : GET /api/components?name=<ref>&channel=<channelID>\nX-Myr-Token: <token>
REST -> REST : Vérifier session (ENF12)
REST -> Service : List(channelID)
Service -> Fabric : ListModelRecords(channelID)
Fabric --> Service : []*Model3D (tous les assets)
Service --> REST : []*Model3D
REST -> REST : Filtrer par Name/ID contenant <ref>\nExclure les modules (IsModule()==true)
REST --> Client : 200 [{id, name, category, status, ...}]

Client -> REST : GET /api/modules?name=<ref>&channel=<channelID>
REST -> Service : ListModules(channelID)
Service -> Fabric : ListModelRecords(channelID)
Fabric --> Service : []*Model3D
Service -> Service : Filtrer IsModule()==true
Service --> REST : []*Model3D modules
REST -> REST : Filtrer par Name/ID contenant <ref>
REST --> Client : 200 [{id, name, status, versions, ...}]
@enduml
```

## Règles métier déclenchées

| Règle | Description | Point d'application |
|-------|-------------|---------------------|
| **RM22** | Contrôle d'accès par rôle — lecture autorisée dès `Lecteur` | Middleware `requireAuth`/`requireRole` dans les handlers REST |

## Exigences non-fonctionnelles

- **ENF12** : Authentification par session (token opaque `X-Myr-Token`) obligatoire pour accéder à `/api/components` et `/api/modules`

## Notes d'implémentation

**Endpoints REST utilisés :**
- `GET /api/components` → `handleComponents()` (handlers.go:~441 — filtre les modules)
- `GET /api/modules` → `handleModules()` (handlers.go:~965)

**Filtrage actuel :** `handleComponents()` filtre les assets avec `m.IsModule() == true` (exclusion des modules — handlers.go:~441). Le filtrage par `name` est côté handler (boucle sur les résultats) — pas de requête filtrée côté Fabric. Pour les gros réseaux, un index de recherche côté chaincode est à envisager.

**Accès Visiteur :** La spec D4 indique que l'accès public (sans session) devrait être possible pour les assets publics. Le code actuel applique `requireAuth` systématiquement — l'accès visiteur reste à implémenter (écart documenté dans l'analyse D4).

**Commande CLI équivalente (alias limité) :** `myr model get <id>` (méthode `Get`) est l'équivalent direct d'une recherche par UUID/référence exacte. `myr model list [--channel <id>]` (méthode `List`) permet de parcourir les assets pour une recherche par nom partiel, en l'absence de méthode de filtre serveur dédiée dans `ModelService` (même limite que UCCL01). `myr module list` couvre le pendant module de `GET /api/modules?name=<ref>`.

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Étape suivante — conception**
- **Chaincode** : [§ 4. Fonctions chaincode — Store/Read (D3/D4/D6)](../../3-Conception/Chaincode.md#4.%20Fonctions%20chaincode%20—%20Store/Read%20%28D3/D4/D6%29)
- **DC_CLI_Model** : [§ DC — CLI Modèle : Référence des commandes composant / interfaces / module](../../3-Conception/DC_CLI_Model.md#DC%20—%20CLI%20Modèle%20:%20Référence%20des%20commandes%20composant%20/%20interfaces%20/%20module) · [§ 5. Table de correspondance méthode domaine → commande CLI → use case](../../3-Conception/DC_CLI_Model.md#5.%20Table%20de%20correspondance%20méthode%20domaine%20→%20commande%20CLI%20→%20use%20case) · [§ 6. Écarts et points ouverts](../../3-Conception/DC_CLI_Model.md#6.%20Écarts%20et%20points%20ouverts)
- **DC_D8_Recherche** : [§ DC — D8 : Recherche](../../3-Conception/DC_D8_Recherche.md#DC%20—%20D8%20:%20Recherche) · [§ 1. Objectif](../../3-Conception/DC_D8_Recherche.md#1.%20Objectif)

<!-- liens-obsidian:end -->
