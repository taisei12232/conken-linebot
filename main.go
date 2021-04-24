package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/line/line-bot-sdk-go/linebot"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!\n")
}
func lineHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	sa := option.WithCredentialsFile("serviceAccount.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	hoge := make(map[string]interface{})
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	// データ読み取り
	iter := client.Collection("users").Documents(ctx)
	key := os.Getenv("KITCATCH_SECRET")
	token := os.Getenv("KITCATCH_ACCESS_TOKEN")
	bot, err := linebot.New(
		key,
		token,
	)

	if err != nil {
		log.Fatal(err)
	}

	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				replyMessage := message.Text
				if replyMessage == "a" {
					replyMessage = fmt.Sprintf("aaaaa")
				}
				if replyMessage == "base" {
					for {
						doc, err := iter.Next()
						if err == iterator.Done {
							break
						}
						if err != nil {
							log.Fatalf("Failed to iterate: %v", err)
						}
						//if doc.Data()["date"] == replyMessage {
						hoge = doc.Data()
						//}
					}
					replyMessage = fmt.Sprintln(hoge)
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Print(err)
				}
			case *linebot.StickerMessage:
				replyMessage := fmt.Sprintf(
					"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
func getListenPort() string {
	port := os.Getenv("PORT")
	if port != "" {
		return ":" + port
	}
	log.Println("HELLO!")
	return ":3000"
}
func main() {
	server := http.Server{
		Addr: getListenPort(),
	}
	var file *os.File
	_, err := os.Stat("serviceAccount.json")
	if err != nil {
		file, err = os.Create("serviceAccount.json")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		file, err = os.Open("serviceAccount.json")
		if err != nil {
			log.Fatal(err)
		}
	}
	defer file.Close()
	account := os.Getenv("ACCOUNT")
	content := []byte(account)
	file.Write(content)
	http.HandleFunc("/", handler)
	http.HandleFunc("/callback", lineHandler)
	server.ListenAndServe()
}
