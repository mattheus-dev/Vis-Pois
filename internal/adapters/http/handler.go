package http

import (
	"log"
	"net/http"
	
	"vis-pois/internal/adapters/formatter"
	"vis-pois/internal/domain/ports"
)

type CSVHandler struct {
	csvProcessor ports.CSVProcessorPort
	logger       *log.Logger
}

func NewCSVHandler(csvProcessor ports.CSVProcessorPort, logger *log.Logger) *CSVHandler {
	return &CSVHandler{
		csvProcessor: csvProcessor,
		logger:       logger,
	}
}

func (h *CSVHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	h.logger.Println("Received request to process CSV file")

	if r.Method != http.MethodPost {
		h.logger.Printf("Error: Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.logger.Printf("Error parsing multipart form: %v", err)
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Printf("Error retrieving file from form: %v", err)
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	h.logger.Printf("Received file: %s, size: %d bytes", header.Filename, header.Size)

	records, err := h.csvProcessor.ProcessFile(header)
	if err != nil {
		h.logger.Printf("Error processing file: %v", err)
		http.Error(w, "Error processing file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	tableFormatter := formatter.NewTableFormatter()
	tableOutput := tableFormatter.FormatAsTable(records)
	
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tableOutput))
}
