package validators

import (
	"insider/internal/model"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var messageStatuses = map[string]struct{}{
	model.MessageStatusSent:      {},
	model.MessageStatusPending:   {},
	model.MessageStatusDelivered: {},
	model.MessageStatusFailed:    {},
}

// TODO: make better
func validateMessageStatus(fl validator.FieldLevel) bool {
	statuses := fl.Field()
	if statuses.Kind() != reflect.Slice {
		return false
	}

	for _, v := range statuses.Seq2() {
		_, exists := messageStatuses[v.String()]
		if !exists {
			return false
		}
	}

	return true
}
