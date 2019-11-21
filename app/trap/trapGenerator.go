package trap

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/pafrias/2cgaming-api/utils"
)

/**tiers and costs:
1		1-6
2		7-13
3		14-21
4		22-30
5		31-40
6		41-51
7		52+



*/

func (a *App) TestGenerator() http.HandlerFunc {
	// tierBudgets := map[int]int{
	// 	1: 6,
	// 	2: 13,
	// 	3: 21,
	// 	4: 30,
	// 	5: 40,
	// 	6: 51,
	// 	7: 63,
	// }

	return func(res http.ResponseWriter, req *http.Request) {
		ctx := context.TODO()

		trap, _ := a.trapGenerator(ctx, 50)

		data, err := json.Marshal(trap)
		if a.Test500Error(err, res) {
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		res.Write(data)
		fmt.Println("Trap Complete")
	}
}

type c map[string]interface{}
type s []map[string]interface{}

func (a *App) trapGenerator(con context.Context, budget int) ([]interface{}, error) {
	result := []interface{}{}

	// select necessary components
	rows, _ := a.readComponents(con, "gen")
	components, err := utils.ScanRowsToArray(rows)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return nil, err
	}

	// --> select a random effect, deduct cost from budget
	fmt.Print("filtering rows...")
	effects := filterByType(components, "effect")
	if len(effects) > 0 {
		effect, cost := selectRandomComponent(effects)
		budget -= cost
		result = append(result, effect)
	}
	fmt.Println("remaining budget: ", budget)

	// --> select a random target within budget, deduct cost from budget
	targets := filterByBudget(filterByType(components, "target"), budget)
	if len(targets) > 0 {
		target, cost := selectRandomComponent(targets)
		budget -= cost
		result = append(result, target)
	}
	fmt.Println("remaining budget: ", budget)

	// --> select a random trigger within budget, deduct cost from budget
	triggers := filterByBudget(filterByType(components, "trigger"), budget)
	if len(triggers) > 0 {
		trigger, cost := selectRandomComponent(triggers)
		budget -= cost
		result = append(result, trigger)
	}
	fmt.Println("remaining budget: ", budget)

	// select additional components
	// --> 50% chance to select another trigger, deduct...
	// --> 50% chance to select another target, deduct...
	// --> at tier 2 and above, 50% chance to select another component, deduct...
	// ---> if selected, at tier 4 and above, 25% chance to select another component, deduct...
	// ----> if selected, at tier 6 and above, 12.5% chance to select another component, deduct...

	// upgrade component
	// --> count all the component tiers within budget
	// --> filter and count all upgrades within budget
	// ---> select a random tier or upgrade from those options

	// continue upgrading until budget === 0

	// marshall the result to JSON and return

	return result, nil
}

func filterByType(slice s, t string) s {
	filtered := s{}

	for _, val := range slice {
		if val["type"] == t {
			filtered = append(filtered, val)
		}
	}
	return filtered
}

func filterByBudget(slice s, budget int) s {
	filtered := s{}
	var str string
	var ok bool
	var cost int

	for _, val := range slice {
		str, ok = val["costp"].(string)
		if ok {
			str = strings.Split(str, "\t")[0]
			c, _ := strconv.ParseInt(str, 0, 0)
			cost = int(c)
		} else {
			cost = val["cost"].(int)
		}
		if cost <= budget {
			filtered = append(filtered, val)
		}
	}
	return filtered
}

func selectRandomComponent(slice []map[string]interface{}) (map[string]interface{}, int) {
	var i int = len(slice)
	i = rand.Intn(i)
	result := slice[i]
	var cost int
	str, ok := result["costp"].(string)
	if ok {
		x := strings.Split(str, "\t")[0]
		y, _ := strconv.ParseInt(x, 0, 0)
		cost = int(y)
	} else {
		cost = result["cost"].(int)
	}
	return result, cost
}
