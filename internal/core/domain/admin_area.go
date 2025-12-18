package domain

type AdminArea struct {
	ID         int            `json:"id"`
	Name       string         `json:"name"`
	ISOCode    string         `json:"iso_code"`
	AdminLevel int32          `json:"admin_level"`
	ParentID   *int           `json:"parent_id"`
	ParentCode *string        `json:"parent_code"`
	Geometry   map[string]any `json:"geom"`
}
