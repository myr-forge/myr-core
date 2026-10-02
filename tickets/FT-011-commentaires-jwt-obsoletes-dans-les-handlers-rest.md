---
id: FT-011
titre: "Commentaires JWT obsolètes dans les handlers REST"
type: dette
statut: ouvert
severite: mineure
detecte: 2026-10-02
maj: 2026-10-02
composants: [adapters/in/rest]
uc: []
rm: []
enf: []
tags:
  - ticket
  - ticket/dette
  - statut/ouvert
  - severite/mineure
---
# FT-011 — Commentaires JWT obsolètes dans les handlers REST

> **Dette technique** · sévérité **mineure** · statut **Ouvert** · détecté le 2026-10-02
> Source : Lecture du code

## Constat

Plusieurs commentaires décrivent une authentification « JWT Bearer » qui n'a jamais existé : l'authentification REST repose uniquement sur un jeton de session opaque (`X-Myr-Token`).

## Cause

Reliquat d'une conception abandonnée.

## Impact

Induit en erreur les lecteurs du code (humains et agents) sur le mécanisme d'authentification.

## Preuves

`handlers.go:295` (« 1) JWT Bearer (nouveau) »), `handlers.go:1707`, `handlers_network.go:147`.

## Piste de correction (à valider par le PO)

Réécrire ces commentaires pour décrire le jeton de session opaque.

## Critères de clôture

- [ ] Plus aucune mention de JWT dans `adapters/in/rest`

## Liens

- **Specs** : [DC_D1_Auth_Identity](../specs/3-Conception/DC_D1_Auth_Identity.md)
- **Code** : `adapters/in/rest/handlers.go:295` · `adapters/in/rest/handlers.go:1707` · `adapters/in/rest/handlers_network.go:147`

## Historique

- 2026-10-02 — Ticket créé
