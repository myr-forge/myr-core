---
categorie: Composant Ecriture
titre: "Ajout d'un composant Physique"
probabilite: 3
impact: 5
importance: 15
etat: analyse
tags:
  - couche/analyse
  - type/use-case
  - famille/UCCE
  - domaine/model
  - uc/UCCE01
  - rm/RM01
  - rm/RM02
  - rm/RM03
  - rm/RM04
  - rm/RM05
  - rm/RM07
  - rm/RM16
  - rm/RM19
  - enf/ENF12
  - enf/ENF30
---

# Ajout d'un composant Physique

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Concepteur" as C

rectangle "Application MYR" {
    usecase "Ajouter un composant physique" as UC1
    usecase "Vérifier le hash (anti-plagiat SHA-256)" as UC2
    usecase "Analyser similarité SCM" as UC3
    usecase "Vérifier compatibilité de licence" as UC4
    usecase "Enregistrer sur la blockchain" as UC5
}

C --> UC1
UC1 ..> UC2 : <<include>>
UC1 ..> UC3 : <<include>>
UC1 ..> UC4 : <<extend>> (si ParentID)
UC1 ..> UC5 : <<include>>

@enduml
```

## Contexte

L'ajout d'un composant physique est le use case fondamental du domaine D3. Il permet à un **Concepteur** d'enregistrer une création originale (pièce CAO, fichier 3D/STL/STEP/OBJ…) sur la blockchain Fabric, établissant ainsi son droit d'auteur de façon immuable.

L'entité produite est un `Model3D` avec `Hash` renseigné et `WorkspaceInstances` vide — c'est la définition d'un **Composant** dans la terminologie MYR.

La catégorie `base` est la seule catégorie qui ne requiert pas de `ParentID`. C'est aussi la seule catégorie soumise à la vérification anti-plagiat complète (SHA-256 + SCM > 50%). Les 7 autres catégories impliquent toutes un `ParentID` et un fichier parent existant sur le réseau.

**Écart code connu (E4) :** `service.AddFull()` calcule le SHA-256 du fichier soumis mais ne compare pas ce hash avec les assets déjà inscrits sur Fabric. La vérification complète (RM01) est une cible à implémenter.

## Pré-conditions

- Le Concepteur est authentifié avec le rôle `contributor` (session REST valide).
- Un canal Fabric est opérationnel et accessible.
- Le Concepteur dispose d'un fichier 3D/CAO valide (STL, STEP, OBJ, ou format natif).
- Pour une catégorie non-`base` : l'asset parent existe sur le canal et son ID est connu.

## Scénario

**Étape initiale :** `POST /api/components` est appelée (ou l'équivalent CLI `myr model add`) avec le fichier 3D et les métadonnées

### Flux nominal — Composant base nouveau

1. Le fichier 3D est transmis (champ `file` ou `stl` — multipart/form-data).
2. Les métadonnées sont transmises : nom, description, licence, tags, catégorie (`base`).
3. Le REST Handler valide les champs : `name` non vide (max 256), `owner_id` valide, `channel_id` valide.
4. Le handler appelle `service.AddFull(AddRequest{...})`.
5. Le service calcule le SHA-256 du fichier : `sha256:<hex>`.
6. **[Cible RM01]** Le service interroge Fabric (`ListModelRecords`) et compare le hash avec tous les assets existants.
7. Si aucun doublon : le service analyse la similarité SCM (seuil 50%) avec les assets existants.
8. Le service construit le `Model3D` avec un UUID généré (`generateID()`) et le hash — le fichier transmis n'est jamais conservé par Myr, qui priorise la traçabilité (savoir où le composant existe, voir UCCL03) plutôt que l'hébergement.
9. Le service soumet la transaction `StoreModel` sur Fabric (`blockchain.StoreModelRecord(m)`).
10. Fabric valide la transaction et ancre le bloc.
11. L'API retourne `201 Created` avec le `Model3D` JSON (ID, name, hash, blockID…).

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 9 déduit(s)
> - 🟡 [TestAddFull_AllFields](../../../docs/tests/domain-model/TestAddFull_AllFields.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_AlwaysDraft_NoBlockchainWrite](../../../docs/tests/domain-model/TestAddFull_AlwaysDraft_NoBlockchainWrite.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_MissingFile_ReturnsError](../../../docs/tests/domain-model/TestAddFull_MissingFile_ReturnsError.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_NilBlockchain_StillWorks](../../../docs/tests/domain-model/TestAddFull_NilBlockchain_StillWorks.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_NilDraftStore_Rejected](../../../docs/tests/domain-model/TestAddFull_NilDraftStore_Rejected.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_NoFile_NoVersionNoHash](../../../docs/tests/domain-model/TestAddFull_NoFile_NoVersionNoHash.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_STLAndSTEP_BothAccepted](../../../docs/tests/domain-model/TestAddFull_STLAndSTEP_BothAccepted.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_WithFile_HashAndVersionCreated](../../../docs/tests/domain-model/TestAddFull_WithFile_HashAndVersionCreated.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_WithParent](../../../docs/tests/domain-model/TestAddFull_WithParent.md) — déduit : teste `ModelService.AddFull`
<!-- tests-obsidian:end -->

### Flux alternatif — Import depuis un format CAO non natif (STL, STEP, OBJ)

1. Le fichier est dans un format supporté mais non natif.
2. Le système accepte le fichier et calcule son SHA-256 normalement.
3. Les métadonnées géométriques extractibles automatiquement (dimensions, volume) sont pré-renseignées selon le format.
4. Le Concepteur complète les métadonnées non extractibles (description, licence).
5. La transaction est soumise normalement (flux nominal à partir de l'étape 7).

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)

1. Le Concepteur transmet `draft: true` (ou `myr model add --draft`) en plus des champs habituels.
2. Le service crée le `Model3D` avec `Status: draft` — la transaction Fabric initiale est tout de même soumise (comme pour un module en `draft`, RM16) mais l'asset reste modifiable par convention métier tant qu'il n'est pas soumis.
3. Le Concepteur peut alors ajouter ou modifier des interfaces (UCCE06, UCAM03) — ces changements restent en brouillon local (`InterfaceStore`, ADR-02) et n'émettent aucune transaction Fabric supplémentaire.
4. Quand le Concepteur est prêt, il appelle l'action de soumission (`POST /api/components/:id/submit` ou `myr model submit <id>` — voir `specs/3-Conception/DC_CLI_Model.md`) : le service relit le brouillon, embarque `Interfaces` dans le `Model3D`, soumet une dernière transaction Fabric et passe `Status` à `submitted`.
5. Une fois `submitted`, l'asset est immuable (règle 7, RM19) — toute nouvelle interface exige un fork (voir UCCE06).

> Ce flux est **optionnel** : le flux nominal (sans `draft: true`) reste inchangé — un composant est créé et soumis en une seule transaction, comme aujourd'hui. `draft` sert uniquement au cas où le Concepteur veut affiner les interfaces après création, sans les connaître entièrement à l'import du fichier.

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

### Flux alternatif — Catégorie dérivée (non-`base`)

1. Le Concepteur sélectionne une catégorie dérivée (`amelioration`, `variation`, `adaptation`, `derivation`, `extension`, `regression`, `decoupage`).
2. Il renseigne le `parent_id` (UUID de l'asset parent sur le canal).
3. Si une `license_id` est fournie : le service vérifie la compatibilité de licence avec le parent (`CheckLicenseCompatibility`). En cas d'incompatibilité : erreur retournée avant soumission Fabric.
4. La vérification anti-plagiat SHA-256 + SCM n'est **pas** déclenchée pour les catégories dérivées.
5. La transaction est soumise avec `ParentID` renseigné.

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux erreur — Hash déjà existant (RM01)

1. La comparaison des hashes (étape 7) détecte un doublon exact.
2. Le service retourne une erreur avant toute soumission Fabric.
3. L'API retourne `409 Conflict` : `{ "error": "Composant déjà existant — risque de plagiat. Contacter l'administration." }`.

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux erreur — Similarité SCM > 50% (RM01)

1. L'analyse SCM (étape 8) détecte une similarité structurelle supérieure à 50% avec un asset existant.
2. Le service retourne une erreur avant toute soumission Fabric.
3. L'API retourne `409 Conflict` : `{ "error": "Similarité trop élevée avec un composant existant. Contacter l'administration." }`.

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux erreur — Incompatibilité de licence (RM03)

1. `CheckLicenseCompatibility(parentLicenseID, req.LicenseID)` retourne `Compatible: false`.
2. Le service retourne l'erreur avant toute soumission Fabric.
3. L'API retourne `422 Unprocessable Entity` : `{ "error": "Incompatibilité de licence : <raison>." }`.

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux erreur — Échec endorsement Fabric

1. `blockchain.StoreModelRecord(m)` retourne une erreur Fabric (nœud indisponible, politique non satisfaite).
2. Le fichier uploadé sur IPFS reste (orphelin temporaire — acceptable).
3. L'API retourne `500 Internal Server Error` : `{ "error": "Erreur blockchain : <message>. Aucune donnée enregistrée." }`.

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

## Post-conditions

- Le `Model3D` est inscrit sur la blockchain Fabric (immuable — RM07).
- Un UUID unique est attribué au composant (généré par le service, non par le client — RM04).
- Le fichier 3D est stocké dans IPFS avec une référence dans `Versions[0].Hash`.
- Le droit d'auteur est enregistré via `OwnerID`.

## Diagramme de séquence

```plantuml
@startuml
participant "Client\n(CLI ou API REST)" as Browser
participant "REST Handler\n(adapters/in/rest/)" as REST
participant "Model Service\n(domain/model/)" as ModelSvc
database "IPFS\n(adapters/out/ipfs/)" as IPFS
database "Fabric\n(adapters/out/fabric/)" as Fabric

Browser -> REST : POST /api/components\n(multipart: name, file, category, license_id, parent_id, tags)
REST -> REST : ParseMultipartForm(32 MB)
REST -> REST : validateFields(name, owner_id, channel_id)

alt Champ invalide (name vide, owner_id malformé…)
    REST --> Browser : 400 Bad Request { error }
else Champs valides
    REST -> ModelSvc : AddFull(AddRequest{FilePath, Name, Category, ParentID, LicenseID, ...})

    alt ParentID && LicenseID renseignés
        ModelSvc -> Fabric : GetModelRecord(ParentID)
        Fabric --> ModelSvc : parentModel
        ModelSvc -> ModelSvc : CheckLicenseCompatibility(parent.LicenseID, req.LicenseID)
        alt Incompatible
            ModelSvc --> REST : ErrLicenseIncompatibility
            REST --> Browser : 422 Incompatibilité de licence
        end
    end

    ModelSvc -> ModelSvc : hashFile(FilePath) → sha256:<hex>

    note over ModelSvc : [Cible RM01 — non implémenté]\nComparer hash avec Fabric ListModelRecords
    ModelSvc -> Fabric : ListModelRecords(channelID)
    Fabric --> ModelSvc : existingAssets[]

    alt Hash doublon détecté
        ModelSvc --> REST : ErrHashDuplicate
        REST --> Browser : 409 Composant déjà existant
    else
        ModelSvc -> ModelSvc : analyseSCM(file, existingAssets)
        alt Similarité > 50%
            ModelSvc --> REST : ErrSimilarityTooHigh
            REST --> Browser : 409 Similarité trop élevée
        else
            ModelSvc -> IPFS : Upload(FilePath)
            IPFS --> ModelSvc : storageRef
            ModelSvc -> ModelSvc : buildModel3D(UUID, hash, storageRef)
            ModelSvc -> Fabric : StoreModelRecord(model3D)
            alt Échec Fabric
                Fabric --> ModelSvc : ErrEndorsement
                ModelSvc --> REST : ErrBlockchain
                REST --> Browser : 500 Erreur blockchain
            else Succès
                Fabric --> ModelSvc : blockID
                ModelSvc --> REST : Model3D{ID, Name, Hash, BlockID}
                REST --> Browser : 201 Created { model3D }
            end
        end
    end
end

@enduml
```

## Règles métier déclenchées

| Règle | Description |
|-------|-------------|
| **RM01** | Anti-plagiat obligatoire pour tout asset `base` : SHA-256 + SCM > 50% → rejet |
| **RM02** | Catégorie obligatoire parmi les 8 types : `base`, `amelioration`, `variation`, `adaptation`, `derivation`, `extension`, `regression`, `decoupage` |
| **RM03** | Si `ParentID != ""` et `LicenseID != ""` : vérification de compatibilité de licence obligatoire |
| **RM04** | UUID généré par le service, jamais par le client |
| **RM05** | `ParentID` obligatoire pour tout asset non-`base` |
| **RM07** | Validation complète côté serveur avant toute soumission blockchain |
| **RM16** | `draft: true` optionnel → composant créé en `draft` (soumission différée) ; par défaut, un composant est directement `submitted` |
| **RM19** | Une fois `submitted`, l'asset est immuable — toute modification ultérieure (ex. ajouter une interface, UCCE06) exige un fork |

## Exigences non-fonctionnelles

| ID | Exigence |
|----|---------|
| **EF10** | Import de fichiers 3D (STL, STEP, OBJ minimum) |
| **EF11** | Vérification anti-plagiat obligatoire avant enregistrement |
| **ENF12** | Contrôle du rôle `contributor` côté serveur avant toute écriture |
| **ENF30** | En cas d'échec blockchain, l'état local (draft) est conservé intact |

## Notes d'implémentation

**Route existante :** `POST /api/components` → `handler.createAsset()` → `service.AddFull()` → `fabric.StoreModelRecord()`.

**Commande CLI équivalente (existante, à étendre) :** `myr model add <file> --name <nom> --channel <id> --category base [--description <texte>] [--tags <a,b>] [--owner-id <id>] [--draft]` (voir `specs/3-Conception/DC_CLI_Model.md` § 3.1). Aujourd'hui `adapters/in/cli/model.go` ne câble que `--name`/`--channel`/`--tags` via `modelSvc.Add()` — un sous-ensemble de `AddRequest`. Cible : basculer sur `modelSvc.AddFull(AddRequest{...})`, la même méthode que le handler REST `createAsset()`, pour exposer aussi `--category`, `--parent`, `--license` et `--draft` (RM16). Une fois câblée, la commande déclenche les mêmes règles (anti-plagiat RM01, compatibilité de licence RM03) que le flux REST — seul le canal de sortie change (texte terminal vs JSON HTTP). La soumission différée d'un composant en brouillon s'effectue via `myr model submit <id>` (§ 3.6bis).

**Écart E4 (RM01 incomplet) :** `service.AddFull()` calcule le SHA-256 mais ne compare pas avec les assets existants. À implémenter : appel `blockchain.ListModelRecords(channelID)` suivi d'une comparaison de hashes avant l'enregistrement blockchain (`StoreModelRecord`).

**Écart E1 (catégorie `decoupage` absente) :** `domain/model/entity.go` ne définit pas `CategoryDecoupage`. À ajouter : `CategoryDecoupage Category = "decoupage"`. Cette catégorie est la seule qui transforme un composant en module (UCAM05).

**Écart E8 (`Status` composant, cf. `specs/2-Analyse/Analyse_des_besoins.md` § Écarts structurels connus) :** `domain/model/entity.go` ne traite `Status` que pour les modules aujourd'hui. À corriger pour permettre `AddFull(AddRequest{Draft: true, ...})` de créer un composant avec `Status: draft`, et pour exposer une méthode de soumission (réutilisation de la logique de `SubmitModule`, généralisée à tout `Model3D`) qui embarque `Interfaces` et passe `Status` à `submitted`.

**Analyse SCM :** L'algorithme de similarité structurelle (SCM > 50%) est un service domaine indépendant à créer dans `domain/model/` — il n'est pas encore implémenté.

**Multipart fields acceptés :** `name`, `description`, `owner_id`, `channel_id`, `category`, `parent_id`, `license_id`, `tags` (CSV), `links` (JSON array), `file` ou `stl` (fichier binaire), `thumbnail` (data-URL).

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Étape suivante — conception**
- **API_REST** : [§ 4. D3/D4/D6 — Composants (ressource unique, ADR-11)](../../3-Conception/API_REST.md#4.%20D3/D4/D6%20—%20Composants%20%28ressource%20unique,%20ADR-11%29)
- **Architecture_Composition** : [§ Architecture — Composition (D3/D5/D6 : composant, assemblage, module)](../../3-Conception/Architecture_Composition.md#Architecture%20—%20Composition%20%28D3/D5/D6%20:%20composant,%20assemblage,%20module%29)
- **Chaincode** : [§ 4. Fonctions chaincode — Store/Read (D3/D4/D6)](../../3-Conception/Chaincode.md#4.%20Fonctions%20chaincode%20—%20Store/Read%20%28D3/D4/D6%29)
- **DC_CLI_Model** : [§ DC — CLI Modèle : Référence des commandes composant / interfaces / module](../../3-Conception/DC_CLI_Model.md#DC%20—%20CLI%20Modèle%20:%20Référence%20des%20commandes%20composant%20/%20interfaces%20/%20module) · [§ 2. Arbre de commandes](../../3-Conception/DC_CLI_Model.md#2.%20Arbre%20de%20commandes) · [§ 3.1 `myr model add`](../../3-Conception/DC_CLI_Model.md#3.1%20`myr%20model%20add`) · [§ 3.6bis `myr model submit` (RM16/RM17/RM19 — soumission, composant ou assemblage)](../../3-Conception/DC_CLI_Model.md#3.6bis%20`myr%20model%20submit`%20%28RM16/RM17/RM19%20—%20soumission,%20composant%20ou%20assemblage%29) · [§ 5. Table de correspondance méthode domaine → commande CLI → use case](../../3-Conception/DC_CLI_Model.md#5.%20Table%20de%20correspondance%20méthode%20domaine%20→%20commande%20CLI%20→%20use%20case) · [§ 6. Écarts et points ouverts](../../3-Conception/DC_CLI_Model.md#6.%20Écarts%20et%20points%20ouverts)
- **Sequence_soumission_asset** : [§ Séquence — Soumission d'un composant (D3)](../../3-Conception/Sequence_soumission_asset.md#Séquence%20—%20Soumission%20d'un%20composant%20%28D3%29)

<!-- liens-obsidian:end -->
