package restc

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
)

func DefaultParseResponse(request *Request, response *Response) (any, error) {
	var (
		content any
		err     error
	)

	contentType := response.ContentType()
	contentType = strings.TrimSpace(strings.Split(contentType, ";")[0])
	switch contentType {
	case TypeApplicationJSON:
		content = request.GetResponseType()
		err = json.Unmarshal(response.Bytes(), &content)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrParseJSON, err)
		}
	case TypeApplicationXML, TypeTextXML:
		content = request.GetResponseType()
		err = xml.Unmarshal(response.Bytes(), &content)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrParseXML, err)
		}
	default:
		return nil, fmt.Errorf("%w: '%s'", ErrUnexpectedType, contentType)
	}

	return content, nil
}

func DefaultParseError(request *Request, response *Response) (any, error) {
	var (
		content any
		err     error
	)

	contentType := response.ContentType()
	contentType = strings.TrimSpace(strings.Split(contentType, ";")[0])
	switch contentType {
	case TypeApplicationJSON:
		content = request.GetErrorRespType()
		err = json.Unmarshal(response.Bytes(), &content)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrParseJSON, err)
		}
	case TypeApplicationXML, TypeTextXML:
		content = request.GetErrorRespType()
		err = xml.Unmarshal(response.Bytes(), &content)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrParseXML, err)
		}
	case TypeTextHTML:
		content, err = getBodyConcatainedText(response.Bytes())
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrParseHTML, err)
		}
	default:
		return nil, fmt.Errorf("%w: '%s'", ErrUnexpectedType, contentType)
	}

	return content, nil
}
