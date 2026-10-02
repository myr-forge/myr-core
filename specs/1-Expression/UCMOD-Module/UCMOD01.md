---
categorie: Module
titre: "Créer un Module"
probabilite: 3
impact: 5
importance: 15
etat: relire
tags:
  - couche/expression
  - type/use-case
  - famille/UCMOD
  - domaine/model
  - uc/UCMOD01
---

# Créer un Module

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Concepteur" as C

rectangle "Application MYR" {
    usecase "Créer un module (draft)" as UC1
    usecase "Créer liaisons entre composants" as UC2
}

C --> UC1
UC1 ..> UC2 : <<include>>

@enduml
```

## Contexte

Assemblage de plusieurs composants ou de Modules existants selon les compatibilités de leurs interfaces pour créer un Module à part entière.

Un module nouvellement créé est en état **draft** (brouillon) : il existe localement mais n'est pas encore ancré sur la blockchain. L'état draft est modifiable à volonté — le concepteur peut ajouter, retirer ou reconfigurer des liaisons autant de fois que nécessaire avant de publier. Un nom par défaut lui est attribué à la création (`<identité du concepteur>_<date>_<heure>`), modifiable ensuite (voir UCMOD03).

Lorsque le module est prêt, la soumission à la blockchain (voir UCMOD06) est une étape distincte et explicite. Elle crée une **ModuleVersion** immuable — snapshot figé et horodaté de l'assemblage, non modifiable après publication. Toute évolution ultérieure nécessite la création d'une nouvelle version.

Un module peut aussi être créé à partir d'un module déjà soumis (`--parent-id`), pour reprendre un travail existant sans reconstruire manuellement l'assemblage — sa composition (instances et liaisons internes) est alors dupliquée dans le nouveau brouillon.

## Pré-conditions

- Être connecté au réseau
- Avoir les droits de création
- Avoir des composants ou modules existants

## Scénario

**Étape initiale :** `myr module create --name <nom> --channel <id> [--owner-id <id>] [--description <texte>] [--license <id>]` est exécutée (ou l'appel API équivalent), pour le compte du Concepteur

### Flux nominal — Module créé (draft)

1. Le module est initialisé en état **draft**
2. Les liaisons entre les interfaces compatibles des composants sont créées (`myr model instance add` puis `myr model link add`, voir UCAM01)
3. Le module est nommé et configuré
4. Le module est enregistré en état **draft**

### Flux alternatif — Dérivation d'un module existant

1. `myr module create --parent-id <id> [--license <id>]` référence un module déjà soumis
2. Si une licence est fournie, sa compatibilité avec celle du module parent est vérifiée
3. Le nouveau module est créé en état **draft**, avec la composition complète du module parent (instances et liaisons) déjà dupliquée dedans

### Flux erreur — Aucune liaison créée

1. Le service refuse de nommer ou sauvegarder un module sans au moins un assemblage
2. Message d'erreur : "Ajoutez au moins une liaison entre composants"

### Flux erreur — Licence incompatible avec le module parent

1. Une dérivation est demandée avec une licence incompatible avec celle du module parent
2. Le service refuse la création : "Incompatibilité de licence"

## Post-conditions

- Le module est en état **draft**
- Les interfaces libres du module sont visibles
- Le module n'est pas encore visible sur le réseau (soumission requise — voir UCMOD06)

## Diagramme d'activités

```plantuml
@startuml
skin rose
title Créer un Module
start
:Créer le module (myr module create);
if (Liaisons entre composants créées?) then (oui)
  :Créer les liaisons entre interfaces compatibles (myr model link add);
  :Nommer et configurer le module;
  :Enregistrer le module en état draft;
  stop
else (non)
  :Refuser — "Ajoutez au moins une liaison entre composants";
  stop
endif
@enduml
```

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCMOD — Module](../../Carte_des_specs.md#UCMOD%20—%20Module)
- [UCMOD01 — couche analyse](../../2-Analyse/UCMOD-Module/UCMOD01.md)
- [Traçabilité UCMOD01 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCMOD01)

**Exigences fonctionnelles couvertes**
- [EF16 — Vérifier la compatibilité de licence lors d'une dérivation](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)
- [EF26 — Assembler plusieurs composants en module (état draft), y compris par dérivation d'un module existant (composition dupliquée depuis un `parent_id`)](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Use cases cités**
- [UCAM01 — Liaison entre interfaces](../UCAM-Assemblage_Module/UCAM01.md)
- [UCMOD03 — Modifier les métadonnées d'un Module](UCMOD03.md)
- [UCMOD06 — Soumettre un module à la blockchain](UCMOD06.md)

**Cité par**
- [Expression_des_besoins_Intro](../Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../Matrice_Tracabilite.md)
- [UCMOD03 (expression)](UCMOD03.md)
- [UCMOD06 (expression)](UCMOD06.md)
- [UCMOD08 (expression)](UCMOD08.md)
- [todo (expression)](../todo.md)
- [Analyse_des_besoins](../../2-Analyse/Analyse_des_besoins.md)
- [UCAM05 (analyse)](../../2-Analyse/UCAM-Assemblage_Module/UCAM05.md)
- [UCDEV02 (analyse)](../../2-Analyse/UCDEV-Developpement/UCDEV02.md)
- [UCMOD02 (analyse)](../../2-Analyse/UCMOD-Module/UCMOD02.md)
- [UCMOD03 (analyse)](../../2-Analyse/UCMOD-Module/UCMOD03.md)
- [UCMOD06 (analyse)](../../2-Analyse/UCMOD-Module/UCMOD06.md)
- [UCMOD07 (analyse)](../../2-Analyse/UCMOD-Module/UCMOD07.md)
- [UCMOD08 (analyse)](../../2-Analyse/UCMOD-Module/UCMOD08.md)
- [todo (analyse)](../../2-Analyse/todo.md)
- [API_REST](../../3-Conception/API_REST.md)
- [Architecture_Composition](../../3-Conception/Architecture_Composition.md)
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)
- [Sequence_soumission_module](../../3-Conception/Sequence_soumission_module.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
