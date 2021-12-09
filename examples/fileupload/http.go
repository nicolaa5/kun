// Code generated by kun; DO NOT EDIT.
// github.com/RussellLuo/kun

package fileupload

import (
	"context"
	"net/http"

	"github.com/RussellLuo/kok/pkg/httpcodec"
	httpoption "github.com/RussellLuo/kok/pkg/httpoption2"
	"github.com/RussellLuo/kok/pkg/oas2"
	"github.com/go-chi/chi"
	kithttp "github.com/go-kit/kit/transport/http"
)

func NewHTTPRouter(svc Service, codecs httpcodec.Codecs, opts ...httpoption.Option) chi.Router {
	r := chi.NewRouter()
	options := httpoption.NewOptions(opts...)

	r.Method("GET", "/api", oas2.Handler(OASv2APIDoc, options.ResponseSchema()))

	var codec httpcodec.Codec
	var validator httpoption.Validator
	var kitOptions []kithttp.ServerOption

	codec = codecs.EncodeDecoder("Upload")
	validator = options.RequestValidator("Upload")
	r.Method(
		"POST", "/upload",
		kithttp.NewServer(
			MakeEndpointOfUpload(svc),
			decodeUploadRequest(codec, validator),
			httpcodec.MakeResponseEncoder(codec, 204),
			append(kitOptions,
				kithttp.ServerErrorEncoder(httpcodec.MakeErrorEncoder(codec)),
			)...,
		),
	)

	return r
}

func NewHTTPRouterWithOAS(svc Service, codecs httpcodec.Codecs, schema oas2.Schema) chi.Router {
	return NewHTTPRouter(svc, codecs, httpoption.ResponseSchema(schema))
}

func decodeUploadRequest(codec httpcodec.Codec, validator httpoption.Validator) kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var _req UploadRequest

		if err := codec.DecodeRequestBody(r, &_req); err != nil {
			return nil, err
		}

		if err := validator.Validate(&_req); err != nil {
			return nil, err
		}

		return &_req, nil
	}
}
