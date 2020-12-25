package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "regexp"
  "strings"

  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
  "gopkg.in/yaml.v2"
)

type conf struct {
  APIKey string `yaml:"API_KEY"`
}

func main() {
  c, err := readConf("conf.yaml")
  if err != nil {
    log.Fatal(err)
  }

  bot, err := tgbotapi.NewBotAPI(c.APIKey)
  if err != nil {
    log.Panic(err)
  }

  log.Printf("Authorized on account %s", bot.Self.UserName)

  u := tgbotapi.NewUpdate(0)
  u.Timeout = 60
  updates, _ := bot.GetUpdatesChan(u)

  for update := range updates {
    if update.Message == nil {
      continue
    }

    if update.Message.Text == "/start" {
      reply := "Hoi! Stuur mij een zelfstandig naamwoord en ik reageer met een lidwoord. Je kan [hier](www.github.com/sjagoori/dehetbot) lezen hoe ik in elkaar zit."
      message := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
      message.ParseMode = "markdown"
      bot.Send(message)
    } else {
      reply := getLidwoord(loadPage("https://www.welklidwoord.nl/"+update.Message.Text), update.Message.Text)
      message := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
      message.ParseMode = "markdown"
      bot.Send(message)
    }
  }
}

func getLidwoord(content string, word string) (lidwoord string) {
  re := regexp.MustCompile(`(?m)<span(?: [^>]*)?>.?.?.<\/span>`)
  rm := re.FindString(content)
  a := strings.Trim(rm, "<span>")
  b := strings.Trim(a, "</span>")

  if b == "" {
    return "Helaas, we zijn nog niet zo slim is het wel een zelfstandig naamwoord?"
  }

  return b + " " + word
}

func loadPage(url string) (content string) {
  resp, err := http.Get(url)
  if err != nil {
    return "Unable to fetch"
  }
  defer resp.Body.Close()

  responseData, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return "Unable to read"
  }
  responseString := string(responseData)

  return responseString
}

func readConf(filename string) (*conf, error) {
  buf, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  c := &conf{}
  err = yaml.Unmarshal(buf, c)
  if err != nil {
    return nil, fmt.Errorf("in file %q: %v", filename, err)
  }
  return c, nil
}
