package web

import (
	"net/http"
	"net/url"
)

func isHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func redirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	if isHTMXRequest(r) {
		w.Header().Set("HX-Redirect", url)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	http.Redirect(w, r, url, code)
}

func browserURL(r *http.Request) (*url.URL, error) {
	cu := r.Header.Get("HX-Current-URL")
	if cu != "" {
		return url.Parse(cu)
	}

	return r.URL, nil
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	err := app.html.render(w, 200, nil, "base", "pages/home.tmpl")
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func (app *application) gopher(w http.ResponseWriter, r *http.Request) {
	width := 100
	err := app.html.render(w, http.StatusOK, width, "partial:image:gopher")
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}
