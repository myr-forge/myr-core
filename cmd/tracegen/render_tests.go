package main

import (
	"fmt"
	"sort"
	"strings"
)

// renderTests produit une note par test et la matrice de couverture.
func (m *model) renderTests() map[string]string {
	out := map[string]string{}
	for _, t := range m.cb.tests {
		out[m.testNote(t)] = m.testPage(t)
	}
	out[matrixPath] = m.matrixPage()
	return out
}

// reqLink retourne le lien vers la section d'une exigence.
func (w *noteWriter) reqLink(r *req) string {
	return w.link(r.Label(), r.Doc.Path, r.Sec.Title)
}

func (m *model) testPage(t *test) string {
	reqs := m.reqsOfTest(t)
	tags := []string{"couche/test", "type/test", "hexa/" + t.Subject.Layer}
	for _, id := range sortedIDs(t.IDs) {
		switch {
		case strings.HasPrefix(id, "RM"):
			tags = append(tags, "rm/"+id)
		case strings.HasPrefix(id, "ENF"):
			tags = append(tags, "enf/"+id)
		case m.ucByID[id] != nil:
			tags = append(tags, "uc/"+id)
		}
	}
	if len(reqs) == 0 {
		tags = append(tags, "test/sans-exigence")
	}
	w := m.note(m.testNote(t), t.Name, tags)
	w.f("**Package testé** : %s · **Source** : %s\n", w.link(t.Subject.Path, pkgNote(t.Subject), ""), w.loc(t.File, t.Line))
	if t.Brief != "" {
		w.line("> " + strings.TrimSpace(t.Brief) + "\n")
	}
	if t.Input != "" {
		w.line("- **Entrée** : " + strings.TrimSpace(t.Input))
	}
	if t.Expect != "" {
		w.line("- **Attendu** : " + strings.TrimSpace(t.Expect))
	}
	if t.Input != "" || t.Expect != "" {
		w.line("")
	}
	w.line("## Exigences couvertes\n")
	if len(reqs) == 0 {
		w.line("_Aucune exigence rattachée : le test ne cite aucun identifiant et ne teste aucun élément cité par une exigence._\n")
	}
	for _, rh := range reqs {
		mark, reason := "🟡", strings.Replace(rh.H.Reason, "déduit", "déduite", 1)
		if rh.H.Explicit {
			mark = "✅"
		}
		w.f("- %s %s — %s", mark, w.reqLink(rh.R), reason)
	}
	if len(reqs) > 0 {
		w.line("")
	}
	var codeLinks []string
	for _, rn := range t.Routes {
		codeLinks = append(codeLinks, w.link("REST "+rn.Title, routeNotePath(rn), ""))
	}
	for _, c := range t.Cmds {
		codeLinks = append(codeLinks, w.link("CLI "+c.Name(), cmdNote(c), ""))
	}
	for _, f := range t.Fns {
		codeLinks = append(codeLinks, w.link(f.Label(), fnNote(f), ""))
	}
	w.line("## Code testé\n")
	if len(codeLinks) == 0 {
		w.line("_Aucune route, commande ni fonction de port détectée dans le test._")
	}
	for _, l := range codeLinks {
		w.line("- " + l)
	}
	return w.b.String()
}

func (m *model) matrixPage() string {
	w := m.note(matrixPath, "Matrice de couverture des tests", []string{"couche/test", "type/tracabilite"})
	w.line("Pour chaque exigence des specs (règle métier, exigence non fonctionnelle, exigence fonctionnelle, flux d'un use case d'analyse) : les tests qui la couvrent.\n")
	w.line("- ✅ **explicite** : le nom ou le commentaire du test cite l'identifiant de l'exigence — couverture voulue et vérifiable ;")
	w.line("- 🟡 **déduite** : le test exerce une fonction de port, une route ou une commande que l'exigence cite — couverture probable, à confirmer en relisant le test ;")
	w.line("- ❌ **non couverte** : aucun test rattaché.\n")
	w.line("Une exigence fonctionnelle (EF) hérite des tests des flux de ses use cases. Chaque exigence ouvre sa section dans les specs ; chaque test ouvre sa note (code source, exigences couvertes).\n")
	groups := []struct {
		title string
		reqs  []*req
	}{{"Règles métier", m.rms}, {"Exigences non fonctionnelles", m.enfs}, {"Exigences fonctionnelles", m.efs}, {"Flux des use cases (analyse)", m.flows}}
	w.line("## Synthèse\n")
	w.line("| Catégorie | Exigences | ✅ Explicite | 🟡 Déduite seulement | ❌ Non couverte | Couverture explicite |\n|-----------|-----------|-------------|---------------------|----------------|---------------------|")
	for _, g := range groups {
		ex, de, no := 0, 0, 0
		for _, r := range g.reqs {
			c := m.cov[r]
			switch {
			case len(c.Explicit) > 0:
				ex++
			case len(c.Deduced) > 0:
				de++
			default:
				no++
			}
		}
		pct := 0
		if len(g.reqs) > 0 {
			pct = ex * 100 / len(g.reqs)
		}
		w.f("| %s | %d | %d | %d | %d | %d %% |", w.link(g.title, matrixPath, g.title), len(g.reqs), ex, de, no, pct)
	}
	orphans := map[string][]*test{}
	nOrphans := 0
	for _, t := range m.cb.tests {
		if len(m.reqsOfTest(t)) == 0 {
			orphans[t.Subject.Path] = append(orphans[t.Subject.Path], t)
			nOrphans++
		}
	}
	w.f("\n%d tests au total, dont %d rattachés à aucune exigence (%s).\n", len(m.cb.tests), nOrphans, w.link("détail", matrixPath, "Tests sans exigence"))
	for _, g := range groups {
		w.line("## " + g.title + "\n")
		w.line("| Exigence | Statut | Tests explicites | Tests déduits |\n|----------|--------|------------------|---------------|")
		for _, r := range g.reqs {
			c := m.cov[r]
			var ex []string
			for _, h := range c.Explicit {
				ex = append(ex, w.link(h.T.Name, m.testNote(h.T), ""))
			}
			de := "—"
			if n := len(c.Deduced); n > 0 {
				var ls []string
				for i, h := range c.Deduced {
					if i == 3 {
						break
					}
					ls = append(ls, w.link(h.T.Name, m.testNote(h.T), ""))
				}
				de = fmt.Sprintf("%d : %s", n, strings.Join(ls, " · "))
				if n > 3 {
					de += " …"
				}
			}
			w.f("| %s | %s | %s | %s |", w.reqLink(r), c.status(), dash(strings.Join(ex, " · ")), de)
		}
		w.line("")
	}
	w.line("## Tests sans exigence\n")
	w.line("Tests qui ne citent aucun identifiant et n'exercent aucun élément cité par une exigence, par package.\n")
	w.line("| Package | Tests sans exigence | Tests au total |\n|---------|--------------------|----------------|")
	for _, p := range m.cb.pkgs {
		if len(p.Tests) > 0 {
			w.f("| %s | %d | %d |", w.link(p.Path, pkgNote(p), ""), len(orphans[p.Path]), len(p.Tests))
		}
	}
	w.line("")
	w.line("## Identifiants cités par des tests mais absents des specs\n")
	missing := map[string][]*test{}
	for _, t := range m.cb.tests {
		for id := range t.IDs {
			if strings.HasPrefix(id, "UC") {
				if m.ucByID[id] == nil {
					missing[id] = append(missing[id], t)
				}
			} else if m.reqByID[id] == nil {
				missing[id] = append(missing[id], t)
			}
		}
	}
	if len(missing) == 0 {
		w.line("_Aucun._")
		return w.b.String()
	}
	w.line("| Identifiant | Tests |\n|-------------|-------|")
	for _, id := range sortedKeys(missing) {
		ts := missing[id]
		sort.Slice(ts, func(i, j int) bool { return ts[i].Name < ts[j].Name })
		var ls []string
		for _, t := range ts {
			ls = append(ls, w.link(t.Name, m.testNote(t), ""))
		}
		w.f("| %s | %s |", code(id), strings.Join(ls, " · "))
	}
	return w.b.String()
}
