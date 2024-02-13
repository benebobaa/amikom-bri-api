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

func GetSenderParamResetPassword(toEmail, resetLink string) (string, string, []string) {
	subject := "Password Reset for Amikom Pedia"

	content := fmt.Sprintf(`
        <h1>Hello %s,</h1>

        <p>We received a request to reset your password for your Amikom Pedia account. To proceed with the password reset, please click the following link:</p>

        <a href="%s" target="_blank">%s</a>

        <p>This link is valid for a single use and will expire in 5 minutes. If you did not request a password reset, please disregard this email. Your account security is important to us.</p>

        <p>Thank you for choosing Amikom Pedia!</p>
    `, toEmail, resetLink, resetLink)

	to := []string{toEmail}

	return subject, content, to
}

func GetSenderParamTransferNotification(toEmail string, amount int64) (string, string, []string) {
	subject := "Transfer Notification from Amikom Pedia"

	content := fmt.Sprintf(`
        <h1>Hello %s,</h1>

        <p>We want to inform you that a transfer has been made from your Amikom Pedia account.</p>

        <p>Transfer Details:</p>
        <ul>
            <li>Amount: %d $</li>
        </ul>

        <p>If you authorized this transfer, you can ignore this email. If you did not initiate this transfer or have any concerns, please contact our support team immediately.</p>

        <p>Thank you for using Amikom Pedia!</p>
    `, toEmail, amount)

	to := []string{toEmail}

	return subject, content, to
}

func GetReceiverParamTransferNotification(receiverEmail, senderName string, amount int64) (string, string, []string) {
	subject := "Incoming Transfer Notification from Amikom Pedia"

	content := fmt.Sprintf(`
        <h1>Hello %s,</h1>

        <p>You have received a transfer from %s through Amikom Pedia.</p>

        <p>Transfer Details:</p>
        <ul>
            <li>Amount: %d $</li>
            <!-- You can add more details like date, transaction ID, etc. as needed -->
        </ul>

        <p>If you were not expecting this transfer or have any concerns, please contact our support team immediately.</p>

        <p>Thank you for using Amikom Pedia!</p>
    `, receiverEmail, senderName, amount)

	to := []string{receiverEmail}

	return subject, content, to
}
