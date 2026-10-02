---
categorie: Composant Ecriture
titre: "Supprimer un Composant"
probabilite: 3
impact: 4
importance: 12
etat: analyse
tags:
  - couche/analyse
  - type/use-case
  - famille/UCCE
  - domaine/model
  - uc/UCCE07
  - rm/RM06
  - rm/RM08
  - rm/RM14
  - enf/ENF12
---

# Supprimer un Composant

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Concepteur" as C

rectangle "Application MYR" {
    usecase "Supprimer un composant" as UC1
    usecase "Retirer en cascade les connexions" as UC2
}

C --> UC1
UC1 ..> UC2 : <<include>>

@enduml
```

## Contexte

Un composant qui n'a plus d'usage pour son propriétaire (brouillon abandonné, ou composant déjà soumis mais devenu obsolète) doit pouvoir être supprimé. Aucune opération de suppression n'existe sur la blockchain Fabric (RM06, RM08) : un composant déjà soumis y reste inscrit en permanence, de façon immuable.

La suppression a donc un effet différent selon l'état du composant :
- **Brouillon (`draft`)** : le composant est réellement retiré du stockage local (`DraftStore`) — il n'a jamais existé sur la blockchain.
- **Soumis (`submitted`)** : le composant est masqué localement (il disparaît des listes retournées par `GET /api/components`) mais son enregistrement sur la blockchain n'est jamais modifié ni retiré. Il reste consultable directement par son identifiant (`GET /api/components/:id`) — par exemple parce qu'un autre asset le référence comme parent, ou qu'un client en a gardé une référence directe.

Dans les deux cas, les connexions impliquant ce composant sont retirées en cascade (même mécanisme que RM14/UCAM08).

## Pré-conditions

- Le Concepteur est authentifié avec le rôle `contributor`
- Le Concepteur est propriétaire du composant (`OwnerID` correspond à son identité)
- Le composant existe (brouillon local ou soumis)

## Scénario

**Étape initiale :** `DELETE /api/components/:id` est appelée (ou l'équivalent CLI `myr model remove`)

### Flux nominal — Suppression d'un brouillon

1. Le composant est retrouvé dans le stockage local (`DraftStore`)
2. Les connexions impliquant ce composant sont retirées en cascade
3. Le composant est retiré du stockage local
4. L'API retourne `204 No Content`
5. Le composant n'apparaît plus dans aucune liste ni lecture directe — il n'a jamais existé sur la blockchain

### Flux nominal — Suppression (masquage) d'un composant soumis

1. Le composant est retrouvé sur la blockchain (déjà soumis)
2. Les connexions impliquant ce composant sont retirées en cascade
3. L'identifiant du composant est enregistré dans le registre local de masquage — aucune écriture blockchain n'est effectuée
4. L'API retourne `204 No Content`
5. Le composant disparaît de `GET /api/components` mais reste accessible via `GET /api/components/:id` et reste inscrit sur la blockchain, inchangé

### Flux erreur — Composant introuvable

1. Aucun composant ne correspond à l'identifiant fourni
2. L'API retourne `404 Not Found`

### Flux erreur — Droits insuffisants

1. L'identité n'est pas propriétaire du composant
2. L'API retourne `403 Forbidden`

## Post-conditions

- Le composant n'apparaît plus dans `GET /api/components`
- Si le composant était un brouillon : il n'existe plus nulle part
- Si le composant était soumis : son enregistrement blockchain reste inchangé et consultable par identifiant direct
- Les connexions impliquant le composant sont retirées

## Diagramme de séquence

```plantuml
@startuml
participant "Client\n(CLI ou API REST)" as Client
participant "REST Handler\n(adapters/in/rest/)" as REST
participant "Model Service\n(domain/model/)" as Service
database "LocalStorage\n(adapters/out/localstorage/)" as Local
database "Fabric\n(adapters/out/fabric/)" as Fabric

Client -> REST : DELETE /api/components/:id
REST -> Service : Remove(id)
Service -> Local : GetDraft(id)

alt Composant en brouillon
    Local --> Service : *Model3D (draft)
    Service -> Service : removeConnectionsFor(id)
    Service -> Local : RemoveDraft(id)
    Local --> Service : OK
else Composant soumis
    Local --> Service : introuvable
    Service -> Service : removeConnectionsFor(id)
    Service -> Local : HideAsset(id)
    Local --> Service : OK
    note right : Fabric n'est jamais appelé —\naucune suppression sur le ledger (RM08)
end

Service --> REST : nil (succès)
REST --> Client : 204 No Content
@enduml
```

## Règles métier déclenchées

| Règle | Description | Point d'application |
|-------|-------------|---------------------|
| **RM06** | Immuabilité des transactions — aucune suppression sur le ledger | `Remove()` n'appelle jamais d'opération de suppression blockchain |
| **RM08** | Masquage local, ledger jamais modifié | `Remove()` — brouillon retiré, asset soumis masqué via le registre local |
| **RM14** (généralisée) | Retrait en cascade des connexions impliquant l'asset supprimé | `removeConnectionsFor()` |

## Exigences non-fonctionnelles

- **ENF12** : Contrôle d'ownership côté serveur obligatoire

## Notes d'implémentation

**Endpoint REST utilisé :**
- `DELETE /api/components/:id` → `service.Remove()` (handlers.go — `deleteComponent`)

**Registre de masquage :** le port `RemovedAssetStore` (`domain/model/ports.go`) est implémenté par le store local consolidé (`adapters/out/localstorage/json_blockchain.go`, au même titre que `DraftStore`/`ConnectionStore`). `List()`/`ListModules()` excluent les identifiants masqués ; `Get()`/`GetModule()` restent volontairement non filtrés, pour rester résolvables par identifiant direct (ex. vérification de licence d'un composant dérivé dont le parent a été supprimé).

**Commande CLI équivalente :** `myr model remove <id>` appelle la même méthode `Remove()`.

**Voir aussi UCMOD08** (« Supprimer un Module ») — comportement strictement identique, RM08 s'applique de façon générique à tout asset, composant ou module.

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCCE — Composant Ecriture](../../Carte_des_specs.md#UCCE%20—%20Composant%20Ecriture)
- [UCCE07 — couche expression](../../1-Expression/UCCE-Composant_Ecriture/UCCE07.md)
- [Traçabilité UCCE07 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCCE07)

**Exigences fonctionnelles couvertes**
- [EF60 — Supprimer un composant (masquage local des listes, ledger jamais modifié)](../../1-Expression/Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Use cases cités**
- [UCAM08 — Retirer une instance de composant d'un Module](../UCAM-Assemblage_Module/UCAM08.md)
- [UCMOD08 — Supprimer un Module](../UCMOD-Module/UCMOD08.md)

**Règles métier**
- [RM06 — Immuabilité des transactions](../../1-Expression/Regles_Metier.md#2.%20Blockchain%20et%20immuabilité)
- [RM08 — Masquage local, ledger jamais modifié](../../1-Expression/Regles_Metier.md#2.%20Blockchain%20et%20immuabilité)
- [RM14 — Suppression en cascade des connexions](../../1-Expression/Regles_Metier.md#4.%20Composition%20d'un%20Module%20%28instances%29)

**Exigences non fonctionnelles**
- [ENF12 — Contrôle d'accès par rôle](../../1-Expression/Exigences_Non_Fonctionnelles.md#Tableau%20des%20exigences%20non-fonctionnelles)

**Cité par**
- [Matrice_Tracabilite](../../1-Expression/Matrice_Tracabilite.md)
- [UCMOD08 (expression)](../../1-Expression/UCMOD-Module/UCMOD08.md)
- [UCAM05 (analyse)](../UCAM-Assemblage_Module/UCAM05.md)
- [UCMOD08 (analyse)](../UCMOD-Module/UCMOD08.md)
- [API_REST](../../3-Conception/API_REST.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
