---
id: FT-028
titre: "Mot de passe serveur en clair dans le script de déploiement"
type: securite
statut: ouvert
severite: majeure
detecte: 2026-10-02
maj: 2026-10-02
composants: [scripts]
uc: []
rm: []
enf: [ENF11]
tags:
  - ticket
  - ticket/securite
  - statut/ouvert
  - severite/majeure
  - enf/ENF11
---
# FT-028 — Mot de passe serveur en clair dans le script de déploiement

> **Sécurité** · sévérité **majeure** · statut **Ouvert** · détecté le 2026-10-02
> Source : install_TODO.md ; note « Ticket » du vault (mot de passe partagé en clair)

## Constat

`scripts/deploy.ps1` (appelé par `make deploy`) contient le mot de passe SSH du serveur en dur (`$RemotePass`). Le dossier `scripts/` n'est pas versionné, mais le mot de passe a aussi été partagé en clair dans une conversation.

## Cause

Script historique ; `deploy_api.ps1` et `deploy_cli.ps1` utilisent déjà une clé SSH.

## Impact

Compromission possible du compte `fabricadmin` du serveur de production.

## Preuves

`scripts/deploy.ps1:12`.

## Piste de correction (à valider par le PO)

Changer le mot de passe du compte sur le serveur ; faire pointer `make deploy` vers les scripts à clé SSH (ou retirer `deploy.ps1`) ; désactiver l'authentification SSH par mot de passe si possible.

## Critères de clôture

- [ ] Mot de passe changé
- [ ] Aucun secret en clair dans les scripts
- [ ] `make deploy` utilise la clé SSH

## Liens

- **Exigences non fonctionnelles** : `ENF11` (tags `enf/…`)
- **Specs** : [Deploiement](../specs/3-Conception/Deploiement.md) · [Securite](../specs/3-Conception/Securite.md)
- **Code** : `scripts/deploy.ps1:12` · `scripts/deploy_api.ps1` _(dossier `scripts/` non versionné : liens valables en local uniquement)_
- **Tickets liés** : [FT-029 — make deploy suppose un fabric.env présent sur le serveur](FT-029-make-deploy-suppose-un-fabric-env-present.md)

## Historique

- 2026-10-02 — Ticket créé
