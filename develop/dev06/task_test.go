package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"
)

func TestMainFunction(t *testing.T) {
	inputFile := "input.txt"
	inputBytes, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Fatalf("Error reading input file: %v", err)
	}
	testTable := []struct {
		name      string
		command   []string
		expectRes string
	}{
		{
			name:      "OK1",
			command:   []string{"-f", "1"},
			expectRes: "Name:\nItem1:\nItem2:\nItem3:\nItem4:\nItem5:\nItem6:\nline\n",
		},
		{
			name:      "OK2",
			command:   []string{"-s", "-f", "2"},
			expectRes: "Company:\nCompany1:\nCompany2:\nCompany3:\nCompany4:\nCompany5:\nCompany6:\n",
		},
		{
			name:      "OK3",
			command:   []string{"-d", ":", "-f", "3"},
			expectRes: "\tPrice\n\t60000$\n\t70000$\n\t90000$\n\t90000$\n\tGanre1\n\t41000$\nline\n",
		},
		{
			name:      "OK4",
			command:   []string{"-f", "5"},
			expectRes: "Multiplayer\nYes\nYes\nYes\nYes\n\nYes\nline\n",
		},
	}

	for _, testCase := range testTable {
		cmd := exec.Command("go", "run", "task.go")
		cmd.Args = append(cmd.Args, testCase.command...)
		cmd.Stdin = bytes.NewReader(inputBytes)
		var output bytes.Buffer
		cmd.Stdout = &output
		err = cmd.Run()
		if err != nil {
			t.Fatalf("Command execution failed: %v", err)
		}
		if output.String() != testCase.expectRes {
			t.Errorf("Name: %s\nExpected:\n %s\n Got:\n %s", testCase.name, testCase.expectRes, output.String())
		}
	}
}
