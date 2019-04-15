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
	return append(base, []byte{}...)
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
		resB := b
		for {
			b = readWsPacket(b)
			if len(b) == 0 {
				break
			}
		}
		i, err := c.Write(resB)
		fmt.Println(i)
		if err != nil {
			panic(err)
		}
		break
	}
}

func ByteToBinaryDigit(b byte) string {
	var result = ""
	for i := 7; i >= 0; i-- {
		result = result + fmt.Sprint(refBit(b, uint(i)))
	}
	return result
}

func refBit(target byte, num uint) int {
	return (int(target) >> num) & 1
}

func readWsPacket(b []byte) []byte {
	if len(b) == 0 {
		return []byte{}
	}
	firstByte := b[0]
	fin := refBit(firstByte, 7)
	rsv1 := refBit(firstByte, 6)
	rsv2 := refBit(firstByte, 5)
	rsv3 := refBit(firstByte, 4)
	//本当は16進数に変換する必要がある
	opCode := refBit(firstByte, 3)*2*2*2 + refBit(firstByte, 2)*2*2 + refBit(firstByte, 1)*2 + refBit(firstByte, 0)*1
	secondByte := b[1]
	mask := refBit(firstByte, 7)
	//payloadのバイト長を取得
	payloadLength := refBit(secondByte, 6)*2*2*2*2*2*2 + refBit(secondByte, 5)*2*2*2*2*2 + refBit(secondByte, 4)*2*2*2*2 + refBit(secondByte, 3)*2*2*2 + refBit(secondByte, 2)*2*2 + refBit(secondByte, 1)*2 + refBit(secondByte, 0)*1
	// payloadの長さが7ビットで表せるかチェック
	if payloadLength > 128 {
	}
	// payloadがmaskされているかチェック。RFCの規定ではクライアントからサーバーに送る際は必ずmaskするので、一旦スキップ
	//if mask == 1 {}
	// maskフラグが立っていれば,maskKeyを取得する必要がある
	maskKey := b[2:6]
	//次のpayloadの開始地点がわからないので、ここで切る必要がある？
	payloadEnd := 6 + payloadLength
	payload := b[6:payloadEnd]
	fmt.Printf("fin:%d\nrsv:%d\nrsv2:%d\nrsv3:%d\nopCode:%d\nmask:%d\npayloadLen:%d\n", fin, rsv1, rsv2, rsv3, opCode, mask, payloadLength)
	fmt.Printf("payload: %s\n", string(convertPayload(int(payloadLength), maskKey, payload)))
	return b[payloadEnd:]
}

func convertPayload(payloadLen int, maskKey []byte, maskedPayload []byte) []byte {
	var result = []byte{}
	//1バイトずつpayloadを切り出して変換していく
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
		fmt.Println("new connection")
		if err != nil {
			panic(err)
		}
		wsConnection(conn)
	}
}
