---
id: FT-004
titre: "Échec de l'auto-enregistrement CA non journalisé"
type: anomalie
statut: ouvert
severite: majeure
detecte: 2026-07-17
maj: 2026-10-02
composants: [adapters/in/rest]
uc: [UCA01]
rm: []
enf: [ENF11]
tags:
  - ticket
  - ticket/anomalie
  - statut/ouvert
  - severite/majeure
  - domaine/identity
  - uc/UCA01
  - enf/ENF11
---
# FT-004 — Échec de l'auto-enregistrement CA non journalisé

> **Anomalie** · sévérité **majeure** · statut **Ouvert** · détecté le 2026-07-17
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

- [ ] Un échec d'`AutoRegister` produit une ligne de log exploitable
- [ ] Aucun secret ni donnée personnelle superflue dans le log (ENF11)

## Liens

- **Use cases** : [UCA01](../specs/2-Analyse/UCA-Compte_et_Acces/UCA01.md)
- **Exigences non fonctionnelles** : `ENF11` (tags `enf/…`)
- **Code** : `adapters/in/rest/handlers_identity.go:248`
- **Fonctions** : [IdentityService.AutoRegister](../docs/code/fonctions/identity.IdentityService.AutoRegister.md)
- **Tickets liés** : [FT-003 — Le client CA du serveur REST ignore le profil réseau actif](FT-003-le-client-ca-du-serveur-rest-ignore.md)

## Historique

- 2026-07-17 — Écart documenté dans la roadmap
- 2026-10-02 — Ticket créé ; constat revérifié au commit `2aa69c1`
