package api

import (
	"github.com/zhikh23/itsreg-bots/internal/application/dto"
	botsv1 "github.com/zhikh23/itsreg-proto/gen/go/bots/proto"
)

func messageFromDto(msg dto.Message) *botsv1.Message {
	options := make([]*botsv1.Button, len(msg.Options))

	for i, opt := range msg.Options {
		options[i].Text = opt
	}

	return &botsv1.Message{
		Text:    msg.Text,
		Buttons: options,
	}
}

func messagesFromDtos(msgs []dto.Message) []*botsv1.Message {
	res := make([]*botsv1.Message, len(msgs))

	for i, m := range msgs {
		res[i] = messageFromDto(m)
	}

	return res
}
