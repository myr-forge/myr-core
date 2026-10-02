---
categorie: Module
titre: "Visualiser les composants d'un Module"
probabilite: 2
impact: 5
importance: 10
etat: analyse
tags:
  - couche/analyse
  - type/use-case
  - famille/UCMOD
  - domaine/model
  - uc/UCMOD04
  - rm/RM13
  - enf/ENF12
---

# Visualiser les composants d'un Module

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Client\n(dépôt GUI externe)" as U

rectangle "API myr" {
    usecase "Obtenir la composition\nd'un Module" as UC1
    usecase "Obtenir les interfaces\nexposées" as UC3
}

U --> UC1
UC1 ..> UC3 : <<include>>

@enduml
```

## Contexte

Un Module est structuré comme un graphe d'instances de Composants et de sous-Modules reliés par des Liaisons. Tout client authentifié doit pouvoir obtenir la composition interne d'un Module soumis pour en comprendre la structure, évaluer la réutilisabilité, préparer une commande ou planifier une intégration.

Cette lecture est une opération en **lecture seule** qui ne modifie pas l'état du Module. Elle s'appuie sur trois sources de données :
- `Model3D.WorkspaceInstances` : liste des instances d'assets composant le module
- `Model3D.Assemblies` : IDs des Liaisons internes
- `GetModuleInterfaces()` : calcule les interfaces exposées (non connectées en interne) de façon récursive — chaque instance du module contribue ses propres interfaces non connectées en interne à la liste exposée, même si plusieurs instances partagent le même asset sous-jacent : le décompte se fait par `(instance, interface)`, jamais dédupliqué au niveau de l'asset

> La présentation (navigation hiérarchique, fil d'Ariane, rendu 3D) relève du dépôt GUI externe — hors périmètre de ce document. Seuls les contrats REST et CLI ci-dessous font partie de `myr` — le CLI n'est pas un citoyen de seconde zone par rapport à l'API.

## Pré-conditions

- Utilisateur authentifié (tout rôle)
- Un Module identifié (ID connu du client)
- Le Module est accessible (état `submitted` sur la blockchain, ou `draft` si le demandeur est le propriétaire)

## Scénario

### Flux nominal — Composition du module retournée

1. Le client appelle `GET /api/modules/:id`
2. La réponse contient :
   - La liste des `WorkspaceInstances` (composants et sous-modules composant le module)
   - Les Liaisons internes (connexions entre instances)
3. Pour chaque instance, le nom et la catégorie de l'asset référencé sont inclus

### Flux alternatif — Module contenant des sous-modules (hiérarchie)

1. Le Module contient des instances référençant d'autres Modules (sous-modules)
2. Le client appelle récursivement `GET /api/modules/:subModuleID` pour obtenir la composition de chaque sous-module

### Flux alternatif — Consultation des interfaces exposées

1. Le client appelle `GET /api/modules/:id/interfaces`
2. Service : `GetModuleInterfaces(id)` calcule récursivement les interfaces non connectées en interne
3. La liste des interfaces exposées est retournée avec leurs attributs (catégorie, type, direction, valeurs, unité)

### Flux erreur — Module introuvable

1. L'ID du Module ne correspond à aucun record sur la blockchain
2. Le système retourne une erreur 404 "Module introuvable"

### Flux erreur — Module en état draft appartenant à un autre utilisateur

1. Le Module cible est en état `draft` et l'utilisateur n'est pas le propriétaire
2. Le serveur retourne `403 Forbidden`

## Post-conditions

- La composition complète du Module est retournée (instances + liaisons + interfaces exposées)
- Aucune modification de l'état du Module

## Diagramme de séquence

```plantuml
@startuml
participant "Client\n(dépôt GUI externe)" as Client
participant "REST Handler\n(adapters/in/rest/)" as REST
participant "Model Service\n(domain/model/)" as Service
database "LocalStorage\n(adapters/out/localstorage/)" as Local
database "Fabric\n(adapters/out/fabric/)" as Fabric

Client -> REST : GET /api/modules/:id
REST -> Service : GetModule(id)
Service -> Fabric : GetModelRecord(id, "")

alt Module trouvé
    Fabric --> Service : *Model3D
    Service --> REST : *Model3D
    REST --> Client : 200 moduleDTO\n(WorkspaceInstances + Assemblies)
else Module introuvable
    Fabric --> Service : ErrNotFound
    REST --> Client : 404 "Module introuvable"
end

Client -> REST : GET /api/modules/:id/interfaces
REST -> Service : GetModuleInterfaces(id)
Service -> Fabric : GetModelRecord(id, "")
Service -> Service : getModuleInterfacesInto()\n[récursif sur WorkspaceInstances]
Service -> Local : ListInterfacesForAsset(subAssetID) ×n
Local --> Service : []*AssetInterface
Service -> Service : Filtrer interfaces connectées en interne\n(via ListConnections + Assemblies)
Service --> REST : []*AssetInterface exposées
REST --> Client : 200 [{id, category, type, direction, valueMin, valueMax, unit}, ...]

Client -> REST : GET /api/modules/:id/connections
REST -> Service : GetModule(id)
Service -> Fabric : GetModelRecord(id)
Service -> Local : ListConnections()
Service -> Service : Filtrer connexions dans Assemblies
REST --> Client : 200 [Connection ...]
@enduml
```

## Règles métier déclenchées

| Règle | Description | Point d'application |
|-------|-------------|---------------------|
| **RM13** | Slot virtuel garanti — affiché comme point de connexion libre | `GetModuleInterfaces()` inclut les slots virtuels |

## Exigences non-fonctionnelles

- **ENF12** : Accès en lecture restreint aux utilisateurs authentifiés (sauf si la configuration réseau autorise les Visiteurs)

## Notes d'implémentation

**Endpoints REST utilisés :**
- `GET /api/modules/:id` → chargement principal (handlers.go:~1202)
- `GET /api/modules/:id/interfaces` → interfaces exposées via `GetModuleInterfaces()` (handlers.go:~1032)
- `GET /api/modules/:id/connections` → liaisons internes (handlers.go:~1047)

**Récursivité `GetModuleInterfaces()` :** La fonction `getModuleInterfacesInto()` (service.go:~640) parcourt récursivement les `WorkspaceInstances` pour calculer les interfaces non connectées en interne. Un cache par `assetID` évite les appels blockchain redondants — mais ce cache sert uniquement à réutiliser le résultat déjà calculé pour un asset donné (résolution récursive de ses propres sous-instances) ; il ne doit jamais réduire le nombre d'instances traitées au niveau du module courant : quand plusieurs `WorkspaceInstances` du module référencent le même asset, chacune contribue séparément ses interfaces non connectées en interne à la liste exposée. Les modules profondément imbriqués peuvent générer de nombreux appels — un mécanisme de profondeur maximale est à envisager pour les cas extrêmes.

**Commande CLI équivalente :** `myr module get <id>` (méthode `GetModule`) est le strict équivalent en lecture seule de `GET /api/modules/:id`. `myr module interfaces <id>` (méthode `GetModuleInterfaces`) couvre `GET /api/modules/:id/interfaces`. Voir `specs/3-Conception/DC_CLI_Model.md` § 5.

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCMOD — Module](../../Carte_des_specs.md#UCMOD%20—%20Module)
- [UCMOD04 — couche expression](../../1-Expression/UCMOD-Module/UCMOD04.md)
- [Traçabilité UCMOD04 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCMOD04)

**Exigences fonctionnelles couvertes**
- [EF29 — Visualiser la composition d'un module](../../1-Expression/Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Règles métier**
- [RM13 — Slot virtuel garanti](../../1-Expression/Regles_Metier.md#3.%20Interfaces%20et%20liaisons)

**Exigences non fonctionnelles**
- [ENF12 — Contrôle d'accès par rôle](../../1-Expression/Exigences_Non_Fonctionnelles.md#Tableau%20des%20exigences%20non-fonctionnelles)

**Documents cités**
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)

**Cité par**
- [Expression_des_besoins_Intro](../../1-Expression/Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../../1-Expression/Matrice_Tracabilite.md)
- [todo (expression)](../../1-Expression/todo.md)
- [Analyse_des_besoins](../Analyse_des_besoins.md)
- [UCDEV02 (analyse)](../UCDEV-Developpement/UCDEV02.md)
- [todo (analyse)](../todo.md)
- [API_REST](../../3-Conception/API_REST.md)
- [Architecture_Composition](../../3-Conception/Architecture_Composition.md)
- [Chaincode](../../3-Conception/Chaincode.md)
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)
- [DC_D1_Auth_Identity](../../3-Conception/DC_D1_Auth_Identity.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
