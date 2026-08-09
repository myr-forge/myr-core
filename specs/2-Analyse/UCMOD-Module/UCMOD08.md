---
categorie: Module
titre: "Supprimer un Module"
probabilite: 3
impact: 4
importance: 12
etat: analyse
---

# Supprimer un Module

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Concepteur" as C

rectangle "Application MYR" {
    usecase "Supprimer un module" as UC1
    usecase "Retirer en cascade les connexions" as UC2
}

C --> UC1
UC1 ..> UC2 : <<include>>

@enduml
```

## Contexte

Un module (brouillon abandonné, ou module déjà soumis mais devenu obsolète) doit pouvoir être supprimé par son propriétaire. Ce comportement est strictement identique à la suppression d'un composant (UCCE07) : `RemoveModule()` délègue au même `Remove()` que les composants, et RM08 s'applique de façon générique à tout asset — composant ou module.

Un module en état `draft` est réellement retiré (`DraftStore`). Un module déjà `submitted` est masqué localement (disparaît de `GET /api/modules`) mais son enregistrement blockchain — y compris la ou les `ModuleVersion` déjà ancrées — n'est jamais modifié. Le Concepteur choisit ainsi librement entre :
- **Supprimer** un module soumis dont il n'a plus l'usage (il disparaît de ses listes, sans perte d'intégrité du ledger).
- **Conserver** un module soumis pour le dériver plus tard sans reconstruire sa composition (voir UCMOD01, flux « Dérivation d'un Module existant »).

## Pré-conditions

- Le Concepteur est authentifié avec le rôle `contributor`
- Le Concepteur est propriétaire du module (`OwnerID` correspond à son identité)
- Le module existe (brouillon local ou soumis)

## Scénario

**Étape initiale :** `DELETE /api/modules/:id` est appelée (ou l'équivalent CLI `myr module remove`)

### Flux nominal — Suppression d'un module en brouillon

1. Le module est retrouvé dans le stockage local (`DraftStore`)
2. Les connexions internes du module (`Assemblies`) et toute connexion externe l'impliquant sont retirées en cascade
3. Le module est retiré du stockage local
4. L'API retourne `204 No Content`

### Flux nominal — Suppression (masquage) d'un module soumis

1. Le module est retrouvé sur la blockchain (déjà soumis, avec ses `ModuleVersion`)
2. Les connexions impliquant le module (en tant qu'instance dans un autre module hôte) sont retirées en cascade
3. L'identifiant du module est enregistré dans le registre local de masquage — aucune écriture blockchain, aucune modification de `ModuleVersion`
4. L'API retourne `204 No Content`
5. Le module disparaît de `GET /api/modules` mais reste accessible via `GET /api/modules/:id` — notamment pour toute dérivation ultérieure (`parent_id`, UCMOD01) ou pour un module hôte qui l'a déjà instancié

### Flux erreur — Module introuvable

1. Aucun module ne correspond à l'identifiant fourni
2. L'API retourne `404 Not Found`

### Flux erreur — Droits insuffisants

1. L'identité n'est pas propriétaire du module
2. L'API retourne `403 Forbidden`

## Post-conditions

- Le module n'apparaît plus dans `GET /api/modules`
- Si le module était un brouillon : il n'existe plus nulle part
- Si le module était soumis : ses `ModuleVersion` et son enregistrement blockchain restent inchangés et consultables par identifiant direct
- Les connexions impliquant le module sont retirées

## Diagramme de séquence

```plantuml
@startuml
participant "Client\n(CLI ou API REST)" as Client
participant "REST Handler\n(adapters/in/rest/)" as REST
participant "Model Service\n(domain/model/)" as Service
database "LocalStorage\n(adapters/out/localstorage/)" as Local
database "Fabric\n(adapters/out/fabric/)" as Fabric

Client -> REST : DELETE /api/modules/:id
REST -> Service : RemoveModule(id)
Service -> Service : Remove(id) — voir UCCE07, même mécanisme
Service -> Local : GetDraft(id)

alt Module en brouillon
    Local --> Service : *Model3D (draft)
    Service -> Service : removeConnectionsFor(id)
    Service -> Local : RemoveDraft(id)
else Module soumis
    Local --> Service : introuvable
    Service -> Service : removeConnectionsFor(id)
    Service -> Local : HideAsset(id)
    note right : Fabric n'est jamais appelé —\nModuleVersion déjà ancrées inchangées (RM06/RM08)
end

Service --> REST : nil (succès)
REST --> Client : 204 No Content
@enduml
```

## Règles métier déclenchées

| Règle | Description | Point d'application |
|-------|-------------|---------------------|
| **RM06** | Immuabilité des transactions — aucune suppression sur le ledger, `ModuleVersion` incluses | `Remove()` n'appelle jamais d'opération de suppression blockchain |
| **RM08** | Masquage local, ledger jamais modifié | `RemoveModule()` → `Remove()` — identique à UCCE07 |
| **RM14/RM15** (généralisées) | Retrait en cascade des connexions impliquant l'asset supprimé | `removeConnectionsFor()` |

## Exigences non-fonctionnelles

- **ENF12** : Contrôle d'ownership côté serveur obligatoire
- **ENF28** : Immuabilité des `ModuleVersion` déjà ancrées, y compris après suppression locale du module

## Notes d'implémentation

**Endpoint REST utilisé :**
- `DELETE /api/modules/:id` → `service.RemoveModule()` (handlers.go — `handleModule`, méthode DELETE)

**Registre de masquage :** identique à UCCE07 — même port `RemovedAssetStore`, même store local (`adapters/out/localstorage/json_blockchain.go`), même exclusion de `ListModules()` / non-filtrage de `GetModule()`.

**Commande CLI équivalente :** `myr module remove <id>` appelle `RemoveModule()`.

**Voir aussi UCMOD01** (flux « Dérivation d'un Module existant ») : un module masqué reste un point de départ valide pour une dérivation, puisque `GetModule(id)` continue de le résoudre.
