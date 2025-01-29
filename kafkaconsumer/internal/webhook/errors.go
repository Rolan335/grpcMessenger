package webhook

import "errors"

var ErrWebhookCall = errors.New("failed to call webhook")
var ErrWebhookRetriesExceeded = errors.New("webhook retries exceeded")