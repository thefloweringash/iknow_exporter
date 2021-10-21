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

type CumulativeStats map[string]struct {
	Started                int `json:"cumulative_items_started"`
	TimeMillis             int `json:"cumulative_items_study_time_millis"`
	TotalLogHalflifeMillis int `json:"cumulative_total_log_halflife_millis"`
	Checkpoint1            int `json:"cumulative_items_reached_checkpoint_1"`
	Checkpoint2            int `json:"cumulative_items_reached_checkpoint_2"`
	Checkpoint3            int `json:"cumulative_items_reached_checkpoint_3"`
}

func (self IknowClient) GetCumulativeStats() (CumulativeStats, error) {
	req, err := http.NewRequest("GET", "https://iknow.jp/api/v2/statistics/learning_engine/cumulative?application_domains[]=items", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+self.Secret)

	resp, err := self.client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("API returned non-200 response: %d", resp.StatusCode))
	}

	var result CumulativeStats
	decoder := json.NewDecoder(resp.Body)

	if err = decoder.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
