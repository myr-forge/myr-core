---
id: FT-040
titre: "Fonctions propres aux outils extérieurs implémentées dans myr-core"
type: dette
statut: a-trancher
severite: majeure
detecte: 2026-11-08
maj: 2026-11-09
composants: [domain/model, adapters/in/rest, adapters/in/cli, adapters/out/localstorage, adapters/out/webimage]
uc: [UCAM01, UCAM02, UCAM03, UCCE01, UCCE06, UCMOD01, UCMOD03, UCMOD06]
rm: [RM13, RM16, RM17]
enf: []
tags:
  - ticket
  - ticket/dette
  - statut/a-trancher
  - severite/majeure
  - domaine/model
  - uc/UCAM01
  - uc/UCAM02
  - uc/UCAM03
  - uc/UCCE01
  - uc/UCCE06
  - uc/UCMOD01
  - uc/UCMOD03
  - uc/UCMOD06
  - rm/RM13
  - rm/RM16
  - rm/RM17
---
# FT-040 — Fonctions propres aux outils extérieurs implémentées dans myr-core

> **Dette technique** · sévérité **majeure** · statut **À trancher (PO)** · détecté le 2026-11-08
> Source : décision du PO pendant R-03 étape 3 (immuabilité d'un asset soumis)

## Constat

L'API de `myr-core` doit apporter le strict nécessaire au domaine ; les aides à la conception d'un modèle relèvent des outils qui l'entourent (`myr-web`, plugins CAO). Plusieurs de ces aides sont pourtant implémentées dans le cœur, avec leurs règles, leurs routes, leurs commandes et leur stockage :

| Fonction | Ce que le cœur en porte |
|----------|-------------------------|
| **Slots virtuels** (point de connexion proposé avant qu'une interface existe) | `AssetInterface.Virtual`, `Service.EnsureVirtualSlot`, `Service.ConnectVirtualToPhysical` (`domain/model/service.go`) ; route `/api/virtual-connect` ; commande `myr model link connect-virtual` ; règle RM13 « Slot virtuel garanti » |
| **Disposition à l'écran** (positions x/y des instances d'un module) | `WorkspaceInstance.X`/`Y` (`domain/model/entity.go`), écrits sur le ledger avec le module ; `Service.UpdateInstancePosition` ; `PATCH /api/components/{id}/instances/{instanceID}` |
| **Brouillons** (état modifiable avant soumission) | port `DraftStore` (`domain/model/ports.go`) et son stockage JSON local ; statut `draft` ; `AddFull` et `CreateModule` créent un brouillon, `Submit`/`SubmitModule` l'engagent ; interfaces en cours d'édition (`InterfaceStore`, ADR-02) ; RM16 ; commande `myr model submit`, route `POST /api/components/{id}/submit` |
| **Miniatures** (image d'aperçu) | port `ThumbnailStore`, `Service.SaveThumbnail`/`GetThumbnail`/`RegenerateThumbnail` ; dérivation côté serveur depuis l'`og:image` d'un lien (port `OGImageFetcher`, adapter `adapters/out/webimage` et sa protection SSRF) ; routes `/api/components/{id}/thumbnail[/regenerate]` ; commande `myr model thumbnail` |

Restent dans le cœur : le vocabulaire de référence des interfaces (`/api/refs`) et les connexions libres (`/api/connections`).

## Cause

Ces fonctions ont été ajoutées au fil des besoins de `myr-web`, dans le dépôt qui sert l'API, plutôt que dans l'outil qui les utilise.

## Impact

- Le cœur porte une logique et un stockage d'ergonomie qui ne servent pas le domaine ; chaque évolution de l'interface graphique touche le cœur, ses specs et ses tests.
- Ces fonctions entrent en conflit avec l'immuabilité du ledger (ADR-12) : la disposition d'un module soumis ne peut plus bouger ; la matérialisation d'un slot virtuel ajoute une interface à un composant déjà soumis ; un brouillon réouvert sous le même identifiant qu'un asset soumis ne peut plus être soumis.
- Surface d'attaque évitable : la dérivation de miniature fait télécharger au serveur des URL choisies par l'utilisateur (SSRF, corrigé par R-02).

## Preuves

Voir le tableau ci-dessus ; specs concernées : RM13, RM16 (`Regles_Metier.md`), ADR-02 (`Conception_intro.md`), UCAM01–03, UCCE01, UCCE06, UCMOD01, UCMOD03, UCMOD06 et `API_REST.md`.

## Piste de correction (à valider par le PO)

Le cœur ne reçoit plus que des assets complets, enregistrés en une écriture : un composant avec ses interfaces, un module avec ses instances et ses liaisons. La préparation (brouillon, slots virtuels, disposition, miniature) se fait dans l'outil, qui soumet le résultat.

**À trancher avant tout retrait :**
1. **Où vivent les brouillons et les miniatures ?** `myr-web` est une application cliente : stockage dans le navigateur (propre à un poste), backend propre à `myr-web`, ou autre outil.
2. **Ordre du retrait.** `myr-web` doit d'abord porter ces fonctions ; le cœur les retire ensuite, avec une fenêtre de transition coordonnée (comme pour les routes `/api/modules`, FT-022).

Retrait côté cœur, une fois tranché : ports `DraftStore`, `ThumbnailStore`, `OGImageFetcher` et leurs adapters (dont `adapters/out/webimage`) ; `Virtual`, `X`/`Y` et les méthodes associées ; routes et commandes correspondantes ; specs (RM13 supprimée, RM16 et ADR-02 révisés, use cases réécrits en assets complets) et règles d'agent associées.

## Inventaire du 2026-11-09 — ce que l'API fait à la place du client

Demandé par le PO : lister ce que `myr-core` devine, complète ou facilite au lieu de s'en tenir au strict minimum, avec une proposition pour chaque point, à trancher.

| # | Fonction | Ce que le cœur fait à la place du client | Où | Proposition |
|---|----------|------------------------------------------|----|-------------|
| 1 | Slots virtuels (tableau ci-dessus) | En plus : à la matérialisation, il **déduit** catégorie, type, sens inversé et valeurs de l'interface physique ; une **lecture** (`GET …/interfaces` d'un module, `GetModuleInterfaces`) crée un slot ; après `DELETE /api/interfaces/{id}`, c'est le **handler REST** qui recrée le slot (la CLI ne le fait pas) | `service.go` `ConnectVirtualToPhysical`, `getModuleInterfacesInto` ; `handlers.go` `handleInterface` | Retirer avec les slots (décision du 2026-11-08) |
| 2 | Positions x/y (tableau ci-dessus) | — | — | Retirer (décision du 2026-11-08) |
| 3 | Brouillons (tableau ci-dessus) | — | — | Retirer (décision du 2026-11-08), stockage à trancher |
| 4 | Miniatures (tableau ci-dessus) | En plus : à la **création**, sans miniature fournie, le serveur télécharge l'`og:image` du premier lien | `handlers.go` `applyThumbnailFallback` | Retirer avec les miniatures ; supprimer d'abord ce téléchargement implicite |
| 5 | Catégorie par défaut | `POST /api/components` et `/batch` posent `base` si la catégorie manque ; la CLI la laisse vide | `handlers.go` `createAsset`, `createComponentsBatch` ; `service.go` `Add` | Catégorie obligatoire (RM02), `400` sinon, en REST comme en CLI |
| 6 | Canal par défaut | Sans `channel_id`, repli sur le canal de la session, puis sur celui de `fabric.env` | `handlers.go` `channelFor` | Refuser (`400`) une écriture sans canal explicite |
| 7 | Création par lot | `POST /api/components/batch` crée un asset par fichier et **tire le nom du nom de fichier** | `handlers.go` `createComponentsBatch` | Retirer : le client appelle `POST /api/components` par fichier, avec un nom |
| 8 | Copie de composition à la dérivation | Un asset dérivé d'un module reprend une copie de ses instances et liaisons | `service.go` `AddFull`, `copyModuleComposition` | À trancher : garder (UCMOD01 le décrit) ou laisser le client recopier |
| 9 | Arborescence | `GET /api/components/{id}/tree` calcule ancêtres, descendants et liaisons | `handlers.go` `handleComponentTree` | À trancher : vue dérivée en lecture seule, que le client peut reconstruire depuis `parent_id` |
| 10 | Interfaces exposées d'un module | Calcule les interfaces libres d'un module, récursivement | `service.go` `GetModuleInterfaces` | Garder (notion métier, sert UCREC02), sans création de slot à la lecture (point 1) |
| 11 | Vérification d'intégrité | `POST …/verify` compare l'empreinte du ledger à elle-même : répond toujours `ok` (R-18) | `service.go` `Verify` | Retirer, ou recevoir l'empreinte du fichier du client |
| 12 | Accès invité | `POST /api/identity/guest` ne crée plus rien (RM49) | `handlers_identity.go` | Retirer après la transition de myr-web |
| 13 | Alias `/api/modules` | Redirige vers `/api/components` | `server.go` | Retrait daté : FT-022 |
| 14 | Réseau actif par session | `PUT /api/networks/active` mémorise le choix du client | `handlers_network.go` | Garder : choix explicite du client |
| 15 | Liste en mode dégradé | Ledger injoignable : renvoie les brouillons seuls, avec un avertissement | `service.go` `List` | Garder : signalé explicitement |

**Contrôles absents — le cœur accepte au lieu de refuser :**

| # | Constat | Où | Proposition |
|---|---------|----|-------------|
| A | Une interface peut être ajoutée à un asset inexistant, et une instance référencer un asset inexistant | `service.go` `AddInterface` (`requireDraftCarrier` ne refuse qu'un porteur trouvé et soumis), `AddAssetToWorkspace` | Refuser un asset inconnu |
| B | `POST …/interfaces` accepte un `id` fourni par le client : avec l'identifiant d'une interface existante, il **l'écrase et la rattache à l'asset de l'appelant**, quel que soit son propriétaire | `handlers.go` `handleComponentInterfaces`, `service.go` `AddInterface`, enregistrement par identifiant | Identifiant toujours généré par le serveur |
| C | `PATCH /api/interfaces/{id}` accepte `asset_id` dans le corps : l'interface change de porteur, y compris vers un asset soumis | `handlers.go` `handleInterface`, `service.go` `UpdateInterface` | Porteur non modifiable |

B et C aggravent FT-001 (propriété non contrôlée côté serveur).

## Critères de clôture

- [ ] Stockage des brouillons et des miniatures tranché, ordre de retrait coordonné avec `myr-web`
- [ ] Aucune des quatre fonctions dans `domain/`, `adapters/` ni l'API ; specs alignées
- [ ] Création d'un composant et d'un module complets en une écriture, testée

## Liens

- **Use cases** : [UCAM01 (analyse)](../specs/2-Analyse/UCAM-Assemblage_Module/UCAM01.md) · [UCCE01 (analyse)](../specs/2-Analyse/UCCE-Composant_Ecriture/UCCE01.md) · [UCMOD06 (analyse)](../specs/2-Analyse/UCMOD-Module/UCMOD06.md)
- **Règles métier** : `RM13`, `RM16`, `RM17` (tags `rm/…`)
- **Specs** : [Conception_intro](../specs/3-Conception/Conception_intro.md) (ADR-02, ADR-12) · [API_REST](../specs/3-Conception/API_REST.md)
- **Code** : `domain/model/service.go` · `domain/model/ports.go` · `adapters/out/webimage/fetcher.go`
- **Tickets liés** : [FT-012 — Un module soumis peut être soumis à nouveau sans fork](FT-012-un-module-soumis-peut-etre-soumis-a.md) · [FT-022 — Date de retrait des routes legacy /api/modules](FT-022-date-de-retrait-des-routes-legacy-api.md)

## Historique

- 2026-11-08 — Ticket créé ; décision du PO : slots virtuels, disposition à l'écran, brouillons et miniatures relèvent des outils extérieurs ; vocabulaire de référence et connexions libres restent dans le cœur. Stockage et ordre du retrait à trancher
- 2026-11-09 — Décision du PO (R-15) : `myr-core` ne devine pas. Une connexion en double, auparavant dédoublonnée en silence (un identifiant jamais enregistré était renvoyé), est désormais refusée (`409`, RM50) ; créer une nouvelle interface reste une demande explicite de l'outil. Le PO estime que d'autres fonctions de l'API ne devraient pas exister dans le cœur : dresser l'inventaire des fonctions qui devinent ou facilitent à la place du client, et le trancher avec lui
- 2026-11-09 — Inventaire dressé (§ Inventaire du 2026-11-09) : 15 fonctions et 3 contrôles absents, propositions à trancher par le PO
- 2026-11-09 — Décisions du PO sur l'inventaire : A, B, C corrigés (commit `3c48a30`) ; 5 et 6 (catégorie et canal obligatoires) retenus ; 4, 7 et 11 à retirer ; 12 à retirer après la transition de myr-web ; 8 et 9 gardés — la copie de composition (8) ne fait que remplir un brouillon et partira avec les brouillons (point 3), l'arborescence (9) est une lecture de la filiation, cœur de la traçabilité
- 2026-11-09 — Déployé en production : API et CLI `b5ed57f` depuis `feature/ledger-invariants`, après répétition sur myr-test (42 OK, 0 NOK) et 19 tests d'intégration Fabric OK ; sauvegarde `~/.Myr/server-data.bak-20261109-185157`, anciens binaires conservés en `.prev` (contrôles A à C)
- 2026-11-09 — Points 11, 7, 4, 5 et 6 corrigés sur `feature/ledger-invariants`, non déployés : `verify` retiré en REST et en CLI (`bf2c832`), création par lot retirée (`10847b5`), plus de téléchargement d'`og:image` à la création (`08ee6bb`), catégorie et canal obligatoires, `Service.Add` retiré (`51f12f3`). Changements cassants pour myr-web : mise en production à coordonner. `FileStoragePort` n'a plus d'utilisateur dans le domaine : retrait à trancher. Point 12 après la transition de myr-web
- 2026-11-09 — Point 12 fait (commit `89fa3ca`) : `POST /api/identity/guest` retiré, sur décision du PO, sans attendre myr-web. `FileStoragePort` retiré (commit `e5e6b10`)
