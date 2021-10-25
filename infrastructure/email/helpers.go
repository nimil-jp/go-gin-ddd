package email

import (
	"bytes"
	templateHtml "html/template"

	"github.com/pkg/errors"
	"jaytaylor.com/html2text"
)

func htmlToPlain(body Body) (string, error) {
	html, err := body.Html()
	if err != nil {
		return "", err
	}
	return html2text.FromString(html)
}

func setHtmlTemplate(template *templateHtml.Template, data interface{}) (string, error) {
	var out bytes.Buffer
	err := template.Execute(&out, data)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return out.String(), nil
}

// func setPlainTemplate(template *templateText.Template, data interface{}) (string, error) {
// 	var out bytes.Buffer
// 	err := template.Execute(&out, data)
// 	if err != nil {
// 		return "", errors.WithStack(err)
// 	}
// 	return out.String(), nil
// }
