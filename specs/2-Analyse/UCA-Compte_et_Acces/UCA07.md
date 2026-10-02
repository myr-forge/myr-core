---
categorie: Compte et Accès
titre: "Vérification du rôle attribué"
probabilite: 1
impact: 1
importance: 1
etat: analyse
tags:
  - couche/analyse
  - type/use-case
  - famille/UCA
  - domaine/identity
  - domaine/role
  - uc/UCA07
  - rm/RM22
---

# Vérification du rôle attribué

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Client\n(dépôt GUI externe)" as U

rectangle "API myr" {
    usecase "Mettre en cache le rôle\nreçu à la connexion" as UC1
}

U --> UC1

@enduml
```

## Contexte

Comme pour UCA04, **il n'existe aucun endpoint permettant à un client de (re)demander son rôle courant au serveur.** Le rôle n'est communiqué qu'une seule fois : dans la réponse de `POST /api/identity/session` ou `POST /api/identity/guest` (UCA02), sous la forme `{role, pseudo, channel, ...}`. Le client doit le mettre en cache lui-même pour toute la durée de vie de sa session (jusqu'à 7 jours).

> La restitution de cette information à l'utilisateur (infobulle, libellé traduit, etc.) est un choix d'interface qui relève du dépôt GUI externe — hors périmètre de ce document. Seul le fait que le rôle soit une donnée serveur, non falsifiable côté client, fait partie de `myr`.

## Pré-conditions

- Le client s'est connecté au moins une fois (UCA02) et a reçu un `role` dans la réponse

## Scénario

### Flux nominal — Rôle connu depuis la connexion

1. Le client lit le champ `role` reçu lors de `POST /api/identity/session` ou `POST /api/identity/guest`
2. Il n'y a rien d'autre à faire — cette valeur reste valable tant que le token n'est pas expiré ou révoqué

### Flux alternatif — Vérification indirecte par tentative d'action

1. Le client tente une action nécessitant une permission donnée
2. Si `HTTP 403` est retourné, le client peut en déduire que son rôle en cache ne porte pas (ou plus) cette permission — sans que le serveur ne lui indique explicitement son rôle actuel

## Post-conditions

- Le client dispose (ou non) d'une information de rôle en cache — aucune source de vérité consultable à la demande côté serveur

## Diagramme de séquence

```plantuml
@startuml
participant "Client\n(dépôt GUI externe)" as Client
participant "REST Handler\n(adapters/in/rest/handlers_identity.go)" as REST

Client -> REST : POST /api/identity/session\n{name, secret, org_id}
REST --> Client : HTTP 200 {token, role, pseudo, channel}
Client -> Client : mettre en cache {role, pseudo, channel}

note over Client : Aucun endpoint ultérieur ne permet\nde revérifier ce rôle auprès du serveur
@enduml
```

## Règles métier déclenchées

- **RM22** — Le rôle communiqué au client provient exclusivement du serveur (session REST) — aucun mécanisme ne permet à un client de le modifier lui-même.

## Notes d'implémentation

**Pas d'endpoint de consultation :** aucune route REST ne retourne le rôle courant en dehors de la réponse initiale de connexion. Voir aussi l'écart similaire documenté dans UCA04.

**Rôle figé, potentiellement obsolète :** le rôle mis en cache par le client ne reflète que l'état au moment de la connexion. Si un administrateur modifie le rôle de l'identité entre-temps (`myr identity set-role`, UCA08), le client n'en sera informé qu'à sa prochaine connexion — voir aussi l'écart documenté dans UCA02 sur le rôle de session actuellement figé à `"contributor"`.

**Statut d'implémentation :**
- Transmission du rôle à la connexion : **opérationnelle** (`POST /api/identity/session`, `POST /api/identity/guest`)
- Endpoint de re-consultation du rôle courant : **inexistant** — écart réel, pas une simplification de cette spec

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCA — Compte et Acces](../../Carte_des_specs.md#UCA%20—%20Compte%20et%20Acces)
- [UCA07 — couche expression](../../1-Expression/UCA-Compte_et_Acces/UCA07.md)
- [Traçabilité UCA07 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCA07)

**Exigences fonctionnelles couvertes**
- [EF05 — Contrôler les accès selon le rôle attribué](../../1-Expression/Matrice_Tracabilite.md#1.%20Exigences%20fonctionnelles%20et%20UC%20couvrant)

**Use cases cités**
- [UCA02 — Se Connecter](UCA02.md)
- [UCA04 — Vérification de la connexion](UCA04.md)
- [UCA08 — Demander un rôle](UCA08.md)

**Règles métier**
- [RM22 — Changement de rôle réservé à l'administrateur](../../1-Expression/Regles_Metier.md#6.%20Compte%20et%20accès)

**Cité par**
- [Expression_des_besoins_Intro](../../1-Expression/Expression_des_besoins_Intro.md)
- [Matrice_Tracabilite](../../1-Expression/Matrice_Tracabilite.md)
- [todo (expression)](../../1-Expression/todo.md)
- [Analyse_des_besoins](../Analyse_des_besoins.md)
- [UCA02 (analyse)](UCA02.md)
- [todo (analyse)](../todo.md)
- [DC_CLI_Identity](../../3-Conception/DC_CLI_Identity.md)
- [todo (conception)](../../3-Conception/todo.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
