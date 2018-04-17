package milterclient

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Negotiation sends SMFIC_OPTNEG
// should be the first Packet after tcp connect
func (session *MilterSession) Negotiation() *Message {
	return &Message{Code: byte(SmficOptNeg), Data: []byte{}}
}

// Macro sends SMFIC_MACRO
func (session *MilterSession) Macro(commandCode byte) *Message {
	data := &bytes.Buffer{}
	data.WriteByte(commandCode)
	for k, v := range session.Macros {
		data.Write(EncodeCString(k))
		data.Write(EncodeCString(v))
	}
	data.WriteByte(0)
	data.WriteByte(0)

	return &Message{Code: byte(SmficMacro), Data: data.Bytes()}
}

// Connect sends SMFIC_CONNECT
func (session *MilterSession) Connect(hostname string, ipv6 bool, ip string) *Message {
	data := &bytes.Buffer{}
	data.Write(EncodeCString(hostname))
	if ipv6 == true {
		data.WriteByte('6')
	} else {
		data.WriteByte('4')
	}
	// write Port number as uint16:
	binary.Write(data, binary.BigEndian, uint16(65535))
	data.Write(EncodeCString(ip))

	return &Message{Code: byte(SmficConnect), Data: data.Bytes()}
}

// MailFrom sends SMFIC_MAIL
func (session *MilterSession) MailFrom(from string) *Message {
	data := &bytes.Buffer{}
	data.Write(EncodeCString(fmt.Sprintf("<%v>", from)))
	data.Write(EncodeCString(""))
	data.WriteByte(0)

	return &Message{Code: byte(SmficMail), Data: data.Bytes()}
}

// RcptTo sends SMFIC_RCPT
func (session *MilterSession) RcptTo(rcpt string) *Message {
	data := &bytes.Buffer{}
	data.Write(EncodeCString(fmt.Sprintf("<%v>", rcpt)))
	data.Write(EncodeCString(""))
	data.WriteByte(0)

	return &Message{Code: byte(SmficRcpt), Data: data.Bytes()}
}

// Header sends SMFIC_HEADER (multiple)
func (session *MilterSession) Header(k, v string) *Message {
	data := &bytes.Buffer{}
	data.Write(EncodeCString(k))
	data.Write(EncodeCString(v))
	data.WriteByte(0)

	return &Message{Code: byte(SmficHeader), Data: data.Bytes()}
}

// EndOfHeader sends SMFIC_EOH
func (session *MilterSession) EndOfHeader() *Message {
	return &Message{Code: byte(SmficEOH), Data: []byte{}}
}

// Body sends SMFIC_BODY (multiple)
func (session *MilterSession) Body(body []byte) (*Message, []byte) {
	if len(body) > MilterChunkSize {
		thisMsg, remainingBody := body[0:MilterChunkSize], body[MilterChunkSize:]
		return &Message{Code: byte(SmficBody), Data: thisMsg}, remainingBody
	}
	return &Message{Code: byte(SmficBody), Data: body}, nil

}

// EndOfBody sends SMFIC_BODYEOB
func (session *MilterSession) EndOfBody() *Message {
	return &Message{Code: byte(SmficBodyEOB), Data: []byte{}}
}

// Quit sends SMFIC_QUIT
func (session *MilterSession) Quit() *Message {
	return &Message{Code: byte(SmficQuit), Data: []byte{}}
}
