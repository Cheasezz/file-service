package grpcsrv

import (
	"io"

	"github.com/Cheasezz/fileService/internal/core"
	file "github.com/Cheasezz/fileService/proto"
)

func (s *server) Upload(stream file.File_UploadServer) error {
	var totalSize uint64

	req, err := stream.Recv()
	if err != nil {
		return toGRPCErr(err)
	}

	fileInfo := req.GetInfo()
	if fileInfo == nil {
		return toGRPCErr(core.ErrFirstMessage)
	}

	f, err := s.service.CreateFile(fileInfo.GetClientUuid(), fileInfo.GetName())
	if err != nil {
		return toGRPCErr(err)
	}
	defer f.Close()

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&file.UploadResp{Name: fileInfo.GetName(), Size: totalSize})
		}
		if err != nil {
			return toGRPCErr(err)
		}

		n, _ := f.Write(req.GetChunk().GetData())
		totalSize += uint64(n)
	}
}
