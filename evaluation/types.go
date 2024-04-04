package main

// Define a struct to represent your JSON data
type PromptAnswer struct {
	Prompt     string `json:"prompt"`
	Answer     string `json:"answer"`
	Evaluation string `json:"evaluation"`
	Section    string `json:"section"`
}

type JSONData struct {
	Model          string         `json:"model"`
	Recommendation string         `json:"recommendation"`
	Sections       []PromptAnswer `json:"sections"`
}
