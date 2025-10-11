package models

type Equipment struct {
	EquipmentId   int    `json:"equipment_id,omitempty"`
	EquipmentName string `json:"equipment_name" form:"equipment_name"`
	Manufacturer  string `json:"manufacturer" form:"manufacturer"`
	Description   string `json:"description" form:"description"`
	ImageURL      string `json:"image_url,omitempty"`
}
