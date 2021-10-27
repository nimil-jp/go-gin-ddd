package email

import (
	_ "embed"
	"html/template"
)

//go:embed template/default.gohtml
var defaultHTMLTemplateString string
var defaultHTMLTemplate *template.Template

func init() {
	var err error
	defaultHTMLTemplate, err = template.New("defaultHtml").Parse(defaultHTMLTemplateString)
	if err != nil {
		panic(err)
	}
}

type DefaultBody struct {
	Title string
	Body  string
}

func (d DefaultBody) Subject() string {
	return d.Title
}

func (d DefaultBody) HTML() (string, error) {
	return setHTMLTemplate(defaultHTMLTemplate, d)
}

func (d DefaultBody) Plain() (string, error) {
	return htmlToPlain(d)
}
