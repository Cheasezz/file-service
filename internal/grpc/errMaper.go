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
	case errors.Is(err, core.ErrFirstMessageFileInfo):
		fallthrough
	case errors.Is(err, core.ErrInvalidUUID):
		fallthrough
	case errors.Is(err, core.ErrFirstMessageUserInfo):
		fallthrough
	case errors.Is(err, core.ErrEmptyFileMeta):
		fallthrough
	case errors.Is(err, core.ErrInvalidFileName):
		fallthrough
	case errors.Is(err, core.ErrEmptyChunk):
		code = codes.InvalidArgument
	case errors.Is(err, core.ErrFileNotFound):
		code = codes.NotFound
	default:
		return status.Error(codes.Internal, core.ErrInternal.Error())
	}

	return status.Error(code, err.Error())
}
