// Copyright 2020 Daniel G. Taylor
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
package openapi

import "fmt"

// ErrorDetail provides details about a specific error.
type ErrorDetail struct {
	// Message is a human-readable explanation of the error.
	Message string `json:"message,omitempty" doc:"Error message text"`

	// Location is a path-like string indicating where the error occurred.
	// It typically begins with `path`, `query`, `header`, or `body`. Example:
	// `body.items[3].tags` or `path.thing-id`.
	Location string `json:"location,omitempty" doc:"Where the error occurred, e.g. 'body.items[3].tags' or 'path.thing-id'"`

	// Value is the value at the given location, echoed back to the client
	// to help with debugging. This can be useful for e.g. validating that
	// the client didn't send extra whitespace or help when the client
	// did not log an outgoing request.
	Value any `json:"value,omitempty" doc:"The value at the given location"`
}

// Error returns the error message / satisfies the `error` interface. If a
// location and value are set, they will be included in the error message,
// otherwise just the message is returned.
func (e *ErrorDetail) Error() string {
	if e.Location == "" && e.Value == nil {
		return e.Message
	}
	return fmt.Sprintf("%s (%s: %v)", e.Message, e.Location, e.Value)
}

// ErrorDetail satisfies the `ErrorDetailer` interface.
func (e *ErrorDetail) ErrorDetail() *ErrorDetail {
	return e
}
