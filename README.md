## What is it
 This is demonstration of grpc client side streaming.
## How it works for client
- Client call grpc method for upload file;
- Send file name in first message for stream; 
- In Cycle send 32kb chunks of file for stream;
- If EOF then break cycle and close stream;
- Print uploaded file name and size;
 ## How it works for server
 - Server starts;
 - Create base directory for store files.    
 It is created relative to the home directory. (os.UserHomeDir);
 - Server receive first message from the client stream.  
   If its file name then create file in base directory;
 - In cycle receive chunks and write into file;
 - If EOF then break cycle and close stream with response (downloaded file name);
## How it start
Clone repo:

```bash
git clone https://github.com/Cheasezz/file-service.git
```

Change directory:

```bash
cd file-service
```

Start server:

```bash
make server
```

Start client (in other terminal on repo directory):

```bash
make client
```

The file will be in **your_home_directory/.fileService/uploads**
