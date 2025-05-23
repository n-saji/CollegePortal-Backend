package service

import (
	"CollegeAdministration/models"
	"CollegeAdministration/utils"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (ac *Service) ValidateLogin(email, password string) error {

	reg := regexp.MustCompile(`^[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z]{2,}$`)
	ok := reg.MatchString(email)
	if !ok {
		return fmt.Errorf("wrong email format")
	} else if password == ":password" {
		return fmt.Errorf("password cant be empty ")
	} else if len(password) < 8 {
		return fmt.Errorf("password length insufficient")
	}
	return nil
}

func (s *Service) ValidateEmail(email string) error {
	ok, _ := regexp.MatchString("@gmail.com", email)
	if !ok || len(email) < 11 {
		return fmt.Errorf("wrong email format")
	}
	return nil
}
func (s *Service) ValidatePassword(password string) error {

	if password == "" {
		return fmt.Errorf("password cant be empty ")
	} else if len(password) < 8 {
		return fmt.Errorf("password length insufficient")
	}
	return nil
}
func (ac *Service) CheckEmailExist(email string) (bool, error) {
	ok, err := ac.daos.CheckForEmail(email)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}
func (ac *Service) CheckCredentials(email, password string) error {

	hashed_password, err1 := ac.daos.FetchPasswordUsingEmailID(email)
	if err1 != nil {
		return err1
	}
	if hashed_password == "" {
		return fmt.Errorf("invalid email id")
	}
	err4 := bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(password))
	if err4 != nil {
		return fmt.Errorf("invalid credentials")
	}

	return nil
}

func (ac *Service) GetTokenAfterLogging(account_id string) (uuid.UUID, error) {
	var token_table models.Token_generator
	account_uuid, err := uuid.Parse(account_id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing uuid")
	}
	token_table.AccountId = account_uuid

	token := uuid.New()
	token_table.Token = token
	token_table.IsValid = true
	tn := time.Now()
	token_table.ValidFrom = tn.Unix()
	validtill := tn.Add(time.Hour * 24)
	token_table.ValidTill = validtill.Unix()

	err = ac.daos.InsertToken(&token_table)
	if err != nil {
		return uuid.Nil, fmt.Errorf("token insertion failed")
	}
	return token_table.Token, nil
}

func (ac *Service) CheckTokenValidity(token uuid.UUID) (bool, error) {

	err1 := ac.CheckTokenExpiry(token)
	if err1 != nil {
		return false, err1
	}
	status, err := ac.daos.GetTokenStatus(token)
	if err != nil {
		return status, err
	}
	return status, nil
}

func (ac *Service) CheckTokenExpiry(token uuid.UUID) error {

	token_details, err := ac.daos.GetTokenStored(token)
	if err != nil {
		return err
	}
	time_now := time.Now()
	diff := token_details.ValidTill - time_now.Unix()
	if diff < 0 {
		err2 := ac.daos.SetTokenFalse(token)
		if err2 != nil {
			return err2
		}
		return fmt.Errorf("token expired! Generate new token")
	}

	return nil
}

func (s *Service) CheckTokenWithCookie(token string) error {

	token_id, err4 := uuid.Parse(token)
	if err4 != nil {
		return fmt.Errorf("error parsing uuid")
	}
	status, err1 := s.CheckTokenValidity(token_id)
	if err1 != nil {
		return err1
	}
	if !status {
		return fmt.Errorf("token expired")
	}
	return nil
}

func (s *Service) ValidateInstructorDetails(iid *models.InstructorDetails) (bool, error) {

	for _, each_letter := range iid.Department {
		if each_letter == '-' {
			return false, fmt.Errorf("wrong format! remove -")
		}
	}
	return true, nil
}

func (s *Service) DisableToken(token string) error {

	parsedToken, err := uuid.Parse(token)
	if err != nil {
		return fmt.Errorf("unable to invalidate token, err- %s", err)
	}
	err = s.daos.DisableToken(parsedToken)
	if err != nil {
		return fmt.Errorf("db error while invalidating token, err - %s", err)
	}
	return nil
}

func (s *Service) GetAccountByToken(token string) (string, error) {

	parsedToken, err := uuid.Parse(token)
	if err != nil {
		return "", fmt.Errorf("unable to parse token, err- %s", err)
	}
	acc, err := s.daos.GetAccountByToken(parsedToken)
	if err != nil {
		return "", fmt.Errorf("db error while fetching account, err - %s", err)
	}
	return acc.AccountId.String(), nil
}

func (s *Service) VerifyAccountWithOTP(email string, otp string) error {

	res, err := s.daos.GetOTPByEmailIdAndOTP(email, otp)
	if err != nil {
		return fmt.Errorf("db error while fetching otp with account, err - %s", err)
	}
	if res == nil {
		return fmt.Errorf("no otp found")
	}
	if res.IsUsed {
		return fmt.Errorf("otp already used")
	}

	if res.ExpiresAt < time.Now().Local().Unix() {
		return fmt.Errorf("otp expired")
	}

	err = s.daos.DeleteOTPByEmailId(email)
	if err != nil {
		return fmt.Errorf("db error while deleting otp, err - %s", err)
	}

	id, err := s.daos.GetIDUsingEmail(email)
	if err != nil {
		return fmt.Errorf("db error while fetching account id, err - %s", err)
	}

	err = s.daos.UpdateAccountStatusAsTrue(id)
	if err != nil {
		return fmt.Errorf("db error while updating account status, err - %s", err)
	}

	return nil
}

func (s *Service) GenerateOTPAndStore(email string) error {

	acc, err := s.daos.GetIDUsingEmail(email)
	if err != nil {
		log.Println("error fetching account by email")
		return err
	}
	accnt_uuid, err := uuid.Parse(acc)
	if err != nil {
		log.Println("error parsing account id")
		return err
	}

	acnt, err := s.daos.GetAccountByID(accnt_uuid)
	if err != nil {
		log.Println("error fetching account by id")
		return err
	}

	old_otps, err := s.daos.GetOTPByAccountID(acc)
	if err != nil {
		log.Println("error fetching old otps")
		return err
	}

	for _, each := range old_otps {
		if each.IsUsed {
			continue
		}
		if each.ExpiresAt > time.Now().Local().Unix() {
			return fmt.Errorf("otp already sent")
		}
		err = s.daos.DeleteOTPByAccountId(acc)
		if err != nil {
			log.Println("error deleting old otp")
			return err
		}
	}

	otp := utils.GenerateOTP(6)
	err = s.daos.InsertOTP(&models.OTP{
		ID:        uuid.New(),
		EmailId:   email,
		AccountID: acnt.Id,
		OTPCode:   otp,
		CreatedAt: time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		IsUsed:    false,
	})
	if err != nil {
		log.Println("error storing otp")
		return err
	}
	go utils.SendAccountCreationOTP(acnt.Name, acnt.Info.Credentials.EmailId, otp)
	return nil
}

func (ac *Service) CreateResetPasswordToken(account_id string) (uuid.UUID, error) {
	token_table := &models.Token_generator{}
	account_uuid, err := uuid.Parse(account_id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing uuid")
	}
	token_table.AccountId = account_uuid

	token := uuid.New()
	token_table.Token = token
	token_table.IsValid = true
	tn := time.Now()
	token_table.ValidFrom = tn.Unix()
	validtill := tn.Add(time.Minute * 15)
	token_table.ValidTill = validtill.Unix()

	err = ac.daos.InsertToken(token_table)
	if err != nil {
		return uuid.Nil, fmt.Errorf("token insertion failed")
	}
	return token_table.Token, nil
}

func (s *Service) ValidatePasswordReq(req *models.ResetPasswordReq) error {

	if req == nil {
		return fmt.Errorf("request is nil")
	}

	if req.EmailId == "" {
		return fmt.Errorf("email_id is required")
	}
	if req.Token == "" {
		return fmt.Errorf("token is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

func (s *Service) DisableAllTokensForAccount(account_id string) error {
	account_uuid, err := uuid.Parse(account_id)
	if err != nil {
		return fmt.Errorf("error parsing uuid")
	}
	err = s.daos.DisableAllTokensForAccount(account_uuid)
	if err != nil {
		return fmt.Errorf("db error while disabling all tokens for account, err - %s", err)
	}
	return nil
}