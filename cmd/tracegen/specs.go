package main

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// doc est un fichier Markdown de specs/.
type doc struct {
	Path      string // relatif à la racine, séparateur /
	Layer     string // expression, analyse, conception, transverse
	Name      string // nom de fichier sans extension
	Label     string // libellé d'affichage, ex. « UCA01 (analyse) »
	UC        string // identifiant du use case décrit, vide sinon
	Family    string // famille du use case, ex. UCA
	FamilyDir string // dossier de la famille, ex. UCA-Compte_et_Acces
	Title     string // titre du use case (frontmatter titre) ou H1
	crlf      bool
	front     []string // lignes du frontmatter, sans les délimiteurs
	body      []string // corps sans les blocs générés
	Sections  []*section
	Cites     cites // citations de tout le document
}

// section est une portion de document introduite par un titre.
type section struct {
	Level      int
	Title      string
	start, end int // lignes de body : [start, end[, titre compris
	Cites      cites
}

// cites rassemble ce qu'un texte cite : identifiants d'exigences et de use
// cases, routes, commandes, fonctions, annotations de relecture.
type cites struct {
	UCs, RMs, ENFs, EFs, Notes map[string]bool
	Routes                     []routeCite
	Cmds                       []cmdCite
	Funcs                      []funcCite
}

type routeCite struct{ Method, Path string }

type cmdCite struct {
	Words  []string // mots qui suivent « myr »
	Quoted bool     // citée en code (backticks ou bloc de code)
}

// funcCite est une citation de méthode : Qual (interface ou récepteur, ex.
// IdentityService, identitySvc) peut être vide pour un nom seul.
type funcCite struct{ Qual, Name string }

var (
	ucRe        = regexp.MustCompile(`\bUC([A-Z]+)(\d{2})(?:\s*[–-]\s*(?:UC[A-Z]+)?(\d{2})\b)?((?:/\d{2}\b)*)`)
	reqRe       = regexp.MustCompile(`\b(RM|ENF|EF)(\d{2})\b`)
	routeRe     = regexp.MustCompile(`(?:\b(GET|POST|PUT|PATCH|DELETE)\s+)?(/api(?:/[A-Za-z0-9_\-{}:.]*)*/?|/metrics\b)`)
	cmdRe       = regexp.MustCompile(`\bmyr((?:\s+[a-z][a-z0-9-]*)+)`)
	qualFuncRe  = regexp.MustCompile(`\b([A-Za-z][A-Za-z0-9]*)\.([A-Z][A-Za-z0-9]*)\b`)
	bareFuncRe  = regexp.MustCompile("(?:`([A-Z][A-Za-z0-9]*)(?:\\([^`]*\\))?`|\\b([A-Z][A-Za-z0-9]*)\\()")
	noteRe      = regexp.MustCompile(`(?:^|[\s(])#(question|remarque|incoherence|ecart)\b`)
	reqTitleRe  = regexp.MustCompile(`^(RM|ENF|EF)(\d{2})\s+—\s+(.*)$`)
	coveringRe  = regexp.MustCompile(`Use cases couvrants\*{0,2}\s*:\s*(.*)$`)
	scopeLineRe = regexp.MustCompile(`Use cases\s*:\s*(.*)$`)
	ucFileRe    = regexp.MustCompile(`^UC[A-Z]+\d{2}$`)
)

func newCites() cites {
	return cites{UCs: map[string]bool{}, RMs: map[string]bool{}, ENFs: map[string]bool{}, EFs: map[string]bool{}, Notes: map[string]bool{}}
}

// merge ajoute à c les citations de o.
func (c *cites) merge(o cites) {
	for _, p := range []struct{ dst, src map[string]bool }{{c.UCs, o.UCs}, {c.RMs, o.RMs}, {c.ENFs, o.ENFs}, {c.EFs, o.EFs}, {c.Notes, o.Notes}} {
		for k := range p.src {
			p.dst[k] = true
		}
	}
	c.Routes = append(c.Routes, o.Routes...)
	c.Cmds = append(c.Cmds, o.Cmds...)
	c.Funcs = append(c.Funcs, o.Funcs...)
}

// expandUCs retourne les identifiants de use case cités dans s, plages et
// suffixes compris : « UCA01–03 » cite UCA01, UCA02 et UCA03 ;
// « UCAM01–03/05/07 » y ajoute UCAM05 et UCAM07. Une plage peut couvrir des
// numéros sans use case : l'appelant ne retient que les use cases existants.
func expandUCs(s string) []string {
	var out []string
	for _, m := range ucRe.FindAllStringSubmatch(s, -1) {
		out = append(out, "UC"+m[1]+m[2])
		if m[3] != "" {
			from, to := atoi(m[2]), atoi(m[3])
			for n := from + 1; n <= to && n-from < 20; n++ {
				out = append(out, "UC"+m[1]+pad2(n))
			}
		}
		for _, suf := range strings.Split(m[4], "/") {
			if suf != "" {
				out = append(out, "UC"+m[1]+suf)
			}
		}
	}
	return out
}

func atoi(s string) int {
	n := 0
	for _, r := range s {
		n = n*10 + int(r-'0')
	}
	return n
}

func pad2(n int) string {
	if n < 10 {
		return "0" + string(rune('0'+n))
	}
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}

// extract relève les citations d'un ensemble de lignes.
func extract(lines []string) cites {
	c := newCites()
	inFence := false
	for _, l := range lines {
		if isFence(l) {
			inFence = !inFence
			continue
		}
		for _, id := range expandUCs(l) {
			c.UCs[id] = true
		}
		for _, m := range reqRe.FindAllStringSubmatch(l, -1) {
			id := m[1] + m[2]
			switch m[1] {
			case "RM":
				c.RMs[id] = true
			case "ENF":
				c.ENFs[id] = true
			case "EF":
				c.EFs[id] = true
			}
		}
		for _, m := range routeRe.FindAllStringSubmatch(l, -1) {
			p := strings.TrimRight(m[2], ".:")
			c.Routes = append(c.Routes, routeCite{Method: m[1], Path: p})
		}
		quoted := quotedSpans(l)
		for _, m := range cmdRe.FindAllStringSubmatchIndex(l, -1) {
			end, inSpan := m[3], false
			for _, s := range quoted {
				if m[0] > s[0] && m[0] < s[1] {
					// Citation en code : seuls les mots du même span forment la commande.
					end, inSpan = min(end, s[1]), true
				}
			}
			words := strings.Fields(l[m[2]:max(end, m[2])])
			if len(words) > 0 {
				c.Cmds = append(c.Cmds, cmdCite{Words: words, Quoted: inFence || inSpan})
			}
		}
		for _, m := range qualFuncRe.FindAllStringSubmatch(l, -1) {
			c.Funcs = append(c.Funcs, funcCite{Qual: m[1], Name: m[2]})
		}
		for _, m := range bareFuncRe.FindAllStringSubmatch(l, -1) {
			name := m[1]
			if name == "" {
				name = m[2]
			}
			c.Funcs = append(c.Funcs, funcCite{Name: name})
		}
		if !inFence {
			for _, m := range noteRe.FindAllStringSubmatch(stripInlineCode(l), -1) {
				c.Notes[m[1]] = true
			}
		}
	}
	return c
}

// quotedSpans retourne les intervalles [début, fin[ des spans de code d'une ligne.
func quotedSpans(l string) [][2]int {
	var spans [][2]int
	start := -1
	for i, r := range l {
		if r != '`' {
			continue
		}
		if start < 0 {
			start = i
		} else {
			spans = append(spans, [2]int{start, i})
			start = -1
		}
	}
	return spans
}

// loadSpecs lit tous les fichiers Markdown de specs/.
func loadSpecs(root string) ([]*doc, error) {
	var docs []*doc
	err := filepath.WalkDir(filepath.Join(root, "specs"), func(p string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".md") {
			return err
		}
		rel, _ := filepath.Rel(root, p)
		dc, err := readDoc(root, filepath.ToSlash(rel))
		if err != nil {
			return err
		}
		docs = append(docs, dc)
		return nil
	})
	sort.Slice(docs, func(i, j int) bool { return docs[i].Path < docs[j].Path })
	return docs, err
}

func readDoc(root, rel string) (*doc, error) {
	raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		return nil, err
	}
	text := string(raw)
	d := &doc{Path: rel, crlf: strings.Contains(text, "\r\n")}
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) > 0 && lines[0] == "---" {
		for i := 1; i < len(lines); i++ {
			if lines[i] == "---" {
				d.front = lines[1:i]
				lines = lines[i+1:]
				break
			}
		}
	}
	d.body = removeBlocks(lines)

	parts := strings.Split(rel, "/")
	d.Name = strings.TrimSuffix(parts[len(parts)-1], ".md")
	switch {
	case strings.HasPrefix(rel, "specs/1-Expression/"):
		d.Layer = "expression"
	case strings.HasPrefix(rel, "specs/2-Analyse/"):
		d.Layer = "analyse"
	case strings.HasPrefix(rel, "specs/3-Conception/"):
		d.Layer = "conception"
	default:
		d.Layer = "transverse"
	}
	d.Label = d.Name
	if ucFileRe.MatchString(d.Name) && len(parts) == 4 {
		d.UC = d.Name
		d.FamilyDir = parts[2]
		d.Family = strings.SplitN(d.FamilyDir, "-", 2)[0]
		d.Label = d.UC + " (" + d.Layer + ")"
	} else if d.Name == "todo" {
		d.Label = "todo (" + d.Layer + ")"
	}
	d.Title = frontValue(d.front, "titre")
	d.parseSections()
	return d, nil
}

// frontValue retourne la valeur d'une clé simple du frontmatter, sans guillemets.
func frontValue(front []string, key string) string {
	for _, l := range front {
		if v, ok := strings.CutPrefix(l, key+":"); ok {
			return strings.Trim(strings.TrimSpace(v), `"'`)
		}
	}
	return ""
}

// tags retourne les tags du frontmatter (liste YAML « tags: » à tirets).
func (d *doc) tags() []string {
	var out []string
	in := false
	for _, l := range d.front {
		switch {
		case l == "tags:":
			in = true
		case in && strings.HasPrefix(l, "  - "):
			out = append(out, strings.TrimPrefix(l, "  - "))
		case in:
			in = false
		}
	}
	return out
}

func (d *doc) parseSections() {
	inFence := false
	var cur *section
	for i, l := range d.body {
		if isFence(l) {
			inFence = !inFence
		}
		if inFence {
			continue
		}
		m := headingRe.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		if cur != nil {
			cur.end = i
		}
		cur = &section{Level: len(m[1]), Title: m[2], start: i}
		d.Sections = append(d.Sections, cur)
		if cur.Level == 1 && d.Title == "" {
			d.Title = cur.Title
		}
	}
	if cur != nil {
		cur.end = len(d.body)
	}
	d.Cites = extract(d.body)
	for _, s := range d.Sections {
		s.Cites = extract(d.body[s.start:s.end])
	}
}

// scope retourne les use cases que le document déclare dans son périmètre
// (ligne « Use cases : … » de son en-tête).
func (d *doc) scope() []string {
	for _, l := range d.body {
		if strings.HasPrefix(strings.TrimSpace(l), "#") && !strings.HasPrefix(strings.TrimSpace(l), "# ") {
			break
		}
		if m := scopeLineRe.FindStringSubmatch(l); m != nil && strings.HasPrefix(strings.TrimSpace(l), ">") {
			return uniqueSorted(expandUCs(m[1]))
		}
	}
	return nil
}

// useCase réunit les deux couches d'un use case.
type useCase struct {
	ID, Title, Family, FamilyDir string
	Expr, Ana                    *doc
}

func (u *useCase) FamilyLabel() string {
	return strings.ReplaceAll(strings.Replace(u.FamilyDir, "-", " — ", 1), "_", " ")
}

// req est une exigence à couvrir par des tests : règle métier, exigence non
// fonctionnelle, exigence fonctionnelle ou flux d'un use case (analyse).
type req struct {
	ID    string // RM03, ENF12, EF01 — vide pour un flux
	Kind  string // RM, ENF, EF, flux
	Title string // titre de la section
	UC    string // use case du flux
	Doc   *doc
	Sec   *section
	UCs   []string // use cases couvrants (EF)
}

// Label est le libellé d'une exigence dans les tableaux et les listes.
func (r *req) Label() string {
	if r.Kind == "flux" {
		return r.UC + " · " + r.Sec.Title
	}
	return r.Sec.Title
}

// collectReqs relève les exigences des specs.
func collectReqs(docs []*doc) (rms, enfs, efs, flows []*req) {
	for _, d := range docs {
		for _, s := range d.Sections {
			if m := reqTitleRe.FindStringSubmatch(s.Title); m != nil && d.Layer == "expression" && d.UC == "" {
				r := &req{ID: m[1] + m[2], Kind: m[1], Title: m[3], Doc: d, Sec: s}
				switch m[1] {
				case "RM":
					rms = append(rms, r)
				case "ENF":
					enfs = append(enfs, r)
				case "EF":
					for _, l := range d.body[s.start:s.end] {
						if c := coveringRe.FindStringSubmatch(l); c != nil {
							r.UCs = uniqueSorted(expandUCs(c[1]))
						}
					}
					efs = append(efs, r)
				}
				continue
			}
			if d.Layer == "analyse" && d.UC != "" && strings.HasPrefix(s.Title, "Flux ") {
				flows = append(flows, &req{Kind: "flux", Title: s.Title, UC: d.UC, Doc: d, Sec: s})
			}
		}
	}
	return
}

func uniqueSorted(in []string) []string {
	set := map[string]bool{}
	for _, s := range in {
		set[s] = true
	}
	return keys(set)
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
