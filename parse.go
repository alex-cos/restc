package restc

import (
	"encoding/json"
	"errors"
	"fmt"
)

func DefaultParseResponse(request *Request, response *Response) (any, error) {
	var (
		content any
		err     error
	)

	switch response.ContentType() {
	case TypeApplicationJSON:
		content = request.GetResponseType()
		err = json.Unmarshal(response.Bytes(), &content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %w", err)
		}

	default:
		return nil, errors.New("unexpected response type")
	}

	return content, nil
}

func DefaultParseError(request *Request, response *Response) (any, error) {
	var (
		content any
		err     error
	)

	switch response.ContentType() {
	case TypeApplicationJSON:
		content = request.GetErrorRespType()
		err := json.Unmarshal(response.Bytes(), &content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %w", err)
		}
	case TypeTextHTML:
		content, err = getBodyConcatainedText(response.Bytes())
		if err != nil {
			return nil, fmt.Errorf("failed to parse HTML response: %w", err)
		}
	default:
		return nil, errors.New("unexpected response type")
	}

	return content, nil
}
