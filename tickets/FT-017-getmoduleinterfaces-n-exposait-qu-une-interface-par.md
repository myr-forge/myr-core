---
id: FT-017
titre: "GetModuleInterfaces n'exposait qu'une interface par asset partagé"
type: anomalie
statut: resolu
severite: majeure
detecte: 2026-08-16
maj: 2026-10-02
composants: [domain/model]
uc: [UCMOD04, UCAM02]
rm: [RM13]
enf: []
tags:
  - ticket
  - ticket/anomalie
  - statut/resolu
  - severite/majeure
  - domaine/model
  - uc/UCMOD04
  - uc/UCAM02
  - rm/RM13
---
# FT-017 — GetModuleInterfaces n'exposait qu'une interface par asset partagé

> **Anomalie** · sévérité **majeure** · statut **Résolu** · détecté le 2026-08-16
> Source : Remontée myr-web (note « Ticket » du vault, 16/08)

## Constat

Un module contenant N instances du même composant n'exposait qu'une seule interface (4 instances → 1 interface exposée) au lieu d'une par `(instance, interface)` libre.

## Cause

Dédoublonnage `seen[inst.AssetID]` dans `getModuleInterfacesInto` : les instances suivantes d'un même asset étaient ignorées.

## Impact

Composition de produits faussée côté myr-web.

## Preuves

Relevés console myr-web : 4 instances → 1 interface exposée ; 1 instance → 1 interface exposée.

## Piste de correction

Corrigé : suppression du dédoublonnage ; filtre « interne » scindé en `internalGlobal` (liaison sans instance) et `internalByInstance` (liaison avec instance explicite).

## Critères de clôture

- [x] Une interface exposée par instance libre
- [x] Tests ajoutés et verts

## Liens

- **Use cases** : [UCMOD04 (analyse)](../specs/2-Analyse/UCMOD-Module/UCMOD04.md) · [UCMOD04 (expression)](../specs/1-Expression/UCMOD-Module/UCMOD04.md) · [UCAM02 (analyse)](../specs/2-Analyse/UCAM-Assemblage_Module/UCAM02.md) · [UCAM02 (expression)](../specs/1-Expression/UCAM-Assemblage_Module/UCAM02.md)
- **Règles métier** : [RM13](../specs/1-Expression/Regles_Metier.md)
- **Code** : [domain/model/service.go:994](../domain/model/service.go)
- **Fonctions** : [ModelService.GetModuleInterfaces](../docs/code/fonctions/model.ModelService.GetModuleInterfaces.md)
- **Tickets liés** : [FT-016 — Interfaces exposées d'un module sans identité d'instance](FT-016-interfaces-exposees-d-un-module-sans-identite.md)

## Historique

- 2026-08-16 — Remonté par myr-web
- 2026-10-01 — Corrigé — commit `6867439` ; tests `TestGetModuleInterfaces_MultipleInstancesSameAsset_AllExposed` et `TestGetModuleInterfaces_InternalConnection_ScopedToInstance`
