package ctftime

import (
	"context"
	"fmt"
)

type BriefTeam struct {
	Aliases  []string `json:"aliases"`
	Country  string   `json:"country"`
	Academic bool     `json:"academic"`
	ID       int      `json:"id"`
	Name     string   `json:"name"`
}

type TeamsResult struct {
	Limit  int         `json:"limit"`
	Result []BriefTeam `json:"result"`
	Offset int         `json:"offset"`
}

type TeamRating struct {
	OrganizerPoints *float64 `json:"organizer_points,omitempty"`
	RatingPoints    *float64 `json:"rating_points,omitempty"`
	RatingPlace     *int     `json:"rating_place,omitempty"`
	CountryPlace    *int     `json:"country_place,omitempty"`
}

type Team struct {
	Academic     bool                  `json:"academic"`
	PrimaryAlias string                `json:"primary_alias"`
	Name         string                `json:"name"`
	Rating       map[string]TeamRating `json:"rating"`
	Logo         string                `json:"logo"`
	Country      string                `json:"country"`
	ID           int                   `json:"id"`
	Aliases      []string              `json:"aliases"`
}

func GetTeams(ctx context.Context, args ...int) (*TeamsResult, error) {
	limit := 100
	if len(args) > 0 {
		limit = args[0]
	}
	client := NewClient("https://ctftime.org/api")
	var teams TeamsResult
	err := client.Get(ctx, fmt.Sprintf("/v1/teams/?limit=%d", limit), &teams)
	if err != nil {
		return nil, err
	}
	return &teams, nil
}

func GetTeamByID(ctx context.Context, id int) (*Team, error) {
	client := NewClient("https://ctftime.org/api")
	var team Team
	err := client.Get(ctx, fmt.Sprintf("/v1/teams/%d/", id), &team)
	if err != nil {
		return nil, err
	}
	return &team, nil
}
