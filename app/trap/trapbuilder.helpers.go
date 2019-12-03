package trap

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/pafrias/2cgaming-api/utils"
)

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

func calculateComponentCost(val map[string]interface{}) int {
	var result int
	cost, ok := val["cost"].([]int)
	if ok {
		result = cost[1]
	} else {
		result, _ = val["cost"].(int)
	}
	return result
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
	}
	return 7
}

func filterByType(slice []map[string]interface{}, target string) []map[string]interface{} {
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

func filterByCost(slice []map[string]interface{}, budget int) []map[string]interface{} {
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

func filterByCompID(slice []map[string]interface{}, target []int) []map[string]interface{} {
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

func selectComponent(slice []map[string]interface{}, budget int) ([]map[string]interface{}, map[string]interface{}, int) {
	i := rand.Intn(len(slice))
	component := slice[i]
	component["upgrades"] = []map[string]interface{}{}
	budget -= calculateComponentCost(component)
	newSlice := append(append([]map[string]interface{}{}, slice[:i]...), slice[i+1:]...)

	return newSlice, component, budget
}

func purchaseEffectTier(effect map[string]interface{}, budget int) (map[string]interface{}, int) {
	tier, ok := effect["tier"].(int)
	if !ok {
		fmt.Println("effect tier error")
		return effect, budget
	} else if tier == 7 {
		return effect, budget
	} else if tier > 7 {
		fmt.Println("tier is greater than 7...")
		return effect, budget
	}

	costs, ok := effect["cost"].([]int)
	if !ok {
		fmt.Println("effect cost error")
		return effect, budget
	}

	currentCost := costs[tier-1]
	options := []int{}
	var dif, i int

	for _, val := range costs[tier:] {
		dif = val - currentCost
		if dif <= budget {
			options = append(options, val-currentCost)
		}
	}

	if len(options) > 0 {
		i = rand.Intn(len(options))
		fmt.Printf("Tier for component(%v) increased from %v to %v\n", effect["name"], tier, tier+1+i)
		effect["tier"] = tier + 1 + i
		budget -= options[i]
	}

	return effect, budget
}
