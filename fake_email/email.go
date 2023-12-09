package fakeemail

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"path/filepath"
	"strings"
)

func Start() {
	var (
		mx  string
		err error
	)

	if len(from) < 1 {
		fmt.Println("mail from ?")
		flagCmd.PrintDefaults()
		return
	}

	if len(to) < 1 {
		fmt.Println("mail to ?")
		flagCmd.PrintDefaults()
		return
	}

	if len(attachment) > 0 && len(attachmentName) < 1 {
		attachmentName = filepath.Base(attachment)
	}

	msg := composeMimeMail(to, from, subject, body, attachment, attachmentName)

	if len(host) < 1 {
		mx, err = getMXRecord(to)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		mx = host
	}

	err = smtp.SendMail(mx+":"+port, nil, from, []string{to}, msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Send to", to, "success!")
}

func getMXRecord(to string) (mx string, err error) {
	var e *mail.Address
	e, err = mail.ParseAddress(to)
	if err != nil {
		return
	}

	domain := strings.Split(e.Address, "@")[1]

	var mxs []*net.MX
	mxs, err = net.LookupMX(domain)
	if err != nil {
		return
	}

	for _, x := range mxs {
		mx = x.Host
		return
	}

	return
}

func formatEmailAddress(addr string) string {
	e, err := mail.ParseAddress(addr)
	if err != nil {
		return addr
	}
	return e.String()
}

func encodeRFC2047(str string) string {
	addr := mail.Address{Address: str}
	return strings.Trim(addr.String(), " <>")
}

func makeBoundary() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%X%X%X%X%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func readFile(fileName string) []byte {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func composeMimeMail(to, from, subject, body, attachment, attachmentName string) []byte {
	message := ""
	boundary := ""

	message += fmt.Sprintf("From: %s\r\n", formatEmailAddress(from))
	message += fmt.Sprintf("To: %s\r\n", formatEmailAddress(to))
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	if len(attachment) > 0 {
		boundary = makeBoundary()
		message += fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary)
		message += "\r\n"
		message += fmt.Sprintf("--%s\r\n", boundary)
	}

	message += "Content-Type: text/plain; charset=\"utf-8\"\r\n"
	message += "Content-Transfer-Encoding: base64\r\n"

	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))
	if len(attachment) > 0 {
		message += "\r\n"
		message += fmt.Sprintf("--%s\r\n", boundary)
		message += "Content-Type: text/plain; charset=\"utf-8\"\r\n"
		message += "Content-Transfer-Encoding: base64\r\n"
		message += fmt.Sprintf("Content-Disposition: attachment; filename=%s\r\n", attachmentName)
		message += "\r\n" + base64.StdEncoding.EncodeToString(readFile(attachment))
		message += "\r\n" + fmt.Sprintf("--%s--", boundary)
	}

	return []byte(message)
}
