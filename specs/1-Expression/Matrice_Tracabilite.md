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
> [!check]- Tests — 0 explicite(s) · 47 déduit(s)
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
> - … et 32 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF02 — Authentifier une identité (enrôlement CA + token de session opaque)

- **Use cases couvrants** : UCA02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 30 déduit(s)
> - 🟡 [TestChannels_PUT_SwitchesChannel](../../docs/tests/adapters-in-rest/TestChannels_PUT_SwitchesChannel.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - 🟡 [TestHandleIdentityGuest_DisallowedReturnsForbidden](../../docs/tests/adapters-in-rest/TestHandleIdentityGuest_DisallowedReturnsForbidden.md) — déduit : via UCA02 · Flux nominal — Accès invité (réseau ouvert à la lecture publique)
> - 🟡 [TestHandleIdentityGuest_RM49_PublicRead_NoSession](../../docs/tests/adapters-in-rest/TestHandleIdentityGuest_RM49_PublicRead_NoSession.md) — déduit : via UCA02 · Flux nominal — Accès invité (réseau ouvert à la lecture publique)
> - 🟡 [TestHandleIdentityPolicy_BlockchainPrivate](../../docs/tests/adapters-in-rest/TestHandleIdentityPolicy_BlockchainPrivate.md) — déduit : via UCA02 · Flux nominal — Accès invité (réseau ouvert à la lecture publique)
> - 🟡 [TestHandleIdentityPolicy_BlockchainPublic](../../docs/tests/adapters-in-rest/TestHandleIdentityPolicy_BlockchainPublic.md) — déduit : via UCA02 · Flux nominal — Accès invité (réseau ouvert à la lecture publique)
> - 🟡 [TestHandleIdentityPolicy_RequireApproval](../../docs/tests/adapters-in-rest/TestHandleIdentityPolicy_RequireApproval.md) — déduit : via UCA02 · Flux nominal — Accès invité (réseau ouvert à la lecture publique)
> - 🟡 [TestHandleIdentityRequest_InvalidJSON](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_InvalidJSON.md) — déduit : via UCA02 · Flux alternatif — Accès invité refusé (réseau privé)
> - 🟡 [TestHandleIdentityRequest_InvalidName_Rejected](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_InvalidName_Rejected.md) — déduit : via UCA02 · Flux alternatif — Accès invité refusé (réseau privé)
> - 🟡 [TestHandleIdentityRequest_MissingFields](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_MissingFields.md) — déduit : via UCA02 · Flux alternatif — Accès invité refusé (réseau privé)
> - 🟡 [TestHandleIdentityRequest_RegistrationClosed_Rejected](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_RegistrationClosed_Rejected.md) — déduit : via UCA02 · Flux alternatif — Accès invité refusé (réseau privé)
> - 🟡 [TestHandleIdentityRequest_ServiceError](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_ServiceError.md) — déduit : via UCA02 · Flux alternatif — Accès invité refusé (réseau privé)
> - 🟡 [TestHandleIdentityRequest_Success](../../docs/tests/adapters-in-rest/TestHandleIdentityRequest_Success.md) — déduit : via UCA02 · Flux alternatif — Accès invité refusé (réseau privé)
> - 🟡 [TestHandleIdentitySession_CannotSignIn_Forbidden](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_CannotSignIn_Forbidden.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - 🟡 [TestHandleIdentitySession_EnrollError](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_EnrollError.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - 🟡 [TestHandleIdentitySession_InvalidJSON](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_InvalidJSON.md) — déduit : via UCA02 · Flux nominal — Connexion avec secret CA
> - … et 15 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
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
> [!check]- Tests — 0 explicite(s) · 22 déduit(s)
> - 🟡 [TestChannels_PUT_SwitchesChannel](../../docs/tests/adapters-in-rest/TestChannels_PUT_SwitchesChannel.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentityGuest_DisallowedReturnsForbidden](../../docs/tests/adapters-in-rest/TestHandleIdentityGuest_DisallowedReturnsForbidden.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentityGuest_RM49_PublicRead_NoSession](../../docs/tests/adapters-in-rest/TestHandleIdentityGuest_RM49_PublicRead_NoSession.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_CannotSignIn_Forbidden](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_CannotSignIn_Forbidden.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_EnrollError](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_EnrollError.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_InvalidJSON](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_InvalidJSON.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MethodNotAllowed.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_MissingFields](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_MissingFields.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_NoIdentityService](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_NoIdentityService.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_RoleFromCertificate](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_RoleFromCertificate.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestHandleIdentitySession_Success](../../docs/tests/adapters-in-rest/TestHandleIdentitySession_Success.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestRequireAuth_BlockchainWithValidToken](../../docs/tests/adapters-in-rest/TestRequireAuth_BlockchainWithValidToken.md) — déduit : via UCA07 · Flux nominal — Rôle connu depuis la connexion
> - 🟡 [TestCreate_NoTable_Rejected](../../docs/tests/domain-role/TestCreate_NoTable_Rejected.md) — déduit : via UCA05 · Flux nominal — Action autorisée
> - 🟡 [TestDelete_Success](../../docs/tests/domain-role/TestDelete_Success.md) — déduit : via UCA05 · Flux nominal — Action autorisée
> - 🟡 [TestHasPermission_DefaultTable](../../docs/tests/domain-role/TestHasPermission_DefaultTable.md) — déduit : via UCA05 · Flux nominal — Action autorisée
> - … et 7 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF06 — Consulter les assets possédés par l'utilisateur

- **Use cases couvrants** : UCA06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 84 déduit(s)
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
> - … et 69 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
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
> [!check]- Tests — 0 explicite(s) · 139 déduit(s)
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
> - … et 124 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF11 — Enregistrer un composant numérique sur la blockchain

- **Use cases couvrants** : UCCE03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 50 déduit(s)
> - 🟡 [TestFabricAdapter_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_SubmitError.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestFabricAdapter_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_Success.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestFabricBlockchain_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_SubmitError.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestFabricBlockchain_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_Success.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestIntegration_ModelService_SubmittedAsset](../../docs/tests/adapters-out-fabric/TestIntegration_ModelService_SubmittedAsset.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestIntegration_OwnerNotSigner_Rejected](../../docs/tests/adapters-out-fabric/TestIntegration_OwnerNotSigner_Rejected.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestIntegration_PendingOrReader_Rejected](../../docs/tests/adapters-out-fabric/TestIntegration_PendingOrReader_Rejected.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestIntegration_RM01_DuplicateHash_Rejected](../../docs/tests/adapters-out-fabric/TestIntegration_RM01_DuplicateHash_Rejected.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestIntegration_RM05_ParentMustExist](../../docs/tests/adapters-out-fabric/TestIntegration_RM05_ParentMustExist.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestIntegration_RM06_WriteOnce](../../docs/tests/adapters-out-fabric/TestIntegration_RM06_WriteOnce.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestIntegration_RM19_DraftRejected](../../docs/tests/adapters-out-fabric/TestIntegration_RM19_DraftRejected.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestIntegration_RM48_WritePermissionFollowsTable](../../docs/tests/adapters-out-fabric/TestIntegration_RM48_WritePermissionFollowsTable.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestIntegration_SharedIdentity_ChaincodeRefuses](../../docs/tests/adapters-out-fabric/TestIntegration_SharedIdentity_ChaincodeRefuses.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestIntegration_SharedIdentity_ClientRefuses](../../docs/tests/adapters-out-fabric/TestIntegration_SharedIdentity_ClientRefuses.md) — déduit : via UCCE03 · Flux erreur — Échec endorsement Fabric
> - … et 35 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF12 — Configurer les métadonnées d'un composant

- **Use cases couvrants** : UCCE02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 17 déduit(s)
> - 🟡 [TestFabricAdapter_GetModelRecord_EvaluateError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_EvaluateError.md) — déduit : via UCCE02 · Flux erreur — Asset introuvable
> - 🟡 [TestFabricAdapter_GetModelRecord_InvalidJSON_ReturnsEmptyModel](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_InvalidJSON_ReturnsEmptyModel.md) — déduit : via UCCE02 · Flux erreur — Asset introuvable
> - 🟡 [TestFabricAdapter_GetModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_Success.md) — déduit : via UCCE02 · Flux erreur — Asset introuvable
> - 🟡 [TestFabricBlockchain_GetModelRecord_EvaluateError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_EvaluateError.md) — déduit : via UCCE02 · Flux erreur — Asset introuvable
> - 🟡 [TestFabricBlockchain_GetModelRecord_InvalidJSON_ReturnsError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_InvalidJSON_ReturnsError.md) — déduit : via UCCE02 · Flux erreur — Asset introuvable
> - 🟡 [TestFabricBlockchain_GetModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_Success.md) — déduit : via UCCE02 · Flux erreur — Asset introuvable
> - 🟡 [TestFabricBlockchain_UpdateModelMetadata](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_UpdateModelMetadata.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestIntegration_GetModelRecord_NotFound](../../docs/tests/adapters-out-fabric/TestIntegration_GetModelRecord_NotFound.md) — déduit : via UCCE02 · Flux erreur — Asset introuvable
> - 🟡 [TestIntegration_ModelService_SubmittedAsset](../../docs/tests/adapters-out-fabric/TestIntegration_ModelService_SubmittedAsset.md) — déduit : via UCCE02 · Flux erreur — Asset introuvable
> - 🟡 [TestIntegration_RM06_WriteOnce](../../docs/tests/adapters-out-fabric/TestIntegration_RM06_WriteOnce.md) — déduit : via UCCE02 · Flux erreur — Asset introuvable
> - 🟡 [TestIntegration_StoreAndGet_UserSigned](../../docs/tests/adapters-out-fabric/TestIntegration_StoreAndGet_UserSigned.md) — déduit : via UCCE02 · Flux erreur — Asset introuvable
> - 🟡 [TestIntegration_UpdateMetadata_OwnerOnly](../../docs/tests/adapters-out-fabric/TestIntegration_UpdateMetadata_OwnerOnly.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestLicense_UnknownLicense_Rejected](../../docs/tests/domain-model/TestLicense_UnknownLicense_Rejected.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestUpdateAsset_Draft_StaysLocal](../../docs/tests/domain-model/TestUpdateAsset_Draft_StaysLocal.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - 🟡 [TestUpdateAsset_RM03_IncompatibleLicense_Rejected](../../docs/tests/domain-model/TestUpdateAsset_RM03_IncompatibleLicense_Rejected.md) — déduit : via UCCE02 · Flux nominal — Configuration réussie
> - … et 2 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF13 — Faire évoluer un composant (amélioration, dérivation, extension)

- **Use cases couvrants** : UCCE04, UCCE05

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 63 déduit(s)
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
> - 🟡 [TestIntegration_ModelService_SubmittedAsset](../../docs/tests/adapters-out-fabric/TestIntegration_ModelService_SubmittedAsset.md) — déduit : via UCCE04 · Flux erreur — Asset parent introuvable
> - 🟡 [TestIntegration_OwnerNotSigner_Rejected](../../docs/tests/adapters-out-fabric/TestIntegration_OwnerNotSigner_Rejected.md) — déduit : via UCCE04 · Flux erreur — Échec endorsement Fabric
> - 🟡 [TestIntegration_PendingOrReader_Rejected](../../docs/tests/adapters-out-fabric/TestIntegration_PendingOrReader_Rejected.md) — déduit : via UCCE04 · Flux erreur — Échec endorsement Fabric
> - … et 48 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF14 — Ajouter une interface à un composant existant

- **Use cases couvrants** : UCCE06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 110 déduit(s)
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
> - … et 95 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF15 — Vérifier l'unicité d'un composant (anti-plagiat SHA-256 + SCM > 50 %)

- **Use cases couvrants** : UCCE01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 139 déduit(s)
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
> - … et 124 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF16 — Vérifier la compatibilité de licence lors d'une dérivation

- **Use cases couvrants** : UCCE04, UCMOD01, UCMOD06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 181 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponent_AddInstance](../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponent_AddTwoInstancesSequentially](../../docs/tests/adapters-in-rest/TestComponent_AddTwoInstancesSequentially.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponent_DELETE_NoContent](../../docs/tests/adapters-in-rest/TestComponent_DELETE_NoContent.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - … et 166 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF60 — Supprimer un composant (masquage local des listes, ledger jamais modifié)

- **Use cases couvrants** : UCCE07

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 78 déduit(s)
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
> - … et 63 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Composants — Lecture

#### EF17 — Rechercher et filtrer les composants disponibles sur le réseau

- **Use cases couvrants** : UCCL01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 84 déduit(s)
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
> - … et 69 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Composition (instances)

#### EF18 — Créer des liaisons entre interfaces compatibles de composants

- **Use cases couvrants** : UCAM01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 39 déduit(s)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCAM01 · Flux — Liaison devenue incompatible après modification
> - 🟡 [TestDomainErrors_AssemblyRules_MappedToHTTP](../../docs/tests/adapters-in-rest/TestDomainErrors_AssemblyRules_MappedToHTTP.md) — déduit : via UCAM01 · Flux alternatif — Liaison via slot virtuel (UCAM03)
> - 🟡 [TestInterface_DELETE](../../docs/tests/adapters-in-rest/TestInterface_DELETE.md) — déduit : via UCAM01 · Flux — Liaison devenue incompatible après modification
> - 🟡 [TestInterface_PATCH_NotFound](../../docs/tests/adapters-in-rest/TestInterface_PATCH_NotFound.md) — déduit : via UCAM01 · Flux — Liaison devenue incompatible après modification
> - 🟡 [TestInterface_PATCH_OK](../../docs/tests/adapters-in-rest/TestInterface_PATCH_OK.md) — déduit : via UCAM01 · Flux — Liaison devenue incompatible après modification
> - 🟡 [TestUpdateConnection](../../docs/tests/adapters-out-localstorage/TestUpdateConnection.md) — déduit : via UCAM01 · Flux — Liaison devenue incompatible après modification
> - 🟡 [TestUpdateConnection_NotFound](../../docs/tests/adapters-out-localstorage/TestUpdateConnection_NotFound.md) — déduit : via UCAM01 · Flux — Liaison devenue incompatible après modification
> - 🟡 [TestAddAssemblyLink](../../docs/tests/domain-model/TestAddAssemblyLink.md) — déduit : via UCAM01 · Flux nominal A — Liaison directe
> - 🟡 [TestAddAssemblyLink_MissingStore](../../docs/tests/domain-model/TestAddAssemblyLink_MissingStore.md) — déduit : via UCAM01 · Flux nominal A — Liaison directe
> - 🟡 [TestAddAssemblyLink_RM09_InterfaceUsedOnce](../../docs/tests/domain-model/TestAddAssemblyLink_RM09_InterfaceUsedOnce.md) — déduit : via UCAM01 · Flux nominal A — Liaison directe
> - 🟡 [TestAddAssemblyLink_RM09_NoInstance_InterfaceUsedOnce](../../docs/tests/domain-model/TestAddAssemblyLink_RM09_NoInstance_InterfaceUsedOnce.md) — déduit : via UCAM01 · Flux nominal A — Liaison directe
> - 🟡 [TestAddAssemblyLink_RM09_SameInterfaceOtherInstance_OK](../../docs/tests/domain-model/TestAddAssemblyLink_RM09_SameInterfaceOtherInstance_OK.md) — déduit : via UCAM01 · Flux nominal A — Liaison directe
> - 🟡 [TestAddAssemblyLink_RM09_SameInterfaceTwoInstances_OK](../../docs/tests/domain-model/TestAddAssemblyLink_RM09_SameInterfaceTwoInstances_OK.md) — déduit : via UCAM01 · Flux nominal A — Liaison directe
> - 🟡 [TestAddAssemblyLink_RM09_SelfLink_Rejected](../../docs/tests/domain-model/TestAddAssemblyLink_RM09_SelfLink_Rejected.md) — déduit : via UCAM01 · Flux nominal A — Liaison directe
> - 🟡 [TestAddAssemblyLink_RM10_IncompatibleInterfaces_Rejected](../../docs/tests/domain-model/TestAddAssemblyLink_RM10_IncompatibleInterfaces_Rejected.md) — déduit : via UCAM01 · Flux nominal A — Liaison directe
> - … et 24 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF19 — Visualiser les interfaces physiques d'un composant

- **Use cases couvrants** : UCAM02

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 101 déduit(s)
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
> - … et 86 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF20 — Définir une interface sur un composant

- **Use cases couvrants** : UCAM03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 125 déduit(s)
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
> - … et 110 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
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
> [!check]- Tests — 0 explicite(s) · 59 déduit(s)
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCAM05 · Flux nominal — Conversion réussie
> - 🟡 [TestFabricAdapter_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_SubmitError.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestFabricAdapter_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_Success.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestFabricBlockchain_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_SubmitError.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestFabricBlockchain_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_Success.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestIntegration_ModelService_SubmittedAsset](../../docs/tests/adapters-out-fabric/TestIntegration_ModelService_SubmittedAsset.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestIntegration_OwnerNotSigner_Rejected](../../docs/tests/adapters-out-fabric/TestIntegration_OwnerNotSigner_Rejected.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestIntegration_PendingOrReader_Rejected](../../docs/tests/adapters-out-fabric/TestIntegration_PendingOrReader_Rejected.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestIntegration_RM01_DuplicateHash_Rejected](../../docs/tests/adapters-out-fabric/TestIntegration_RM01_DuplicateHash_Rejected.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestIntegration_RM05_ParentMustExist](../../docs/tests/adapters-out-fabric/TestIntegration_RM05_ParentMustExist.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - 🟡 [TestIntegration_RM06_WriteOnce](../../docs/tests/adapters-out-fabric/TestIntegration_RM06_WriteOnce.md) — déduit : via UCAM05 · Flux erreur — Erreur blockchain à la soumission
> - … et 44 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF23 — Choisir un asset d'accroche (fastener) pour une liaison

- **Use cases couvrants** : UCAM07

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 97 déduit(s)
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
> - … et 82 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF24 — Retirer une instance de composant d'un module (avec cascade des connexions)

- **Use cases couvrants** : UCAM08

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 14 déduit(s)
> - 🟡 [TestFabricAdapter_GetModelRecord_EvaluateError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_EvaluateError.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestFabricAdapter_GetModelRecord_InvalidJSON_ReturnsEmptyModel](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_InvalidJSON_ReturnsEmptyModel.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestFabricAdapter_GetModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_GetModelRecord_Success.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestFabricBlockchain_GetModelRecord_EvaluateError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_EvaluateError.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestFabricBlockchain_GetModelRecord_InvalidJSON_ReturnsError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_InvalidJSON_ReturnsError.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestFabricBlockchain_GetModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_GetModelRecord_Success.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestIntegration_GetModelRecord_NotFound](../../docs/tests/adapters-out-fabric/TestIntegration_GetModelRecord_NotFound.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestIntegration_ModelService_SubmittedAsset](../../docs/tests/adapters-out-fabric/TestIntegration_ModelService_SubmittedAsset.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestIntegration_RM06_WriteOnce](../../docs/tests/adapters-out-fabric/TestIntegration_RM06_WriteOnce.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestIntegration_StoreAndGet_UserSigned](../../docs/tests/adapters-out-fabric/TestIntegration_StoreAndGet_UserSigned.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestIntegration_UpdateMetadata_OwnerOnly](../../docs/tests/adapters-out-fabric/TestIntegration_UpdateMetadata_OwnerOnly.md) — déduit : via UCAM08 · Flux erreur — Module non trouvé
> - 🟡 [TestModuleComposition_RM19_Submitted_Rejected](../../docs/tests/domain-model/TestModuleComposition_RM19_Submitted_Rejected.md) — déduit : via UCAM08 · Flux nominal — Retrait sans liaisons actives
> - 🟡 [TestWorkspace_AddRemoveInstance](../../docs/tests/domain-model/TestWorkspace_AddRemoveInstance.md) — déduit : via UCAM08 · Flux nominal — Retrait sans liaisons actives
> - 🟡 [TestWorkspace_NilStore_ReturnsError](../../docs/tests/domain-model/TestWorkspace_NilStore_ReturnsError.md) — déduit : via UCAM08 · Flux nominal — Retrait sans liaisons actives
<!-- tests-obsidian:end -->

#### EF25 — Garantir un slot virtuel disponible sur chaque asset

- **Use cases couvrants** : UCAM03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 125 déduit(s)
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
> - … et 110 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
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
> [!check]- Tests — 0 explicite(s) · 124 déduit(s)
> - 🟡 [TestAPIIntegration_AssetLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_AssetLifecycle.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestAPIIntegration_ConcurrentAssetCreation](../../docs/tests/adapters-in-rest/TestAPIIntegration_ConcurrentAssetCreation.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestAPIIntegration_InterfaceLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_InterfaceLifecycle.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestAPIIntegration_ModuleLifecycle](../../docs/tests/adapters-in-rest/TestAPIIntegration_ModuleLifecycle.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponentInterfaces_GET](../../docs/tests/adapters-in-rest/TestComponentInterfaces_GET.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponentInterfaces_POST_BadJSON](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_BadJSON.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponentInterfaces_POST_Created](../../docs/tests/adapters-in-rest/TestComponentInterfaces_POST_Created.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponentTree_GET](../../docs/tests/adapters-in-rest/TestComponentTree_GET.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponentTree_MethodNotAllowed](../../docs/tests/adapters-in-rest/TestComponentTree_MethodNotAllowed.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponent_AddAssembly](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponent_AddAssembly_MissingConnectionID](../../docs/tests/adapters-in-rest/TestComponent_AddAssembly_MissingConnectionID.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponent_AddInstance](../../docs/tests/adapters-in-rest/TestComponent_AddInstance.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponent_AddInstance_MissingAssetID](../../docs/tests/adapters-in-rest/TestComponent_AddInstance_MissingAssetID.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponent_AddTwoInstancesSequentially](../../docs/tests/adapters-in-rest/TestComponent_AddTwoInstancesSequentially.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - 🟡 [TestComponent_DELETE_NoContent](../../docs/tests/adapters-in-rest/TestComponent_DELETE_NoContent.md) — déduit : via UCMOD01 · Flux alternatif — Dérivation d'un Module existant (reprise sans reconstruction)
> - … et 109 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF27 — Soumettre un module à la blockchain (ModuleVersion immuable)

- **Use cases couvrants** : UCMOD06

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 36 déduit(s)
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestFabricAdapter_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_SubmitError.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestFabricAdapter_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricAdapter_StoreModelRecord_Success.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_ReadOnly_Rejected.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestFabricBlockchain_StoreModelRecord_SubmitError](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_SubmitError.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestFabricBlockchain_StoreModelRecord_Success](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_StoreModelRecord_Success.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestIntegration_ModelService_SubmittedAsset](../../docs/tests/adapters-out-fabric/TestIntegration_ModelService_SubmittedAsset.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestIntegration_OwnerNotSigner_Rejected](../../docs/tests/adapters-out-fabric/TestIntegration_OwnerNotSigner_Rejected.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestIntegration_PendingOrReader_Rejected](../../docs/tests/adapters-out-fabric/TestIntegration_PendingOrReader_Rejected.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestIntegration_RM01_DuplicateHash_Rejected](../../docs/tests/adapters-out-fabric/TestIntegration_RM01_DuplicateHash_Rejected.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestIntegration_RM05_ParentMustExist](../../docs/tests/adapters-out-fabric/TestIntegration_RM05_ParentMustExist.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - 🟡 [TestIntegration_RM06_WriteOnce](../../docs/tests/adapters-out-fabric/TestIntegration_RM06_WriteOnce.md) — déduit : via UCMOD06 · Flux nominal — Soumission réussie
> - … et 21 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF28 — Ajouter un module existant à l'espace de travail

- **Use cases couvrants** : UCMOD02 — #incoherence : cette ligne référençait aussi UCMOD03 et UCMOD05 ; UCMOD03 ne couvre pas cet EF (il documente la modification des métadonnées d'un module, sans rapport avec l'ajout d'un module existant comme instance) et UCMOD05 ne correspond à aucun fichier existant dans `specs/1-Expression/UCMOD-Module/` ni `specs/2-Analyse/UCMOD-Module/` — origine de ces deux références à clarifier avant de les retirer définitivement

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 30 déduit(s)
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestFabricBlockchain_UpdateModelMetadata](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_UpdateModelMetadata.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestIntegration_UpdateMetadata_OwnerOnly](../../docs/tests/adapters-out-fabric/TestIntegration_UpdateMetadata_OwnerOnly.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestAddAssemblyToModule_RM51_ConnectionNotInModule_Rejected](../../docs/tests/domain-model/TestAddAssemblyToModule_RM51_ConnectionNotInModule_Rejected.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestAddAssemblyToModule_RM51_InstanceLink_OK](../../docs/tests/domain-model/TestAddAssemblyToModule_RM51_InstanceLink_OK.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestAddAssetToWorkspace_R21_Concurrent_NoLostUpdate](../../docs/tests/domain-model/TestAddAssetToWorkspace_R21_Concurrent_NoLostUpdate.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestAddAssetToWorkspace_RM52_Cycle_Rejected](../../docs/tests/domain-model/TestAddAssetToWorkspace_RM52_Cycle_Rejected.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestCreateModule_ParentID_CopiesComposition](../../docs/tests/domain-model/TestCreateModule_ParentID_CopiesComposition.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestGetModuleInterfaces_InternalConnection_ScopedToInstance](../../docs/tests/domain-model/TestGetModuleInterfaces_InternalConnection_ScopedToInstance.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestGetModuleInterfaces_InternalConnectionsConsumed](../../docs/tests/domain-model/TestGetModuleInterfaces_InternalConnectionsConsumed.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestGetModuleInterfaces_MultipleInstancesSameAsset_AllExposed](../../docs/tests/domain-model/TestGetModuleInterfaces_MultipleInstancesSameAsset_AllExposed.md) — déduit : via UCMOD02 · Flux nominal — Module ajouté (première instance)
> - 🟡 [TestLicense_UnknownLicense_Rejected](../../docs/tests/domain-model/TestLicense_UnknownLicense_Rejected.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - … et 15 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF29 — Visualiser la composition d'un module

- **Use cases couvrants** : UCMOD04

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 10 déduit(s)
> - 🟡 [TestLegacyModulesAlias_GET_DelegatesToComponents](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_GET_DelegatesToComponents.md) — déduit : via UCMOD04 · Flux nominal — Composition du module retournée
> - 🟡 [TestLegacyModulesAlias_ListRoute](../../docs/tests/adapters-in-rest/TestLegacyModulesAlias_ListRoute.md) — déduit : via UCMOD04 · Flux nominal — Composition du module retournée
> - 🟡 [TestModules_RegenerateThumbnail_OK](../../docs/tests/adapters-in-rest/TestModules_RegenerateThumbnail_OK.md) — déduit : via UCMOD04 · Flux nominal — Composition du module retournée
> - 🟡 [TestModules_Verify_OK](../../docs/tests/adapters-in-rest/TestModules_Verify_OK.md) — déduit : via UCMOD04 · Flux nominal — Composition du module retournée
> - 🟡 [TestGetModuleInterfaces_InternalConnection_ScopedToInstance](../../docs/tests/domain-model/TestGetModuleInterfaces_InternalConnection_ScopedToInstance.md) — déduit : via UCMOD04 · Flux alternatif — Consultation des interfaces exposées
> - 🟡 [TestGetModuleInterfaces_InternalConnectionsConsumed](../../docs/tests/domain-model/TestGetModuleInterfaces_InternalConnectionsConsumed.md) — déduit : via UCMOD04 · Flux alternatif — Consultation des interfaces exposées
> - 🟡 [TestGetModuleInterfaces_MultipleInstancesSameAsset_AllExposed](../../docs/tests/domain-model/TestGetModuleInterfaces_MultipleInstancesSameAsset_AllExposed.md) — déduit : via UCMOD04 · Flux alternatif — Consultation des interfaces exposées
> - 🟡 [TestGetModuleInterfaces_NilStore](../../docs/tests/domain-model/TestGetModuleInterfaces_NilStore.md) — déduit : via UCMOD04 · Flux alternatif — Consultation des interfaces exposées
> - 🟡 [TestGetModuleInterfaces_RM52_LegacyCycle_Terminates](../../docs/tests/domain-model/TestGetModuleInterfaces_RM52_LegacyCycle_Terminates.md) — déduit : via UCMOD04 · Flux alternatif — Consultation des interfaces exposées
> - 🟡 [TestGetModuleInterfaces_SimpleComponent](../../docs/tests/domain-model/TestGetModuleInterfaces_SimpleComponent.md) — déduit : via UCMOD04 · Flux alternatif — Consultation des interfaces exposées
<!-- tests-obsidian:end -->

#### EF61 — Modifier les métadonnées d'un module (nom, description, licence, tags, liens)

- **Use cases couvrants** : UCMOD03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 7 déduit(s)
> - 🟡 [TestFabricBlockchain_UpdateModelMetadata](../../docs/tests/adapters-out-fabric/TestFabricBlockchain_UpdateModelMetadata.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestIntegration_UpdateMetadata_OwnerOnly](../../docs/tests/adapters-out-fabric/TestIntegration_UpdateMetadata_OwnerOnly.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestLicense_UnknownLicense_Rejected](../../docs/tests/domain-model/TestLicense_UnknownLicense_Rejected.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestUpdateAsset_Draft_StaysLocal](../../docs/tests/domain-model/TestUpdateAsset_Draft_StaysLocal.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestUpdateAsset_RM03_IncompatibleLicense_Rejected](../../docs/tests/domain-model/TestUpdateAsset_RM03_IncompatibleLicense_Rejected.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestUpdateAsset_RM19_SubmittedLicense_Rejected](../../docs/tests/domain-model/TestUpdateAsset_RM19_SubmittedLicense_Rejected.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
> - 🟡 [TestUpdateAsset_SubmittedMetadata_UpdatedOnLedger](../../docs/tests/domain-model/TestUpdateAsset_SubmittedMetadata_UpdatedOnLedger.md) — déduit : via UCMOD03 · Flux nominal — Renommage du module
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
> [!check]- Tests — 0 explicite(s) · 78 déduit(s)
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
> - … et 63 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
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
> [!check]- Tests — 0 explicite(s) · 78 déduit(s)
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
> - … et 63 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Recherche

#### EF39 — Rechercher un asset par référence ou filtre

- **Use cases couvrants** : UCREC01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 82 déduit(s)
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
> - … et 67 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
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
> - 🟡 [TestIntegration_VerifyAndList](../../docs/tests/adapters-out-fabric/TestIntegration_VerifyAndList.md) — déduit : via UCREC02 · Flux erreur — Blockchain indisponible
<!-- tests-obsidian:end -->

#### EF41 — Consulter l'arbre de versions d'un composant

- **Use cases couvrants** : UCREC03

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 80 déduit(s)
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
> - … et 65 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
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
> [!check]- Tests — 0 explicite(s) · 78 déduit(s)
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
> - … et 63 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### EF47 — Gérer les versions SCM d'un modèle 3D

- **Use cases couvrants** : UCAUT04

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 78 déduit(s)
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
> - … et 63 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Interface Graphique, Paramètres, Documentation

EF48–EF54 : retirés — use cases 100 % frontend (UCIG, UCPAR, UCDOC), sans contrat REST propre à `myr`, désormais spécifiés dans le dépôt GUI externe

### Développement autour de MYR

#### EF55 — Exposer une API REST pour les intégrations tierces

- **Use cases couvrants** : UCDEV01

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!check]- Tests — 0 explicite(s) · 90 déduit(s)
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
> - … et 75 autre(s) : voir la [matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
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
