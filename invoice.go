package swiss_qr_invoice

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/creasty/defaults"
)

// Invoice contains all necessary information for the generation of an invoice.
type Invoice struct {
	ReceiverIBAN    string `yaml:"receiver_iban" default:"CH44 3199 9123 0008 8901 2"`
	IsQrIBAN bool `yaml:"is_qr_iban" default:"true"`
	ReceiverName    string `yaml:"receiver_name" default:"Robert Schneider AG"`
	ReceiverStreet  string `yaml:"receiver_street" default:"Rue du Lac"`
	ReceiverNumber  string `yaml:"receiver_number" default:"12"`
	ReceiverZIPCode string `yaml:"receiver_zip_code" default:"2501"`
	ReceiverPlace   string `yaml:"receiver_place" default:"Biel"`
	ReceiverCountry string `yaml:"receiver_country" default:"CH"`
	PayeeName       string `yaml:"receiver_name" default:"Pia-Maria Rutschmann-Schnyder"`
	PayeeStreet     string `yaml:"receiver_street" default:"Grosse Marktgasse"`
	PayeeNumber     string `yaml:"receiver_number" default:"28"`
	PayeeZIPCode    string `yaml:"receiver_zip_code" default:"9400"`
	PayeePlace      string `yaml:"receiver_place" default:"Rorschach"`
	PayeeCountry    string `yaml:"payee_country" default:"CH"`
	Reference       string `yaml:"reference" default:"21 00000 00003 13947 14300 09017"`
	AdditionalInfo  string `yaml:"additional_info" default:"Rechnung Nr. 3139 für Gartenarbeiten"`
	Amount          string `yaml:"amount" default:"3949.75"`
	Currency        string `yaml:"currency" default:"CHF"`
}

// NewInvoice returns a new invoice optional with the default values.
func NewInvoice(useDefaults bool) (*Invoice, error) {
	rsl := &Invoice{}
	if useDefaults {
		if err := defaults.Set(rsl); err != nil {
			return nil, err
		}
	}
	return rsl, nil
}

// SaveAsPDF generates the invoice and save it as a PDF.
func (i Invoice) SaveAsPDF(path string) error {
	return nil
}

// noPayee returns true if no fields of the payee are set
func (i Invoice) noPayee() bool {
	return i.PayeeName == "" && i.PayeeStreet == "" && i.PayeeZIPCode == "" && i.PayeePlace == ""
}

func (i Invoice) qrContent() (string, error) {
	qrTpl := `
SPC
0200
1
{{ .iban }}
S
{{ .inv.ReceiverName }}
{{ .inv.ReceiverStreet }}
{{ .inv.ReceiverNumber }}
{{ .inv.ReceiverZIPCode }}
{{ .inv.ReceiverPlace }}
{{ .inv.ReceiverCountry }}







{{ .inv.Amount }}
{{ .inv.Currency }}
S
{{ .inv.PayeeName }}
{{ .inv.PayeeStreet }}
{{ .inv.PayeeNumber }}
{{ .inv.PayeeZIPCode }}
{{ .inv.PayeePlace }}
{{ .inv.PayeeCountry }}
{{ .refType }}
{{ .inv.Reference }}
{{ .inv.AdditionalInfo }}
EPD
`
	refType := "QRR"
	if !i.IsQrIBAN {
		refType = "SCOR"
	}
	data := map[string]interface{}{
		"inv":     i,
		"iban":    strings.Replace(i.ReceiverIBAN, " ", "", -1),
		"refType": refType,
	}
	tpl, err := template.New("qr-content").Parse(qrTpl)
	if err != nil {
		return "", fmt.Errorf("error while creating qr template: %s", err)
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error while applying data to qr template: %s", err)
	}
	return buf.String(), nil
}
