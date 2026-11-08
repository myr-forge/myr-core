---
id: FT-039
titre: "Tests d'intégration REST en échec depuis l'authentification obligatoire"
type: anomalie
statut: en-cours
severite: majeure
detecte: 2026-11-08
maj: 2026-11-08
composants: [adapters/in/rest]
uc: [UCDEV01, UCMOD06, UCCE06]
rm: [RM13, RM17]
enf: []
tags:
  - ticket
  - ticket/anomalie
  - statut/en-cours
  - severite/majeure
  - uc/UCDEV01
  - uc/UCMOD06
  - uc/UCCE06
  - rm/RM13
  - rm/RM17
---
# FT-039 — Tests d'intégration REST en échec depuis l'authentification obligatoire

> **Anomalie** · sévérité **majeure** · statut **En cours** · détecté le 2026-11-08
> Source : exécution de `go test -tags integration ./adapters/in/rest/integration/` pendant R-03

## Constat

7 des 10 tests d'intégration de l'API REST échouaient : toutes leurs requêtes recevaient une 401 (`GET /api/licenses` renvoyait alors l'objet d'erreur au lieu de la liste). Une fois l'authentification rétablie, trois autres écarts apparaissaient : création de composant en 500 (« stockage de brouillons non configuré »), création de module en 400, interface « toujours présente » après suppression.

## Cause

- Les tests datent d'avant l'authentification obligatoire des routes de données : le serveur de test n'avait pas de service d'identité, aucune session ne pouvait être ouverte.
- Le serveur de test suivait un ancien câblage de `cmd/api/main.go` (sans stockage des brouillons ni des assets retirés).
- Deux attentes ne suivaient plus le contrat : création d'un module par `POST /api/modules` en JSON (depuis l'ADR-11, un module est un composant créé par `POST /api/components` en multipart, avec instances et assemblages), et liste d'interfaces sans le slot virtuel toujours présent (RM13).

Ces échecs sont antérieurs à R-03 : ils sont identiques au commit `dfa6ebb`.

## Impact

Les tests d'intégration ne détectaient plus aucune régression de l'API REST de bout en bout.

## Preuves

`adapters/in/rest/integration/api_integration_test.go` — `go test -count=1 -tags integration ./adapters/in/rest/integration/` au commit `dfa6ebb` : 7 FAIL, dont `création: got 401, want 201`.

## Piste de correction (à valider par le PO)

Ouvrir une vraie session dans chaque test : identité réelle (`domain/identity`) sur une autorité de certification en mémoire qui émet de vrais certificats X.509 lus par le parseur de l'adapter Fabric (comme `memBC` remplace la blockchain) ; câblage aligné sur `cmd/api/main.go` ; attentes alignées sur l'ADR-11 et RM13.

## Critères de clôture

- [x] Les 10 tests d'intégration REST passent
- [x] Chaque test passe par une connexion réelle (`POST /api/identity/session`) avec rôle lu dans le certificat
- [ ] Les tests d'intégration exécutés à chaque modification (intégration continue, FT-027)

## Liens

- **Use cases** : [UCDEV01 (analyse)](../specs/2-Analyse/UCDEV-Developpement/UCDEV01.md) · [UCMOD06 (analyse)](../specs/2-Analyse/UCMOD-Module/UCMOD06.md) · [UCCE06 (analyse)](../specs/2-Analyse/UCCE-Composant_Ecriture/UCCE06.md)
- **Règles métier** : `RM13`, `RM17` (tags `rm/…`)
- **Code** : [adapters/in/rest/integration/api_integration_test.go](../adapters/in/rest/integration/api_integration_test.go)
- **Tickets liés** : [FT-027 — Intégration continue absente du dépôt](FT-027-integration-continue-absente-du-depot.md)

## Historique

- 2026-11-08 — Ticket créé ; échecs constatés à l'identique au commit `dfa6ebb` (antérieurs à R-03)
- 2026-11-08 — Corrigé — commit `353fece` : session réelle sur une CA en mémoire, câblage aligné sur `cmd/api`, attentes alignées sur l'ADR-11 et RM13 ; les 10 tests passent. Branche `feature/ledger-invariants`, pas encore fusionnée
