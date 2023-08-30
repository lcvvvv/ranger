package ranger

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

var raw = `["1","2","3","4","5","6","67","8"]`

func TestJson(t *testing.T) {
	var r = Ranger[int]{}
	err := json.Unmarshal([]byte(raw), &r)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(r.Value())

	buf, err := json.Marshal(&r)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(buf))

}
