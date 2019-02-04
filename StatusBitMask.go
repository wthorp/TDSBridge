package main

type StatusBitMask byte

const (
	Normal                  StatusBitMask = 0x00
	EndOfMessage            StatusBitMask = 0x01
	IgnoreEvent             StatusBitMask = 0x02
	MultiPartMessage        StatusBitMask = 0x04
	ResetConnection         StatusBitMask = 0x08
	ResetConnectionSkipTran StatusBitMask = 0x10
)
