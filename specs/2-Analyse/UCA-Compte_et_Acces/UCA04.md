---
categorie: Compte et Accès
titre: "Vérification de la connexion"
probabilite: 2
impact: 3
importance: 6
etat: analyse
tags:
  - couche/analyse
  - type/use-case
  - famille/UCA
  - domaine/identity
  - domaine/role
  - uc/UCA04
  - rm/RM22
  - enf/ENF12
---

# Vérification de la connexion

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Client\n(dépôt GUI externe)" as U

rectangle "API myr" {
    usecase "Utiliser le token reçu\nà la connexion" as UC1
    usecase "Constater un rejet\n(401) sur une requête protégée" as UC2
}

U --> UC1
U --> UC2

@enduml
```

## Contexte

**Il n'existe aucun endpoint qui répond « suis-je connecté ? / quel est mon rôle actuel ? ».** Il n'y a pas d'équivalent de `GET /api/auth/me` dans le code réel.

Le seul moment où le client apprend son rôle, son pseudo et son canal est la réponse de `POST /api/identity/session` ou `POST /api/identity/guest` (UCA02) — il doit les mettre en cache lui-même. Ensuite, la seule façon de savoir si la session est toujours valide est de tenter un appel protégé et d'observer la réponse :

- `HTTP 200` (ou tout code de succès métier) → le token est toujours valide
- `HTTP 401` → le token est absent, invalide ou expiré — il faut se reconnecter (UCA02)

> La restitution visuelle de cet état (indicateur coloré, tooltip au survol, etc.) est un choix d'interface qui relève du dépôt GUI externe — hors périmètre de ce document.

## Pré-conditions

- Le client dispose (ou non) d'un token obtenu via UCA02

## Scénario

### Flux nominal — Token toujours valide

1. Le client envoie une requête protégée quelconque avec `X-Myr-Token: <token>`
2. `requireAuth` retrouve la session dans le store et vérifie qu'elle n'a pas expiré
3. La requête aboutit normalement — le client en déduit que sa session est active

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux nominal — Aucun token disponible

1. Le client n'a pas (ou plus) de token en cache
2. Il doit passer par UCA02 (connexion ou accès invité) avant tout appel protégé

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Flux erreur — Token expiré ou invalide

1. Le client envoie une requête protégée avec un token absent du store ou expiré
2. `requireAuth` répond `HTTP 401` avec `{"error": "authentification requise"}`
3. Le client doit se reconnecter (UCA02) pour obtenir un nouveau token — il n'y a pas de rafraîchissement automatique (pas de refresh token, le token opaque est valide 7 jours en bloc)

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

## Post-conditions

- Le client sait, après une tentative de requête, si son token est encore valide ou non
- Aucune information de session n'est accessible autrement que par cette tentative ou par ce qui a été mis en cache à la connexion

## Diagramme de séquence

```plantuml
@startuml
participant "Client\n(dépôt GUI externe)" as Client
participant "REST Handler" as REST
participant "requireAuth\nmiddleware\n(handlers.go)" as Middleware
participant "Session Store\n(adapters/in/rest/session.go)" as Sessions

Client -> REST : GET /api/... (endpoint protégé quelconque)\nX-Myr-Token: <token>
REST -> Middleware : requireAuth
Middleware -> Sessions : get(token)

alt session valide et non expirée
  Sessions --> Middleware : *myrSession
  Middleware -> REST : ctx avec session
  REST --> Client : HTTP 200 (résultat normal)
else token absent, inconnu ou expiré
  Sessions --> Middleware : nil
  Middleware --> Client : HTTP 401 {error: "authentification requise"}
end
@enduml
```

## Règles métier déclenchées

- **RM22** — La validité de la session est vérifiée côté serveur à chaque requête protégée (`requireAuth`) ; le client ne peut pas décider seul qu'une session est valide.

## Exigences non-fonctionnelles

- **ENF12** — La validation de session est strictement côté serveur.

## Notes d'implémentation

**Pas d'endpoint dédié :** `requireAuth` (`adapters/in/rest/handlers.go`) est un middleware appliqué à des routes métier — il n'existe pas de route « neutre » dont le seul rôle serait de vérifier/retourner l'état de connexion.

`requireAuth` valide la session à chaque requête protégée, sans exception ni contournement.

**Statut d'implémentation :**
- Vérification de session sur requête protégée : **opérationnelle** (`requireAuth`)
- Endpoint de consultation de session (« whoami ») : **inexistant** — écart réel, pas une simplification de cette spec

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Étape suivante — conception**
- **DC_CLI_Identity** : [§ DC — CLI Identité & Session : Référence des commandes `myr identity` / `myr session`](../../3-Conception/DC_CLI_Identity.md#DC%20—%20CLI%20Identité%20&%20Session%20:%20Référence%20des%20commandes%20`myr%20identity`%20/%20`myr%20session`) · [§ 2. Arbre de commandes](../../3-Conception/DC_CLI_Identity.md#2.%20Arbre%20de%20commandes)

<!-- liens-obsidian:end -->
