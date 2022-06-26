package config

import (
	"encoding/json"
	"io"
	"os"
)

type config struct {
	Formatters []struct {
		Name string `json:"name"`
		On   bool   `json:"on"`
	} `json:"formatters"`
}

func ReadConfig() (*config, error) {
	m := &config{}

	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make([]byte, 64) // ?
	str := ""
	for {
		fl, err := file.Read(data)
		if err == io.EOF {

			break
		}
		str += string(data[:fl])
	}

	dt := []byte(str)
	err = json.Unmarshal(dt, m)
	if err != nil {
		return nil, err
	}

	return m, nil
}
