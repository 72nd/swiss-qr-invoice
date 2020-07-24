package swiss_qr_invoice

import (
	wrapper "github.com/72nd/gopdf-wrapper"
	"github.com/72nd/gopdf-wrapper/fonts"
	"github.com/72nd/swiss-qr-invoice/assets"
	"github.com/signintech/gopdf"
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

	if err := renderBasics(doc); err != nil {
		return nil, err
	}
	if err := receivePayment(doc, inv); err != nil {
		return nil, err
	}

	return doc, nil
}

func renderBasics(doc *wrapper.Doc) error {
	doc.AddLine(0, yTop, 210, yTop, 0.1, wrapper.SolidLine)
	doc.AddLine(62, yTop, 62, yBottom, 0.1, wrapper.SolidLine)
	scissors, err := assets.Scissors()
	if err != nil {
		return err
	}
	img, err := gopdf.ImageHolderByBytes(scissors)
	doc.ImageByHolder(img, 59.5, yTop + 10, nil)
	return err
}

func receivePayment(doc *wrapper.Doc, inv Invoice) error {
	doc.AddFormattedText(5, yTop + 5, "Empfangsschein", 11, "bold")
	doc.AddFormattedText(5, yTop + 12, "Konto / Zahlbar an", 6, "bold")
	return nil
}
