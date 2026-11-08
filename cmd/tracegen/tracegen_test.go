package main

import (
	"reflect"
	"strings"
	"testing"
)

// / @brief  anchor encode un titre de section pour un lien Obsidian
// / @input  titre avec espaces, parenthèses, backticks et tiret long
// / @expect espaces → %20, parenthèses → %28/%29, le reste inchangé
func TestAnchor_EscapesSpacesAndParens(t *testing.T) {
	got := anchor("Flux nominal — Accès (`AllowAutoRegister=true`)")
	want := "Flux%20nominal%20—%20Accès%20%28`AllowAutoRegister=true`%29"
	if got != want {
		t.Errorf("anchor = %q, attendu %q", got, want)
	}
}

// / @brief  relLink calcule le chemin relatif entre deux fichiers du dépôt
// / @input  une analyse de use case vers une note de test, une conception vers elle-même
// / @expect ../../../docs/tests/… et le seul nom de fichier
func TestRelLink_BetweenFiles(t *testing.T) {
	cases := [][3]string{
		{"specs/2-Analyse/UCA-Compte_et_Acces/UCA01.md", "docs/tests/domain-identity/T.md", "../../../docs/tests/domain-identity/T.md"},
		{"specs/3-Conception/API_REST.md", "specs/3-Conception/API_REST.md", "API_REST.md"},
		{"specs/Carte_des_specs.md", "specs/1-Expression/Regles_Metier.md", "1-Expression/Regles_Metier.md"},
	}
	for _, c := range cases {
		if got := relLink(c[0], c[1]); got != c[2] {
			t.Errorf("relLink(%q, %q) = %q, attendu %q", c[0], c[1], got, c[2])
		}
	}
}

// / @brief  expandUCs développe les plages et les suffixes de use cases
// / @input  « UCA01–03, UCAM01–03/05/07 et UCMOD06 »
// / @expect UCA01 à UCA03, UCAM01 à UCAM03, UCAM05, UCAM07, UCMOD06
func TestExpandUCs_RangesAndSuffixes(t *testing.T) {
	got := uniqueSorted(expandUCs("UCA01–03, UCAM01–03/05/07 et UCMOD06"))
	want := []string{"UCA01", "UCA02", "UCA03", "UCAM01", "UCAM02", "UCAM03", "UCAM05", "UCAM07", "UCMOD06"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expandUCs = %v, attendu %v", got, want)
	}
}

// / @brief  removeBlocks retire les sections générées sans toucher au reste du texte
// / @input  texte avec un bloc de tests et deux lignes vides consécutives ailleurs
// / @expect le bloc disparaît, les deux lignes vides hors bloc sont conservées
func TestRemoveBlocks_KeepsOtherText(t *testing.T) {
	in := []string{"### RM01 — Règle", "texte", "", testsBegin, "> [!warning] Tests", testsEnd, "", "### RM02 — Autre", "a", "", "", "b"}
	got := removeBlocks(in)
	want := []string{"### RM01 — Règle", "texte", "", "### RM02 — Autre", "a", "", "", "b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("removeBlocks = %q, attendu %q", got, want)
	}
}

// / @brief  extract relève routes, commandes, méthodes et annotations d'un texte
// / @input  lignes citant une route avec méthode, une commande en code suivie de prose, une méthode qualifiée, une annotation et une annotation citée en code
// / @expect route POST /api/identity/request ; commande « identity request » limitée au span de code ; IdentityService.Enroll ; annotation remarque seule
func TestExtract_Citations(t *testing.T) {
	c := extract([]string{
		"Le client appelle `POST /api/identity/request` puis `myr identity request --pseudo x` avec le secret.",
		"Voir `IdentityService.Enroll` #remarque à préciser — l'annotation `#question` n'est qu'un exemple.",
	})
	if len(c.Routes) != 1 || c.Routes[0] != (routeCite{Method: "POST", Path: "/api/identity/request"}) {
		t.Errorf("routes = %+v", c.Routes)
	}
	if len(c.Cmds) != 1 || strings.Join(c.Cmds[0].Words, " ") != "identity request" || !c.Cmds[0].Quoted {
		t.Errorf("commandes = %+v", c.Cmds)
	}
	found := false
	for _, f := range c.Funcs {
		if f == (funcCite{Qual: "IdentityService", Name: "Enroll"}) {
			found = true
		}
	}
	if !found {
		t.Errorf("IdentityService.Enroll non relevé : %+v", c.Funcs)
	}
	if !c.Notes["remarque"] || c.Notes["question"] {
		t.Errorf("annotations = %v, attendu remarque seule", c.Notes)
	}
}
