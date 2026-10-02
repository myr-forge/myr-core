---
id: FT-008
titre: "Wallets non chiffrés au repos"
type: risque
statut: a-trancher
severite: majeure
detecte: 2026-10-02
maj: 2026-10-02
composants: [domain/identity]
uc: [UCA02]
rm: []
enf: [ENF10]
tags:
  - ticket
  - ticket/risque
  - statut/a-trancher
  - severite/majeure
  - domaine/identity
  - uc/UCA02
  - enf/ENF10
---
# FT-008 — Wallets non chiffrés au repos

> **Risque** · sévérité **majeure** · statut **À trancher (PO)** · détecté le 2026-10-02
> Source : ADR-03 (Conception_intro.md) ; Securite.md §4

## Constat

Les wallets (certificat X.509 et clé privée) sont des fichiers PEM en clair sous `~/.Myr/wallets/` (permissions `0600`). ENF10 exige un chiffrement au repos. La variable `WALLET_ENCRYPT_KEY`, citée dans d'anciennes versions de la documentation, n'est lue nulle part dans le code.

## Cause

Décision non tranchée (ADR-03) : chiffrer les wallets ou réviser ENF10.

## Impact

Une compromission du compte système du serveur expose les clés privées de toutes les identités hébergées. La roadmap (Post-V1) évoque encore une « rotation de `WALLET_ENCRYPT_KEY` » de « wallets SQLite », deux notions obsolètes.

## Preuves

Aucune occurrence de `WALLET_ENCRYPT_KEY` dans le code Go (hors `vendor/`).

## Piste de correction (à valider par le PO)

À trancher par le PO : chiffrement au repos avec une clé dérivée d'un secret d'exploitation (et procédure de rotation), ou révision d'ENF10 acceptant des PEM en clair sous permissions restreintes comme risque assumé.

## Critères de clôture

- [ ] Décision consignée dans ADR-03
- [ ] ENF10 et Securite.md alignés sur la décision
- [ ] Mention obsolète retirée de la roadmap

## Liens

- **Use cases** : [UCA02](../specs/2-Analyse/UCA-Compte_et_Acces/UCA02.md)
- **Exigences non fonctionnelles** : `ENF10` (tags `enf/…`)
- **Specs** : [Conception_intro — ADR-03](../specs/3-Conception/Conception_intro.md) · [Securite](../specs/3-Conception/Securite.md)
- **Code** : `domain/identity/service.go`
- **Tickets liés** : [FT-035 — Roadmap et tableaux d'état d'implémentation obsolètes](FT-035-roadmap-et-tableaux-d-etat-d-implementation.md)

## Historique

- 2026-10-02 — Ticket créé
