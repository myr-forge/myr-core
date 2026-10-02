---
categorie: Module
titre: "Visualiser les composants d'un Module"
probabilite: 2
impact: 5
importance: 10
tags:
  - couche/expression
  - type/use-case
  - famille/UCMOD
  - domaine/model
  - uc/UCMOD04
---

# Visualiser les composants d'un Module

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Utilisateur" as U

rectangle "Application MYR" {
    usecase "Visualiser les composants d'un module" as UC1
    usecase "Explorer la hiérarchie" as UC2
}

U --> UC1
UC1 .> UC2 : <<extend>>

@enduml
```

## Contexte

Un Module est composé de composants visualisables. L'utilisateur peut explorer la structure interne d'un module pour en comprendre la composition.

## Pré-conditions

- Être connecté au réseau
- Avoir un identifiant de module

## Scénario

**Étape initiale :** `myr module get <id>` est exécutée (ou l'appel API équivalent), en lecture seule

### Flux nominal — Composition retournée

1. La composition du module (composants constitutifs, liaisons internes) est retournée
2. `myr module interfaces <id>` complète la vue avec les interfaces exposées
3. Les sous-modules éventuels peuvent être explorés récursivement (`myr module get <sousModuleID>`)

## Post-conditions

- La composition complète du module est visible

## Diagramme d'activités

```plantuml
@startuml
skin rose
title Visualiser les composants d'un Module
start
:Transmettre l'identifiant du module (myr module get);
:Retourner la composition (composants, liaisons internes);
if (Sous-modules à explorer?) then (oui)
  :Explorer récursivement les sous-modules;
  stop
else (non)
  stop
endif
@enduml
```

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCMOD — Module](../../Carte_des_specs.md#UCMOD%20—%20Module)
- [UCMOD04 — couche analyse](../../2-Analyse/UCMOD-Module/UCMOD04.md)
- [Traçabilité UCMOD04 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCMOD04)

**Exigences fonctionnelles couvertes**
- [EF29 — Visualiser la composition d'un module](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Cité par**
- [Expression_des_besoins_Intro](../Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../Matrice_Tracabilite.md)
- [todo (expression)](../todo.md)
- [Analyse_des_besoins](../../2-Analyse/Analyse_des_besoins.md)
- [UCDEV02 (analyse)](../../2-Analyse/UCDEV-Developpement/UCDEV02.md)
- [todo (analyse)](../../2-Analyse/todo.md)
- [API_REST](../../3-Conception/API_REST.md)
- [Architecture_Composition](../../3-Conception/Architecture_Composition.md)
- [Chaincode](../../3-Conception/Chaincode.md)
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)
- [DC_D1_Auth_Identity](../../3-Conception/DC_D1_Auth_Identity.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
