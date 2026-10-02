---
id: FT-030
titre: "Parité CLI et REST incomplète"
type: ecart
statut: ouvert
severite: majeure
detecte: 2026-10-02
maj: 2026-10-02
composants: [adapters/in/rest]
uc: [UCADM01, UCADM03, UCADM04, UCPI01, UCPI02, UCADM07]
rm: []
enf: []
tags:
  - ticket
  - ticket/ecart
  - statut/ouvert
  - severite/majeure
  - domaine/network
  - domaine/channel
  - domaine/payment
  - domaine/role
  - uc/UCADM01
  - uc/UCADM03
  - uc/UCADM04
  - uc/UCPI01
  - uc/UCPI02
  - uc/UCADM07
---
# FT-030 — Parité CLI et REST incomplète

> **Écart spec ↔ code** · sévérité **majeure** · statut **Ouvert** · détecté le 2026-10-02
> Source : Tableau d'état d'implémentation ; roadmap

## Constat

Le CLI et l'API REST doivent exposer les mêmes capacités. Ce n'est pas le cas :
- `network` : le REST n'expose que la liste et le profil actif (pas de création, mise à jour, activation, suppression, test, ajout de pair) ;
- `channel` : aucune route REST (organisations, nœuds) ;
- `payment` : aucune route REST, `adapters/out/payment/` vide ;
- `role` : aucune gestion des rôles en REST (seulement le contrôle `requireRole`).

## Cause

Développement mené d'abord côté CLI.

## Impact

Un client graphique (myr-web) ou un script tiers ne peut pas administrer le réseau ni gérer les rôles ; les deux surfaces divergent.

## Preuves

Notes de package [adapters/in/rest](../docs/code/adapters-in-rest.md) et [adapters/in/cli](../docs/code/adapters-in-cli.md) (routes et commandes, avec les méthodes du domaine atteintes).

## Piste de correction (à valider par le PO)

Ajouter les routes manquantes dans `handlers_<domaine>.go` en appelant les mêmes méthodes de service que le CLI, domaine par domaine (priorité à fixer par le PO). Exception assumée : le démantèlement d'un réseau reste CLI uniquement (RM28).

## Critères de clôture

- [ ] Chaque commande CLI d'administration a son équivalent REST (hors exceptions documentées)

## Liens

- **Use cases** : [UCADM01 (analyse)](../specs/2-Analyse/UCADM-Administration/UCADM01.md) · [UCADM01 (expression)](../specs/1-Expression/UCADM-Administration/UCADM01.md) · [UCADM03 (analyse)](../specs/2-Analyse/UCADM-Administration/UCADM03.md) · [UCADM03 (expression)](../specs/1-Expression/UCADM-Administration/UCADM03.md) · [UCADM04 (analyse)](../specs/2-Analyse/UCADM-Administration/UCADM04.md) · [UCADM04 (expression)](../specs/1-Expression/UCADM-Administration/UCADM04.md) · [UCPI01 (analyse)](../specs/2-Analyse/UCPI-Propriete_Intellectuelle/UCPI01.md) · [UCPI01 (expression)](../specs/1-Expression/UCPI-Propriete_Intellectuelle/UCPI01.md) · [UCPI02 (analyse)](../specs/2-Analyse/UCPI-Propriete_Intellectuelle/UCPI02.md) · [UCPI02 (expression)](../specs/1-Expression/UCPI-Propriete_Intellectuelle/UCPI02.md) · [UCADM07 (analyse)](../specs/2-Analyse/UCADM-Administration/UCADM07.md) · [UCADM07 (expression)](../specs/1-Expression/UCADM-Administration/UCADM07.md)
- **Specs** : [API_REST](../specs/3-Conception/API_REST.md) · [DC_CLI_Admin](../specs/3-Conception/DC_CLI_Admin.md)
- **Code** : [adapters/in/rest/server.go](../adapters/in/rest/server.go)
- **Fonctions** : [NetworkService.Add](../docs/code/fonctions/network.NetworkService.Add.md) · [ChannelService.AddOrganisation](../docs/code/fonctions/channel.ChannelService.AddOrganisation.md) · [RoleService.Create](../docs/code/fonctions/role.RoleService.Create.md)

## Historique

- 2026-10-02 — Ticket créé
