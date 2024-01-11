package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"html/template"
	"io/ioutil"
	"log"
	"os"
)

type Tables struct {
	Tables []Table `json:"tables"`
}

type Table struct {
	TableName string   `json:"table_name"`
	Columns   []Column `json:"columns"`
}

type Column struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func readTableSchemas(fileName string) Tables {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		os.Exit(1)
	}

	var tables Tables
	err = json.Unmarshal(file, &tables)
	if err != nil {
		fmt.Println("Error unmarshalling JSON data:", err)
		os.Exit(1)
	}

	return tables
}

func generateQuery(tables Tables, question string) string {
	tmpl, err := template.ParseFiles("prompt.txt")
	if err != nil {
		log.Fatalf("Error parsing template: %s", err)
	}

	jsonData, err := json.Marshal(tables)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
	}

	// Convert the JSON bytes to a string and print
	jsonString := string(jsonData)

	input := map[string]string{
		"Question": question,
		"Tables":   jsonString,
	}

	// Execute the template, writing to os.Stdout or another writer.
	var buf bytes.Buffer

	err = tmpl.Execute(&buf, input)
	if err != nil {
		log.Fatalf("Error executing template: %s", err)
	}

	prompt := buf.String()

	llm, err := openai.New()
	if err != nil {
		log.Fatal(err)
	}

	completion, err := llm.Call(context.Background(), prompt, llms.WithMaxTokens(2048))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(completion)

	return completion
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: query 'question for the data'")
		return
	}

	question := os.Args[1]
	tables := readTableSchemas("tables.json")

	generateQuery(tables, question)
}
