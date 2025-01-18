package types

import "encoding/json"

type ProcessingResult struct {
	Email     string
	RequestId int

	IsValid bool

	IsNonpersonal bool
	IsDisposable  bool

	MX       string
	HasMX    bool
	HasSPF   bool
	HasDMARC bool
	HasDKIM  bool

	Handshake     int
	HandshakeName string
}

type ProcessingResults struct {
	RequestId int
	Results   []ProcessingResult
}

func (results *ProcessingResults) SerializeEmailResults() (string, error) {
	jsonData, err := json.Marshal(results)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func DeserializeEmailResults(data string) (*ProcessingResults, error) {
	var results ProcessingResults
	err := json.Unmarshal([]byte(data), &results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}
