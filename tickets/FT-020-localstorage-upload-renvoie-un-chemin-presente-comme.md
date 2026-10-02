---
id: FT-020
titre: "LocalStorage.Upload renvoie un chemin présenté comme un hash"
type: dette
statut: ouvert
severite: mineure
detecte: 2026-10-02
maj: 2026-10-02
composants: [adapters/out/localstorage]
uc: [UCCL03]
rm: [RM42]
enf: []
tags:
  - ticket
  - ticket/dette
  - statut/ouvert
  - severite/mineure
  - domaine/model
  - uc/UCCL03
  - rm/RM42
---
# FT-020 — LocalStorage.Upload renvoie un chemin présenté comme un hash

> **Dette technique** · sévérité **mineure** · statut **Ouvert** · détecté le 2026-10-02
> Source : todo.md (travaux traçabilité du fichier source)

## Constat

`FileStoragePort.Upload` est documenté comme retournant un hash ; l'implémentation `LocalStorage` retourne le chemin local du fichier.

## Cause

Nommage hérité ; sans effet fonctionnel aujourd'hui (le hash de comparaison est `Model3D.Hash`, calculé séparément).

## Impact

Risque de confusion pour une future implémentation ou un futur appelant qui comparerait cette valeur à un hash.

## Preuves

`adapters/out/localstorage/store.go` ; documentation de `FileStoragePort` dans `domain/model/ports.go`.

## Piste de correction (à valider par le PO)

Renommer la valeur de retour (« référence de stockage ») dans le port et sa documentation.

## Critères de clôture

- [ ] Contrat du port sans ambiguïté entre référence et hash

## Liens

- **Use cases** : [UCCL03](../specs/1-Expression/UCCL-Composant_Lecture/UCCL03.md)
- **Règles métier** : `RM42` (tags `rm/…`)
- **Code** : `adapters/out/localstorage/store.go` · `domain/model/ports.go`
- **Fonctions** : [FileStoragePort.Upload](../docs/code/fonctions/model.FileStoragePort.Upload.md)

## Historique

- 2026-10-02 — Ticket créé
