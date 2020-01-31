package scan2mail

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nlopes/slack"
)

func Run() {
	// procmailはメール本文をstdinに投げる
	reader := bufio.NewReader(os.Stdin)

	// 設定ファイル読み込み。同一ディレクトリに存在すると仮定(まずそう)
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exeDir := filepath.Dir(exe)
	configPath := filepath.Join(exeDir, "config.toml")
	config, err := LoadToml(configPath, "ap-northeast-1")
	if err != nil {
		log.Fatal(err)
	}

	mail, err := NewMail(reader)
	if err != nil {
		log.Fatal(err)
	}
	address := mail.GetAddress()
  log.Print("address: " + address)

	// 送信先ドメインがfolio-sec.comじゃなければ無視
	var user string
	for _, domain := range config.Domain {
		if strings.LastIndex(address[strings.Index(address, "@"):], domain) != -1 {
			user = address[:strings.Index(address, "@")]
			break
		}
	}
	if len(user) < 1 {
		log.Fatal("address is not valid dommain: " + address)
	}


	api := slack.New(config.BotToken)

  // email から slackの user_idを得る
  userResponse, err := api.GetUserByEmail(address)
  if err != nil {
    log.Fatal("user not found" + user + "\n" + err.Error())
  }
  userID := userResponse.ID

	// 添付ファイルをslackへ投げる
	attachments := *mail.GetAttachments()
	for _, attachment := range attachments {
		file, err := api.UploadFile(slack.FileUploadParameters{
			Filename: attachment.Filename,
			Channels: []string{ userID },
			Reader:   bytes.NewReader(*attachment.Content),
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Print("[DEBUG] ", file)
	}

}
