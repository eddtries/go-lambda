package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type Programs struct {
	Programs []Program `json:"programs"`
}

type Program struct {
	Name     string        `json:"name"`
	Path     string        `json:"path"`
	Language string        `json:"language"`
	Interval time.Duration `json:"interval"`
	Output   bool          `json:"output"`
}

func readSettings() Programs {
	content, err := os.ReadFile("./settings.json")
	if err != nil {
		fmt.Println(err)
	}

	var programs Programs
	err = json.Unmarshal(content, &programs)
	if err != nil {
		fmt.Println(err)
	}

	return programs
}

func (program *Program) createTicker() *time.Ticker {
	return time.NewTicker(program.Interval * time.Second)
}

func (program *Program) run() {
	switch program.Language {
	case "Python3":
		{
			out, err := exec.Command("python3", program.Path).Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Ran " + program.Name + " at " + time.Now().String())
			if program.Output {
				fmt.Println(program.Name + ": " + string(out))
			}
		}
	case "Go":
		{
			out, err := exec.Command("go", "run", program.Path).Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Ran " + program.Name + " at " + time.Now().String())
			if program.Output {
				fmt.Println(program.Name + ": " + string(out))
			}
		}
	}
}

func schedule(program Program, done <-chan bool) *time.Ticker {
	ticker := program.createTicker()
	go func() {
		for {
			select {
			case <-ticker.C:
				program.run()
			case <-done:
				return
			}
		}
	}()
	return ticker
}

func main() {
	programs := readSettings()
	done := make(chan bool)

	for _, program := range programs.Programs {
		schedule(program, done)
	}
	time.Sleep(1 * time.Hour)
	close(done)
}
