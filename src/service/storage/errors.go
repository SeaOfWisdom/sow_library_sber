package storage

import "errors"

var (
	ErrParticipantNotExists     = errors.New("participant does not exist")
	ErrParticipantAlreadyExists = errors.New("participant already exists")
	ErrSomethingWentWrong       = errors.New("something went wrong")
	ErrWorkNotExists            = errors.New("work does not exist")
)
