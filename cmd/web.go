package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// 硬编码的服务器连接信息
var (
	serverIP       = "192.9.253.209"
	serverUsername = "root"
	serverPassword = "tgb12345ujm678"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "启动 Web SSH 服务",
	Long:  `启动一个 Web 服务,允许通过浏览器执行 SSH 命令并返回结果`,
	Run:   runWebServer,
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().IntP("port", "p", 8080, "Web 服务器端口")
}

func runWebServer(cmd *cobra.Command, args []string) {
	port, _ := cmd.Flags().GetInt("port")

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ws", handleWebSocket)

	log.Printf("Web SSH 服务器正在监听端口 %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket 升级失败:", err)
		return
	}
	defer conn.Close()

	log.Println("新的 WebSocket 连接已建立")

	var sshClient *ssh.Client
	var sshSession *ssh.Session
	var stdin io.WriteCloser

	// 直接使用硬编码的连接信息建立 SSH 连接
	sshClient, sshSession, err = connectSSH(serverIP, serverUsername, serverPassword)
	if err != nil {
		log.Printf("SSH 连接失败: %v", err)
		sendMessage(conn, "connection_status", fmt.Sprintf("SSH 连接失败: %v", err))
		return
	}
	defer sshClient.Close()
	defer sshSession.Close()

	stdin, err = sshSession.StdinPipe()
	if err != nil {
		log.Printf("获取 stdin pipe 失败: %v", err)
		sendMessage(conn, "connection_status", fmt.Sprintf("获取 stdin pipe 失败: %v", err))
		return
	}

	log.Println("SSH 连接成功")
	sendMessage(conn, "connection_status", "SSH 连接成功")

	go handleSSHSession(sshSession, conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取 WebSocket 消息失败:", err)
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("解析 JSON 消息失败:", err)
			sendMessage(conn, "connection_status", "解析消息失败")
			continue
		}

		if msg.Type == "input" {
			_, err := stdin.Write([]byte(msg.Data))
			if err != nil {
				log.Printf("写入 SSH 会话失败: %v", err)
			}
		}
	}

	log.Println("WebSocket 连接已关闭")
}

func connectSSH(ip, username, password string) (*ssh.Client, *ssh.Session, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	client, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		session.Close()
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}

func handleSSHSession(session *ssh.Session, wsConn *websocket.Conn) {
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Printf("无法获取 stdout pipe: %v", err)
		return
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		log.Printf("无法获取 stderr pipe: %v", err)
		return
	}

	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := stdout.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("读取 stdout 失败: %v", err)
				}
				break
			}
			sendMessage(wsConn, "output", string(buf[:n]))
		}
	}()

	go func() {
		for {
			buf := make([]byte, 1024)
			n, err := stderr.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("读取 stderr 失败: %v", err)
				}
				break
			}
			sendMessage(wsConn, "output", string(buf[:n]))
		}
	}()

	if err := session.Shell(); err != nil {
		log.Printf("启动 shell 失败: %v", err)
		return
	}

	if err := session.Wait(); err != nil {
		log.Printf("等待 session 结束时发生错误: %v", err)
	}
}

func sendMessage(conn *websocket.Conn, msgType string, data string) {
	msg := Message{Type: msgType, Data: data}
	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("发送 WebSocket 消息失败: %v", err)
	}
}
