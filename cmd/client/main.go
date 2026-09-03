package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/Cheasezz/fileService/config"
	"github.com/Cheasezz/fileService/pkg/logger"
	file "github.com/Cheasezz/fileService/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var filePath = flag.String("path", "", "File path for upload")

func main() {
	var opts []grpc.DialOption
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)

	cfg := config.MustLoad()

	log := logger.New(cfg.Env)

	log.Info("starting application")

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("127.0.0.1:"+strconv.Itoa(cfg.GRPC.Port), opts...)
	if err != nil {
		log.Error("fail to dial: %v", err)
	}
	defer conn.Close()

	client := file.NewFileClient(conn)

	flag.Parse()
	err = upload(client, *filePath)
	if err != nil {
		log.Error("Error while upload: ", err)
	}

	log.Info("File correct uploaded")
	<-exit

	log.Info("Work is over")
}

func upload(fc file.FileClient, filePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := fc.Upload(ctx)
	if err != nil {
		return err
	}

	// First message with filename
	err = stream.Send(&file.UploadReq{
		Payload: &file.UploadReq_Info{
			Info: &file.FileInfo{
				ClientUuid: uuid.NewString(),
				Name:       filepath.Base(filePath),
			},
		},
	})
	if err != nil {
		return err
	}

	// Open file end send chunks
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Chunk is 32kb
	buf := make([]byte, 32*1024)
	for {
		n, err := f.Read(buf)
		if n > 0 {
			sendErr := stream.Send(&file.UploadReq{
				Payload: &file.UploadReq_Chunk{
					// chunk is buf[:n] cuz buffer may be not full at the end.
					Chunk: buf[:n],
				},
			})
			if sendErr != nil {
				return sendErr
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}

	status, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	fmt.Printf("Uploaded file: %s, size: %d\n", status.GetName(), status.GetSize())

	return nil
}
