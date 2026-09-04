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
		fallthrough
	case errors.Is(err, core.ErrFileNotFound):
		code = codes.InvalidArgument
	default:
		return status.Error(codes.Internal, core.ErrInternal.Error())
	}

	return status.Error(code, err.Error())
}
