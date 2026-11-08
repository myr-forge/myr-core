package main

import (
	"path"
	"regexp"
	"strings"
)

// Marqueurs des sections générées insérées dans les specs. Le texte qui suit
// le marqueur d'ouverture rappelle qu'elles ne s'éditent pas à la main.
const (
	testsBegin = "<!-- tests-obsidian:begin — tests rattachés à cette exigence, section générée : ne pas l'éditer à la main -->"
	testsEnd   = "<!-- tests-obsidian:end -->"
	liensBegin = "<!-- liens-obsidian:begin — section de traçabilité générée à partir des références du document : ne pas l'éditer à la main -->"
	liensEnd   = "<!-- liens-obsidian:end -->"
)

// generatedNote ouvre chaque note de docs/ et signale qu'elle est régénérée.
const generatedNote = "> Note de traçabilité Obsidian générée automatiquement à partir des specs et du code — ne pas l'éditer à la main."

var headingRe = regexp.MustCompile(`^(#{1,6})\s+(.*?)\s*$`)

// anchor encode un titre de section pour un lien Obsidian (fichier.md#Titre) :
// seuls les espaces et les parenthèses sont échappés.
func anchor(title string) string {
	r := strings.NewReplacer(" ", "%20", "(", "%28", ")", "%29")
	return r.Replace(title)
}

// escPath échappe les espaces d'un chemin de fichier dans un lien Markdown.
func escPath(p string) string {
	return strings.ReplaceAll(p, " ", "%20")
}

// relLink retourne le chemin relatif de to vu depuis le fichier from (chemins
// relatifs à la racine du dépôt, séparateur /).
func relLink(from, to string) string {
	fromDir := strings.Split(path.Dir(from), "/")
	if path.Dir(from) == "." {
		fromDir = nil
	}
	toParts := strings.Split(to, "/")
	i := 0
	for i < len(fromDir) && i < len(toParts)-1 && fromDir[i] == toParts[i] {
		i++
	}
	var out []string
	for range fromDir[i:] {
		out = append(out, "..")
	}
	out = append(out, toParts[i:]...)
	return strings.Join(out, "/")
}

// link formate un lien Markdown de from vers to, avec un fragment facultatif.
func link(text, from, to, frag string) string {
	target := escPath(relLink(from, to))
	if frag != "" {
		target += "#" + anchor(frag)
	}
	return "[" + text + "](" + target + ")"
}

// code formate une valeur en code Markdown.
func code(s string) string { return "`" + s + "`" }

// removeBlocks retire du texte les blocs générés (marqueurs inclus) et la
// ligne vide qui les précédait, pour que leur régénération reste idempotente.
func removeBlocks(lines []string) []string {
	var out []string
	in, removed := false, false
	for _, l := range lines {
		t := strings.TrimSpace(l)
		switch {
		case strings.HasPrefix(t, "<!-- tests-obsidian:begin"), strings.HasPrefix(t, "<!-- liens-obsidian:begin"):
			in = true
			for len(out) > 0 && strings.TrimSpace(out[len(out)-1]) == "" {
				out = out[:len(out)-1]
			}
			out = append(out, "")
		case t == testsEnd || t == liensEnd:
			in, removed = false, true
		case !in:
			// Seules les lignes vides qui suivent immédiatement un bloc retiré
			// sont fusionnées ; le reste du texte est conservé tel quel.
			if removed && t == "" {
				continue
			}
			removed = false
			out = append(out, l)
		}
	}
	return out
}

// isFence indique si la ligne ouvre ou ferme un bloc de code.
func isFence(l string) bool {
	t := strings.TrimSpace(l)
	return strings.HasPrefix(t, "```") || strings.HasPrefix(t, "~~~")
}

// stripInlineCode remplace le contenu des spans de code par des espaces, pour
// ignorer ce qui y est cité (annotations de relecture mentionnées en exemple).
func stripInlineCode(s string) string {
	var b strings.Builder
	in := false
	for _, r := range s {
		if r == '`' {
			in = !in
			b.WriteRune(' ')
			continue
		}
		if in {
			b.WriteRune(' ')
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
