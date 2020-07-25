package swiss_qr_invoice

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
	if err := ioutil.WriteFile(path, raw, 0644); err != nil {
		return err
	}
	return nil
}

// SaveAsPDF generates the invoice and save it as a PDF.
func (i Invoice) SaveAsPDF(path string) error {
	doc, err := getDoc(i)
	if err != nil {
		return err
	}
	return doc.WritePdf(path)
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
