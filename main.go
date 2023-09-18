package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/akamensky/argparse"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)



type Finder struct {
	line   int
	col    int
	result ast.Node
	root   ast.Node
}

var (
	re = regexp.MustCompile(`\[([0-9+])\]`)
)

func pathToJSON6901(path string) string {
	path = strings.Replace(path, "$.", "/", 1)
	path = re.ReplaceAllString(path, "/$1")
	path = strings.ReplaceAll(path, ".", "/")
	return path
}

func NewFinder(line int, col int) *Finder {
	return &Finder{
		line: line,
		col: col,
		result: nil,
		root: nil,
	}
}

func (f *Finder) GetResult() ast.Node {
	return f.result
}

func (f *Finder) GetRoot() ast.Node {
	return f.root
}

func (f *Finder) Analyze(docs []*ast.DocumentNode) error {
	for _, doc := range docs {
		f.root = doc
		ast.Walk(f, doc)
		if f.GetResult() != nil {
			return nil
		}
	}
	return fmt.Errorf("no node found at line %d, col %d", f.line, f.col)
}

func (f *Finder) Visit(n ast.Node) ast.Visitor {
	length := len(n.GetToken().Value)
	pos := n.GetToken().Position
	if (pos.Line == f.line) && (pos.Column <= f.col) && (pos.Column + length > f.col) {
		f.result = n
	}
	return f
}


func readFile(filepath string, stdin string) ([]byte, string) {
	var (
		path string
		err error
		file io.Reader
	)

	if filepath != "" {
		path = filepath
		file, err = os.Open(filepath)
		if err != nil {
			log.Fatalf("failed to open file %s: %s", filepath, err)
		}
	} else {
		path = stdin
		file = os.Stdin
	}

	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read input: %s", err)
	}

	return content, path
}

func resolve(target string, n ast.Node) {
	root := path.Dir(target)
	file := n.GetToken().Value
	res := path.Join(root, file)
	fmt.Printf("%s\n", res)
}


func jsonPathAtNode(target ast.Node) {
	fmt.Printf("%s\n", target.GetPath())
}

func patchPathAtNode(target ast.Node) {
	fmt.Printf("%s\n", pathToJSON6901(target.GetPath()))
}

func main() {

	p := argparse.NewParser("kustomize-yaml-helper", "inspect kustomization file")
	filepath := p.String("f", "file",  &argparse.Options{
		Required: false,
		Help: "path to input file",
	})
	stdin := p.String("s", "stdin", &argparse.Options{
		Required: false,
		Help: "input file name, actual content is read from stdin",
	})

	line := p.Int("l", "line",  &argparse.Options{
		Required: true,
		Help: "inspect YAML at given line",
	})

	col := p.Int("c", "col",  &argparse.Options{
		Required: true,
		Help: "inspect YAML at given column",
	})

	opts := []string{"resolve", "json-path", "patch-path"}
	action := p.Selector("a", "action", opts, &argparse.Options{
		Required: true,
		Help: "select action to perform",
		Default: "resolve",
	})

	err := p.Parse(os.Args)
	if err != nil {
		fmt.Print(p.Usage(err))
		os.Exit(1)
	}

	if (*filepath == "" && *stdin == "") || (*filepath != "" && *stdin != "") {
		log.Fatalf("you must provide either --file=<path> or --stdin=<filename> arugment %v %v", filepath, stdin)
	}

	content, path := readFile(*filepath, *stdin)

	astNode, err := parser.ParseBytes(content, parser.ParseComments)
	if err != nil {
		log.Fatalf("could not parse yaml-file: %s", err)
	}

	finder := NewFinder(*line, *col)
	if err = finder.Analyze(astNode.Docs); err != nil {
		log.Fatalf("%s\n", err)
	}

	switch *action {
		case "resolve":
		resolve(path, finder.GetResult())
		case "json-path":
		jsonPathAtNode(finder.GetResult())
		case "patch-path":
		patchPathAtNode(finder.GetResult())
	}
}
