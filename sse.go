package globalrouter

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

var errEmptySSEData = errors.New("empty SSE data payload")

type StreamEvent[T any] struct {
	Event string
	Data  T
}

type SSEStream[T any] struct {
	response *http.Response
	scanner  *bufio.Scanner
	done     bool

	event string
	data  []string
}

func newSSEStream[T any](res *http.Response) *SSEStream[T] {
	scanner := bufio.NewScanner(res.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	return &SSEStream[T]{
		response: res,
		scanner:  scanner,
	}
}

func (s *SSEStream[T]) Next() (StreamEvent[T], error) {
	var zero StreamEvent[T]
	if s.done {
		return zero, io.EOF
	}
	for s.scanner.Scan() {
		line := strings.TrimRight(s.scanner.Text(), "\r")
		if line == "" {
			if len(s.data) == 0 {
				s.event = ""
				continue
			}
			event, err := s.dispatch()
			if errors.Is(err, errEmptySSEData) {
				continue
			}
			if err == io.EOF {
				s.done = true
			}
			if err != nil {
				return zero, err
			}
			return event, nil
		}
		if strings.HasPrefix(line, ":") {
			continue
		}
		field, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		value = strings.TrimPrefix(value, " ")
		switch field {
		case "event":
			s.event = value
		case "data":
			s.data = append(s.data, value)
		}
	}
	if err := s.scanner.Err(); err != nil {
		s.done = true
		return zero, err
	}
	if len(s.data) > 0 {
		event, err := s.dispatch()
		if errors.Is(err, errEmptySSEData) {
			s.done = true
			return zero, io.EOF
		}
		if err != nil {
			s.done = true
			return zero, err
		}
		return event, nil
	}
	s.done = true
	return zero, io.EOF
}

func (s *SSEStream[T]) Close() error {
	s.done = true
	if s.response == nil || s.response.Body == nil {
		return nil
	}
	return s.response.Body.Close()
}

func (s *SSEStream[T]) dispatch() (StreamEvent[T], error) {
	var zero StreamEvent[T]
	payload := strings.Join(s.data, "\n")
	eventName := s.event
	s.event = ""
	s.data = nil
	if payload == "" {
		return zero, errEmptySSEData
	}
	if payload == "[DONE]" {
		return zero, io.EOF
	}
	var data T
	decoder := json.NewDecoder(bytes.NewBufferString(payload))
	if err := decoder.Decode(&data); err != nil {
		return zero, err
	}
	return StreamEvent[T]{Event: eventName, Data: data}, nil
}
