---
id: FT-009
titre: "Registrar CA unique par organisation"
type: risque
statut: ouvert
severite: mineure
detecte: 2026-10-02
maj: 2026-10-02
composants: [adapters/out/fabric]
uc: [UCA01]
rm: []
enf: []
tags:
  - ticket
  - ticket/risque
  - statut/ouvert
  - severite/mineure
  - domaine/identity
  - uc/UCA01
---
# FT-009 — Registrar CA unique par organisation

> **Risque** · sévérité **mineure** · statut **Ouvert** · détecté le 2026-10-02
> Source : ADR-07 (Conception_intro.md) ; roadmap § Compléments Post-V1

## Constat

Une seule identité administrateur CA par profil réseau signe tous les enregistrements d'une organisation.

## Cause

Choix de conception assumé (ADR-07) ; aucune mitigation entreprise.

## Impact

Point de défaillance unique à l'échelle d'une organisation : perte ou compromission de cette clé = plus aucun enregistrement possible, ou enregistrements frauduleux.

## Preuves

ADR-07.

## Piste de correction (à valider par le PO)

Évaluer rotation de clé, HSM, ou registrar de secours par organisation.

## Critères de clôture

- [ ] Mitigation choisie et documentée, ou risque explicitement accepté dans Securite.md

## Liens

- **Use cases** : [UCA01 (analyse)](../specs/2-Analyse/UCA-Compte_et_Acces/UCA01.md) · [UCA01 (expression)](../specs/1-Expression/UCA-Compte_et_Acces/UCA01.md)
- **Specs** : [Conception_intro — ADR-07](../specs/3-Conception/Conception_intro.md)
- **Code** : [adapters/out/fabric/ca_client.go](../adapters/out/fabric/ca_client.go)
- **Fonctions** : [CAPort.Register](../docs/code/fonctions/identity.CAPort.Register.md)

## Historique

- 2026-10-02 — Ticket créé
