package milterclient

import "fmt"

// Message represents command sent from milter client
type Message struct {
	Code byte
	Data []byte
}

const (
	// MilterVersion we claim to speak
	MilterVersion = 2
	// MilterChunkSize defines how large a SmficBody body can be
	MilterChunkSize = 65535
	// SmficAbort milter command code Abort
	SmficAbort = 'A'
	// SmficBody milter command code Body
	SmficBody = 'B'
	// SmficConnect milter command code Connection information
	SmficConnect = 'C'
	// SmficMacro milter command code Define macro
	SmficMacro = 'D'
	// SmficBodyEOB milter command code final body chunk (End)
	SmficBodyEOB = 'E'
	// SmficHelo milter command code HELO/EHLO
	SmficHelo = 'H'
	// SmficHeader milter command code Header
	SmficHeader = 'L'
	// SmficMail milter command code MAIL from
	SmficMail = 'M'
	// SmficEOH milter command code End Of Header EOH
	SmficEOH = 'N'
	// SmficOptNeg milter command code Option negotation
	SmficOptNeg = 'O'
	// SmficRcpt milter command code RCPT to
	SmficRcpt = 'R'
	// SmficQuit milter command code QUIT
	SmficQuit = 'Q'
	// SmficData milter command code DATA
	SmficData = 'T'
	// SmfifAddHdrs Message Modification
	SmfifAddHdrs = 0x01
	// SmfifChgBody Message Modification
	SmfifChgBody = 0x02
	// SmfifAddRcpt Message Modification
	SmfifAddRcpt = 0x04
	// SmfifDelRcpt Message Modification
	SmfifDelRcpt = 0x08
	// SmfifChgHdrs Message Modification
	SmfifChgHdrs = 0x10
	// SmfifQuarantine Message Handling
	SmfifQuarantine = 0x20

	// SmfiV2Acts A bitmask of all actions supporting in protocol version 2.
	SmfiV2Acts = 0x3f

	// From sendmail's include/libmilter/mfdef.h
	// What the MTA can send/filter wants in protocol

	// SmfipNoConnect MTA should not send connect info
	SmfipNoConnect = 0x01
	// SmfipNoHelo MTA should not send HELO info
	SmfipNoHelo = 0x02
	// SmfipNoMail MTA should not send MAIL info
	SmfipNoMail = 0x04
	// SmfipNoRcpt MTA should not send RCPT info
	SmfipNoRcpt = 0x08
	// SmfipNoBody MTA should not send body
	SmfipNoBody = 0x10
	// SmfipNoHdrs MTA should not send headers
	SmfipNoHdrs = 0x20
	// SmfipNoEoh MTA should not send EOH
	SmfipNoEoh = 0x40

	// SmfiV2Prot A bitmask of all supported protocol steps in protocol version 2.
	SmfiV2Prot = 0x7f

	// Acceptable response commands/codes - actions (replies):

	// SmfirAddRcpt MTA must Add recipient
	SmfirAddRcpt = '+'
	// SmfirDelRcpt MTA must Delete recipient
	SmfirDelRcpt = '-'
	// SmfirAccept MTA must Accept Mail
	SmfirAccept = 'a'
	// SmfirReplBody MTA must replace body (chunk)
	SmfirReplBody = 'b'
	// SmfirContinue MTA must Continue
	SmfirContinue = 'c'
	// SmfirDiscard MTA must discard
	SmfirDiscard = 'd'
	// SmfirConnFail MTA must cause a connection failure
	SmfirConnFail = 'f'
	// SmfirAddHeader MTA must add header
	SmfirAddHeader = 'h'
	// SmfirInsHeader MTA must insert header
	SmfirInsHeader = 'i'
	// SmfirChgHeader MTA must change header
	SmfirChgHeader = 'm'
	// SmfirProgress MTA must progress
	SmfirProgress = 'p'
	// SmfirQuarantine MTA must quarantine
	SmfirQuarantine = 'q'
	// SmfirReject MTA must reject
	SmfirReject = 'r'
	// SmfirTempfail MTA must tempfail
	SmfirTempfail = 't'
	// SmfirReplycode MTA must reply code
	SmfirReplycode = 'y'
)

func (m *Message) String() string {
	return fmt.Sprintf("{%c: %v}", m.Code, string(m.Data))
}
