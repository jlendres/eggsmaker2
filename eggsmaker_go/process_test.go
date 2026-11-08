package main

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestStreamCommandDisplaysOutputBeforeProcessFinishes(t *testing.T) {
	test.NewTempApp(t)
	console := widget.NewMultiLineEntry()
	state := &AppState{Log: &LogView{entry: console}}

	finished := make(chan error, 1)
	go func() {
		finished <- state.streamCommand(context.Background(), "test", exec.Command("sh", "-c", "printf partial; sleep 0.7; printf done"))
	}()

	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) {
		if strings.Contains(console.Text, "partial") {
			if strings.Contains(console.Text, "Proceso finalizado") {
				t.Fatal("process finished before partial output was displayed")
			}
			if err := <-finished; err != nil {
				t.Fatalf("streamCommand() error = %v", err)
			}
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("partial output was not displayed while the process was running")
}

func TestLogViewAppendKeepsCursorAtEnd(t *testing.T) {
	test.NewTempApp(t)
	console := widget.NewMultiLineEntry()
	console.SetText("previous\n")
	logView := &LogView{entry: console}
	logView.moveCursorToEnd()
	console.CursorRow = 0
	console.CursorColumn = 0

	logView.Append("next\nlast")

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if console.Text == "previous\nnext\nlast" &&
			console.CursorRow == 2 && console.CursorColumn == 4 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("log append left cursor at (%d, %d), want (2, 4)", console.CursorRow, console.CursorColumn)
}
