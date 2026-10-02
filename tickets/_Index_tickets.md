---
tags:
  - ticket/index
---
# Tickets — faits techniques de myr-core

Un fichier = un fait technique : anomalie, écart entre specs et code, risque, faille de sécurité, dette, incohérence documentaire ou décision produit en attente. 
Chaque ticket est relié par des liens aux use cases, règles métier, documents de conception, fichiers de code et notes de fonction concernés, et porte des tags qui le rendent filtrable dans Obsidian.

Les specs décrivent le comportement attendu ; l'écart entre ce comportement et le code, et son suivi, vivent ici.

Vue dynamique : [Tickets.base](Tickets.base) (Obsidian Bases) · Modèle : [_Modele_ticket](_Modele_ticket.md) · Voir aussi : [Carte des specs](../specs/Carte_des_specs.md) · [Traçabilité use cases ↔ code](../docs/code/Tracabilite_UC_Code.md)

## Synthèse

| Statut | Nombre |
|--------|--------|
| Ouvert | 20 |
| À vérifier | 2 |
| À trancher (PO) | 14 |
| Différé | 1 |
| Résolu | 1 |

## Critiques à traiter en premier

| Ticket | Type | Sévérité | Statut | Use cases |
|--------|------|----------|--------|-----------|
| [FT-001 — Propriété des assets non contrôlée côté serveur](FT-001-propriete-des-assets-non-controlee-cote-serveur.md) | Sécurité | critique | Ouvert | UCCE01, UCCE02, UCCE07, UCMOD03, UCMOD04, UCMOD08, UCPI07 |
| [FT-002 — Rôle de session REST codé en dur à contributor](FT-002-role-de-session-rest-code-en-dur.md) | Sécurité | critique | Ouvert | UCA02, UCA05, UCA07 |
| [FT-003 — Le client CA du serveur REST ignore le profil réseau actif](FT-003-le-client-ca-du-serveur-rest-ignore.md) | Anomalie | critique | À trancher (PO) | UCA01 |
| [FT-024 — Chaincode myrcc absent du canal de production](FT-024-chaincode-myrcc-absent-du-canal-de-production.md) | Anomalie | critique | À vérifier | UCCL01, UCCE01, UCMOD06 |

## Ouvert

| Ticket | Type | Sévérité | Statut | Use cases |
|--------|------|----------|--------|-----------|
| [FT-001 — Propriété des assets non contrôlée côté serveur](FT-001-propriete-des-assets-non-controlee-cote-serveur.md) | Sécurité | critique | Ouvert | UCCE01, UCCE02, UCCE07, UCMOD03, UCMOD04, UCMOD08, UCPI07 |
| [FT-002 — Rôle de session REST codé en dur à contributor](FT-002-role-de-session-rest-code-en-dur.md) | Sécurité | critique | Ouvert | UCA02, UCA05, UCA07 |
| [FT-004 — Échec de l'auto-enregistrement CA non journalisé](FT-004-echec-de-l-auto-enregistrement-ca-non.md) | Anomalie | majeure | Ouvert | UCA01 |
| [FT-012 — Un module soumis peut être soumis à nouveau sans fork](FT-012-un-module-soumis-peut-etre-soumis-a.md) | Écart spec ↔ code | majeure | Ouvert | UCMOD06 |
| [FT-014 — Catégorie d'asset decoupage absente](FT-014-categorie-d-asset-decoupage-absente.md) | Écart spec ↔ code | majeure | Ouvert | UCAM05, UCAM09 |
| [FT-018 — Identifiant de bloc des versions de module simulé](FT-018-identifiant-de-bloc-des-versions-de-module.md) | Anomalie | majeure | Ouvert | UCMOD06 |
| [FT-019 — Deux tests rouges dans domain/model](FT-019-deux-tests-rouges-dans-domain-model.md) | Anomalie | majeure | Ouvert | UCREC03, UCCE01 |
| [FT-027 — Intégration continue absente du dépôt](FT-027-integration-continue-absente-du-depot.md) | Incohérence documentaire | majeure | Ouvert | — |
| [FT-028 — Mot de passe serveur en clair dans le script de déploiement](FT-028-mot-de-passe-serveur-en-clair-dans.md) | Sécurité | majeure | Ouvert | — |
| [FT-030 — Parité CLI et REST incomplète](FT-030-parite-cli-et-rest-incomplete.md) | Écart spec ↔ code | majeure | Ouvert | UCADM01, UCADM03, UCADM04, UCPI01, UCPI02, UCADM07 |
| [FT-031 — Les specs citent des commandes et routes qui n'existent pas](FT-031-les-specs-citent-des-commandes-et-routes.md) | Incohérence documentaire | majeure | Ouvert | UCAM08, UCMOD01, UCMOD04, UCDEV02, UCADM07, UCA03, UCA04 |
| [FT-009 — Registrar CA unique par organisation](FT-009-registrar-ca-unique-par-organisation.md) | Risque | mineure | Ouvert | UCA01 |
| [FT-011 — Commentaires JWT obsolètes dans les handlers REST](FT-011-commentaires-jwt-obsoletes-dans-les-handlers-rest.md) | Dette technique | mineure | Ouvert | — |
| [FT-015 — Critère tag absent des interfaces](FT-015-critere-tag-absent-des-interfaces.md) | Écart spec ↔ code | mineure | Ouvert | UCAM01, UCREC02 |
| [FT-020 — LocalStorage.Upload renvoie un chemin présenté comme un hash](FT-020-localstorage-upload-renvoie-un-chemin-presente-comme.md) | Dette technique | mineure | Ouvert | UCCL03 |
| [FT-026 — Erreur Fabric indisponible renvoyée en 500 au lieu de 503](FT-026-erreur-fabric-indisponible-renvoyee-en-500-au.md) | Anomalie | mineure | Ouvert | UCCL01 |
| [FT-029 — make deploy suppose un fabric.env présent sur le serveur](FT-029-make-deploy-suppose-un-fabric-env-present.md) | Anomalie | mineure | Ouvert | — |
| [FT-032 — API_REST.md diverge des routes réellement servies](FT-032-api-rest-md-diverge-des-routes-reellement.md) | Incohérence documentaire | mineure | Ouvert | UCCE02, UCMOD03 |
| [FT-034 — Modèle de domaine : ligne RM22 décrivant un écart](FT-034-modele-de-domaine-ligne-rm22-decrivant-un.md) | Incohérence documentaire | mineure | Ouvert | UCA02 |
| [FT-035 — Roadmap et tableaux d'état d'implémentation obsolètes](FT-035-roadmap-et-tableaux-d-etat-d-implementation.md) | Dette technique | mineure | Ouvert | UCA02 |

## À vérifier

| Ticket | Type | Sévérité | Statut | Use cases |
|--------|------|----------|--------|-----------|
| [FT-024 — Chaincode myrcc absent du canal de production](FT-024-chaincode-myrcc-absent-du-canal-de-production.md) | Anomalie | critique | À vérifier | UCCL01, UCCE01, UCMOD06 |
| [FT-025 — Le serveur signale un nœud Fabric déconnecté](FT-025-le-serveur-signale-un-nud-fabric-deconnecte.md) | Anomalie | majeure | À vérifier | UCADM02 |

## À trancher (PO)

| Ticket | Type | Sévérité | Statut | Use cases |
|--------|------|----------|--------|-----------|
| [FT-003 — Le client CA du serveur REST ignore le profil réseau actif](FT-003-le-client-ca-du-serveur-rest-ignore.md) | Anomalie | critique | À trancher (PO) | UCA01 |
| [FT-005 — Aucune approbation manuelle des demandes de compte](FT-005-aucune-approbation-manuelle-des-demandes-de-compte.md) | Écart spec ↔ code | majeure | À trancher (PO) | UCA01, UCA08 |
| [FT-006 — Demandes de compte stockées sur un seul nœud](FT-006-demandes-de-compte-stockees-sur-un-seul.md) | Décision à prendre | majeure | À trancher (PO) | UCA01 |
| [FT-008 — Wallets non chiffrés au repos](FT-008-wallets-non-chiffres-au-repos.md) | Risque | majeure | À trancher (PO) | UCA02 |
| [FT-013 — Anti-plagiat sans comparaison avec les assets existants](FT-013-anti-plagiat-sans-comparaison-avec-les-assets.md) | Écart spec ↔ code | majeure | À trancher (PO) | UCCE01, UCPI06 |
| [FT-016 — Interfaces exposées d'un module sans identité d'instance](FT-016-interfaces-exposees-d-un-module-sans-identite.md) | Écart spec ↔ code | majeure | À trancher (PO) | UCMOD04, UCAM02 |
| [FT-021 — Identifiant d'un asset décomposé en module](FT-021-identifiant-d-un-asset-decompose-en-module.md) | Décision à prendre | majeure | À trancher (PO) | UCAM05, UCAM09 |
| [FT-036 — Commissions : composition seule ou aussi chaîne de dérivation](FT-036-commissions-composition-seule-ou-aussi-chaine-de.md) | Décision à prendre | majeure | À trancher (PO) | UCPI02, UCPI04, UCPI01 |
| [FT-007 — Lecture du catalogue sans compte](FT-007-lecture-du-catalogue-sans-compte.md) | Décision à prendre | mineure | À trancher (PO) | UCA01, UCCL01 |
| [FT-010 — UCA04 et UCA07 s'appuient sur une commande qui ne couvre pas leur besoin](FT-010-uca04-et-uca07-s-appuient-sur-une.md) | Incohérence documentaire | mineure | À trancher (PO) | UCA04, UCA07 |
| [FT-022 — Date de retrait des routes legacy /api/modules](FT-022-date-de-retrait-des-routes-legacy-api.md) | Décision à prendre | mineure | À trancher (PO) | UCMOD01, UCMOD04, UCMOD06 |
| [FT-033 — UCMOD05 référencé mais inexistant](FT-033-ucmod05-reference-mais-inexistant.md) | Incohérence documentaire | mineure | À trancher (PO) | UCMOD02, UCMOD03 |
| [FT-037 — Commission sur le canal de fabrication externe](FT-037-commission-sur-le-canal-de-fabrication-externe.md) | Décision à prendre | mineure | À trancher (PO) | UCAUT01, UCPI01 |
| [FT-038 — Algorithme de différence géométrique des versions 3D](FT-038-algorithme-de-difference-geometrique-des-versions-3d.md) | Décision à prendre | mineure | À trancher (PO) | UCAUT04 |

## Différé

| Ticket | Type | Sévérité | Statut | Use cases |
|--------|------|----------|--------|-----------|
| [FT-023 — Décomposition STEP assistée : prérequis non réunis](FT-023-decomposition-step-assistee-prerequis-non-reunis.md) | Écart spec ↔ code | mineure | Différé | UCAM09 |

## Résolu

| Ticket | Type | Sévérité | Statut | Use cases |
|--------|------|----------|--------|-----------|
| [FT-017 — GetModuleInterfaces n'exposait qu'une interface par asset partagé](FT-017-getmoduleinterfaces-n-exposait-qu-une-interface-par.md) | Anomalie | majeure | Résolu | UCMOD04, UCAM02 |

## Fonctionnement

### Créer un ticket

1. Copier [_Modele_ticket](_Modele_ticket.md) sous le nom `FT-NNN-titre-court.md`, avec le numéro suivant le plus grand existant (jamais de réutilisation d'un numéro).
2. Remplir le frontmatter et les sections. Un ticket décrit un seul fait, vérifiable : constat, preuve (fichier et ligne, log, commande), impact.
3. Relier : use cases, règles (RM), exigences (ENF), specs, fichiers de code, notes de fonction (`docs/code/fonctions/`), tickets liés.
4. Ajouter la ligne du ticket dans le tableau de son statut ci-dessus (la vue [Tickets.base](Tickets.base) se met à jour seule).

### Faire vivre un ticket

- Changer `statut` (et le tag `statut/…`) et `maj` dans le frontmatter, puis ajouter une ligne datée à l'**Historique** (décision prise, commit de correction, vérification serveur).
- Un ticket n'est jamais supprimé : une fois corrigé, il passe à `resolu` avec le commit en historique ; abandonné ou sans objet, il passe à `ferme` avec la raison.
- Un message de commit qui corrige un ticket cite son identifiant (ex. `fix(identity): rôle de session lu depuis le certificat (FT-002)`).

### Valeurs

| Champ | Valeurs |
|-------|---------|
| `type` | `anomalie` (Anomalie), `ecart` (Écart spec ↔ code), `securite` (Sécurité), `risque` (Risque), `dette` (Dette technique), `incoherence` (Incohérence documentaire), `decision` (Décision à prendre) |
| `statut` | `ouvert` (Ouvert), `a-trancher` (À trancher (PO)), `a-verifier` (À vérifier), `en-cours` (En cours), `differe` (Différé), `resolu` (Résolu), `ferme` (Fermé) |
| `severite` | `critique` (sécurité, perte de données ou service inutilisable), `majeure` (fonction ou règle métier non assurée), `mineure` (gêne, documentation, dette) |

### Retrouver les tickets dans Obsidian

- Par tag : `#ticket`, `#ticket/securite`, `#statut/a-trancher`, `#severite/critique`… et les tags partagés avec les specs et le code : `#uc/UCA01`, `#rm/RM19`, `#domaine/identity`. Le tag `#uc/UCA01` réunit ainsi le use case, les fonctions qui le réalisent et les tickets qui le concernent.
- Par rétroliens : un use case, une note de fonction ou un document de conception affiche dans son panneau *Rétroliens* les tickets qui le citent.
- Par la vue [Tickets.base](Tickets.base) : tableaux filtrés par statut et sévérité.
