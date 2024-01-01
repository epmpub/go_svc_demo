package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Message struct {
	OS      string
	Type    string
	Command string
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	// fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

	var message Message
	err := json.Unmarshal(msg.Payload(), &message)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// COMMON:
	// Print Recv JSON
	// log.Printf("Name: %s\n", message.OS)
	// log.Printf("Type: %s\n", message.Type)
	// log.Printf("Command: %s\n", message.Command)

	//TODO:OS DETECT
	log.Println(runtime.GOOS)

	if runtime.GOOS == message.OS {
		if message.Type == "bat" {
			log.Println("win bat cmd recv")
			publish_log(client, get_hostname()+"win bat cmd recv")
		} else if message.Type == "powershell" {
			log.Println("win powershell cmd recv")
		} else if message.Type == "pwsh" {
			log.Println("win pwsh cmd recv")
		} else if message.Type == "powershell_script" {
			log.Println("win powershell_script cmd recv")
		} else if message.Type == "bash" {
			log.Println("linux cmd")
			publish_log(client, get_hostname()+"linux bash cmd recv")
		} else {
			log.Println("Bad command type")
			publish_log(client, get_hostname()+":Bad command type")

		}
	} else {
		log.Println("Bad OS parameter")
		publish_log(client, get_hostname()+"Bad OS parameter")

	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func random_string() string {
	rand.Seed(time.Now().Unix())

	str := "abcdefghijklmnopqrstuvwxyz"

	shuff := []rune(str)

	// Shuffling the string
	rand.Shuffle(len(shuff), func(i, j int) {
		shuff[i], shuff[j] = shuff[j], shuff[i]
	})

	// Displaying the random string
	return string(shuff)
}

func set_work(client mqtt.Client) {
	subscript_opera(client)
	now := time.Now()
	publish_live(client, fmt.Sprintf("Time:%s,Hostname:%s is alive", now.Format("2006-01-02 15:04:05"), get_hostname()))
}

func get_client() mqtt.Client {
	var broker = "broker.emqx.io"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(random_string())
	opts.SetUsername("emqx")
	opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}

func publish_live(client mqtt.Client, msg string) {
	text := fmt.Sprintf("Message %s", msg)
	token := client.Publish("topic/live", 0, false, text)
	token.Wait()
	time.Sleep(time.Second)
}

func publish_log(client mqtt.Client, msg string) {
	text := fmt.Sprintf("Message: %s", msg)
	token := client.Publish("topic/log", 0, false, text)
	token.Wait()
	time.Sleep(time.Second)
}

func subscript_opera(client mqtt.Client) {
	topic := "topic/cmd"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	// fmt.Print(".")
}

func get_hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("get hostname err")
	}
	return hostname + "(" + runtime.GOOS + ")"
}
