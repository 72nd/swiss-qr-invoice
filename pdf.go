package swissqrinvoice

import (
	"fmt"

	wrapper "github.com/72nd/gopdf-wrapper"
	"github.com/72nd/gopdf-wrapper/fonts"
	"github.com/72nd/swiss-qr-invoice/assets"
	"github.com/signintech/gopdf"
	"github.com/skip2/go-qrcode"
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

	// Recycling Part
	if err := renderBasics(doc); err != nil {
		return nil, err
	}
	if err := receivingInformation(doc, inv); err != nil {
		return nil, err
	}
	receivingAmount(doc, inv)
	receivingOffice(doc, inv)

	// Payment Part
	if err := paymentBasics(doc, inv); err != nil {
		return nil, err
	}
	paymentAmount(doc, inv)
	if err := paymentInformation(doc, inv); err != nil {
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
	doc.ImageByHolder(img, 60.25, yTop+10, &gopdf.Rect{W: 3.5, H: 5.8})
	return err
}

func receivingInformation(doc *wrapper.Doc, inv Invoice) error {
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
	if inv.ReceiverStreet != "" {
		address := fmt.Sprintf("%s %s", inv.ReceiverStreet, inv.ReceiverNumber)
		doc.AddSizedText(5, yReceiverBase+doc.LineHeight(8)*recCnt, address, 8)
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
	if inv.Reference == "" {
		yPayeeBase -= doc.LineHeight(6) + doc.LineHeight(8)
	}
	doc.AddFormattedText(5, yPayeeBase, "Zahlbar durch", 6, "bold")
	yPayeeBase += doc.LineHeight(8)
	if inv.noPayee() {
		emptyFields(doc, 5, yPayeeBase, 57, yPayeeBase+20)
		return nil
	}
	pyeCnt := 0.0
	if inv.PayeeName != "" {
		doc.AddSizedText(5, yPayeeBase, inv.PayeeName, 8)
		pyeCnt++
	}
	if inv.PayeeStreet != "" {
		address := fmt.Sprintf("%s %s", inv.PayeeStreet, inv.PayeeNumber)
		doc.AddSizedText(5, yPayeeBase+doc.LineHeight(8)*pyeCnt, address, 8)
		pyeCnt++
	}
	if inv.PayeeZIPCode != "" || inv.PayeePlace != "" {
		doc.AddSizedText(5, yPayeeBase+doc.LineHeight(8)*pyeCnt, fmt.Sprintf("%s %s", inv.PayeeZIPCode, inv.PayeePlace), 8)
	}

	return nil
}

func receivingAmount(doc *wrapper.Doc, inv Invoice) {
	yAmountBase := yTop + 68
	doc.AddFormattedText(5, yAmountBase, "Währung", 6, "bold")
	doc.AddFormattedText(18, yAmountBase, "Betrag", 6, "bold")
	doc.AddSizedText(5, yAmountBase+doc.LineHeight(9), inv.Currency, 8)
	if inv.Amount != "" {
		doc.AddSizedText(18, yAmountBase+doc.LineHeight(9), inv.Amount, 8)
	} else {
		emptyFields(doc, 27, yAmountBase, 27+30, yAmountBase+10)
	}
}

func receivingOffice(doc *wrapper.Doc, inv Invoice) {
	yReceivingBase := yTop + 82
	text := "Annahmestelle"
	doc.AddFormattedText(40.5, yReceivingBase, text, 6, "bold")
}

func paymentBasics(doc *wrapper.Doc, inv Invoice) error {
	doc.AddFormattedText(67, yTop+5, "Zahlteil", 11, "bold")

	content, err := inv.qrContent()
	if err != nil {
		return err
	}
	qr, err := qrcode.New(content, qrcode.Medium)
	if err != nil {
		return err
	}
	qr.DisableBorder = true
	png, err := qr.PNG(512)
	if err != nil {
		return err
	}
	qrImg, err := gopdf.ImageHolderByBytes(png)
	if err != nil {
		return err
	}
	doc.ImageByHolder(qrImg, 67, yTop+17, &gopdf.Rect{W: 46, H: 46})

	cross, err := assets.CHCross()
	if err != nil {
		return err
	}
	crossImg, err := gopdf.ImageHolderByBytes(cross)
	if err != nil {
		return err
	}
	doc.ImageByHolder(crossImg, 67+19.4, yTop+17+19.4, &gopdf.Rect{W: 7.28, H: 7.28})

	return nil
}

func paymentAmount(doc *wrapper.Doc, inv Invoice) {
	yAmountBase := yTop + 68
	doc.AddFormattedText(67, yAmountBase, "Währung", 8, "bold")
	doc.AddFormattedText(83, yAmountBase, "Betrag", 8, "bold")
	doc.AddSizedText(67, yAmountBase+doc.LineHeight(13), inv.Currency, 10)
	if inv.Amount != "" {
		doc.AddSizedText(83, yAmountBase+doc.LineHeight(13), inv.Amount, 10)
	} else {
		emptyFields(doc, 77, yAmountBase+doc.LineHeight(13), 77+40, yAmountBase+doc.LineHeight(13)+15)
	}
}

func paymentInformation(doc *wrapper.Doc, inv Invoice) error {
	doc.AddFormattedText(118, yTop+12, "Konto / Zahlbar an", 8, "bold")

	yReceiverBase := yTop + 12 + doc.LineHeight(8)
	recCnt := 0.0
	if inv.ReceiverIBAN != "" {
		doc.AddSizedText(118, yReceiverBase, inv.ReceiverIBAN, 10)
		recCnt++
	}
	if inv.ReceiverName != "" {
		doc.AddSizedText(118, yReceiverBase+doc.LineHeight(10)*recCnt, inv.ReceiverName, 10)
		recCnt++
	}
	if inv.ReceiverStreet != "" {
		address := fmt.Sprintf("%s %s", inv.ReceiverStreet, inv.ReceiverNumber)
		doc.AddSizedText(118, yReceiverBase+doc.LineHeight(10)*recCnt, address, 10)
		recCnt++
	}
	if inv.ReceiverZIPCode != "" || inv.ReceiverPlace != "" {
		doc.AddSizedText(118, yReceiverBase+doc.LineHeight(10)*recCnt, fmt.Sprintf("%s %s", inv.ReceiverZIPCode, inv.ReceiverPlace), 10)
		recCnt++
	}

	yReferenceBase := yReceiverBase + doc.LineHeight(10)*recCnt + doc.LineHeight(11)
	if inv.Reference != "" {
		doc.AddFormattedText(118, yReferenceBase, "Referenz", 8, "bold")
		doc.AddSizedText(118, yReferenceBase+doc.LineHeight(8), inv.Reference, 10)
	}

	yAdditionalBase := yReferenceBase + doc.LineHeight(10) + doc.LineHeight(8) + doc.LineHeight(11)
	if inv.Reference == "" {
		yAdditionalBase -= doc.LineHeight(8) + doc.LineHeight(10)
	}
	if inv.AdditionalInfo != "" {
		doc.AddFormattedText(118, yAdditionalBase, "Zusätzliche Informationen", 8, "bold")
		doc.AddSizedText(118, yAdditionalBase+doc.LineHeight(8), inv.AdditionalInfo, 10)
	}

	yPayeeBase := yAdditionalBase + doc.LineHeight(10) + doc.LineHeight(8) + doc.LineHeight(11)
	if inv.AdditionalInfo == "" {
		yPayeeBase -= doc.LineHeight(8) + doc.LineHeight(10)
	}
	doc.AddFormattedText(118, yPayeeBase, "Zahlbar durch", 8, "bold")
	yPayeeBase += doc.LineHeight(8)
	if inv.noPayee() {
		emptyFields(doc, 118, yPayeeBase+doc.LineHeight(8), 118+65, yPayeeBase+doc.LineHeight(8)+25)
		return nil
	}
	pyeCnt := 0.0
	if inv.PayeeName != "" {
		doc.AddSizedText(118, yPayeeBase, inv.PayeeName, 10)
		pyeCnt++
	}
	if inv.PayeeStreet != "" {
		address := fmt.Sprintf("%s %s", inv.PayeeStreet, inv.PayeeNumber)
		doc.AddSizedText(118, yPayeeBase+doc.LineHeight(10)*pyeCnt, address, 10)
		pyeCnt++
	}
	if inv.PayeeZIPCode != "" || inv.PayeePlace != "" {
		doc.AddSizedText(118, yPayeeBase+doc.LineHeight(10)*pyeCnt, fmt.Sprintf("%s %s", inv.PayeeZIPCode, inv.PayeePlace), 10)
	}

	return nil
}

func emptyFields(doc *wrapper.Doc, x1, y1, x2, y2 float64) {
	doc.AddLine(x1, y1+0.125, x1+3, y1+0.125, 0.25, wrapper.SolidLine)
	doc.AddLine(x1+0.125, y1, x1+0.125, y1+3, 0.25, wrapper.SolidLine)
	doc.AddLine(x2, y1+0.125, x2-3, y1+0.125, 0.25, wrapper.SolidLine)
	doc.AddLine(x2-0.125, y1, x2-0.125, y1+3, 0.25, wrapper.SolidLine)
	doc.AddLine(x1, y2-0.125, x1+3, y2-0.125, 0.25, wrapper.SolidLine)
	doc.AddLine(x1+0.125, y2, x1+0.125, y2-3, 0.25, wrapper.SolidLine)
	doc.AddLine(x2, y2-0.125, x2-3, y2-0.125, 0.25, wrapper.SolidLine)
	doc.AddLine(x2-0.125, y2, x2-0.125, y2-3, 0.25, wrapper.SolidLine)
}
