package email

import (
	"bytes"
	templateHtml "html/template"

	"jaytaylor.com/html2text"

	"github.com/nimil-jp/gin-utils/errors"
)

func htmlToPlain(body Body) (string, error) {
	html, err := body.HTML()
	if err != nil {
		return "", err
	}
	return html2text.FromString(html)
}

func setHTMLTemplate(template *templateHtml.Template, data interface{}) (string, error) {
	var out bytes.Buffer
	err := template.Execute(&out, data)
	if err != nil {
		return "", errors.NewUnexpected(err)
	}
	return out.String(), nil
}

// func setPlainTemplate(template *templateText.Template, data interface{}) (string, error) {
// 	var out bytes.Buffer
// 	err := template.Execute(&out, data)
// 	if err != nil {
// 		return "", errors.NewUnexpected(err)
// 	}
// 	return out.String(), nil
// }
