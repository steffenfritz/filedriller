package filedriller

import (
	"log"
	"os"
	"strconv"

	"github.com/richardlehane/siegfried"
)

func siegfriedIdent(s *siegfried.Siegfried, inFile string) (bool, string) {
	var oneFile string
	var resultBool bool
	f, err := os.Open(inFile)
	if err != nil {
		log.Println(err)
		return resultBool, err.Error()
	} else {
		defer f.Close()
		fi, _ := f.Stat()
		if fi.Size() == 0 {
			return resultBool, inFile + ",,,,,,,,"
		}
		ids, err := s.Identify(f, "", "")
		if err != nil {
			ret := inFile + " : " + err.Error()
			log.Println(ret)
			return resultBool, ret
		}
		for _, id := range ids {
			values := id.Values()
			for _, value := range values {
				oneFile += "\"" + value + "\"" + ","
			}
			oneFile = "\"" + inFile + "\",\"" + strconv.Itoa(int(fi.Size())) + "\"," + oneFile[:len(oneFile)-1] // remove last comma
		}
	}

	return true, oneFile

}
