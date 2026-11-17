//go:build solution

package pointers

func Swap(a, b *int) {
	*a, *b = *b, *a
}

func Clamp(p *int, lo, hi int) {
	if p == nil {
		return
	}
	*p = max(lo, min(*p, hi))
}

func ApplyPatch(p *Profile, patch ProfilePatch) {
	if patch.Name != nil {
		p.Name = *patch.Name
	}
	if patch.Phone != nil {
		p.Phone = *patch.Phone
	}
	if patch.Points != nil {
		p.Points = *patch.Points
	}
}

// Cheapest 要取 slice 元素本身的位址 &items[i]；range 的 it 只是元素的複本。
func Cheapest(items []Item) *Item {
	var best *Item
	for i := range items {
		if best == nil || items[i].Price < best.Price {
			best = &items[i]
		}
	}
	return best
}

func PushBack(head *Node, name string) *Node {
	n := &Node{Name: name}
	if head == nil {
		return n
	}
	last := head
	for last.Next != nil {
		last = last.Next
	}
	last.Next = n
	return head
}

func Names(head *Node) []string {
	var names []string
	for n := head; n != nil; n = n.Next {
		names = append(names, n.Name)
	}
	return names
}

// Remove 要分開處理「移除第一個人」，因為那時回傳的第一個人會變。
func Remove(head *Node, name string) *Node {
	if head == nil {
		return nil
	}
	if head.Name == name {
		return head.Next
	}
	for prev := head; prev.Next != nil; prev = prev.Next {
		if prev.Next.Name == name {
			prev.Next = prev.Next.Next
			break
		}
	}
	return head
}

// Duplicate 複製 struct 只會複製 Note 這個位址，要另外配置一個新的字串給複本。
func Duplicate(o Order, newID int) Order {
	dup := o
	dup.ID = newID
	if o.Note != nil {
		dup.Note = new(*o.Note)
	}
	return dup
}
