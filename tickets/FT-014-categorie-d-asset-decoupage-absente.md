---
id: FT-014
titre: "Catégorie d'asset decoupage absente"
type: ecart
statut: ouvert
severite: majeure
detecte: 2026-10-02
maj: 2026-10-02
composants: [domain/model]
uc: [UCAM05, UCAM09]
rm: [RM02]
enf: []
tags:
  - ticket
  - ticket/ecart
  - statut/ouvert
  - severite/majeure
  - domaine/model
  - uc/UCAM05
  - uc/UCAM09
  - rm/RM02
---
# FT-014 — Catégorie d'asset decoupage absente

> **Écart spec ↔ code** · sévérité **majeure** · statut **Ouvert** · détecté le 2026-10-02
> Source : Écart E1 (Analyse_des_besoins.md)

## Constat

RM02 définit 8 catégories d'asset ; le domaine n'en déclare que 7 (`base`, `amelioration`, `variation`, `adaptation`, `derivation`, `extension`, `regression`). `decoupage` manque.

## Cause

Catégorie ajoutée aux specs après l'écriture de l'entité.

## Impact

Bloque la transformation d'un composant en module (UCAM05) et la décomposition assistée (UCAM09).

## Preuves

`domain/model/entity.go:14-20`.

## Piste de correction (à valider par le PO)

Ajouter la constante `decoupage` et l'accepter partout où les catégories sont validées (domaine, CLI, REST, chaincode).

## Critères de clôture

- [ ] Catégorie `decoupage` acceptée de bout en bout
- [ ] Test de création d'un asset `decoupage`

## Liens

- **Use cases** : [UCAM05](../specs/2-Analyse/UCAM-Assemblage_Module/UCAM05.md) · [UCAM09](../specs/1-Expression/UCAM-Assemblage_Module/UCAM09.md)
- **Règles métier** : `RM02` (tags `rm/…`)
- **Code** : `domain/model/entity.go:14`
- **Fonctions** : [ModelService.AddFull](../docs/code/fonctions/model.ModelService.AddFull.md)
- **Tickets liés** : [FT-021 — Identifiant d'un asset décomposé en module](FT-021-identifiant-d-un-asset-decompose-en-module.md) · [FT-023 — Décomposition STEP assistée : prérequis non réunis](FT-023-decomposition-step-assistee-prerequis-non-reunis.md)

## Historique

- 2026-10-02 — Ticket créé ; constat vérifié au commit `2aa69c1`
