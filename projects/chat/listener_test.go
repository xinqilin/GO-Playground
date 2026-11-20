package main

import (
	"bufio"
	"context"
	"io"
	"log/slog"
	"net"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
)

// pipeListener 是假的 net.Listener：dial 用 net.Pipe 在記憶體裡接一條線，一頭交給 Accept。
// 不用開 port，而且 net.Pipe 是用 channel 做的，在 synctest 的泡泡裡卡住時算「durably blocked」。
type pipeListener struct {
	conns chan net.Conn
	done  chan struct{}
	once  sync.Once
}

func newPipeListener() *pipeListener {
	return &pipeListener{conns: make(chan net.Conn), done: make(chan struct{})}
}

func (l *pipeListener) Accept() (net.Conn, error) {
	select {
	case c := <-l.conns:
		return c, nil
	case <-l.done:
		return nil, net.ErrClosed
	}
}

func (l *pipeListener) Close() error {
	l.once.Do(func() { close(l.done) })
	return nil
}

func (l *pipeListener) Addr() net.Addr { return pipeAddr{} }

type pipeAddr struct{}

func (pipeAddr) Network() string { return "pipe" }
func (pipeAddr) String() string  { return "pipe" }

// dial 接一條新的線，回傳用戶端那一頭。
func (l *pipeListener) dial() net.Conn {
	server, client := net.Pipe()
	l.conns <- server
	return client
}

// chatTest 在 synctest 的泡泡裡跑一個聊天室，並把每個人送出和收到的東西記成逐字稿。
type chatTest struct {
	ln     *pipeListener
	cancel context.CancelFunc
	served chan error
	users  []*user
	log    strings.Builder
}

// user 是一個連進來的人。背景的 goroutine 一直讀，讀到的行先存在 got。
type user struct {
	name string
	conn net.Conn
	mu   sync.Mutex
	got  []string
}

func start() *chatTest {
	ctx, cancel := context.WithCancel(context.Background())
	ct := &chatTest{ln: newPipeListener(), cancel: cancel, served: make(chan error, 1)}
	go func() { ct.served <- NewServer(slog.New(slog.DiscardHandler)).Serve(ctx, ct.ln) }()
	return ct
}

// connect 讓 name 連進來。readLimit 是他最多讀幾行就不再讀（模擬跟不上的人），0 表示一直讀。
func (ct *chatTest) connect(name string, readLimit int) *user {
	u := &user{name: name, conn: ct.ln.dial()}
	ct.users = append(ct.users, u)
	go func() {
		sc := bufio.NewScanner(u.conn)
		for n := 1; sc.Scan(); n++ {
			u.add(sc.Text())
			if n == readLimit {
				return
			}
		}
		u.add("（連線結束）")
	}()
	ct.note(name + " 連線")
	ct.settle()
	return u
}

func (u *user) add(line string) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.got = append(u.got, line)
}

func (u *user) take() []string {
	u.mu.Lock()
	defer u.mu.Unlock()
	got := u.got
	u.got = nil
	return got
}

// send 讓 u 送出一行，等訊息都送到之後記下每個人收到什麼。
func (ct *chatTest) send(u *user, line string) {
	ct.note(u.name + " 送出：" + line)
	io.WriteString(u.conn, line+"\n") // 對方斷線時會失敗，逐字稿會顯示原因
	ct.settle()
}

// settle 等泡泡裡每個 goroutine 都停下來，再依連線的順序記下每個人收到的行。
func (ct *chatTest) settle() {
	synctest.Wait()
	for _, u := range ct.users {
		for _, line := range u.take() {
			ct.note(u.name + " 收到：" + line)
		}
	}
}

func (ct *chatTest) note(line string) {
	ct.log.WriteString(line + "\n")
}

// shutdown 取消 ctx（就像按下 Ctrl-C），等 Serve 回傳。
func (ct *chatTest) shutdown(t *testing.T) {
	ct.note("伺服器關閉")
	ct.cancel()
	if err := <-ct.served; err != nil {
		t.Errorf("Serve 回傳 %v，想要 nil", err)
	}
	ct.settle()
}
