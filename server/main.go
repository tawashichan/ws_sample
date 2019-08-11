package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
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

const (
	textHeadUpgrade = "HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\n"
)

const (
	crlf          = "\r\n"
	colonAndSpace = ": "
	commaAndSpace = ", "
)

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
	//defer c.Close()
	if err != nil {
		panic(err)
	}
	for {
		// ws frameは小さくて8バイトなので、一旦8バイトで区切る
		buf := make([]byte,8)
		size,err := c.Read(buf)
		if size == 0 {
			continue
		} 
		if err != nil {
			panic(err)
		}
		wsPacket := readWsPacket(buf)
		fmt.Println(string(wsPacket))
		c.Write([]byte{129, 2, 101, 101})
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
	mask := refBit(secondByte, 7)
	//payloadのバイト長を取得
	payloadLength := refBit(secondByte, 6)*2*2*2*2*2*2 + refBit(secondByte, 5)*2*2*2*2*2 + refBit(secondByte, 4)*2*2*2*2 + refBit(secondByte, 3)*2*2*2 + refBit(secondByte, 2)*2*2 + refBit(secondByte, 1)*2 + refBit(secondByte, 0)*1
	// payloadの長さが7ビットで表せるかチェック
	if payloadLength > 128 {
	}
	// payloadがmaskされているかチェック。RFCの規定ではクライアントからサーバーに送る際は必ずmaskするので、一旦スキップ
	if mask == 1 {
		// maskフラグが立っていれば,maskKeyを取得する必要がある
		maskKey := b[2:6]
		//次のpayloadの開始地点がわからないので、ここで切る必要がある？
		payloadEnd := 6 + payloadLength
		rawPayload := b[6:payloadEnd]
		fmt.Printf("fin:%d\nrsv:%d\nrsv2:%d\nrsv3:%d\nopCode:%d\nmask:%d\npayloadLen:%d\n", fin, rsv1, rsv2, rsv3, opCode, mask, payloadLength)
		payload := convertPayload(int(payloadLength), maskKey, rawPayload)
		//fmt.Printf("payload: %s\n", string(convertPayload(int(payloadLength), maskKey, payload)))
		return payload//b[payloadEnd:]
	} else {
		//次のpayloadの開始地点がわからないので、ここで切る必要がある？
		payloadEnd := 2 + payloadLength
		payload := b[2:payloadEnd]
		fmt.Printf("fin:%d\nrsv:%d\nrsv2:%d\nrsv3:%d\nopCode:%d\nmask:%d\npayloadLen:%d\n", fin, rsv1, rsv2, rsv3, opCode, mask, payloadLength)
		fmt.Printf("payload: %s\n", string(payload))
		return payload//b[payloadEnd:]
	}
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
		if err != nil {
			panic(err)
		}
		go wsConnection(conn)
	}

}
