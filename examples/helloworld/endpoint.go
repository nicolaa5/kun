// Code generated by kun; DO NOT EDIT.
// github.com/RussellLuo/kun

package helloworld

import (
	"context"

	"github.com/RussellLuo/kok/pkg/httpoption"
	"github.com/RussellLuo/validating/v2"
	"github.com/go-kit/kit/endpoint"
)

type SayHelloRequest struct {
	Name string `json:"name"`
}

// ValidateSayHelloRequest creates a validator for SayHelloRequest.
func ValidateSayHelloRequest(newSchema func(*SayHelloRequest) validating.Schema) httpoption.Validator {
	return httpoption.FuncValidator(func(value interface{}) error {
		req := value.(*SayHelloRequest)
		return httpoption.Validate(newSchema(req))
	})
}

type SayHelloResponse struct {
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (r *SayHelloResponse) Body() interface{} { return r }

// Failed implements endpoint.Failer.
func (r *SayHelloResponse) Failed() error { return r.Err }

// MakeEndpointOfSayHello creates the endpoint for s.SayHello.
func MakeEndpointOfSayHello(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*SayHelloRequest)
		message, err := s.SayHello(
			ctx,
			req.Name,
		)
		return &SayHelloResponse{
			Message: message,
			Err:     err,
		}, nil
	}
}
