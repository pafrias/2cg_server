package trap

import "github.com/pafrias/2cgaming-api/db/models"

type upgrade struct {
	ID          int32                 `json:"_id,omitempty" form:"-"`
	Name        string                `json:"name" form:",required"`
	Type        string                `json:"type" form:",required"`
	Component   models.JsonNullString `json:"component,omitempty" form:"-"`
	ComponentID models.JsonNullInt32  `json:"componentID,omitempty" form:"component_id,omitempty"`
	Text        string                `json:"text" form:",required"`
	Cost        int32                 `json:"cost" form:",required"`
	Max         int32                 `json:"max" form:",required"`
}

type component struct {
	ID   int32                 `json:"_id,omitempty"`
	Name string                `json:"name" form:"blasdlkfsj,required"`
	Type models.JsonNullString `json:"type" form:",required"`
	Text string                `json:"text,omitempty" form:",required"`
	Cost models.JsonNullInt32  `json:"cost,omitempty" form:",omitempty"`
	P1   models.JsonNullString `json:"param1,omitempty"`
	P2   models.JsonNullString `json:"param2,omitempty"`
	P4   models.JsonNullString `json:"param4,omitempty"`
	P3   models.JsonNullString `json:"param3,omitempty"`
}

type shortComponent struct {
	ID   string `json:"_id,omitempty"`
	Name string `json:"name"`
	Type string `json:"type"`
}
