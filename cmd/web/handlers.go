package web

import (
	"net/http"
	"net/url"
	"time"
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
	err := render(app.html, w, r, 200, noTemplateData, "base", "pages/home.tmpl")
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func (app *application) gopher(w http.ResponseWriter, r *http.Request) {
	width := 100
	err := render(app.html, w, r, http.StatusOK, templateData[struct{ Width int }]{Page: struct{ Width int }{Width: width}}, "partial:image:gopher")
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	password := r.FormValue("password")

	if password != app.password {
		app.logger.Info("login failed", "password", password)
		redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	c := cookieTemplate
	token, err := newSessionToken(app.jwtSecret)
	c.Value = token
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, http.StatusText(500), 500)
		return
	}
	c.Expires = time.Now().Add(24 * time.Hour)
	http.SetCookie(w, &c)

}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	c := cookieTemplate
	c.Value = ""
	// Expire in the past
	c.Expires = time.Now().Add(-1 * time.Hour)
	http.SetCookie(w, &c)
	redirect(w, r, "/", http.StatusSeeOther)
}
