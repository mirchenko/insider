package services

import (
	"insider/internal/sender"
)

type SenderService struct {
	sender *sender.Sender
}

func NewSenderService(sender *sender.Sender) *SenderService {
	return &SenderService{
		sender: sender,
	}
}

func (s *SenderService) Start() error {
	return s.sender.Start()
}

func (s *SenderService) Stop() {
	s.sender.Stop()
}
