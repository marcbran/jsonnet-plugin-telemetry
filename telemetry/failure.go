package telemetry

import (
	"errors"
	"net/http"

	"github.com/marcbran/jpoet/pkg/jpoet"
)

func FailureFor(err error) *jpoet.Failure {
	var f *jpoet.Failure
	if errors.As(err, &f) {
		return &jpoet.Failure{Code: f.Code, Reason: f.Reason, Message: err.Error()}
	}
	return jpoet.NewFailure(http.StatusBadGateway, err.Error())
}
