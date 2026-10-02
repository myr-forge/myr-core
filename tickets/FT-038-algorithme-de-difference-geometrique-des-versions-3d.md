---
id: FT-038
titre: "Algorithme de différence géométrique des versions 3D"
type: decision
statut: a-trancher
severite: mineure
detecte: 2026-10-02
maj: 2026-10-02
composants: [domain/model]
uc: [UCAUT04]
rm: []
enf: []
tags:
  - ticket
  - ticket/decision
  - statut/a-trancher
  - severite/mineure
  - domaine/model
  - uc/UCAUT04
---
# FT-038 — Algorithme de différence géométrique des versions 3D

> **Décision à prendre** · sévérité **mineure** · statut **À trancher (PO)** · détecté le 2026-10-02
> Source : DC_D9_Automatisation.md (question ouverte, annotation #remarque)

## Constat

La gestion de versions des modèles 3D (UCAUT04) suppose un calcul de différence géométrique (`GeometryDelta`). Non tranché : librairie Go, outil CAO externe appelé en sous-processus, ou hors périmètre v1 (diff des métadonnées et interfaces seulement). Piste notée : stocker un format paramétrique (arbre de construction) plutôt que le fichier 3D final.

## Cause

Question technique ouverte ; toute librairie serait une nouvelle technologie à valider.

## Impact

Bloque la conception détaillée d'UCAUT04 (post-V1).

## Preuves

`specs/3-Conception/DC_D9_Automatisation.md`, question ouverte pour le PO.

## Piste de correction (à valider par le PO)

Décision PO ; si une technologie est retenue, la documenter dans specs/3-Conception avant implémentation.

## Critères de clôture

- [ ] Décision consignée dans DC_D9_Automatisation.md

## Liens

- **Use cases** : [UCAUT04](../specs/2-Analyse/UCAUT-Automatisation/UCAUT04.md)
- **Specs** : [DC_D9_Automatisation](../specs/3-Conception/DC_D9_Automatisation.md)
- **Code** : aucun — la gestion de versions des modèles 3D (UCAUT04) n'est pas développée

## Historique

- 2026-10-02 — Ticket créé
