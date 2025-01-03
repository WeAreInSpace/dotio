package packet

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
)

/*
Dot I/O Packet format
packet data length + packet id + packet data
*/

type Outgoing struct {
	Conn *net.TCPConn
}

func (og *Outgoing) Write() OutgoingBuffer {
	return OutgoingBuffer{
		buffer: new(bytes.Buffer),
		packet: new(bytes.Buffer),
		conn:   og.Conn,
	}
}

type OutgoingBuffer struct {
	conn   *net.TCPConn
	buffer *bytes.Buffer
	packet *bytes.Buffer
}

// 4
func (og *OutgoingBuffer) WriteInt32(number int32) {
	binary.Write(og.buffer, binary.BigEndian, number)
}

// 4
func WriteInt32(number int32) []byte {
	tempByte := new(bytes.Buffer)

	binary.Write(tempByte, binary.BigEndian, number)

	return tempByte.Bytes()
}

// 4
func (og *OutgoingBuffer) WriteFloat32(number float32) {
	binary.Write(og.buffer, binary.BigEndian, number)
}

// 8
func (og *OutgoingBuffer) WriteInt64(number int64) {
	binary.Write(og.buffer, binary.BigEndian, number)
}

// 8
func (og *OutgoingBuffer) WriteFloat64(number float64) {
	binary.Write(og.buffer, binary.BigEndian, number)
}

// length
func (og *OutgoingBuffer) WriteString(str string) {
	length := WriteInt32(int32(len(str)))

	og.buffer.Write(length)
	og.buffer.Write([]byte(str))
}

// length
func (og *OutgoingBuffer) WriteByteArray(bytesArray *bytes.Buffer) {
	length := WriteInt32(int32(bytesArray.Len()))

	og.buffer.Write(length)
	og.buffer.ReadFrom(bytesArray)
}

// 1
func (og *OutgoingBuffer) WriteBoolean(boolean bool) {
	if boolean {
		binary.Write(og.buffer, binary.BigEndian, int8(1))
	} else {
		binary.Write(og.buffer, binary.BigEndian, int8(0))
	}
}

func (og *OutgoingBuffer) Sent(id []byte) error {
	og.packet.Write(id)
	for _, buffer := range og.buffer.Bytes() {
		og.packet.WriteByte(buffer)
	}

	packetLength := WriteInt32(int32(og.packet.Len()))

	_, err := og.conn.Write(packetLength)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return err
	}

	_, err = og.conn.Write(og.packet.Bytes())
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return err
	}

	return nil
}
