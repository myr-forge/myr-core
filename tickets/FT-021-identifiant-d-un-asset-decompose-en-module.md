---
id: FT-021
titre: "Identifiant d'un asset décomposé en module"
type: decision
statut: a-trancher
severite: majeure
detecte: 2026-08-12
maj: 2026-10-02
composants: [domain/model]
uc: [UCAM05, UCAM09]
rm: [RM39]
enf: []
tags:
  - ticket
  - ticket/decision
  - statut/a-trancher
  - severite/majeure
  - domaine/model
  - uc/UCAM05
  - uc/UCAM09
  - rm/RM39
---
# FT-021 — Identifiant d'un asset décomposé en module

> **Décision à prendre** · sévérité **majeure** · statut **À trancher (PO)** · détecté le 2026-08-12
> Source : ADR-11 (Conception_intro.md), point ouvert ; roadmap § Fusion composant/module

## Constat

Quand un composant est transformé en module (décomposé en sous-pièces), deux sémantiques s'opposent : garder le même identifiant (proposition de myr-web) ou créer systématiquement un nouvel identifiant (ce que spécifient UCAM05, UCAM09 et RM39).

## Cause

Point laissé ouvert par ADR-11.

## Impact

Bloque l'implémentation finale de UCAM05/UCAM09 et le contrat avec myr-web.

## Preuves

ADR-11, « point ouvert PO ».

## Piste de correction (à valider par le PO)

Trancher entre les deux options en tenant compte de l'immuabilité (un asset soumis ne change pas — RM19) et de la filiation (RM39).

## Critères de clôture

- [ ] Décision consignée dans ADR-11
- [ ] UCAM05/UCAM09/RM39 alignés

## Liens

- **Use cases** : [UCAM05 (analyse)](../specs/2-Analyse/UCAM-Assemblage_Module/UCAM05.md) · [UCAM05 (expression)](../specs/1-Expression/UCAM-Assemblage_Module/UCAM05.md) · [UCAM09 (expression)](../specs/1-Expression/UCAM-Assemblage_Module/UCAM09.md)
- **Règles métier** : [RM39](../specs/1-Expression/Regles_Metier.md)
- **Specs** : [Conception_intro — ADR-11](../specs/3-Conception/Conception_intro.md)
- **Fonctions** : [ModelService.AddFull](../docs/code/fonctions/model.ModelService.AddFull.md) · [ModelService.AddAssetToWorkspace](../docs/code/fonctions/model.ModelService.AddAssetToWorkspace.md)
- **Tickets liés** : [FT-014 — Catégorie d'asset decoupage absente](FT-014-categorie-d-asset-decoupage-absente.md) · [FT-022 — Date de retrait des routes legacy /api/modules](FT-022-date-de-retrait-des-routes-legacy-api.md)

## Historique

- 2026-08-12 — Point ouvert par ADR-11
- 2026-10-02 — Ticket créé
