package ports

import (
	"mime/multipart"

	"vis-pois/internal/domain/entities"
)

type CSVProcessorPort interface {
	ProcessFile(file *multipart.FileHeader) ([]entities.Record, error)
}
