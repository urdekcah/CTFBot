package ctftime

import (
	"context"
	"fmt"
)

type CTFScore struct {
	TeamID int     `json:"team_id"`
	Points float64 `json:"points"`
	Place  int     `json:"place"`
}

type CTFResult struct {
	Title  string     `json:"title"`
	Scores []CTFScore `json:"scores"`
	Time   int        `json:"time"`
}

func GetResults(ctx context.Context, args ...int) ([]CTFResult, error) {
	limit := 10
	if len(args) > 0 {
		limit = args[0]
	}
	client := NewClient("https://ctftime.org/api")
	var results []CTFResult
	err := client.Get(ctx, fmt.Sprintf("/v1/results/?limit=%d", limit), &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func GetResultsByYear(ctx context.Context, year int, args ...int) ([]CTFResult, error) {
	limit := 10
	if len(args) > 0 {
		limit = args[0]
	}
	client := NewClient("https://ctftime.org/api")
	var results []CTFResult
	err := client.Get(ctx, fmt.Sprintf("/v1/results/%d/?limit=%d", year, limit), &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
