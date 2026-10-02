---
tags:
  - couche/conception
  - type/conception
  - rm/RM01
  - rm/RM08
  - rm/RM14
  - rm/RM17
  - rm/RM19
  - rm/RM23
  - rm/RM24
  - rm/RM30
  - rm/RM31
  - rm/RM32
  - rm/RM33
  - rm/RM39
  - rm/RM40
  - rm/RM41
  - rm/RM42
  - enf/ENF18
  - enf/ENF25
  - relecture/incoherence
  - relecture/remarque
---
# API REST — Contrat d'interface

> Phase 3 — Arrington | Implémentation : `adapters/in/rest/`

---

## 1. Principes

- **Base URL :** `http(s)://<host>:<port>/api/`
- **Format :** JSON — `Content-Type: application/json`
- **Auth :** token opaque de session — en-tête `X-Myr-Token: <token>` ( voir `specs/3-Conception/DC_D1_Auth_Identity.md`)
- **Erreurs :** `{ "error": "<code>", "message": "<description>" }`
- **Succès :** objet ou tableau JSON directement (pas d'enveloppe `{ data: ... }`)

### Niveaux d'accès

| Niveau | Middleware | Condition |
|--------|-----------|-----------|
| Public | aucun | Toujours accessible |
| Auth | `requireAuth` | Token de session valide requis (`X-Myr-Token`) |
| Write | `requireRole(rbac.PermWrite)` | Session valide + permission `write` (rôle `contributor` par défaut) |
| Admin | `requireRole(rbac.PermAdmin)` | Session valide + permission `admin` uniquement |

> Pour les routes mixtes (lecture libre, écriture protégée), le middleware vérifie la permission uniquement sur `POST`, `PUT`, `PATCH`, `DELETE`.

**Codes d'erreur identité/session :**

| Code HTTP | Cas |
|-----------|-----|
| 400 | Corps invalide, champs requis manquants |
| 401 | Secret d'enrôlement CA invalide, token de session absent/expiré |
| 403 | Permission insuffisante, accès invité refusé (réseau privé) |
| 429 | Trop de tentatives (`authLimiter`, 10/min/IP) |

---

## 2. D1 — Identité & RBAC

| Méthode | Route | Accès | Description |
|---------|-------|-------|-------------|
| GET | `/api/identity/policy` | Public | Politique du réseau actif (AllowAutoGuest, AllowAutoRegister) |
| POST | `/api/identity/enroll` | Public | Enrôlement Fabric CA avec secret d'enrollment (sans créer de session REST) |
| POST | `/api/identity/session` | Public | (Ré-)enrôlement + création d'une session REST (token opaque) |
| POST | `/api/identity/guest` | Public | Token `reader` automatique (si `AllowAutoGuest=true`) |
| POST | `/api/identity/request` | Public | Demande d'accès (auto-enregistrement si `AllowAutoRegister=true`, sinon mise en attente `pending` jusqu'à `POST /api/identity/requests/{id}/approve`) |
| GET | `/api/identity/requests` | Admin | Liste des demandes de compte en attente — expose des données personnelles (email, message libre), réservée à `identity.admin` |
| POST | `/api/identity/requests/{id}/approve` | Admin | Approuver une demande `pending` — transforme la demande en identité CA active (UCA01) |
| GET | `/api/identity/wallets` | Auth | Wallets locaux — le répertoire de wallets est partagé par tous les utilisateurs enrôlés sur le nœud : une session sans `identity.admin` ne voit que son propre wallet (filtré par pseudo), jamais celui des autres |
| GET | `/api/identity/status` | Auth | Statut CA d'un wallet (`?handle=pseudo@org`) |

> Pas de CRUD REST pour le RBAC (`domain/role`) — gestion des rôles CLI uniquement (`myr role ...`). Le REST ne fait que consulter les permissions via `requireRole`.

---

## 3. D2 — Administration

| Méthode | Route | Accès | Description |
|---------|-------|-------|-------------|
| GET | `/api/channels` | Public (lecture) | Liste les canaux Fabric disponibles |
| PUT | `/api/channels` | Auth | Mettre à jour la config canal |
| GET | `/api/networks` | Public | Liste des profils réseau configurés |
| GET/PUT | `/api/networks/active` | Auth | Réseau actif de la session |
| GET | `/api/admin/sessions` | Admin | Sessions actives (identifiant permettant de cibler `DELETE /api/admin/sessions/{token}` — voir `specs/roadmap_dev.md` § Écarts Identité & Session, E4 pour le point ouvert sur la forme de cet identifiant) |
| DELETE | `/api/admin/sessions/{token}` | Admin | Invalider une session (token complet requis) |

> **Ajout d'organisation, de nœud (UCADM01/03/04) :** exposés en **CLI uniquement** (`myr org add`, `myr node add/remove`) — l'API REST n'expose jamais ces opérations d'infrastructure (voir `DC_CLI_Admin.md`).

---

## 4. D3/D4/D6 — Composants (ressource unique, ADR-11)

> **Fusion composant/module (`Conception_intro.md` ADR-11) :** `component` est la seule ressource REST exposée pour `Model3D` — qu'il porte un fichier ressource (`Hash`), des instances (`WorkspaceInstances`), ou les deux à la fois. `componentDTO` porte tous les champs (`category`, `hash`, `parent_id`, `tags`, `links`, `instances`, `assemblies`, `versions`, `module_versions`) — vides quand non pertinents, jamais absents. La distinction affichée « composant » / « produit » se dérive de `instances.length > 0`, elle n'est stockée nulle part. Les anciennes routes `/api/modules/*` (§6) sont retirées de ce contrat ; une fenêtre d'alias temporaire vers les routes ci-dessous reste un point ouvert pour le PO (durée et date de coupure à coordonner avec `myr-web`, voir ADR-11).

| Méthode | Route | Accès | Description |
|---------|-------|-------|-------------|
| GET | `/api/components` | Public (DC-D1-06) | Liste de tous les composants, décomposés ou non (filtres: q, categories, owner_id, parent_id, hash, tags, channel, status, limit — `?uses_component={id}` pour les assets qui intègrent un composant donné, `DC_D8_Recherche.md` §4, UCREC04) |
| POST | `/api/components` | Contributor | Créer un composant (`AddFull`), un fichier CAO par appel ; `parent_id` optionnel pour dériver d'un asset existant (composant ou déjà décomposé), avec copie de sa composition le cas échéant (UCMOD01) |
| POST | `/api/components/batch` | Contributor | Créer plusieurs composants en un seul appel, un fichier CAO par composant (UCCE01, voir note ci-dessous) — sans garantie transactionnelle : chaque fichier réussit ou échoue indépendamment |
| GET | `/api/components/{id}` | Auth | Détail d'un composant, décomposé ou non |
| PATCH | `/api/components/{id}` | Contributor | Modifier les métadonnées (patch : name, description, tags, locations, license_id — UCCE02/UCMOD03, mêmes champs quel que soit le nombre d'instances) — #incoherence la méthode documentée ici était PUT, mais le handler (`adapters/in/rest/handlers.go`, `patchComponent`) répond uniquement à PATCH ; corrigé pour refléter le comportement réel, à confirmer que PUT n'était pas l'intention initiale |
| DELETE | `/api/components/{id}` | Contributor | Supprimer (UCCE07/UCMOD08) — brouillon retiré du stockage local ; asset déjà soumis masqué localement (disparaît des listes), son enregistrement blockchain (y compris ses `ModuleVersion` éventuelles) n'est jamais modifié (RM08) |
| GET | `/api/components/{id}/interfaces` | Auth | Interfaces physiques directes si l'asset n'a aucune instance, interfaces exposées (non connectées en interne) s'il en a (UCAM02, UCMOD04) |
| GET | `/api/components/{id}/compatible` | Auth | Composants compatibles (type d'interface, catégorie, tag, sens — `DC_D8_Recherche.md` §2, UCREC02) |
| GET | `/api/components/{id}/versions` | Auth | Arbre de versions / historique de dérivation (`DC_D8_Recherche.md` §3, UCREC03) |
| GET | `/api/components/{id}/instances` | Auth | Instances de l'asset (UCMOD04) — liste vide si l'asset n'a jamais été décomposé |
| POST | `/api/components/{id}/instances` | Contributor | Ajouter un composant existant comme instance (`AddAssetToWorkspace`) — applicable à n'importe quel `id`, y compris un asset qui n'a encore aucune instance ; brouillon local, sans effet blockchain avant soumission (UCMOD01, UCAM05) |
| DELETE | `/api/components/{id}/instances/{instanceId}` | Contributor | Retirer une instance (cascade sur ses connexions, RM14 — UCAM08) |
| GET | `/api/components/{id}/assemblies` | Auth | Connexions internes de l'asset (UCMOD04) — liste vide si l'asset n'a aucune instance |
| GET | `/api/components/{id}/bom` | Auth | Export BOM (Bill of Materials, `DC_D8_Recherche.md` §5, UCREC05) |
| POST | `/api/components/{id}/submit` | Contributor | Soumettre à la blockchain — route unique dispatchant en interne vers `SubmitModule` (applique RM17, crée une `ModuleVersion`) si l'asset a des `Assemblies`, vers `Submit` sinon (voir ADR-11, Conséquence — le point de soumission doit dispatcher en interne) |
| POST | `/api/components/{id}/decompose/preview` | Contributor | Analyser le fichier STEP du composant, retourner une proposition de découpage (sous-pièces, connexions candidates) sous un `decomposition_id` temporaire — aucune entité créée (UCAM09, RM40) |
| POST | `/api/components/{id}/decompose/commit` | Contributor | Matérialiser une proposition de découpage retenue (`decomposition_id`), éventuellement corrigée : sous-composants en draft, module `decoupage`, liaisons compatibles (UCAM09, RM39/RM41) — crée un nouvel asset distinct référençant `{id}` par `parent_id`, ne modifie jamais `{id}` lui-même (post-condition UCAM09, voir ADR-11 pour le point ouvert PO sur cette sémantique d'identité) |

> **Accès visiteur (EF17, UCCL01) :** seule la liste (`GET /api/components`) est concernée par l'accès public — voir `DC_D1_Auth_Identity.md` DC-D1-06. Avant ADR-11, une incohérence d'accès existait entre `GET /api/components` (Public) et `GET /api/modules` (Auth) pour deux endpoints qui ne se distinguaient que par le contenu de l'asset résolu ; la fusion en une seule route supprime cette incohérence de fait — `GET /api/components` est Public quel que soit le nombre d'instances de l'asset.

> **`POST /api/components/batch` (UCCE01) :** multipart, champ répété `files` (un ou plusieurs fichiers CAO). Les autres champs (`owner_id`, `channel_id`, `category`, `parent_id`, `license_id`, `tags`) sont partagés par tous les fichiers du lot. Le nom de chaque composant créé est dérivé du nom de fichier (sans extension) — pas de nom distinct par fichier dans ce contrat. La réponse liste un résultat par fichier (`{filename, component}` en cas de succès, `{filename, error}` en cas d'échec) : aucune atomicité n'est garantie entre les fichiers d'un même lot (pas de base relationnelle, pas de transaction inter-blockchain). Champ répété optionnel `thumbnails` : une valeur (data URL base64, chaîne vide acceptée) par fichier, associée par **position** — le i-ème `thumbnails` correspond au i-ème `files`, jamais par nom de fichier. `myr-core` ne rendant jamais lui-même un fichier 3D (§ Séparation des dépôts), c'est au client (`myr-web`) de rendre chaque fichier et de fournir sa miniature ici ; sans valeur à une position donnée, ce composant reste sans miniature (repli `og:image` uniquement si un `locations`/`links` est déclaré, RM42).

> **`POST /api/components/{id}/decompose/preview` puis `.../commit` (UCAM09) :** deux appels distincts pour ne jamais engager un résultat non validé — voir `specs/1-Expression/UCAM-Assemblage_Module/UCAM09.md` pour le détail du scénario et de la forme des corps de requête/réponse (`parts`, `suggested_connections`, `decomposition_id`). `preview` reste une requête HTTP bloquante bornée par un budget de temps serveur (pas de job asynchrone à interroger) — voir `DC_CLI_Model.md` DC-CLIM-04 et §8 pour le point encore ouvert sur la technologie d'analyse STEP dont dépend le temps de traitement réel.

---

## 5. D5 — Connexions et Interfaces (Composition de Module)

| Méthode | Route | Accès | Description |
|---------|-------|-------|-------------|
| GET | `/api/connections` | Auth | Liste des connexions |
| POST | `/api/connections` | Contributor | Créer une connexion simple |
| DELETE | `/api/connections/{id}` | Contributor | Supprimer une connexion |
| POST | `/api/assembly-links` | Contributor | Créer une liaison interface→interface |
| POST | `/api/virtual-connect` | Contributor | Relier interface virtuelle → physique |
| GET | `/api/interfaces/{id}` | Auth | Détail d'une interface |
| PUT | `/api/interfaces/{id}` | Contributor | Modifier une interface |
| DELETE | `/api/interfaces/{id}` | Contributor | Supprimer une interface |
| GET | `/api/refs` | Auth | Vocabulaire de référence (catégories, types, unités) |
| POST | `/api/refs/categories` | Contributor | Ajouter une catégorie |
| POST | `/api/refs/types` | Contributor | Ajouter un type |
| POST | `/api/refs/units` | Contributor | Ajouter une unité |

---

## 6. D6 — Modules (retiré, fusionné dans §4 — ADR-11)

La famille `/api/modules/*` n'est plus une ressource distincte : toutes les routes qu'elle exposait (liste, détail, métadonnées, instances, soumission, interfaces exposées, assemblages, BOM, filtre `uses_component`) sont désormais servies par `/api/components/*` (§4 ci-dessus), qui répond identiquement qu'un asset ait ou non des instances. La question d'accès visiteur restée ouverte pour `GET /api/modules` (`Auth`) est résolue par la fusion : `GET /api/components` est Public (§4), sans distinction selon le contenu de l'asset.

---

## 7. Licences

| Méthode | Route | Accès | Description |
|---------|-------|-------|-------------|
| GET | `/api/licenses` | Auth | Catalogue des licences |
| GET | `/api/licenses/{id}` | Auth | Détail d'une licence |
| POST | `/api/licenses/check` | Auth | Vérifie la compatibilité `{ parent_license_id, proposed_license_id }` (règle 8) — #remarque route implémentée (`handleLicense`) mais absente de ce tableau jusqu'ici |
| POST | `/api/licenses/check-product` | Auth | Vérifie la compatibilité d'un module assemblant plusieurs composants `{ component_license_ids[], proposed_module_license_id }` (règle 8) — #remarque idem |

---

## 8. Miniatures

| Méthode | Route | Accès | Description |
|---------|-------|-------|-------------|
| POST | `/api/components/{id}/thumbnail` | Contributor | Sauvegarder une miniature STL (data URL base64), quel que soit le nombre d'instances de l'asset — #incoherence route absente du routeur (`server.go`) et du handler (`handleComponent` ne reconnaît que les suffixes `/interfaces`, `/tree` et `/thumbnail/regenerate`) ; la miniature d'un composant se fixe uniquement au moment de sa création via le champ `thumbnail` du formulaire multipart de `POST /api/components` |
| GET | `/api/components/{id}/thumbnail` | Auth | Récupérer la miniature — #incoherence idem, non routée ; la miniature est aujourd'hui exposée via le champ `thumbnail` de `componentDTO` (`GET /api/components`, `GET /api/components/{id}`) |
| POST | `/api/components/{id}/thumbnail/regenerate` | Contributor | Redériver la miniature depuis la seule source durable accessible côté serveur : l'og:image du premier emplacement externe enregistré (`Locations`) — un modèle 3D sans emplacement n'a pas de source régénérable côté serveur (le rendu 3D est produit par le client GUI, pas par `myr`, voir § Séparation des dépôts) ; répond `{ "thumbnail": "<dataURL>" }` (jamais le DTO complet de l'asset) ; échoue en 500 aussi bien pour l'absence d'emplacement externe que pour un échec réseau, un timeout ou une balise `og:image` absente sur la page cible — ces cas ne sont distingués par aucun code ni champ d'erreur dédié, seul le texte du message diffère ; échoue en 503 si l'asset n'est ni un brouillon local ni joignable sur une blockchain configurée — même handler Go `regenerateThumbnail` quel que soit le nombre d'instances de l'asset (ADR-11) — #remarque aucun test (unitaire ou intégration) ne couvre `adapters/out/webimage.Fetcher` (le scraping HTTP réel de la balise `og:image`) ; seuls des tests avec `OGImageFetcher`/service mockés existent (`domain/model/tests/`, `adapters/in/rest/tests/`, `adapters/in/cli/model_test.go`) — aucune preuve qu'un client ait déjà exercé ce chemin contre une vraie page web |

---

## 8bis. Vérification d'intégrité blockchain

Myr ne conserve jamais de copie du fichier ressource transmis — seule son empreinte (`Model3D.Hash`) est retenue à la soumission (RM01). Il n'y a donc rien à relocaliser ni à revérifier dans un stockage propre à Myr ; `verify` ne porte que sur la cohérence des métadonnées inscrites sur la blockchain. Pour la traçabilité des emplacements où un composant ou produit existe réellement, voir §8ter (RM42).

| Méthode | Route | Accès | Description |
|---------|-------|-------|-------------|
| POST | `/api/components/{id}/verify` | Contributor | Vérifie la cohérence blockchain des métadonnées d'un asset déjà soumis (aucune vérification requise pour un brouillon local, jamais encore ancré), quel que soit son nombre d'instances — même handler Go `verifyAsset` (ADR-11). Répond toujours `200` : `{ "ok": true }` en cas de succès, `{ "ok": false, "reason": "blockchain_integrity_failed" }` en cas d'échec. Échoue en `503` uniquement si la blockchain est indisponible et l'asset n'est pas un brouillon local |

---

## 8ter. Traçabilité et alerte des emplacements externes (RM42)

| Méthode | Route | Accès | Description |
|---------|-------|-------|-------------|
| POST | `/api/components/{id}/locations/check` | Contributor | Vérifie l'accessibilité de chaque emplacement externe enregistré (`Locations`) par une requête sur son URL, quel que soit le nombre d'instances de l'asset — même handler Go `checkLocations` (ADR-11). Met à jour le statut et la date de vérification de chaque emplacement (`LocationCheck`, hors blockchain — voir `Modele_Domaine.md` §3) sans jamais modifier `Model3D` lui-même, qu'il soit brouillon ou déjà soumis (RM19 non concerné). Répond toujours `200` avec le détail par emplacement : `{ "locations": [ { "url": "...", "status": "reachable"\|"unreachable", "checked_at": "..." } ] }` ; `{ "locations": [] }` si l'asset ne porte aucun emplacement externe. Ne vérifie que l'accessibilité de l'URL, jamais le contenu qui y est exposé — Myr ne conservant aucune copie du fichier ressource (§8bis), il n'a lui-même rien à comparer à `Model3D.Hash` |

---

## 9. Infrastructure

| Méthode | Route | Accès | Description |
|---------|-------|-------|-------------|
| GET | `/api/ping` | Public | Liveness (sans appel Fabric) |
| GET | `/api/status` | Public | Statut serveur |
| GET | `/api/health` | Public | Health check (load balancer) |
| GET | `/metrics` | Public | Métriques Prometheus |

---

## 10. Headers de sécurité (appliqués à toutes les réponses)

```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Referrer-Policy: strict-origin-when-cross-origin
Content-Security-Policy: default-src 'none'
```

> `myr` sert exclusivement du JSON (pas de GUI embarquée) — voir `specs/3-Conception/Securite.md`.

---

## 11. D7 — Paiement

> Détail des entités (`Order`, `AssetPrice`) et règles métier (RM23, RM24, RM30–RM33) dans `DC_D7_Payment.md`.

| Méthode | Route | Domaine | UC |
|---------|-------|---------|-----|
| POST | `/api/orders` | D7 | UCPI01 |
| GET | `/api/orders/{id}` | D7 | UCPI01 |
| POST | `/api/orders/{id}/deliver` | D7 | UCAUT01 (`DC_D9_Automatisation.md` §2) |
| GET/POST | `/api/assets/{id}/price` | D7 | UCPI04 |
| POST | `/api/assets/{id}/transfer` | D7 | UCPI07 |
| POST | `/api/assets/{id}/clone` | D7 | UCPI08 |

## 12. D9 — Automatisation avancée

> Détail dans `DC_D9_Automatisation.md` §5. Le plugin CAO (UCAUT03) n'a pas de route dédiée — il consomme `POST /api/components` comme tout client API (DC-D9-02).

| Méthode | Route | Domaine | UC |
|---------|-------|---------|-----|
| POST | `/api/components/{id}/versions` | D9 | UCAUT04 |
| GET | `/api/components/{id}/versions/diff` | D9 | UCAUT04 |

---

## 13. Spec OpenAPI générée automatiquement

Les tableaux ci-dessus décrivent le contrat cible de l'API (y compris des routes non encore câblées, ex. §11 D7 Paiement — voir § « CLAUDE.md n'est pas une spec » sur la distinction attendu / avancement). En complément, une **spec OpenAPI 3 est générée automatiquement à partir des commentaires Go** au-dessus de chaque handler REST — elle documente précisément ce que le serveur expose **réellement** à un instant donné (paramètres, schémas de requête/réponse, codes HTTP, exigence d'authentification), pour tout développeur tiers qui construit un logiciel consommant l'API (GUI externe, plugin CAO, script d'intégration).

### Mécanisme

- **Outil :** `swaggo/swag` — parse des commentaires `@Summary`, `@Param`, `@Success`, `@Router`, etc. placés directement au-dessus de chaque fonction handler dans `adapters/in/rest/*.go`, ainsi que le bloc d'annotations générales (`@title`, `@BasePath`, `@securityDefinitions`) au-dessus de `func main()` dans `cmd/api/main.go`.
- **Aucune dépendance ajoutée à `go.mod`/`vendor/` :** le CLI `swag` s'exécute via `go run github.com/swaggo/swag/cmd/swag@<version>` — un outil de développement au même titre que `gofmt`, jamais compilé dans le binaire serveur. La génération produit uniquement `api/swagger.json` et `api/swagger.yaml` (flag `--outputTypes json,yaml` — pas de fichier Go `docs.go`, qui obligerait sinon à importer le package `swaggo/swag` au runtime).
- **Commande :** `make docs-api` (cible du `Makefile`, version de l'outil pinée dans la variable `SWAG_VERSION`).
- **Authentification documentée :** en-tête `X-Myr-Token` (schéma `apiKey`, nommé `MyrToken`) — jamais de JWT (voir règle 23 du domaine identité).

### Convention d'annotation

- Chaque fonction handler qui répond directement à une requête HTTP porte un bloc `// @Summary ... // @Router /chemin [méthode]` juste au-dessus de sa déclaration — sans ligne vide entre le commentaire et le `func`, faute de quoi `swag` ne l'associe pas à la fonction.
- Les corps de requête décrits par une struct Go nommée au niveau paquet (ex. `connectionDTO`, `assemblyLinkRequest`, `walletDTO`) sont référencés par leur type — le schéma JSON exact apparaît dans la spec générée. Les corps décrits par une struct anonyme locale à la fonction (ex. `var body struct{...}` dans `patchComponent`) ne sont pas nommables par `swag` : ils apparaissent dans la spec comme un `object` générique, le détail des champs restant uniquement dans la description textuelle du commentaire — limite connue de l'outil, pas du domaine.
- **Cas des handlers multi-routes** (ex. `handleModule`, qui route une dizaine de sous-ressources d'un module derrière un seul point d'entrée Go) : `swag` associe un commentaire à une fonction, pas à une route — un seul bloc `@Description`/`@Success` couvre alors plusieurs lignes `@Router`, avec une précision par sous-route moindre que pour les handlers dédiés à une seule route. Diviser ces fonctions uniquement pour affiner la documentation n'a pas été fait ici : cela reviendrait à refactorer du code qui fonctionne pour un besoin de documentation, hors du périmètre demandé — un futur découpage de ces handlers (s'il est motivé par autre chose que la doc) affinera la spec générée sans changement de convention.

### Autorité en cas de divergence

En cas d'écart entre les tableaux `§1`–`§12` de ce document et `api/swagger.json`/`api/swagger.yaml` : la spec générée fait foi de l'implémentation réelle (elle est dérivée du code), les tableaux ci-dessus font foi de l'intention cible. Les écarts déjà identifiés lors de la mise en place de la génération automatique sont marqués `#incoherence` / `#remarque` dans ce document (§4, §6, §7, §8) — à trancher par le product owner.

### Régénération

`api/swagger.json`/`.yaml` ne sont pas régénérés automatiquement à chaque modification d'un handler — exécuter `make docs-api` après tout changement de signature de route (nouveau paramètre, nouveau code retour, nouvelle route). Aucune vérification CI ne détecte aujourd'hui une spec désynchronisée du code (`make ci` ne couvre que ENF18/ENF25) ; l'ajout d'un tel contrôle est une évolution possible, non traitée ici faute de validation explicite.

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs](../Carte_des_specs.md)

**Use cases cités**
- UCA01 — Création d'un compte : [expression](../1-Expression/UCA-Compte_et_Acces/UCA01.md) · [analyse](../2-Analyse/UCA-Compte_et_Acces/UCA01.md)
- UCADM01 — Ajouter une organisation au réseau : [expression](../1-Expression/UCADM-Administration/UCADM01.md) · [analyse](../2-Analyse/UCADM-Administration/UCADM01.md)
- UCADM03 — Ajouter un nœud à un réseau existant : [expression](../1-Expression/UCADM-Administration/UCADM03.md) · [analyse](../2-Analyse/UCADM-Administration/UCADM03.md)
- UCADM04 — Retirer un nœud d'un réseau existant : [expression](../1-Expression/UCADM-Administration/UCADM04.md) · [analyse](../2-Analyse/UCADM-Administration/UCADM04.md)
- UCAM02 — Visualiser les interfaces physiques de composants : [expression](../1-Expression/UCAM-Assemblage_Module/UCAM02.md) · [analyse](../2-Analyse/UCAM-Assemblage_Module/UCAM02.md)
- UCAM05 — Transformation d'un composant en module : [expression](../1-Expression/UCAM-Assemblage_Module/UCAM05.md) · [analyse](../2-Analyse/UCAM-Assemblage_Module/UCAM05.md)
- UCAM08 — Retirer une instance de composant d'un Module : [expression](../1-Expression/UCAM-Assemblage_Module/UCAM08.md) · [analyse](../2-Analyse/UCAM-Assemblage_Module/UCAM08.md)
- UCAM09 — Décomposition assistée d'un composant assemblage : [expression](../1-Expression/UCAM-Assemblage_Module/UCAM09.md)
- UCAUT01 — Fabrication/Livraison d'un Composant : [expression](../1-Expression/UCAUT-Automatisation/UCAUT01.md) · [analyse](../2-Analyse/UCAUT-Automatisation/UCAUT01.md)
- UCAUT03 — Ajouter un modèle 3D depuis un logiciel CAO : [expression](../1-Expression/UCAUT-Automatisation/UCAUT03.md) · [analyse](../2-Analyse/UCAUT-Automatisation/UCAUT03.md)
- UCAUT04 — Gestion SCM d'un modèle 3D : [expression](../1-Expression/UCAUT-Automatisation/UCAUT04.md) · [analyse](../2-Analyse/UCAUT-Automatisation/UCAUT04.md)
- UCCE01 — Ajout d'un composant Physique : [expression](../1-Expression/UCCE-Composant_Ecriture/UCCE01.md) · [analyse](../2-Analyse/UCCE-Composant_Ecriture/UCCE01.md)
- UCCE02 — Configurer un Composant : [expression](../1-Expression/UCCE-Composant_Ecriture/UCCE02.md) · [analyse](../2-Analyse/UCCE-Composant_Ecriture/UCCE02.md)
- UCCE03 — Ajout d'un composant Numérique : [expression](../1-Expression/UCCE-Composant_Ecriture/UCCE03.md) · [analyse](../2-Analyse/UCCE-Composant_Ecriture/UCCE03.md)
- UCCE07 — Supprimer un Composant : [expression](../1-Expression/UCCE-Composant_Ecriture/UCCE07.md) · [analyse](../2-Analyse/UCCE-Composant_Ecriture/UCCE07.md)
- UCCL01 — Faire une recherche par filtre : [expression](../1-Expression/UCCL-Composant_Lecture/UCCL01.md) · [analyse](../2-Analyse/UCCL-Composant_Lecture/UCCL01.md)
- UCMOD01 — Créer un Module : [expression](../1-Expression/UCMOD-Module/UCMOD01.md) · [analyse](../2-Analyse/UCMOD-Module/UCMOD01.md)
- UCMOD03 — Modifier les métadonnées d'un Module : [expression](../1-Expression/UCMOD-Module/UCMOD03.md) · [analyse](../2-Analyse/UCMOD-Module/UCMOD03.md)
- UCMOD04 — Visualiser les composants d'un Module : [expression](../1-Expression/UCMOD-Module/UCMOD04.md) · [analyse](../2-Analyse/UCMOD-Module/UCMOD04.md)
- UCMOD08 — Supprimer un Module : [expression](../1-Expression/UCMOD-Module/UCMOD08.md) · [analyse](../2-Analyse/UCMOD-Module/UCMOD08.md)
- UCPI01 — Commander un Module complet : [expression](../1-Expression/UCPI-Propriete_Intellectuelle/UCPI01.md) · [analyse](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI01.md)
- UCPI04 — Définir un prix sur un Composant proprietaire : [expression](../1-Expression/UCPI-Propriete_Intellectuelle/UCPI04.md) · [analyse](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI04.md)
- UCPI07 — Transfert de propriété intellectuelle : [expression](../1-Expression/UCPI-Propriete_Intellectuelle/UCPI07.md) · [analyse](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI07.md)
- UCPI08 — Cloner un Composant sur un réseau exterieur : [expression](../1-Expression/UCPI-Propriete_Intellectuelle/UCPI08.md) · [analyse](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI08.md)
- UCREC02 — Rechercher les Composants compatibles : [expression](../1-Expression/UCREC-Recherche/UCREC02.md) · [analyse](../2-Analyse/UCREC-Recherche/UCREC02.md)
- UCREC03 — Rechercher les versions des Composants : [expression](../1-Expression/UCREC-Recherche/UCREC03.md) · [analyse](../2-Analyse/UCREC-Recherche/UCREC03.md)
- UCREC04 — Rechercher les Modules qui utilisent un Composant : [expression](../1-Expression/UCREC-Recherche/UCREC04.md) · [analyse](../2-Analyse/UCREC-Recherche/UCREC04.md)
- UCREC05 — Exporter BOM Module : [expression](../1-Expression/UCREC-Recherche/UCREC05.md) · [analyse](../2-Analyse/UCREC-Recherche/UCREC05.md)

**Règles métier**
- [RM01 — Anti-plagiat obligatoire](../1-Expression/Regles_Metier.md#1.%20Assets%20et%20composants)
- [RM08 — Masquage local, ledger jamais modifié](../1-Expression/Regles_Metier.md#2.%20Blockchain%20et%20immuabilité)
- [RM14 — Suppression en cascade des connexions](../1-Expression/Regles_Metier.md#4.%20Composition%20d'un%20Module%20%28instances%29)
- [RM17 — Assemblage requis pour soumission (module uniquement)](../1-Expression/Regles_Metier.md#5.%20Modules%20et%20cycle%20de%20vie%20des%20assets)
- [RM19 — Fork d'un asset soumis](../1-Expression/Regles_Metier.md#5.%20Modules%20et%20cycle%20de%20vie%20des%20assets)
- [RM23 — Distribution automatique des commissions](../1-Expression/Regles_Metier.md#7.%20Propriété%20intellectuelle%20et%20commissions)
- [RM24 — Répartition proportionnelle multi-auteurs](../1-Expression/Regles_Metier.md#7.%20Propriété%20intellectuelle%20et%20commissions)
- [RM30 — Calcul automatique du prix d'un module](../1-Expression/Regles_Metier.md#7.%20Propriété%20intellectuelle%20et%20commissions)
- [RM31 — Modification de prix — effet sur les commandes futures uniquement](../1-Expression/Regles_Metier.md#7.%20Propriété%20intellectuelle%20et%20commissions)
- [RM32 — Asset à prix nul — librement disponible](../1-Expression/Regles_Metier.md#7.%20Propriété%20intellectuelle%20et%20commissions)
- [RM33 — Devise unique par réseau](../1-Expression/Regles_Metier.md#7.%20Propriété%20intellectuelle%20et%20commissions)
- [RM39 — Filiation d'un découpage](../1-Expression/Regles_Metier.md#5.%20Modules%20et%20cycle%20de%20vie%20des%20assets)
- [RM40 — Proposition de découpage non engageante](../1-Expression/Regles_Metier.md#5.%20Modules%20et%20cycle%20de%20vie%20des%20assets)
- [RM41 — Compatibilité toujours vérifiée pour une connexion suggérée](../1-Expression/Regles_Metier.md#5.%20Modules%20et%20cycle%20de%20vie%20des%20assets)
- [RM42 — Traçabilité et alerte des emplacements externes](../1-Expression/Regles_Metier.md#1.%20Assets%20et%20composants)

**Exigences non fonctionnelles**
- [ENF18 — Isolation du domaine métier](../1-Expression/Exigences_Non_Fonctionnelles.md#Tableau%20des%20exigences%20non-fonctionnelles)
- [ENF25 — Licence du code source](../1-Expression/Exigences_Non_Fonctionnelles.md#Tableau%20des%20exigences%20non-fonctionnelles)

**Documents cités**
- [Conception_intro](Conception_intro.md)
- [DC_CLI_Admin](DC_CLI_Admin.md)
- [DC_CLI_Model](DC_CLI_Model.md)
- [DC_D1_Auth_Identity](DC_D1_Auth_Identity.md)
- [DC_D7_Payment](DC_D7_Payment.md)
- [DC_D8_Recherche](DC_D8_Recherche.md)
- [DC_D9_Automatisation](DC_D9_Automatisation.md)
- [Modele_Domaine](Modele_Domaine.md)
- [Securite](Securite.md)
- [roadmap_dev](../roadmap_dev.md)

**Cité par**
- [Conception_intro](Conception_intro.md)
- [DC_D1_Auth_Identity](DC_D1_Auth_Identity.md)

<!-- liens-obsidian:end -->
