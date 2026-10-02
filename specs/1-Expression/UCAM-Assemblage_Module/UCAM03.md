---
categorie: Assemblage Module
titre: "Créer une interface sur un composant"
probabilite: 3
impact: 5
importance: 15
etat: relu
tags:
  - couche/expression
  - type/use-case
  - famille/UCAM
  - domaine/model
  - uc/UCAM03
  - rm/RM13
---
# Créer une interface sur un composant

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Concepteur" as C

rectangle "Application MYR" {
    usecase "Créer une interface (connexion virtuelle)" as UC1
    usecase "Créer une interface (attributs explicites)" as UC2
    usecase "Définir les attributs de l'interface" as UC3
    usecase "Enregistrer l'interface (draft)" as UC4
}

C --> UC1
C --> UC2
UC1 ..> UC3 : <<include>>
UC2 ..> UC3 : <<include>>
UC3 ..> UC4 : <<include>>

@enduml
```

## Contexte

Chaque composant possède toujours un **slot d'interface virtuel** (`Virtual=true`) : dès qu'un slot virtuel est matérialisé (relié à une interface physique existante), un nouveau slot virtuel est recréé automatiquement sur le composant (RM13).

Une interface peut être créée de deux façons :

- **Connexion depuis un slot virtuel** : le slot virtuel d'un composant est relié à une interface physique existante d'un autre composant. Le service déduit automatiquement la plupart des propriétés de la nouvelle interface à partir de la cible — seul le tag doit être renseigné explicitement.
- **Création manuelle** : tous les attributs de la nouvelle interface (catégorie, sens, tag, type, valeur/plage, unité) sont renseignés explicitement.

## Pré-conditions

- Être connecté au réseau
- Avoir au moins un composant existant
- Avoir les droits de création ou d'édition sur le composant (défini dans le rôle)

## Scénario

### Flux nominal A — Connexion depuis un slot virtuel (déduction automatique)

**Étape initiale :** `myr model link connect-virtual --virtual-iface <id> --physical-iface <id>` est exécutée (ou l'appel API équivalent)

1. Le service déduit les propriétés de la nouvelle interface à partir de la cible :
   - **catégorie** : identique à la cible
   - **sens** : inversé (sortie → entrée ; entrée → sortie ; bidirectionnel → bidirectionnel)
   - **type** : identique à la cible
   - **valeur/unité** : inférée depuis la cible (ex : cible sortie 3–6 V → nouvelle interface entrée 3,3 V)
2. Le tag doit être renseigné explicitement (`--tag <tag>`) — jamais déduit
3. La nouvelle interface est créée sur le composant source (état draft)
4. La liaison est enregistrée
5. Le slot virtuel est automatiquement recréé sur le composant (RM13)

### Flux nominal B — Création manuelle avec attributs explicites

**Étape initiale :** `myr model interface add <assetID> --category <cat> --type <type> --direction <in|out|bidir> [--value-min --value-max --unit] --tag <tag> [--name <label>]` est exécutée

1. L'interface est créée et enregistrée dans l'état draft

### Flux alternatif A2 — Ajustement des valeurs déduites

1. Les valeurs déduites (flux A) sont surchargées par des flags explicites avant validation (ex : affiner la plage de valeur)
2. L'interface est créée avec les valeurs ajustées

### Flux erreur — Incompatibilité de connexion virtuelle

1. Le service détecte une incompatibilité entre le slot virtuel et l'interface physique ciblée (catégories différentes)
2. La connexion est refusée

## Post-conditions

- La nouvelle interface est enregistrée sur le composant
- Le slot d'interface virtuel reste disponible sur le composant (`Virtual=true` recréé)
- La liaison est enregistrée (si flux A)

## Diagramme d'activités

```plantuml
@startuml
skin rose
title Créer une interface sur un composant
start
if (Méthode de création?) then (connexion virtuelle)
  :Transmettre virtual-iface et physical-iface cible;
  if (Catégorie compatible?) then (oui)
    :Déduire catégorie, sens inversé, type, valeur/unité;
    :Renseigner le tag explicitement (jamais déduit);
    :Créer l'interface sur le composant source (draft);
    :Enregistrer la liaison;
    :Recréer automatiquement le slot virtuel;
    stop
  else (non)
    :Refuser la connexion;
    stop
  endif
else (attributs explicites)
  :Transmettre les attributs (catégorie, sens, tag, type, valeur, unité);
  :Enregistrer l'interface (draft);
  stop
endif
@enduml
```

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCAM — Assemblage Module](../../Carte_des_specs.md#UCAM%20—%20Assemblage%20Module)
- [UCAM03 — couche analyse](../../2-Analyse/UCAM-Assemblage_Module/UCAM03.md)
- [Traçabilité UCAM03 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCAM03)

**Exigences fonctionnelles couvertes**
- [EF20 — Définir une interface sur un composant](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)
- [EF25 — Garantir un slot virtuel disponible sur chaque asset](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Règles métier**
- [RM13 — Slot virtuel garanti](../Regles_Metier.md#3.%20Interfaces%20et%20liaisons)

**Cité par**
- [Expression_des_besoins_Intro](../Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../Matrice_Tracabilite.md)
- [UCAM05 (expression)](UCAM05.md)
- [todo (expression)](../todo.md)
- [Analyse_des_besoins](../../2-Analyse/Analyse_des_besoins.md)
- [UCAM01 (analyse)](../../2-Analyse/UCAM-Assemblage_Module/UCAM01.md)
- [UCAM02 (analyse)](../../2-Analyse/UCAM-Assemblage_Module/UCAM02.md)
- [UCAM05 (analyse)](../../2-Analyse/UCAM-Assemblage_Module/UCAM05.md)
- [UCCE01 (analyse)](../../2-Analyse/UCCE-Composant_Ecriture/UCCE01.md)
- [UCCE06 (analyse)](../../2-Analyse/UCCE-Composant_Ecriture/UCCE06.md)
- [UCDEV02 (analyse)](../../2-Analyse/UCDEV-Developpement/UCDEV02.md)
- [todo (analyse)](../../2-Analyse/todo.md)
- [Architecture_Composition](../../3-Conception/Architecture_Composition.md)
- [Conception_intro](../../3-Conception/Conception_intro.md)
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)
- [Sequence_soumission_asset](../../3-Conception/Sequence_soumission_asset.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
