package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/Cheasezz/fileService/config"
	"github.com/Cheasezz/fileService/pkg/logger"
	file "github.com/Cheasezz/fileService/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	filePath = flag.String("path", "", "File path for upload")
	dirPath  = flag.String("dir", "", "Dir path for upload files (no nested files)")
	savePath string
	clientID string
)

type LocalFile struct {
	Path, Name, Hash string
}

func collectFiles(paths []string) ([]LocalFile, error) {
	files := make([]LocalFile, 0, len(paths))

	for _, path := range paths {
		hash, err := hashFile(path)
		if err != nil {
			return nil, err
		}

		files = append(files, LocalFile{
			Name: filepath.Base(path),
			Path: path,
			Hash: hash,
		})
	}

	return files, nil
}

func collectFilesFromDir(dir string) ([]LocalFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("cant read dir: %w", err)
	}

	paths := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			paths = append(paths, filepath.Join(dir, e.Name()))
		}
	}

	return collectFiles(paths)
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("cant open file %s: %w", path, err)
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", fmt.Errorf("cant calculate hash for file %s: %w", path, err)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func main() {
	var opts []grpc.DialOption
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)

	cfg := config.MustLoad()

	log := logger.New(cfg.Env)

	log.Info("starting application")

	// --- Init client id
	clientID = uuid.NewString()
	// ---

	// --- Create dir for downloaded files
	userDir, err := os.UserHomeDir()
	if err != nil {
		log.Error("cant find user dir: %v", err)
		return
	}

	savePath = filepath.Join(userDir, ".fileService", "download")

	err = os.MkdirAll(savePath, 0o755)
	if err != nil {
		log.Error("cant create save dir: %v", err)
		return
	}
	// ---

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("127.0.0.1:"+strconv.Itoa(cfg.GRPC.Port), opts...)
	if err != nil {
		log.Error("fail to dial: %v", err)
		return
	}
	defer conn.Close()

	client := file.NewFileClient(conn)

	flag.Parse()

	fileInfo := &file.FileInfo{
		Client: &file.Client{
			Uuid: clientID,
		},
		Name: filepath.Base(*filePath),
	}

	err = upload(client, *filePath)
	if err != nil {
		log.Error("Error while upload: ", err)
	} else {
		log.Info("File correct uploaded")
	}

	err = getAllFilesNames(client)
	if err != nil {
		log.Error("Cant get all files names: %v", err)
	}

	err = download(client, fileInfo)
	if err != nil {
		log.Error("Error while download: ", err)
	} else {
		log.Info("File correct download")
	}

	files, err := collectFilesFromDir(*dirPath)
	if err != nil {
		log.Error("Error while collect files from dir: ", err)
	}

	err = syncAndUpload(client, files)
	if err != nil {
		log.Error("Error while syncAndUpload: ", err)
	} else {
		log.Info("Files syncronized and uploaded on server")
	}

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
				Client: &file.Client{
					Uuid: clientID,
				},
				Name: filepath.Base(filePath),
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
					Chunk: &file.Chunk{
						Data: buf[:n],
					},
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

	fmt.Printf("Uploaded file: %s, bytes: %d\n", status.GetName(), status.GetSize())

	return nil
}

func download(fc file.FileClient, fi *file.FileInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := fc.Download(ctx, fi)
	if err != nil {
		return err
	}

	f, err := os.Create(savePath + "/" + fi.Name)
	if err != nil {
		return fmt.Errorf("cant creale file in save dir: %v", err)
	}
	defer f.Close()

	var totalSize uint64
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error while receiving chunks: %v", err)
		}

		n, err := f.Write(req.GetData())
		if err != nil {
			return fmt.Errorf("error write chunk to file: %v", err)
		}

		totalSize += uint64(n)
	}

	fmt.Printf("Downloaded file: %s, bytes: %d\n", fi.Name, totalSize)
	return nil
}

func getAllFilesNames(fc file.FileClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	files, err := fc.GetAllFilesNames(ctx, &file.Client{Uuid: clientID})
	if err != nil {
		return fmt.Errorf("error while get all files names: %v", err)
	}

	if len(files.GetNames()) == 0 {
		fmt.Println("No one file uploaded on server")
		return nil
	}

	fmt.Println("--- Uploaded files in server ---")
	for _, file := range files.GetNames() {
		fmt.Println(file)
	}
	fmt.Println("--- ---")

	return nil
}

func checkFiles(fc file.FileClient, files []LocalFile) (<-chan LocalFile, <-chan error) {
	toUpload := make(chan LocalFile, len(files))
	errCh := make(chan error, 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)

	go func() {
		defer close(toUpload)
		defer close(errCh)
		defer cancel()

		stream, err := fc.CheckFiles(ctx)
		if err != nil {
			errCh <- fmt.Errorf("cant open CheckFiles stream: %v", err)
			return
		}

		if err := stream.Send(&file.CheckFilesReq{
			Payload: &file.CheckFilesReq_Client{
				Client: &file.Client{
					Uuid: clientID,
				},
			},
		}); err != nil {
			errCh <- fmt.Errorf("cand send clien id: %w", err)
			return
		}

		go func() {
			for _, f := range files {
				if err := stream.Send(&file.CheckFilesReq{
					Payload: &file.CheckFilesReq_Meta{
						Meta: &file.FileMeta{
							Name: f.Name,
							Path: f.Path,
							Hash: f.Hash,
						},
					},
				}); err != nil {
					errCh <- fmt.Errorf("cant check file %s: %w", f.Name, err)
					return
				}
			}
			if err := stream.CloseSend(); err != nil {
				errCh <- fmt.Errorf("cant closeSend in checkFiles: %w", err)
				return
			}
		}()

		for {
			decision, err := stream.Recv()

			if err == io.EOF {
				return
			}
			if err != nil {
				errCh <- fmt.Errorf("error recv decision: %v", err)
				return
			}
			fmt.Printf("File decision %s: %t\n", decision.GetFilename(), decision.GetNeedUpload())
			if decision.GetNeedUpload() {
				toUpload <- LocalFile{Name: decision.GetFilename(), Path: decision.GetFilepath()}
			}
		}
	}()

	return toUpload, errCh
}

func syncAndUpload(fc file.FileClient, files []LocalFile) error {
	var errs []error
	const workers = 3
	uploadErrs := make(chan error, len(files))
	var wg sync.WaitGroup

	toUpload, syncErrCh := checkFiles(fc, files)

	for range workers {
		wg.Go(func() {
			for f := range toUpload {
				err := upload(fc, f.Path)
				if err != nil {
					uploadErrs <- fmt.Errorf("upload %s: %w", f.Name, err)
				}
			}
		})
	}

	wg.Wait()
	close(uploadErrs)

	if syncErr := <-syncErrCh; syncErr != nil {
		errs = append(errs, syncErr)
	}

	for err := range uploadErrs {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
