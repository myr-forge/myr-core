---
categorie: Assemblage Module
titre: "Visualiser les interfaces physiques de composants"
probabilite: 3
impact: 5
importance: 15
etat: analyse
tags:
  - couche/analyse
  - type/use-case
  - famille/UCAM
  - domaine/model
  - uc/UCAM02
  - rm/RM11
  - rm/RM13
  - enf/ENF12
---

# Visualiser les interfaces physiques de composants

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Concepteur" as C
actor "Consommateur" as CL

rectangle "API myr" {
    usecase "Obtenir les interfaces physiques" as UC1
}

C --> UC1
CL --> UC1

@enduml
```

## Contexte

Chaque composant ou module expose ses interfaces physiques (`AssetInterface`) comme points de connexion. Cette lecture est le prérequis de toute opération de liaison (UCAM01) ou de création d'interface (UCAM03) côté client.

Les interfaces sont persistées dans l'`InterfaceStore` local tant que l'asset qui les porte est en brouillon (ADR-02, `specs/3-Conception/Conception_intro.md`) — récupérées via `GET /api/components/:id/interfaces` ou `GET /api/modules/:id/interfaces`. Elles ne rejoignent la blockchain (`Model3D.Interfaces`) qu'à la soumission de l'asset.

Pour un module, les interfaces exposées sont calculées dynamiquement : seules les interfaces des sous-composants **non reliées en interne** sont exposées (`GetModuleInterfaces` — calcul récursif avec cache). Chaque instance du module contribue ses propres interfaces non reliées en interne, y compris lorsque plusieurs instances partagent le même asset sous-jacent : le cache de résolution récursive par asset ne doit jamais faire disparaître une instance de la liste exposée.

Un slot virtuel (`Virtual: true`) est toujours présent sur chaque asset (RM13), permettant au client de proposer la création d'une nouvelle interface (voir UCAM03).

> La représentation visuelle (icônes, grisage, tooltip) relève du dépôt GUI externe — hors périmètre de ce document. Le contrat REST ci-dessous ainsi que son équivalent CLI (`myr model interface list`, voir Notes d'implémentation) font partie de `myr`.

## Pré-conditions

- L'utilisateur est authentifié (rôle **Lecteur** minimum)
- Un composant ou module est identifié (ID connu du client)
- L'`InterfaceStore` est configuré côté serveur

## Scénario

### Flux nominal — Interfaces d'un composant simple

1. Le client appelle `GET /api/components/:id/interfaces`
2. Le handler appelle `service.ListInterfacesForAsset(assetID)`, qui lit le brouillon local (`InterfaceStore`) si l'asset n'est pas encore soumis, sinon `Model3D.Interfaces` via `blockchain.GetModelRecord()`
3. Si aucune interface n'existe encore, `EnsureVirtualSlot(assetID)` crée un slot virtuel (RM13) et le persiste localement (brouillon)
4. La liste des interfaces est retournée : chaque interface contient `{ id, asset_id, name, category, tag, type, direction, value_min, value_max, is_range, unit, virtual }`

### Flux nominal — Interfaces d'un module

1. Le client appelle `GET /api/modules/:id/interfaces`
2. Le handler appelle `service.GetModuleInterfaces(id)` — calcul récursif
3. Le service parcourt les `WorkspaceInstances` du module et collecte les interfaces de chaque sous-composant — une instance dont l'asset est partagé avec une autre instance du même module contribue quand même ses propres interfaces, séparément
4. Seules les interfaces non présentes dans une connexion interne au module sont retournées (interfaces "exposées")
5. Un slot virtuel propre au module est garanti si aucune interface directe n'existe

### Flux erreur — InterfaceStore non configuré

1. Le service `ListInterfacesForAsset()` retourne `nil, nil` (pas d'erreur — juste une liste vide)

### Flux erreur — Asset introuvable

1. `GET /api/components/:id/interfaces` pour un ID inexistant
2. Le handler retourne HTTP 404 ou une liste vide selon l'état du store

## Post-conditions

- Les interfaces physiques du composant ou module sont retournées au client
- Le slot virtuel (`Virtual: true`) est toujours présent (RM13)
- L'état du système est inchangé (lecture seule)

## Diagramme de séquence

```plantuml
@startuml
participant "Client\n(CLI ou API REST)" as Client
participant "REST Handler\n(adapters/in/rest/)" as REST
participant "Model Service\n(domain/model/)" as Service
database "LocalStorage\n(adapters/out/localstorage/)" as Local

alt Composant simple
    Client -> REST : GET /api/components/:id/interfaces
    REST -> Service : ListInterfacesForAsset(assetID)
    Service -> Local : ifaceStore.ListInterfacesForAsset(assetID)
    Local --> Service : []*AssetInterface

    alt Aucune interface (première ouverture)
        Service -> Service : EnsureVirtualSlot(assetID) [RM13]
        Service -> Local : ifaceStore.SaveInterface({Virtual: true})
        Service -> Local : ifaceStore.ListInterfacesForAsset(assetID)
        Local --> Service : [*AssetInterface{Virtual:true}]
    end

    Service --> REST : []*AssetInterface
    REST --> Client : 200 [{ id, category, tag, type,\n  direction, value_min, value_max,\n  is_range, unit, virtual }]

else Module (interfaces exposées)
    Client -> REST : GET /api/modules/:id/interfaces
    REST -> Service : GetModuleInterfaces(moduleID)

    loop Pour chaque WorkspaceInstance du module
        Service -> Local : getModuleInterfacesInto(inst.AssetID, cache)
        Local --> Service : []*AssetInterface
    end

    Service -> Service : filtrer interfaces internes\n(présentes dans connexions du module)
    Service --> REST : []*AssetInterface exposées
    REST --> Client : 200 [interfaces exposées]
end
@enduml
```

## Règles métier déclenchées

| Règle | Description | État code |
|-------|-------------|-----------|
| **RM13** | Slot virtuel garanti : si aucune interface définie, un slot `Virtual: true` est créé automatiquement | Implémenté (`EnsureVirtualSlot`) |
| **RM11** | Les 6 attributs d'interface (`Category`, `Tag`, `Type`, `Direction`, `ValueMin/Max`, `Unit`) doivent être affichés | Partiel — `Tag` absent du modèle code (E2) |

## Exigences non-fonctionnelles

- **ENF12** — Lecture authentifiée côté serveur (rôle `reader` minimum)
- Temps de réponse `GET /api/components/:id/interfaces` : `< 200 ms` (opération locale)

## Notes d'implémentation

**Champ Tag manquant (E2) :** La struct `AssetInterface` dans `domain/model/entity.go` ne contient pas de champ `Tag`. Selon la spec §6.4 et RM11, ce champ est obligatoire. L'ajouter avant toute implémentation de l'UI d'affichage pour éviter une migration de données ultérieure.

**Routes REST concernées :**
- `GET /api/components/:id/interfaces` → `handleComponentInterfaces()` — retourne `ListInterfacesForAsset`
- `GET /api/modules/:id/interfaces` → sous-route dans `handleModule()` — retourne `GetModuleInterfaces`

**Interfaces exposées d'un module :** Le calcul est récursif dans `getModuleInterfacesInto()`. Un cache `map[string][]*AssetInterface` évite les appels blockchain redondants pour la résolution récursive d'un même asset, mais ne réduit jamais le nombre d'instances traitées au niveau du module : deux instances du module référençant le même asset exposent chacune leurs propres interfaces non connectées en interne. Les interfaces internes (présentes dans `m.Assemblies` comme `FromIfaceID` ou `ToIfaceID`) sont exclues du résultat.

**Statut d'utilisation d'une interface :** une interface est "utilisée" si son `id` apparaît dans `connection.from_iface_id` ou `connection.to_iface_id` d'une connexion non incompatible — c'est un croisement que le client peut effectuer localement à partir des réponses de `GET .../interfaces` et `GET /api/connections`, sans appel serveur supplémentaire.

**Vocabulaire de référence :** Les catégories, types et unités disponibles sont chargés via `GET /api/refs` — `handleRefs()` → `GetRefs()`. Ce vocabulaire est extensible par l'administrateur.

**Commande CLI équivalente :** `myr model interface list <assetID>` (voir `specs/3-Conception/DC_CLI_Model.md` § 3.3). Elle appelle `ModelService.ListInterfacesForAsset(assetID)` — ou `ModelService.GetModuleInterfaces(id)` si l'ID désigne un module, la commande détectant le type via `Get`/`GetModule` — soit les mêmes méthodes de service que respectivement `GET /api/components/:id/interfaces` et `GET /api/modules/:id/interfaces`. Le comportement (slot virtuel garanti RM13, calcul récursif des interfaces exposées d'un module) est strictement identique ; seul le canal de sortie change : une ligne par interface (`id`, `category`, `type`, `direction`, valeur/plage, `unit`, `virtual`) en texte terminal plutôt qu'un tableau JSON.

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCAM — Assemblage Module](../../Carte_des_specs.md#UCAM%20—%20Assemblage%20Module)
- [UCAM02 — couche expression](../../1-Expression/UCAM-Assemblage_Module/UCAM02.md)
- [Traçabilité UCAM02 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCAM02)

**Exigences fonctionnelles couvertes**
- [EF19 — Visualiser les interfaces physiques d'un composant](../../1-Expression/Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Use cases cités**
- [UCAM01 — Liaison entre interfaces](UCAM01.md)
- [UCAM03 — Créer une interface sur un composant](UCAM03.md)

**Règles métier**
- [RM11 — Critères de compatibilité d'interfaces](../../1-Expression/Regles_Metier.md#3.%20Interfaces%20et%20liaisons)
- [RM13 — Slot virtuel garanti](../../1-Expression/Regles_Metier.md#3.%20Interfaces%20et%20liaisons)

**Exigences non fonctionnelles**
- [ENF12 — Contrôle d'accès par rôle](../../1-Expression/Exigences_Non_Fonctionnelles.md#Tableau%20des%20exigences%20non-fonctionnelles)

**Documents cités**
- [Conception_intro](../../3-Conception/Conception_intro.md)
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)

**Cité par**
- [Expression_des_besoins_Intro](../../1-Expression/Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../../1-Expression/Matrice_Tracabilite.md)
- [UCAM01 (expression)](../../1-Expression/UCAM-Assemblage_Module/UCAM01.md)
- [todo (expression)](../../1-Expression/todo.md)
- [Analyse_des_besoins](../Analyse_des_besoins.md)
- [UCDEV02 (analyse)](../UCDEV-Developpement/UCDEV02.md)
- [UCMOD03 (analyse)](../UCMOD-Module/UCMOD03.md)
- [todo (analyse)](../todo.md)
- [API_REST](../../3-Conception/API_REST.md)
- [Architecture_Composition](../../3-Conception/Architecture_Composition.md)
- [DC_CLI_Model](../../3-Conception/DC_CLI_Model.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
