package qual

import (
	"errors"
	"math"
)

var ErrEmptyScores = errors.New("qual: empty scores")

var ErrInvalidPhred = errors.New("qual: phred score out of range [0,93]")

var ErrWindowTooLarge = errors.New("qual: window larger than data")

func PhredToProb(q int) float64 {
	if q < 0 {
		return 1.0
	}
	return math.Pow(10, float64(-q)/10.0)
}

func ProbToPhred(p float64) int {
	if p <= 0 {
		return 93
	}
	if p >= 1 {
		return 0
	}
	q := int(-10 * math.Log10(p))
	if q > 93 {
		q = 93
	}
	if q < 0 {
		q = 0
	}
	return q
}

func PhredFromASCII(c byte) int {
	return int(c) - 33
}

func ASCIIFromPhred(q int) byte {
	return byte(q + 33)
}

func MeanQuality(scores []int) (float64, error) {
	if len(scores) == 0 {
		return 0, ErrEmptyScores
	}
	sum := 0
	for _, s := range scores {
		sum += s
	}
	return float64(sum) / float64(len(scores)), nil
}

func MeanErrorRate(scores []int) (float64, error) {
	if len(scores) == 0 {
		return 0, ErrEmptyScores
	}
	var sumP float64
	for _, q := range scores {
		sumP += PhredToProb(q)
	}
	return sumP / float64(len(scores)), nil
}

func TrimBWA(scores []int, threshold int) int {
	if len(scores) == 0 {
		return 0
	}
	n := len(scores)
	bestEnd := n
	maxSum := 0
	runSum := 0
	for i := n - 1; i >= 0; i-- {
		runSum += threshold - scores[i]
		if runSum < 0 {
			runSum = 0
		}
		if runSum > maxSum {
			maxSum = runSum
			bestEnd = i
		}
	}
	return bestEnd
}

func SlidingWindowMean(scores []int, windowSize int) ([]float64, error) {
	if len(scores) == 0 {
		return nil, ErrEmptyScores
	}
	if windowSize > len(scores) {
		return nil, ErrWindowTooLarge
	}
	if windowSize <= 0 {
		windowSize = 1
	}
	n := len(scores) - windowSize + 1
	means := make([]float64, n)
	sum := 0
	for i := 0; i < windowSize; i++ {
		sum += scores[i]
	}
	means[0] = float64(sum) / float64(windowSize)
	for i := 1; i < n; i++ {
		sum += scores[i+windowSize-1] - scores[i-1]
		means[i] = float64(sum) / float64(windowSize)
	}
	return means, nil
}

func CountAbove(scores []int, threshold int) int {
	n := 0
	for _, s := range scores {
		if s >= threshold {
			n++
		}
	}
	return n
}

func CountBelow(scores []int, threshold int) int {
	n := 0
	for _, s := range scores {
		if s < threshold {
			n++
		}
	}
	return n
}

func Validate(scores []int) error {
	for _, s := range scores {
		if s < 0 || s > 93 {
			return ErrInvalidPhred
		}
	}
	return nil
}

func PercentAbove(scores []int, threshold int) float64 {
	if len(scores) == 0 {
		return 0
	}
	return float64(CountAbove(scores, threshold)) / float64(len(scores)) * 100
}
