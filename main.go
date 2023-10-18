package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

type program struct {
	path     string
	language string
}

type schedule struct {
	program  program
	interval time.Duration
}

func (program *program) setSchedule(interval time.Duration) schedule {
	return schedule{program: *program, interval: interval}
}

func (schedule *schedule) createTicker() *time.Ticker {
	return time.NewTicker(schedule.interval)
}

func (program *program) run() {
	switch program.language {
	case "Python3":
		{
			out, err := exec.Command("python3", program.path).Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
		}
	case "Go":
		{
			out, err := exec.Command("go", "run", program.path).Output()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(out))
		}
	}
}

func main() {
	pythonProgram1 := program{path: "programs/hello_world.py", language: "Python3"}
	scheduledProgram1 := pythonProgram1.setSchedule(5 * time.Second)
	scheduledProgramTicker1 := scheduledProgram1.createTicker()

	pythonProgram2 := program{path: "programs/project_euler_0001.py", language: "Python3"}
	scheduledProgram2 := pythonProgram2.setSchedule(3 * time.Second)
	scheduledProgramTicker2 := scheduledProgram2.createTicker()

	goProgram1 := program{path: "programs/bye_world.go", language: "Go"}
	scheduledProgram3 := goProgram1.setSchedule(6 * time.Second)
	schedledProgramTicker3 := scheduledProgram3.createTicker()

	for {
		select {
		case <-scheduledProgramTicker1.C:
			scheduledProgram1.program.run()
		case <-scheduledProgramTicker2.C:
			scheduledProgram2.program.run()
		case <-schedledProgramTicker3.C:
			scheduledProgram3.program.run()
		}
	}
}
