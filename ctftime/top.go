package ctftime

import (
	"context"
	"fmt"
)

type TopTeam struct {
	TeamName string  `json:"team_name"`
	Points   float64 `json:"points"`
	TeamID   int     `json:"team_id"`
}

func GetTopTeams(ctx context.Context, args ...int) (map[string][]TopTeam, error) {
	limit := 10
	if len(args) > 0 {
		limit = args[0]
	}
	client := NewClient("https://ctftime.org/api")
	teamsData := make(map[string][]TopTeam)
	err := client.Get(ctx, fmt.Sprintf("/v1/top/?limit=%d", limit), &teamsData)
	if err != nil {
		return nil, err
	}
	return teamsData, nil
}

func GetTopTeamsByYear(ctx context.Context, year int, args ...int) (map[string][]TopTeam, error) {
	limit := 10
	if len(args) > 0 {
		limit = args[0]
	}
	client := NewClient("https://ctftime.org/api")
	teamsData := make(map[string][]TopTeam)
	err := client.Get(ctx, fmt.Sprintf("/v1/top/%d/?limit=%d", year, limit), &teamsData)
	if err != nil {
		return nil, err
	}
	return teamsData, nil
}

func GetTopTeamsByCountry(ctx context.Context, country string, args ...int) (map[string][]TopTeam, error) {
	limit := 10
	if len(args) > 0 {
		limit = args[0]
	}
	client := NewClient("https://ctftime.org/api")
	teamsData := make(map[string][]TopTeam)
	err := client.Get(ctx, fmt.Sprintf("/v1/top-by-country/%s/?limit=%d", country, limit), &teamsData)
	if err != nil {
		return nil, err
	}
	return teamsData, nil
}
