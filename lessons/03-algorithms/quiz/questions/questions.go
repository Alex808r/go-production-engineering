package questions

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Quiz struct {
	Country string `json:"country"`
	Capital string `json:"capital"`
}

func LoadQuestions() ([]Quiz, error) {
	jsonData, err := os.ReadFile("quiz.json")

	if err != nil {
		return nil, fmt.Errorf("error reading quiz data: %w", err)
	}

	var list []Quiz

	if err := json.Unmarshal(jsonData, &list); err != nil {
		return nil, fmt.Errorf("error unmarshalling quiz data: %w", err)
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("quiz must contain at least one question")
	}
	for i := range list {
		list[i].Country = strings.TrimSpace(list[i].Country)
		list[i].Capital = strings.TrimSpace(list[i].Capital)
		if list[i].Country == "" || list[i].Capital == "" {
			return nil, fmt.Errorf("question %d: country and capital must be non-empty", i+1)
		}
	}
	return list, nil
}
