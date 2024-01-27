package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jamirhan/project_sibling/tclient"
	"github.com/jamirhan/project_sibling/yagptclient"
)

type Steward struct {
	controller *tclient.Controller
}

var (
	brother_name = "Ахмадшах"
	prompt = "Ты - мой брат. Отвечай так, как будто ты в общем чате и отвечаешь на сообщение."
	unexpected_error_text = "ща, пока не могу отвечать"
	telegram_api_url = "https://api.telegram.org"
)

func (s *Steward) HandleNewMessage(m tclient.Message) {
	if s.controller == nil {
		panic("controller is empty")
	}
	if !strings.Contains(m.Text, brother_name) {
		return
	}

	client := yagptclient.ClientImpl{
		Token: os.Getenv("YAGPT_TOKEN"),
		Endpoint: yagptclient.DefaultEndpoint,
		FolderID: os.Getenv("FOLDER_ID"),
	}

	resp, err := client.GenerateResponse(context.Background(), prompt, m.Text)
	if err != nil {
		s.controller.SendMessage(m.Chat.ID, m.MessageID, unexpected_error_text)
		return
	}

	err = s.controller.SendMessage(m.Chat.ID, m.MessageID, resp)
	if err != nil {
		panic(err)
	}
}

func main() {
	var controller *tclient.Controller
	var err error

	controller, err = tclient.CreateController(context.Background(), telegram_api_url, os.Getenv("TELEGRAM_TOKEN"), func(tclient.Chat) tclient.ChatSteward {
		return &Steward{
			controller: controller,
		}
	})

	if err != nil {
		panic(err)
	}

	controller.Start()
}
