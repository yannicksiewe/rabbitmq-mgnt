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
	tlsConfig        *tls.Config
	defaultLogin     = os.Getenv("ADMINUSER")
	defaultPassword  = os.Getenv("ADMINPASSWORD")
	loginName        = os.Getenv("LOGINACCOUNT")
	loginPassword    = os.Getenv("LOGINPASSWORD")
	mqEndpoint       = os.Getenv("MQENDPOINT")
	secureConnection = true
)

func handleError(err error, msg string) {
	if err != nil {
		log.Println("%s: %s", msg, err)
	} else {
		fmt.Println("Successfully Connected to our RabbitMQ Instance")
	}
}

func connection() *rabbithole.Client {
	if secureConnection == false {
		amqpServerURL := []string{"http://", mqEndpoint, ":15672"}
		rmqc, _ := rabbithole.NewClient(strings.Join(amqpServerURL, ""), defaultLogin, defaultPassword)

		return rmqc
	} else {
		fmt.Println("Starting with HTTPS")
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		amqpServerURL := []string{"https://", mqEndpoint}
		rmqc, _ := rabbithole.NewTLSClient(strings.Join(amqpServerURL, ""), defaultLogin, defaultPassword, transport)

		return rmqc
	}
}

// Update user password and vhost Setting
func updateuserpassword() error {
	fmt.Println("Start password updated")
	rmqc := connection()
	resp, err := rmqc.PutUser(loginName, rabbithole.UserSettings{Password: loginPassword, Tags: "administrator"})
	handleError(err, "Failed to Update User Password")
	fmt.Println(resp)
	if err != nil {
		return err
	}

	resp2, err2 := rmqc.UpdatePermissionsIn("/", loginName, rabbithole.Permissions{Configure: ".*", Write: ".*", Read: ".*"})
	handleError(err2, "Failed to update permission")
	fmt.Println(resp2)
	if err2 != nil {
		return err2
	}

	return nil
}

// Try to connect to the RabbitMq instance with the username and password provided
// If the connection is successful, terminate and inform the user that their login is working
// If the login failed, log in to the instance with the administrator's login and update the user's password
func main() {
	amqpProtocol := "amqp://"
	if secureConnection == true {
		amqpProtocol = "amqps://"
		fmt.Println("Start with secure amqps protocol")
	}
	amqpServerURL := []string{amqpProtocol, loginName, ":", loginPassword, "@", mqEndpoint, ":5671/"}
	conn, err := amqp.Dial(strings.Join(amqpServerURL, ""))
	if err == nil {
		log.Println("Password Already Updated")

		return
	} else if conn == nil && err != nil {
		fmt.Println(err)
		fmt.Println("Connection Failed")

		return
	} else {
		fmt.Println(err)
		fmt.Println("Password Need to be update...")
	}

	defer conn.Close()

	err = updateuserpassword()
	if err == nil {
		fmt.Println("Password update success")
	}
}
