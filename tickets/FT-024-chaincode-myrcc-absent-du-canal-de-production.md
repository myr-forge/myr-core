---
id: FT-024
titre: "Chaincode myrcc absent du canal de production"
type: anomalie
statut: a-verifier
severite: critique
detecte: 2026-07-29
maj: 2026-10-02
composants: [chaincode, adapters/out/fabric]
uc: [UCCL01, UCCE01, UCMOD06]
rm: [RM06]
enf: [ENF02]
tags:
  - ticket
  - ticket/anomalie
  - statut/a-verifier
  - severite/critique
  - domaine/model
  - uc/UCCL01
  - uc/UCCE01
  - uc/UCMOD06
  - rm/RM06
  - enf/ENF02
---
# FT-024 — Chaincode myrcc absent du canal de production

> **Anomalie** · sévérité **critique** · statut **À vérifier** · détecté le 2026-07-29
> Source : roadmap_dev.md § Écarts Infrastructure — Chaincode ; écart E6

## Constat

Le 2026-07-29, sur le serveur de production (réseau `diy-network`, canal `sandbox`), `GET /api/components` et `GET /api/modules` répondaient 500 : `no peers available to evaluate chaincode myrcc in channel sandbox`. Le profil réseau pointait vers le chaincode `myrcc` sans qu'aucun cycle de vie Fabric (install, approve, commit) n'ait été exécuté.

Depuis, le code du chaincode a été ajouté au dépôt (commit `e98b814` : `chaincode/main.go`, `chaincode/model/contract.go` avec `storeModel`, `getModel`, `listModels`, `verifyModel`, entité de 27 champs) — l'écart E6 (entité à 7 champs) est donc probablement résorbé, mais le déploiement sur le serveur n'est pas confirmé.

## Cause

`myr network create` provisionne l'infrastructure Fabric mais n'installe aucun chaincode (`install.md` §7).

## Impact

Toute lecture ou écriture blockchain échoue : listing, soumission, vérification.

## Preuves

Logs `journalctl --user -u myr` et `docker logs peer0.org1.diy-network.com` rapportés dans la roadmap.

## Piste de correction (à valider par le PO)

Vérifier sur le serveur : `peer lifecycle chaincode querycommitted -C sandbox`. Si absent : construire l'image (`chaincode/Dockerfile`), puis package → install → approve (Org1MSP) → commit sur `sandbox`, et seulement ensuite pointer `chaincode_name` du profil réseau. Vérifier l'alignement de l'entité chaincode avec `Model3D`.

## Critères de clôture

- [ ] Chaincode `myrcc` commité sur `sandbox`
- [ ] `GET /api/components` répond 200 en production
- [ ] Entité chaincode alignée avec `Model3D` (E6 clos)
- [ ] Procédure de déploiement documentée (install.md)

## Liens

- **Use cases** : [UCCL01 (analyse)](../specs/2-Analyse/UCCL-Composant_Lecture/UCCL01.md) · [UCCL01 (expression)](../specs/1-Expression/UCCL-Composant_Lecture/UCCL01.md) · [UCCE01 (analyse)](../specs/2-Analyse/UCCE-Composant_Ecriture/UCCE01.md) · [UCCE01 (expression)](../specs/1-Expression/UCCE-Composant_Ecriture/UCCE01.md) · [UCMOD06 (analyse)](../specs/2-Analyse/UCMOD-Module/UCMOD06.md) · [UCMOD06 (expression)](../specs/1-Expression/UCMOD-Module/UCMOD06.md)
- **Règles métier** : [RM06](../specs/1-Expression/Regles_Metier.md)
- **Exigences non fonctionnelles** : [ENF02](../specs/1-Expression/Exigences_Non_Fonctionnelles.md)
- **Specs** : [Chaincode](../specs/3-Conception/Chaincode.md) · [Deploiement](../specs/3-Conception/Deploiement.md)
- **Code** : [chaincode/model/contract.go](../chaincode/model/contract.go) · [chaincode/model/entity.go](../chaincode/model/entity.go)
- **Fonctions** : [BlockchainPort.StoreModelRecord](../docs/code/fonctions/model.BlockchainPort.StoreModelRecord.md) · [BlockchainPort.ListModelRecords](../docs/code/fonctions/model.BlockchainPort.ListModelRecords.md)
- **Tickets liés** : [FT-025 — Le serveur signale un nœud Fabric déconnecté](FT-025-le-serveur-signale-un-nud-fabric-deconnecte.md) · [FT-026 — Erreur Fabric indisponible renvoyée en 500 au lieu de 503](FT-026-erreur-fabric-indisponible-renvoyee-en-500-au.md) · [FT-018 — Identifiant de bloc des versions de module simulé](FT-018-identifiant-de-bloc-des-versions-de-module.md)

## Historique

- 2026-07-29 — Constaté sur le serveur de production
- — — Code du chaincode ajouté (commit `e98b814`)
- 2026-10-02 — Ticket créé ; déploiement à revérifier sur le serveur
