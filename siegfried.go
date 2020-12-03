package filedriller

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/richardlehane/siegfried"
)

func siegfriedIdent(s *siegfried.Siegfried, inFile string) string {
	var oneFile string
	f, err := os.Open(inFile)
	if err != nil {
		log.Println(err)
		return "err"
	} else {
		defer f.Close()
		fi, _ := f.Stat()
		if fi.Size() == 0 {
			return inFile + ",,,,,,,"
		}
		ids, err := s.Identify(f, "", "")
		if err != nil {
			log.Println(err)
			return "err"
		}
		for _, id := range ids {
			values := id.Values()
			for _, value := range values {
				oneFile += "\"" + value + "\"" + ","
			}
			oneFile = "\"" + inFile + "\",\"" + strconv.Itoa(int(fi.Size())) + "\"," + oneFile[:len(oneFile)-1] // remove last comma
		}
	}

	return oneFile

}
