package sync

import (
	"sync"
)

var (
	waitgrp sync.WaitGroup
)

func Add(count int) {
	waitgrp.Add(count)
}

func Done() {
	waitgrp.Done()
}

func Wait() {
	waitgrp.Wait()
}
