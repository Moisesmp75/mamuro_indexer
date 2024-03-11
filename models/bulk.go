package models

type DocZinc struct {
	Document  string  `json:"index"`
	ListEmail []Email `json:"records"`
}