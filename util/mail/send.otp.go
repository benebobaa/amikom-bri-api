package mail

import "fmt"

func GetSenderParamEmailVerification(toEmail, verificationLink string) (string, string, []string) {
	subject := "Email Verification for Amikom Pedia"

	content := fmt.Sprintf(`
        <h1>Hello %s,</h1>

        <p>Welcome to Amikom Pedia! To complete your registration, please click the following link to verify your email:</p>

        <a href="%s" target="_blank">%s</a>

        <p>This link is valid for a single use and will expire in 5 minutes. If you did not attempt to create an account with Amikom Pedia, please disregard this email. Your account security is important to us.</p>

        <p>Thank you for choosing Amikom Pedia!</p>
    `, toEmail, verificationLink, verificationLink)

	to := []string{toEmail}

	return subject, content, to
}
