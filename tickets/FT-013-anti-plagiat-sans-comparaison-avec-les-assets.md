---
id: FT-013
titre: "Anti-plagiat sans comparaison avec les assets existants"
type: ecart
statut: a-trancher
severite: majeure
detecte: 2026-10-02
maj: 2026-10-02
composants: [domain/model]
uc: [UCCE01, UCPI06]
rm: [RM01]
enf: [ENF29]
tags:
  - ticket
  - ticket/ecart
  - statut/a-trancher
  - severite/majeure
  - domaine/model
  - uc/UCCE01
  - uc/UCPI06
  - rm/RM01
  - enf/ENF29
---
# FT-013 — Anti-plagiat sans comparaison avec les assets existants

> **Écart spec ↔ code** · sévérité **majeure** · statut **À trancher (PO)** · détecté le 2026-10-02
> Source : Écart E4 (Analyse_des_besoins.md) ; roadmap § Anti-plagiat structurel

## Constat

À l'ajout d'un composant de type `base`, l'empreinte SHA-256 du fichier est calculée et stockée, mais elle n'est comparée à aucun asset existant, et aucune analyse de similarité (> 50 %) n'est effectuée. RM01 rend cette vérification obligatoire avant tout ajout.

## Cause

`AddFull` (`service.go:85`) se limite au calcul du hash. L'algorithme de similarité n'est pas choisi (décision PO bloquante citée par la roadmap Beta).

## Impact

Un même fichier peut être enregistré plusieurs fois comme création originale ; la protection de la propriété intellectuelle repose sur ce contrôle.

## Preuves

`domain/model/service.go` : dans `AddFull`, le hash est affecté (`m.Hash = hash`) sans recherche d'un asset portant le même hash.

## Piste de correction (à valider par le PO)

Étape 1, sans décision PO : rejeter un hash identique à un asset existant. Étape 2 : choisir l'algorithme de similarité (librairie Go, service externe ou développement propre — nouvelle technologie à valider) et l'intégrer derrière un port.

## Critères de clôture

- [ ] Hash identique rejeté avec contact administration (RM01)
- [ ] Algorithme de similarité choisi et documenté dans specs/3-Conception
- [ ] Tests couvrant les deux cas

## Liens

- **Use cases** : [UCCE01](../specs/2-Analyse/UCCE-Composant_Ecriture/UCCE01.md) · [UCPI06](../specs/2-Analyse/UCPI-Propriete_Intellectuelle/UCPI06.md)
- **Règles métier** : `RM01` (tags `rm/…`)
- **Exigences non fonctionnelles** : `ENF29` (tags `enf/…`)
- **Specs** : [Sequence_soumission_asset](../specs/3-Conception/Sequence_soumission_asset.md)
- **Code** : `domain/model/service.go:85`
- **Fonctions** : [ModelService.AddFull](../docs/code/fonctions/model.ModelService.AddFull.md)

## Historique

- 2026-10-02 — Ticket créé ; constat vérifié au commit `2aa69c1`
