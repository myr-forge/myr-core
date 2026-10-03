---
tags:
  - couche/transverse
  - type/carte
---
# Carte des specs

> Note de traçabilité Obsidian générée automatiquement à partir des specs et du code — ne pas l'éditer à la main.

Point d'entrée pour parcourir les specs de `myr-core` : chaque use case y est relié à ses deux couches (expression du besoin, analyse), aux documents de conception qui le couvrent et aux domaines Go concernés. Chaque fichier de `specs/` se termine par une section **Liens** (générée) et porte des tags dans son frontmatter.

## Comment naviguer

| Tag | Regroupe |
|-----|----------|
| `couche/expression` · `couche/analyse` · `couche/conception` · `couche/code` | les documents d'une couche (les notes de code vivent dans `docs/code/`) |
| `famille/UCMOD` … | les use cases d'une famille |
| `uc/UCMOD04` … | tout ce qui concerne un use case : ses deux couches, la conception qui le déclare dans son périmètre, les packages Go qui le réalisent |
| `rm/RM13` · `enf/ENF12` … | les use cases, documents de conception et packages Go (tests) qui appliquent une règle métier ou une exigence non fonctionnelle |
| `domaine/model` … | specs et code d'un domaine Go (`domain/<domaine>/`) |
| `hexa/domaine` · `hexa/adapter-in` · `hexa/adapter-out` · `hexa/point-entree` | les packages Go d'une couche de l'architecture hexagonale |
| `relecture/question` · `relecture/remarque` · `relecture/incoherence` · `relecture/ecart` | les documents qui portent une annotation de relecture ouverte |

Vues transverses : [Expression_des_besoins_Intro](1-Expression/Expression_des_besoins_Intro.md) · [Regles_Metier](1-Expression/Regles_Metier.md) · [Exigences_Non_Fonctionnelles](1-Expression/Exigences_Non_Fonctionnelles.md) · [Matrice_Tracabilite](1-Expression/Matrice_Tracabilite.md) · [Analyse_des_besoins](2-Analyse/Analyse_des_besoins.md) · [Conception_intro](3-Conception/Conception_intro.md) · [Architecture_Hexagonale](3-Conception/Architecture_Hexagonale.md) · [Architecture du code](../docs/code/Architecture.md) · [Traçabilité use cases ↔ code](../docs/code/Tracabilite_UC_Code.md)

## Use cases

### UCA — Compte et Acces

Domaines : `identity`, `role`

| UC | Titre | Expression | Analyse | Conception |
|----|-------|------------|---------|------------|
| UCA01 | Création d'un compte | [UCA01 (expression)](1-Expression/UCA-Compte_et_Acces/UCA01.md) | [UCA01 (analyse)](2-Analyse/UCA-Compte_et_Acces/UCA01.md) | [API_REST](3-Conception/API_REST.md) · [DC_CLI_Identity](3-Conception/DC_CLI_Identity.md) · [DC_D1_Auth_Identity](3-Conception/DC_D1_Auth_Identity.md) |
| UCA02 | Se Connecter | [UCA02 (expression)](1-Expression/UCA-Compte_et_Acces/UCA02.md) | [UCA02 (analyse)](2-Analyse/UCA-Compte_et_Acces/UCA02.md) | [DC_CLI_Identity](3-Conception/DC_CLI_Identity.md) |
| UCA03 | Se Déconnecter | [UCA03 (expression)](1-Expression/UCA-Compte_et_Acces/UCA03.md) | [UCA03 (analyse)](2-Analyse/UCA-Compte_et_Acces/UCA03.md) | — |
| UCA04 | Vérification de la connexion | [UCA04 (expression)](1-Expression/UCA-Compte_et_Acces/UCA04.md) | [UCA04 (analyse)](2-Analyse/UCA-Compte_et_Acces/UCA04.md) | [DC_CLI_Identity](3-Conception/DC_CLI_Identity.md) |
| UCA05 | Vérification des accès du rôle attribué | [UCA05 (expression)](1-Expression/UCA-Compte_et_Acces/UCA05.md) | [UCA05 (analyse)](2-Analyse/UCA-Compte_et_Acces/UCA05.md) | — |
| UCA06 | Vérifier les possessions | [UCA06 (expression)](1-Expression/UCA-Compte_et_Acces/UCA06.md) | [UCA06 (analyse)](2-Analyse/UCA-Compte_et_Acces/UCA06.md) | [DC_D8_Recherche](3-Conception/DC_D8_Recherche.md) · [Securite](3-Conception/Securite.md) |
| UCA07 | Vérification du rôle attribué | [UCA07 (expression)](1-Expression/UCA-Compte_et_Acces/UCA07.md) | [UCA07 (analyse)](2-Analyse/UCA-Compte_et_Acces/UCA07.md) | [DC_CLI_Identity](3-Conception/DC_CLI_Identity.md) |
| UCA08 | Demander un rôle | [UCA08 (expression)](1-Expression/UCA-Compte_et_Acces/UCA08.md) | [UCA08 (analyse)](2-Analyse/UCA-Compte_et_Acces/UCA08.md) | [DC_CLI_Identity](3-Conception/DC_CLI_Identity.md) |

### UCADM — Administration

Domaines : `channel`, `identity`, `network`

| UC | Titre | Expression | Analyse | Conception |
|----|-------|------------|---------|------------|
| UCADM01 | Ajouter une organisation au réseau | [UCADM01 (expression)](1-Expression/UCADM-Administration/UCADM01.md) | [UCADM01 (analyse)](2-Analyse/UCADM-Administration/UCADM01.md) | [API_REST](3-Conception/API_REST.md) · [DC_CLI_Admin](3-Conception/DC_CLI_Admin.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D2_Administration](3-Conception/DC_D2_Administration.md) · [DC_D7_Payment](3-Conception/DC_D7_Payment.md) |
| UCADM02 | Créer un réseau indépendant | [UCADM02 (expression)](1-Expression/UCADM-Administration/UCADM02.md) | [UCADM02 (analyse)](2-Analyse/UCADM-Administration/UCADM02.md) | [DC_CLI_Admin](3-Conception/DC_CLI_Admin.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D2_Administration](3-Conception/DC_D2_Administration.md) |
| UCADM03 | Ajouter un nœud à un réseau existant | [UCADM03 (expression)](1-Expression/UCADM-Administration/UCADM03.md) | [UCADM03 (analyse)](2-Analyse/UCADM-Administration/UCADM03.md) | [API_REST](3-Conception/API_REST.md) · [DC_CLI_Admin](3-Conception/DC_CLI_Admin.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D2_Administration](3-Conception/DC_D2_Administration.md) |
| UCADM04 | Retirer un nœud d'un réseau existant | [UCADM04 (expression)](1-Expression/UCADM-Administration/UCADM04.md) | [UCADM04 (analyse)](2-Analyse/UCADM-Administration/UCADM04.md) | [API_REST](3-Conception/API_REST.md) · [DC_CLI_Admin](3-Conception/DC_CLI_Admin.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D2_Administration](3-Conception/DC_D2_Administration.md) |
| UCADM05 | Démanteler un réseau (dev/test uniquement) | [UCADM05 (expression)](1-Expression/UCADM-Administration/UCADM05.md) | [UCADM05 (analyse)](2-Analyse/UCADM-Administration/UCADM05.md) | [DC_CLI_Admin](3-Conception/DC_CLI_Admin.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D2_Administration](3-Conception/DC_D2_Administration.md) |
| UCADM06 | Attribuer des rôles à une organisation | [UCADM06 (expression)](1-Expression/UCADM-Administration/UCADM06.md) | [UCADM06 (analyse)](2-Analyse/UCADM-Administration/UCADM06.md) | [DC_CLI_Admin](3-Conception/DC_CLI_Admin.md) · [DC_D2_Administration](3-Conception/DC_D2_Administration.md) |
| UCADM07 | Gérer les rôles | [UCADM07 (expression)](1-Expression/UCADM-Administration/UCADM07.md) | [UCADM07 (analyse)](2-Analyse/UCADM-Administration/UCADM07.md) | [DC_CLI_Admin](3-Conception/DC_CLI_Admin.md) · [DC_D2_Administration](3-Conception/DC_D2_Administration.md) |

### UCAM — Assemblage Module

Domaines : `model`

| UC | Titre | Expression | Analyse | Conception |
|----|-------|------------|---------|------------|
| UCAM01 | Liaison entre interfaces | [UCAM01 (expression)](1-Expression/UCAM-Assemblage_Module/UCAM01.md) | [UCAM01 (analyse)](2-Analyse/UCAM-Assemblage_Module/UCAM01.md) | [Architecture_Composition](3-Conception/Architecture_Composition.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [Sequence_soumission_module](3-Conception/Sequence_soumission_module.md) |
| UCAM02 | Visualiser les interfaces physiques de composants | [UCAM02 (expression)](1-Expression/UCAM-Assemblage_Module/UCAM02.md) | [UCAM02 (analyse)](2-Analyse/UCAM-Assemblage_Module/UCAM02.md) | [API_REST](3-Conception/API_REST.md) · [Architecture_Composition](3-Conception/Architecture_Composition.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) |
| UCAM03 | Créer une interface sur un composant | [UCAM03 (expression)](1-Expression/UCAM-Assemblage_Module/UCAM03.md) | [UCAM03 (analyse)](2-Analyse/UCAM-Assemblage_Module/UCAM03.md) | [Architecture_Composition](3-Conception/Architecture_Composition.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [Sequence_soumission_asset](3-Conception/Sequence_soumission_asset.md) |
| UCAM05 | Transformation d'un composant en module | [UCAM05 (expression)](1-Expression/UCAM-Assemblage_Module/UCAM05.md) | [UCAM05 (analyse)](2-Analyse/UCAM-Assemblage_Module/UCAM05.md) | [API_REST](3-Conception/API_REST.md) · [Architecture_Composition](3-Conception/Architecture_Composition.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) |
| UCAM07 | Choisir un asset d'accroche (Fastener) | [UCAM07 (expression)](1-Expression/UCAM-Assemblage_Module/UCAM07.md) | [UCAM07 (analyse)](2-Analyse/UCAM-Assemblage_Module/UCAM07.md) | [Architecture_Composition](3-Conception/Architecture_Composition.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [Sequence_soumission_module](3-Conception/Sequence_soumission_module.md) |
| UCAM08 | Retirer une instance de composant d'un Module | [UCAM08 (expression)](1-Expression/UCAM-Assemblage_Module/UCAM08.md) | [UCAM08 (analyse)](2-Analyse/UCAM-Assemblage_Module/UCAM08.md) | [API_REST](3-Conception/API_REST.md) · [Architecture_Composition](3-Conception/Architecture_Composition.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [Sequence_soumission_module](3-Conception/Sequence_soumission_module.md) |
| UCAM09 | Décomposition assistée d'un composant assemblage | [UCAM09 (expression)](1-Expression/UCAM-Assemblage_Module/UCAM09.md) | — | [API_REST](3-Conception/API_REST.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [Modele_Domaine](3-Conception/Modele_Domaine.md) |

### UCAUT — Automatisation

Domaines : `model`, `payment`, `role`

| UC | Titre | Expression | Analyse | Conception |
|----|-------|------------|---------|------------|
| UCAUT01 | Fabrication/Livraison d'un Composant | [UCAUT01 (expression)](1-Expression/UCAUT-Automatisation/UCAUT01.md) | [UCAUT01 (analyse)](2-Analyse/UCAUT-Automatisation/UCAUT01.md) | [API_REST](3-Conception/API_REST.md) · [Chaincode](3-Conception/Chaincode.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D7_Payment](3-Conception/DC_D7_Payment.md) · [DC_D9_Automatisation](3-Conception/DC_D9_Automatisation.md) · [Modele_Domaine](3-Conception/Modele_Domaine.md) |
| UCAUT02 | Commande en ligne de Asset | [UCAUT02 (expression)](1-Expression/UCAUT-Automatisation/UCAUT02.md) | [UCAUT02 (analyse)](2-Analyse/UCAUT-Automatisation/UCAUT02.md) | [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D9_Automatisation](3-Conception/DC_D9_Automatisation.md) |
| UCAUT03 | Ajouter un modèle 3D depuis un logiciel CAO | [UCAUT03 (expression)](1-Expression/UCAUT-Automatisation/UCAUT03.md) | [UCAUT03 (analyse)](2-Analyse/UCAUT-Automatisation/UCAUT03.md) | [API_REST](3-Conception/API_REST.md) · [DC_D9_Automatisation](3-Conception/DC_D9_Automatisation.md) |
| UCAUT04 | Gestion SCM d'un modèle 3D | [UCAUT04 (expression)](1-Expression/UCAUT-Automatisation/UCAUT04.md) | [UCAUT04 (analyse)](2-Analyse/UCAUT-Automatisation/UCAUT04.md) | [API_REST](3-Conception/API_REST.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D9_Automatisation](3-Conception/DC_D9_Automatisation.md) |

### UCCE — Composant Ecriture

Domaines : `model`

| UC | Titre | Expression | Analyse | Conception |
|----|-------|------------|---------|------------|
| UCCE01 | Ajout d'un composant Physique | [UCCE01 (expression)](1-Expression/UCCE-Composant_Ecriture/UCCE01.md) | [UCCE01 (analyse)](2-Analyse/UCCE-Composant_Ecriture/UCCE01.md) | [API_REST](3-Conception/API_REST.md) · [Architecture_Composition](3-Conception/Architecture_Composition.md) · [Chaincode](3-Conception/Chaincode.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [Sequence_soumission_asset](3-Conception/Sequence_soumission_asset.md) |
| UCCE02 | Configurer un Composant | [UCCE02 (expression)](1-Expression/UCCE-Composant_Ecriture/UCCE02.md) | [UCCE02 (analyse)](2-Analyse/UCCE-Composant_Ecriture/UCCE02.md) | [API_REST](3-Conception/API_REST.md) · [Architecture_Composition](3-Conception/Architecture_Composition.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) |
| UCCE03 | Ajout d'un composant Numérique | [UCCE03 (expression)](1-Expression/UCCE-Composant_Ecriture/UCCE03.md) | [UCCE03 (analyse)](2-Analyse/UCCE-Composant_Ecriture/UCCE03.md) | [API_REST](3-Conception/API_REST.md) · [Architecture_Composition](3-Conception/Architecture_Composition.md) · [Chaincode](3-Conception/Chaincode.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) |
| UCCE04 | Améliorer un Composant | [UCCE04 (expression)](1-Expression/UCCE-Composant_Ecriture/UCCE04.md) | [UCCE04 (analyse)](2-Analyse/UCCE-Composant_Ecriture/UCCE04.md) | [Architecture_Composition](3-Conception/Architecture_Composition.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) |
| UCCE05 | Créer une extension de Composant | [UCCE05 (expression)](1-Expression/UCCE-Composant_Ecriture/UCCE05.md) | [UCCE05 (analyse)](2-Analyse/UCCE-Composant_Ecriture/UCCE05.md) | [Architecture_Composition](3-Conception/Architecture_Composition.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) |
| UCCE06 | Ajouter une interface à un Composant déjà créé | [UCCE06 (expression)](1-Expression/UCCE-Composant_Ecriture/UCCE06.md) | [UCCE06 (analyse)](2-Analyse/UCCE-Composant_Ecriture/UCCE06.md) | [Architecture_Composition](3-Conception/Architecture_Composition.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [Sequence_soumission_asset](3-Conception/Sequence_soumission_asset.md) |
| UCCE07 | Supprimer un Composant | [UCCE07 (expression)](1-Expression/UCCE-Composant_Ecriture/UCCE07.md) | [UCCE07 (analyse)](2-Analyse/UCCE-Composant_Ecriture/UCCE07.md) | [API_REST](3-Conception/API_REST.md) |

### UCCL — Composant Lecture

Domaines : `model`

| UC | Titre | Expression | Analyse | Conception |
|----|-------|------------|---------|------------|
| UCCL01 | Faire une recherche par filtre | [UCCL01 (expression)](1-Expression/UCCL-Composant_Lecture/UCCL01.md) | [UCCL01 (analyse)](2-Analyse/UCCL-Composant_Lecture/UCCL01.md) | [API_REST](3-Conception/API_REST.md) · [Chaincode](3-Conception/Chaincode.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D1_Auth_Identity](3-Conception/DC_D1_Auth_Identity.md) |
| UCCL03 | Vérifier la validité des emplacements externes d'un composant ou module | [UCCL03 (expression)](1-Expression/UCCL-Composant_Lecture/UCCL03.md) | — | [Architecture_Composition](3-Conception/Architecture_Composition.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [Modele_Domaine](3-Conception/Modele_Domaine.md) |

### UCDEV — Developpement

Domaines : —

| UC | Titre | Expression | Analyse | Conception |
|----|-------|------------|---------|------------|
| UCDEV01 | Utilisation de l'API | [UCDEV01 (expression)](1-Expression/UCDEV-Developpement/UCDEV01.md) | [UCDEV01 (analyse)](2-Analyse/UCDEV-Developpement/UCDEV01.md) | — |
| UCDEV02 | Utilisation du CLI | [UCDEV02 (expression)](1-Expression/UCDEV-Developpement/UCDEV02.md) | [UCDEV02 (analyse)](2-Analyse/UCDEV-Developpement/UCDEV02.md) | [DC_CLI_Admin](3-Conception/DC_CLI_Admin.md) |

### UCMOD — Module

Domaines : `model`

| UC | Titre | Expression | Analyse | Conception |
|----|-------|------------|---------|------------|
| UCMOD01 | Créer un Module | [UCMOD01 (expression)](1-Expression/UCMOD-Module/UCMOD01.md) | [UCMOD01 (analyse)](2-Analyse/UCMOD-Module/UCMOD01.md) | [API_REST](3-Conception/API_REST.md) · [Architecture_Composition](3-Conception/Architecture_Composition.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [Sequence_soumission_module](3-Conception/Sequence_soumission_module.md) |
| UCMOD02 | Ajouter un Module existant | [UCMOD02 (expression)](1-Expression/UCMOD-Module/UCMOD02.md) | [UCMOD02 (analyse)](2-Analyse/UCMOD-Module/UCMOD02.md) | [Architecture_Composition](3-Conception/Architecture_Composition.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [Sequence_soumission_module](3-Conception/Sequence_soumission_module.md) |
| UCMOD03 | Modifier les métadonnées d'un Module | [UCMOD03 (expression)](1-Expression/UCMOD-Module/UCMOD03.md) | [UCMOD03 (analyse)](2-Analyse/UCMOD-Module/UCMOD03.md) | [API_REST](3-Conception/API_REST.md) · [Architecture_Composition](3-Conception/Architecture_Composition.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) |
| UCMOD04 | Visualiser les composants d'un Module | [UCMOD04 (expression)](1-Expression/UCMOD-Module/UCMOD04.md) | [UCMOD04 (analyse)](2-Analyse/UCMOD-Module/UCMOD04.md) | [API_REST](3-Conception/API_REST.md) · [Architecture_Composition](3-Conception/Architecture_Composition.md) · [Chaincode](3-Conception/Chaincode.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D1_Auth_Identity](3-Conception/DC_D1_Auth_Identity.md) |
| UCMOD06 | Soumettre un module à la blockchain | [UCMOD06 (expression)](1-Expression/UCMOD-Module/UCMOD06.md) | [UCMOD06 (analyse)](2-Analyse/UCMOD-Module/UCMOD06.md) | [Architecture_Composition](3-Conception/Architecture_Composition.md) · [Chaincode](3-Conception/Chaincode.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [Sequence_soumission_module](3-Conception/Sequence_soumission_module.md) |
| UCMOD07 | Lister ses Modules en brouillon | [UCMOD07 (expression)](1-Expression/UCMOD-Module/UCMOD07.md) | [UCMOD07 (analyse)](2-Analyse/UCMOD-Module/UCMOD07.md) | — |
| UCMOD08 | Supprimer un Module | [UCMOD08 (expression)](1-Expression/UCMOD-Module/UCMOD08.md) | [UCMOD08 (analyse)](2-Analyse/UCMOD-Module/UCMOD08.md) | [API_REST](3-Conception/API_REST.md) |

### UCPI — Propriete Intellectuelle

Domaines : `model`, `payment`

| UC | Titre | Expression | Analyse | Conception |
|----|-------|------------|---------|------------|
| UCPI01 | Commander un Module complet | [UCPI01 (expression)](1-Expression/UCPI-Propriete_Intellectuelle/UCPI01.md) | [UCPI01 (analyse)](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI01.md) | [API_REST](3-Conception/API_REST.md) · [Chaincode](3-Conception/Chaincode.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D7_Payment](3-Conception/DC_D7_Payment.md) · [DC_D9_Automatisation](3-Conception/DC_D9_Automatisation.md) · [Modele_Domaine](3-Conception/Modele_Domaine.md) |
| UCPI02 | Recevoir une commission sur l'utilisation d'un Module | [UCPI02 (expression)](1-Expression/UCPI-Propriete_Intellectuelle/UCPI02.md) | [UCPI02 (analyse)](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI02.md) | [Chaincode](3-Conception/Chaincode.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D7_Payment](3-Conception/DC_D7_Payment.md) · [Modele_Domaine](3-Conception/Modele_Domaine.md) |
| UCPI03 | [RECLASSIFIÉ] Paramètres de langue de l'interface | [UCPI03 (expression)](1-Expression/UCPI-Propriete_Intellectuelle/UCPI03.md) | [UCPI03 (analyse)](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI03.md) | [DC_D7_Payment](3-Conception/DC_D7_Payment.md) |
| UCPI04 | Définir un prix sur un Composant proprietaire | [UCPI04 (expression)](1-Expression/UCPI-Propriete_Intellectuelle/UCPI04.md) | [UCPI04 (analyse)](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI04.md) | [API_REST](3-Conception/API_REST.md) · [Chaincode](3-Conception/Chaincode.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D7_Payment](3-Conception/DC_D7_Payment.md) · [DC_D8_Recherche](3-Conception/DC_D8_Recherche.md) |
| UCPI05 | Définir un prix sur un Module proprietaire | [UCPI05 (expression)](1-Expression/UCPI-Propriete_Intellectuelle/UCPI05.md) | [UCPI05 (analyse)](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI05.md) | [Chaincode](3-Conception/Chaincode.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D7_Payment](3-Conception/DC_D7_Payment.md) · [DC_D8_Recherche](3-Conception/DC_D8_Recherche.md) |
| UCPI06 | Déclarer un composant similaire | [UCPI06 (expression)](1-Expression/UCPI-Propriete_Intellectuelle/UCPI06.md) | [UCPI06 (analyse)](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI06.md) | [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D7_Payment](3-Conception/DC_D7_Payment.md) |
| UCPI07 | Transfert de propriété intellectuelle | [UCPI07 (expression)](1-Expression/UCPI-Propriete_Intellectuelle/UCPI07.md) | [UCPI07 (analyse)](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI07.md) | [API_REST](3-Conception/API_REST.md) · [Chaincode](3-Conception/Chaincode.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D7_Payment](3-Conception/DC_D7_Payment.md) |
| UCPI08 | Cloner un Composant sur un réseau exterieur | [UCPI08 (expression)](1-Expression/UCPI-Propriete_Intellectuelle/UCPI08.md) | [UCPI08 (analyse)](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI08.md) | [API_REST](3-Conception/API_REST.md) · [Chaincode](3-Conception/Chaincode.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D7_Payment](3-Conception/DC_D7_Payment.md) |
| UCPI09 | Cloner un Module sur un réseau exterieur | [UCPI09 (expression)](1-Expression/UCPI-Propriete_Intellectuelle/UCPI09.md) | [UCPI09 (analyse)](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI09.md) | [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D7_Payment](3-Conception/DC_D7_Payment.md) |
| UCPI10 | Norme de conception écoconception | [UCPI10 (expression)](1-Expression/UCPI-Propriete_Intellectuelle/UCPI10.md) | [UCPI10 (analyse)](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI10.md) | [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D7_Payment](3-Conception/DC_D7_Payment.md) |
| UCPI11 | Modifier le prix d'un asset | [UCPI11 (expression)](1-Expression/UCPI-Propriete_Intellectuelle/UCPI11.md) | [UCPI11 (analyse)](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI11.md) | [DC_CLI_Model](3-Conception/DC_CLI_Model.md) |

### UCREC — Recherche

Domaines : `model`

| UC | Titre | Expression | Analyse | Conception |
|----|-------|------------|---------|------------|
| UCREC01 | Rechercher une référence existante | [UCREC01 (expression)](1-Expression/UCREC-Recherche/UCREC01.md) | [UCREC01 (analyse)](2-Analyse/UCREC-Recherche/UCREC01.md) | [Chaincode](3-Conception/Chaincode.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D8_Recherche](3-Conception/DC_D8_Recherche.md) |
| UCREC02 | Rechercher les Composants compatibles | [UCREC02 (expression)](1-Expression/UCREC-Recherche/UCREC02.md) | [UCREC02 (analyse)](2-Analyse/UCREC-Recherche/UCREC02.md) | [API_REST](3-Conception/API_REST.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D8_Recherche](3-Conception/DC_D8_Recherche.md) |
| UCREC03 | Rechercher les versions des Composants | [UCREC03 (expression)](1-Expression/UCREC-Recherche/UCREC03.md) | [UCREC03 (analyse)](2-Analyse/UCREC-Recherche/UCREC03.md) | [API_REST](3-Conception/API_REST.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D8_Recherche](3-Conception/DC_D8_Recherche.md) |
| UCREC04 | Rechercher les Modules qui utilisent un Composant | [UCREC04 (expression)](1-Expression/UCREC-Recherche/UCREC04.md) | [UCREC04 (analyse)](2-Analyse/UCREC-Recherche/UCREC04.md) | [API_REST](3-Conception/API_REST.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D8_Recherche](3-Conception/DC_D8_Recherche.md) |
| UCREC05 | Exporter BOM Module | [UCREC05 (expression)](1-Expression/UCREC-Recherche/UCREC05.md) | [UCREC05 (analyse)](2-Analyse/UCREC-Recherche/UCREC05.md) | [API_REST](3-Conception/API_REST.md) · [DC_CLI_Model](3-Conception/DC_CLI_Model.md) · [DC_D8_Recherche](3-Conception/DC_D8_Recherche.md) |

## Règles métier

| RM | Règle | Use cases qui la citent |
|----|-------|------------------------|
| [RM01](1-Expression/Regles_Metier.md#RM01%20—%20Anti-plagiat%20obligatoire) | Anti-plagiat obligatoire | [UCAUT03](2-Analyse/UCAUT-Automatisation/UCAUT03.md) · [UCAUT04](2-Analyse/UCAUT-Automatisation/UCAUT04.md) · [UCCE01](2-Analyse/UCCE-Composant_Ecriture/UCCE01.md) · [UCCE03](1-Expression/UCCE-Composant_Ecriture/UCCE03.md) · [UCCL03](1-Expression/UCCL-Composant_Lecture/UCCL03.md) · [UCDEV01](2-Analyse/UCDEV-Developpement/UCDEV01.md) · [UCDEV02](2-Analyse/UCDEV-Developpement/UCDEV02.md) · [UCPI06](1-Expression/UCPI-Propriete_Intellectuelle/UCPI06.md) |
| [RM02](1-Expression/Regles_Metier.md#RM02%20—%20Catégorie%20d'asset%20obligatoire) | Catégorie d'asset obligatoire | [UCAM05](2-Analyse/UCAM-Assemblage_Module/UCAM05.md) · [UCAUT03](2-Analyse/UCAUT-Automatisation/UCAUT03.md) · [UCCE01](2-Analyse/UCCE-Composant_Ecriture/UCCE01.md) · [UCCE03](2-Analyse/UCCE-Composant_Ecriture/UCCE03.md) · [UCCE04](2-Analyse/UCCE-Composant_Ecriture/UCCE04.md) · [UCCE05](2-Analyse/UCCE-Composant_Ecriture/UCCE05.md) · [UCCE06](2-Analyse/UCCE-Composant_Ecriture/UCCE06.md) · [UCREC03](2-Analyse/UCREC-Recherche/UCREC03.md) |
| [RM03](1-Expression/Regles_Metier.md#RM03%20—%20Compatibilité%20de%20licence) | Compatibilité de licence | [UCCE01](2-Analyse/UCCE-Composant_Ecriture/UCCE01.md) · [UCCE02](2-Analyse/UCCE-Composant_Ecriture/UCCE02.md) · [UCCE03](2-Analyse/UCCE-Composant_Ecriture/UCCE03.md) · [UCCE04](2-Analyse/UCCE-Composant_Ecriture/UCCE04.md) · [UCCE05](2-Analyse/UCCE-Composant_Ecriture/UCCE05.md) · [UCMOD01](2-Analyse/UCMOD-Module/UCMOD01.md) · [UCMOD03](2-Analyse/UCMOD-Module/UCMOD03.md) · [UCPI04](1-Expression/UCPI-Propriete_Intellectuelle/UCPI04.md) · [UCPI08](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI08.md) · [UCPI09](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI09.md) |
| [RM04](1-Expression/Regles_Metier.md#RM04%20—%20UUID%20unique) | UUID unique | [UCAUT03](2-Analyse/UCAUT-Automatisation/UCAUT03.md) · [UCCE01](2-Analyse/UCCE-Composant_Ecriture/UCCE01.md) · [UCCE03](2-Analyse/UCCE-Composant_Ecriture/UCCE03.md) · [UCMOD01](2-Analyse/UCMOD-Module/UCMOD01.md) |
| [RM05](1-Expression/Regles_Metier.md#RM05%20—%20ParentID%20obligatoire%20pour%20les%20dérivés) | ParentID obligatoire pour les dérivés | [UCAM05](2-Analyse/UCAM-Assemblage_Module/UCAM05.md) · [UCAUT03](2-Analyse/UCAUT-Automatisation/UCAUT03.md) · [UCCE01](2-Analyse/UCCE-Composant_Ecriture/UCCE01.md) · [UCCE03](2-Analyse/UCCE-Composant_Ecriture/UCCE03.md) · [UCCE04](2-Analyse/UCCE-Composant_Ecriture/UCCE04.md) · [UCCE05](2-Analyse/UCCE-Composant_Ecriture/UCCE05.md) · [UCPI06](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI06.md) · [UCREC03](2-Analyse/UCREC-Recherche/UCREC03.md) |
| [RM06](1-Expression/Regles_Metier.md#RM06%20—%20Immuabilité%20des%20transactions) | Immuabilité des transactions | [UCADM01](2-Analyse/UCADM-Administration/UCADM01.md) · [UCADM04](1-Expression/UCADM-Administration/UCADM04.md) · [UCADM05](2-Analyse/UCADM-Administration/UCADM05.md) · [UCCE07](1-Expression/UCCE-Composant_Ecriture/UCCE07.md) · [UCMOD08](1-Expression/UCMOD-Module/UCMOD08.md) |
| [RM07](1-Expression/Regles_Metier.md#RM07%20—%20Validation%20préalable%20obligatoire) | Validation préalable obligatoire | [UCADM01](2-Analyse/UCADM-Administration/UCADM01.md) · [UCADM02](2-Analyse/UCADM-Administration/UCADM02.md) · [UCADM03](2-Analyse/UCADM-Administration/UCADM03.md) · [UCADM04](2-Analyse/UCADM-Administration/UCADM04.md) · [UCAM05](2-Analyse/UCAM-Assemblage_Module/UCAM05.md) · [UCAUT01](2-Analyse/UCAUT-Automatisation/UCAUT01.md) · [UCAUT02](2-Analyse/UCAUT-Automatisation/UCAUT02.md) · [UCCE01](2-Analyse/UCCE-Composant_Ecriture/UCCE01.md) · [UCCE02](2-Analyse/UCCE-Composant_Ecriture/UCCE02.md) · [UCCE03](2-Analyse/UCCE-Composant_Ecriture/UCCE03.md) · [UCCE04](2-Analyse/UCCE-Composant_Ecriture/UCCE04.md) · [UCCE05](2-Analyse/UCCE-Composant_Ecriture/UCCE05.md) · [UCDEV01](2-Analyse/UCDEV-Developpement/UCDEV01.md) · [UCDEV02](2-Analyse/UCDEV-Developpement/UCDEV02.md) · [UCMOD03](2-Analyse/UCMOD-Module/UCMOD03.md) · [UCMOD06](2-Analyse/UCMOD-Module/UCMOD06.md) · [UCPI01](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI01.md) · [UCPI02](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI02.md) · [UCPI06](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI06.md) · [UCPI07](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI07.md) · [UCPI08](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI08.md) · [UCPI09](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI09.md) · [UCPI10](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI10.md) |
| [RM08](1-Expression/Regles_Metier.md#RM08%20—%20Masquage%20local,%20ledger%20jamais%20modifié) | Masquage local, ledger jamais modifié | [UCADM02](2-Analyse/UCADM-Administration/UCADM02.md) · [UCADM03](2-Analyse/UCADM-Administration/UCADM03.md) · [UCADM05](2-Analyse/UCADM-Administration/UCADM05.md) · [UCAM05](2-Analyse/UCAM-Assemblage_Module/UCAM05.md) · [UCAM08](2-Analyse/UCAM-Assemblage_Module/UCAM08.md) · [UCCE04](2-Analyse/UCCE-Composant_Ecriture/UCCE04.md) · [UCCE05](2-Analyse/UCCE-Composant_Ecriture/UCCE05.md) · [UCCE07](2-Analyse/UCCE-Composant_Ecriture/UCCE07.md) · [UCDEV01](2-Analyse/UCDEV-Developpement/UCDEV01.md) · [UCMOD08](2-Analyse/UCMOD-Module/UCMOD08.md) |
| [RM09](1-Expression/Regles_Metier.md#RM09%20—%20Interface%20à%20usage%20unique) | Interface à usage unique | [UCAM01](1-Expression/UCAM-Assemblage_Module/UCAM01.md) · [UCAM03](2-Analyse/UCAM-Assemblage_Module/UCAM03.md) · [UCAM07](2-Analyse/UCAM-Assemblage_Module/UCAM07.md) |
| [RM10](1-Expression/Regles_Metier.md#RM10%20—%20Vérification%20de%20compatibilité%20automatique) | Vérification de compatibilité automatique | [UCAM01](1-Expression/UCAM-Assemblage_Module/UCAM01.md) · [UCAM03](2-Analyse/UCAM-Assemblage_Module/UCAM03.md) · [UCAM07](2-Analyse/UCAM-Assemblage_Module/UCAM07.md) · [UCAM09](1-Expression/UCAM-Assemblage_Module/UCAM09.md) · [UCMOD01](2-Analyse/UCMOD-Module/UCMOD01.md) · [UCREC02](1-Expression/UCREC-Recherche/UCREC02.md) |
| [RM11](1-Expression/Regles_Metier.md#RM11%20—%20Critères%20de%20compatibilité%20d'interfaces) | Critères de compatibilité d'interfaces | [UCAM01](1-Expression/UCAM-Assemblage_Module/UCAM01.md) · [UCAM02](2-Analyse/UCAM-Assemblage_Module/UCAM02.md) · [UCAM03](2-Analyse/UCAM-Assemblage_Module/UCAM03.md) · [UCAM07](2-Analyse/UCAM-Assemblage_Module/UCAM07.md) · [UCAM09](1-Expression/UCAM-Assemblage_Module/UCAM09.md) · [UCCE06](2-Analyse/UCCE-Composant_Ecriture/UCCE06.md) · [UCMOD01](2-Analyse/UCMOD-Module/UCMOD01.md) · [UCREC02](1-Expression/UCREC-Recherche/UCREC02.md) |
| [RM12](1-Expression/Regles_Metier.md#RM12%20—%20Persistance%20des%20liaisons%20incompatibles) | Persistance des liaisons incompatibles | [UCAM01](2-Analyse/UCAM-Assemblage_Module/UCAM01.md) · [UCCE06](2-Analyse/UCCE-Composant_Ecriture/UCCE06.md) |
| [RM13](1-Expression/Regles_Metier.md#RM13%20—%20Slot%20virtuel%20garanti) | Slot virtuel garanti | [UCAM01](2-Analyse/UCAM-Assemblage_Module/UCAM01.md) · [UCAM02](2-Analyse/UCAM-Assemblage_Module/UCAM02.md) · [UCAM03](1-Expression/UCAM-Assemblage_Module/UCAM03.md) · [UCAM05](2-Analyse/UCAM-Assemblage_Module/UCAM05.md) · [UCCE06](2-Analyse/UCCE-Composant_Ecriture/UCCE06.md) · [UCMOD01](2-Analyse/UCMOD-Module/UCMOD01.md) · [UCMOD02](2-Analyse/UCMOD-Module/UCMOD02.md) · [UCMOD04](2-Analyse/UCMOD-Module/UCMOD04.md) |
| [RM14](1-Expression/Regles_Metier.md#RM14%20—%20Suppression%20en%20cascade%20des%20connexions) | Suppression en cascade des connexions | [UCAM08](2-Analyse/UCAM-Assemblage_Module/UCAM08.md) · [UCCE07](2-Analyse/UCCE-Composant_Ecriture/UCCE07.md) · [UCMOD08](2-Analyse/UCMOD-Module/UCMOD08.md) |
| [RM15](1-Expression/Regles_Metier.md#RM15%20—%20Instance%20indépendante) | Instance indépendante | [UCAM08](1-Expression/UCAM-Assemblage_Module/UCAM08.md) · [UCMOD02](2-Analyse/UCMOD-Module/UCMOD02.md) · [UCMOD08](2-Analyse/UCMOD-Module/UCMOD08.md) · [UCREC05](1-Expression/UCREC-Recherche/UCREC05.md) |
| [RM16](1-Expression/Regles_Metier.md#RM16%20—%20État%20draft) | État draft | [UCAM09](1-Expression/UCAM-Assemblage_Module/UCAM09.md) · [UCCE01](2-Analyse/UCCE-Composant_Ecriture/UCCE01.md) · [UCCE06](2-Analyse/UCCE-Composant_Ecriture/UCCE06.md) · [UCMOD01](2-Analyse/UCMOD-Module/UCMOD01.md) · [UCREC04](2-Analyse/UCREC-Recherche/UCREC04.md) |
| [RM17](1-Expression/Regles_Metier.md#RM17%20—%20Assemblage%20requis%20pour%20soumission%20%28module%20uniquement%29) | Assemblage requis pour soumission (module uniquement) | [UCAM05](2-Analyse/UCAM-Assemblage_Module/UCAM05.md) · [UCMOD06](1-Expression/UCMOD-Module/UCMOD06.md) · [UCPI05](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI05.md) |
| [RM18](1-Expression/Regles_Metier.md#RM18%20—%20ModuleVersion%20immuable%20%28module%20uniquement%29) | ModuleVersion immuable (module uniquement) | [UCAM05](2-Analyse/UCAM-Assemblage_Module/UCAM05.md) · [UCAUT04](2-Analyse/UCAUT-Automatisation/UCAUT04.md) · [UCMOD06](2-Analyse/UCMOD-Module/UCMOD06.md) |
| [RM19](1-Expression/Regles_Metier.md#RM19%20—%20Fork%20d'un%20asset%20soumis) | Fork d'un asset soumis | [UCAM05](2-Analyse/UCAM-Assemblage_Module/UCAM05.md) · [UCAM08](2-Analyse/UCAM-Assemblage_Module/UCAM08.md) · [UCCE01](2-Analyse/UCCE-Composant_Ecriture/UCCE01.md) · [UCCE06](2-Analyse/UCCE-Composant_Ecriture/UCCE06.md) · [UCMOD01](2-Analyse/UCMOD-Module/UCMOD01.md) · [UCMOD03](2-Analyse/UCMOD-Module/UCMOD03.md) · [UCMOD06](2-Analyse/UCMOD-Module/UCMOD06.md) · [UCREC04](2-Analyse/UCREC-Recherche/UCREC04.md) |
| [RM20](1-Expression/Regles_Metier.md#RM20%20—%20Identité%20=%20enrôlement%20CA,%20pas%20un%20compte%20séparé) | Identité = enrôlement CA, pas un compte séparé | [UCA01](2-Analyse/UCA-Compte_et_Acces/UCA01.md) · [UCA02](2-Analyse/UCA-Compte_et_Acces/UCA02.md) |
| [RM21](1-Expression/Regles_Metier.md#RM21%20—%20Rôle%20Lecteur%20par%20défaut%20à%20l'auto-enregistrement) | Rôle Lecteur par défaut à l'auto-enregistrement | [UCA01](2-Analyse/UCA-Compte_et_Acces/UCA01.md) |
| [RM22](1-Expression/Regles_Metier.md#RM22%20—%20Changement%20de%20rôle%20réservé%20à%20l'administrateur) | Changement de rôle réservé à l'administrateur | [UCA02](2-Analyse/UCA-Compte_et_Acces/UCA02.md) · [UCA04](2-Analyse/UCA-Compte_et_Acces/UCA04.md) · [UCA05](2-Analyse/UCA-Compte_et_Acces/UCA05.md) · [UCA07](2-Analyse/UCA-Compte_et_Acces/UCA07.md) · [UCA08](2-Analyse/UCA-Compte_et_Acces/UCA08.md) · [UCAUT01](2-Analyse/UCAUT-Automatisation/UCAUT01.md) · [UCAUT02](2-Analyse/UCAUT-Automatisation/UCAUT02.md) · [UCAUT03](2-Analyse/UCAUT-Automatisation/UCAUT03.md) · [UCAUT04](2-Analyse/UCAUT-Automatisation/UCAUT04.md) · [UCPI01](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI01.md) · [UCPI04](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI04.md) · [UCPI05](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI05.md) · [UCPI06](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI06.md) · [UCPI07](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI07.md) · [UCPI08](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI08.md) · [UCPI09](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI09.md) · [UCPI10](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI10.md) · [UCPI11](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI11.md) · [UCREC01](2-Analyse/UCREC-Recherche/UCREC01.md) |
| [RM23](1-Expression/Regles_Metier.md#RM23%20—%20Distribution%20automatique%20des%20commissions) | Distribution automatique des commissions | [UCAUT01](2-Analyse/UCAUT-Automatisation/UCAUT01.md) · [UCAUT02](2-Analyse/UCAUT-Automatisation/UCAUT02.md) · [UCPI01](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI01.md) · [UCPI02](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI02.md) · [UCPI03](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI03.md) · [UCPI04](1-Expression/UCPI-Propriete_Intellectuelle/UCPI04.md) · [UCPI05](1-Expression/UCPI-Propriete_Intellectuelle/UCPI05.md) |
| [RM24](1-Expression/Regles_Metier.md#RM24%20—%20Répartition%20proportionnelle%20multi-auteurs) | Répartition proportionnelle multi-auteurs | [UCAUT01](2-Analyse/UCAUT-Automatisation/UCAUT01.md) · [UCPI02](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI02.md) · [UCPI03](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI03.md) · [UCPI04](1-Expression/UCPI-Propriete_Intellectuelle/UCPI04.md) · [UCPI05](1-Expression/UCPI-Propriete_Intellectuelle/UCPI05.md) |
| [RM25](1-Expression/Regles_Metier.md#RM25%20—%20Transfert%20de%20propriété%20définitif) | Transfert de propriété définitif | [UCPI03](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI03.md) · [UCPI07](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI07.md) |
| [RM26](1-Expression/Regles_Metier.md#RM26%20—%20Traçabilité%20du%20clonage%20inter-réseaux) | Traçabilité du clonage inter-réseaux | [UCPI03](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI03.md) · [UCPI07](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI07.md) · [UCPI08](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI08.md) · [UCPI09](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI09.md) |
| [RM27](1-Expression/Regles_Metier.md#RM27%20—%20Nombre%20minimum%20de%20nœuds%20actifs) | Nombre minimum de nœuds actifs | [UCADM04](1-Expression/UCADM-Administration/UCADM04.md) |
| [RM28](1-Expression/Regles_Metier.md#RM28%20—%20Démantèlement%20réseau%20:%20opération%20d'infrastructure%20locale) | Démantèlement réseau : opération d'infrastructure locale | [UCADM05](1-Expression/UCADM-Administration/UCADM05.md) |
| [RM29](1-Expression/Regles_Metier.md#RM29%20—%20Taux%20de%20commission%20défini%20par%20le%20réseau) | Taux de commission défini par le réseau | [UCPI04](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI04.md) · [UCPI05](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI05.md) |
| [RM30](1-Expression/Regles_Metier.md#RM30%20—%20Calcul%20automatique%20du%20prix%20d'un%20module) | Calcul automatique du prix d'un module | [UCPI04](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI04.md) · [UCPI05](1-Expression/UCPI-Propriete_Intellectuelle/UCPI05.md) |
| [RM31](1-Expression/Regles_Metier.md#RM31%20—%20Modification%20de%20prix%20—%20effet%20sur%20les%20commandes%20futures%20uniquement) | Modification de prix — effet sur les commandes futures uniquement | [UCPI04](1-Expression/UCPI-Propriete_Intellectuelle/UCPI04.md) · [UCPI05](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI05.md) · [UCPI11](1-Expression/UCPI-Propriete_Intellectuelle/UCPI11.md) |
| [RM32](1-Expression/Regles_Metier.md#RM32%20—%20Asset%20à%20prix%20nul%20—%20librement%20disponible) | Asset à prix nul — librement disponible | [UCPI04](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI04.md) · [UCPI05](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI05.md) · [UCPI11](1-Expression/UCPI-Propriete_Intellectuelle/UCPI11.md) |
| [RM33](1-Expression/Regles_Metier.md#RM33%20—%20Devise%20unique%20par%20réseau) | Devise unique par réseau | [UCPI04](1-Expression/UCPI-Propriete_Intellectuelle/UCPI04.md) · [UCPI05](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI05.md) · [UCPI11](1-Expression/UCPI-Propriete_Intellectuelle/UCPI11.md) |
| [RM34](1-Expression/Regles_Metier.md#RM34%20—%20Rôle%20admin%20protégé) | Rôle admin protégé | [UCADM06](2-Analyse/UCADM-Administration/UCADM06.md) · [UCADM07](2-Analyse/UCADM-Administration/UCADM07.md) |
| [RM35](1-Expression/Regles_Metier.md#RM35%20—%20Révocation%20en%20cascade%20à%20la%20suppression%20d'un%20rôle) | Révocation en cascade à la suppression d'un rôle | [UCADM07](2-Analyse/UCADM-Administration/UCADM07.md) |
| [RM36](1-Expression/Regles_Metier.md#RM36%20—%20Nom%20de%20rôle%20unique) | Nom de rôle unique | [UCADM07](2-Analyse/UCADM-Administration/UCADM07.md) |
| [RM37](1-Expression/Regles_Metier.md#RM37%20—%20Multi-rôles%20par%20organisation) | Multi-rôles par organisation | [UCADM06](2-Analyse/UCADM-Administration/UCADM06.md) |
| [RM39](1-Expression/Regles_Metier.md#RM39%20—%20Filiation%20d'un%20découpage) | Filiation d'un découpage | [UCAM09](1-Expression/UCAM-Assemblage_Module/UCAM09.md) |
| [RM40](1-Expression/Regles_Metier.md#RM40%20—%20Proposition%20de%20découpage%20non%20engageante) | Proposition de découpage non engageante | [UCAM09](1-Expression/UCAM-Assemblage_Module/UCAM09.md) |
| [RM41](1-Expression/Regles_Metier.md#RM41%20—%20Compatibilité%20toujours%20vérifiée%20pour%20une%20connexion%20suggérée) | Compatibilité toujours vérifiée pour une connexion suggérée | [UCAM09](1-Expression/UCAM-Assemblage_Module/UCAM09.md) |
| [RM42](1-Expression/Regles_Metier.md#RM42%20—%20Traçabilité%20et%20alerte%20des%20emplacements%20externes) | Traçabilité et alerte des emplacements externes | — |

## Exigences non fonctionnelles

| ENF | Exigence | Use cases qui la citent |
|-----|----------|------------------------|
| [ENF01](1-Expression/Exigences_Non_Fonctionnelles.md#ENF01%20—%20Temps%20de%20réponse%20des%20endpoints%20REST) | Temps de réponse des endpoints REST | [UCA02](2-Analyse/UCA-Compte_et_Acces/UCA02.md) · [UCAUT01](2-Analyse/UCAUT-Automatisation/UCAUT01.md) · [UCAUT02](2-Analyse/UCAUT-Automatisation/UCAUT02.md) · [UCCL01](2-Analyse/UCCL-Composant_Lecture/UCCL01.md) · [UCMOD07](2-Analyse/UCMOD-Module/UCMOD07.md) · [UCPI01](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI01.md) |
| [ENF02](1-Expression/Exigences_Non_Fonctionnelles.md#ENF02%20—%20Temps%20de%20soumission%20d'une%20transaction%20Fabric) | Temps de soumission d'une transaction Fabric | [UCA06](2-Analyse/UCA-Compte_et_Acces/UCA06.md) · [UCAUT01](2-Analyse/UCAUT-Automatisation/UCAUT01.md) · [UCMOD06](2-Analyse/UCMOD-Module/UCMOD06.md) · [UCPI02](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI02.md) |
| [ENF03](1-Expression/Exigences_Non_Fonctionnelles.md#ENF03%20—%20Génération%20d'une%20BOM%20module) | Génération d'une BOM module | [UCPI05](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI05.md) · [UCPI09](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI09.md) · [UCREC05](2-Analyse/UCREC-Recherche/UCREC05.md) |
| [ENF05](1-Expression/Exigences_Non_Fonctionnelles.md#ENF05%20—%20Disponibilité%20du%20serveur%20Myr%20%28instance%20unique%29) | Disponibilité du serveur Myr (instance unique) | [UCADM01](2-Analyse/UCADM-Administration/UCADM01.md) |
| [ENF06](1-Expression/Exigences_Non_Fonctionnelles.md#ENF06%20—%20Disponibilité%20du%20réseau%20blockchain%20%28multi-nœuds%29) | Disponibilité du réseau blockchain (multi-nœuds) | [UCADM04](2-Analyse/UCADM-Administration/UCADM04.md) |
| [ENF07](1-Expression/Exigences_Non_Fonctionnelles.md#ENF07%20—%20Reprise%20après%20redémarrage%20du%20serveur) | Reprise après redémarrage du serveur | — |
| [ENF08](1-Expression/Exigences_Non_Fonctionnelles.md#ENF08%20—%20Chiffrement%20des%20communications) | Chiffrement des communications | — |
| [ENF09](1-Expression/Exigences_Non_Fonctionnelles.md#ENF09%20—%20Durée%20de%20vie%20des%20tokens%20de%20session) | Durée de vie des tokens de session | — |
| [ENF10](1-Expression/Exigences_Non_Fonctionnelles.md#ENF10%20—%20Chiffrement%20des%20wallets%20Fabric) | Chiffrement des wallets Fabric | — |
| [ENF11](1-Expression/Exigences_Non_Fonctionnelles.md#ENF11%20—%20Secrets%20absents%20des%20logs) | Secrets absents des logs | — |
| [ENF12](1-Expression/Exigences_Non_Fonctionnelles.md#ENF12%20—%20Contrôle%20d'accès%20par%20rôle) | Contrôle d'accès par rôle | [UCA01](2-Analyse/UCA-Compte_et_Acces/UCA01.md) · [UCA02](2-Analyse/UCA-Compte_et_Acces/UCA02.md) · [UCA03](2-Analyse/UCA-Compte_et_Acces/UCA03.md) · [UCA04](2-Analyse/UCA-Compte_et_Acces/UCA04.md) · [UCA05](2-Analyse/UCA-Compte_et_Acces/UCA05.md) · [UCA08](2-Analyse/UCA-Compte_et_Acces/UCA08.md) · [UCAM01](2-Analyse/UCAM-Assemblage_Module/UCAM01.md) · [UCAM02](2-Analyse/UCAM-Assemblage_Module/UCAM02.md) · [UCAM03](2-Analyse/UCAM-Assemblage_Module/UCAM03.md) · [UCAM05](2-Analyse/UCAM-Assemblage_Module/UCAM05.md) · [UCAM07](2-Analyse/UCAM-Assemblage_Module/UCAM07.md) · [UCAM08](2-Analyse/UCAM-Assemblage_Module/UCAM08.md) · [UCAUT01](2-Analyse/UCAUT-Automatisation/UCAUT01.md) · [UCAUT02](2-Analyse/UCAUT-Automatisation/UCAUT02.md) · [UCAUT03](2-Analyse/UCAUT-Automatisation/UCAUT03.md) · [UCAUT04](2-Analyse/UCAUT-Automatisation/UCAUT04.md) · [UCCE01](2-Analyse/UCCE-Composant_Ecriture/UCCE01.md) · [UCCE02](2-Analyse/UCCE-Composant_Ecriture/UCCE02.md) · [UCCE03](2-Analyse/UCCE-Composant_Ecriture/UCCE03.md) · [UCCE04](2-Analyse/UCCE-Composant_Ecriture/UCCE04.md) · [UCCE05](2-Analyse/UCCE-Composant_Ecriture/UCCE05.md) · [UCCE06](2-Analyse/UCCE-Composant_Ecriture/UCCE06.md) · [UCCE07](2-Analyse/UCCE-Composant_Ecriture/UCCE07.md) · [UCCL01](2-Analyse/UCCL-Composant_Lecture/UCCL01.md) · [UCDEV01](2-Analyse/UCDEV-Developpement/UCDEV01.md) · [UCMOD01](2-Analyse/UCMOD-Module/UCMOD01.md) · [UCMOD02](2-Analyse/UCMOD-Module/UCMOD02.md) · [UCMOD03](2-Analyse/UCMOD-Module/UCMOD03.md) · [UCMOD04](2-Analyse/UCMOD-Module/UCMOD04.md) · [UCMOD06](2-Analyse/UCMOD-Module/UCMOD06.md) · [UCMOD07](2-Analyse/UCMOD-Module/UCMOD07.md) · [UCMOD08](2-Analyse/UCMOD-Module/UCMOD08.md) · [UCPI01](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI01.md) · [UCPI04](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI04.md) · [UCPI05](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI05.md) · [UCPI06](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI06.md) · [UCPI07](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI07.md) · [UCPI08](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI08.md) · [UCPI09](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI09.md) · [UCPI10](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI10.md) · [UCPI11](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI11.md) · [UCREC01](2-Analyse/UCREC-Recherche/UCREC01.md) · [UCREC02](2-Analyse/UCREC-Recherche/UCREC02.md) · [UCREC03](2-Analyse/UCREC-Recherche/UCREC03.md) · [UCREC04](2-Analyse/UCREC-Recherche/UCREC04.md) · [UCREC05](2-Analyse/UCREC-Recherche/UCREC05.md) |
| [ENF13](1-Expression/Exigences_Non_Fonctionnelles.md#ENF13%20—%20Protection%20contre%20l'injection) | Protection contre l'injection | — |
| [ENF14](1-Expression/Exigences_Non_Fonctionnelles.md#ENF14%20—%20Nombre%20d'assets%20par%20réseau) | Nombre d'assets par réseau | [UCAUT02](2-Analyse/UCAUT-Automatisation/UCAUT02.md) · [UCCL01](2-Analyse/UCCL-Composant_Lecture/UCCL01.md) |
| [ENF15](1-Expression/Exigences_Non_Fonctionnelles.md#ENF15%20—%20Taille%20maximale%20d'un%20fichier%20CAO) | Taille maximale d'un fichier CAO | [UCAUT03](2-Analyse/UCAUT-Automatisation/UCAUT03.md) |
| [ENF16](1-Expression/Exigences_Non_Fonctionnelles.md#ENF16%20—%20Utilisateurs%20simultanés%20par%20instance) | Utilisateurs simultanés par instance | [UCAUT02](2-Analyse/UCAUT-Automatisation/UCAUT02.md) |
| [ENF17](1-Expression/Exigences_Non_Fonctionnelles.md#ENF17%20—%20Extension%20horizontale) | Extension horizontale | — |
| [ENF18](1-Expression/Exigences_Non_Fonctionnelles.md#ENF18%20—%20Isolation%20du%20domaine%20métier) | Isolation du domaine métier | [UCA05](2-Analyse/UCA-Compte_et_Acces/UCA05.md) · [UCA06](2-Analyse/UCA-Compte_et_Acces/UCA06.md) · [UCADM01](2-Analyse/UCADM-Administration/UCADM01.md) · [UCADM02](2-Analyse/UCADM-Administration/UCADM02.md) · [UCADM03](2-Analyse/UCADM-Administration/UCADM03.md) · [UCADM04](2-Analyse/UCADM-Administration/UCADM04.md) · [UCADM05](2-Analyse/UCADM-Administration/UCADM05.md) · [UCADM06](2-Analyse/UCADM-Administration/UCADM06.md) · [UCADM07](2-Analyse/UCADM-Administration/UCADM07.md) · [UCAM01](2-Analyse/UCAM-Assemblage_Module/UCAM01.md) · [UCAM03](2-Analyse/UCAM-Assemblage_Module/UCAM03.md) · [UCAM05](2-Analyse/UCAM-Assemblage_Module/UCAM05.md) · [UCAM07](2-Analyse/UCAM-Assemblage_Module/UCAM07.md) · [UCAM08](2-Analyse/UCAM-Assemblage_Module/UCAM08.md) · [UCMOD01](2-Analyse/UCMOD-Module/UCMOD01.md) · [UCPI06](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI06.md) |
| [ENF19](1-Expression/Exigences_Non_Fonctionnelles.md#ENF19%20—%20Couverture%20de%20tests%20unitaires) | Couverture de tests unitaires | — |
| [ENF20](1-Expression/Exigences_Non_Fonctionnelles.md#ENF20%20—%20Durée%20de%20déploiement%20d'une%20mise%20à%20jour) | Durée de déploiement d'une mise à jour | — |
| [ENF21](1-Expression/Exigences_Non_Fonctionnelles.md#ENF21%20—%20Compatibilité%20multi-réseaux%20blockchain) | Compatibilité multi-réseaux blockchain | — |
| [ENF23](1-Expression/Exigences_Non_Fonctionnelles.md#ENF23%20—%20Plateformes%20serveur) | Plateformes serveur | — |
| [ENF25](1-Expression/Exigences_Non_Fonctionnelles.md#ENF25%20—%20Licence%20du%20code%20source) | Licence du code source | — |
| [ENF26](1-Expression/Exigences_Non_Fonctionnelles.md#ENF26%20—%20Compatibilité%20de%20licence%20des%20assets) | Compatibilité de licence des assets | — |
| [ENF27](1-Expression/Exigences_Non_Fonctionnelles.md#ENF27%20—%20Protection%20des%20données%20personnelles%20%28RGPD%29) | Protection des données personnelles (RGPD) | [UCPI04](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI04.md) · [UCPI05](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI05.md) · [UCPI11](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI11.md) |
| [ENF28](1-Expression/Exigences_Non_Fonctionnelles.md#ENF28%20—%20Immuabilité%20des%20transactions%20blockchain) | Immuabilité des transactions blockchain | [UCMOD02](2-Analyse/UCMOD-Module/UCMOD02.md) · [UCMOD06](2-Analyse/UCMOD-Module/UCMOD06.md) · [UCMOD08](2-Analyse/UCMOD-Module/UCMOD08.md) |
| [ENF29](1-Expression/Exigences_Non_Fonctionnelles.md#ENF29%20—%20Anti-plagiat%20obligatoire) | Anti-plagiat obligatoire | — |
| [ENF30](1-Expression/Exigences_Non_Fonctionnelles.md#ENF30%20—%20Intégrité%20en%20cas%20d'échec%20blockchain) | Intégrité en cas d'échec blockchain | [UCAM05](2-Analyse/UCAM-Assemblage_Module/UCAM05.md) · [UCAUT01](2-Analyse/UCAUT-Automatisation/UCAUT01.md) · [UCAUT03](2-Analyse/UCAUT-Automatisation/UCAUT03.md) · [UCCE01](2-Analyse/UCCE-Composant_Ecriture/UCCE01.md) · [UCCE02](2-Analyse/UCCE-Composant_Ecriture/UCCE02.md) · [UCCE03](2-Analyse/UCCE-Composant_Ecriture/UCCE03.md) · [UCCE04](2-Analyse/UCCE-Composant_Ecriture/UCCE04.md) · [UCCE05](2-Analyse/UCCE-Composant_Ecriture/UCCE05.md) · [UCMOD03](2-Analyse/UCMOD-Module/UCMOD03.md) · [UCMOD06](2-Analyse/UCMOD-Module/UCMOD06.md) · [UCPI01](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI01.md) · [UCPI02](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI02.md) · [UCPI07](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI07.md) · [UCPI08](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI08.md) · [UCPI09](2-Analyse/UCPI-Propriete_Intellectuelle/UCPI09.md) |
| [ENF31](1-Expression/Exigences_Non_Fonctionnelles.md#ENF31%20—%20Validation%20avant%20soumission) | Validation avant soumission | [UCMOD01](2-Analyse/UCMOD-Module/UCMOD01.md) · [UCMOD06](2-Analyse/UCMOD-Module/UCMOD06.md) |
| [ENF32](1-Expression/Exigences_Non_Fonctionnelles.md#ENF32%20—%20Traçabilité%20et%20intégrité%20du%20fichier%20source) | Traçabilité et intégrité du fichier source | — |

## Documents de conception

| Document | Périmètre déclaré | Domaines |
|----------|-------------------|----------|
| [API_REST](3-Conception/API_REST.md) | — | — |
| [Architecture_Composition](3-Conception/Architecture_Composition.md) | UCAM01, UCAM02, UCAM03, UCAM05, UCAM07, UCAM08, UCCE01, UCCE02, UCCE03, UCCE04, UCCE05, UCCE06, UCMOD01, UCMOD02, UCMOD03, UCMOD04, UCMOD06 | `model` |
| [Architecture_Hexagonale](3-Conception/Architecture_Hexagonale.md) | — | — |
| [Chaincode](3-Conception/Chaincode.md) | — | — |
| [Conception_intro](3-Conception/Conception_intro.md) | — | — |
| [DC_CLI_Admin](3-Conception/DC_CLI_Admin.md) | UCADM01, UCADM02, UCADM03, UCADM04, UCADM05, UCDEV02 | `channel`, `identity`, `network` |
| [DC_CLI_Identity](3-Conception/DC_CLI_Identity.md) | UCA01, UCA02, UCA04, UCA07, UCA08 | `identity`, `role` |
| [DC_CLI_Model](3-Conception/DC_CLI_Model.md) | UCAM01, UCAM02, UCAM03, UCAM05, UCAM07, UCAM08, UCCE01, UCCE02, UCCE03, UCCE04, UCCE05, UCCE06, UCCL01, UCMOD01, UCMOD02, UCMOD03, UCMOD04, UCMOD06, UCREC01, UCREC02, UCREC03, UCREC04, UCREC05 | `model` |
| [DC_D1_Auth_Identity](3-Conception/DC_D1_Auth_Identity.md) | — | — |
| [DC_D2_Administration](3-Conception/DC_D2_Administration.md) | UCADM01, UCADM02, UCADM03, UCADM04, UCADM05, UCADM06, UCADM07 | `channel`, `identity`, `network` |
| [DC_D7_Payment](3-Conception/DC_D7_Payment.md) | UCPI01, UCPI02, UCPI03, UCPI04, UCPI05, UCPI06, UCPI07, UCPI08, UCPI09, UCPI10 | `payment` |
| [DC_D8_Recherche](3-Conception/DC_D8_Recherche.md) | UCREC01, UCREC02, UCREC03, UCREC04, UCREC05 | `model` |
| [DC_D9_Automatisation](3-Conception/DC_D9_Automatisation.md) | UCAUT01, UCAUT02, UCAUT03, UCAUT04 | `model`, `payment`, `role` |
| [Deploiement](3-Conception/Deploiement.md) | — | — |
| [Modele_Domaine](3-Conception/Modele_Domaine.md) | — | — |
| [Securite](3-Conception/Securite.md) | — | — |
| [Sequence_soumission_asset](3-Conception/Sequence_soumission_asset.md) | UCCE01, UCCE06 | `model` |
| [Sequence_soumission_module](3-Conception/Sequence_soumission_module.md) | UCAM01, UCAM07, UCAM08, UCMOD01, UCMOD02, UCMOD06 | `model` |

## Autres documents

- [Exigences_Non_Fonctionnelles](1-Expression/Exigences_Non_Fonctionnelles.md) — expression, exigences
- [Expression_des_besoins_Intro](1-Expression/Expression_des_besoins_Intro.md) — expression, introduction
- [Matrice_Tracabilite](1-Expression/Matrice_Tracabilite.md) — expression, tracabilite
- [Regles_Metier](1-Expression/Regles_Metier.md) — expression, regles-metier
- [todo (expression)](1-Expression/todo.md) — expression, todo
- [Analyse_des_besoins](2-Analyse/Analyse_des_besoins.md) — analyse, introduction
- [todo (analyse)](2-Analyse/todo.md) — analyse, todo
- [todo (conception)](3-Conception/todo.md) — conception, todo
- [licence_choix](licence_choix.md) — transverse, licence
- [roadmap_dev](roadmap_dev.md) — transverse, roadmap
