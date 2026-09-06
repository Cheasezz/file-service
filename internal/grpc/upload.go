package grpcsrv

import (
	"io"

	"github.com/Cheasezz/fileService/internal/core"
	file "github.com/Cheasezz/fileService/proto"
)

func (s *server) Upload(stream file.File_UploadServer) error {
	var totalSize uint64
	var userID, filename string

	req, err := stream.Recv()
	if err != nil {
		return toGRPCErr(err)
	}

	fileInfo := req.GetInfo()
	if fileInfo == nil {
		return toGRPCErr(core.ErrFirstMessageFileInfo)
	}

	userID, filename = fileInfo.GetClient().GetUuid(), fileInfo.GetName()

	fw, err := s.service.CreateFile(userID, filename)
	if err != nil {
		return toGRPCErr(err)
	}
	defer fw.File.Close()

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			err = s.service.FinishUpload(fw, userID, filename)
			if err != nil {
				return toGRPCErr(err)
			}
			return stream.SendAndClose(&file.UploadResp{Name: filename, Size: totalSize})
		}
		if err != nil {
			return toGRPCErr(err)
		}

		n, _ := fw.Write(req.GetChunk().GetData())
		totalSize += uint64(n)
	}
}
