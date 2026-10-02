---
categorie: Composant Ecriture
titre: "Ajout d'un composant Numérique"
probabilite: 3
impact: 5
importance: 15
tags:
  - couche/expression
  - type/use-case
  - famille/UCCE
  - domaine/model
  - uc/UCCE03
  - rm/RM01
---

# Ajout d'un composant Numérique

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Concepteur" as C

rectangle "Application MYR" {
    usecase "Ajouter un composant numérique" as UC1
    usecase "Vérifier le hash" as UC2
    usecase "Enregistrer sur la blockchain" as UC3
}

C --> UC1
UC1 ..> UC2 : <<include>>
UC1 ..> UC3 : <<include>>

@enduml
```

## Contexte

Ajout d'un élément numérique (logiciel, firmware, driver...) comme composant dans le système MYR.

## Pré-conditions

- Être connecté au réseau
- Avoir les droits de création de composants
- Disposer du fichier numérique à intégrer

## Scénario

**Étape initiale :** `myr model add firmware.bin --name "Firmware v2" --channel greenchannel --category base` est exécutée (ou l'appel API équivalent), pour le compte du Concepteur

### Flux nominal — Composant numérique nouveau

1. Le type "Numérique" (Software) est désigné
2. Le fichier numérique est importé
3. Le système vérifie le hash du fichier
4. Les métadonnées sont renseignées (nom, licence, version, auteur)
5. La transaction est soumise sur la blockchain

### Flux erreur — Doublon détecté

1. Erreur métier : "Composant numérique déjà enregistré" (RM01)

## Post-conditions

- Le composant numérique est enregistré sur le réseau
- Un UUID unique lui est attribué

## Diagramme d'activités

```plantuml
@startuml
skin rose
title Ajout d'un composant Numérique
start
:Transmettre le fichier numérique et les métadonnées (myr model add);
:Vérifier le hash du fichier;
if (Doublon détecté?) then (oui)
  :Retourner l'erreur "Composant numérique déjà enregistré";
  stop
else (non)
  :Soumettre la transaction sur la blockchain;
  stop
endif
@enduml
```

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCCE — Composant Ecriture](../../Carte_des_specs.md#UCCE%20—%20Composant%20Ecriture)
- [UCCE03 — couche analyse](../../2-Analyse/UCCE-Composant_Ecriture/UCCE03.md)
- [Traçabilité UCCE03 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCCE03)

**Exigences fonctionnelles couvertes**
- [EF11 — Enregistrer un composant numérique sur la blockchain](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Règles métier**
- [RM01 — Anti-plagiat obligatoire](../Regles_Metier.md#1.%20Assets%20et%20composants)

**Cité par**
- [Expression_des_besoins_Intro](../Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../Matrice_Tracabilite.md)
- [UCAUT03 (expression)](../UCAUT-Automatisation/UCAUT03.md)
- [todo (expression)](../todo.md)
- [Analyse_des_besoins](../../2-Analyse/Analyse_des_besoins.md)
- [UCAUT03 (analyse)](../../2-Analyse/UCAUT-Automatisation/UCAUT03.md)
- [UCCE06 (analyse)](../../2-Analyse/UCCE-Composant_Ecriture/UCCE06.md)
- [todo (analyse)](../../2-Analyse/todo.md)
- [API_REST](../../3-Conception/API_REST.md)
- [Architecture_Composition](../../3-Conception/Architecture_Composition.md)
- [Chaincode](../../3-Conception/Chaincode.md)
- [Conception_intro](../../3-Conception/Conception_intro.md)
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
