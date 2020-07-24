package swiss_qr_invoice

import "github.com/creasty/defaults"

// Invoice contains all necessary information for the generation of an invoice.
type Invoice struct {
	ReceiverIBAN    string `yaml:"receiver_iban" default:"CH44 3199 9123 0008 8901 2"`
	ReceiverName    string `yaml:"receiver_name" default:"Robert Schneider AG"`
	ReceiverAddress string `yaml:"receiver_address" default:"Rue du Lac 1268"`
	ReceiverZIPCode string `yaml:"receiver_zip_code" default:"2501"`
	ReceiverPlace   string `yaml:"receiver_place" default:"Biel"`
	PayeeName       string `yaml:"receiver_name" default:"Pia-Maria Rutschmann-Schnyder"`
	PayeeAddress    string `yaml:"receiver_address" default:"Grosse Marktgasse 28"`
	PayeeZIPCode    string `yaml:"receiver_zip_code" default:"9400"`
	PayeePlace      string `yaml:"receiver_place" default:"Rorschach"`
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

