package research

import (
	"encoding/hex"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
)

var FromTo *sync.Map
var LastCallee *sync.Map

func init() {
	FromTo = &sync.Map{}
	LastCallee = &sync.Map{}
	os.RemoveAll("delegateTx")
	os.MkdirAll("delegateTx", os.ModePerm)
	os.MkdirAll("code", os.ModePerm)
}

func CaptureDelegateCall(caller, callee string, code []byte, blockNum *big.Int, tx int) {
	caller = strings.ToLower(caller)
	callee = strings.ToLower(callee)
	//if v, ok := LastCallee.Load(caller); ok {
	//	c := v.(string)
	//	if c == callee {
	//		return
	//	}
	//}
	//LastCallee.Store(caller, callee)
	key := caller + callee
	if _, ok2 := FromTo.Load(key); !ok2 {
		//log.Println(caller, "->", callee)
		appendFile(caller, callee, blockNum, tx)
		writeCode(caller, code)
		FromTo.Store(key, 1)
	}
}

func appendFile(caller, callee string, blockNum *big.Int, tx int) {
	f, err := os.OpenFile("delegateTx/"+caller,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	text := callee + "," + blockNum.String() + "," + strconv.Itoa(tx) + "\n"
	if _, err := f.WriteString(text); err != nil {
		panic(err)
	}
}

func writeCode(caller string, code []byte) {
	if PathExists(caller) {
		return
	}
	v := hex.EncodeToString(code)
	os.WriteFile("code/"+caller, []byte(v), 0644)
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
