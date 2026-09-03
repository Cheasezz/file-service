package grpcsrv

import (
	"errors"

	"github.com/Cheasezz/fileService/internal/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toGRPCErr(err error) error {
	var code codes.Code

	switch {
	case errors.Is(err, core.ErrEmptyFileName):
		fallthrough
	case errors.Is(err, core.ErrFirstMessage):
		fallthrough
	case errors.Is(err, core.ErrInvalidUUID):
		code = codes.InvalidArgument
	default:
		code = codes.Internal
	}

	return status.Error(code, err.Error())
}
