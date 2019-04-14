package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

type Bit4 struct {
	first  bool
	second bool
	third  bool
	fourth bool
}

type WsPacket struct {
	Fin           bool
	RSV1          bool
	RSV2          bool
	RSV3          bool
	OpCode        Bit4
	Mask          bool
	PayloadLength uint
	Payload       []byte
}

func (w WsPacket) ToByte() []byte {
	var base = []byte{}
	//payload := []byte("huga")
	return append(base, []byte("huga")...)
}

func Upgrade(conn net.Conn) (net.Conn, error) {
	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		return nil, err
	}
	wsKey := request.Header["Sec-Websocket-Key"][0]
	response := http.Response{
		StatusCode: 101,
		Header: map[string][]string{
			"Upgrade":              {"websocket"},
			"Connection":           {"Upgrade"},
			"Sec-WebSocket-Accept": {genSecWebsocketAccept(wsKey)},
		},
	}
	if err := response.Write(conn); err != nil {
		return nil, err
	}
	return conn, err
}

func wsConnection(conn net.Conn) {
	c, err := Upgrade(conn)
	defer c.Close()
	if err != nil {
		panic(err)
	}
	for {
		b, err := ioutil.ReadAll(c)
		if err != nil {
			panic(err)
		}
		readWsPacket(b)
		res := WsPacket{}
		c.Write(res.ToByte())
		break
	}
}

func byteToBynaryDigit(b byte) {
	fmt.Printf("%d%d%d%d%d%d%d%d\n",
		b&127,
		b&64,
		b&32,
		b&16,
		b&8,
		b&4,
		b&2,
		b&1,
	)
}

func readWsPacket(b []byte) {
	if len(b) == 0 {
		return
	}
	firstByte := b[0]
	fin := firstByte & 1
	rsv1 := firstByte & 2
	rsv2 := firstByte & 4
	rsv3 := firstByte & 8
	opCode := (firstByte&127)*2*2*2 + (firstByte&64)*2*2 + (firstByte&32)*2 + (firstByte&16)*1
	secondByte := b[1]
	mask := secondByte & 127
	payloadLength := (secondByte&64)*2*2*2*2*2*2 + (secondByte&32)*2*2*2*2*2 + (secondByte&16)*2*2*2*2 + (secondByte&8)*2*2*2 + (secondByte&4)*2*2 + (secondByte&2)*2 + (secondByte&1)*1
	//byteToBynaryDigit(secondByte)
	// payloadの長さが7ビットで表せるかチェック
	if payloadLength > 128 {
	}
	//if mask == 1 {}
	maskKey := b[2:6]
	payload := b[6:]
	fmt.Printf("fin:%d\nrsv:%d\nrsv2:%d\nrsv3:%d\nopCode:%d\nmask:%d\npayloadLen:%d\n", fin, rsv1, rsv2, rsv3, opCode, mask, payloadLength)
	fmt.Println(string(unMaskPayload(int(payloadLength), maskKey, payload)))
}

func unMaskPayload(payloadLen int, maskKey []byte, maskedPayload []byte) []byte {
	var result = []byte{}
	for i := 0; i < len(maskedPayload); i++ {
		result = append(result, maskedPayload[i]^maskKey[i%4])
	}
	return result
}

func genSecWebsocketAccept(nonce string) string {
	base := nonce + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1.Sum([]byte(base))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	fmt.Println("start websocket server")
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		wsConnection(conn)
	}
}
