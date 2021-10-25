package email

import (
	_ "embed"
	"html/template"
)

//go:embed template/default.gohtml
var defaultHtmlTemplateString string
var defaultHtmlTemplate *template.Template

func init() {
	var err error
	defaultHtmlTemplate, err = template.New("defaultHtml").Parse(defaultHtmlTemplateString)
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

func (d DefaultBody) Html() (string, error) {
	return setHtmlTemplate(defaultHtmlTemplate, d)
}

func (d DefaultBody) Plain() (string, error) {
	return htmlToPlain(d)
}
