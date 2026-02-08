package bezierscore

import (
	"errors"
	"math"
)

var (
	ParticipantCountOutOfRangeErr = errors.New("participantCount must be at least 2")
	ScoreMinOutOfRangeErr         = errors.New("scoreMin must be at least 1")
	ScoreMaxOutOfRangeErr         = errors.New("scoreMax must be more than scoreMin")
	CoefficientOutOfRangeErr      = errors.New("coeff must be between 0 and 1 inclusive")
	ExponentOutOfRangeErr         = errors.New("exp must be at least 1")
)

func bezier(from, to, control, alpha float64) float64 {
	return (from * math.Pow(1.0-alpha, 2)) + (alpha * control * 2 * (1.0 - alpha)) + (math.Pow(alpha, 2) * to)
}

type System struct {
	participantCount   uint
	upperBound         float64
	lowerBound         float64
	controlCoefficient float64
	exponent           float64
}

func New(participantCount uint, scoreMin, scoreMax, coeff, exp float64) (*System, error) {
	if participantCount < 2 {
		return nil, ParticipantCountOutOfRangeErr
	}

	if scoreMin < 1 {
		return nil, ScoreMinOutOfRangeErr
	}

	if scoreMax <= scoreMin {
		return nil, ScoreMaxOutOfRangeErr
	}

	if coeff < 0 || coeff > 1 {
		return nil, CoefficientOutOfRangeErr
	}

	if exp < 1 {
		return nil, ExponentOutOfRangeErr
	}

	return &System{
		participantCount:   participantCount,
		upperBound:         scoreMin,
		lowerBound:         scoreMax,
		controlCoefficient: coeff,
		exponent:           exp,
	}, nil
}

func (s *System) alpha(position uint) float64 {
	numerator := 1.0 - float64(position-1)
	denominator := float64(s.participantCount - 1)
	return 1.0 - (numerator / denominator)
}

func (s *System) control() float64 {
	middle := (s.lowerBound + s.upperBound) / 2.0
	return ((1 - s.controlCoefficient) * middle) + (s.controlCoefficient * s.lowerBound)
}

// Score returns the computed Bezier score for any given position in a leaderboard.
//
// position must be at least 1, and at most the participantCount. A value of 1 means first place, and a value of
// participantCount means last place.
//
// See https://dresswithpockets.github.io/2025/10/14/scoring-system.html
//
// example:
//
//	participantCount := 500
//	scoreMin         := 1000.0
//	scoreMax         := 100000.0
//	coeff            := 0.5
//	exp              := 1.33
//	system := bezierscore.New(participantCount, scoreMin, scoreMax, coeff, exp)
//
//	firstPlace := system.Score(1)
//	secondPlace := system.Score(2)
//	lastPlace := system.Score(participantCount)
func (s *System) Score(position uint) (score float64, ok bool) {
	if position == 0 || position > s.participantCount {
		return 0, false
	}

	alpha := s.alpha(position)
	score = bezier(s.lowerBound, s.upperBound, s.control(), alpha)
	return score, true
}

// ScoreAll computes the Bezier score for every index in buf.
//
// len(buf) must equal participantCount.
//
// example:
//
//	participantCount := 500
//	scoreMin         := 1000.0
//	scoreMax         := 100000.0
//	coeff            := 0.5
//	exp              := 1.33
//	system := bezierscore.New(participantCount, scoreMin, scoreMax, coeff, exp)
//
//	buf := make([]float, participantCount)
//	_ := system.ScoreAll(buf)
func (s *System) ScoreAll(buf []float64) (ok bool) {
	if uint(len(buf)) != s.participantCount {
		return false
	}

	for idx := uint(0); idx < uint(len(buf)); idx++ {
		buf[idx], _ = s.Score(idx + 1)
	}

	return true
}

/*

Copyright 2026 dresswithpockets

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/
