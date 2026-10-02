---
id: FT-010
titre: "UCA04 et UCA07 s'appuient sur une commande qui ne couvre pas leur besoin"
type: incoherence
statut: a-trancher
severite: mineure
detecte: 2026-10-02
maj: 2026-10-02
composants: [specs]
uc: [UCA04, UCA07]
rm: []
enf: []
tags:
  - ticket
  - ticket/incoherence
  - statut/a-trancher
  - severite/mineure
  - domaine/identity
  - uc/UCA04
  - uc/UCA07
---
# FT-010 — UCA04 et UCA07 s'appuient sur une commande qui ne couvre pas leur besoin

> **Incohérence documentaire** · sévérité **mineure** · statut **À trancher (PO)** · détecté le 2026-10-02
> Source : Annotations #remarque de UCA04 et UCA07 (expression) ; DC_CLI_Identity.md

## Constat

UCA04 (vérification de connexion) et UCA07 (vérification du rôle) citent `myr identity status` / `GET /api/identity/status`, qui ne renvoie que le statut d'enrôlement CA (`pending`/`active`/`suspended`) — ni réseau, organisation ou peer (UCA04), ni rôle RBAC (UCA07). L'analyse de UCA07 affirme en outre qu'aucun endpoint ne permet de reconsulter le rôle après la connexion.

## Cause

Use cases rédigés avant que le périmètre de `GetStatus` soit fixé.

## Impact

Ces use cases paraissent couverts alors qu'ils ne le sont pas ; les deux couches de specs se contredisent.

## Preuves

Annotations `#remarque` dans `specs/1-Expression/UCA-Compte_et_Acces/UCA04.md` et `UCA07.md`.

## Piste de correction (à valider par le PO)

Soit modéliser les capacités manquantes (détails de connexion, rôle courant), soit réduire le périmètre des deux use cases au statut CA.

## Critères de clôture

- [ ] Périmètre tranché
- [ ] Les deux couches de specs cohérentes, annotations levées

## Liens

- **Use cases** : [UCA04 (analyse)](../specs/2-Analyse/UCA-Compte_et_Acces/UCA04.md) · [UCA04 (expression)](../specs/1-Expression/UCA-Compte_et_Acces/UCA04.md) · [UCA07 (analyse)](../specs/2-Analyse/UCA-Compte_et_Acces/UCA07.md) · [UCA07 (expression)](../specs/1-Expression/UCA-Compte_et_Acces/UCA07.md)
- **Specs** : [DC_CLI_Identity](../specs/3-Conception/DC_CLI_Identity.md)
- **Fonctions** : [IdentityService.GetStatus](../docs/code/fonctions/identity.IdentityService.GetStatus.md)
- **Tickets liés** : [FT-002 — Rôle de session REST codé en dur à contributor](FT-002-role-de-session-rest-code-en-dur.md)

## Historique

- 2026-10-02 — Ticket créé à partir des annotations de relecture
