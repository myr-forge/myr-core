---
categorie: Module
titre: "Ajouter un Module existant"
probabilite: 3
impact: 5
importance: 15
etat: analyse
tags:
  - couche/analyse
  - type/use-case
  - famille/UCMOD
  - domaine/model
  - uc/UCMOD02
  - rm/RM13
  - rm/RM15
  - enf/ENF12
  - enf/ENF28
---

# Ajouter un Module existant

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Concepteur" as C
actor "Consommateur" as CL

rectangle "Application MYR" {
    usecase "Ajouter un module existant\ncomme instance" as UC1
    usecase "Rechercher le module" as UC2
    usecase "Créer seconde instance\nindépendante" as UC3
}

C --> UC1
CL --> UC1
UC1 ..> UC2 : <<include>>
UC1 .> UC3 : <<extend>>

@enduml
```

## Contexte

Un Concepteur ou Consommateur souhaite réutiliser un Module déjà existant (contrôleur, caméra, visserie, sous-assemblage…) comme instance dans son propre module hôte. Le Module source peut provenir :
- du réseau blockchain (Module soumis, visible publiquement)
- d'une boutique partenaire liée au réseau

La clé de cette opération est la gestion des **instances indépendantes** (RM15) : le même Module peut être instancié plusieurs fois dans un module hôte, chaque instance ayant ses propres Liaisons sans affecter les autres. Cette règle est symétrique avec l'ajout de Composants (UCAM01).

## Pré-conditions

- Identité authentifiée avec rôle **Concepteur** (`contributor`) ou **Consommateur** (`consumer`)
- Un module hôte existe, en état `draft` (voir UCMOD01)
- Le Module cible est accessible sur la blockchain ou une boutique partenaire (état `submitted`)

## Scénario

**Étape initiale :** Le Module source est identifié par référence ou filtre (voir UCREC01/UCCL01), puis `POST /api/modules/:hostModuleID/instances` est appelée avec son identifiant (ou l'équivalent CLI `myr model instance add`)

### Flux nominal — Module ajouté (première instance)

1. Le Module source est récupéré depuis la blockchain (`GET /api/modules/:id`)
2. `POST /api/modules/:hostModuleID/instances` est appelée avec `{asset_id: moduleSource.ID}`
3. Service : `AddAssetToWorkspace(hostModuleID, assetID)` crée une `WorkspaceInstance` indépendante
4. Le slot virtuel est garanti automatiquement (RM13)
5. Le Module source est instancié dans le module hôte, avec ses interfaces exposées

### Flux alternatif — Module déjà instancié dans le module hôte (RM15)

1. Le Module cible possède déjà au moins une `WorkspaceInstance` dans le module hôte
2. `POST /api/modules/:hostModuleID/instances` avec `{asset_id: sourceID, force: true}` crée une nouvelle `WorkspaceInstance` avec un nouvel ID unique
3. La nouvelle instance est indépendante — ses futures Liaisons n'impactent pas la première instance
4. Les deux instances sont consultables via `GET /api/modules/:hostModuleID/instances`

### Flux erreur — Module introuvable sur la blockchain

1. L'ID ou la référence transmise ne correspond à aucun asset sur la blockchain
2. Message : "Module introuvable — vérifiez la référence"
3. Le module hôte reste inchangé

### Flux erreur — Module non soumis (état draft d'un autre utilisateur)

1. Le Module cible existe mais est en état `draft` — il n'est pas visible sur le réseau
2. Le système retourne une erreur d'accès interdit
3. Message : "Ce module n'est pas encore publié sur le réseau"

## Post-conditions

- Une nouvelle `WorkspaceInstance` est créée dans le module hôte avec un ID unique
- L'instance référence `AssetID` du Module source (pas une copie)
- Le module hôte reste en état `draft`
- Chaque instance dispose d'au moins un slot virtuel (RM13)
- Les Liaisons entre instances restent indépendantes (RM15)

## Diagramme de séquence

```plantuml
@startuml
participant "Client\n(CLI ou API REST)" as Client
participant "REST Handler\n(adapters/in/rest/)" as REST
participant "Model Service\n(domain/model/)" as Service
database "LocalStorage\n(adapters/out/localstorage/)" as Local
database "Fabric\n(adapters/out/fabric/)" as Fabric

Client -> REST : GET /api/modules/:sourceID
REST -> Service : GetModule(sourceID)
Service -> Fabric : GetModelRecord(sourceID, channelID)
Fabric --> Service : *Model3D (submitted)
Service --> REST : *Model3D
REST --> Client : 200 moduleDTO

Client -> REST : POST /api/modules/:hostModuleID/instances\n{asset_id: sourceID}
REST -> Service : AddAssetToWorkspace(hostModuleID, sourceID)
Service -> Fabric : GetModelRecord(hostModuleID)
Fabric --> Service : *Model3D (draft)

alt Module source déjà instancié dans le module hôte (RM15)
    Service -> Service : Détecter instance existante\n(WorkspaceInstances contient sourceID)
    REST --> Client : 200 + flag "already_present"\navec liste instances existantes
    Client -> REST : POST /api/modules/:hostModuleID/instances\n{asset_id: sourceID, force: true}
    REST -> Service : AddAssetToWorkspace(hostModuleID, sourceID)
end

Service -> Service : Créer WorkspaceInstance{ID: newUUID, AssetID: sourceID}
Service -> Fabric : StoreModelRecord(hostModule)
Fabric --> Service : OK
Service -> Service : EnsureVirtualSlot(sourceID) — RM13
Service --> REST : *Model3D (draft mis à jour)
REST --> Client : 200 moduleDTO
@enduml
```

## Règles métier déclenchées

| Règle | Description | Point d'application |
|-------|-------------|---------------------|
| **RM13** | Slot virtuel garanti pour chaque asset instancié | `EnsureVirtualSlot()` dans `AddAssetToWorkspace()` |
| **RM15** | Seconde instance indépendante si le Module est déjà instancié dans le module hôte | Détection dans `AddAssetToWorkspace()` |

## Exigences non-fonctionnelles

- **ENF12** : Rôle vérifié côté serveur — seuls Concepteur et Consommateur peuvent ajouter un Module comme instance
- **ENF28** : Le Module source reste immuable — seule une référence (`AssetID`) est copiée, pas les données

## Notes d'implémentation

**Endpoints REST utilisés :**
- `GET /api/modules/:id` → récupère le Module source depuis la blockchain
- `POST /api/modules/:id/instances` → ajoute une instance au module hôte (handlers.go:~1127)

**Commande CLI équivalente :** `myr module list` (méthode `ListModules`) pour rechercher le Module source, puis `myr model instance add <hostModuleID> <sourceModuleID>` (méthode `AddAssetToWorkspace`) pour l'instancier dans le module hôte — chaque appel crée une nouvelle `WorkspaceInstance` indépendante, y compris pour une seconde instance du même Module (RM15). Une fois une Liaison créée, `myr module add-assembly <hostModuleID> <connID>` (méthode `AddAssemblyToModule`) la rattache au Module hôte. Voir `specs/3-Conception/DC_CLI_Model.md` § 5.

**RM15 — Détection de doublon :** `AddAssetToWorkspace()` inspecte `m.WorkspaceInstances` avant de créer l'instance et retourne un indicateur `already_present` dans la réponse si le même asset y figure déjà, sans bloquer l'ajout — à charge du client (script, CLI ou GUI) de décider s'il force une seconde instance.

**Note architecture :** L'instance créée est une référence (`AssetID`), pas une copie de l'entité — les modifications du Module source sur la blockchain n'affectent pas les snapshots `ModuleVersion` déjà soumis, mais sont reflétées à la prochaine lecture du module hôte.

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCMOD — Module](../../Carte_des_specs.md#UCMOD%20—%20Module)
- [UCMOD02 — couche expression](../../1-Expression/UCMOD-Module/UCMOD02.md)
- [Traçabilité UCMOD02 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCMOD02)

**Exigences fonctionnelles couvertes**
- [EF28 — Ajouter un module existant à l'espace de travail](../../1-Expression/Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Use cases cités**
- [UCAM01 — Liaison entre interfaces](../UCAM-Assemblage_Module/UCAM01.md)
- [UCCL01 — Faire une recherche par filtre](../UCCL-Composant_Lecture/UCCL01.md)
- [UCMOD01 — Créer un Module](UCMOD01.md)
- [UCREC01 — Rechercher une référence existante](../UCREC-Recherche/UCREC01.md)

**Règles métier**
- [RM13 — Slot virtuel garanti](../../1-Expression/Regles_Metier.md#3.%20Interfaces%20et%20liaisons)
- [RM15 — Instance indépendante](../../1-Expression/Regles_Metier.md#4.%20Composition%20d'un%20Module%20%28instances%29)

**Exigences non fonctionnelles**
- [ENF12 — Contrôle d'accès par rôle](../../1-Expression/Exigences_Non_Fonctionnelles.md#Tableau%20des%20exigences%20non-fonctionnelles)
- [ENF28 — Immuabilité des transactions blockchain](../../1-Expression/Exigences_Non_Fonctionnelles.md#Tableau%20des%20exigences%20non-fonctionnelles)

**Documents cités**
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)

**Cité par**
- [Expression_des_besoins_Intro](../../1-Expression/Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../../1-Expression/Matrice_Tracabilite.md)
- [UCAM08 (expression)](../../1-Expression/UCAM-Assemblage_Module/UCAM08.md)
- [todo (expression)](../../1-Expression/todo.md)
- [Analyse_des_besoins](../Analyse_des_besoins.md)
- [UCDEV02 (analyse)](../UCDEV-Developpement/UCDEV02.md)
- [UCMOD03 (analyse)](UCMOD03.md)
- [todo (analyse)](../todo.md)
- [Architecture_Composition](../../3-Conception/Architecture_Composition.md)
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)
- [Sequence_soumission_module](../../3-Conception/Sequence_soumission_module.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
