package main

import (
	"bufio"
	"cmp"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"
)

const (
	idleTimeout  = 5 * time.Minute  // 多久沒送任何東西就斷線
	writeTimeout = 10 * time.Second // 寫一行最多等多久
	maxLine      = 512              // 一行（連同換行）最多幾個 byte
	maxTries     = 3                // 暱稱最多試幾次
)

// Server 是聊天室的伺服器。
type Server struct {
	room *Room
	log  *slog.Logger
}

// NewServer 建立一個伺服器，裡面有一個空的聊天室。
func NewServer(log *slog.Logger) *Server {
	return &Server{room: NewRoom(), log: log}
}

// Serve 接受連線，每個連線交給一個 goroutine，直到 ctx 取消或 ln 出錯。
// 回傳前會請所有人離開，並等每個連線都結束。ctx 取消時回傳 nil。
func (s *Server) Serve(ctx context.Context, ln net.Listener) error {
	stop := context.AfterFunc(ctx, func() { ln.Close() }) // 讓卡在 Accept 的迴圈醒來
	defer stop()
	var wg sync.WaitGroup
	defer wg.Wait()
	for {
		conn, err := ln.Accept()
		if err != nil {
			s.room.Close("* 伺服器要關了，再見。")
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		wg.Go(func() { s.handle(ctx, conn) })
	}
}

// handle 處理一個連線：先問暱稱，再一行一行讀，直到對方離開或被請出去。
// 寫的部分交給另一個 goroutine（writeLoop），讀和寫互不等待。
func (s *Server) handle(ctx context.Context, conn net.Conn) {
	addr := conn.RemoteAddr().String()
	s.log.Info("連線", "addr", addr)
	sc := bufio.NewScanner(conn)
	sc.Buffer(make([]byte, 0, maxLine), maxLine)

	// 還在輸入暱稱時伺服器要關了：直接斷線。進了聊天室之後改由 Room.Close 道別。
	stop := context.AfterFunc(ctx, func() { conn.Close() })
	m, err := s.join(conn, sc)
	stop()
	if err != nil {
		conn.Close()
		s.log.Info("離線", "addr", addr, "err", err)
		return
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		s.writeLoop(m, conn)
	}()

	err = s.chat(m, conn, sc)
	why := "離開了"
	switch {
	case errors.Is(err, os.ErrDeadlineExceeded):
		s.room.Tell(m, "* 太久沒說話，先斷線了。")
		why = "太久沒說話，斷線了"
	case errors.Is(err, bufio.ErrTooLong):
		s.room.Tell(m, "* 一行太長了，先斷線了。")
		why = "送了太長的訊息，斷線了"
	}
	s.room.Leave(m, why) // 已經被請出去的話什麼都不做
	<-done               // writeLoop 寫完剩下的訊息、關掉連線之後才算結束
	s.log.Info("離線", "addr", addr, "err", err)
}

// join 問暱稱，最多問 maxTries 次。這時 writeLoop 還沒開始，直接寫給對方。
func (s *Server) join(conn net.Conn, sc *bufio.Scanner) (*member, error) {
	prompt := "歡迎來到午餐聊天室！請輸入暱稱："
	for range maxTries {
		if err := writeLine(conn, prompt); err != nil {
			return nil, err
		}
		line, err := readLine(conn, sc)
		if err != nil {
			return nil, err
		}
		m, err := s.room.Join(strings.TrimSpace(line))
		if err == nil || errors.Is(err, errClosed) {
			return m, err
		}
		prompt = "* " + err.Error() + "，請再輸入一次："
	}
	writeLine(conn, "* 試太多次了，再見。")
	return nil, errors.New("暱稱試太多次")
}

// chat 一行一行讀 m 送來的東西：一般的話交給其他人，/ 開頭的是指令。
// 對方 /quit 時回傳 nil，其他情況回傳讀取的錯誤。
func (s *Server) chat(m *member, conn net.Conn, sc *bufio.Scanner) error {
	for {
		line, err := readLine(conn, sc)
		if err != nil {
			return err
		}
		line = strings.TrimSpace(clean(line))
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "/") {
			s.room.Say(m, line)
			continue
		}
		cmd, arg, _ := strings.Cut(line, " ")
		switch cmd {
		case "/quit":
			return nil
		case "/who":
			s.room.Who(m)
		case "/nick":
			if err := s.room.Rename(m, strings.TrimSpace(arg)); err != nil {
				s.room.Tell(m, "* "+err.Error())
			}
		case "/help":
			s.room.Tell(m, "* 指令：/who 看誰在線上、/nick 新暱稱 改名、/quit 離開")
		default:
			s.room.Tell(m, "* 沒有 "+cmd+" 這個指令，輸入 /help 看指令")
		}
	}
}

// writeLoop 把 m.out 裡的訊息依序寫給對方。out 被關閉（離開、被請出去、伺服器要關了）時，
// 先把剩下的寫完再關掉連線；handler 卡在讀取的 Scan 會因此結束。
func (s *Server) writeLoop(m *member, conn net.Conn) {
	defer conn.Close()
	for line := range m.out {
		if err := writeLine(conn, line); err != nil {
			why := "離開了"
			if errors.Is(err, os.ErrDeadlineExceeded) {
				why = "跟不上，斷線了"
			}
			s.room.Leave(m, why)
			return
		}
	}
}

// readLine 讀一行，最多等 idleTimeout。telnet 送來的 \r\n，ScanLines 會把 \r 一起拿掉。
func readLine(conn net.Conn, sc *bufio.Scanner) (string, error) {
	conn.SetReadDeadline(time.Now().Add(idleTimeout))
	if sc.Scan() {
		return sc.Text(), nil
	}
	return "", cmp.Or(sc.Err(), io.EOF) // Err 在讀到結尾時是 nil
}

// writeLine 寫一行給對方，最多等 writeTimeout。
func writeLine(conn net.Conn, line string) error {
	conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	_, err := io.WriteString(conn, line+"\n")
	return err
}

// clean 拿掉控制字元，把不合法的 UTF-8 換成 U+FFFD。不這麼做的話，有人可以送出
// ESC 開頭的控制碼，改掉其他人終端機的顏色、甚至清掉整個畫面。
func clean(s string) string {
	return strings.Map(func(c rune) rune {
		if unicode.IsControl(c) {
			return -1
		}
		return c
	}, strings.ToValidUTF8(s, "�"))
}
