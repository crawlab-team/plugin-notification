package main

import (
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/matcornic/hermes"
	"gopkg.in/gomail.v2"
	"net/mail"
	"runtime/debug"
	"strconv"
)

func GetMailContent(s *NotificationSetting, m *models.ModelMap) (content string) {
	return GetTaskEmailMarkdownContent(m)
}

func SendMail(s *NotificationSetting, m *models.ModelMap) error {
	// hermes instance
	h := hermes.Hermes{
		Theme: new(hermes.Default),
		Product: hermes.Product{
			Name:      "Crawlab Team",
			Copyright: "© 2021 Crawlab-Team",
		},
	}

	// config
	port, _ := strconv.Atoi(s.Mail.Port)
	password := s.Mail.Password // test password: ALWVDPRHBEXOENXD
	SMTPUser := s.Mail.User
	smtpConfig := smtpAuthentication{
		Server:         s.Mail.Server,
		Port:           port,
		SenderEmail:    s.Mail.SenderEmail,
		SenderIdentity: s.Mail.SenderIdentity,
		SMTPPassword:   password,
		SMTPUser:       SMTPUser,
	}
	options := sendOptions{
		To:      m.User.Email,
		Subject: s.Title,
	}

	// content
	content := GetMailContent(s, m)

	// email instance
	email := hermes.Email{
		Body: hermes.Body{
			Name:         m.User.Username,
			FreeMarkdown: hermes.Markdown(content + GetFooter()),
		},
	}

	// generate html
	html, err := h.GenerateHTML(email)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// generate text
	text, err := h.GeneratePlainText(email)
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	// send the email
	if err := send(smtpConfig, options, html, text); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return err
	}

	return nil
}

type smtpAuthentication struct {
	Server         string
	Port           int
	SenderEmail    string
	SenderIdentity string
	SMTPUser       string
	SMTPPassword   string
}

// sendOptions are options for sending an email
type sendOptions struct {
	To      string
	Subject string
	Cc      string
}

// send email
func send(smtpConfig smtpAuthentication, options sendOptions, htmlBody string, txtBody string) error {

	if smtpConfig.Server == "" {
		return errors.New("SMTP server config is empty")
	}
	if smtpConfig.Port == 0 {
		return errors.New("SMTP port config is empty")
	}

	if smtpConfig.SMTPUser == "" {
		return errors.New("SMTP user is empty")
	}

	if smtpConfig.SenderIdentity == "" {
		return errors.New("SMTP sender identity is empty")
	}

	if smtpConfig.SenderEmail == "" {
		return errors.New("SMTP sender email is empty")
	}

	if options.To == "" {
		return errors.New("no receiver emails configured")
	}

	from := mail.Address{
		Name:    smtpConfig.SenderIdentity,
		Address: smtpConfig.SenderEmail,
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from.String())
	m.SetHeader("To", options.To)
	m.SetHeader("Subject", options.Subject)
	if options.Cc != "" {
		m.SetHeader("Cc", options.Cc)
	}

	m.SetBody("text/plain", txtBody)
	m.AddAlternative("text/html", htmlBody)

	d := gomail.NewDialer(smtpConfig.Server, smtpConfig.Port, smtpConfig.SMTPUser, smtpConfig.SMTPPassword)

	return d.DialAndSend(m)
}

func GetFooter() string {
	return `
[Github](https://github.com/crawlab-team/crawlab) | [Documentation](http://docs.crawlab.cn) | [Docker](https://hub.docker.com/r/tikazyq/crawlab)
`
}

func GetTaskEmailMarkdownContent(m *models.ModelMap) string {
	n := m.Node
	s := m.Spider
	t := m.Task
	ts := m.TaskStat
	errMsg := ""
	statusMsg := fmt.Sprintf(`<span style="color:green">%s</span>`, t.Status)
	if t.Status == constants.TaskStatusError {
		errMsg = " with errors"
		statusMsg = fmt.Sprintf(`<span style="color:red">%s</span>`, t.Status)
	}
	return fmt.Sprintf(`
Your task has finished%s. Please find the task info below.

|Key:|Value|
|--: | :--|
|**Task ID:** | %s|
|**Task Status:** | %s|
|**Task Param:** | %s|
|**Spider ID:** | %s|
|**Spider Name:** | %s|
|**Node:** | %s|
|**Create Time:** | %s|
|**Start Time:** | %s|
|**Finish Time:** | %s|
|**Wait Duration:** | %d sec|
|**Runtime Duration:** | %d sec|
|**Total Duration:** | %d sec|
|**Number of Results:** | %d|
|**Error:** | <span style="color:red">%s</span>|

Please login to Crawlab to view the details.
`,
		errMsg,
		t.Id,
		statusMsg,
		t.Param,
		s.Id.Hex(),
		s.Name,
		n.Name,
		utils.GetLocalTimeString(ts.CreateTs),
		utils.GetLocalTimeString(ts.StartTs),
		utils.GetLocalTimeString(ts.EndTs),
		ts.WaitDuration/1000,
		ts.RuntimeDuration/1000,
		ts.TotalDuration/1000,
		ts.ResultCount,
		t.Error,
	)
}
