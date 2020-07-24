package swiss_qr_invoice

import (
	wrapper "github.com/72nd/gopdf-wrapper"
	"github.com/72nd/gopdf-wrapper/fonts"
)

const (
	yTop    = 192.0
	yBottom = 297.0
)

func getDoc(inv Invoice) (*wrapper.Doc, error) {
	doc, err := wrapper.NewDoc(12, 1)
	if err != nil {
		return nil, err
	}
	liberation, err := fonts.NewLiberationSansFamily()
	if err != nil {
		return nil, err
	}
	doc.SetFontFamily(*liberation)
	doc.AddPage()

	renderBasics(doc)

	return doc, nil
}

func renderBasics(doc *wrapper.Doc) {
	doc.AddLine(0, yTop, 210, yTop, 0.1, wrapper.SolidLine)
	doc.AddLine(62, yTop, 62, yBottom, 0.1, wrapper.SolidLine)
}
