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

### Flux nominal — Aucun token disponible

1. Le client n'a pas (ou plus) de token en cache
2. Il doit passer par UCA02 (connexion ou accès invité) avant tout appel protégé

### Flux erreur — Token expiré ou invalide

1. Le client envoie une requête protégée avec un token absent du store ou expiré
2. `requireAuth` répond `HTTP 401` avec `{"error": "authentification requise"}`
3. Le client doit se reconnecter (UCA02) pour obtenir un nouveau token — il n'y a pas de rafraîchissement automatique (pas de refresh token, le token opaque est valide 7 jours en bloc)

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

**Navigation**
- [Carte des specs › UCA — Compte et Acces](../../Carte_des_specs.md#UCA%20—%20Compte%20et%20Acces)
- [UCA04 — couche expression](../../1-Expression/UCA-Compte_et_Acces/UCA04.md)
- [Traçabilité UCA04 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCA04)

**Exigences fonctionnelles couvertes**
- [EF04 — Vérifier la validité d'une session active](../../1-Expression/Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Use cases cités**
- [UCA02 — Se Connecter](UCA02.md)

**Règles métier**
- [RM22 — Changement de rôle réservé à l'administrateur](../../1-Expression/Regles_Metier.md#6.%20Compte%20et%20accès)

**Exigences non fonctionnelles**
- [ENF12 — Contrôle d'accès par rôle](../../1-Expression/Exigences_Non_Fonctionnelles.md#Tableau%20des%20exigences%20non-fonctionnelles)

**Cité par**
- [Expression_des_besoins_Intro](../../1-Expression/Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../../1-Expression/Matrice_Tracabilite.md)
- [todo (expression)](../../1-Expression/todo.md)
- [Analyse_des_besoins](../Analyse_des_besoins.md)
- [UCA02 (analyse)](UCA02.md)
- [UCA07 (analyse)](UCA07.md)
- [todo (analyse)](../todo.md)
- [DC_CLI_Identity](../../3-Conception/DC_CLI_Identity.md)
- [todo (conception)](../../3-Conception/todo.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
