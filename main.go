package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

/*
type Pages struct{}

func (p *Pages) home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Called")

	document, _ := template.ParseFiles("index.html")
	document.Execute(w, nil)
}
*/

func send(socket io.Writer, data []byte) {
	length := make([]byte, 4)
	binary.LittleEndian.PutUint32(length, uint32(len(data)))

	data = append(length, data...)
	socket.Write(data)
}

func recv(socket io.Reader) []byte {
	buffer := make([]byte, 4)
	socket.Read(buffer)
	length := binary.LittleEndian.Uint32(buffer)

	buffer = make([]byte, length)
	socket.Read(buffer)
	return buffer
}

func contactServerJSON(message map[string]string) (map[string]string, error) {
	remote_addr := new(net.TCPAddr)
	remote_addr.IP = []byte{127, 0, 0, 1}
	remote_addr.Port = 50508

	connection, err := net.DialTCP("tcp4", nil, remote_addr)
	if err != nil {
		return nil, err
	}

	json_message, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	send(connection, json_message)
	json_response := recv(connection)

	response := *new(map[string]string)
	json.Unmarshal(json_response, &response)

	return response, nil
}

func generateSignature(payload string) (signature, public_key string) {
	req := map[string]string{"command": "generate", "payload": payload}
	response, _ := contactServerJSON(req)

	return response["signature"], response["public key"]
}

func verifySignature(payload, signature, public_key string) bool {
	req := map[string]string{"command": "verify", "payload": payload, "signature": signature, "public key": public_key}
	response, _ := contactServerJSON(req)

	return response["is valid"] == "True"
}

func blacklistSignature(payload, signature, public_key string) {
	req := map[string]string{"payload": payload, "signature": signature, "public key": public_key}
	contactServerJSON(req)
}

func main() {
	signature, key := generateSignature("hello")

	b := verifySignature("herro", signature, key)

	fmt.Println(b)
	fmt.Println(signature)
	fmt.Println(key)
}
