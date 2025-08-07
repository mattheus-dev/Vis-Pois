package formatter

import (
	"bytes"
	"fmt"
	"strings"
	"vis-pois/internal/domain/entities"
)

type TableFormatter struct {
	columns []column
}

type column struct {
	header string
	width  int
	getter func(record entities.Record) string
}

func NewTableFormatter() *TableFormatter {
	return &TableFormatter{
		columns: []column{
			{header: "Index", width: 5, getter: func(record entities.Record) string { return record.ID }},
			{header: "Name", width: 20, getter: func(record entities.Record) string { return record.Name }},
			{header: "Description", width: 30, getter: func(record entities.Record) string { return record.Description }},
			{header: "Brand", width: 15, getter: func(record entities.Record) string { return record.Brand }},
			{header: "Category", width: 15, getter: func(record entities.Record) string { return record.Category }},
			{header: "Price", width: 10, getter: func(record entities.Record) string { return fmt.Sprintf("%.2f", record.Price) }},
			{header: "Currency", width: 10, getter: func(record entities.Record) string { return "BRL" }},
			{header: "Stock", width: 10, getter: func(record entities.Record) string { return fmt.Sprintf("%d", record.Stock) }},
			{header: "EAN", width: 15, getter: func(record entities.Record) string { return record.Barcode }},
			{header: "Color", width: 10, getter: func(record entities.Record) string { return record.Color }},
			{header: "Size", width: 10, getter: func(record entities.Record) string { 
				if record.Dimensions != "" {
					return record.Dimensions 
				}
				return "N/A"
			}},
			{header: "Availability", width: 15, getter: func(record entities.Record) string { 
				if record.Stock > 0 {
					return "Disponível"
				}
				return "Indisponível"
			}},
			{header: "InternalID", width: 15, getter: func(record entities.Record) string { return record.SKU }},
		},
	}
}

func (t *TableFormatter) FormatAsTable(records []entities.Record) string {
	var buf bytes.Buffer
	
	// Escrever linha de cabeçalho
	t.writeHeaderRow(&buf)
	
	// Escrever linha de separação
	t.writeSeparatorRow(&buf)
	
	// Escrever linhas de dados
	for i, record := range records {
		t.writeDataRow(&buf, record, i+1)
	}
	
	return buf.String()
}

func (t *TableFormatter) writeHeaderRow(buf *bytes.Buffer) {
	for i, col := range t.columns {
		if i > 0 {
			buf.WriteString(" | ")
		}
		format := fmt.Sprintf("%%-%ds", col.width)
		buf.WriteString(fmt.Sprintf(format, truncateString(col.header, col.width)))
	}
	buf.WriteString("\n")
}

func (t *TableFormatter) writeSeparatorRow(buf *bytes.Buffer) {
	for i, col := range t.columns {
		if i > 0 {
			buf.WriteString("-+-")
		}
		buf.WriteString(strings.Repeat("-", col.width))
	}
	buf.WriteString("\n")
}

func (t *TableFormatter) writeDataRow(buf *bytes.Buffer, record entities.Record, index int) {
	for i, col := range t.columns {
		if i > 0 {
			buf.WriteString(" | ")
		}
		
		var value string
		if i == 0 {
			// Index column shows the row number
			value = fmt.Sprintf("%d", index)
		} else {
			value = col.getter(record)
		}
		
		format := fmt.Sprintf("%%-%ds", col.width)
		buf.WriteString(fmt.Sprintf(format, truncateString(value, col.width)))
	}
	buf.WriteString("\n")
}

func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}
