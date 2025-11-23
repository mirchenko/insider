package services

import (
	"insider/internal/sender"
)

const (
	SenderStatusStarted = "started"
	SenderStatusStopped = "stopped"
	SenderStatusError   = "error"
)

type SenderService struct {
	sender *sender.Sender
}

func NewSenderService(sender *sender.Sender) *SenderService {
	return &SenderService{
		sender: sender,
	}
}

func (s *SenderService) Toggle() (string, error) {
	if s.sender.IsStarted() {
		s.sender.Stop()
		return SenderStatusStopped, nil
	}

	if err := s.sender.Start(); err != nil {
		return "", err
	}

	return SenderStatusStarted, nil
}
