package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
	"user_manager/src/entity"
)

const (
	EventUserInserted = "user_inserted"
	EventUserUpdated  = "user_updated"
	EventUserDeleted  = "user_deleted"
)

var subscribedEndpoints []string

type NotifyService struct {}

type payload struct {
	Event       string
	Timestamp   string
	User       *entity.User
}
func (notifyService NotifyService) SendEvent(event string, user *entity.User) error {
	payloadStruct := payload{
		Event:     event,
		Timestamp: time.Now().Format(time.RFC3339),
		User:      user,
	}

	payloadString, err := json.Marshal(payloadStruct)
	if err != nil {
		return err
	}
	payloadBuffer := bytes.NewBuffer(payloadString)

	for _, endpoint := range subscribedEndpoints {
		_, err = http.Post(endpoint, "application/json", payloadBuffer)
		if err != nil {
			return err
		}

		// todo: create a retry buffer if the request fails
	}

	return nil
}

func (notifyService NotifyService) Subscribe(endpoint string) error {
	// todo: validate endpoint

	subscribedEndpoints = append(subscribedEndpoints, endpoint)
	return nil
}