package messenger

import "errors"

var ErrChatNotFound = errors.New("chat not found")
var ErrUserDoesNotExist = errors.New("user doesn't exist")
var ErrProhibited = errors.New("prohibited. Only creator can send")

var ErrInvalidSessionUUID = errors.New("invalid session UUID provided")
var ErrInvalidChatUUID = errors.New("invalid chat UUID provided")