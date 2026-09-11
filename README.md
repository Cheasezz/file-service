# File Service

A file storage and synchronization service built with **Go** and **gRPC**.

The main goal of the project is to explore and demonstrate the different gRPC communication patterns in a practical application.

## Goals

- Learn and practice all four gRPC communication patterns.
- Work with client, server and bidirectional streaming.
- Implement chunked file transfer.
- Practice stream lifecycle and `io.EOF` handling.
- Work with goroutines and channels around bidirectional streams.
- Implement file synchronization using SHA-256 hashes.
- Separate transport, service and filesystem layers.

## gRPC Api

The service provides four gRPC methods:

| Method | Type | Description |
|---|---|---|
| `Upload` | Client streaming | Upload a file in chunks |
| `Download` | Server streaming | Download a file in chunks |
| `GetAllFilesNames` | Unary | Get stored file names |
| `CheckFiles` | Bidirectional streaming | Compare local and server file hashes |

## gRPC Streaming

The project contains examples of all gRPC communication patterns:

```text
Unary
Client ───── Request ─────► Server
Client ◄──── Response ───── Server


Client Streaming
Client ──► message ──► message ──► message ──► Server
Client ◄──────────────────────────── Response


Server Streaming
Client ───── Request ─────► Server
Client ◄── chunk ◄── chunk ◄── chunk ◄── Server


Bidirectional Streaming
Client ──► metadata ──► metadata ──► metadata ──► Server
Client ◄── decision ◄── decision ◄── decision ◄── Server
```

## Key Features

- **32 KB chunks** are used for file transfer.
- Files are written directly to the filesystem while receiving chunks.
- SHA-256 hashes are calculated incrementally during upload.
- File downloads are streamed directly to the destination file.
- The bidirectional `CheckFiles` stream can send and receive data concurrently.
- Files sync. Only files that are missing or have a different hash are uploaded.
- Multiple workers can upload files in parallel.
- Basic gRPC error mapping is implemented:
  - invalid client data → `InvalidArgument`
  - missing file → `NotFound`
  - unexpected errors → `Internal`
- File names are validated to prevent path traversal.

## Architecture

```text
Client
  │
  ▼
gRPC Transport
  │
  ▼
Service
  │
  ▼
Filesystem Repository
```

The gRPC layer is responsible for transport and protocol-specific errors, while the service layer contains application logic and the repository handles filesystem operations.

## Project Structure

```text
cmd/
├── client/
└── server/

internal/
├── core/       # Domain types and errors
├── grpc/       # gRPC handlers and error mapping
├── service/    # Application logic
└── repo/       # Filesystem operations

proto/
└── fileService.proto

pkg/
└── logger/
```

## Running

Clone the repository:

```bash
git clone https://github.com/Cheasezz/file-service.git
cd file-service
```

Start the server:

```bash
make server
```

In another terminal, start the client:

```bash
make client
```

Remove file-service folder for download/upload:

```bash
make clear
```

The server stores uploaded files in **your-os-user-home-dir/.fileService/uploads**.  
The client stores downloaded files in **your-os-user-home-dir/.fileService/download**.

## Project Scope

This is primarily a **learning and demonstration project focused on gRPC streaming** rather than a production-ready file storage system.

The filesystem is intentionally used as the storage layer so that the main focus remains on:

- gRPC streaming
- protobuf message design
- chunked data transfer
- stream lifecycle
- concurrency
- file synchronization
- error handling
