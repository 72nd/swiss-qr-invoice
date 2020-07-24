package swiss_qr_invoice

type Invoice struct {
}

func NewInvoice(useDefaults bool) Invoice {
	return Invoice{}
}

func (i Invoice) SaveAsPDF(path string) error {
	return nil
}
