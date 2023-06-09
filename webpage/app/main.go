package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"text/template"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/goombaio/namegenerator"
)

type InstanceStatus int

const (
	Unknown InstanceStatus = iota
	Idle
	Processing
	Killed
)

type TopicPublisher struct {
	Topic *pubsub.Topic
}

type InstanceData struct {
	Name string
}

type InstanceStatusMessage struct {
	Name           string
	RequestCount   int
	InstanceStatus InstanceStatus
}

type HelloData struct {
	RequestCount   int
	ActiveRequests int
	Name           string
	Deleted        bool
}

type Channels struct {
	StartRequest    chan bool
	FinishedRequest chan bool
}

func main() {
	var topicName, projectId string
	if projectId = os.Getenv("PROJECT_ID"); projectId == "" {
		log.Fatalln("No project id given.")
	}
	if topicName = os.Getenv("TOPIC_NAME"); topicName == "" {
		log.Fatalln("No topic name given.")
	}

	topic := setupPubsub(projectId, topicName)
	tp := &TopicPublisher{
		Topic: topic,
	}
	h := &HelloData{
		RequestCount:   0,
		ActiveRequests: 0,
		Deleted:        false,
		Name:           randName(),
	}

	fmt.Println("Starting instance ..")
	tp.writeMessage(h)
	channels := initChannels(h, tp)
	handleRequests(h.Name, channels)
	servePages()

}

func initChannels(h *HelloData, tp *TopicPublisher) Channels {
	terminateSignal := make(chan os.Signal, 1)
	channels := Channels{
		StartRequest:    make(chan bool),
		FinishedRequest: make(chan bool),
	}

	signal.Notify(terminateSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	go handleChannels(h, tp, channels, terminateSignal)
	return channels
}

func handleChannels(h *HelloData, tp *TopicPublisher, channels Channels, terminate <-chan os.Signal) {
	intervalString := os.Getenv("MESSAGE_INTERVAL")
	var interval int
	var err error
	if interval, err = strconv.Atoi(intervalString); err != nil {
		interval = 1
	}
	messageTimer := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-messageTimer.C:
			tp.writeMessage(h)

		case <-channels.StartRequest:
			h.ActiveRequests++
			tp.writeMessage(h)

		case <-channels.FinishedRequest:
			h.RequestCount++
			h.ActiveRequests--

		case <-terminate:
			h.Deleted = true
			tp.writeMessage(h)
			fmt.Println(" Exiting...", h)
			os.Exit(0)
		}
	}
}

func handleRequests(name string, channels Channels) {
	http.HandleFunc("/health", healthCheck())
	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("./images"))))
	tmpl := template.Must(template.ParseFiles("index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := InstanceData{
			Name: name,
		}
		channels.StartRequest <- true
		if err := tmpl.Execute(w, data); err != nil {
			log.Fatalln("Failed template execution %v", err)
		}
		channels.FinishedRequest <- true
	})
}

func setupPubsub(projectId string, topicName string) *pubsub.Topic {
	var client *pubsub.Client
	var err error
	ctx := context.Background()
	if client, err = pubsub.NewClient(ctx, projectId); err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
	}
	return client.Topic(topicName)
}

func (tp *TopicPublisher) publish(msg string) error {
	ctx := context.Background()
	result := tp.Topic.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	_, err := result.Get(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (tp *TopicPublisher) writeMessage(data *HelloData) {

	message := &InstanceStatusMessage{
		RequestCount: data.RequestCount,
		Name:         data.Name,
	}

	if data.Deleted {
		message.InstanceStatus = Killed
	} else if data.ActiveRequests > 0 {
		message.InstanceStatus = Processing
	} else {
		message.InstanceStatus = Idle
	}

	b, _ := json.Marshal(message)
	log.Println("Publish: ", string(b))
	if err := tp.publish(string(b)); err != nil {
		log.Fatalf("Failed to publish: %v", err)
	}
}

func healthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func servePages() {
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("Receive execute listen and serve failed %v", err)
	}
}

func randName() string {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	name := nameGenerator.Generate()
	return name
}
