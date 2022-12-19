package server

import (
	"encoding/json"
	"fmt"
	"github.com/Raeein/gmc"
	"github.com/Raeein/gmc/webadvisor"
	"io"
	"log"
)

func validateRegistration(webadvisorService *webadvisor.WebAdvisor, body io.ReadCloser) error {
	data, err := decode(body)
	if err != nil {
		return err
	}
	fmt.Println(data)
	err = webadvisorService.Exists(data)
	if err != nil {
		return err
	}
	return nil
}

func decode(body io.ReadCloser) (gmc.Section, error) {
	data := gmc.Section{}
	err := json.NewDecoder(body).Decode(&data)
	if err != nil {
		log.Println(err)
		return gmc.Section{}, err
	}
	fmt.Println(data)
	return data, nil
}
