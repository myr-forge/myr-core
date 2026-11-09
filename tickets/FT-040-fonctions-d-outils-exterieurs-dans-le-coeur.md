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
