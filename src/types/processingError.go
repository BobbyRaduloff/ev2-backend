package types

import "encoding/json"

type ProcessingError struct {
	Email string
	Error string
}

type ProcessingErrors struct {
	Errors []ProcessingError
}

func (errors *ProcessingErrors) SerializeProcessingErrors() (string, error) {
	jsonData, err := json.Marshal(errors)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func DeserializeProcessingErrors(data string) (*ProcessingErrors, error) {
	var errors ProcessingErrors
	err := json.Unmarshal([]byte(data), &errors)
	if err != nil {
		return nil, err
	}

	return &errors, nil
}
