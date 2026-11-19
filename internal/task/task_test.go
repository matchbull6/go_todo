package task_test

import (
	"encoding/json"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/matchbull6/go_todo/internal/task"
)

var (
	t0 = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	t1 = t0.Add(time.Hour)
)

func TestParsePriority(t *testing.T) {
	tests := []struct {
		in      string
		want    task.Priority
		wantErr bool
	}{
		{in: "low", want: task.Low},
		{in: "medium", want: task.Medium},
		{in: "HIGH", want: task.High},
		{in: "-", wantErr: true},
		{in: "", wantErr: true},
		{in: "hihg", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := task.ParsePriority(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriority_String(t *testing.T) {
	tests := []struct {
		p    task.Priority
		want string
	}{
		{task.None, "-"},
		{task.Low, "low"},
		{task.High, "high"},
		{task.Priority(42), "42"},
		{task.Priority(-1), "-1"},
	}
	for _, tt := range tests {
		if got := tt.p.String(); got != tt.want {
			t.Errorf("Priority(%d).String() = %q, want %q", int(tt.p), got, tt.want)
		}
	}
}

func TestPriority_JSON(t *testing.T) {
	data, err := json.Marshal(task.Task{ID: 1, Title: "x", Priority: task.High})
	if err != nil {
		t.Fatal(err)
	}
	var got task.Task
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Priority != task.High {
		t.Errorf("round trip: got %v, want high (json %s)", got.Priority, data)
	}

	if err := json.Unmarshal([]byte(`{"priority":"urgent"}`), &got); err == nil {
		t.Error("unknown priority: want error")
	}
}

func TestList_Add(t *testing.T) {
	var l task.List

	got, err := l.Add("  Task1  ", task.High, t0)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	want := task.Task{ID: 1, Title: "Task1", Priority: task.High, CreatedAt: t0}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}

	got, err = l.Add("Task2", task.None, t1)
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if got.ID != 2 {
		t.Errorf("second ID = %d, want 2", got.ID)
	}
	if len(l.Tasks) != 2 || l.Tasks[0] != want {
		t.Errorf("Tasks = %+v", l.Tasks)
	}
}

func TestList_Add_EmptyTitle(t *testing.T) {
	var l task.List
	_, err := l.Add("   ", task.None, t0)
	if !errors.Is(err, task.ErrEmptyTitle) {
		t.Errorf("err = %v, want ErrEmptyTitle", err)
	}
	if len(l.Tasks) != 0 || l.LastID != 0 {
		t.Errorf("list changed: %+v", l)
	}
}

func TestList_Complete(t *testing.T) {
	var l task.List
	l.Add("a", task.None, t0)

	got, err := l.Complete(1, t0)
	if err != nil {
		t.Fatalf("Complete: %v", err)
	}
	if !got.Done || !got.DoneAt.Equal(t0) {
		t.Errorf("got %+v, want done at %v", got, t0)
	}
	if !l.Tasks[0].Done {
		t.Error("task in list not marked done")
	}

	got, err = l.Complete(1, t1)
	if err != nil {
		t.Fatalf("second Complete: %v", err)
	}
	if !got.DoneAt.Equal(t0) {
		t.Errorf("second Complete changed DoneAt to %v", got.DoneAt)
	}

	if _, err := l.Complete(99, t0); !errors.Is(err, task.ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestList_Delete(t *testing.T) {
	var l task.List
	for _, title := range []string{"a", "b", "c"} {
		l.Add(title, task.None, t0)
	}

	got, err := l.Delete(2)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if got.Title != "b" {
		t.Errorf("deleted %q, want b", got.Title)
	}
	if ids := taskIDs(l); !slices.Equal(ids, []int{1, 3}) {
		t.Errorf("IDs = %v, want [1 3]", ids)
	}

	if _, err := l.Delete(2); !errors.Is(err, task.ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestList_IDsNotReused(t *testing.T) {
	var l task.List
	l.Add("a", task.None, t0)
	l.Add("b", task.None, t0)
	l.Delete(2)

	got, _ := l.Add("c", task.None, t0)
	if got.ID != 3 {
		t.Errorf("ID = %d, want 3", got.ID)
	}
}

func taskIDs(l task.List) []int {
	var ids []int
	for _, t := range l.Tasks {
		ids = append(ids, t.ID)
	}
	return ids
}
