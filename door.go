package main

import (
	"fmt"
	"github.com/stianeikeland/go-rpio"
	"log"
	"time"
)

// FAKE Indicates whether the software is running on the real system or not.
var FAKE = false

// State holds the current state of the door.
type State int

const (
	unknown State = iota
	closed
	closing
	open
	opening
	holding
)

// Door is a structure that defines the interface
type Door struct {
	state State
	// time we detected door open
	openTime time.Time

	controlPin rpio.Pin
	sensorPin  rpio.Pin
}

// NewDoor is a constructor
func NewDoor(controlPinNum, sensorPinNum int) (*Door, error) {
	err := rpio.Open()
	door := &Door{
		openTime:   time.Time{},
		state:       unknown,
		controlPin: rpio.Pin(controlPinNum),
		sensorPin:  rpio.Pin(sensorPinNum),
	}
	if err == nil {
		door.controlPin.Output()
		door.sensorPin.Input()
		door.sensorPin.PullUp()
	} else {
		FAKE = true
		log.Printf("FAKE!!!")
	}

	return door, nil
}

func (door *Door) cleanup() {
	if FAKE {
		return
	}
	rpio.Close()
}

func (door *Door) toggle() {
	log.Printf("toggle_door")
	if FAKE {
		return
	}
	door.controlPin.Low()
	time.Sleep(1000 * time.Millisecond)
	door.controlPin.High()
}

func (door *Door) isOpen() bool {
	if FAKE {
		switch door.state {
		case open:
			return true
		case opening:
			door.state = open
			return false
		case closed:
			return false
		case closing:
			door.state = closed
			return false
		case holding:
			return true
		case unknown:
			return false
		}
	}
	return door.sensorPin.Read() == rpio.Low
}

func (door *Door) status() bool {
	if door.isOpen() {
		if door.openTime.IsZero() {
			door.openTime = time.Now().Round(time.Second)
		}
		return true
	}
	door.openTime = time.Time{} //reset to 0
	return false
}

func (door *Door) monitor(closetime float64) {
	for true {
		if door.status() && (time.Since(door.openTime).Minutes() > closetime && door.state != holding) {
			log.Printf("Door has been open too long, auto-closing.")
			door.toggle()
		}
		time.Sleep(20000 * time.Millisecond)
	}
}

func (door *Door) check() (string, bool) {
	if door.status() {
		duration := time.Now().Round(time.Second).Sub(door.openTime)
		if door.state == holding {
			return fmt.Sprintf("Door has been Holding for %v\n", duration), true
		}
		return fmt.Sprintf("Door has been Open for %v\n", duration), true
	}
	return "Door is Closed", false
}

func (door *Door) open() chan string {
	c := make(chan string)
	go func(c chan string) {
		if door.isOpen() {
			c <- "Open"
		} else {
			door.state = opening
			c <- "Opening"
			door.toggle()
			for !door.isOpen() {
				time.Sleep(1000 * time.Millisecond)
			}
			door.state = open
			c <- "Open"
		}
		close(c)
	}(c)
	return c
}

func (door *Door) close() chan string {
	c := make(chan string)
	go func(c chan string) {
		if !door.isOpen() {
			c <- "Closed"
		} else {
			door.state = closing
			c <- "Closing..."
			door.toggle()
			for door.isOpen() {
				time.Sleep(1000 * time.Millisecond)
			}
			door.state = closed
			c <- "Closed"
		}
		close(c)
	}(c)
	return c
}

func (door *Door) hold() (string, bool) {
	if !door.isOpen() {
		return "Closed", false
	}
	door.state = holding
	return "Holding", true
}
