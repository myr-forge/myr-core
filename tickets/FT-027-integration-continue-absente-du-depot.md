---
id: FT-027
titre: "Intégration continue absente du dépôt"
type: incoherence
statut: ouvert
severite: majeure
detecte: 2026-10-02
maj: 2026-10-02
composants: [ci]
uc: []
rm: []
enf: [ENF18, ENF19, ENF25]
tags:
  - ticket
  - ticket/incoherence
  - statut/ouvert
  - severite/majeure
  - enf/ENF18
  - enf/ENF19
  - enf/ENF25
---
# FT-027 — Intégration continue absente du dépôt

> **Incohérence documentaire** · sévérité **majeure** · statut **Ouvert** · détecté le 2026-10-02
> Source : Lecture du dépôt

## Constat

ADR-10 décrit une intégration continue GitHub Actions (`.github/workflows/ci.yml` : compilation, `go test`, `go vet`, `make ci`, tests d'intégration, déploiement SSH). Ce fichier n'existe pas dans le dépôt.

## Cause

Pipeline décidé mais jamais ajouté.

## Impact

Les tests rouges (FT-019), l'isolation hexagonale (ENF18) et la conformité des licences (ENF25) ne sont vérifiés que si quelqu'un lance `make test` / `make ci` à la main.

## Preuves

Aucun dossier `.github/workflows` dans le dépôt.

## Piste de correction (à valider par le PO)

Ajouter `.github/workflows/ci.yml` conforme à ADR-10 (secret SSH pour le déploiement), en commençant par build, `go vet`, `go test`, `make ci`.

## Critères de clôture

- [ ] Pipeline exécuté sur chaque push et pull request
- [ ] Échec bloquant si tests ou `make ci` échouent

## Liens

- **Exigences non fonctionnelles** : `ENF18`, `ENF19`, `ENF25` (tags `enf/…`)
- **Specs** : [Conception_intro — ADR-10](../specs/3-Conception/Conception_intro.md)
- **Code** : `Makefile` · `scripts/ci/check-domain-imports.sh`
- **Tickets liés** : [FT-019 — Deux tests rouges dans domain/model](FT-019-deux-tests-rouges-dans-domain-model.md)

## Historique

- 2026-10-02 — Ticket créé
