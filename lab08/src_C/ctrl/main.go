package ctrl

import "fmt"

func GetCtrl(data []byte) []byte {
	var sum uint16 = 0
	for _, elem := range data {
		sum += uint16(elem)
	}
	sum = 65535 - sum
	return []byte{byte(sum & 255), byte((sum >> 8) & 255)}
}

func CheckCtrl(data []byte, sumBytes []byte) bool {
	sum := uint16(sumBytes[0]) + (uint16(sumBytes[1]) << 8)
	if len(data) > 18 {
		panic("Data is too big for checking control sum")
	}
	var check uint32 = 0
	for _, elem := range data {
		check += uint32(elem)
	}
	check += uint32(sum)
	return check == 65535
}

func test(name string, correctData []byte, corruptedData []byte, expected bool) {
	if CheckCtrl(corruptedData, GetCtrl(correctData)) != expected {
		panic(name + " failed")
	}
	fmt.Println(name + " passed")
}

func main() {
	test("Test OK simple", []byte{1}, []byte{1}, true)
	test("Test Err simple", []byte{1}, []byte{0}, false)

	big := []byte{255, 255, 255, 255, 255, 255, 255, 255, 255}
	corrupter := append([]byte{}, big...)
	corrupter[3] = 1
	test("Test OK big", big, big, true)
	test("Test Err big", big, corrupter, false)

	random := []byte{124, 32, 7, 234, 23, 57, 195, 234, 122}
	randomCorr := []byte{54, 32, 76, 234, 0, 3}
	test("Test OK random", random, random, true)
	test("Test Err random", random, randomCorr, false)
}
