package hdrwatch

import (
	"strconv"
	"strings"
	"time"
)

type MovieMetadata struct {
	Name string
	Year int
	Tags string
}

func extractMovieMetadata(fileName string) MovieMetadata {
	fileNameSplitted := strings.Split(fileName, ".")

	yearFound := false
	movieName := ""
	movieTags := ""
	movieYear := 0
	result := MovieMetadata{}

	for i := len(fileNameSplitted) - 1; i >= 0; i-- {
		// Try to guess year field
		if yearFound == false {
			if len(fileNameSplitted[i]) == 4 {
				year, err := strconv.Atoi(fileNameSplitted[i])
				if err != nil {
					// Field not a year
					continue
				} else {
					if year >= 1950 && year <= time.Now().Local().Year()+1 {
						yearFound = true
						movieYear = year
					}
				}
			}

		} else {
			movieName = strings.Join(fileNameSplitted[0:i+1], " ")
			movieTags = strings.Join(fileNameSplitted[i+2:], " ")
			break
		}
	}

	if len(movieName) == 0 {
		movieName = fileName
	}
	result.Name = movieName
	result.Year = movieYear
	result.Tags = movieTags
	return result
}
