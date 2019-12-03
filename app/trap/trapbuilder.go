package trap

import (
	"fmt"
	"math/rand"
)

func buildRandomizedTrap(components []map[string]interface{}, ups []map[string]interface{}, budget int) (trapTemplate, error) {
	var i int
	trap := trapTemplate{}
	tier := calculateTierfromBudget(budget)
	componentCount := 1
	chosenIDs := []int{}

	effects := filterByType(components, "effect")
	effects, effect, budget := selectComponent(effects, budget)
	effect["tier"] = 1
	if id, ok := effect["id"].(int); ok {
		chosenIDs = append(chosenIDs, id)
	}
	trap.Effects = append(trap.Effects, effect)

	// chance to select additional effect for tier 2 and above (total 4 at highest tiers)
	i = rand.Intn(8)
	for j := 3; i > j && tier > j-2; j += 2 {
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

	// handle tier upgrades for trap effects
	upgrades := filterByCompID(ups, chosenIDs)
	for budget > 0 {
		upgrades = filterByCost(upgrades, budget)
		i = rand.Intn(len(upgrades)) //  + tierCount
		upgrade := upgrades[i]
		id, _ := upgrade["id"].(int)
		cost := calculateComponentCost(upgrade)
		limit, ok := upgrade["max"].(int)
		if !ok || limit == 0 {
			limit = 10
		}
		maximum := budget / cost

		if maximum > limit {
			maximum = limit
		}

		i = rand.Intn(maximum) + 1
		budget -= i * cost
		upgrade["applications"] = i

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
