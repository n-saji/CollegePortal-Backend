package daos

import (
	"CollegeAdministration/models"
	"fmt"

	"github.com/google/uuid"
)

func (ac *Daos) InsertIntoMessages(req *models.Messages) error {

	err := ac.dbConn.Model(models.Messages{}).Create(req).Error
	if err != nil {
		return err
	}
	return nil
}

func (ac *Daos) UpdateMessageStatusForAccountId(account_id uuid.UUID) error {

	err := ac.dbConn.Model(models.Messages{}).Where("account_id = ?", account_id).Update("is_read", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (ac *Daos) UpdateMessageStatusforMsgId(msg_id uuid.UUID) error {

	err := ac.dbConn.Model(models.Messages{}).Where("id = ?", msg_id).Update("is_read", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (ac *Daos) GetActiveMessagesForAccountId(account_id uuid.UUID) ([]models.Messages, error) {

	var messages []models.Messages
	err := ac.dbConn.Model(models.Messages{}).Where("account_id = ? and is_read = false", account_id).Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (ac *Daos) DeleteMessageByAccountId(account_id uuid.UUID) error {

	fmt.Println("Deleting messages for account_id: ", account_id)
	q := ac.dbConn.Model(models.Messages{}).Where("account_id = ?", account_id).Delete(&models.Messages{})
	fmt.Println(q.RowsAffected)
	if q.Error != nil {
		return q.Error
	}
	return nil
}
