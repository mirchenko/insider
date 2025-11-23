package validators

import (
	"insider/internal/model"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func TestValidateMessageStatus_Valid(t *testing.T) {
	v := validator.New()
	if err := v.RegisterValidation("messageStatus", validateMessageStatus); err != nil {
		t.Fatalf("failed to register validation: %v", err)
	}
	type payload struct {
		Status []string `validate:"messageStatus"`
	}
	p := payload{Status: []string{model.MessageStatusSent, model.MessageStatusPending, model.MessageStatusDelivered, model.MessageStatusFailed}}
	if err := v.Struct(p); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestValidateMessageStatus_InvalidValue(t *testing.T) {
	v := validator.New()
	if err := v.RegisterValidation("messageStatus", validateMessageStatus); err != nil {
		t.Fatalf("failed to register validation: %v", err)
	}
	type payload struct {
		Status []string `validate:"messageStatus"`
	}
	p := payload{Status: []string{"unknown"}}
	if err := v.Struct(p); err == nil {
		t.Fatalf("expected validation to fail for invalid value")
	}
}

func TestValidateMessageStatus_NonSlice(t *testing.T) {
	v := validator.New()
	if err := v.RegisterValidation("messageStatus", validateMessageStatus); err != nil {
		t.Fatalf("failed to register validation: %v", err)
	}
	type payload struct {
		Status string `validate:"messageStatus"`
	}
	p := payload{Status: model.MessageStatusSent}
	if err := v.Struct(p); err == nil {
		t.Fatalf("expected validation to fail for non-slice field")
	}
}
func TestRegisterValidators_WithGinBinding(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if err := RegisterValidators(); err != nil {
		t.Fatalf("RegisterValidators error: %v", err)
	}

	r.GET("/test", func(c *gin.Context) {
		var req model.ListMessagesRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(400, gin.H{"err": err.Error()})
			return
		}
		c.Status(200)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test?status=sent&status=failed", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/test?status=bad", nil)
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Fatalf("expected 400 for invalid status, got %d", w.Code)
	}
}
