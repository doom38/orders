package tmpl

import "assets"
import "strings"
import "log"
import "path"
import "fmt"
import "html/template"

// LoadFromAssets load all templates file '.tmpl.html' stored in 'assets'
// package.
func LoadFromAssets(dirs ...string) (t Templater, err error) {
	t = make(Templater)
	var files []string
	var data []byte
	for _, dir := range dirs {
		files, err = assets.AssetDir(dir)
		if err != nil {
			err = fmt.Errorf("Fail to list templates in dir %q: %v", dir, err)
			return
		}
		for _, file := range files {
			if strings.HasSuffix(file, ".tmpl.html") {
				p := path.Join(dir, file)

				data, err = assets.Asset(p)
				if err != nil {
					err = fmt.Errorf("%q: %v", p, err)
					return
				}

				// Fail fast when the template is well formed:
				_, err = template.New("test").Parse(string(data))
				if err != nil {
					err = fmt.Errorf("%q: %v", p, err)
					return
				}

				t[file] = data
				log.Println("TMPL", p)
			}
		}
	}
	return
}
