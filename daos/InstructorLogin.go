package daos

import (
	"CollegeAdministration/models"
	"fmt"

	"github.com/google/uuid"
)

func (ac *Daos) CheckIDPresent(id uuid.UUID) error {

	var id_exits string
	err := ac.dbConn.Select("id").Table("instructor_logins").Where("id = ?", id).Find(&id_exits).Error
	if err != nil {
		return err
	}
	if id_exits != "" {
		return fmt.Errorf("you have already created, please log out ")
	}
	return nil
}
func (ac *Daos) CreateInstructorLogin(il models.InstructorLogin) error {

	err := ac.dbConn.Table("instructor_logins").Create(&il).Error

	if err != nil {
		return err
	}
	return nil
}

// func (ac *AdminstrationCloud) CheckLoginExits(email, password string) (bool, error) {

// 	var iid string
// 	err := ac.dbConn.Select("id").Table("instructor_logins").Where("email_id = ? AND password = ?", email, password).Find(&iid).Error

// 	if iid == "" {
// 		return false, err
// 	}
// 	return true, nil
// }

func (ac Daos) CheckForEmail(email string) (bool, error) {

	var count int64
	err := ac.dbConn.Select("count(*)").Table("instructor_logins").Where("email_id = ?", email).Find(&count).Error

	if err != nil {
		return false, err
	}
	if count != 0 {
		return true, nil
	}
	return false, nil
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

func (ac *Daos) DeleteInstructorLogin(instructor_id uuid.UUID) error {

	err := ac.dbConn.Where("id = ?", instructor_id).Delete(models.InstructorLogin{}).Error

	if err != nil {
		return err
	}
	return nil
}

func (ac *Daos) GetIDUsingEmail(email string) (string, error) {
	var instructor_id string
	err := ac.dbConn.Model(models.InstructorLogin{}).Select("id").Where("email_id = ?", email).Find(&instructor_id).Error

	if err != nil {
		return "", err
	}
	return instructor_id, nil
}

func (ac *Daos) FetchPasswordUsingEmailID(email string) (string, error) {
	var password string
	err := ac.dbConn.Model(models.InstructorLogin{}).Select("password").Where("email_id = ?", email).Find(&password).Error

	if err != nil {
		return "", err
	}
	return password, nil
}

func (ac *Daos) FetchCredentialsUsingID(id uuid.UUID) (*models.InstructorLogin, error) {

	var credentials *models.InstructorLogin
	id_string := id.String()
	err := ac.dbConn.Model(models.InstructorLogin{}).Select("*").Where("id = ?", id_string).Find(&credentials).Error

	if err != nil {
		return nil, err
	}
	return credentials, nil
}

// func (ac *AdministrationCloud) GetCredentialsForInstructor(id string) (*models.InstructorLogin, error) {
// 	credentials := &models.InstructorLogin{}

// 	err := ac.dbConn.Model(models.InstructorLogin{}).Select("*").Where("id = ?", id).Find(&credentials).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return credentials, nil
// }

func (ac *Daos) UpdateCredentials(cred *models.InstructorLogin) error {

	err := ac.dbConn.Save(&cred).Error
	if err != nil {
		return err
	}

	return nil
}

func (ac *Daos) GetAccountByToken(token uuid.UUID) (*models.Token_generator, error) {
	var account models.Token_generator
	err := ac.dbConn.Model(models.Token_generator{}).Select("account_id").Where("token = ? and is_valid = true", token).Find(&account).Error

	if err != nil {
		return nil, err
	}
	return &account, nil
}
