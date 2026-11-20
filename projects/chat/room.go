package main

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// outboxSize 是每個人最多累積幾行還沒寫出去的訊息。再多就表示他跟不上，請他離開。
const outboxSize = 16

// maxNick 是暱稱最多幾個字。
const maxNick = 12

var (
	errNickTaken = errors.New("這個暱稱有人用了")
	errClosed    = errors.New("聊天室要關了")
)

// Room 記得誰在線上，把訊息交給每一個人。所有欄位都由 mu 保護。
type Room struct {
	mu      sync.Mutex
	members map[string]*member // 暱稱 → 成員
	closed  bool               // 伺服器要關了，不再讓人進來
}

// member 是線上的一個人。out 是還沒寫給他的訊息，由他自己的 writeLoop 取出來寫。
// 只有 Room 會送東西進 out、也只有 Room 會關掉它，而且都在持有 mu 的時候，
// 所以不會送進已經關閉的 channel，也不會關兩次。
type member struct {
	nick string
	out  chan string
}

// NewRoom 建立一個空的聊天室。
func NewRoom() *Room {
	return &Room{members: map[string]*member{}}
}

// Join 讓 nick 進來。暱稱不合規定、有人用了、或聊天室要關了的時候回傳錯誤。
func (r *Room) Join(nick string) (*member, error) {
	if err := checkNick(nick); err != nil {
		return nil, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	switch {
	case r.closed:
		return nil, errClosed
	case r.members[nick] != nil:
		return nil, errNickTaken
	}
	m := &member{nick: nick, out: make(chan string, outboxSize)}
	r.members[nick] = m
	r.broadcastLocked(m, "* "+nick+" 進來了")
	r.tellLocked(m, "* 歡迎，"+nick+"！目前在線上："+r.namesLocked()+"。輸入 /help 看指令。")
	return m, nil
}

// Say 把 m 說的話交給其他人。
func (r *Room) Say(m *member, text string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.hasLocked(m) {
		r.broadcastLocked(m, m.nick+"："+text)
	}
}

// Tell 只告訴 m 一個人。
func (r *Room) Tell(m *member, line string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.hasLocked(m) {
		r.tellLocked(m, line)
	}
}

// Who 告訴 m 目前誰在線上。
func (r *Room) Who(m *member) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.hasLocked(m) {
		r.tellLocked(m, "* 在線上："+r.namesLocked())
	}
}

// Rename 把 m 改名為 nick，並告訴所有人。
func (r *Room) Rename(m *member, nick string) error {
	if err := checkNick(nick); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.hasLocked(m) {
		return nil
	}
	if r.members[nick] != nil {
		return errNickTaken
	}
	old := m.nick
	delete(r.members, old)
	m.nick = nick
	r.members[nick] = m
	r.broadcastLocked(nil, "* "+old+" 改名為 "+nick)
	return nil
}

// Leave 讓 m 離開，告訴其他人「m why」。m 已經不在時什麼都不做，所以重複呼叫也沒關係。
func (r *Room) Leave(m *member, why string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.leaveLocked(m, why)
}

// Close 告訴所有人 notice，然後讓每個人離開；之後不再讓人進來。
func (r *Room) Close(notice string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.closed = true
	r.broadcastLocked(nil, notice)
	for _, m := range r.members {
		r.leaveLocked(m, "")
	}
}

// hasLocked 回報 m 還在不在。被請出去的人，他的 handler 可能還沒發現，還會呼叫 Say；
// 這時不能再送東西進他已經關閉的 out。
func (r *Room) hasLocked(m *member) bool {
	return r.members[m.nick] == m
}

// broadcastLocked 把 line 交給 from 以外的每一個人（from 是 nil 時交給所有人）。
// 只放進待送區、不等待，所以持有鎖的時間很短，一個人卡住不會拖累其他人。
func (r *Room) broadcastLocked(from *member, line string) {
	var slow []*member
	for _, m := range r.members {
		if m != from && !offer(m, line) {
			slow = append(slow, m)
		}
	}
	// map 的順序是隨機的；排序之後，好幾個人同時跟不上時，公告的順序也固定。
	slices.SortFunc(slow, func(a, b *member) int { return strings.Compare(a.nick, b.nick) })
	for _, m := range slow {
		r.leaveLocked(m, "跟不上，斷線了")
	}
}

func (r *Room) tellLocked(m *member, line string) {
	if !offer(m, line) {
		r.leaveLocked(m, "跟不上，斷線了")
	}
}

func (r *Room) leaveLocked(m *member, why string) {
	if !r.hasLocked(m) {
		return
	}
	delete(r.members, m.nick)
	close(m.out) // 他的 writeLoop 會把剩下的寫完，再關掉連線
	if why != "" {
		r.broadcastLocked(nil, "* "+m.nick+" "+why)
	}
}

func (r *Room) namesLocked() string {
	names := slices.Sorted(maps.Keys(r.members))
	return strings.Join(names, "、") + "（" + strconv.Itoa(len(names)) + " 人）"
}

// offer 把 line 放進 m 的待送區，不等待；放不下時回傳 false。
func offer(m *member, line string) bool {
	select {
	case m.out <- line:
		return true
	default:
		return false
	}
}

// checkNick 檢查暱稱：1 到 maxNick 個字，只能是文字和數字。
func checkNick(nick string) error {
	switch n := utf8.RuneCountInString(nick); {
	case n == 0:
		return errors.New("暱稱不能是空的")
	case n > maxNick:
		return fmt.Errorf("暱稱最多 %d 個字", maxNick)
	}
	for _, c := range nick {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return fmt.Errorf("暱稱只能用文字和數字，不能有 %q", c)
		}
	}
	return nil
}
