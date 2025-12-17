package domain

type AdminArea struct {
	ID       int            `json:"id"`
	Name     string         `json:"name"`
	ISOCode  string         `json:"iso_code"`
	Geometry map[string]any `json:"geom"`
}
