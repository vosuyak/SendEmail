package main

import (
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type properties struct {
	Username string
	Password string
	Hostname string
	Port     string
}

func init() {
	godotenv.Load(".env")
}

//----------
// Handlers
//----------

func sendEmail(c echo.Context) error {
	tmpEmail := new(email)
	if err := c.Bind(tmpEmail); err != nil {
		return err
	}

	// Set up authentication information.
	auth := smtp.PlainAuth("", os.Getenv("USERNAME"), os.Getenv("PASSWORD"), os.Getenv("HOSTNAME"))

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	response := true
	to := []string{tmpEmail.To}
	msg := []byte("To: " + tmpEmail.To + "\r\n" +
		"Subject: " + tmpEmail.Subject + "!\r\n" +
		"\r\n" +
		tmpEmail.Body + "\r\n")

	err := smtp.SendMail(os.Getenv("HOSTNAME")+":"+os.Getenv("PORT"), auth, os.Getenv("USERNAME"), to, msg)
	if err != nil {
		response = false
		log.Print("ERROR while attempting to send email: ", err)
	}

	return c.JSON(http.StatusCreated, response)
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())

	// Routes
	e.POST("/sendEmail", sendEmail)

	// Start server
	e.Logger.Fatal(e.Start(":9090"))
}
