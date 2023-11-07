cp ../../../gateway/gw.proto .
protoc --go_out=. gw.proto
