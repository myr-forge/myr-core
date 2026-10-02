---
id: FT-006
titre: "Demandes de compte stockées sur un seul nœud"
type: decision
statut: a-trancher
severite: majeure
detecte: 2026-10-02
maj: 2026-10-02
composants: [adapters/out/localstorage]
uc: [UCA01]
rm: []
enf: []
tags:
  - ticket
  - ticket/decision
  - statut/a-trancher
  - severite/majeure
  - domaine/identity
  - uc/UCA01
---
# FT-006 — Demandes de compte stockées sur un seul nœud

> **Décision à prendre** · sévérité **majeure** · statut **À trancher (PO)** · détecté le 2026-10-02
> Source : Annotations #incoherence de UCA01 (expression et analyse)

## Constat

Les demandes de compte sont stockées dans un fichier local du serveur qui les reçoit (`request_store.go`). Les relectures de UCA01 relèvent que la création de compte doit rester accessible indépendamment du nœud, pour qu'une panne de serveur ne fasse perdre ni données ni accès.

## Cause

Choix de conception : `AccountRequest` est classée « hors blockchain (brouillon) » et persistée par l'adapter `localstorage`.

## Impact

Point de défaillance unique, contraire au principe de décentralisation : une demande déposée sur un nœud est invisible depuis les autres et perdue avec lui.

## Preuves

Annotations `#incoherence` : `specs/1-Expression/UCA-Compte_et_Acces/UCA01.md` (« chaque création de compte doit être accessible indépendamment du nœud ») et `specs/2-Analyse/UCA-Compte_et_Acces/UCA01.md` (« pourquoi pas partagé sur le réseau décentralisé ? »).

## Piste de correction (à valider par le PO)

À trancher par le PO : (a) inscrire les demandes sur le réseau (ledger, avec la question des données personnelles et du RGPD — ENF27) ; (b) les répliquer entre nœuds ; (c) assumer le stockage local et l'écrire comme risque accepté dans les specs.

## Critères de clôture

- [ ] Décision consignée (ADR ou spec UCA01)
- [ ] Annotations #incoherence de UCA01 levées

## Liens

- **Use cases** : [UCA01](../specs/2-Analyse/UCA-Compte_et_Acces/UCA01.md)
- **Specs** : [Conception_intro](../specs/3-Conception/Conception_intro.md)
- **Code** : `adapters/out/localstorage/request_store.go`
- **Fonctions** : [RequestStore.Save](../docs/code/fonctions/identity.RequestStore.Save.md)
- **Tickets liés** : [FT-005 — Aucune approbation manuelle des demandes de compte](FT-005-aucune-approbation-manuelle-des-demandes-de-compte.md)

## Historique

- 2026-10-02 — Ticket créé à partir des annotations de relecture
