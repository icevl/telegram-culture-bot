package main

type Channel struct {
	Title     string `json:"title"`
	ChannelID int64  `json:"channel_id"`
	Prompt    string `json:"prompt"`
	MinMins   int    `json:"min_mins"`
	MaxMins   int    `json:"max_mins"`
	NextTime  int64  `json:"next_time"`
}

type SaveData struct {
	Channel *Channel
	Next    int64
}
