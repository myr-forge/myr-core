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
> [!check]- Tests — 0 explicite(s) · 42 déduit(s)
> - 🟡 [TestRunIdentityDecision_Approve_OK](../../docs/tests/adapters-in-cli/TestRunIdentityDecision_Approve_OK.md) — déduit : via UCA01 · Flux alternatif — Auto-enregistrement avec approbation (`RequireApproval=true`)
> - 🟡 [TestRunIdentityDecision_NotPending_Rejected](../../docs/tests/adapters-in-cli/TestRunIdentityDecision_NotPending_Rejected.md) — déduit : via UCA01 · Flux alternatif — Auto-enregistrement avec approbation (`RequireApproval=true`)
> - 🟡 [TestRunIdentityRegister_NominalCase](../../docs/tests/adapters-in-cli/TestRunIdentityRegister_NominalCase.md) — déduit : via UCA01 · Flux alternatif — Création de compte par l'administrateur
> - 🟡 [TestRunIdentityRequest_AutoRegister_Approved](../../docs/tests/adapters-in-cli/TestRunIdentityRequest_AutoRegister_Approved.md) — déduit : via UCA01 · Flux erreur — Réseau fermé à l'auto-enregistrement
> - 🟡 [TestRunIdentityRequest_NoNetwork_Rejected](../../docs/tests/adapters-in-cli/TestRunIdentityRequest_NoNetwork_Rejected.md) — déduit : via UCA01 · Flux erreur — Réseau fermé à l'auto-enregistrement
> - 🟡 [TestRunIdentityRequest_RequireApproval_Pending](../../docs/tests/adapters-in-cli/TestRunIdentityRequest_RequireApproval_Pending.md) — déduit : via UCA01 · Flux erreur — Réseau fermé à l'auto-enregistrement
> - 🟡 [TestRunIdentityRequests_NominalCase](../../docs/tests/adapters-in-cli/TestRunIdentityRequests_NominalCase.md) — déduit : via UCA01 · Flux alternatif — Auto-enregistrement avec approbation (`RequireApproval=true`)
> - 🟡 [TestHandleIdentityAccounts_Admin_Created](../../docs/tests/adapters-in-rest/TestHandleIdentityAccounts_Admin_Created.md) — déduit : via UCA01 · Flux alternatif — Création de compte par l'administrateur
> - 🟡 [TestHandleIdentityAccounts_InvalidNameOrRole_Rejected](../../docs/tests/adapters-in-rest/TestHandleIdentityAccounts_InvalidNameOrRole_Rejected.md) — déduit : via UCA01 · Flux alternatif — Création de compte par l'administrateur
> - 🟡 [TestHandleIdentityAccounts_NonAdmin_Forbidden](../../docs/tests/adapters-in-rest/TestHandleIdentityAccounts_NonAdmin_Forbidden.md) — déduit : via UCA01 · Flux alternatif — Création de compte par l'administrateur
> - 🟡 [TestHandleIdentityRequestDecision_Approve_OK](../../docs/tests/adapters-in-rest/TestHandleIdentityRequestDecision_Approve_OK.md) — déduit : via UCA01 · Flux alternatif — Auto-enregistrement avec approbation (`RequireApproval=true`)
> - 🟡 [TestHandleIdentityRequestDecision_Errors](../../docs/tests/adapters-in-rest/TestHandleIdentityRequestDecision_Errors.md) — déduit : via UCA01 · Flux alternatif — Auto-enregistrement avec approbation (`RequireApproval=true`)
> - 🟡 [TestHandleIdentityRequestDecision_NonAdmin_Forbidden](../../docs/tests/adapters-in-rest/TestHandleIdentityRequestDecision_NonAdmin_Forbidden.md) — déduit : via UCA01 · Flux alternatif — Auto-enregistrement avec approbation (`RequireApproval=true`)
> - 🟡 [TestHandleIdentityRequestDecision_Reject_OK](../../docs/tests/adapters-in-rest/TestHandleIdentityRequestDecision_Reject_OK.md) — déduit : via UCA01 · Flux alternatif — Auto-enregistrement avec approbation (`RequireApproval=true`)
> - 🟡 [TestHandleIdentityRequest_InvalidJSON](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_InvalidJSON.md) — déduit : via UCA01 · Flux nominal — Auto-enregistrement sans approbation (`AllowAutoRegister=true`, `RequireApproval=false`)
> - … et 27 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF02 — Authentifier une identité (enrôlement CA + token de session opaque)

- **Use cases couvrants** : UCA02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 27 déduit(s)
> - 🟡 [TestChannels_PUT_SwitchesChannel](../../docs/tests/adapters-in-rest/TestChannels_PUT_SwitchesChannel.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - 🟡 [TestHandleIdentityGuest_AllowedDeliversToken](../../docs/tests/adapters-in-rest/TestHandleIdentityGuest_AllowedDeliversToken.md) — déduit : via UCA02 · Flux nominal — Accès invité (réseau public)
> - 🟡 [TestHandleIdentityGuest_DisallowedReturnsForbidden](../../docs/tests/adapters-in-rest/TestHandleIdentityGuest_DisallowedReturnsForbidden.md) — déduit : via UCA02 · Flux nominal — Accès invité (réseau public)
> - 🟡 [TestHandleIdentityRequest_InvalidJSON](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_InvalidJSON.md) — déduit : via UCA02 · Flux alternatif — Accès invité refusé (réseau privé)
> - 🟡 [TestHandleIdentityRequest_InvalidName_Rejected](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_InvalidName_Rejected.md) — déduit : via UCA02 · Flux alternatif — Accès invité refusé (réseau privé)
> - 🟡 [TestHandleIdentityRequest_MissingFields](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_MissingFields.md) — déduit : via UCA02 · Flux alternatif — Accès invité refusé (réseau privé)
> - 🟡 [TestHandleIdentityRequest_RegistrationClosed_Rejected](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_RegistrationClosed_Rejected.md) — déduit : via UCA02 · Flux alternatif — Accès invité refusé (réseau privé)
> - 🟡 [TestHandleIdentityRequest_ServiceError](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_ServiceError.md) — déduit : via UCA02 · Flux alternatif — Accès invité refusé (réseau privé)
> - 🟡 [TestHandleIdentityRequest_Success](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_Success.md) — déduit : via UCA02 · Flux alternatif — Accès invité refusé (réseau privé)
> - 🟡 [TestHandleIdentitySession_CannotSignIn_Forbidden](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_CannotSignIn_Forbidden.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - 🟡 [TestHandleIdentitySession_EnrollError](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_EnrollError.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - 🟡 [TestHandleIdentitySession_InvalidJSON](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_InvalidJSON.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - 🟡 [TestHandleIdentitySession_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MethodNotAllowed.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - 🟡 [TestHandleIdentitySession_MissingFields](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MissingFields.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - 🟡 [TestHandleIdentitySession_NoIdentityService](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_NoIdentityService.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - … et 12 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
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
> [!check]- Tests — 0 explicite(s) · 14 déduit(s)
> - 🟡 [TestChannels_PUT_SwitchesChannel](../../docs/tests/adapters-in-rest/TestChannels_PUT_SwitchesChannel.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentityGuest_AllowedDeliversToken](../../docs/tests/adapters-in-rest/TestHandleIdentityGuest_AllowedDeliversToken.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentityGuest_DisallowedReturnsForbidden](../../docs/tests/adapters-in-rest/TestHandleIdentityGuest_DisallowedReturnsForbidden.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_CannotSignIn_Forbidden](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_CannotSignIn_Forbidden.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_EnrollError](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_EnrollError.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_InvalidJSON](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_InvalidJSON.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MethodNotAllowed.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_MissingFields](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MissingFields.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_NoIdentityService](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_NoIdentityService.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_RoleFromCertificate](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_RoleFromCertificate.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_Success](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_Success.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestRequireAuth_BlockchainWithValidToken](../../docs/tests/adapters-in-rest/TestRequireAuth_BlockchainWithValidToken.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHasPermission_BuiltinDefaults](../../docs/tests/domain-role/TestHasPermission_BuiltinDefaults.md) — déduit : via UCA05 · Flux nominal — Action autorisée
> - 🟡 [TestHasPermission_UnknownRole](../../docs/tests/domain-role/TestHasPermission_UnknownRole.md) — déduit : via UCA05 · Flux nominal — Action autorisée
<!-- tests-obsidian:end -->

#### EF06 — Consulter les assets possédés par l'utilisateur

- **Use cases couvrants** : UCA06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 78 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCA06 · Flux nominal — Composants possédés retournés
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
> - … et 63 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
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
> [!check]- Tests — 0 explicite(s) · 4 déduit(s)
> - 🟡 [TestRunNetworkDestroy_NominalCase](../../docs/tests/adapters-in-cli/TestRunNetworkDestroy_NominalCase.md) — déduit : via UCADM05 · Flux erreur — Flag --confirm absent
> - 🟡 [TestRunNetworkDestroy_NotFound](../../docs/tests/adapters-in-cli/TestRunNetworkDestroy_NotFound.md) — déduit : via UCADM05 · Flux erreur — Flag --confirm absent
> - 🟡 [TestRunNetworkDestroy_ProductionRefused](../../docs/tests/adapters-in-cli/TestRunNetworkDestroy_ProductionRefused.md) — déduit : via UCADM05 · Flux erreur — Flag --confirm absent
> - 🟡 [TestRunNetworkDestroy_WithoutConfirm](../../docs/tests/adapters-in-cli/TestRunNetworkDestroy_WithoutConfirm.md) — déduit : via UCADM05 · Flux erreur — Flag --confirm absent
<!-- tests-obsidian:end -->

### Composants — Création et édition

#### EF10 — Enregistrer un composant physique sur la blockchain

- **Use cases couvrants** : UCCE01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 116 déduit(s)
> - 🟡 [TestResolveOwner](../../docs/tests/adapters-in-cli/TestResolveOwner.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestRunModelAdd_NominalCase](../../docs/tests/adapters-in-cli/TestRunModelAdd_NominalCase.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestRunModelAdd_ServiceError_Rejected](../../docs/tests/adapters-in-cli/TestRunModelAdd_ServiceError_Rejected.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestRunModelSubmit_NominalCase](../../docs/tests/adapters-in-cli/TestRunModelSubmit_NominalCase.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestRunModelSubmit_ServiceError_Rejected](../../docs/tests/adapters-in-cli/TestRunModelSubmit_ServiceError_Rejected.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestRunModelSubmit_WithNote](../../docs/tests/adapters-in-cli/TestRunModelSubmit_WithNote.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - … et 101 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF11 — Enregistrer un composant numérique sur la blockchain

- **Use cases couvrants** : UCCE03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 32 déduit(s)
> - 🟡 [TestFabricAdapter_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_SubmitError.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestFabricAdapter_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_Success.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestFabricBlockchain_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_SubmitError.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestFabricBlockchain_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_Success.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestIntegration_StoreAndGet](../../docs/tests/adapters-out-fabric/TestIntegration_StoreAndGet.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestIntegration_VerifyIntegrity](../../docs/tests/adapters-out-fabric/TestIntegration_VerifyIntegrity.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestNetworkPool_BlockchainForSigner](../../docs/tests/adapters-out-fabric/TestNetworkPool_BlockchainForSigner.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestAddFull_AllFields](../../docs/tests/domain-model/TestAddFull_AllFields.md) — déduit : via UCCE03 · Flux nominal — Composant numérique nouveau
> - 🟡 [TestAddFull_AlwaysDraft_NoBlockchainWrite](../../docs/tests/domain-model/TestAddFull_AlwaysDraft_NoBlockchainWrite.md) — déduit : via UCCE03 · Flux nominal — Composant numérique nouveau
> - 🟡 [TestAddFull_MissingFile_ReturnsError](../../docs/tests/domain-model/TestAddFull_MissingFile_ReturnsError.md) — déduit : via UCCE03 · Flux nominal — Composant numérique nouveau
> - 🟡 [TestAddFull_NilBlockchain_StillWorks](../../docs/tests/domain-model/TestAddFull_NilBlockchain_StillWorks.md) — déduit : via UCCE03 · Flux nominal — Composant numérique nouveau
> - 🟡 [TestAddFull_NilDraftStore_Rejected](../../docs/tests/domain-model/TestAddFull_NilDraftStore_Rejected.md) — déduit : via UCCE03 · Flux nominal — Composant numérique nouveau
> - 🟡 [TestAddFull_NoFile_NoVersionNoHash](../../docs/tests/domain-model/TestAddFull_NoFile_NoVersionNoHash.md) — déduit : via UCCE03 · Flux nominal — Composant numérique nouveau
> - 🟡 [TestAddFull_STLAndSTEP_BothAccepted](../../docs/tests/domain-model/TestAddFull_STLAndSTEP_BothAccepted.md) — déduit : via UCCE03 · Flux nominal — Composant numérique nouveau
> - … et 17 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF12 — Configurer les métadonnées d'un composant

- **Use cases couvrants** : UCCE02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 16 déduit(s)
> - 🟡 [TestFabricAdapter_GetModelRecord_EvaluateError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_EvaluateError.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestFabricAdapter_GetModelRecord_InvalidJSON_ReturnsEmptyModel](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_InvalidJSON_ReturnsEmptyModel.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestFabricAdapter_GetModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_Success.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestFabricAdapter_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_SubmitError.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestFabricAdapter_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_Success.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestFabricBlockchain_GetModelRecord_EvaluateError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_EvaluateError.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestFabricBlockchain_GetModelRecord_InvalidJSON_ReturnsError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_InvalidJSON_ReturnsError.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestFabricBlockchain_GetModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_Success.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestFabricBlockchain_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_SubmitError.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestFabricBlockchain_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_Success.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestIntegration_GetModelRecord_NotFound](../../docs/tests/adapters-out-fabric/TestIntegration_GetModelRecord_NotFound.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestIntegration_StoreAndGet](../../docs/tests/adapters-out-fabric/TestIntegration_StoreAndGet.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestIntegration_VerifyIntegrity](../../docs/tests/adapters-out-fabric/TestIntegration_VerifyIntegrity.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestNetworkPool_BlockchainForSigner](../../docs/tests/adapters-out-fabric/TestNetworkPool_BlockchainForSigner.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - … et 1 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF13 — Faire évoluer un composant (amélioration, dérivation, extension)

- **Use cases couvrants** : UCCE04, UCCE05

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 44 déduit(s)
> - 🟡 [TestFabricAdapter_GetModelRecord_EvaluateError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_EvaluateError.md) — déduit : via UCCE04 · Flux erreur — Asset parent introuvable
> - 🟡 [TestFabricAdapter_GetModelRecord_InvalidJSON_ReturnsEmptyModel](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_InvalidJSON_ReturnsEmptyModel.md) — déduit : via UCCE04 · Flux erreur — Asset parent introuvable
> - 🟡 [TestFabricAdapter_GetModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_Success.md) — déduit : via UCCE04 · Flux erreur — Asset parent introuvable
> - 🟡 [TestFabricAdapter_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_SubmitError.md) — déduit : via UCCE04 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestFabricAdapter_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_Success.md) — déduit : via UCCE04 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestFabricBlockchain_GetModelRecord_EvaluateError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_EvaluateError.md) — déduit : via UCCE04 · Flux erreur — Asset parent introuvable
> - 🟡 [TestFabricBlockchain_GetModelRecord_InvalidJSON_ReturnsError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_InvalidJSON_ReturnsError.md) — déduit : via UCCE04 · Flux erreur — Asset parent introuvable
> - 🟡 [TestFabricBlockchain_GetModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_Success.md) — déduit : via UCCE04 · Flux erreur — Asset parent introuvable
> - 🟡 [TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected.md) — déduit : via UCCE04 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestFabricBlockchain_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_SubmitError.md) — déduit : via UCCE04 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestFabricBlockchain_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_Success.md) — déduit : via UCCE04 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestIntegration_GetModelRecord_NotFound](../../docs/tests/adapters-out-fabric/TestIntegration_GetModelRecord_NotFound.md) — déduit : via UCCE04 · Flux erreur — Asset parent introuvable
> - 🟡 [TestIntegration_StoreAndGet](../../docs/tests/adapters-out-fabric/TestIntegration_StoreAndGet.md) — déduit : via UCCE04 · Flux erreur — Asset parent introuvable
> - 🟡 [TestIntegration_VerifyIntegrity](../../docs/tests/adapters-out-fabric/TestIntegration_VerifyIntegrity.md) — déduit : via UCCE04 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestNetworkPool_BlockchainForSigner](../../docs/tests/adapters-out-fabric/TestNetworkPool_BlockchainForSigner.md) — déduit : via UCCE04 · Flux erreur — Échec endorsement Fabric
> - … et 29 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF14 — Ajouter une interface à un composant existant

- **Use cases couvrants** : UCCE06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 101 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCCE06 · Flux alternatif — Mise à jour d'une interface existante (UpdateInterface)
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCCE06 · Flux alternatif — Mise à jour d'une interface existante (UpdateInterface)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCCE06 · Flux alternatif — Mise à jour d'une interface existante (UpdateInterface)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCCE06 · Flux alternatif — Mise à jour d'une interface existante (UpdateInterface)
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCCE06 · Flux alternatif — Mise à jour d'une interface existante (UpdateInterface)
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCCE06 · Flux alternatif — Mise à jour d'une interface existante (UpdateInterface)
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCCE06 · Flux alternatif — Mise à jour d'une interface existante (UpdateInterface)
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCCE06 · Flux alternatif — Mise à jour d'une interface existante (UpdateInterface)
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCCE06 · Flux alternatif — Mise à jour d'une interface existante (UpdateInterface)
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : via UCCE06 · Flux alternatif — Mise à jour d'une interface existante (UpdateInterface)
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : via UCCE06 · Flux alternatif — Mise à jour d'une interface existante (UpdateInterface)
> - 🟡 [TestComponent_AddInstance](../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : via UCCE06 · Flux alternatif — Mise à jour d'une interface existante (UpdateInterface)
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : via UCCE06 · Flux alternatif — Mise à jour d'une interface existante (UpdateInterface)
> - 🟡 [TestComponent_AddTwoInstancesSequentially](../../docs/tests/adapters-in-rest/TestComponent_AddTwoInstancesSequentially.md) — déduit : via UCCE06 · Flux alternatif — Mise à jour d'une interface existante (UpdateInterface)
> - 🟡 [TestComponent_DELETE_NoContent](../../docs/tests/adapters-in-rest/TestComponent_DELETE_NoContent.md) — déduit : via UCCE06 · Flux alternatif — Mise à jour d'une interface existante (UpdateInterface)
> - … et 86 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF15 — Vérifier l'unicité d'un composant (anti-plagiat SHA-256 + SCM > 50 %)

- **Use cases couvrants** : UCCE01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 116 déduit(s)
> - 🟡 [TestResolveOwner](../../docs/tests/adapters-in-cli/TestResolveOwner.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestRunModelAdd_NominalCase](../../docs/tests/adapters-in-cli/TestRunModelAdd_NominalCase.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestRunModelAdd_ServiceError_Rejected](../../docs/tests/adapters-in-cli/TestRunModelAdd_ServiceError_Rejected.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestRunModelSubmit_NominalCase](../../docs/tests/adapters-in-cli/TestRunModelSubmit_NominalCase.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestRunModelSubmit_ServiceError_Rejected](../../docs/tests/adapters-in-cli/TestRunModelSubmit_ServiceError_Rejected.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestRunModelSubmit_WithNote](../../docs/tests/adapters-in-cli/TestRunModelSubmit_WithNote.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCCE01 · Flux alternatif — Création en brouillon (soumission différée, RM16/RM19)
> - … et 101 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF16 — Vérifier la compatibilité de licence lors d'une dérivation

- **Use cases couvrants** : UCCE04, UCMOD01, UCMOD06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 80 déduit(s)
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestFabricAdapter_GetModelRecord_EvaluateError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_EvaluateError.md) — déduit : via UCCE04 · Flux erreur — Asset parent introuvable
> - 🟡 [TestFabricAdapter_GetModelRecord_InvalidJSON_ReturnsEmptyModel](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_InvalidJSON_ReturnsEmptyModel.md) — déduit : via UCCE04 · Flux erreur — Asset parent introuvable
> - 🟡 [TestFabricAdapter_GetModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_Success.md) — déduit : via UCCE04 · Flux erreur — Asset parent introuvable
> - 🟡 [TestFabricAdapter_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_SubmitError.md) — déduit : via UCCE04 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestFabricAdapter_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_Success.md) — déduit : via UCCE04 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestFabricBlockchain_GetModelRecord_EvaluateError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_EvaluateError.md) — déduit : via UCCE04 · Flux erreur — Asset parent introuvable
> - 🟡 [TestFabricBlockchain_GetModelRecord_InvalidJSON_ReturnsError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_InvalidJSON_ReturnsError.md) — déduit : via UCCE04 · Flux erreur — Asset parent introuvable
> - 🟡 [TestFabricBlockchain_GetModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_Success.md) — déduit : via UCCE04 · Flux erreur — Asset parent introuvable
> - 🟡 [TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected.md) — déduit : via UCCE04 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestFabricBlockchain_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_SubmitError.md) — déduit : via UCCE04 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestFabricBlockchain_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_Success.md) — déduit : via UCCE04 · Flux erreur — Échec endorsement Fabric
> - … et 65 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF60 — Supprimer un composant (masquage local des listes, ledger jamais modifié)

- **Use cases couvrants** : UCCE07

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 72 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCCE07 · Flux nominal — Suppression (masquage) d'un composant soumis
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
> - … et 57 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Composants — Lecture

#### EF17 — Rechercher et filtrer les composants disponibles sur le réseau

- **Use cases couvrants** : UCCL01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 78 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCCL01 · Flux nominal — Résultats trouvés
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
> - … et 63 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Composition (instances)

#### EF18 — Créer des liaisons entre interfaces compatibles de composants

- **Use cases couvrants** : UCAM01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 26 déduit(s)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCAM01 · Flux — Liaison devenue incompatible après modification
> - 🟡 [TestInterface_DELETE](../../docs/tests/adapters-in-rest/TestInterface_DELETE.md) — déduit : via UCAM01 · Flux — Liaison devenue incompatible après modification
> - 🟡 [TestInterface_PATCH_NotFound](../../docs/tests/adapters-in-rest/TestInterface_PATCH_NotFound.md) — déduit : via UCAM01 · Flux — Liaison devenue incompatible après modification
> - 🟡 [TestInterface_PATCH_OK](../../docs/tests/adapters-in-rest/TestInterface_PATCH_OK.md) — déduit : via UCAM01 · Flux — Liaison devenue incompatible après modification
> - 🟡 [TestUpdateConnection](../../docs/tests/adapters-out-localstorage/TestUpdateConnection.md) — déduit : via UCAM01 · Flux — Liaison devenue incompatible après modification
> - 🟡 [TestUpdateConnection_NotFound](../../docs/tests/adapters-out-localstorage/TestUpdateConnection_NotFound.md) — déduit : via UCAM01 · Flux — Liaison devenue incompatible après modification
> - 🟡 [TestAddAssemblyLink](../../docs/tests/domain-model/TestAddAssemblyLink.md) — déduit : via UCAM01 · Flux nominal A — Liaison directe
> - 🟡 [TestAddAssemblyLink_MissingStore](../../docs/tests/domain-model/TestAddAssemblyLink_MissingStore.md) — déduit : via UCAM01 · Flux nominal A — Liaison directe
> - 🟡 [TestAddAssemblyLink_UnknownInterface](../../docs/tests/domain-model/TestAddAssemblyLink_UnknownInterface.md) — déduit : via UCAM01 · Flux nominal A — Liaison directe
> - 🟡 [TestConnectVirtualToPhysical_AssemblyLinkCreated](../../docs/tests/domain-model/TestConnectVirtualToPhysical_AssemblyLinkCreated.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - 🟡 [TestConnectVirtualToPhysical_BidirStaysBidir](../../docs/tests/domain-model/TestConnectVirtualToPhysical_BidirStaysBidir.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - 🟡 [TestConnectVirtualToPhysical_InBecomesOut](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InBecomesOut.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - 🟡 [TestConnectVirtualToPhysical_InheritsPhysicalName](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InheritsPhysicalName.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - 🟡 [TestConnectVirtualToPhysical_InheritsPhysicalValues](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InheritsPhysicalValues.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - 🟡 [TestConnectVirtualToPhysical_NilConnStore_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NilConnStore_Error.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - … et 11 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF19 — Visualiser les interfaces physiques d'un composant

- **Use cases couvrants** : UCAM02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 89 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCAM02 · Flux nominal — Interfaces d'un composant simple
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
> - … et 74 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF20 — Définir une interface sur un composant

- **Use cases couvrants** : UCAM03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 107 déduit(s)
> - 🟡 [TestRunModelInterfaceAdd_NominalCase](../../docs/tests/adapters-in-cli/TestRunModelInterfaceAdd_NominalCase.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestRunModelInterfaceAdd_ServiceError_Rejected](../../docs/tests/adapters-in-cli/TestRunModelInterfaceAdd_ServiceError_Rejected.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestRunModelLinkConnectVirtual_MissingVirtual_Rejected](../../docs/tests/adapters-in-cli/TestRunModelLinkConnectVirtual_MissingVirtual_Rejected.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestRunModelLinkConnectVirtual_NominalCase](../../docs/tests/adapters-in-cli/TestRunModelLinkConnectVirtual_NominalCase.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - … et 92 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
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
> [!check]- Tests — 0 explicite(s) · 39 déduit(s)
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - 🟡 [TestFabricAdapter_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_SubmitError.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestFabricAdapter_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_Success.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestFabricBlockchain_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_SubmitError.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestFabricBlockchain_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_Success.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestIntegration_StoreAndGet](../../docs/tests/adapters-out-fabric/TestIntegration_StoreAndGet.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestIntegration_VerifyIntegrity](../../docs/tests/adapters-out-fabric/TestIntegration_VerifyIntegrity.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestNetworkPool_BlockchainForSigner](../../docs/tests/adapters-out-fabric/TestNetworkPool_BlockchainForSigner.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestCreateModule_ParentID_CopiesComposition](../../docs/tests/domain-model/TestCreateModule_ParentID_CopiesComposition.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - 🟡 [TestCreateModule_ParentID_NotFound_Rejected](../../docs/tests/domain-model/TestCreateModule_ParentID_NotFound_Rejected.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - 🟡 [TestCreateModule_ParentID_RM03_LicenseIncompatible_Rejected](../../docs/tests/domain-model/TestCreateModule_ParentID_RM03_LicenseIncompatible_Rejected.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - … et 24 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF23 — Choisir un asset d'accroche (fastener) pour une liaison

- **Use cases couvrants** : UCAM07

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 81 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCAM07 · Flux nominal — Asset d'accroche précisé
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
> - … et 66 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF24 — Retirer une instance de composant d'un module (avec cascade des connexions)

- **Use cases couvrants** : UCAM08

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 10 déduit(s)
> - 🟡 [TestFabricAdapter_GetModelRecord_EvaluateError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_EvaluateError.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestFabricAdapter_GetModelRecord_InvalidJSON_ReturnsEmptyModel](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_InvalidJSON_ReturnsEmptyModel.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestFabricAdapter_GetModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_Success.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestFabricBlockchain_GetModelRecord_EvaluateError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_EvaluateError.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestFabricBlockchain_GetModelRecord_InvalidJSON_ReturnsError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_InvalidJSON_ReturnsError.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestFabricBlockchain_GetModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_Success.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestIntegration_GetModelRecord_NotFound](../../docs/tests/adapters-out-fabric/TestIntegration_GetModelRecord_NotFound.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestIntegration_StoreAndGet](../../docs/tests/adapters-out-fabric/TestIntegration_StoreAndGet.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestWorkspace_AddRemoveInstance](../../docs/tests/domain-model/TestWorkspace_AddRemoveInstance.md) — déduit : via UCAM08 · Flux nominal — Retrait sans liaisons actives
> - 🟡 [TestWorkspace_NilStore_ReturnsError](../../docs/tests/domain-model/TestWorkspace_NilStore_ReturnsError.md) — déduit : via UCAM08 · Flux nominal — Retrait sans liaisons actives
<!-- tests-obsidian:end -->

#### EF25 — Garantir un slot virtuel disponible sur chaque asset

- **Use cases couvrants** : UCAM03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 107 déduit(s)
> - 🟡 [TestRunModelInterfaceAdd_NominalCase](../../docs/tests/adapters-in-cli/TestRunModelInterfaceAdd_NominalCase.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestRunModelInterfaceAdd_ServiceError_Rejected](../../docs/tests/adapters-in-cli/TestRunModelInterfaceAdd_ServiceError_Rejected.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestRunModelLinkConnectVirtual_MissingVirtual_Rejected](../../docs/tests/adapters-in-cli/TestRunModelLinkConnectVirtual_MissingVirtual_Rejected.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestRunModelLinkConnectVirtual_NominalCase](../../docs/tests/adapters-in-cli/TestRunModelLinkConnectVirtual_NominalCase.md) — déduit : via UCAM03 · Flux A — Connexion virtuelle (déduction automatique)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : via UCAM03 · Flux B — Attributs explicites
> - … et 92 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
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
> [!check]- Tests — 0 explicite(s) · 39 déduit(s)
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCMOD01 · Flux nominal — Module créé en état draft
> - 🟡 [TestConnectVirtualToPhysical_AssemblyLinkCreated](../../docs/tests/domain-model/TestConnectVirtualToPhysical_AssemblyLinkCreated.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - 🟡 [TestConnectVirtualToPhysical_BidirStaysBidir](../../docs/tests/domain-model/TestConnectVirtualToPhysical_BidirStaysBidir.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - 🟡 [TestConnectVirtualToPhysical_InBecomesOut](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InBecomesOut.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - 🟡 [TestConnectVirtualToPhysical_InheritsPhysicalName](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InheritsPhysicalName.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - 🟡 [TestConnectVirtualToPhysical_InheritsPhysicalValues](../../docs/tests/domain-model/TestConnectVirtualToPhysical_InheritsPhysicalValues.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - 🟡 [TestConnectVirtualToPhysical_NilConnStore_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NilConnStore_Error.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - 🟡 [TestConnectVirtualToPhysical_NilIfaceStore_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NilIfaceStore_Error.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - 🟡 [TestConnectVirtualToPhysical_NotVirtual_Error](../../docs/tests/domain-model/TestConnectVirtualToPhysical_NotVirtual_Error.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - 🟡 [TestConnectVirtualToPhysical_OutBecomesIn](../../docs/tests/domain-model/TestConnectVirtualToPhysical_OutBecomesIn.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - 🟡 [TestConnectVirtualToPhysical_UserOverridesName](../../docs/tests/domain-model/TestConnectVirtualToPhysical_UserOverridesName.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - 🟡 [TestConnectVirtualToPhysical_UserOverridesValues](../../docs/tests/domain-model/TestConnectVirtualToPhysical_UserOverridesValues.md) — déduit : via UCMOD01 · Flux alternatif — Liaison via interface virtuelle
> - … et 24 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF27 — Soumettre un module à la blockchain (ModuleVersion immuable)

- **Use cases couvrants** : UCMOD06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 24 déduit(s)
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestFabricAdapter_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_SubmitError.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestFabricAdapter_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_Success.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestFabricBlockchain_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_SubmitError.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestFabricBlockchain_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_Success.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestIntegration_StoreAndGet](../../docs/tests/adapters-out-fabric/TestIntegration_StoreAndGet.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestIntegration_VerifyIntegrity](../../docs/tests/adapters-out-fabric/TestIntegration_VerifyIntegrity.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestNetworkPool_BlockchainForSigner](../../docs/tests/adapters-out-fabric/TestNetworkPool_BlockchainForSigner.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestCheckModuleLicense_AllCompatible](../../docs/tests/domain-model/TestCheckModuleLicense_AllCompatible.md) — déduit : via UCMOD06 · Flux alternatif — Module avec composants dérivés (licences parentales)
> - 🟡 [TestCheckModuleLicense_EmptyProposed](../../docs/tests/domain-model/TestCheckModuleLicense_EmptyProposed.md) — déduit : via UCMOD06 · Flux alternatif — Module avec composants dérivés (licences parentales)
> - 🟡 [TestCheckModuleLicense_Incompatible](../../docs/tests/domain-model/TestCheckModuleLicense_Incompatible.md) — déduit : via UCMOD06 · Flux alternatif — Module avec composants dérivés (licences parentales)
> - … et 9 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF28 — Ajouter un module existant à l'espace de travail

- **Use cases couvrants** : UCMOD02 — #incoherence : cette ligne référençait aussi UCMOD03 et UCMOD05 ; UCMOD03 ne couvre pas cet EF (il documente la modification des métadonnées d'un module, sans rapport avec l'ajout d'un module existant comme instance) et UCMOD05 ne correspond à aucun fichier existant dans `specs/1-Expression/UCMOD-Module/` ni `specs/2-Analyse/UCMOD-Module/` — origine de ces deux références à clarifier avant de les retirer définitivement

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 27 déduit(s)
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestFabricAdapter_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_SubmitError.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestFabricAdapter_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_Success.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestFabricBlockchain_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_SubmitError.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestFabricBlockchain_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_Success.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestIntegration_StoreAndGet](../../docs/tests/adapters-out-fabric/TestIntegration_StoreAndGet.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestIntegration_VerifyIntegrity](../../docs/tests/adapters-out-fabric/TestIntegration_VerifyIntegrity.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestNetworkPool_BlockchainForSigner](../../docs/tests/adapters-out-fabric/TestNetworkPool_BlockchainForSigner.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestCreateModule_ParentID_CopiesComposition](../../docs/tests/domain-model/TestCreateModule_ParentID_CopiesComposition.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestGetModuleInterfaces_InternalConnection_ScopedToInstance](../../docs/tests/domain-model/TestGetModuleInterfaces_InternalConnection_ScopedToInstance.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestGetModuleInterfaces_InternalConnectionsConsumed](../../docs/tests/domain-model/TestGetModuleInterfaces_InternalConnectionsConsumed.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - … et 12 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF29 — Visualiser la composition d'un module

- **Use cases couvrants** : UCMOD04

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 9 déduit(s)
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
> [!check]- Tests — 0 explicite(s) · 9 déduit(s)
> - 🟡 [TestFabricAdapter_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_SubmitError.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestFabricAdapter_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_Success.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestFabricBlockchain_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_SubmitError.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestFabricBlockchain_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_Success.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestIntegration_StoreAndGet](../../docs/tests/adapters-out-fabric/TestIntegration_StoreAndGet.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestIntegration_VerifyIntegrity](../../docs/tests/adapters-out-fabric/TestIntegration_VerifyIntegrity.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestNetworkPool_BlockchainForSigner](../../docs/tests/adapters-out-fabric/TestNetworkPool_BlockchainForSigner.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestUpdateAsset_Draft_StaysLocal](../../docs/tests/domain-model/TestUpdateAsset_Draft_StaysLocal.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
<!-- tests-obsidian:end -->

#### EF62 — Supprimer un module (masquage local des listes, ledger jamais modifié)

- **Use cases couvrants** : UCMOD08

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 4 déduit(s)
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCMOD08 · Flux nominal — Suppression (masquage) d'un module soumis
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCMOD08 · Flux nominal — Suppression (masquage) d'un module soumis
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCMOD08 · Flux nominal — Suppression (masquage) d'un module soumis
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCMOD08 · Flux nominal — Suppression (masquage) d'un module soumis
<!-- tests-obsidian:end -->

#### EF63 — Lister ses modules en brouillon, filtrés par propriétaire et par statut

- **Use cases couvrants** : UCMOD07

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 9 déduit(s)
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCMOD07 · Flux nominal — Modules en brouillon d'un utilisateur
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCMOD07 · Flux nominal — Modules en brouillon d'un utilisateur
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCMOD07 · Flux nominal — Modules en brouillon d'un utilisateur
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCMOD07 · Flux nominal — Modules en brouillon d'un utilisateur
> - 🟡 [TestIsModule_DraftComponent_NotMisclassified](../../docs/tests/domain-model/TestIsModule_DraftComponent_NotMisclassified.md) — déduit : via UCMOD07 · Flux nominal — Modules en brouillon d'un utilisateur
> - 🟡 [TestModule_CreateGetList](../../docs/tests/domain-model/TestModule_CreateGetList.md) — déduit : via UCMOD07 · Flux nominal — Modules en brouillon d'un utilisateur
> - 🟡 [TestModule_ListModules_BlockchainUnreachable_DegradedReturnsDraftModules](../../docs/tests/domain-model/TestModule_ListModules_BlockchainUnreachable_DegradedReturnsDraftModules.md) — déduit : via UCMOD07 · Flux nominal — Modules en brouillon d'un utilisateur
> - 🟡 [TestModule_ListNilStore_ReturnsNil](../../docs/tests/domain-model/TestModule_ListNilStore_ReturnsNil.md) — déduit : via UCMOD07 · Flux nominal — Modules en brouillon d'un utilisateur
> - 🟡 [TestRemoveModule_Submitted_HidesLocally](../../docs/tests/domain-model/TestRemoveModule_Submitted_HidesLocally.md) — déduit : via UCMOD07 · Flux nominal — Modules en brouillon d'un utilisateur
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
> [!check]- Tests — 0 explicite(s) · 72 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCPI06 · Flux traitement admin — Signalement examiné (acteur : Administrateur)
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
> - … et 57 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
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
> [!check]- Tests — 0 explicite(s) · 72 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCPI10 · Flux B — Vérification automatique (déclenchée à la soumission d'un composant)
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
> - … et 57 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Recherche

#### EF39 — Rechercher un asset par référence ou filtre

- **Use cases couvrants** : UCREC01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 76 déduit(s)
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
> - … et 61 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF40 — Identifier les composants compatibles entre eux (interfaces)

- **Use cases couvrants** : UCREC02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 6 déduit(s)
> - 🟡 [TestFabricAdapter_ListModelRecords_EvaluateError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_ListModelRecords_EvaluateError.md) — déduit : via UCREC02 · Flux erreur — Blockchain indisponible
> - 🟡 [TestFabricAdapter_ListModelRecords_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_ListModelRecords_Success.md) — déduit : via UCREC02 · Flux erreur — Blockchain indisponible
> - 🟡 [TestFabricBlockchain_ListModelRecords_EvaluateError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_ListModelRecords_EvaluateError.md) — déduit : via UCREC02 · Flux erreur — Blockchain indisponible
> - 🟡 [TestFabricBlockchain_ListModelRecords_NoPeersAvailable_ErrBlockchainUnreachable](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_ListModelRecords_NoPeersAvailable_ErrBlockchainUnreachable.md) — déduit : via UCREC02 · Flux erreur — Blockchain indisponible
> - 🟡 [TestFabricBlockchain_ListModelRecords_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_ListModelRecords_Success.md) — déduit : via UCREC02 · Flux erreur — Blockchain indisponible
> - 🟡 [TestIntegration_ListModelRecords](../../docs/tests/adapters-out-fabric/TestIntegration_ListModelRecords.md) — déduit : via UCREC02 · Flux erreur — Blockchain indisponible
<!-- tests-obsidian:end -->

#### EF41 — Consulter l'arbre de versions d'un composant

- **Use cases couvrants** : UCREC03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 74 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCREC03 · Flux nominal — Arbre de versions retourné
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
> - … et 59 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF42 — Identifier tous les modules qui intègrent un composant donné

- **Use cases couvrants** : UCREC04

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 5 déduit(s)
> - 🟡 [TestIsModule_DraftComponent_NotMisclassified](../../docs/tests/domain-model/TestIsModule_DraftComponent_NotMisclassified.md) — déduit : via UCREC04 · Flux nominal — Modules trouvés
> - 🟡 [TestModule_CreateGetList](../../docs/tests/domain-model/TestModule_CreateGetList.md) — déduit : via UCREC04 · Flux nominal — Modules trouvés
> - 🟡 [TestModule_ListModules_BlockchainUnreachable_DegradedReturnsDraftModules](../../docs/tests/domain-model/TestModule_ListModules_BlockchainUnreachable_DegradedReturnsDraftModules.md) — déduit : via UCREC04 · Flux nominal — Modules trouvés
> - 🟡 [TestModule_ListNilStore_ReturnsNil](../../docs/tests/domain-model/TestModule_ListNilStore_ReturnsNil.md) — déduit : via UCREC04 · Flux nominal — Modules trouvés
> - 🟡 [TestRemoveModule_Submitted_HidesLocally](../../docs/tests/domain-model/TestRemoveModule_Submitted_HidesLocally.md) — déduit : via UCREC04 · Flux nominal — Modules trouvés
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
> [!check]- Tests — 0 explicite(s) · 72 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCAUT03 · Flux nominal — Intégration réussie (nouveau composant)
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
> - … et 57 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF47 — Gérer les versions SCM d'un modèle 3D

- **Use cases couvrants** : UCAUT04

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 72 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCAUT04 · Flux alternatif — Consultation de l'historique
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
> - … et 57 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Interface Graphique, Paramètres, Documentation

EF48–EF54 : retirés — use cases 100 % frontend (UCIG, UCPAR, UCDOC), sans contrat REST propre à `myr`, désormais spécifiés dans le dépôt GUI externe

### Développement autour de MYR

#### EF55 — Exposer une API REST pour les intégrations tierces

- **Use cases couvrants** : UCDEV01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 84 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCDEV01 · Flux nominal — Lecture de données (query)
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
> - … et 69 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF56 — Proposer un CLI d'administration serveur

- **Use cases couvrants** : UCDEV02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 55 déduit(s)
> - 🟡 [TestResolveOwner](../../docs/tests/adapters-in-cli/TestResolveOwner.md) — déduit : via UCDEV02 · Flux nominal — Commandes opérationnelles existantes
> - 🟡 [TestRunModelAdd_NominalCase](../../docs/tests/adapters-in-cli/TestRunModelAdd_NominalCase.md) — déduit : via UCDEV02 · Flux nominal — Commandes opérationnelles existantes
> - 🟡 [TestRunModelAdd_ServiceError_Rejected](../../docs/tests/adapters-in-cli/TestRunModelAdd_ServiceError_Rejected.md) — déduit : via UCDEV02 · Flux nominal — Commandes opérationnelles existantes
> - 🟡 [TestRunModelList_BlockchainUnreachable_NoDrafts_Rejected](../../docs/tests/adapters-in-cli/TestRunModelList_BlockchainUnreachable_NoDrafts_Rejected.md) — déduit : via UCDEV02 · Flux nominal — Commandes opérationnelles existantes
> - 🟡 [TestRunModelList_BlockchainUnreachable_WithDrafts_ShowsWarningAndTable](../../docs/tests/adapters-in-cli/TestRunModelList_BlockchainUnreachable_WithDrafts_ShowsWarningAndTable.md) — déduit : via UCDEV02 · Flux nominal — Commandes opérationnelles existantes
> - 🟡 [TestRunModelList_Empty](../../docs/tests/adapters-in-cli/TestRunModelList_Empty.md) — déduit : via UCDEV02 · Flux nominal — Commandes opérationnelles existantes
> - 🟡 [TestRunModelList_NominalCase](../../docs/tests/adapters-in-cli/TestRunModelList_NominalCase.md) — déduit : via UCDEV02 · Flux nominal — Commandes opérationnelles existantes
> - 🟡 [TestRunModelVerify_Compromised](../../docs/tests/adapters-in-cli/TestRunModelVerify_Compromised.md) — déduit : via UCDEV02 · Flux nominal — Vérification d'intégrité d'un modèle
> - 🟡 [TestRunModelVerify_FileMissing](../../docs/tests/adapters-in-cli/TestRunModelVerify_FileMissing.md) — déduit : via UCDEV02 · Flux nominal — Vérification d'intégrité d'un modèle
> - 🟡 [TestRunModelVerify_HashMismatch](../../docs/tests/adapters-in-cli/TestRunModelVerify_HashMismatch.md) — déduit : via UCDEV02 · Flux nominal — Vérification d'intégrité d'un modèle
> - 🟡 [TestRunModelVerify_OK](../../docs/tests/adapters-in-cli/TestRunModelVerify_OK.md) — déduit : via UCDEV02 · Flux nominal — Vérification d'intégrité d'un modèle
> - 🟡 [TestRunNetworkActivate_NominalCase](../../docs/tests/adapters-in-cli/TestRunNetworkActivate_NominalCase.md) — déduit : via UCDEV02 · Flux nominal — Commandes d'administration réseau (UCADM01–05)
> - 🟡 [TestRunNetworkActivate_NotFound](../../docs/tests/adapters-in-cli/TestRunNetworkActivate_NotFound.md) — déduit : via UCDEV02 · Flux nominal — Commandes d'administration réseau (UCADM01–05)
> - 🟡 [TestRunNetworkAdd_InvalidPeer_Rejected](../../docs/tests/adapters-in-cli/TestRunNetworkAdd_InvalidPeer_Rejected.md) — déduit : via UCDEV02 · Flux nominal — Commandes d'administration réseau (UCADM01–05)
> - 🟡 [TestRunNetworkAdd_NoApproval_NotSet](../../docs/tests/adapters-in-cli/TestRunNetworkAdd_NoApproval_NotSet.md) — déduit : via UCDEV02 · Flux nominal — Commandes d'administration réseau (UCADM01–05)
> - … et 40 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF59 — Modifier le prix d'un asset (effet commandes futures uniquement)

- **Use cases couvrants** : UCPI11

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->
