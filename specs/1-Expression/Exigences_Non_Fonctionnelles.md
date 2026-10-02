---
tags:
  - couche/expression
  - type/exigences
---
# Exigences non-fonctionnelles — Myr System

Ce document structure les contraintes non-fonctionnelles du projet Myr selon les catégories ISO 25010. Pour chaque exigence, un critère d'acceptance mesurable est défini. La section 2.3 de l'introduction présente ces mêmes contraintes sous forme textuelle ; ce document en est la version formalisée et vérifiable.



## Tableau des exigences non-fonctionnelles

| ID | Catégorie | Exigence | Critère d'acceptance | Priorité |
|-|--||-||
| **Performance** | | | | |
| ENF01 | Performance | Temps de réponse des endpoints REST | ≤ 500 ms au 95e percentile, hors opérations blockchain | HAUTE |
| ENF02 | Performance | Temps de soumission d'une transaction Fabric | ≤ 30 s (endorsement + commit), réseau local | HAUTE |
| ENF03 | Performance | Génération d'une BOM module | ≤ 5 s pour un module contenant jusqu'à 100 composants | MOYENNE |
| **Disponibilité** | | | | |
| ENF05 | Disponibilité | Disponibilité du serveur Myr (instance unique) | ≥ 99 % sur 30 jours glissants, hors maintenance planifiée | HAUTE |
| ENF06 | Disponibilité | Disponibilité du réseau blockchain (multi-nœuds) | ≥ 99,9 % — le réseau reste opérationnel si un nœud sur trois est indisponible | HAUTE |
| ENF07 | Disponibilité | Reprise après redémarrage du serveur | Le serveur est opérationnel en ≤ 60 s après redémarrage | MOYENNE |
| **Sécurité** | | | | |
| ENF08 | Sécurité | Chiffrement des communications | Toutes les communications serveur ↔ client en HTTPS (TLS 1.2 minimum) | HAUTE |
| ENF09 | Sécurité | Durée de vie des tokens de session | Expiration ≤ 24 h ; révocation possible côté serveur | HAUTE |
| ENF10 | Sécurité | Chiffrement des wallets Fabric | Wallets stockés chiffrés au repos | HAUTE |
| ENF11 | Sécurité | Secrets absents des logs | Aucun secret d'enrôlement CA, clé privée ni certificat X.509 ne doit apparaître dans les logs applicatifs | HAUTE |
| ENF12 | Sécurité | Contrôle d'accès par rôle | Toute action est vérifiée côté serveur selon le rôle de l'utilisateur authentifié — le client ne peut pas élever ses droits | HAUTE |
| ENF13 | Sécurité | Protection contre l'injection | Les entrées utilisateur sont validées avant toute soumission blockchain ou requête base de données | HAUTE |
| **Scalabilité** | | | | |
| ENF14 | Scalabilité | Nombre d'assets par réseau | ≤ 10 000 assets sans dégradation mesurable des temps de réponse (ENF01) | MOYENNE |
| ENF15 | Scalabilité | Taille maximale d'un fichier CAO | Jusqu'à 100 Mo par fichier asset (stocké via le dépôt distribué 3D) | MOYENNE |
| ENF16 | Scalabilité | Utilisateurs simultanés par instance | 50 utilisateurs simultanés sans dégradation des temps de réponse (ENF01) | MOYENNE |
| ENF17 | Scalabilité | Extension horizontale | L'architecture supporte plusieurs instances Myr en parallèle via sessions Redis partagées (`REDIS_URL`) | BASSE |
| **Maintenabilité** | | | | |
| ENF18 | Maintenabilité | Isolation du domaine métier | Le domaine (`domain/`) ne contient aucune dépendance directe à Fabric, Redis, SQLite ou IPFS — vérifié par revue de code | HAUTE |
| ENF19 | Maintenabilité | Couverture de tests unitaires | ≥ 80 % des services domaine couverts par des tests unitaires (`go test ./domain/...`) | MOYENNE |
| ENF20 | Maintenabilité | Durée de déploiement d'une mise à jour | `make deploy` complet ≤ 10 min sur le serveur cible (compilation Linux + transfert + redémarrage) | BASSE |
| ENF21 | Maintenabilité | Compatibilité multi-réseaux blockchain | L'adapter Fabric est interchangeable sans modification du domaine — un nouvel adapter blockchain peut être branché en implémentant les ports `out` existants | HAUTE |
| **Portabilité** | | | | |
| ENF23 | Portabilité | Plateformes serveur | Binaires disponibles pour Linux amd64 et Windows amd64 (produits par `make api` et `make cli`) | HAUTE |
| **Conformité légale** | | | | |
| ENF25 | Conformité | Licence du code source | Code source publié sous AGPL 3.0 ; toute contribution ou extension doit respecter cette licence | HAUTE |
| ENF26 | Conformité | Compatibilité de licence des assets | La compatibilité entre licence parent et licence dérivé est vérifiée automatiquement avant toute soumission d'un asset dérivé (`ParentID != ""`) | HAUTE |
| ENF27 | Conformité | Protection des données personnelles (RGPD) | Les données personnelles (email, profil) sont stockées uniquement en base de données locale ; aucun transfert vers un tiers sans consentement explicite | MOYENNE |
| **Fiabilité** | | | | |
| ENF28 | Fiabilité | Immuabilité des transactions blockchain | Aucune opération de suppression n'est implémentée sur la blockchain ; `ErrNotSupported` est retourné si une telle opération est tentée | HAUTE |
| ENF29 | Fiabilité | Anti-plagiat obligatoire | Tout asset de type `base` passe une vérification SHA-256 + similarité SCM > 50 % avant enregistrement — aucune exception possible | HAUTE |
| ENF30 | Fiabilité | Intégrité en cas d'échec blockchain | Si une transaction blockchain échoue, l'état local (draft) est conservé intact — aucune perte de données côté serveur | HAUTE |
| ENF31 | Fiabilité | Validation avant soumission | Toutes les données sont validées côté serveur avant soumission blockchain — les erreurs de validation ne génèrent pas de transaction partielle | HAUTE |
| ENF32 | Fiabilité | Traçabilité et intégrité du fichier source | Pour tout asset possédant un fichier ressource, l'emplacement de stockage effectif et la correspondance de son hash sont vérifiables à la demande, quelle que soit la technologie de stockage sous-jacente (adapter `out/` local, IPFS, ou autre) ; toute divergence (fichier introuvable, hash différent) est signalée explicitement, jamais masquée silencieusement | HAUTE |



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

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs](../Carte_des_specs.md)

**Application des exigences**
- **ENF01** — Temps de réponse des endpoints REST : [UCA02 (analyse)](../2-Analyse/UCA-Compte_et_Acces/UCA02.md) · [UCAUT01 (analyse)](../2-Analyse/UCAUT-Automatisation/UCAUT01.md) · [UCAUT02 (analyse)](../2-Analyse/UCAUT-Automatisation/UCAUT02.md) · [UCCL01 (analyse)](../2-Analyse/UCCL-Composant_Lecture/UCCL01.md) · [UCMOD07 (analyse)](../2-Analyse/UCMOD-Module/UCMOD07.md) · [UCPI01 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI01.md) · [roadmap_dev](../roadmap_dev.md)
- **ENF02** — Temps de soumission d'une transaction Fabric : [UCA06 (analyse)](../2-Analyse/UCA-Compte_et_Acces/UCA06.md) · [UCAUT01 (analyse)](../2-Analyse/UCAUT-Automatisation/UCAUT01.md) · [UCMOD06 (analyse)](../2-Analyse/UCMOD-Module/UCMOD06.md) · [UCPI02 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI02.md) · [roadmap_dev](../roadmap_dev.md)
- **ENF03** — Génération d'une BOM module : [Analyse_des_besoins](../2-Analyse/Analyse_des_besoins.md) · [UCPI05 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI05.md) · [UCPI09 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI09.md) · [UCREC05 (analyse)](../2-Analyse/UCREC-Recherche/UCREC05.md) · [DC_D8_Recherche](../3-Conception/DC_D8_Recherche.md) · [roadmap_dev](../roadmap_dev.md)
- **ENF05** — Disponibilité du serveur Myr (instance unique) : [UCADM01 (analyse)](../2-Analyse/UCADM-Administration/UCADM01.md) · [roadmap_dev](../roadmap_dev.md)
- **ENF06** — Disponibilité du réseau blockchain (multi-nœuds) : [UCADM04 (analyse)](../2-Analyse/UCADM-Administration/UCADM04.md) · [roadmap_dev](../roadmap_dev.md)
- **ENF07** — Reprise après redémarrage du serveur : [roadmap_dev](../roadmap_dev.md)
- **ENF08** — Chiffrement des communications : [roadmap_dev](../roadmap_dev.md)
- **ENF09** — Durée de vie des tokens de session : [roadmap_dev](../roadmap_dev.md)
- **ENF10** — Chiffrement des wallets Fabric : [Conception_intro](../3-Conception/Conception_intro.md) · [Securite](../3-Conception/Securite.md) · [roadmap_dev](../roadmap_dev.md)
- **ENF11** — Secrets absents des logs : [roadmap_dev](../roadmap_dev.md)
- **ENF12** — Contrôle d'accès par rôle : [Analyse_des_besoins](../2-Analyse/Analyse_des_besoins.md) · [UCA01 (analyse)](../2-Analyse/UCA-Compte_et_Acces/UCA01.md) · [UCA02 (analyse)](../2-Analyse/UCA-Compte_et_Acces/UCA02.md) · [UCA03 (analyse)](../2-Analyse/UCA-Compte_et_Acces/UCA03.md) · [UCA04 (analyse)](../2-Analyse/UCA-Compte_et_Acces/UCA04.md) · [UCA05 (analyse)](../2-Analyse/UCA-Compte_et_Acces/UCA05.md) · [UCA08 (analyse)](../2-Analyse/UCA-Compte_et_Acces/UCA08.md) · [UCAM01 (analyse)](../2-Analyse/UCAM-Assemblage_Module/UCAM01.md) · [UCAM02 (analyse)](../2-Analyse/UCAM-Assemblage_Module/UCAM02.md) · [UCAM03 (analyse)](../2-Analyse/UCAM-Assemblage_Module/UCAM03.md) · [UCAM05 (analyse)](../2-Analyse/UCAM-Assemblage_Module/UCAM05.md) · [UCAM07 (analyse)](../2-Analyse/UCAM-Assemblage_Module/UCAM07.md) · [UCAM08 (analyse)](../2-Analyse/UCAM-Assemblage_Module/UCAM08.md) · [UCAUT01 (analyse)](../2-Analyse/UCAUT-Automatisation/UCAUT01.md) · [UCAUT02 (analyse)](../2-Analyse/UCAUT-Automatisation/UCAUT02.md) · [UCAUT03 (analyse)](../2-Analyse/UCAUT-Automatisation/UCAUT03.md) · [UCAUT04 (analyse)](../2-Analyse/UCAUT-Automatisation/UCAUT04.md) · [UCCE01 (analyse)](../2-Analyse/UCCE-Composant_Ecriture/UCCE01.md) · [UCCE02 (analyse)](../2-Analyse/UCCE-Composant_Ecriture/UCCE02.md) · [UCCE03 (analyse)](../2-Analyse/UCCE-Composant_Ecriture/UCCE03.md) · [UCCE04 (analyse)](../2-Analyse/UCCE-Composant_Ecriture/UCCE04.md) · [UCCE05 (analyse)](../2-Analyse/UCCE-Composant_Ecriture/UCCE05.md) · [UCCE06 (analyse)](../2-Analyse/UCCE-Composant_Ecriture/UCCE06.md) · [UCCE07 (analyse)](../2-Analyse/UCCE-Composant_Ecriture/UCCE07.md) · [UCCL01 (analyse)](../2-Analyse/UCCL-Composant_Lecture/UCCL01.md) · [UCDEV01 (analyse)](../2-Analyse/UCDEV-Developpement/UCDEV01.md) · [UCMOD01 (analyse)](../2-Analyse/UCMOD-Module/UCMOD01.md) · [UCMOD02 (analyse)](../2-Analyse/UCMOD-Module/UCMOD02.md) · [UCMOD03 (analyse)](../2-Analyse/UCMOD-Module/UCMOD03.md) · [UCMOD04 (analyse)](../2-Analyse/UCMOD-Module/UCMOD04.md) · [UCMOD06 (analyse)](../2-Analyse/UCMOD-Module/UCMOD06.md) · [UCMOD07 (analyse)](../2-Analyse/UCMOD-Module/UCMOD07.md) · [UCMOD08 (analyse)](../2-Analyse/UCMOD-Module/UCMOD08.md) · [UCPI01 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI01.md) · [UCPI04 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI04.md) · [UCPI05 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI05.md) · [UCPI06 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI06.md) · [UCPI07 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI07.md) · [UCPI08 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI08.md) · [UCPI09 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI09.md) · [UCPI10 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI10.md) · [UCPI11 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI11.md) · [UCREC01 (analyse)](../2-Analyse/UCREC-Recherche/UCREC01.md) · [UCREC02 (analyse)](../2-Analyse/UCREC-Recherche/UCREC02.md) · [UCREC03 (analyse)](../2-Analyse/UCREC-Recherche/UCREC03.md) · [UCREC04 (analyse)](../2-Analyse/UCREC-Recherche/UCREC04.md) · [UCREC05 (analyse)](../2-Analyse/UCREC-Recherche/UCREC05.md) · [Securite](../3-Conception/Securite.md) · [roadmap_dev](../roadmap_dev.md)
- **ENF13** — Protection contre l'injection : [roadmap_dev](../roadmap_dev.md)
- **ENF14** — Nombre d'assets par réseau : [UCAUT02 (analyse)](../2-Analyse/UCAUT-Automatisation/UCAUT02.md) · [UCCL01 (analyse)](../2-Analyse/UCCL-Composant_Lecture/UCCL01.md) · [roadmap_dev](../roadmap_dev.md)
- **ENF15** — Taille maximale d'un fichier CAO : [UCAUT03 (analyse)](../2-Analyse/UCAUT-Automatisation/UCAUT03.md) · [roadmap_dev](../roadmap_dev.md)
- **ENF16** — Utilisateurs simultanés par instance : [UCAUT02 (analyse)](../2-Analyse/UCAUT-Automatisation/UCAUT02.md) · [roadmap_dev](../roadmap_dev.md)
- **ENF17** — Extension horizontale : [roadmap_dev](../roadmap_dev.md)
- **ENF18** — Isolation du domaine métier : [Analyse_des_besoins](../2-Analyse/Analyse_des_besoins.md) · [UCA05 (analyse)](../2-Analyse/UCA-Compte_et_Acces/UCA05.md) · [UCA06 (analyse)](../2-Analyse/UCA-Compte_et_Acces/UCA06.md) · [UCADM01 (analyse)](../2-Analyse/UCADM-Administration/UCADM01.md) · [UCADM02 (analyse)](../2-Analyse/UCADM-Administration/UCADM02.md) · [UCADM03 (analyse)](../2-Analyse/UCADM-Administration/UCADM03.md) · [UCADM04 (analyse)](../2-Analyse/UCADM-Administration/UCADM04.md) · [UCADM05 (analyse)](../2-Analyse/UCADM-Administration/UCADM05.md) · [UCADM06 (analyse)](../2-Analyse/UCADM-Administration/UCADM06.md) · [UCADM07 (analyse)](../2-Analyse/UCADM-Administration/UCADM07.md) · [UCAM01 (analyse)](../2-Analyse/UCAM-Assemblage_Module/UCAM01.md) · [UCAM03 (analyse)](../2-Analyse/UCAM-Assemblage_Module/UCAM03.md) · [UCAM05 (analyse)](../2-Analyse/UCAM-Assemblage_Module/UCAM05.md) · [UCAM07 (analyse)](../2-Analyse/UCAM-Assemblage_Module/UCAM07.md) · [UCAM08 (analyse)](../2-Analyse/UCAM-Assemblage_Module/UCAM08.md) · [UCMOD01 (analyse)](../2-Analyse/UCMOD-Module/UCMOD01.md) · [UCPI06 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI06.md) · [API_REST](../3-Conception/API_REST.md) · [Architecture_Hexagonale](../3-Conception/Architecture_Hexagonale.md) · [Conception_intro](../3-Conception/Conception_intro.md) · [DC_CLI_Admin](../3-Conception/DC_CLI_Admin.md) · [DC_D2_Administration](../3-Conception/DC_D2_Administration.md) · [Securite](../3-Conception/Securite.md) · [roadmap_dev](../roadmap_dev.md)
- **ENF19** — Couverture de tests unitaires : [roadmap_dev](../roadmap_dev.md)
- **ENF20** — Durée de déploiement d'une mise à jour : [roadmap_dev](../roadmap_dev.md)
- **ENF21** — Compatibilité multi-réseaux blockchain : [roadmap_dev](../roadmap_dev.md)
- **ENF23** — Plateformes serveur : [roadmap_dev](../roadmap_dev.md)
- **ENF25** — Licence du code source : [API_REST](../3-Conception/API_REST.md) · [Architecture_Hexagonale](../3-Conception/Architecture_Hexagonale.md) · [Conception_intro](../3-Conception/Conception_intro.md) · [roadmap_dev](../roadmap_dev.md)
- **ENF26** — Compatibilité de licence des assets : [roadmap_dev](../roadmap_dev.md)
- **ENF27** — Protection des données personnelles (RGPD) : [UCPI04 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI04.md) · [UCPI05 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI05.md) · [UCPI11 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI11.md) · [Securite](../3-Conception/Securite.md) · [roadmap_dev](../roadmap_dev.md)
- **ENF28** — Immuabilité des transactions blockchain : [UCMOD02 (analyse)](../2-Analyse/UCMOD-Module/UCMOD02.md) · [UCMOD06 (analyse)](../2-Analyse/UCMOD-Module/UCMOD06.md) · [UCMOD08 (analyse)](../2-Analyse/UCMOD-Module/UCMOD08.md) · [Securite](../3-Conception/Securite.md) · [roadmap_dev](../roadmap_dev.md)
- **ENF29** — Anti-plagiat obligatoire : [roadmap_dev](../roadmap_dev.md)
- **ENF30** — Intégrité en cas d'échec blockchain : [Analyse_des_besoins](../2-Analyse/Analyse_des_besoins.md) · [UCAM05 (analyse)](../2-Analyse/UCAM-Assemblage_Module/UCAM05.md) · [UCAUT01 (analyse)](../2-Analyse/UCAUT-Automatisation/UCAUT01.md) · [UCAUT03 (analyse)](../2-Analyse/UCAUT-Automatisation/UCAUT03.md) · [UCCE01 (analyse)](../2-Analyse/UCCE-Composant_Ecriture/UCCE01.md) · [UCCE02 (analyse)](../2-Analyse/UCCE-Composant_Ecriture/UCCE02.md) · [UCCE03 (analyse)](../2-Analyse/UCCE-Composant_Ecriture/UCCE03.md) · [UCCE04 (analyse)](../2-Analyse/UCCE-Composant_Ecriture/UCCE04.md) · [UCCE05 (analyse)](../2-Analyse/UCCE-Composant_Ecriture/UCCE05.md) · [UCMOD03 (analyse)](../2-Analyse/UCMOD-Module/UCMOD03.md) · [UCMOD06 (analyse)](../2-Analyse/UCMOD-Module/UCMOD06.md) · [UCPI01 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI01.md) · [UCPI02 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI02.md) · [UCPI07 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI07.md) · [UCPI08 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI08.md) · [UCPI09 (analyse)](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI09.md) · [Conception_intro](../3-Conception/Conception_intro.md) · [roadmap_dev](../roadmap_dev.md)
- **ENF31** — Validation avant soumission : [UCMOD01 (analyse)](../2-Analyse/UCMOD-Module/UCMOD01.md) · [UCMOD06 (analyse)](../2-Analyse/UCMOD-Module/UCMOD06.md) · [roadmap_dev](../roadmap_dev.md)
- **ENF32** — Traçabilité et intégrité du fichier source

**Cité par**
- [Expression_des_besoins_Intro](Expression_des_besoins_Intro.md)
- [Conception_intro](../3-Conception/Conception_intro.md)
- [roadmap_dev](../roadmap_dev.md)

<!-- liens-obsidian:end -->
