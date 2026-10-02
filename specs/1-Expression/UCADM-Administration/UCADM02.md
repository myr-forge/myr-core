---
categorie: Administration
titre: "Créer un réseau indépendant"
etat : "a lire"
tags:
  - couche/expression
  - type/use-case
  - famille/UCADM
  - domaine/channel
  - domaine/identity
  - domaine/network
  - uc/UCADM02
---
# Créer un réseau indépendant

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Administrateur" as ADM

rectangle "Infrastructure MYR" {
    usecase "Créer un réseau indépendant" as UC1
    usecase "Configurer les nœuds et identités" as UC2
    usecase "Démarrer et synchroniser le réseau" as UC3
}

ADM --> UC1
UC1 ..> UC2 : <<include>>
UC1 ..> UC3 : <<include>>

@enduml
```

## Contexte

Un administrateur d'un organisme peut créer un nouveau réseau MYR indépendant. Il choisit la technologie blockchain sous-jacente selon les besoins de son déploiement (performance, confidentialité, disponibilité des nœuds, contraintes légales ou obsolescence d'une technologie). Le réseau nécessite un nombre minimum de nœuds pour garantir sa résilience et sa disponibilité.

## Pré-conditions

- Avoir les droits d'administration de l'infrastructure
- Disposer d'au moins 3 serveurs accessibles avec une adresse réseau stable
- L'environnement serveur est prêt à accueillir des nœuds pour la technologie blockchain choisie

## Scénario

**Étape initiale :** L'administrateur initie la création du réseau

### Flux nominal — Réseau créé

1. Il choisit la technologie blockchain à utiliser pour ce réseau
2. Il déclare les organisations participantes et leurs rôles dans le réseau
3. Il configure la topologie du réseau (nœuds, politiques de consensus, règles d'accès)
4. Il génère et distribue les identités cryptographiques des participants
5. Il démarre le réseau et vérifie que les nœuds sont synchronisés
6. Il fournit aux organisations les informations de connexion leur permettant de rejoindre le réseau

### Flux alternatif — Import d'un profil de connexion existant

1. Au lieu de configurer le réseau manuellement, l'administrateur importe un fichier de profil de connexion (JSON ou YAML)
2. Le système détecte le format et la technologie blockchain depuis le profil, et pré-remplit automatiquement les champs (organisations, nœuds, politiques)
3. L'administrateur vérifie les informations importées et corrige si nécessaire
4. Le réseau est créé à partir du profil importé

## Post-conditions

- Le réseau indépendant est opérationnel avec le nombre minimum de nœuds synchronisés
- Les organisations déclarées peuvent s'y connecter

## Diagramme d'activités

```plantuml
@startuml
skin rose
title Créer un réseau indépendant
start
:Initier la création du réseau;
if (Import d'un profil de connexion existant?) then (oui)
  :Importer le fichier de profil (JSON ou YAML);
  :Détecter la technologie blockchain et parser le profil;
  :Vérifier et corriger les informations importées;
else (non)
  :Choisir la technologie blockchain;
  :Déclarer les organisations participantes et leurs rôles;
  :Configurer la topologie (nœuds, politiques de consensus, règles d'accès);
endif
:Générer et distribuer les identités cryptographiques;
:Démarrer le réseau;
:Vérifier la synchronisation des nœuds;
:Fournir les informations de connexion aux organisations;
stop
@enduml
```

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCADM — Administration](../../Carte_des_specs.md#UCADM%20—%20Administration)
- [UCADM02 — couche analyse](../../2-Analyse/UCADM-Administration/UCADM02.md)
- [Traçabilité UCADM02 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCADM02)

**Exigences fonctionnelles couvertes**
- [EF07 — Créer un réseau blockchain indépendant](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Cité par**
- [Expression_des_besoins_Intro](../Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../Matrice_Tracabilite.md)
- [UCADM05 (expression)](UCADM05.md)
- [UCDEV02 (expression)](../UCDEV-Developpement/UCDEV02.md)
- [todo (expression)](../todo.md)
- [Analyse_des_besoins](../../2-Analyse/Analyse_des_besoins.md)
- [UCADM01 (analyse)](../../2-Analyse/UCADM-Administration/UCADM01.md)
- [UCADM03 (analyse)](../../2-Analyse/UCADM-Administration/UCADM03.md)
- [UCADM04 (analyse)](../../2-Analyse/UCADM-Administration/UCADM04.md)
- [UCADM05 (analyse)](../../2-Analyse/UCADM-Administration/UCADM05.md)
- [UCDEV02 (analyse)](../../2-Analyse/UCDEV-Developpement/UCDEV02.md)
- [todo (analyse)](../../2-Analyse/todo.md)
- [DC_CLI_Admin](../../3-Conception/DC_CLI_Admin.md)
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)
- [DC_D2_Administration](../../3-Conception/DC_D2_Administration.md)
- [todo (conception)](../../3-Conception/todo.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
