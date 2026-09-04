package grpcsrv

import (
	"io"

	file "github.com/Cheasezz/fileService/proto"
)

func (s *server) Download(req *file.FileInfo, stream file.File_DownloadServer) error {
	reader, err := s.service.OpenFile(req.GetClientUuid(), req.GetName())
	if err != nil {
		return toGRPCErr(err)
	}
	defer reader.Close()

	buf := make([]byte, 32*1024)
	for {
		n, err := reader.Read(buf)

		if n > 0 {
			sendErr := stream.Send(&file.Chunk{Data: buf[:n]})
			if sendErr != nil {
				return sendErr
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return toGRPCErr(err)
		}
	}

	return nil
}
