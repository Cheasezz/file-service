package grpcsrv

import (
	"io"

	file "github.com/Cheasezz/fileService/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) Upload(stream file.File_UploadServer) error {
	var totalSize uint64

	req, err := stream.Recv()
	if err != nil {
		return err
	}

	fileInfo := req.GetInfo()
	if fileInfo == nil {
		return status.Error(codes.InvalidArgument, "first message must be file info")
	}

	f, err := s.service.CreateFile(fileInfo.GetName())
	if err != nil {
		return status.Error(codes.Internal, "something went wrong")
	}
	defer f.Close()

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&file.UploadResp{Name: fileInfo.GetName(), Size: totalSize})
		}
		if err != nil {
			return status.Error(codes.Internal, "something went wrong")
		}

		n, _ := f.Write(req.GetChunk())
		totalSize += uint64(n)
	}
}
