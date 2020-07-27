// Package swissqrinvoice  Generate Swiss QR Invoices as described in [this standard](https://www.paymentstandards.ch/dam/downloads/ig-qr-bill-de.pdf) and [the style guide](https://www.paymentstandards.ch/dam/downloads/style-guide-de.pdf). The library uses [gopdf](https://github.com/signintech/gopdf) via the [gopdf-wrapper](https://github.com/72nd/gopdf-wrapper).
package swissqrinvoice

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"

	wrapper "github.com/72nd/gopdf-wrapper"
	"github.com/creasty/defaults"
	"github.com/signintech/gopdf"
	"gopkg.in/yaml.v2"
)

// Invoice contains all necessary information for the generation of an invoice.
type Invoice struct {
	ReceiverIBAN    string `yaml:"receiver_iban" default:"CH44 3199 9123 0008 8901 2"`
	IsQrIBAN        bool   `yaml:"is_qr_iban" default:"true"`
	ReceiverName    string `yaml:"receiver_name" default:"Robert Schneider AG"`
	ReceiverStreet  string `yaml:"receiver_street" default:"Rue du Lac"`
	ReceiverNumber  string `yaml:"receiver_number" default:"12"`
	ReceiverZIPCode string `yaml:"receiver_zip_code" default:"2501"`
	ReceiverPlace   string `yaml:"receiver_place" default:"Biel"`
	ReceiverCountry string `yaml:"receiver_country" default:"CH"`
	PayeeName       string `yaml:"payee_name" default:"Pia-Maria Rutschmann-Schnyder"`
	PayeeStreet     string `yaml:"payee_street" default:"Grosse Marktgasse"`
	PayeeNumber     string `yaml:"payee_number" default:"28"`
	PayeeZIPCode    string `yaml:"payee_zip_code" default:"9400"`
	PayeePlace      string `yaml:"payee_place" default:"Rorschach"`
	PayeeCountry    string `yaml:"payee_country" default:"CH"`
	Reference       string `yaml:"reference" default:"21 00000 00003 13947 14300 09017"`
	AdditionalInfo  string `yaml:"additional_info" default:"Rechnung Nr. 3139 für Gartenarbeiten"`
	Amount          string `yaml:"amount" default:"3 949.75"`
	Currency        string `yaml:"currency" default:"CHF"`
}

// New returns a new invoice optional with the default values.
func New(useDefaults bool) (*Invoice, error) {
	rsl := &Invoice{}
	if useDefaults {
		if err := defaults.Set(rsl); err != nil {
			return nil, err
		}
	}
	return rsl, nil
}

// OpenInvoice opens a invoice config YAML file with the given path.
func OpenInvoice(path string) (*Invoice, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var inv Invoice
	if err := yaml.Unmarshal(raw, &inv); err != nil {
		return nil, err
	}
	return &inv, nil
}

// Save saves the invoice as a YAML file.
func (i Invoice) Save(path string) error {
	raw, err := yaml.Marshal(i)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, raw, 0644)
}

// SaveAsPDF generates the invoice and save it as a PDF.
func (i Invoice) SaveAsPDF(path string) error {
	doc, err := getDoc(i)
	if err != nil {
		return err
	}
	return doc.WritePdf(path)
}

// SaveQrConent saves the content of the QR code to a text file for debugging.
func (i Invoice) SaveQrConent(path string) error {
	raw, err := i.qrContent()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, []byte(raw), 0644)
}

// GoPdf returns the invoice as a gopdf.GoPdf element. This can be used to further
// customizing the invoice.
func (i Invoice) GoPdf() (*gopdf.GoPdf, error) {
	doc, err := getDoc(i)
	if err != nil {
		return nil, err
	}
	return &doc.GoPdf, nil
}

// Doc returns the invoice as a gopdf_wrapper.Doc element. This can be used to further
// customizing the invoice. in contrary to GoPdf() the wrapper doc has more convince
// functions than using gopdf directly.
func (i Invoice) Doc() (*wrapper.Doc, error) {
	doc, err := getDoc(i)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// noPayee returns true if no fields of the payee are set
func (i Invoice) noPayee() bool {
	return i.PayeeName == "" && i.PayeeStreet == "" && i.PayeeZIPCode == "" && i.PayeePlace == ""
}

func (i Invoice) qrContent() (string, error) {
	qrTpl := `SPC\r
0200\r
1\r
{{ .iban }}\r
S\r
{{ .inv.ReceiverName }}\r
{{ .inv.ReceiverStreet }}\r
{{ .inv.ReceiverNumber }}\r
{{ .inv.ReceiverZIPCode }}\r
{{ .inv.ReceiverPlace }}\r
{{ .inv.ReceiverCountry }}\r
\r
\r
\r
\r
\r
\r
\r
{{ .amount }}\r
{{ .inv.Currency }}\r
{{ .payeeAdrType }}\r
{{ .inv.PayeeName }}\r
{{ .inv.PayeeStreet }}\r
{{ .inv.PayeeNumber }}\r
{{ .inv.PayeeZIPCode }}\r
{{ .inv.PayeePlace }}\r
{{ .inv.PayeeCountry }}\r
{{ .refType }}\r
{{ .reference }}\r
{{ .inv.AdditionalInfo }}\r
EPD
`
	refType := "QRR"
	if !i.IsQrIBAN {
		refType = "SCOR"
	}
	payeeAdrType := "S"
	if i.noPayee() {
		payeeAdrType = ""
	}
	data := map[string]interface{}{
		"inv":          i,
		"iban":         strings.Replace(i.ReceiverIBAN, " ", "", -1),
		"amount":       strings.Replace(i.Amount, " ", "", -1),
		"payeeAdrType": payeeAdrType,
		"reference":    strings.Replace(i.Reference, " ", "", -1),
		"refType":      refType,
	}
	tpl, err := template.New("qr-content").Parse(qrTpl)
	if err != nil {
		return "", fmt.Errorf("error while creating qr template: %s", err)
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error while applying data to qr template: %s", err)
	}
	rsl := strings.Replace(buf.String(), "\\r", "\r", -1)
	return rsl, nil
}
