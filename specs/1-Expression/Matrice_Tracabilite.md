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
