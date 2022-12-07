package endpoints

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/zalgonoise/dns/store/encoder"
	"github.com/zalgonoise/logx/attr"
	"github.com/zalgonoise/x/spanner"
)

var (
	ErrInvalidBody   = errors.New("invalid body")
	ErrInvalidJSON   = errors.New("body contains invalid JSON")
	ErrInternal      = errors.New("internal error")
	ErrResponseWrite = errors.New("failed to write response")
)

// ContextResponseEncoder is a type used to identify encoders in Contexts
type ContextResponseEncoder string

// ResponseEncoderKey is the common key used by this package to store an
// EncodeDecoder in a Context
const ResponseEncoderKey ContextResponseEncoder = "encoder"

// enc returns the EncodeDecoder from the input Context `ctx`, or creates a new
// JSON encoder if the context doesn't have one
func enc(ctx context.Context) encoder.EncodeDecoder {
	var enc encoder.EncodeDecoder = nil

	v := ctx.Value(ResponseEncoderKey)
	if v != nil {
		if e, ok := v.(encoder.EncodeDecoder); ok {
			enc = e
		}
	}
	if enc == nil {
		return encoder.New("json")
	}
	return enc
}

// HTTPWriter writes an object to a http.ResponseWriter as a HTTP response
type HTTPWritter interface {
	// WriteHTTP writes the contents of the object to the http.ResponseWriter `w`
	WriteHTTP(ctx context.Context, w http.ResponseWriter)
}

// httpResponse is a generic type for HTTP responses, both OK and not-OK
type HttpResponse[T any] struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
	Error   error  `json:"error,omitempty"`
	Data    *T     `json:"data,omitempty"`
	Status  int    `json:"-"`
}

// newResponse builds a httpResponse object for type T derived from the input data (*T)
//
// It takes in an int `status`, a string `message` and either an error or data.
//   - if the error is nil, it is assumed that it is an OK response and the provided *T `data`
//     is stored in the response, and the error ignored
//   - if the error is not nil, it is assumed that the response is not-OK, so the *T `data` is
//     ignored, and the error stored in the response
func NewResponse[T any](status int, message string, err error, data *T) HttpResponse[T] {
	if err != nil {
		return HttpResponse[T]{
			Success: false,
			Message: message,
			Error:   err,
			Data:    nil,
			Status:  status,
		}
	}

	return HttpResponse[T]{
		Success: true,
		Message: message,
		Error:   nil,
		Data:    data,
		Status:  status,
	}
}

// WriteHTTP writes the contents of the object to the http.ResponseWriter `w`
func (r HttpResponse[T]) WriteHTTP(ctx context.Context, w http.ResponseWriter) {
	enc := enc(ctx)

	w.WriteHeader(r.Status)
	response, _ := enc.Encode(r.Data)
	_, _ = w.Write(response)
}

// readBody reads the data in the Body of *http.Request `r` as a bytes buffer,
// and attempts to decode it into an object of type T by creating a new pointer of
// this type and decoding the buffer into it
//
// Returns a pointer to the object T and an error
func readBody[T any](ctx context.Context, r *http.Request) (*T, error) {
	ctx, s := spanner.Start(ctx, "http.readBody", attr.String("for_type", fmt.Sprintf("%T", *new(T))))
	defer s.End()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		s.Event("error reading body", attr.New("error", err.Error()))
		return nil, fmt.Errorf("%w: %v", ErrInvalidBody, err)
	}
	item := new(T)
	err = enc(ctx).Decode(b, item)
	if err != nil {
		s.Event("error decoding buffer", attr.New("error", err.Error()), attr.String("buffer", string(b)))
		return nil, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}
	s.Event("decoded request body", attr.Ptr("item", item))
	return item, nil
}