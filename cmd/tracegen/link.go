package main

import (
	"regexp"
	"sort"
	"strings"
)

// model relie les specs au code : résolution des citations, couverture des
// exigences par les tests, code sollicité par chaque use case.
type model struct {
	commit  string
	absRoot string // racine absolue, pour les liens d'ouverture dans VS Code
	docs    []*doc
	ucs     []*useCase
	ucByID  map[string]*useCase
	rms     []*req
	enfs    []*req
	efs     []*req
	flows   []*req
	reqByID map[string]*req
	cb      *codebase
	cov     map[*req]*coverage
	ucCode  map[string]*ucCode
	// citations des specs sans correspondance dans le code : citation → docs
	strayCmds   map[string]map[*doc]bool
	strayRoutes map[string]map[*doc]bool
	testPaths   map[*test]string
}

// hit est un test rattaché à une exigence, avec la raison du rattachement.
type hit struct {
	T        *test
	Reason   string
	Explicit bool
}

type coverage struct{ Explicit, Deduced []hit }

func (c *coverage) all() []hit { return append(append([]hit{}, c.Explicit...), c.Deduced...) }

func (c *coverage) status() string {
	switch {
	case len(c.Explicit) > 0:
		return "✅ explicite"
	case len(c.Deduced) > 0:
		return "🟡 déduite"
	}
	return "❌ non couverte"
}

// resolved est ce qu'un texte cite et qui existe dans le code.
type resolved struct {
	Routes []routeHit
	Cmds   []*command
	Fns    []*fn
}

type routeHit struct {
	Cite string // citation telle qu'écrite, ex. POST /api/identity/request
	Reg  *route
}

// ucCode est le code sollicité par un use case.
type ucCode struct {
	Routes []routeHit
	Cmds   []*command
	Fns    []*fn
	Why    map[*fn][]string
}

// genericNames sont des noms de méthode trop courants pour qu'une citation
// sans qualificatif désigne une méthode de port précise.
var genericNames = map[string]bool{
	"Get": true, "List": true, "Add": true, "Remove": true, "Delete": true, "Update": true, "Save": true,
	"Create": true, "Load": true, "FindAll": true, "FindByID": true, "FindByName": true, "Close": true,
	"Start": true, "Stop": true, "Run": true, "Execute": true, "String": true, "Error": true, "Clear": true,
	"Current": true, "Logout": true, "Activate": true, "Submit": true, "Register": true,
}

var svcSuffixRe = regexp.MustCompile(`(?i)(svc|service)$`)

// proseWords sont des mots courants qui, après « myr », signalent une phrase
// (« relancez myr avec… ») et non une commande.
var proseWords = map[string]bool{
	"avec": true, "qui": true, "est": true, "et": true, "ne": true, "pour": true, "sur": true, "dans": true,
	"le": true, "la": true, "les": true, "un": true, "une": true, "des": true, "du": true, "de": true,
	"en": true, "se": true, "sans": true, "peut": true, "doit": true, "a": true, "ou": true, "par": true,
}

func buildModel(docs []*doc, cb *codebase, commit string) *model {
	m := &model{commit: commit, docs: docs, cb: cb, ucByID: map[string]*useCase{}, reqByID: map[string]*req{},
		cov: map[*req]*coverage{}, ucCode: map[string]*ucCode{},
		strayCmds: map[string]map[*doc]bool{}, strayRoutes: map[string]map[*doc]bool{}}
	for _, d := range docs {
		if d.UC == "" {
			continue
		}
		u := m.ucByID[d.UC]
		if u == nil {
			u = &useCase{ID: d.UC, Family: d.Family, FamilyDir: d.FamilyDir}
			m.ucByID[d.UC] = u
			m.ucs = append(m.ucs, u)
		}
		if d.Layer == "expression" {
			u.Expr = d
			if d.Title != "" {
				u.Title = d.Title
			}
		} else {
			u.Ana = d
			if u.Title == "" {
				u.Title = d.Title
			}
		}
	}
	sort.Slice(m.ucs, func(i, j int) bool { return m.ucs[i].ID < m.ucs[j].ID })
	m.rms, m.enfs, m.efs, m.flows = collectReqs(docs)
	for _, r := range append(append(append([]*req{}, m.rms...), m.enfs...), m.efs...) {
		m.reqByID[r.ID] = r
	}
	for _, d := range docs {
		res := m.resolve(d.Cites, d)
		for _, rh := range res.Routes {
			addDoc(&rh.Reg.Note.Docs, d)
		}
		for _, c := range res.Cmds {
			addDoc(&c.Docs, d)
		}
		for _, f := range res.Fns {
			addDoc(&f.Docs, d)
		}
	}
	for _, r := range m.rms {
		m.cov[r] = m.cover(r)
	}
	for _, r := range m.enfs {
		m.cov[r] = m.cover(r)
	}
	for _, r := range m.flows {
		m.cov[r] = m.cover(r)
	}
	for _, r := range m.efs {
		m.cov[r] = m.coverEF(r)
	}
	for _, u := range m.ucs {
		m.ucCode[u.ID] = m.codeOf(u)
	}
	for _, f := range cb.fns {
		for _, rn := range cb.routeNotes {
			for _, c := range rn.calls() {
				if c == f {
					f.Routes = append(f.Routes, rn)
				}
			}
		}
		for _, c := range cb.cmds {
			for _, x := range c.Calls {
				if x == f {
					f.Cmds = append(f.Cmds, c)
				}
			}
		}
	}
	return m
}

func addDoc(list *[]*doc, d *doc) {
	for _, x := range *list {
		if x == d {
			return
		}
	}
	*list = append(*list, d)
}

// resolve rattache des citations au code. Les citations sans correspondance
// sont mémorisées pour le document d (s'il est fourni).
func (m *model) resolve(c cites, d *doc) resolved {
	var res resolved
	seenR := map[string]bool{}
	for _, rc := range c.Routes {
		cite := strings.TrimSpace(rc.Method + " " + rc.Path)
		if seenR[cite] {
			continue
		}
		seenR[cite] = true
		if r := m.cb.resolveRoute(rc.Path); r != nil {
			res.Routes = append(res.Routes, routeHit{Cite: cite, Reg: r})
		} else if d != nil && strings.HasPrefix(rc.Path, "/api/") && !strings.HasSuffix(rc.Path, "/") {
			// un préfixe (/api/admin/) désigne une famille de routes, pas une route
			stray(m.strayRoutes, rc.Path, d)
		}
	}
	sort.Slice(res.Routes, func(i, j int) bool { return res.Routes[i].Cite < res.Routes[j].Cite })
	seenC := map[*command]bool{}
	for _, cc := range c.Cmds {
		cmd := m.cb.resolveCmd(cc.Words)
		if cmd != nil && len(cmd.Words)-1 < len(cc.Words) && len(cmd.Children) > 0 {
			cmd = nil // mot suivant attendu parmi les sous-commandes : citation sans correspondance
		}
		if cmd == nil {
			if d != nil && cc.Quoted && !proseWords[cc.Words[0]] {
				stray(m.strayCmds, "myr "+strings.Join(cc.Words, " "), d)
			}
			continue
		}
		if !seenC[cmd] {
			seenC[cmd] = true
			res.Cmds = append(res.Cmds, cmd)
		}
	}
	sort.Slice(res.Cmds, func(i, j int) bool { return res.Cmds[i].Name() < res.Cmds[j].Name() })
	fns := map[*fn]bool{}
	for _, fc := range c.Funcs {
		for _, f := range m.resolveFunc(fc) {
			fns[f] = true
		}
	}
	res.Fns = sortFns(fns)
	return res
}

func stray(set map[string]map[*doc]bool, cite string, d *doc) {
	if set[cite] == nil {
		set[cite] = map[*doc]bool{}
	}
	set[cite][d] = true
}

// resolveFunc rattache une citation de méthode aux méthodes de port :
// Interface.Méthode, récepteur nommé d'après un domaine (identitySvc.Méthode)
// ou nom seul, s'il est unique parmi tous les ports et non générique.
func (m *model) resolveFunc(c funcCite) []*fn {
	var out []*fn
	if c.Qual != "" {
		for _, p := range m.cb.ports {
			if p.Name == c.Qual {
				for _, f := range p.Methods {
					if f.Name == c.Name {
						out = append(out, f)
					}
				}
			}
		}
		if len(out) > 0 {
			return out
		}
		dom := strings.ToLower(svcSuffixRe.ReplaceAllString(c.Qual, ""))
		if dom == "rbac" {
			dom = "role"
		}
		if dom == "" || dom == strings.ToLower(c.Qual) {
			return nil
		}
		for _, p := range m.cb.ports {
			if p.In && p.Pkg.Name == dom {
				for _, f := range p.Methods {
					if f.Name == c.Name {
						out = append(out, f)
					}
				}
			}
		}
		return out
	}
	if genericNames[c.Name] {
		return nil
	}
	for _, f := range m.cb.fns {
		if f.Name == c.Name {
			out = append(out, f)
		}
	}
	if len(out) != 1 {
		return nil
	}
	return out
}

// cover calcule les tests d'une règle, d'une exigence non fonctionnelle ou
// d'un flux : explicites s'ils citent l'identifiant, déduits s'ils exercent
// une route, une commande ou une fonction que la section cite.
func (m *model) cover(r *req) *coverage {
	c := &coverage{}
	explicit := map[*test]bool{}
	if r.ID != "" {
		for _, t := range m.cb.tests {
			if t.IDs[r.ID] {
				explicit[t] = true
				c.Explicit = append(c.Explicit, hit{T: t, Explicit: true, Reason: "explicite (le test cite " + r.ID + ")"})
			}
		}
	}
	reasons := map[*test][]string{}
	res := m.resolve(r.Sec.Cites, nil)
	seenNote := map[*routeNote]bool{}
	for _, rh := range res.Routes {
		if seenNote[rh.Reg.Note] {
			continue
		}
		seenNote[rh.Reg.Note] = true
		for _, t := range rh.Reg.Note.Tests {
			reasons[t] = append(reasons[t], "route "+code(rh.Reg.Note.Title))
		}
	}
	for _, cmd := range res.Cmds {
		for _, t := range cmd.Tests {
			reasons[t] = append(reasons[t], "commande "+code(cmd.Name()))
		}
	}
	for _, f := range res.Fns {
		for _, t := range f.Tests {
			reasons[t] = append(reasons[t], code(f.Label()))
		}
	}
	for t, rs := range reasons {
		if !explicit[t] {
			c.Deduced = append(c.Deduced, hit{T: t, Reason: "déduit : teste " + strings.Join(rs, ", ")})
		}
	}
	sortHits(c.Explicit)
	sortHits(c.Deduced)
	return c
}

// coverEF : une exigence fonctionnelle hérite des tests des flux de ses use cases.
func (m *model) coverEF(r *req) *coverage {
	c := &coverage{}
	seen := map[*test]bool{}
	for _, t := range m.cb.tests {
		if t.IDs[r.ID] {
			seen[t] = true
			c.Explicit = append(c.Explicit, hit{T: t, Explicit: true, Reason: "explicite (le test cite " + r.ID + ")"})
		}
	}
	for _, uc := range r.UCs {
		for _, fl := range m.flows {
			if fl.UC != uc {
				continue
			}
			for _, h := range m.cov[fl].all() {
				if !seen[h.T] {
					seen[h.T] = true
					c.Deduced = append(c.Deduced, hit{T: h.T, Reason: "déduit : via " + fl.Label()})
				}
			}
		}
	}
	sortHits(c.Explicit)
	sortHits(c.Deduced)
	return c
}

func sortHits(h []hit) {
	sort.Slice(h, func(i, j int) bool {
		if h[i].T.Subject.Path != h[j].T.Subject.Path {
			return h[i].T.Subject.Path < h[j].T.Subject.Path
		}
		return h[i].T.Name < h[j].T.Name
	})
}

// codeOf retourne le code sollicité par un use case : routes et commandes
// citées par ses specs, fonctions citées ou appelées par ces routes et commandes.
func (m *model) codeOf(u *useCase) *ucCode {
	uc := &ucCode{Why: map[*fn][]string{}}
	seenR := map[string]bool{}
	seenC := map[*command]bool{}
	why := func(f *fn, reason string) {
		for _, r := range uc.Why[f] {
			if r == reason {
				return
			}
		}
		uc.Why[f] = append(uc.Why[f], reason)
	}
	for _, d := range []*doc{u.Expr, u.Ana} {
		if d == nil {
			continue
		}
		res := m.resolve(d.Cites, nil)
		for _, rh := range res.Routes {
			if !seenR[rh.Cite] {
				seenR[rh.Cite] = true
				uc.Routes = append(uc.Routes, rh)
			}
			for _, f := range rh.Reg.Calls {
				why(f, "appelée par la route "+code(rh.Cite)+" citée dans la spec")
			}
		}
		for _, c := range res.Cmds {
			if !seenC[c] {
				seenC[c] = true
				uc.Cmds = append(uc.Cmds, c)
			}
			for _, f := range c.Calls {
				why(f, "appelée par la commande "+code(c.Name())+" citée dans la spec")
			}
		}
		for _, f := range res.Fns {
			why(f, "citée dans la spec ("+d.Layer+")")
		}
	}
	sort.Slice(uc.Routes, func(i, j int) bool { return uc.Routes[i].Cite < uc.Routes[j].Cite })
	sort.Slice(uc.Cmds, func(i, j int) bool { return uc.Cmds[i].Name() < uc.Cmds[j].Name() })
	set := map[*fn]bool{}
	for f, rs := range uc.Why {
		set[f] = true
		sort.SliceStable(rs, func(i, j int) bool { return reasonRank(rs[i]) < reasonRank(rs[j]) })
	}
	uc.Fns = sortFns(set)
	return uc
}

func reasonRank(r string) int {
	switch {
	case strings.HasPrefix(r, "citée"):
		return 0
	case strings.Contains(r, "commande"):
		return 1
	}
	return 2
}

// ucsOfFn retourne les use cases dont le code sollicité contient f.
func (m *model) ucsOfFn(f *fn) []string {
	var out []string
	for _, u := range m.ucs {
		if len(m.ucCode[u.ID].Why[f]) > 0 {
			out = append(out, u.ID)
		}
	}
	return out
}

// ucsOfRoute et ucsOfCmd : use cases dont une spec cite la route ou la commande.
func (m *model) ucsOfRoute(rn *routeNote) []string {
	var out []string
	for _, u := range m.ucs {
		for _, rh := range m.ucCode[u.ID].Routes {
			if rh.Reg.Note == rn {
				out = append(out, u.ID)
				break
			}
		}
	}
	return out
}

func (m *model) ucsOfCmd(c *command) []string {
	var out []string
	for _, u := range m.ucs {
		for _, x := range m.ucCode[u.ID].Cmds {
			if x == c {
				out = append(out, u.ID)
				break
			}
		}
	}
	return out
}

// reqsOfTest retourne les exigences qu'un test couvre, avec la raison.
type reqHit struct {
	R *req
	H hit
}

func (m *model) reqsOfTest(t *test) []reqHit {
	var out []reqHit
	for _, group := range [][]*req{m.rms, m.enfs, m.flows, m.efs} {
		for _, r := range group {
			for _, h := range m.cov[r].all() {
				if h.T == t {
					out = append(out, reqHit{R: r, H: h})
				}
			}
		}
	}
	return out
}
