package main

import (
	"bytes"
	"os/exec"
	"testing"
)

func TestInputErrors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "ErrorInTableNumber",
			input:  "testdata/errordata/test1.txt",
			output: "!",
		},
		{
			name:   "ErrorInOpenAndCloseTime",
			input:  "testdata/errordata/test2.txt",
			output: "a b",
		},
		{
			name:   "ErrorInHourPrice",
			input:  "testdata/errordata/test3.txt",
			output: "abc",
		},
		{
			name:   "ErrorInClientName",
			input:  "testdata/errordata/test4.txt",
			output: "10:30 1 Alex",
		},
		{
			name:   "ErrorInEventTime",
			input:  "testdata/errordata/test5.txt",
			output: "abc 1 alex",
		},
		{
			name:   "ErrorInOrderOfEvent",
			input:  "testdata/errordata/test6.txt",
			output: "10:20 2 alex 1",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", "main.go", tc.input)
			var buf bytes.Buffer
			cmd.Stdout = &buf
			err := cmd.Run()
			if err == nil {
				t.Fatalf("Expected an error")
			}
			got := buf.String()
			expected := tc.output
			if got != expected {
				t.Errorf("Output did not match an expected.\nGot: %s\nExpected: %s", got, expected)
			}
		})
	}
}

func TestClub(t *testing.T) {
	for _, tc := range []struct {
		name   string
		input  string
		output string
	}{
		{
			name:  "TestFromExample",
			input: "testdata/test1.txt",
			output: `09:00
08:48 1 client1
08:48 13 NotOpenYet
09:41 1 client1
09:48 1 client2
09:52 3 client1
09:52 13 ICanWaitNoLonger!
09:54 2 client1 1
10:25 2 client2 2
10:58 1 client3
10:59 2 client3 3
11:30 1 client4
11:35 2 client4 2
11:35 13 PlaceIsBusy
11:45 3 client4
12:33 4 client1
12:33 12 client4 1
12:43 4 client2
15:52 4 client4
19:00 11 client3
19:00
1 70 05:58
2 30 02:18
3 90 08:01
`,
		},
		{
			name:  "TestAlphabetOrder",
			input: "testdata/test2.txt",
			output: `09:00
10:00 1 rudolf
10:10 2 rudolf 1
10:20 1 alex
10:30 2 alex 2
15:00 11 alex
15:00 11 rudolf
15:00
1 500 04:50
2 500 04:30
3 0 00:00
`,
		},
		{
			name:  "TestClientUnknownSecondEvent",
			input: "testdata/test3.txt",
			output: `09:00
10:00 1 rudolf
10:10 2 alex 1
10:10 13 ClientUnknown
15:00 11 rudolf
15:00
1 0 00:00
`,
		},
		{
			name:  "TestClientUnknownFourthEvent",
			input: "testdata/test4.txt",
			output: `09:00
10:00 1 alex
10:10 4 anton
10:10 13 ClientUnknown
15:00 11 alex
15:00
1 0 00:00
`,
		},
		{
			name:  "TestYouShallNotPass",
			input: "testdata/test5.txt",
			output: `09:00
10:00 1 rudolf
10:10 1 rudolf
10:10 13 YouShallNotPass
15:00 11 rudolf
15:00
1 0 00:00
`,
		},
		{
			name:  "TestSeatOnHisOwnPlace",
			input: "testdata/test6.txt",
			output: `09:00
10:00 1 alex
10:10 1 rudolf
10:20 2 alex 1
10:30 2 rudolf 2
11:00 2 alex 1
11:00 13 PlaceIsBusy
15:00 11 alex
15:00 11 rudolf
15:00
1 500 04:40
2 500 04:30
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", "main.go", tc.input)
			var buf bytes.Buffer
			cmd.Stdout = &buf
			err := cmd.Run()
			if err != nil {
				t.Fatalf("Unxpected error")
			}
			got := buf.String()
			expected := tc.output
			if got != expected {
				t.Errorf("Output did not match an expected.\nGot: %s\nExpected: %s", got, expected)
			}
		})
	}
}
