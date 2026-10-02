---
id: FT-022
titre: "Date de retrait des routes legacy /api/modules"
type: decision
statut: a-trancher
severite: mineure
detecte: 2026-08-12
maj: 2026-10-02
composants: [adapters/in/rest]
uc: [UCMOD01, UCMOD04, UCMOD06]
rm: []
enf: []
tags:
  - ticket
  - ticket/decision
  - statut/a-trancher
  - severite/mineure
  - domaine/model
  - uc/UCMOD01
  - uc/UCMOD04
  - uc/UCMOD06
---
# FT-022 — Date de retrait des routes legacy /api/modules

> **Décision à prendre** · sévérité **mineure** · statut **À trancher (PO)** · détecté le 2026-08-12
> Source : ADR-11 ; roadmap § Fusion composant/module

## Constat

Depuis la fusion composant/module, les routes `/api/modules/*` sont redirigées vers les handlers `/api/components/*` (`legacyModulesAlias`), en réponse au format `componentDTO`. La date de leur retrait n'est pas fixée.

## Cause

Fenêtre de transition à coordonner avec myr-web.

## Impact

Deux familles de routes à maintenir et documenter ; les specs citent encore massivement `/api/modules/*`.

## Preuves

`adapters/in/rest/server.go:25-31` et `:132-133`.

## Piste de correction (à valider par le PO)

Convenir d'une date avec myr-web, signaler la dépréciation (en-tête `Deprecation`), puis retirer les routes et mettre à jour les specs.

## Critères de clôture

- [ ] Date convenue
- [ ] Routes retirées ou dépréciation signalée
- [ ] Specs mises à jour

## Liens

- **Use cases** : [UCMOD01 (analyse)](../specs/2-Analyse/UCMOD-Module/UCMOD01.md) · [UCMOD01 (expression)](../specs/1-Expression/UCMOD-Module/UCMOD01.md) · [UCMOD04 (analyse)](../specs/2-Analyse/UCMOD-Module/UCMOD04.md) · [UCMOD04 (expression)](../specs/1-Expression/UCMOD-Module/UCMOD04.md) · [UCMOD06 (analyse)](../specs/2-Analyse/UCMOD-Module/UCMOD06.md) · [UCMOD06 (expression)](../specs/1-Expression/UCMOD-Module/UCMOD06.md)
- **Specs** : [API_REST](../specs/3-Conception/API_REST.md)
- **Code** : [adapters/in/rest/server.go:25](../adapters/in/rest/server.go) · [adapters/in/rest/server.go:132](../adapters/in/rest/server.go)
- **Tickets liés** : [FT-021 — Identifiant d'un asset décomposé en module](FT-021-identifiant-d-un-asset-decompose-en-module.md) · [FT-031 — Les specs citent des commandes et routes qui n'existent pas](FT-031-les-specs-citent-des-commandes-et-routes.md)

## Historique

- 2026-08-12 — Fenêtre de transition ouverte par ADR-11
- 2026-10-02 — Ticket créé
