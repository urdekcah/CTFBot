package ctftime

import (
	"context"
	"fmt"
)

type Organizer struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type CTFDuration struct {
	Hours int `json:"hours"`
	Day   int `json:"day"`
}

type CTFEvent struct {
	Organizers    []Organizer `json:"organizers"`
	OnSite        bool        `json:"onsite"`
	Finish        string      `json:"finish"`
	Description   string      `json:"description"`
	Weight        float64     `json:"weight"`
	Title         string      `json:"title"`
	URL           string      `json:"url"`
	IsVotableNow  bool        `json:"is_votable_now"`
	Restrictions  string      `json:"restrictions"`
	Format        string      `json:"format"`
	Start         string      `json:"start"`
	Participants  int         `json:"participants"`
	CTFTimeURL    string      `json:"ctftime_url"`
	Location      string      `json:"location"`
	LiveFeed      string      `json:"live_feed"`
	PublicVotable bool        `json:"public_votable"`
	Duration      CTFDuration `json:"duration"`
	Logo          string      `json:"logo"`
	FormatID      int         `json:"format_id"`
	ID            int         `json:"id"`
	CTFID         int         `json:"ctf_id"`
}

func GetEvents(ctx context.Context, args ...int) ([]CTFEvent, error) {
	limit := 100
	if len(args) > 0 {
		limit = args[0]
	}
	client := NewClient("https://ctftime.org/api")
	var events []CTFEvent
	err := client.Get(ctx, fmt.Sprintf("/v1/events/?limit=%d", limit), &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func GetEventsByPeriod(ctx context.Context, start int, args ...int) ([]CTFEvent, error) {
	if start < 0 {
		return nil, fmt.Errorf("invalid start time")
	}

	finish, limit := -1, 100
	if len(args) > 0 {
		finish = args[0]
	}
	if len(args) > 1 {
		limit = args[1]
	}

	if finish < start {
		return nil, fmt.Errorf("invalid finish time")
	}

	var url string
	if finish == -1 {
		url = fmt.Sprintf("/v1/events/?limit=%d&start=%d", start, limit)
	} else {
		url = fmt.Sprintf("/v1/events/?limit=%d&start=%d&finish=%d", start, finish, limit)
	}

	client := NewClient("https://ctftime.org/api")
	var events []CTFEvent
	if err := client.Get(ctx, url, &events); err != nil {
		return nil, err
	}
	return events, nil
}

func GetSpecificEvent(ctx context.Context, id int) (*CTFEvent, error) {
	client := NewClient("https://ctftime.org/api")
	var event CTFEvent
	err := client.Get(ctx, fmt.Sprintf("/v1/events/%d/", id), &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
