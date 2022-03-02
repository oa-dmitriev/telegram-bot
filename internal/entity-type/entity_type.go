package entitytype

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type EntityType interface {
	EntityType() string
}

type entityType string

func (et entityType) EntityType() string {
	return string(et)
}

const (
	Invalid  entityType = "invalid"
	Command  entityType = "command"
	Message  entityType = "message"
	CallBack entityType = "callback"
)

func New(u *tgbotapi.Update) EntityType {
	if u == nil || u.Message == nil && u.CallbackQuery == nil {
		return Invalid
	}
	if u.CallbackQuery != nil {
		return CallBack
	}
	if u.Message.IsCommand() {
		return Command
	}
	return Message
}

func isEqual(first, second EntityType) bool {
	if first == nil || second == nil {
		return first == second
	}
	return first.EntityType() == second.EntityType()
}

func IsInvalid(et EntityType) bool {
	return isEqual(et, Invalid)
}

func IsCommand(et EntityType) bool {
	return isEqual(et, Command)
}

func IsMessage(et EntityType) bool {
	return isEqual(et, Message)
}

func IsCallback(et EntityType) bool {
	return isEqual(et, CallBack)
}
