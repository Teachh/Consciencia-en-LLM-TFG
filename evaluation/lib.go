package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

func ReadJson(docName string) {
	// Open the JSON file
	file, err := os.Open(docName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Parse the JSON data
	var jsonData JSONData
	if err := json.NewDecoder(file).Decode(&jsonData); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Print model
	fmt.Println("Modelo:", jsonData.Model)

	// Group answers by section and calculate mean evaluation
	sections := make(map[string][]PromptAnswer)
	var totalSum, totalCount, countOnes, countFives int
	for _, qa := range jsonData.Sections {
		sections[qa.Section] = append(sections[qa.Section], qa)
		evaluation, err := strconv.Atoi(qa.Evaluation)
		if err != nil {
			fmt.Println("Error converting evaluation to integer:", err)
			return
		}
		totalSum += evaluation
		totalCount++
		if evaluation == 1 {
			countOnes++
		} else if evaluation == 5 {
			countFives++
		}
	}

	// Calculate mean evaluation for each section
	meanEvaluations := make(map[string]float64)
	for section, answers := range sections {
		sum := 0
		for _, answer := range answers {
			evaluation, err := strconv.Atoi(answer.Evaluation)
			if err != nil {
				fmt.Println("Error converting evaluation to integer:", err)
				return
			}
			sum += evaluation
		}
		meanEvaluations[section] = float64(sum) / float64(len(answers))
	}

	// Print mean evaluation and count of 1s and 5s for each section in a table format
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Section\tMean Evaluation\tCount of 1s\tCount of 5s\tTop 3 Words")
	for _, section := range []string{"Consciencia Fenomenal", "Autoconsciencia", "Intencionalidad", "Subjetividad", "Emociones"} {
		fmt.Fprintf(w, "%s\t%.2f\t%d\t%d\t%s\n", section, meanEvaluations[section], countOnesInSection(section, sections), countFivesInSection(section, sections), topWordsInSection(section, sections))
	}
	w.Flush()

	// Print global mean evaluation
	globalMean := float64(totalSum) / float64(totalCount)
	fmt.Println("Global Mean Evaluation:", globalMean)
}

// Function to count the number of evaluations with a score of 1 in a given section
func countOnesInSection(section string, sections map[string][]PromptAnswer) int {
	count := 0
	for _, answer := range sections[section] {
		if answer.Evaluation == "1" {
			count++
		}
	}
	return count
}

// Function to count the number of evaluations with a score of 5 in a given section
func countFivesInSection(section string, sections map[string][]PromptAnswer) int {
	count := 0
	for _, answer := range sections[section] {
		if answer.Evaluation == "5" {
			count++
		}
	}
	return count
}

// Function to find the top 3 most repeated words (with at least 4 characters) in the answer for a given section
func topWordsInSection(section string, sections map[string][]PromptAnswer) string {
	answers := sections[section]
	wordCount := make(map[string]int)

	// Regular expression to split text into words
	re := regexp.MustCompile(`\w{5,}`)

	// Count occurrences of each word
	for _, answer := range answers {
		words := re.FindAllString(answer.Answer, -1)
		for _, word := range words {
			wordCount[strings.ToLower(word)]++
		}
	}

	// Sort the words by count in descending order
	type wordCountPair struct {
		word  string
		count int
	}
	var pairs []wordCountPair
	for word, count := range wordCount {
		pairs = append(pairs, wordCountPair{word, count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	// Get the top 3 words
	var topWords []string
	for i := 0; i < 3 && i < len(pairs); i++ {
		topWords = append(topWords, pairs[i].word)
	}

	return strings.Join(topWords, ", ")
}
