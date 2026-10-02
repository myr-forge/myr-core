---
categorie: Propriété Intellectuelle
titre: "Cloner un Composant sur un réseau exterieur"
probabilite: 1
impact: 2
importance: 2
etat: relire
tags:
  - couche/expression
  - type/use-case
  - famille/UCPI
  - domaine/model
  - domaine/payment
  - uc/UCPI08
---

# Cloner un Composant sur un réseau exterieur

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Concepteur" as C

rectangle "Application MYR" {
    usecase "Cloner un composant sur un réseau externe" as UC1
    usecase "Vérifier compatibilité de licence" as UC2
}

C --> UC1
UC1 ..> UC2 : <<include>>

@enduml
```

## Contexte

Un composant peut être cloné vers un réseau MYR externe sous réserve de compatibilité de licence.

Conformément au principe de parité CLI/REST, le clonage d'un composant vers un réseau externe doit pouvoir être déclenché en CLI, au même titre que via l'interface graphique.

## Pré-conditions

- Être connecté au réseau source
- Avoir les droits de clonage (licence compatible)
- Réseau de destination accessible

## Scénario

**Étape initiale :** `myr model clone <id> --target-network <id>` est exécutée (ou l'appel API équivalent)

### Flux nominal — Clonage autorisé

1. Le réseau de destination est transmis
2. Le système vérifie la compatibilité de licence
3. La transaction de clonage est soumise sur les deux réseaux

### Flux alternatif — Réseau cible déjà connu (profil de connexion existant)

1. Le réseau cible est déjà référencé dans les profils de connexion locaux
2. Le clonage est initié directement, sans saisie manuelle des informations du réseau cible
3. Le composant est publié sur le réseau cible avec le même UUID et les métadonnées originales

### Flux erreur — Licence incompatible

1. Erreur métier : "La licence du composant ne permet pas le clonage vers ce réseau"

## Post-conditions

- Le composant est disponible sur le réseau de destination
- La traçabilité de l'origine est conservée

## Diagramme d'activités

```plantuml
@startuml
skin rose
title Cloner un Composant sur un réseau extérieur
start
:Transmettre l'identifiant du composant et le réseau cible (myr model clone);
:Vérifier la compatibilité de licence;
if (Licence compatible?) then (oui)
  :Soumettre la transaction de clonage sur les deux réseaux;
  stop
else (non)
  :Retourner "La licence du composant ne permet pas le clonage vers ce réseau";
  stop
endif
@enduml
```

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCPI — Propriete Intellectuelle](../../Carte_des_specs.md#UCPI%20—%20Propriete%20Intellectuelle)
- [UCPI08 — couche analyse](../../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI08.md)
- [Traçabilité UCPI08 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCPI08)

**Exigences fonctionnelles couvertes**
- [EF36 — Cloner un composant sur un réseau externe](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Cité par**
- [Expression_des_besoins_Intro](../Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../Matrice_Tracabilite.md)
- [todo (expression)](../todo.md)
- [Analyse_des_besoins](../../2-Analyse/Analyse_des_besoins.md)
- [UCPI09 (analyse)](../../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI09.md)
- [todo (analyse)](../../2-Analyse/todo.md)
- [API_REST](../../3-Conception/API_REST.md)
- [Chaincode](../../3-Conception/Chaincode.md)
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)
- [DC_D7_Payment](../../3-Conception/DC_D7_Payment.md)
- [todo (conception)](../../3-Conception/todo.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
