---
id: FT-042
titre: "Interfaces à plusieurs contraintes de valeur"
type: decision
statut: a-trancher
severite: majeure
detecte: 2026-11-09
maj: 2026-11-09
composants: [domain/model, chaincode/model, adapters/in/rest, adapters/in/cli]
uc: [UCAM01, UCAM03, UCCE06]
rm: [RM11]
enf: []
tags:
  - ticket
  - ticket/decision
  - statut/a-trancher
  - severite/majeure
  - domaine/model
  - uc/UCAM01
  - uc/UCAM03
  - uc/UCCE06
  - rm/RM11
---
# FT-042 — Interfaces à plusieurs contraintes de valeur

> **Décision à prendre** · sévérité **majeure** · statut **À trancher (PO)** · détecté le 2026-11-09
> Source : revue de code (R-13) et décision du PO du 2026-11-09

## Constat

Une interface (`AssetInterface`) ne porte qu'une contrainte de valeur : `ValueMin`, `ValueMax`, `IsRange` et `Unit`. Deux interfaces d'unités différentes (sortie 5 V, entrée [3 ; 6] mA) voient leurs valeurs comparées comme si l'unité était la même.

Décision du PO : une interface peut porter une ou plusieurs contraintes de valeur, par exemple une tension ET un courant (RM11).

## Cause

Le modèle n'a qu'une valeur par interface, dans le domaine ([entity.go](../domain/model/entity.go), `AssetInterface`) comme dans le chaincode ([chaincode/model/entity.go](../chaincode/model/entity.go), même structure). `ifaceValuesOverlap` ([service.go](../domain/model/service.go)) compare les deux valeurs sans tenir compte de leurs unités.

## Impact

- La compatibilité est fausse entre interfaces d'unités différentes : une paire est acceptée ou refusée à tort.
- Une alimentation 5 V / 2 A ne peut pas être décrite.
- Le changement touche le ledger. Les interfaces sont enregistrées dans `Model3D.Interfaces` par le chaincode, dont la structure `AssetInterface` perdrait un champ qu'elle ne connaît pas : le chaincode doit être mis à jour (nouvelle séquence). Les interfaces déjà soumises gardent leur forme actuelle et doivent rester lisibles.

## Preuves

- `ifaceValuesOverlap` dans [service.go](../domain/model/service.go) : seules `ValueMin`/`ValueMax` sont comparées ; l'unité ne sert qu'à savoir s'il y a une contrainte.
- [chaincode/model/entity.go](../chaincode/model/entity.go) : `AssetInterface` à valeur unique.

## Piste de correction (à valider par le PO)

- `AssetInterface.Constraints` : liste de contraintes `{value_min, value_max, is_range, unit}`. Les champs actuels se lisent comme une première contrainte, pour les assets déjà soumis.
- Compatibilité : pour chaque contrainte présente des deux côtés, les valeurs se chevauchent (valeur simple = [v ; v], RM11).

**À trancher :**
1. **Correspondance entre contraintes** : même unité exacte (V avec V ; mV et V ne correspondent pas, myr-core ne convertit rien), ou même grandeur avec conversion (mV ↔ V).
2. **Contrainte présente d'un seul côté** (alimentation 5 V / 2 A ↔ consommateur 5 V qui n'indique pas de courant) : ignorée (paire compatible) ou paire refusée.
3. **Aucune unité commune** (V d'un côté, mA de l'autre) : paire compatible faute de contrainte comparable, ou refusée.
4. **Forme de l'API** : `constraints: [...]` seul, ou champs actuels conservés un temps pour myr-web.

## Critères de clôture

- [ ] Points 1 à 4 tranchés ; RM11 mise à jour
- [ ] Domaine, chaincode, REST, CLI et guide myr-web alignés ; chaincode mis à jour sur le réseau de test, puis en production
- [ ] Tests RM11 : plusieurs contraintes, contrainte présente d'un seul côté, aucune unité commune, interface soumise avant le changement

## Liens

- **Use cases** : [UCAM01 (analyse)](../specs/2-Analyse/UCAM-Assemblage_Module/UCAM01.md) · [UCAM03 (analyse)](../specs/2-Analyse/UCAM-Assemblage_Module/UCAM03.md) · [UCCE06 (analyse)](../specs/2-Analyse/UCCE-Composant_Ecriture/UCCE06.md)
- **Règles métier** : [RM11](../specs/1-Expression/Regles_Metier.md)
- **Code** : [domain/model/entity.go](../domain/model/entity.go) · [domain/model/service.go](../domain/model/service.go) · [chaincode/model/entity.go](../chaincode/model/entity.go)
- **Tickets liés** : [FT-015 — Critère tag absent des interfaces](FT-015-critere-tag-absent-des-interfaces.md) · [FT-040 — Fonctions propres aux outils extérieurs implémentées dans myr-core](FT-040-fonctions-d-outils-exterieurs-dans-le-coeur.md)

## Historique

- 2026-11-09 — Ticket créé ; décision du PO : une interface peut porter une ou plusieurs contraintes de valeur (ex. tension ET courant). Points 1 à 4 à trancher
