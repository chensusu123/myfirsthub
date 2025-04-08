package randfuncs

import (
	"fmt"
	"testing"
)

func TestRandom(t *testing.T) {

	circle := 0
	ret := make(map[int32]int32)
	for circle < 1000000 {
		num := RandByWeightV3(nil, map[int32]int32{0: 50, 13: 50, 21: 50, 3: 50, 4: 50, 50: 50, 16: 50, 17: 50, 11: 50, 10: 50}, false)
		ret[num]++
		circle++
	}

	fmt.Println(ret)
}

func TestMap(t *testing.T) {
	circle := 0

	retList := make([]map[int32]int32, 10)
	for i := 0; i < 10; i++ {
		retList[i] = make(map[int32]int32)
	}
	in := map[int32]int32{0: 50, 13: 50, 21: 50, 3: 50, 4: 50, 50: 50, 16: 50, 17: 50, 11: 50, 10: 50}
	for circle < 1000000 {
		var i int32
		for k := range in {
			retList[i][k]++
			i++
		}
		circle++
	}
	for i := 0; i < 10; i++ {
		fmt.Println(i, retList[i])
	}

}
