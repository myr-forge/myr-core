---
id: FT-001
titre: "Propriété des assets non contrôlée côté serveur"
type: securite
statut: ouvert
severite: critique
detecte: 2026-10-02
maj: 2026-10-02
composants: [adapters/in/rest, domain/model]
uc: [UCCE01, UCCE02, UCCE07, UCMOD03, UCMOD04, UCMOD08, UCPI07]
rm: [RM08, RM23, RM24]
enf: [ENF12]
tags:
  - ticket
  - ticket/securite
  - statut/ouvert
  - severite/critique
  - domaine/model
  - uc/UCCE01
  - uc/UCCE02
  - uc/UCCE07
  - uc/UCMOD03
  - uc/UCMOD04
  - uc/UCMOD08
  - uc/UCPI07
  - rm/RM08
  - rm/RM23
  - rm/RM24
  - enf/ENF12
---
# FT-001 — Propriété des assets non contrôlée côté serveur

> **Sécurité** · sévérité **critique** · statut **Ouvert** · détecté le 2026-10-02
> Source : Analyse du code (création des tickets) ; annotation #incoherence de Modele_Domaine.md

## Constat

À la création d'un asset, le propriétaire est lu dans le formulaire envoyé par le client (`owner_id`) au lieu d'être déduit de l'identité de la session : `POST /api/components` et l'import par lot prennent `r.FormValue("owner_id")` tel quel. Aucun contrôle ne compare ensuite le propriétaire de l'asset à l'identité appelante avant une modification, une suppression ou une soumission — ni dans `adapters/in/rest`, ni dans `domain/model`. Le seul usage de `OwnerID` côté REST est un filtre de listing (`?owner_id=`).

## Cause

Le domaine ne reçoit jamais l'identité de l'appelant : les méthodes de `ModelService` prennent l'`OwnerID` comme une donnée parmi d'autres. Le contrôle d'accès REST (`requireRole`) vérifie une permission (`write`), pas la propriété d'une ressource.

## Impact

- Tout porteur de la permission `write` peut créer un asset au nom d'un autre utilisateur, et modifier, supprimer ou soumettre les assets d'autrui.
- La chaîne de propriété, sur laquelle reposent les commissions (RM23/RM24) et le transfert de propriété (UCPI07), devient falsifiable.
- Le comportement attendu par UCMOD04 (403 sur le brouillon d'un autre utilisateur) n'est pas assuré.

## Preuves

- `handlers.go:640` et `handlers.go:796` : `ownerID := strings.TrimSpace(r.FormValue("owner_id"))`.
- Aucune occurrence de comparaison `OwnerID` / session dans `adapters/in/rest/*.go` ni `domain/model/service.go` (seule comparaison : le filtre de listing, `handlers.go:529`).

## Piste de correction (à valider par le PO)

Faire porter l'identité appelante jusqu'au domaine (paramètre explicite ou contexte) pour que la règle de propriété vive dans `domain/model` et s'applique identiquement au CLI et au REST : `OwnerID` imposé à la création, contrôle de propriété sur toute mutation d'un asset. Côté REST, ne plus accepter `owner_id` du client pour une création.

## Critères de clôture

- [ ] `OwnerID` d'un nouvel asset dérivé de l'identité de session, jamais du client
- [ ] Modification, suppression et soumission refusées (403) pour un asset dont l'appelant n'est pas propriétaire
- [ ] Règle implémentée dans le domaine, couverte par des tests unitaires (`domain/model/tests`) et REST

## Liens

- **Use cases** : [UCCE01](../specs/2-Analyse/UCCE-Composant_Ecriture/UCCE01.md) · [UCCE02](../specs/2-Analyse/UCCE-Composant_Ecriture/UCCE02.md) · [UCCE07](../specs/2-Analyse/UCCE-Composant_Ecriture/UCCE07.md) · [UCMOD03](../specs/2-Analyse/UCMOD-Module/UCMOD03.md) · [UCMOD04](../specs/2-Analyse/UCMOD-Module/UCMOD04.md) · [UCMOD08](../specs/2-Analyse/UCMOD-Module/UCMOD08.md) · [UCPI07](../specs/2-Analyse/UCPI-Propriete_Intellectuelle/UCPI07.md)
- **Règles métier** : `RM08`, `RM23`, `RM24` (tags `rm/…`)
- **Exigences non fonctionnelles** : `ENF12` (tags `enf/…`)
- **Specs** : [Modele_Domaine](../specs/3-Conception/Modele_Domaine.md) · [Securite](../specs/3-Conception/Securite.md)
- **Code** : `adapters/in/rest/handlers.go:640` · `adapters/in/rest/handlers.go:796` · `adapters/in/rest/handlers.go:529`
- **Fonctions** : [ModelService.AddFull](../docs/code/fonctions/model.ModelService.AddFull.md) · [ModelService.Remove](../docs/code/fonctions/model.ModelService.Remove.md) · [ModelService.Submit](../docs/code/fonctions/model.ModelService.Submit.md)
- **Tickets liés** : [FT-002 — Rôle de session REST codé en dur à contributor](FT-002-role-de-session-rest-code-en-dur.md)

## Historique

- 2026-10-02 — Ticket créé ; constat vérifié dans le code au commit `2aa69c1`
