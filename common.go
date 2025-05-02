package main

import (
	"encoding/json"
	"fmt"
)

func displayJSON(jsonObj any) error {
	output, err := json.Marshal(jsonObj)
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
