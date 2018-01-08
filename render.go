package htmlrender

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	// ContentTypeHTML value
	ContentTypeHTML = "text/html"
	// ContentType header constant.
	ContentType = "Content-Type"
	// Default character encoding.
	defaultCharset = "UTF-8"
)

// Options is a struct for specifying configuration options for the render.Render object.
type Options struct {
	// Directory to load templates. Default is "views".
	Directory string
	// Asset function to use in place of directory. Defaults to nil.
	Asset func(name string) ([]byte, error)
	// AssetNames function to use in place of directory. Defaults to nil.
	AssetNames func() []string
	// Funcs is a slice of FuncMaps to apply to the template upon compilation. This is useful for helper functions. Defaults to [].
	Funcs []template.FuncMap
	// Appends the given character set to the Content-Type header. Default is "UTF-8".
	Charset string
}

// Render is a service that provides functions for easily writing HTML templates out to a HTTP Response.
type Render struct {
	// Customize Secure with an Options struct.
	opt             Options
	templates       *template.Template
	compiledCharset string
}

// New constructs a new Render instance with the supplied options.
func New(options ...Options) *Render {
	var o Options
	if len(options) == 0 {
		o = Options{}
	} else {
		o = options[0]
	}

	r := Render{
		opt: o,
	}

	r.prepareOptions()
	r.compileTemplates()

	return &r
}

func (r *Render) prepareOptions() {
	// Fill in the defaults if need be.
	if len(r.opt.Charset) == 0 {
		r.opt.Charset = defaultCharset
	}
	r.compiledCharset = "; charset=" + r.opt.Charset

	if len(r.opt.Directory) == 0 {
		r.opt.Directory = "views"
	}
}

func (r *Render) compileTemplates() {
	if r.opt.Asset == nil || r.opt.AssetNames == nil {
		r.compileTemplatesFromDir()
		return
	}
	r.compileTemplatesFromAsset()
}

func (r *Render) compileTemplatesFromAsset() {
	dir := r.opt.Directory
	r.templates = template.New(dir)
	r.templates.Delims("{{", "}}")
	for _, path := range r.opt.AssetNames() {
		if !strings.HasPrefix(path, dir) {
			continue
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			panic(err)
		}

		ext := ""
		if strings.Contains(rel, ".") {
			ext = "." + strings.Join(strings.Split(rel, ".")[1:], ".")
		}
		extension := ".html"
		if ext == extension {
			fmt.Println(path)
			buf, err := r.opt.Asset(path)
			if err != nil {
				panic(err)
			}

			name := (rel[0 : len(rel)-len(ext)])
			tmpl := r.templates.New(filepath.ToSlash(name))

			// Add our funcmaps.
			for _, funcs := range r.opt.Funcs {
				tmpl = tmpl.Funcs(funcs)
			}

			// Break out if this parsing fails. We don't want any silent server starts.
			template.Must(tmpl.Parse(string(buf)))
		}
	}
}

func (r *Render) compileTemplatesFromDir() {
	dir := r.opt.Directory
	r.templates = template.New(dir)
	r.templates.Delims("{{", "}}")

	// Walk the supplied directory and compile any files that match our extension list.
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error { // nolint
		// Fix same-extension-dirs bug: some dir might be named to: "local.html".
		// These dirs should be excluded as they are not valid golang templates, but files under
		// them should be treat as normal.
		// If is a dir, return immediately (dir is not a valid golang template).
		if info == nil || info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		ext := ""
		extension := ".html"

		if strings.Contains(rel, ".") {
			ext = filepath.Ext(rel)
		}

		if ext == extension {
			buf, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}

			name := (rel[0 : len(rel)-len(ext)])
			tmpl := r.templates.New(filepath.ToSlash(name))

			// Add our funcmaps.
			for _, funcs := range r.opt.Funcs {
				tmpl = tmpl.Funcs(funcs)
			}

			// Break out if this parsing fails. We don't want any silent server starts.
			template.Must(tmpl.Parse(string(buf)))
		}
		return nil
	})
}

// TemplateLookup is a wrapper around template.Lookup and returns
// the template with the given name that is associated with t, or nil
// if there is no such template.
func (r *Render) TemplateLookup(t string) *template.Template {
	return r.templates.Lookup(t)
}

// Render is the generic function called by HTML, and can be called by custom implementations.
func (r *Render) Render(w http.ResponseWriter, e Engine, data interface{}) error {
	err := e.Render(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return err
}

// HTML builds up the response from the specified template and bindings.
func (r *Render) HTML(w http.ResponseWriter, status int, name string, binding interface{}) error {

	head := Head{
		ContentType: ContentTypeHTML + r.compiledCharset,
		Status:      status,
	}

	h := HTML{
		Head:      head,
		Name:      name,
		Templates: r.templates,
	}

	return r.Render(w, h, binding)
}
