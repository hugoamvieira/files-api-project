package main

import (
	"io/ioutil"
	"math"
	"regexp"
	"strings"
)

// Internal struct that holds data for calculation
type stats struct {
	alphanumericCharsCountPerFile []int
	wordlengths                   []int
}

// Stats is the struct that holds all the statistics the user gets when it hits the `/info` endpoint
type Stats struct {
	FileCount                   int64              `json:"file_count"`
	AlphanumericCharsPercentage map[string]float64 `json:"alphanum_chars_percentage"`
	AlphanumericCharsStdDev     float64            `json:"alphanum_chars_stddev"`
	WordLenAvg                  float64            `json:"word_length_avg"`
	WordLenStdDev               float64            `json:"word_length_stddev"`
	BytesTotal                  int64              `json:"bytes_total"`
	internal                    *stats
}

// GetStats populates the `Stats` struct with data and returns it to the caller
func GetStats(path string) (*Stats, error) {
	s := new(Stats)
	s.AlphanumericCharsPercentage = make(map[string]float64)

	s.internal = new(stats)
	path = RootPath + path
	if strings.HasSuffix(path, "/") {
		path = strings.TrimRight(path, "/")
	}

	fi, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var allFileData string
	for _, f := range fi {
		s.BytesTotal += f.Size()
		if !f.IsDir() {
			s.FileCount++
			fileBytes, err := ioutil.ReadFile(path + "/" + f.Name())
			if err != nil {
				return nil, err
			}

			allFileData += string(fileBytes) + "\n"
			s.getAlphanumCharCount(f.Name(), string(fileBytes))
		}
	}

	s.getWordCount(allFileData)

	_, s.AlphanumericCharsStdDev = calcAvgAndStdDev(s.internal.alphanumericCharsCountPerFile)
	s.WordLenAvg, s.WordLenStdDev = calcAvgAndStdDev(s.internal.wordlengths)

	return s, nil
}

func (s *Stats) getAlphanumCharCount(filename string, data string) {
	if len(data) == 0 {
		return
	}

	// Assuming POSIX definition of alphanumeric (A-Z, a-z, 0-9)
	m := regexp.MustCompile(`[a-zA-Z0-9]`)

	alphanumCount := len(m.FindAllString(data, -1))
	if alphanumCount == 0 {
		return
	}

	s.AlphanumericCharsPercentage[filename] = float64(alphanumCount) / float64(len(data))
	s.internal.alphanumericCharsCountPerFile = append(s.internal.alphanumericCharsCountPerFile, alphanumCount)
}

func (s *Stats) getWordCount(data string) {
	// Assuming a word is defined by character(s) surrounded by whitespace...
	words := strings.Fields(data)

	for _, word := range words {
		s.internal.wordlengths = append(s.internal.wordlengths, len(word))
	}
}

func calcAvgAndStdDev(collection []int) (float64, float64) {
	var sum float64
	for _, n := range collection {
		sum += float64(n)
	}

	if sum == 0 {
		return -1, -1
	}
	avg := sum / float64(len(collection))

	var diffSum float64
	for _, n := range collection {
		sum += float64(n)
		diffSum += math.Pow((float64(n) - avg), 2)
	}

	return avg, math.Sqrt(diffSum / float64(len(collection)))
}
