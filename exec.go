package main

import (
	"bufio"
	"io"
	"os/exec"
)

func printer(logf io.Writer, reader io.ReadCloser, done chan bool) string {
	ret := ""
	scanner := bufio.NewScanner(reader)
	go func() {
		for scanner.Scan() {
			logf.Write([]byte(scanner.Text() + "\n"))
			ret += scanner.Text() + "\n"
		}
		done <- true
	}()
	return ret
}

func execute(logf io.Writer, c string, p ...string) error {
	cmd := exec.Command(c, p...)
	reader_o, _ := cmd.StdoutPipe()
	reader_e, _ := cmd.StderrPipe()
	done_o := make(chan bool)
	done_e := make(chan bool)
	printer(logf, reader_o, done_o)
	printer(logf, reader_e, done_e)
	cmd.Start()
	<-done_o
	<-done_e
	return cmd.Wait()
}

func ExecuteCommand(c string, p ...string) error {
	return execute(io.Discard, c, p...)
}
