package main

import (
	"fmt"
)

type AllHeader struct {
	bPayload []byte
}

func (hdr AllHeader) Length() uint32 {
	return uint32(hdr.bPayload[3])*0x01000000 +
		uint32(hdr.bPayload[2])*0x00010000 +
		uint32(hdr.bPayload[1])*0x00000100 +
		uint32(hdr.bPayload[0])*0x00000001
}

func NewAllHeader(bPayload []byte) *AllHeader {
	h := new(AllHeader)
	h.bPayload = bPayload
	return h
}

type TDSHeader struct {
	Buffer []byte
}

type HeaderType int

const (
	SQLBatch                  HeaderType = 1
	PreTD7Login               HeaderType = 2
	RPC                       HeaderType = 3
	TabularResult             HeaderType = 4
	AttentionSignal           HeaderType = 6
	BulkLoadData              HeaderType = 7
	TransactionManagerRequest HeaderType = 14
	TDS7Login                 HeaderType = 16
	SSPIMessage               HeaderType = 17
	PreLoginMessage           HeaderType = 18
	Unknown                   HeaderType = 0xFF
)

const HeaderSize int = 8

func newTDSHeader(bPacket []byte) *TDSHeader {
	hdr := new(TDSHeader)
	hdr.Buffer = make([]byte, HeaderSize, HeaderSize)
	copy(hdr.Buffer, bPacket[:HeaderSize])
	return hdr
}

func (hdr TDSHeader) Type() HeaderType {
	return HeaderType(hdr.Buffer[0])
}

func (hdr TDSHeader) StatusBitMask() StatusBitMask {
	return StatusBitMask(hdr.Buffer[1])
}

func (hdr TDSHeader) LengthIncludingHeader() int {
	return int(hdr.Buffer[2])*0x100 + int(hdr.Buffer[3])
}

func (hdr TDSHeader) PayloadSize() int {
	return hdr.LengthIncludingHeader() - HeaderSize
}

func (hdr TDSHeader) Get(idx int) byte {
	return hdr.Buffer[idx]
}

func (hdr TDSHeader) Set(idx int, value byte) {
	hdr.Buffer[idx] = value
}

func (hdr TDSHeader) String() string {
	return fmt.Sprintf("TDSHeader[Type=%s;StatusBitMask=%d;LengthIncludingHeader=%d;PayloadSize=%d]", hdr.Type(),
		hdr.StatusBitMask(), hdr.LengthIncludingHeader(), hdr.PayloadSize())
}
