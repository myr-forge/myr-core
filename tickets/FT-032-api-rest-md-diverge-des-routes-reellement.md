---
id: FT-032
titre: "API_REST.md diverge des routes réellement servies"
type: incoherence
statut: ouvert
severite: mineure
detecte: 2026-10-02
maj: 2026-10-02
composants: [specs, adapters/in/rest]
uc: [UCCE02, UCMOD03]
rm: []
enf: []
tags:
  - ticket
  - ticket/incoherence
  - statut/ouvert
  - severite/mineure
  - domaine/model
  - uc/UCCE02
  - uc/UCMOD03
---
# FT-032 — API_REST.md diverge des routes réellement servies

> **Incohérence documentaire** · sévérité **mineure** · statut **Ouvert** · détecté le 2026-10-02
> Source : Annotations #incoherence et #remarque de API_REST.md

## Constat

Le document de conception de l'API porte plusieurs écarts annotés :
- modification d'un asset documentée en `PUT`, alors que le handler ne répond qu'à `PATCH` (tableau corrigé, intention initiale à confirmer) ;
- routes `/api/licenses/check` et `/api/licenses/check-product` implémentées mais longtemps absentes du tableau ;
- routes `GET`/`POST /api/components/{id}/thumbnail` annotées « non routées » — elles sont désormais servies (`handlers.go:914-925`, commit `15a14f9`) : annotation devenue obsolète.

## Cause

Documentation tenue à la main en parallèle de la spec générée (`api/swagger.json`).

## Impact

Les clients de l'API (myr-web) peuvent se fier à une description fausse.

## Preuves

Annotations dans `specs/3-Conception/API_REST.md` (§4, §6, §7, §8).

## Piste de correction (à valider par le PO)

Trancher chaque écart (PO), mettre à jour les tableaux, lever les annotations ; régénérer `api/swagger.*` (`make docs-api`).

## Critères de clôture

- [ ] Plus d'annotation ouverte dans API_REST.md
- [ ] Tableaux cohérents avec `api/swagger.json`

## Liens

- **Use cases** : [UCCE02 (analyse)](../specs/2-Analyse/UCCE-Composant_Ecriture/UCCE02.md) · [UCCE02 (expression)](../specs/1-Expression/UCCE-Composant_Ecriture/UCCE02.md) · [UCMOD03 (analyse)](../specs/2-Analyse/UCMOD-Module/UCMOD03.md) · [UCMOD03 (expression)](../specs/1-Expression/UCMOD-Module/UCMOD03.md)
- **Specs** : [API_REST](../specs/3-Conception/API_REST.md)
- **Code** : [adapters/in/rest/handlers.go:914](../adapters/in/rest/handlers.go)
- **Tickets liés** : [FT-022 — Date de retrait des routes legacy /api/modules](FT-022-date-de-retrait-des-routes-legacy-api.md)

## Historique

- 2026-10-02 — Ticket créé
