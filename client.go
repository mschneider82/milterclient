package milterclient

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

// MilterSession keeps session state during MTA communication
type MilterSession struct {
	Sock   io.ReadWriteCloser
	Macros map[string]string
}

// ReadResponses reads all Packets and Print them
// Also sends Milter Reply Packetcode back into a channel
func (c *MilterSession) ReadResponses(done chan byte) error {
	for {
		m, err := c.ReadPacket()
		if err != nil {
			//fmt.Printf("Error on Packet: %v \n", err)
			return err
		}
		//fmt.Printf("RecvMsg: %v \n", m)
		if m.Code != SmfirAddRcpt &&
			m.Code != SmfirDelRcpt &&
			m.Code != SmfirReplBody &&
			m.Code != SmfirConnFail &&
			m.Code != SmfirAddHeader &&
			m.Code != SmfirInsHeader &&
			m.Code != SmfirChgHeader &&
			m.Code != SmfirProgress {
			done <- m.Code
		}
	}
}

// WriteMessages writes messages with timeout
func (c *MilterSession) WriteMessages(messages []*Message, timeoutSecs int, done chan byte) (byte, error) {
	var err error
	var code byte
	for _, m := range messages {
		// fmt.Printf("Writing: %v\n", m)
		if m.Code != SmficOptNeg {
			err = c.WritePacket(c.Macro(m.Code))
			if err != nil {
				return 0, err
			}
		}

		err = c.WritePacket(m)
		if err != nil {
			// fmt.Printf("Error %v", err)
			return 0, err
		}

		if m.Code != SmficMacro && m.Code != SmficQuit {
		innerloop:
			for {
				select {
				case code = <-done:
					//fmt.Printf("...\n")
					if code == SmfirTempfail {
						//fmt.Printf("Error tempfail\n")
						return code, nil
					}
					break innerloop
				case <-time.After(time.Duration(timeoutSecs) * time.Second):
					return 0, fmt.Errorf("Timeout after %v seconds", timeoutSecs)
				}
			}
		}
	}
	return code, nil
}

// ReadPacket reads incoming milter packet
func (c *MilterSession) ReadPacket() (*Message, error) {
	// read packet length
	var length uint32
	if err := binary.Read(c.Sock, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	// read packet data
	data := make([]byte, length)
	if _, err := io.ReadFull(c.Sock, data); err != nil {
		return nil, err
	}

	// prepare response data
	message := Message{
		Code: data[0],
		Data: data[1:],
	}

	return &message, nil
}

// WritePacket sends a milter response packet to socket stream
func (c *MilterSession) WritePacket(msg *Message) error {
	buffer := bufio.NewWriter(c.Sock)

	// calculate and write response length
	length := uint32(len(msg.Data) + 1)
	if err := binary.Write(buffer, binary.BigEndian, length); err != nil {
		return err
	}

	// write response code
	if err := buffer.WriteByte(msg.Code); err != nil {
		return err
	}

	// write response data
	if _, err := buffer.Write(msg.Data); err != nil {
		return err
	}

	// flush data to network socket stream
	if err := buffer.Flush(); err != nil {
		return err
	}

	return nil
}
