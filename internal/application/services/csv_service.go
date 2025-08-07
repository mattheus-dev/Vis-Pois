package services

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"vis-pois/internal/domain/entities"
)

type CSVService struct {
	logger *log.Logger
}

func NewCSVService(logger *log.Logger) *CSVService {
	return &CSVService{
		logger: logger,
	}
}

func (s *CSVService) ProcessFile(file *multipart.FileHeader) ([]entities.Record, error) {
	s.logger.Println("Starting to process CSV file:", file.Filename)

	tempDir := os.TempDir()
	filePath := filepath.Join(tempDir, file.Filename)

	if err := s.saveFile(file, filePath); err != nil {
		s.logger.Printf("Error saving file: %v", err)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}
	s.logger.Println("File saved successfully at:", filePath)

	records, err := s.readCSV(filePath)
	if err != nil {
		s.logger.Printf("Error reading CSV: %v", err)
		return nil, fmt.Errorf("failed to process CSV: %w", err)
	}

	s.logger.Printf("Successfully processed %d records from CSV file", len(records))
	return records, nil
}

func (s *CSVService) saveFile(fileHeader *multipart.FileHeader, filePath string) error {
	s.logger.Println("Saving file to:", filePath)

	uploadedFile, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer uploadedFile.Close()

	destFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, uploadedFile); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return nil
}

func (s *CSVService) readCSV(filePath string) ([]entities.Record, error) {
	s.logger.Println("Starting to read CSV from:", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stats: %w", err)
	}
	if info.Size() == 0 {
		return nil, fmt.Errorf("file is empty")
	}

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading header: %w", err)
		}
		return nil, fmt.Errorf("file has no content")
	}

	headerLine := scanner.Text()
	var separator string
	if strings.Contains(headerLine, ";") {
		separator = ";"
	} else {
		separator = ","
	}

	headers := strings.Split(headerLine, separator)
	columnMap := make(map[string]int)

	for i, header := range headers {
		columnMap[strings.ToLower(strings.TrimSpace(header))] = i
	}

	s.logger.Printf("Found %d columns in CSV header", len(headers))

	var records []entities.Record
	lineNum := 1

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		fields := strings.Split(line, separator)
		if len(fields) < len(columnMap) {
			s.logger.Printf("Warning: Line %d is malformed, skipping. Expected %d fields, got %d", lineNum, len(columnMap), len(fields))
			continue
		}

		record, err := s.parseComplexRecord(fields, columnMap, lineNum)
		if err != nil {
			s.logger.Printf("Warning: Error parsing line %d: %v", lineNum, err)
			continue
		}

		records = append(records, record)
		s.logger.Printf("Processed record ID: %s, Name: %s", record.ID, record.Name)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no valid records found in file")
	}

	return records, nil
}

func (s *CSVService) parseComplexRecord(fields []string, columnMap map[string]int, lineNum int) (entities.Record, error) {
	record := entities.Record{}
	var err error

	getString := func(key string) string {
		if idx, ok := columnMap[key]; ok && idx < len(fields) {
			return strings.TrimSpace(fields[idx])
		}
		return ""
	}

	getFloat := func(key string) (float64, error) {
		if idx, ok := columnMap[key]; ok && idx < len(fields) {
			value := strings.TrimSpace(fields[idx])
			if value == "" {
				return 0, nil
			}
			return strconv.ParseFloat(value, 64)
		}
		return 0, nil
	}

	getInt := func(key string) (int, error) {
		if idx, ok := columnMap[key]; ok && idx < len(fields) {
			value := strings.TrimSpace(fields[idx])
			if value == "" {
				return 0, nil
			}
			return strconv.Atoi(value)
		}
		return 0, nil
	}

	getBool := func(key string) bool {
		if idx, ok := columnMap[key]; ok && idx < len(fields) {
			value := strings.ToLower(strings.TrimSpace(fields[idx]))
			return value == "true" || value == "yes" || value == "1"
		}
		return false
	}

	record.ID = getString("id")
	if record.ID == "" {
		return record, fmt.Errorf("missing required field 'id' in line %d", lineNum)
	}

	record.Name = getString("name")
	if record.Name == "" {
		return record, fmt.Errorf("missing required field 'name' in line %d", lineNum)
	}

	record.Price, err = getFloat("price")
	if err != nil {
		return record, fmt.Errorf("invalid price format in line %d: %w", lineNum, err)
	}

	record.Stock, err = getInt("stock")
	if err != nil {
		return record, fmt.Errorf("invalid stock format in line %d: %w", lineNum, err)
	}

	record.Category = getString("category")
	record.Subcategory = getString("subcategory")
	record.Brand = getString("brand")
	record.Description = getString("description")
	record.ImageURL = getString("image_url")
	record.Dimensions = getString("dimensions")
	record.Color = getString("color")
	record.Material = getString("material")
	record.CountryOfOrigin = getString("country_of_origin")
	record.Manufacturer = getString("manufacturer")
	record.SKU = getString("sku")
	record.Barcode = getString("barcode")

	record.Weight, _ = getFloat("weight")
	record.TaxRate, _ = getFloat("tax_rate")
	record.DiscountPercentage, _ = getFloat("discount_percentage")
	record.Rating, _ = getFloat("rating")
	record.ReviewCount, _ = getInt("review_count")

	record.IsActive = getBool("is_active")

	record.CreatedAt = getString("created_at")
	record.UpdatedAt = getString("updated_at")

	return record, nil
}

func (s *CSVService) parseRecord(fields []string, lineNum int) (entities.Record, error) {
	price, err := strconv.ParseFloat(strings.TrimSpace(fields[2]), 64)
	if err != nil {
		return entities.Record{}, fmt.Errorf("invalid price format in line %d: %w", lineNum, err)
	}

	stock, err := strconv.Atoi(strings.TrimSpace(fields[3]))
	if err != nil {
		return entities.Record{}, fmt.Errorf("invalid stock format in line %d: %w", lineNum, err)
	}

	return entities.Record{
		ID:    strings.TrimSpace(fields[0]),
		Name:  strings.TrimSpace(fields[1]),
		Price: price,
		Stock: stock,
	}, nil
}
