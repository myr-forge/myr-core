---
tags:
  - couche/conception
  - type/conception
  - domaine/model
  - uc/UCAM01
  - uc/UCAM02
  - uc/UCAM03
  - uc/UCAM05
  - uc/UCAM07
  - uc/UCAM08
  - uc/UCCE01
  - uc/UCCE02
  - uc/UCCE03
  - uc/UCCE04
  - uc/UCCE05
  - uc/UCCE06
  - uc/UCCL01
  - uc/UCMOD01
  - uc/UCMOD02
  - uc/UCMOD03
  - uc/UCMOD04
  - uc/UCMOD06
  - uc/UCREC01
  - uc/UCREC02
  - uc/UCREC03
  - uc/UCREC04
  - uc/UCREC05
  - rm/RM01
  - rm/RM03
  - rm/RM09
  - rm/RM10
  - rm/RM11
  - rm/RM13
  - rm/RM14
  - rm/RM15
  - rm/RM16
  - rm/RM17
  - rm/RM19
  - rm/RM27
  - rm/RM39
  - rm/RM40
  - rm/RM41
---
# DC — CLI Modèle : Référence des commandes composant / interfaces / module

> Phase 3 — Arrington | Use cases : UCCE01–06, UCAM01–03/05/07/08, UCMOD01–06, UCCL01, UCREC01–05 | Outil : `myr` (`bin/myr-cli`)

---

## 1. Objectif

Ce document définit le **contrat d'interface CLI** couvrant les capacités du domaine `model` (composants, interfaces, liaisons, instances, modules). Il joue pour ces use cases le rôle que `DC_CLI_Admin.md` joue pour l'administration réseau (UCADM01–05) : référence unique du nommage des commandes, afin que les use cases de `specs/1-Expression/` et `specs/2-Analyse/` pointent vers une nomenclature stable.

**Parité CLI/REST :** voir [Architecture_Hexagonale.md](Architecture_Hexagonale.md) (CLI Handler et REST Handler consomment tous deux `ModelService`) — ce document liste les commandes `myr model`, chacune ayant un appel API REST équivalent (flux nominal des deux côtés).

**Qui exécute ces commandes ?** Conformément au principe d'exécution distante, le CLI ne tourne jamais sur le poste d'un Concepteur ou d'un Consommateur — uniquement sur le serveur, via SSH. Les commandes ci-dessous sont donc typiquement exécutées par l'administrateur du serveur **pour le compte d'une identité** (`--owner-id`, `--as`), à des fins de script, d'import en masse, de support ou de restauration ; un client de l'API REST (interface graphique tierce, plugin, boutique partenaire…) peut exécuter la même action directement pour son propre compte.

---

## 2. Arbre de commandes

```
myr model
├── add <file>                          — publier un composant                            UCCE01/03/04/05
│                                          (--draft pour soumission différée, RM16)
├── submit <id>                         — soumettre un composant en brouillon              UCCE01, RM16/RM19
├── get <id>                            — afficher un composant                            UCCL01
├── list                                — lister les composants                            UCCL01
├── verify <id>                         — vérifier l'intégrité
├── location check <id>                 — vérifier l'accessibilité des emplacements externes UCCL03
├── update <id>                         — modifier les métadonnées                         UCCE02
├── remove <id>                         — retirer un composant
├── children <parentID>                 — lister les dérivés d'un composant                UCREC03
├── thumbnail set <id> <fichier>        — associer une miniature
├── thumbnail get <id>                  — récupérer la miniature
├── thumbnail regenerate <id>           — redériver la miniature depuis le lien source (og:image)
├── license list                        — lister le catalogue de licences
├── license get <id>                    — afficher une licence
├── license check                       — vérifier une compatibilité de licence            UCCE04
├── interface
│   ├── add <assetID>                   — définir une interface                            UCAM03, UCCE06
│   ├── update <id>                     — modifier une interface
│   ├── remove <id>                     — supprimer une interface
│   ├── list <assetID>                  — visualiser les interfaces d'un asset             UCAM02
│   └── get <id>                        — afficher une interface
├── ref
│   ├── list                            — afficher le vocabulaire (catégories/types/unités)
│   ├── add-category <cat>              — étendre le vocabulaire
│   ├── add-type <cat> <type>
│   └── add-unit <cat> <unit>
├── link
│   ├── add                             — créer une liaison entre deux interfaces          UCAM01, UCAM07
│   ├── connect-virtual                 — relier un slot virtuel à une interface physique   UCAM03
│   ├── remove <id>                     — supprimer une liaison
│   └── list                            — visualiser les liaisons d'un module               UCAM02
├── instance
│   ├── add <assetID> <subAssetID>      — ajouter un composant existant comme instance d'un asset UCMOD01, UCAM05
│   └── remove <assetID> <instanceID>   — retirer une instance (cascade)                    UCAM08
├── assembly
│   ├── add <assetID> <connID>          — rattacher une liaison à l'asset                   UCMOD02
│   └── remove <assetID> <connID>       — détacher une liaison de l'asset
└── decompose
    ├── preview <assetID>               — analyser un composant STEP, proposer un découpage UCAM09
    └── commit <assetID> --from <id>    — matérialiser la proposition retenue               UCAM09
```

**Fusion composant/module (ADR-11, `Conception_intro.md`) :** `myr module` n'est plus un arbre de commandes distinct — `Model3D` est une seule ressource, `myr model` en est l'unique point d'entrée CLI, qu'un asset ait ou non des instances. `add <assetID> <subAssetID>` s'applique à n'importe quel `assetID` existant, y compris un composant qui n'a encore aucune instance ; il n'y a plus de commande `create` séparée pour « créer un module » — un module se crée comme n'importe quel composant (`myr model add`), sa nature d'assemblage vient ensuite de l'ajout d'instances. `interfaces`/`list`/`get`/`submit`/`remove` d'un asset ayant des instances passent par les mêmes commandes `myr model interface list`/`list`/`get`/`submit`/`remove` qu'un composant simple — `submit` route en interne vers la logique `SubmitModule` (RM17, `ModuleVersion`) dès que l'asset a des `Assemblies`, voir ADR-11. Seul `assembly add/remove` (ex-`module add-assembly/remove-assembly`) reste nommé distinctement sous `myr model`, faute d'équivalent générique existant pour rattacher une connexion déjà créée à la liste `Assemblies` d'un asset.

Voir § 5 pour la table de correspondance complète entre méthode du service domaine, commande CLI et use case.

---

## 3. Détail des commandes les plus citées

### 3.1 `myr model add`

```
myr model add <file> --name <nom> --channel <id> [--category <cat>] [--parent <id>] [--license <id>] [--description <texte>] [--tags <a,b>] [--owner-id <id>] [--draft]
```

Appelle `modelSvc.AddFull(AddRequest{...})`, qui expose `Category`, `ParentID`, `LicenseID`, `Description` en plus de `Name`/`Channel`/`Tags` — le strict équivalent CLI de `POST /api/components` (voir UCCE01, UCCE03, UCCE04, UCCE05 en `2-Analyse`).

**`--draft` (RM16, ADR-02 `Conception_intro.md`) :** optionnel, `false` par défaut — le comportement nominal (une seule transaction Fabric, composant `submitted` immédiatement) reste inchangé. Si `--draft` est passé, le composant est créé `Status: draft` : ses interfaces (`myr model interface add/update`, UCCE06/UCAM03) restent alors en brouillon local jusqu'à `myr model submit` (§ 3.6bis).

### 3.2 `myr model interface add`

```
myr model interface add <assetID> --category <ELEC|MECA|HYD|...> --type <type> --direction <in|out|bidir> [--value-min <f>] [--value-max <f>] [--unit <u>] [--name <label>]
```

Appelle `ModelService.AddInterface(*AssetInterface)`. Équivalent CLI de UCAM03 (« Créer une interface ») et UCCE06 (« Ajouter une interface »).

### 3.3 `myr model interface list`

```
myr model interface list <assetID>
```

Appelle `ModelService.ListInterfacesForAsset(assetID)` (ou `ModelService.GetModuleInterfaces(id)` si l'ID désigne un module — la commande détecte le type via `Get`/`GetModule`). Équivalent CLI de UCAM02 (« Visualiser les interfaces »). Sortie texte : une ligne par interface (`id`, `category`, `type`, `direction`, valeur/plage, `unit`, `virtual`).

### 3.4 `myr model link add`

```
myr model link add --from <ifaceID> --to <ifaceID> [--fastener <assetID>] [--label <texte>] [--from-instance <id>] [--to-instance <id>]
```

Appelle `ModelService.AddAssemblyLink(fromIfaceID, toIfaceID, label, fromInstanceID, toInstanceID, fastenerAssetID)`. Équivalent CLI de UCAM01 (« Liaison entre interfaces ») et UCAM07 (« Asset d'accroche » — flag `--fastener`). La vérification de compatibilité (RM10/RM11) est faite par le service — comportement strictement identique quel que soit le canal (CLI ou REST), les interfaces étant désignées par identifiant explicite.

### 3.5 `myr model instance add` / `remove`

```
myr model instance add <assetID> <subAssetID>
myr model instance remove <assetID> <instanceID>
```

Appellent respectivement `AddAssetToWorkspace` et `RemoveAssetFromWorkspace`. `add` s'applique à n'importe quel `assetID` existant (ADR-11) — une invocation ajoute un composant comme instance, qu'il s'agisse du premier assemblage (UCMOD01, UCAM05) ou d'un ajout à une composition déjà en cours ; un placement en lot reste un script shell côté appelant (boucle sur cette commande), pas une commande dédiée. `remove` est l'équivalent CLI de UCAM08 (retrait en cascade — la cascade des connexions est gérée par le service, RM15).

### 3.6 `myr model assembly add` / `remove` (ex-`myr module add-assembly`/`remove-assembly`, ADR-11)

```
myr model assembly add <assetID> <connID>
myr model assembly remove <assetID> <connID>
```

Appellent `AddAssemblyToModule(assetID, connID)` / `RemoveAssemblyFromModule(assetID, connID)` — rattachent/détachent une connexion déjà créée (`myr model link add`) à la liste `Assemblies` de l'asset. Équivalent CLI de UCMOD02.

Il n'existe plus de commande `create` séparée : un asset destiné à devenir un assemblage se crée exactement comme n'importe quel composant (`myr model add`), sa nature d'assemblage venant ensuite de l'ajout d'instances (§3.5) — voir ADR-11, point 3.

### 3.6bis `myr model submit` (RM16/RM17/RM19 — soumission, composant ou assemblage)

```
myr model submit <assetID> [--note <texte>]
```

Équivalent CLI unique de UCCE01 (« Flux alternatif — Création en brouillon »), UCCE06 (fork) et UCMOD06 (soumission d'un module). Committe l'état courant du brouillon (interfaces incluses, `Model3D.Interfaces`) sur Fabric en une transaction et passe `Status` à `submitted`. **Conséquence de la fusion (ADR-11) :** cette commande route en interne vers la logique `SubmitModule` (vérification RM17 — `len(Assemblies) > 0` — et création d'une `ModuleVersion` horodatée, flag `--note`) si `m.IsModule()` est vrai, vers la logique `Submit` générique sinon ; ce n'est plus à l'appelant de savoir laquelle des deux méthodes domaine invoquer. **Écart de conception (E8, `specs/2-Analyse/Analyse_des_besoins.md` § Écarts structurels connus) :** `Submit` généralise déjà la logique de persistance de `SubmitModule` — reste à en faire le point de dispatch unique décrit ci-dessus plutôt que deux méthodes séparées exposées séparément ; voir § 6 point 4.

### 3.7 Transformation composant → module (UCAM05, catégorie `decoupage`)

```
myr model add --category decoupage --parent-id <composantOrigine> --name <nom> [...]
myr model instance add <nouvelID> <subAssetID>   # répété pour chaque sous-composant
myr model link add --from <ifaceID> --to <ifaceID> ...   # répété pour chaque liaison
myr model submit <nouvelID>
```

Aucune commande dédiée : depuis la fusion (ADR-11), la transformation manuelle décrite par UCAM05 est la composition des commandes génériques ci-dessus — `myr model add` avec `Category = decoupage` (cf. écart E1 dans `Architecture_Composition.md`, constante `CategoryDecoupage` encore absente d'`entity.go`) crée un **nouvel** asset distinct, `parent_id` référençant le composant d'origine, qui n'est lui-même ni modifié ni supprimé (post-condition UCAM09, RM39). Il n'y a donc plus d'ancienne commande `to-module` à documenter séparément. Ce point (nouvel id systématique vs. mutation du composant d'origine en place) reste néanmoins un point ouvert pour le PO — voir ADR-11.

### 3.8 `myr model decompose preview` / `commit` (UCAM09 — décomposition assistée d'un composant STEP)

```
myr model decompose preview <assetID>
myr model decompose commit <assetID> --from <decomposition_id>
```

`preview` analyse le fichier STEP/STP du composant `<assetID>` et retourne, sous un `decomposition_id` temporaire, une proposition de sous-pièces et de connexions candidates — sans créer aucune entité (RM40). `commit` reprend cette proposition, éventuellement corrigée côté appelant (sous-pièces retirées/renommées, connexions rejetées), et matérialise sous-composants (draft), module de catégorie `decoupage` et liaisons compatibles (RM39/RM41) — même sémantique d'identité que §3.7 (nouvel asset distinct, composant d'origine inchangé). Ce couple de commandes n'a pas encore de méthode `ModelService` dédiée — voir § 6 point 5.

---

## 4. Format de sortie et erreurs

Mêmes conventions que `DC_CLI_Admin.md` § 7 : succès sur stdout, erreurs sur stderr via `RunE`, tableaux alignés au tabwriter. Les messages d'erreur métier (anti-plagiat RM01, licence incompatible RM03, interface déjà utilisée RM09, seuil de nœuds RM27 non applicable ici…) sont ceux renvoyés par le service domaine — identiques à ceux de l'API REST, seul le canal de sortie change (texte terminal vs JSON HTTP).

---

## 5. Table de correspondance méthode domaine → commande CLI → use case

| Méthode `ModelService` (`domain/model/port_in.go`) | Commande CLI | Use case(s) |
|---|---|---|
| `Add` / `AddFull` (avec `Draft bool`) | `myr model add [--draft]` | UCCE01, UCCE03, UCCE04, UCCE05 |
| `Submit` (généralise `SubmitModule` à tout `Model3D`) | `myr model submit` | UCCE01 (brouillon), UCCE06 (RM16/RM19) |
| `Get` | `myr model get` | UCCL01, UCREC01 |
| `List` | `myr model list` | UCCL01 |
| `Verify` | `myr model verify` | — (intégrité, hors périmètre user) |
| `CheckLocations` | `myr model location check` | UCCL03 |
| `UpdateAsset` | `myr model update` | UCCE02 |
| `Remove` | `myr model remove` | — |
| `GetChildren` | `myr model children` | UCREC03, UCREC04 |
| `SaveThumbnail` / `GetThumbnail` | `myr model thumbnail set/get` | — |
| `RegenerateThumbnail` | `myr model thumbnail regenerate` | — |
| `AddInterface` / `UpdateInterface` / `RemoveInterface` | `myr model interface add/update/remove` | UCAM03, UCCE06 |
| `ListInterfacesForAsset` / `GetInterface` | `myr model interface list/get` | UCAM02 |
| `EnsureVirtualSlot` | (appelé automatiquement par `interface list`) | UCAM03 (RM13) |
| `GetRefs` / `AddRefCategory` / `AddRefType` / `AddRefUnit` | `myr model ref list/add-category/add-type/add-unit` | UCAM03 |
| `AddConnection` / `AddAssemblyLink` / `RemoveConnection` / `ListConnections` | `myr model link add/remove/list` | UCAM01, UCAM07 |
| `ConnectVirtualToPhysical` | `myr model link connect-virtual` | UCAM03 |
| `CreateModule` / `GetModule` / `ListModules` / `RemoveModule` (ADR-11 : retirés en tant que commandes CLI distinctes, remplacés par `Add`/`AddFull`, `Get`, `List`, `Remove` déjà génériques) | `myr model add/get/list/remove` | UCMOD01, UCMOD02, UCMOD04 |
| `AddAssemblyToModule` / `RemoveAssemblyFromModule` | `myr model assembly add/remove` | UCMOD02 |
| `SubmitModule` (ADR-11 : point de dispatch interne de `myr model submit`, plus une commande distincte) | `myr model submit` | UCMOD06 |
| `GetModuleInterfaces` (ADR-11 : point de dispatch interne de `myr model interface list`, plus une commande distincte) | `myr model interface list` | UCAM02 |
| `AddAssetToWorkspace` / `RemoveAssetFromWorkspace` | `myr model instance add/remove` | UCMOD01, UCAM05 (add) · UCAM08 (remove) |
| `ListLicenses` / `GetLicense` / `CheckLicenseCompatibility` / `CheckModuleLicenseCompatibility` | `myr model license list/get/check` | UCCE04, UCMOD06 |
| (à concevoir — analyse STEP, aucune persistance) | `myr model decompose preview` | UCAM09, RM40 |
| (à concevoir — s'appuie sur `CreateModule`/`AddAssetToWorkspace`/`AddAssemblyLink`) | `myr model decompose commit` | UCAM09, RM39, RM41 |

Recherches et export (UCCL01, UCREC01–05) se combinent à partir de `List`, `Get`, `GetChildren`, `GetModuleInterfaces` côté service. Le point ouvert de savoir si `ModelService` doit exposer une méthode de filtre serveur dédiée (`Search(criteria)`), plutôt que de laisser le filtrage au client sur le résultat de `List`, est documenté en § 6 point 2 ; la commande `myr model search --filter <critère>` dépend de cette décision. Le détail des méthodes pour UCREC02–05 (`FindCompatibleAssets`, `GetLineage`, `FindModulesUsingComponent`, `ResolveBOM`) est conçu dans `DC_D8_Recherche.md`, pas ici.

---

## 6. Écarts et points ouverts

| # | Écart / question | Impact |
|---|---|---|
| 1 | Depuis la fusion composant/module (ADR-11, `Conception_intro.md`), UCAM05 ne nécessite plus de méthode `ModelService` dédiée ni de commande CLI dédiée (§3.7) — elle se compose entièrement de méthodes déjà génériques (`AddFull`, `AddAssetToWorkspace`, `AddAssemblyLink`, `Submit`). Le seul écart restant est E1 (`Architecture_Composition.md`) : la constante `CategoryDecoupage` est absente d'`entity.go`. | UCAM05 non exposable tant que E1 n'est pas résolu, quel que soit le canal (REST ou CLI) — ce n'est pas un écart spécifique au CLI, et il ne dépend plus de la fusion ADR-11 (déjà réglée côté conception). |
| 2 | `ModelService` n'expose aucune méthode de filtre serveur (`Search(criteria)`) — `List(channelID)` retourne tout le canal, à charge du dépôt GUI externe de filtrer. Reste à trancher si le filtrage doit devenir un comportement serveur. | UCCL01 et UCREC01–05 : la commande `myr model search` ne peut être qu'un alias de `list` tant que cette décision n'est pas prise et le filtre remonté côté domaine. |
| 3 | Tarification, commission, transfert de PI, clonage inter-réseau, écoconception (UCPI01/02/04/05/06/07/08/09/10/11) et automatisation (UCAUT01/02/04) n'ont aucun port domaine ni entité correspondante (`Price`, `Commission`, `Transfer`…absents de `domain/model` et `domain/payment`). Documentés dans les UC concernés comme commandes CLI de niveau 2 : le nom de commande est proposé, mais dépend d'abord de la conception du domaine (hors périmètre de ce document). | Pas d'implémentation CLI possible avant modélisation du domaine correspondant. |
| 4 | `myr model submit` (§ 3.6bis) n'a pas de méthode `ModelService` dédiée : elle nécessite un changement domaine (voir `specs/2-Analyse/Analyse_des_besoins.md` § Écarts structurels connus, E8) — ajouter `Status`/`Draft` à `AddRequest` et une méthode `Submit` généralisant `SubmitModule` à tout `Model3D`. | UCCE01 (flux brouillon) et UCCE06 (fork) dépendent de ce changement domaine, contrairement aux autres commandes de ce document qui n'exigent qu'un adaptateur CLI sur des méthodes `ModelService` déjà définies. |
| 5 | `myr model decompose preview/commit` (UCAM09, § 3.8) n'a pas non plus de méthode `ModelService` dédiée — `preview` dépend en plus d'une analyse géométrique du fichier STEP (arbre d'assemblage, détection de contacts/coaxialité) et du repère géométrique à ajouter sur `AssetInterface` (écart E9, `specs/2-Analyse/Analyse_des_besoins.md` § Écarts structurels connus). Orientation retenue pour la technologie d'analyse : une librairie Go native — mais aucune librairie mature ne couvre aujourd'hui à la fois le parsing STEP AP214/AP242 et la détection géométrique de contacts (voir § 8 ci-dessous, point ouvert). Mécanisme retenu pour `preview` : requête HTTP bloquante avec budget de temps serveur (pas de job asynchrone), cohérent avec le reste du contrat REST. | UCAM09 non exposable tant que E1 (catégorie `decoupage`), E9 (repère géométrique) et le choix concret de la librairie/l'approche de parsing (§ 8) ne sont pas tranchés. |

---

## 7. Décisions de conception

| ID | Décision | Raison |
|----|---------|--------|
| DC-CLIM-01 | Un seul groupe `myr model` porte composants, interfaces, liaisons et instances ; `myr module` reste séparé pour les opérations propres aux modules (création, soumission) | Miroir de la distinction domaine `Model3D` (composant/module unifié) vs. use cases (UCCE/UCAM d'un côté, UCMOD de l'autre) — évite un groupe `myr model` démesuré tout en gardant `model add/get/list/verify` stables (rétrocompatibilité de la commande existante) |
| DC-CLIM-02 | Les commandes CLI n'ajoutent aucune vérification propre — elles délèguent entièrement au service domaine | Garantit que le comportement (RM01, RM03, RM09-11, RM13-15) est strictement identique quel que soit le canal (CLI ou REST), conformément au principe de parité fonctionnelle |
| DC-CLIM-03 | Le CLI/API ne fait que des actions brutes et directes — les identifiants (interface, instance, asset) sont fournis explicitement en argument/flag, jamais par sélection interactive | `interface list` / `link list` permettent de retrouver les identifiants nécessaires avant d'agir ; toute ergonomie de sélection visuelle relève exclusivement d'un client externe (dépôt GUI) |
| DC-CLIM-04 | `myr model decompose preview` (UCAM09, § 3.8) reste une requête bloquante avec un budget de temps serveur — pas de mécanisme job + polling | Cohérence avec le reste du contrat CLI/REST, qui n'a aucun autre pattern asynchrone ; à revoir si le temps de traitement réel sur de gros assemblages dépasse le budget en pratique, une fois la technologie d'analyse STEP choisie (§ 8) |

---

## 8. Point ouvert — technologie d'analyse STEP (UCAM09)

Orientation retenue : une librairie **Go native**, cohérente avec le stack 100 % Go actuel (pas de dépendance cgo/binaire externe). Recherche effectuée sur l'écosystème disponible : aucune librairie Go mature ne couvre aujourd'hui les deux besoins d'UCAM09 :

1. **Parsing de la structure d'assemblage** (fichier texte STEP Part 21 / schéma EXPRESS AP214-AP242 : occurrences de produit, transformations relatives) — écrire un parseur Go dédié à ce sous-ensemble du format est raisonnable (format texte documenté, pas de dépendance à un noyau géométrique).
2. **Détection géométrique des contacts/coaxialités entre sous-pièces** (§2.2 de la proposition d'origine) — nécessite une représentation B-rep et des calculs de géométrie solide qu'aucune librairie Go connue ne fournit ; c'est typiquement le rôle d'un noyau CAO (OpenCASCADE et équivalents), aujourd'hui hors de la table des technologies autorisées.

**Ce point reste ouvert** : soit un parseur Go maison se limite dans un premier temps à la structure d'assemblage et à une détection de contact approximative (recouvrement de boîtes englobantes plutôt qu'analyse de faces), avec un score de confiance revu à la baisse en conséquence ; soit le besoin de précision impose de revisiter l'option d'un noyau géométrique externe malgré le coût d'intégration. Décision à prendre avant toute implémentation (règle 20 CLAUDE.md) et à documenter ici une fois tranchée.

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Étape suivante — code, par section**
- [§ 1. Objectif](DC_CLI_Model.md#1.%20Objectif) → [CLI myr model](../../docs/code/commandes/CLI%20myr-model.md)
- [§ 2. Arbre de commandes](DC_CLI_Model.md#2.%20Arbre%20de%20commandes) → [CLI myr model](../../docs/code/commandes/CLI%20myr-model.md) · [CLI myr model add](../../docs/code/commandes/CLI%20myr-model-add.md) · [CLI myr model interface list](../../docs/code/commandes/CLI%20myr-model-interface-list.md) · [ModelService.SubmitModule](../../docs/code/fonctions/model.ModelService.SubmitModule.md)
- [§ 3.1 `myr model add`](DC_CLI_Model.md#3.1%20`myr%20model%20add`) → [REST /api/components/](../../docs/code/routes/REST%20api-components.md) · [CLI myr model add](../../docs/code/commandes/CLI%20myr-model-add.md) · [CLI myr model interface add](../../docs/code/commandes/CLI%20myr-model-interface-add.md) · [CLI myr model submit](../../docs/code/commandes/CLI%20myr-model-submit.md) · [ModelService.AddFull](../../docs/code/fonctions/model.ModelService.AddFull.md)
- [§ 3.2 `myr model interface add`](DC_CLI_Model.md#3.2%20`myr%20model%20interface%20add`) → [CLI myr model interface add](../../docs/code/commandes/CLI%20myr-model-interface-add.md) · [ModelService.AddInterface](../../docs/code/fonctions/model.ModelService.AddInterface.md)
- [§ 3.3 `myr model interface list`](DC_CLI_Model.md#3.3%20`myr%20model%20interface%20list`) → [CLI myr model interface list](../../docs/code/commandes/CLI%20myr-model-interface-list.md) · [InterfaceStore.ListInterfacesForAsset](../../docs/code/fonctions/model.InterfaceStore.ListInterfacesForAsset.md) · [ModelService.Get](../../docs/code/fonctions/model.ModelService.Get.md) · [ModelService.GetModule](../../docs/code/fonctions/model.ModelService.GetModule.md) · [ModelService.GetModuleInterfaces](../../docs/code/fonctions/model.ModelService.GetModuleInterfaces.md) · [ModelService.ListInterfacesForAsset](../../docs/code/fonctions/model.ModelService.ListInterfacesForAsset.md)
- [§ 3.4 `myr model link add`](DC_CLI_Model.md#3.4%20`myr%20model%20link%20add`) → [CLI myr model link add](../../docs/code/commandes/CLI%20myr-model-link-add.md) · [ModelService.AddAssemblyLink](../../docs/code/fonctions/model.ModelService.AddAssemblyLink.md)
- [§ 3.5 `myr model instance add` / `remove`](DC_CLI_Model.md#3.5%20`myr%20model%20instance%20add`%20/%20`remove`) → [CLI myr model instance add](../../docs/code/commandes/CLI%20myr-model-instance-add.md) · [CLI myr model instance remove](../../docs/code/commandes/CLI%20myr-model-instance-remove.md) · [ModelService.AddAssetToWorkspace](../../docs/code/fonctions/model.ModelService.AddAssetToWorkspace.md) · [ModelService.RemoveAssetFromWorkspace](../../docs/code/fonctions/model.ModelService.RemoveAssetFromWorkspace.md)
- [§ 3.6 `myr model assembly add` / `remove` (ex-`myr module add-assembly`/`remove-assembly`, ADR-11)](DC_CLI_Model.md#3.6%20`myr%20model%20assembly%20add`%20/%20`remove`%20%28ex-`myr%20module%20add-assembly`/`remove-assembly`,%20ADR-11%29) → [CLI myr model assembly add](../../docs/code/commandes/CLI%20myr-model-assembly-add.md) · [CLI myr model assembly remove](../../docs/code/commandes/CLI%20myr-model-assembly-remove.md) · [CLI myr model link add](../../docs/code/commandes/CLI%20myr-model-link-add.md) · [CLI myr model add](../../docs/code/commandes/CLI%20myr-model-add.md) · [ModelService.AddAssemblyToModule](../../docs/code/fonctions/model.ModelService.AddAssemblyToModule.md) · [ModelService.RemoveAssemblyFromModule](../../docs/code/fonctions/model.ModelService.RemoveAssemblyFromModule.md)
- [§ 3.6bis `myr model submit` (RM16/RM17/RM19 — soumission, composant ou assemblage)](DC_CLI_Model.md#3.6bis%20`myr%20model%20submit`%20%28RM16/RM17/RM19%20—%20soumission,%20composant%20ou%20assemblage%29) → [CLI myr model submit](../../docs/code/commandes/CLI%20myr-model-submit.md) · [ModelService.Submit](../../docs/code/fonctions/model.ModelService.Submit.md) · [ModelService.SubmitModule](../../docs/code/fonctions/model.ModelService.SubmitModule.md)
- [§ 3.7 Transformation composant → module (UCAM05, catégorie `decoupage`)](DC_CLI_Model.md#3.7%20Transformation%20composant%20→%20module%20%28UCAM05,%20catégorie%20`decoupage`%29) → [CLI myr model add](../../docs/code/commandes/CLI%20myr-model-add.md) · [CLI myr model instance add](../../docs/code/commandes/CLI%20myr-model-instance-add.md) · [CLI myr model link add](../../docs/code/commandes/CLI%20myr-model-link-add.md) · [CLI myr model submit](../../docs/code/commandes/CLI%20myr-model-submit.md)
- [§ 3.8 `myr model decompose preview` / `commit` (UCAM09 — décomposition assistée d'un composant STEP)](DC_CLI_Model.md#3.8%20`myr%20model%20decompose%20preview`%20/%20`commit`%20%28UCAM09%20—%20décomposition%20assistée%20d'un%20composant%20STEP%29) → [CLI myr model](../../docs/code/commandes/CLI%20myr-model.md)
- [§ 5. Table de correspondance méthode domaine → commande CLI → use case](DC_CLI_Model.md#5.%20Table%20de%20correspondance%20méthode%20domaine%20→%20commande%20CLI%20→%20use%20case) → [CLI myr model add](../../docs/code/commandes/CLI%20myr-model-add.md) · [CLI myr model submit](../../docs/code/commandes/CLI%20myr-model-submit.md) · [CLI myr model get](../../docs/code/commandes/CLI%20myr-model-get.md) · [CLI myr model list](../../docs/code/commandes/CLI%20myr-model-list.md) · [CLI myr model verify](../../docs/code/commandes/CLI%20myr-model-verify.md) · [CLI myr model](../../docs/code/commandes/CLI%20myr-model.md) · [CLI myr model update](../../docs/code/commandes/CLI%20myr-model-update.md) · [CLI myr model remove](../../docs/code/commandes/CLI%20myr-model-remove.md) · [CLI myr model children](../../docs/code/commandes/CLI%20myr-model-children.md) · [CLI myr model thumbnail set](../../docs/code/commandes/CLI%20myr-model-thumbnail-set.md) · [CLI myr model thumbnail regenerate](../../docs/code/commandes/CLI%20myr-model-thumbnail-regenerate.md) · [CLI myr model interface add](../../docs/code/commandes/CLI%20myr-model-interface-add.md) · [CLI myr model interface list](../../docs/code/commandes/CLI%20myr-model-interface-list.md) · [CLI myr model ref list](../../docs/code/commandes/CLI%20myr-model-ref-list.md) · [CLI myr model link add](../../docs/code/commandes/CLI%20myr-model-link-add.md) · [CLI myr model link connect-virtual](../../docs/code/commandes/CLI%20myr-model-link-connect-virtual.md) · [CLI myr model assembly add](../../docs/code/commandes/CLI%20myr-model-assembly-add.md) · [CLI myr model instance add](../../docs/code/commandes/CLI%20myr-model-instance-add.md) · [CLI myr model license list](../../docs/code/commandes/CLI%20myr-model-license-list.md) · [ConnectionStore.ListConnections](../../docs/code/fonctions/model.ConnectionStore.ListConnections.md) · [ConnectionStore.RemoveConnection](../../docs/code/fonctions/model.ConnectionStore.RemoveConnection.md) · [FileStoragePort.Verify](../../docs/code/fonctions/model.FileStoragePort.Verify.md) · [InterfaceStore.AddRefCategory](../../docs/code/fonctions/model.InterfaceStore.AddRefCategory.md) · [InterfaceStore.AddRefType](../../docs/code/fonctions/model.InterfaceStore.AddRefType.md) · [InterfaceStore.AddRefUnit](../../docs/code/fonctions/model.InterfaceStore.AddRefUnit.md) · [InterfaceStore.GetInterface](../../docs/code/fonctions/model.InterfaceStore.GetInterface.md) · [InterfaceStore.GetRefs](../../docs/code/fonctions/model.InterfaceStore.GetRefs.md) · [InterfaceStore.ListInterfacesForAsset](../../docs/code/fonctions/model.InterfaceStore.ListInterfacesForAsset.md) · [InterfaceStore.RemoveInterface](../../docs/code/fonctions/model.InterfaceStore.RemoveInterface.md) · [ModelService.Add](../../docs/code/fonctions/model.ModelService.Add.md) · [ModelService.AddAssemblyLink](../../docs/code/fonctions/model.ModelService.AddAssemblyLink.md) · [ModelService.AddAssemblyToModule](../../docs/code/fonctions/model.ModelService.AddAssemblyToModule.md) · [ModelService.AddAssetToWorkspace](../../docs/code/fonctions/model.ModelService.AddAssetToWorkspace.md) · [ModelService.AddConnection](../../docs/code/fonctions/model.ModelService.AddConnection.md) · [ModelService.AddFull](../../docs/code/fonctions/model.ModelService.AddFull.md) · [ModelService.AddInterface](../../docs/code/fonctions/model.ModelService.AddInterface.md) · [ModelService.AddRefCategory](../../docs/code/fonctions/model.ModelService.AddRefCategory.md) · [ModelService.AddRefType](../../docs/code/fonctions/model.ModelService.AddRefType.md) · [ModelService.AddRefUnit](../../docs/code/fonctions/model.ModelService.AddRefUnit.md) · [ModelService.CheckLicenseCompatibility](../../docs/code/fonctions/model.ModelService.CheckLicenseCompatibility.md) · [ModelService.CheckModuleLicenseCompatibility](../../docs/code/fonctions/model.ModelService.CheckModuleLicenseCompatibility.md) · [ModelService.ConnectVirtualToPhysical](../../docs/code/fonctions/model.ModelService.ConnectVirtualToPhysical.md) · [ModelService.CreateModule](../../docs/code/fonctions/model.ModelService.CreateModule.md) · [ModelService.EnsureVirtualSlot](../../docs/code/fonctions/model.ModelService.EnsureVirtualSlot.md) · [ModelService.Get](../../docs/code/fonctions/model.ModelService.Get.md) · [ModelService.GetChildren](../../docs/code/fonctions/model.ModelService.GetChildren.md) · [ModelService.GetInterface](../../docs/code/fonctions/model.ModelService.GetInterface.md) · [ModelService.GetLicense](../../docs/code/fonctions/model.ModelService.GetLicense.md) · [ModelService.GetModule](../../docs/code/fonctions/model.ModelService.GetModule.md) · [ModelService.GetModuleInterfaces](../../docs/code/fonctions/model.ModelService.GetModuleInterfaces.md) · [ModelService.GetRefs](../../docs/code/fonctions/model.ModelService.GetRefs.md) · [ModelService.GetThumbnail](../../docs/code/fonctions/model.ModelService.GetThumbnail.md) · [ModelService.List](../../docs/code/fonctions/model.ModelService.List.md) · [ModelService.ListConnections](../../docs/code/fonctions/model.ModelService.ListConnections.md) · [ModelService.ListInterfacesForAsset](../../docs/code/fonctions/model.ModelService.ListInterfacesForAsset.md) · [ModelService.ListLicenses](../../docs/code/fonctions/model.ModelService.ListLicenses.md) · [ModelService.ListModules](../../docs/code/fonctions/model.ModelService.ListModules.md) · [ModelService.RegenerateThumbnail](../../docs/code/fonctions/model.ModelService.RegenerateThumbnail.md) · [ModelService.Remove](../../docs/code/fonctions/model.ModelService.Remove.md) · [ModelService.RemoveAssemblyFromModule](../../docs/code/fonctions/model.ModelService.RemoveAssemblyFromModule.md) · [ModelService.RemoveAssetFromWorkspace](../../docs/code/fonctions/model.ModelService.RemoveAssetFromWorkspace.md) · [ModelService.RemoveConnection](../../docs/code/fonctions/model.ModelService.RemoveConnection.md) · [ModelService.RemoveInterface](../../docs/code/fonctions/model.ModelService.RemoveInterface.md) · [ModelService.RemoveModule](../../docs/code/fonctions/model.ModelService.RemoveModule.md) · [ModelService.SaveThumbnail](../../docs/code/fonctions/model.ModelService.SaveThumbnail.md) · [ModelService.Submit](../../docs/code/fonctions/model.ModelService.Submit.md) · [ModelService.SubmitModule](../../docs/code/fonctions/model.ModelService.SubmitModule.md) · [ModelService.UpdateAsset](../../docs/code/fonctions/model.ModelService.UpdateAsset.md) · [ModelService.UpdateInterface](../../docs/code/fonctions/model.ModelService.UpdateInterface.md) · [ModelService.Verify](../../docs/code/fonctions/model.ModelService.Verify.md) · [ThumbnailStore.GetThumbnail](../../docs/code/fonctions/model.ThumbnailStore.GetThumbnail.md) · [ThumbnailStore.SaveThumbnail](../../docs/code/fonctions/model.ThumbnailStore.SaveThumbnail.md)
- [§ 6. Écarts et points ouverts](DC_CLI_Model.md#6.%20Écarts%20et%20points%20ouverts) → [CLI myr model](../../docs/code/commandes/CLI%20myr-model.md) · [CLI myr model submit](../../docs/code/commandes/CLI%20myr-model-submit.md) · [ModelService.AddAssemblyLink](../../docs/code/fonctions/model.ModelService.AddAssemblyLink.md) · [ModelService.AddAssetToWorkspace](../../docs/code/fonctions/model.ModelService.AddAssetToWorkspace.md) · [ModelService.AddFull](../../docs/code/fonctions/model.ModelService.AddFull.md) · [ModelService.List](../../docs/code/fonctions/model.ModelService.List.md) · [ModelService.Submit](../../docs/code/fonctions/model.ModelService.Submit.md) · [ModelService.SubmitModule](../../docs/code/fonctions/model.ModelService.SubmitModule.md) · [PaymentGatewayPort.Transfer](../../docs/code/fonctions/payment.PaymentGatewayPort.Transfer.md)
- [§ 7. Décisions de conception](DC_CLI_Model.md#7.%20Décisions%20de%20conception) → [CLI myr model](../../docs/code/commandes/CLI%20myr-model.md)

<!-- liens-obsidian:end -->
