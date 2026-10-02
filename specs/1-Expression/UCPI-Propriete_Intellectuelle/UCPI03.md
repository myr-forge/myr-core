---
categorie: Propriété Intellectuelle
titre: "[RECLASSIFIÉ] Paramètres de langue de l'interface"
etat: reclassifié
tags:
  - couche/expression
  - type/use-case
  - famille/UCPI
  - domaine/model
  - domaine/payment
  - uc/UCPI03
---

# UCPI03 — Reclassifié

> **Ce use case a été reclassifié. Cet identifiant est conservé pour maintenir la continuité de la numérotation.**

## Décision de reclassification

| Champ | Valeur |
|-------|--------|
| Titre original | "Passage en version Anglaise / Chinoise" |
| Motif | Les paramètres de langue relèvent de la configuration de l'interface utilisateur, pas de la propriété intellectuelle |
| Date de reclassification | 2026 |
| Reclassifié vers | UCPAR01 — Passage en version Anglaise (supprimé, voir note ci-dessous) |
| | UCPAR02 — Passage en version Chinoise (supprimé, voir note ci-dessous) |

## Pourquoi la numérotation n'a pas été remaniée

Renuméroter UCPI04→UCPI03, UCPI05→UCPI04… aurait cassé toutes les références existantes (matrice de traçabilité, règles métier, liens entre UC). Ce stub préserve la numérotation tout en documentant explicitement la décision, conformément à la méthode Arrington.

## Retrait ultérieur de UCPAR

UCPAR01/02 étaient des use cases 100 % frontend (changement de langue de l'interface, aucune règle métier). Ils ont été supprimés du dépôt `myr` lors de la séparation des dépôts — l'i18n de l'interface relève désormais du dépôt GUI externe. Cet identifiant UCPI03 reste donc définitivement « reclassifié puis retiré » — aucun lien de renvoi actif.

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCPI — Propriete Intellectuelle](../../Carte_des_specs.md#UCPI%20—%20Propriete%20Intellectuelle)
- [UCPI03 — couche analyse](../../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI03.md)
- [Traçabilité UCPI03 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCPI03)

**Use cases cités**
- [UCPI04 — Définir un prix sur un Composant proprietaire](UCPI04.md)
- [UCPI05 — Définir un prix sur un Module proprietaire](UCPI05.md)

**Cité par**
- [Expression_des_besoins_Intro](../Expression_des_besoins_Intro.md)
- [todo (expression)](../todo.md)
- [Analyse_des_besoins](../../2-Analyse/Analyse_des_besoins.md)
- [todo (analyse)](../../2-Analyse/todo.md)
- [DC_D7_Payment](../../3-Conception/DC_D7_Payment.md)
- [todo (conception)](../../3-Conception/todo.md)

<!-- liens-obsidian:end -->
