package main

import (
	"html/template"
	"net/http"
	"fmt"
	"bytes"

)

var templates = template.Must(template.New("t").ParseGlob("templates/**/*.html"))

var errorTemplate = `
<html>
 <body>
     <h1>Error rendering template %s</h1>
     <p>%s</p>
 </body>
</html>
`
var layoutFunc = template.FuncMap{
	"ground" : func() (string, error) {
		return  "", fmt.Errorf("ground called inapropriately")
	},
	"yield" : func() (string, error) {
	    return  "", fmt.Errorf("yield called inapropriately")
	},

}

var layout = template.Must(template.New("layout.html").Funcs(layoutFunc).ParseFiles("templates/layout.html"), )


func FlashMessages(m string) string {
	var s = "FlashSuccess"
	var  i = "FlashInfo"
	var w = "FlashWarning"
	var d = "FlashDanger"
	switch msg := m; msg {
	case "User created":
		return s
	case "Signed in":
		return s
	case "User info":
		return i
	case "User warning":
		return w
	case "User Danger":
		return d
	default:
		return ""
	}
    return ""
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	data["CurrentUser"] = RequestUser(r)
	//
	flash := r.URL.Query().Get("flash")
	flashTemplate := FlashMessages(flash)
	if flashTemplate != "" {
		data[flashTemplate] = r.URL.Query().Get("flash")
	}
funcs := template.FuncMap{
	"yield": func() (template.HTML, error) {
		buf := bytes.NewBuffer(nil)
		err := templates.ExecuteTemplate(buf, name, data)
		return template.HTML(buf.String()), err
	},
	"ground": func() (template.HTML, error) {
		buf := bytes.NewBuffer(nil)
		err := templates.ExecuteTemplate(buf, name, data)
		return template.HTML(buf.String()), err
	},

}


layoutClone, _ := layout.Clone()
template.Must(layoutClone.New("ground.html").Funcs(funcs).ParseFiles("templates/ground.html"))

layoutClone.Funcs(funcs)
err := layoutClone.Execute(w,data)
if err != nil {
	http.Error(
		w,
			fmt.Sprintf(errorTemplate,name, err),
				http.StatusInternalServerError,
	)
}
}
