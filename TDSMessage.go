package main

import (
	"bytes"
	"fmt"
	"reflect"
)

type TDSMessage struct {
	Packets []TDSPacket
}
type AttentionMessage TDSMessage
type RPCRequestMessage TDSMessage
type SQLBatchMessage TDSMessage

func (msg SQLBatchMessage) GetBatchText() string {
	bPayload := TDSMessage(msg).AssemblePayload()
	allHeader := NewAllHeader(bPayload)
	iHeaderLength := int(allHeader.Length())
	return string(bPayload[iHeaderLength : len(bPayload)-iHeaderLength])
}

func newTDSMessage(firstPacket TDSPacket) *TDSMessage {
	msg := new(TDSMessage)
	msg.Packets = append(msg.Packets, firstPacket)
	return msg
}

func (msg TDSMessage) IsComplete() bool {
	if len(msg.Packets) == 0 {
		return false
	}
	return (msg.Packets[len(msg.Packets)-1].header.StatusBitMask() & EndOfMessage) == EndOfMessage
}

func (msg TDSMessage) HasIgnoreBitSet() bool {
	if len(msg.Packets) == 0 {
		return false
	}
	return (msg.Packets[len(msg.Packets)-1].header.StatusBitMask() & IgnoreEvent) == IgnoreEvent
}

func (msg TDSMessage) AssemblePayload() []byte {
	lPayLoad := make([]byte, 4096*4, 4096*4)
	for i := 0; i < len(msg.Packets); i++ {
		lPayLoad = append(lPayLoad, msg.Packets[i].payload...)
	}
	return lPayLoad
}

func CreateFromFirstPacket(firstPacket TDSPacket) TDSMessage {
	var msg TDSMessage
	switch firstPacket.header.Type() {
	case SQLBatch:
		msg = TDSMessage(*new(SQLBatchMessage))
	case AttentionSignal:
		msg = TDSMessage(*new(AttentionMessage))
	case RPC:
		msg = TDSMessage(*new(RPCRequestMessage))
	default:
		msg = *new(TDSMessage)
	}
	msg.Packets = append(msg.Packets, firstPacket)
	return msg
}

func (msg TDSMessage) String() string {
	if msg.IsComplete() {
		var buffer bytes.Buffer
		buffer.WriteString(reflect.TypeOf(msg).String())
		buffer.WriteString(fmt.Sprintf("[#Packets=%d;IsComplete=%t;HasIgnoreBitSet=%t;TotalPayloadSize=%d",
			len(msg.Packets), msg.IsComplete(), msg.HasIgnoreBitSet(), len(msg.AssemblePayload())))
		for i := 0; i < len(msg.Packets); i++ {
			buffer.WriteString(fmt.Sprintf("\n\t[P%d[%s]]", i, msg.Packets[i].String()))
		}
		buffer.WriteString("]")
		return buffer.String()
	}
	return reflect.TypeOf(msg).String() + "{Incomplete message}"
}
