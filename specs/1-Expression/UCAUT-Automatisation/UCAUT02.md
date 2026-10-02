---
categorie: Automatisation
titre: "Commande en ligne de Asset"
probabilite: 5
impact: 3
importance: 15
etat: relire
tags:
  - couche/expression
  - type/use-case
  - famille/UCAUT
  - domaine/model
  - domaine/payment
  - domaine/role
  - uc/UCAUT02
---

# Commande en ligne de asset

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Consommateur" as CL
actor "Developpeur" as D
actor "Manufactureur" as M

rectangle "Application MYR" {
    usecase "Commander un asset en ligne" as UC1
    usecase "Récupérer le prix via l'API" as UC2
    usecase "Transmettre la commande" as UC3
}

CL --> UC1
D --> UC2
M --> UC3
UC1 ..> UC2 : <<include>>
UC1 ..> UC3 : <<include>>

@enduml
```

## Contexte

Commande en ligne d'un asset via l'interface MYR ou une boutique partenaire avec récupération du prix via l'API.

Conformément au principe de parité CLI/REST, une commande passée par une boutique partenaire via l'API doit pouvoir être reproduite en CLI pour le compte d'un consommateur, au même titre que via l'interface graphique.

## Pré-conditions

- Être connecté au réseau
- asset disponible à la commande
- Adresse de livraison renseignée

## Scénario

**Étape initiale :** Un client (boutique partenaire, script, interface graphique tierce...) appelle l'API pour commander un asset

### Flux nominal — Commande passée

1. Le prix de l'asset est récupéré via l'API
2. La commande est confirmée par le client
3. La transaction est enregistrée sur la blockchain
4. La commande est transmise à la boutique ou au manufactureur

## Post-conditions

- La commande est enregistrée et en cours de traitement
- L'utilisateur reçoit une confirmation

## Diagramme d'activités

```plantuml
@startuml
skin rose
title Commande en ligne de Asset
start
:Récupérer le prix de l'asset via l'API;
:Confirmer la commande;
:Enregistrer la transaction sur la blockchain;
:Transmettre la commande à la boutique ou au manufactureur;
:Retourner une confirmation au client;
stop
@enduml
```

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCAUT — Automatisation](../../Carte_des_specs.md#UCAUT%20—%20Automatisation)
- [UCAUT02 — couche analyse](../../2-Analyse/UCAUT-Automatisation/UCAUT02.md)
- [Traçabilité UCAUT02 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCAUT02)

**Exigences fonctionnelles couvertes**
- [EF45 — Automatiser la commande en ligne d'un asset](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Cité par**
- [Expression_des_besoins_Intro](../Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../Matrice_Tracabilite.md)
- [todo (expression)](../todo.md)
- [Analyse_des_besoins](../../2-Analyse/Analyse_des_besoins.md)
- [todo (analyse)](../../2-Analyse/todo.md)
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)
- [DC_D9_Automatisation](../../3-Conception/DC_D9_Automatisation.md)
- [todo (conception)](../../3-Conception/todo.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
