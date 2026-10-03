---
categorie: Propriété Intellectuelle
titre: "Recevoir une commission sur l'utilisation d'un Module"
probabilite: 5
impact: 5
importance: 25
etat: analyse
tags:
  - couche/analyse
  - type/use-case
  - famille/UCPI
  - domaine/model
  - domaine/payment
  - uc/UCPI02
  - rm/RM07
  - rm/RM23
  - rm/RM24
  - enf/ENF02
  - enf/ENF30
  - relecture/question
---

# Recevoir une commission sur l'utilisation d'un Module

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Consommateur" as CL
actor "Concepteur" as C
actor "Smart Contract\n(système)" as SC

rectangle "Myr System" {
    usecase "Passer une commande" as UC1
    usecase "Confirmer la livraison" as UC2
    usecase "Calculer les commissions" as UC3
    usecase "Distribuer les commissions" as UC4
    usecase "Recevoir une commission" as UC5
}

CL --> UC1
SC --> UC3
SC --> UC4
C --> UC5
UC1 ..> UC2 : <<extend>> (UCAUT01)
UC2 ..> UC3 : <<include>>
UC3 ..> UC4 : <<include>>
UC4 ..> UC5 : <<include>>

@enduml
```

## Contexte

Lorsqu'un module est commandé et livré, le smart contract Fabric calcule automatiquement les commissions dues à chaque auteur impliqué dans la chaîne de propriété du module (composants constitutifs, dérivations #question — voir `specs/1-Expression/Regles_Metier.md` RM23 : le périmètre exact "composition" vs "dérivation" n'est pas tranché). La distribution est atomique avec la confirmation de livraison — RM23 impose que les commissions soient distribuées **par** le smart contract, sans intervention humaine.

Ce use case est déclenché par UCAUT01 (livraison confirmée) et non directement par le consommateur. Le concepteur est un acteur passif : il reçoit, il ne demande pas.

## Pré-conditions

- Une commande est passée (UCPI01) et enregistrée sur la blockchain avec statut `pending`
- Le module est soumis (`Status = submitted`) avec un prix et des commissions définies sur ses composants
- Le manufactureur ou la boutique a confirmé la livraison (déclencheur externe — UCAUT01)
- Chaque auteur concerné dispose d'un wallet actif sur le réseau

## Scénario

**Étape initiale :** Le manufactureur ou la boutique confirme la livraison sur la blockchain (voir UCAUT01)

### Flux nominal — Commission distribuée à un auteur unique

1. Le smart contract reçoit l'événement de livraison confirmée (`OrderDelivered`)
2. Il lit les métadonnées du module sur la blockchain (OwnerID, prix, liste des composants)
3. Il calcule la commission due à l'auteur unique selon le taux défini sur le module
4. Une transaction de paiement individuelle est soumise : `from=escrow, to=ownerID, amount=commission, modelID`
5. La commission est créditée sur le wallet de l'auteur
6. L'auteur reçoit une notification de commission reçue (événement blockchain)

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux alternatif — Répartition entre plusieurs co-auteurs (RM24)

1. Le smart contract analyse les composants constitutifs du module
2. Pour chaque composant, il lit le `OwnerID` et le taux de commission défini
3. Il calcule la part proportionnelle de chaque auteur selon la contribution de chaque composant au prix total
4. Une transaction de paiement individuelle est créée par auteur (`for each ownerID : Transfer(escrow → ownerID, montant)`)
5. Les transactions sont soumises de façon atomique dans un seul bloc Fabric
6. Chaque co-auteur reçoit sa part et est notifié individuellement
7. L'historique complet des distributions est enregistré sur la blockchain

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux erreur A — Wallet d'un auteur inexistant ou inactif

1. Le smart contract détecte que le wallet d'un auteur n'est plus actif
2. La commission concernée est mise en séquestre (`escrow`) avec un identifiant de bénéficiaire
3. L'auteur peut réclamer sa commission ultérieurement après réactivation de son wallet
4. La livraison est confirmée malgré ce cas — les autres commissions sont distribuées normalement

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux erreur B — Échec blockchain lors de la distribution

1. Une transaction de commission échoue (endorsement refusé)
2. Toutes les transactions de la distribution sont annulées (atomicité garantie par Fabric)
3. La commande reste en statut `pending` — la livraison n'est pas confirmée
4. Une alerte est émise vers l'administrateur du réseau

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

## Post-conditions

- Chaque auteur concerné a reçu sa commission sur son wallet blockchain
- Chaque transaction de commission est enregistrée individuellement et immuablement sur la blockchain
- La commande passe au statut `delivered`
- L'historique des commissions est consultable par chaque auteur via son historique de paiements

## Diagramme de séquence

```plantuml
@startuml
participant "Client\n(CLI ou API REST)" as Browser
participant "REST Handler\n(adapters/in/rest/)" as REST
participant "Payment Service\n(domain/payment/)" as PaySvc
participant "Model Service\n(domain/model/)" as ModelSvc
database "Fabric\n(adapters/out/fabric/)" as Fabric
database "Smart Contract\n(chaincode/)" as CC

note over CC : Déclenché par UCAUT01\n(confirmation de livraison)

Fabric -> CC : Événement OrderDelivered\n{orderID, moduleID, buyerID}

CC -> CC : QueryModule(moduleID)\n→ composants, prix, OwnerIDs
CC -> CC : CalculerCommissions()\nselon taux par composant (RM24)

alt Co-auteurs multiples
    loop Pour chaque auteur
        CC -> Fabric : Transfer(escrow → ownerID_i, montant_i)\n{modelID, orderID, type=commission}
        Fabric --> CC : tx confirmée
    end
else Auteur unique
    CC -> Fabric : Transfer(escrow → ownerID, montant)\n{modelID, orderID, type=commission}
    Fabric --> CC : tx confirmée
end

CC -> Fabric : UpdateOrderStatus(orderID, delivered)
Fabric --> CC : OK

note over Fabric : Événements émis — les\nwallets des auteurs sont\nmis à jour

Browser -> REST : GET /api/payments/history
REST -> PaySvc : GetHistory(identityID)
PaySvc -> Fabric : QueryPaymentHistory(identityID)
Fabric --> PaySvc : liste des paiements reçus
PaySvc --> REST : [Payment{...}]
REST --> Browser : 200 [{commission, amount, date, modelID}]

@enduml
```

## Règles métier déclenchées

| Règle | Description | Détail |
|-------|-------------|--------|
| RM23 | Commissions distribuées automatiquement à la livraison | Le smart contract, pas le serveur, exécute la distribution |
| RM24 | Répartition proportionnelle par auteur | Calculée à partir du taux défini sur chaque composant constitutif |
| RM07 | Transactions Fabric immuables | Distribution atomique — tout ou rien |

## Exigences non-fonctionnelles

| ENF | Description |
|-----|-------------|
| ENF02 | Transactions blockchain confirmées en ≤ 5 s sous charge normale |
| ENF30 | En cas d'échec blockchain, aucune commission partielle — rollback complet |

## Notes d'implémentation

- **Non implémenté** : Le smart contract chaincode ne contient aucune logique de commission — seule l'entité `Model3D` est présente dans `chaincode/model/entity.go`
- **Non implémenté** : L'endpoint REST `GET /api/payments/history` n'existe pas — le service `GetHistory()` dans `domain/payment/service.go` est prêt mais non exposé
- **À créer** : Fonction chaincode `DistributeCommissions(orderID)` dans `chaincode/`
- **À créer** : Mécanisme d'escrow et de séquestre pour les wallets inactifs
- **À créer** : Route `GET /api/payments/history` dans `adapters/in/rest/handlers_payment.go`
- L'atomicité de la distribution multi-auteurs est garantie par le mécanisme de transaction Fabric (un seul bloc pour N transfers)
- La notion de "taux de commission" par composant n'est pas encore modélisée dans les entités (`Model3D`, `chaincode/model/entity.go`) — champ `CommissionRate float64` à ajouter
- **Parité CLI/REST :** conformément au principe de parité, la consultation des commissions reçues devrait être exposée en CLI. `GetHistory(identityID)` existe déjà côté service mais n'est exposé ni en REST ni en CLI ; surtout, aucune notion de `CommissionRate` ni de distribution proportionnelle multi-auteurs n'existe dans le domaine (voir notes ci-dessus). Une commande CLI (par ex. `myr payment commissions <identityID>`) ne pourra être ajoutée qu'une fois ces éléments conçus au niveau domaine.

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Étape suivante — conception**
- **Chaincode** : [§ 5. Fonctions chaincode — Commissions (D7, à implémenter)](../../3-Conception/Chaincode.md#5.%20Fonctions%20chaincode%20—%20Commissions%20%28D7,%20à%20implémenter%29)
- **DC_CLI_Model** : [§ 6. Écarts et points ouverts](../../3-Conception/DC_CLI_Model.md#6.%20Écarts%20et%20points%20ouverts)
- **DC_D7_Payment** : [§ DC — D7 : Propriété Intellectuelle & Paiements](../../3-Conception/DC_D7_Payment.md#DC%20—%20D7%20:%20Propriété%20Intellectuelle%20&%20Paiements) · [§ 3. Entités à concevoir (RM23/RM24)](../../3-Conception/DC_D7_Payment.md#3.%20Entités%20à%20concevoir%20%28RM23/RM24%29) · [§ 4. Algorithme de distribution des commissions (RM23/RM24)](../../3-Conception/DC_D7_Payment.md#4.%20Algorithme%20de%20distribution%20des%20commissions%20%28RM23/RM24%29) · [§ 7. Écarts code → specs](../../3-Conception/DC_D7_Payment.md#7.%20Écarts%20code%20→%20specs)
- **Modele_Domaine** : [§ Agrégat Paiement (D7)](../../3-Conception/Modele_Domaine.md#Agrégat%20Paiement%20%28D7%29)

<!-- liens-obsidian:end -->
