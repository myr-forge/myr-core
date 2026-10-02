---
categorie: Module
titre: "Créer un Module"
probabilite: 3
impact: 5
importance: 15
etat: analyse
tags:
  - couche/analyse
  - type/use-case
  - famille/UCMOD
  - domaine/model
  - uc/UCMOD01
  - rm/RM03
  - rm/RM04
  - rm/RM10
  - rm/RM11
  - rm/RM13
  - rm/RM16
  - rm/RM19
  - enf/ENF12
  - enf/ENF18
  - enf/ENF31
---

# Créer un Module

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Concepteur" as C

rectangle "Application MYR" {
    usecase "Créer un module (draft)" as UC1
    usecase "Créer liaisons entre composants" as UC2
    usecase "Nommer et configurer le module" as UC3
}

C --> UC1
UC1 ..> UC2 : <<include>>
UC1 ..> UC3 : <<include>>

@enduml
```

## Contexte

Le Concepteur compose un Module en assemblant plusieurs Composants ou Modules existants via leurs Interfaces physiques. Cette composition est un état local côté serveur — hors blockchain, modifiable librement par actions directes sur le module (`AddAssetToWorkspace`, `AddAssemblyLink`).

Un Module nouvellement créé est toujours initialisé en état **draft** (RM16) : il existe localement mais n'est pas encore ancré sur la blockchain. Tant qu'il reste en état `draft`, le Concepteur peut ajouter, retirer ou reconfigurer des Liaisons (`Connection`) autant de fois que nécessaire. Un nom par défaut lui est attribué à la création (`<identité du Concepteur>_<date>_<heure>`), modifiable ensuite (voir UCMOD03).

La soumission à la blockchain est une étape distincte et explicite (UCMOD06). Elle crée une `ModuleVersion` immuable horodatée — snapshot figé et ancré, non modifiable après publication.

Un Module peut aussi être créé à partir d'un Module déjà soumis, pour reprendre son travail sans reconstruire manuellement l'assemblage (voir flux alternatif « Dérivation d'un Module existant » ci-dessous).

## Pré-conditions

- Identité authentifiée avec rôle **Concepteur** (`contributor`)
- Au moins un Composant ou Module existant disponible dans le système (en local ou sur la blockchain)
- Le module n'est pas encore créé (première création)

## Scénario

**Étape initiale :** `POST /api/modules` est appelée (ou l'équivalent CLI `myr module create`), pour le compte du Concepteur

### Flux nominal — Module créé en état draft

1. Le système appelle `POST /api/modules` → `service.CreateModule()` → initialise `Model3D` avec `Status=draft`, `Assemblies=[]`, `WorkspaceInstances=[]`
2. Des Composants sont ajoutés comme instances (`POST /api/modules/:id/instances`) — chaque ajout crée une `WorkspaceInstance` indépendante et garantit au moins un slot virtuel (RM13)
3. Des Liaisons sont créées entre Interfaces compatibles — le service vérifie la compatibilité (RM10/RM11 : catégorie + type + sens + plages de valeurs) et refuse toute liaison incompatible
4. Le module est nommé et configuré (nom, description, licence) via `PUT /api/modules/:id`
5. Le module est enregistré localement en état **draft**

### Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)

1. `POST /api/modules` est appelée avec `parent_id` renseigné, référençant un Module déjà soumis
2. Si `license_id` est fourni : le service vérifie la compatibilité avec la licence du Module parent (RM03), comme pour un composant dérivé (UCCE02, UCCE04)
3. Le nouveau Module est créé en état **draft** avec `ParentID` renseigné
4. La composition complète du Module parent est dupliquée dans le nouveau brouillon : chaque `WorkspaceInstance` est clonée avec un nouvel identifiant d'instance, et chaque Liaison interne (`Assemblies`) est reconstruite entre les instances clonées correspondantes
5. Le Concepteur retrouve immédiatement l'assemblage complet du Module parent, prêt à être modifié, sans avoir ajouté une seule instance ni recréé une seule Liaison manuellement

### Flux erreur — Dérivation avec licence incompatible

1. `parent_id` et `license_id` sont tous deux renseignés
2. La licence proposée est incompatible avec celle du Module parent (RM03)
3. Le système refuse la création : "Incompatibilité de licence : <raison>"
4. Aucun Module n'est créé, aucune composition n'est dupliquée

### Flux alternatif — Liaison via interface virtuelle

1. Une connexion est demandée depuis un slot virtuel vers une Interface physique (`POST /api/virtual-connect`)
2. Le système appelle `ConnectVirtualToPhysical()` — matérialise l'interface virtuelle avec les attributs de l'interface physique cible (direction opposée)
3. Un nouveau slot virtuel est automatiquement recréé pour l'asset (invariant RM13)
4. La Liaison est créée et le module reste en état `draft`

### Flux erreur — Aucune liaison créée (sauvegarde bloquée)

1. Une tentative de finalisation du module est effectuée sans qu'aucune Liaison n'ait été créée
2. Le système détecte `len(m.Assemblies) == 0` — refus côté service
3. Message retourné : "Ajoutez au moins une liaison entre composants"

### Flux erreur — Erreur de persistance locale

1. Le service ne peut pas persister le `Model3D` (store indisponible)
2. Message d'erreur : "Impossible de créer le module — réessayez"
3. Aucune entrée créée — état du système inchangé

## Post-conditions

- Le module existe en état **draft** dans le store local (`adapters/out/localstorage/`)
- Le module possède un ID unique (UUID généré côté serveur — RM04)
- Les `WorkspaceInstances` ajoutées sont persistées
- Le module n'est pas visible sur le réseau (soumission requise — UCMOD06)
- Chaque asset instancié dispose d'au moins un slot virtuel (RM13)
- En cas de dérivation (`ParentID` renseigné) : le nouveau module référence son Module parent, et sa composition initiale (instances + liaisons internes) est une copie indépendante de celle du parent — modifier le nouveau brouillon n'affecte jamais le Module parent déjà soumis

## Diagramme de séquence

```plantuml
@startuml
participant "Client\n(CLI ou API REST)" as Client
participant "REST Handler\n(adapters/in/rest/)" as REST
participant "Model Service\n(domain/model/)" as Service
database "LocalStorage\n(adapters/out/localstorage/)" as Local
database "Fabric\n(adapters/out/fabric/)" as Fabric

Client -> REST : POST /api/modules\n{name, description, channelID, licenseID}
REST -> Service : CreateModule(ModuleRequest)
Service -> Service : Générer UUID (RM04)\nInitialiser Status=draft
Service -> Fabric : StoreModelRecord(m)
Fabric --> Service : OK
Service --> REST : *Model3D (draft)
REST --> Client : 201 moduleDTO

Client -> REST : POST /api/modules/:id/instances\n{asset_id}
REST -> Service : AddAssetToWorkspace(moduleID, assetID)
Service -> Fabric : GetModelRecord(moduleID)
Fabric --> Service : *Model3D
Service -> Service : Créer WorkspaceInstance
Service -> Fabric : StoreModelRecord(m)
Service -> Service : EnsureVirtualSlot(assetID) — RM13
REST --> Client : 200 moduleDTO (instances mises à jour)

Client -> REST : POST /api/connections\n{fromIfaceID, toIfaceID, ...}
REST -> Service : AddAssemblyLink(...)
Service -> Service : Vérifier compatibilité\nifacesCompatible() — RM10/RM11

alt Interfaces compatibles
    Service -> Local : SaveConnection(conn)
    Local --> Service : OK
    REST --> Client : 200 Connection
else Interfaces incompatibles
    Service --> REST : erreur "interfaces incompatibles"
    REST --> Client : 400 Bad Request
end

Client -> REST : POST /api/modules/:id/assemblies\n{connection_id}
REST -> Service : AddAssemblyToModule(moduleID, connID)
Service -> Fabric : GetModelRecord(moduleID)
Service -> Service : Ajouter connID à Assemblies
Service -> Fabric : StoreModelRecord(m)
REST --> Client : 200 moduleDTO (draft)
@enduml
```

### Diagramme de séquence — Dérivation d'un Module existant

```plantuml
@startuml
participant "Client\n(CLI ou API REST)" as Client
participant "REST Handler\n(adapters/in/rest/)" as REST
participant "Model Service\n(domain/model/)" as Service
database "LocalStorage\n(adapters/out/localstorage/)" as Local
database "Fabric\n(adapters/out/fabric/)" as Fabric

Client -> REST : POST /api/modules\n{name, parent_id, license_id?, ...}
REST -> Service : CreateModule(ModuleRequest{ParentID, ...})
Service -> Fabric : GetModelRecord(parentID)
Fabric --> Service : *Model3D (module parent, submitted)

alt license_id fourni et incompatible avec le parent (RM03)
    Service --> REST : erreur "incompatibilité de licence"
    REST --> Client : 422 Unprocessable Entity
else Compatible ou licence non fournie
    Service -> Service : Générer UUID (RM04)\nInitialiser Status=draft, ParentID
    loop pour chaque WorkspaceInstance du parent
        Service -> Service : Cloner l'instance (nouvel InstanceID)
    end
    loop pour chaque Liaison interne du parent (Assemblies)
        Service -> Local : ListConnections() / SaveConnection(clone)
        Local --> Service : nouvelle Connection (instances remappées)
    end
    Service -> Local : SaveDraft(m) — instances et Assemblies clonés
    Local --> Service : OK
    Service --> REST : *Model3D (draft, composition dupliquée)
    REST --> Client : 201 moduleDTO
end
@enduml
```

## Règles métier déclenchées

| Règle | Description | Point d'application |
|-------|-------------|---------------------|
| **RM03** | Compatibilité de licence si `ParentID` et `LicenseID` sont renseignés à la création | `CreateModule()`, même vérification que `AddFull()` |
| **RM04** | UUID généré par le système, jamais par le client | `generateID()` dans `CreateModule()` |
| **RM10** | Vérification de compatibilité à chaque création de Liaison | `ifacesCompatible()` dans `AddAssemblyLink()` |
| **RM11** | 5 critères : catégorie + tag (manquant E2) + type + sens + plages | `ifacesCompatible()` — tag absent du code |
| **RM13** | Au moins un slot virtuel garanti par asset instancié | `EnsureVirtualSlot()` dans `AddAssetToWorkspace()` |
| **RM16** | Module initialisé en état `draft` obligatoirement | `Status: ModuleDraft` dans `CreateModule()` |

## Exigences non-fonctionnelles

- **ENF12** : Rôle Concepteur vérifié côté serveur avant toute opération d'écriture
- **ENF18** : Le domaine ne connaît que des interfaces — aucune dépendance Fabric dans `domain/model/`
- **ENF31** : Validation des données avant toute soumission

## Notes d'implémentation

**Endpoints REST utilisés :**
- `POST /api/modules` → crée le module, avec `parent_id` optionnel pour une dérivation (handlers.go:~1296)
- `POST /api/modules/:id/instances` → ajoute un asset comme instance (handlers.go:~1127)
- `POST /api/modules/:id/assemblies` → associe une Liaison au module (handlers.go:~1078)
- `PATCH /api/modules/:id` → met à jour nom/description/licence (voir UCMOD03)

**Commande CLI équivalente :** `myr module create --name <nom> --channel <id> [--owner-id <id>] [--description <texte>] [--license <id>] [--parent-id <id>]` appelle le même `CreateModule(ModuleRequest)` que `POST /api/modules`. L'ajout de composants comme instances et la création de liaisons se poursuivent avec `myr model instance add <moduleID> <assetID>` et `myr model link add` (mêmes méthodes `AddAssetToWorkspace` / `AddAssemblyLink`, mêmes vérifications RM10/RM11/RM13 côté service, quel que soit le canal). Voir `specs/3-Conception/DC_CLI_Model.md` § 3.6 et § 5. Ces commandes sont exécutées par l'administrateur du serveur via SSH, pour le compte du Concepteur (principe d'exécution distante).

**Dérivation (`parent_id`) — pourquoi une copie de composition et pas seulement une référence :** un simple lien de filiation (`ParentID` renseigné sans copie) documenterait la provenance mais laisserait le Concepteur reconstruire manuellement chaque instance et chaque Liaison — ce qui ne répond pas au besoin de reprendre un travail existant. `CreateModule()` clone donc `WorkspaceInstances` (nouveaux identifiants d'instance, pour ne jamais entrer en collision avec ceux du parent) et les Liaisons internes du parent (nouvelles `Connection`, mêmes attributs, instances remappées) dans le nouveau brouillon. Le Module parent, déjà soumis, n'est jamais modifié par cette opération — seule sa composition est lue.

**Écart E2 à noter :** Le champ `Tag` est absent de `AssetInterface` dans `domain/model/entity.go`. La vérification RM11 ne comporte donc que 4 critères dans le code actuel (catégorie + type + sens + valeurs). Le champ `Tag` doit être ajouté pour la conformité complète à RM11.

**Écart E5 (RM19) :** `AddAssemblyToModule()` ne vérifie pas `Status != ModuleSubmitted` — un module soumis reste modifiable dans le code actuel. Le comportement attendu (fork obligatoire pour tout module soumis) est décrit dans UCMOD06 et doit être corrigé.

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Étape suivante — conception**
- **API_REST** : [§ 4. D3/D4/D6 — Composants (ressource unique, ADR-11)](../../3-Conception/API_REST.md#4.%20D3/D4/D6%20—%20Composants%20%28ressource%20unique,%20ADR-11%29)
- **Architecture_Composition** : [§ Architecture — Composition (D3/D5/D6 : composant, assemblage, module)](../../3-Conception/Architecture_Composition.md#Architecture%20—%20Composition%20%28D3/D5/D6%20:%20composant,%20assemblage,%20module%29)
- **DC_CLI_Model** : [§ DC — CLI Modèle : Référence des commandes composant / interfaces / module](../../3-Conception/DC_CLI_Model.md#DC%20—%20CLI%20Modèle%20:%20Référence%20des%20commandes%20composant%20/%20interfaces%20/%20module) · [§ 2. Arbre de commandes](../../3-Conception/DC_CLI_Model.md#2.%20Arbre%20de%20commandes) · [§ 3.5 `myr model instance add` / `remove`](../../3-Conception/DC_CLI_Model.md#3.5%20`myr%20model%20instance%20add`%20/%20`remove`) · [§ 5. Table de correspondance méthode domaine → commande CLI → use case](../../3-Conception/DC_CLI_Model.md#5.%20Table%20de%20correspondance%20méthode%20domaine%20→%20commande%20CLI%20→%20use%20case)
- **Sequence_soumission_module** : [§ Séquence — Composition et soumission d'un module (D5/D6)](../../3-Conception/Sequence_soumission_module.md#Séquence%20—%20Composition%20et%20soumission%20d'un%20module%20%28D5/D6%29)

<!-- liens-obsidian:end -->
