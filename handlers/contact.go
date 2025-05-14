package handlers

import (
	"fmt"
	"github.com/hismama/personal_website/config"
	"github.com/hismama/personal_website/utils"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

type ContactPageData struct {
	Success bool
	Error   string
}

func Contact(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:

		err := utils.RenderTemplate(w, r, "contact.gohtml", ContactPageData{})
		utils.Check(err, w)

	case http.MethodPost:

		parseErr := r.ParseForm()
		if parseErr != nil {
			err := utils.RenderTemplate(w, r, "contact.gohtml", ContactPageData{
				Error: parseErr.Error(),
			})
			utils.Check(err, w)
			return
		}

		//Bot Dump
		if r.FormValue("phone") != "" {
			w.WriteHeader(http.StatusOK)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		message := r.FormValue("message")

		if name == "" || email == "" || message == "" {
			http.Error(w, "Name and Email and Message are required", http.StatusBadRequest)
			return
		}

		_, err := mail.ParseAddress(email)
		if err != nil {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}

		log.Printf("SMTPUser: [%s], ToEmail: [%s]", config.SMTPUser, config.ToEmail)
		if err := sendEmail(name, email, message); err != nil {
			http.Error(w, "Failed to send email: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = utils.RenderTemplate(w, r, "contact.gohtml", ContactPageData{
			Success: true,
		})
		utils.Check(err, w)
		return

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

}

// buildMessage constructs header and body of eMail to be more legitimate
func buildMessage(fromName, fromEmail, body string) []byte {
	headers := map[string]string{
		"From":                      fmt.Sprintf("%s <%s>", "Personal Website Contact", config.SMTPUser),
		"To":                        config.ToEmail,
		"Reply-To":                  fmt.Sprintf("%s <%s>", fromName, fromEmail),
		"Subject":                   fmt.Sprintf("Contact form message from %s", fromName),
		"Date":                      time.Now().Format(time.RFC1123Z),
		"Message-ID":                fmt.Sprintf("<%d-%s>", time.Now().UnixNano(), config.SMTPUser),
		"MIME-Version":              "1.0",
		"Content-Type":              "text/plain; charset=\"utf-8\"",
		"Content-Transfer-Encoding": "quoted-printable",
	}
	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n" + body + "\r\n")
	return []byte(msg.String())
}

// sendEmail takes in name, visitorEmail, and message and tries to send via the
// constants for SMTP server. Uses buildMessage for construction.
func sendEmail(name, visitorEmail, message string) error {
	auth := smtp.PlainAuth("", config.SMTPUser, config.SMTPPass, config.SMTPServer)
	return smtp.SendMail(config.SMTPServer+":"+config.SMTPPort, auth, config.SMTPUser,
		[]string{config.ToEmail}, buildMessage(name, visitorEmail, message))
}
