package trap

import (
	"github.com/pafrias/2cgaming-api/db/models"
)

type postUpgrade struct {
	ID          int32                 `json:"_id,omitempty" form:"-"`
	Name        string                `json:"name" form:",required"`
	Type        string                `json:"type" form:",required"`
	Component   models.JsonNullString `json:"component,omitempty" form:"-"`
	ComponentID models.JsonNullInt32  `json:"componentID,omitempty" form:"component_id,omitempty"`
	Text        string                `json:"text" form:",required"`
	Cost        int32                 `json:"cost" form:",required"`
	Max         int32                 `json:"max" form:",required"`
}

type postComponent struct {
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

//used for caching
// type componentStore struct {
// 	sync.RWMutex
// 	components map[uint16]component
// }

// type upgradeStore struct {
// 	sync.RWMutex
// 	upgrades map[uint16]upgrade
// }

type trapTemplate struct {
	Tier     int
	Upgrades []upgrade   `json:"upgrades,omitempty"`
	Triggers []component `json:"triggers,omitempty"`
	Targets  []component `json:"targets,omitempty"`
	Effects  []component `json:"effects,omitempty"`
}

type component interface {
	getID() int
	getCost() int
	getUpgrades() []upgrade
	getComponentID() int
}

type effect struct {
	ID       int
	Name     string
	Cost     []int
	Tier     int       `json:"tier"`
	Upgrades []upgrade `json:"upgrades"`
	isDone   bool
}

func (e effect) getCost() int {
	return e.Cost[0]
}

func (e effect) getID() int {
	return e.ID
}

func (e effect) getUpgrades() []upgrade {
	return e.Upgrades
}

func (e effect) getComponentID() int {
	return 0
}

// targets and triggers
type nonEffect struct {
	ID       int
	Name     string
	Cost     int
	Upgrades []upgrade `json:"upgrades"`
}

func (n nonEffect) getCost() int {
	return n.Cost
}
func (n nonEffect) getID() int {
	return n.ID
}

func (n nonEffect) getUpgrades() []upgrade {
	return n.Upgrades
}

func (n nonEffect) getComponentID() int {
	return 0
}

type upgrade struct {
	ID, cID                 int
	Name, cType             string
	Applications, Cost, Max int
}

func (u upgrade) getCost() int {
	return u.Cost
}
func (u upgrade) getID() int {
	return u.ID
}
func (u upgrade) getUpgrades() []upgrade {
	return nil
}

func (u upgrade) getComponentID() int {
	return u.cID
}
