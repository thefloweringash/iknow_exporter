package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type IknowClient struct {
	Secret string
	client http.Client
}

func (self IknowClient) fetch(url string, result interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+self.Secret)

	resp, err := self.client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("API returned non-200 response: %d", resp.StatusCode))
	}

	decoder := json.NewDecoder(resp.Body)

	if err = decoder.Decode(&result); err != nil {
		return err
	}

	return nil
}

type CumulativeStats map[string]struct {
	Started                int `json:"cumulative_items_started"`
	TimeMillis             int `json:"cumulative_items_study_time_millis"`
	TotalLogHalflifeMillis int `json:"cumulative_total_log_halflife_millis"`
	Checkpoint1            int `json:"cumulative_items_reached_checkpoint_1"`
	Checkpoint2            int `json:"cumulative_items_reached_checkpoint_2"`
	Checkpoint3            int `json:"cumulative_items_reached_checkpoint_3"`
}

func (self IknowClient) GetCumulativeStats() (CumulativeStats, error) {
	var result CumulativeStats
	err := self.fetch("https://iknow.jp/api/v2/statistics/learning_engine/cumulative?application_domains[]=items", &result)
	return result, err
}

type AggregateResponse struct {
	Goals []GoalStats `json:"goals"`
}

type GoalStats struct {
	GoalId int `json:"goal_id"`
	Items  struct {
		EligibleItemsCount int `json:"eligible_items_count"`
		SkippedItemsCount  int `json:"skipped_items_count"`
		StudiedItemsCount  int `json:"studied_items_count"`
	} `json:"items"`
}

func (self IknowClient) GetAggregateStats() ([]GoalStats, error) {
	var result AggregateResponse
	err := self.fetch("https://iknow.jp/api/v2/goals/enrolled/memories/aggregate", &result)
	return result.Goals, err
}
