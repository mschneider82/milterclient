package milterclient
    import "github.com/mschneider82/milterclient"

    Package milterclient is an implementation of the milter protocol basicly
    in the role of a MTA. It can be used to Unittest a milter or can be
    modified and implemented in a MTA.

CONSTANTS

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

FUNCTIONS

func DecodeCStrings(data []byte) []string
    DecodeCStrings splits c style strings into golang slice

func EncodeCString(data string) []byte
    EncodeCString encodes a strinc to []byte with ending 0

func GenMtaID(length int) string
    GenMtaID generates an random ID. vocals are removed to prevent dirty
    words which could be negative in spam score

func ReadCString(data []byte) string
    ReadCString reads and returs c style string from []byte

func ReadMessage(r io.Reader) (msg *mail.Message, headers *MIMEHeaderOrdered, err error)
    ReadMessage reads an eml file and parses the header and body

func SendEml(eml io.Reader, milterHostPort, from, to, sendingMtaHostname, sendingMtaIP, mtaMsgID string, ipv6 bool, timeoutSecs int) (byte, error)
    SendEml will send and eml to a Milter empty sendingMtaHostname will be
    mapped to localhost empty sendingMtaIP string will detect IP empty
    mtaMsgID will generate ID

TYPES

type MIMEHeader map[string][]string
    A MIMEHeader represents a MIME-style header mapping keys to sets of
    values.

func (h MIMEHeader) Add(key, value string)
    Add adds the key, value pair to the header. It appends to any existing
    values associated with key.

func (h MIMEHeader) Del(key string)
    Del deletes the values associated with key.

func (h MIMEHeader) Get(key string) string
    Get gets the first value associated with the given key. It is case
    insensitive; CanonicalMIMEHeaderKey is used to canonicalize the provided
    key. If there are no values associated with the key, Get returns "". To
    access multiple values of a key, or to use non-canonical keys, access
    the map directly.

func (h MIMEHeader) Set(key, value string)
    Set sets the header entries associated with key to the single element
    value. It replaces any existing values associated with key.

type MIMEHeaderOrdered struct {
    Keys   []string
    Values []string
}
    MIMEHeaderOrdered stores Header Keys and Values in the correct order
    textprotos MIMEHeader is not ordered.

type Message struct {
    Code byte
    Data []byte
}
    Message represents command sent from milter client

func (m *Message) String() string

type MilterSession struct {
    Sock   io.ReadWriteCloser
    Macros map[string]string
}
    MilterSession keeps session state during MTA communication

func (session *MilterSession) Body(body []byte) (*Message, []byte)
    Body sends SMFIC_BODY (multiple)

func (session *MilterSession) Connect(hostname string, ipv6 bool, ip string) *Message
    Connect sends SMFIC_CONNECT

func (session *MilterSession) EndOfBody() *Message
    EndOfBody sends SMFIC_BODYEOB

func (session *MilterSession) EndOfHeader() *Message
    EndOfHeader sends SMFIC_EOH

func (session *MilterSession) Header(k, v string) *Message
    Header sends SMFIC_HEADER (multiple)

func (session *MilterSession) Macro(commandCode byte) *Message
    Macro sends SMFIC_MACRO

func (session *MilterSession) MailFrom(from string) *Message
    MailFrom sends SMFIC_MAIL

func (session *MilterSession) Negotiation() *Message
    Negotiation sends SMFIC_OPTNEG should be the first Packet after tcp
    connect

func (session *MilterSession) Quit() *Message
    Quit sends SMFIC_QUIT

func (session *MilterSession) RcptTo(rcpt string) *Message
    RcptTo sends SMFIC_RCPT

func (c *MilterSession) ReadPacket() (*Message, error)
    ReadPacket reads incoming milter packet

func (c *MilterSession) ReadResponses(done chan byte) error
    ReadResponses reads all Packets and Print them Also sends Milter Reply
    Packetcode back into a channel

func (c *MilterSession) WriteMessages(messages []*Message, timeoutSecs int, done chan byte) (byte, error)
    WriteMessages writes messages with timeout

func (c *MilterSession) WritePacket(msg *Message) error
    WritePacket sends a milter response packet to socket stream
