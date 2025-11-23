package services

import (
	"insider/internal/sender"
)

const (
	SenderStatusStarted = "started"
	SenderStatusStopped = "stopped"
	SenderStatusError   = "error"
)

type SenderService interface {
	Toggle() (string, error)
}

type SenderServiceImpl struct {
	sender sender.Sender
}

func NewSenderService(sender sender.Sender) *SenderServiceImpl {
	return &SenderServiceImpl{
		sender: sender,
	}
}

func (s *SenderServiceImpl) Toggle() (string, error) {
	if s.sender.IsStarted() {
		s.sender.Stop()
		return SenderStatusStopped, nil
	}

	if err := s.sender.Start(); err != nil {
		return SenderStatusError, err
	}

	return SenderStatusStarted, nil
}
