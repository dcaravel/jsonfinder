package json

import (
	"fmt"
	"sort"
)

type BreadcrumbEntryType int

const (
	Map BreadcrumbEntryType = iota
	Slice
	String
)

type BreadcrumbEntry struct {
	Typ     BreadcrumbEntryType
	Val     any
	Context map[string]string
}

type Breadcrumb struct {
	entries []*BreadcrumbEntry
}

func (b *Breadcrumb) Len() int {
	return len(b.entries)
}

// Push returns true if an object was pushed
func (b *Breadcrumb) Push(value any, typ BreadcrumbEntryType) {
	b.entries = append(b.entries, &BreadcrumbEntry{Val: value, Typ: typ})
}

// peek returns the last entry from the stack
func (b *Breadcrumb) peek() *BreadcrumbEntry {
	return b.entries[len(b.entries)-1]
}

func (b *Breadcrumb) AddContext(key, value string) {
	if b.Len() <= 0 {
		return
	}

	e := b.peek()
	if e.Context == nil {
		e.Context = map[string]string{}
	}

	e.Context[key] = value

}

func (b *Breadcrumb) RemoveContext(key string) {
	if b.Len() <= 1 {
		return
	}
	e := b.peek()
	delete(e.Context, key)
}

// AddContextToParent adds a key and value to the second to last entry in the breadcrumb stack
// When printing the breadcrumbs the context should come with it
func (b *Breadcrumb) AddContextToParent(key, value string) {
	if b.Len() <= 1 {
		return
	}

	e := b.entries[len(b.entries)-2]
	if e.Context == nil {
		e.Context = map[string]string{}
	}

	e.Context[key] = value
}

func (b *Breadcrumb) RemoveContextFromParent(key string) {
	if b.Len() <= 1 {
		return
	}

	e := b.entries[len(b.entries)-2]
	delete(e.Context, key)
}

func (b *Breadcrumb) Pop() any {
	if b.Len() <= 0 {
		return ""
	}

	v := b.peek()
	b.entries = b.entries[:len(b.entries)-1]
	return v.Val
}

func (b *Breadcrumb) AllStringNoIndex() []string {
	cp := make([]string, 0, len(b.entries))
	for _, e := range b.entries {
		if e.Typ == Slice {
			continue
		}
		v := fmt.Sprintf("%v", e.Val)
		cp = append(cp, v)
	}
	return cp
}

func (b *Breadcrumb) AllWithContext() []string {
	cp := make([]string, 0, len(b.entries))
	for _, e := range b.entries {
		context := ""
		if len(e.Context) > 0 {
			context = "["
			started := false
			keys := make([]string, 0, len(e.Context))
			for k := range e.Context {
				keys = append(keys, k)
			}
			sort.Slice(keys, func(i, j int) bool {
				l := keys[i]
				r := keys[j]
				if l == "i" {
					return true
				}
				if r == "i" {
					return false
				}
				return l < r
			})
			for _, k := range keys {
				v := e.Context[k]
				if started {
					context += ","
				}
				started = true
				context += fmt.Sprintf("%s:%s", k, v)
			}
			context += "]"
		}
		cp = append(cp, fmt.Sprintf("%v%s", e.Val, context))
	}
	return cp
}

func (b *Breadcrumb) Entries() []*BreadcrumbEntry {
	out := make([]*BreadcrumbEntry, 0, len(b.entries))
	out = append(out, b.entries...)
	return out
}
