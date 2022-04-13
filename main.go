package main

import (
	"crypto/tls"
	"fmt"
	"github.com/michaelklishin/rabbit-hole"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	tlsConfig       *tls.Config
	defaultLogin    = os.Getenv("ADMINUSER")
	defaultPassword = os.Getenv("ADMINPASSWORD")
	loginName       = os.Getenv("LOGINACCOUNT")
	loginPassword   = os.Getenv("LOGINPASSWORD")
	mqEndpoint      = os.Getenv("MQENDPOINT")
)

func handleError(err error, msg string) {
	if err != nil {
		log.Println("%s: %s", msg, err)
	} else {
		fmt.Println("Successfully Connected to our RabbitMQ Instance")
	}
}

// Update user password function
func updateuserpassword() error {
	fmt.Println("Start password updated")
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	amqpServerURL := []string{"https://", mqEndpoint}
	rmqc, _ := rabbithole.NewTLSClient(strings.Join(amqpServerURL, ""), defaultLogin, defaultPassword, transport)
	resp, err := rmqc.PutUser(loginName, rabbithole.UserSettings{Password: loginPassword, Tags: "administrator"})
	handleError(err, "Failed to Update User Password")
	fmt.Println(resp)

	return nil
}

// Try to connect to the RabbitMq instance with the username and password provided
// If the connection is successful, terminate and inform the user that their login is working
// If the login failed, log in to the instance with the administrator's login and update the user's password
func main() {
	amqpServerURL := []string{"amqps://", loginName, ":", loginPassword, "@", mqEndpoint, ":5671/"}
	conn, err := amqp.Dial(strings.Join(amqpServerURL, ""))
	if err == nil {
		log.Println("Password Already Updated")
		return
	} else {
		fmt.Println("Password Need to be update...")
	}
	defer conn.Close()

	updateuserpassword()

	fmt.Println("Password update success")
}
