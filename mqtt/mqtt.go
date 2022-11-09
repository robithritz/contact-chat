package mqtt

import (
	"contact-chat/database"
	"contact-chat/models/chats"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var roomWildcardTopic = "rooms/#"
var readWildcardTopic = "read/#"
var receivedWildcardTopic = "received/#"
var userTopic = "users/"
var LMQACKTopic = "LMQACK/#"

var client mqtt.Client

type RoomNewMessagePayload struct {
	MessageId       uuid.UUID `json:"messageId"`
	MessageType     string    `json:"messageType"`
	MessageContent  string    `json:"messageContent"`
	Targets         string    `json:"targets"`
	CreatorID       uint8     `json:"creatorId"`
	CreatorUsername string    `json:"creatorUsername"`
	CreatedAt       time.Time `json:"createdAt"`
}

type RoomTreatedMessagePayload struct {
	RoomId uint8 `json:"roomId"`
	RoomNewMessagePayload
	CreatedAt time.Time `json:"createdAt"`
}

type ReadMessagePayload struct {
	MessageId uuid.UUID `json:"messageId"`
	ReaderId  uint8     `json:"readerId"`
	ReadAt    time.Time `json:"readAt"`
}

type ReceivedMessagePayload struct {
	MessageId  uuid.UUID `json:"messageId"`
	ReceiverId uint8     `json:"receiverId"`
	ReceivedAt time.Time `json:"receivedAt"`
}

var messagePubHandler mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {

	fmt.Printf("MQTT - Received message : %s , from topic : %s \n", m.Payload(), m.Topic())

}

var connectedHandler mqtt.OnConnectHandler = func(c mqtt.Client) {
	fmt.Println("MQTT - Client connected to broker")

	client.Subscribe(roomWildcardTopic, 0, roomsWildcardHandler)
	client.Subscribe(readWildcardTopic, 0, readWildcardHandler)
	client.Subscribe(receivedWildcardTopic, 0, receivedWildcardHandler)
	client.Subscribe(LMQACKTopic, 0, LMQACKHandler)
}

var lostConnectionHandler mqtt.ConnectionLostHandler = func(c mqtt.Client, err error) {
	fmt.Printf("MQTT - Connection lost %s \n", err.Error())
}

var readWildcardHandler mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
	start := time.Now()

	topicSlice := strings.Split(m.Topic(), "/")
	fmt.Printf("MQTT - Received message : %s | topic : %s \n", m.Payload(), m.Topic())

	roomId := topicSlice[1]

	data := ReadMessagePayload{}
	if err := json.Unmarshal(m.Payload(), &data); err != nil {
		fmt.Printf("Error %s", err.Error())
	}

	fmt.Printf("MQTT - room %s read message ACK %s from %d \n", roomId, data.MessageId, data.ReaderId)

	saveRead := &database.TChatReaders{
		MessageId: data.MessageId,
		TargetID:  uint(data.ReaderId),
		ReadAt:    data.ReadAt,
	}

	err := chats.ReadMessage(saveRead)
	if err != nil {
		return
	}

	var recordFound database.TChats

	result := database.DB.First(&recordFound, "id = ?", data.MessageId)
	if result.RowsAffected == 1 {
		senderId := strconv.Itoa(int(recordFound.SenderID))
		newMQID := uuid.New()

		payload := echo.Map{
			"messageId": data.MessageId,
			"infoType":  "read",
			"target":    data.ReaderId,
			"datetime":  data.ReadAt,
		}

		payloadStringified, _ := json.Marshal(payload)
		targetTopic := "users/" + senderId + "/chatinfo/" + newMQID.String()

		logMQ := &database.LMQ{
			MQID:        newMQID,
			TargetTopic: targetTopic,
			TargetID:    recordFound.SenderID,
			Payload:     string(payloadStringified),
		}
		err := chats.SaveMQLog(logMQ)
		if err != nil {
			return
		}

		Publish(targetTopic, payload, 0)

		fmt.Printf("send read status to message owner | %s \n", time.Since(start))
	}

}

var receivedWildcardHandler mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
	start := time.Now()

	topicSlice := strings.Split(m.Topic(), "/")
	fmt.Printf("MQTT - Received message : %s | topic : %s \n", m.Payload(), m.Topic())

	roomId := topicSlice[1]

	data := ReceivedMessagePayload{}
	if err := json.Unmarshal(m.Payload(), &data); err != nil {
		fmt.Printf("Error %s", err.Error())
	}

	fmt.Printf("MQTT - room %s received message ACK %s from %d \n", roomId, data.MessageId, data.ReceiverId)

	saveSent := &database.TChatSents{
		MessageId: data.MessageId,
		TargetID:  uint(data.ReceiverId),
		SentAt:    data.ReceivedAt,
	}

	err := chats.ReceivedMessage(saveSent)
	if err != nil {
		return
	}

	var recordFound database.TChats

	result := database.DB.First(&recordFound, "id = ?", data.MessageId)
	if result.RowsAffected == 1 {
		senderId := strconv.Itoa(int(recordFound.SenderID))
		newMQID := uuid.New()

		payload := echo.Map{
			"messageId": data.MessageId,
			"infoType":  "received",
			"target":    data.ReceiverId,
			"datetime":  data.ReceivedAt,
		}

		payloadStringified, _ := json.Marshal(payload)
		targetTopic := "users/" + senderId + "/chatinfo/" + newMQID.String()

		logMQ := &database.LMQ{
			MQID:        newMQID,
			TargetTopic: targetTopic,
			TargetID:    recordFound.SenderID,
			Payload:     string(payloadStringified),
		}
		err := chats.SaveMQLog(logMQ)
		if err != nil {
			return
		}

		Publish(targetTopic, payload, 0)

		fmt.Printf("send received status to message owner | %s \n", time.Since(start))
	}

}

var roomsWildcardHandler mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
	topicSlice := strings.Split(m.Topic(), "/")
	fmt.Printf("MQTT - Received message : %s | topic : %s  \n", m.Payload(), m.Topic())

	roomId, err := strconv.Atoi(topicSlice[1])
	if err != nil {
		return
	}
	if len(topicSlice) == 3 {
		start := time.Now()
		// user send new message to be treated
		data := RoomNewMessagePayload{}
		if err := json.Unmarshal(m.Payload(), &data); err != nil {
			fmt.Printf("Error %s", err.Error())
		}

		fmt.Printf("%s \n", time.Since(start))
		stringified, _ := json.Marshal(data)
		fmt.Printf("MQTT - New message to room %d | %s \n", roomId, stringified)

		message := &database.TChats{
			ID:             data.MessageId,
			RoomID:         uint(roomId),
			SenderID:       uint(data.CreatorID),
			MessageType:    data.MessageType,
			MessageContent: data.MessageContent,
		}

		err = chats.SaveMessage(message)
		if err != nil {
			return
		}

		fmt.Printf("Saving chat | %s \n", time.Since(start))

		listTargets := strings.Split(data.Targets, ",")

		var logMQs []database.LMQ

		var publishTargetTopics []string

		for _, targetId := range listTargets {
			// Saving log MQ Messages to determined received and read

			targetIdInt, _ := strconv.Atoi(targetId)
			newMQID := uuid.New()
			targetTopic := "users/" + targetId + "/chats/" + newMQID.String()

			logMQ := &database.LMQ{
				MQID:        newMQID,
				TargetTopic: targetTopic,
				TargetID:    uint(targetIdInt),
				Payload:     string(m.Payload()),
			}

			logMQs = append(logMQs, *logMQ)

			publishTargetTopics = append(publishTargetTopics, targetTopic)
		}

		err := chats.SaveMQLogBulk(logMQs)
		if err != nil {
			return
		}

		fmt.Printf("saving bulk log | %s \n", time.Since(start))

		for idx, targetId := range listTargets {
			// Sending to all targets's own topic
			Publish(publishTargetTopics[idx], data, 1)
			fmt.Printf("published to topic %s message %s to %s | %s \n", publishTargetTopics[idx], data.MessageId, targetId, time.Since(start))
		}

	}

}

var LMQACKHandler mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
	topicSlice := strings.Split(m.Topic(), "/")
	fmt.Printf("MQTT - LMQACK %s | topic : %s \n", m.Payload(), m.Topic())

	lmqId := topicSlice[1]

	data := ReceivedMessagePayload{}
	if err := json.Unmarshal(m.Payload(), &data); err != nil {
		fmt.Printf("Error %s", err.Error())
	}

	database.DB.Where("mq_id = ?", lmqId).Delete(&database.LMQ{})

	fmt.Printf("LMQ Deleted Successfully %s \n", lmqId)

}

func ConnectBroker() {
	var broker = "test.mosquitto.org"
	var port = 8080
	options := mqtt.NewClientOptions()
	options.AddBroker(fmt.Sprintf("ws://%s:%d", broker, port))
	options.SetClientID("contact-chat-server-1")
	options.SetConnectionLostHandler(lostConnectionHandler)
	options.SetOnConnectHandler(connectedHandler)
	client = mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func Subscribe(topic string) {
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("MQTT - Subscribed to topic %s", topic)
}

func Publish(topic string, payload any, qos int) {

	marshaled, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("jzon invalid")
	}
	client.Publish(topic, byte(qos), false, marshaled)
	// token.Wait()
	// if token.Error() != nil {
	// fmt.Printf("Publish Error %s", token.Error())
	// }

}
