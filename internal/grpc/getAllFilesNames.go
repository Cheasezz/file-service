package grpcsrv

import (
	"context"

	file "github.com/Cheasezz/fileService/proto"
)

func (s *server) GetAllFilesNames(ctx context.Context, req *file.Client) (*file.Files, error) {
	filesNames, err := s.service.GetAllFilesNames(req.GetUuid())
	if err != nil {
		return nil, toGRPCErr(err)
	}

	return &file.Files{Names: filesNames}, nil
}
