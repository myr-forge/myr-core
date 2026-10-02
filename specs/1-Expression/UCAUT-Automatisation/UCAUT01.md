---
categorie: Automatisation
titre: "Fabrication/Livraison d'un Composant"
probabilite: 5
impact: 4
importance: 20
etat: relire
tags:
  - couche/expression
  - type/use-case
  - famille/UCAUT
  - domaine/model
  - domaine/payment
  - domaine/role
  - uc/UCAUT01
---

# Fabrication/Livraison d'un Composant

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Consommateur" as CL
actor "Manufactureur" as M
actor "Concepteur" as C

rectangle "Application MYR" {
    usecase "Commander une fabrication" as UC1
    usecase "Transmettre CAO au manufactureur" as UC2
    usecase "Livrer le composant" as UC3
    usecase "Distribuer les commissions" as UC4
}

CL --> UC1
UC1 ..> UC2 : <<include>>
M --> UC3
UC3 ..> UC4 : <<include>>
C --> UC4

@enduml
```

## Contexte

Passage par un fabricant externe agréé par le réseau qui se chargera de créer et livrer le(s) composant(s) à l'adresse définie.

Conformément au principe de parité CLI/REST, la confirmation de livraison par un manufactureur (et la distribution automatique des commissions qui s'ensuit) doit être déclenchable en CLI pour son compte, au même titre que via l'interface graphique ou l'API REST.

## Pré-conditions

- Être connecté au réseau
- Composant commandé avec fichier CAO disponible sur la blockchain
- Manufactureur agréé disponible sur le réseau
- Adresse de livraison renseignée

## Scénario

**Étape initiale :** Une commande de fabrication est déclenchée

### Flux nominal — Fabrication et livraison réussies

1. Le système identifie le manufactureur agréé disponible
2. Le fichier CAO et les spécifications sont transmis via la blockchain
3. Le manufactureur produit le composant
4. Le composant est livré à l'adresse définie
5. Le statut de commande est mis à jour sur le réseau

### Flux alternatif — Fabricant partenaire industriel externe (sans nœud réseau)

1. Le composant est confié à un fabricant partenaire déjà établi (ex. imprimeur 3D ou façonnier industriel), intégré à Myr sans qu'il opère lui-même de nœud ni de compte sur le réseau blockchain
2. Le fichier CAO, les spécifications et l'adresse de livraison lui sont transmis via l'intégration propre à ce partenaire (son API commerciale existante)
3. Le partenaire produit et livre le composant selon son propre processus
4. Sa confirmation de livraison est reçue par le réseau Myr via cette même intégration, qui enregistre le statut de commande et déclenche la distribution des commissions exactement comme dans le flux nominal

### Flux alternatif — Commande en lot (quantité > 1)

1. Le consommateur saisit une quantité supérieure à 1 dans sa commande
2. Le système vérifie si la capacité de fabrication disponible peut absorber le lot
3. Si la capacité est suffisante chez un seul manufactureur : la commande est acceptée comme lot unique
4. Si la capacité est partielle : le système propose de répartir la commande entre plusieurs manufactureurs
5. L'utilisateur accepte la répartition
6. Les ordres de fabrication sont transmis en parallèle sur la blockchain à chaque manufactureur concerné

### Flux erreur — Aucun manufactureur disponible

1. L'utilisateur est notifié et mis en liste d'attente

## Post-conditions

- Le composant est fabriqué et livré
- La commande est enregistrée sur la blockchain
- Les commissions des auteurs sont distribuées automatiquement

## Diagrammes

### Flux de commande — boutique partenaire ou API MYR vers manufactureur

```plantuml
@startuml
skin rose
title fonctionnement de l'échange
:Consommateur: --> (website) :order
:Consommateur: --> (myr) :order
(website) --> (myr) : get CAO
(myr) --> (blockchain) : get CAO
(blockchain) ..> (manufacturer) : build product
(blockchain) ..> (shop) : buy product
(shop) ..> :Consommateur: :deliver
(manufacturer) ..> :Consommateur: :deliver
@enduml
```

### Diagramme d'activités

```plantuml
@startuml
skin rose
title Fabrication/Livraison d'un Composant
start
:Déclencher une commande de fabrication;
if (Manufactureur agréé disponible?) then (oui)
  if (Quantité > 1 et capacité partielle?) then (oui)
    :Proposer une répartition entre plusieurs manufactureurs;
    :L'utilisateur accepte la répartition;
    :Transmettre les ordres de fabrication en parallèle sur la blockchain;
  else (non)
    :Identifier le manufactureur disponible;
    :Transmettre le fichier CAO et les spécifications via la blockchain;
  endif
  :Produire le(s) composant(s);
  :Livrer à l'adresse définie;
  :Mettre à jour le statut de commande sur le réseau;
  :Distribuer automatiquement les commissions des auteurs;
  stop
else (non)
  :Notifier l'utilisateur;
  :Mettre l'utilisateur en liste d'attente;
  stop
endif
@enduml
```

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCAUT — Automatisation](../../Carte_des_specs.md#UCAUT%20—%20Automatisation)
- [UCAUT01 — couche analyse](../../2-Analyse/UCAUT-Automatisation/UCAUT01.md)
- [Traçabilité UCAUT01 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCAUT01)

**Exigences fonctionnelles couvertes**
- [EF31 — Distribuer automatiquement les commissions aux auteurs à la livraison](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)
- [EF44 — Automatiser la fabrication et la livraison d'un composant](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Cité par**
- [Expression_des_besoins_Intro](../Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../Matrice_Tracabilite.md)
- [todo (expression)](../todo.md)
- [Analyse_des_besoins](../../2-Analyse/Analyse_des_besoins.md)
- [UCAUT02 (analyse)](../../2-Analyse/UCAUT-Automatisation/UCAUT02.md)
- [UCPI01 (analyse)](../../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI01.md)
- [UCPI02 (analyse)](../../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI02.md)
- [UCPI04 (analyse)](../../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI04.md)
- [UCPI05 (analyse)](../../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI05.md)
- [todo (analyse)](../../2-Analyse/todo.md)
- [API_REST](../../3-Conception/API_REST.md)
- [Chaincode](../../3-Conception/Chaincode.md)
- [Conception_intro](../../3-Conception/Conception_intro.md)
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)
- [DC_D7_Payment](../../3-Conception/DC_D7_Payment.md)
- [DC_D9_Automatisation](../../3-Conception/DC_D9_Automatisation.md)
- [Modele_Domaine](../../3-Conception/Modele_Domaine.md)
- [todo (conception)](../../3-Conception/todo.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
