//go:build tools
// +build tools

package hwcalendar

import (
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
	_ "github.com/onsi/ginkgo/ginkgo"
	_ "github.com/pressly/goose"
)
