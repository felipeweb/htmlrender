package htmlrender

import (
	"html/template"
	"net/http"
)

// Engine is the generic interface for all responses.
type Engine interface {
	Render(http.ResponseWriter, interface{}) error
}

// Head defines the basic ContentType and Status fields.
type Head struct {
	ContentType string
	Status      int
}

// HTML built-in renderer.
type HTML struct {
	Head
	Name      string
	Templates *template.Template
}

// Write outputs the header content.
func (h Head) Write(w http.ResponseWriter) {
	w.Header().Set(ContentType, h.ContentType)
	w.WriteHeader(h.Status)
}

// Render a HTML response.
func (h HTML) Render(w http.ResponseWriter, binding interface{}) error {
	// Retrieve a buffer from the pool to write to.
	out := bufPool.Get()
	err := h.Templates.ExecuteTemplate(out, h.Name, binding)
	if err != nil {
		return err
	}

	h.Head.Write(w)
	_, err = out.WriteTo(w)
	if err != nil {
		return err
	}

	// Return the buffer to the pool.
	bufPool.Put(out)
	return nil
}
