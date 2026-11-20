package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

// update 是這個 package 的測試自己加的旗標：go test ./projects/chat -run TestSession -update
var update = flag.Bool("update", false, "用目前的逐字稿覆寫 testdata/session.golden")

// TestSession 跑一段完整的對話，逐字稿和 golden file 比對。
func TestSession(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ct := start()
		a := ct.connect("A", 0)
		ct.send(a, "小美")
		b := ct.connect("B", 0)
		ct.send(b, "小美")
		ct.send(b, "阿 明")
		ct.send(b, "阿明")
		ct.send(a, "今天中午吃什麼？")
		ct.send(b, "排骨飯，誰要一起訂？")
		ct.send(b, "/who")
		ct.send(b, "/nick 明哥")
		ct.send(a, "/order 排骨飯")
		ct.connect("C", 0)
		ct.send(a, "/quit")
		ct.shutdown(t)

		golden := filepath.Join("testdata", "session.golden")
		if *update {
			if err := os.WriteFile(golden, []byte(ct.log.String()), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		want, err := os.ReadFile(golden)
		if err != nil {
			t.Fatalf("%v（第一次請加 -update 產生）", err)
		}
		if got := ct.log.String(); got != string(want) {
			t.Errorf("逐字稿和 %s 不一樣：\n--- got\n%s--- want\n%s", golden, got, want)
		}
	})
}

// TestIdle：5 分鐘沒送任何東西就斷線，送了東西就重新計時。時間是假的，測試立刻跑完。
func TestIdle(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ct := start()
		a := ct.connect("A", 0)
		ct.send(a, "小美")
		b := ct.connect("B", 0)
		ct.send(b, "阿明")
		ct.log.Reset()

		time.Sleep(idleTimeout - time.Second)
		ct.send(a, "還有人在嗎？")    // 小美重新計時，阿明沒有
		time.Sleep(time.Second) // 阿明剛好滿 5 分鐘
		ct.settle()
		time.Sleep(idleTimeout - time.Second) // 小美也滿 5 分鐘
		ct.settle()
		ct.shutdown(t)

		want := `A 送出：還有人在嗎？
B 收到：小美：還有人在嗎？
A 收到：* 阿明 太久沒說話，斷線了
B 收到：* 太久沒說話，先斷線了。
B 收到：（連線結束）
A 收到：* 太久沒說話，先斷線了。
A 收到：（連線結束）
伺服器關閉
`
		if got := ct.log.String(); got != want {
			t.Errorf("逐字稿：\n%s\n想要：\n%s", got, want)
		}
	})
}

// TestSlowUser：阿明讀完歡迎訊息就不再讀。他的 writeLoop 卡在第 1 行，待送區再放 outboxSize 行，
// 第 outboxSize+2 行放不下，他就被請出去。
func TestSlowUser(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ct := start()
		a := ct.connect("A", 0)
		ct.send(a, "小美")
		b := ct.connect("B", 2) // 只讀「請輸入暱稱」和「歡迎」兩行
		ct.send(b, "阿明")
		ct.log.Reset()

		for i := 1; i <= outboxSize+2; i++ {
			ct.send(a, "第 "+strconv.Itoa(i)+" 份")
			kicked := strings.Contains(ct.log.String(), "* 阿明 跟不上，斷線了")
			if want := i == outboxSize+2; kicked != want {
				t.Fatalf("送出第 %d 行後，被請出去 = %v，想要 %v", i, kicked, want)
			}
		}

		ct.shutdown(t)
	})
}

// TestWriteTimeout：一行都不讀的人，writeLoop 剛好等 writeTimeout 就放棄。
func TestWriteTimeout(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ct := start()
		a := ct.connect("A", 0)
		ct.send(a, "小美")
		b := ct.connect("B", 2)
		ct.send(b, "阿明")
		ct.send(a, "有人在嗎？") // 阿明的 writeLoop 卡在這一行
		ct.log.Reset()

		time.Sleep(writeTimeout - time.Nanosecond)
		ct.settle()
		if ct.log.Len() != 0 {
			t.Errorf("還沒逾時就有動靜：\n%s", ct.log.String())
		}
		time.Sleep(time.Nanosecond)
		ct.settle()
		if got, want := ct.log.String(), "A 收到：* 阿明 跟不上，斷線了\n"; got != want {
			t.Errorf("逾時後：%q，想要 %q", got, want)
		}
		if _, err := b.conn.Read(make([]byte, 1)); err != io.EOF {
			t.Errorf("阿明那一頭讀到 %v，想要 io.EOF（伺服器已經關掉連線）", err)
		}
		ct.shutdown(t)
	})
}

// TestLineLimit：一行連同換行最多 maxLine 個 byte。
func TestLineLimit(t *testing.T) {
	for _, tt := range []struct {
		size   int // 不含換行
		kicked bool
	}{
		{maxLine - 1, false},
		{maxLine, true},
	} {
		synctest.Test(t, func(t *testing.T) {
			ct := start()
			a := ct.connect("A", 0)
			ct.send(a, "小美")
			b := ct.connect("B", 0)
			ct.send(b, "阿明")
			ct.log.Reset()
			ct.send(a, strings.Repeat("x", tt.size))
			kicked := strings.Contains(ct.log.String(), "A 收到：* 一行太長了，先斷線了。\nA 收到：（連線結束）\nB 收到：* 小美 送了太長的訊息，斷線了\n")
			if kicked != tt.kicked {
				t.Errorf("送出 %d bytes：被請出去 = %v，想要 %v\n%s", tt.size, kicked, tt.kicked, ct.log.String())
			}
			ct.shutdown(t)
		})
	}
}

// TestConcurrentSay：五個人同時各說三句，每個人都要收到其他人的每一句，而且同一個人說的順序不變。
func TestConcurrentSay(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ct := start()
		names := []string{"一號", "二號", "三號", "四號", "五號"}
		var users []*user
		for i, name := range names {
			u := ct.connect(strconv.Itoa(i), 0)
			ct.send(u, name)
			users = append(users, u)
		}
		for _, u := range users {
			u.take()
		}

		var wg sync.WaitGroup
		for _, u := range users {
			wg.Go(func() {
				for i := 1; i <= 3; i++ {
					io.WriteString(u.conn, strconv.Itoa(i)+"\n")
				}
			})
		}
		wg.Wait()
		synctest.Wait()

		for i, u := range users {
			next := map[string]int{} // 說話的人 → 下一句應該是第幾句
			got := u.take()
			for _, line := range got {
				who, n, _ := strings.Cut(line, "：")
				if want := strconv.Itoa(next[who] + 1); n != want {
					t.Errorf("%s 收到 %q，想要 %s 的第 %s 句", names[i], line, who, want)
				}
				next[who]++
			}
			if len(got) != 3*(len(names)-1) {
				t.Errorf("%s 收到 %d 行，想要 %d 行", names[i], len(got), 3*(len(names)-1))
			}
		}
		ct.shutdown(t)
	})
}

func TestCheckNick(t *testing.T) {
	tests := []struct{ nick, wantErr string }{
		{"小美", ""},
		{"Amy2", ""},
		{"一二三四五六七八九十一二", ""},
		{"", "暱稱不能是空的"},
		{"一二三四五六七八九十一二三", "暱稱最多 12 個字"},
		{"阿 明", `暱稱只能用文字和數字，不能有 ' '`},
		{"小美\x1b", `暱稱只能用文字和數字，不能有 '\x1b'`},
		{"\xff", `暱稱只能用文字和數字，不能有 '�'`},
	}
	for _, tt := range tests {
		err := checkNick(tt.nick)
		if got := fmt.Sprint(err); (err == nil) != (tt.wantErr == "") || (err != nil && got != tt.wantErr) {
			t.Errorf("checkNick(%q) = %v，想要 %q", tt.nick, err, tt.wantErr)
		}
	}
}

// TestRenameSame：改成自己現在的暱稱什麼都不做，不能回「有人用了」，也不用通知大家。
func TestRenameSame(t *testing.T) {
	r := NewRoom()
	m, err := r.Join("小明")
	if err != nil {
		t.Fatal(err)
	}
	<-m.out // 歡迎訊息
	if err := r.Rename(m, "小明"); err != nil {
		t.Errorf("Rename 成自己的暱稱 = %v，想要 nil", err)
	}
	select {
	case line := <-m.out:
		t.Errorf("不應該廣播，卻收到 %q", line)
	default:
	}
}

func Example_clean() {
	for _, s := range []string{
		"\x1b[2J\x1b[31m大家好\x1b[0m", // 清掉畫面、改成紅字
		"便當\a\a\a到了",                // 響鈴
		"排骨\xff飯",                   // 不合法的 UTF-8
	} {
		fmt.Printf("%q\n", clean(s))
	}
	// Output:
	// "[2J[31m大家好[0m"
	// "便當到了"
	// "排骨�飯"
}
