package gox

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"regexp"
	"strings"
)

var populateRe = regexp.MustCompile(`<meta name="populate" content="(\w+:\w+)" ?\/?>`)

func (g *Gox) AddTemplates(dir fs.FS, names map[string]string) {
	g.templateLock.Lock()
	defer g.templateLock.Unlock()

	for name, path := range names {
		tmpl := template.New(name)

		tmpl.Funcs(template.FuncMap{
			"include": func(name string, data any) template.HTML {
				t, ok := g.templates[name]
				if !ok {
					g.goxLogger.Errorf("template %s does not exist", name)
					return ""
				}

				buf := &bytes.Buffer{}

				err := t.ExecuteTemplate(buf, name, data)
				if err != nil {
					g.goxLogger.Errorf("failed to execute template %s: %s", name, err)
					return ""
				}

				return template.HTML(buf.String())
			},
		})

		f, err := dir.Open(path)
		if err != nil {
			g.goxLogger.Errorf("failed to open template %s: %s", name, err)
			continue
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			g.goxLogger.Errorf("failed to read template %s: %s", name, err)
			continue
		}

		tmpl, err = tmpl.Parse(string(data))
		if err != nil {
			g.goxLogger.Errorf("failed to parse template %s: %s", name, err)
			continue
		}

		g.templates[name] = tmpl
	}
}

func (g *Gox) RenderTemplate(name string, data any) (string, error) {
	g.templateLock.Lock()
	defer g.templateLock.Unlock()

	t, ok := g.templates[name]
	if !ok {
		return "", fmt.Errorf("template %s does not exist", name)
	}

	buf := &bytes.Buffer{}

	err := t.Execute(buf, data)
	if err != nil {
		return "", err
	}

	d := buf.String()

	if populateRe.MatchString(d) {
		matches := populateRe.FindAllStringSubmatch(d, -1)
		name := matches[0][1]

		parts := strings.Split(name, ":")

		populateName := parts[0]
		field := parts[1]

		t, ok := g.templates[populateName]
		if !ok {
			return "", fmt.Errorf("template %s does not exist", name)
		}

		buf := &bytes.Buffer{}

		err := t.Execute(buf, map[string]template.HTML{
			field: template.HTML(d),
		})
		if err != nil {
			return "", err
		}

		return buf.String(), nil
	}

	return d, nil
}
