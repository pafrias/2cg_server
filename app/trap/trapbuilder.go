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

	effects := []map[string]interface{}{}
	for _, val := range filterByType(components, "effect") {
		str, ok := val["costp"].(string)
		if ok {
			arr, err := costStringToArray(str)
			if err == nil {
				val["cost"] = arr
				delete(val, "costp")
				effects = append(effects, val)
			}
		}
	}

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

	// handle higher tiers for trap effects
	// should tier/upgrade be 50/50 chance?
	upgrades := filterByCompID(ups, chosenIDs)
	for budget > 0 {
		i = rand.Intn(2)
		if i == 0 { // purchase higher tier
			fmt.Println("adding tier")
			i = rand.Intn(len(trap.Effects))
			effect, budget = purchaseEffectTier(trap.Effects[i], budget)
			trap.Effects[i] = effect
		} else { // purchase upgrade
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
				i = rand.Intn(l)

				dest, _ := arr[i]["upgrades"].([]map[string]interface{})
				arr[i]["upgrades"] = append(dest, upgrade)
			}
		}
	}

	return trap, nil

}
