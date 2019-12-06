package trap

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/pafrias/2cgaming-api/utils"
)

func buildRandomizedTrap(componentRows, upgradeRows *sql.Rows, budget int) (trapTemplate, error) {
	var trap trapTemplate
	trap.Tier = calculateTierfromBudget(budget)
	selected := []int{}

	targets, triggers, effects := scanComponentsToStruct(componentRows)
	effects = filterByCost(effects, budget)
	var c component
	var cost int
	effects, c, cost = purchaseComponent(effects, budget)
	budget -= cost
	selected = append(selected, c.getID())
	trap.Effects = append(trap.Effects, c)

	// chance to select additional effect for tier 2 and above (total 4 at highest tiers)
	var i int
	for tierRequired := 2; trap.Tier >= tierRequired; tierRequired += 2 {
		i = rand.Intn(6 + tierRequired*2)
		if i < trap.Tier {
			effects = filterByCost(effects, budget)
			if len(effects) != 0 {
				effects, c, cost = purchaseComponent(effects, budget)
				budget -= cost
				selected = append(selected, c.getID())
				trap.Effects = append(trap.Effects, c)
			}
		}
	}

	targets = filterByCost(targets, budget)
	targets, c, cost = purchaseComponent(targets, budget)
	budget -= cost
	selected = append(selected, c.getID())
	trap.Targets = append(trap.Targets, c)

	if i = rand.Intn(3); i > 0 && budget > 1 { // chance to select additional target
		budget--
		targets = filterByCost(targets, budget)
		_, c, cost = purchaseComponent(targets, budget)
		budget -= cost
		selected = append(selected, c.getID())
		trap.Targets = append(trap.Targets, c)
	}

	triggers = filterByCost(triggers, budget)
	triggers, c, cost = purchaseComponent(triggers, budget)
	budget -= cost
	selected = append(selected, c.getID())
	trap.Triggers = append(trap.Triggers, c)

	if i = rand.Intn(3); i > 0 && budget > 1 { // chance to select additional trigger
		budget--
		triggers = filterByCost(triggers, budget)
		_, c, cost = purchaseComponent(triggers, budget)
		budget -= cost
		selected = append(selected, c.getID())
		trap.Triggers = append(trap.Triggers, c)
	}

	upgrades := scanUpgradesToStruct(upgradeRows)
	upgrades = filterByCompID(upgrades, selected)
	tiersComplete := false

	for timer := 100; budget > 0 && timer > 0; timer-- {
		i = rand.Intn(2)
		if i == 1 && !tiersComplete { // purchase higher tier
			i = rand.Intn(len(trap.Effects))
			c, cost = purchaseEffectTier(trap.Effects[i], budget)
			if cost == 0 {
				tiersComplete, _ = utils.Every(trap.Effects, func(val component) bool {
					eff, _ := val.(effect)
					return eff.isDone
				})
			}
			budget -= cost
			trap.Effects[i] = c
		} else { // purchase upgrade
			upgrades = filterByCost(upgrades, budget)
			i = rand.Intn(len(upgrades))
			u, _ := upgrades[i].(upgrade)

			var focus []upgrade
			var upgradeType string = u.cType
			switch upgradeType {
			case "universal":
				focus = trap.Upgrades
			case "trigger":
				i = rand.Intn(len(trap.Triggers))
				focus = trap.Triggers[i].getUpgrades()
			case "target":
				i = rand.Intn(len(trap.Targets))
				focus = trap.Targets[i].getUpgrades()
			case "effect":
				i = rand.Intn(len(trap.Effects))
				focus = trap.Effects[i].getUpgrades()
			}

			index, _ := utils.Any(focus, func(val upgrade) bool { return val.ID == u.ID })
			if index != -1 {
				u, cost = purchaseUpgrade(focus[index], budget, true)
				focus[index] = u
			} else {
				u, cost = purchaseUpgrade(u, budget, false)
				focus = append(focus, u)
			}
			if cost == 0 {
				fmt.Printf("Budget = %v\n\t%+v\n", budget, upgrades)
				continue
			}
			budget -= cost

			switch upgradeType {
			case "universal":
				trap.Upgrades = focus
			case "trigger":
				val, _ := trap.Triggers[i].(nonEffect)
				val.Upgrades = focus
				trap.Triggers[i] = val
			case "target":
				val, _ := trap.Targets[i].(nonEffect)
				val.Upgrades = focus
				trap.Targets[i] = val
			case "effect":
				val, _ := trap.Effects[i].(effect)
				val.Upgrades = focus
				trap.Effects[i] = val
			}

		}
	}
	return trap, nil
}
