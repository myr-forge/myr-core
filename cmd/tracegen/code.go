package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

const modulePath = "myr-core"

// pkg est un package Go du module (hors tests, vendor et chaincode).
type pkg struct {
	Path       string // ex. domain/identity
	Name       string // nom Go du package
	Layer      string // domaine, adapter-in, adapter-out, point-entree
	Files      []fileStat
	Imports    map[string]bool
	ImportedBy map[string]bool
	Ports      []*port
	Tests      []*test

	fset    *token.FileSet
	files   []*fileCtx
	funcs   map[string]*funcCtx            // fonctions de package
	methods map[string]map[string]*funcCtx // type → méthode
	fields  map[string]map[string]typeRef  // type struct → champ → type
	vars    map[string]typeRef             // variables de package
}

type fileStat struct {
	Path  string
	Lines int
}

type fileCtx struct {
	path    string
	ast     *ast.File
	imports map[string]string // alias → chemin du package importé
}

type funcCtx struct {
	pkg  *pkg
	file *fileCtx
	decl *ast.FuncDecl
}

// typeRef est une expression de type et le fichier qui permet de la résoudre.
type typeRef struct {
	expr ast.Expr
	file *fileCtx
}

// port est une interface déclarée par un domaine : d'entrée (port_in.go) ou
// de sortie (ports.go).
type port struct {
	Name    string
	In      bool
	Pkg     *pkg
	File    string
	Line    int
	Methods []*fn
	Impls   []impl
}

type impl struct {
	Pkg  *pkg
	Type string
}

// fn est une méthode de port : l'unité de traçabilité du code.
type fn struct {
	Port              *port
	Name, Sig, Doc    string
	DeclFile          string
	DeclLine          int
	ImplFile          string
	ImplLine          int
	Routes            []*routeNote
	Cmds              []*command
	Tests             []*test
	Docs              []*doc
	nParams, nResults int
}

// Key identifie la fonction : domaine.Interface.Méthode.
func (f *fn) Key() string { return f.Port.Pkg.Name + "." + f.Port.Name + "." + f.Name }

// Label est le nom d'affichage : Interface.Méthode.
func (f *fn) Label() string { return f.Port.Name + "." + f.Name }

// route est un enregistrement de route dans le routeur REST.
type route struct {
	Pattern string
	File    string
	Line    int
	Handler string
	HFile   string
	HLine   int
	Calls   []*fn
	Note    *routeNote
}

// routeNote regroupe les enregistrements qui partagent une note (/api/x et /api/x/).
type routeNote struct {
	Slug, Title string
	Regs        []*route
	Tests       []*test
	Docs        []*doc
}

func (r *routeNote) calls() []*fn {
	set := map[*fn]bool{}
	for _, reg := range r.Regs {
		for _, f := range reg.Calls {
			set[f] = true
		}
	}
	return sortFns(set)
}

// command est une commande cobra du CLI.
type command struct {
	Var      string
	Words    []string // ex. [myr identity request]
	File     string
	Line     int
	Parent   *command
	Children []*command
	Calls    []*fn
	runFuncs map[string]bool // fonctions du package appelées directement par RunE
	Tests    []*test
	Docs     []*doc
	body     *ast.FuncLit
	fc       *fileCtx
	pkg      *pkg
}

func (c *command) Name() string { return strings.Join(c.Words, " ") }

// test est une fonction de test Go.
type test struct {
	Name, File           string
	Line                 int
	Subject              *pkg
	Brief, Input, Expect string
	IDs                  map[string]bool // exigences et use cases cités par le nom ou le commentaire
	calls                map[string]bool // méthodes appelées
	idents               map[string]bool // fonctions appelées
	literals             []string        // chemins de route en littéral
	Fns                  []*fn
	Routes               []*routeNote
	Cmds                 []*command
}

// codebase est le modèle du code analysé.
type codebase struct {
	pkgs       []*pkg
	byPath     map[string]*pkg
	ports      []*port
	fns        []*fn
	routes     []*route
	routeNotes []*routeNote
	cmds       []*command
	tests      []*test
}

var testIDRe = regexp.MustCompile(`\b(RM|ENF|EF)\d{2}\b`)

// loadCode analyse les packages Go de domain/, adapters/ et cmd/.
func loadCode(root string) (*codebase, error) {
	cb := &codebase{byPath: map[string]*pkg{}}
	// Un seul FileSet pour tout le code : les positions restent valables quel
	// que soit le fichier dont elles proviennent.
	fset := token.NewFileSet()
	type testFile struct {
		rel  string
		dir  string
		file *ast.File
		fset *token.FileSet
	}
	var testFiles []testFile
	for _, top := range []string{"domain", "adapters", "cmd"} {
		err := filepath.WalkDir(filepath.Join(root, top), func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if name := d.Name(); name == "testdata" || name == "vendor" {
					return filepath.SkipDir
				}
				return nil
			}
			if !strings.HasSuffix(p, ".go") {
				return nil
			}
			rel, _ := filepath.Rel(root, p)
			rel = filepath.ToSlash(rel)
			dir := filepath.ToSlash(filepath.Dir(rel))
			f, err := parser.ParseFile(fset, p, nil, parser.ParseComments)
			if err != nil {
				return err
			}
			if strings.HasSuffix(p, "_test.go") {
				testFiles = append(testFiles, testFile{rel: rel, dir: dir, file: f, fset: fset})
				return nil
			}
			pk := cb.byPath[dir]
			if pk == nil {
				pk = &pkg{Path: dir, Name: f.Name.Name, Layer: layerOf(dir), Imports: map[string]bool{}, ImportedBy: map[string]bool{},
					funcs: map[string]*funcCtx{}, methods: map[string]map[string]*funcCtx{}, fields: map[string]map[string]typeRef{}, vars: map[string]typeRef{}}
				cb.byPath[dir] = pk
				cb.pkgs = append(cb.pkgs, pk)
			}
			raw, _ := os.ReadFile(p)
			pk.Files = append(pk.Files, fileStat{Path: rel, Lines: bytes.Count(raw, []byte("\n"))})
			pk.addFile(rel, f, fset)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	sort.Slice(cb.pkgs, func(i, j int) bool { return cb.pkgs[i].Path < cb.pkgs[j].Path })
	for _, pk := range cb.pkgs {
		sort.Slice(pk.Files, func(i, j int) bool { return pk.Files[i].Path < pk.Files[j].Path })
		for imp := range pk.Imports {
			if dep := cb.byPath[imp]; dep != nil {
				dep.ImportedBy[pk.Path] = true
			} else {
				delete(pk.Imports, imp)
			}
		}
	}
	cb.loadPorts()
	cb.loadRoutes()
	cb.loadCommands()
	for _, tf := range testFiles {
		cb.loadTests(tf.rel, tf.dir, tf.file, tf.fset)
	}
	cb.linkTests()
	return cb, nil
}

func layerOf(dir string) string {
	switch {
	case strings.HasPrefix(dir, "domain/"):
		return "domaine"
	case strings.HasPrefix(dir, "adapters/in/"):
		return "adapter-in"
	case strings.HasPrefix(dir, "adapters/out/"):
		return "adapter-out"
	default:
		return "point-entree"
	}
}

func (pk *pkg) addFile(rel string, f *ast.File, fset *token.FileSet) {
	pk.fset = fset
	fc := &fileCtx{path: rel, ast: f, imports: map[string]string{}}
	for _, im := range f.Imports {
		p, _ := strconv.Unquote(im.Path.Value)
		alias := p[strings.LastIndex(p, "/")+1:]
		if im.Name != nil {
			alias = im.Name.Name
		}
		fc.imports[alias] = p
		if strings.HasPrefix(p, modulePath+"/") {
			pk.Imports[strings.TrimPrefix(p, modulePath+"/")] = true
		}
	}
	pk.files = append(pk.files, fc)
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			ctx := &funcCtx{pkg: pk, file: fc, decl: d}
			if d.Recv == nil {
				pk.funcs[d.Name.Name] = ctx
				continue
			}
			t := recvType(d.Recv.List[0].Type)
			if pk.methods[t] == nil {
				pk.methods[t] = map[string]*funcCtx{}
			}
			pk.methods[t][d.Name.Name] = ctx
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					if st, ok := s.Type.(*ast.StructType); ok {
						m := map[string]typeRef{}
						for _, fl := range st.Fields.List {
							for _, n := range fl.Names {
								m[n.Name] = typeRef{fl.Type, fc}
							}
						}
						pk.fields[s.Name.Name] = m
					}
				case *ast.ValueSpec:
					for _, n := range s.Names {
						if s.Type != nil {
							pk.vars[n.Name] = typeRef{s.Type, fc}
						}
					}
				}
			}
		}
	}
}

func recvType(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.StarExpr:
		return recvType(t.X)
	case *ast.Ident:
		return t.Name
	case *ast.IndexExpr:
		return recvType(t.X)
	}
	return ""
}

func (cb *codebase) pos(pk *pkg, p token.Pos) int { return pk.fset.Position(p).Line }

// loadPorts relève les interfaces des domaines et leurs implémentations.
func (cb *codebase) loadPorts() {
	for _, pk := range cb.pkgs {
		if pk.Layer != "domaine" {
			continue
		}
		for _, fc := range pk.files {
			base := filepath.Base(fc.path)
			if base != "port_in.go" && base != "ports.go" {
				continue
			}
			for _, decl := range fc.ast.Decls {
				gd, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}
				for _, spec := range gd.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					it, ok := ts.Type.(*ast.InterfaceType)
					if !ok {
						continue
					}
					p := &port{Name: ts.Name.Name, In: base == "port_in.go", Pkg: pk, File: fc.path, Line: cb.pos(pk, ts.Pos())}
					for _, m := range it.Methods.List {
						ft, ok := m.Type.(*ast.FuncType)
						if !ok || len(m.Names) == 0 {
							continue
						}
						f := &fn{Port: p, Name: m.Names[0].Name, DeclFile: fc.path, DeclLine: cb.pos(pk, m.Pos()),
							Sig: m.Names[0].Name + strings.TrimPrefix(printNode(pk.fset, ft), "func"),
							Doc: oneLine(m.Doc.Text()), nParams: countFields(ft.Params), nResults: countFields(ft.Results)}
						p.Methods = append(p.Methods, f)
						cb.fns = append(cb.fns, f)
					}
					pk.Ports = append(pk.Ports, p)
					cb.ports = append(cb.ports, p)
				}
			}
		}
	}
	// Implémentations : types dont les méthodes couvrent toutes celles du port
	// (même nom, même nombre de paramètres et de résultats).
	for _, p := range cb.ports {
		for _, pk := range cb.pkgs {
			for _, t := range sortedKeys(pk.methods) {
				ms := pk.methods[t]
				ok := len(p.Methods) > 0
				for _, f := range p.Methods {
					m := ms[f.Name]
					if m == nil || countFields(m.decl.Type.Params) != f.nParams || countFields(m.decl.Type.Results) != f.nResults {
						ok = false
						break
					}
				}
				if !ok {
					continue
				}
				p.Impls = append(p.Impls, impl{Pkg: pk, Type: t})
				for _, f := range p.Methods {
					if f.ImplFile != "" && (!p.In || pk != p.Pkg) {
						continue
					}
					m := ms[f.Name]
					f.ImplFile, f.ImplLine = m.file.path, cb.pos(pk, m.decl.Pos())
				}
			}
		}
	}
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func countFields(fl *ast.FieldList) int {
	if fl == nil {
		return 0
	}
	n := 0
	for _, f := range fl.List {
		if len(f.Names) == 0 {
			n++
		} else {
			n += len(f.Names)
		}
	}
	return n
}

func printNode(fset *token.FileSet, n ast.Node) string {
	var b bytes.Buffer
	_ = printer.Fprint(&b, fset, n)
	return b.String()
}

func oneLine(s string) string { return strings.Join(strings.Fields(s), " ") }

// portOf résout une expression de type vers un port du module.
func (cb *codebase) portOf(t typeRef) *port {
	e := t.expr
	if s, ok := e.(*ast.StarExpr); ok {
		e = s.X
	}
	var path, name string
	switch x := e.(type) {
	case *ast.SelectorExpr:
		id, ok := x.X.(*ast.Ident)
		if !ok {
			return nil
		}
		path, name = strings.TrimPrefix(t.file.imports[id.Name], modulePath+"/"), x.Sel.Name
	case *ast.Ident:
		path, name = filepath.ToSlash(filepath.Dir(t.file.path)), x.Name
	default:
		return nil
	}
	for _, p := range cb.ports {
		if p.Pkg.Path == path && p.Name == name {
			return p
		}
	}
	return nil
}

// uniqueIn retourne la méthode de port d'entrée nommée name si elle est la
// seule de ce nom parmi tous les ports d'entrée.
func (cb *codebase) uniqueIn(name string) *fn {
	var found *fn
	for _, f := range cb.fns {
		if f.Port.In && f.Name == name {
			if found != nil {
				return nil
			}
			found = f
		}
	}
	return found
}

// calls relève les méthodes de port appelées par un corps de fonction, en
// suivant les fonctions et méthodes du même package qu'il appelle.
type callWalker struct {
	cb     *codebase
	pk     *pkg
	seen   map[*ast.BlockStmt]bool
	found  map[*fn]bool
	direct map[string]bool // fonctions du package appelées au premier niveau
	depth  int
}

func (cb *codebase) collectCalls(pk *pkg, fc *fileCtx, recv *ast.FieldList, params *ast.FieldList, body *ast.BlockStmt) (map[*fn]bool, map[string]bool) {
	w := &callWalker{cb: cb, pk: pk, seen: map[*ast.BlockStmt]bool{}, found: map[*fn]bool{}, direct: map[string]bool{}}
	w.walk(fc, recv, params, body)
	return w.found, w.direct
}

func (w *callWalker) walk(fc *fileCtx, recv, params *ast.FieldList, body *ast.BlockStmt) {
	if body == nil || w.seen[body] || w.depth > 8 {
		return
	}
	w.seen[body] = true
	w.depth++
	defer func() { w.depth-- }()
	env := map[string]typeRef{}
	recvName, recvT := "", ""
	if recv != nil && len(recv.List) > 0 && len(recv.List[0].Names) > 0 {
		recvName, recvT = recv.List[0].Names[0].Name, recvType(recv.List[0].Type)
	}
	addFields := func(fl *ast.FieldList) {
		if fl == nil {
			return
		}
		for _, f := range fl.List {
			for _, n := range f.Names {
				env[n.Name] = typeRef{f.Type, fc}
			}
		}
	}
	addFields(params)
	ast.Inspect(body, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncLit:
			addFields(x.Type.Params)
		case *ast.ValueSpec:
			for _, id := range x.Names {
				if x.Type != nil {
					env[id.Name] = typeRef{x.Type, fc}
				}
			}
		case *ast.CallExpr:
			switch fun := x.Fun.(type) {
			case *ast.Ident:
				if f := w.pk.funcs[fun.Name]; f != nil {
					if w.depth == 1 {
						w.direct[fun.Name] = true
					}
					w.walk(f.file, nil, f.decl.Type.Params, f.decl.Body)
				}
			case *ast.SelectorExpr:
				if id, ok := fun.X.(*ast.Ident); ok && id.Name == recvName {
					if m := w.pk.methods[recvT][fun.Sel.Name]; m != nil {
						w.walk(m.file, m.decl.Recv, m.decl.Type.Params, m.decl.Body)
						return true
					}
				}
				w.resolveCall(fun, env, recvName, recvT, true)
			}
		case *ast.SelectorExpr:
			// Valeur de méthode (decide := svc.ApproveRequest) : rattachée
			// seulement si le type du récepteur est connu.
			w.resolveCall(x, env, recvName, recvT, false)
		}
		return true
	})
}

// resolveCall rattache x.M à une méthode de port : par le type du récepteur
// quand il est connu, sinon — pour un appel — par le seul nom s'il est unique.
func (w *callWalker) resolveCall(sel *ast.SelectorExpr, env map[string]typeRef, recvName, recvT string, call bool) {
	var t *typeRef
	switch x := sel.X.(type) {
	case *ast.Ident:
		if r, ok := env[x.Name]; ok {
			t = &r
		} else if r, ok := w.pk.vars[x.Name]; ok {
			t = &r
		}
	case *ast.SelectorExpr:
		if id, ok := x.X.(*ast.Ident); ok && id.Name == recvName {
			if r, ok := w.pk.fields[recvT][x.Sel.Name]; ok {
				t = &r
			}
		}
	case *ast.CallExpr:
		var callee *funcCtx
		switch f := x.Fun.(type) {
		case *ast.Ident:
			callee = w.pk.funcs[f.Name]
		case *ast.SelectorExpr:
			if id, ok := f.X.(*ast.Ident); ok && id.Name == recvName {
				callee = w.pk.methods[recvT][f.Sel.Name]
			}
		}
		if callee != nil && callee.decl.Type.Results != nil && len(callee.decl.Type.Results.List) > 0 {
			r := typeRef{callee.decl.Type.Results.List[0].Type, callee.file}
			t = &r
		}
	}
	if t != nil {
		if p := w.cb.portOf(*t); p != nil {
			for _, f := range p.Methods {
				if f.Name == sel.Sel.Name {
					w.found[f] = true
				}
			}
			return
		}
	}
	if !call {
		return
	}
	if f := w.cb.uniqueIn(sel.Sel.Name); f != nil && t == nil {
		w.found[f] = true
	}
}

// loadRoutes relève les routes enregistrées par le routeur REST.
func (cb *codebase) loadRoutes() {
	pk := cb.byPath["adapters/in/rest"]
	if pk == nil {
		return
	}
	for _, fc := range pk.files {
		ast.Inspect(fc.ast, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || len(call.Args) != 2 {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || (sel.Sel.Name != "HandleFunc" && sel.Sel.Name != "Handle") {
				return true
			}
			if id, ok := sel.X.(*ast.Ident); !ok || id.Name != "mux" {
				return true
			}
			lit, ok := call.Args[0].(*ast.BasicLit)
			if !ok {
				return true
			}
			pattern, _ := strconv.Unquote(lit.Value)
			r := &route{Pattern: pattern, File: fc.path, Line: cb.pos(pk, call.Pos()), Handler: handlerName(call.Args[1])}
			if m := pk.methods["Handler"][r.Handler]; m != nil {
				r.HFile, r.HLine = m.file.path, cb.pos(pk, m.decl.Pos())
				found, _ := cb.collectCalls(pk, m.file, m.decl.Recv, m.decl.Type.Params, m.decl.Body)
				r.Calls = sortFns(found)
			}
			cb.routes = append(cb.routes, r)
			return true
		})
	}
	sort.SliceStable(cb.routes, func(i, j int) bool { return cb.routes[i].Pattern < cb.routes[j].Pattern })
	notes := map[string]*routeNote{}
	for _, r := range cb.routes {
		slug := strings.ReplaceAll(strings.Trim(r.Pattern, "/"), "/", "-")
		rn := notes[slug]
		if rn == nil {
			title := strings.TrimSuffix(r.Pattern, "/")
			rn = &routeNote{Slug: slug, Title: title}
			notes[slug] = rn
			cb.routeNotes = append(cb.routeNotes, rn)
		}
		rn.Regs = append(rn.Regs, r)
		r.Note = rn
	}
}

// handlerName retourne le handler d'un enregistrement de route : le premier
// sélecteur handle… (ou …Handler) de l'expression.
func handlerName(e ast.Expr) string {
	name, fallback := "", ""
	ast.Inspect(e, func(n ast.Node) bool {
		s, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if name == "" && len(s.Sel.Name) > 6 && strings.HasPrefix(s.Sel.Name, "handle") && unicode.IsUpper(rune(s.Sel.Name[6])) {
			name = s.Sel.Name
		}
		if fallback == "" && strings.HasSuffix(s.Sel.Name, "Handler") {
			fallback = s.Sel.Name
		}
		return true
	})
	if name == "" {
		return fallback
	}
	return name
}

// resolveRoute retourne l'enregistrement de route qui sert le chemin p :
// motif exact, sinon le plus long motif de sous-arbre (terminé par /).
func (cb *codebase) resolveRoute(p string) *route {
	if i := strings.IndexAny(p, "?#"); i >= 0 {
		p = p[:i]
	}
	var best *route
	for _, r := range cb.routes {
		if r.Pattern == p {
			return r
		}
		if strings.HasSuffix(r.Pattern, "/") && strings.HasPrefix(p, r.Pattern) && (best == nil || len(r.Pattern) > len(best.Pattern)) {
			best = r
		}
	}
	return best
}

// loadCommands relève l'arbre des commandes cobra du CLI.
func (cb *codebase) loadCommands() {
	pk := cb.byPath["adapters/in/cli"]
	if pk == nil {
		return
	}
	byVar := map[string]*command{}
	use := map[*command]string{}
	for _, fc := range pk.files {
		for _, decl := range fc.ast.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.VAR {
				continue
			}
			for _, spec := range gd.Specs {
				vs := spec.(*ast.ValueSpec)
				for i, n := range vs.Names {
					if i >= len(vs.Values) {
						continue
					}
					cl := cobraLit(vs.Values[i])
					if cl == nil {
						continue
					}
					c := &command{Var: n.Name, File: fc.path, Line: cb.pos(pk, n.Pos()), runFuncs: map[string]bool{}, fc: fc, pkg: pk}
					for _, el := range cl.Elts {
						kv, ok := el.(*ast.KeyValueExpr)
						if !ok {
							continue
						}
						key, _ := kv.Key.(*ast.Ident)
						if key == nil {
							continue
						}
						switch key.Name {
						case "Use":
							if lit, ok := kv.Value.(*ast.BasicLit); ok {
								s, _ := strconv.Unquote(lit.Value)
								use[c] = strings.Fields(s + " x")[0]
							}
						case "RunE", "Run":
							if fl, ok := kv.Value.(*ast.FuncLit); ok {
								c.body = fl
							}
						}
					}
					byVar[n.Name] = c
				}
			}
		}
	}
	for _, fc := range pk.files {
		ast.Inspect(fc.ast, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || sel.Sel.Name != "AddCommand" {
				return true
			}
			id, ok := sel.X.(*ast.Ident)
			if !ok || byVar[id.Name] == nil {
				return true
			}
			parent := byVar[id.Name]
			for _, a := range call.Args {
				if aid, ok := a.(*ast.Ident); ok && byVar[aid.Name] != nil && byVar[aid.Name].Parent == nil {
					child := byVar[aid.Name]
					child.Parent = parent
					parent.Children = append(parent.Children, child)
				}
			}
			return true
		})
	}
	for _, c := range byVar {
		var words []string
		for x := c; x != nil; x = x.Parent {
			words = append([]string{use[x]}, words...)
		}
		if words[0] != "myr" {
			continue
		}
		c.Words = words
		if c.body != nil {
			found, direct := cb.collectCalls(pk, c.fc, nil, c.body.Type.Params, c.body.Body)
			c.Calls, c.runFuncs = sortFns(found), direct
		}
		cb.cmds = append(cb.cmds, c)
	}
	sort.Slice(cb.cmds, func(i, j int) bool { return cb.cmds[i].Name() < cb.cmds[j].Name() })
}

func cobraLit(e ast.Expr) *ast.CompositeLit {
	if u, ok := e.(*ast.UnaryExpr); ok {
		e = u.X
	}
	cl, ok := e.(*ast.CompositeLit)
	if !ok {
		return nil
	}
	if sel, ok := cl.Type.(*ast.SelectorExpr); ok && sel.Sel.Name == "Command" {
		return cl
	}
	return nil
}

// resolveCmd retourne la commande la plus longue désignée par les mots qui
// suivent « myr ».
func (cb *codebase) resolveCmd(words []string) *command {
	var best *command
	for _, c := range cb.cmds {
		w := c.Words[1:]
		if len(w) == 0 || len(w) > len(words) {
			continue
		}
		match := true
		for i := range w {
			if w[i] != words[i] {
				match = false
				break
			}
		}
		if match && (best == nil || len(c.Words) > len(best.Words)) {
			best = c
		}
	}
	return best
}

// loadTests relève les fonctions de test d'un fichier.
func (cb *codebase) loadTests(rel, dir string, f *ast.File, fset *token.FileSet) {
	subject := strings.TrimSuffix(strings.TrimSuffix(dir, "/tests"), "/integration")
	pk := cb.byPath[subject]
	if pk == nil {
		return
	}
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Recv != nil || !strings.HasPrefix(fd.Name.Name, "Test") || fd.Body == nil {
			continue
		}
		t := &test{Name: fd.Name.Name, File: rel, Line: fset.Position(fd.Pos()).Line, Subject: pk,
			IDs: map[string]bool{}, calls: map[string]bool{}, idents: map[string]bool{}}
		comment := ""
		if fd.Doc != nil {
			comment = fd.Doc.Text()
		}
		t.parseComment(comment)
		for _, m := range testIDRe.FindAllString(t.Name+" "+comment, -1) {
			t.IDs[m] = true
		}
		for _, id := range expandUCs(t.Name + " " + comment) {
			t.IDs[id] = true
		}
		ast.Inspect(fd.Body, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.CallExpr:
				switch fun := x.Fun.(type) {
				case *ast.SelectorExpr:
					t.calls[fun.Sel.Name] = true
				case *ast.Ident:
					t.idents[fun.Name] = true
				}
			case *ast.BasicLit:
				if x.Kind == token.STRING {
					if s, err := strconv.Unquote(x.Value); err == nil && (strings.HasPrefix(s, "/api/") || s == "/metrics") {
						t.literals = append(t.literals, s)
					}
				}
			}
			return true
		})
		pk.Tests = append(pk.Tests, t)
		cb.tests = append(cb.tests, t)
	}
}

var tagLineRe = regexp.MustCompile(`^[/\s]*@(brief|input|expect|details)\s*(.*)$`)

// parseComment lit le commentaire structuré @brief / @input / @expect.
func (t *test) parseComment(c string) {
	cur := ""
	for _, l := range strings.Split(c, "\n") {
		if m := tagLineRe.FindStringSubmatch(l); m != nil {
			cur = m[1]
			switch cur {
			case "brief":
				t.Brief = m[2]
			case "input":
				t.Input = m[2]
			case "expect":
				t.Expect = m[2]
			case "details":
				t.Brief = strings.TrimSpace(t.Brief + " @details " + m[2])
			}
			continue
		}
		l = strings.TrimSpace(strings.TrimLeft(l, "/ "))
		if l == "" || cur == "" {
			continue
		}
		switch cur {
		case "brief", "details":
			t.Brief += " " + l
		case "input":
			t.Input += " " + l
		case "expect":
			t.Expect += " " + l
		}
	}
}

// linkTests rattache chaque test aux fonctions, routes et commandes qu'il exerce.
func (cb *codebase) linkTests() {
	for _, t := range cb.tests {
		fns := map[*fn]bool{}
		for _, p := range cb.ports {
			mine := p.Pkg == t.Subject
			for _, im := range p.Impls {
				if im.Pkg == t.Subject {
					mine = true
				}
			}
			if !mine {
				continue
			}
			for _, f := range p.Methods {
				if t.calls[f.Name] {
					fns[f] = true
				}
			}
		}
		t.Fns = sortFns(fns)
		for _, f := range t.Fns {
			f.Tests = append(f.Tests, t)
		}
		// Un chemin en littéral n'exerce une route que dans un test de
		// l'adapter REST (ailleurs, c'est une simple donnée de test).
		routes := map[*routeNote]bool{}
		for _, l := range t.literals {
			if r := cb.resolveRoute(l); r != nil && t.Subject.Path == "adapters/in/rest" {
				routes[r.Note] = true
			}
		}
		for _, rn := range cb.routeNotes {
			if routes[rn] {
				t.Routes = append(t.Routes, rn)
				rn.Tests = append(rn.Tests, t)
			}
		}
		if t.Subject.Path == "adapters/in/cli" {
			for _, c := range cb.cmds {
				for f := range c.runFuncs {
					if t.idents[f] {
						t.Cmds = append(t.Cmds, c)
						c.Tests = append(c.Tests, t)
						break
					}
				}
			}
		}
	}
}

func sortFns(set map[*fn]bool) []*fn {
	out := make([]*fn, 0, len(set))
	for f := range set {
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key() < out[j].Key() })
	return out
}
