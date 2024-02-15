package onesignal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type OneSignal struct {
	ApiUrl  string
	AppId   string
	AuthKey string
}

func NewOneSignal(apiUrl, appId, authKey string) *OneSignal {
	return &OneSignal{
		ApiUrl:  apiUrl,
		AppId:   appId,
		AuthKey: authKey,
	}
}

func (o *OneSignal) SendNotification(title, description string) error {

	payload := map[string]interface{}{
		"app_id":            o.AppId,
		"included_segments": []string{"All"},
		"headings": map[string]string{
			"en": title,
			"id": title,
		},
		"contents": map[string]string{
			"en": description,
			"id": description,
		},
		"name": "M-REPORT",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}

	req, err := http.NewRequest("POST", o.ApiUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", o.AuthKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return err
	}
	defer response.Body.Close()

	log.Printf("Response Status: %s Notification Sent", response.Status)

	return nil
}
