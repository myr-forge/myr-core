---
tags:
  - couche/expression
  - type/tracabilite
  - relecture/incoherence
---
# Matrice de traçabilité — Exigences fonctionnelles ↔ Use Cases

Ce document croise les exigences fonctionnelles (EF) du projet Myr avec les use cases (UC) qui les couvrent. Il permet de vérifier qu'aucune exigence n'est orpheline et qu'aucun UC n'est sans motivation fonctionnelle.

---

## 1. Exigences fonctionnelles et UC couvrant

Les exigences sont dérivées des objectifs du projet (section 2.1), des contraintes (section 2.3) et de l'ensemble des use cases.

### Compte et Accès

#### EF01 — Permettre à un visiteur de créer un compte sur un réseau

- **Use cases couvrants** : UCA01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 10 déduit(s)
> - 🟡 [TestHandleIdentityRequest_InvalidJSON](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_InvalidJSON.md) — déduit : via UCA01 · Flux nominal — Auto-enregistrement (`AllowAutoRegister=true`)
> - 🟡 [TestHandleIdentityRequest_MissingFields](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_MissingFields.md) — déduit : via UCA01 · Flux nominal — Auto-enregistrement (`AllowAutoRegister=true`)
> - 🟡 [TestHandleIdentityRequest_ServiceError](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_ServiceError.md) — déduit : via UCA01 · Flux nominal — Auto-enregistrement (`AllowAutoRegister=true`)
> - 🟡 [TestHandleIdentityRequest_Success](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_Success.md) — déduit : via UCA01 · Flux nominal — Auto-enregistrement (`AllowAutoRegister=true`)
> - 🟡 [TestAutoRegister_NoCA](../../docs/tests/domain-identity/TestAutoRegister_NoCA.md) — déduit : via UCA01 · Flux nominal — Auto-enregistrement (`AllowAutoRegister=true`)
> - 🟡 [TestAutoRegister_RoleForwarded](../../docs/tests/domain-identity/TestAutoRegister_RoleForwarded.md) — déduit : via UCA01 · Flux nominal — Auto-enregistrement (`AllowAutoRegister=true`)
> - 🟡 [TestAutoRegister_Success](../../docs/tests/domain-identity/TestAutoRegister_Success.md) — déduit : via UCA01 · Flux nominal — Auto-enregistrement (`AllowAutoRegister=true`)
> - 🟡 [TestHandleIdentityRequests_Admin_EmptyList](../../docs/tests/adapters-in-rest/TestHandleIdentityRequests_Admin_EmptyList.md) — déduit : via UCA01 · Flux alternatif — Demande en attente (`AllowAutoRegister=false`)
> - 🟡 [TestHandleIdentityRequests_Admin_ReturnsList](../../docs/tests/adapters-in-rest/TestHandleIdentityRequests_Admin_ReturnsList.md) — déduit : via UCA01 · Flux alternatif — Demande en attente (`AllowAutoRegister=false`)
> - 🟡 [TestHandleIdentityRequests_NonAdmin_Forbidden](../../docs/tests/adapters-in-rest/TestHandleIdentityRequests_NonAdmin_Forbidden.md) — déduit : via UCA01 · Flux alternatif — Demande en attente (`AllowAutoRegister=false`)
<!-- tests-obsidian:end -->

#### EF02 — Authentifier une identité (enrôlement CA + token de session opaque)

- **Use cases couvrants** : UCA02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 10 déduit(s)
> - 🟡 [TestChannels_PUT_SwitchesChannel](../../docs/tests/adapters-in-rest/TestChannels_PUT_SwitchesChannel.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - 🟡 [TestHandleIdentitySession_EnrollError](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_EnrollError.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - 🟡 [TestHandleIdentitySession_InvalidJSON](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_InvalidJSON.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - 🟡 [TestHandleIdentitySession_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MethodNotAllowed.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - 🟡 [TestHandleIdentitySession_MissingFields](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MissingFields.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - 🟡 [TestHandleIdentitySession_NoIdentityService](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_NoIdentityService.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - 🟡 [TestHandleIdentitySession_Success](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_Success.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - 🟡 [TestRequireAuth_BlockchainWithValidToken](../../docs/tests/adapters-in-rest/TestRequireAuth_BlockchainWithValidToken.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - 🟡 [TestHandleIdentityGuest_AllowedDeliversToken](../../docs/tests/adapters-in-rest/TestHandleIdentityGuest_AllowedDeliversToken.md) — déduit : via UCA02 · Flux nominal — Accès invité (réseau public)
> - 🟡 [TestHandleIdentityGuest_DisallowedReturnsForbidden](../../docs/tests/adapters-in-rest/TestHandleIdentityGuest_DisallowedReturnsForbidden.md) — déduit : via UCA02 · Flux nominal — Accès invité (réseau public)
<!-- tests-obsidian:end -->

#### EF03 — Déconnecter un utilisateur et invalider sa session

- **Use cases couvrants** : UCA03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF04 — Vérifier la validité d'une session active

- **Use cases couvrants** : UCA04

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF05 — Contrôler les accès selon le rôle attribué

- **Use cases couvrants** : UCA05, UCA07

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 12 déduit(s)
> - 🟡 [TestHasPermission_BuiltinDefaults](../../docs/tests/domain-role/TestHasPermission_BuiltinDefaults.md) — déduit : via UCA05 · Flux nominal — Action autorisée
> - 🟡 [TestHasPermission_UnknownRole](../../docs/tests/domain-role/TestHasPermission_UnknownRole.md) — déduit : via UCA05 · Flux nominal — Action autorisée
> - 🟡 [TestChannels_PUT_SwitchesChannel](../../docs/tests/adapters-in-rest/TestChannels_PUT_SwitchesChannel.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentityGuest_AllowedDeliversToken](../../docs/tests/adapters-in-rest/TestHandleIdentityGuest_AllowedDeliversToken.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentityGuest_DisallowedReturnsForbidden](../../docs/tests/adapters-in-rest/TestHandleIdentityGuest_DisallowedReturnsForbidden.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_EnrollError](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_EnrollError.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_InvalidJSON](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_InvalidJSON.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MethodNotAllowed.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_MissingFields](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MissingFields.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_NoIdentityService](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_NoIdentityService.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_Success](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_Success.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestRequireAuth_BlockchainWithValidToken](../../docs/tests/adapters-in-rest/TestRequireAuth_BlockchainWithValidToken.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
<!-- tests-obsidian:end -->

#### EF06 — Consulter les assets possédés par l'utilisateur

- **Use cases couvrants** : UCA06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 69 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - 🟡 [TestComponent_AddInstance](../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - 🟡 [TestComponent_AddTwoInstancesSequentially](../../docs/tests/adapters-in-rest/TestComponent_AddTwoInstancesSequentially.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - 🟡 [TestComponent_DELETE_NoContent](../../docs/tests/adapters-in-rest/TestComponent_DELETE_NoContent.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - 🟡 [TestComponent_DELETE_ServiceError](../../docs/tests/adapters-in-rest/TestComponent_DELETE_ServiceError.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - … et 54 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Administration

#### EF07 — Créer un réseau blockchain indépendant

- **Use cases couvrants** : UCADM02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF08 — Gérer les organisations membres d'un réseau (ajout, mise à jour)

- **Use cases couvrants** : UCADM01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF09 — Étendre un réseau avec de nouveaux nœuds (peer ou orderer)

- **Use cases couvrants** : UCADM03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF57 — Retirer administrativement un nœud d'un réseau existant

- **Use cases couvrants** : UCADM04

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF58 — Démanteler un réseau de test (CLI uniquement — jamais via REST)

- **Use cases couvrants** : UCADM05

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 2 déduit(s)
> - 🟡 [TestDelete_NotFound](../../docs/tests/domain-network/TestDelete_NotFound.md) — déduit : via UCADM05 · Flux nominal — Réseau démantelé avec succès
> - 🟡 [TestDelete_Success](../../docs/tests/domain-network/TestDelete_Success.md) — déduit : via UCADM05 · Flux nominal — Réseau démantelé avec succès
<!-- tests-obsidian:end -->

### Composants — Création et édition

#### EF10 — Enregistrer un composant physique sur la blockchain

- **Use cases couvrants** : UCCE01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 78 déduit(s)
> - 🟡 [TestAddFull_AllFields](../../docs/tests/domain-model/TestAddFull_AllFields.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAddFull_AlwaysDraft_NoBlockchainWrite](../../docs/tests/domain-model/TestAddFull_AlwaysDraft_NoBlockchainWrite.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAddFull_MissingFile_ReturnsError](../../docs/tests/domain-model/TestAddFull_MissingFile_ReturnsError.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAddFull_NilBlockchain_StillWorks](../../docs/tests/domain-model/TestAddFull_NilBlockchain_StillWorks.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAddFull_NilDraftStore_Rejected](../../docs/tests/domain-model/TestAddFull_NilDraftStore_Rejected.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAddFull_NoFile_NoVersionNoHash](../../docs/tests/domain-model/TestAddFull_NoFile_NoVersionNoHash.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAddFull_STLAndSTEP_BothAccepted](../../docs/tests/domain-model/TestAddFull_STLAndSTEP_BothAccepted.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAddFull_WithFile_HashAndVersionCreated](../../docs/tests/domain-model/TestAddFull_WithFile_HashAndVersionCreated.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAddFull_WithParent](../../docs/tests/domain-model/TestAddFull_WithParent.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - … et 63 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF11 — Enregistrer un composant numérique sur la blockchain

- **Use cases couvrants** : UCCE03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 9 déduit(s)
> - 🟡 [TestAddFull_AllFields](../../docs/tests/domain-model/TestAddFull_AllFields.md) — déduit : via UCCE03 · Flux nominal — Composant numérique nouveau
> - 🟡 [TestAddFull_AlwaysDraft_NoBlockchainWrite](../../docs/tests/domain-model/TestAddFull_AlwaysDraft_NoBlockchainWrite.md) — déduit : via UCCE03 · Flux nominal — Composant numérique nouveau
> - 🟡 [TestAddFull_MissingFile_ReturnsError](../../docs/tests/domain-model/TestAddFull_MissingFile_ReturnsError.md) — déduit : via UCCE03 · Flux nominal — Composant numérique nouveau
> - 🟡 [TestAddFull_NilBlockchain_StillWorks](../../docs/tests/domain-model/TestAddFull_NilBlockchain_StillWorks.md) — déduit : via UCCE03 · Flux nominal — Composant numérique nouveau
> - 🟡 [TestAddFull_NilDraftStore_Rejected](../../docs/tests/domain-model/TestAddFull_NilDraftStore_Rejected.md) — déduit : via UCCE03 · Flux nominal — Composant numérique nouveau
> - 🟡 [TestAddFull_NoFile_NoVersionNoHash](../../docs/tests/domain-model/TestAddFull_NoFile_NoVersionNoHash.md) — déduit : via UCCE03 · Flux nominal — Composant numérique nouveau
> - 🟡 [TestAddFull_STLAndSTEP_BothAccepted](../../docs/tests/domain-model/TestAddFull_STLAndSTEP_BothAccepted.md) — déduit : via UCCE03 · Flux nominal — Composant numérique nouveau
> - 🟡 [TestAddFull_WithFile_HashAndVersionCreated](../../docs/tests/domain-model/TestAddFull_WithFile_HashAndVersionCreated.md) — déduit : via UCCE03 · Flux nominal — Composant numérique nouveau
> - 🟡 [TestAddFull_WithParent](../../docs/tests/domain-model/TestAddFull_WithParent.md) — déduit : via UCCE03 · Flux nominal — Composant numérique nouveau
<!-- tests-obsidian:end -->

#### EF12 — Configurer les métadonnées d'un composant

- **Use cases couvrants** : UCCE02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 1 déduit(s)
> - 🟡 [TestUpdateAsset_Draft_StaysLocal](../../docs/tests/domain-model/TestUpdateAsset_Draft_StaysLocal.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
<!-- tests-obsidian:end -->

#### EF13 — Faire évoluer un composant (amélioration, dérivation, extension)

- **Use cases couvrants** : UCCE04, UCCE05

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 10 déduit(s)
> - 🟡 [TestAddFull_AllFields](../../docs/tests/domain-model/TestAddFull_AllFields.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestAddFull_AlwaysDraft_NoBlockchainWrite](../../docs/tests/domain-model/TestAddFull_AlwaysDraft_NoBlockchainWrite.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestAddFull_MissingFile_ReturnsError](../../docs/tests/domain-model/TestAddFull_MissingFile_ReturnsError.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestAddFull_NilBlockchain_StillWorks](../../docs/tests/domain-model/TestAddFull_NilBlockchain_StillWorks.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestAddFull_NilDraftStore_Rejected](../../docs/tests/domain-model/TestAddFull_NilDraftStore_Rejected.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestAddFull_NoFile_NoVersionNoHash](../../docs/tests/domain-model/TestAddFull_NoFile_NoVersionNoHash.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestAddFull_STLAndSTEP_BothAccepted](../../docs/tests/domain-model/TestAddFull_STLAndSTEP_BothAccepted.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestAddFull_WithFile_HashAndVersionCreated](../../docs/tests/domain-model/TestAddFull_WithFile_HashAndVersionCreated.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestAddFull_WithParent](../../docs/tests/domain-model/TestAddFull_WithParent.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestListInterfacesForAsset](../../docs/tests/adapters-out-localstorage/TestListInterfacesForAsset.md) — déduit : via UCCE05 · Flux nominal — Extension créée avec succès
<!-- tests-obsidian:end -->

#### EF14 — Ajouter une interface à un composant existant

- **Use cases couvrants** : UCCE06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 84 déduit(s)
> - 🟡 [TestGetRefs_EmptyCategories_ReturnsDefaults](../../docs/tests/adapters-out-localstorage/TestGetRefs_EmptyCategories_ReturnsDefaults.md) — déduit : via UCCE06 · Flux nominal — Composant en brouillon : interface ajoutée localement
> - 🟡 [TestGetRefs_EmptyFile_ReturnsDefaults](../../docs/tests/adapters-out-localstorage/TestGetRefs_EmptyFile_ReturnsDefaults.md) — déduit : via UCCE06 · Flux nominal — Composant en brouillon : interface ajoutée localement
> - 🟡 [TestSaveInterface_And_GetInterface](../../docs/tests/adapters-out-localstorage/TestSaveInterface_And_GetInterface.md) — déduit : via UCCE06 · Flux nominal — Composant en brouillon : interface ajoutée localement
> - 🟡 [TestSaveInterface_Update](../../docs/tests/adapters-out-localstorage/TestSaveInterface_Update.md) — déduit : via UCCE06 · Flux nominal — Composant en brouillon : interface ajoutée localement
> - 🟡 [TestConnectVirtualToPhysical_AssemblyLinkCreated](../../docs/tests/domain-model/TestConnectVirtualToPhysical_AssemblyLinkCreated.md) — déduit : via UCCE06 · Flux alternatif — Interface virtuelle (slot de connexion non encore typé)
> - 🟡 [TestConnectVirtualToPhysical_BidirStaysBidir](../../docs/tests/domain-model/TestConnectVirtualToPhysical_BidirStaysBidir.md) — déduit : via UCCE06 · Flux alternatif — Interface virtuelle (slot de connexion non encore typé)
> - 🟡 [TestConnectVirtualToPhysical_InBecomesOut](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InBecomesOut.md) — déduit : via UCCE06 · Flux alternatif — Interface virtuelle (slot de connexion non encore typé)
> - 🟡 [TestConnectVirtualToPhysical_InheritsPhysicalName](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InheritsPhysicalName.md) — déduit : via UCCE06 · Flux alternatif — Interface virtuelle (slot de connexion non encore typé)
> - 🟡 [TestConnectVirtualToPhysical_InheritsPhysicalValues](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InheritsPhysicalValues.md) — déduit : via UCCE06 · Flux alternatif — Interface virtuelle (slot de connexion non encore typé)
> - 🟡 [TestConnectVirtualToPhysical_NilConnStore_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NilConnStore_Error.md) — déduit : via UCCE06 · Flux alternatif — Interface virtuelle (slot de connexion non encore typé)
> - 🟡 [TestConnectVirtualToPhysical_NilIfaceStore_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NilIfaceStore_Error.md) — déduit : via UCCE06 · Flux alternatif — Interface virtuelle (slot de connexion non encore typé)
> - 🟡 [TestConnectVirtualToPhysical_NotVirtual_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NotVirtual_Error.md) — déduit : via UCCE06 · Flux alternatif — Interface virtuelle (slot de connexion non encore typé)
> - 🟡 [TestConnectVirtualToPhysical_OutBecomesIn](../../docs/tests/domain-model/TestConnectVirtualToPhysical_OutBecomesIn.md) — déduit : via UCCE06 · Flux alternatif — Interface virtuelle (slot de connexion non encore typé)
> - 🟡 [TestConnectVirtualToPhysical_UserOverridesName](../../docs/tests/domain-model/TestConnectVirtualToPhysical_UserOverridesName.md) — déduit : via UCCE06 · Flux alternatif — Interface virtuelle (slot de connexion non encore typé)
> - 🟡 [TestConnectVirtualToPhysical_UserOverridesValues](../../docs/tests/domain-model/TestConnectVirtualToPhysical_UserOverridesValues.md) — déduit : via UCCE06 · Flux alternatif — Interface virtuelle (slot de connexion non encore typé)
> - … et 69 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF15 — Vérifier l'unicité d'un composant (anti-plagiat SHA-256 + SCM > 50 %)

- **Use cases couvrants** : UCCE01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 78 déduit(s)
> - 🟡 [TestAddFull_AllFields](../../docs/tests/domain-model/TestAddFull_AllFields.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAddFull_AlwaysDraft_NoBlockchainWrite](../../docs/tests/domain-model/TestAddFull_AlwaysDraft_NoBlockchainWrite.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAddFull_MissingFile_ReturnsError](../../docs/tests/domain-model/TestAddFull_MissingFile_ReturnsError.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAddFull_NilBlockchain_StillWorks](../../docs/tests/domain-model/TestAddFull_NilBlockchain_StillWorks.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAddFull_NilDraftStore_Rejected](../../docs/tests/domain-model/TestAddFull_NilDraftStore_Rejected.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAddFull_NoFile_NoVersionNoHash](../../docs/tests/domain-model/TestAddFull_NoFile_NoVersionNoHash.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAddFull_STLAndSTEP_BothAccepted](../../docs/tests/domain-model/TestAddFull_STLAndSTEP_BothAccepted.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAddFull_WithFile_HashAndVersionCreated](../../docs/tests/domain-model/TestAddFull_WithFile_HashAndVersionCreated.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAddFull_WithParent](../../docs/tests/domain-model/TestAddFull_WithParent.md) — déduit : via UCCE01 · Flux nominal — Composant base nouveau
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - … et 63 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF16 — Vérifier la compatibilité de licence lors d'une dérivation

- **Use cases couvrants** : UCCE04, UCMOD01, UCMOD06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 28 déduit(s)
> - 🟡 [TestAddFull_AllFields](../../docs/tests/domain-model/TestAddFull_AllFields.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestAddFull_AlwaysDraft_NoBlockchainWrite](../../docs/tests/domain-model/TestAddFull_AlwaysDraft_NoBlockchainWrite.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestAddFull_MissingFile_ReturnsError](../../docs/tests/domain-model/TestAddFull_MissingFile_ReturnsError.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestAddFull_NilBlockchain_StillWorks](../../docs/tests/domain-model/TestAddFull_NilBlockchain_StillWorks.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestAddFull_NilDraftStore_Rejected](../../docs/tests/domain-model/TestAddFull_NilDraftStore_Rejected.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestAddFull_NoFile_NoVersionNoHash](../../docs/tests/domain-model/TestAddFull_NoFile_NoVersionNoHash.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestAddFull_STLAndSTEP_BothAccepted](../../docs/tests/domain-model/TestAddFull_STLAndSTEP_BothAccepted.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestAddFull_WithFile_HashAndVersionCreated](../../docs/tests/domain-model/TestAddFull_WithFile_HashAndVersionCreated.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestAddFull_WithParent](../../docs/tests/domain-model/TestAddFull_WithParent.md) — déduit : via UCCE04 · Flux nominal — Amélioration réussie (mêmes interfaces)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestCreateModule_ParentID_CopiesComposition](../../docs/tests/domain-model/TestCreateModule_ParentID_CopiesComposition.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - … et 13 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF60 — Supprimer un composant (masquage local des listes, ledger jamais modifié)

- **Use cases couvrants** : UCCE07

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 69 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - 🟡 [TestComponent_AddInstance](../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - 🟡 [TestComponent_AddTwoInstancesSequentially](../../docs/tests/adapters-in-rest/TestComponent_AddTwoInstancesSequentially.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - 🟡 [TestComponent_DELETE_NoContent](../../docs/tests/adapters-in-rest/TestComponent_DELETE_NoContent.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - 🟡 [TestComponent_DELETE_ServiceError](../../docs/tests/adapters-in-rest/TestComponent_DELETE_ServiceError.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - … et 54 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Composants — Lecture

#### EF17 — Rechercher et filtrer les composants disponibles sur le réseau

- **Use cases couvrants** : UCCL01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 70 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - 🟡 [TestComponent_AddInstance](../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - 🟡 [TestComponent_AddTwoInstancesSequentially](../../docs/tests/adapters-in-rest/TestComponent_AddTwoInstancesSequentially.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - 🟡 [TestComponent_DELETE_NoContent](../../docs/tests/adapters-in-rest/TestComponent_DELETE_NoContent.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - 🟡 [TestComponent_DELETE_ServiceError](../../docs/tests/adapters-in-rest/TestComponent_DELETE_ServiceError.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - … et 55 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Composition (instances)

#### EF18 — Créer des liaisons entre interfaces compatibles de composants

- **Use cases couvrants** : UCAM01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 20 déduit(s)
> - 🟡 [TestAddAssemblyLink](../../docs/tests/domain-model/TestAddAssemblyLink.md) — déduit : via UCAM01 · Flux nominal A — Liaison directe
> - 🟡 [TestAddAssemblyLink_MissingStore](../../docs/tests/domain-model/TestAddAssemblyLink_MissingStore.md) — déduit : via UCAM01 · Flux nominal A — Liaison directe
> - 🟡 [TestAddAssemblyLink_UnknownInterface](../../docs/tests/domain-model/TestAddAssemblyLink_UnknownInterface.md) — déduit : via UCAM01 · Flux nominal A — Liaison directe
> - 🟡 [TestConnectVirtualToPhysical_AssemblyLinkCreated](../../docs/tests/domain-model/TestConnectVirtualToPhysical_AssemblyLinkCreated.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - 🟡 [TestConnectVirtualToPhysical_BidirStaysBidir](../../docs/tests/domain-model/TestConnectVirtualToPhysical_BidirStaysBidir.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - 🟡 [TestConnectVirtualToPhysical_InBecomesOut](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InBecomesOut.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - 🟡 [TestConnectVirtualToPhysical_InheritsPhysicalName](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InheritsPhysicalName.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - 🟡 [TestConnectVirtualToPhysical_InheritsPhysicalValues](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InheritsPhysicalValues.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - 🟡 [TestConnectVirtualToPhysical_NilConnStore_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NilConnStore_Error.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - 🟡 [TestConnectVirtualToPhysical_NilIfaceStore_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NilIfaceStore_Error.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - 🟡 [TestConnectVirtualToPhysical_NotVirtual_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NotVirtual_Error.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - 🟡 [TestConnectVirtualToPhysical_OutBecomesIn](../../docs/tests/domain-model/TestConnectVirtualToPhysical_OutBecomesIn.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - 🟡 [TestConnectVirtualToPhysical_UserOverridesName](../../docs/tests/domain-model/TestConnectVirtualToPhysical_UserOverridesName.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - 🟡 [TestConnectVirtualToPhysical_UserOverridesValues](../../docs/tests/domain-model/TestConnectVirtualToPhysical_UserOverridesValues.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCAM01 · Flux — Liaison devenue incompatible après modification
> - … et 5 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF19 — Visualiser les interfaces physiques d'un composant

- **Use cases couvrants** : UCAM02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 80 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - 🟡 [TestComponent_AddInstance](../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - 🟡 [TestComponent_AddTwoInstancesSequentially](../../docs/tests/adapters-in-rest/TestComponent_AddTwoInstancesSequentially.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - 🟡 [TestComponent_DELETE_NoContent](../../docs/tests/adapters-in-rest/TestComponent_DELETE_NoContent.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - 🟡 [TestComponent_DELETE_ServiceError](../../docs/tests/adapters-in-rest/TestComponent_DELETE_ServiceError.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - … et 65 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF20 — Définir une interface sur un composant

- **Use cases couvrants** : UCAM03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 85 déduit(s)
> - 🟡 [TestSaveInterface_And_GetInterface](../../docs/tests/adapters-out-localstorage/TestSaveInterface_And_GetInterface.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestSaveInterface_Update](../../docs/tests/adapters-out-localstorage/TestSaveInterface_Update.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestAddAssemblyLink](../../docs/tests/domain-model/TestAddAssemblyLink.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestAddAssemblyLink_MissingStore](../../docs/tests/domain-model/TestAddAssemblyLink_MissingStore.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestAddAssemblyLink_UnknownInterface](../../docs/tests/domain-model/TestAddAssemblyLink_UnknownInterface.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_AssemblyLinkCreated](../../docs/tests/domain-model/TestConnectVirtualToPhysical_AssemblyLinkCreated.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_BidirStaysBidir](../../docs/tests/domain-model/TestConnectVirtualToPhysical_BidirStaysBidir.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_InBecomesOut](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InBecomesOut.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_InheritsPhysicalName](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InheritsPhysicalName.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_InheritsPhysicalValues](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InheritsPhysicalValues.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_NilConnStore_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NilConnStore_Error.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_NilIfaceStore_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NilIfaceStore_Error.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_NotVirtual_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NotVirtual_Error.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_OutBecomesIn](../../docs/tests/domain-model/TestConnectVirtualToPhysical_OutBecomesIn.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_UserOverridesName](../../docs/tests/domain-model/TestConnectVirtualToPhysical_UserOverridesName.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - … et 70 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF21 — *(retiré)* — sélection multiple/jauge de progression pour l'ajout d'instances : ergonomie 100 % frontend, sans logique domaine propre (le placement unitaire `AddAssetToWorkspace` reste couvert par EF26/UCMOD01 et UCAM05)

- **Use cases couvrants** : —

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF22 — Transformer un composant en module (découpage en sous-systèmes)

- **Use cases couvrants** : UCAM05

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 8 déduit(s)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - 🟡 [TestCreateModule_ParentID_CopiesComposition](../../docs/tests/domain-model/TestCreateModule_ParentID_CopiesComposition.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - 🟡 [TestCreateModule_ParentID_NotFound_Rejected](../../docs/tests/domain-model/TestCreateModule_ParentID_NotFound_Rejected.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - 🟡 [TestCreateModule_ParentID_RM03_LicenseIncompatible_Rejected](../../docs/tests/domain-model/TestCreateModule_ParentID_RM03_LicenseIncompatible_Rejected.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
<!-- tests-obsidian:end -->

#### EF23 — Choisir un asset d'accroche (fastener) pour une liaison

- **Use cases couvrants** : UCAM07

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 76 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - 🟡 [TestAssemblyLinks_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestAssemblyLinks_MethodNotAllowed.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - 🟡 [TestAssemblyLinks_POST_Created](../../docs/tests/adapters-in-rest/TestAssemblyLinks_POST_Created.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - 🟡 [TestAssemblyLinks_POST_MissingToIface](../../docs/tests/adapters-in-rest/TestAssemblyLinks_POST_MissingToIface.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - 🟡 [TestComponent_AddInstance](../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - … et 61 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF24 — Retirer une instance de composant d'un module (avec cascade des connexions)

- **Use cases couvrants** : UCAM08

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 1 déduit(s)
> - 🟡 [TestRemoveConnection](../../docs/tests/adapters-out-localstorage/TestRemoveConnection.md) — déduit : via UCAM08 · Flux nominal — Retrait avec liaisons en cascade
<!-- tests-obsidian:end -->

#### EF25 — Garantir un slot virtuel disponible sur chaque asset

- **Use cases couvrants** : UCAM03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 85 déduit(s)
> - 🟡 [TestSaveInterface_And_GetInterface](../../docs/tests/adapters-out-localstorage/TestSaveInterface_And_GetInterface.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestSaveInterface_Update](../../docs/tests/adapters-out-localstorage/TestSaveInterface_Update.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestAddAssemblyLink](../../docs/tests/domain-model/TestAddAssemblyLink.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestAddAssemblyLink_MissingStore](../../docs/tests/domain-model/TestAddAssemblyLink_MissingStore.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestAddAssemblyLink_UnknownInterface](../../docs/tests/domain-model/TestAddAssemblyLink_UnknownInterface.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_AssemblyLinkCreated](../../docs/tests/domain-model/TestConnectVirtualToPhysical_AssemblyLinkCreated.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_BidirStaysBidir](../../docs/tests/domain-model/TestConnectVirtualToPhysical_BidirStaysBidir.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_InBecomesOut](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InBecomesOut.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_InheritsPhysicalName](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InheritsPhysicalName.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_InheritsPhysicalValues](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InheritsPhysicalValues.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_NilConnStore_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NilConnStore_Error.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_NilIfaceStore_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NilIfaceStore_Error.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_NotVirtual_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NotVirtual_Error.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_OutBecomesIn](../../docs/tests/domain-model/TestConnectVirtualToPhysical_OutBecomesIn.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestConnectVirtualToPhysical_UserOverridesName](../../docs/tests/domain-model/TestConnectVirtualToPhysical_UserOverridesName.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - … et 70 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF64 — Proposer un découpage automatique (sous-pièces + connexions candidates) d'un composant STEP en amont d'une transformation composant → module

- **Use cases couvrants** : UCAM09

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Modules

#### EF26 — Assembler plusieurs composants en module (état draft), y compris par dérivation d'un module existant (composition dupliquée depuis un `parent_id`)

- **Use cases couvrants** : UCMOD01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 19 déduit(s)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestCreateModule_ParentID_CopiesComposition](../../docs/tests/domain-model/TestCreateModule_ParentID_CopiesComposition.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestCreateModule_ParentID_NotFound_Rejected](../../docs/tests/domain-model/TestCreateModule_ParentID_NotFound_Rejected.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestCreateModule_ParentID_RM03_LicenseIncompatible_Rejected](../../docs/tests/domain-model/TestCreateModule_ParentID_RM03_LicenseIncompatible_Rejected.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestConnectVirtualToPhysical_AssemblyLinkCreated](../../docs/tests/domain-model/TestConnectVirtualToPhysical_AssemblyLinkCreated.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - 🟡 [TestConnectVirtualToPhysical_BidirStaysBidir](../../docs/tests/domain-model/TestConnectVirtualToPhysical_BidirStaysBidir.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - 🟡 [TestConnectVirtualToPhysical_InBecomesOut](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InBecomesOut.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - 🟡 [TestConnectVirtualToPhysical_InheritsPhysicalName](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InheritsPhysicalName.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - 🟡 [TestConnectVirtualToPhysical_InheritsPhysicalValues](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InheritsPhysicalValues.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - 🟡 [TestConnectVirtualToPhysical_NilConnStore_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NilConnStore_Error.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - 🟡 [TestConnectVirtualToPhysical_NilIfaceStore_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NilIfaceStore_Error.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - … et 4 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF27 — Soumettre un module à la blockchain (ModuleVersion immuable)

- **Use cases couvrants** : UCMOD06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 5 déduit(s)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
<!-- tests-obsidian:end -->

#### EF28 — Ajouter un module existant à l'espace de travail

- **Use cases couvrants** : UCMOD02 — #incoherence : cette ligne référençait aussi UCMOD03 et UCMOD05 ; UCMOD03 ne couvre pas cet EF (il documente la modification des métadonnées d'un module, sans rapport avec l'ajout d'un module existant comme instance) et UCMOD05 ne correspond à aucun fichier existant dans `specs/1-Expression/UCMOD-Module/` ni `specs/2-Analyse/UCMOD-Module/` — origine de ces deux références à clarifier avant de les retirer définitivement

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 5 déduit(s)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
<!-- tests-obsidian:end -->

#### EF29 — Visualiser la composition d'un module

- **Use cases couvrants** : UCMOD04

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 10 déduit(s)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCMOD04 · Flux nominal — Composition du module retournée
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCMOD04 · Flux nominal — Composition du module retournée
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCMOD04 · Flux nominal — Composition du module retournée
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCMOD04 · Flux nominal — Composition du module retournée
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCMOD04 · Flux nominal — Composition du module retournée
> - 🟡 [TestGetModuleInterfaces_InternalConnection_ScopedToInstance](../../docs/tests/domain-model/TestGetModuleInterfaces_InternalConnection_ScopedToInstance.md) — déduit : via UCMOD04 · Flux alternatif — Consultation des interfaces exposées
> - 🟡 [TestGetModuleInterfaces_InternalConnectionsConsumed](../../docs/tests/domain-model/TestGetModuleInterfaces_InternalConnectionsConsumed.md) — déduit : via UCMOD04 · Flux alternatif — Consultation des interfaces exposées
> - 🟡 [TestGetModuleInterfaces_MultipleInstancesSameAsset_AllExposed](../../docs/tests/domain-model/TestGetModuleInterfaces_MultipleInstancesSameAsset_AllExposed.md) — déduit : via UCMOD04 · Flux alternatif — Consultation des interfaces exposées
> - 🟡 [TestGetModuleInterfaces_NilStore](../../docs/tests/domain-model/TestGetModuleInterfaces_NilStore.md) — déduit : via UCMOD04 · Flux alternatif — Consultation des interfaces exposées
> - 🟡 [TestGetModuleInterfaces_SimpleComponent](../../docs/tests/domain-model/TestGetModuleInterfaces_SimpleComponent.md) — déduit : via UCMOD04 · Flux alternatif — Consultation des interfaces exposées
<!-- tests-obsidian:end -->

#### EF61 — Modifier les métadonnées d'un module (nom, description, licence, tags, liens)

- **Use cases couvrants** : UCMOD03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 1 déduit(s)
> - 🟡 [TestUpdateAsset_Draft_StaysLocal](../../docs/tests/domain-model/TestUpdateAsset_Draft_StaysLocal.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
<!-- tests-obsidian:end -->

#### EF62 — Supprimer un module (masquage local des listes, ledger jamais modifié)

- **Use cases couvrants** : UCMOD08

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 5 déduit(s)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCMOD08 · Flux nominal — Suppression (masquage) d'un module soumis
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCMOD08 · Flux nominal — Suppression (masquage) d'un module soumis
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCMOD08 · Flux nominal — Suppression (masquage) d'un module soumis
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCMOD08 · Flux nominal — Suppression (masquage) d'un module soumis
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCMOD08 · Flux nominal — Suppression (masquage) d'un module soumis
<!-- tests-obsidian:end -->

#### EF63 — Lister ses modules en brouillon, filtrés par propriétaire et par statut

- **Use cases couvrants** : UCMOD07

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 5 déduit(s)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCMOD07 · Flux nominal — Modules en brouillon d'un utilisateur
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCMOD07 · Flux nominal — Modules en brouillon d'un utilisateur
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCMOD07 · Flux nominal — Modules en brouillon d'un utilisateur
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCMOD07 · Flux nominal — Modules en brouillon d'un utilisateur
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCMOD07 · Flux nominal — Modules en brouillon d'un utilisateur
<!-- tests-obsidian:end -->

### Propriété Intellectuelle et Rémunération

#### EF30 — Commander un module complet (fabrication ou achat en stock)

- **Use cases couvrants** : UCPI01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF31 — Distribuer automatiquement les commissions aux auteurs à la livraison

- **Use cases couvrants** : UCPI02, UCAUT01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF32 — Définir un prix sur un composant propriétaire

- **Use cases couvrants** : UCPI04

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF33 — Définir un prix sur un module propriétaire

- **Use cases couvrants** : UCPI05

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF34 — Signaler un composant similaire à un existant

- **Use cases couvrants** : UCPI06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 69 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - 🟡 [TestComponent_AddInstance](../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - 🟡 [TestComponent_AddTwoInstancesSequentially](../../docs/tests/adapters-in-rest/TestComponent_AddTwoInstancesSequentially.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - 🟡 [TestComponent_DELETE_NoContent](../../docs/tests/adapters-in-rest/TestComponent_DELETE_NoContent.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - 🟡 [TestComponent_DELETE_ServiceError](../../docs/tests/adapters-in-rest/TestComponent_DELETE_ServiceError.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - … et 54 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF35 — Transférer la propriété intellectuelle d'un asset

- **Use cases couvrants** : UCPI07

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF36 — Cloner un composant sur un réseau externe

- **Use cases couvrants** : UCPI08

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF37 — Cloner un module sur un réseau externe

- **Use cases couvrants** : UCPI09

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF38 — Respecter et vérifier une norme d'écoconception

- **Use cases couvrants** : UCPI10

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 69 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - 🟡 [TestComponent_AddInstance](../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - 🟡 [TestComponent_AddTwoInstancesSequentially](../../docs/tests/adapters-in-rest/TestComponent_AddTwoInstancesSequentially.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - 🟡 [TestComponent_DELETE_NoContent](../../docs/tests/adapters-in-rest/TestComponent_DELETE_NoContent.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - 🟡 [TestComponent_DELETE_ServiceError](../../docs/tests/adapters-in-rest/TestComponent_DELETE_ServiceError.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - … et 54 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Recherche

#### EF39 — Rechercher un asset par référence ou filtre

- **Use cases couvrants** : UCREC01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 75 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCREC01 · Flux nominal — Référence trouvée (UUID exact)
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCREC01 · Flux nominal — Référence trouvée (UUID exact)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCREC01 · Flux nominal — Référence trouvée (UUID exact)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCREC01 · Flux nominal — Référence trouvée (UUID exact)
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCREC01 · Flux nominal — Référence trouvée (UUID exact)
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCREC01 · Flux nominal — Référence trouvée (UUID exact)
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCREC01 · Flux nominal — Référence trouvée (UUID exact)
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCREC01 · Flux nominal — Référence trouvée (UUID exact)
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCREC01 · Flux nominal — Référence trouvée (UUID exact)
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : via UCREC01 · Flux nominal — Référence trouvée (UUID exact)
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : via UCREC01 · Flux nominal — Référence trouvée (UUID exact)
> - 🟡 [TestComponent_AddInstance](../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : via UCREC01 · Flux nominal — Référence trouvée (UUID exact)
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : via UCREC01 · Flux nominal — Référence trouvée (UUID exact)
> - 🟡 [TestComponent_AddTwoInstancesSequentially](../../docs/tests/adapters-in-rest/TestComponent_AddTwoInstancesSequentially.md) — déduit : via UCREC01 · Flux nominal — Référence trouvée (UUID exact)
> - 🟡 [TestComponent_DELETE_NoContent](../../docs/tests/adapters-in-rest/TestComponent_DELETE_NoContent.md) — déduit : via UCREC01 · Flux nominal — Référence trouvée (UUID exact)
> - … et 60 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF40 — Identifier les composants compatibles entre eux (interfaces)

- **Use cases couvrants** : UCREC02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 2 déduit(s)
> - 🟡 [TestListInterfacesForAsset](../../docs/tests/adapters-out-localstorage/TestListInterfacesForAsset.md) — déduit : via UCREC02 · Flux nominal — Composants compatibles trouvés
> - 🟡 [TestList_MergesDraftsAndBlockchain](../../docs/tests/domain-model/TestList_MergesDraftsAndBlockchain.md) — déduit : via UCREC02 · Flux nominal — Composants compatibles trouvés
<!-- tests-obsidian:end -->

#### EF41 — Consulter l'arbre de versions d'un composant

- **Use cases couvrants** : UCREC03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 71 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - 🟡 [TestComponent_AddInstance](../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - 🟡 [TestComponent_AddTwoInstancesSequentially](../../docs/tests/adapters-in-rest/TestComponent_AddTwoInstancesSequentially.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - 🟡 [TestComponent_DELETE_NoContent](../../docs/tests/adapters-in-rest/TestComponent_DELETE_NoContent.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - 🟡 [TestComponent_DELETE_ServiceError](../../docs/tests/adapters-in-rest/TestComponent_DELETE_ServiceError.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - … et 56 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF42 — Identifier tous les modules qui intègrent un composant donné

- **Use cases couvrants** : UCREC04

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF43 — Exporter la BOM (Bill Of Materials) d'un module

- **Use cases couvrants** : UCREC05

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Automatisation

#### EF44 — Automatiser la fabrication et la livraison d'un composant

- **Use cases couvrants** : UCAUT01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF45 — Automatiser la commande en ligne d'un asset

- **Use cases couvrants** : UCAUT02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF46 — Intégrer un modèle 3D depuis un logiciel CAO (plugin)

- **Use cases couvrants** : UCAUT03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 69 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - 🟡 [TestComponent_AddInstance](../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - 🟡 [TestComponent_AddTwoInstancesSequentially](../../docs/tests/adapters-in-rest/TestComponent_AddTwoInstancesSequentially.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - 🟡 [TestComponent_DELETE_NoContent](../../docs/tests/adapters-in-rest/TestComponent_DELETE_NoContent.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - 🟡 [TestComponent_DELETE_ServiceError](../../docs/tests/adapters-in-rest/TestComponent_DELETE_ServiceError.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - … et 54 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF47 — Gérer les versions SCM d'un modèle 3D

- **Use cases couvrants** : UCAUT04

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 69 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - 🟡 [TestComponent_AddInstance](../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - 🟡 [TestComponent_AddTwoInstancesSequentially](../../docs/tests/adapters-in-rest/TestComponent_AddTwoInstancesSequentially.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - 🟡 [TestComponent_DELETE_NoContent](../../docs/tests/adapters-in-rest/TestComponent_DELETE_NoContent.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - 🟡 [TestComponent_DELETE_ServiceError](../../docs/tests/adapters-in-rest/TestComponent_DELETE_ServiceError.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - … et 54 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Interface Graphique, Paramètres, Documentation

EF48–EF54 : retirés — use cases 100 % frontend (UCIG, UCPAR, UCDOC), sans contrat REST propre à `myr`, désormais spécifiés dans le dépôt GUI externe

### Développement autour de MYR

#### EF55 — Exposer une API REST pour les intégrations tierces

- **Use cases couvrants** : UCDEV01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 79 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - 🟡 [TestChannels_PUT_SwitchesChannel](../../docs/tests/adapters-in-rest/TestChannels_PUT_SwitchesChannel.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - 🟡 [TestComponent_AddInstance](../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - 🟡 [TestComponent_AddTwoInstancesSequentially](../../docs/tests/adapters-in-rest/TestComponent_AddTwoInstancesSequentially.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - 🟡 [TestComponent_DELETE_NoContent](../../docs/tests/adapters-in-rest/TestComponent_DELETE_NoContent.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - … et 64 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF56 — Proposer un CLI d'administration serveur

- **Use cases couvrants** : UCDEV02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF59 — Modifier le prix d'un asset (effet commandes futures uniquement)

- **Use cases couvrants** : UCPI11

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->
