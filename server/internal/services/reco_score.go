package services

import (
	"math"
	"sort"

	"leblanc/server/internal/models"
)

// DrinkScore represents a drink with its recommendation score
type DrinkScore struct {
	DrinkID string
	Score   float64
}

// ScoreDrinks scores drinks based on optional preferences (caffeine, temp, sweetness)
func ScoreDrinks(drinks []models.Drink, caffeine, temp string, sweetness int) []DrinkScore {
	var scores []DrinkScore

	for _, drink := range drinks {
		// Initialize preference score
		prefScore := 1.0

		// Apply caffeine preference if specified
		if caffeine != "" && drink.Caffeine != caffeine {
			prefScore *= 0.7
		}

		// Apply temperature preference if specified
		if temp != "" && drink.Temp != "either" && drink.Temp != temp {
			prefScore *= 0.7
		}

		// Apply sweetness preference if specified
		if sweetness > 0 {
			sweetnessDiff := math.Abs(float64(drink.Sweetness - sweetness))
			sweetnessScore := 1.0 - (sweetnessDiff / 10.0)
			if sweetnessScore < 0 {
				sweetnessScore = 0
			}
			prefScore *= (0.5 + 0.5*sweetnessScore)
		}

		scores = append(scores, DrinkScore{
			DrinkID: drink.ID.Hex(),
			Score:   math.Round(prefScore*1000) / 1000,
		})
	}

	// Sort by score descending
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	// Return top 5
	if len(scores) > 5 {
		scores = scores[:5]
	}

	return scores
}
