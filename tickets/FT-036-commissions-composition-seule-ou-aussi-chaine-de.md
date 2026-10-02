---
id: FT-036
titre: "Commissions : composition seule ou aussi chaîne de dérivation"
type: decision
statut: a-trancher
severite: majeure
detecte: 2026-10-02
maj: 2026-10-02
composants: [domain/payment]
uc: [UCPI02, UCPI04, UCPI01]
rm: [RM23, RM24]
enf: []
tags:
  - ticket
  - ticket/decision
  - statut/a-trancher
  - severite/majeure
  - domaine/payment
  - uc/UCPI02
  - uc/UCPI04
  - uc/UCPI01
  - rm/RM23
  - rm/RM24
---
# FT-036 — Commissions : composition seule ou aussi chaîne de dérivation

> **Décision à prendre** · sévérité **majeure** · statut **À trancher (PO)** · détecté le 2026-10-02
> Source : Annotations #question de Regles_Metier.md, UCPI02, UCPI04

## Constat

Le périmètre de la « chaîne de propriété » rémunérée à la livraison n'est pas tranché : uniquement les auteurs des composants du module livré (composition — seul cas couvert par l'algorithme de `DC_D7_Payment.md` §4), ou aussi les auteurs de la lignée de dérivation de chaque composant. Les deux couches de specs divergent.

## Cause

Question produit ouverte.

## Impact

Bloque l'implémentation de RM23/RM24, le schéma de l'entité `Commission` et l'algorithme de distribution.

## Preuves

Annotations `#question` dans `Regles_Metier.md` (RM23/RM24), `UCPI02.md` (analyse), `UCPI04.md` (expression et analyse).

## Piste de correction (à valider par le PO)

Décision PO ; puis mettre à jour RM23/RM24, DC_D7_Payment §3-§4 et lever les annotations.

## Critères de clôture

- [ ] Décision consignée dans Regles_Metier.md
- [ ] DC_D7_Payment.md aligné
- [ ] Annotations levées

## Liens

- **Use cases** : [UCPI02 (analyse)](../specs/2-Analyse/UCPI-Propriete_Intellectuelle/UCPI02.md) · [UCPI02 (expression)](../specs/1-Expression/UCPI-Propriete_Intellectuelle/UCPI02.md) · [UCPI04 (analyse)](../specs/2-Analyse/UCPI-Propriete_Intellectuelle/UCPI04.md) · [UCPI04 (expression)](../specs/1-Expression/UCPI-Propriete_Intellectuelle/UCPI04.md) · [UCPI01 (analyse)](../specs/2-Analyse/UCPI-Propriete_Intellectuelle/UCPI01.md) · [UCPI01 (expression)](../specs/1-Expression/UCPI-Propriete_Intellectuelle/UCPI01.md)
- **Règles métier** : [RM23](../specs/1-Expression/Regles_Metier.md) · [RM24](../specs/1-Expression/Regles_Metier.md)
- **Specs** : [DC_D7_Payment](../specs/3-Conception/DC_D7_Payment.md)
- **Code** : [domain/payment/entity.go](../domain/payment/entity.go)
- **Tickets liés** : [FT-037 — Commission sur le canal de fabrication externe](FT-037-commission-sur-le-canal-de-fabrication-externe.md)

## Historique

- 2026-10-02 — Ticket créé
