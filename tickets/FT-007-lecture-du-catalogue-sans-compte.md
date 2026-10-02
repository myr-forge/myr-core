---
id: FT-007
titre: "Lecture du catalogue sans compte"
type: decision
statut: a-trancher
severite: mineure
detecte: 2026-10-02
maj: 2026-10-02
composants: [adapters/in/rest]
uc: [UCA01, UCCL01]
rm: []
enf: [ENF12]
tags:
  - ticket
  - ticket/decision
  - statut/a-trancher
  - severite/mineure
  - domaine/identity
  - domaine/role
  - uc/UCA01
  - uc/UCCL01
  - enf/ENF12
---
# FT-007 — Lecture du catalogue sans compte

> **Décision à prendre** · sévérité **mineure** · statut **À trancher (PO)** · détecté le 2026-10-02
> Source : Annotation #remarque de UCA01 (expression) ; roadmap (UCCL01) ; note « TODO MYR-CORE (backend) » du vault

## Constat

Le besoin exprimé : si l'administrateur du réseau l'autorise, la consultation du catalogue (composants, modules) doit être possible sans compte, comme la vitrine d'un site. La roadmap note que UCCL01 (lecture publique) est bloqué par l'authentification.

## Cause

Les routes de lecture exigent une session (`auth`/`contrib` dans `server.go`) ; aucun réglage réseau n'ouvre la lecture aux visiteurs.

## Impact

Impossible de proposer une vitrine publique du réseau ; contradiction possible avec ENF12 tant que la règle n'est pas écrite.

## Preuves

`adapters/in/rest/server.go` (enveloppes `auth`/`contrib` sur `/api/components`). Note du vault : « On n'a pas besoin de créer un compte pour la lecture, mais il faut que l'administrateur le valide ».

## Piste de correction (à valider par le PO)

Ajouter au profil réseau un réglage « lecture publique » (désactivé par défaut) et un rôle visiteur sans permission d'écriture, appliqué par le RBAC ; préciser UCA01/UCCL01 et ENF12.

## Critères de clôture

- [ ] Règle écrite dans les specs (UCA01, UCCL01, ENF12)
- [ ] Lecture anonyme possible uniquement quand le réseau l'autorise
- [ ] Écriture toujours refusée aux visiteurs

## Liens

- **Use cases** : [UCA01](../specs/2-Analyse/UCA-Compte_et_Acces/UCA01.md) · [UCCL01](../specs/2-Analyse/UCCL-Composant_Lecture/UCCL01.md)
- **Exigences non fonctionnelles** : `ENF12` (tags `enf/…`)
- **Code** : `adapters/in/rest/server.go:123`
- **Fonctions** : [ModelService.List](../docs/code/fonctions/model.ModelService.List.md) · [ModelService.Get](../docs/code/fonctions/model.ModelService.Get.md)

## Historique

- 2026-10-02 — Ticket créé
