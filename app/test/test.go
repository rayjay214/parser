package main

import (
    "fmt"
    "github.com/rayjay214/parser/common"
    "io/ioutil"
)

func main() {
    writer := common.NewWriter()

    var buf []byte
    buf, _ = ioutil.ReadFile("a.amr")

    writer.Write(buf)

    fmt.Println(writer.Bytes())
}
