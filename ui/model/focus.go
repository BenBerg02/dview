package model

type Focus int

const (
	FocusTables Focus = iota
	FocusData
	FocusContext
	FocusGlobal
)
