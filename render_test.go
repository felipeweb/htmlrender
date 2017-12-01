package htmlrender

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTMLBad(t *testing.T) {
	render := New(Options{
		Directory: "testdata/basic",
	})

	var err error
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err = render.HTML(w, http.StatusOK, "nope", nil)
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	h.ServeHTTP(res, req)

	expectNotNil(t, err)
	expect(t, res.Code, 500)
	expect(t, res.Body.String(), "html/template: \"nope\" is undefined\n")
}

func TestHTMLBasic(t *testing.T) {
	render := New(Options{
		Directory: "testdata/basic",
	})

	var err error
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err = render.HTML(w, http.StatusOK, "hello", "gophers")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	h.ServeHTTP(res, req)

	expectNil(t, err)
	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentTypeHTML+"; charset=UTF-8")
	expect(t, res.Body.String(), "<h1>Hello gophers</h1>")
}

func TestHTMLFuncs(t *testing.T) {
	render := New(Options{
		Directory: "testdata/custom_funcs",
		Funcs: []template.FuncMap{
			{
				"myCustomFunc": func() string {
					return "My custom function"
				},
			},
		},
	})

	var err error
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err = render.HTML(w, http.StatusOK, "index", "gophers")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	h.ServeHTTP(res, req)

	expectNil(t, err)
	expect(t, res.Body.String(), "My custom function")
}

func TestHTMLNested(t *testing.T) {
	render := New(Options{
		Directory: "testdata/basic",
	})

	var err error
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err = render.HTML(w, http.StatusOK, "admin/index", "gophers")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	h.ServeHTTP(res, req)

	expectNil(t, err)
	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentTypeHTML+"; charset=UTF-8")
	expect(t, res.Body.String(), "<h1>Admin gophers</h1>")
}

func TestHTMLBadPath(t *testing.T) {
	render := New(Options{
		Directory: "../../../../../../../../../../../../../../../../testdata/basic",
	})

	var err error
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err = render.HTML(w, http.StatusOK, "hello", "gophers")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	h.ServeHTTP(res, req)

	expectNotNil(t, err)
	expect(t, res.Code, 500)
}

func TestHTMLDefaultCharset(t *testing.T) {
	render := New(Options{
		Directory: "testdata/basic",
	})

	var err error
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err = render.HTML(w, http.StatusOK, "hello", "gophers")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foo", nil)
	h.ServeHTTP(res, req)

	expectNil(t, err)
	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentTypeHTML+"; charset=UTF-8")

	expect(t, res.Body.String(), "<h1>Hello gophers</h1>")
}

func TestHTMLNoRace(t *testing.T) {
	// This test used to fail if run with -race
	render := New(Options{
		Directory: "testdata/basic",
	})

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := render.HTML(w, http.StatusOK, "hello", "gophers")
		expectNil(t, err)
	})

	done := make(chan bool)
	doreq := func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/foo", nil)

		h.ServeHTTP(res, req)

		expect(t, res.Code, 200)
		expect(t, res.Header().Get(ContentType), ContentTypeHTML+"; charset=UTF-8")
		expect(t, res.Body.String(), "<h1>Hello gophers</h1>")
		done <- true
	}
	// Run two requests to check there is no race condition
	go doreq()
	go doreq()
	<-done
	<-done
}
