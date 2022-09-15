../bin/protoc.exe \
  --proto_path=. \
  --plugin=protoc-gen-go=../bin/protoc-gen-go.exe \
  --go_out=../ \
  *.proto
