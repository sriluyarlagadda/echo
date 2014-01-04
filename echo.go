//package echo is used to make an ICMP Echo Request
package echo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func OnesCompliment(x int32) uint16 {
	var sum_16 uint16
	if x >= 0 {
		var overflow uint16 = uint16(x >> 16)
		x_16 := uint16(x & 65535)
		sum_16 = x_16 + overflow
	} else {
		value := OnesCompliment(-x)
		return ^(value)
	}

	return sum_16
}

//calculate the checksum of the given payload
func CheckSum(payload []byte) uint16 {
	reader := bytes.NewReader(payload)
	payloadInt_16 := make([]uint16, len(payload)/2)

	//convert slice of bytes to slice of ints
	binary.Read(reader, binary.LittleEndian, payloadInt_16)

	var sum int32 = 0
	for _, value := range payloadInt_16 {
		sum += int32(value)
	}
	checkSum := OnesCompliment(-sum)
	return checkSum
}

const MESSAGE_TYPE byte = 8

type EchoMessage struct {
	messageType    byte
	identifier     uint16
	sequenceNumber uint16
	data           []byte
	iaddr          *net.IPAddr
}

//Create a new echo Message
func NewMessage(identifier, sequencenumber uint16) *EchoMessage {
	return &EchoMessage{
		messageType:    8,
		identifier:     identifier,
		sequenceNumber: sequencenumber,
	}
}

//set the data and ipaddress to send the data to
func (m *EchoMessage) Set(ipAddr string, data []byte) error {
	address, err := net.ResolveIPAddr("ip4", ipAddr)
	if err != nil {
		return err
	}
	m.iaddr = address
	m.data = data
	return nil
}

//send the data!!
func (m *EchoMessage) Send() (responseData []byte, err error) {

	//start creating the icmp echo packet
	payloadLength := len(m.data)

	//8 bytes is the size of the control header
	echoMessagePacket := make([]byte, 8+payloadLength)

	//setting message type
	echoMessagePacket[0] = byte(m.messageType)

	echoMessagePacket[1] = byte(0)

	//checksum
	echoMessagePacket[2] = byte(0)
	echoMessagePacket[3] = byte(0)

	//identifier for the message
	echoMessagePacket[4] = byte(m.identifier & 255)
	echoMessagePacket[5] = byte(m.identifier >> 8)

	//sequence number
	echoMessagePacket[6] = byte(m.sequenceNumber & 255)
	echoMessagePacket[7] = byte(m.sequenceNumber >> 8)

	echoMessagePacket = append(echoMessagePacket, m.data...)

	//caclulate checksum
	var checkSum uint16 = CheckSum(echoMessagePacket[0:])
	var byte1 byte = byte(checkSum & 255)
	var byte2 byte = byte(checkSum >> 8)

	echoMessagePacket[2] = byte1
	echoMessagePacket[3] = byte2

	packagePayLoadLength := len(echoMessagePacket)

	var localAddr string
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {

			if ipnet.IP.IsLoopback() != true && ipnet.IP.To4() != nil {
				localAddr = ipnet.IP.String()
				fmt.Println(localAddr)
			}

		}
	}

	laddr, err := net.ResolveIPAddr("ip", localAddr)
	if err != nil {
		return nil, err

	}

	ipConnection, err := net.DialIP("ip4:icmp", laddr, m.iaddr)
	if err != nil {
		return nil, err
	}
	defer ipConnection.Close()

	//write the message to the connection
	_, err = ipConnection.Write(echoMessagePacket)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, packagePayLoadLength+24)
	_, err = ipConnection.Read(buf)
	if err != nil {
		return nil, err
	}

	//the message recieved contains the ip header
	//remove the header and extract the data from the response packet
	icmpMessage := getIPPayLoad(buf)

	data := icmpMessage[8:]
	return data, nil
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func getIPPayLoad(buf []byte) []byte {

	//get the first 4 bytes of the payload
	internetHeaderLength := (buf[0] & 15)
	ipHeaderLength := int(internetHeaderLength) * 4

	//return payload of ip protocol
	return buf[ipHeaderLength:]
}
