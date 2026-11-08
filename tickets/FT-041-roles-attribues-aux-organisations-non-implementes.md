---
id: FT-041
titre: "Rôles attribués aux organisations non implémentés"
type: ecart
statut: ouvert
severite: mineure
detecte: 2026-11-08
maj: 2026-11-08
composants: [domain/role, chaincode/model, adapters/in/cli]
uc: [UCADM06, UCADM07]
rm: [RM34, RM35, RM37, RM48]
enf: []
tags:
  - ticket
  - ticket/ecart
  - statut/ouvert
  - severite/mineure
  - domaine/role
  - uc/UCADM06
  - uc/UCADM07
  - rm/RM34
  - rm/RM35
  - rm/RM37
  - rm/RM48
---
# FT-041 — Rôles attribués aux organisations non implémentés

> **Écart spec ↔ code** · sévérité **mineure** · statut **Ouvert** · détecté le 2026-11-08
> Source : décision du PO du 2026-11-08 sur l'incohérence UCADM06/07 ↔ ADR-12 (R-03 étape 4)

## Constat

Le PO a tranché : un rôle s'attribue à une identité (certificat, `Myr.role`) et à une organisation (UCADM06). Les droits effectifs d'une identité sont l'union des deux (RM37, RM48), à réévaluer lors de la phase alpha. Seule l'attribution par identité est implémentée.

## Cause

La table des rôles publiée sur le ledger (`chaincode/model/roles.go`, `adapters/out/fabric/role_store.go`) ne porte que les rôles et leurs permissions. Aucune liaison organisation ↔ rôles n'existe : pas de `myr org role assign|remove`. Le contrôle d'écriture du chaincode (`requirePermission`) et le contrôle REST (`requireRole`) ne lisent que le rôle du certificat.

## Impact

Fonctionnalité absente, sans risque de sécurité : sans liaison, une identité n'a que les droits de son propre rôle, ce qui est plus restrictif que la cible.

## Preuves

- `chaincode/model/rules.go` (`requirePermission`) : seul `signer.role` est consulté.
- `adapters/in/rest/handlers.go` (`hasPermission`) : seul `sess.Role` est consulté.
- `adapters/in/cli/org.go` : aucune sous-commande `role`.

## Piste de correction (à valider par le PO)

Publier dans la même table du ledger, pour chaque organisation (MSP), la liste de ses rôles (`admin` exclu, RM34). Retirer un rôle de la table le retire aussi des organisations (RM35). Le chaincode accorderait une permission si le rôle du certificat ou l'un des rôles de l'organisation du signataire la porte. Côté serveur, `RoleService` exposerait une vérification prenant l'identité et son organisation, utilisée par `requireRole`. CLI : `myr org role assign|remove`, avec un équivalent REST si myr-web en a besoin.

## Critères de clôture

- [ ] `myr org role assign|remove` publient la liaison sur le ledger (tous les canaux)
- [ ] Chaincode et serveur appliquent l'union des droits (tests RM37)
- [ ] Suppression d'un rôle propagée aux organisations (test RM35)

## Liens

- **Use cases** : [UCADM06 (analyse)](../specs/2-Analyse/UCADM-Administration/UCADM06.md) · [UCADM07 (analyse)](../specs/2-Analyse/UCADM-Administration/UCADM07.md)
- **Règles métier** : [RM37, RM48](../specs/1-Expression/Regles_Metier.md)
- **Specs** : [Chaincode §4.3](../specs/3-Conception/Chaincode.md)
- **Code** : [chaincode/model/roles.go](../chaincode/model/roles.go) · [adapters/out/fabric/role_store.go](../adapters/out/fabric/role_store.go)

## Historique

- 2026-11-08 — Ticket créé ; décision PO (attribution par organisation et par identité) consignée dans UCADM06/07, RM48 et l'ADR-12
