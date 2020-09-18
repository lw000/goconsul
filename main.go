package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// 抽象观察者
type IObserver interface {
	Notify()
}

// 抽象被观察者
type ISubject interface {
	Start() error
	Stop()
	AddObject(observers ...IObserver)
	NotifyObservers()
}

// 观察者
type Observer struct {
	Name string
}

func (o *Observer) Notify() {
	time.Sleep(time.Second * time.Duration(rand.Intn(5)))
	log.Println("notify:", o.Name)
}

// 被观察者
type Subject struct {
	wg        sync.WaitGroup
	observers []IObserver
}

func (s *Subject) Start() error {

	return nil
}

func (s *Subject) Stop() {
	s.wg.Wait()
}

func (s *Subject) AddObservers(observers ...IObserver) {
	s.observers = append(s.observers, observers...)
}

func (s *Subject) NotifyObservers() {
	for k := range s.observers {
		// k := k
		s.wg.Add(1)
		go func(o IObserver) {
			defer func() {
				if x := recover(); x != nil {
					log.Println(x)
				}
				s.wg.Done()
			}()
			o.Notify()
		}(s.observers[k])
	}
}

func main() {
	var s Subject
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 1000; i++ {
		s.AddObservers(&Observer{fmt.Sprintf("%d", i)})
	}
	s.NotifyObservers()

	s.Stop()
}
