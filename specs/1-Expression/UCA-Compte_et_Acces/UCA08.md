---
categorie: Compte et Accès
titre: "Demander un rôle"
probabilite: 4
impact: 4
importance: 16
etat: relire
tags:
  - couche/expression
  - type/use-case
  - famille/UCA
  - domaine/identity
  - domaine/role
  - uc/UCA08
---

# Demander un rôle

## Diagramme d'acteurs

```plantuml
@startuml
left to right direction

actor "Utilisateur" as U
actor "Administrateur" as ADM

rectangle "Application MYR" {
    usecase "Demander un rôle" as UC1
    usecase "Traiter la demande" as UC2
}

U --> UC1
ADM --> UC2
UC1 .> UC2 : <<extend>>

@enduml
```

## Contexte

Un utilisateur connecté (rôle Lecteur ou autre) peut demander un rôle supplémentaire depuis son profil. La demande est traitée automatiquement si le rôle est configuré en auto-distribution sur ce réseau, ou soumise à l'administrateur dans le cas contraire.

## Pré-conditions

- Être connecté
- Rôle cible différent du rôle déjà détenu

## Scénario

**Étape initiale :** Le client transmet une demande de rôle pour l'identité connectée (`POST /api/identity/roles/request` — rôle souhaité)

### Flux nominal — Attribution automatique

1. Le rôle souhaité est transmis
2. Le rôle cible est configuré en auto-distribution sur ce réseau
3. Le rôle est attribué immédiatement — la réponse indique le nouveau rôle actif

### Flux alternatif — Validation manuelle par l'administrateur

1. Le rôle souhaité est transmis
2. Le rôle cible nécessite une validation admin
3. La demande est enregistrée en attente (`202 Accepted`)
4. Le rôle courant est conservé jusqu'à la décision de l'administrateur (`myr identity set-role`)

### Flux erreur — Rôle déjà attribué

1. Le rôle demandé est déjà détenu par l'identité
2. Réponse d'erreur : « Vous possédez déjà ce rôle » (`400 Bad Request`)

## Post-conditions

- Rôle attribué immédiatement (auto-distribution), ou demande en attente de validation (manuel)

## Diagramme d'activités

```plantuml
@startuml
skin rose
title Demander un rôle
start
:Transmettre la demande de rôle (POST /api/identity/roles/request);
if (Rôle déjà attribué?) then (oui)
  :Retourner l'erreur "Vous possédez déjà ce rôle" (400);
  stop
else (non)
  if (Auto-distribution activée pour ce rôle?) then (oui)
    :Attribuer le rôle immédiatement;
    stop
  else (non)
    :Enregistrer la demande en attente (202);
    stop
  endif
endif
@enduml
```

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs › UCA — Compte et Acces](../../Carte_des_specs.md#UCA%20—%20Compte%20et%20Acces)
- [UCA08 — couche analyse](../../2-Analyse/UCA-Compte_et_Acces/UCA08.md)
- [Traçabilité UCA08 vers le code](../../../docs/code/Tracabilite_UC_Code.md#UCA08)

**Cité par**
- [Expression_des_besoins_Intro](../Expression_des_besoins_Intro.md)
- [UCA01 (expression)](UCA01.md)
- [todo (expression)](../todo.md)
- [Analyse_des_besoins](../../2-Analyse/Analyse_des_besoins.md)
- [UCA01 (analyse)](../../2-Analyse/UCA-Compte_et_Acces/UCA01.md)
- [UCA02 (analyse)](../../2-Analyse/UCA-Compte_et_Acces/UCA02.md)
- [UCA05 (analyse)](../../2-Analyse/UCA-Compte_et_Acces/UCA05.md)
- [UCA07 (analyse)](../../2-Analyse/UCA-Compte_et_Acces/UCA07.md)
- [todo (analyse)](../../2-Analyse/todo.md)
- [DC_CLI_Identity](../../3-Conception/DC_CLI_Identity.md)
- [todo (conception)](../../3-Conception/todo.md)
- [roadmap_dev](../../roadmap_dev.md)

<!-- liens-obsidian:end -->
