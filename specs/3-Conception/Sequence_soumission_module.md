---
tags:
  - couche/conception
  - type/sequence
  - domaine/model
  - uc/UCAM01
  - uc/UCAM07
  - uc/UCAM08
  - uc/UCMOD01
  - uc/UCMOD02
  - uc/UCMOD06
  - rm/RM09
  - rm/RM10
  - rm/RM11
  - rm/RM12
  - rm/RM14
  - rm/RM15
  - rm/RM16
  - rm/RM17
  - rm/RM18
  - rm/RM19
---
# Séquence — Composition et soumission d'un module (D5/D6)

> Phase 3 — Arrington | Use cases : UCMOD01, UCMOD02, UCMOD06, UCAM01, UCAM07, UCAM08 | Domaine : `domain/model`

---

## 1. Objectif

Diagramme de séquence couvrant le cycle complet d'un module : création (`draft`, RM16), composition (instances + liaisons, mutable et locale), et soumission (`SubmitModule`, RM17/RM18). Complète `Architecture_Composition.md` §3 (diagramme d'états) et §9 (ADR-05).

---

## 2. Création et composition (état `draft`, entièrement local)

```plantuml
@startuml
actor "Concepteur" as C
participant "adapter in\n(CLI ou REST)" as In
participant "ModelService\n(domain/model)" as Svc
participant "ConnectionStore\n(adapters/out/localstorage)" as ConnS
participant "InterfaceStore" as IfaceS

C -> In : myr module create --name <nom>\nPOST /api/modules
In -> Svc : CreateModule(req ModuleRequest)
Svc -> Svc : Status = draft (RM16)
Svc --> In : *Model3D (module, draft, Assemblies: [])
In --> C : module créé

C -> In : myr model instance add <moduleID> <assetID>\nPOST /api/modules/:id/instances
In -> Svc : AddAssetToWorkspace(moduleID, assetID)
Svc -> Svc : créer WorkspaceInstance\n(indépendante — RM15)
Svc --> In : *Model3D (WorkspaceInstances += 1)
In --> C : instance ajoutée

C -> In : myr model link add\nPOST /api/connections
In -> Svc : AddAssemblyLink(fromIfaceID, toIfaceID, ...)
Svc -> IfaceS : GetInterface(fromIfaceID), GetInterface(toIfaceID)
IfaceS --> Svc : AssetInterface x2
Svc -> Svc : ifacesCompatible()\n(RM10 — vérification automatique,\n5 critères RM11, voir Architecture_Composition.md §5)
alt interfaces compatibles
  Svc -> Svc : vérifier usage unique (RM09)
  Svc -> ConnS : SaveConnection(Connection{Incompatible: false})
else interfaces incompatibles
  Svc -> ConnS : SaveConnection(Connection{Incompatible: true})
  note right: RM12 — jamais supprimée\nautomatiquement
end
ConnS --> Svc : ok
Svc --> In : *Connection
In --> C : liaison créée

C -> In : myr model instance remove <moduleID> <instanceID>\nDELETE /api/modules/:id/instances/:instanceID
In -> Svc : RemoveAssetFromWorkspace(moduleID, instanceID)
Svc -> ConnS : ListConnections() puis suppression en cascade\nde toutes les Connection référençant instanceID
note right: RM14 — cascade obligatoire
ConnS --> Svc : ok
Svc --> In : *Model3D (WorkspaceInstances -= 1)
In --> C : instance retirée
@enduml
```

---

## 3. Soumission (`SubmitModule` — RM17/RM18, transaction unique)

```plantuml
@startuml
actor "Concepteur" as C
participant "adapter in\n(CLI ou REST)" as In
participant "ModelService" as Svc
participant "ConnectionStore" as ConnS
participant "BlockchainPort\n(adapters/out/fabric)" as BC

C -> In : myr module submit <moduleID> [--note]\nPOST /api/modules/:id/submit
In -> Svc : SubmitModule(moduleID, note)

Svc -> Svc : len(Assemblies) > 0 ?\n(RM17)
alt Assemblies vide
  Svc --> In : erreur RM17
  In --> C : 4xx — module vide, soumission refusée
else Assemblies non vide
  Svc -> ConnS : ListConnections() (snapshot des liaisons du module)
  ConnS --> Svc : []Connection
  Svc -> Svc : construire ModuleVersion{\n  Assemblies: [...connID],\n  Hash: hash(snapshot),\n  CreatedAt: now}
  note over Svc
    RM18 — ModuleVersion immuable,
    horodatée, hashée
  end note
  Svc -> Svc : embarquer Interfaces des composants\ndu module (ADR-02, comme un asset standard)
  Svc -> BC : StoreModelRecord(Model3D{\n  Status: submitted,\n  ModuleVersions: [...prev, newVersion]})
  BC --> Svc : BlockID
  Svc --> In : *Model3D (submitted)
  In --> C : module soumis — immuable (RM19)
end
@enduml
```

---

## 4. Notes

- La composition (§2) reste entièrement locale et mutable — aucune transaction Fabric n'est déclenchée avant `SubmitModule` (ADR-05). C'est ce qui garantit les exigences de temps de réponse UCAM (`< 200 ms` / `< 300 ms`, `Conception_intro.md` §6 ADR-02).
- Une fois `submitted`, toute modification du module doit passer par un fork (RM19) — état de cette garde : `specs/2-Analyse/Analyse_des_besoins.md` § Écarts structurels connus, E5. Ce diagramme ne représente pas le chemin de fork, qui reste à concevoir au niveau service.
- Le retrait en cascade (§2, dernier bloc) est la seule opération de composition qui modifie un état déjà persisté localement (les `Connection` liées) — elle reste néanmoins purement locale tant que le module n'est pas soumis.

<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->
## Liens

**Navigation**
- [Carte des specs](../Carte_des_specs.md)

**Use cases cités**
- UCAM01 — Liaison entre interfaces : [expression](../1-Expression/UCAM-Assemblage_Module/UCAM01.md) · [analyse](../2-Analyse/UCAM-Assemblage_Module/UCAM01.md)
- UCAM07 — Choisir un asset d'accroche (Fastener) : [expression](../1-Expression/UCAM-Assemblage_Module/UCAM07.md) · [analyse](../2-Analyse/UCAM-Assemblage_Module/UCAM07.md)
- UCAM08 — Retirer une instance de composant d'un Module : [expression](../1-Expression/UCAM-Assemblage_Module/UCAM08.md) · [analyse](../2-Analyse/UCAM-Assemblage_Module/UCAM08.md)
- UCMOD01 — Créer un Module : [expression](../1-Expression/UCMOD-Module/UCMOD01.md) · [analyse](../2-Analyse/UCMOD-Module/UCMOD01.md)
- UCMOD02 — Ajouter un Module existant : [expression](../1-Expression/UCMOD-Module/UCMOD02.md) · [analyse](../2-Analyse/UCMOD-Module/UCMOD02.md)
- UCMOD06 — Soumettre un module à la blockchain : [expression](../1-Expression/UCMOD-Module/UCMOD06.md) · [analyse](../2-Analyse/UCMOD-Module/UCMOD06.md)

**Règles métier**
- [RM09 — Interface à usage unique](../1-Expression/Regles_Metier.md#3.%20Interfaces%20et%20liaisons)
- [RM10 — Vérification de compatibilité automatique](../1-Expression/Regles_Metier.md#3.%20Interfaces%20et%20liaisons)
- [RM11 — Critères de compatibilité d'interfaces](../1-Expression/Regles_Metier.md#3.%20Interfaces%20et%20liaisons)
- [RM12 — Persistance des liaisons incompatibles](../1-Expression/Regles_Metier.md#3.%20Interfaces%20et%20liaisons)
- [RM14 — Suppression en cascade des connexions](../1-Expression/Regles_Metier.md#4.%20Composition%20d'un%20Module%20%28instances%29)
- [RM15 — Instance indépendante](../1-Expression/Regles_Metier.md#4.%20Composition%20d'un%20Module%20%28instances%29)
- [RM16 — État draft](../1-Expression/Regles_Metier.md#5.%20Modules%20et%20cycle%20de%20vie%20des%20assets)
- [RM17 — Assemblage requis pour soumission (module uniquement)](../1-Expression/Regles_Metier.md#5.%20Modules%20et%20cycle%20de%20vie%20des%20assets)
- [RM18 — ModuleVersion immuable (module uniquement)](../1-Expression/Regles_Metier.md#5.%20Modules%20et%20cycle%20de%20vie%20des%20assets)
- [RM19 — Fork d'un asset soumis](../1-Expression/Regles_Metier.md#5.%20Modules%20et%20cycle%20de%20vie%20des%20assets)

**Documents cités**
- [Analyse_des_besoins](../2-Analyse/Analyse_des_besoins.md)
- [Architecture_Composition](Architecture_Composition.md)
- [Conception_intro](Conception_intro.md)

**Cité par**
- [Conception_intro](Conception_intro.md)

<!-- liens-obsidian:end -->
