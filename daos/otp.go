package daos

import "CollegeAdministration/models"

func (db *Daos) InsertOTP(otp *models.OTP) error {
	err := db.dbConn.Table("otps").Create(&otp).Error
	if err != nil {
		return err
	}
	return nil
}
func (db *Daos) GetOTPByAccountID(accountID string) ([]*models.OTP, error) {
	var otp []*models.OTP
	err := db.dbConn.Table("otps").Where("account_id = ?", accountID).Find(&otp).Error
	if err != nil {
		return nil, err
	}
	return otp, nil
}
func (db *Daos) UpdateOTP(otp *models.OTP) error {
	err := db.dbConn.Table("otps").Where("account_id = ?", otp.AccountID).Updates(otp).Error
	if err != nil {
		return err
	}
	return nil
}
func (db *Daos) DeleteOTPByAccountId(account_id string) error {
	err := db.dbConn.Table("otps").Where("account_id = ?", account_id).Update("is_used", false).Error
	if err != nil {
		return err
	}
	return nil
}
func (db *Daos) DeleteOTPByEmailId(email_id string) error {
	err := db.dbConn.Table("otps").Where("email_id = ?", email_id).Update("is_used", true).Error
	if err != nil {
		return err
	}
	return nil
}
func (db *Daos) GetOTPByAccountIdAndOTP(accountID string, otp string) (*models.OTP, error) {
	var otpModel *models.OTP
	err := db.dbConn.Table("otps").Where("account_id = ? AND otp_code = ?", accountID, otp).First(&otpModel).Error
	if err != nil {
		return nil, err
	}
	return otpModel, nil
}

func (db *Daos) GetOTPByEmailIdAndOTP(emailId string, otp string) (*models.OTP, error) {
	var otpModel *models.OTP
	err := db.dbConn.Table("otps").Where("email_id = ? AND otp_code = ?", emailId, otp).First(&otpModel).Error
	if err != nil {
		return nil, err
	}
	return otpModel, nil
}
