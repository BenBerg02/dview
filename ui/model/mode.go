package model

type Mode int

const (
	ModeNone Mode = iota
	ModeEdit
	ModeMenu
	ModeStructure
	ModeSearch
	ModeConfirm
	ModeInsert
)
