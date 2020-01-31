package scan2mail

import (
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
)

// Mail は 添付付きメール
type Mail struct {
	address     string
	attachments *[]Attachment
}

// GetAddress は送信先アドレスを取得する
func (m Mail) GetAddress() string {
	return m.address
}

// GetAttachments は添付ファイルを取得する
func (m Mail) GetAttachments() *[]Attachment {
	return m.attachments
}

// NewMail はメールオブジェクトを返す。引数はメール本体のbyte列のbytes.Reader
func NewMail(r io.Reader) (*Mail, error) {
	message, err := mail.ReadMessage(r)
	if err != nil {
		return nil, err
	}
	mail := &Mail{}
	mail.address, err = getEmailAddress(message)
	if err != nil {
		return nil, err
	}
	mail.attachments, err = newAttachments(message)
	if err != nil {
		return nil, err
	}
	return mail, nil

}

// Attachment は添付ファイル
type Attachment struct {
	Filename    string
	ContentType string
	Content     *[]byte
}

// https://tools.ietf.org/html/rfc2822#section-3.4
// address = name-addr <addr-spec> or addr-spec
// addr-spec = user@domain.name
func getEmailAddress(m *mail.Message) (string, error) {
	address := m.Header.Get("X-Original-To")
	if len(address) < 3 {
		address = m.Header.Get("To")
	}
	if len(address) < 3 {
		return "", errors.New("Mial don't have To address")
	}
	if strings.LastIndex(address, "<") != -1 {
		return address[strings.Index(address, "<")+1 : len(address)-1], nil
	}
	return address, nil
}

func newAttachments(m *mail.Message) (*[]Attachment, error) {
	header := m.Header
	mediaType, params, err := mime.ParseMediaType(header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	// multipartじゃないやつはスキャンデータ付いてない
	if !strings.HasPrefix(mediaType, "multipart/") {
		return nil, errors.New("This mail have not Attachment")
	}
	mimeReader := multipart.NewReader(m.Body, params["boundary"])

	// 添付ファイルリスト作成
	// scanner でスキャンされたデータは PDF or jpeg だと思っている。tiffの可能性も?
	var attachments []Attachment
	for {
    part, err := mimeReader.NextPart()
		if err == io.EOF {
			break
		}
    // base64でなかったら添付ファイルではない
    if part.Header.Get("Content-Transfer-Encoding") != "base64" {
      log.Println(part.Header.Get("Content-Transfer-Encoding"))
      continue
    }
		if err != nil {
			log.Fatal(err)
		}
		decoder := base64.NewDecoder(base64.StdEncoding, part)
		content, err := ioutil.ReadAll(decoder)
		if err != nil {
			log.Fatal(err)
		}
		attachments = append(attachments, Attachment{
			Filename: part.FileName(),
			Content:  &content,
		})
	}
	return &attachments, nil
}
