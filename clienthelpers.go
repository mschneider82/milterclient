package milterclient

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"time"
)

// GenMtaID generates an random ID. vocals are removed to prevent dirty words which could be negative in spam score
func GenMtaID(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("bcdfghjklmnpqrstvwxyzBCDFGHJKLMNPQRSTVWXYZ")

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// SendEml will send and eml to a Milter
// empty sendingMtaHostname will be mapped to localhost
// empty sendingMtaIP string will detect IP
// empty mtaMsgID will generate ID
func SendEml(eml io.Reader, milterHostPort, from, to, sendingMtaHostname, sendingMtaIP, mtaMsgID string, ipv6 bool, timeoutSecs int) (byte, error) {

	msg, headers, _ := ReadMessage(eml)

	body, _ := ioutil.ReadAll(msg.Body)

	done := make(chan byte)

	conn, err := net.Dial("tcp", milterHostPort)

	if err != nil {
		return 0, err
	}

	if sendingMtaHostname == "" {
		sendingMtaHostname = "localhost"
	}

	if sendingMtaIP == "" {
		sendingMtaIP = conn.LocalAddr().(*net.TCPAddr).IP.String()
	}

	if mtaMsgID == "" {
		mtaMsgID = GenMtaID(12)
	}
	//fmt.Printf("MessageId: %s\n", mtaMsgID)

	Session := &MilterSession{Sock: conn, Macros: map[string]string{"i": mtaMsgID}}

	go func() {
		err1 := Session.ReadResponses(done)
		if err1 != nil {
			fmt.Printf("Error Reading: %v\n", err1)
		}
	}()

	messages := []*Message{
		Session.Negotiation(),
		Session.Connect(sendingMtaHostname, ipv6, sendingMtaIP),
		Session.MailFrom(from),
		Session.RcptTo(to),
	}

	for i, key := range headers.Keys {
		value := headers.Values[i]
		messages = append(messages, Session.Header(key, value))
	}

	messages = append(messages, Session.EndOfHeader())

	var remainingBody = body
	var m *Message
	for remainingBody != nil {
		//fmt.Printf("remainingBody len: %v\n", len(remainingBody))
		m, remainingBody = Session.Body(remainingBody)
		messages = append(messages, m)
	}
	messages = append(messages, Session.EndOfBody())
	//messages = append(messages, Session.Quit())
	var lastCode byte
	lastCode, err = Session.WriteMessages(messages, timeoutSecs, done)

	return lastCode, err
}
