package trap

import "github.com/pafrias/2cgaming-api/db/models"

type Upgrade struct {
	Cost,
	Max int `form:",required"`
	ComponentID models.JsonNullInt32 `json:"component_id" form:",required"`
	Name,
	Text,
	Type string `form:",required"`
}

// used for trap building logic
type trapTemplate struct {
	Tier     int
	Upgrades []upgrade   `json:"upgrades"`
	Triggers []component `json:"triggers"`
	Targets  []component `json:"targets"`
	Effects  []component `json:"effects"`
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

// DEPRECATED
/*
type shortComponent struct {
	ID   string `json:"_id,omitempty"`
	Name string `json:"name"`
	Type string `json:"type"`
}
*/
