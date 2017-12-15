package config

import (
	"encoding/json"
	"log"
	"os"
)

type Section struct {
	Name     string    `json:"name"`
	Host     string    `json:"host"`
	Port     string    `json:"port"`
	Username string    `json:"username"`
	Passwd   string    `json:"passwd"`
	Commands []Command `json:"commands"`
}

type Command struct {
	Cmd         string `json:"command"`
	Commandfile string `json:"command_file"`
}

var Configurations []Section

func LoadConfig(filename string) bool {
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("warning: %v", err)
		return false
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&Configurations); err != nil {
		log.Printf("Decode %v failed", filename)
		return false
	}
	return true
}
