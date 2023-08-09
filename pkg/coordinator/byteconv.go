package coordinator

import (
	"bytes"
	"encoding/binary"
)

// []uint64和[]byte互转代码
func uint64ToByte(uint64Slice []uint64) []byte {
	buf := new(bytes.Buffer)
	for _, v := range uint64Slice {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			return nil
		}
	}
	return buf.Bytes()
}

func byteToUint64(byteSlice []byte) ([]uint64, error) {
	var uint64Slice []uint64

	buf := bytes.NewBuffer(byteSlice)
	for buf.Len() > 0 {
		var v uint64
		err := binary.Read(buf, binary.LittleEndian, &v)
		if err != nil {
			return nil, err
		}
		uint64Slice = append(uint64Slice, v)
	}
	return uint64Slice, nil
}
