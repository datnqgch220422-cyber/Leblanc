package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GraphQLRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

type GraphQLResponse struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []string    `json:"errors,omitempty"`
}

// Handler creates a Gin handler for GraphQL
func Handler() gin.HandlerFunc {
	resolver := &Resolver{}

	return func(c *gin.Context) {
		var req GraphQLRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, GraphQLResponse{
				Errors: []string{err.Error()},
			})
			return
		}

		ctx := context.Background()
		result, err := executeQuery(ctx, resolver, req.Query, req.Variables)
		if err != nil {
			c.JSON(http.StatusOK, GraphQLResponse{
				Errors: []string{err.Error()},
			})
			return
		}

		c.JSON(http.StatusOK, GraphQLResponse{Data: result})
	}
}

// Simple query executor (this is a simplified version - in production, use a proper GraphQL library)
func executeQuery(ctx context.Context, resolver *Resolver, query string, variables map[string]interface{}) (map[string]interface{}, error) {
	// Basic routing based on the presence of operation names; responses are wrapped
	// in a field keyed by the operation to match the GraphQL spec and frontend expectations.
	if contains(query, "query") {
		if contains(query, "drink(") {
			id, _ := variables["id"].(string)
			drink, err := resolver.Drink(ctx, id)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{"drink": drink}, nil
		}
		if contains(query, "drinks") {
			drinks, err := resolver.Drinks(ctx)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{"drinks": drinks}, nil
		}
	}

	if contains(query, "mutation") {
		if contains(query, "createBooking") {
			var input CreateBookingInput
			if inputData, ok := variables["input"].(map[string]interface{}); ok {
				jsonData, _ := json.Marshal(inputData)
				_ = json.Unmarshal(jsonData, &input)
				booking, err := resolver.CreateBooking(ctx, input)
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{"createBooking": booking}, nil
			}
		}
		if contains(query, "recommendFromFeatures") {
			var emotionFit EmotionFitInput
			if emotionData, ok := variables["emotionFit"].(map[string]interface{}); ok {
				jsonData, _ := json.Marshal(emotionData)
				_ = json.Unmarshal(jsonData, &emotionFit)
			}

			var caffeine, temp *string
			var sweetness *int

			if c, ok := variables["caffeine"].(string); ok {
				caffeine = &c
			}
			if t, ok := variables["temp"].(string); ok {
				temp = &t
			}
			if s, ok := variables["sweetness"].(float64); ok {
				sInt := int(s)
				sweetness = &sInt
			}

			scores, err := resolver.RecommendFromFeatures(ctx, emotionFit, caffeine, temp, sweetness)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{"recommendFromFeatures": scores}, nil
		}
	}

	return nil, fmt.Errorf("unsupported query")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
