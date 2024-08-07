package main

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/rayjay214/parser/app/gw_client/myproto/gwpb"
	"github.com/rayjay214/parser/protocol/common"
	"io"
	"io/ioutil"
	"net"
	"time"
)

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	_, err := r.Read(buf[:])
	if err != nil {
		fmt.Println("read failed ", err)
		return
	}
	reader := common.NewReader(buf)
	pblen, _ := reader.ReadUint32()
	reader.ReadUint32() //msg type

	pbbuf, _ := reader.Read(int(pblen))

	resp := &gwpb.InnerResp{}
	proto.Unmarshal(pbbuf, resp)

	marshaler := jsonpb.Marshaler{}
	jsResp, _ := marshaler.MarshalToString(resp)
	fmt.Println(jsResp)
}

func main() {
	conn, err := net.Dial("unix", "/home/dev/gw/slxkgw/gw.sock")
	if err != nil {
		fmt.Println("Dial error", err)
		return
	}
	defer conn.Close()
	go reader(conn)

	req := &gwpb.InnerReq{}
	data, err := ioutil.ReadFile("req.json")
	if err != nil {
		fmt.Println("read failed")
		return
	}
	err = jsonpb.UnmarshalString(string(data), req)
	if err != nil {
		fmt.Println("parse json failed", err)
		return
	}

	pbbuf, _ := proto.Marshal(req)

	writer := common.NewWriter()
	writer.WriteUint32(uint32(len(pbbuf)))
	writer.WriteUint32(0)
	writer.Write(pbbuf)

	buf := writer.Bytes()

	_, err = conn.Write(buf)
	if err != nil {
		fmt.Println("Write error:", err)
		return
	}
	time.Sleep(2 * time.Second)

}
