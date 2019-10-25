package model

// A Component is a trigger, target or effect of a trap
// does not work yet
type Component struct {
	ID     int    `json:"_id,omitempty"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Text   string `json:"text"`
	Cost   int    `json:"cost,omitempty"`
	Param1 string `json:"param1,omitempty"`
	Param2 string `json:"param2,omitempty"`
	Param3 string `json:"param3,omitempty"`
	Param4 string `json:"param4,omitempty"`
}

// An Upgrade is a trigger, target or effect of a trap
type Upgrade struct {
	ID          int    `json:"_id,omitempty"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Text        string `json:"text"`
	Cost        int    `json:"cost"`
	ComponentID int    `json:"component_id"`
}
