package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// os.Interrupt
func main() {
	c := make(chan os.Signal, 10)
	//signal.Notify(c, syscall.SIGTERM)
	//signal.Notify(c, syscall.SIGKILL)
	//signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGQUIT)

	fmt.Println("rec:", <-c)
}
