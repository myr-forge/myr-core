---
id: FT-016
titre: "Interfaces exposées d'un module sans identité d'instance"
type: ecart
statut: a-trancher
severite: majeure
detecte: 2026-08-16
maj: 2026-10-02
composants: [domain/model, adapters/in/rest]
uc: [UCMOD04, UCAM02]
rm: [RM13]
enf: []
tags:
  - ticket
  - ticket/ecart
  - statut/a-trancher
  - severite/majeure
  - domaine/model
  - uc/UCMOD04
  - uc/UCAM02
  - rm/RM13
---
# FT-016 — Interfaces exposées d'un module sans identité d'instance

> **Écart spec ↔ code** · sévérité **majeure** · statut **À trancher (PO)** · détecté le 2026-08-16
> Source : Limite résiduelle du correctif FT-017 (roadmap § Bugs bloquants)

## Constat

Quand plusieurs instances d'un même asset exposent leurs interfaces, `GET .../interfaces` renvoie N entrées `AssetInterface` portant le même `id` : rien dans la réponse ne dit à quelle instance appartient chacune.

## Cause

`AssetInterface` est rattachée à l'asset (`ID`, `AssetID`), pas à l'instance.

## Impact

Le client (myr-web) doit recouper avec `WorkspaceInstances` pour savoir quelle instance relier ; une liaison peut viser la mauvaise instance.

## Preuves

`domain/model/entity.go` (structure `AssetInterface`) ; test `TestGetModuleInterfaces_MultipleInstancesSameAsset_AllExposed` (3 entrées de même id).

## Piste de correction (à valider par le PO)

À trancher : ajouter un `instance_id` à chaque interface exposée (DTO dédié côté REST ou champ calculé), sans changer le stockage des interfaces de l'asset.

## Critères de clôture

- [ ] Contrat d'API décidé et documenté (API_REST.md)
- [ ] Chaque interface exposée identifiable par `(instance_id, interface_id)`
- [ ] myr-web informé du changement

## Liens

- **Use cases** : [UCMOD04](../specs/2-Analyse/UCMOD-Module/UCMOD04.md) · [UCAM02](../specs/2-Analyse/UCAM-Assemblage_Module/UCAM02.md)
- **Règles métier** : `RM13` (tags `rm/…`)
- **Specs** : [API_REST](../specs/3-Conception/API_REST.md)
- **Code** : `domain/model/entity.go`
- **Fonctions** : [ModelService.GetModuleInterfaces](../docs/code/fonctions/model.ModelService.GetModuleInterfaces.md)
- **Tickets liés** : [FT-017 — GetModuleInterfaces n'exposait qu'une interface par asset partagé](FT-017-getmoduleinterfaces-n-exposait-qu-une-interface-par.md)

## Historique

- 2026-08-16 — Identifiée lors du correctif FT-017
- 2026-10-02 — Ticket créé
