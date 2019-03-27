package network

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

// --------------
// | head | protobuf_data |
// --------------

type HeadUnit struct {
	From         uint8
	CodecType    uint8
	MsgLen       uint32
	MsgID        uint32
	Version      uint16
	ClientExtral uint32
	Extral       uint64
}

// func (p *HeadUnit) String() string {
// 	return fmt.Sprintf()
// }

func (p *HeadUnit) Encode() ([]byte, error) {
	nBytesOut := make([]byte, 24)
	nBytesOut[0] = p.From
	nBytesOut[1] = p.CodecType
	binary.LittleEndian.PutUint32(nBytesOut[2:6], p.MsgLen)
	binary.LittleEndian.PutUint32(nBytesOut[6:10], p.MsgID)
	binary.LittleEndian.PutUint16(nBytesOut[10:12], p.Version)
	binary.LittleEndian.PutUint32(nBytesOut[12:16], p.ClientExtral)
	binary.LittleEndian.PutUint64(nBytesOut[16:24], p.Extral)
	return nBytesOut, nil
}

func (p *HeadUnit) Decode(nByteIn []byte) {
	p.From = nByteIn[0]
	p.CodecType = nByteIn[1]
	p.MsgLen = binary.LittleEndian.Uint32(nByteIn[2:6])
	p.MsgID = binary.LittleEndian.Uint32(nByteIn[6:10])
	p.Version = binary.LittleEndian.Uint16(nByteIn[10:12])
	p.ClientExtral = binary.LittleEndian.Uint32(nByteIn[12:16])
	p.Extral = binary.LittleEndian.Uint64(nByteIn[16:24])
}

type MsgParser struct {
	lenMsgLen    int
	minMsgLen    uint32
	maxMsgLen    uint32
	littleEndian bool
}

func NewLiveMsgParser() *MsgParser {
	p := new(MsgParser)
	p.lenMsgLen = 2
	p.minMsgLen = 1
	p.maxMsgLen = 10240
	p.littleEndian = true
	return p
}

// It's dangerous to call the method on reading or writing
func (p *MsgParser) SetMsgLen(lenMsgLen int, minMsgLen uint32, maxMsgLen uint32) {
	if lenMsgLen == 1 || lenMsgLen == 2 || lenMsgLen == 4 {
		p.lenMsgLen = lenMsgLen
	}
	if minMsgLen != 0 {
		p.minMsgLen = minMsgLen
	}
	if maxMsgLen != 0 {
		p.maxMsgLen = maxMsgLen
	}

	var max uint32
	switch p.lenMsgLen {
	case 1:
		max = math.MaxUint8
	case 2:
		max = math.MaxUint16
	case 4:
		max = math.MaxUint32
	}
	if p.minMsgLen > max {
		p.minMsgLen = max
	}
	if p.maxMsgLen > max {
		p.maxMsgLen = max
	}
}

// It's dangerous to call the method on reading or writing
func (p *MsgParser) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

// goroutine safe
func (p *MsgParser) Read(conn *TCPConn) ([]byte, error) {

	bufMsgHead := make([]byte, 24)
	if _, err := io.ReadFull(conn, bufMsgHead); err != nil {
		return nil, err
	}

	nUnit := &HeadUnit{}
	nUnit.Decode(bufMsgHead)
	// parse len
	var msgLen uint32 = nUnit.MsgLen
	// check len
	if msgLen > math.MaxUint32 {
		return nil, errors.New("message too long")
	} else if msgLen < 0 {
		return nil, errors.New("message too short")
	}

	msgData := make([]byte, msgLen+24)
	if msgLen > 0 {
		// data
		if _, err := io.ReadFull(conn, msgData[24:]); err != nil {
			return nil, err
		}
	}
	copy(msgData[0:24], bufMsgHead)
	return msgData, nil
}

// goroutine safe
func (p *MsgParser) Write(conn *TCPConn, args ...[]byte) error {
	// get len
	var msgLen uint32
	for i := 0; i < len(args); i++ {
		msgLen += uint32(len(args[i]))
	}

	// check len
	if msgLen > math.MaxUint32 {
		return errors.New("message too long")
	} else if msgLen < 0 {
		return errors.New("message too short")
	}

	msg := make([]byte, msgLen)
	// write data
	l := 0
	for i := 0; i < len(args); i++ {
		copy(msg[l:], args[i])
		l += len(args[i])
	}
	conn.Write(msg)
	return nil
}
