package controllers

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Befous/backend.befous.com/models"
	"github.com/Befous/backend.befous.com/utils"
	"gopkg.in/gomail.v2"
)

func generateOTP() string {
	otp := make([]byte, 6)
	_, err := rand.Read(otp)
	if err != nil {
		log.Fatal(err)
	}
	for i := range otp {
		otp[i] = '0' + (otp[i] % 10)
	}
	return string(otp)
}

func sendEmail(recipient, otp string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	senderEmail := os.Getenv("EMAIL_SENDER")
	senderPassword := os.Getenv("EMAIL_PASSWORD")

	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; text-align: center; padding: 20px; }
				h2 { color: #333; }
				.otp-code { font-size: 24px; font-weight: bold; color: #ff6600; margin: 10px 0; }
				.note { font-size: 14px; color: #555; margin-top: 20px; }
			</style>
		</head>
		<body>
			<h2>Kode Login</h2>
			<p>Ini adalah kode persetujuan loginmu:</p>
			<p class="otp-code">%s</p>
			<p class="note">
				Jika kamu tidak mengajukan permintaan ini, segera ganti kata sandi akunmu untuk mencegah akses tanpa izin di kemudian hari. 
				Baca <a href="https://support.google.com/accounts/answer/32040" target="_blank">Protecting Your Account</a> untuk tips mengenai kekuatan kata sandi.
			</p>
		</body>
		</html>
	`, otp)

	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", recipient)
	m.SetHeader("Reply-To", senderEmail)
	m.SetHeader("Subject", "Kode Login: "+otp)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)

	return d.DialAndSend(m)
}

func KirimGmail(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")
	if email == "" {
		utils.WriteJSONResponse(w, http.StatusBadRequest, models.Pesan{
			Message: "Email is required",
		})
		return
	}
	otp := generateOTP()
	err := sendEmail(email, otp)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, models.Pesan{
			Message: "Failed to send OTP",
		})
		return
	}
	utils.WriteJSONResponse(w, http.StatusOK, models.Pesan{
		Message: "OTP has been sent successfully",
	})
}
