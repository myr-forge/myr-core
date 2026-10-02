---
id: FT-037
titre: "Commission sur le canal de fabrication externe"
type: decision
statut: a-trancher
severite: mineure
detecte: 2026-10-02
maj: 2026-10-02
composants: [domain/payment]
uc: [UCAUT01, UCPI01]
rm: [RM24, RM29]
enf: []
tags:
  - ticket
  - ticket/decision
  - statut/a-trancher
  - severite/mineure
  - domaine/payment
  - uc/UCAUT01
  - uc/UCPI01
  - rm/RM24
  - rm/RM29
---
# FT-037 — Commission sur le canal de fabrication externe

> **Décision à prendre** · sévérité **mineure** · statut **À trancher (PO)** · détecté le 2026-10-02
> Source : ADR-09 (Conception_intro.md), point ouvert ; roadmap Post-V1

## Constat

Un partenaire industriel externe (canal `external_adapter`) prélève probablement sa propre marge en amont du prix suivi par `AssetPrice`. Il reste à décider si le taux de commission du réseau (RM29) et la répartition (RM24) s'appliquent à l'identique sur ce canal.

## Cause

Point ouvert d'ADR-09.

## Impact

Bloque la conception financière du canal externe.

## Preuves

ADR-09.

## Piste de correction (à valider par le PO)

Décision PO, puis mise à jour de DC_D7_Payment (DC-D7-09) et RM29.

## Critères de clôture

- [ ] Décision consignée dans ADR-09

## Liens

- **Use cases** : [UCAUT01](../specs/2-Analyse/UCAUT-Automatisation/UCAUT01.md) · [UCPI01](../specs/2-Analyse/UCPI-Propriete_Intellectuelle/UCPI01.md)
- **Règles métier** : `RM24`, `RM29` (tags `rm/…`)
- **Specs** : [Conception_intro — ADR-09](../specs/3-Conception/Conception_intro.md) · [DC_D7_Payment](../specs/3-Conception/DC_D7_Payment.md)
- **Code** : `domain/payment/entity.go` · `domain/payment/service.go`
- **Tickets liés** : [FT-036 — Commissions : composition seule ou aussi chaîne de dérivation](FT-036-commissions-composition-seule-ou-aussi-chaine-de.md)

## Historique

- 2026-10-02 — Ticket créé
