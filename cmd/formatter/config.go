package main

type config struct {
	Formatters []struct {
		Name      string        `json:"name"`
		On        bool          `json:"on"`
		Parameter []interface{} `json:"parameter"`
	} `json:"formatters"`
}
