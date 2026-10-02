---
categorie: Assemblage Module
titre: "Retirer une instance de composant d'un Module"
probabilite: 4
impact: 4
importance: 16
etat: relire
tags:
  - couche/expression
  - type/use-case
  - famille/UCAM
  - domaine/model
  - uc/UCAM08
  - rm/RM15
---

# Retirer une instance de composant d'un Module

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Concepteur" as C

rectangle "Application MYR" {
    usecase "Retirer une instance de composant d'un module" as UC1
    usecase "Supprimer les liaisons en cascade" as UC2
}

C --> UC1
UC1 ..> UC2 : <<include>>

@enduml
```

## Contexte

Une instance de composant ou de module ajoutée à un module hôte peut en être retirée par une action directe (`RemoveAssetFromWorkspace` au niveau du domaine). Cette opération supprime **uniquement l'instance** : le composant reste disponible sur le réseau blockchain et peut être rajouté à tout moment.

Toutes les liaisons impliquant cette instance sont supprimées automatiquement en cascade — elles ne peuvent pas rester orphelines.

## Pré-conditions

- Connaître l'identifiant du module hôte et de l'instance à retirer (`myr module get <id>` pour lister ses instances)

## Scénario

**Étape initiale :** `myr model instance remove <moduleID> <instanceID>` est exécutée (ou l'appel API équivalent), sans étape de confirmation interactive

### Flux nominal — Retrait sans liaisons actives

1. L'instance ciblée n'a aucune liaison enregistrée
2. Elle est retirée immédiatement du module

### Flux nominal — Retrait avec liaisons en cascade

1. L'instance ciblée possède des liaisons avec d'autres instances du module
2. Toutes les liaisons impliquant cette instance sont supprimées en cascade (RM15)
3. L'instance est retirée du module

## Post-conditions

- Le composant n'est plus une instance du module hôte
- Toutes ses liaisons sont supprimées (cascade)
- Le composant reste disponible sur le réseau et peut être réajouté (voir UCMOD02)
- **Aucune transaction blockchain n'est émise** — Fabric ne supporte pas la suppression d'asset

## Diagramme

### Effet cascade sur les liaisons

```plantuml
@startuml
skin rose
title Retrait du composant C2 → cascade sur ses liaisons

(C1) -- (C2) : liaison L1 ← supprimée
(C2) -- (C3) : liaison L2 ← supprimée
(C3) -- (C4) : liaison L3 ← conservée

note bottom of (C2) : instance retirée
@enduml
```

### Diagramme d'activités

```plantuml
@startuml
skin rose
title Retirer une instance de composant d'un Module
start
:Transmettre moduleID et instanceID (myr model instance remove);
if (Instance a des liaisons actives?) then (oui)
  :Supprimer toutes les liaisons en cascade (RM15);
  :Retirer l'instance du module;
  stop
else (non)
  :Retirer l'instance du module immédiatement;
  stop
endif
@enduml
```

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCAM — Assemblage Module](../../Carte_des_specs.md#UCAM%20—%20Assemblage%20Module)
- [UCAM08 — couche analyse](../../2-Analyse/UCAM-Assemblage_Module/UCAM08.md)
- [Traçabilité UCAM08 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCAM08)

**Exigences fonctionnelles couvertes**
- [EF24 — Retirer une instance de composant d'un module (avec cascade des connexions)](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Use cases cités**
- [UCMOD02 — Ajouter un Module existant](../UCMOD-Module/UCMOD02.md)

**Règles métier**
- [RM15 — Instance indépendante](../Regles_Metier.md#4.%20Composition%20d'un%20Module%20%28instances%29)

**Cité par**
- [Expression_des_besoins_Intro](../Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../Matrice_Tracabilite.md)
- [todo (expression)](../todo.md)
- [UCCE07 (analyse)](../../2-Analyse/UCCE-Composant_Ecriture/UCCE07.md)
- [UCDEV02 (analyse)](../../2-Analyse/UCDEV-Developpement/UCDEV02.md)
- [todo (analyse)](../../2-Analyse/todo.md)
- [API_REST](../../3-Conception/API_REST.md)
- [Architecture_Composition](../../3-Conception/Architecture_Composition.md)
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)
- [Sequence_soumission_module](../../3-Conception/Sequence_soumission_module.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
