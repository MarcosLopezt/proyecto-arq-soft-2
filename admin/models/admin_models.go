package models

type Instancia struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
	Estado string `json:"estado"`
	URL    string `json:"url,omitempty"`
}
