package main

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/anmit007/parago"
)

const (
	applicants  = 50000
	selected    = 11
	trials      = 1000
	skillWeight = 0.95
	luckWeight  = 0.05
)

type candidate struct {
	score float64
	index int
}

func LuckOrSkill() {
	start := time.Now()
	source := rand.NewSource(time.Now().UnixNano())
	results, errs := parago.Map(
		make([]int, trials),
		func(_ int) (int, error) {
			return processTrial(source), nil
		},
		parago.WithWorkers(10),
		parago.WithContext(context.Background()),
	)

	if len(errs) > 0 {
		fmt.Println("Errors occurred:", errs)
		return
	}

	totalCommon := 0
	for _, res := range results {
		totalCommon += res
	}

	average := float64(totalCommon) / float64(trials)
	fmt.Printf("\nOn average, %.1f out of 11 candidates selected with luck factor would also be selected based purely on skill\n", average)
	fmt.Printf("Processing time: %v\n", time.Since(start))
}

func processTrial(source rand.Source) int {
	localRand := rand.New(rand.NewSource(source.Int63()))
	skills := make([]float64, applicants)
	lucks := make([]float64, applicants)
	overalls := make([]float64, applicants)

	for j := 0; j < applicants; j++ {
		skills[j] = localRand.Float64() * 100
		lucks[j] = localRand.Float64() * 100
		overalls[j] = skillWeight*skills[j] + luckWeight*lucks[j]
	}
	topOverall := rankCandidates(overalls, selected)
	topSkill := rankCandidates(skills, selected)

	return countIntersection(topOverall, topSkill)
}

func rankCandidates(scores []float64, topN int) []int {
	candidates := make([]candidate, len(scores))
	for i, score := range scores {
		candidates[i] = candidate{score: score, index: i}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	result := make([]int, topN)
	for i := 0; i < topN; i++ {
		result[i] = candidates[i].index
	}
	return result
}

func countIntersection(a, b []int) int {
	set := make(map[int]struct{}, len(a))
	for _, v := range a {
		set[v] = struct{}{}
	}

	count := 0
	for _, v := range b {
		if _, exists := set[v]; exists {
			count++
		}
	}
	return count
}
