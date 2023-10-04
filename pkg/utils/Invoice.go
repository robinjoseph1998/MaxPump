package utils

import (
	"MAXPUMP1/pkg/domain/entity"
	"fmt"
	"log"

	"github.com/jung-kurt/gofpdf"

	"github.com/google/uuid"
)

func generateInvoiceNumber() string {
	uuid := uuid.New()
	return fmt.Sprintf("INV-%s", uuid)
}

func InvoiceGenerator(Order *entity.Order, Address *entity.Address, User *entity.User) error {
	var invoiceData *InvoiceModel
	FirstName := User.FirstName
	LastName := User.LastName
	FullName := FirstName + LastName

	invoiceData = &InvoiceModel{
		InvoiceNumber:  generateInvoiceNumber(),
		Date:           Order.Date_Of_Delivered,
		BillingName:    FullName,
		BillingAddress: Address.HouseName,
		District:       Address.District,
		Pincode:        Address.Pincode,
		Landmark:       Address.Landmark,
		TotalPrice:     Order.TotalPrice,
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetTextColor(0, 0, 255)
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "MAXPUMP", "", 0, "C", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 10, fmt.Sprintf("Invoice Number: %s", invoiceData.InvoiceNumber))
	pdf.Ln(10)
	pdf.Cell(0, 10, fmt.Sprintf("Date: %s", invoiceData.Date))
	pdf.Ln(10)
	pdf.Cell(0, 10, fmt.Sprintf("Billing Name: %s", invoiceData.BillingName))
	pdf.Ln(10)
	pdf.Cell(0, 10, fmt.Sprintf("Billing Address: %s", invoiceData.BillingAddress))
	pdf.Ln(10)
	pdf.Cell(0, 10, fmt.Sprintf("pincode: %s", invoiceData.Pincode))
	pdf.Ln(10)
	pdf.Cell(0, 10, fmt.Sprintf("Landmark: %s", invoiceData.Landmark))
	pdf.Ln(10)
	pdf.Cell(0, 10, fmt.Sprintf("District: %s", invoiceData.District))
	pdf.Ln(10)
	pdf.Cell(0, 10, fmt.Sprintf("Total Price: %.2f", invoiceData.TotalPrice))
	err := pdf.OutputFileAndClose("Invoice.pdf")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PDF generated successfully.")
	return nil
}
