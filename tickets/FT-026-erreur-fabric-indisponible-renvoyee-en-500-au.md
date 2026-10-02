---
id: FT-026
titre: "Erreur Fabric indisponible renvoyée en 500 au lieu de 503"
type: anomalie
statut: ouvert
severite: mineure
detecte: 2026-07-29
maj: 2026-10-02
composants: [adapters/in/rest, adapters/out/fabric]
uc: [UCCL01]
rm: []
enf: [ENF05]
tags:
  - ticket
  - ticket/anomalie
  - statut/ouvert
  - severite/mineure
  - domaine/model
  - uc/UCCL01
  - enf/ENF05
---
# FT-026 — Erreur Fabric indisponible renvoyée en 500 au lieu de 503

> **Anomalie** · sévérité **mineure** · statut **Ouvert** · détecté le 2026-07-29
> Source : roadmap_dev.md § Écarts Infrastructure — Chaincode

## Constat

Une erreur Fabric `FailedPrecondition` (par exemple chaincode absent) est renvoyée en 500. Seules `ErrBlockchainUnavailable` et `ErrBlockchainUnreachable` sont traduites en 503.

## Cause

`internalErr` (`handlers.go:2060-2067`) ne reconnaît pas cette erreur, et l'adapter Fabric ne la convertit pas en erreur de domaine.

## Impact

Le client ne peut pas distinguer une panne d'infrastructure (réessayer plus tard) d'un bug serveur.

## Preuves

`adapters/in/rest/handlers.go:2060-2067`.

## Piste de correction (à valider par le PO)

Traduire dans l'adapter Fabric les erreurs gRPC d'indisponibilité (`FailedPrecondition`, `Unavailable`) en `ErrBlockchainUnavailable`.

## Critères de clôture

- [ ] Chaincode absent → 503
- [ ] Test de l'adapter couvrant la conversion

## Liens

- **Use cases** : [UCCL01 (analyse)](../specs/2-Analyse/UCCL-Composant_Lecture/UCCL01.md) · [UCCL01 (expression)](../specs/1-Expression/UCCL-Composant_Lecture/UCCL01.md)
- **Exigences non fonctionnelles** : [ENF05](../specs/1-Expression/Exigences_Non_Fonctionnelles.md)
- **Code** : [adapters/in/rest/handlers.go:2060](../adapters/in/rest/handlers.go) · [adapters/out/fabric/blockchain.go](../adapters/out/fabric/blockchain.go)
- **Fonctions** : [BlockchainPort.ListModelRecords](../docs/code/fonctions/model.BlockchainPort.ListModelRecords.md)
- **Tickets liés** : [FT-024 — Chaincode myrcc absent du canal de production](FT-024-chaincode-myrcc-absent-du-canal-de-production.md)

## Historique

- 2026-07-29 — Constaté lors de l'incident chaincode
- 2026-10-02 — Ticket créé ; constat revérifié au commit `2aa69c1`
