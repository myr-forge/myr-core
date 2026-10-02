---
categorie: Compte et Accès
titre: "Vérification des accès du rôle attribué"
probabilite: 1
impact: 3
importance: 3
etat: relire
tags:
  - couche/expression
  - type/use-case
  - famille/UCA
  - domaine/identity
  - domaine/role
  - uc/UCA05
---

# Vérification des accès du rôle attribué

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Utilisateur" as U
actor "Administrateur" as ADM

rectangle "Application MYR" {
    usecase "Effectuer une action" as UC1
    usecase "Vérifier les droits du rôle" as UC2
    usecase "Attribuer un rôle" as UC3
}

U --> UC1
UC1 ..> UC2 : <<include>>
ADM --> UC3

@enduml
```

## Contexte

Vérifier que l'utilisateur peut réaliser les actions autorisées par son rôle et ne peut pas réaliser celles qui lui sont refusées.

## Pré-conditions

- Être connecté au réseau
- Rôle attribué à l'utilisateur

## Scénario

**Étape initiale :** L'utilisateur tente d'effectuer une action

### Flux nominal — Action autorisée

1. L'action est dans les droits du rôle attribué
2. Le système exécute l'action

### Flux erreur — Action non autorisée

1. L'action n'est pas dans les droits du rôle attribué
2. Message d'erreur : "Vous n'avez pas les droits nécessaires pour cette action"

## Post-conditions

- Les droits du rôle sont respectés

## Diagrammes

### Rôles disponibles dans le système

```plantuml
@startuml
:Administrateur:
:Lecteur:
:Concepteur:
:Consommateur:
:Manufactureur:
:Developpeur:
@enduml
```

### Diagramme d'activités

```plantuml
@startuml
skin rose
title Vérification des accès du rôle attribué
start
:Tenter d'effectuer une action;
if (Action dans les droits du rôle?) then (oui)
  :Exécuter l'action;
  stop
else (non)
  :Afficher "Vous n'avez pas les droits nécessaires pour cette action";
  stop
endif
@enduml
```

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCA — Compte et Acces](../../Carte_des_specs.md#UCA%20—%20Compte%20et%20Acces)
- [UCA05 — couche analyse](../../2-Analyse/UCA-Compte_et_Acces/UCA05.md)
- [Traçabilité UCA05 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCA05)

**Exigences fonctionnelles couvertes**
- [EF05 — Contrôler les accès selon le rôle attribué](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Cité par**
- [Expression_des_besoins_Intro](../Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../Matrice_Tracabilite.md)
- [todo (expression)](../todo.md)
- [Analyse_des_besoins](../../2-Analyse/Analyse_des_besoins.md)
- [UCA08 (analyse)](../../2-Analyse/UCA-Compte_et_Acces/UCA08.md)
- [todo (analyse)](../../2-Analyse/todo.md)
- [todo (conception)](../../3-Conception/todo.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
