package forms

import (
	"errors"
	errorsPkg "github.com/pkg/errors"
)

type ErrorHandler struct {
	err error
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

func (h *ErrorHandler) Handle(err error, wrapText ...string) {
	if err == nil {
		return
	}

	if len(wrapText) > 0 {
		err = errorsPkg.Wrap(err, wrapText[0])
	}
	if h.err != nil {
		err = errors.Join(err, h.err)
	}

	h.err = err
}

func (h *ErrorHandler) GetErrorText() string {
	errText := ""
	if h.err != nil {
		errText = h.err.Error()
		h.err = nil
	}

	return errText
}
