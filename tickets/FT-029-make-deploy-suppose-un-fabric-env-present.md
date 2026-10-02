---
id: FT-029
titre: "make deploy suppose un fabric.env présent sur le serveur"
type: anomalie
statut: ouvert
severite: mineure
detecte: 2026-10-02
maj: 2026-10-02
composants: [scripts, cmd/api]
uc: []
rm: []
enf: [ENF07]
tags:
  - ticket
  - ticket/anomalie
  - statut/ouvert
  - severite/mineure
  - enf/ENF07
---
# FT-029 — make deploy suppose un fabric.env présent sur le serveur

> **Anomalie** · sévérité **mineure** · statut **Ouvert** · détecté le 2026-10-02
> Source : Documentation des scripts de déploiement

## Constat

`deploy.ps1` passe toujours `--env-file` au service sans jamais créer `fabric.env` : si le fichier est absent, `myr-api` s'arrête au démarrage (`log.Fatalf`). `deploy_api.ps1` démarre au contraire en mode dégradé sans Fabric. Les trois scripts ne sont reliés par rien et un seul est appelé par `make deploy`.

## Cause

Scripts écrits à des moments différents.

## Impact

Un déploiement sur un serveur neuf échoue sans message clair ; le comportement dépend du script choisi.

## Preuves

`scripts/deploy.ps1`, `scripts/deploy_api.ps1`.

## Piste de correction (à valider par le PO)

Unifier les scripts (un seul chemin de déploiement), vérifier la présence de `fabric.env` avant de passer `--env-file`, documenter sa création dans install.md.

## Critères de clôture

- [ ] Un seul script de déploiement de référence
- [ ] Absence de `fabric.env` détectée avec un message explicite

## Liens

- **Exigences non fonctionnelles** : [ENF07](../specs/1-Expression/Exigences_Non_Fonctionnelles.md)
- **Specs** : [Deploiement](../specs/3-Conception/Deploiement.md)
- **Tickets liés** : [FT-028 — Mot de passe serveur en clair dans le script de déploiement](FT-028-mot-de-passe-serveur-en-clair-dans.md) · [FT-025 — Le serveur signale un nœud Fabric déconnecté](FT-025-le-serveur-signale-un-nud-fabric-deconnecte.md)

## Historique

- 2026-10-02 — Ticket créé
