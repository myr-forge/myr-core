---
categorie: Composant Ecriture
titre: "Configurer un Composant"
probabilite: 3
impact: 5
importance: 15
etat: analyse
tags:
  - couche/analyse
  - type/use-case
  - famille/UCCE
  - domaine/model
  - uc/UCCE02
  - rm/RM03
  - rm/RM07
  - enf/ENF12
  - enf/ENF30
---

# Configurer un Composant

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Concepteur" as C

rectangle "Application MYR" {
    usecase "Configurer un composant" as UC1
    usecase "Modifier nom et description" as UC2
    usecase "Sélectionner la licence" as UC3
    usecase "Vérifier compatibilité de licence" as UC4
    usecase "Mettre à jour sur la blockchain" as UC5
}

C --> UC1
UC1 ..> UC2 : <<include>>
UC1 ..> UC3 : <<include>>
UC1 ..> UC4 : <<extend>> (si licence modifiée et ParentID présent)
UC1 ..> UC5 : <<include>>

@enduml
```

## Contexte

Après création, un composant peut être reconfiguré par son propriétaire : modification du nom, de la description, de la licence ou des tags. Cette opération ne modifie pas le fichier 3D ni le hash — elle met à jour les métadonnées de l'asset.

La mise à jour est soumise à la blockchain Fabric via `PATCH /api/components/:id` → `service.UpdateAsset()` → `blockchain.StoreModelRecord()`. Chaque appel produit un nouveau bloc sur le ledger : l'historique des configurations est traçable.

Cette opération est généralisée à tout asset, y compris un module (voir UCMOD03) : le même service `UpdateAsset()`, les mêmes champs modifiables (`name`, `description`, `license_id`, `tags`, `links`) et la même vérification RM03 s'appliquent, seule la route change (`PATCH /api/modules/:id`).

**Contrainte clé :** Si l'asset a un `ParentID` et que la nouvelle licence est modifiée, la compatibilité de licence avec le parent doit être re-vérifiée (RM03) — `UpdateAsset()` applique cette vérification au même titre que `AddFull()` à la création.

## Pré-conditions

- Le Concepteur est authentifié avec le rôle `contributor` (session REST valide).
- L'asset existe sur le canal Fabric et son ID est connu.
- Le Concepteur est propriétaire de l'asset (`OwnerID` correspond à son identité).

## Scénario

**Étape initiale :** `PATCH /api/components/<uuid>` est appelée (ou l'équivalent CLI `myr model update`) avec un ou plusieurs champs à modifier

### Flux nominal — Configuration réussie

1. Un ou plusieurs champs sont transmis : nom, description, licence, tags, liens (JSON body ou form).
2. Le REST Handler valide les champs (name max 256, description max 10000, IDs valides).
3. Le handler appelle `service.UpdateAsset(UpdateRequest{ID, Name, Description, LicenseID, Tags, Links})`.
4. Le service récupère l'asset existant : `blockchain.GetModelRecord(req.ID, "")`.
5. Le service applique le patch partiel : seuls les champs non vides de `UpdateRequest` écrasent les valeurs actuelles.
6. Si `LicenseID` est modifié et que l'asset a un `ParentID` : le service vérifie la compatibilité de licence avec le parent.
7. Le service soumet la transaction : `blockchain.StoreModelRecord(m)`.
8. Fabric valide et ancre le nouveau bloc.
9. L'API retourne `200 OK` avec le `Model3D` mis à jour.

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 1 déduit(s)
> - 🟡 [TestUpdateAsset_Draft_StaysLocal](../../../docs/tests/domain-model/TestUpdateAsset_Draft_StaysLocal.md) — déduit : teste `ModelService.UpdateAsset`
<!-- tests-obsidian:end -->

### Flux alternatif — Patch partiel (mise à jour d'un seul champ)

1. Seule la licence est transmise (ex : passe de CC BY à CC BY-SA).
2. Seul le champ `license_id` est envoyé dans la requête.
3. Le service applique uniquement ce champ — les autres restent inchangés (`UpdateAsset` est un patch partiel).
4. La compatibilité de licence avec le parent est vérifiée si applicable.
5. La transaction est soumise normalement.

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 1 déduit(s)
> - 🟡 [TestUpdateAsset_Draft_StaysLocal](../../../docs/tests/domain-model/TestUpdateAsset_Draft_StaysLocal.md) — déduit : teste `ModelService.UpdateAsset`
<!-- tests-obsidian:end -->

### Flux erreur — Incompatibilité de licence avec le parent (RM03)

1. L'asset a un `ParentID` et la nouvelle licence est incompatible avec la licence du parent.
2. `CheckLicenseCompatibility(parent.LicenseID, newLicenseID)` retourne `Compatible: false`.
3. Le service retourne l'erreur avant toute soumission Fabric.
4. L'API retourne `422 Unprocessable Entity` : `{ "error": "Incompatibilité de licence : <raison>." }`.

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux erreur — Asset introuvable

1. `blockchain.GetModelRecord(req.ID, "")` retourne une erreur (ID inexistant sur le canal).
2. L'API retourne `404 Not Found` : `{ "error": "asset introuvable" }`.

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux erreur — Échec endorsement Fabric

1. `blockchain.StoreModelRecord(m)` retourne une erreur Fabric.
2. L'état de l'asset en mémoire n'est pas persisté.
3. L'API retourne `500 Internal Server Error` : `{ "error": "Erreur blockchain : <message>." }`.

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

## Post-conditions

- Les métadonnées mises à jour sont inscrites sur la blockchain (nouveau bloc — immuable).
- L'historique des configurations est traçable via les blocs successifs du ledger.
- Les informations mises à jour (nom, licence) sont visibles par les autres utilisateurs du réseau.

## Diagramme de séquence

```plantuml
@startuml
participant "Client\n(CLI ou API REST)" as Browser
participant "REST Handler\n(adapters/in/rest/)" as REST
participant "Model Service\n(domain/model/)" as ModelSvc
database "Fabric\n(adapters/out/fabric/)" as Fabric

Browser -> REST : PATCH /api/components/<uuid>\n{ name?, description?, license_id?, tags?, links? }
REST -> REST : validateFields(name, description longueur)

alt Champ invalide
    REST --> Browser : 400 Bad Request { error }
else Champs valides
    REST -> ModelSvc : UpdateAsset(UpdateRequest{ID, Name?, Description?, LicenseID?, Tags?, Links?})
    ModelSvc -> Fabric : GetModelRecord(req.ID, "")

    alt Asset introuvable
        Fabric --> ModelSvc : ErrNotFound
        ModelSvc --> REST : ErrNotFound
        REST --> Browser : 404 Asset introuvable
    else Asset trouvé
        Fabric --> ModelSvc : model3D (existant)
        ModelSvc -> ModelSvc : patchPartiel(model3D, req)

        alt LicenseID modifié && ParentID présent
            ModelSvc -> Fabric : GetModelRecord(ParentID)
            Fabric --> ModelSvc : parentModel
            ModelSvc -> ModelSvc : CheckLicenseCompatibility(parent.LicenseID, newLicenseID)
            alt Incompatible
                ModelSvc --> REST : ErrLicenseIncompatibility
                REST --> Browser : 422 Incompatibilité de licence
            end
        end

        ModelSvc -> Fabric : StoreModelRecord(model3D_updated)

        alt Échec Fabric
            Fabric --> ModelSvc : ErrEndorsement
            ModelSvc --> REST : ErrBlockchain
            REST --> Browser : 500 Erreur blockchain
        else Succès
            Fabric --> ModelSvc : ok
            ModelSvc --> REST : Model3D{updated}
            REST --> Browser : 200 OK { model3D }
        end
    end
end

@enduml
```

## Règles métier déclenchées

| Règle | Description |
|-------|-------------|
| **RM03** | Si `ParentID != ""` et `LicenseID` modifié : re-vérification de compatibilité de licence obligatoire |
| **RM07** | Validation complète côté serveur avant toute soumission blockchain |

### Héritage de licences Commercial / Non-Commercial

La compatibilité de licence suit ces règles :
- Un asset Commercial peut dériver d'un asset Non-Commercial (si la licence du parent l'autorise).
- Un asset Non-Commercial peut dériver d'un asset Commercial.
- La compatibilité exacte dépend du catalogue de licences (`ListLicenses()` / `CheckLicenseCompatibility()`).

```plantuml
@startuml
skin rose
note "Asset parent → Asset dérivé\nC = Commercial  NC = Non-Commercial" as N

(Asset1.2 NC) <-- (Asset1.1 C)
(Asset2.2 C)  <-- (Asset2.1 NC)
(Asset3.3 C)  <-- (Asset3.2 NC)
(Asset3.2 NC) <-- (Asset3.1 C)
(Asset4.3 NC) <-- (Asset4.2 C)
(Asset4.2 C)  <-- (Asset4.1 NC)
@enduml
```

## Exigences non-fonctionnelles

| ID | Exigence |
|----|---------|
| **EF12** | La configuration d'un composant (nom, licence) est modifiable par son propriétaire |
| **ENF12** | Contrôle du rôle `contributor` côté serveur avant toute écriture |
| **ENF30** | En cas d'échec blockchain, l'état du composant reste celui du dernier enregistrement valide |

## Notes d'implémentation

**Route existante :** `PATCH /api/components/:id` → `handler.patchComponent()` → `service.UpdateAsset()` → `fabric.StoreModelRecord()`.

**Commande CLI équivalente :** `myr model update <id> --description <texte> --license <id> --tags <a,b>` (voir `specs/3-Conception/DC_CLI_Model.md` § 5), appelant la même méthode `service.UpdateAsset(UpdateRequest{...})` que le handler REST `updateAsset()`, avec le même comportement de patch partiel et la même vérification de compatibilité de licence (RM03).

**Patch partiel :** `service.UpdateAsset()` applique uniquement les champs non-vides de `UpdateRequest`. Les champs `Tags` et `Links` sont des slices — si `nil`, ils ne sont pas écrasés ; si `[]string{}` (slice vide), ils effacent les valeurs existantes.

**Vérification de licence manquante dans UpdateAsset :** `service.UpdateAsset()` ne vérifie pas actuellement la compatibilité de licence lors d'une modification. À ajouter : si `req.LicenseID != ""` et que l'asset a un `ParentID`, appeler `CheckLicenseCompatibility(parent.LicenseID, req.LicenseID)` avant `StoreModelRecord`.

**Catalogue de licences :** `service.ListLicenses()` et `service.GetLicense(id)` sont déjà implémentés dans `domain/model/`. L'interface REST d'exposition du catalogue (`GET /api/licenses`) est à créer.

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Étape suivante — conception**
- **API_REST** : [§ 4. D3/D4/D6 — Composants (ressource unique, ADR-11)](../../3-Conception/API_REST.md#4.%20D3/D4/D6%20—%20Composants%20%28ressource%20unique,%20ADR-11%29)
- **Architecture_Composition** : [§ Architecture — Composition (D3/D5/D6 : composant, assemblage, module)](../../3-Conception/Architecture_Composition.md#Architecture%20—%20Composition%20%28D3/D5/D6%20:%20composant,%20assemblage,%20module%29)
- **DC_CLI_Model** : [§ DC — CLI Modèle : Référence des commandes composant / interfaces / module](../../3-Conception/DC_CLI_Model.md#DC%20—%20CLI%20Modèle%20:%20Référence%20des%20commandes%20composant%20/%20interfaces%20/%20module) · [§ 2. Arbre de commandes](../../3-Conception/DC_CLI_Model.md#2.%20Arbre%20de%20commandes) · [§ 5. Table de correspondance méthode domaine → commande CLI → use case](../../3-Conception/DC_CLI_Model.md#5.%20Table%20de%20correspondance%20méthode%20domaine%20→%20commande%20CLI%20→%20use%20case)

<!-- liens-obsidian:end -->
