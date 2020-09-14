package template

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/markbates/pkger"
)

func New(name string, framework string, typescript bool) {
	os.Mkdir(name, 0744)
	switch framework {
	case "none", "blank":
		handler, _ := pkger.Open("/templates/blank.ts")
		copy(handler, path.Join(name, "main.ts"), !typescript)

	case "react":
		handler, _ := pkger.Open("/templates/react.tsx")
		copy(handler, path.Join(name, "main.ts"), !typescript)
	}
}

func formatTemplate(name string) string {
	vars := struct {
		Title     string
		ScriptSrc string
	}{
		Title:     name,
		ScriptSrc: "main.js",
	}
	handler, _ := pkger.Open("/templates/index.html")
	b, _ := ioutil.ReadAll(handler)
	tmpl, _ := template.New("index").Parse(string(b))
	buf := bytes.NewBufferString("")
	tmpl.Execute(buf, vars)
	return buf.String()
}

func copy(h io.Reader, dest string, convertFromTS bool) {
	data, _ := ioutil.ReadAll(h)
	if convertFromTS {
		result := api.Transform(string(data), api.TransformOptions{
			Format: api.FormatESModule,
		})
		ioutil.WriteFile(dest, result.JS, 0744)
	} else {
		ioutil.WriteFile(dest, data, 0744)
	}
}
