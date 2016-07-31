// stub service entry point

package stubman

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

const StaticPath = `static`
const viewsDir = `views`
const StubmanPathPrefix = `stubman`

type pathConcat struct {
	prefix string
}

// fullpath append prefix and path
func (p *pathConcat) fullPath(path string) string {
	buf := bytes.NewBufferString(p.prefix)
	buf.WriteString(path)

	//	if Debug {
	fmt.Println(`Generated path for Stubman: `, buf.String())
	//	}

	return buf.String()
}

func init() {
	InitTemplates()
}

// AddGuiHandlers add all handlers for income requests that come to stub service
func AddStubmanCrudHandlers(prefix string, mux *http.ServeMux) {
	pcat := pathConcat{prefix}

	// list all stubs
	mux.HandleFunc(pcat.fullPath(`/`), func(w http.ResponseWriter, req *http.Request) {
		repo := NewStubRepo(nil)
		models, err := repo.FindAll()

		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			w.WriteHeader(500)

			return
		}

		RenderPage(`index.tpl`, models, w)
	})

	pathRegId := regexp.MustCompile(`\d+$`)
	// edit
	mux.HandleFunc(pcat.fullPath(`/edit/`), func(w http.ResponseWriter, req *http.Request) {
		id := pathRegId.FindString(req.URL.Path)
		if id == `` {
			w.Write([]byte(`Not Found`))
			w.WriteHeader(404)

			return
		}

		repo := NewStubRepo(nil)
		idNum, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			w.WriteHeader(400)

			return
		}

		model, err := repo.Find(idNum)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			w.WriteHeader(400)

			return
		}

		if model.Id == 0 {
			w.Write([]byte(`Not Found`))
			w.WriteHeader(404)

			return
		}

		RenderPage(`edit.tpl`, model, w)
	})
}
