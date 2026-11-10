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
> [!check]- Tests — 1 explicite(s) · 0 déduit(s)
> - ✅ [TestIntegration_RM01_DuplicateHash_Rejected](../../docs/tests/adapters-out-fabric/TestIntegration_RM01_DuplicateHash_Rejected.md) — explicite (le test cite RM01)
<!-- tests-obsidian:end -->

### RM02 — Catégorie d'asset obligatoire

- **Déclencheur** : Création ou dérivation d'un asset
- **Condition** : Toujours
- **Conséquence** : Chaque asset doit appartenir à l'une des 8 catégories : `base`, `amélioration`, `variation`, `adaptation`, `dérivation`, `extension`, `régression`, `découpage`. La catégorie est donnée par le demandeur : le système n'en choisit aucune par défaut, et refuse une création sans catégorie ou avec une catégorie inconnue
- **Use cases** : UCCE01–06, UCAM05

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 2 explicite(s) · 0 déduit(s)
> - ✅ [TestComponents_POST_RM02_NoDefaultCategoryOrChannel](../../docs/tests/adapters-in-rest/TestComponents_POST_RM02_NoDefaultCategoryOrChannel.md) — explicite (le test cite RM02)
> - ✅ [TestAddFull_RM02_CategoryRequired_Rejected](../../docs/tests/domain-model/TestAddFull_RM02_CategoryRequired_Rejected.md) — explicite (le test cite RM02)
<!-- tests-obsidian:end -->

### RM03 — Compatibilité de licence

- **Déclencheur** : Soumission d'un asset avec `ParentID != ""` et `LicenseID != ""`
- **Condition** : Toujours
- **Conséquence** : La licence de l'asset dérivé doit être compatible avec celle du parent. Si incompatible → soumission rejetée
- **Use cases** : UCCE04, UCMOD06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 4 explicite(s) · 0 déduit(s)
> - ✅ [TestRunModelAdd_ServiceError_Rejected](../../docs/tests/adapters-in-cli/TestRunModelAdd_ServiceError_Rejected.md) — explicite (le test cite RM03)
> - ✅ [TestDomainErrors_SubmittedAndLicense_MappedToHTTP](../../docs/tests/adapters-in-rest/TestDomainErrors_SubmittedAndLicense_MappedToHTTP.md) — explicite (le test cite RM03)
> - ✅ [TestCreateModule_ParentID_RM03_LicenseIncompatible_Rejected](../../docs/tests/domain-model/TestCreateModule_ParentID_RM03_LicenseIncompatible_Rejected.md) — explicite (le test cite RM03)
> - ✅ [TestUpdateAsset_RM03_IncompatibleLicense_Rejected](../../docs/tests/domain-model/TestUpdateAsset_RM03_IncompatibleLicense_Rejected.md) — explicite (le test cite RM03)
<!-- tests-obsidian:end -->

### RM04 — UUID unique

- **Déclencheur** : Enregistrement blockchain réussi
- **Condition** : Toujours
- **Conséquence** : Un UUID unique est généré et attribué à l'asset. Il est immuable et ne peut pas être modifié ultérieurement
- **Use cases** : UCCE01–06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM05 — ParentID obligatoire pour les dérivés

- **Déclencheur** : Création d'un asset non-`base`
- **Condition** : Catégorie ≠ `base`
- **Conséquence** : L'asset dérivé doit référencer un `ParentID` valide. Un asset `base` n'a pas de parent
- **Use cases** : UCCE04, UCCE05

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 1 explicite(s) · 0 déduit(s)
> - ✅ [TestIntegration_RM05_ParentMustExist](../../docs/tests/adapters-out-fabric/TestIntegration_RM05_ParentMustExist.md) — explicite (le test cite RM05)
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
> [!check]- Tests — 1 explicite(s) · 0 déduit(s)
> - ✅ [TestIntegration_RM06_WriteOnce](../../docs/tests/adapters-out-fabric/TestIntegration_RM06_WriteOnce.md) — explicite (le test cite RM06)
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
> [!check]- Tests — 2 explicite(s) · 0 déduit(s)
> - ✅ [TestRemoveModule_Submitted_HidesLocally](../../docs/tests/domain-model/TestRemoveModule_Submitted_HidesLocally.md) — explicite (le test cite RM08)
> - ✅ [TestRemove_Submitted_HidesLocally](../../docs/tests/domain-model/TestRemove_Submitted_HidesLocally.md) — explicite (le test cite RM08)
<!-- tests-obsidian:end -->

## 3. Interfaces et liaisons

### RM09 — Interface à usage unique

- **Déclencheur** : Création d'une liaison entre deux interfaces
- **Condition** : L'interface source ou cible est déjà engagée dans une liaison sur l'instance désignée, ou la liaison relierait une interface à elle-même sur la même instance
- **Conséquence** : La liaison est refusée. L'usage unique porte sur l'interface d'une instance donnée : la même interface d'un composant reste disponible dans ses autres instances (RM15), et deux instances d'un même composant peuvent se relier par cette interface. Une liaison qui ne désigne pas d'instance engage l'interface elle-même
- **Use cases** : UCAM01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 8 explicite(s) · 0 déduit(s)
> - ✅ [TestRunModelLinkAdd_RM10_Refused_ErrorReturned](../../docs/tests/adapters-in-cli/TestRunModelLinkAdd_RM10_Refused_ErrorReturned.md) — explicite (le test cite RM09)
> - ✅ [TestDomainErrors_AssemblyRules_MappedToHTTP](../../docs/tests/adapters-in-rest/TestDomainErrors_AssemblyRules_MappedToHTTP.md) — explicite (le test cite RM09)
> - ✅ [TestAddAssemblyLink_RM09_InterfaceUsedOnce](../../docs/tests/domain-model/TestAddAssemblyLink_RM09_InterfaceUsedOnce.md) — explicite (le test cite RM09)
> - ✅ [TestAddAssemblyLink_RM09_NoInstance_InterfaceUsedOnce](../../docs/tests/domain-model/TestAddAssemblyLink_RM09_NoInstance_InterfaceUsedOnce.md) — explicite (le test cite RM09)
> - ✅ [TestAddAssemblyLink_RM09_SameInterfaceOtherInstance_OK](../../docs/tests/domain-model/TestAddAssemblyLink_RM09_SameInterfaceOtherInstance_OK.md) — explicite (le test cite RM09)
> - ✅ [TestAddAssemblyLink_RM09_SameInterfaceTwoInstances_OK](../../docs/tests/domain-model/TestAddAssemblyLink_RM09_SameInterfaceTwoInstances_OK.md) — explicite (le test cite RM09)
> - ✅ [TestAddAssemblyLink_RM09_SelfLink_Rejected](../../docs/tests/domain-model/TestAddAssemblyLink_RM09_SelfLink_Rejected.md) — explicite (le test cite RM09)
> - ✅ [TestConnectVirtualToPhysical_RM09_InterfaceInUse_NotMaterialized](../../docs/tests/domain-model/TestConnectVirtualToPhysical_RM09_InterfaceInUse_NotMaterialized.md) — explicite (le test cite RM09)
<!-- tests-obsidian:end -->

### RM10 — Vérification de compatibilité automatique

- **Déclencheur** : Création de toute liaison
- **Condition** : Toujours
- **Conséquence** : Le système appelle `ifacesCompatible` pour vérifier la cohérence de la paire avant d'enregistrer la liaison. Une liaison incompatible est refusée par le service. Avec un asset d'accroche, c'est l'accroche qui doit être compatible avec chacune des deux interfaces (UCAM07). L'incompatibilité post-création est couverte par RM12
- **Use cases** : UCAM01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 4 explicite(s) · 0 déduit(s)
> - ✅ [TestRunModelLinkAdd_RM10_Refused_ErrorReturned](../../docs/tests/adapters-in-cli/TestRunModelLinkAdd_RM10_Refused_ErrorReturned.md) — explicite (le test cite RM10)
> - ✅ [TestDomainErrors_AssemblyRules_MappedToHTTP](../../docs/tests/adapters-in-rest/TestDomainErrors_AssemblyRules_MappedToHTTP.md) — explicite (le test cite RM10)
> - ✅ [TestAddAssemblyLink_RM10_IncompatibleInterfaces_Rejected](../../docs/tests/domain-model/TestAddAssemblyLink_RM10_IncompatibleInterfaces_Rejected.md) — explicite (le test cite RM10)
> - ✅ [TestConnectVirtualToPhysical_RM09_InterfaceInUse_NotMaterialized](../../docs/tests/domain-model/TestConnectVirtualToPhysical_RM09_InterfaceInUse_NotMaterialized.md) — explicite (le test cite RM10)
<!-- tests-obsidian:end -->

### RM11 — Critères de compatibilité d'interfaces

- **Déclencheur** : Appel de `ifacesCompatible`
- **Condition** : Toujours
- **Conséquence** : Deux interfaces sont compatibles si : (1) même catégorie (Électrique, Mécanique, Numérique…), (2) même type (USB-C, UART…), y compris quand il n'est renseigné que d'un côté, (3) sens complémentaires (♂/♀ ou bidirectionnel), (4) valeurs compatibles : une interface porte zéro, une ou plusieurs contraintes de valeur, chacune avec sa valeur ou sa plage et son unité (ex. une alimentation : 5 V et 0 à 2 A). Chaque contrainte dont l'unité de base figure des deux côtés doit se chevaucher, préfixes convertis (1 V = 1 000 mV, RM53) ; une valeur simple vaut l'intervalle [v ; v] et 0 est une vraie valeur (0 V, 0 bar). Une contrainte présente d'un seul côté n'impose rien ; deux interfaces qui portent chacune des contraintes sans aucune unité de base commune sont incompatibles. Une interface sans contrainte n'en impose aucune
- **Use cases** : UCAM01, UCAM03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 2 explicite(s) · 0 déduit(s)
> - ✅ [TestAddAssemblyLink_RM10_IncompatibleInterfaces_Rejected](../../docs/tests/domain-model/TestAddAssemblyLink_RM10_IncompatibleInterfaces_Rejected.md) — explicite (le test cite RM11)
> - ✅ [TestAddAssemblyLink_RM11_ValueRanges](../../docs/tests/domain-model/TestAddAssemblyLink_RM11_ValueRanges.md) — explicite (le test cite RM11)
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
- **Conséquence** : Chaque asset dispose en permanence d'au moins un slot virtuel (`Virtual: true`). Dès qu'un slot virtuel est matérialisé en interface, un nouveau slot virtuel est recréé automatiquement. Le slot n'est matérialisé que si la liaison demandée est permise (RM09, RM10) : un refus le laisse virtuel
- **Use cases** : UCAM02, UCAM03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 1 explicite(s) · 0 déduit(s)
> - ✅ [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — explicite (le test cite RM13)
<!-- tests-obsidian:end -->

### RM50 — Liaison directe unique

- **Déclencheur** : Création d'une liaison directe entre deux assets, sans interface
- **Condition** : Une liaison directe relie déjà le même asset source au même asset cible
- **Conséquence** : La création est refusée. La liaison existante reste consultable dans la liste des liaisons : le système ne la renvoie pas à la place d'une nouvelle et n'en crée pas de doublon. Une liaison de sens inverse (cible vers source) est une autre liaison
- **Use cases** : UCAM01

---

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 3 explicite(s) · 0 déduit(s)
> - ✅ [TestConnections_GET_RM50_ListsAll](../../docs/tests/adapters-in-rest/TestConnections_GET_RM50_ListsAll.md) — explicite (le test cite RM50)
> - ✅ [TestDomainErrors_AssemblyRules_MappedToHTTP](../../docs/tests/adapters-in-rest/TestDomainErrors_AssemblyRules_MappedToHTTP.md) — explicite (le test cite RM50)
> - ✅ [TestAddConnection_RM50_Duplicate_Rejected](../../docs/tests/domain-model/TestAddConnection_RM50_Duplicate_Rejected.md) — explicite (le test cite RM50)
<!-- tests-obsidian:end -->

### RM53 — Unité d'une contrainte de valeur

- **Déclencheur** : Ajout ou modification d'une interface portant une contrainte de valeur
- **Condition** : Toujours
- **Conséquence** : L'unité est lue ainsi : (1) une unité connue telle quelle du vocabulaire des unités (ex. `min`, `Pa`, `mol`, `bar`) vaut facteur 1 ; (2) sinon, un préfixe SI usuel — `n`, `µ` ou `u`, `m`, `c`, `d`, `k`, `M`, `G`, la casse comptant (`m` milli, `M` méga) — suivi d'une unité connue applique son facteur (`mV` = 0,001 V, `kW` = 1 000 W, `mm` = 0,001 m) ; (3) sinon, l'unité est inconnue et l'interface est refusée — le système ne devine pas. Deux contraintes ne se comparent que sur la même unité de base : aucune conversion entre unités différentes (`bar` et `Pa`, `L` et `m³`, `°C` et `K`)
- **Use cases** : UCCE06, UCAM03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

## 4. Composition d'un Module (instances)

### RM14 — Suppression en cascade des connexions

- **Déclencheur** : Retrait d'une instance d'un module (`RemoveAssetFromWorkspace`)
- **Condition** : Toujours
- **Conséquence** : Toutes les connexions de l'instance retirée sont supprimées automatiquement. Les assets liés restent dans le module mais leurs interfaces concernées redeviennent libres
- **Use cases** : UCAM08

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 3 déduit(s)
> - 🟡 [TestModuleComposition_RM19_Submitted_Rejected](../../docs/tests/domain-model/TestModuleComposition_RM19_Submitted_Rejected.md) — déduit : teste `ModelService.RemoveAssetFromWorkspace`
> - 🟡 [TestWorkspace_AddRemoveInstance](../../docs/tests/domain-model/TestWorkspace_AddRemoveInstance.md) — déduit : teste `ModelService.RemoveAssetFromWorkspace`
> - 🟡 [TestWorkspace_NilStore_ReturnsError](../../docs/tests/domain-model/TestWorkspace_NilStore_ReturnsError.md) — déduit : teste `ModelService.RemoveAssetFromWorkspace`
<!-- tests-obsidian:end -->

### RM15 — Instance indépendante

- **Déclencheur** : Ajout d'un module déjà instancié dans le module hôte
- **Condition** : Module déjà instancié
- **Conséquence** : Une seconde instance est créée avec ses propres connexions, indépendantes de la première
- **Use cases** : UCMOD02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 1 explicite(s) · 0 déduit(s)
> - ✅ [TestRunModelInstanceRemove_NominalCase](../../docs/tests/adapters-in-cli/TestRunModelInstanceRemove_NominalCase.md) — explicite (le test cite RM15)
<!-- tests-obsidian:end -->

### RM51 — Assemblage interne au module

- **Déclencheur** : Rattachement d'une liaison à un asset (`myr model assembly add`, `POST /api/components/{id}/assemblies`)
- **Condition** : Toujours
- **Conséquence** : La liaison doit exister et relier deux instances de cet asset : par leurs identifiants d'instance, ou, pour une liaison directe qui n'en désigne pas, par les assets de ses instances. Sinon le rattachement est refusé et l'asset reste inchangé — un composant ne devient jamais module par un rattachement refusé
- **Use cases** : UCMOD01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 5 explicite(s) · 70 déduit(s)
> - ✅ [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — explicite (le test cite RM51)
> - ✅ [TestDomainErrors_AssemblyRules_MappedToHTTP](../../docs/tests/adapters-in-rest/TestDomainErrors_AssemblyRules_MappedToHTTP.md) — explicite (le test cite RM51)
> - ✅ [TestAddAssemblyToModule_RM51_Component_StaysComponent](../../docs/tests/domain-model/TestAddAssemblyToModule_RM51_Component_StaysComponent.md) — explicite (le test cite RM51)
> - ✅ [TestAddAssemblyToModule_RM51_ConnectionNotInModule_Rejected](../../docs/tests/domain-model/TestAddAssemblyToModule_RM51_ConnectionNotInModule_Rejected.md) — explicite (le test cite RM51)
> - ✅ [TestAddAssemblyToModule_RM51_InstanceLink_OK](../../docs/tests/domain-model/TestAddAssemblyToModule_RM51_InstanceLink_OK.md) — explicite (le test cite RM51)
> - 🟡 [TestRunModelAssemblyAdd_NominalCase](../../docs/tests/adapters-in-cli/TestRunModelAssemblyAdd_NominalCase.md) — déduit : teste commande `myr model assembly add`
> - 🟡 [TestRunModelAssemblyAdd_ServiceError_Rejected](../../docs/tests/adapters-in-cli/TestRunModelAssemblyAdd_ServiceError_Rejected.md) — déduit : teste commande `myr model assembly add`
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : teste route `/api/components`
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : teste route `/api/components`
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : teste route `/api/components`
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : teste route `/api/components`
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : teste route `/api/components`
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : teste route `/api/components`
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : teste route `/api/components`
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : teste route `/api/components`
> - … et 60 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM52 — Composition sans cycle

- **Déclencheur** : Ajout d'une instance à un module (`myr model instance add`, `POST /api/components/{id}/instances`)
- **Condition** : L'asset ajouté est le module lui-même, ou le contient, même indirectement
- **Conséquence** : L'ajout est refusé. La composition d'un module reste un arbre fini, que le système parcourt sans fin de récursion (interfaces exposées, arborescence)
- **Use cases** : UCMOD02

---

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 3 explicite(s) · 71 déduit(s)
> - ✅ [TestDomainErrors_AssemblyRules_MappedToHTTP](../../docs/tests/adapters-in-rest/TestDomainErrors_AssemblyRules_MappedToHTTP.md) — explicite (le test cite RM52)
> - ✅ [TestAddAssetToWorkspace_RM52_Cycle_Rejected](../../docs/tests/domain-model/TestAddAssetToWorkspace_RM52_Cycle_Rejected.md) — explicite (le test cite RM52)
> - ✅ [TestGetModuleInterfaces_RM52_LegacyCycle_Terminates](../../docs/tests/domain-model/TestGetModuleInterfaces_RM52_LegacyCycle_Terminates.md) — explicite (le test cite RM52)
> - 🟡 [TestRunModelInstanceAdd_NominalCase](../../docs/tests/adapters-in-cli/TestRunModelInstanceAdd_NominalCase.md) — déduit : teste commande `myr model instance add`
> - 🟡 [TestRunModelInstanceAdd_ServiceError_Rejected](../../docs/tests/adapters-in-cli/TestRunModelInstanceAdd_ServiceError_Rejected.md) — déduit : teste commande `myr model instance add`
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : teste route `/api/components`
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : teste route `/api/components`
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : teste route `/api/components`
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : teste route `/api/components`
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : teste route `/api/components`
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : teste route `/api/components`
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : teste route `/api/components`
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : teste route `/api/components`
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : teste route `/api/components`
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : teste route `/api/components`
> - … et 59 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

## 5. Modules et cycle de vie des assets

> RM16 et RM19 sont génériques à tout `Model3D` (composant ou module) — RM17 et RM18 restent spécifiques au module (une notion d'« assemblage » ou de « ModuleVersion » n'a pas de sens pour un composant simple). Voir `specs/3-Conception/Conception_intro.md` ADR-02.

### RM16 — État draft

- **Déclencheur** : Création d'un asset (composant ou module)
- **Condition** : Un module est **toujours** créé en `draft` (RM17 l'exige avant soumission). Un composant est créé **directement `submitted`** par défaut, sauf si `draft: true` est explicitement demandé (`myr model add --draft`)
- **Conséquence** : L'asset en `draft` n'est pas visible sur le réseau comme définitif. Ses interfaces et attributs mutables peuvent être modifiés librement (localement) tant qu'il n'est pas soumis (ADR-02)
- **Use cases** : UCMOD01, UCCE01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 1 explicite(s) · 3 déduit(s)
> - ✅ [TestIntegration_RM19_DraftRejected](../../docs/tests/adapters-out-fabric/TestIntegration_RM19_DraftRejected.md) — explicite (le test cite RM16)
> - 🟡 [TestResolveOwner](../../docs/tests/adapters-in-cli/TestResolveOwner.md) — déduit : teste commande `myr model add`
> - 🟡 [TestRunModelAdd_NominalCase](../../docs/tests/adapters-in-cli/TestRunModelAdd_NominalCase.md) — déduit : teste commande `myr model add`
> - 🟡 [TestRunModelAdd_ServiceError_Rejected](../../docs/tests/adapters-in-cli/TestRunModelAdd_ServiceError_Rejected.md) — déduit : teste commande `myr model add`
<!-- tests-obsidian:end -->

### RM17 — Assemblage requis pour soumission (module uniquement)

- **Déclencheur** : Appel de `SubmitModule`
- **Condition** : Toujours
- **Conséquence** : Le module doit contenir au moins une liaison entre composants. Si aucune liaison → soumission rejetée avec message explicite. Ne s'applique pas à un composant : sa soumission ne requiert aucun assemblage
- **Use cases** : UCMOD06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 1 explicite(s) · 9 déduit(s)
> - ✅ [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — explicite (le test cite RM17)
> - 🟡 [TestCreateModule_ParentID_CopiesComposition](../../docs/tests/domain-model/TestCreateModule_ParentID_CopiesComposition.md) — déduit : teste `ModelService.SubmitModule`
> - 🟡 [TestCreateModule_ParentID_RM03_LicenseIncompatible_Rejected](../../docs/tests/domain-model/TestCreateModule_ParentID_RM03_LicenseIncompatible_Rejected.md) — déduit : teste `ModelService.SubmitModule`
> - 🟡 [TestCreateModule_RM19_ForkOfSubmitted_Submittable](../../docs/tests/domain-model/TestCreateModule_RM19_ForkOfSubmitted_Submittable.md) — déduit : teste `ModelService.SubmitModule`
> - 🟡 [TestModule_NotFound_ReturnsError](../../docs/tests/domain-model/TestModule_NotFound_ReturnsError.md) — déduit : teste `ModelService.SubmitModule`
> - 🟡 [TestModule_Submit](../../docs/tests/domain-model/TestModule_Submit.md) — déduit : teste `ModelService.SubmitModule`
> - 🟡 [TestModule_Submit_NoAssembly](../../docs/tests/domain-model/TestModule_Submit_NoAssembly.md) — déduit : teste `ModelService.SubmitModule`
> - 🟡 [TestRemoveModule_Submitted_HidesLocally](../../docs/tests/domain-model/TestRemoveModule_Submitted_HidesLocally.md) — déduit : teste `ModelService.SubmitModule`
> - 🟡 [TestService_NilBlockchain_ErrBlockchainUnavailable](../../docs/tests/domain-model/TestService_NilBlockchain_ErrBlockchainUnavailable.md) — déduit : teste `ModelService.SubmitModule`
> - 🟡 [TestSubmitModule_RM19_AlreadySubmitted_Rejected](../../docs/tests/domain-model/TestSubmitModule_RM19_AlreadySubmitted_Rejected.md) — déduit : teste `ModelService.SubmitModule`
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
- **Conséquence** : Un asset soumis est immuable sur le ledger : il n'est jamais re-soumis, et ni sa licence, ni ses interfaces, ni sa composition (instances, liaisons) ne changent. Seules ses métadonnées descriptives — nom, description, tags, liens — restent modifiables, par son seul propriétaire. Toute autre évolution, y compris une nouvelle version d'un module, passe par un nouvel asset dérivé (fork), créé en état `draft` avec `ParentID` référençant l'asset d'origine
- **Use cases** : UCMOD06, UCCE06, UCCE02, UCMOD03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 12 explicite(s) · 0 déduit(s)
> - ✅ [TestDomainErrors_SubmittedAndLicense_MappedToHTTP](../../docs/tests/adapters-in-rest/TestDomainErrors_SubmittedAndLicense_MappedToHTTP.md) — explicite (le test cite RM19)
> - ✅ [TestIntegration_ModelService_SubmittedAsset](../../docs/tests/adapters-out-fabric/TestIntegration_ModelService_SubmittedAsset.md) — explicite (le test cite RM19)
> - ✅ [TestIntegration_RM06_WriteOnce](../../docs/tests/adapters-out-fabric/TestIntegration_RM06_WriteOnce.md) — explicite (le test cite RM19)
> - ✅ [TestIntegration_RM19_DraftRejected](../../docs/tests/adapters-out-fabric/TestIntegration_RM19_DraftRejected.md) — explicite (le test cite RM19)
> - ✅ [TestIntegration_UpdateMetadata_OwnerOnly](../../docs/tests/adapters-out-fabric/TestIntegration_UpdateMetadata_OwnerOnly.md) — explicite (le test cite RM19)
> - ✅ [TestAddFull_RM19_DerivedFromSubmittedModule_CopiesComposition](../../docs/tests/domain-model/TestAddFull_RM19_DerivedFromSubmittedModule_CopiesComposition.md) — explicite (le test cite RM19)
> - ✅ [TestConnectVirtualToPhysical_RM19_SubmittedCarrier_Rejected](../../docs/tests/domain-model/TestConnectVirtualToPhysical_RM19_SubmittedCarrier_Rejected.md) — explicite (le test cite RM19)
> - ✅ [TestCreateModule_RM19_ForkOfSubmitted_Submittable](../../docs/tests/domain-model/TestCreateModule_RM19_ForkOfSubmitted_Submittable.md) — explicite (le test cite RM19)
> - ✅ [TestModuleComposition_RM19_Submitted_Rejected](../../docs/tests/domain-model/TestModuleComposition_RM19_Submitted_Rejected.md) — explicite (le test cite RM19)
> - ✅ [TestSubmitModule_RM19_AlreadySubmitted_Rejected](../../docs/tests/domain-model/TestSubmitModule_RM19_AlreadySubmitted_Rejected.md) — explicite (le test cite RM19)
> - ✅ [TestSubmit_RM19_AlreadySubmitted_Rejected](../../docs/tests/domain-model/TestSubmit_RM19_AlreadySubmitted_Rejected.md) — explicite (le test cite RM19)
> - ✅ [TestUpdateAsset_RM19_SubmittedLicense_Rejected](../../docs/tests/domain-model/TestUpdateAsset_RM19_SubmittedLicense_Rejected.md) — explicite (le test cite RM19)
<!-- tests-obsidian:end -->

### RM54 — Seul le propriétaire modifie son asset

- **Déclencheur** : Modification, suppression ou soumission d'un asset — métadonnées, licence, interfaces, instances, liaisons rattachées
- **Condition** : L'identité qui agit n'est pas le propriétaire de l'asset, ou aucune identité n'agit
- **Conséquence** : L'opération est refusée, administrateur compris. Le propriétaire d'un asset est l'identité qui le crée (`pseudo@org`), et non celle que la demande désignerait ; sans identité, rien n'est créé ni modifié
- **Use cases** : UCCE02, UCCE06, UCCE07, UCMOD02, UCMOD03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 5 explicite(s) · 0 déduit(s)
> - ✅ [TestDomainErrors_RM54_NotOwner_Forbidden](../../docs/tests/adapters-in-rest/TestDomainErrors_RM54_NotOwner_Forbidden.md) — explicite (le test cite RM54)
> - ✅ [TestOwnership_RM54_CreatorIsOwner](../../docs/tests/domain-model/TestOwnership_RM54_CreatorIsOwner.md) — explicite (le test cite RM54)
> - ✅ [TestOwnership_RM54_NoIdentity_Rejected](../../docs/tests/domain-model/TestOwnership_RM54_NoIdentity_Rejected.md) — explicite (le test cite RM54)
> - ✅ [TestOwnership_RM54_OtherIdentity_Rejected](../../docs/tests/domain-model/TestOwnership_RM54_OtherIdentity_Rejected.md) — explicite (le test cite RM54)
> - ✅ [TestOwnership_RM54_SubmittedMetadata_OtherIdentity_Rejected](../../docs/tests/domain-model/TestOwnership_RM54_SubmittedMetadata_OtherIdentity_Rejected.md) — explicite (le test cite RM54)
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
> [!check]- Tests — 0 explicite(s) · 10 déduit(s)
> - 🟡 [TestChannels_PUT_SwitchesChannel](../../docs/tests/adapters-in-rest/TestChannels_PUT_SwitchesChannel.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_CannotSignIn_Forbidden](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_CannotSignIn_Forbidden.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_EnrollError](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_EnrollError.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_InvalidJSON](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_InvalidJSON.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MethodNotAllowed.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_MissingFields](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MissingFields.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_NoIdentityService](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_NoIdentityService.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_RoleFromCertificate](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_RoleFromCertificate.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_Success](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_Success.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestRequireAuth_BlockchainWithValidToken](../../docs/tests/adapters-in-rest/TestRequireAuth_BlockchainWithValidToken.md) — déduit : teste route `/api/identity/session`
<!-- tests-obsidian:end -->

### RM21 — Rôle Lecteur par défaut à l'auto-enregistrement

- **Déclencheur** : `POST /api/identity/request` avec auto-enregistrement réseau actif (`AllowAutoRegister`)
- **Condition** : Le profil réseau ne définit pas explicitement un autre rôle (`AutoRegisterRole`)
- **Conséquence** : Une identité auto-enregistrée reçoit par défaut le rôle **Lecteur** (`reader`), sauf si le profil réseau configure explicitement un autre rôle initial
- **Use cases** : UCA01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 2 explicite(s) · 6 déduit(s)
> - ✅ [TestHandleIdentitySession_ReaderCannotWrite](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_ReaderCannotWrite.md) — explicite (le test cite RM21)
> - ✅ [TestCreateAccount_RM46_CustomAndDefaultRole](../../docs/tests/domain-identity/TestCreateAccount_RM46_CustomAndDefaultRole.md) — explicite (le test cite RM21)
> - 🟡 [TestHandleIdentityRequest_InvalidJSON](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_InvalidJSON.md) — déduit : teste route `/api/identity/request`
> - 🟡 [TestHandleIdentityRequest_InvalidName_Rejected](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_InvalidName_Rejected.md) — déduit : teste route `/api/identity/request`
> - 🟡 [TestHandleIdentityRequest_MissingFields](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_MissingFields.md) — déduit : teste route `/api/identity/request`
> - 🟡 [TestHandleIdentityRequest_RegistrationClosed_Rejected](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_RegistrationClosed_Rejected.md) — déduit : teste route `/api/identity/request`
> - 🟡 [TestHandleIdentityRequest_ServiceError](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_ServiceError.md) — déduit : teste route `/api/identity/request`
> - 🟡 [TestHandleIdentityRequest_Success](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_Success.md) — déduit : teste route `/api/identity/request`
<!-- tests-obsidian:end -->

### RM22 — Changement de rôle réservé à l'administrateur

- **Déclencheur** : `myr identity set-role`
- **Condition** : Toujours
- **Conséquence** : Seul l'administrateur peut modifier le rôle (`Myr.role`) d'une identité existante auprès de la CA. Le nouveau rôle ne s'applique qu'au prochain ré-enrôlement de l'identité — ce n'est pas immédiat
- **Use cases** : UCA08

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 1 déduit(s)
> - 🟡 [TestRunIdentitySetRole_NominalCase](../../docs/tests/adapters-in-cli/TestRunIdentitySetRole_NominalCase.md) — déduit : teste commande `myr identity set-role`
<!-- tests-obsidian:end -->

### RM43 — Inscription selon les règles du réseau

- **Déclencheur** : Demande d'accès (`POST /api/identity/request`, `myr identity request`)
- **Condition** : Toujours
- **Conséquence** : Les règles d'inscription sont propres à chaque réseau et fixées par son administrateur. Sans auto-enregistrement, aucune demande n'est acceptée et rien n'est conservé : seul l'administrateur crée des comptes (`myr identity register`, `POST /api/identity/accounts`), actifs d'emblée avec le rôle qu'il choisit. Avec auto-enregistrement, l'identité est créée dès la demande avec le rôle par défaut du réseau (RM21) : active si le réseau n'exige pas d'approbation, en attente sinon
- **Use cases** : UCA01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 5 explicite(s) · 12 déduit(s)
> - ✅ [TestHandleIdentityRequest_RegistrationClosed_Rejected](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_RegistrationClosed_Rejected.md) — explicite (le test cite RM43)
> - ✅ [TestAccountRequest_FullFlow](../../docs/tests/domain-identity/TestAccountRequest_FullFlow.md) — explicite (le test cite RM43)
> - ✅ [TestRequestAccess_NoApproval_Active](../../docs/tests/domain-identity/TestRequestAccess_NoApproval_Active.md) — explicite (le test cite RM43)
> - ✅ [TestRequestAccess_RegistrationClosed_Rejected](../../docs/tests/domain-identity/TestRequestAccess_RegistrationClosed_Rejected.md) — explicite (le test cite RM43)
> - ✅ [TestRequestAccess_RequireApproval_Pending](../../docs/tests/domain-identity/TestRequestAccess_RequireApproval_Pending.md) — explicite (le test cite RM43)
> - 🟡 [TestRunIdentityRegister_NominalCase](../../docs/tests/adapters-in-cli/TestRunIdentityRegister_NominalCase.md) — déduit : teste commande `myr identity register`
> - 🟡 [TestRunIdentityRequest_AutoRegister_Approved](../../docs/tests/adapters-in-cli/TestRunIdentityRequest_AutoRegister_Approved.md) — déduit : teste commande `myr identity request`
> - 🟡 [TestRunIdentityRequest_NoNetwork_Rejected](../../docs/tests/adapters-in-cli/TestRunIdentityRequest_NoNetwork_Rejected.md) — déduit : teste commande `myr identity request`
> - 🟡 [TestRunIdentityRequest_RequireApproval_Pending](../../docs/tests/adapters-in-cli/TestRunIdentityRequest_RequireApproval_Pending.md) — déduit : teste commande `myr identity request`
> - 🟡 [TestHandleIdentityAccounts_Admin_Created](../../docs/tests/adapters-in-rest/TestHandleIdentityAccounts_Admin_Created.md) — déduit : teste route `/api/identity/accounts`
> - 🟡 [TestHandleIdentityAccounts_InvalidNameOrRole_Rejected](../../docs/tests/adapters-in-rest/TestHandleIdentityAccounts_InvalidNameOrRole_Rejected.md) — déduit : teste route `/api/identity/accounts`
> - 🟡 [TestHandleIdentityAccounts_NonAdmin_Forbidden](../../docs/tests/adapters-in-rest/TestHandleIdentityAccounts_NonAdmin_Forbidden.md) — déduit : teste route `/api/identity/accounts`
> - 🟡 [TestHandleIdentityRequest_InvalidJSON](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_InvalidJSON.md) — déduit : teste route `/api/identity/request`
> - 🟡 [TestHandleIdentityRequest_InvalidName_Rejected](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_InvalidName_Rejected.md) — déduit : teste route `/api/identity/request`
> - 🟡 [TestHandleIdentityRequest_MissingFields](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_MissingFields.md) — déduit : teste route `/api/identity/request`
> - … et 2 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM44 — Décision sur une demande en attente

- **Déclencheur** : `myr identity approve|reject`, `POST /api/identity/requests/{id}/approve|reject`
- **Condition** : La demande est en attente ; l'appelant dispose de la permission d'administration des identités
- **Conséquence** : Approuver active l'identité ; refuser la suspend — rien n'est supprimé, ni l'identité ni la demande. Une demande déjà décidée ne peut plus l'être. La décision prend effet au prochain enrôlement de l'identité, c'est-à-dire à sa prochaine connexion
- **Use cases** : UCA01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 8 explicite(s) · 8 déduit(s)
> - ✅ [TestHandleIdentityRequestDecision_Approve_OK](../../docs/tests/adapters-in-rest/TestHandleIdentityRequestDecision_Approve_OK.md) — explicite (le test cite RM44)
> - ✅ [TestHandleIdentityRequestDecision_Reject_OK](../../docs/tests/adapters-in-rest/TestHandleIdentityRequestDecision_Reject_OK.md) — explicite (le test cite RM44)
> - ✅ [TestAccountRequest_FullFlow](../../docs/tests/domain-identity/TestAccountRequest_FullFlow.md) — explicite (le test cite RM44)
> - ✅ [TestApproveRequest_NominalCase](../../docs/tests/domain-identity/TestApproveRequest_NominalCase.md) — explicite (le test cite RM44)
> - ✅ [TestDecideRequest_NotFoundOrNotPending_Rejected](../../docs/tests/domain-identity/TestDecideRequest_NotFoundOrNotPending_Rejected.md) — explicite (le test cite RM44)
> - ✅ [TestRejectRequest_NominalCase](../../docs/tests/domain-identity/TestRejectRequest_NominalCase.md) — explicite (le test cite RM44)
> - ✅ [TestRejectRequest_RM47_RoleRequest_IdentityUnchanged](../../docs/tests/domain-identity/TestRejectRequest_RM47_RoleRequest_IdentityUnchanged.md) — explicite (le test cite RM44)
> - ✅ [TestRequestRole_RM47_Approval_PendingThenGranted](../../docs/tests/domain-identity/TestRequestRole_RM47_Approval_PendingThenGranted.md) — explicite (le test cite RM44)
> - 🟡 [TestRunIdentityDecision_Approve_OK](../../docs/tests/adapters-in-cli/TestRunIdentityDecision_Approve_OK.md) — déduit : teste commande `myr identity approve`
> - 🟡 [TestRunIdentityDecision_NotPending_Rejected](../../docs/tests/adapters-in-cli/TestRunIdentityDecision_NotPending_Rejected.md) — déduit : teste commande `myr identity approve`
> - 🟡 [TestHandleIdentityRequestDecision_Errors](../../docs/tests/adapters-in-rest/TestHandleIdentityRequestDecision_Errors.md) — déduit : teste route `/api/identity/requests`
> - 🟡 [TestHandleIdentityRequestDecision_NonAdmin_Forbidden](../../docs/tests/adapters-in-rest/TestHandleIdentityRequestDecision_NonAdmin_Forbidden.md) — déduit : teste route `/api/identity/requests`
> - 🟡 [TestHandleIdentityRequests_Admin_EmptyList](../../docs/tests/adapters-in-rest/TestHandleIdentityRequests_Admin_EmptyList.md) — déduit : teste route `/api/identity/requests`
> - 🟡 [TestHandleIdentityRequests_Admin_ReturnsList](../../docs/tests/adapters-in-rest/TestHandleIdentityRequests_Admin_ReturnsList.md) — déduit : teste route `/api/identity/requests`
> - 🟡 [TestHandleIdentityRequests_NonAdmin_Forbidden](../../docs/tests/adapters-in-rest/TestHandleIdentityRequests_NonAdmin_Forbidden.md) — déduit : teste route `/api/identity/requests`
> - … et 1 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM45 — Connexion réservée aux identités actives dotées d'un rôle

- **Déclencheur** : Connexion (`POST /api/identity/session`)
- **Condition** : Toujours
- **Conséquence** : Le statut et le rôle de l'identité sont lus dans le certificat émis à la connexion, jamais fournis par le client. Une identité en attente, suspendue, sans statut reconnu ou sans rôle n'ouvre aucune session. La session porte le rôle du certificat
- **Use cases** : UCA02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 4 explicite(s) · 8 déduit(s)
> - ✅ [TestHandleIdentitySession_CannotSignIn_Forbidden](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_CannotSignIn_Forbidden.md) — explicite (le test cite RM45)
> - ✅ [TestHandleIdentitySession_RoleFromCertificate](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_RoleFromCertificate.md) — explicite (le test cite RM45)
> - ✅ [TestAccountRequest_FullFlow](../../docs/tests/domain-identity/TestAccountRequest_FullFlow.md) — explicite (le test cite RM45)
> - ✅ [TestWalletEntry_CanSignIn](../../docs/tests/domain-identity/TestWalletEntry_CanSignIn.md) — explicite (le test cite RM45)
> - 🟡 [TestChannels_PUT_SwitchesChannel](../../docs/tests/adapters-in-rest/TestChannels_PUT_SwitchesChannel.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_EnrollError](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_EnrollError.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_InvalidJSON](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_InvalidJSON.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MethodNotAllowed.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_MissingFields](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MissingFields.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_NoIdentityService](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_NoIdentityService.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestHandleIdentitySession_Success](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_Success.md) — déduit : teste route `/api/identity/session`
> - 🟡 [TestRequireAuth_BlockchainWithValidToken](../../docs/tests/adapters-in-rest/TestRequireAuth_BlockchainWithValidToken.md) — déduit : teste route `/api/identity/session`
<!-- tests-obsidian:end -->

### RM46 — Rôle attribué toujours défini dans la table des rôles

- **Déclencheur** : Attribution d'un rôle à une identité — création de compte (`myr identity register`), inscription (rôle par défaut `--auto-role` du réseau), changement de rôle (`myr identity set-role`), demande de rôle et son approbation (`myr identity request-role`, `myr identity approve`)
- **Condition** : Toujours
- **Conséquence** : Le rôle attribué doit exister dans la table des rôles du réseau ; un rôle inconnu est refusé, jamais remplacé par un autre — y compris à l'approbation d'une demande de rôle dont le rôle a été retiré de la table depuis. Le rôle attribué par défaut aux inscriptions ne peut porter aucune permission d'administration (`admin`, `identity.admin`, `role.admin`, `network.admin`) : une inscription automatique ne donne jamais de droit d'administration — ce contrôle s'applique dès la configuration du réseau
- **Use cases** : UCA01, UCA08

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 11 explicite(s) · 6 déduit(s)
> - ✅ [TestHandleIdentityAccounts_InvalidNameOrRole_Rejected](../../docs/tests/adapters-in-rest/TestHandleIdentityAccounts_InvalidNameOrRole_Rejected.md) — explicite (le test cite RM46)
> - ✅ [TestApproveRequest_RM46_RoleRemovedSinceRequest_Rejected](../../docs/tests/domain-identity/TestApproveRequest_RM46_RoleRemovedSinceRequest_Rejected.md) — explicite (le test cite RM46)
> - ✅ [TestCreateAccount_RM46_CustomAndDefaultRole](../../docs/tests/domain-identity/TestCreateAccount_RM46_CustomAndDefaultRole.md) — explicite (le test cite RM46)
> - ✅ [TestCreateAccount_RM46_UnknownRole_Rejected](../../docs/tests/domain-identity/TestCreateAccount_RM46_UnknownRole_Rejected.md) — explicite (le test cite RM46)
> - ✅ [TestRequestAccess_RM46_InvalidDefaultRole_Rejected](../../docs/tests/domain-identity/TestRequestAccess_RM46_InvalidDefaultRole_Rejected.md) — explicite (le test cite RM46)
> - ✅ [TestRequestRole_RM47_Invalid_Rejected](../../docs/tests/domain-identity/TestRequestRole_RM47_Invalid_Rejected.md) — explicite (le test cite RM46)
> - ✅ [TestSetRole_RM46_UnknownRole_Rejected](../../docs/tests/domain-identity/TestSetRole_RM46_UnknownRole_Rejected.md) — explicite (le test cite RM46)
> - ✅ [TestAutoRegisterRole_RM46_Invalid_Rejected](../../docs/tests/domain-network/TestAutoRegisterRole_RM46_Invalid_Rejected.md) — explicite (le test cite RM46)
> - ✅ [TestAssignableByDefault_RM46_AdminPermission_Rejected](../../docs/tests/domain-role/TestAssignableByDefault_RM46_AdminPermission_Rejected.md) — explicite (le test cite RM46)
> - ✅ [TestAssignable_RM46_KnownRole_OK](../../docs/tests/domain-role/TestAssignable_RM46_KnownRole_OK.md) — explicite (le test cite RM46)
> - ✅ [TestAssignable_RM46_UnknownRole_Rejected](../../docs/tests/domain-role/TestAssignable_RM46_UnknownRole_Rejected.md) — explicite (le test cite RM46)
> - 🟡 [TestRunIdentityDecision_Approve_OK](../../docs/tests/adapters-in-cli/TestRunIdentityDecision_Approve_OK.md) — déduit : teste commande `myr identity approve`
> - 🟡 [TestRunIdentityDecision_NotPending_Rejected](../../docs/tests/adapters-in-cli/TestRunIdentityDecision_NotPending_Rejected.md) — déduit : teste commande `myr identity approve`
> - 🟡 [TestRunIdentityRegister_NominalCase](../../docs/tests/adapters-in-cli/TestRunIdentityRegister_NominalCase.md) — déduit : teste commande `myr identity register`
> - 🟡 [TestRunIdentityRequestRole_RM47_GrantedOrPending](../../docs/tests/adapters-in-cli/TestRunIdentityRequestRole_RM47_GrantedOrPending.md) — déduit : teste commande `myr identity request-role`
> - … et 2 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### RM47 — Demande de rôle, attribution automatique ou approbation

- **Déclencheur** : Demande d'un rôle par une identité inscrite (`myr identity request-role`, `POST /api/identity/roles/request`)
- **Condition** : L'identité est active ; le rôle demandé figure dans la table des rôles (RM46) et diffère de son rôle actuel ; aucune autre demande de rôle de l'identité n'est en attente
- **Conséquence** : Chaque rôle de la table porte un réglage d'attribution. En attribution automatique, le rôle est attribué aussitôt ; sinon, la demande attend la décision de l'administrateur : approuver attribue le rôle, refuser laisse l'identité inchangée (RM44). Un rôle portant une permission d'administration n'est jamais attribué automatiquement, quel que soit son réglage. Le nouveau rôle prend effet au prochain enrôlement de l'identité, c'est-à-dire à sa prochaine connexion
- **Use cases** : UCA08

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 11 explicite(s) · 2 déduit(s)
> - ✅ [TestRunIdentityRequestRole_RM47_GrantedOrPending](../../docs/tests/adapters-in-cli/TestRunIdentityRequestRole_RM47_GrantedOrPending.md) — explicite (le test cite RM47)
> - ✅ [TestHandleIdentityRoleRequest_RM47_GrantedOrPending](../../docs/tests/adapters-in-rest/TestHandleIdentityRoleRequest_RM47_GrantedOrPending.md) — explicite (le test cite RM47)
> - ✅ [TestIntegration_RoleService_PublishedAndApplied](../../docs/tests/adapters-out-fabric/TestIntegration_RoleService_PublishedAndApplied.md) — explicite (le test cite RM47)
> - ✅ [TestRejectRequest_RM47_RoleRequest_IdentityUnchanged](../../docs/tests/domain-identity/TestRejectRequest_RM47_RoleRequest_IdentityUnchanged.md) — explicite (le test cite RM47)
> - ✅ [TestRequestRole_RM47_AdminRole_NeverAutomatic](../../docs/tests/domain-identity/TestRequestRole_RM47_AdminRole_NeverAutomatic.md) — explicite (le test cite RM47)
> - ✅ [TestRequestRole_RM47_Approval_PendingThenGranted](../../docs/tests/domain-identity/TestRequestRole_RM47_Approval_PendingThenGranted.md) — explicite (le test cite RM47)
> - ✅ [TestRequestRole_RM47_AutoGrant_Granted](../../docs/tests/domain-identity/TestRequestRole_RM47_AutoGrant_Granted.md) — explicite (le test cite RM47)
> - ✅ [TestRequestRole_RM47_Invalid_Rejected](../../docs/tests/domain-identity/TestRequestRole_RM47_Invalid_Rejected.md) — explicite (le test cite RM47)
> - ✅ [TestCreate_RM47_AutoGrantAdminRole_Rejected](../../docs/tests/domain-role/TestCreate_RM47_AutoGrantAdminRole_Rejected.md) — explicite (le test cite RM47)
> - ✅ [TestGrantedOnRequest_RM47](../../docs/tests/domain-role/TestGrantedOnRequest_RM47.md) — explicite (le test cite RM47)
> - ✅ [TestUpdate_RM47_AutoGrantWithAdminPermission_Rejected](../../docs/tests/domain-role/TestUpdate_RM47_AutoGrantWithAdminPermission_Rejected.md) — explicite (le test cite RM47)
> - 🟡 [TestSessionIdentity](../../docs/tests/adapters-in-cli/TestSessionIdentity.md) — déduit : teste commande `myr identity request-role`
> - 🟡 [TestHandleIdentityRoleRequest_Errors](../../docs/tests/adapters-in-rest/TestHandleIdentityRoleRequest_Errors.md) — déduit : teste route `/api/identity/roles/request`
<!-- tests-obsidian:end -->

### RM48 — Table des rôles partagée par le réseau, rôle administrateur immuable

- **Déclencheur** : Création, modification ou suppression d'un rôle (`myr role create|update|delete`) ; tout contrôle d'accès
- **Condition** : Toujours
- **Conséquence** : La table des rôles et de leurs permissions est publiée sur le ledger, un exemplaire par canal du réseau, et appliquée par le chaincode comme par chaque serveur. Seule une identité dont le rôle porte `role.admin` la modifie. Tous les rôles se modifient ou se suppriment, sauf `admin`, qui porte toujours toutes les permissions. Tant qu'aucune table n'est publiée, la table initiale s'applique (`reader`, `contributor`, `auditor`, `admin`) ; si la table devient illisible, un serveur applique la dernière table lue, à défaut la table initiale. Une identité dont le rôle a été supprimé n'a plus aucune permission. Les droits effectifs d'une identité sont l'union des permissions de son rôle (certificat) et de celles des rôles attribués à son organisation (RM37, UCADM06)
- **Use cases** : UCA05, UCADM06, UCADM07

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 14 explicite(s) · 0 déduit(s)
> - ✅ [TestFabricRoleStore_RM48_PublishedOnEveryChannel](../../docs/tests/adapters-out-fabric/TestFabricRoleStore_RM48_PublishedOnEveryChannel.md) — explicite (le test cite RM48)
> - ✅ [TestIntegration_PendingOrReader_Rejected](../../docs/tests/adapters-out-fabric/TestIntegration_PendingOrReader_Rejected.md) — explicite (le test cite RM48)
> - ✅ [TestIntegration_RM48_RolePublishedOnEveryChannel](../../docs/tests/adapters-out-fabric/TestIntegration_RM48_RolePublishedOnEveryChannel.md) — explicite (le test cite RM48)
> - ✅ [TestIntegration_RM48_RoleTableGuards](../../docs/tests/adapters-out-fabric/TestIntegration_RM48_RoleTableGuards.md) — explicite (le test cite RM48)
> - ✅ [TestIntegration_RM48_WritePermissionFollowsTable](../../docs/tests/adapters-out-fabric/TestIntegration_RM48_WritePermissionFollowsTable.md) — explicite (le test cite RM48)
> - ✅ [TestIntegration_RoleService_PublishedAndApplied](../../docs/tests/adapters-out-fabric/TestIntegration_RoleService_PublishedAndApplied.md) — explicite (le test cite RM48)
> - ✅ [TestDelete_RM48_Admin_Rejected](../../docs/tests/domain-role/TestDelete_RM48_Admin_Rejected.md) — explicite (le test cite RM48)
> - ✅ [TestDelete_RM48_DefaultRoleDeletable](../../docs/tests/domain-role/TestDelete_RM48_DefaultRoleDeletable.md) — explicite (le test cite RM48)
> - ✅ [TestTable_RM48_AdminAlwaysImmutable](../../docs/tests/domain-role/TestTable_RM48_AdminAlwaysImmutable.md) — explicite (le test cite RM48)
> - ✅ [TestTable_RM48_NothingPublished_DefaultRoles](../../docs/tests/domain-role/TestTable_RM48_NothingPublished_DefaultRoles.md) — explicite (le test cite RM48)
> - ✅ [TestTable_RM48_Unreachable_LastKnownTable](../../docs/tests/domain-role/TestTable_RM48_Unreachable_LastKnownTable.md) — explicite (le test cite RM48)
> - ✅ [TestTable_RM48_Unreachable_NeverRead_DefaultRoles](../../docs/tests/domain-role/TestTable_RM48_Unreachable_NeverRead_DefaultRoles.md) — explicite (le test cite RM48)
> - ✅ [TestUpdate_RM48_Admin_Rejected](../../docs/tests/domain-role/TestUpdate_RM48_Admin_Rejected.md) — explicite (le test cite RM48)
> - ✅ [TestUpdate_RM48_DefaultRoleEditable](../../docs/tests/domain-role/TestUpdate_RM48_DefaultRoleEditable.md) — explicite (le test cite RM48)
<!-- tests-obsidian:end -->

### RM49 — Lecture publique du catalogue, ouverte par l'administrateur

- **Déclencheur** : Lecture du catalogue sans session — liste et détail des composants et modules (`GET /api/components`, `/api/components/{id}` et ses sous-ressources), référentiels d'interfaces (`/api/refs`), catalogue des licences (`/api/licenses`)
- **Condition** : L'administrateur a ouvert la lecture publique du réseau (`myr network add|update --auto-guest`) ; elle est fermée par défaut
- **Conséquence** : Un visiteur sans compte lit le catalogue publié — les assets soumis et leurs sous-ressources — sans session ni jeton. Les brouillons lui restent inconnus (404) et il n'écrit jamais. Réseau fermé : toute lecture exige une session (401). Un jeton présent mais invalide est refusé, même lecture ouverte, pour que le client apprenne l'expiration de sa session. Aucune session invitée n'est créée
- **Use cases** : UCCL01, UCA02

---

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 3 explicite(s) · 84 déduit(s)
> - ✅ [TestPublicRead_RM49_Closed_SessionRequired](../../docs/tests/adapters-in-rest/TestPublicRead_RM49_Closed_SessionRequired.md) — explicite (le test cite RM49)
> - ✅ [TestPublicRead_RM49_Open_PublishedOnly](../../docs/tests/adapters-in-rest/TestPublicRead_RM49_Open_PublishedOnly.md) — explicite (le test cite RM49)
> - ✅ [TestPublicRead_RM49_WritesAndOtherRoutes_Refused](../../docs/tests/adapters-in-rest/TestPublicRead_RM49_WritesAndOtherRoutes_Refused.md) — explicite (le test cite RM49)
> - 🟡 [TestRunNetworkAdd_InvalidPeer_Rejected](../../docs/tests/adapters-in-cli/TestRunNetworkAdd_InvalidPeer_Rejected.md) — déduit : teste commande `myr network add`
> - 🟡 [TestRunNetworkAdd_NoApproval_NotSet](../../docs/tests/adapters-in-cli/TestRunNetworkAdd_NoApproval_NotSet.md) — déduit : teste commande `myr network add`
> - 🟡 [TestRunNetworkAdd_NominalCase](../../docs/tests/adapters-in-cli/TestRunNetworkAdd_NominalCase.md) — déduit : teste commande `myr network add`
> - 🟡 [TestRunNetworkAdd_RequireApproval](../../docs/tests/adapters-in-cli/TestRunNetworkAdd_RequireApproval.md) — déduit : teste commande `myr network add`
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : teste route `/api/components`
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : teste route `/api/components`
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : teste route `/api/components`
> - 🟡 [TestAPIIntegration_LicenseCatalog](../../docs/tests/adapters-in-rest/TestAPIIntegration_LicenseCatalog.md) — déduit : teste route `/api/licenses`
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : teste route `/api/components`
> - 🟡 [TestAPIIntegration_RefsCategoryPersists](../../docs/tests/adapters-in-rest/TestAPIIntegration_RefsCategoryPersists.md) — déduit : teste route `/api/refs`
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : teste route `/api/components`
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : teste route `/api/components`
> - … et 72 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
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
> [!check]- Tests — 0 explicite(s) · 4 déduit(s)
> - 🟡 [TestRunNetworkDestroy_NominalCase](../../docs/tests/adapters-in-cli/TestRunNetworkDestroy_NominalCase.md) — déduit : teste commande `myr network destroy`
> - 🟡 [TestRunNetworkDestroy_NotFound](../../docs/tests/adapters-in-cli/TestRunNetworkDestroy_NotFound.md) — déduit : teste commande `myr network destroy`
> - 🟡 [TestRunNetworkDestroy_ProductionRefused](../../docs/tests/adapters-in-cli/TestRunNetworkDestroy_ProductionRefused.md) — déduit : teste commande `myr network destroy`
> - 🟡 [TestRunNetworkDestroy_WithoutConfirm](../../docs/tests/adapters-in-cli/TestRunNetworkDestroy_WithoutConfirm.md) — déduit : teste commande `myr network destroy`
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
| [RM09](#RM09%20—%20Interface%20à%20usage%20unique) | Interface à usage unique dans une liaison, par instance ; jamais reliée à elle-même | Interfaces |
| [RM10](#RM10%20—%20Vérification%20de%20compatibilité%20automatique) | Vérification de compatibilité automatique à la création | Interfaces |
| [RM11](#RM11%20—%20Critères%20de%20compatibilité%20d'interfaces) | Critères de compatibilité : catégorie + type + sens + contraintes de valeur (unités de base communes chevauchantes, aucune en commun = incompatible) | Interfaces |
| [RM12](#RM12%20—%20Persistance%20des%20liaisons%20incompatibles) | Liaison incompatible persistante (`Incompatible:true`, non supprimée automatiquement) | Interfaces |
| [RM13](#RM13%20—%20Slot%20virtuel%20garanti) | Slot virtuel garanti et recréé automatiquement à chaque matérialisation | Interfaces |
| [RM50](#RM50%20—%20Liaison%20directe%20unique) | Liaison directe en double refusée, jamais renvoyée ni dupliquée | Interfaces |
| [RM53](#RM53%20—%20Unité%20d'une%20contrainte%20de%20valeur) | Unité lue dans le vocabulaire, sinon préfixe SI usuel + unité connue, sinon refusée ; comparaison sur la même unité de base | Interfaces |
| [RM14](#RM14%20—%20Suppression%20en%20cascade%20des%20connexions) | Suppression en cascade des connexions au retrait d'une instance | Composition (instances) |
| [RM15](#RM15%20—%20Instance%20indépendante) | Seconde instance indépendante si module déjà instancié dans le module hôte | Composition (instances) |
| [RM51](#RM51%20—%20Assemblage%20interne%20au%20module) | Seule une liaison existante entre deux instances de l'asset s'y rattache | Composition (instances) |
| [RM52](#RM52%20—%20Composition%20sans%20cycle) | Un module ne se contient pas lui-même, même indirectement | Composition (instances) |
| [RM16](#RM16%20—%20État%20draft) | État draft générique (module : toujours ; composant : optionnel via `--draft`) | Modules |
| [RM17](#RM17%20—%20Assemblage%20requis%20pour%20soumission%20%28module%20uniquement%29) | Au moins un assemblage requis pour soumettre (module uniquement) | Modules |
| [RM18](#RM18%20—%20ModuleVersion%20immuable%20%28module%20uniquement%29) | ModuleVersion immuable horodatée à la soumission (module uniquement) | Modules |
| [RM19](#RM19%20—%20Fork%20d'un%20asset%20soumis) | Asset soumis immuable (jamais re-soumis, licence, interfaces et composition figées) ; seules ses métadonnées descriptives changent, par son propriétaire ; toute autre évolution crée un fork | Modules |
| [RM54](#RM54%20—%20Seul%20le%20propriétaire%20modifie%20son%20asset) | Seul le propriétaire (identité créatrice) modifie, supprime ou soumet son asset, administrateur compris | Modules |
| [RM39](#RM39%20—%20Filiation%20d'un%20découpage) | Filiation par `ParentID` du composant d'origine pour tous les éléments d'un découpage (manuel ou assisté) | Modules |
| [RM40](#RM40%20—%20Proposition%20de%20découpage%20non%20engageante) | Une proposition de découpage automatique ne crée aucune entité tant qu'elle n'est pas validée | Modules |
| [RM41](#RM41%20—%20Compatibilité%20toujours%20vérifiée%20pour%20une%20connexion%20suggérée) | Connexion suggérée soumise au même contrôle de compatibilité qu'une liaison manuelle, sans dérogation | Modules |
| [RM20](#RM20%20—%20Identité%20=%20enrôlement%20CA,%20pas%20un%20compte%20séparé) | Identité = enrôlement CA — pas de compte email/mot de passe séparé | Compte |
| [RM21](#RM21%20—%20Rôle%20Lecteur%20par%20défaut%20à%20l'auto-enregistrement) | Rôle Lecteur par défaut à l'auto-enregistrement, sauf rôle explicite | Compte |
| [RM22](#RM22%20—%20Changement%20de%20rôle%20réservé%20à%20l'administrateur) | Changement de rôle réservé à l'administrateur (`myr identity set-role`) | Compte |
| [RM43](#RM43%20—%20Inscription%20selon%20les%20règles%20du%20réseau) | Inscription selon les règles du réseau : fermé → refus ; ouvert → identité active, ou en attente si approbation exigée | Compte |
| [RM44](#RM44%20—%20Décision%20sur%20une%20demande%20en%20attente) | Approuver active, refuser suspend — rien n'est supprimé, effet au prochain enrôlement | Compte |
| [RM45](#RM45%20—%20Connexion%20réservée%20aux%20identités%20actives%20dotées%20d'un%20rôle) | Connexion réservée aux identités actives dotées d'un rôle, lus dans le certificat | Compte |
| [RM46](#RM46%20—%20Rôle%20attribué%20toujours%20défini%20dans%20la%20table%20des%20rôles) | Rôle attribué toujours présent dans la table des rôles ; rôle d'inscription par défaut sans droit d'administration | Compte |
| [RM47](#RM47%20—%20Demande%20de%20rôle,%20attribution%20automatique%20ou%20approbation) | Demande de rôle attribuée aussitôt si le rôle est en attribution automatique, sinon soumise à l'administrateur ; jamais automatique pour un rôle d'administration | Compte |
| [RM48](#RM48%20—%20Table%20des%20rôles%20partagée%20par%20le%20réseau,%20rôle%20administrateur%20immuable) | Table des rôles publiée sur le ledger et appliquée par tout le réseau ; `admin` immuable ; table initiale tant qu'aucune n'est publiée | Compte |
| [RM49](#RM49%20—%20Lecture%20publique%20du%20catalogue,%20ouverte%20par%20l'administrateur) | Lecture du catalogue publié sans compte ni session si l'administrateur l'ouvre ; brouillons jamais visibles d'un visiteur | Compte |
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
