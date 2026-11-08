---
id: FT-004
titre: "Échec de l'auto-enregistrement CA non journalisé"
type: anomalie
statut: en-cours
severite: majeure
detecte: 2026-07-17
maj: 2026-11-08
composants: [adapters/in/rest]
uc: [UCA01]
rm: []
enf: [ENF11]
tags:
  - ticket
  - ticket/anomalie
  - statut/en-cours
  - severite/majeure
  - domaine/identity
  - uc/UCA01
  - enf/ENF11
---
# FT-004 — Échec de l'auto-enregistrement CA non journalisé

> **Anomalie** · sévérité **majeure** · statut **En cours** · détecté le 2026-07-17
> Source : roadmap_dev.md § Écart connexe — échec AutoRegister avalé sans trace

## Constat

Quand `AutoRegister` échoue, le handler retombe sur une demande `pending` sans rien journaliser. Le commentaire annonce « logguer l'erreur », mais aucun appel de log n'existe.

## Cause

`handlers_identity.go:248-262` : la branche d'échec (`err3 != nil`) ne contient qu'un commentaire.

## Impact

L'échec est invisible côté serveur comme côté client : le diagnostic de FT-003 a nécessité un accès SSH direct aux conteneurs.

## Preuves

`adapters/in/rest/handlers_identity.go:248-262`.

## Piste de correction (à valider par le PO)

Journaliser l'échec sans données sensibles (pseudo et cause, jamais le secret), par exemple `log.Printf("auto-register CA échoué pour %s : %v", pseudo, err)`. Corrigeable indépendamment de FT-003.

## Critères de clôture

- [x] Un échec d'`AutoRegister` produit une ligne de log exploitable
- [x] Aucun secret ni donnée personnelle superflue dans le log (ENF11)

## Liens

- **Use cases** : [UCA01](../specs/2-Analyse/UCA-Compte_et_Acces/UCA01.md)
- **Exigences non fonctionnelles** : `ENF11` (tags `enf/…`)
- **Code** : `adapters/in/rest/handlers_identity.go:248`
- **Fonctions** : [IdentityService.AutoRegister](../docs/code/fonctions/identity.IdentityService.AutoRegister.md)
- **Tickets liés** : [FT-003 — Le client CA du serveur REST ignore le profil réseau actif](FT-003-le-client-ca-du-serveur-rest-ignore.md)

## Historique

- 2026-07-17 — Écart documenté dans la roadmap
- 2026-10-02 — Ticket créé ; constat revérifié au commit `2aa69c1`
- 2026-11-08 — Corrigé — commit `e2724c9` : `AutoRegister` est remplacé par `RequestAccess`, qui ne masque plus l'échec CA — la demande est conservée au statut `failed` et l'erreur est retournée (500, journalisée par `internalErr` ; message d'erreur en CLI) ; commit `e8498f3` : le message nomme l'identité `pseudo@org`, sans secret ni e-mail (test `TestRequestAccess_CAError_Failed`). La cause racine FT-003 reste ouverte. Branche `feature/ledger-invariants`, ni fusionnée ni déployée
- 2026-11-08 — Repassé en cours : la correction est sur la branche `feature/ledger-invariants`, pas encore fusionnée dans `develop` (règle du cycle de vie, `_Index_tickets.md`) — passera à `resolu` à la fusion
