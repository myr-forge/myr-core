---
id: FT-002
titre: "Rôle de session REST codé en dur à contributor"
type: securite
statut: ouvert
severite: critique
detecte: 2026-07-17
maj: 2026-10-02
composants: [adapters/in/rest]
uc: [UCA02, UCA05, UCA07]
rm: [RM21, RM22]
enf: [ENF12]
tags:
  - ticket
  - ticket/securite
  - statut/ouvert
  - severite/critique
  - domaine/identity
  - domaine/role
  - uc/UCA02
  - uc/UCA05
  - uc/UCA07
  - rm/RM21
  - rm/RM22
  - enf/ENF12
---
# FT-002 — Rôle de session REST codé en dur à contributor

> **Sécurité** · sévérité **critique** · statut **Ouvert** · détecté le 2026-07-17
> Source : roadmap_dev.md § Écarts Identité & Session ; écart E3 (Analyse_des_besoins.md)

## Constat

À la connexion (`POST /api/identity/session`), la session REST est créée avec le rôle `"contributor"` en dur, quelle que soit l'identité. Le rôle porté par le certificat CA (`Myr.role`) n'est jamais lu, et un changement de rôle (`myr identity set-role`) n'a aucun effet sur les sessions.

## Cause

`handlers_identity.go:391` : `h.sessions.create(req.Name, "contributor", channel)` — annoté `#question` dans le code (ligne 388). La roadmap situe ce bug dans `domain/auth/service.go`, fichier qui n'existe pas : l'emplacement réel est ce handler.

## Impact

- Tout utilisateur connecté obtient les permissions `contributor` (écriture), y compris un compte qui devrait être `reader` (RM21) : élévation de privilège systématique.
- Le RBAC dynamique (`domain/role`) est contourné pour la session REST ; un administrateur ne peut pas restreindre un compte.

## Preuves

`adapters/in/rest/handlers_identity.go:388-391`.

## Piste de correction (à valider par le PO)

Lire le rôle depuis l'attribut `Myr.role` du certificat enrôlé (via `domain/identity`) au moment de créer la session, avec repli fermé (`reader`) si l'attribut est absent. Décider si un changement de rôle invalide les sessions ouvertes ou prend effet au prochain ré-enrôlement (RM22).

## Critères de clôture

- [ ] Le rôle de session reflète le rôle CA de l'identité
- [ ] Un compte sans rôle CA obtient `reader`
- [ ] Test REST couvrant un compte `reader` refusé en écriture

## Liens

- **Use cases** : [UCA02 (analyse)](../specs/2-Analyse/UCA-Compte_et_Acces/UCA02.md) · [UCA02 (expression)](../specs/1-Expression/UCA-Compte_et_Acces/UCA02.md) · [UCA05 (analyse)](../specs/2-Analyse/UCA-Compte_et_Acces/UCA05.md) · [UCA05 (expression)](../specs/1-Expression/UCA-Compte_et_Acces/UCA05.md) · [UCA07 (analyse)](../specs/2-Analyse/UCA-Compte_et_Acces/UCA07.md) · [UCA07 (expression)](../specs/1-Expression/UCA-Compte_et_Acces/UCA07.md)
- **Règles métier** : [RM21](../specs/1-Expression/Regles_Metier.md) · [RM22](../specs/1-Expression/Regles_Metier.md)
- **Exigences non fonctionnelles** : [ENF12](../specs/1-Expression/Exigences_Non_Fonctionnelles.md)
- **Specs** : [DC_D1_Auth_Identity](../specs/3-Conception/DC_D1_Auth_Identity.md) · [Modele_Domaine](../specs/3-Conception/Modele_Domaine.md)
- **Code** : [adapters/in/rest/handlers_identity.go:388](../adapters/in/rest/handlers_identity.go)
- **Fonctions** : [IdentityService.GetStatus](../docs/code/fonctions/identity.IdentityService.GetStatus.md) · [RoleService.HasPermission](../docs/code/fonctions/role.RoleService.HasPermission.md)
- **Tickets liés** : [FT-001 — Propriété des assets non contrôlée côté serveur](FT-001-propriete-des-assets-non-controlee-cote-serveur.md) · [FT-034 — Modèle de domaine : ligne RM22 décrivant un écart](FT-034-modele-de-domaine-ligne-rm22-decrivant-un.md)

## Historique

- 2026-07-17 — Écart documenté dans la roadmap (E3)
- 2026-10-02 — Ticket créé ; constat revérifié au commit `2aa69c1`
