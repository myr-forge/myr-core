---
id: FT-012
titre: "Un module soumis peut être soumis à nouveau sans fork"
type: ecart
statut: ouvert
severite: majeure
detecte: 2026-10-02
maj: 2026-10-02
composants: [domain/model]
uc: [UCMOD06]
rm: [RM19, RM06]
enf: [ENF28]
tags:
  - ticket
  - ticket/ecart
  - statut/ouvert
  - severite/majeure
  - domaine/model
  - uc/UCMOD06
  - rm/RM19
  - rm/RM06
  - enf/ENF28
---
# FT-012 — Un module soumis peut être soumis à nouveau sans fork

> **Écart spec ↔ code** · sévérité **majeure** · statut **Ouvert** · détecté le 2026-10-02
> Source : Écart E5 (Analyse_des_besoins.md) ; roadmap § Bugs bloquants

## Constat

`SubmitModule` ne vérifie pas que le module est encore en brouillon : appelé sur un module déjà `submitted`, il ajoute une nouvelle `ModuleVersion` et réécrit l'enregistrement. RM19 impose qu'un asset soumis soit immuable et que toute évolution passe par un fork.

## Cause

`service.go:838-863` : contrôles limités à l'existence de l'asset, à la blockchain et à la présence d'au moins un assemblage (RM17).

## Impact

L'immuabilité promise par la soumission n'est pas garantie par le domaine ; un module publié peut changer de composition sous le même identifiant.

## Preuves

`domain/model/service.go:838-863` : aucune lecture de `m.Status` ni du drapeau brouillon retourné par `getAsset`.

## Piste de correction (à valider par le PO)

Refuser `SubmitModule` (et les mutations de composition : ajout/retrait d'instance, de liaison) sur un asset déjà soumis, avec une erreur métier dédiée ; vérifier le même invariant pour la soumission d'un composant.

## Critères de clôture

- [ ] Soumission d'un module déjà soumis refusée
- [ ] Mutations de composition refusées sur un module soumis
- [ ] Tests `TestSubmitModule_AlreadySubmitted_Rejected` et équivalents

## Liens

- **Use cases** : [UCMOD06](../specs/2-Analyse/UCMOD-Module/UCMOD06.md)
- **Règles métier** : `RM19`, `RM06` (tags `rm/…`)
- **Exigences non fonctionnelles** : `ENF28` (tags `enf/…`)
- **Specs** : [Sequence_soumission_module](../specs/3-Conception/Sequence_soumission_module.md)
- **Code** : `domain/model/service.go:838`
- **Fonctions** : [ModelService.SubmitModule](../docs/code/fonctions/model.ModelService.SubmitModule.md) · [ModelService.Submit](../docs/code/fonctions/model.ModelService.Submit.md) · [ModelService.AddAssetToWorkspace](../docs/code/fonctions/model.ModelService.AddAssetToWorkspace.md)
- **Tickets liés** : [FT-018 — Identifiant de bloc des versions de module simulé](FT-018-identifiant-de-bloc-des-versions-de-module.md)

## Historique

- 2026-10-02 — Ticket créé ; constat vérifié au commit `2aa69c1`
