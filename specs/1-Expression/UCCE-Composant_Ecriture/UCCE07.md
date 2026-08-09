---
categorie: Composant Ecriture
titre: "Supprimer un Composant"
probabilite: 3
impact: 4
importance: 12
etat: relire
---

# Supprimer un Composant

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Concepteur" as C

rectangle "Application MYR" {
    usecase "Supprimer un composant" as UC1
}

C --> UC1

@enduml
```

## Contexte

Un composant qui n'a plus d'usage doit pouvoir être supprimé par son propriétaire, qu'il soit encore en brouillon ou déjà soumis. Un composant soumis ne peut jamais être retiré de la blockchain (RM06) : le supprimer le masque des listes du système, sans jamais toucher à son enregistrement blockchain.

## Pré-conditions

- Être connecté au réseau
- Être propriétaire du composant

## Scénario

**Étape initiale :** `myr model remove <id>` est exécutée (ou l'appel API équivalent `DELETE /api/components/:id`)

### Flux nominal — Composant en brouillon

1. Le composant, jamais soumis, est retiré du stockage local
2. Il n'existe plus nulle part

### Flux nominal — Composant déjà soumis

1. Le composant est masqué des listes retournées par le système
2. Son enregistrement sur la blockchain reste inchangé, consultable par identifiant direct

## Post-conditions

- Le composant n'apparaît plus dans les listes
- Son enregistrement blockchain (s'il existait) n'est jamais modifié

## Diagramme d'activités

```plantuml
@startuml
skin rose
title Supprimer un Composant
start
:Demander la suppression (myr model remove);
if (Composant déjà soumis?) then (oui)
  :Masquer localement — ledger inchangé;
else (non, brouillon)
  :Retirer du stockage local;
endif
stop
@enduml
```
