package grest

import "log"

// Try panics when err != nil
func Try(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
