package zine

import (
	"bytes"
	"encoding/binary"
)

// defaultHeadLen len[4] + msgID[4] 8个字节
const defaultHeadLen uint32 = 8

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) HeadLen() uint32 {
	return defaultHeadLen
}

func (d *DataPack) Pack(msg *Message) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.LittleEndian, msg.MsgLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, msg.MsgID()); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, msg.MsgData()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (d *DataPack) UnPack(data []byte) (*Message, error) {
	buf := bytes.NewReader(data)
	msg := &Message{}
	if err := binary.Read(buf, binary.LittleEndian, &msg.Len); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}
	return msg, nil
}
