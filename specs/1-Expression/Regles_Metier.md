---
tags:
  - couche/expression
  - type/regles-metier
  - relecture/question
---
# Règles métier — Myr System

Ce document centralise les règles de gestion métier du projet Myr. Chaque règle indique son déclencheur, sa condition et sa conséquence observable.

Ces règles complètent les use cases : elles régissent ce que le système DOIT faire indépendamment du scénario emprunté. Elles servent de référence pour les tests de validation.

---

## 1. Assets et composants

### RM01 — Anti-plagiat obligatoire

- **Déclencheur** : Soumission d'un asset de type `base`
- **Condition** : Toujours
- **Conséquence** : Le système calcule le SHA-256 du fichier et effectue une analyse de similarité SCM. Si hash identique ou similarité > 50 % avec un asset existant → soumission rejetée, contact administration requis
- **Use cases** : UCCE01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM02 — Catégorie d'asset obligatoire

- **Déclencheur** : Création ou dérivation d'un asset
- **Condition** : Toujours
- **Conséquence** : Chaque asset doit appartenir à l'une des 8 catégories : `base`, `amélioration`, `variation`, `adaptation`, `dérivation`, `extension`, `régression`, `découpage`
- **Use cases** : UCCE01–06, UCAM05

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM03 — Compatibilité de licence

- **Déclencheur** : Soumission d'un asset avec `ParentID != ""` et `LicenseID != ""`
- **Condition** : Toujours
- **Conséquence** : La licence de l'asset dérivé doit être compatible avec celle du parent. Si incompatible → soumission rejetée
- **Use cases** : UCCE04, UCMOD06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 2 explicite(s) · 12 déduit(s)
> - ✅ [TestRunModelAdd_ServiceError_Rejected](../../docs/tests/adapters-in-cli/TestRunModelAdd_ServiceError_Rejected.md) — explicite (le test cite RM03)
> - ✅ [TestCreateModule_ParentID_RM03_LicenseIncompatible_Rejected](../../docs/tests/domain-model/TestCreateModule_ParentID_RM03_LicenseIncompatible_Rejected.md) — explicite (le test cite RM03)
> - 🟡 [TestAddFull_AllFields](../../docs/tests/domain-model/TestAddFull_AllFields.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_AlwaysDraft_NoBlockchainWrite](../../docs/tests/domain-model/TestAddFull_AlwaysDraft_NoBlockchainWrite.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_MissingFile_ReturnsError](../../docs/tests/domain-model/TestAddFull_MissingFile_ReturnsError.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_NilBlockchain_StillWorks](../../docs/tests/domain-model/TestAddFull_NilBlockchain_StillWorks.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_NilDraftStore_Rejected](../../docs/tests/domain-model/TestAddFull_NilDraftStore_Rejected.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_NoFile_NoVersionNoHash](../../docs/tests/domain-model/TestAddFull_NoFile_NoVersionNoHash.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_STLAndSTEP_BothAccepted](../../docs/tests/domain-model/TestAddFull_STLAndSTEP_BothAccepted.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_WithFile_HashAndVersionCreated](../../docs/tests/domain-model/TestAddFull_WithFile_HashAndVersionCreated.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_WithParent](../../docs/tests/domain-model/TestAddFull_WithParent.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestCreateModule_ParentID_CopiesComposition](../../docs/tests/domain-model/TestCreateModule_ParentID_CopiesComposition.md) — déduit : teste `ModelService.CreateModule`
> - 🟡 [TestCreateModule_ParentID_NotFound_Rejected](../../docs/tests/domain-model/TestCreateModule_ParentID_NotFound_Rejected.md) — déduit : teste `ModelService.CreateModule`
> - 🟡 [TestUpdateAsset_Draft_StaysLocal](../../docs/tests/domain-model/TestUpdateAsset_Draft_StaysLocal.md) — déduit : teste `ModelService.UpdateAsset`
<!-- tests-obsidian:end -->

### RM04 — UUID unique

- **Déclencheur** : Enregistrement blockchain réussi
- **Condition** : Toujours
- **Conséquence** : Un UUID unique est généré et attribué à l'asset. Il est immuable et ne peut pas être modifié ultérieurement
- **Use cases** : UCCE01–06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 3 déduit(s)
> - 🟡 [TestCreateModule_ParentID_CopiesComposition](../../docs/tests/domain-model/TestCreateModule_ParentID_CopiesComposition.md) — déduit : teste `ModelService.CreateModule`
> - 🟡 [TestCreateModule_ParentID_NotFound_Rejected](../../docs/tests/domain-model/TestCreateModule_ParentID_NotFound_Rejected.md) — déduit : teste `ModelService.CreateModule`
> - 🟡 [TestCreateModule_ParentID_RM03_LicenseIncompatible_Rejected](../../docs/tests/domain-model/TestCreateModule_ParentID_RM03_LicenseIncompatible_Rejected.md) — déduit : teste `ModelService.CreateModule`
<!-- tests-obsidian:end -->

### RM05 — ParentID obligatoire pour les dérivés

- **Déclencheur** : Création d'un asset non-`base`
- **Condition** : Catégorie ≠ `base`
- **Conséquence** : L'asset dérivé doit référencer un `ParentID` valide. Un asset `base` n'a pas de parent
- **Use cases** : UCCE04, UCCE05

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 9 déduit(s)
> - 🟡 [TestAddFull_AllFields](../../docs/tests/domain-model/TestAddFull_AllFields.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_AlwaysDraft_NoBlockchainWrite](../../docs/tests/domain-model/TestAddFull_AlwaysDraft_NoBlockchainWrite.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_MissingFile_ReturnsError](../../docs/tests/domain-model/TestAddFull_MissingFile_ReturnsError.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_NilBlockchain_StillWorks](../../docs/tests/domain-model/TestAddFull_NilBlockchain_StillWorks.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_NilDraftStore_Rejected](../../docs/tests/domain-model/TestAddFull_NilDraftStore_Rejected.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_NoFile_NoVersionNoHash](../../docs/tests/domain-model/TestAddFull_NoFile_NoVersionNoHash.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_STLAndSTEP_BothAccepted](../../docs/tests/domain-model/TestAddFull_STLAndSTEP_BothAccepted.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_WithFile_HashAndVersionCreated](../../docs/tests/domain-model/TestAddFull_WithFile_HashAndVersionCreated.md) — déduit : teste `ModelService.AddFull`
> - 🟡 [TestAddFull_WithParent](../../docs/tests/domain-model/TestAddFull_WithParent.md) — déduit : teste `ModelService.AddFull`
<!-- tests-obsidian:end -->

### RM42 — Traçabilité et alerte des emplacements externes

- **Déclencheur** : Consultation d'un composant ou module portant un ou plusieurs emplacements externes (`Locations` non vide)
- **Condition** : Toujours
- **Conséquence** : Myr ne conserve jamais de copie du fichier ressource — priorité à la traçabilité (savoir où un composant ou produit existe) plutôt qu'à l'hébergement. Le système sait, à la demande, vérifier l'accessibilité de chaque emplacement externe déclaré (boutique, dépôt de fichiers tiers…) au moyen d'une requête sur son URL. Le résultat (accessible / inaccessible) et la date de vérification sont mis à jour individuellement pour chaque emplacement — jamais agrégés en un statut global qui masquerait lequel a échoué. Cette vérification porte uniquement sur l'accessibilité de l'emplacement, jamais sur le contenu qui y est exposé : Myr ne conservant aucune copie du fichier, `Model3D.Hash` (calculé une fois à la soumission, RM01) ne peut être recomparé qu'au contenu que le demandeur récupère lui-même depuis un emplacement — jamais automatiquement par Myr, qui n'a lui-même rien à comparer
- **Use cases** : UCCL03

---

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

## 2. Blockchain et immuabilité

### RM06 — Immuabilité des transactions

- **Déclencheur** : Toute soumission blockchain
- **Condition** : Toujours
- **Conséquence** : Une transaction soumise ne peut pas être modifiée ou annulée. Il n'existe aucune opération de suppression sur la blockchain
- **Use cases** : UCCE01, UCMOD06, UCPI07–09

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 3 déduit(s)
> - 🟡 [TestRemove_Draft_RemovesLocally](../../docs/tests/domain-model/TestRemove_Draft_RemovesLocally.md) — déduit : teste `ModelService.Remove`
> - 🟡 [TestRemove_Submitted_HidesLocally](../../docs/tests/domain-model/TestRemove_Submitted_HidesLocally.md) — déduit : teste `ModelService.Remove`
> - 🟡 [TestRemove_UnknownID_Rejected](../../docs/tests/domain-model/TestRemove_UnknownID_Rejected.md) — déduit : teste `ModelService.Remove`
<!-- tests-obsidian:end -->

### RM07 — Validation préalable obligatoire

- **Déclencheur** : Avant toute soumission blockchain
- **Condition** : Toujours
- **Conséquence** : Toutes les données (métadonnées, licences, interfaces, UUID) sont validées côté serveur avant soumission. Une validation échouée ne génère aucune transaction
- **Use cases** : UCCE01, UCMOD06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 1 explicite(s) · 0 déduit(s)
> - ✅ [TestAddOrganisation_InvalidMSPID_Rejected](../../docs/tests/domain-channel/TestAddOrganisation_InvalidMSPID_Rejected.md) — explicite (le test cite RM07)
<!-- tests-obsidian:end -->

### RM08 — Masquage local, ledger jamais modifié

- **Déclencheur** : Suppression d'un asset (composant ou module) par son propriétaire
- **Condition** : Toujours
- **Conséquence** : Aucune opération de suppression n'existe sur la blockchain : un asset déjà soumis y reste inscrit en permanence. Le système retire uniquement l'asset des listes retournées au demandeur (recherche, ateliers, catalogue) — sa consultation directe par identifiant reste possible. Aucune erreur non contrôlée (`panic`) ne doit se produire, qu'il s'agisse d'un brouillon (retiré du stockage local) ou d'un asset déjà soumis (masqué)
- **Use cases** : —

---

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 2 explicite(s) · 2 déduit(s)
> - ✅ [TestRemoveModule_Submitted_HidesLocally](../../docs/tests/domain-model/TestRemoveModule_Submitted_HidesLocally.md) — explicite (le test cite RM08)
> - ✅ [TestRemove_Submitted_HidesLocally](../../docs/tests/domain-model/TestRemove_Submitted_HidesLocally.md) — explicite (le test cite RM08)
> - 🟡 [TestRemove_Draft_RemovesLocally](../../docs/tests/domain-model/TestRemove_Draft_RemovesLocally.md) — déduit : teste `ModelService.Remove`
> - 🟡 [TestRemove_UnknownID_Rejected](../../docs/tests/domain-model/TestRemove_UnknownID_Rejected.md) — déduit : teste `ModelService.Remove`
<!-- tests-obsidian:end -->

## 3. Interfaces et liaisons

### RM09 — Interface à usage unique

- **Déclencheur** : Tentative de liaison d'une interface déjà engagée
- **Condition** : Interface source ou cible déjà utilisée dans une liaison
- **Conséquence** : La liaison est refusée
- **Use cases** : UCAM01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 1 explicite(s) · 3 déduit(s)
> - ✅ [TestRunModelLinkAdd_MissingFrom_Rejected](../../docs/tests/adapters-in-cli/TestRunModelLinkAdd_MissingFrom_Rejected.md) — explicite (le test cite RM09)
> - 🟡 [TestAddAssemblyLink](../../docs/tests/domain-model/TestAddAssemblyLink.md) — déduit : teste `ModelService.AddAssemblyLink`
> - 🟡 [TestAddAssemblyLink_MissingStore](../../docs/tests/domain-model/TestAddAssemblyLink_MissingStore.md) — déduit : teste `ModelService.AddAssemblyLink`
> - 🟡 [TestAddAssemblyLink_UnknownInterface](../../docs/tests/domain-model/TestAddAssemblyLink_UnknownInterface.md) — déduit : teste `ModelService.AddAssemblyLink`
<!-- tests-obsidian:end -->

### RM10 — Vérification de compatibilité automatique

- **Déclencheur** : Création de toute liaison
- **Condition** : Toujours
- **Conséquence** : Le système appelle `ifacesCompatible` pour vérifier la cohérence de la paire avant d'enregistrer la liaison. Une liaison incompatible est refusée par le service. L'incompatibilité post-création est couverte par RM12
- **Use cases** : UCAM01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 2 explicite(s) · 3 déduit(s)
> - ✅ [TestRunModelLinkAdd_Incompatible](../../docs/tests/adapters-in-cli/TestRunModelLinkAdd_Incompatible.md) — explicite (le test cite RM10)
> - ✅ [TestRunModelLinkAdd_MissingFrom_Rejected](../../docs/tests/adapters-in-cli/TestRunModelLinkAdd_MissingFrom_Rejected.md) — explicite (le test cite RM10)
> - 🟡 [TestAddAssemblyLink](../../docs/tests/domain-model/TestAddAssemblyLink.md) — déduit : teste `ModelService.AddAssemblyLink`
> - 🟡 [TestAddAssemblyLink_MissingStore](../../docs/tests/domain-model/TestAddAssemblyLink_MissingStore.md) — déduit : teste `ModelService.AddAssemblyLink`
> - 🟡 [TestAddAssemblyLink_UnknownInterface](../../docs/tests/domain-model/TestAddAssemblyLink_UnknownInterface.md) — déduit : teste `ModelService.AddAssemblyLink`
<!-- tests-obsidian:end -->

### RM11 — Critères de compatibilité d'interfaces

- **Déclencheur** : Appel de `ifacesCompatible`
- **Condition** : Toujours
- **Conséquence** : Deux interfaces sont compatibles si : (1) même catégorie (Électrique, Mécanique, Numérique…), (2) même tag (Câble, vis…) si renseigné, (3) même type (USB-C, UART…) si renseigné, (4) sens complémentaires (♂/♀ ou bidirectionnel), (5) plages de valeurs se chevauchant (`ValueMin`/`ValueMax`)
- **Use cases** : UCAM01, UCAM03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 1 explicite(s) · 0 déduit(s)
> - ✅ [TestRunModelLinkAdd_Incompatible](../../docs/tests/adapters-in-cli/TestRunModelLinkAdd_Incompatible.md) — explicite (le test cite RM11)
<!-- tests-obsidian:end -->

### RM12 — Persistance des liaisons incompatibles

- **Déclencheur** : Liaison devenant incompatible après modification d'un asset
- **Condition** : Toujours
- **Conséquence** : La liaison n'est pas supprimée automatiquement. Elle passe à `Incompatible: true` et doit rester accessible jusqu'à suppression manuelle par l'utilisateur
- **Use cases** : UCAM01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM13 — Slot virtuel garanti

- **Déclencheur** : Ajout d'un asset à un module
- **Condition** : Toujours
- **Conséquence** : Chaque asset dispose en permanence d'au moins un slot virtuel (`Virtual: true`). Dès qu'un slot virtuel est matérialisé en interface, un nouveau slot virtuel est recréé automatiquement
- **Use cases** : UCAM02, UCAM03

---

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 16 déduit(s)
> - 🟡 [TestConnectVirtualToPhysical_AssemblyLinkCreated](../../docs/tests/domain-model/TestConnectVirtualToPhysical_AssemblyLinkCreated.md) — déduit : teste `ModelService.ConnectVirtualToPhysical`
> - 🟡 [TestConnectVirtualToPhysical_BidirStaysBidir](../../docs/tests/domain-model/TestConnectVirtualToPhysical_BidirStaysBidir.md) — déduit : teste `ModelService.ConnectVirtualToPhysical`
> - 🟡 [TestConnectVirtualToPhysical_InBecomesOut](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InBecomesOut.md) — déduit : teste `ModelService.ConnectVirtualToPhysical`
> - 🟡 [TestConnectVirtualToPhysical_InheritsPhysicalName](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InheritsPhysicalName.md) — déduit : teste `ModelService.ConnectVirtualToPhysical`
> - 🟡 [TestConnectVirtualToPhysical_InheritsPhysicalValues](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InheritsPhysicalValues.md) — déduit : teste `ModelService.ConnectVirtualToPhysical`
> - 🟡 [TestConnectVirtualToPhysical_NilConnStore_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NilConnStore_Error.md) — déduit : teste `ModelService.ConnectVirtualToPhysical`
> - 🟡 [TestConnectVirtualToPhysical_NilIfaceStore_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NilIfaceStore_Error.md) — déduit : teste `ModelService.ConnectVirtualToPhysical`
> - 🟡 [TestConnectVirtualToPhysical_NotVirtual_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NotVirtual_Error.md) — déduit : teste `ModelService.ConnectVirtualToPhysical`
> - 🟡 [TestConnectVirtualToPhysical_OutBecomesIn](../../docs/tests/domain-model/TestConnectVirtualToPhysical_OutBecomesIn.md) — déduit : teste `ModelService.ConnectVirtualToPhysical`
> - 🟡 [TestConnectVirtualToPhysical_UserOverridesName](../../docs/tests/domain-model/TestConnectVirtualToPhysical_UserOverridesName.md) — déduit : teste `ModelService.ConnectVirtualToPhysical`
> - 🟡 [TestConnectVirtualToPhysical_UserOverridesValues](../../docs/tests/domain-model/TestConnectVirtualToPhysical_UserOverridesValues.md) — déduit : teste `ModelService.ConnectVirtualToPhysical`
> - 🟡 [TestGetModuleInterfaces_InternalConnection_ScopedToInstance](../../docs/tests/domain-model/TestGetModuleInterfaces_InternalConnection_ScopedToInstance.md) — déduit : teste `ModelService.GetModuleInterfaces`
> - 🟡 [TestGetModuleInterfaces_InternalConnectionsConsumed](../../docs/tests/domain-model/TestGetModuleInterfaces_InternalConnectionsConsumed.md) — déduit : teste `ModelService.GetModuleInterfaces`
> - 🟡 [TestGetModuleInterfaces_MultipleInstancesSameAsset_AllExposed](../../docs/tests/domain-model/TestGetModuleInterfaces_MultipleInstancesSameAsset_AllExposed.md) — déduit : teste `ModelService.GetModuleInterfaces`
> - 🟡 [TestGetModuleInterfaces_NilStore](../../docs/tests/domain-model/TestGetModuleInterfaces_NilStore.md) — déduit : teste `ModelService.GetModuleInterfaces`
> - … et 1 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

## 4. Composition d'un Module (instances)

### RM14 — Suppression en cascade des connexions

- **Déclencheur** : Retrait d'une instance d'un module (`RemoveAssetFromWorkspace`)
- **Condition** : Toujours
- **Conséquence** : Toutes les connexions de l'instance retirée sont supprimées automatiquement. Les assets liés restent dans le module mais leurs interfaces concernées redeviennent libres
- **Use cases** : UCAM08

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM15 — Instance indépendante

- **Déclencheur** : Ajout d'un module déjà instancié dans le module hôte
- **Condition** : Module déjà instancié
- **Conséquence** : Une seconde instance est créée avec ses propres connexions, indépendantes de la première
- **Use cases** : UCMOD02

---

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 1 explicite(s) · 0 déduit(s)
> - ✅ [TestRunModelInstanceRemove_NominalCase](../../docs/tests/adapters-in-cli/TestRunModelInstanceRemove_NominalCase.md) — explicite (le test cite RM15)
<!-- tests-obsidian:end -->

## 5. Modules et cycle de vie des assets

> RM16 et RM19 sont génériques à tout `Model3D` (composant ou module) — RM17 et RM18 restent spécifiques au module (une notion d'« assemblage » ou de « ModuleVersion » n'a pas de sens pour un composant simple). Voir `specs/3-Conception/Conception_intro.md` ADR-02.

### RM16 — État draft

- **Déclencheur** : Création d'un asset (composant ou module)
- **Condition** : Un module est **toujours** créé en `draft` (RM17 l'exige avant soumission). Un composant est créé **directement `submitted`** par défaut, sauf si `draft: true` est explicitement demandé (`myr model add --draft`)
- **Conséquence** : L'asset en `draft` n'est pas visible sur le réseau comme définitif. Ses interfaces et attributs mutables peuvent être modifiés librement (localement) tant qu'il n'est pas soumis (ADR-02)
- **Use cases** : UCMOD01, UCCE01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 3 déduit(s)
> - 🟡 [TestCreateModule_ParentID_CopiesComposition](../../docs/tests/domain-model/TestCreateModule_ParentID_CopiesComposition.md) — déduit : teste `ModelService.CreateModule`
> - 🟡 [TestCreateModule_ParentID_NotFound_Rejected](../../docs/tests/domain-model/TestCreateModule_ParentID_NotFound_Rejected.md) — déduit : teste `ModelService.CreateModule`
> - 🟡 [TestCreateModule_ParentID_RM03_LicenseIncompatible_Rejected](../../docs/tests/domain-model/TestCreateModule_ParentID_RM03_LicenseIncompatible_Rejected.md) — déduit : teste `ModelService.CreateModule`
<!-- tests-obsidian:end -->

### RM17 — Assemblage requis pour soumission (module uniquement)

- **Déclencheur** : Appel de `SubmitModule`
- **Condition** : Toujours
- **Conséquence** : Le module doit contenir au moins une liaison entre composants. Si aucune liaison → soumission rejetée avec message explicite. Ne s'applique pas à un composant : sa soumission ne requiert aucun assemblage
- **Use cases** : UCMOD06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM18 — ModuleVersion immuable (module uniquement)

- **Déclencheur** : Soumission réussie d'un module
- **Condition** : Toujours
- **Conséquence** : Une `ModuleVersion` est créée avec un hash de l'assemblage et un horodatage. Ce snapshot est immuable. Toute modification ultérieure exige la création d'une nouvelle version. Pour un composant, l'équivalent est une entrée `Versions[]` (déjà utilisée par UCCE02) — pas de `ModuleVersion`
- **Use cases** : UCMOD06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM19 — Fork d'un asset soumis

- **Déclencheur** : Modification d'un asset soumis (composant ou module)
- **Condition** : Toujours
- **Conséquence** : Un asset soumis ne peut pas être modifié directement (immuabilité Fabric, règle 7). Une nouvelle version (fork) doit être créée en état `draft`, avec `ParentID` référençant l'asset d'origine
- **Use cases** : UCMOD06, UCCE06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 1 déduit(s)
> - 🟡 [TestUpdateAsset_Draft_StaysLocal](../../docs/tests/domain-model/TestUpdateAsset_Draft_StaysLocal.md) — déduit : teste `ModelService.UpdateAsset`
<!-- tests-obsidian:end -->

### RM39 — Filiation d'un découpage

- **Déclencheur** : Transformation d'un composant en module de catégorie `découpage`, manuelle (UCAM05) ou assistée par une proposition automatique (UCAM09)
- **Condition** : Toujours
- **Conséquence** : Les sous-composants créés et le module `découpage` référencent tous `ParentID` vers le composant d'origine (généralisation de RM05 à cette transformation). Le composant d'origine n'est ni supprimé ni modifié
- **Use cases** : UCAM05, UCAM09

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM40 — Proposition de découpage non engageante

- **Déclencheur** : Analyse d'un composant en vue d'un découpage automatique (UCAM09)
- **Condition** : Toujours
- **Conséquence** : L'analyse ne crée aucune entité persistée (sous-composant, module, interface, liaison) : elle retourne une proposition temporaire, valable uniquement pour sa relecture par le Concepteur. Rien n'existe côté serveur tant que la proposition n'a pas été validée explicitement
- **Use cases** : UCAM09

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM41 — Compatibilité toujours vérifiée pour une connexion suggérée

- **Déclencheur** : Validation d'une proposition de découpage dont une connexion candidate a été conservée (UCAM09)
- **Condition** : Toujours
- **Conséquence** : La connexion passe par le même contrôle de compatibilité qu'une liaison créée manuellement (RM10/RM11) — aucune dérogation pour une suggestion automatique. Si elle échoue, elle n'est pas créée et est rapportée en avertissement ; le score de confiance de la suggestion n'a aucun statut privilégié dans ce contrôle
- **Use cases** : UCAM09

---

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

## 6. Compte et accès

### RM20 — Identité = enrôlement CA, pas un compte séparé

- **Déclencheur** : Connexion (`POST /api/identity/session`)
- **Condition** : Toujours
- **Conséquence** : Il n'existe pas de « création de compte » distincte de l'identité blockchain : se connecter, c'est enrôler (ou ré-enrôler) l'identité auprès de la Fabric CA avec le secret fourni. Aucun compte email/mot de passe n'est créé au préalable
- **Use cases** : UCA01, UCA02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 8 déduit(s)
> - 🟡 [TestChannels_PUT_SwitchesChannel](../../docs/tests/adapters-in-rest/TestChannels_PUT_SwitchesChannel.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_EnrollError](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_EnrollError.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_InvalidJSON](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_InvalidJSON.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MethodNotAllowed.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_MissingFields](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MissingFields.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_NoIdentityService](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_NoIdentityService.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_Success](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_Success.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestRequireAuth_BlockchainWithValidToken](../../docs/tests/adapters-in-rest/TestRequireAuth_BlockchainWithValidToken.md) — déduit : teste route `/api/identity/session`
<!-- tests-obsidian:end -->

### RM21 — Rôle Lecteur par défaut à l'auto-enregistrement

- **Déclencheur** : `POST /api/identity/request` avec auto-enregistrement réseau actif (`AllowAutoRegister`)
- **Condition** : Le profil réseau ne définit pas explicitement un autre rôle (`AutoRegisterRole`)
- **Conséquence** : Une identité auto-enregistrée reçoit par défaut le rôle **Lecteur** (`reader`), sauf si le profil réseau configure explicitement un autre rôle initial
- **Use cases** : UCA01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 4 déduit(s)
> - 🟡 [TestHandleIdentityRequest_InvalidJSON](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_InvalidJSON.md) — déduit : teste route `/api/identity/request`
> - 🟡 [TestHandleIdentityRequest_MissingFields](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_MissingFields.md) — déduit : teste route `/api/identity/request`
> - 🟡 [TestHandleIdentityRequest_ServiceError](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_ServiceError.md) — déduit : teste route `/api/identity/request`
> - 🟡 [TestHandleIdentityRequest_Success](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_Success.md) — déduit : teste route `/api/identity/request`
<!-- tests-obsidian:end -->

### RM22 — Changement de rôle réservé à l'administrateur

- **Déclencheur** : `myr identity set-role`
- **Condition** : Toujours
- **Conséquence** : Seul l'administrateur peut modifier le rôle (`Myr.role`) d'une identité existante auprès de la CA. Le nouveau rôle ne s'applique qu'au prochain ré-enrôlement de l'identité — ce n'est pas immédiat
- **Use cases** : UCA08

---

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

## 7. Propriété intellectuelle et commissions

### RM23 — Distribution automatique des commissions

- **Déclencheur** : Livraison d'un composant ou module commandé
- **Condition** : Toujours
- **Conséquence** : Le smart contract calcule et distribue les commissions automatiquement à chaque auteur dans la chaîne de propriété. Chaque transaction de commission est enregistrée individuellement sur la blockchain
- **Use cases** : UCPI02, UCAUT01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM24 — Répartition proportionnelle multi-auteurs

- **Déclencheur** : Module commandé avec plusieurs auteurs
- **Condition** : Plusieurs concepteurs impliqués dans le module
- **Conséquence** : Chaque auteur reçoit une commission proportionnelle à son apport. Les transactions sont indépendantes par auteur
- **Use cases** : UCPI02

#question RM23/RM24 — « chaîne de propriété » est ambiguë sur son périmètre : couvre-t-elle uniquement les auteurs des composants directement constitutifs d'un module livré (composition — c'est le seul cas couvert par l'algorithme conçu dans `specs/3-Conception/DC_D7_Payment.md` §4, qui parcourt récursivement les composants du module sans jamais remonter un `ParentID`), ou aussi les auteurs de la lignée de dérivation de chaque composant (dérivation — un composant `base` toucherait alors une part quand ses dérivés sont vendus, comportement évoqué dans `specs/2-Analyse/UCPI-Propriete_Intellectuelle/UCPI02.md` et `UCPI04.md`) ? Les deux couches de specs divergent actuellement sur ce point. À trancher par le product owner avant toute implémentation de RM23/RM24 — le choix conditionne le schéma de l'entité `Commission` (`DC_D7_Payment.md` §3) et l'algorithme de distribution.

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM25 — Transfert de propriété définitif

- **Déclencheur** : Transfert accepté par le destinataire
- **Condition** : Toujours
- **Conséquence** : La propriété est transférée de manière immuable. L'ancien propriétaire perd immédiatement les droits d'édition
- **Use cases** : UCPI07

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM26 — Traçabilité du clonage inter-réseaux

- **Déclencheur** : Clonage d'un asset sur un réseau externe
- **Condition** : Toujours
- **Conséquence** : Le clone conserve l'UUID et les métadonnées originales. La transaction de clonage est enregistrée sur les deux réseaux pour assurer la traçabilité de l'origine
- **Use cases** : UCPI08, UCPI09

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM29 — Taux de commission défini par le réseau

- **Déclencheur** : Calcul des commissions à la livraison (RM23)
- **Condition** : Toujours
- **Conséquence** : Le taux de commission (`commission_rate`) est une propriété du réseau, définie par l'administrateur (défaut : 10 %). Il s'applique uniformément à tous les assets du réseau. Un auteur ne peut pas définir son propre taux — il accepte le taux du réseau en publiant sur celui-ci
- **Use cases** : UCPI02, UCADM02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM30 — Calcul automatique du prix d'un module

- **Déclencheur** : Commande d'un module (UCPI01)
- **Condition** : Le module n'a pas de prix défini par l'auteur
- **Conséquence** : Si l'auteur n'a pas défini de prix via UCPI05, le prix du module est calculé automatiquement comme la somme des prix unitaires de ses composants constitutifs. Si un composant est lui-même sans prix, il est compté à 0
- **Use cases** : UCPI01, UCPI05

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM31 — Modification de prix — effet sur les commandes futures uniquement

- **Déclencheur** : Mise à jour du prix d'un asset (UCPI11)
- **Condition** : Toujours
- **Conséquence** : La modification d'un prix ne s'applique qu'aux commandes passées après la mise à jour. Les commandes en cours ou livrées conservent le prix enregistré à leur création (`OrderItem.UnitPrice`). Le prix n'est pas stocké sur la blockchain Fabric — il est local et mutable
- **Use cases** : UCPI11, UCPI04, UCPI05

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM32 — Asset à prix nul — librement disponible

- **Déclencheur** : Commande ou utilisation d'un asset dont `AssetPrice.Price = 0` ou sans `AssetPrice` défini
- **Condition** : Toujours
- **Conséquence** : L'asset est disponible sans frais. Aucune commission n'est générée pour ses auteurs. La liberté d'accès n'est pas liée à la licence AGPL du code Myr — un asset physique peut être libre de droits commerciaux tout en étant documenté sur le réseau
- **Use cases** : UCPI04, UCPI02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM33 — Devise unique par réseau

- **Déclencheur** : Définition ou comparaison de prix entre assets d'un même réseau
- **Condition** : Toujours
- **Conséquence** : Tous les prix et commissions d'un réseau sont exprimés dans la devise définie par l'administrateur à la création du réseau. Aucune conversion de devise n'est effectuée par le système. Un asset cloné sur un réseau étranger (UCPI08/09) adopte la devise du réseau cible
- **Use cases** : UCPI04, UCPI05, UCADM02

---

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

## 8. Administration réseau

### RM27 — Nombre minimum de nœuds actifs

- **Déclencheur** : Demande de retrait d'un nœud du canal (UCADM04)
- **Condition** : Le retrait provoquerait un passage sous 3 nœuds actifs sur le canal
- **Conséquence** : L'opération est refusée ; aucune transaction n'est soumise. Message d'erreur explicite indiquant le nombre de nœuds actifs courant
- **Use cases** : UCADM04

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 3 explicite(s) · 0 déduit(s)
> - ✅ [TestRunNodeRemove_RM27_MinNodesRequired_Rejected](../../docs/tests/adapters-in-cli/TestRunNodeRemove_RM27_MinNodesRequired_Rejected.md) — explicite (le test cite RM27)
> - ✅ [TestRemoveNode_RM27_MinNodesRequired_Rejected](../../docs/tests/domain-channel/TestRemoveNode_RM27_MinNodesRequired_Rejected.md) — explicite (le test cite RM27)
> - ✅ [TestCreate_DefaultsNumPeers](../../docs/tests/domain-network/TestCreate_DefaultsNumPeers.md) — explicite (le test cite RM27)
<!-- tests-obsidian:end -->

### RM28 — Démantèlement réseau : opération d'infrastructure locale

- **Déclencheur** : Commande `myr network destroy` (UCADM05)
- **Condition** : Toujours
- **Conséquence** : L'opération est purement locale (arrêt de processus OS, suppression de fichiers). Aucune transaction Fabric n'est soumise. RM06 et RM08 ne s'appliquent pas — la blockchain n'est pas impliquée. Confirmation explicite (`--confirm`) obligatoire. Refusée sur tout réseau marqué `IsProduction: true`.
- **Use cases** : UCADM05

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM34 — Rôle admin protégé

- **Déclencheur** : Tentative d'édition ou de suppression du rôle `admin`
- **Condition** : ID du rôle = `"admin"`
- **Conséquence** : L'opération est refusée par le service domaine. Le rôle `admin` ne peut pas non plus être attribué à une organisation tierce via UCADM06.
- **Use cases** : UCADM06, UCADM07

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM35 — Révocation en cascade à la suppression d'un rôle

- **Déclencheur** : Suppression d'un rôle (UCADM07)
- **Condition** : Toujours
- **Conséquence** : Toutes les liaisons organisation ↔ rôle référençant ce rôle sont supprimées automatiquement avant la suppression du rôle. Les droits correspondants sont révoqués immédiatement pour les organisations concernées.
- **Use cases** : UCADM07

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM36 — Nom de rôle unique

- **Déclencheur** : Création d'un rôle (UCADM07)
- **Condition** : Un rôle avec le même nom existe déjà
- **Conséquence** : La création est refusée. Le nom d'un rôle est unique dans le système (insensible à la casse).
- **Use cases** : UCADM07

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM37 — Multi-rôles par organisation

- **Déclencheur** : Attribution d'un rôle à une organisation (UCADM06)
- **Condition** : L'organisation possède déjà un ou plusieurs rôles
- **Conséquence** : L'organisation peut posséder plusieurs rôles simultanément. Ses droits effectifs sont l'union des droits de tous ses rôles actifs.
- **Use cases** : UCADM06

---

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

## 9. Résumé — index des règles

| ID | Règle (résumé) | Domaine |
|----|---------------|---------|
| [RM01](#RM01%20—%20Anti-plagiat%20obligatoire) | Anti-plagiat obligatoire (SHA-256 + SCM > 50 %) | Assets |
| [RM02](#RM02%20—%20Catégorie%20d'asset%20obligatoire) | Catégorie d'asset obligatoire (8 types) | Assets |
| [RM03](#RM03%20—%20Compatibilité%20de%20licence) | Compatibilité de licence pour tout asset dérivé | Assets |
| [RM04](#RM04%20—%20UUID%20unique) | UUID unique et immuable à l'enregistrement | Assets |
| [RM05](#RM05%20—%20ParentID%20obligatoire%20pour%20les%20dérivés) | ParentID obligatoire pour les assets non-`base` | Assets |
| [RM42](#RM42%20—%20Traçabilité%20et%20alerte%20des%20emplacements%20externes) | Aucune copie du fichier ressource conservée par Myr ; vérification à la demande de l'accessibilité de chaque emplacement externe, statut individuel par emplacement | Assets |
| [RM06](#RM06%20—%20Immuabilité%20des%20transactions) | Transactions blockchain immuables, pas de suppression | Blockchain |
| [RM07](#RM07%20—%20Validation%20préalable%20obligatoire) | Validation complète avant toute soumission blockchain | Blockchain |
| [RM08](#RM08%20—%20Masquage%20local,%20ledger%20jamais%20modifié) | ErrNotSupported retourné sur tentative de suppression | Blockchain |
| [RM09](#RM09%20—%20Interface%20à%20usage%20unique) | Interface à usage unique dans une liaison | Interfaces |
| [RM10](#RM10%20—%20Vérification%20de%20compatibilité%20automatique) | Vérification de compatibilité automatique à la création | Interfaces |
| [RM11](#RM11%20—%20Critères%20de%20compatibilité%20d'interfaces) | Critères de compatibilité : catégorie + tag + type + sens + plage de valeurs | Interfaces |
| [RM12](#RM12%20—%20Persistance%20des%20liaisons%20incompatibles) | Liaison incompatible persistante (`Incompatible:true`, non supprimée automatiquement) | Interfaces |
| [RM13](#RM13%20—%20Slot%20virtuel%20garanti) | Slot virtuel garanti et recréé automatiquement à chaque matérialisation | Interfaces |
| [RM14](#RM14%20—%20Suppression%20en%20cascade%20des%20connexions) | Suppression en cascade des connexions au retrait d'une instance | Composition (instances) |
| [RM15](#RM15%20—%20Instance%20indépendante) | Seconde instance indépendante si module déjà instancié dans le module hôte | Composition (instances) |
| [RM16](#RM16%20—%20État%20draft) | État draft générique (module : toujours ; composant : optionnel via `--draft`) | Modules |
| [RM17](#RM17%20—%20Assemblage%20requis%20pour%20soumission%20%28module%20uniquement%29) | Au moins un assemblage requis pour soumettre (module uniquement) | Modules |
| [RM18](#RM18%20—%20ModuleVersion%20immuable%20%28module%20uniquement%29) | ModuleVersion immuable horodatée à la soumission (module uniquement) | Modules |
| [RM19](#RM19%20—%20Fork%20d'un%20asset%20soumis) | Toute modification d'un asset soumis (composant ou module) crée un fork | Modules |
| [RM39](#RM39%20—%20Filiation%20d'un%20découpage) | Filiation par `ParentID` du composant d'origine pour tous les éléments d'un découpage (manuel ou assisté) | Modules |
| [RM40](#RM40%20—%20Proposition%20de%20découpage%20non%20engageante) | Une proposition de découpage automatique ne crée aucune entité tant qu'elle n'est pas validée | Modules |
| [RM41](#RM41%20—%20Compatibilité%20toujours%20vérifiée%20pour%20une%20connexion%20suggérée) | Connexion suggérée soumise au même contrôle de compatibilité qu'une liaison manuelle, sans dérogation | Modules |
| [RM20](#RM20%20—%20Identité%20=%20enrôlement%20CA,%20pas%20un%20compte%20séparé) | Identité = enrôlement CA — pas de compte email/mot de passe séparé | Compte |
| [RM21](#RM21%20—%20Rôle%20Lecteur%20par%20défaut%20à%20l'auto-enregistrement) | Rôle Lecteur par défaut à l'auto-enregistrement, sauf rôle explicite | Compte |
| [RM22](#RM22%20—%20Changement%20de%20rôle%20réservé%20à%20l'administrateur) | Changement de rôle réservé à l'administrateur (`myr identity set-role`) | Compte |
| [RM23](#RM23%20—%20Distribution%20automatique%20des%20commissions) | Commissions distribuées automatiquement à la livraison | PI |
| [RM24](#RM24%20—%20Répartition%20proportionnelle%20multi-auteurs) | Répartition proportionnelle par auteur | PI |
| [RM25](#RM25%20—%20Transfert%20de%20propriété%20définitif) | Transfert de propriété définitif et immuable | PI |
| [RM26](#RM26%20—%20Traçabilité%20du%20clonage%20inter-réseaux) | Traçabilité UUID préservée lors du clonage inter-réseaux | PI |
| [RM27](#RM27%20—%20Nombre%20minimum%20de%20nœuds%20actifs) | Retrait d'un nœud refusé si le canal passerait sous 3 nœuds actifs | Administration réseau |
| [RM28](#RM28%20—%20Démantèlement%20réseau%20:%20opération%20d'infrastructure%20locale) | Démantèlement réseau = infrastructure locale uniquement, hors blockchain, confirmation obligatoire | Administration réseau |
| [RM34](#RM34%20—%20Rôle%20admin%20protégé) | Rôle `admin` protégé — ni modifiable, ni supprimable, ni attribuable à une organisation tierce | Administration réseau |
| [RM35](#RM35%20—%20Révocation%20en%20cascade%20à%20la%20suppression%20d'un%20rôle) | Suppression d'un rôle → révocation automatique sur toutes les organisations | Administration réseau |
| [RM36](#RM36%20—%20Nom%20de%20rôle%20unique) | Nom de rôle unique dans le système (insensible à la casse) | Administration réseau |
| [RM37](#RM37%20—%20Multi-rôles%20par%20organisation) | Une organisation peut avoir plusieurs rôles — droits effectifs = union de tous ses rôles | Administration réseau |
| [RM29](#RM29%20—%20Taux%20de%20commission%20défini%20par%20le%20réseau) | Taux de commission défini par le réseau (défaut 10 %), uniforme pour tous les assets | PI |
| [RM30](#RM30%20—%20Calcul%20automatique%20du%20prix%20d'un%20module) | Prix module = somme des composants si non défini par l'auteur | PI |
| [RM31](#RM31%20—%20Modification%20de%20prix%20—%20effet%20sur%20les%20commandes%20futures%20uniquement) | Modification de prix applicable aux commandes futures uniquement | PI |
| [RM32](#RM32%20—%20Asset%20à%20prix%20nul%20—%20librement%20disponible) | Asset à prix nul → libre accès, aucune commission | PI |
| [RM33](#RM33%20—%20Devise%20unique%20par%20réseau) | Devise unique par réseau, définie par l'admin, non modifiable par l'auteur | PI |
