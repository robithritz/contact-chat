const clientId = 1; // robith harus nya sih login dulu yaa
const client = mqtt.connect('ws://test.mosquitto.org:8080', {
    clientId: clientId,
    clean: false,
    properties: {
        sessionExpiryInterval: 7600
    }
}) // you add a ws:// url here

const clientUsername = "robith";
const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InJvYml0aCIsImZpcnN0TmFtZSI6IlJvYml0aCIsImxhc3ROYW1lIjoiUml0eiIsImltYWdlVVJMIjoiIiwicGhvbmVOdW1iZXIiOiIifQ.UiEsOpdGtu82-_fmEBOAGYjT48Oib2EXemswL6ywXcw";
const btnStartMQTT = document.getElementById('startMQTT');
btnStartMQTT.addEventListener('click', startMQTT);
const pushButton = document.getElementById('pushMessage');
const inputEl = document.getElementById('messageInput');
let roomId = "";

const helper = new Helper(client, clientId);
var startTime = performance.now();

client.on('connect', function(conack) {
    console.log("MQTT connected");

    /**
     * Subscribe to user own topic
     * 
     */

     client.subscribe(`users/${clientId}/#`, {
        qos: 2
    }, function(err) {
        if(!err) {
            console.log(`Subscribed to topic users/${clientId}/#`);
        }
    });
});

client.on('disconnect', function() {
    console.log("MQTT Disconnected!!");
})

client.on('message', function(topic, payload) {
    helper.olahMessage(topic, payload, roomId);
});

// called when the client loses its connection
function onConnectionLost(responseObject) {
    if (responseObject.errorCode !== 0) {
        console.log("onConnectionLost:"+responseObject.errorMessage);
      }else{
        console.log("connection lost");
      }
}

inputEl.addEventListener('keydown', function(ev) {
    if(ev.key == "Enter") {
        pushButton.click()
    }
})

pushButton.addEventListener('click', function() {
    
    const isi = inputEl.value;

    const chatData = {
        messageId: uuid.v1(),
        messageType: "text",
        messageContent: isi,
        creatorId: clientId,
        targets: `6`, // ambil dari room detail nanti nya, room detail isi nya ada list participants, dan sender client id gk usah dikirim
        creatorUsername: clientUsername,
        createdAt: dayjs().format()
    }
    helper.appendChat(clientId, clientUsername, chatData.messageContent, chatData.createdAt, chatData.messageId);

    client.publish(`rooms/${roomId}/chats`, JSON.stringify(chatData), {
        qos: 0
    })

    inputEl.value = "";
});

async function startMQTT() {
    // user 1 robith
    let result = await fetch('http://localhost:8888/rooms', {
        method: "POST",
        headers: {
            'Authorization': token,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            roomType: 'private',
            participants: [6, clientId]
        })
    })
    
    const finalResult = await result.json()
    console.log(finalResult)
    roomId = finalResult['roomId'];
    console.log(`Room ${roomId} created`);

    const chatData = {
        messageId: uuid.v1(),
        messageType: "text",
        messageContent: "halo semua!!!",
        creatorId: clientId,
        targets: `6`, // ambil dari room detail nanti nya, room detail isi nya ada list participants, dan sender client id gk usah dikirim
        creatorUsername: clientUsername,
        createdAt: dayjs().format()
    }
    helper.appendChat(clientId, clientUsername, chatData.messageContent, chatData.createdAt, chatData.messageId);

    client.publish(`rooms/${roomId}/chats`, JSON.stringify(chatData), {
        qos: 0
    })
}

async function newRoom(topic, data) {
    console.log(`new room alert ${data.roomId}`);

    roomId = data.roomId;
}