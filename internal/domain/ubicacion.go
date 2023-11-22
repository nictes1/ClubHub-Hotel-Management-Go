package domain

// Location representa una ubicación geográfica de la franquicia o compañía.
type Location struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Address   string  `json:"address"`
	ZipCode   string  `json:"zip_code"`
	Latitude  float64 `json:"latitude,omitempty"`  // Latitud para ubicaciones geográficas
	Longitude float64 `json:"longitude,omitempty"` // Longitud para ubicaciones geográficas
}
