package trap

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/pafrias/2cgaming-api/utils"
)

type trapTemplate struct {
	triggers,
	targets,
	effects,
	upgrades []map[string]interface{}
}
type builderFunc func(components []map[string]interface{}, upgrades []map[string]interface{}, budget int) (trapTemplate, error)

func createBuilderFunc() builderFunc {
	calculateComponentCost := func(val map[string]interface{}) int {
		var cost int
		str, ok := val["costp"].(string)
		if ok {
			x := strings.Split(str, "\t")[1]
			y, _ := strconv.ParseInt(x, 0, 0)
			cost = int(y)
		} else {
			cost = val["cost"].(int)
		}
		return cost
	}

	calculateTrapTier := func(budget int) int {
		switch {
		case budget < 7:
			return 1
		case budget < 14:
			return 2
		case budget < 22:
			return 3
		case budget < 31:
			return 4
		case budget < 41:
			return 5
		case budget < 52:
			return 6
		}
		return 7
	}

	filterByCost := func(slice []map[string]interface{}, budget int) []map[string]interface{} {
		results := []map[string]interface{}{}
		testFunc := func(val map[string]interface{}) bool {
			return calculateComponentCost(val) <= budget
		}
		filtered, err := utils.Filter(slice, testFunc)
		if err != nil {
			fmt.Println(err)
			return results
		}
		for _, val := range filtered {
			if v, ok := val.(map[string]interface{}); ok {
				results = append(results, v)
			}
		}
		return results
	}

	filterByType := func(slice []map[string]interface{}, target string) []map[string]interface{} {
		var results = []map[string]interface{}{}
		testFunc := func(val map[string]interface{}) bool {
			return val["type"] == target
		}
		filtered, err := utils.Filter(slice, testFunc)
		if err != nil {
			fmt.Println(err)
			return results
		}
		for _, val := range filtered {
			if v, ok := val.(map[string]interface{}); ok {
				results = append(results, v)
			}
		}
		return results
	}

	filterByCompID := func(slice []map[string]interface{}, target []int) []map[string]interface{} {
		var results = []map[string]interface{}{}
		testFunc := func(val map[string]interface{}) bool {
			id := val["component_id"]
			if id == nil {
				return true
			} else if i, ok := id.(int); ok {
				has, err := utils.SliceHas(target, i)
				if err == nil {
					return has
				}
			} else {
				fmt.Println("it was not ok")
			}
			return false
		}
		filtered, err := utils.Filter(slice, testFunc)
		if err != nil {
			fmt.Println(err)
			return results
		}
		for _, val := range filtered {
			if v, ok := val.(map[string]interface{}); ok {
				results = append(results, v)
			}
		}
		return results
	}

	// heavy lifter, removes a component from the slice, returns slice, component, and budget
	selectComponent := func(slice []map[string]interface{}, budget int) ([]map[string]interface{}, map[string]interface{}, int) {
		i := rand.Intn(len(slice))
		component := slice[i]
		component["upgrades"] = []map[string]interface{}{}
		budget -= calculateComponentCost(component)
		newSlice := append(append([]map[string]interface{}{}, slice[:i]...), slice[i+1:]...)

		return newSlice, component, budget
	}

	return func(cmps []map[string]interface{}, ups []map[string]interface{}, budget int) (trapTemplate, error) {
		fmt.Println("Building Trap...")
		var i int
		trap := trapTemplate{}
		tier := calculateTrapTier(budget)
		componentCount := 1
		chosenIDs := []int{}

		// select components
		// possible error for empty array after filter
		effects := filterByType(cmps, "effect")
		effects, effect, budget := selectComponent(effects, budget)
		effect["tier"] = 1
		if id, ok := effect["id"].(int); ok {
			chosenIDs = append(chosenIDs, id)
		}
		trap.effects = append(trap.effects, effect)
		i = rand.Intn(8)
		for j := 3; i > j && tier > j-2; j += 2 { // chance to select additional effect for tier 2 and above (total 4 at highest tiers)
			effects := filterByCost(effects, budget)
			effects, effect, budget = selectComponent(effects, budget)
			componentCount++
			effect["tier"] = 1
			if id, ok := effect["id"].(int); ok {
				chosenIDs = append(chosenIDs, id)
			}
			trap.effects = append(trap.effects, effect)
		}
		fmt.Println("Effects Selected => remaining budget: ", budget)

		targets := filterByType(cmps, "target")
		targets = filterByCost(targets, budget)
		targets, target, budget := selectComponent(targets, budget)
		if id, ok := target["id"].(int); ok {
			chosenIDs = append(chosenIDs, id)
		}
		trap.targets = append(trap.targets, target)
		if i = rand.Intn(2); i > 0 { // chance to select additional target
			targets = filterByCost(targets, budget)
			targets, target, budget = selectComponent(targets, budget)
			if id, ok := target["id"].(int); ok {
				chosenIDs = append(chosenIDs, id)
			}
			trap.targets = append(trap.targets, target)
		}
		fmt.Println("Targets Selected => remaining budget: ", budget)

		triggers := filterByType(cmps, "trigger")
		triggers = filterByCost(triggers, budget)
		triggers, trigger, budget := selectComponent(triggers, budget)
		if id, ok := trigger["id"].(int); ok {
			chosenIDs = append(chosenIDs, id)
		}
		trap.triggers = append(trap.triggers, trigger)
		if i = rand.Intn(2); i > 0 { // chance to select additional trigger
			triggers = filterByCost(triggers, budget)
			triggers, trigger, budget = selectComponent(triggers, budget)
			if id, ok := trigger["id"].(int); ok {
				chosenIDs = append(chosenIDs, id)
			}
			trap.triggers = append(trap.triggers, trigger)
		}
		fmt.Println("Triggers Selected => remaining budget: ", budget)

		fmt.Println("chosen components: ", chosenIDs)

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
			if i >= len(ups) {           // buy a higher tier of effect
			} else { // buy an upgrade, check for repeats
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
				fmt.Printf("Cost of upgrade(%v) is %v\n Can buy %v times", id, cost, maximum)

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
					trap.upgrades = append(trap.upgrades, upgrade)
				} else {
					var arr *[]map[string]interface{}
					switch t {
					case "trigger":
						arr = &trap.triggers
					case "target":
						arr = &trap.targets
					case "effect":
						arr = &trap.effects
					}
					i = rand.Intn(len(arr))
				} if t == "trigger" {
					
					slice, _ := trap.triggers[i]["upgrades"].([]map[string]interface{})

				} else if t == "target" {
					i = rand.Intn(len(trap.targets))
					target = trap.targets[i]

				} else if t == "effect" {
					i = rand.Intn(len(trap.effects))
					effect = trap.effects[i]

				}
			}

			budget -= 1
		}

		return trap, nil

	}

	// marshall the result to JSON and return

}
