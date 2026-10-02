---
tags:
  - couche/expression
  - type/tracabilite
  - relecture/incoherence
---
# Matrice de traçabilité — Exigences fonctionnelles ↔ Use Cases

Ce document croise les exigences fonctionnelles (EF) du projet Myr avec les use cases (UC) qui les couvrent. Il permet de vérifier qu'aucune exigence n'est orpheline et qu'aucun UC n'est sans motivation fonctionnelle.

---

## 1. Exigences fonctionnelles et UC couvrant

Les exigences sont dérivées des objectifs du projet (section 2.1), des contraintes (section 2.3) et de l'ensemble des use cases.


| ID                                               | Exigence fonctionnelle                                                  | UC couvrant               |
| ------------------------------------------------ | ----------------------------------------------------------------------- | ------------------------- |
| **Compte et Accès**                             |                                                                         |                           |
| EF01                                             | Permettre à un visiteur de créer un compte sur un réseau             | UCA01                     |
| EF02                                             | Authentifier une identité (enrôlement CA + token de session opaque)   | UCA02                     |
| EF03                                             | Déconnecter un utilisateur et invalider sa session                     | UCA03                     |
| EF04                                             | Vérifier la validité d'une session active                             | UCA04                     |
| EF05                                             | Contrôler les accès selon le rôle attribué                          | UCA05, UCA07              |
| EF06                                             | Consulter les assets possédés par l'utilisateur                       | UCA06                     |
| **Administration**                               |                                                                         |                           |
| EF07                                             | Créer un réseau blockchain indépendant                               | UCADM02                   |
| EF08                                             | Gérer les organisations membres d'un réseau (ajout, mise à jour)     | UCADM01                   |
| EF09                                             | Étendre un réseau avec de nouveaux nœuds (peer ou orderer)           | UCADM03                   |
| EF57                                             | Retirer administrativement un nœud d'un réseau existant               | UCADM04                   |
| EF58                                             | Démanteler un réseau de test (CLI uniquement — jamais via REST)        | UCADM05                   |
| **Composants — Création et édition**          |                                                                         |                           |
| EF10                                             | Enregistrer un composant physique sur la blockchain                     | UCCE01                    |
| EF11                                             | Enregistrer un composant numérique sur la blockchain                   | UCCE03                    |
| EF12                                             | Configurer les métadonnées d'un composant                             | UCCE02                    |
| EF13                                             | Faire évoluer un composant (amélioration, dérivation, extension)     | UCCE04, UCCE05            |
| EF14                                             | Ajouter une interface à un composant existant                          | UCCE06                    |
| EF15                                             | Vérifier l'unicité d'un composant (anti-plagiat SHA-256 + SCM > 50 %) | UCCE01                    |
| EF16                                             | Vérifier la compatibilité de licence lors d'une dérivation           | UCCE04, UCMOD01, UCMOD06  |
| EF60                                             | Supprimer un composant (masquage local des listes, ledger jamais modifié) | UCCE07                |
| **Composants — Lecture**                        |                                                                         |                           |
| EF17                                             | Rechercher et filtrer les composants disponibles sur le réseau         | UCCL01                    |
| **Composition (instances)**                       |                                                                         |                           |
| EF18                                             | Créer des liaisons entre interfaces compatibles de composants          | UCAM01                    |
| EF19                                             | Visualiser les interfaces physiques d'un composant                      | UCAM02                    |
| EF20                                             | Définir une interface sur un composant                                 | UCAM03                    |
| EF21                                             | *(retiré)* — sélection multiple/jauge de progression pour l'ajout d'instances : ergonomie 100 % frontend, sans logique domaine propre (le placement unitaire `AddAssetToWorkspace` reste couvert par EF26/UCMOD01 et UCAM05) | —                         |
| EF22                                             | Transformer un composant en module (découpage en sous-systèmes)       | UCAM05                    |
| EF23                                             | Choisir un asset d'accroche (fastener) pour une liaison                 | UCAM07                    |
| EF24                                             | Retirer une instance de composant d'un module (avec cascade des connexions) | UCAM08                    |
| EF25                                             | Garantir un slot virtuel disponible sur chaque asset                    | UCAM03                    |
| EF64                                             | Proposer un découpage automatique (sous-pièces + connexions candidates) d'un composant STEP en amont d'une transformation composant → module | UCAM09 |
| **Modules**                                      |                                                                         |                           |
| EF26                                             | Assembler plusieurs composants en module (état draft), y compris par dérivation d'un module existant (composition dupliquée depuis un `parent_id`) | UCMOD01 |
| EF27                                             | Soumettre un module à la blockchain (ModuleVersion immuable)           | UCMOD06                   |
| EF28                                             | Ajouter un module existant à l'espace de travail                       | UCMOD02 — #incoherence : cette ligne référençait aussi UCMOD03 et UCMOD05 ; UCMOD03 ne couvre pas cet EF (il documente la modification des métadonnées d'un module, sans rapport avec l'ajout d'un module existant comme instance) et UCMOD05 ne correspond à aucun fichier existant dans `specs/1-Expression/UCMOD-Module/` ni `specs/2-Analyse/UCMOD-Module/` — origine de ces deux références à clarifier avant de les retirer définitivement |
| EF29                                             | Visualiser la composition d'un module                                   | UCMOD04                   |
| EF61                                             | Modifier les métadonnées d'un module (nom, description, licence, tags, liens) | UCMOD03            |
| EF62                                             | Supprimer un module (masquage local des listes, ledger jamais modifié) | UCMOD08                   |
| EF63                                             | Lister ses modules en brouillon, filtrés par propriétaire et par statut | UCMOD07                   |
| **Propriété Intellectuelle et Rémunération** |                                                                         |                           |
| EF30                                             | Commander un module complet (fabrication ou achat en stock)             | UCPI01                    |
| EF31                                             | Distribuer automatiquement les commissions aux auteurs à la livraison  | UCPI02, UCAUT01           |
| EF32                                             | Définir un prix sur un composant propriétaire                         | UCPI04                    |
| EF33                                             | Définir un prix sur un module propriétaire                            | UCPI05                    |
| EF34                                             | Signaler un composant similaire à un existant                          | UCPI06                    |
| EF35                                             | Transférer la propriété intellectuelle d'un asset                    | UCPI07                    |
| EF36                                             | Cloner un composant sur un réseau externe                              | UCPI08                    |
| EF37                                             | Cloner un module sur un réseau externe                                 | UCPI09                    |
| EF38                                             | Respecter et vérifier une norme d'écoconception                       | UCPI10                    |
| **Recherche**                                    |                                                                         |                           |
| EF39                                             | Rechercher un asset par référence ou filtre                           | UCREC01                   |
| EF40                                             | Identifier les composants compatibles entre eux (interfaces)            | UCREC02                   |
| EF41                                             | Consulter l'arbre de versions d'un composant                            | UCREC03                   |
| EF42                                             | Identifier tous les modules qui intègrent un composant donné          | UCREC04                   |
| EF43                                             | Exporter la BOM (Bill Of Materials) d'un module                         | UCREC05                   |
| **Automatisation**                               |                                                                         |                           |
| EF44                                             | Automatiser la fabrication et la livraison d'un composant               | UCAUT01                   |
| EF45                                             | Automatiser la commande en ligne d'un asset                             | UCAUT02                   |
| EF46                                             | Intégrer un modèle 3D depuis un logiciel CAO (plugin)                 | UCAUT03                   |
| EF47                                             | Gérer les versions SCM d'un modèle 3D                                 | UCAUT04                   |
| **Interface Graphique, Paramètres, Documentation** | EF48–EF54 : retirés — use cases 100 % frontend (UCIG, UCPAR, UCDOC), sans contrat REST propre à `myr`, désormais spécifiés dans le dépôt GUI externe | —                         |
| **Développement autour de MYR**                 |                                                                         |                           |
| EF55                                             | Exposer une API REST pour les intégrations tierces                     | UCDEV01                   |
| EF56                                             | Proposer un CLI d'administration serveur                                | UCDEV02                   |
| EF59                                             | Modifier le prix d'un asset (effet commandes futures uniquement)        | UCPI11                    |

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs](../Carte_des_specs.md)

**Use cases cités**
- UCA01 — Création d'un compte : [expression](UCA-Compte_et_Acces/UCA01.md) · [analyse](../2-Analyse/UCA-Compte_et_Acces/UCA01.md)
- UCA02 — Se Connecter : [expression](UCA-Compte_et_Acces/UCA02.md) · [analyse](../2-Analyse/UCA-Compte_et_Acces/UCA02.md)
- UCA03 — Se Déconnecter : [expression](UCA-Compte_et_Acces/UCA03.md) · [analyse](../2-Analyse/UCA-Compte_et_Acces/UCA03.md)
- UCA04 — Vérification de la connexion : [expression](UCA-Compte_et_Acces/UCA04.md) · [analyse](../2-Analyse/UCA-Compte_et_Acces/UCA04.md)
- UCA05 — Vérification des accès du rôle attribué : [expression](UCA-Compte_et_Acces/UCA05.md) · [analyse](../2-Analyse/UCA-Compte_et_Acces/UCA05.md)
- UCA06 — Vérifier les possessions : [expression](UCA-Compte_et_Acces/UCA06.md) · [analyse](../2-Analyse/UCA-Compte_et_Acces/UCA06.md)
- UCA07 — Vérification du rôle attribué : [expression](UCA-Compte_et_Acces/UCA07.md) · [analyse](../2-Analyse/UCA-Compte_et_Acces/UCA07.md)
- UCADM01 — Ajouter une organisation au réseau : [expression](UCADM-Administration/UCADM01.md) · [analyse](../2-Analyse/UCADM-Administration/UCADM01.md)
- UCADM02 — Créer un réseau indépendant : [expression](UCADM-Administration/UCADM02.md) · [analyse](../2-Analyse/UCADM-Administration/UCADM02.md)
- UCADM03 — Ajouter un nœud à un réseau existant : [expression](UCADM-Administration/UCADM03.md) · [analyse](../2-Analyse/UCADM-Administration/UCADM03.md)
- UCADM04 — Retirer un nœud d'un réseau existant : [expression](UCADM-Administration/UCADM04.md) · [analyse](../2-Analyse/UCADM-Administration/UCADM04.md)
- UCADM05 — Démanteler un réseau (dev/test uniquement) : [expression](UCADM-Administration/UCADM05.md) · [analyse](../2-Analyse/UCADM-Administration/UCADM05.md)
- UCAM01 — Liaison entre interfaces : [expression](UCAM-Assemblage_Module/UCAM01.md) · [analyse](../2-Analyse/UCAM-Assemblage_Module/UCAM01.md)
- UCAM02 — Visualiser les interfaces physiques de composants : [expression](UCAM-Assemblage_Module/UCAM02.md) · [analyse](../2-Analyse/UCAM-Assemblage_Module/UCAM02.md)
- UCAM03 — Créer une interface sur un composant : [expression](UCAM-Assemblage_Module/UCAM03.md) · [analyse](../2-Analyse/UCAM-Assemblage_Module/UCAM03.md)
- UCAM05 — Transformation d'un composant en module : [expression](UCAM-Assemblage_Module/UCAM05.md) · [analyse](../2-Analyse/UCAM-Assemblage_Module/UCAM05.md)
- UCAM07 — Choisir un asset d'accroche (Fastener) : [expression](UCAM-Assemblage_Module/UCAM07.md) · [analyse](../2-Analyse/UCAM-Assemblage_Module/UCAM07.md)
- UCAM08 — Retirer une instance de composant d'un Module : [expression](UCAM-Assemblage_Module/UCAM08.md) · [analyse](../2-Analyse/UCAM-Assemblage_Module/UCAM08.md)
- UCAM09 — Décomposition assistée d'un composant assemblage : [expression](UCAM-Assemblage_Module/UCAM09.md)
- UCAUT01 — Fabrication/Livraison d'un Composant : [expression](UCAUT-Automatisation/UCAUT01.md) · [analyse](../2-Analyse/UCAUT-Automatisation/UCAUT01.md)
- UCAUT02 — Commande en ligne de Asset : [expression](UCAUT-Automatisation/UCAUT02.md) · [analyse](../2-Analyse/UCAUT-Automatisation/UCAUT02.md)
- UCAUT03 — Ajouter un modèle 3D depuis un logiciel CAO : [expression](UCAUT-Automatisation/UCAUT03.md) · [analyse](../2-Analyse/UCAUT-Automatisation/UCAUT03.md)
- UCAUT04 — Gestion SCM d'un modèle 3D : [expression](UCAUT-Automatisation/UCAUT04.md) · [analyse](../2-Analyse/UCAUT-Automatisation/UCAUT04.md)
- UCCE01 — Ajout d'un composant Physique : [expression](UCCE-Composant_Ecriture/UCCE01.md) · [analyse](../2-Analyse/UCCE-Composant_Ecriture/UCCE01.md)
- UCCE02 — Configurer un Composant : [expression](UCCE-Composant_Ecriture/UCCE02.md) · [analyse](../2-Analyse/UCCE-Composant_Ecriture/UCCE02.md)
- UCCE03 — Ajout d'un composant Numérique : [expression](UCCE-Composant_Ecriture/UCCE03.md) · [analyse](../2-Analyse/UCCE-Composant_Ecriture/UCCE03.md)
- UCCE04 — Améliorer un Composant : [expression](UCCE-Composant_Ecriture/UCCE04.md) · [analyse](../2-Analyse/UCCE-Composant_Ecriture/UCCE04.md)
- UCCE05 — Créer une extension de Composant : [expression](UCCE-Composant_Ecriture/UCCE05.md) · [analyse](../2-Analyse/UCCE-Composant_Ecriture/UCCE05.md)
- UCCE06 — Ajouter une interface à un Composant déjà créé : [expression](UCCE-Composant_Ecriture/UCCE06.md) · [analyse](../2-Analyse/UCCE-Composant_Ecriture/UCCE06.md)
- UCCE07 — Supprimer un Composant : [expression](UCCE-Composant_Ecriture/UCCE07.md) · [analyse](../2-Analyse/UCCE-Composant_Ecriture/UCCE07.md)
- UCCL01 — Faire une recherche par filtre : [expression](UCCL-Composant_Lecture/UCCL01.md) · [analyse](../2-Analyse/UCCL-Composant_Lecture/UCCL01.md)
- UCDEV01 — Utilisation de l'API : [expression](UCDEV-Developpement/UCDEV01.md) · [analyse](../2-Analyse/UCDEV-Developpement/UCDEV01.md)
- UCDEV02 — Utilisation du CLI : [expression](UCDEV-Developpement/UCDEV02.md) · [analyse](../2-Analyse/UCDEV-Developpement/UCDEV02.md)
- UCMOD01 — Créer un Module : [expression](UCMOD-Module/UCMOD01.md) · [analyse](../2-Analyse/UCMOD-Module/UCMOD01.md)
- UCMOD02 — Ajouter un Module existant : [expression](UCMOD-Module/UCMOD02.md) · [analyse](../2-Analyse/UCMOD-Module/UCMOD02.md)
- UCMOD03 — Modifier les métadonnées d'un Module : [expression](UCMOD-Module/UCMOD03.md) · [analyse](../2-Analyse/UCMOD-Module/UCMOD03.md)
- UCMOD04 — Visualiser les composants d'un Module : [expression](UCMOD-Module/UCMOD04.md) · [analyse](../2-Analyse/UCMOD-Module/UCMOD04.md)
- UCMOD06 — Soumettre un module à la blockchain : [expression](UCMOD-Module/UCMOD06.md) · [analyse](../2-Analyse/UCMOD-Module/UCMOD06.md)
- UCMOD07 — Lister ses Modules en brouillon : [expression](UCMOD-Module/UCMOD07.md) · [analyse](../2-Analyse/UCMOD-Module/UCMOD07.md)
- UCMOD08 — Supprimer un Module : [expression](UCMOD-Module/UCMOD08.md) · [analyse](../2-Analyse/UCMOD-Module/UCMOD08.md)
- UCPI01 — Commander un Module complet : [expression](UCPI-Propriete_Intellectuelle/UCPI01.md) · [analyse](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI01.md)
- UCPI02 — Recevoir une commission sur l'utilisation d'un Module : [expression](UCPI-Propriete_Intellectuelle/UCPI02.md) · [analyse](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI02.md)
- UCPI04 — Définir un prix sur un Composant proprietaire : [expression](UCPI-Propriete_Intellectuelle/UCPI04.md) · [analyse](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI04.md)
- UCPI05 — Définir un prix sur un Module proprietaire : [expression](UCPI-Propriete_Intellectuelle/UCPI05.md) · [analyse](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI05.md)
- UCPI06 — Déclarer un composant similaire : [expression](UCPI-Propriete_Intellectuelle/UCPI06.md) · [analyse](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI06.md)
- UCPI07 — Transfert de propriété intellectuelle : [expression](UCPI-Propriete_Intellectuelle/UCPI07.md) · [analyse](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI07.md)
- UCPI08 — Cloner un Composant sur un réseau exterieur : [expression](UCPI-Propriete_Intellectuelle/UCPI08.md) · [analyse](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI08.md)
- UCPI09 — Cloner un Module sur un réseau exterieur : [expression](UCPI-Propriete_Intellectuelle/UCPI09.md) · [analyse](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI09.md)
- UCPI10 — Norme de conception écoconception : [expression](UCPI-Propriete_Intellectuelle/UCPI10.md) · [analyse](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI10.md)
- UCPI11 — Modifier le prix d'un asset : [expression](UCPI-Propriete_Intellectuelle/UCPI11.md) · [analyse](../2-Analyse/UCPI-Propriete_Intellectuelle/UCPI11.md)
- UCREC01 — Rechercher une référence existante : [expression](UCREC-Recherche/UCREC01.md) · [analyse](../2-Analyse/UCREC-Recherche/UCREC01.md)
- UCREC02 — Rechercher les Composants compatibles : [expression](UCREC-Recherche/UCREC02.md) · [analyse](../2-Analyse/UCREC-Recherche/UCREC02.md)
- UCREC03 — Rechercher les versions des Composants : [expression](UCREC-Recherche/UCREC03.md) · [analyse](../2-Analyse/UCREC-Recherche/UCREC03.md)
- UCREC04 — Rechercher les Modules qui utilisent un Composant : [expression](UCREC-Recherche/UCREC04.md) · [analyse](../2-Analyse/UCREC-Recherche/UCREC04.md)
- UCREC05 — Exporter BOM Module : [expression](UCREC-Recherche/UCREC05.md) · [analyse](../2-Analyse/UCREC-Recherche/UCREC05.md)

**Cité par**
- [Expression_des_besoins_Intro](Expression_des_besoins_Intro.md)
- [roadmap_dev](../roadmap_dev.md)

<!-- liens-obsidian:end -->
