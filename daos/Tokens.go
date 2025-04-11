package daos

import (
	"CollegeAdministration/models"

	"github.com/google/uuid"
)

func (ac *Daos) GetAccountByToken(token uuid.UUID) (*models.Token_generator, error) {
	var account models.Token_generator
	err := ac.dbConn.Model(models.Token_generator{}).Select("account_id").Where("token = ? and is_valid = true", token).Find(&account).Error

	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (ac *Daos) GetTokenStored(token uuid.UUID) (*models.Token_generator, error) {

	var toke_details models.Token_generator
	err := ac.dbConn.Model(toke_details).Where("token = ?", token).Find(&toke_details).Error

	if err != nil {
		return nil, err
	}
	return &toke_details, nil

}
func (ac *Daos) SetTokenFalse(token uuid.UUID) error {

	err := ac.dbConn.Model(models.Token_generator{}).Where("token = ?", token).Update("is_valid", false).Error

	if err != nil {
		return err
	}
	return nil

}

func (ac *Daos) InsertToken(tg models.Token_generator) error {

	err := ac.dbConn.Table("token_generators").Create(&tg).Error
	if err != nil {
		return err
	}
	return nil
}

func (ac *Daos) GetTokenStatus(token uuid.UUID) (bool, error) {
	var status bool
	err := ac.dbConn.Model(models.Token_generator{}).Select("is_valid").Where("token = ?", token).Find(&status).Error

	if err != nil {
		return status, err
	} else {
		return status, nil
	}
}

func (ac *Daos) DeleteToken(token uuid.UUID) error {
	err := ac.dbConn.Where("token = ?", token).Delete(models.Token_generator{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (ac *Daos) DeleteTokenByAccountId(account_id uuid.UUID) error {
	err := ac.dbConn.Where("account_id = ?", account_id).Delete(models.Token_generator{}).Error
	if err != nil {
		return err
	}
	return nil
}
