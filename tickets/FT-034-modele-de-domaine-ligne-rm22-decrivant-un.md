---
id: FT-034
titre: "Modèle de domaine : ligne RM22 décrivant un écart"
type: incoherence
statut: ouvert
severite: mineure
detecte: 2026-10-02
maj: 2026-10-02
composants: [specs]
uc: [UCA02]
rm: [RM22]
enf: []
tags:
  - ticket
  - ticket/incoherence
  - statut/ouvert
  - severite/mineure
  - domaine/identity
  - uc/UCA02
  - rm/RM22
---
# FT-034 — Modèle de domaine : ligne RM22 décrivant un écart

> **Incohérence documentaire** · sévérité **mineure** · statut **Ouvert** · détecté le 2026-10-02
> Source : Annotations #incoherence de Modele_Domaine.md

## Constat

Dans `Modele_Domaine.md`, la ligne RM22 décrit un défaut de synchronisation du rôle de session (FT-002), pas la règle RM22 elle-même (changement de rôle réservé à l'administrateur, effectif au prochain ré-enrôlement). Le diagramme annote aussi la relation « owns (OwnerID) » comme non validée par rapport à la session (FT-001).

## Cause

Écarts d'implémentation consignés dans une spec, qui doit décrire le comportement cible.

## Impact

La spec mélange comportement cible et état du code.

## Preuves

`specs/3-Conception/Modele_Domaine.md`, annotations `#incoherence`.

## Piste de correction (à valider par le PO)

Réécrire la ligne RM22 selon la règle cible, retirer l'annotation « non validé vs session » une fois FT-001 et FT-002 suivis par ces tickets.

## Critères de clôture

- [ ] Modele_Domaine.md décrit uniquement le comportement cible
- [ ] Annotations levées

## Liens

- **Use cases** : [UCA02](../specs/2-Analyse/UCA-Compte_et_Acces/UCA02.md)
- **Règles métier** : `RM22` (tags `rm/…`)
- **Specs** : [Modele_Domaine](../specs/3-Conception/Modele_Domaine.md)
- **Code** : `adapters/in/rest/handlers_identity.go:388` · `adapters/in/rest/handlers.go:640`
- **Tickets liés** : [FT-001 — Propriété des assets non contrôlée côté serveur](FT-001-propriete-des-assets-non-controlee-cote-serveur.md) · [FT-002 — Rôle de session REST codé en dur à contributor](FT-002-role-de-session-rest-code-en-dur.md)

## Historique

- 2026-10-02 — Ticket créé
