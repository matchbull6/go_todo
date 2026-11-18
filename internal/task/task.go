// Package task is the task list model, without I/O.
package task

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Priority of a task; None means not set.
type Priority int

const (
	None Priority = iota
	Low
	Medium
	High
)

var priorityNames = []string{"-", "low", "medium", "high"}

func (p Priority) String() string {
	if 0 <= p && int(p) < len(priorityNames) {
		return priorityNames[p]
	}
	return strconv.Itoa(int(p))
}

func (p Priority) MarshalText() ([]byte, error) { return []byte(p.String()), nil }
func (p *Priority) UnmarshalText(b []byte) error {
	parsed, err := ParsePriority(string(b))
	if err != nil {
		return err
	}
	*p = parsed
	return nil
}

// ParsePriority is case-insensitive and rejects "-".
func ParsePriority(s string) (Priority, error) {
	switch idx := slices.Index(priorityNames, strings.ToLower(s)); idx {
	case -1, 0:
		return None, fmt.Errorf("invalid priority %q: want low, medium or high", s)
	default:
		return Priority(idx), nil
	}
}

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Priority  Priority  `json:"priority,omitzero"`
	Done      bool      `json:"done"`
	DoneAt    time.Time `json:"done_at,omitzero"`
	CreatedAt time.Time `json:"created_at"`
}

// List is the stored task list. IDs are never reused.
type List struct {
	NextID int    `json:"next_id"`
	Tasks  []Task `json:"tasks"`
}

var ErrEmptyTitle = errors.New("task: empty title")
var ErrNotFound = errors.New("task: not found")

func (l *List) Add(title string, p Priority, now time.Time) (Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, ErrEmptyTitle
	}
	l.NextID++
	task := Task{ID: l.NextID, Title: title, Priority: p, CreatedAt: now}
	l.Tasks = append(l.Tasks, task)
	return task, nil
}

func (l *List) index(id int) (int, error) {
	for i, v := range l.Tasks {
		if v.ID == id {
			return i, nil
		}
	}
	return -1, fmt.Errorf("%w: id %d", ErrNotFound, id)
}

// Complete keeps the original DoneAt if the task is already done.
func (l *List) Complete(id int, now time.Time) (Task, error) {
	i, err := l.index(id)
	if err != nil {
		return Task{}, err
	}
	if !l.Tasks[i].Done {
		l.Tasks[i].Done = true
		l.Tasks[i].DoneAt = now
	}
	return l.Tasks[i], nil
}

func (l *List) Delete(id int) (Task, error) {
	i, err := l.index(id)
	if err != nil {
		return Task{}, err
	}
	task := l.Tasks[i]
	l.Tasks = slices.Delete(l.Tasks, i, i+1)
	return task, nil
}
