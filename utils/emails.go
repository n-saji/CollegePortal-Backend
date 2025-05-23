package utils

import (
	"CollegeAdministration/config"
	"context"
	"fmt"
	"os"
	"time"

	mg "github.com/mailgun/mailgun-go/v4"
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
		"University Portal <notifications@nikhilsaji.me>",
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

func SendResetPasswordEmail(emailId string, token string, accountId, name string) error {
	m := mg.NewMessage(
		"University Portal <notifications@nikhilsaji.me>",
		"Reset Your University Portal Password",
		"",
	)
	m.SetTemplate("reset password mail")
	m.AddRecipient(emailId)
	m.AddVariable("user_name", name)
	m.AddVariable("reset_link", fmt.Sprintf("%s/reset-password/%s/%s/%s", config.FRONTEND_URL, token, accountId, emailId))

	err := SendMessage(m)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return err
	}
	fmt.Println("Reset password email sent successfully")
	return nil
}
