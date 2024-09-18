package template

import (
	"bytes"
	_ "embed"
	"github.com/pkg/errors"
	templatePkg "github.com/yanodincov/k8s-forwarder/pkg/template"
	"html/template"
)

//go:embed select.gohtml
var selectTemplateString string
var selectTemplate *template.Template

type SelectTemplateData struct {
	Header   string
	Question string
	Filter   string
	Footer   string

	Options []SelectTemplateOption

	TotalPages  int
	CurrentPage int
	PageSymbol  string

	HoverSymbol   string
	AccentColor   string
	ActiveColor   string
	InactiveColor string
}

type SelectTemplateOption struct {
	Text       string
	IsSelected bool
	IsHovered  bool
	IsSpecial  bool
}

func RenderSelectTemplate(data SelectTemplateData) (string, error) {
	w := bytes.NewBuffer(nil)
	if err := selectTemplate.Execute(w, data); err != nil {
		return "", errors.Wrap(err, "render select template")
	}

	return w.String(), nil
}

func init() {
	selectTemplate = template.Must(template.
		New("select").
		Funcs(templatePkg.GetCliTemplateFns()).
		Parse(selectTemplateString),
	)
}
