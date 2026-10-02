---
categorie: Propriété Intellectuelle
titre: "Transfert de propriété intellectuelle"
probabilite: 1
impact: 3
importance: 3
etat: relire
tags:
  - couche/expression
  - type/use-case
  - famille/UCPI
  - domaine/model
  - domaine/payment
  - uc/UCPI07
---

# Transfert de propriété intellectuelle

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Propriétaire" as P
actor "Destinataire" as DEST

rectangle "Application MYR" {
    usecase "Transférer la propriété d'un asset" as UC1
    usecase "Accepter le transfert" as UC2
    usecase "Enregistrer sur la blockchain" as UC3
}

P --> UC1
DEST --> UC2
UC1 ..> UC2 : <<include>>
UC2 ..> UC3 : <<include>>

@enduml
```

## Contexte

Un propriétaire peut transférer la propriété intellectuelle d'un composant ou module à un autre utilisateur ou organisation.

Conformément au principe de parité CLI/REST, le transfert de propriété d'un asset doit pouvoir être initié en CLI pour le compte d'un propriétaire, au même titre que via l'interface graphique.

## Pré-conditions

- Être connecté au réseau
- Être propriétaire du composant/module à transférer

## Scénario

**Étape initiale :** `myr model transfer <id> <destinataireID>` est exécutée (ou l'appel API équivalent), pour le compte du propriétaire

### Flux nominal — Transfert réussi

1. L'identifiant du destinataire est transmis
2. La transaction de transfert est soumise sur la blockchain
3. Le destinataire est notifié et accepte le transfert

### Flux alternatif — Transfert vers un réseau externe

1. L'identifiant du destinataire appartient à un utilisateur sur un réseau externe
2. Le système génère un token de transfert signé cryptographiquement
3. Le destinataire est notifié sur son réseau et accepte le transfert via le token
4. La propriété est transférée avec enregistrement de la transaction sur les deux réseaux

## Post-conditions

- La propriété est transférée et enregistrée de manière immuable sur la blockchain
- L'ancien propriétaire perd les droits d'édition

## Diagramme d'activités

```plantuml
@startuml
skin rose
title Transfert de propriété intellectuelle
start
:Transmettre l'identifiant du destinataire (myr model transfer);
if (Destinataire sur un réseau externe?) then (oui)
  :Générer un token de transfert signé cryptographiquement;
  :Notifier le destinataire sur son réseau avec le token;
  :Le destinataire accepte via le token;
  :Enregistrer la transaction sur les deux réseaux;
  stop
else (non)
  :Soumettre la transaction de transfert sur la blockchain;
  :Notifier le destinataire;
  :Le destinataire accepte le transfert;
  stop
endif
@enduml
```

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCPI — Propriete Intellectuelle](../../Carte_des_specs.md#UCPI%20—%20Propriete%20Intellectuelle)
- [UCPI07 — couche analyse](../../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI07.md)
- [Traçabilité UCPI07 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCPI07)

**Exigences fonctionnelles couvertes**
- [EF35 — Transférer la propriété intellectuelle d'un asset](../Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Cité par**
- [Expression_des_besoins_Intro](../Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../Matrice_Tracabilite.md)
- [todo (expression)](../todo.md)
- [Analyse_des_besoins](../../2-Analyse/Analyse_des_besoins.md)
- [todo (analyse)](../../2-Analyse/todo.md)
- [API_REST](../../3-Conception/API_REST.md)
- [Chaincode](../../3-Conception/Chaincode.md)
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)
- [DC_D7_Payment](../../3-Conception/DC_D7_Payment.md)
- [todo (conception)](../../3-Conception/todo.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
