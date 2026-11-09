---
id: FT-019
titre: "Deux tests rouges dans domain/model"
type: anomalie
statut: en-cours
severite: majeure
detecte: 2026-10-02
maj: 2026-11-09
composants: [domain/model]
uc: [UCREC03, UCCE01]
rm: []
enf: [ENF19]
tags:
  - ticket
  - ticket/anomalie
  - statut/en-cours
  - severite/majeure
  - domaine/model
  - uc/UCREC03
  - uc/UCCE01
  - enf/ENF19
---
# FT-019 — Deux tests rouges dans domain/model

> **Anomalie** · sévérité **majeure** · statut **En cours** · détecté le 2026-10-02
> Source : `go test ./...` au commit 2aa69c1

## Constat

`go test ./...` échoue sur deux tests de `domain/model/tests` :
- `TestGetChildren_Success` — « attendu 1 enfant, got 0 » (`service_test.go:1557`) ;
- `TestAddFull_NilDraftStore_Rejected` — « AddFull sans DraftStore doit retourner une erreur » (`service_test.go:2401`).

Ces échecs sont antérieurs aux travaux récents (une note de travail en signalait 5, dont ceux-ci).

## Cause

À analyser : changement de comportement de `GetChildren` et de `AddFull` (création directe soumise sans passer par le brouillon) non répercuté dans les tests, ou régression.

## Impact

La suite de tests n'est pas verte : toute nouvelle régression du domaine passe inaperçue au milieu des échecs connus.

## Preuves

Sortie de `go test ./domain/model/tests -run "TestGetChildren_Success|TestAddFull_NilDraftStore_Rejected"`.

## Piste de correction (à valider par le PO)

Pour chaque test : décider si le comportement attendu (specs UCREC03, UCCE01) ou le code fait foi, puis corriger l'un ou l'autre.

## Critères de clôture

- [ ] `go test ./...` vert

## Liens

- **Use cases** : [UCREC03](../specs/2-Analyse/UCREC-Recherche/UCREC03.md) · [UCCE01](../specs/2-Analyse/UCCE-Composant_Ecriture/UCCE01.md)
- **Exigences non fonctionnelles** : `ENF19` (tags `enf/…`)
- **Code** : `domain/model/tests/service_test.go:1557` · `domain/model/tests/service_test.go:2401`
- **Fonctions** : [ModelService.GetChildren](../docs/code/fonctions/model.ModelService.GetChildren.md) · [ModelService.AddFull](../docs/code/fonctions/model.ModelService.AddFull.md)
- **Tickets liés** : [FT-027 — Intégration continue absente du dépôt](FT-027-integration-continue-absente-du-depot.md)

## Historique

- 2026-10-02 — Ticket créé ; échecs reproduits au commit `2aa69c1`
- 2026-11-09 — `TestAddFull_NilDraftStore_Rejected` corrigé (commit `51f12f3`) : depuis que toute création passe par un brouillon (`9cdfb54`), le test injectait par erreur un stockage de brouillons et ne testait plus son absence ; il le fait de nouveau. Reste `TestGetChildren_Success` : `GetChildren` ne lit que le ledger, alors que l'enfant créé par `AddFull` est un brouillon — comportement attendu à trancher (brouillons inclus ou non)
- 2026-11-09 — Corrigé sur `feature/ledger-invariants` (commit `03c9ad8`) : décision du PO, `GetChildren` ne liste que les dérivés soumis (UCREC03) ; le test soumet le parent et l'enfant. `go test ./...` passe en entier
