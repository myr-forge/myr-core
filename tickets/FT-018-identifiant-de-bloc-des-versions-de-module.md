---
id: FT-018
titre: "Identifiant de bloc des versions de module simulé"
type: anomalie
statut: ouvert
severite: majeure
detecte: 2026-10-02
maj: 2026-10-02
composants: [domain/model, adapters/out/fabric]
uc: [UCMOD06]
rm: [RM18]
enf: [ENF28]
tags:
  - ticket
  - ticket/anomalie
  - statut/ouvert
  - severite/majeure
  - domaine/model
  - uc/UCMOD06
  - rm/RM18
  - enf/ENF28
---
# FT-018 — Identifiant de bloc des versions de module simulé

> **Anomalie** · sévérité **majeure** · statut **Ouvert** · détecté le 2026-10-02
> Source : Lecture du code (création des tickets)

## Constat

À la soumission d'un module, la `ModuleVersion` reçoit un `BlockID` tiré au hasard (`P` suivi de 8 caractères hexadécimaux) par `simModuleBlockID()`, et non l'identifiant de la transaction ou du bloc Fabric réellement produit.

## Cause

`service.go:851` appelle `simModuleBlockID()` (`service.go:1116`) avant même l'écriture blockchain ; `StoreModelRecord` ne renvoie aucun identifiant de transaction.

## Impact

La traçabilité affichée (« inscrit au bloc X ») est fictive : impossible de retrouver la transaction d'une version sur le ledger.

## Preuves

`domain/model/service.go:850-851` et `:1116-1120`.

## Piste de correction (à valider par le PO)

Faire renvoyer l'identifiant de transaction par le port blockchain (`StoreModelRecord`) et l'enregistrer dans la `ModuleVersion` ; le simulateur ne doit subsister que dans les adapters de test.

## Critères de clôture

- [ ] `BlockID` issu de la transaction Fabric
- [ ] Plus d'identifiant simulé dans le domaine

## Liens

- **Use cases** : [UCMOD06](../specs/2-Analyse/UCMOD-Module/UCMOD06.md)
- **Règles métier** : `RM18` (tags `rm/…`)
- **Exigences non fonctionnelles** : `ENF28` (tags `enf/…`)
- **Specs** : [Sequence_soumission_module](../specs/3-Conception/Sequence_soumission_module.md)
- **Code** : `domain/model/service.go:851` · `domain/model/service.go:1116`
- **Fonctions** : [ModelService.SubmitModule](../docs/code/fonctions/model.ModelService.SubmitModule.md) · [BlockchainPort.StoreModelRecord](../docs/code/fonctions/model.BlockchainPort.StoreModelRecord.md)
- **Tickets liés** : [FT-012 — Un module soumis peut être soumis à nouveau sans fork](FT-012-un-module-soumis-peut-etre-soumis-a.md) · [FT-024 — Chaincode myrcc absent du canal de production](FT-024-chaincode-myrcc-absent-du-canal-de-production.md)

## Historique

- 2026-10-02 — Ticket créé ; constat vérifié au commit `2aa69c1`
