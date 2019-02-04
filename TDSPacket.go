package main

import (
	"fmt"
)

type TDSPacket struct {
	header  TDSHeader
	payload []byte
}

func newTDSPacket(bBuffer []byte) *TDSPacket {
	pkt := new(TDSPacket)
	pkt.header = *newTDSHeader(bBuffer)
	iPayloadSize := pkt.header.LengthIncludingHeader() - HeaderSize
	pkt.payload = make([]byte, iPayloadSize, iPayloadSize)
	copy(pkt.payload, bBuffer[HeaderSize:])
	return pkt
}

func newTDSPacket2(bHeader []byte, bPayload []byte, iPayloadSize int) *TDSPacket {
	pkt := new(TDSPacket)
	pkt.header = *newTDSHeader(bHeader)
	pkt.payload = make([]byte, iPayloadSize, iPayloadSize)
	copy(pkt.payload, bPayload)
	return pkt
}

func (pkt TDSPacket) String() string {
	return fmt.Sprintf("TDSPacket [Header=%s]", pkt.header)
}
