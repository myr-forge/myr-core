---
categorie: Assemblage Module
titre: "Liaison entre interfaces"
probabilite: 4
impact: 5
importance: 20
etat: relire
tags:
  - couche/expression
  - type/use-case
  - famille/UCAM
  - domaine/model
  - uc/UCAM01
  - rm/RM09
  - rm/RM10
  - rm/RM11
---

# Liaison entre interfaces

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Concepteur" as C

rectangle "Application MYR" {
    usecase "Créer une liaison" as UC1
    usecase "Vérifier compatibilité des interfaces" as UC2
    usecase "Choisir un asset d'accroche" as UC3
}

C --> UC1
UC1 ..> UC2 : <<include>>
UC1 .> UC3 : <<extend>>

@enduml
```

## Contexte

Création de nouvelles interfaces d'un composant par la création de liaisons entre composants. Le composant doit garder ses interfaces créées par liaisons.

Une liaison peut être **directe** (les interfaces se connectent sans intermédiaire) ou **via un asset d'accroche** (vis, câble, connecteur…). L'asset d'accroche est lui-même un composant existant sur le réseau — son identifiant (`FastenerAssetID`) est enregistré sur la liaison. Voir UCAM07.

## Pré-conditions

- Être connecté au réseau
- Avoir au moins deux composants avec des interfaces existants
- Les interfaces à relier doivent être compatibles (catégorie, sens, tag, type, valeur)

## Scénario

**Étape initiale :** Sur le serveur (SSH), l'administrateur exécute `myr model link add --from <ifaceID_source> --to <ifaceID_cible> [--fastener <assetID>]` pour le compte du Concepteur — même effet via l'API REST équivalente (identifiants d'interfaces obtenus via `myr model interface list`, voir UCAM02)

### Flux nominal — Liaison directe

1. Les identifiants des interfaces source et cible sont transmis
2. Le service vérifie la compatibilité de la paire d'interfaces (catégorie, sens, tag, type, valeur — RM10/RM11)
3. Aucun asset d'accroche n'est précisé : la liaison est enregistrée directement (`FastenerAssetID` vide)
4. Les interfaces du composant résultant sont créées et conservées

### Flux nominal — Liaison via asset d'accroche

1. Étapes 1 à 2 identiques au flux nominal précédent
2. Un asset d'accroche est précisé (`--fastener <assetID>`, voir UCAM07)
3. La liaison est enregistrée avec le `FastenerAssetID` de l'asset choisi

### Flux erreur — Interfaces incompatibles

1. La paire d'interfaces transmise ne respecte pas les critères de compatibilité (RM10/RM11)
2. Le service refuse la création de la liaison — message d'erreur métier explicite

### Flux erreur — Interface déjà utilisée

1. L'interface cible est déjà engagée dans une liaison existante (RM09)
2. Le service refuse la création : "Cette interface est déjà utilisée dans une liaison"

### Flux — Liaison devenue incompatible après modification

Ce scénario survient lorsqu'un asset impliqué dans une liaison est remplacé par une version dont les interfaces ont changé (ex. : passage à une version améliorée avec un type de port différent).

1. La modification de l'asset est enregistrée
2. Le système détecte que la liaison existante n'est plus compatible avec les nouvelles interfaces
3. La liaison n'est **pas supprimée** — elle passe en état `incompatible`
4. Elle reste consultable, marquée `Incompatible: true`
5. Elle peut être supprimée manuellement, ou les assets peuvent être adaptés

## Post-conditions

- La liaison est enregistrée entre les deux composants
- Les interfaces libres du composant résultant sont visibles
- Une liaison devenue incompatible après modification reste présente, marquée `Incompatible: true` — elle n'est jamais supprimée automatiquement

## Diagrammes

### Principe d'une liaison entre deux interfaces compatibles

```plantuml
@startuml
skin rose
(C2) <-- (L2♀) : interface
(C1) <-- (L2♂) : interface

(L2♂) -- (L2♀) : liaison
@enduml
```

### Exemple concret — Raspberry Pi, caméra, écran, boitier

```plantuml
@startuml
title Exemple de liaison d'interface
skin rose

component rpi {
  port "trou" as T1
  port "trou" as T2
  port "trou" as T3
  port "trou" as T4
  port "port_camera"  as PC
  port "port_ecran"   as PE
}

component camera {
  port "port"   as PCAM
}

component ecran {
  port "port"   as PECR
}

component CableCSI {
  port "port"   as PCSI1
  port "port"   as PCSI2
}

component boitier {
  port "port"   as PB1
  port "port"   as PB2
  port "port"   as PB3
  port "port"   as PB4
}
together {
  component vis1 {
    port "male"   as PM1
  }
  component vis2 {
    port "male"   as PM2
  }
  component vis3 {
    port "male"   as PM3
  }
  component vis4 {
    port "male"   as PM4
  }
}

PM1 --> T1 : MECA
T1 -down-> PB1 : MECA

PM2 --> T2 : MECA
T2 --> PB2 : MECA

PM3 --> T3 : MECA
T3 --> PB3 : MECA

PM4 --> T4 : MECA
T4 --> PB4 : MECA

PE --> PCSI1 : ELEC
PCSI2 --> PECR : ELEC

PC --> PCAM : ELEC

@enduml
```

### Diagramme d'activités

```plantuml
@startuml
skin rose
title Liaison entre interfaces
start
:Transmettre les identifiants des interfaces source et cible (myr model link add);
if (Interface cible déjà utilisée?) then (oui)
  :Refuser — "Cette interface est déjà utilisée dans une liaison";
  stop
else (non)
  :Vérifier la compatibilité de la paire d'interfaces (RM10/RM11);
  if (Compatible?) then (non)
    :Refuser la création de la liaison;
    stop
  else (oui)
    if (Asset d'accroche précisé (--fastener)?) then (oui)
      :Enregistrer la liaison avec FastenerAssetID;
      stop
    else (non)
      :Enregistrer la liaison directe;
      stop
    endif
  endif
endif
@enduml
```

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCAM — Assemblage Module](../../Carte_des_specs.md#UCAM%20—%20Assemblage%20Module)
- [UCAM01 — couche analyse](../../2-Analyse/UCAM-Assemblage_Module/UCAM01.md)
- [Traçabilité UCAM01 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCAM01)

**Exigences fonctionnelles couvertes**
- [EF18 — Créer des liaisons entre interfaces compatibles de composants](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Use cases cités**
- [UCAM02 — Visualiser les interfaces physiques de composants](UCAM02.md)
- [UCAM07 — Choisir un asset d'accroche (Fastener)](UCAM07.md)

**Règles métier**
- [RM09 — Interface à usage unique](../Regles_Metier.md#3.%20Interfaces%20et%20liaisons)
- [RM10 — Vérification de compatibilité automatique](../Regles_Metier.md#3.%20Interfaces%20et%20liaisons)
- [RM11 — Critères de compatibilité d'interfaces](../Regles_Metier.md#3.%20Interfaces%20et%20liaisons)

**Cité par**
- [Expression_des_besoins_Intro](../Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../Matrice_Tracabilite.md)
- [UCAM05 (expression)](UCAM05.md)
- [UCAM07 (expression)](UCAM07.md)
- [UCCE06 (expression)](../UCCE-Composant_Ecriture/UCCE06.md)
- [UCMOD01 (expression)](../UCMOD-Module/UCMOD01.md)
- [UCMOD02 (expression)](../UCMOD-Module/UCMOD02.md)
- [todo (expression)](../todo.md)
- [Analyse_des_besoins](../../2-Analyse/Analyse_des_besoins.md)
- [UCAM02 (analyse)](../../2-Analyse/UCAM-Assemblage_Module/UCAM02.md)
- [UCAM05 (analyse)](../../2-Analyse/UCAM-Assemblage_Module/UCAM05.md)
- [UCAM07 (analyse)](../../2-Analyse/UCAM-Assemblage_Module/UCAM07.md)
- [UCCE05 (analyse)](../../2-Analyse/UCCE-Composant_Ecriture/UCCE05.md)
- [UCCE06 (analyse)](../../2-Analyse/UCCE-Composant_Ecriture/UCCE06.md)
- [UCDEV02 (analyse)](../../2-Analyse/UCDEV-Developpement/UCDEV02.md)
- [UCMOD02 (analyse)](../../2-Analyse/UCMOD-Module/UCMOD02.md)
- [UCMOD03 (analyse)](../../2-Analyse/UCMOD-Module/UCMOD03.md)
- [UCMOD06 (analyse)](../../2-Analyse/UCMOD-Module/UCMOD06.md)
- [todo (analyse)](../../2-Analyse/todo.md)
- [Architecture_Composition](../../3-Conception/Architecture_Composition.md)
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)
- [Sequence_soumission_module](../../3-Conception/Sequence_soumission_module.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
