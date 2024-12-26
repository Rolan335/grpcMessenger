package serviceErrors

import "errors"

var ErrReadOnlyChat = errors.New("chat is readonly")
var ErrChatNotFound = errors.New("chat not found")
var ErrUserDoesNotExist = errors.New("user doesn't exist")
var ErrProhibited = errors.New("prohibited. Only creator can send")

var ErrInvalidUuid = errors.New("invalid uuid provided")
var ErrNoMessage = errors.New("Cannot send empty messsage")