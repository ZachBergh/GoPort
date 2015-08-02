package main

import (
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os/exec"
	"sync"
)

func main() {
	app := cli.NewApp()
	app.Name = "goport"
	app.Usage = "Check your linux's port"
	app.Action = func(c *cli.Context) {
		cmd := `netstat -an | grep LISTEN | egrep '0.0.0.0|:::' | awk '/^tcp/ {print $4}' | awk -F: '{print $2$4}' | sort -n`
		cmdType := "bash"
		out, err := run(cmd, cmdType)
		if err != nil {
			fmt.Printf("[Error]" + cmdType + ":" + err.Error())
		}
		fmt.Println(out)
	}
}

func run(cmd, tp string) (result string, err error) {
	wg := new(sync.WaitGroup)
	wg.Add(1)

	var c *exec.Cmd
	switch tp {
	case "bash":
		c = exec.Command("/bin/sh", "-c", cmd)
		break
	}
	stdout, err := c.StdoutPipe()
	if err != nil {
		return
	}

	stderr, err := c.StderrPipe()
	if err != nil {
		return
	}

	if err = c.Start(); err != nil {
		return
	}

	bytesErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		return
	}

	if len(bytesErr) != 0 {
		err = errors.New("Stderr's Byte != 0")
		return
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return
	}

	if err = c.Wait(); err != nil {
		return
	}
	fmt.Printf("stdout: %s", bytes)
	result = string(bytes)
	wg.Done()
	return
}
