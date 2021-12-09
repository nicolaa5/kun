// Code generated by kun; DO NOT EDIT.
// github.com/RussellLuo/kun

package profilesvc

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

	codec = codecs.EncodeDecoder("DeleteAddress")
	validator = options.RequestValidator("DeleteAddress")
	r.Method(
		"DELETE", "/profiles/{id}/addresses/{addressID}",
		kithttp.NewServer(
			MakeEndpointOfDeleteAddress(svc),
			decodeDeleteAddressRequest(codec, validator),
			httpcodec.MakeResponseEncoder(codec, 200),
			append(kitOptions,
				kithttp.ServerErrorEncoder(httpcodec.MakeErrorEncoder(codec)),
			)...,
		),
	)

	codec = codecs.EncodeDecoder("DeleteProfile")
	validator = options.RequestValidator("DeleteProfile")
	r.Method(
		"DELETE", "/profiles/{id}",
		kithttp.NewServer(
			MakeEndpointOfDeleteProfile(svc),
			decodeDeleteProfileRequest(codec, validator),
			httpcodec.MakeResponseEncoder(codec, 200),
			append(kitOptions,
				kithttp.ServerErrorEncoder(httpcodec.MakeErrorEncoder(codec)),
			)...,
		),
	)

	codec = codecs.EncodeDecoder("GetAddress")
	validator = options.RequestValidator("GetAddress")
	r.Method(
		"GET", "/profiles/{id}/addresses/{addressID}",
		kithttp.NewServer(
			MakeEndpointOfGetAddress(svc),
			decodeGetAddressRequest(codec, validator),
			httpcodec.MakeResponseEncoder(codec, 200),
			append(kitOptions,
				kithttp.ServerErrorEncoder(httpcodec.MakeErrorEncoder(codec)),
			)...,
		),
	)

	codec = codecs.EncodeDecoder("GetAddresses")
	validator = options.RequestValidator("GetAddresses")
	r.Method(
		"GET", "/profiles/{id}/addresses",
		kithttp.NewServer(
			MakeEndpointOfGetAddresses(svc),
			decodeGetAddressesRequest(codec, validator),
			httpcodec.MakeResponseEncoder(codec, 200),
			append(kitOptions,
				kithttp.ServerErrorEncoder(httpcodec.MakeErrorEncoder(codec)),
			)...,
		),
	)

	codec = codecs.EncodeDecoder("GetProfile")
	validator = options.RequestValidator("GetProfile")
	r.Method(
		"GET", "/profiles/{id}",
		kithttp.NewServer(
			MakeEndpointOfGetProfile(svc),
			decodeGetProfileRequest(codec, validator),
			httpcodec.MakeResponseEncoder(codec, 200),
			append(kitOptions,
				kithttp.ServerErrorEncoder(httpcodec.MakeErrorEncoder(codec)),
			)...,
		),
	)

	codec = codecs.EncodeDecoder("PatchProfile")
	validator = options.RequestValidator("PatchProfile")
	r.Method(
		"PATCH", "/profiles/{id}",
		kithttp.NewServer(
			MakeEndpointOfPatchProfile(svc),
			decodePatchProfileRequest(codec, validator),
			httpcodec.MakeResponseEncoder(codec, 200),
			append(kitOptions,
				kithttp.ServerErrorEncoder(httpcodec.MakeErrorEncoder(codec)),
			)...,
		),
	)

	codec = codecs.EncodeDecoder("PostAddress")
	validator = options.RequestValidator("PostAddress")
	r.Method(
		"POST", "/profiles/{id}/addresses",
		kithttp.NewServer(
			MakeEndpointOfPostAddress(svc),
			decodePostAddressRequest(codec, validator),
			httpcodec.MakeResponseEncoder(codec, 200),
			append(kitOptions,
				kithttp.ServerErrorEncoder(httpcodec.MakeErrorEncoder(codec)),
			)...,
		),
	)

	codec = codecs.EncodeDecoder("PostProfile")
	validator = options.RequestValidator("PostProfile")
	r.Method(
		"POST", "/profiles",
		kithttp.NewServer(
			MakeEndpointOfPostProfile(svc),
			decodePostProfileRequest(codec, validator),
			httpcodec.MakeResponseEncoder(codec, 200),
			append(kitOptions,
				kithttp.ServerErrorEncoder(httpcodec.MakeErrorEncoder(codec)),
			)...,
		),
	)

	codec = codecs.EncodeDecoder("PutProfile")
	validator = options.RequestValidator("PutProfile")
	r.Method(
		"PUT", "/profiles/{id}",
		kithttp.NewServer(
			MakeEndpointOfPutProfile(svc),
			decodePutProfileRequest(codec, validator),
			httpcodec.MakeResponseEncoder(codec, 200),
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

func decodeDeleteAddressRequest(codec httpcodec.Codec, validator httpoption.Validator) kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var _req DeleteAddressRequest

		id := []string{chi.URLParam(r, "id")}
		if err := codec.DecodeRequestParam("id", id, &_req.Id); err != nil {
			return nil, err
		}

		addressID := []string{chi.URLParam(r, "addressID")}
		if err := codec.DecodeRequestParam("addressID", addressID, &_req.AddressID); err != nil {
			return nil, err
		}

		if err := validator.Validate(&_req); err != nil {
			return nil, err
		}

		return &_req, nil
	}
}

func decodeDeleteProfileRequest(codec httpcodec.Codec, validator httpoption.Validator) kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var _req DeleteProfileRequest

		id := []string{chi.URLParam(r, "id")}
		if err := codec.DecodeRequestParam("id", id, &_req.Id); err != nil {
			return nil, err
		}

		if err := validator.Validate(&_req); err != nil {
			return nil, err
		}

		return &_req, nil
	}
}

func decodeGetAddressRequest(codec httpcodec.Codec, validator httpoption.Validator) kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var _req GetAddressRequest

		id := []string{chi.URLParam(r, "id")}
		if err := codec.DecodeRequestParam("id", id, &_req.Id); err != nil {
			return nil, err
		}

		addressID := []string{chi.URLParam(r, "addressID")}
		if err := codec.DecodeRequestParam("addressID", addressID, &_req.AddressID); err != nil {
			return nil, err
		}

		if err := validator.Validate(&_req); err != nil {
			return nil, err
		}

		return &_req, nil
	}
}

func decodeGetAddressesRequest(codec httpcodec.Codec, validator httpoption.Validator) kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var _req GetAddressesRequest

		id := []string{chi.URLParam(r, "id")}
		if err := codec.DecodeRequestParam("id", id, &_req.Id); err != nil {
			return nil, err
		}

		if err := validator.Validate(&_req); err != nil {
			return nil, err
		}

		return &_req, nil
	}
}

func decodeGetProfileRequest(codec httpcodec.Codec, validator httpoption.Validator) kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var _req GetProfileRequest

		id := []string{chi.URLParam(r, "id")}
		if err := codec.DecodeRequestParam("id", id, &_req.Id); err != nil {
			return nil, err
		}

		if err := validator.Validate(&_req); err != nil {
			return nil, err
		}

		return &_req, nil
	}
}

func decodePatchProfileRequest(codec httpcodec.Codec, validator httpoption.Validator) kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var _req PatchProfileRequest

		if err := codec.DecodeRequestBody(r, &_req); err != nil {
			return nil, err
		}

		id := []string{chi.URLParam(r, "id")}
		if err := codec.DecodeRequestParam("id", id, &_req.Id); err != nil {
			return nil, err
		}

		if err := validator.Validate(&_req); err != nil {
			return nil, err
		}

		return &_req, nil
	}
}

func decodePostAddressRequest(codec httpcodec.Codec, validator httpoption.Validator) kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var _req PostAddressRequest

		if err := codec.DecodeRequestBody(r, &_req); err != nil {
			return nil, err
		}

		id := []string{chi.URLParam(r, "id")}
		if err := codec.DecodeRequestParam("id", id, &_req.Id); err != nil {
			return nil, err
		}

		if err := validator.Validate(&_req); err != nil {
			return nil, err
		}

		return &_req, nil
	}
}

func decodePostProfileRequest(codec httpcodec.Codec, validator httpoption.Validator) kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var _req PostProfileRequest

		if err := codec.DecodeRequestBody(r, &_req); err != nil {
			return nil, err
		}

		if err := validator.Validate(&_req); err != nil {
			return nil, err
		}

		return &_req, nil
	}
}

func decodePutProfileRequest(codec httpcodec.Codec, validator httpoption.Validator) kithttp.DecodeRequestFunc {
	return func(_ context.Context, r *http.Request) (interface{}, error) {
		var _req PutProfileRequest

		if err := codec.DecodeRequestBody(r, &_req); err != nil {
			return nil, err
		}

		id := []string{chi.URLParam(r, "id")}
		if err := codec.DecodeRequestParam("id", id, &_req.Id); err != nil {
			return nil, err
		}

		if err := validator.Validate(&_req); err != nil {
			return nil, err
		}

		return &_req, nil
	}
}
