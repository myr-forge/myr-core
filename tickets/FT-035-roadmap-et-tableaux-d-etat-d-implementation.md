---
id: FT-035
titre: "Roadmap et tableaux d'état d'implémentation obsolètes"
type: dette
statut: ouvert
severite: mineure
detecte: 2026-10-02
maj: 2026-10-02
composants: [specs]
uc: [UCA02]
rm: []
enf: []
tags:
  - ticket
  - ticket/dette
  - statut/ouvert
  - severite/mineure
  - uc/UCA02
---
# FT-035 — Roadmap et tableaux d'état d'implémentation obsolètes

> **Dette technique** · sévérité **mineure** · statut **Ouvert** · détecté le 2026-10-02
> Source : Comparaison roadmap_dev.md ↔ code au commit 2aa69c1

## Constat

`specs/roadmap_dev.md` contient des constats dépassés :
- bug « rôle initial » situé dans `domain/auth/service.go`, fichier inexistant (le défaut réel est FT-002) ;
- « Chaincode : entity seulement (7 champs) » et « IPFS : adapter vide », alors que le chaincode (4 fonctions, 27 champs) et l'adapter IPFS existent ;
- UCA02 décrit comme « connexion JWT » (il n'y a jamais eu de JWT) ;
- fusion composant/module (ADR-11) : retrait de `moduleDTO`, de l'arbre `myr module` et dispatch de la soumission déjà réalisés, mais cases non cochées ;
- items « Post-V1 » obsolètes : révocation JWT, rotation de `WALLET_ENCRYPT_KEY` des « wallets SQLite », migration SQLite → PostgreSQL (aucune base SQL par principe) ;
- sections UI (UCIG, UCPAR, UCDOC) qui relèvent du dépôt myr-web.

## Cause

Roadmap tenue à la main, qui mélange état d'avancement et backlog.

## Impact

Une lecture de la roadmap conduit à de mauvaises priorités ; les faits techniques y sont dispersés.

## Preuves

Constats vérifiés dans le code (voir les tickets liés).

## Piste de correction (à valider par le PO)

Rafraîchir la roadmap et la limiter au backlog fonctionnel ; renvoyer les faits techniques vers ce dossier `tickets/` (un lien par ticket) plutôt que de les décrire en double.

## Critères de clôture

- [ ] Roadmap cohérente avec le code
- [ ] Sections « Écarts » remplacées par des liens vers les tickets

## Liens

- **Use cases** : [UCA02](../specs/2-Analyse/UCA-Compte_et_Acces/UCA02.md)
- **Specs** : [roadmap_dev](../specs/roadmap_dev.md)
- **Code** : `chaincode/model/contract.go` · `adapters/out/ipfs/storage.go` · `adapters/in/rest/server.go:25` · `adapters/in/rest/handlers_identity.go:391`
- **Tickets liés** : [FT-002 — Rôle de session REST codé en dur à contributor](FT-002-role-de-session-rest-code-en-dur.md) · [FT-008 — Wallets non chiffrés au repos](FT-008-wallets-non-chiffres-au-repos.md) · [FT-024 — Chaincode myrcc absent du canal de production](FT-024-chaincode-myrcc-absent-du-canal-de-production.md) · [FT-031 — Les specs citent des commandes et routes qui n'existent pas](FT-031-les-specs-citent-des-commandes-et-routes.md)

## Historique

- 2026-10-02 — Ticket créé
