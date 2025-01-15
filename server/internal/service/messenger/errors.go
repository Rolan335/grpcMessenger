package messenger

import "errors"

var ErrChatNotFound = errors.New("chat not found")
var ErrUserDoesNotExist = errors.New("user doesn't exist")
var ErrProhibited = errors.New("prohibited. Only creator can send")

var ErrInvalidSessionUuid = errors.New("invalid session uuid provided")
var ErrInvalidChatUuid = errors.New("invalid chat uuid provided")