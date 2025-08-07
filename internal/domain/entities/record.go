package entities

type Record struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	Stock          int     `json:"stock"`
	Category       string  `json:"category"`
	Subcategory    string  `json:"subcategory"`
	Brand          string  `json:"brand"`
	Description    string  `json:"description"`
	ImageURL       string  `json:"image_url"`
	Weight         float64 `json:"weight"`
	Dimensions     string  `json:"dimensions"`
	Color          string  `json:"color"`
	Material       string  `json:"material"`
	CountryOfOrigin string `json:"country_of_origin"`
	Manufacturer   string  `json:"manufacturer"`
	SKU            string  `json:"sku"`
	Barcode        string  `json:"barcode"`
	TaxRate        float64 `json:"tax_rate"`
	DiscountPercentage float64 `json:"discount_percentage"`
	Rating         float64 `json:"rating"`
	ReviewCount    int     `json:"review_count"`
	IsActive       bool    `json:"is_active"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}
