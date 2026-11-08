package main

import (
	"fmt"
	"sort"
	"strings"
)

const matrixPath = "docs/tests/Matrice_Couverture_Tests.md"

// maxBlockTests borne le nombre de tests listés dans un bloc d'exigence.
const maxBlockTests = 15

// Chemins des notes de code et de test.
func pkgSlug(p *pkg) string             { return strings.ReplaceAll(p.Path, "/", "-") }
func pkgNote(p *pkg) string             { return "docs/code/" + pkgSlug(p) + ".md" }
func fnNote(f *fn) string               { return "docs/code/fonctions/" + f.Key() + ".md" }
func routeNotePath(r *routeNote) string { return "docs/code/routes/REST " + r.Slug + ".md" }
func cmdNote(c *command) string {
	return "docs/code/commandes/CLI " + strings.Join(c.Words, "-") + ".md"
}

// testNote retourne la note d'un test ; deux tests homonymes du même package
// (tests unitaires et d'intégration) sont distingués par leur fichier.
func (m *model) testNote(t *test) string {
	if m.testPaths == nil {
		m.testPaths = map[*test]string{}
		used := map[string]bool{}
		for _, x := range m.cb.tests {
			p := "docs/tests/" + pkgSlug(x.Subject) + "/" + x.Name + ".md"
			if used[p] {
				base := strings.TrimSuffix(x.File[strings.LastIndex(x.File, "/")+1:], "_test.go")
				p = "docs/tests/" + pkgSlug(x.Subject) + "/" + x.Name + " (" + base + ").md"
			}
			used[p] = true
			m.testPaths[x] = p
		}
	}
	return m.testPaths[t]
}

// renderSpecs met à jour les specs : tags dérivés du texte, blocs de tests des
// exigences, blocs de liens. Retourne le contenu par chemin.
func (m *model) renderSpecs() map[string]string {
	out := map[string]string{}
	reqsByDoc := map[*doc][]*req{}
	for _, group := range [][]*req{m.rms, m.enfs, m.efs, m.flows} {
		for _, r := range group {
			reqsByDoc[r.Doc] = append(reqsByDoc[r.Doc], r)
		}
	}
	for _, d := range m.docs {
		if d.Name == "Carte_des_specs" && d.Layer == "transverse" {
			continue
		}
		body := append([]string{}, d.body...)
		// Blocs de tests : insérés de la dernière section à la première pour
		// garder les indices valides.
		rs := reqsByDoc[d]
		sort.Slice(rs, func(i, j int) bool { return rs[i].Sec.start > rs[j].Sec.start })
		for _, r := range rs {
			body = insertBefore(body, r.Sec.end, m.testsBlock(r))
		}
		if lb := m.liensBlock(d); lb != nil {
			for len(body) > 0 && strings.TrimSpace(body[len(body)-1]) == "" {
				body = body[:len(body)-1]
			}
			body = append(append(body, ""), lb...)
		}
		out[d.Path] = joinDoc(m.updateTags(d), body, d.crlf)
	}
	out["specs/Carte_des_specs.md"] = m.renderCarte()
	return out
}

// insertBefore insère un bloc avant la ligne at, séparé par des lignes vides.
func insertBefore(body []string, at int, block []string) []string {
	head := append([]string{}, body[:at]...)
	for len(head) > 0 && strings.TrimSpace(head[len(head)-1]) == "" {
		head = head[:len(head)-1]
	}
	head = append(head, "")
	head = append(head, block...)
	if at < len(body) {
		head = append(head, "")
	}
	return append(head, body[at:]...)
}

func joinDoc(front, body []string, crlf bool) string {
	var lines []string
	if front != nil {
		lines = append(append([]string{"---"}, front...), "---")
	}
	lines = append(lines, body...)
	s := strings.Join(lines, "\n") + "\n"
	if crlf {
		s = strings.ReplaceAll(s, "\n", "\r\n")
	}
	return s
}

// testsBlock produit le bloc des tests d'une exigence.
func (m *model) testsBlock(r *req) []string {
	c := m.cov[r]
	from := r.Doc.Path
	lines := []string{testsBegin}
	all := c.all()
	if len(all) == 0 {
		lines = append(lines, "> [!warning] Tests — aucun test ne couvre cette exigence",
			"> "+link("Matrice de couverture", from, matrixPath, ""))
		return append(lines, testsEnd)
	}
	lines = append(lines, fmt.Sprintf("> [!check]- Tests — %d explicite(s) · %d déduit(s)", len(c.Explicit), len(c.Deduced)))
	for i, h := range all {
		if i == maxBlockTests {
			lines = append(lines, fmt.Sprintf("> - … et %d autre(s) : voir la %s", len(all)-maxBlockTests, link("matrice de couverture", from, matrixPath, "")))
			break
		}
		mark := "🟡"
		if h.Explicit {
			mark = "✅"
		}
		lines = append(lines, fmt.Sprintf("> - %s %s — %s", mark, link(h.T.Name, from, m.testNote(h.T), ""), h.Reason))
	}
	return append(lines, testsEnd)
}

// liensBlock produit la section Liens d'un use case ou d'un document de
// conception — nil pour les autres documents.
func (m *model) liensBlock(d *doc) []string {
	var lines []string
	switch {
	case d.UC != "" && d.Layer == "expression":
		lines = append(lines, "**Étape suivante — analyse**")
		if u := m.ucByID[d.UC]; u.Ana != nil {
			lines = append(lines, "- "+link(d.UC+" — analyse", d.Path, u.Ana.Path, ""))
		} else {
			lines = append(lines, "- _Aucune analyse ne traite ce use case : la chaîne s'arrête à l'expression._")
		}
	case d.UC != "" && d.Layer == "analyse":
		lines = append(lines, "**Étape suivante — conception**")
		conc := m.conceptionOf(d.UC)
		for _, cd := range conc {
			var secs []string
			for _, s := range cd.Sections {
				if s.Cites.UCs[d.UC] {
					secs = append(secs, link("§ "+s.Title, d.Path, cd.Path, s.Title))
				}
			}
			lines = append(lines, "- **"+cd.Name+"** : "+strings.Join(secs, " · "))
		}
		if len(conc) == 0 {
			lines = append(lines, "- _Aucun document de conception ne traite ce use case : la chaîne s'arrête à l'analyse._")
			if refs := m.codeRefs(m.resolve(d.Cites, nil)); len(refs) > 0 {
				lines = append(lines, "", "**Code cité par cette analyse (sans conception : non relié)**")
				for _, r := range refs {
					lines = append(lines, "- "+r)
				}
			}
		}
	case d.Layer == "conception" && d.Name != "todo" && !strings.HasSuffix(d.Name, "_intro"):
		lines = append(lines, "**Étape suivante — code, par section**")
		n := 0
		for _, s := range d.Sections {
			res := m.resolve(s.Cites, nil)
			links := m.codeLinks(d.Path, res)
			if len(links) == 0 {
				continue
			}
			n++
			lines = append(lines, "- "+link("§ "+s.Title, d.Path, d.Path, s.Title)+" → "+strings.Join(links, " · "))
		}
		if n == 0 {
			lines = append(lines, "- _Aucune route, commande ni fonction du code n'est citée par ce document._")
		}
	default:
		return nil
	}
	return append(append([]string{liensBegin, "## Liens", ""}, lines...), "", liensEnd)
}

// conceptionOf retourne les documents de conception dont une section cite le use case.
func (m *model) conceptionOf(uc string) []*doc {
	var out []*doc
	for _, d := range m.docs {
		if d.Layer != "conception" || d.Name == "todo" {
			continue
		}
		for _, s := range d.Sections {
			if s.Cites.UCs[uc] {
				out = append(out, d)
				break
			}
		}
	}
	return out
}

// codeLinks retourne les liens vers les notes du code résolu.
func (m *model) codeLinks(from string, res resolved) []string {
	var out []string
	seen := map[*routeNote]bool{}
	for _, rh := range res.Routes {
		if !seen[rh.Reg.Note] {
			seen[rh.Reg.Note] = true
			out = append(out, link("REST "+rh.Reg.Note.Title, from, routeNotePath(rh.Reg.Note), ""))
		}
	}
	for _, c := range res.Cmds {
		out = append(out, link("CLI "+c.Name(), from, cmdNote(c), ""))
	}
	for _, f := range res.Fns {
		out = append(out, link(f.Label(), from, fnNote(f), ""))
	}
	return out
}

// codeRefs retourne le code résolu sous forme de texte, sans lien.
func (m *model) codeRefs(res resolved) []string {
	var out []string
	seen := map[*routeNote]bool{}
	for _, rh := range res.Routes {
		if !seen[rh.Reg.Note] {
			seen[rh.Reg.Note] = true
			out = append(out, "REST "+code(rh.Reg.Note.Title))
		}
	}
	for _, c := range res.Cmds {
		out = append(out, "CLI "+code(c.Name()))
	}
	for _, f := range res.Fns {
		out = append(out, "Fonction "+code(f.Label()))
	}
	return out
}

// updateTags recalcule les tags dérivés du texte ; les autres tags (couche,
// type, famille, domaine…) sont saisis à la main et conservés :
//   - relecture/ : annotations de relecture ouvertes, pour tout document ;
//   - rm/ et enf/ : règles et exigences citées, pour un use case ou un document
//     de conception — ni les documents qui les définissent, ni les vues
//     transverses (introductions, roadmap) ne les portent ;
//   - uc/ : périmètre déclaré d'un document de conception.
func (m *model) updateTags(d *doc) []string {
	managed := map[string][]string{}
	prefixes := []string{"relecture/"}
	if d.UC != "" || d.Layer == "conception" {
		prefixes = []string{"rm/", "enf/", "relecture/"}
		for _, id := range keys(d.Cites.RMs) {
			if m.reqByID[id] != nil {
				managed["rm/"] = append(managed["rm/"], "rm/"+id)
			}
		}
		for _, id := range keys(d.Cites.ENFs) {
			if m.reqByID[id] != nil {
				managed["enf/"] = append(managed["enf/"], "enf/"+id)
			}
		}
	}
	for _, n := range keys(d.Cites.Notes) {
		managed["relecture/"] = append(managed["relecture/"], "relecture/"+n)
	}
	if d.Layer == "conception" {
		prefixes = append([]string{"uc/"}, prefixes...)
		for _, uc := range m.scopeOf(d) {
			managed["uc/"] = append(managed["uc/"], "uc/"+uc)
		}
	}
	isManaged := func(t string) bool {
		for _, p := range prefixes {
			if strings.HasPrefix(t, p) {
				return true
			}
		}
		return false
	}
	var tags []string
	for _, t := range d.tags() {
		if !isManaged(t) {
			tags = append(tags, t)
		}
	}
	for _, p := range prefixes {
		tags = append(tags, managed[p]...)
	}
	// Reconstruit le frontmatter : clés existantes conservées, liste tags remplacée.
	var front []string
	written := false
	inTags := false
	for _, l := range d.front {
		if l == "tags:" {
			inTags, written = true, true
			front = append(front, l)
			for _, t := range tags {
				front = append(front, "  - "+t)
			}
			continue
		}
		if inTags && strings.HasPrefix(l, "  - ") {
			continue
		}
		inTags = false
		front = append(front, l)
	}
	if !written && len(tags) > 0 {
		front = append(front, "tags:")
		for _, t := range tags {
			front = append(front, "  - "+t)
		}
	}
	if d.front == nil && len(tags) == 0 {
		return nil
	}
	return front
}

// renderCarte produit la carte des specs : use cases par famille, règles,
// exigences non fonctionnelles, documents de conception.
func (m *model) renderCarte() string {
	const from = "specs/Carte_des_specs.md"
	var b strings.Builder
	w := func(format string, a ...any) { fmt.Fprintf(&b, format+"\n", a...) }
	w("---\ntags:\n  - couche/transverse\n  - type/carte\n---")
	w("# Carte des specs\n")
	w("%s\n", generatedNote)
	w("Point d'entrée pour parcourir les specs de `myr-core` : chaque use case y est relié à ses deux couches (expression du besoin, analyse), aux documents de conception qui le couvrent et aux domaines Go concernés. Les use cases et les documents de conception se terminent par une section **Liens** (générée), et chaque fichier de `specs/` porte des tags dans son frontmatter.\n")
	w("## Comment naviguer\n")
	w("| Tag | Regroupe |\n|-----|----------|")
	w("| `couche/expression` · `couche/analyse` · `couche/conception` · `couche/code` | les documents d'une couche (les notes de code vivent dans `docs/code/`) |")
	w("| `famille/UCMOD` … | les use cases d'une famille |")
	w("| `uc/UCMOD04` … | tout ce qui concerne un use case : ses deux couches, la conception qui le déclare dans son périmètre, le code qui le réalise |")
	w("| `rm/RM13` · `enf/ENF12` … | les documents et les tests qui citent une règle métier ou une exigence non fonctionnelle |")
	w("| `domaine/model` … | specs et code d'un domaine Go (`domain/<domaine>/`) |")
	w("| `hexa/domaine` · `hexa/adapter-in` · `hexa/adapter-out` · `hexa/point-entree` | les packages Go d'une couche de l'architecture hexagonale |")
	w("| `relecture/question` · `relecture/remarque` · `relecture/incoherence` · `relecture/ecart` | les documents qui portent une annotation de relecture ouverte |\n")
	var views []string
	for _, name := range []string{"Expression_des_besoins_Intro", "Regles_Metier", "Exigences_Non_Fonctionnelles", "Matrice_Tracabilite", "Analyse_des_besoins", "Conception_intro", "Architecture_Hexagonale"} {
		if d := m.docByName(name); d != nil {
			views = append(views, link(name, from, d.Path, ""))
		}
	}
	views = append(views, link("Architecture du code", from, "docs/code/Architecture.md", ""), link("Traçabilité use cases ↔ code", from, "docs/code/Tracabilite_UC_Code.md", ""), link("Matrice de couverture des tests", from, matrixPath, ""))
	w("Vues transverses : %s\n", strings.Join(views, " · "))
	w("## Use cases\n")
	fam := ""
	for _, u := range m.ucs {
		if u.Family != fam {
			fam = u.Family
			w("### %s\n", u.FamilyLabel())
			w("Domaines : %s\n", m.familyDomains(fam))
			w("| UC | Titre | Expression | Analyse | Conception |\n|----|-------|------------|---------|------------|")
		}
		expr, ana := "—", "—"
		if u.Expr != nil {
			expr = link(u.Expr.Label, from, u.Expr.Path, "")
		}
		if u.Ana != nil {
			ana = link(u.Ana.Label, from, u.Ana.Path, "")
		}
		var conc []string
		for _, d := range m.conceptionOf(u.ID) {
			conc = append(conc, link(d.Name, from, d.Path, ""))
		}
		w("| %s | %s | %s | %s | %s |", u.ID, u.Title, expr, ana, dash(strings.Join(conc, " · ")))
		if next := m.nextUC(u); next == nil || next.Family != fam {
			w("")
		}
	}
	for _, g := range []struct {
		title, col string
		reqs       []*req
	}{{"Règles métier", "RM | Règle", m.rms}, {"Exigences non fonctionnelles", "ENF | Exigence", m.enfs}} {
		w("## %s\n", g.title)
		w("| %s | Use cases qui la citent |\n|----|-------|------------------------|", g.col)
		for _, r := range g.reqs {
			w("| %s | %s | %s |", link(r.ID, from, r.Doc.Path, r.Sec.Title), r.Title, dash(strings.Join(m.ucsCiting(from, r.ID), " · ")))
		}
		w("")
	}
	w("## Documents de conception\n")
	w("| Document | Périmètre déclaré | Domaines |\n|----------|-------------------|----------|")
	for _, d := range m.docs {
		if d.Layer != "conception" || d.Name == "todo" {
			continue
		}
		var doms []string
		for _, t := range d.tags() {
			if v, ok := strings.CutPrefix(t, "domaine/"); ok {
				doms = append(doms, code(v))
			}
		}
		w("| %s | %s | %s |", link(d.Name, from, d.Path, ""), dash(strings.Join(m.scopeOf(d), ", ")), dash(strings.Join(doms, ", ")))
	}
	w("\n## Autres documents\n")
	for _, d := range m.docs {
		if d.UC != "" || d.Path == from || (d.Layer == "conception" && d.Name != "todo") {
			continue
		}
		typ := ""
		for _, t := range d.tags() {
			if v, ok := strings.CutPrefix(t, "type/"); ok {
				typ = ", " + v
			}
		}
		w("- %s — %s%s", link(d.Label, from, d.Path, ""), d.Layer, typ)
	}
	return b.String()
}

func dash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

func (m *model) docByName(name string) *doc {
	for _, d := range m.docs {
		if d.Name == name {
			return d
		}
	}
	return nil
}

func (m *model) nextUC(u *useCase) *useCase {
	for i, x := range m.ucs {
		if x == u && i+1 < len(m.ucs) {
			return m.ucs[i+1]
		}
	}
	return nil
}

// familyDomains réunit les tags domaine/ des use cases d'une famille.
func (m *model) familyDomains(fam string) string {
	set := map[string]bool{}
	for _, u := range m.ucs {
		if u.Family != fam {
			continue
		}
		for _, d := range []*doc{u.Expr, u.Ana} {
			if d == nil {
				continue
			}
			for _, t := range d.tags() {
				if v, ok := strings.CutPrefix(t, "domaine/"); ok {
					set[v] = true
				}
			}
		}
	}
	var out []string
	for _, k := range keys(set) {
		out = append(out, code(k))
	}
	return dash(strings.Join(out, ", "))
}

// ucsCiting retourne les liens vers les use cases qui citent l'identifiant :
// l'analyse si elle le cite, sinon l'expression.
func (m *model) ucsCiting(from, id string) []string {
	var out []string
	cited := func(d *doc) bool {
		return d != nil && (d.Cites.RMs[id] || d.Cites.ENFs[id])
	}
	for _, u := range m.ucs {
		switch {
		case cited(u.Ana):
			out = append(out, link(u.ID, from, u.Ana.Path, ""))
		case cited(u.Expr):
			out = append(out, link(u.ID, from, u.Expr.Path, ""))
		}
	}
	return out
}

// scopeOf retourne le périmètre déclaré d'un document, réduit aux use cases existants.
func (m *model) scopeOf(d *doc) []string {
	var out []string
	for _, uc := range d.scope() {
		if m.ucByID[uc] != nil {
			out = append(out, uc)
		}
	}
	return out
}
