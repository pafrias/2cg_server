package models

type Upgrade struct {
	ID          int32          `json:"_id,omitempty" form:"-"`
	Name        string         `json:"name" form:",required"`
	Type        string         `json:"type" form:",required"`
	Component   JsonNullString `json:"component,omitempty" form:"-"`
	ComponentID JsonNullInt32  `json:"componentID,omitempty" form:"component_id,omitempty"`
	Text        string         `json:"text" form:",required"`
	Cost        int32          `json:"cost" form:",required"`
	Max         int32          `json:"max" form:",required"`
}

type Component struct {
	ID   int32          `json:"_id,omitempty"`
	Name string         `json:"name" form:"blasdlkfsj,required"`
	Type JsonNullString `json:"type" form:",required"`
	Text string         `json:"text,omitempty" form:",required"`
	Cost JsonNullInt32  `json:"cost,omitempty" form:",omitempty"`
	P1   JsonNullString `json:"param1,omitempty"`
	P2   JsonNullString `json:"param2,omitempty"`
	P4   JsonNullString `json:"param4,omitempty"`
	P3   JsonNullString `json:"param3,omitempty"`
}
