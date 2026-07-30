package web

import (
	"bytes"
	"html/template"
	"io/fs"
	"net/http"
	"time"
)

type htmlRenderer struct {
	templateFS      fs.FS
	sharedTemplates *template.Template
}

type templateData[TPage any] struct {
	Authed bool
	Page   TPage
}

// Eventually, if we want multiple type parameters:
// type templateData[TPage any, TFlash any] struct {
//     Authed bool
//     Page   TPage
//     Flash  TFlash
// }

// seems you can't use a generic in a const
var noTemplateData templateData[struct{}] = templateData[struct{}]{}

// The newHTMLRenderer function creates a new htmlRenderer containing a shared
// set of parsed templates with support for any custom template functions.
func newHTMLRenderer(templateFS fs.FS, sharedTemplateFiles ...string) (*htmlRenderer, error) {
	funcs := template.FuncMap{
		"now": time.Now,
		// Other custom template functions go here...
	}

	sharedTemplates, err := template.New("").Funcs(funcs).ParseFS(templateFS, sharedTemplateFiles...)
	if err != nil {
		return nil, err
	}

	r := &htmlRenderer{
		templateFS:      templateFS,
		sharedTemplates: sharedTemplates,
	}

	return r, nil
}

// The render method clones the shared template set, optionally parses additional
// templates, executes the named template with the supplied data, and writes the
// response.

// No longer a method since otherwise we can't use the generic
// func (h *htmlRenderer) render[Tpage any](w http.ResponseWriter, status int, data templateData[Tpage], templateName string, additionalTemplateFiles ...string) error {
func render[Tpage any](h *htmlRenderer, w http.ResponseWriter, r *http.Request, status int, data templateData[Tpage], templateName string, additionalTemplateFiles ...string) error {
	ts, err := h.sharedTemplates.Clone()
	if err != nil {
		return err
	}

	if len(additionalTemplateFiles) > 0 {
		ts, err = ts.ParseFS(h.templateFS, additionalTemplateFiles...)
		if err != nil {
			return err
		}
	}

	authed, _ := r.Context().Value(authedContextKey).(bool)

	buf := new(bytes.Buffer)
	// Execute the named template, passing in any dynamic data, and write to a buffer.
	data.Authed = authed
	err = ts.ExecuteTemplate(buf, templateName, data)
	if err != nil {
		return err
	}

	// https://www.alexedwards.net/blog/how-i-use-htmx-with-go#:~:text=finish%20this%20up.-,Setting%20the%20Vary%20header,-Because%20we%27re%20sending
	w.Header().Add("Vary", "HX-Request")
	w.WriteHeader(status)
	buf.WriteTo(w)

	return nil
}
