package algorithmia

import (
	"fmt"
)

//Response from the API for an asynchronous request (output=void)
type AsyncResponse struct {
	AsyncProtocol string `json:"async_protocol"`
	RequestId     string `json:"request_id"`
	Error         *Err   `json:"error"` //never set!
}

func (resp *AsyncResponse) String() string {
	return fmt.Sprintf("AsyncResponse(async_protocol=%q, request_id=%q)", resp.AsyncProtocol, resp.RequestId)
}
