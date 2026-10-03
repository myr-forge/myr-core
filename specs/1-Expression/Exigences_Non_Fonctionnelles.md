---
tags:
  - couche/expression
  - type/exigences
---
# Exigences non-fonctionnelles — Myr System

Ce document structure les contraintes non-fonctionnelles du projet Myr selon les catégories ISO 25010. Pour chaque exigence, un critère d'acceptance mesurable est défini. La section 2.3 de l'introduction présente ces mêmes contraintes sous forme textuelle ; ce document en est la version formalisée et vérifiable.

## Liste des exigences non-fonctionnelles

### Performance

#### ENF01 — Temps de réponse des endpoints REST

- **Catégorie** : Performance
- **Critère d'acceptance** : ≤ 500 ms au 95e percentile, hors opérations blockchain
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF02 — Temps de soumission d'une transaction Fabric

- **Catégorie** : Performance
- **Critère d'acceptance** : ≤ 30 s (endorsement + commit), réseau local
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF03 — Génération d'une BOM module

- **Catégorie** : Performance
- **Critère d'acceptance** : ≤ 5 s pour un module contenant jusqu'à 100 composants
- **Priorité** : MOYENNE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Disponibilité

#### ENF05 — Disponibilité du serveur Myr (instance unique)

- **Catégorie** : Disponibilité
- **Critère d'acceptance** : ≥ 99 % sur 30 jours glissants, hors maintenance planifiée
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF06 — Disponibilité du réseau blockchain (multi-nœuds)

- **Catégorie** : Disponibilité
- **Critère d'acceptance** : ≥ 99,9 % — le réseau reste opérationnel si un nœud sur trois est indisponible
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF07 — Reprise après redémarrage du serveur

- **Catégorie** : Disponibilité
- **Critère d'acceptance** : Le serveur est opérationnel en ≤ 60 s après redémarrage
- **Priorité** : MOYENNE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Sécurité

#### ENF08 — Chiffrement des communications

- **Catégorie** : Sécurité
- **Critère d'acceptance** : Toutes les communications serveur ↔ client en HTTPS (TLS 1.2 minimum)
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF09 — Durée de vie des tokens de session

- **Catégorie** : Sécurité
- **Critère d'acceptance** : Expiration ≤ 24 h ; révocation possible côté serveur
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF10 — Chiffrement des wallets Fabric

- **Catégorie** : Sécurité
- **Critère d'acceptance** : Wallets stockés chiffrés au repos
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF11 — Secrets absents des logs

- **Catégorie** : Sécurité
- **Critère d'acceptance** : Aucun secret d'enrôlement CA, clé privée ni certificat X.509 ne doit apparaître dans les logs applicatifs
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF12 — Contrôle d'accès par rôle

- **Catégorie** : Sécurité
- **Critère d'acceptance** : Toute action est vérifiée côté serveur selon le rôle de l'utilisateur authentifié — le client ne peut pas élever ses droits
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF13 — Protection contre l'injection

- **Catégorie** : Sécurité
- **Critère d'acceptance** : Les entrées utilisateur sont validées avant toute soumission blockchain ou requête base de données
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Scalabilité

#### ENF14 — Nombre d'assets par réseau

- **Catégorie** : Scalabilité
- **Critère d'acceptance** : ≤ 10 000 assets sans dégradation mesurable des temps de réponse (ENF01)
- **Priorité** : MOYENNE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF15 — Taille maximale d'un fichier CAO

- **Catégorie** : Scalabilité
- **Critère d'acceptance** : Jusqu'à 100 Mo par fichier asset (stocké via le dépôt distribué 3D)
- **Priorité** : MOYENNE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF16 — Utilisateurs simultanés par instance

- **Catégorie** : Scalabilité
- **Critère d'acceptance** : 50 utilisateurs simultanés sans dégradation des temps de réponse (ENF01)
- **Priorité** : MOYENNE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF17 — Extension horizontale

- **Catégorie** : Scalabilité
- **Critère d'acceptance** : L'architecture supporte plusieurs instances Myr en parallèle via sessions Redis partagées (`REDIS_URL`)
- **Priorité** : BASSE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Maintenabilité

#### ENF18 — Isolation du domaine métier

- **Catégorie** : Maintenabilité
- **Critère d'acceptance** : Le domaine (`domain/`) ne contient aucune dépendance directe à Fabric, Redis, SQLite ou IPFS — vérifié par revue de code
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF19 — Couverture de tests unitaires

- **Catégorie** : Maintenabilité
- **Critère d'acceptance** : ≥ 80 % des services domaine couverts par des tests unitaires (`go test ./domain/...`)
- **Priorité** : MOYENNE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF20 — Durée de déploiement d'une mise à jour

- **Catégorie** : Maintenabilité
- **Critère d'acceptance** : `make deploy` complet ≤ 10 min sur le serveur cible (compilation Linux + transfert + redémarrage)
- **Priorité** : BASSE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF21 — Compatibilité multi-réseaux blockchain

- **Catégorie** : Maintenabilité
- **Critère d'acceptance** : L'adapter Fabric est interchangeable sans modification du domaine — un nouvel adapter blockchain peut être branché en implémentant les ports `out` existants
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Portabilité

#### ENF23 — Plateformes serveur

- **Catégorie** : Portabilité
- **Critère d'acceptance** : Binaires disponibles pour Linux amd64 et Windows amd64 (produits par `make api` et `make cli`)
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Conformité légale

#### ENF25 — Licence du code source

- **Catégorie** : Conformité
- **Critère d'acceptance** : Code source publié sous AGPL 3.0 ; toute contribution ou extension doit respecter cette licence
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF26 — Compatibilité de licence des assets

- **Catégorie** : Conformité
- **Critère d'acceptance** : La compatibilité entre licence parent et licence dérivé est vérifiée automatiquement avant toute soumission d'un asset dérivé (`ParentID != ""`)
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF27 — Protection des données personnelles (RGPD)

- **Catégorie** : Conformité
- **Critère d'acceptance** : Les données personnelles (email, profil) sont stockées uniquement en base de données locale ; aucun transfert vers un tiers sans consentement explicite
- **Priorité** : MOYENNE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

### Fiabilité

#### ENF28 — Immuabilité des transactions blockchain

- **Catégorie** : Fiabilité
- **Critère d'acceptance** : Aucune opération de suppression n'est implémentée sur la blockchain ; `ErrNotSupported` est retourné si une telle opération est tentée
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF29 — Anti-plagiat obligatoire

- **Catégorie** : Fiabilité
- **Critère d'acceptance** : Tout asset de type `base` passe une vérification SHA-256 + similarité SCM > 50 % avant enregistrement — aucune exception possible
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF30 — Intégrité en cas d'échec blockchain

- **Catégorie** : Fiabilité
- **Critère d'acceptance** : Si une transaction blockchain échoue, l'état local (draft) est conservé intact — aucune perte de données côté serveur
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF31 — Validation avant soumission

- **Catégorie** : Fiabilité
- **Critère d'acceptance** : Toutes les données sont validées côté serveur avant soumission blockchain — les erreurs de validation ne génèrent pas de transaction partielle
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

#### ENF32 — Traçabilité et intégrité du fichier source

- **Catégorie** : Fiabilité
- **Critère d'acceptance** : Pour tout asset possédant un fichier ressource, l'emplacement de stockage effectif et la correspondance de son hash sont vérifiables à la demande, quelle que soit la technologie de stockage sous-jacente (adapter `out/` local, IPFS, ou autre) ; toute divergence (fichier introuvable, hash différent) est signalée explicitement, jamais masquée silencieusement
- **Priorité** : HAUTE

<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->
> [!warning] Tests — aucun test ne couvre cette exigence
> [Matrice de couverture](../../docs/tests/Matrice_Couverture_Tests.md)
<!-- tests-obsidian:end -->

## Récapitulatif par priorité

| Priorité | Nombre | IDs |
|-|--|--|
| HAUTE | 19 | ENF01, ENF02, ENF05, ENF06, ENF08–ENF13, ENF18, ENF21, ENF23, ENF25–ENF26, ENF28–ENF32 |
| MOYENNE | 7 | ENF03, ENF07, ENF14–ENF16, ENF19, ENF27 |
| BASSE | 2 | ENF17, ENF20 |

## Correspondance avec la section 2.3 de l'introduction

| Contrainte textuelle (section 2.3) | ENF correspondant |
|-|-|
| Immuabilité blockchain | ENF28, ENF31 |
| Vérification anti-plagiat | ENF29 |
| Compatibilité de licence | ENF26 |
| Licence Open-Source AGPL 3.0 | ENF25 |
| Architecture hexagonale | ENF18, ENF21 |
| Compatibilité multi-réseaux | ENF21 |
