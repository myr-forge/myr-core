---
id: FT-033
titre: "UCMOD05 référencé mais inexistant"
type: incoherence
statut: a-trancher
severite: mineure
detecte: 2026-10-02
maj: 2026-10-02
composants: [specs]
uc: [UCMOD02, UCMOD03]
rm: []
enf: []
tags:
  - ticket
  - ticket/incoherence
  - statut/a-trancher
  - severite/mineure
  - domaine/model
  - uc/UCMOD02
  - uc/UCMOD03
---
# FT-033 — UCMOD05 référencé mais inexistant

> **Incohérence documentaire** · sévérité **mineure** · statut **À trancher (PO)** · détecté le 2026-10-02
> Source : Annotation #incoherence de Matrice_Tracabilite.md (EF28) ; roadmap § Modules

## Constat

La matrice de traçabilité (EF28) et la roadmap (V1 : « UCMOD05 — ajouter un module existant via plugin navigateur ») citent UCMOD05, qui ne correspond à aucun fichier de specs. EF28 citait aussi UCMOD03, sans rapport avec cette exigence.

## Cause

Origine non retrouvée (use case supprimé ou jamais rédigé).

## Impact

Exigence EF28 partiellement orpheline ; numérotation trompeuse.

## Preuves

Annotation sur la ligne EF28 de `specs/1-Expression/Matrice_Tracabilite.md`.

## Piste de correction (à valider par le PO)

Décider : rédiger UCMOD05 (ajout via plugin navigateur) ou retirer les références.

## Critères de clôture

- [ ] Références à UCMOD05 résolues
- [ ] Annotation levée

## Liens

- **Use cases** : [UCMOD02 (analyse)](../specs/2-Analyse/UCMOD-Module/UCMOD02.md) · [UCMOD02 (expression)](../specs/1-Expression/UCMOD-Module/UCMOD02.md) · [UCMOD03 (analyse)](../specs/2-Analyse/UCMOD-Module/UCMOD03.md) · [UCMOD03 (expression)](../specs/1-Expression/UCMOD-Module/UCMOD03.md)
- **Specs** : [Matrice_Tracabilite](../specs/1-Expression/Matrice_Tracabilite.md) · [roadmap_dev](../specs/roadmap_dev.md)

## Historique

- 2026-10-02 — Ticket créé
