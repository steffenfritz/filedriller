package filedriller

import "log"

func e(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
