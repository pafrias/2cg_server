package trap

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	slice "github.com/pafrias/array-utils"
)

func scanComponentsToStruct(rows *sql.Rows) ([]component, []component, []component) {
	var targets, triggers, effects []component

	var id int
	var cost sql.NullInt32
	var name, param1, cType string
	var costs []int
	var err error
	for rows.Next() {
		err = rows.Scan(&id, &name, &cost, &param1, &cType)
		if err == nil {
			switch cType {
			case "target":
				targets = append(targets, nonEffect{id, name, int(cost.Int32), []upgrade{}})
			case "trigger":
				triggers = append(triggers, nonEffect{id, name, int(cost.Int32), []upgrade{}})
			case "effect":
				costs, _ = costStringToArray(param1)
				effects = append(effects, effect{id, name, costs, 1, []upgrade{}, false})
			}
		} else {
			fmt.Println(err)
			panic(err)
		}
	}
	return targets, triggers, effects
}

func scanUpgradesToStruct(rows *sql.Rows) []component {
	var upgrades []component
	var err error
	for rows.Next() {
		var u upgrade
		var x sql.NullInt64
		err = rows.Scan(&u.ID, &u.Name, &u.cType, &x, &u.Cost, &u.Max)
		if err == nil {
			if x.Valid {
				u.cID = int(x.Int64)
			} else {
				u.cID = 0
			}
			upgrades = append(upgrades, u)
		} else {
			fmt.Println("error scanning: ", err)
			fmt.Println(u)
			panic(err)
		}
	}
	return upgrades
}

func costStringToArray(str string) ([]int, error) {
	result := []int{}
	arr := strings.Split(str, "\t")
	for _, val := range arr[1:] {
		num, _ := strconv.ParseInt(val, 0, 0)
		n := int(num)
		result = append(result, n)
	}
	return result, nil
}

func filterComponents(src []component, testFunc func(val component) bool) []component {
	var results []component
	filtered, err := slice.Filter(src, testFunc)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	for _, val := range filtered {
		if v, ok := val.(component); ok {
			results = append(results, v)
		}
	}
	return results
}

func filterByCost(src []component, budget int) []component {
	testFunc := func(val component) bool {
		return val.getCost() <= budget
	}
	return filterComponents(src, testFunc)
}

func filterByCompID(src []component, target []int) []component {
	testFunc := func(val component) bool {
		id := val.getComponentID()
		if id == 0 {
			return true
		}
		has, err := slice.Has(target, id)
		if err != nil {
			return false
		}
		return has
	}
	return filterComponents(src, testFunc)
}

func purchaseComponent(src []component, budget int) ([]component, component, int) {
	i := rand.Intn(len(src))
	result := src[i]
	cost := result.getCost()
	var remainingComponents []component
	remainingComponents = append(remainingComponents, src[:i]...)
	remainingComponents = append(remainingComponents, src[i+1:]...)

	return remainingComponents, result, cost
}

func purchaseEffectTier(c component, budget int) (effect, int) {
	eff, ok := c.(effect)
	if !ok {
		fmt.Println("effect error")
		return eff, 0
	} else if eff.isDone {
		return eff, 0
	}
	var min, max int = eff.Tier, 7
	var currentCost int = eff.Cost[min-1]
	for index, tierCost := range eff.Cost {
		if tierCost-currentCost > budget {
			max = index
			break
		}
	}
	if min >= max {
		eff.isDone = true
		fmt.Println(eff.getID(), " is done")
		return eff, 0
	}
	var options []int = eff.Cost[min:max]
	i := rand.Intn(len(options))
	eff.Tier += (i + 1)
	cost := (options[i] - currentCost)
	return eff, cost
}

func purchaseUpgrade(u upgrade, budget int, has bool) (upgrade, int) {
	var cost int
	if !has { // did not have upgrade
		maximum := calculateUpgradeMax(u, budget)
		applications := rand.Intn(maximum) + 1 // perhaps skew based on tier towards lower or higher numbers
		cost = applications * u.Cost
		u.Applications = applications
	} else if u.Max != u.Applications { // already had upgrade
		fmt.Println(u.Applications)
		u.Applications++
		cost = u.Cost
	}
	return u, cost
}

func calculateUpgradeMax(u upgrade, budget int) int {
	max := u.Max
	limit := budget / u.Cost
	if max == 0 {
		max = 5
	}
	if max > limit {
		return limit
	}
	return max
}

func calculateTierfromBudget(budget int) int {
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
	default:
		return 7
	}
}
