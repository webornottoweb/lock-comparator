package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/wontw/lock-comparator/structure"
)

func main() {
	dataLeft, err := ioutil.ReadFile("composer_l.lock")
	if err != nil {
		fmt.Println(err)
		return
	}

	dataRight, err := ioutil.ReadFile("composer_r.lock")
	if err != nil {
		fmt.Println(err)
		return
	}

	var lockFileLeft structure.LockFile
	var lockFileRight structure.LockFile

	err = json.Unmarshal(dataLeft, &lockFileLeft)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(dataRight, &lockFileRight)
	if err != nil {
		fmt.Println(err)
		return
	}

	var result string
	result = lockFileLeft.Merge(lockFileRight).String()

	ioutil.WriteFile("composer_o.lock", []byte(result), 777)
}
