// cmd/tracegen/main.go — génère la traçabilité specs ↔ code ↔ tests.
//
// Outil de développement, exécuté par « make docs-trace » (go run, rien n'est
// compilé dans les binaires livrés). Il lit specs/ et le code Go (domain/,
// adapters/, cmd/), puis :
//   - met à jour dans specs/ les sections générées (blocs « Tests » des
//     exigences, sections « Liens »), les tags dérivés du texte et la carte
//     des specs (specs/Carte_des_specs.md) ;
//   - régénère les notes Obsidian de docs/code/ et docs/tests/ (non versionnées).
//
// Usage :
//
//	go run ./cmd/tracegen           # met à jour specs/ et docs/
//	go run ./cmd/tracegen -check    # n'écrit rien, échoue si specs/ n'est pas à jour
//
// Les règles de rattachement sont décrites dans specs/3-Conception/Tracabilite.md.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	root := flag.String("root", ".", "racine du dépôt myr-core")
	check := flag.Bool("check", false, "n'écrit rien ; code de sortie 1 si les sections générées de specs/ ne sont pas à jour")
	flag.Parse()

	abs, err := filepath.Abs(*root)
	if err != nil {
		log.Fatal(err)
	}
	docs, err := loadSpecs(abs)
	if err != nil {
		log.Fatalf("lecture des specs : %v", err)
	}
	cb, err := loadCode(abs)
	if err != nil {
		log.Fatalf("analyse du code : %v", err)
	}
	m := buildModel(docs, cb, commitOf(abs))
	m.absRoot = filepath.ToSlash(abs)

	specs := m.renderSpecs()
	stale := 0
	for _, p := range sortedKeys(specs) {
		file := filepath.Join(abs, filepath.FromSlash(p))
		old, _ := os.ReadFile(file)
		if string(old) == specs[p] {
			continue
		}
		stale++
		if *check {
			fmt.Println("à régénérer :", p)
			continue
		}
		if err := os.WriteFile(file, []byte(specs[p]), 0o644); err != nil {
			log.Fatal(err)
		}
	}
	if *check {
		if stale > 0 {
			fmt.Printf("%d fichier(s) de specs/ à régénérer : make docs-trace\n", stale)
			os.Exit(1)
		}
		fmt.Println("specs/ à jour")
		return
	}

	notes := m.renderCode()
	for p, s := range m.renderTests() {
		notes[p] = s
	}
	removed := cleanGenerated(abs, notes)
	for _, p := range sortedKeys(notes) {
		file := filepath.Join(abs, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(file, []byte(notes[p]), 0o644); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("specs/ : %d fichier(s) mis à jour · docs/ : %d note(s) écrite(s), %d obsolète(s) supprimée(s)\n", stale, len(notes), removed)
	fmt.Printf("%d use cases · %d exigences · %d flux · %d fonctions · %d routes · %d commandes · %d tests\n",
		len(m.ucs), len(m.rms)+len(m.enfs)+len(m.efs), len(m.flows), len(cb.fns), len(cb.routeNotes), len(cb.cmds), len(cb.tests))
}

// commitOf retourne le commit analysé, signalé si l'arbre de travail diffère.
func commitOf(root string) string {
	out, err := exec.Command("git", "-C", root, "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return "inconnu"
	}
	c := strings.TrimSpace(string(out))
	dirty, _ := exec.Command("git", "-C", root, "status", "--porcelain", "--", "domain", "adapters", "cmd", "specs").Output()
	if len(strings.TrimSpace(string(dirty))) > 0 {
		c += " + modifications locales"
	}
	return c
}

// cleanGenerated supprime de docs/code et docs/tests les notes générées qui
// ne seront pas réécrites (test supprimé, fonction renommée…). Seuls les
// fichiers portant la mention de génération sont concernés.
func cleanGenerated(root string, keep map[string]string) int {
	n := 0
	for _, dir := range []string{"docs/code", "docs/tests"} {
		var files []string
		_ = filepath.WalkDir(filepath.Join(root, dir), func(p string, d os.DirEntry, err error) error {
			if err == nil && !d.IsDir() && strings.HasSuffix(p, ".md") {
				files = append(files, p)
			}
			return nil
		})
		sort.Strings(files)
		for _, p := range files {
			rel, _ := filepath.Rel(root, p)
			if _, ok := keep[filepath.ToSlash(rel)]; ok {
				continue
			}
			if raw, err := os.ReadFile(p); err == nil && strings.Contains(string(raw), generatedNote) {
				if os.Remove(p) == nil {
					n++
				}
			}
		}
	}
	return n
}
