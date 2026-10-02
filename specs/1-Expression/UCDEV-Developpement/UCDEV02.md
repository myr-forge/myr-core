---
categorie: Développement autour de MYR
titre: "Utilisation du CLI"
probabilite: 1
importance: 0
etat: relire
tags:
  - couche/expression
  - type/use-case
  - famille/UCDEV
  - uc/UCDEV02
---

# Utilisation du CLI

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Administrateur" as ADM

rectangle "CLI MYR (Serveur)" {
    usecase "Administrer le réseau (UCADM01-05)" as UC1
    usecase "Exécuter une action domaine\npour le compte d'une identité\n(asset, interface, liaison, module...)" as UC2
}

ADM --> UC1
ADM --> UC2

@enduml
```

## Contexte

Le CLI `myr` n'est pas un outil d'administration réseau au sens strict : il expose la **totalité** des capacités du domaine (`domain/`), au même titre que l'API REST — le CLI n'est pas un citoyen de seconde zone par rapport à l'API. Toute action métier exposée par l'API (créer un composant, définir une interface, créer une liaison, assembler et soumettre un module, rechercher un asset…) a donc une commande `myr` équivalente, en plus des commandes propres à l'administration réseau (`myr network/org/node`, UCADM01–05).

Ce qui distingue le CLI de l'API REST n'est pas le périmètre fonctionnel mais **qui peut l'exécuter et où** : conformément au principe d'exécution distante, `myr` (CLI et serveur d'API) ne tourne jamais sur le poste d'un utilisateur final — uniquement sur le serveur, via une connexion SSH de l'administrateur. Un Concepteur ou un Consommateur n'ouvre donc jamais lui-même un terminal `myr` ; c'est l'administrateur du serveur qui, le cas échéant, exécute une commande domaine **pour le compte** d'une identité (scripts d'import en masse, migration, support, restauration après incident).

Le Développeur (UCDEV01), qui n'a pas d'accès SSH au serveur, interagit avec MYR exclusivement via l'API REST — un client tiers (interface graphique, boutique partenaire, plugin CAO...) fait de même. Ce n'est pas une limitation fonctionnelle du CLI, mais une conséquence du modèle d'exécution distante.

## Pré-conditions

- Être connecté sur le serveur distant (SSH)

## Scénario

**Étape initiale :** L'administrateur ouvre un terminal sur le serveur distant

### Flux nominal — Commande d'administration réseau

1. L'administrateur saisit une commande d'administration (ex : `myr network add`, `myr org add`, `myr node add` — voir UCADM01–05)
2. La commande est transmise au réseau
3. Le résultat est affiché en sortie texte dans le terminal

### Flux nominal — Commande domaine (parité API/CLI)

1. L'administrateur saisit une commande domaine pour le compte d'une identité (ex : `myr model add`, `myr model interface add`, `myr model link add`, `myr module submit` — voir les use cases UCCE/UCAM/UCMOD correspondants)
2. La commande appelle le même service domaine que le handler REST équivalent — les mêmes règles métier s'appliquent (anti-plagiat, compatibilité de licence, compatibilité d'interfaces…)
3. Le résultat est affiché en sortie texte dans le terminal

### Flux erreur — Commande invalide

1. Un message d'aide est affiché avec la syntaxe correcte

### Flux erreur — Règle métier violée

1. Le service domaine rejette l'action (ex : hash déjà existant, interface déjà utilisée, licence incompatible)
2. Le message d'erreur retourné par le service est affiché tel quel — identique à celui que recevrait un client REST

## Post-conditions

- La commande CLI a été exécutée avec succès, avec le même effet et les mêmes règles métier que l'action équivalente via l'API REST

## Diagramme d'activités

```plantuml
@startuml
skin rose
title Utilisation du CLI
start
:Ouvrir un terminal sur le serveur distant;
:Saisir une commande MYR CLI;
:Transmettre la commande au réseau;
if (Commande valide?) then (oui)
  :Afficher le résultat en sortie texte;
  stop
else (non)
  :Afficher un message d'aide avec la syntaxe correcte;
  stop
endif
@enduml
```

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCDEV — Developpement](../../Carte_des_specs.md#UCDEV%20—%20Developpement)
- [UCDEV02 — couche analyse](../../2-Analyse/UCDEV-Developpement/UCDEV02.md)
- [Traçabilité UCDEV02 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCDEV02)

**Exigences fonctionnelles couvertes**
- [EF56 — Proposer un CLI d'administration serveur](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Use cases cités**
- [UCADM01 — Ajouter une organisation au réseau](../UCADM-Administration/UCADM01.md)
- [UCADM02 — Créer un réseau indépendant](../UCADM-Administration/UCADM02.md)
- [UCADM03 — Ajouter un nœud à un réseau existant](../UCADM-Administration/UCADM03.md)
- [UCADM04 — Retirer un nœud d'un réseau existant](../UCADM-Administration/UCADM04.md)
- [UCADM05 — Démanteler un réseau (dev/test uniquement)](../UCADM-Administration/UCADM05.md)
- [UCDEV01 — Utilisation de l'API](UCDEV01.md)

**Cité par**
- [Expression_des_besoins_Intro](../Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../Matrice_Tracabilite.md)
- [todo (expression)](../todo.md)
- [Analyse_des_besoins](../../2-Analyse/Analyse_des_besoins.md)
- [todo (analyse)](../../2-Analyse/todo.md)
- [DC_CLI_Admin](../../3-Conception/DC_CLI_Admin.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
