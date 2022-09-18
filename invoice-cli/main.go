package main

import (
	"os"

	invoice "github.com/72nd/swiss-qr-invoice"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "invoice-cli",
		Usage: "create swiss QR-Code invoices as PDF",
		Action: func(c *cli.Context) error {
			_ = cli.ShowCommandHelp(c, c.Command.Name)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "generate PDF invoice based on `FILE`",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "input",
						Aliases: []string{"i"},
						Usage:   "path to input config file"},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "destination path for PDF invoice"},
				},
				Action: func(c *cli.Context) error {
					if c.String("input") == "" {
						logrus.Fatal("please specify the input file with the -i flag")
					}
					if c.String("output") == "" {
						logrus.Fatal("please specify the destination path for the PDF with the -o flag")
					}
					inv, err := invoice.OpenInvoice(c.String("input"))
					if err != nil {
						logrus.Fatal(err)
					}
					if err := inv.SaveAsPDF(c.String("output")); err != nil {
						logrus.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:  "new",
				Usage: "create new config file in `FILE`",
				Action: func(c *cli.Context) error {
					if c.Args().Len() != 1 {
						logrus.Fatal("please specify the output file as an argument")
					}
					inv, err := invoice.New(true)
					if err != nil {
						logrus.Fatal(err)
					}
					if err := inv.Save(c.Args().First()); err != nil {
						logrus.Fatal(err)
					}
					return nil
				},
			},
			{
				Name:  "qr-debug",
				Usage: "saves the content of the QR-Code as a text file for debugging",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "input",
						Aliases: []string{"i"},
						Usage:   "path to input config file"},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "destination path for PDF invoice"},
				},

				Action: func(c *cli.Context) error {
					if c.String("input") == "" {
						logrus.Fatal("please specify the input file with the -i flag")
					}
					if c.String("output") == "" {
						logrus.Fatal("please specify the destination path for the text file with the -o flag")
					}
					inv, err := invoice.OpenInvoice(c.String("input"))
					if err != nil {
						logrus.Fatal(err)
					}
					if err := inv.SaveQrConent(c.String("output")); err != nil {
						logrus.Fatal(err)
					}
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal("cli app error: ", err)
	}
}
