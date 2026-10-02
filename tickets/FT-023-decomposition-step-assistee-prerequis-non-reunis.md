---
id: FT-023
titre: "Décomposition STEP assistée : prérequis non réunis"
type: ecart
statut: differe
severite: mineure
detecte: 2026-10-02
maj: 2026-10-02
composants: [domain/model]
uc: [UCAM09]
rm: [RM39, RM40, RM41]
enf: []
tags:
  - ticket
  - ticket/ecart
  - statut/differe
  - severite/mineure
  - domaine/model
  - uc/UCAM09
  - rm/RM39
  - rm/RM40
  - rm/RM41
---
# FT-023 — Décomposition STEP assistée : prérequis non réunis

> **Écart spec ↔ code** · sévérité **mineure** · statut **Différé** · détecté le 2026-10-02
> Source : Écart E9 (Analyse_des_besoins.md) ; roadmap § Compléments Post-V1 ; DC_CLI_Model.md §8

## Constat

UCAM09 (proposer un découpage automatique d'un composant STEP) est mis de côté : aucune librairie Go mature ne couvre à la fois la lecture STEP AP214/AP242 et la détection géométrique de contacts, et `AssetInterface` n'a aucun repère géométrique (position, orientation, entité STEP d'origine).

## Cause

Absence de solution technique validée ; écarts E1 (FT-014) et E9.

## Impact

Fonctionnalité reportée après la V1.

## Preuves

`specs/3-Conception/DC_CLI_Model.md` §8 (point ouvert) ; commande `myr model decompose` citée par les specs mais inexistante.

## Piste de correction (à valider par le PO)

Spike limité dans le temps avant toute reprise ; périmètre réduit à une v1 « structure seule » (sous-pièces et filiation RM39, sans suggestion de connexions).

## Critères de clôture

- [ ] Spike réalisé et conclusions documentées
- [ ] Repère géométrique ajouté à `AssetInterface` si la piste est retenue

## Liens

- **Use cases** : [UCAM09](../specs/1-Expression/UCAM-Assemblage_Module/UCAM09.md)
- **Règles métier** : `RM39`, `RM40`, `RM41` (tags `rm/…`)
- **Specs** : [DC_CLI_Model](../specs/3-Conception/DC_CLI_Model.md)
- **Code** : `domain/model/entity.go`
- **Tickets liés** : [FT-014 — Catégorie d'asset decoupage absente](FT-014-categorie-d-asset-decoupage-absente.md)

## Historique

- 2026-10-02 — Ticket créé (différé post-V1)
