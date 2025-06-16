package main

import "time"

type Event struct {
	Type      string    `json:"type"`
	Actor     Actor     `json:"actor"`
	Repo      Repo      `json:"repo"`
	Payload   Payload   `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

type Actor struct {
	Login string `json:"login"`
}

type Repo struct {
	Name string `json:"name"`
}

type Payload struct {
	Ref     string   `json:"ref"`
	RefType string   `json:"ref_type"`
	Commits []Commit `json:"commits"`
	Action  string   `json:"action"`
}

type Commit struct {
	SHA string `json:"sha"`
}
