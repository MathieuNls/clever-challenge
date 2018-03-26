package main

import (
	"bytes"
	"sync"
)

type BufferLock struct {
	buffer bytes.Buffer
	mutex sync.Mutex
}
