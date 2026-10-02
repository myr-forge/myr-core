---
categorie: Administration
titre: "Démanteler un réseau (dev/test uniquement)"
tags:
  - couche/expression
  - type/use-case
  - famille/UCADM
  - domaine/channel
  - domaine/identity
  - domaine/network
  - uc/UCADM05
  - rm/RM28
---

# Démanteler un réseau (dev/test uniquement)

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Administrateur" as ADM

rectangle "Infrastructure locale (hors blockchain)" {
    usecase "Confirmer le démantèlement" as UC1
    usecase "Arrêter les processus Fabric" as UC2
    usecase "Supprimer les artefacts cryptographiques" as UC3
    usecase "Supprimer les données ledger locales" as UC4
    usecase "Supprimer le profil de connexion" as UC5
}

ADM --> UC1
UC1 ..> UC2 : <<include>>
UC1 ..> UC3 : <<include>>
UC1 ..> UC4 : <<include>>
UC1 ..> UC5 : <<include>>

@enduml
```

## Contexte

Lors du développement et des tests, un administrateur peut avoir besoin de démanteler entièrement un réseau Fabric pour en recréer un propre — par exemple après une série de tests de validation du processus de création réseau (UCADM02) ou pour repartir d'un état vierge.

Cette opération est **hors périmètre blockchain** : elle arrête des processus système et supprime des fichiers locaux. Elle ne soumet aucune transaction Fabric et ne constitue pas une violation de l'immuabilité du ledger (RM28). Les données blockchain que le réseau avait produites disparaissent simplement avec la destruction de l'infrastructure locale.

**Cette opération est irréversible. Elle n'est disponible que via le CLI admin — jamais via l'API REST.**

## Pré-conditions

- Être administrateur avec accès à l'infrastructure locale
- Le réseau cible n'est pas marqué comme réseau de production (`IsProduction: false`)
- Commande exécutée via le CLI admin (`myr`)
- Flag `--confirm` explicitement fourni

## Scénario

**Étape initiale :** L'administrateur exécute `myr network destroy <id> --confirm`.

### Flux nominal — Réseau démantelé avec succès

1. L'administrateur fournit l'ID du réseau et le flag `--confirm`
2. Le système vérifie que le réseau n'est pas marqué en production
3. Les processus Fabric sont arrêtés (peer, orderer, CA) sur la machine locale
4. Les données ledger locales sont supprimées (`/var/hyperledger/production/` ou équivalent)
5. Les artefacts cryptographiques sont supprimés (certificats, clés, channel artifacts, wallet)
6. Le profil de connexion est supprimé du stockage local (`NetworkStore.Delete()`)
7. Confirmation : `Réseau "<nom>" démantelé. Toutes les données locales ont été supprimées.`

### Flux erreur — Confirmation absente

1. L'administrateur exécute la commande sans le flag `--confirm`
2. Le système affiche un avertissement et refuse d'agir
3. Message :
   ```
   ATTENTION : Cette opération est irréversible.
   Elle supprimera toutes les données locales du réseau "<nom>".
   Utilisez --confirm pour confirmer : myr network destroy <id> --confirm
   ```

### Flux erreur — Réseau marqué en production

1. Le profil de connexion porte le flag `IsProduction: true`
2. Le système refuse l'opération, même avec `--confirm`
3. Message : `Erreur : "<nom>" est marqué comme réseau de production. Démantèlement refusé.`

### Flux erreur — Processus Fabric non arrêtables

1. Un ou plusieurs processus Fabric ne répondent pas à l'arrêt
2. Le système signale les processus concernés
3. Message : `Avertissement : processus <pid> (peer) non arrêté. Suppression des données poursuivie.`
4. L'opération continue — les fichiers sont supprimés, le profil retiré

## Post-conditions

- Les processus Fabric sont arrêtés sur la machine locale
- Toutes les données locales (ledger, certificats, artefacts cryptographiques) sont supprimées
- Le profil de connexion est retiré de la liste des réseaux dans myr-api
- Opération irréversible — aucune donnée récupérable

## Diagramme d'activités

```plantuml
@startuml
skin rose
title Démanteler un réseau (dev/test uniquement)
start
if (Flag --confirm présent?) then (non)
  :Afficher avertissement et instructions;
  stop
else (oui)
  if (Réseau marqué IsProduction = true?) then (oui)
    :Erreur : démantèlement refusé sur réseau de production;
    stop
  else (non)
    :Arrêter les processus Fabric (peer, orderer, CA);
    :Supprimer les données ledger locales;
    :Supprimer les artefacts cryptographiques;
    :Supprimer le profil de connexion (NetworkStore.Delete);
    :Confirmer "Réseau démantelé — données supprimées";
    stop
  endif
endif
@enduml
```

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCADM — Administration](../../Carte_des_specs.md#UCADM%20—%20Administration)
- [UCADM05 — couche analyse](../../2-Analyse/UCADM-Administration/UCADM05.md)
- [Traçabilité UCADM05 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCADM05)

**Exigences fonctionnelles couvertes**
- [EF58 — Démanteler un réseau de test (CLI uniquement — jamais via REST)](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Use cases cités**
- [UCADM02 — Créer un réseau indépendant](UCADM02.md)

**Règles métier**
- [RM28 — Démantèlement réseau : opération d'infrastructure locale](../Regles_Metier.md#8.%20Administration%20réseau)

**Cité par**
- [Expression_des_besoins_Intro](../Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../Matrice_Tracabilite.md)
- [UCDEV02 (expression)](../UCDEV-Developpement/UCDEV02.md)
- [Analyse_des_besoins](../../2-Analyse/Analyse_des_besoins.md)
- [UCDEV02 (analyse)](../../2-Analyse/UCDEV-Developpement/UCDEV02.md)
- [DC_CLI_Admin](../../3-Conception/DC_CLI_Admin.md)
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)
- [DC_D2_Administration](../../3-Conception/DC_D2_Administration.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
