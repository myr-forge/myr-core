---
id: FT-031
titre: "Les specs citent des commandes et routes qui n'existent pas"
type: incoherence
statut: ouvert
severite: majeure
detecte: 2026-10-02
maj: 2026-10-02
composants: [specs]
uc: [UCAM08, UCMOD01, UCMOD04, UCDEV02, UCADM07, UCA03, UCA04]
rm: []
enf: []
tags:
  - ticket
  - ticket/incoherence
  - statut/ouvert
  - severite/majeure
  - uc/UCAM08
  - uc/UCMOD01
  - uc/UCMOD04
  - uc/UCDEV02
  - uc/UCADM07
  - uc/UCA03
  - uc/UCA04
---
# FT-031 — Les specs citent des commandes et routes qui n'existent pas

> **Incohérence documentaire** · sévérité **majeure** · statut **Ouvert** · détecté le 2026-10-02
> Source : Traçabilité générée (docs/code/Tracabilite_UC_Code.md)

## Constat

La comparaison automatique specs ↔ code relève des commandes et routes citées sans équivalent dans le code. Principaux cas :
- `myr module …` (create, get, interfaces, add-assembly…), cité par une dizaine de use cases, alors que l'arbre `myr module` a été retiré lors de la fusion composant/module (ADR-11) ;
- `myr role edit`, `myr org role`, `myr network sync`, `myr peer` ;
- `/api/auth/me` (UCA04), `/api/identity/logout` (UCA03), `/api/identity/requests/{id}/approve` (UCA01).

Une partie correspond à des fonctionnalités non encore développées (commandes, paiement, transferts), une autre à des citations obsolètes.

## Cause

Specs non mises à jour après des renommages ou des retraits (notamment ADR-11).

## Impact

Un lecteur suit une commande ou une route qui n'existe pas ; la traçabilité use case → code est rompue.

## Preuves

Section « Références des specs sans correspondance dans le code » de [Tracabilite_UC_Code](../docs/code/Tracabilite_UC_Code.md#R%C3%A9f%C3%A9rences%20des%20specs%20sans%20correspondance%20dans%20le%20code).

## Piste de correction (à valider par le PO)

Pour chaque citation : la remplacer par la commande ou route actuelle (ex. `myr module get` → `myr model get`), ou, s'il s'agit d'une fonctionnalité à venir, la conserver comme comportement cible et ouvrir le ticket de développement correspondant.

## Critères de clôture

- [ ] Plus aucune citation `myr module`
- [ ] Chaque citation restante correspond à une commande existante ou à un ticket de développement

## Liens

- **Use cases** : [UCAM08](../specs/2-Analyse/UCAM-Assemblage_Module/UCAM08.md) · [UCMOD01](../specs/2-Analyse/UCMOD-Module/UCMOD01.md) · [UCMOD04](../specs/2-Analyse/UCMOD-Module/UCMOD04.md) · [UCDEV02](../specs/2-Analyse/UCDEV-Developpement/UCDEV02.md) · [UCADM07](../specs/2-Analyse/UCADM-Administration/UCADM07.md) · [UCA03](../specs/2-Analyse/UCA-Compte_et_Acces/UCA03.md) · [UCA04](../specs/2-Analyse/UCA-Compte_et_Acces/UCA04.md)
- **Specs** : [DC_CLI_Model](../specs/3-Conception/DC_CLI_Model.md) · [DC_CLI_Admin](../specs/3-Conception/DC_CLI_Admin.md)
- **Code** : `adapters/in/cli/root.go:47`
- **Tickets liés** : [FT-022 — Date de retrait des routes legacy /api/modules](FT-022-date-de-retrait-des-routes-legacy-api.md) · [FT-035 — Roadmap et tableaux d'état d'implémentation obsolètes](FT-035-roadmap-et-tableaux-d-etat-d-implementation.md)

## Historique

- 2026-10-02 — Ticket créé
