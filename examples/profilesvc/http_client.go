// Code generated by kok; DO NOT EDIT.
// github.com/RussellLuo/kok

package profilesvc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	httpcodec "github.com/RussellLuo/kok/pkg/codec/httpv2"
)

type HTTPClient struct {
	codecs     httpcodec.Codecs
	httpClient *http.Client
	scheme     string
	host       string
	pathPrefix string
}

func NewHTTPClient(codecs httpcodec.Codecs, httpClient *http.Client, baseURL string) (*HTTPClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &HTTPClient{
		codecs:     codecs,
		httpClient: httpClient,
		scheme:     u.Scheme,
		host:       u.Host,
		pathPrefix: strings.TrimSuffix(u.Path, "/"),
	}, nil
}

func (c *HTTPClient) DeleteAddress(ctx context.Context, profileID string, addressID string) (err error) {
	codec := c.codecs.EncodeDecoder("DeleteAddress")

	path := fmt.Sprintf("/profiles/%s/addresses/%s",
		codec.EncodeRequestParam("profileID", profileID),
		codec.EncodeRequestParam("addressID", addressID),
	)
	u := &url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   c.pathPrefix + path,
	}

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusNoContent {
		return nil
	} else {
		var respErr error
		err := codec.DecodeFailureResponse(resp.Body, &respErr)
		if err == nil {
			err = respErr
		}
		return err
	}
}

func (c *HTTPClient) DeleteProfile(ctx context.Context, id string) (err error) {
	codec := c.codecs.EncodeDecoder("DeleteProfile")

	path := fmt.Sprintf("/profiles/%s",
		codec.EncodeRequestParam("id", id),
	)
	u := &url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   c.pathPrefix + path,
	}

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusNoContent {
		return nil
	} else {
		var respErr error
		err := codec.DecodeFailureResponse(resp.Body, &respErr)
		if err == nil {
			err = respErr
		}
		return err
	}
}

func (c *HTTPClient) GetAddress(ctx context.Context, profileID string, addressID string) (address Address, err error) {
	codec := c.codecs.EncodeDecoder("GetAddress")

	path := fmt.Sprintf("/profiles/%s/addresses/%s",
		codec.EncodeRequestParam("profileID", profileID),
		codec.EncodeRequestParam("addressID", addressID),
	)
	u := &url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   c.pathPrefix + path,
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return Address{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Address{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusNoContent {
		var respBody struct {
			Address Address `json:"address"`
		}
		err := codec.DecodeSuccessResponse(resp.Body, &respBody)
		if err != nil {
			return Address{}, err
		}
		return respBody.Address, nil
	} else {
		var respErr error
		err := codec.DecodeFailureResponse(resp.Body, &respErr)
		if err == nil {
			err = respErr
		}
		return Address{}, err
	}
}

func (c *HTTPClient) GetAddresses(ctx context.Context, id string) (addresses []Address, err error) {
	codec := c.codecs.EncodeDecoder("GetAddresses")

	path := fmt.Sprintf("/profiles/%s/addresses",
		codec.EncodeRequestParam("id", id),
	)
	u := &url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   c.pathPrefix + path,
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusNoContent {
		var respBody struct {
			Addresses []Address `json:"addresses"`
		}
		err := codec.DecodeSuccessResponse(resp.Body, &respBody)
		if err != nil {
			return nil, err
		}
		return respBody.Addresses, nil
	} else {
		var respErr error
		err := codec.DecodeFailureResponse(resp.Body, &respErr)
		if err == nil {
			err = respErr
		}
		return nil, err
	}
}

func (c *HTTPClient) GetProfile(ctx context.Context, id string) (profile Profile, err error) {
	codec := c.codecs.EncodeDecoder("GetProfile")

	path := fmt.Sprintf("/profiles/%s",
		codec.EncodeRequestParam("id", id),
	)
	u := &url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   c.pathPrefix + path,
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return Profile{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Profile{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusNoContent {
		var respBody struct {
			Profile Profile `json:"profile"`
		}
		err := codec.DecodeSuccessResponse(resp.Body, &respBody)
		if err != nil {
			return Profile{}, err
		}
		return respBody.Profile, nil
	} else {
		var respErr error
		err := codec.DecodeFailureResponse(resp.Body, &respErr)
		if err == nil {
			err = respErr
		}
		return Profile{}, err
	}
}

func (c *HTTPClient) PatchProfile(ctx context.Context, id string, profile Profile) (err error) {
	codec := c.codecs.EncodeDecoder("PatchProfile")

	path := fmt.Sprintf("/profiles/%s",
		codec.EncodeRequestParam("id", id),
	)
	u := &url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   c.pathPrefix + path,
	}

	reqBody := struct {
		Profile Profile `json:"profile"`
	}{
		Profile: profile,
	}
	reqBodyReader, headers, err := codec.EncodeRequestBody(&reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", u.String(), reqBodyReader)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusNoContent {
		return nil
	} else {
		var respErr error
		err := codec.DecodeFailureResponse(resp.Body, &respErr)
		if err == nil {
			err = respErr
		}
		return err
	}
}

func (c *HTTPClient) PostAddress(ctx context.Context, profileID string, address Address) (err error) {
	codec := c.codecs.EncodeDecoder("PostAddress")

	path := fmt.Sprintf("/profiles/%s/addresses",
		codec.EncodeRequestParam("profileID", profileID),
	)
	u := &url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   c.pathPrefix + path,
	}

	reqBody := struct {
		Address Address `json:"address"`
	}{
		Address: address,
	}
	reqBodyReader, headers, err := codec.EncodeRequestBody(&reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u.String(), reqBodyReader)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusNoContent {
		return nil
	} else {
		var respErr error
		err := codec.DecodeFailureResponse(resp.Body, &respErr)
		if err == nil {
			err = respErr
		}
		return err
	}
}

func (c *HTTPClient) PostProfile(ctx context.Context, profile Profile) (err error) {
	codec := c.codecs.EncodeDecoder("PostProfile")

	path := "/profiles"
	u := &url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   c.pathPrefix + path,
	}

	reqBody := struct {
		Profile Profile `json:"profile"`
	}{
		Profile: profile,
	}
	reqBodyReader, headers, err := codec.EncodeRequestBody(&reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u.String(), reqBodyReader)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusNoContent {
		return nil
	} else {
		var respErr error
		err := codec.DecodeFailureResponse(resp.Body, &respErr)
		if err == nil {
			err = respErr
		}
		return err
	}
}

func (c *HTTPClient) PutProfile(ctx context.Context, id string, profile Profile) (err error) {
	codec := c.codecs.EncodeDecoder("PutProfile")

	path := fmt.Sprintf("/profiles/%s",
		codec.EncodeRequestParam("id", id),
	)
	u := &url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   c.pathPrefix + path,
	}

	reqBody := struct {
		Profile Profile `json:"profile"`
	}{
		Profile: profile,
	}
	reqBodyReader, headers, err := codec.EncodeRequestBody(&reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", u.String(), reqBodyReader)
	if err != nil {
		return err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusNoContent {
		return nil
	} else {
		var respErr error
		err := codec.DecodeFailureResponse(resp.Body, &respErr)
		if err == nil {
			err = respErr
		}
		return err
	}
}