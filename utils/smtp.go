package utils

import (
	"backend/internal/config"
	"fmt"
	"log"
	"net/smtp"
	"sync"
	"time"
)

// EmailTask represents an email to be sent
type EmailTask struct {
	To      string
	URL     string
	Config  *config.SMTPConfig
	Created time.Time
}

var (
	emailQueue    = make(chan EmailTask, 100) // Buffer for up to 100 emails
	workerStarted = false
	workerMutex   sync.Mutex
)

// StartEmailWorker starts the background worker that processes emails
func StartEmailWorker() {
	workerMutex.Lock()
	defer workerMutex.Unlock()

	if workerStarted {
		return
	}

	workerStarted = true
	go processEmails()
}

// processEmails runs in the background and sends emails with delay
func processEmails() {
	var lastSentTime time.Time

	for task := range emailQueue {
		// Calculate time since last email
		timeSinceLastEmail := time.Since(lastSentTime)

		// If less than 30 seconds have passed since the last email, wait
		if !lastSentTime.IsZero() && timeSinceLastEmail < 30*time.Second {
			sleepTime := 30*time.Second - timeSinceLastEmail
			time.Sleep(sleepTime)
		}

		// Send the email
		err := sendEmailDirectly(task.To, task.URL, task.Config)
		if err != nil {
			log.Printf("Ошибка при отправке письма: %v", err)
		} else {
			lastSentTime = time.Now()
		}
	}
}

// SendEmail queues an email to be sent with the required delay
func SendEmail(to string, url string, cfg *config.SMTPConfig) error {
	if cfg == nil {
		return fmt.Errorf("SMTP configuration is nil")
	}

	// Ensure the worker is started
	StartEmailWorker()

	// Add the email to the queue
	task := EmailTask{
		To:      to,
		URL:     url,
		Config:  cfg,
		Created: time.Now(),
	}

	select {
	case emailQueue <- task:
		// Successfully queued
		return nil
	default:
		// Queue is full
		return fmt.Errorf("email queue is full, try again later")
	}
}

// sendEmailDirectly actually sends the email
func sendEmailDirectly(to string, url string, cfg *config.SMTPConfig) error {
	subject := "Подтверждение регистрации"

	body := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Подтверждение регистрации</title>
	</head>
	<body style="font-family: 'Helvetica Neue', Arial, sans-serif; margin: 0; padding: 0; background-color: #f4f7fa; color: #333;">
		<div style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1); margin-top: 20px; margin-bottom: 20px;">
			<!-- Header -->
			<div style="background-color: #1a73e8; padding: 24px; text-align: center;">
				<h1 style="color: #ffffff; margin: 0; font-size: 24px; font-weight: 500;">Подтверждение регистрации</h1>
			</div>
			
			<!-- Content -->
			<div style="padding: 30px 40px;">
				<h2 style="color: #1a73e8; font-weight: 500; margin-top: 0;">Добро пожаловать!</h2>
				<p style="font-size: 16px; line-height: 1.6; color: #555; margin-bottom: 25px;">Благодарим вас за регистрацию. Пожалуйста, нажмите на кнопку ниже, чтобы подтвердить вашу учетную запись.</p>
				
				<!-- Button -->
				<div style="text-align: center; margin: 30px 0;">
					<a href="%s" style="display: inline-block; padding: 12px 28px; font-size: 16px; color: #ffffff; background-color: #1a73e8; text-decoration: none; border-radius: 4px; font-weight: 500; box-shadow: 0 2px 5px rgba(0,0,0,0.1); transition: all 0.2s ease-in-out;">Подтвердить регистрацию</a>
				</div>
				
				<p style="font-size: 14px; line-height: 1.6; color: #777; margin-top: 25px;">Если кнопка не работает, вы можете скопировать и вставить следующую ссылку в адресную строку вашего браузера:</p>
				<p style="font-size: 14px; background-color: #f5f5f5; padding: 10px; border-radius: 4px; word-break: break-all; margin-bottom: 0;"><a href="%s" style="color: #1a73e8; text-decoration: none;">%s</a></p>
			</div>
			
			<!-- Footer -->
			<div style="background-color: #f5f7fa; padding: 20px; text-align: center; border-top: 1px solid #e0e0e0;">
				<p style="font-size: 13px; color: #999; margin: 0;">Это автоматическое сообщение, пожалуйста, не отвечайте на него.</p>
				<p style="font-size: 13px; color: #999; margin: 10px 0 0 0;">&copy; 2025 VSRS-RS. Все права защищены.</p>
			</div>
		</div>
	</body>
	</html>`, url, url, url)

	// Формируем сообщение в формате MIME
	message := []byte("From: " + cfg.From + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-Version: 1.0\n" +
		"Content-Type: text/html; charset=UTF-8\n\n" +
		body)

	// Аутентификация
	auth := smtp.PlainAuth("", cfg.From, cfg.Password, cfg.Host)

	// Отправка письма
	return smtp.SendMail(cfg.Host+":"+cfg.Port, auth, cfg.From, []string{to}, message)
}
