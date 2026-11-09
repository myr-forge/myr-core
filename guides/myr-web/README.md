# Transition myr-web — nouveau fonctionnement de l'API myr-core

Ce guide s'adresse aux développeurs de myr-web, et de tout autre client de l'API REST de myr-core. Il décrit ce qui change pour un client avec les évolutions de novembre 2026 : invariants du ledger (ADR-12), immuabilité des assets soumis, rôles, lecture publique et bornes de l'API. Pour chaque changement, il indique les routes à utiliser.

- **Version de l'API** : `079a5a9` et suivantes, branche `feature/ledger-invariants` (production : `http://89.167.102.193/api/`).
- **Contrat détaillé** : `api/swagger.json` / `api/swagger.yaml` (OpenAPI, généré depuis le code).
- **Règles métier citées** : `specs/1-Expression/Regles_Metier.md`.

---

## 1. En bref

| Domaine | Avant | Maintenant | À faire dans myr-web |
|---------|-------|------------|----------------------|
| Connexion | Rôle de session figé à `contributor` | Rôle et statut lus dans le certificat de l'identité | Afficher le rôle reçu ; gérer `403` pour une identité en attente ou suspendue |
| Visiteur sans compte | `POST /api/identity/guest` délivrait un jeton `reader` de 7 jours | Aucune session : le catalogue publié se lit **sans jeton** si l'administrateur l'a ouvert | Lire `allow_auto_guest` dans `/api/identity/policy`, puis appeler sans `X-Myr-Token` |
| Création d'un asset | — | Toujours un **brouillon**, propriété de l'identité de la session | Ne plus envoyer `owner_id` ; enchaîner sur `/submit` |
| Asset soumis | Modifiable, licence comprise | **Figé** : seuls nom, description, tags et liens changent, par le propriétaire | Désactiver les autres éditions ; proposer « créer une nouvelle version » |
| Nouvelle version | — | Asset dérivé (`parent_id`) ; un module dérivé reprend la composition du parent | Bouton « nouvelle version » → `POST /api/components` avec `parent_id` |
| Licence | Libre | Doit exister au catalogue et être compatible avec celle du parent | Proposer la liste de `/api/licenses` ; gérer `400` |
| Slots virtuels | Matérialisables partout | Refusé sur un asset soumis (`409`) | Relier un asset soumis par ses interfaces physiques |
| Liaisons | Toujours acceptées ; un doublon recevait un identifiant fantôme | Refusées si incompatibles (`400`), si l'interface est déjà engagée sur son instance ou si la liaison existe déjà (`409`) | Désigner les instances ; relire `/api/connections` avant d'en créer ; ne plus réessayer « pour voir » |
| Composition d'un module | Toute liaison, même inconnue, se rattachait | Seule une liaison existante entre deux instances du module se rattache ; un module ne se contient pas lui-même (`400`) | Créer les instances, puis la liaison entre elles, puis la rattacher |
| Rôles | Changement par l'administrateur seulement | Demande de rôle en libre-service, automatique ou sur approbation | Écran « demander un rôle » ; écran d'approbation pour l'administrateur |
| Bornes | Aucune | Taille des corps, débit, délais | Gérer `413` et `429` |

---

## 2. Accès et sessions

### 2.1 Politique du réseau — toujours la première requête

```
GET /api/identity/policy
→ 200 {"mode":"blockchain","allow_auto_guest":false,"allow_auto_register":true,"require_approval":false,"network_name":"diy-network"}
```

| Champ | Sens pour le client |
|-------|---------------------|
| `allow_auto_guest` | Le catalogue publié se lit sans compte (§ 2.3) |
| `allow_auto_register` | Un visiteur peut demander un compte (§ 2.4) ; sinon seul l'administrateur en crée |
| `require_approval` | Un compte demandé reste en attente jusqu'à l'approbation de l'administrateur |

### 2.2 Connexion d'un utilisateur

```
POST /api/identity/session
{"name": "alice", "secret": "<secret d'enrôlement>", "org_id": "Org1MSP", "channel": "sandbox"}   // channel optionnel
→ 200 {"token": "…", "role": "contributor", "pseudo": "alice", "channel": "sandbox"}
```

- Le jeton s'envoie ensuite dans l'en-tête `X-Myr-Token` ; il est valable 7 jours.
- Le **rôle et le statut viennent du certificat** émis par la CA, jamais du client (RM45). Tant que l'identité est en attente (`pending`) ou suspendue, la connexion répond **`403`** avec un message explicite : « compte en attente d'approbation par l'administrateur », « compte suspendu »… À afficher tel quel.
- Un changement de rôle (approbation, demande de rôle) s'applique **à la connexion suivante**. Pour en tenir compte, il suffit de se reconnecter avec le même secret.
- Toute écriture sur le ledger est signée par le wallet de l'utilisateur : une session doit donc être liée à une identité. Sans wallet sur le serveur, une écriture répond `403` (« reconnectez-vous avec votre identité »).

### 2.3 Visiteur sans compte (lecture publique) — changement de comportement

Quand l'administrateur a ouvert la lecture publique (`allow_auto_guest: true`, réglage `myr network update --auto-guest`, fermé par défaut), un visiteur **lit le catalogue publié sans jeton ni session** (RM49) :

| Route (`GET`) | Visiteur |
|---------------|----------|
| `/api/components`, `/api/components/{id}` et ses sous-ressources (`/interfaces`, `/tree`, `/instances`, `/assemblies`, `/thumbnail`) | Assets **soumis** uniquement ; un brouillon répond `404` |
| `/api/modules…` (alias déprécié) | Idem |
| `/api/refs`, `/api/licenses`, `/api/licenses/{id}` | Oui |
| Toute autre route, et toute écriture | `401` |

Ce que myr-web doit faire :

1. Ne plus stocker de « jeton invité ». Un visiteur n'envoie **pas** d'en-tête `X-Myr-Token`, ou l'envoie vide.
2. Un jeton présent mais invalide ou expiré est refusé (`401`), même lecture publique ouverte. Dans ce cas, supprimer le jeton local, puis relire sans jeton ou se reconnecter.
3. `POST /api/identity/guest` est **déprécié**. Il répond encore `200 {"token":"","role":"reader","guest":true,…}` si la lecture publique est ouverte, `403` sinon, et ne crée plus de session. Préférer `allow_auto_guest` de la politique.

Lecture publique fermée (le cas en production aujourd'hui) : toute lecture exige une session.

### 2.4 Demander un compte

```
POST /api/identity/request
{"pseudo": "bob", "email": "bob@exemple.com", "org_id": "Org1MSP", "display_name": "Bob", "message": "…"}
→ 201 {"id": "req-…", "status": "approved" | "pending", "pseudo": "bob", "org_id": "Org1MSP", "secret": "<secret d'enrôlement>"}
```

- `403` si le réseau n'accepte pas les inscriptions (`allow_auto_register: false`).
- Le secret n'est renvoyé qu'une fois : l'utilisateur s'en sert pour se connecter (§ 2.2).
- Avec `status: "pending"`, la connexion répond `403` jusqu'à l'approbation de l'administrateur.

### 2.5 Demander un rôle (nouveau)

```
POST /api/identity/roles/request          (session requise)
{"role": "contributor"}
→ 200 {…, "kind": "role", "role": "contributor", "status": "approved"}   // rôle en attribution automatique : attribué
→ 202 {…, "kind": "role", "role": "contributor", "status": "pending"}    // en attente de l'administrateur
```

| Code | Cas |
|------|-----|
| `400` | Rôle déjà détenu, ou absent de la table des rôles |
| `409` | Une autre demande de rôle de l'utilisateur est déjà en attente |
| `403` | Identité inactive |

Le nouveau rôle prend effet à la connexion suivante (§ 2.2). Un rôle d'administration n'est jamais attribué automatiquement (RM47).

> Il n'existe pas encore de route REST listant les rôles disponibles : la table se gère en CLI (`myr role list`). Si l'écran de demande en a besoin, demander l'ajout de `GET /api/roles` côté myr-core.

### 2.6 Écrans d'administration (permission `identity.admin`)

| Route | Usage |
|-------|-------|
| `GET /api/identity/requests` | Demandes d'inscription **et** de rôle : `kind` vaut `""` (inscription) ou `"role"` ; `role` porte le rôle demandé |
| `POST /api/identity/requests/{id}/approve` | Inscription : l'identité devient active. Demande de rôle : le rôle est attribué (`400` s'il a été retiré de la table entre-temps) |
| `POST /api/identity/requests/{id}/reject` | Inscription : l'identité est suspendue. Demande de rôle : l'identité garde son rôle. `409` si la demande est déjà décidée |
| `POST /api/identity/accounts` | Créer un compte actif : `{"pseudo","org_id","role","email"}` → `201 {"handle","secret"}` |

En production, le seul administrateur est le compte d'exploitation `ops@Org1`.

---

## 3. Cycle de vie d'un asset

### 3.1 Créer, c'est créer un brouillon

```
POST /api/components                       (multipart/form-data, permission write)
name, description, category, parent_id, license_id, tags ("a,b"), links (tableau JSON en chaîne), file, thumbnail (data URL)
→ 201 componentDTO  {"id": "…", "status": "draft", "owner_id": "alice@Org1", …}
```

- Toute création produit un **brouillon** (`status: "draft"`), local au serveur, sans transaction blockchain.
- Le propriétaire est **l'identité de la session** (`pseudo@org`). Le champ `owner_id` envoyé est ignoré : ne plus l'envoyer.
- Le fichier CAO n'est reçu que pour calculer son empreinte SHA-256 ; il n'est jamais conservé (RM42). Taille maximale : 256 Mo (§ 5).
- Une empreinte déjà enregistrée sur le ledger sous un autre asset est refusée par le ledger à la soumission : un même fichier n'est publié qu'une fois.

### 3.2 Soumettre

```
POST /api/components/{id}/submit           {"note": "v1"}   // note optionnelle, versions de module
→ 200 componentDTO {"status": "submitted", …}
```

| Code | Cas |
|------|-----|
| `400` | Module sans aucune liaison (RM17) |
| `403` | Session sans identité signataire |
| `409` | Déjà soumis : un asset n'est jamais re-soumis (RM19) |
| `503` | Blockchain injoignable |

### 3.3 Un asset soumis est figé

| Opération sur un asset `submitted` | Résultat |
|------------------------------------|----------|
| `PATCH /api/components/{id}` avec `name`, `description`, `tags`, `links` | `200` — mis à jour sur le ledger, **par le propriétaire seulement** |
| `PATCH` avec un `license_id` différent | `409` — la licence est figée |
| `POST /submit` | `409` |
| Instances : `POST /api/components/{id}/instances`, `PATCH`/`DELETE …/instances/{instanceId}` | `409` |
| Liaisons internes : `POST …/assemblies`, `DELETE …/assemblies/{connID}` | `409` |
| Interfaces : `POST …/interfaces`, `PATCH`/`DELETE /api/interfaces/{id}` | `409` |
| Slot virtuel : `POST /api/virtual-connect` | `409` (§ 4) |
| `DELETE /api/components/{id}` | Masqué localement : il disparaît de la liste mais reste consultable par son id (le ledger ne supprime rien) |

> **Non-propriétaire** : la modification des métadonnées d'un asset d'autrui est refusée par le ledger, mais l'API répond aujourd'hui `500` au lieu de `403` (ticket FT-001). Il faut donc masquer l'édition quand `owner_id` diffère de l'identité de la session.

Sur un **brouillon**, tous ces champs se modifient librement ; la licence proposée reste contrôlée (§ 3.5).

### 3.4 Nouvelle version : créer un asset dérivé

Un asset soumis n'évolue que par un **asset dérivé** :

```
POST /api/components                       (multipart)
name=Drone v2, parent_id=<id de l'asset soumis>, category=amelioration, license_id=…
→ 201 componentDTO {"status": "draft", "parent_id": "<id>", "instances": […], "assemblies": […]}
```

- Si le parent est un **module** (asset décomposé), le dérivé reprend une copie indépendante de sa composition : instances et liaisons internes, avec de nouveaux identifiants. Le module d'origine n'est jamais modifié.
- Le dérivé se modifie librement, puis se soumet (§ 3.2) : c'est la nouvelle version.
- L'arbre des versions s'affiche via `GET /api/components/{id}/tree`.

### 3.5 Licences

- `license_id` doit exister dans `GET /api/licenses` ; sinon `400` (« licence absente du catalogue »).
- Pour un asset dérivé, la licence doit être compatible avec celle du parent (RM03), sinon `400`. `POST /api/licenses/check` permet de le vérifier avant d'envoyer (session avec permission `write`).

---

## 4. Composition, interfaces et slots virtuels

- La composition d'un module (instances, liaisons, interfaces) ne se modifie que tant qu'il est en brouillon (§ 3.3).
- `POST /api/virtual-connect` matérialise le slot virtuel d'un asset en interface physique : c'est refusé (`409`) si cet asset est soumis, car ce serait ajouter une interface à un asset figé (RM19). Pour relier un composant soumis :
  - créer la liaison entre ses **interfaces physiques existantes** : `POST /api/assembly-links {"from_iface_id","to_iface_id","from_instance_id","to_instance_id","label","fastener_asset_id"}` ;
  - ou créer un **dérivé** du composant (§ 3.4) et lui ajouter l'interface.
- **Si `virtual-connect` est refusé** (`400` ou `409`), le slot reste virtuel : rien n'a été modifié, il n'y a rien à annuler côté client.
- **À anticiper** : slots virtuels, positions x/y des instances, miniatures, voire brouillons, sont des facilités d'interface appelées à relever de myr-web plutôt que du cœur (ticket FT-040, à trancher). Mieux vaut les traiter dès maintenant comme des concepts d'affichage côté client, sans construire de logique métier dessus.

### 4.1 Créer une liaison

Toujours dans cet ordre, sur un module en brouillon :

1. Ajouter les instances : `POST /api/components/{id}/instances {"asset_id"}` — la réponse contient `instances`, dont l'identifiant de la nouvelle instance (la dernière).
2. Créer la liaison entre deux de ces instances : `POST /api/assembly-links` avec `from_instance_id` et `to_instance_id` renseignés.
3. La rattacher au module : `POST /api/components/{id}/assemblies {"connection_id"}`.

| Refus | Code | Cause | Que faire |
|-------|------|-------|-----------|
| interfaces incompatibles | `400` | Catégorie, type, sens ou valeurs ne correspondent pas (RM10/RM11) | Montrer la cause ; proposer un asset d'accroche (`fastener_asset_id`) compatible avec les deux |
| interface reliée à elle-même | `400` | Même interface sur la même instance | — |
| interface déjà engagée | `409` | Cette interface, **sur cette instance**, sert déjà dans une liaison (RM09) ; le message donne l'identifiant de la liaison existante | Proposer une autre instance ou une autre interface, ou supprimer l'ancienne liaison |
| liaison déjà existante | `409` | `POST /api/connections` : une connexion directe relie déjà le même `from` au même `to` (RM50) | La retrouver dans `GET /api/connections` — myr-core ne la renvoie pas à votre place et n'en crée pas de copie |
| liaison étrangère au module | `400` | `POST …/assemblies` : la connexion n'existe pas ou ne relie pas deux instances du module (RM51) | Créer d'abord les instances et la liaison entre elles |
| composition cyclique | `400` | `POST …/instances` : l'asset ajouté est le module ou le contient (RM52) | — |
| asset introuvable | `404` | `POST …/instances` ou `POST …/interfaces` : l'asset n'existe pas | Ne référencer que des assets existants |

Les interfaces : l'identifiant est attribué par le serveur à la création (`POST …/interfaces`, un `id` envoyé est ignoré), et une interface ne change jamais d'asset porteur (`PATCH /api/interfaces/{id}` avec un autre `asset_id` : `400`).

La même interface d'un composant reste libre dans ses **autres** instances : deux vis identiques se relient chacune à leur propre écrou. Une liaison créée sans `from_instance_id`/`to_instance_id` engage l'interface elle-même, quelle que soit l'instance : à éviter dans un module.

### 4.2 Valeurs des interfaces

La compatibilité compare les valeurs (`value_min`, `value_max`, `is_range`, `unit`) ainsi :

- une interface **sans unité** n'impose aucune contrainte de valeur ;
- avec une unité, une **valeur simple** (`is_range: false`) vaut l'intervalle [v ; v] : une sortie 5 V est compatible avec une entrée [3 ; 6] V ;
- **0 est une vraie valeur** : une sortie 0 V n'est pas compatible avec une entrée [3 ; 6] V. Pour « pas de valeur », envoyer une interface sans unité.

---

## 5. Erreurs et bornes

### 5.1 Format des erreurs

```
{"error": "message lisible"}
```

Une erreur interne (`500`) porte une référence : `{"error": "erreur interne (réf. a1b2c3d4)", "ref": "a1b2c3d4"}`, plus l'en-tête `X-Myr-Error-Ref`. Afficher la référence à l'utilisateur : elle permet de retrouver le détail dans le journal du serveur.

| Code | Sens |
|------|------|
| `400` | Requête invalide, licence refusée, rôle inconnu, liaison ou composition refusée (§ 4.1) |
| `401` | Pas de session, ou jeton invalide ou expiré |
| `403` | Droits insuffisants, identité inactive, écriture sans identité signataire |
| `404` | Introuvable — y compris un brouillon demandé par un visiteur, ou un asset inconnu référencé par une instance ou une interface |
| `409` | Asset soumis (figé), demande déjà décidée ou déjà en attente, interface déjà engagée, liaison déjà existante (§ 4.1) |
| `413` | Corps trop volumineux (§ 5.2) |
| `429` | Trop de tentatives (§ 5.2) — « réessayez dans une minute » |
| `503` | Blockchain injoignable |

### 5.2 Bornes

| Borne | Valeur |
|-------|--------|
| Corps JSON | 1 Mo |
| Miniature (`POST /api/components/{id}/thumbnail`) | image de 2 Mo |
| Envoi de fichier (multipart) | 256 Mo par défaut (réglable côté serveur, `MYR_MAX_UPLOAD_MB`) |
| Tentatives sur `/api/identity/session`, `/enroll`, `/request`, `/guest` | 10 par minute et par adresse du client |
| Durée d'une requête | jusqu'à 15 minutes (envoi lent d'un gros fichier) |

Les écritures sur le référentiel d'interfaces (`POST /api/refs/categories|types|units`) demandent maintenant la permission `write`, comme toute écriture du catalogue.

---

## 6. Déploiement de myr-web

**myr-web et myr-core sur le même serveur** (cas actuel) : le proxy inverse relaie `/api/` vers `127.0.0.1:8765`, ce qui ne demande aucun réglage côté myr-core. Pour nginx :

```nginx
client_max_body_size 256m;                                   # au moins MYR_MAX_UPLOAD_MB
location /api/ {
    proxy_pass http://127.0.0.1:8765;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;   # indispensable pour la limite de débit
    proxy_read_timeout 15m;
    proxy_send_timeout 15m;
}
```

**myr-web sur un autre serveur** : l'API ne croit l'en-tête `X-Forwarded-For` que d'un proxy de confiance. L'administrateur de myr-core doit donc ajouter l'adresse du proxy de myr-web à `MYR_TRUSTED_PROXIES` (par exemple `MYR_TRUSTED_PROXIES=127.0.0.1,::1,203.0.113.25`). Sans cela, tous les utilisateurs de myr-web partagent la limite de débit de l'adresse du proxy : 10 connexions par minute pour tout le monde. L'API ne fait pas TLS elle-même ; entre deux serveurs, la placer derrière un proxy TLS.

---

## 7. Routes dépréciées

| Route | Remplacement |
|-------|--------------|
| `/api/modules`, `/api/modules/{id}…` | `/api/components…` (ADR-11 : un module est un composant décomposé) |
| `POST /api/identity/guest` | `GET /api/identity/policy` (`allow_auto_guest`) puis lecture sans jeton |

---

## 8. Liste de contrôle de la migration

- [ ] Lire `/api/identity/policy` au démarrage et adapter l'écran d'accueil (visiteur, inscription, connexion)
- [ ] Visiteur : ne plus demander ni stocker de jeton invité ; appeler le catalogue sans `X-Myr-Token`
- [ ] Gérer `401` (jeton expiré → l'effacer), `403` (compte en attente ou suspendu : afficher le message), `429`, `413`
- [ ] Afficher le rôle reçu à la connexion ; proposer une reconnexion après une approbation ou une demande de rôle
- [ ] Création : ne plus envoyer `owner_id` ; afficher l'état brouillon puis proposer la soumission
- [ ] Asset soumis : limiter l'édition à nom, description, tags et liens, pour le propriétaire seulement ; licence et composition en lecture seule
- [ ] Action « nouvelle version » → `POST /api/components` avec `parent_id`
- [ ] Choix de licence depuis `/api/licenses`, avec vérification de compatibilité pour un dérivé
- [ ] Slots virtuels : ne plus les proposer sur un asset soumis ; relier par les interfaces physiques
- [ ] Liaisons : toujours renseigner `from_instance_id`/`to_instance_id` ; gérer `400` (incompatible, réflexive) et `409` (interface engagée, liaison existante) sans réessai automatique
- [ ] Composition : instances, puis liaison entre elles, puis rattachement ; gérer `400` (liaison étrangère, composition cyclique)
- [ ] Valeurs d'interface : envoyer une unité dès qu'une valeur compte ; ne plus utiliser 0 pour « pas de valeur »
- [ ] Écran « demander un rôle » (`/api/identity/roles/request`) et, pour l'administrateur, l'écran des demandes (`kind`, `role`)
- [ ] Afficher la référence des erreurs `500`
- [ ] Remplacer les appels `/api/modules…` par `/api/components…`
