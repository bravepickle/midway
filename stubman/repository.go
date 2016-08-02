// contains DB repository for stubs db
package stubman

import (
	"database/sql"
	"encoding/json"
	//	"fmt"
	"time"
)

const stubTable = `stub`

type ResponseStub struct {
	Headers []string
	Body    string
}

type RequestStub struct {
	Headers []string
	Body    string
}

type Stub struct {
	Id             int
	Name           string
	RequestUri     string
	RequestMethod  string
	Request        string
	RequestParsed  RequestStub
	Response       string
	ResponseParsed ResponseStub
	Created        time.Time
}

// Parse parses values from Request and Response and puts them to RequestParsed, ResponseParsed accordingly
func (s *Stub) Decode() {
	if s.Request != `` {
		json.Unmarshal([]byte(s.Request), &s.RequestParsed)
	} else {
		s.RequestParsed = RequestStub{}
	}

	if s.Response != `` {
		json.Unmarshal([]byte(s.Response), &s.ResponseParsed)
	} else {
		s.ResponseParsed = ResponseStub{}
	}
}

// Encode encodes to string all structs
func (s *Stub) Encode() {
	var raw []byte
	raw, _ = json.Marshal(s.RequestParsed)
	s.Request = string(raw)

	raw, _ = json.Marshal(s.ResponseParsed)
	s.Response = string(raw)
}

type StubRepo struct {
	Table string
	Conn  *sql.DB
}

func (r *StubRepo) FindAll() ([]Stub, error) {
	var result []Stub

	rows, err := r.Conn.Query("SELECT id, name, request_method, " +
		"request_uri, request, response, created FROM stub")
	if err != nil {
		return result, err
	}

	for rows.Next() {
		model := Stub{}

		if err := rows.Scan(&model.Id, &model.Name,
			&model.RequestMethod, &model.RequestUri,
			&model.Request, &model.Response, &model.Created); err != nil {
			return result, err
		}
		model.Decode()

		//		fmt.Printf("Model: %v\n", model)

		result = append(result, model)
	}

	if err := rows.Err(); err != nil {
		return result, err
	}

	return result, nil
}

// Find find model by ID
func (r *StubRepo) Find(id int) (Stub, error) {
	model := Stub{}

	rows, err := r.Conn.Query("SELECT id, name, request_method, "+
		"request_uri, request, response, created FROM stub WHERE id = $1", id)
	if err != nil {
		return model, err
	}

	for rows.Next() {
		if err := rows.Scan(&model.Id, &model.Name,
			&model.RequestMethod, &model.RequestUri,
			&model.Request, &model.Response, &model.Created); err != nil {
			return model, err
		}
		model.Decode()
	}

	if err := rows.Err(); err != nil {
		return model, err
	}

	return model, nil
}

func NewStubRepo(db *Db) *StubRepo {
	if db == nil {
		db = DefaultDb
	}

	return &StubRepo{Table: stubTable, Conn: db.Connection}
}

func NewNullObjectStub() *Stub {
	return &Stub{
		RequestMethod: `GET`,
		RequestParsed: RequestStub{Headers: []string{`Content-Type: application/json`}},
		Request:       `{"headers": ["Content-Type: application/json"], "body":""}`}
}
