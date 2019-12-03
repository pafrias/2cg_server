package trap

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/pafrias/2cgaming-api/utils"
)

func calculateComponentCost(val map[string]interface{}) int {
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
