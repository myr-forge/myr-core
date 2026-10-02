---
id: FT-005
titre: "Aucune approbation manuelle des demandes de compte"
type: ecart
statut: a-trancher
severite: majeure
detecte: 2026-07-17
maj: 2026-10-02
composants: [domain/identity, adapters/in/rest, adapters/in/cli]
uc: [UCA01, UCA08]
rm: [RM20, RM22]
enf: []
tags:
  - ticket
  - ticket/ecart
  - statut/a-trancher
  - severite/majeure
  - domaine/identity
  - uc/UCA01
  - uc/UCA08
  - rm/RM20
  - rm/RM22
---
# FT-005 — Aucune approbation manuelle des demandes de compte

> **Écart spec ↔ code** · sévérité **majeure** · statut **À trancher (PO)** · détecté le 2026-07-17
> Source : DC_D1_Auth_Identity.md (écarts E2/E3, annotation #ecart) ; note « TODO MYR-CORE (backend) » du vault

## Constat

Une demande de compte (`AccountRequest`) ne passe à `approved` que par l'auto-enregistrement CA. Sans auto-enregistrement (ou quand il échoue, FT-003), elle reste `pending` indéfiniment : ni le CLI ni le REST ne permettent à un administrateur d'approuver une demande en attente. Le dépôt `myr-web` affiche donc des comptes « en attente » qui ne peuvent pas être validés depuis `myr-core`. La demande ne porte par ailleurs aucun champ structuré pour un rôle souhaité (seulement un message libre).

## Cause

`IdentityService` n'expose pas d'opération d'approbation ; `myr identity requests` et `GET /api/identity/requests` ne font que lister.

## Impact

Aucun nouvel utilisateur ne peut obtenir d'identité sans auto-enregistrement ; la demande d'un rôle supplémentaire (UCA08) n'a pas de support structuré.

## Preuves

- `domain/identity/port_in.go` : aucune méthode d'approbation.
- Annotation `#ecart "a trancher"` dans `DC_D1_Auth_Identity.md` (section `AccountRequest`).

## Piste de correction (à valider par le PO)

Ajouter au domaine une approbation (et un refus) de demande par un administrateur (`identity.admin`), exposée à parité en CLI (`myr identity approve`) et REST ; ajouter un champ « rôle souhaité » à `AccountRequest`. À cadrer avec la décision de FT-006 (où vivent les demandes).

## Critères de clôture

- [ ] Un administrateur peut approuver ou refuser une demande en CLI et en REST
- [ ] La demande approuvée produit une identité CA et passe à `approved`
- [ ] Spécification UCA01/UCA08 mise à jour si le flux change

## Liens

- **Use cases** : [UCA01 (analyse)](../specs/2-Analyse/UCA-Compte_et_Acces/UCA01.md) · [UCA01 (expression)](../specs/1-Expression/UCA-Compte_et_Acces/UCA01.md) · [UCA08 (analyse)](../specs/2-Analyse/UCA-Compte_et_Acces/UCA08.md) · [UCA08 (expression)](../specs/1-Expression/UCA-Compte_et_Acces/UCA08.md)
- **Règles métier** : [RM20](../specs/1-Expression/Regles_Metier.md) · [RM22](../specs/1-Expression/Regles_Metier.md)
- **Specs** : [DC_D1_Auth_Identity](../specs/3-Conception/DC_D1_Auth_Identity.md) · [DC_CLI_Identity](../specs/3-Conception/DC_CLI_Identity.md)
- **Code** : [domain/identity/port_in.go](../domain/identity/port_in.go) · [adapters/out/localstorage/request_store.go](../adapters/out/localstorage/request_store.go)
- **Fonctions** : [IdentityService.ListRequests](../docs/code/fonctions/identity.IdentityService.ListRequests.md) · [IdentityService.AutoRegister](../docs/code/fonctions/identity.IdentityService.AutoRegister.md)
- **Tickets liés** : [FT-003 — Le client CA du serveur REST ignore le profil réseau actif](FT-003-le-client-ca-du-serveur-rest-ignore.md) · [FT-006 — Demandes de compte stockées sur un seul nœud](FT-006-demandes-de-compte-stockees-sur-un-seul.md)

## Historique

- 2026-07-17 — Écart E2 documenté
- 2026-10-02 — Ticket créé ; remontée myr-web (comptes en attente) rattachée
