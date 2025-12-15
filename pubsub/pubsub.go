package pubsub

import (
	"sync"

	"gitlab.com/davesaah/fatch/internal/database"
)

type Subscriber chan database.Log

type PubSub struct {
	mu     sync.RWMutex
	topics map[string]map[Subscriber]struct{} // a set of subscribers to a topic
}

func New() *PubSub {
	return &PubSub{
		topics: make(map[string]map[Subscriber]struct{}),
	}
}

func (ps *PubSub) Subscribe(topic string, buffer int) Subscriber {
	ch := make(Subscriber, buffer)

	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, ok := ps.topics[topic]; !ok {
		ps.topics[topic] = make(map[Subscriber]struct{})
	}

	ps.topics[topic][ch] = struct{}{}
	return ch
}

func (ps *PubSub) Unsubscribe(topic string, ch Subscriber) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	subs, ok := ps.topics[topic]
	if !ok {
		return
	}

	if _, ok := subs[ch]; ok {
		delete(subs, ch)
		close(ch)
	}

	// there cannot be an empty topic
	if len(subs) == 0 {
		delete(ps.topics, topic)
	}
}

func (ps *PubSub) Publish(topic string, log database.Log) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	subs, ok := ps.topics[topic]
	if !ok {
		return
	}

	for ch := range subs {
		ch <- log
	}
}
