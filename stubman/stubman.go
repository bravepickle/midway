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
	"time"
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

	return buf.String()
}

func init() {
	InitTemplates()
}

// AddGuiHandlers add all handlers for income requests that come to stub service
func AddStubmanCrudHandlers(prefix string, mux *http.ServeMux) {
	//	if Config.
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

		page := Page{HomePage: true, Data: models}
		RenderPage(`index.tpl`, page, w)
	})

	// create
	mux.HandleFunc(pcat.fullPath(`/create/`), func(w http.ResponseWriter, req *http.Request) {
		if req.Method == `POST` {
			req.ParseForm()
			stub := NewStubFromRequest(req)
			repo := NewStubRepo(nil)

			id, err := repo.Insert(stub)
			if err != nil {
				w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusBadRequest)

				return
			}

			w.Header().Add(`Location`, fmt.Sprintf(pcat.fullPath(`/edit/%d`), id))
			w.WriteHeader(http.StatusFound)

			return
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
			w.WriteHeader(http.StatusNotFound)

			return
		}

		repo := NewStubRepo(nil)
		idNum, err := strconv.Atoi(id)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		model, err := repo.Find(idNum)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		if model.Id == 0 {
			w.Write([]byte(`Not Found`))
			w.WriteHeader(http.StatusNotFound)

			return
		}

		if req.Method == `POST` {
			req.ParseForm()
			stub := NewStubFromRequest(req)
			stub.Id = int64(idNum)

			//			fmt.Printf("Model: %v\n", stub)

			err := repo.Update(stub)
			if err != nil {
				w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusBadRequest)

				return
			}

			w.Header().Add(`Location`, fmt.Sprintf(pcat.fullPath(`/edit/%d`), idNum))
			w.WriteHeader(http.StatusFound)

			return
		}

		page := Page{EditPage: true, Data: model}
		RenderPage(`edit.tpl`, page, w)
	})

	// delete
	mux.HandleFunc(pcat.fullPath(`/delete/`), func(w http.ResponseWriter, req *http.Request) {
		id := pathRegId.FindString(req.URL.Path)
		if id == `` {
			w.Write([]byte(`Not Found`))
			w.WriteHeader(http.StatusNotFound)

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

		deleted, err := repo.Delete(idNum)
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			w.WriteHeader(500)

			return
		}

		if deleted {
			w.Write([]byte(`Not Found`))
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	})
}

func NewStubFromRequest(req *http.Request) *Stub {
	stub := &Stub{Created: time.Now()}
	log.Println(`REQUEST BODY: `, string(req.Form.Get(`request[headers][]`)))

	stub.Name = string(req.Form.Get(`name`))
	stub.RequestMethod = string(req.Form.Get(`request_method`))
	stub.RequestMethod = string(req.Form.Get(`request_uri`))
	stub.RequestParsed.Body = string(req.Form.Get(`request[body]`))
	for _, val := range req.Form[`request[headers][]`] {
		stub.RequestParsed.Headers = append(stub.RequestParsed.Headers, string(val))
	}

	stub.ResponseParsed.Body = string(req.Form.Get(`response[body]`))
	for _, val := range req.Form[`response[headers][]`] {
		stub.ResponseParsed.Headers = append(stub.ResponseParsed.Headers, string(val))
	}

	stub.Encode()

	return stub
}
