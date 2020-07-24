package swiss_qr_invoice

import (
	"fmt"

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
	if err := paymentInformation(doc, inv); err != nil {
		return nil, err
	}
	paymentAmount(doc, inv)
	receivingOffice(doc, inv)
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
	doc.ImageByHolder(img, 59.5, yTop+10, nil)
	return err
}

func paymentInformation(doc *wrapper.Doc, inv Invoice) error {
	doc.AddFormattedText(5, yTop+5, "Empfangsschein", 11, "bold")
	doc.AddFormattedText(5, yTop+12, "Konto / Zahlbar an", 6, "bold")

	yReceiverBase := yTop + 12 + doc.LineHeight(6)
	recCnt := 0.0
	if inv.ReceiverIBAN != "" {
		doc.AddSizedText(5, yReceiverBase, inv.ReceiverIBAN, 8)
		recCnt++
	}
	if inv.ReceiverName != "" {
		doc.AddSizedText(5, yReceiverBase+doc.LineHeight(8)*recCnt, inv.ReceiverName, 8)
		recCnt++
	}
	if inv.ReceiverAddress != "" {
		doc.AddSizedText(5, yReceiverBase+doc.LineHeight(8)*recCnt, inv.ReceiverAddress, 8)
		recCnt++
	}
	if inv.ReceiverZIPCode != "" || inv.ReceiverPlace != "" {
		doc.AddSizedText(5, yReceiverBase+doc.LineHeight(8)*recCnt, fmt.Sprintf("%s %s", inv.ReceiverZIPCode, inv.ReceiverPlace), 8)
		recCnt++
	}

	yReferenceBase := yReceiverBase + doc.LineHeight(8)*recCnt + doc.LineHeight(9)
	if inv.Reference != "" {
		doc.AddFormattedText(5, yReferenceBase, "Referenz", 6, "bold")
		doc.AddSizedText(5, yReferenceBase+doc.LineHeight(6), inv.Reference, 8)
	}

	yPayeeBase := yReferenceBase + doc.LineHeight(9) + doc.LineHeight(6) + doc.LineHeight(8)
	doc.AddFormattedText(5, yPayeeBase, "Zahlbar durch", 6, "bold")
	yPayeeBase += doc.LineHeight(6)
	if inv.Reference == "" {
		yPayeeBase -= doc.LineHeight(6) + doc.LineHeight(8)
	}
	if inv.noPayee() {
		doc.AddText(5, yPayeeBase+doc.LineHeight(6), "not implemented")
		return nil
	}
	pyeCnt := 0.0
	if inv.PayeeName != "" {
		doc.AddSizedText(5, yPayeeBase, inv.PayeeName, 8)
		pyeCnt++
	}
	if inv.PayeeAddress != "" {
		doc.AddSizedText(5, yPayeeBase+doc.LineHeight(8)*pyeCnt, inv.PayeeAddress, 8)
		pyeCnt++
	}
	if inv.PayeeZIPCode != "" || inv.PayeePlace != "" {
		doc.AddSizedText(5, yPayeeBase+doc.LineHeight(8)*pyeCnt, fmt.Sprintf("%s %s", inv.PayeeZIPCode, inv.PayeePlace), 8)
	}

	return nil
}

func paymentAmount(doc *wrapper.Doc, inv Invoice) {
	yAmountBase := yTop + 68
	doc.AddFormattedText(5, yAmountBase, "Währung", 6, "bold")
	doc.AddFormattedText(18, yAmountBase, "Betrag", 6, "bold")
	doc.AddSizedText(5, yAmountBase+doc.LineHeight(9), inv.Currency, 8)
	if inv.Amount != "" {
		doc.AddSizedText(18, yAmountBase+doc.LineHeight(9), inv.Amount, 8)
	} else {
		doc.AddText(18, yAmountBase+doc.LineHeight(9), "not implemented")
	}
}

func receivingOffice(doc *wrapper.Doc, inv Invoice) {
	yReceivingBase := yTop + 82
	text := "Annahmestelle"
	doc.AddFormattedText(40.5, yReceivingBase, text, 6, "bold")
}

