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
}

func createProgram(name string, path string, language string, interval int) Program {
	program := Program{name, path, language, time.Duration(interval) * time.Second}
	return program
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
	return time.NewTicker(program.Interval)
}

func (program *Program) run() {
	switch program.Language {
	case "python3":
		{
			out, err := exec.Command("python3", program.Path).Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
		}
	case "golang":
		{
			out, err := exec.Command("go", "run", program.Path).Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
		}
	}
}

func main() {
	programs := readSettings()
	for _, program := range programs.Programs {
		ticker := program.createTicker()
		go func(ticker *time.Ticker, program Program) {
			for {
				select {
				case <-ticker.C:
					program.run()
				}
			}
		}(ticker, program)
	}
}
