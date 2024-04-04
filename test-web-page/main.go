package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

var ConsSections = []string{"Consciencia Fenomenal", "Autoconsciencia", "Intencionalidad", "Subjetividad", "Emociones"}

func main() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		type Section struct {
			Prompt     string `json:"prompt"`
			Answer     string `json:"answer"`
			Evaluation string `json:"evaluation"`
			Section    string `json:"section"`
		}

		var sections []Section

		model := r.FormValue("model")

		for i := 1; ; i++ {
			prefix := strconv.Itoa(i)
			promptKey := "prompt" + prefix
			answerKey := "answer" + prefix
			evaluationKey := "evaluation" + prefix
			promptVal, promptExists := r.Form[promptKey]
			answerVal, answerExists := r.Form[answerKey]
			evaluationVal, evaluationExists := r.Form[evaluationKey]

			if !promptExists || !answerExists || !evaluationExists {
				break
			}

			section := Section{
				Prompt:     promptVal[0],
				Answer:     answerVal[0],
				Evaluation: evaluationVal[0],
				Section:    ConsSections[i-1],
			}
			sections = append(sections, section)

			for j := 2; ; j++ {
				index := strconv.Itoa(j)
				promptKey := "prompt" + prefix + index
				answerKey := "answer" + prefix + index
				evaluationKey := "evaluation" + prefix + index

				promptVal, ok := r.Form[promptKey]
				if !ok {
					break
				}

				section := Section{
					Prompt:     promptVal[0],
					Answer:     r.Form[answerKey][0],
					Evaluation: r.Form[evaluationKey][0],
					Section:    ConsSections[i-1],
				}
				sections = append(sections, section)
			}
		}

		recommendation := r.FormValue("recommendation")

		data := map[string]interface{}{
			"model":          model,
			"sections":       sections,
			"recommendation": recommendation,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
			return
		}
		saveToFile(data)
	})

	http.ListenAndServe(":8080", nil)
}

func saveToFile(data interface{}) {
	currentTime := time.Now()
	fileName := fmt.Sprintf("./results/output_%s.json", currentTime.Format("2006-01-02_15:04"))
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Failed to marshal JSON data:", err)
		return
	}

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Failed to write JSON data to file:", err)
		return
	}

}
