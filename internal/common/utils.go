package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

func PrettyPrintJSON(body io.ReadCloser) error {
	// Check if the body exists
	if body == nil {
		return fmt.Errorf("body is nil")
	}

	// Read the body
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("error reading body: %v", err)
	}

	// Close the body when done
	defer body.Close()

	// Create a buffer to store the pretty-printed JSON
	var prettyJSON bytes.Buffer

	// Indent the JSON with 4 spaces
	err = json.Indent(&prettyJSON, bodyBytes, "", "    ")
	if err != nil {
		return fmt.Errorf("error formatting JSON: %v", err)
	}

	// Print the pretty JSON
	fmt.Println(prettyJSON.String())

	return nil
}

// PrettyPrint converts any struct to a pretty-printed JSON string
func PrettyPrintStruct(v interface{}) {
	// Marshal the interface with indentation (4 spaces)
	prettyJSON, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		log.Fatalf("Failed to marshal interface: %v", err)
	}

	// Print the formatted JSON
	fmt.Println(string(prettyJSON))
}
