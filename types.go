package main

type Channel struct {
	Title    string `json:"title"`
	ChatID   int64  `json:"chat_id"`
	Prompt   string `json:"prompt"`
	Image    string `json:"image"`
	MinMins  int    `json:"min_mins"`
	MaxMins  int    `json:"max_mins"`
	NextTime int64  `json:"next_time"`
}

type SaveData struct {
	Channel *Channel
	Next    int64
}
