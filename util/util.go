package util

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// common used functions

type Utils struct {
}

// TODO: test this
func ConvertToTime(input int) string {
	// TODO: check if the std go provides  a function like this
	const NUMBEROFSECONFS = 60
	var minutes int = input / NUMBEROFSECONFS
	var hours int = minutes / NUMBEROFSECONFS
	var seconds int = int(minutes % NUMBEROFSECONFS)

	var result string = fmt.Sprintf("%d : %d : %d", hours, minutes, seconds)
	return result
}

func Progress(precent string) (*string, error) {

	result := "0% "
	size := 50

	boundaries, err := strconv.Atoi(precent)

	if nil != err {
		return nil, err
	}

	boundaries = (boundaries / 100) * size

	// create  the bars

	for i := 0; i < size; i++ {
		if i <= boundaries {
			result += "|"
			continue
		}
		result += " "
	}

	// TODO: make % works
	result += fmt.Sprintf(" %s /%%100", precent[0:5])

	return &result, nil
}

// TODO: rename it to open file
func GetStream(path string) (*os.File, error) {

	file, err := os.Open(path)

	if nil != err {
		return nil, err
	}

	return file, nil
}

func FormateTime(dur float64) string {

	duration := time.Duration(int(dur)) * time.Second
	h := duration / time.Hour
	duration -= h * time.Hour
	m := duration / time.Minute
	duration -= m * time.Minute
	s := duration / time.Second
	formatedTime := fmt.Sprintf("%d : %d : %d", int(h), int(m), int(s))

	return formatedTime
}

// wrapper for creating streams
// std::ifstream Util::getStream(std::string path)
// {
//     std::ifstream stream(path);
//     if  (!stream) {
//         throw std::runtime_error("Non - existing PID");
//     }
//     return stream;
// }
