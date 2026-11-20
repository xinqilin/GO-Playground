package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Example_storeFile() {
	dir, err := os.MkdirTemp("", "todoapi")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	path := filepath.Join(dir, "todos.json")

	s, _ := OpenStore(path)
	s.Add(Todo{Title: "訂明天的便當", Due: "2026-11-22", CreatedAt: fixedNow()})
	s.Add(Todo{Title: "繳電話費", CreatedAt: fixedNow()})
	s.Delete(2)
	data, _ := os.ReadFile(path)
	fmt.Println(string(data))

	s, _ = OpenStore(path) // 重新啟動
	t, _ := s.Add(Todo{Title: "倒垃圾", CreatedAt: fixedNow()})
	fmt.Printf("重新啟動後的新編號：%d\n", t.ID)
	// Output:
	// {
	//   "next_id": 3,
	//   "todos": [
	//     {
	//       "id": 1,
	//       "title": "訂明天的便當",
	//       "done": false,
	//       "due": "2026-11-22",
	//       "created_at": "2026-11-21T09:30:00+08:00"
	//     }
	//   ]
	// }
	// 重新啟動後的新編號：3
}

func TestStoreReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "todos.json")
	s, err := OpenStore(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, title := range []string{"訂便當", "繳電話費", "倒垃圾"} {
		if _, err := s.Add(Todo{Title: title}); err != nil {
			t.Fatal(err)
		}
	}
	s.Update(1, func(t *Todo) { t.Done = true })
	s.Delete(3)

	s, err = OpenStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if got, _ := s.Get(1); !got.Done {
		t.Errorf("重新開啟後，第 1 筆的 Done 不見了：%+v", got)
	}
	if _, err := s.Get(3); !errors.Is(err, ErrNotFound) {
		t.Errorf("第 3 筆刪掉了，Get 應該回傳 ErrNotFound，got %v", err)
	}
	if got, _ := s.Add(Todo{Title: "洗碗"}); got.ID != 4 {
		t.Errorf("新的編號是 %d，想要 4（不能重複用被刪掉的 3）", got.ID)
	}
	entries, _ := os.ReadDir(filepath.Dir(path))
	if len(entries) != 1 {
		t.Errorf("目錄裡有 %d 個檔案，暫存檔應該都換名字或刪掉了", len(entries))
	}
}

// TestStoreSaveFails：存檔失敗時，記憶體裡的資料也不能改。
func TestStoreSaveFails(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "data") // 這個目錄還不存在，寫檔一定失敗
	s, err := OpenStore(filepath.Join(dir, "todos.json"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Add(Todo{Title: "訂便當"}); err == nil {
		t.Fatal("目錄不存在，Add 應該失敗")
	}
	if n := len(s.List()); n != 0 {
		t.Errorf("存檔失敗，記憶體裡卻多了 %d 筆", n)
	}

	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if got, err := s.Add(Todo{Title: "訂便當"}); err != nil || got.ID != 1 {
		t.Errorf("目錄建好之後 Add = %+v, %v；編號應該還是 1", got, err)
	}
}

func TestOpenStoreCorrupt(t *testing.T) {
	path := filepath.Join(t.TempDir(), "todos.json")
	if err := os.WriteFile(path, []byte(`{"todos": [`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := OpenStore(path); err == nil || !strings.Contains(err.Error(), path) {
		t.Errorf("檔案壞掉時應該回報檔名，got %v", err)
	}
}
