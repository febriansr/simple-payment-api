package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/febriansr/simple-payment-api/model/app_error"
)

func ReadParseJSON(fileName string, target any) error {
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return app_error.InternalServerError("Failed to read JSON data: " + err.Error())
	}

	err = json.Unmarshal(fileContent, target)
	if err != nil {
		return app_error.InternalServerError("Failed to unmarshal JSON data: " + err.Error())
	}

	return nil
}

func WriteJSON(fileName string, data any) error {
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return app_error.InternalServerError("Failed to marshal JSON data: " + err.Error())
	}

	err = ioutil.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		return app_error.InternalServerError("Failed to write JSON data to file: " + err.Error())
	}

	return nil
}
