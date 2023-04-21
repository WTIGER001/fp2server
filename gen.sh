#!/bin/bash
#!/bin/sh

# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# npm install ts-proto

# Go Part
PREFIX=github.com/wtiger001/f2pserver
protoc --proto_path=proto --go_out=common --go_opt=paths=source_relative models.proto

#rm ../../market-maker-website/projects/api/src/lib/binary/*
protoc --plugin=protoc-gen-ts_proto.cmd --proto_path=proto  --ts_proto_out=../fp2web/src/app/modelp models.proto