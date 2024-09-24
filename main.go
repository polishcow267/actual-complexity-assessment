package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"io"
	"strconv"
	"math"
)

type task struct{
	guessingProbability float64
	results []taskCompletion
}

type taskCompletion struct{
	result int
	tetha int
}

func (t task) printTask(){
	fmt.Println(t.guessingProbability)
	for i := 0; i<len(t.results); i++{
		fmt.Println(t.results[i].result, " - ", t.results[i].tetha)
	}
}

func readFile(filePath string) task{

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
 
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 2
	
	data := task{
		guessingProbability: 0.001388889,
		results: []taskCompletion{},
	}

	for {
		record, e := reader.Read()
		if err == io.EOF {
			break
		}
		if e != nil {
			fmt.Println(e)
			break
		}
		//fmt.Println(record)
		task := parseRecordToTaskCompletion(record)
		data.results = append(data.results, task)
	}
	return data
}

func parseRecordToTaskCompletion(record []string) taskCompletion {
	result, _ := strconv.Atoi(record[0])
	tetha, _ := strconv.Atoi(record[1])
	
	return taskCompletion{result, tetha}
}

func formulaBirnbaum(tetha int, delta float64, guessingProbability float64) float64{
	numerator := math.Exp(1.71*(float64(tetha)-delta))
	denominator := 1 + math.Exp(1.71*(float64(tetha)-delta))
	return guessingProbability + (1 - guessingProbability)*(numerator/denominator)
}

func maximumLikelihoodMethod(data task) float64{
	resultBirnbaum := []float64{}
	resultDelta := -4.0
	for delta:=-4.0; delta<=4; delta+=0.5{
		ratioLikelihood := 0.0
		for _, value := range data.results{
			if value.result == 1 {ratioLikelihood += math.Log(formulaBirnbaum(value.tetha, delta, data.guessingProbability))
			}else {ratioLikelihood += math.Log(1 - formulaBirnbaum(value.tetha, delta, data.guessingProbability))}
		}
		//fmt.Println(ratioLikelihood)
		resultBirnbaum = append(resultBirnbaum, ratioLikelihood)
	}
	maximumLikelihood := resultBirnbaum[0]
	tmpDelta := -4.0
	for i := 1; i<len(resultBirnbaum); i++{
		if resultBirnbaum[i] > maximumLikelihood{
			maximumLikelihood = resultBirnbaum[i]
			resultDelta = tmpDelta
		}
		tmpDelta += 0.5
	}
	return resultDelta
}

func main(){
	filePath := "data.csv"
	data := readFile(filePath)
	//data.printTask()
	fmt.Println(maximumLikelihoodMethod(data))

}