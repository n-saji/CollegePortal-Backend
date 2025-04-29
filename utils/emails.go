package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	mg "github.com/mailgun/mailgun-go/v4"
)

var (
	FRONTEND_URL = os.Getenv("FRONTEND_URL")
)

func SendMessage(m *mg.Message) error {

	apiKey := os.Getenv("MAIL_GUN_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("api key not set")
	}
	domain := os.Getenv("MAIL_GUN_DOMAIN")
	if domain == "" {
		return fmt.Errorf("domain name not set")
	}
	if domain == "" || apiKey == "" {
		return fmt.Errorf("domain or api key not set")
	}

	mg := mg.NewMailgun(domain, apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return err
	}
	fmt.Printf("ID: %s\n", id)
	fmt.Println("Email sent successfully")
	return nil
}

func SendAccountCreationOTP(name, emailId, otp string) error {

	m := mg.NewMessage(
		"University Portal <postmaster@notificationbot.me>",
		"Account Creation OTP",
		"",
	)
	m.SetTemplate("otp creation mail")
	m.AddRecipient(emailId)
	m.AddVariable("otp", otp)
	m.AddVariable("user_name", name)

	err := SendMessage(m)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return err
	}
	fmt.Println("Account creation OTP sent successfully")
	return nil
}

// func SendPasswordReset(emailId string) error {
// 	title := "Password Reset"
// 	message := "Please reset your password by clicking on the link below:\n" +
// 		"http://localhost:8080/reset-password?email=" + emailId
// 	err := SendMessage(title, message, emailId)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func SendAccountCreation(emailId string) error {
// 	title := "Account Creation"
// 	message := "Your account has been created successfully. Please login to your account."
// 	err := SendMessage(title, message, emailId)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func SendAccountPasswordChange(emailId string) error {
// 	title := "Account Password Change"
// 	message := "Your account password has been changed successfully. Please login to your account."
// 	err := SendMessage(title, message, emailId)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
