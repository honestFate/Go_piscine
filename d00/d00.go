package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

const (
	flagMean = 1 << iota
	flagMedian
	flagMode
	flagSD
	flagAll = 1<<iota - 1
)

func main() {
	userFlags := flagParsing()
	if userFlags == 0 {
		fmt.Println("Flag error")
		os.Exit(2)
	}
	data := read_data()
	if len(data) == 0 {
		fmt.Println("Input error")
		os.Exit(2)
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})
	printMetrics(data, userFlags)
}

func truncateToHundreds(f float64) float64 {
	return math.Round(f*100) / 100
}

func printMetrics(data []int32, userFlags int32) {
	if userFlags&flagMean != 0 {
		mean := getMean(data)
		if mean-math.Round(mean) == 0 {
			fmt.Printf("Mean: %.1f\n", mean)
		} else {
			mean = truncateToHundreds(mean)
			fmt.Println("Mean:", mean)
		}
	}
	if userFlags&flagMedian != 0 {
		median := getMedian(data)
		if median-math.Round(median) == 0 {
			fmt.Printf("Median: %.1f\n", median)
		} else {
			median = truncateToHundreds(median)
			fmt.Println("Median:", median)
		}
	}
	if userFlags&flagMode != 0 {
		fmt.Println("Mode:", getMode(data))
	}
	if userFlags&flagSD != 0 {
		sd := getSD(data)
		if sd-math.Round(sd) == 0 {
			fmt.Printf("SD: %.1f\n", sd)
		} else {
			sd = truncateToHundreds(sd)
			fmt.Println("SD:", sd)
		}
	}
}

func flagParsing() int32 {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	var userFlags int32 = 0
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.PanicOnError)
	flagMeanTemp := flag.Bool("mean", false, "usage -mean")
	flagMedianTemp := flag.Bool("median", false, "usage -median")
	flagModeTemp := flag.Bool("mode", false, "usage -mode")
	flagSDTemp := flag.Bool("sd", false, "usage -sd")
	flag.Parse()
	if *flagMeanTemp {
		userFlags += flagMean
	}
	if *flagMedianTemp {
		userFlags += flagMedian
	}
	if *flagModeTemp {
		userFlags += flagMode
	}
	if *flagSDTemp {
		userFlags += flagSD
	}
	if userFlags == 0 {
		userFlags = flagAll
	}
	return userFlags
}

func read_data() []int32 {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("Invalid input:", err)
		}
	}()
	numSl := make([]int32, 0, 64)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		num, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}
		if num < -100000 || num > 100000 {
			panic("Value must be in range [-100000:100000]")
		}
		numSl = append(numSl, int32(num))
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return numSl
}

func getMean(slNum []int32) float64 {
	var mean float64
	for i := 0; i < len(slNum); i++ {
		mean += (float64(slNum[i]) - mean) / float64((i + 1))
	}
	return mean
}

func getMedian(slNum []int32) float64 {
	if len(slNum)%2 == 0 {
		return float64(slNum[len(slNum)/2]+slNum[len(slNum)/2-1]) / 2
	}
	return float64(slNum[len(slNum)/2])
}

func getMode(slNum []int32) int32 {
	var maxStreak int32 = 1
	var currentStreak int32 = 1
	mode := slNum[0]
	for i := 1; i < len(slNum); i++ {
		if slNum[i] == slNum[i-1] {
			currentStreak++
		} else {
			if currentStreak > maxStreak {
				maxStreak = currentStreak
				mode = slNum[i-1]
			}
			currentStreak = 1
		}
	}
	if currentStreak > maxStreak {
		maxStreak = currentStreak
		mode = slNum[len(slNum)-1]
	}
	return mode
}

func getSD(slNum []int32) float64 {
	mean := getMean(slNum)
	var dispersion float64
	for i := 0; i < len(slNum); i++ {
		dispersion += math.Pow(float64(float64(slNum[i])-mean), 2)
	}
	dispersion /= float64(len(slNum))
	return math.Sqrt(dispersion)
}
