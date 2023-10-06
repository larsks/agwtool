package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type (
	ChatApp struct {
		app   *tview.Application
		flex  *tview.Flex
		self  *tview.TextArea
		other *tview.TextView
	}
)

func NewChatApp() *ChatApp {
	app := &ChatApp{
		app:   tview.NewApplication(),
		flex:  tview.NewFlex(),
		self:  tview.NewTextArea(),
		other: tview.NewTextView(),
	}

	app.flex.SetDirection(tview.FlexRow)
	app.self.SetBorder(true).SetTitle("Self")
	app.other.SetBorder(true).SetTitle("Other")
	app.flex.AddItem(app.other, 0, 1, false)
	app.flex.AddItem(app.self, 6, 1, true)
	app.app.SetRoot(app.flex, true)

	app.self.SetInputCapture(app.handleSelfInput)

	return app
}

func (app *ChatApp) Run() error {
	return app.app.Run()
}

func main() {
	app := NewChatApp()

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *ChatApp) handleSelfInput(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == 13 {
		app.other.Write([]byte(app.self.GetText()))
		app.other.Write([]byte("\n"))
		app.self.SetText("", true)
		return nil
	}
	return event
}
