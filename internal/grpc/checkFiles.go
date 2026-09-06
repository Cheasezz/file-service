package grpcsrv

import (
	"io"

	"github.com/Cheasezz/fileService/internal/core"
	file "github.com/Cheasezz/fileService/proto"
)

func (s *server) CheckFiles(stream file.File_CheckFilesServer) error {
	req, err := stream.Recv()
	if err != nil {
		return toGRPCErr(err)
	}
	client := req.GetClient()
	if client == nil {
		return toGRPCErr(core.ErrFirstMessageUserInfo)
	}

	hl, err := s.service.UserHashList(client.GetUuid())
	if err != nil {
		return toGRPCErr(err)
	}

	for {
		fileMeta := req.GetMeta()

		if err == io.EOF {
			break
		}
		if err != nil {
			return toGRPCErr(err)
		}

		isSame := hl.CompareHash(fileMeta.GetName(), fileMeta.GetHash())

		sendErr := stream.Send(&file.SyncDecision{
			Filename:   fileMeta.GetName(),
			NeedUpload: !isSame,
		})
		if sendErr != nil {
			return sendErr
		}

	}
	return nil
}
