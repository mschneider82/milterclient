package milterclient

import (
	"bytes"
	"fmt"
	"github.com/phalaaxx/milter"
	"log"
	"net"
	"net/textproto"
	"os"
	"strings"
	"testing"
)

/* ExtMilter object */
type ExtMilter struct {
	milter.Milter
	multipart bool
	message   *bytes.Buffer
}

// https://github.com/phalaaxx/milter/blob/master/interface.go

func (e *ExtMilter) Connect(name, value string, port uint16, ip net.IP, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

func (e *ExtMilter) MailFrom(name string, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

func (e *ExtMilter) RcptTo(name string, m *milter.Modifier) (milter.Response, error) {
	return milter.RespContinue, nil
}

/* handle headers one by one */
func (e *ExtMilter) Header(name, value string, m *milter.Modifier) (milter.Response, error) {
	// if message has multiple parts set processing flag to true
	if name == "Content-Type" && strings.HasPrefix(value, "multipart/") {
		e.multipart = true
	}
	return milter.RespContinue, nil
}

/* at end of headers initialize message buffer and add headers to it */
func (e *ExtMilter) Headers(headers textproto.MIMEHeader, m *milter.Modifier) (milter.Response, error) {
	// return accept if not a multipart message
	if !e.multipart {
		return milter.RespAccept, nil
	}
	// prepare message buffer
	e.message = new(bytes.Buffer)
	// print headers to message buffer
	for k, vl := range headers {
		for _, v := range vl {
			if _, err := fmt.Fprintf(e.message, "%s: %s\n", k, v); err != nil {
				return nil, err
			}
		}
	}
	if _, err := fmt.Fprintf(e.message, "\n"); err != nil {
		return nil, err
	}
	// continue with milter processing
	return milter.RespContinue, nil
}

// accept body chunk
func (e *ExtMilter) BodyChunk(chunk []byte, m *milter.Modifier) (milter.Response, error) {
	// save chunk to buffer
	if _, err := e.message.Write(chunk); err != nil {
		return nil, err
	}
	return milter.RespContinue, nil
}

/* Body is called when email message body has been sent */
func (e *ExtMilter) Body(m *milter.Modifier) (milter.Response, error) {
	// prepare buffer
	_ = bytes.NewReader(e.message.Bytes())
	// accept message by default
	return milter.RespAccept, nil
}

/* RunServer creates new Milter instance */
func RunServer(socket net.Listener) {
	// declare milter init function
	init := func() (milter.Milter, uint32, uint32) {
		return &ExtMilter{},
			milter.OptAddHeader | milter.OptChangeHeader,
			milter.OptNoRcptTo
	}
	// start server
	_ = milter.RunServer(socket, init)
}

/* main program */
func TestMilterClient(t *testing.T) {

	// parse commandline arguments
	protocol := "tcp"
	address := "127.0.0.1:12349"

	// bind to listening address
	socket, err := net.Listen(protocol, address)
	if err != nil {
		log.Fatal(err)
	}
	//defer socket.Close()

	// run server
	go RunServer(socket)

	// run tests:
	emlFilePath := "testmail.eml"
	eml, err := os.Open(emlFilePath)
	if err != nil {
		t.Errorf("Error opening test eml file %v: %v", emlFilePath, err)
	}
	defer eml.Close()

	msgID := GenMtaID(12)
	last, err := SendEml(eml, "127.0.0.1:12349", "from@unittest.de", "to@unittest.de", "", "", msgID, false, 5)
	if err != nil {
		t.Errorf("Error sending eml to milter: %v", err)
	}

	fmt.Printf("MsgId: %s, Lastmilter code: %s\n", msgID, string(last))
	if last != SmfirAccept {
		t.Errorf("Excepted Accept from Milter, got %v", last)
	}
	socket.Close()
}
