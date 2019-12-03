package trap

import (
	"fmt"
	"math/rand"
)

type trapTemplate struct {
	Upgrades []map[string]interface{} `json:"upgrades,omitempty"`
	Triggers []map[string]interface{} `json:"triggers,omitempty"`
	Targets  []map[string]interface{} `json:"targets,omitempty"`
	Effects  []map[string]interface{} `json:"effects,omitempty"`
}

type builderFunc func(components []map[string]interface{}, upgrades []map[string]interface{}, budget int) (trapTemplate, error)

func buildRandomizedTrap(components []map[string]interface{}, ups []map[string]interface{}, budget int) (trapTemplate, error) {
	var i int
	trap := trapTemplate{}
	tier := calculateTierfromBudget(budget)
	componentCount := 1
	chosenIDs := []int{}

	// select components
	// possible error for empty array after filter?
	effects := filterByType(components, "effect")
	effects, effect, budget := selectComponent(effects, budget)
	effect["tier"] = 1
	if id, ok := effect["id"].(int); ok {
		chosenIDs = append(chosenIDs, id)
	}
	trap.Effects = append(trap.Effects, effect)
	i = rand.Intn(8)
	for j := 3; i > j && tier > j-2; j += 2 { // chance to select additional effect for tier 2 and above (total 4 at highest tiers)
		budget--
		effects := filterByCost(effects, budget)
		effects, effect, budget = selectComponent(effects, budget)
		componentCount++
		effect["tier"] = 1
		if id, ok := effect["id"].(int); ok {
			chosenIDs = append(chosenIDs, id)
		}
		trap.Effects = append(trap.Effects, effect)
	}
	//fmt.Println("Effects Selected\n\t remaining budget: ", budget)

	targets := filterByType(components, "target")
	targets = filterByCost(targets, budget)
	targets, target, budget := selectComponent(targets, budget)
	if id, ok := target["id"].(int); ok {
		chosenIDs = append(chosenIDs, id)
	}
	trap.Targets = append(trap.Targets, target)
	if i = rand.Intn(2); i > 0 { // chance to select additional target
		budget--
		targets = filterByCost(targets, budget)
		targets, target, budget = selectComponent(targets, budget)
		if id, ok := target["id"].(int); ok {
			chosenIDs = append(chosenIDs, id)
		}
		trap.Targets = append(trap.Targets, target)
	}
	//fmt.Println("Targets Selected\n\t remaining budget: ", budget)

	triggers := filterByType(components, "trigger")
	triggers = filterByCost(triggers, budget)
	triggers, trigger, budget := selectComponent(triggers, budget)
	if id, ok := trigger["id"].(int); ok {
		chosenIDs = append(chosenIDs, id)
	}
	trap.Triggers = append(trap.Triggers, trigger)
	if i = rand.Intn(2); i > 0 { // chance to select additional trigger
		budget--
		triggers = filterByCost(triggers, budget)
		triggers, trigger, budget = selectComponent(triggers, budget)
		if id, ok := trigger["id"].(int); ok {
			chosenIDs = append(chosenIDs, id)
		}
		trap.Triggers = append(trap.Triggers, trigger)
	}
	//fmt.Println("Triggers Selected\n\t remaining budget: ", budget)

	fmt.Printf("Trap Components Selected:\n%+v\n", trap)

	// upgrade components
	// --> count all the component tiers within budget
	// --> filter and count all upgrades within budget
	// ---> select a random tier or upgrade from those options
	// continue upgrading until budget === 0
	// tierCount := 6 * componentCount
	upgrades := filterByCompID(ups, chosenIDs)
	for budget > 0 {
		upgrades = filterByCost(upgrades, budget)
		i = rand.Intn(len(upgrades)) //  + tierCount
		upgrade := upgrades[i]
		id, _ := upgrade["id"].(int)
		// get cost and max apps
		cost := calculateComponentCost(upgrade)
		limit, ok := upgrade["max"].(int)
		if !ok {
			fmt.Printf("upgrade id:%v has no max\n", id)
			continue
		}
		if limit == 0 {
			limit = 10
		}
		maximum := budget / cost

		if maximum > limit {
			maximum = limit
		}
		// fmt.Printf("Cost of upgrade(%v) is %v\n\t Can buy %v times\n", id, cost, maximum)

		i = rand.Intn(maximum) + 1
		budget -= i * cost
		upgrade["applications"] = i

		// get type, and apply to random component or trap as a whole
		t, ok := upgrade["type"].(string)
		if !ok {
			fmt.Printf("upgrade id:%v has no type\n", id)
			continue
		}

		// check for duplicate upgrades, add them together if possible
		if t == "universal" {
			trap.Upgrades = append(trap.Upgrades, upgrade)
		} else {
			var arr []map[string]interface{}
			var l int
			switch t {
			case "trigger":
				l = len(trap.Triggers)
				arr = trap.Triggers
			case "target":
				l = len(trap.Targets)
				arr = trap.Targets
			case "effect":
				l = len(trap.Effects)
				arr = trap.Effects
			}
			fmt.Println(arr)
			i = rand.Intn(l)

			dest, _ := arr[i]["upgrades"].([]map[string]interface{})
			dest = append(dest, upgrade)
			arr[i]["upgrades"] = dest
			fmt.Println(arr[i])
		}
	}

	fmt.Println("remaining budget: ", budget)

	return trap, nil

}
