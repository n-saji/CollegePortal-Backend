package service

import (
	"CollegeAdministration/config"
	"CollegeAdministration/models"
	"CollegeAdministration/utils"
	"fmt"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) CreateNewAccount(acc *models.Account) error {

	err := s.ValidateLogin(acc.Info.Credentials.EmailId, acc.Info.Credentials.Password)
	if err != nil {
		log.Println("failed to validate credentials")
		return err
	}
	exists, err := s.CheckEmailExist(acc.Info.Credentials.EmailId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("email already exists")
	}
	password, err := bcrypt.GenerateFromPassword([]byte(acc.Info.Credentials.Password), 10)
	if err != nil {
		log.Println("failed to encrypt password")
		return err
	}
	acc.Info.Credentials.Password = string(password)
	err = s.daos.CreateAccount(acc)
	if err != nil {
		log.Println("failed to create account")
		return err
	}
	allCourses, err := s.RetrieveCA()
	if err != nil || len(allCourses) == 0 {
		log.Println("unable to get any courses")
		return err
	}
	instructorDetails := &models.InstructorDetails{
		Id:              acc.Id,
		InstructorCode:  "-",
		InstructorName:  acc.Name,
		Department:      "Empty Department",
		CourseId:        allCourses[0].Id,
		CourseName:      "Empty Course",
		ClassesEnrolled: models.CourseInfo{},
		Info:            models.Instructor_Info{},
	}
	err = s.daos.InsertInstructorDetails(instructorDetails)
	if err != nil {
		log.Println("failed to insert into instructor details")
		if err != nil {
			err := s.daos.DeleteAccount(acc.Id)
			if err != nil {
				log.Println("failed to revert account creation changes - delete account")
				return err
			}
		}
		return err
	}
	instructorLogin := models.InstructorLogin{
		Id:       acc.Id,
		EmailId:  acc.Info.Credentials.EmailId,
		Password: acc.Info.Credentials.Password,
	}
	err = s.daos.CreateInstructorLogin(instructorLogin)
	if err != nil {
		log.Println("failed to insert into instructor login")

		err := s.daos.DeleteAccount(acc.Id)
		if err != nil {
			log.Println("failed to revert account creation changes - delete account")
			return err
		}
		return err

	}

	go s.GenerateOTPAndStore(acc.Info.Credentials.EmailId)

	go utils.StoreMessages("New Instructor Signed-up!", acc.Name, config.AccountTypeInstructor, "*")

	return nil
}

func (s *Service) VerifyAccountStatusById(acc_id string) (bool, error) {

	account_id_uuid, err := uuid.Parse(acc_id)
	if err != nil {
		log.Println("failed to parse account id")
		return false, err
	}
	accnt, err := s.daos.GetAccountByID(account_id_uuid)
	if err != nil {
		log.Println("failed to get account by id")
		return false, err
	}

	return accnt.Verified, nil
}

func (s *Service) SendResetPasswordMail(emailId string) error {
	// Check if the email exists in the database
	exists, err := s.CheckEmailExist(emailId)
	if err != nil {
		log.Println("failed to check email existence")
		return err
	}
	if !exists {
		return fmt.Errorf("email does not exist")
	}

	// Generate a reset password token

	account_id, err := s.daos.GetIDUsingEmail(emailId)
	if err != nil {
		log.Println("failed to get account by email")
		return err
	}
	account_id_uuid, err := uuid.Parse(account_id)
	if err != nil {
		log.Println("failed to parse account id")
		return err
	}
	token, err := s.CreateResetPasswordToken(account_id)
	if err != nil {
		log.Println("failed to create reset password token")
		return err
	}
	account_name, err := s.daos.GetAccountNameById(account_id_uuid)
	if err != nil {
		log.Println("failed to get account name by id")
		return err
	}

	// Send the reset password email
	err = utils.SendResetPasswordEmail(emailId, token.String(), account_id, account_name.Name)
	if err != nil {
		log.Println("failed to send reset password email")
		return err
	}

	return nil
}

func (s *Service) ResetPassword(req models.ResetPasswordReq) error {
	// Check if the email exists in the database
	err := s.ValidatePasswordReq(&req)
	if err != nil {
		return fmt.Errorf("failed to validate reset password request")
	}

	exists, err := s.CheckEmailExist(req.EmailId)
	if err != nil {
		return fmt.Errorf("failed to check email existence")
	}
	if !exists {
		return fmt.Errorf("email does not exist")
	}

	// Validate the token
	account_id, err := s.daos.GetIDUsingEmail(req.EmailId)
	if err != nil {
		return fmt.Errorf("failed to get account by email")
	}
	
	account_id_uuid, err := uuid.Parse(account_id)
	if err != nil {
		return fmt.Errorf("failed to parse account id")
	}

	if account_id != req.AccountID {
		return fmt.Errorf("account id mismatch")
	}

	token_uuid, err := uuid.Parse(req.Token)
	if err != nil {
		return fmt.Errorf("failed to parse token")
	}

	valid, err := s.CheckTokenValidity(token_uuid)
	if err != nil {
		return fmt.Errorf("token expired or invalid")
	}
	if !valid {
		return fmt.Errorf("token is invalid or expired")
	}

	credentials := &models.InstructorLogin{
		Id:       account_id_uuid,
		Password: req.Password,
	}

	err = s.UpdateInstructorCredentials(credentials)
	if err != nil {
		log.Println("failed to update account password")
		return err
	}

	err = s.DisableAllTokensForAccount(account_id)
	if err != nil {
		log.Println("failed to delete all token")
		return err
	}

	return nil
}
