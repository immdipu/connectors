package atlassian

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/amp-labs/connectors/common"
	"github.com/amp-labs/connectors/common/interpreter"
)

func (*Connector) interpretJSONError(res *http.Response, body []byte) error { //nolint:cyclop
	payload := make(map[string]any)
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("interpretJSONError general: %w %w", interpreter.ErrUnmarshal, err)
	}

	// now we can choose which error response Schema we expect
	var schema common.ErrorDescriptor

	if _, ok := payload["status"]; ok {
		apiError := &ResponseStatusError{}
		if err := json.Unmarshal(body, &apiError); err != nil {
			return fmt.Errorf("interpretJSONError SingleError: %w %w", interpreter.ErrUnmarshal, err)
		}

		schema = apiError
	} else {
		apiError := &ResponseMessagesError{}
		if err := json.Unmarshal(body, &apiError); err != nil {
			return fmt.Errorf("interpretJSONError MessagesError: %w %w", interpreter.ErrUnmarshal, err)
		}

		schema = apiError
	}

	// enhance status code error with response payload
	return schema.CombineErr(interpreter.DefaultStatusCodeMappingToErr(res, body))
}

type ResponseMessagesError struct {
	ErrorMessages   []string          `json:"errorMessages"`
	WarningMessages []string          `json:"warningMessages"`
	Errors          map[string]string `json:"errors"`
}

// CombineErr will produce dynamic error from server response body.
// The base error serves as a main static error on top of stacked errors.
// That static error should be used in conditional decisions. Ex: common.ErrBadRequest.
func (r ResponseMessagesError) CombineErr(base error) error {
	result := base

	if len(r.ErrorMessages) != 0 {
		result = errors.Join(result, errors.New( // nolint:goerr113
			strings.Join(r.ErrorMessages, ","),
		))
	}

	if len(r.WarningMessages) != 0 {
		result = errors.Join(result, errors.New( // nolint:goerr113
			strings.Join(r.WarningMessages, ","),
		))
	}

	if len(r.Errors) != 0 {
		messages := make([]string, 0)
		for k, v := range r.Errors {
			messages = append(messages, fmt.Sprintf("%v:%v", k, v))
		}

		result = errors.Join(result, errors.New( // nolint:goerr113
			strings.Join(messages, ","),
		))
	}

	return result
}

type ResponseStatusError struct {
	Status    int       `json:"status"`
	Error     string    `json:"error"`
	Message   string    `json:"message"`
	Path      string    `json:"path"`
	Timestamp time.Time `json:"timestamp"`
}

func (r ResponseStatusError) CombineErr(base error) error {
	if len(r.Error) == 0 {
		return base
	}

	if len(r.Message) == 0 {
		return fmt.Errorf("%w: %v", base, r.Error)
	}

	return fmt.Errorf("%w: %v - %v", base, r.Error, r.Message)
}