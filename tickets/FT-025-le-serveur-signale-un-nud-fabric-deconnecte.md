---
id: FT-025
titre: "Le serveur signale un nœud Fabric déconnecté"
type: anomalie
statut: a-verifier
severite: majeure
detecte: 2026-08-16
maj: 2026-10-02
composants: [cmd/api, adapters/out/fabric]
uc: [UCADM02]
rm: []
enf: [ENF05, ENF06]
tags:
  - ticket
  - ticket/anomalie
  - statut/a-verifier
  - severite/majeure
  - domaine/network
  - uc/UCADM02
  - enf/ENF05
  - enf/ENF06
---
# FT-025 — Le serveur signale un nœud Fabric déconnecté

> **Anomalie** · sévérité **majeure** · statut **À vérifier** · détecté le 2026-08-16
> Source : Note « Ticket » du vault (remontée myr-web, 16/08)

## Constat

`GET /api/status` renvoie `"connected": false` : le gateway Fabric du serveur n'est pas connecté, alors que myr-web et nginx fonctionnent normalement.

## Cause

À établir — causes possibles : `fabric.env` ou variables `FABRIC_*` absents au démarrage (voir FT-003, FT-029), peer arrêté, chaincode absent (FT-024).

## Impact

Le serveur fonctionne en mode dégradé sans blockchain : aucune soumission ni lecture du ledger.

## Preuves

Relevé myr-web du 2026-08-16.

## Piste de correction (à valider par le PO)

Sur le serveur : `systemctl --user status myr`, `journalctl --user -u myr`, présence de `fabric.env`, `docker ps` des peers ; corriger la cause et documenter le diagnostic.

## Critères de clôture

- [ ] `/api/status` renvoie `connected: true` en production
- [ ] Cause identifiée et consignée dans ce ticket

## Liens

- **Use cases** : [UCADM02 (analyse)](../specs/2-Analyse/UCADM-Administration/UCADM02.md) · [UCADM02 (expression)](../specs/1-Expression/UCADM-Administration/UCADM02.md)
- **Exigences non fonctionnelles** : [ENF05](../specs/1-Expression/Exigences_Non_Fonctionnelles.md) · [ENF06](../specs/1-Expression/Exigences_Non_Fonctionnelles.md)
- **Specs** : [Deploiement](../specs/3-Conception/Deploiement.md)
- **Code** : [cmd/api/main.go:84](../cmd/api/main.go) · [adapters/in/rest/handlers_health.go](../adapters/in/rest/handlers_health.go)
- **Tickets liés** : [FT-003 — Le client CA du serveur REST ignore le profil réseau actif](FT-003-le-client-ca-du-serveur-rest-ignore.md) · [FT-024 — Chaincode myrcc absent du canal de production](FT-024-chaincode-myrcc-absent-du-canal-de-production.md) · [FT-029 — make deploy suppose un fabric.env présent sur le serveur](FT-029-make-deploy-suppose-un-fabric-env-present.md)

## Historique

- 2026-08-16 — Remonté par myr-web
- 2026-10-02 — Ticket créé ; état à revérifier sur le serveur
