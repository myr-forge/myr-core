---
id: FT-015
titre: "Critère tag absent des interfaces"
type: ecart
statut: ouvert
severite: mineure
detecte: 2026-10-02
maj: 2026-10-02
composants: [domain/model]
uc: [UCAM01, UCREC02]
rm: [RM11]
enf: []
tags:
  - ticket
  - ticket/ecart
  - statut/ouvert
  - severite/mineure
  - domaine/model
  - uc/UCAM01
  - uc/UCREC02
  - rm/RM11
---
# FT-015 — Critère tag absent des interfaces

> **Écart spec ↔ code** · sévérité **mineure** · statut **Ouvert** · détecté le 2026-10-02
> Source : Écart E2 (Analyse_des_besoins.md)

## Constat

RM11 définit cinq critères de compatibilité entre interfaces, dont le tag (câble, vis…). `AssetInterface` n'a pas de champ `Tag`, donc ce critère n'est jamais évalué.

## Cause

Champ non ajouté à l'entité.

## Impact

Des interfaces de même catégorie et type mais de tags différents sont jugées compatibles à tort.

## Preuves

`domain/model/entity.go` : structure `AssetInterface` sans champ `Tag`.

## Piste de correction (à valider par le PO)

Ajouter `Tag` à `AssetInterface` et l'intégrer à `ifacesCompatible` (critère ignoré si non renseigné).

## Critères de clôture

- [ ] Champ `Tag` persisté et exposé
- [ ] `ifacesCompatible` applique le critère
- [ ] Test `TestIfacesCompatible_MismatchTag_Rejected`

## Liens

- **Use cases** : [UCAM01](../specs/2-Analyse/UCAM-Assemblage_Module/UCAM01.md) · [UCREC02](../specs/2-Analyse/UCREC-Recherche/UCREC02.md)
- **Règles métier** : `RM11` (tags `rm/…`)
- **Code** : `domain/model/entity.go`
- **Fonctions** : [ModelService.AddAssemblyLink](../docs/code/fonctions/model.ModelService.AddAssemblyLink.md)

## Historique

- 2026-10-02 — Ticket créé
