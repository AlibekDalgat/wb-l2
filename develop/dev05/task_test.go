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
			name:    "OK1",
			command: []string{"-i", "go"},
			expectRes: "drwxrwxr-x  6 alibek alibek  4096 июн 24 23:07  go\n" +
				"drwxrwxr-x 12 alibek alibek  4096 ноя  8 19:47  GolandProjects\n" +
				"drwxrwxr-x  3 alibek alibek  4096 мая 31 13:19  GoProjects\n",
		},
		{
			name:    "OK2",
			command: []string{"-v", "-n", "Project"},
			expectRes: "1:drwxr-xr-x  5 alibek alibek  4096 сен 23 10:06  BurpSuiteCommunity\n" +
				"3:drwxr-xr-x  4 alibek alibek  4096 ноя 17 21:05  Desktop\n" +
				"4:drwxr-xr-x  7 alibek alibek  4096 ноя 15 16:31  Documents\n" +
				"5:drwxr-xr-x  5 alibek alibek 12288 ноя 17 21:04  Downloads\n" +
				"6:drwxrwxr-x  6 alibek alibek  4096 июн 24 23:07  go\n" +
				"9:drwxr-xr-x  3 alibek alibek  4096 июн 29  2022  Music\n" +
				"11:drwxr-xr-x  7 alibek alibek 12288 ноя 17 21:03  Pictures\n" +
				"13:drwxr-xr-x  2 alibek alibek  4096 сен 29  2020  Public\n" +
				"16:drwx------ 17 alibek alibek  4096 июл  2 20:50  snap\n" +
				"17:drwxrwxr-x  3 alibek alibek  4096 июн 29  2021  Steam\n" +
				"18:drwxrwxr-x  2 alibek alibek  4096 июн  9 00:05  tempgit\n" +
				"19:drwxr-xr-x  2 alibek alibek  4096 сен 29  2020  Templates\n" +
				"20:drwxrwxr-x  3 alibek alibek  4096 сен 21  2022  turtle-python\n" +
				"21:drwxrwxr-x  5 alibek alibek  4096 дек  8  2022  venv\n" +
				"22:drwxr-xr-x  2 alibek alibek  4096 окт 31 21:01  Videos\n" +
				"23:drwx------  6 alibek alibek  4096 ноя 15 11:31 'VirtualBox VMs'\n" +
				"24:drwxrwxr-x  3 alibek alibek  4096 ноя 16 16:30  vmware\n",
		},
		{
			name:      "OK3",
			command:   []string{"-c", "Project"},
			expectRes: "7\n",
		},
		{
			name:    "OK4",
			command: []string{"-A", "2", "Go"},
			expectRes: "drwxrwxr-x 12 alibek alibek  4096 ноя  8 19:47  GolandProjects\n" +
				"drwxrwxr-x  3 alibek alibek  4096 мая 31 13:19  GoProjects\n" +
				"drwxr-xr-x  3 alibek alibek  4096 июн 29  2022  Music\n" +
				"drwxrwxr-x  3 alibek alibek  4096 мар  5  2023  PhpstormProjects\n",
		},
		{
			name:    "OK5",
			command: []string{"-C", "2", "-v", "Projects"},
			expectRes: "drwxr-xr-x  5 alibek alibek  4096 сен 23 10:06  BurpSuiteCommunity\n" +
				"drwxrwxr-x  7 alibek alibek  4096 дек 31  2022  CLionProjects\n" +
				"drwxr-xr-x  4 alibek alibek  4096 ноя 17 21:05  Desktop\n" +
				"drwxr-xr-x  7 alibek alibek  4096 ноя 15 16:31  Documents\n" +
				"drwxr-xr-x  5 alibek alibek 12288 ноя 17 21:04  Downloads\n" +
				"drwxrwxr-x  6 alibek alibek  4096 июн 24 23:07  go\n" +
				"drwxrwxr-x 12 alibek alibek  4096 ноя  8 19:47  GolandProjects\n" +
				"drwxrwxr-x  3 alibek alibek  4096 мая 31 13:19  GoProjects\n" +
				"drwxr-xr-x  3 alibek alibek  4096 июн 29  2022  Music\n" +
				"drwxrwxr-x  3 alibek alibek  4096 мар  5  2023  PhpstormProjects\n" +
				"drwxr-xr-x  7 alibek alibek 12288 ноя 17 21:03  Pictures\n" +
				"-rw-rw-r--  1 alibek alibek     0 ноя 18 11:02  Projects2\n" +
				"drwxr-xr-x  2 alibek alibek  4096 сен 29  2020  Public\n" +
				"drwxrwxr-x  5 alibek alibek  4096 июн  6 02:17  PycharmProjects\n" +
				"drwxrwxr-x  4 alibek alibek  4096 янв  6  2023  pythonProject1\n" +
				"drwx------ 17 alibek alibek  4096 июл  2 20:50  snap\n" +
				"drwxrwxr-x  3 alibek alibek  4096 июн 29  2021  Steam\n" +
				"drwxrwxr-x  2 alibek alibek  4096 июн  9 00:05  tempgit\n" +
				"drwxr-xr-x  2 alibek alibek  4096 сен 29  2020  Templates\n" +
				"drwxrwxr-x  3 alibek alibek  4096 сен 21  2022  turtle-python\n" +
				"drwxrwxr-x  5 alibek alibek  4096 дек  8  2022  venv\n" +
				"drwxr-xr-x  2 alibek alibek  4096 окт 31 21:01  Videos\n" +
				"drwx------  6 alibek alibek  4096 ноя 15 11:31 'VirtualBox VMs'\n" +
				"drwxrwxr-x  3 alibek alibek  4096 ноя 16 16:30  vmware\n",
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
