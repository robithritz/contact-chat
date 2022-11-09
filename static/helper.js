class Helper {
    constructor(client, clientId) {
        this.client = client;
        this.clientId = clientId
    }
    receivedMessageACK(messageId, roomId) {
        this.client.publish(`received/${roomId}`, JSON.stringify({
            messageId: messageId,
            receiverId: this.clientId,
            receivedAt: dayjs().format()
        }), {
            qos: 0
        })

        console.log(`push received/${roomId}`)
    }

    readMessageACK(messageId, roomId) {
        this.client.publish(`read/${roomId}`, JSON.stringify({
            messageId: messageId,
            readerId: this.clientId,
            readAt: dayjs().format()
        }), {
            qos: 0
        })

        console.log(`push read/${roomId}`)
    }

    LMQACK(MQID) {
        this.client.publish(`LMQACK/${MQID}`, JSON.stringify({}), {
            qos: 0
        });
        console.log("push LMQACK ", MQID)
    }

    olahMessage(topic, payload, roomId) {
        console.log(`MQTT - Received message from topic ${topic} | ${payload}`);

        const splitTopic = topic.split('/');
        const message = JSON.parse(payload);

        if (topic.indexOf(`users/${this.clientId}/chats`) >= 0) {
            const lmqId = splitTopic[3];
            helper.LMQACK(lmqId);
            helper.receivedMessageACK(message.messageId, roomId);
            helper.appendChat(message.creatorId, message.creatorUsername, message.messageContent, message.createdAt, message.messageId);

            setTimeout(() => {
                helper.readMessageACK(message.messageId, roomId);
            }, 200);
        } else if (topic.indexOf(`users/${this.clientId}/chatinfo`) >= 0) {
            /**
             * push LMQACK
             * change centang 1 jadi 2 pertanda sent/received message
             */

            const lmqId = splitTopic[3];
            const infoType = message.infoType || '';
            const messageId = message.messageId;

            const messageElement = document.getElementById(messageId);
            if (infoType == "received") {
                messageElement.style.color = 'blue';
            } else if (infoType == "read") {
                messageElement.style.color = 'green';
            }
            helper.LMQACK(lmqId);
        } else if (topic.indexOf(`users/${this.clientId}/new-room`) >= 0) {
            const lmqId = splitTopic[3];

            newRoom(topic, message);
            helper.LMQACK(lmqId);
        }
    }

    appendChat(creatorId, creatorUsername, message, time, messageId) {
        let wrapper = document.createElement('div');
        let newDiv = document.createElement('div');
        let h3 = document.createElement('h3')
        let lbl = document.createElement('label');
        h3.innerText = `${creatorUsername} | ${time}`;
        lbl.innerText = message;
        newDiv.appendChild(h3);
        newDiv.appendChild(lbl);

        wrapper.style.display = 'flex';
        wrapper.style.flexDirection = 'row';
        wrapper.id = messageId;
        if (creatorId == this.clientId) {
            wrapper.style.justifyContent = "flex-end";
        } else {
            wrapper.style.justifyContent = "flex-start";
        }

        wrapper.appendChild(newDiv);
        document.body.appendChild(wrapper);
    }
}