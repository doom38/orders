package tmpl

import (
	"fmt"
	"html/template"
	"io"
	"log"
)

// Templater holds all application templates by name.
type Templater map[string][]byte

// Render writes the rendered result in 'w' using 'data' values.
// The rendering is done using the 'names' templates.
func (t Templater) Render(w io.Writer, data interface{}, names ...string) {
	r := template.New("root")
	for _, name := range names {
		_, err := r.Parse(string(t[name]))
		if err != nil {
			log.Printf("Can't parse template %q: %v", name, err)
		}
	}

	err := r.Execute(w, data)
	if err != nil {
		fmt.Fprintf(w, "Can't execute template: %v\n", err)
		log.Printf("Can't execute template: %v\n", err)
	}
}
