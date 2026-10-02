---
id: FT-003
titre: "Le client CA du serveur REST ignore le profil réseau actif"
type: anomalie
statut: a-trancher
severite: critique
detecte: 2026-07-17
maj: 2026-10-02
composants: [cmd/api, adapters/out/fabric]
uc: [UCA01]
rm: []
enf: []
tags:
  - ticket
  - ticket/anomalie
  - statut/a-trancher
  - severite/critique
  - domaine/identity
  - domaine/network
  - uc/UCA01
---
# FT-003 — Le client CA du serveur REST ignore le profil réseau actif

> **Anomalie** · sévérité **critique** · statut **À trancher (PO)** · détecté le 2026-07-17
> Source : roadmap_dev.md § Écarts Identité & Session ; ADR-08 (Conception_intro.md) ; écart E4 (DC_D1_Auth_Identity.md)

## Constat

`myr network update <id> --auto-register` met bien à jour le profil réseau, et `GET /api/identity/policy` affiche `allow_auto_register: true`. Pourtant `POST /api/identity/request` reste `pending` sans secret : aucune requête n'atteint la CA (vérifié sur le conteneur `ca.ca-org1` du serveur de production le 2026-07-17).

## Cause

`cmd/api/main.go:84` construit la configuration Fabric uniquement depuis l'environnement (`fabricadapter.ConfigFromEnv()` : variables `FABRIC_CA_*` ou `fabric.env`). Les champs CA du `NetworkProfile` actif ne sont jamais lus pour créer le client CA (`caPort`, lignes 104-131), alors que `ConfigFromProfile` existe et sert déjà au CLI et au pool multi-réseau. Sans variables d'environnement, `caPort` reste `nil` et `AutoRegister` échoue toujours (« aucun CA configuré »).

## Impact

L'auto-enregistrement (UCA01) est silencieusement inopérant alors que CLI et API affichent une configuration correcte. Un redémarrage du service ne corrige rien.

## Preuves

- `cmd/api/main.go:84` (`ConfigFromEnv`), `:104-131` (construction de `caPort`).
- Constat serveur du 2026-07-17 rapporté dans la roadmap (réseau `net-1783775510080523221`).

## Piste de correction (à valider par le PO)

Deux options décrites par ADR-08, à trancher par le PO : (a) construire la configuration CA depuis `ConfigFromProfile(profilActif)`, l'environnement ne servant que de surcharge (alignement sur le CLI et le `NetworkPool`) ; (b) faire de l'environnement la seule source de vérité et déprécier les champs CA du profil.

## Critères de clôture

- [ ] Décision PO consignée dans ADR-08
- [ ] Auto-enregistrement fonctionnel avec la seule configuration du profil réseau (option a) ou champs CA du profil retirés (option b)
- [ ] Test d'intégration sur le serveur (`scripts/test_remote.ps1`)

## Liens

- **Use cases** : [UCA01 (analyse)](../specs/2-Analyse/UCA-Compte_et_Acces/UCA01.md) · [UCA01 (expression)](../specs/1-Expression/UCA-Compte_et_Acces/UCA01.md)
- **Specs** : [Conception_intro — ADR-08](../specs/3-Conception/Conception_intro.md) · [DC_D1_Auth_Identity](../specs/3-Conception/DC_D1_Auth_Identity.md)
- **Code** : [cmd/api/main.go:84](../cmd/api/main.go) · [cmd/api/main.go:104](../cmd/api/main.go) · [adapters/out/fabric/config.go](../adapters/out/fabric/config.go)
- **Fonctions** : [IdentityService.AutoRegister](../docs/code/fonctions/identity.IdentityService.AutoRegister.md) · [NetworkService.GetActive](../docs/code/fonctions/network.NetworkService.GetActive.md)
- **Tickets liés** : [FT-004 — Échec de l'auto-enregistrement CA non journalisé](FT-004-echec-de-l-auto-enregistrement-ca-non.md) · [FT-005 — Aucune approbation manuelle des demandes de compte](FT-005-aucune-approbation-manuelle-des-demandes-de-compte.md) · [FT-025 — Le serveur signale un nœud Fabric déconnecté](FT-025-le-serveur-signale-un-nud-fabric-deconnecte.md)

## Historique

- 2026-07-17 — Constaté sur le serveur de production
- 2026-10-02 — Ticket créé ; cause revérifiée au commit `2aa69c1`
