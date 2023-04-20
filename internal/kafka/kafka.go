// Package kafka provides some helper
// functions for producing and consuming
// protobuf messages from kafka.
package kafka

import (
	"encoding/json"

	"google.golang.org/protobuf/proto"
)

const (
	// CIDHeader defines the correlation id
	// header name used in kafka messages.
	CIDHeader = "x-cid"

	// MessageNameHeader defines the message name
	// header name used in kafka messages.
	MessageNameHeader = "x-message-name"
)

// HeadersToJSON is a helper function that
// turns a headers map to a json string
func HeadersToJSON(headers map[string]string) string {
	jsonStr, _ := json.Marshal(headers)
	return string(jsonStr)
}

// ProtoToJSON converts a protobuf message to json.
func ProtoToJSON(msg proto.Message) []byte {
	buf, _ := proto.Marshal(msg)

	return buf
}
