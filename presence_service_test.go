package main

import (
	"testing"
	"time"
)

func TestDueSoon(t *testing.T) {
	now := time.Date(2026, 9, 9, 8, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		learners []Learner
		want     []string
	}{
		{"keeps only next 48 hours", []Learner{{ID: "a", Deadline: now.Add(24 * time.Hour)}, {ID: "b", Deadline: now.Add(72 * time.Hour)}, {ID: "c", Deadline: now.Add(-time.Hour)}}, []string{"a"}},
		{"sorts by deadline", []Learner{{ID: "late", Deadline: now.Add(30 * time.Hour)}, {ID: "early", Deadline: now.Add(2 * time.Hour)}}, []string{"early", "late"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dueSoon(tt.learners, now)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d learners, want %d", len(got), len(tt.want))
			}
			for i, id := range tt.want {
				if got[i].ID != id {
					t.Errorf("position %d: got %q, want %q", i, got[i].ID, id)
				}
			}
		})
	}
}
