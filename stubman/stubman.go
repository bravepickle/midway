// stub service entry point

package stubman

import (
	"bytes"
	"fmt"
	//	"io/ioutil"
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
	//	if Config.
	pcat := pathConcat{prefix}

	// static files
	pathRegExt := regexp.MustCompile(`\.\w{2,4}$`)
	mux.HandleFunc(pcat.fullPath(`/static/`), func(w http.ResponseWriter, req *http.Request) {
		ext := pathRegExt.FindString(req.URL.Path)
		if ext == `` {
			ext = `unknown`
		}

		w.Header().Add(`X-Test-Extension`, ext)
		w.WriteHeader(403)
	})

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

		page := Page{HomePage: true, Data: models}
		RenderPage(`index.tpl`, page, w)
	})

	// create
	mux.HandleFunc(pcat.fullPath(`/create/`), func(w http.ResponseWriter, req *http.Request) {
		if req.Method == `POST` {

			//			body, _ := ioutil.ReadAll(req.Body)
			req.ParseForm()
			log.Println(`REQUEST BODY: `, string(req.Form.Get(`request[headers][]`)))
			log.Printf(`REQUEST BODY2: %v`, req.Form)

		}
		model := NewNullObjectStub()
		page := Page{CreatePage: true, Data: model}
		RenderPage(`create.tpl`, page, w)
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

		page := Page{EditPage: true, Data: model}
		RenderPage(`edit.tpl`, page, w)
	})
}
