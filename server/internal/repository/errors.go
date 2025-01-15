package repository

import "errors"

var ErrNotFound = errors.New("not found")
var ErrProhibited = errors.New("prohibited. Only creator can send")
var ErrUserDoesntExist = errors.New("user doesn't exist")