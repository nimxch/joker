package learn

import (
	"fmt"
	"hash/crc32"
	"os"
)

func TestFsync(iter int) {
	dir, _ := os.Getwd()
	file_name := dir + "/wal.log"

	file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	defer file.Close()
	text := fmt.Sprintf("This is aa content {%d}\n", iter)
	checksum := fmt.Sprintf("Checksum: %d\n", CalculateCr32(text))
	_, err = file.Write([]byte(text))
	_, err = file.Write([]byte(checksum))

	if err != nil {
		fmt.Println("Error: ", err)
	}
	file.Sync()
	fmt.Println("Written for Iter : ", iter)
}

const (
	IEEE       = 0xedb88320
	Castagnoli = 0x82f63b78
	Koopman    = 0xeb31d82e
)

func CalculateCr32(data string) uint32 {
	return crc32.Checksum([]byte(data), crc32.MakeTable(Koopman))
}
