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
// example call: 	err := app.html.render(w, 200, nil, "base", "pages/home.tmpl")
func (h *htmlRenderer) render(w http.ResponseWriter, status int, data any, templateName string, additionalTemplateFiles ...string) error {
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

	buf := new(bytes.Buffer)

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
