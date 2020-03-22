// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v7"
)

// State struct
type State struct {
	Client *redis.Client
}

// IsStateful check if app runs on stateful mode
func (s *State) IsStateful() bool {
	if os.Getenv("IS_STATEFUL") == "" {
		return false
	}

	return os.Getenv("IS_STATEFUL") == "true"
}

// IsStateless check if app runs on stateless mode
func (s *State) IsStateless() bool {
	return !s.IsStateful()
}

// Init init state driver
func (s *State) Init() error {
	s.Client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	_, err := s.Client.Ping().Result()

	if err != nil {
		return fmt.Errorf("Unable to connect to redis")
	}

	return nil
}

// Get gets the current state
func (s *State) Get() int {
	value := s.Client.Get("toad__state")

	if value.Val() != "" {
		ivalue, _ := strconv.Atoi(value.Val())

		return ivalue
	}

	return 1
}

// Change changes the current state
func (s *State) Change() {
	s.Client.Set("toad__state", strconv.Itoa(s.Get()+1), 0)
}

// Reset resets the current state
func (s *State) Reset() {
	s.Client.Set("toad__state", strconv.Itoa(1), 0)
}

// HostUp mark current host as up
func (s *State) HostUp() {
	host, _ := os.Hostname()

	s.Client.Set(fmt.Sprintf("toad__host_health__%s", host), "up", 0)
}

// HostDown mark current host as down
func (s *State) HostDown() {
	host, _ := os.Hostname()

	s.Client.Set(fmt.Sprintf("toad__host_health__%s", host), "down", 0)
}

// AllUp mark all hosts as up
func (s *State) AllUp() {
	s.Client.Set("toad__host_health", "up", 0)
}

// AllDown mark all hosts as down
func (s *State) AllDown() {
	s.Client.Set("toad__host_health", "down", 0)
}

// IsDown checks if the current host is down
func (s *State) IsDown() bool {
	host, _ := os.Hostname()

	hostHealth := s.Client.Get(fmt.Sprintf("toad__host_health__%s", host))
	allHealth := s.Client.Get("toad__host_health")

	if hostHealth.Val() == "down" || allHealth.Val() == "down" {
		return true
	}

	return false
}
