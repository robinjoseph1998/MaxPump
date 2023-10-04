package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func SalesReportGenerator(date string, SalesCount int, TotalAmount float64, SalesByProducts map[string]int, PricePerQuantity map[string]float64) error {
	SalesReport := &SalesReportModel{
		Date:             time.Now(),
		TotalSales:       SalesCount,
		TotalAmount:      TotalAmount,
		ProductWiseSales: make(map[string]int),
		PricePerQuantity: make(map[string]float64),
	}
	for product, quantity := range SalesByProducts {
		SalesReport.ProductWiseSales[product] = quantity
	}
	for product, price := range PricePerQuantity {
		SalesReport.PricePerQuantity[product] = price
	}

	//creating a new pdf document and setting its oreantations
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	pdf.SetTextColor(0, 0, 255)
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "MAXPUMP", "", 0, "C", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(10)
	//adding cells to display date sales and total maount
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 7, fmt.Sprintf("Date: %s", date))
	pdf.Ln(10)
	pdf.Cell(0, 7, fmt.Sprintf("Total Sales: %d", SalesReport.TotalSales))
	pdf.Ln(10)
	pdf.Cell(0, 7, fmt.Sprintf("Total Amount: %.2f", SalesReport.TotalAmount))
	pdf.Ln(10)

	//adding the table heasder
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 10, "Date", "1", 0, "", false, 0, "")
	pdf.CellFormat(40, 10, "Product", "1", 0, "", false, 0, "")
	pdf.CellFormat(40, 10, "Quantity", "1", 0, "", false, 0, "")
	pdf.CellFormat(40, 10, "Price Per Quantity", "1", 0, "", false, 0, "")
	pdf.Ln(10)

	// set the font to bold and add a table header with columns "Date","Product" and "Quantity"
	pdf.SetFont("Arial", "", 12)
	dateLayout := "02 January 2006"
	for product, quantity := range SalesByProducts {
		pdf.CellFormat(40, 10, SalesReport.Date.Format(dateLayout), "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 10, product, "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%d units", quantity), "1", 0, "", false, 0, "")
		price, priceExists := PricePerQuantity[product]
		if priceExists {
			priceStr := fmt.Sprintf("%.2f", price)
			pdf.CellFormat(40, 10, priceStr, "1", 0, "", false, 0, "")
		} else {
			pdf.CellFormat(40, 10, "N/A", "1", 0, "", false, 0, "")
		}
		pdf.Ln(10)
	}
	err := pdf.OutputFileAndClose("Report.pdf")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PDF generated successfully.")
	return nil
}
