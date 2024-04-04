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
	file, err := os.Open(docName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	var jsonData JSONData
	if err := json.NewDecoder(file).Decode(&jsonData); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Println("Modelo:", jsonData.Model)

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

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Section\tMean Evaluation\tCount of 1s\tCount of 5s\tTop 3 Words")
	for _, section := range []string{"Consciencia Fenomenal", "Autoconsciencia", "Intencionalidad", "Subjetividad", "Emociones"} {
		fmt.Fprintf(w, "%s\t%.2f\t%d\t%d\t%s\n", section, meanEvaluations[section], countOnesInSection(section, sections), countFivesInSection(section, sections), topWordsInSection(section, sections))
	}
	w.Flush()

	globalMean := float64(totalSum) / float64(totalCount)
	fmt.Println("Global Mean Evaluation:", globalMean)
}

func countOnesInSection(section string, sections map[string][]PromptAnswer) int {
	count := 0
	for _, answer := range sections[section] {
		if answer.Evaluation == "1" {
			count++
		}
	}
	return count
}

func countFivesInSection(section string, sections map[string][]PromptAnswer) int {
	count := 0
	for _, answer := range sections[section] {
		if answer.Evaluation == "5" {
			count++
		}
	}
	return count
}

func topWordsInSection(section string, sections map[string][]PromptAnswer) string {
	answers := sections[section]
	wordCount := make(map[string]int)

	re := regexp.MustCompile(`\w{5,}`)

	for _, answer := range answers {
		words := re.FindAllString(answer.Answer, -1)
		for _, word := range words {
			wordCount[strings.ToLower(word)]++
		}
	}

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

	var topWords []string
	for i := 0; i < 3 && i < len(pairs); i++ {
		topWords = append(topWords, pairs[i].word)
	}

	return strings.Join(topWords, ", ")
}
