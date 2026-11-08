package main

import (
	"fmt"
	"sort"
	"strings"
)

var layerLabels = map[string]string{
	"point-entree": "Point d'entrée (assemblage des adapters)",
	"adapter-in":   "Adapter entrant (déclenche le domaine)",
	"domaine":      "Domaine (cœur métier)",
	"adapter-out":  "Adapter sortant (infrastructure)",
}

var layerOrder = []string{"point-entree", "adapter-in", "domaine", "adapter-out"}

// noteWriter accumule le contenu d'une note de docs/.
type noteWriter struct {
	m    *model
	path string
	b    strings.Builder
}

func (m *model) note(path, title string, tags []string) *noteWriter {
	w := &noteWriter{m: m, path: path}
	w.line("---\ntags:")
	for _, t := range tags {
		w.line("  - " + t)
	}
	w.line("---")
	w.line("# " + title + "\n")
	w.line(generatedNote)
	w.line("> État du code analysé : commit " + code(m.commit) + " — les numéros de ligne s'y rapportent.\n")
	return w
}

func (w *noteWriter) line(s string)                     { w.b.WriteString(s + "\n") }
func (w *noteWriter) f(format string, a ...any)         { w.line(fmt.Sprintf(format, a...)) }
func (w *noteWriter) link(text, to, frag string) string { return link(text, w.path, to, frag) }

// loc formate un emplacement de code avec son lien d'ouverture dans VS Code.
func (w *noteWriter) loc(file string, line int) string {
	if file == "" {
		return "—"
	}
	return fmt.Sprintf("`%s:%d` [↗ VS Code](vscode://file/%s/%s:%d)", file, line, w.m.absRoot, file, line)
}

func (w *noteWriter) pkgLink(p *pkg) string { return w.link(p.Path, pkgNote(p), "") }

func (w *noteWriter) ucList(ids []string) string {
	var out []string
	for _, id := range ids {
		out = append(out, code(id)+" "+w.m.ucByID[id].Title)
	}
	return dash(strings.Join(out, " · "))
}

func docLabels(ds []*doc) string {
	sort.Slice(ds, func(i, j int) bool { return docOrder(ds[i]) < docOrder(ds[j]) })
	var out []string
	for _, d := range ds {
		out = append(out, code(d.Label))
	}
	return dash(strings.Join(out, ", "))
}

// docOrder classe les documents par couche puis par chemin.
func docOrder(d *doc) string {
	rank := map[string]string{"expression": "1", "analyse": "2", "conception": "3", "transverse": "4"}
	return rank[d.Layer] + d.Path
}

func ucTags(ids []string) []string {
	var out []string
	for _, id := range ids {
		out = append(out, "uc/"+id)
	}
	return out
}

// renderCode produit les notes de docs/code.
func (m *model) renderCode() map[string]string {
	out := map[string]string{}
	for _, p := range m.cb.pkgs {
		out[pkgNote(p)] = m.pkgPage(p)
	}
	for _, f := range m.cb.fns {
		out[fnNote(f)] = m.fnPage(f)
	}
	for _, rn := range m.cb.routeNotes {
		out[routeNotePath(rn)] = m.routePage(rn)
	}
	for _, c := range m.cb.cmds {
		out[cmdNote(c)] = m.cmdPage(c)
	}
	out["docs/code/Architecture.md"] = m.architecturePage()
	out["docs/code/Tracabilite_UC_Code.md"] = m.traceabilityPage()
	return out
}

// pkgDocs retourne les specs qui citent le chemin du package ou ses éléments.
func (m *model) pkgDocs(p *pkg) []*doc {
	var out []*doc
	for _, d := range m.docs {
		if strings.Contains(strings.Join(d.body, "\n"), p.Path) {
			addDoc(&out, d)
		}
	}
	for _, pt := range p.Ports {
		for _, f := range pt.Methods {
			for _, d := range f.Docs {
				addDoc(&out, d)
			}
		}
	}
	if p.Path == "adapters/in/rest" {
		for _, rn := range m.cb.routeNotes {
			for _, d := range rn.Docs {
				addDoc(&out, d)
			}
		}
	}
	if p.Path == "adapters/in/cli" {
		for _, c := range m.cb.cmds {
			for _, d := range c.Docs {
				addDoc(&out, d)
			}
		}
	}
	return out
}

// pkgUCs retourne les use cases liés au package : leurs specs citent une de
// ses fonctions, routes ou commandes, ou un de ses tests cite le use case.
func (m *model) pkgUCs(p *pkg) []string {
	set := map[string]bool{}
	for _, pt := range p.Ports {
		for _, f := range pt.Methods {
			for _, u := range m.ucsOfFn(f) {
				set[u] = true
			}
		}
	}
	if p.Path == "adapters/in/rest" {
		for _, rn := range m.cb.routeNotes {
			for _, u := range m.ucsOfRoute(rn) {
				set[u] = true
			}
		}
	}
	if p.Path == "adapters/in/cli" {
		for _, c := range m.cb.cmds {
			for _, u := range m.ucsOfCmd(c) {
				set[u] = true
			}
		}
	}
	for _, t := range p.Tests {
		for id := range t.IDs {
			if m.ucByID[id] != nil {
				set[id] = true
			}
		}
	}
	return keys(set)
}

func (m *model) pkgPage(p *pkg) string {
	ucs := m.pkgUCs(p)
	tags := []string{"couche/code", "type/package", "hexa/" + p.Layer}
	if p.Layer == "domaine" {
		tags = append(tags, "domaine/"+p.Name)
	}
	w := m.note(pkgNote(p), p.Path, append(tags, ucTags(ucs)...))
	w.f("**Couche** : %s — voir %s\n", layerLabels[p.Layer], w.link("Architecture du code", "docs/code/Architecture.md", ""))
	w.line("## Dépendances\n")
	w.f("- **Importe** : %s", w.pkgLinks(keys(p.Imports)))
	w.f("- **Importé par** : %s\n", w.pkgLinks(keys(p.ImportedBy)))
	w.line("## Specs qui décrivent ce package\n")
	w.line(docLabels(m.pkgDocs(p)) + "\n")
	for _, pt := range p.Ports {
		kind := "Port de sortie "
		if pt.In {
			kind = "Port d'entrée "
		}
		w.line("## " + kind + pt.Name + "\n")
		w.f("Déclaré en %s. Implémenté par %s.", code(fmt.Sprintf("%s:%d", pt.File, pt.Line)), w.impls(pt))
		if pt.In {
			w.line("Chaque méthode ouvre sa note de fonction (code source, use cases impactés, tests).\n")
			w.line("| Méthode | Implémentation | Appelée depuis | Specs qui la citent |\n|---------|----------------|----------------|---------------------|")
			for _, f := range pt.Methods {
				impl := "—"
				if f.ImplFile != "" {
					impl = code(fmt.Sprintf("%s:%d", f.ImplFile, f.ImplLine))
				}
				w.f("| %s | %s | %s | %s |", w.link(f.Name, fnNote(f), ""), impl, callers(f), docLabels(f.Docs))
			}
		} else {
			w.line("")
			w.line("| Méthode | Specs qui la citent |\n|---------|---------------------|")
			for _, f := range pt.Methods {
				w.f("| %s | %s |", w.link(f.Name, fnNote(f), ""), docLabels(f.Docs))
			}
		}
		w.line("")
	}
	if p.Path == "adapters/in/rest" {
		w.line("## Routes REST\n")
		w.line("| Route | Handler | Fonctions du domaine appelées |\n|-------|---------|-------------------------------|")
		for _, rn := range m.cb.routeNotes {
			var hs []string
			for _, r := range rn.Regs {
				hs = append(hs, code(r.Handler))
			}
			w.f("| %s | %s | %s |", w.link(rn.Title, routeNotePath(rn), ""), strings.Join(hs, " · "), w.fnLinks(rn.calls()))
		}
		w.line("")
	}
	if p.Path == "adapters/in/cli" {
		w.line("## Commandes CLI\n")
		w.line("| Commande | Définition | Fonctions du domaine appelées |\n|----------|------------|-------------------------------|")
		for _, c := range m.cb.cmds {
			w.f("| %s | %s | %s |", w.link(c.Name(), cmdNote(c), ""), code(fmt.Sprintf("%s:%d", c.File, c.Line)), w.fnLinks(c.Calls))
		}
		w.line("")
	}
	if p.Layer == "adapter-out" {
		var impls []string
		for _, pt := range m.cb.ports {
			for _, im := range pt.Impls {
				if im.Pkg == p {
					kind := "Port de sortie "
					if pt.In {
						kind = "Port d'entrée "
					}
					impls = append(impls, "- "+w.link(pt.Pkg.Name+"."+pt.Name, pkgNote(pt.Pkg), kind+pt.Name)+" — "+code(im.Type))
				}
			}
		}
		w.line("## Ports implémentés\n")
		if len(impls) == 0 {
			impls = []string{"_Aucun port du domaine implémenté._"}
		}
		for _, l := range impls {
			w.line(l)
		}
		w.line("")
	}
	if len(ucs) > 0 {
		w.line("## Use cases liés\n")
		w.line("Use cases dont la spec cite une méthode, une route ou une commande de ce package (ou dont l'identifiant figure dans un test).\n")
		for _, id := range ucs {
			w.f("- %s %s", code(id), m.ucByID[id].Title)
		}
		w.line("")
	}
	w.line("## Tests\n")
	if len(p.Tests) == 0 {
		w.line("_Aucun test._\n")
	} else {
		w.line("| Test | Règles / exigences / UC | Ce qu'il vérifie |\n|------|-------------------------|------------------|")
		for _, t := range sortTests(p.Tests) {
			w.f("| %s — %s | %s | %s |", w.link(t.Name, m.testNote(t), ""), w.loc(t.File, t.Line), dash(strings.Join(sortedIDs(t.IDs), ", ")), cell(t.Brief))
		}
		w.line("")
	}
	w.line("## Fichiers\n")
	w.line("| Fichier | Lignes |\n|---------|--------|")
	for _, fs := range p.Files {
		w.f("| %s | %d |", w.loc(fs.Path, 1), fs.Lines)
	}
	return w.b.String()
}

func cell(s string) string { return dash(strings.ReplaceAll(strings.TrimSpace(s), "|", "\\|")) }

func sortedIDs(ids map[string]bool) []string { return keys(ids) }

func sortTests(ts []*test) []*test {
	out := append([]*test{}, ts...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func (w *noteWriter) pkgLinks(paths []string) string {
	var out []string
	for _, p := range paths {
		out = append(out, w.pkgLink(w.m.cb.byPath[p]))
	}
	return dash(strings.Join(out, " · "))
}

func (w *noteWriter) fnLinks(fs []*fn) string {
	var out []string
	for _, f := range fs {
		out = append(out, w.link(f.Label(), fnNote(f), ""))
	}
	return dash(strings.Join(out, " · "))
}

func (w *noteWriter) impls(pt *port) string {
	var out []string
	for _, im := range pt.Impls {
		if im.Pkg == pt.Pkg {
			out = append(out, code(im.Type)+" (ce package)")
		} else {
			out = append(out, w.pkgLink(im.Pkg)+" ("+code(im.Type)+")")
		}
	}
	if len(out) == 0 {
		return "aucun type du module"
	}
	return strings.Join(out, " · ")
}

// callers liste les routes et commandes qui appellent une fonction.
func callers(f *fn) string {
	var out []string
	for _, c := range f.Cmds {
		out = append(out, "CLI "+code(c.Name()))
	}
	for _, r := range f.Routes {
		out = append(out, "REST "+code(r.Title))
	}
	return dash(strings.Join(out, " · "))
}

func (m *model) fnPage(f *fn) string {
	ucs := m.ucsOfFn(f)
	dir := "port/sortie"
	kind := "Port de sortie "
	if f.Port.In {
		dir, kind = "port/entree", "Port d'entrée "
	}
	tags := append([]string{"couche/code", "type/fonction", "hexa/domaine", "domaine/" + f.Port.Pkg.Name, dir}, ucTags(ucs)...)
	w := m.note(fnNote(f), f.Label(), tags)
	w.f("**Package** : %s · **Port** : %s\n", w.pkgLink(f.Port.Pkg), w.link(kind+f.Port.Name, pkgNote(f.Port.Pkg), kind+f.Port.Name))
	w.line("```go\n" + f.Sig + "\n```\n")
	if f.Doc != "" {
		w.line("> " + f.Doc + "\n")
	}
	w.line("## Code source\n")
	w.f("- **Déclaration** : %s", w.loc(f.DeclFile, f.DeclLine))
	if f.ImplFile != "" {
		label := "Implémentation (domaine)"
		if !f.Port.In {
			label = "Implémentation"
		}
		w.f("- **%s** : %s", label, w.loc(f.ImplFile, f.ImplLine))
	}
	w.line("")
	w.line("## Appelée depuis\n")
	if len(f.Routes)+len(f.Cmds) == 0 {
		w.line("_Aucune route ni commande ne l'appelle (détection par nom de méthode)._")
	}
	for _, r := range f.Routes {
		w.line("- " + w.link("REST "+r.Title, routeNotePath(r), ""))
	}
	for _, c := range f.Cmds {
		w.line("- " + w.link("CLI "+c.Name(), cmdNote(c), ""))
	}
	w.line("")
	w.line("## Use cases impactés\n")
	if len(ucs) == 0 {
		w.line("_Aucun use case ne la sollicite._\n")
	} else {
		w.line("Relation affichée en texte (et portée par les tags `uc/…`) : dans le graphe, un use case rejoint le code par sa conception.\n")
		for _, id := range ucs {
			w.f("- %s %s — _%s_", code(id), m.ucByID[id].Title, strings.Join(m.ucCode[id].Why[f], " ; "))
		}
		w.line("")
	}
	w.line("## Tests qui l'appellent\n")
	if len(f.Tests) == 0 {
		w.line("_Aucun test ne l'appelle._\n")
	}
	for _, t := range sortTests(f.Tests) {
		w.f("- %s — %s — %s", w.link(t.Name, m.testNote(t), ""), w.loc(t.File, t.Line), cell(t.Brief))
	}
	if len(f.Tests) > 0 {
		w.line("")
	}
	w.line("## Specs qui la citent\n")
	w.line(docLabels(f.Docs))
	return w.b.String()
}

func (m *model) routePage(rn *routeNote) string {
	ucs := m.ucsOfRoute(rn)
	w := m.note(routeNotePath(rn), "REST "+rn.Title, append([]string{"couche/code", "type/route", "hexa/adapter-in"}, ucTags(ucs)...))
	w.f("**Package** : %s\n", w.pkgLink(m.cb.byPath["adapters/in/rest"]))
	w.line("## Code source\n")
	for _, r := range rn.Regs {
		w.f("- **Enregistrement de la route** %s : %s", code(r.Pattern), w.loc(r.File, r.Line))
		w.f("- **Handler** %s : %s", code(r.Handler), w.loc(r.HFile, r.HLine))
	}
	w.line("")
	w.line("## Fonctions du domaine appelées\n")
	w.line("Appels trouvés en suivant le handler et les fonctions qu'il appelle dans le package (type du récepteur, sinon nom de méthode unique).\n")
	calls := rn.calls()
	if len(calls) == 0 {
		w.line("_Aucun appel à un port du domaine détecté._")
	}
	for _, f := range calls {
		w.line("- " + w.link(f.Label(), fnNote(f), ""))
	}
	w.line("")
	m.testsSection(w, rn.Tests, "Tests qui l'appellent")
	w.line("## Use cases et specs qui la citent\n")
	w.f("- Use cases : %s", w.ucList(ucs))
	w.f("- Specs : %s", docLabels(rn.Docs))
	return w.b.String()
}

func (m *model) cmdPage(c *command) string {
	ucs := m.ucsOfCmd(c)
	w := m.note(cmdNote(c), "CLI "+c.Name(), append([]string{"couche/code", "type/commande", "hexa/adapter-in"}, ucTags(ucs)...))
	head := "**Package** : " + w.pkgLink(m.cb.byPath["adapters/in/cli"])
	if c.Parent != nil && len(c.Parent.Words) > 0 {
		head += " · **Commande parente** : " + w.link("CLI "+c.Parent.Name(), cmdNote(c.Parent), "")
	}
	w.line(head + "\n")
	if len(c.Children) > 0 {
		var subs []string
		for _, s := range c.Children {
			if len(s.Words) > 0 {
				subs = append(subs, w.link("CLI "+s.Name(), cmdNote(s), ""))
			}
		}
		w.line("**Sous-commandes** : " + strings.Join(subs, " · ") + "\n")
	}
	w.line("## Code source\n")
	w.f("- **Définition** : %s\n", w.loc(c.File, c.Line))
	w.line("## Fonctions du domaine appelées\n")
	if len(c.Calls) == 0 {
		w.line("_Aucun appel à un port du domaine détecté._")
	}
	for _, f := range c.Calls {
		w.line("- " + w.link(f.Label(), fnNote(f), ""))
	}
	w.line("")
	m.testsSection(w, c.Tests, "Tests qui l'exercent")
	w.line("## Use cases et specs qui la citent\n")
	w.f("- Use cases : %s", w.ucList(ucs))
	w.f("- Specs : %s", docLabels(c.Docs))
	return w.b.String()
}

func (m *model) testsSection(w *noteWriter, ts []*test, title string) {
	w.line("## " + title + "\n")
	if len(ts) == 0 {
		w.line("_Aucun test._\n")
		return
	}
	for _, t := range sortTests(ts) {
		w.f("- %s — %s", w.link(t.Name, m.testNote(t), ""), cell(t.Brief))
	}
	w.line("")
}

func mermaidID(p *pkg) string { return strings.NewReplacer("/", "_", "-", "_").Replace(p.Path) }

func (m *model) architecturePage() string {
	w := m.note("docs/code/Architecture.md", "Architecture du code", []string{"couche/code", "type/architecture"})
	w.line("Vue des packages Go de `myr-core` selon l'architecture hexagonale : les points d'entrée assemblent les adapters, les adapters entrants (CLI, REST) appellent les ports d'entrée du domaine, et les adapters sortants implémentent ses ports de sortie.\n")
	w.f("Voir aussi : %s · %s\n", w.link("Carte des specs", "specs/Carte_des_specs.md", ""), w.link("Traçabilité use cases ↔ code", "docs/code/Tracabilite_UC_Code.md", ""))
	w.line("## Dépendances entre packages\n")
	w.line("```mermaid\nflowchart LR")
	for _, l := range layerOrder {
		w.f("  subgraph %s[\"%s\"]", strings.ReplaceAll(l, "-", "_"), layerLabels[l])
		for _, p := range m.cb.pkgs {
			if p.Layer == l {
				w.f("    %s[\"%s\"]", mermaidID(p), p.Path)
			}
		}
		w.line("  end")
	}
	for _, p := range m.cb.pkgs {
		for _, imp := range keys(p.Imports) {
			w.f("  %s --> %s", mermaidID(p), mermaidID(m.cb.byPath[imp]))
		}
	}
	w.line("```\n")
	w.line("## Packages\n")
	w.line("| Package | Couche | Ports | Routes | Commandes | Tests |\n|---------|--------|-------|--------|-----------|-------|")
	for _, l := range layerOrder {
		for _, p := range m.cb.pkgs {
			if p.Layer != l {
				continue
			}
			routes, cmds := 0, 0
			if p.Path == "adapters/in/rest" {
				routes = len(m.cb.routeNotes)
			}
			if p.Path == "adapters/in/cli" {
				cmds = len(m.cb.cmds)
			}
			w.f("| %s | %s | %d | %d | %d | %d |", w.pkgLink(p), layerLabels[l], len(p.Ports), routes, cmds, len(p.Tests))
		}
	}
	return w.b.String()
}

func (m *model) traceabilityPage() string {
	const path = "docs/code/Tracabilite_UC_Code.md"
	w := m.note(path, "Traçabilité use cases ↔ code", []string{"couche/code", "type/tracabilite"})
	w.line("Pour chaque use case : les routes REST et commandes CLI **citées explicitement** par sa spec (couches expression et analyse), et les fonctions du domaine qu'il sollicite — citées dans la spec, ou appelées par une commande ou une route qu'elle cite. Chaque fonction ouvre sa note (code source, autres use cases impactés, tests) ; les liens « ouvrir dans VS Code » ouvrent le fichier à la ligne.\n")
	w.f("Voir aussi : %s · %s\n", w.link("Carte des specs", "specs/Carte_des_specs.md", ""), w.link("Architecture du code", "docs/code/Architecture.md", ""))
	w.line("## Synthèse\n")
	w.line("| UC | Titre | Routes | Commandes | Fonctions |\n|----|-------|--------|-----------|-----------|")
	for _, u := range m.ucs {
		c := m.ucCode[u.ID]
		w.f("| %s | %s | %d | %d | %d |", w.link(u.ID, path, u.ID), u.Title, len(c.Routes), len(c.Cmds), len(c.Fns))
	}
	w.line("")
	fam := ""
	for _, u := range m.ucs {
		if u.Family != fam {
			fam = u.Family
			w.line("## " + u.FamilyLabel() + "\n")
		}
		w.line("### " + u.ID + "\n")
		var layers []string
		if u.Expr != nil {
			layers = append(layers, w.link("expression", u.Expr.Path, ""))
		}
		if u.Ana != nil {
			layers = append(layers, w.link("analyse", u.Ana.Path, ""))
		}
		w.line(u.Title + " — " + strings.Join(layers, " · ") + "\n")
		c := m.ucCode[u.ID]
		if len(c.Routes)+len(c.Cmds)+len(c.Fns) == 0 {
			w.line("_Aucune route, commande ni fonction du code n'est citée par ce use case._\n")
			continue
		}
		for _, rh := range c.Routes {
			w.f("- **REST** %s → %s (%s)", code(rh.Cite), w.link("route "+rh.Reg.Note.Title, routeNotePath(rh.Reg.Note), ""), w.pkgLink(m.cb.byPath["adapters/in/rest"]))
		}
		for _, cm := range c.Cmds {
			w.f("- **CLI** %s — %s (%s)", w.link(cm.Name(), cmdNote(cm), ""), w.loc(cm.File, cm.Line), w.pkgLink(m.cb.byPath["adapters/in/cli"]))
		}
		for _, f := range c.Fns {
			w.f("- **Fonction** %s — %s — _%s_", w.link(f.Label(), fnNote(f), ""), w.loc(f.ImplFile, f.ImplLine), strings.Join(c.Why[f], " ; "))
		}
		w.line("")
	}
	w.line("## Règles métier → tests\n")
	w.line("Tests dont le nom ou le commentaire cite la règle.\n")
	w.line("| RM | Règle | Tests |\n|----|-------|-------|")
	for _, r := range m.rms {
		var ts []string
		for _, h := range m.cov[r].Explicit {
			ts = append(ts, w.link(h.T.Name, m.testNote(h.T), "")+" ("+w.pkgLink(h.T.Subject)+")")
		}
		w.f("| %s | %s | %s |", w.link(r.ID, r.Doc.Path, r.Sec.Title), r.Title, dash(strings.Join(ts, "<br>")))
	}
	w.line("")
	w.line("## Références des specs sans correspondance dans le code\n")
	w.line("Commandes et routes citées par les specs qui ne correspondent à aucune commande cobra ni route enregistrée : écart entre spec et code, ou citation à corriger.\n")
	for _, g := range []struct {
		title string
		set   map[string]map[*doc]bool
	}{{"Commandes CLI", m.strayCmds}, {"Routes REST", m.strayRoutes}} {
		w.line("### " + g.title + "\n")
		if len(g.set) == 0 {
			w.line("_Aucune._\n")
			continue
		}
		w.line("| Citation | Specs |\n|----------|-------|")
		for _, cite := range sortedKeys(g.set) {
			var ds []*doc
			for d := range g.set[cite] {
				ds = append(ds, d)
			}
			sort.Slice(ds, func(i, j int) bool { return docOrder(ds[i]) < docOrder(ds[j]) })
			var links []string
			for _, d := range ds {
				links = append(links, w.link(d.Label, d.Path, ""))
			}
			w.f("| %s | %s |", code(cite), strings.Join(links, " · "))
		}
		w.line("")
	}
	return w.b.String()
}
