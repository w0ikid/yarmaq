package smtpclient

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	// TLS нужен для реального SMTP (SES, Gmail)
	// для mailhog — false
	UseTLS bool
}

type Client struct {
	cfg Config
}

func New(cfg Config) *Client {
	return &Client{cfg: cfg}
}

type Message struct {
	To      []string
	Subject string
	Body    string
	// true = HTML, false = plain text
	IsHTML bool
}

func (c *Client) Send(msg Message) error {
	addr := fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)

	contentType := "text/plain"
	if msg.IsHTML {
		contentType = "text/html"
	}

	headers := strings.Join([]string{
		"From: " + c.cfg.From,
		"To: " + strings.Join(msg.To, ", "),
		"Subject: " + msg.Subject,
		"MIME-Version: 1.0",
		"Content-Type: " + contentType + "; charset=UTF-8",
		"Date: " + time.Now().Format(time.RFC1123Z),
	}, "\r\n")

	body := headers + "\r\n\r\n" + msg.Body

	if c.cfg.UseTLS {
		return c.sendTLS(addr, msg.To, []byte(body))
	}

	// mailhog без auth, без TLS
	return smtp.SendMail(addr, nil, c.cfg.From, msg.To, []byte(body))
}

// sendTLS for real providers (SES, Gmail) with auth and TLS
func (c *Client) sendTLS(addr string, to []string, body []byte) error {
	tlsCfg := &tls.Config{
		ServerName: c.cfg.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsCfg)
	if err != nil {
		return c.sendSTARTTLS(addr, to, body)
	}
	defer conn.Close()

	host, _, _ := net.SplitHostPort(addr)
	cl, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp new client: %w", err)
	}
	defer cl.Close()

	auth := smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)
	if err = cl.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}

	return c.sendMessage(cl, to, body)
}

func (c *Client) sendSTARTTLS(addr string, to []string, body []byte) error {
	auth := smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)
	return smtp.SendMail(addr, auth, c.cfg.From, to, body)
}

func (c *Client) sendMessage(cl *smtp.Client, to []string, body []byte) error {
	if err := cl.Mail(c.cfg.From); err != nil {
		return fmt.Errorf("smtp MAIL FROM: %w", err)
	}
	for _, addr := range to {
		if err := cl.Rcpt(addr); err != nil {
			return fmt.Errorf("smtp RCPT TO %s: %w", addr, err)
		}
	}
	w, err := cl.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA: %w", err)
	}
	defer w.Close()

	if _, err = w.Write(body); err != nil {
		return fmt.Errorf("smtp write body: %w", err)
	}
	return nil
}
