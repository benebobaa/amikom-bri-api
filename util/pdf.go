package util

import (
	"fmt"
	"github.com/benebobaa/amikom-bri-api/domain/entity"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"log"
	"time"
)

type PDFGenerator struct {
	PDF *gofpdf.Fpdf
}

func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{
		PDF: gofpdf.New("P", "mm", "A4", ""),
	}
}

func (p *PDFGenerator) GeneratePdf(entries []entity.Entry, totalIn, totalOut int64) (string, error) {
	p.PDF.AddPage()

	// Set font
	p.PDF.SetFont("Arial", "", 12)

	// Title
	title := "Transaction History"
	p.PDF.SetFont("Arial", "B", 18)
	_, lineHt := p.PDF.GetFontSize()
	p.PDF.CellFormat(0, lineHt, title, "", 1, "C", false, 0, "")

	//Add padding top
	p.PDF.Ln(10)

	// User information
	userInfo := "Generated by: AmikomPedia Financials\nDate: " + time.Now().Format("2006-01-02 15:04:05")
	p.PDF.SetFont("Arial", "", 12)
	p.PDF.MultiCell(0, lineHt, userInfo, "", "L", false)

	//Add padding top
	p.PDF.Ln(10)

	// Account information
	accountInfo := fmt.Sprintf("Username: %s\nEmail: %s\nTotal Income: %d\nTotal Expenses: %d", entries[0].Account.User.Username, entries[0].Account.User.Email, totalIn, totalOut)
	p.PDF.SetFont("Arial", "", 12)
	p.PDF.MultiCell(0, lineHt, accountInfo, "", "L", false)

	//Add padding top
	p.PDF.Ln(10)

	// Add data table headers
	headers := []string{"No", "Account ID", "Amount", "Type", "Date"}
	p.PDF.SetFont("Arial", "", 14)
	for _, header := range headers {
		if header == "No" {
			p.PDF.CellFormat(20, 10, header, "1", 0, "C", false, 0, "")
			continue
		}
		p.PDF.CellFormat(40, 10, header, "1", 0, "C", false, 0, "")
	}

	// Convert data from struct to string
	var data [][]string

	for i, entry := range entries {
		data = append(data, []string{
			fmt.Sprintf("%d", i+1),
			fmt.Sprintf("%d", entry.ID),
			fmt.Sprintf("%d", entry.Amount),
			entry.EntryType,
			entry.Date.Format("2006-01-02"),
		})
	}

	p.PDF.Ln(-1)
	p.PDF.SetFont("Arial", "", 10)
	for _, row := range data {
		for i, col := range row {
			if i == 0 {
				p.PDF.CellFormat(20, 10, col, "1", 0, "C", false, 0, "")
				continue
			}
			p.PDF.CellFormat(40, 10, col, "1", 0, "C", false, 0, "")
		}
		p.PDF.Ln(-1)
	}

	fileName := fmt.Sprintf("public/transaction_pdf/%s-%s.transaction_pdf", entries[0].Account.User.Email, uuid.NewString())

	// Save the PDF file
	err := p.PDF.OutputFileAndClose(fileName)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	fmt.Println("PDF generated successfully.")
	return fileName, nil
}
