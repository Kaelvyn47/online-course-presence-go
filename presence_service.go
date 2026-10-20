package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"
)

// realtime.presence.get is the domain read used by the status endpoint.

type Learner struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Deadline time.Time `json:"deadline"`
}
type OnlineMember struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	LastSeen time.Time `json:"last_seen"`
}
type PresenceData struct {
	Members []OnlineMember `json:"members"`
}
type CourseStatus struct {
	Online  []OnlineMember `json:"online"`
	DueSoon []Learner      `json:"due_soon"`
}

func dueSoon(learners []Learner, now time.Time) []Learner {
	cutoff := now.Add(48 * time.Hour)
	result := make([]Learner, 0)
	for _, learner := range learners {
		if learner.Deadline.After(now) && learner.Deadline.Before(cutoff) {
			result = append(result, learner)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Deadline.Before(result[j].Deadline) })
	return result
}

type service struct {
	client   *InfraiClient
	learners []Learner
	mu       sync.Mutex
}

func (s *service) handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/course/status" {
		http.NotFound(w, r)
		return
	}
	var presence PresenceData
	if err := s.client.Presence("course-presence", &presence); err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	s.mu.Lock()
	due := dueSoon(s.learners, time.Now())
	s.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	_ = jsonEncode(w, CourseStatus{Online: presence.Members, DueSoon: due})
}

func jsonEncode(w http.ResponseWriter, v any) error { return json.NewEncoder(w).Encode(v) }
