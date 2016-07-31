// rendering functionality of templates
package stubman

import (
	"fmt"
	//	"io/ioutil"
	"io"
	"log"
	"path/filepath"
	"text/template"
)

const baseTemplate = `base.tpl`

var parsedTemplates map[string]*template.Template

// RenderPage renders template with params and returns resulting page to given output
func RenderPage(tpl string, p interface{}, w io.Writer) {
	tp, ok := parsedTemplates[tpl]

	if !ok {
		log.Fatalln(fmt.Sprintf(`Failed to find template: %s`, tpl))

		return
	}

	if err := tp.ExecuteTemplate(w, `base`, p); err != nil {
		log.Fatalln(err.Error())
	}
}

func InitTemplates() {
	tplNames := []string{`index.tpl`, `edit.tpl`}
	parsedTemplates = make(map[string]*template.Template)
	sep := string(filepath.Separator)
	viewsPrefix := StubmanPathPrefix + sep + viewsDir + sep

	for _, name := range tplNames {
		parsed, err := template.ParseFiles(viewsPrefix+name, viewsPrefix+baseTemplate)
		parsedTemplates[name] = template.Must(parsed, err)
	}
}
