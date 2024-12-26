package checkuuid

import (
	"github.com/google/uuid"
)

func IsParsed(inUuid ...string) bool {
	for _, v := range inUuid {
		_, err := uuid.Parse(v)
		if err != nil {
			return false
		}
	}
	return true
}
