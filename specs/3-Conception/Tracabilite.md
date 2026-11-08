---
tags:
  - couche/conception
  - type/conception
---
# Traçabilité specs ↔ code ↔ tests

> Phase 3 — Conception | Outil : `cmd/tracegen` (`make docs-trace`)

## 1. Objectif

Relier, sans saisie manuelle, chaque exigence des specs au code qui la réalise et aux tests qui la vérifient, et rendre ces liens navigables dans Obsidian :

- d'un use case vers son analyse, puis vers les documents de conception qui le traitent, puis vers les routes REST, commandes CLI et fonctions de port citées ;
- d'une règle métier (RM), d'une exigence non fonctionnelle (ENF), d'une exigence fonctionnelle (EF) ou d'un flux de use case vers les tests qui la couvrent ;
- d'une fonction, route ou commande vers les use cases qu'elle sert, les specs qui la citent et les tests qui l'exercent.

La traçabilité est **dérivée** des specs et du code : elle n'est jamais saisie à la main, et se régénère après toute modification des specs, du code ou des tests.

## 2. Outil

`cmd/tracegen` est un outil de développement, exécuté par `make docs-trace` (`go run ./cmd/tracegen`). Il n'est compilé dans aucun binaire livré et n'utilise que la bibliothèque standard Go (`go/parser`, `go/ast`, `regexp`) : aucune dépendance n'est ajoutée à `go.mod` ni à `vendor/`. Il lit les fichiers sans exécuter le code ni les tests.

`go run ./cmd/tracegen -check` n'écrit rien et se termine en erreur si une section générée de `specs/` n'est pas à jour — contrôle utilisable en intégration continue.

## 3. Ce qui est généré

### Dans `specs/` (versionné)

| Élément | Où | Contenu |
|---------|----|---------|
| Bloc **Tests** (`tests-obsidian`) | Fin de chaque section RM, ENF (`Regles_Metier.md`, `Exigences_Non_Fonctionnelles.md`), EF (`Matrice_Tracabilite.md`) et de chaque section « Flux … » d'une analyse de use case | Tests qui couvrent l'exigence, explicites puis déduits (au plus 15 listés, le reste renvoie à la matrice de couverture) |
| Section **Liens** (`liens-obsidian`) | Fin de chaque use case et de chaque document de conception (hors `*_intro` et `todo`) | Expression → analyse ; analyse → sections de conception qui citent le use case (à défaut, code cité par l'analyse) ; conception → notes de code citées, section par section |
| Tags dérivés | Frontmatter | `relecture/…` (annotations ouvertes) pour tout document ; `rm/…` et `enf/…` (règles et exigences citées) pour un use case ou un document de conception ; `uc/…` (périmètre déclaré, ligne « Use cases : » de l'en-tête) pour un document de conception |
| Carte des specs | `specs/Carte_des_specs.md` | Use cases par famille avec leurs couches et leur conception, règles et exigences avec les use cases qui les citent, périmètre des documents de conception |

Les sections générées sont encadrées par des marqueurs HTML qui rappellent qu'elles ne s'éditent pas à la main ; tout le reste du document est conservé tel quel. Les tags `couche/`, `type/`, `famille/`, `domaine/` sont saisis à la main et conservés.

### Dans `docs/` (non versionné)

| Note | Contenu |
|------|---------|
| `docs/code/<package>.md` | Couche hexagonale, dépendances, specs qui décrivent le package, ports et implémentations, routes ou commandes, use cases liés, tests, fichiers |
| `docs/code/fonctions/<domaine>.<Port>.<Méthode>.md` | Signature, déclaration et implémentation, routes et commandes appelantes, use cases impactés, tests, specs qui la citent |
| `docs/code/routes/REST <route>.md` · `docs/code/commandes/CLI <commande>.md` | Enregistrement ou définition, fonctions du domaine appelées, tests, use cases et specs qui la citent |
| `docs/code/Architecture.md` · `docs/code/Tracabilite_UC_Code.md` | Graphe des dépendances entre packages ; code sollicité par chaque use case, tests explicites par règle, citations des specs sans correspondance dans le code |
| `docs/tests/<package>/<Test>.md` · `docs/tests/Matrice_Couverture_Tests.md` | Commentaire structuré du test, exigences couvertes, code exercé ; couverture de chaque exigence, tests sans exigence, identifiants cités par des tests mais absents des specs |

Chaque note de `docs/` indique le commit analysé : ses numéros de ligne s'y rapportent.

## 4. Règles de rattachement

### 4.1 Citations relevées dans les specs

Le texte d'un document et de chacune de ses sections (titre compris, sans les sections générées) est analysé ligne à ligne, blocs de code et diagrammes compris :

- **Use cases** : identifiants `UC<famille><nn>`, plages (`UCA01–03`) et suffixes (`UCAM01–03/05/07`) développés ; seuls les use cases existants sont retenus pour le périmètre ;
- **Exigences** : `RMnn`, `ENFnn`, `EFnn` ;
- **Routes** : chemins `/api/…` précédés ou non d'une méthode HTTP (`POST /api/identity/request`) ;
- **Commandes** : `myr` suivi de mots en minuscules (`myr identity request`) ; une citation en code s'arrête à la fin de son span ;
- **Fonctions** : `Port.Méthode` (`IdentityService.Enroll`), récepteur nommé d'après un domaine (`identitySvc.RequestAccess` → domaine `identity`), ou nom seul en code ou suivi d'une parenthèse, retenu s'il désigne une seule méthode de tous les ports et n'est pas un nom générique (`Get`, `List`, `Save`…) ;
- **Annotations** : `#question`, `#remarque`, `#incoherence`, `#ecart`, hors spans de code.

### 4.2 Résolution dans le code

- **Ports** : interfaces de `domain/<domaine>/port_in.go` (ports d'entrée) et `ports.go` (ports de sortie) ; un type implémente un port s'il en porte toutes les méthodes, avec le même nombre de paramètres et de résultats.
- **Routes** : enregistrements `mux.HandleFunc` / `mux.Handle` de `adapters/in/rest` ; une citation désigne le motif exact, sinon le plus long motif de sous-arbre (terminé par `/`) qui la préfixe. `/api/components` et `/api/components/` partagent une note.
- **Commandes** : variables `cobra.Command` de `adapters/in/cli` et leur arbre `AddCommand` ; une citation désigne la plus longue commande qu'elle préfixe — elle est sans correspondance si le mot suivant aurait dû être une sous-commande.
- **Appels** : depuis le handler d'une route ou le `RunE` d'une commande, en suivant les fonctions et méthodes du même package. Un appel `x.M` est rattaché par le type de `x` (paramètre, champ du récepteur, variable de package ou résultat d'une fonction du package), sinon par le nom `M` s'il est unique parmi les ports d'entrée. Une valeur de méthode (`decide := svc.ApproveRequest`) est rattachée si le type est connu.
- **Tests** : le package testé est le dossier du test, sans `/tests` ni `/integration`. Un test exerce une fonction de port qu'il appelle (port du package testé, ou qu'il implémente), une route dont il cite le chemin en littéral, une commande dont il appelle la fonction d'exécution.

### 4.3 Couverture d'une exigence

- **Explicite** : le nom ou le commentaire du test cite l'identifiant de l'exigence (`RMnn`, `ENFnn`, `EFnn`).
- **Déduite** : le test exerce une route, une commande ou une fonction que la section de l'exigence cite — couverture probable, à confirmer en relisant le test.
- Une **EF** hérite des tests des flux des use cases listés dans sa ligne « Use cases couvrants ».
- Un flux de use case n'a pas d'identifiant : sa couverture est seulement déduite.

## 5. Limites

- L'analyse est syntaxique, sans vérification de types : un appel à travers une interface intermédiaire, une fonction d'un autre package ou une variable locale non déclarée explicitement n'est rattaché que par le nom de la méthode.
- Une citation de méthode ambiguë (nom porté par plusieurs ports, ou générique) n'est rattachée que qualifiée (`Port.Méthode`).
- La couverture déduite n'atteste pas qu'un test vérifie l'exigence : seule la couverture explicite le fait. Citer l'identifiant dans le commentaire `@brief` d'un test la rend explicite.
- L'outil ne vérifie ni que les tests passent, ni leur taux de couverture de code.

## 6. Décisions de conception

| ID | Décision | Raison |
|----|----------|--------|
| DC-TRACE-01 | Traçabilité dérivée par un outil du dépôt, jamais saisie à la main | Une traçabilité saisie à la main se désynchronise dès le premier changement de code ou de spec ; l'outil versionné avec le code la reconstruit à l'identique. |
| DC-TRACE-02 | Bibliothèque standard Go seule, exécution par `go run` | Aucune technologie ni dépendance nouvelle ; l'outil se lance là où le code se compile, sans installation. |
| DC-TRACE-03 | Sections générées encadrées de marqueurs, le reste des specs intact | Les specs restent éditées à la main ; seule la partie dérivée est remplacée, et une régénération sans changement ne modifie aucun fichier. |
| DC-TRACE-04 | Notes de `docs/` non versionnées | Elles portent des numéros de ligne et le commit analysé, qui changent à chaque commit ; seules les sections de `specs/`, stables, sont versionnées. |

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Étape suivante — code, par section**
- [§ 4.1 Citations relevées dans les specs](Tracabilite.md#4.1%20Citations%20relevées%20dans%20les%20specs) → [REST /api/identity/request](../../docs/code/routes/REST%20api-identity-request.md) · [CLI myr identity request](../../docs/code/commandes/CLI%20myr-identity-request.md) · [IdentityService.Enroll](../../docs/code/fonctions/identity.IdentityService.Enroll.md) · [IdentityService.RequestAccess](../../docs/code/fonctions/identity.IdentityService.RequestAccess.md)
- [§ 4.2 Résolution dans le code](Tracabilite.md#4.2%20Résolution%20dans%20le%20code) → [REST /api/components](../../docs/code/routes/REST%20api-components.md)

<!-- liens-obsidian:end -->
