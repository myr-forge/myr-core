---
id: FT-042
titre: "Interfaces à plusieurs contraintes de valeur"
type: decision
statut: ouvert
severite: majeure
detecte: 2026-11-09
maj: 2026-11-09
composants: [domain/model, chaincode/model, adapters/in/rest, adapters/in/cli]
uc: [UCAM01, UCAM03, UCCE06]
rm: [RM11, RM53]
enf: []
tags:
  - ticket
  - ticket/decision
  - statut/ouvert
  - severite/majeure
  - domaine/model
  - uc/UCAM01
  - uc/UCAM03
  - uc/UCCE06
  - rm/RM11
  - rm/RM53
---
# FT-042 — Interfaces à plusieurs contraintes de valeur

> **Décision à prendre** · sévérité **majeure** · statut **Ouvert** · détecté le 2026-11-09
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

**Décisions du PO (2026-11-09)**, consignées dans RM11 et RM53 :
1. **Correspondance** : par unité de base, préfixes SI convertis (1 V = 1 000 mV). L'unité est lue dans le vocabulaire telle quelle, sinon comme un préfixe usuel (`n`, `µ`/`u`, `m`, `c`, `d`, `k`, `M`, `G`, casse comprise) suivi d'une unité connue, sinon refusée. Aucune conversion entre unités de base différentes.
2. **Contrainte présente d'un seul côté** : elle n'impose rien, la paire reste compatible.
3. **Aucune unité de base commune** entre deux interfaces qui portent chacune des contraintes : paire incompatible.
4. **API** : `constraints: [{value_min, value_max, is_range, unit}]` remplace les quatre champs actuels ; les interfaces déjà soumises sont relues comme une contrainte unique.

Le développement viendra à la demande du PO : domaine, chaincode (montée de version répétée sur myr-test), REST, CLI (syntaxe à définir dans `DC_CLI_Model.md`), lecture des interfaces déjà soumises, guide myr-web.

## Critères de clôture

- [x] Points 1 à 4 tranchés ; RM11 mise à jour et RM53 créée
- [ ] Domaine, chaincode, REST, CLI et guide myr-web alignés ; chaincode mis à jour sur le réseau de test, puis en production
- [ ] Tests RM11 : plusieurs contraintes, contrainte présente d'un seul côté, aucune unité commune, interface soumise avant le changement

## Liens

- **Use cases** : [UCAM01 (analyse)](../specs/2-Analyse/UCAM-Assemblage_Module/UCAM01.md) · [UCAM03 (analyse)](../specs/2-Analyse/UCAM-Assemblage_Module/UCAM03.md) · [UCCE06 (analyse)](../specs/2-Analyse/UCCE-Composant_Ecriture/UCCE06.md)
- **Règles métier** : [RM11](../specs/1-Expression/Regles_Metier.md)
- **Code** : [domain/model/entity.go](../domain/model/entity.go) · [domain/model/service.go](../domain/model/service.go) · [chaincode/model/entity.go](../chaincode/model/entity.go)
- **Tickets liés** : [FT-015 — Critère tag absent des interfaces](FT-015-critere-tag-absent-des-interfaces.md) · [FT-040 — Fonctions propres aux outils extérieurs implémentées dans myr-core](FT-040-fonctions-d-outils-exterieurs-dans-le-coeur.md)

## Historique

- 2026-11-09 — Ticket créé ; décision du PO : une interface peut porter une ou plusieurs contraintes de valeur (ex. tension ET courant). Points 1 à 4 à trancher
- 2026-11-09 — Décisions du PO sur les points 1 à 4 (voir § Piste) ; RM11 révisée et RM53 créée ; développement à faire sur demande du PO
