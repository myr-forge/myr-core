---
categorie: Administration
titre: "Ajouter un nœud à un réseau existant"
etat: "RELIRE"
tags:
  - couche/expression
  - type/use-case
  - famille/UCADM
  - domaine/channel
  - domaine/identity
  - domaine/network
  - uc/UCADM03
---

# Ajouter un nœud à un réseau existant

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Administrateur" as ADM
actor "Organisation" as ORG

rectangle "Infrastructure MYR" {
    usecase "Ajouter un nœud au réseau" as UC1
    usecase "Valider l'ajout du nœud" as UC2
    usecase "Synchroniser le ledger" as UC3
}

ADM --> UC1
ORG --> UC2
UC1 ..> UC2 : <<include>>
UC2 ..> UC3 : <<include>>

@enduml
```

## Contexte

Un nœud (peer) supplémentaire peut être ajouté à un réseau existant pour étendre sa capacité et sa résilience. Le peer doit être validé selon les règles établies par le réseau.

## Pré-conditions

- Être administrateur du réseau
- Réseau existant et opérationnel
- Serveur disponible avec IP fixe publique (port 7051 ouvert)

## Scénario

**Étape initiale :** L'administrateur configure le nouveau peer sur son serveur

### Flux nominal — Peer ajouté avec succès

1. L'administrateur soumet une demande d'ajout au réseau
2. Les organisations existantes valident l'ajout selon la politique du réseau
3. Le peer est synchronisé avec le ledger existant
4. Le peer est déclaré actif et participe au consensus

### Flux alternatif — Nœud de type orderer

1. L'administrateur choisit d'ajouter un nœud orderer (et non un peer)
2. Il configure les paramètres spécifiques aux orderers (consensus Raft, block cutting parameters)
3. Le nœud orderer rejoint le canal système du réseau
4. Confirmation : "Nœud orderer ajouté au réseau"

### Flux erreur — Peer déconnecté temporairement

1. Le peer se désynchronise pendant la déconnexion
2. À la reconnexion (avant délai maximum), le peer se resynchronise automatiquement

### Flux erreur — Peer déconnecté trop longtemps

1. Le délai maximum de déconnexion est dépassé
2. Le peer est marqué comme mort et supprimé du réseau
3. À la reconnexion, un message propose de reconstruire un nouveau peer

## Post-conditions

- Le nouveau nœud est actif et synchronisé avec le réseau
- Il contribue à la résilience et au consensus du réseau

## Diagramme d'activités

```plantuml
@startuml
skin rose
title Ajouter un nœud à un réseau existant
start
if (Type de nœud : orderer?) then (oui)
  :Configurer le nœud orderer (Raft, block cutting parameters);
  :Rejoindre le canal système du réseau;
  :Confirmer "Nœud orderer ajouté";
  stop
else (non)
  :Configurer le nouveau peer sur le serveur;
  :Soumettre une demande d'ajout au réseau;
  :Les organisations existantes valident l'ajout;
  :Synchroniser le peer avec le ledger existant;
  :Déclarer le peer actif;
endif
if (Peer se déconnecte?) then (oui)
  if (Délai maximum dépassé?) then (oui)
    :Marquer le peer comme mort;
    :Supprimer le peer du réseau;
    :Proposer de reconstruire un nouveau peer à la reconnexion;
    stop
  else (non)
    :Resynchroniser automatiquement à la reconnexion;
    stop
  endif
else (non)
  :Peer actif et participant au consensus;
  stop
endif
@enduml
```

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCADM — Administration](../../Carte_des_specs.md#UCADM%20—%20Administration)
- [UCADM03 — couche analyse](../../2-Analyse/UCADM-Administration/UCADM03.md)
- [Traçabilité UCADM03 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCADM03)

**Exigences fonctionnelles couvertes**
- [EF09 — Étendre un réseau avec de nouveaux nœuds (peer ou orderer)](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Cité par**
- [Expression_des_besoins_Intro](../Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../Matrice_Tracabilite.md)
- [UCADM04 (expression)](UCADM04.md)
- [UCDEV02 (expression)](../UCDEV-Developpement/UCDEV02.md)
- [todo (expression)](../todo.md)
- [Analyse_des_besoins](../../2-Analyse/Analyse_des_besoins.md)
- [UCADM04 (analyse)](../../2-Analyse/UCADM-Administration/UCADM04.md)
- [UCADM05 (analyse)](../../2-Analyse/UCADM-Administration/UCADM05.md)
- [UCDEV02 (analyse)](../../2-Analyse/UCDEV-Developpement/UCDEV02.md)
- [todo (analyse)](../../2-Analyse/todo.md)
- [API_REST](../../3-Conception/API_REST.md)
- [Conception_intro](../../3-Conception/Conception_intro.md)
- [DC_CLI_Admin](../../3-Conception/DC_CLI_Admin.md)
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)
- [DC_D2_Administration](../../3-Conception/DC_D2_Administration.md)
- [todo (conception)](../../3-Conception/todo.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
