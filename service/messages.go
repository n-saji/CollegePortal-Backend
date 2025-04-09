package service

import (
	"CollegeAdministration/config"
	"CollegeAdministration/utils"
	"fmt"

	"github.com/google/uuid"
)

func (s *Service) UpdateMessageStatusAsRead(account_id string) (string, error) {

	account_id_uuid, err := uuid.Parse(account_id)
	if err != nil {
		return "", err
	}
	err = s.daos.UpdateMessageStatusForAccountId(account_id_uuid)
	if err != nil {
		return "", fmt.Errorf("failed to update message status as read for %s: %v", account_id, err)
	}
	return account_id, nil
}

func (s *Service) SendMessageAsBroadCast(broadcastType, msg string) error {
	if broadcastType == "event" {
		utils.StoreMessages("Test", msg, config.AccountTypeInstructor, "")
		return nil
	} else {
		utils.SendMessageToConnectedClients(msg, "Admin", "Test")
		return nil
	}
}
