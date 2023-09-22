package main

type Url struct {
	Id    int
	Link  string `json:"url"`
	Short string
	Ttl   int
}
