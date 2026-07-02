FROM golang:1.26

RUN apt update

# Правильный синтаксис установки без протоколов и двоеточий
RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# RUN apt install -y protobuf-compiler
# RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest