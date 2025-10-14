package workstyle

import (
	"math"
	"sort"
)

type MetricScores struct {
	Creativity    float64 `json:"creativity"`
	Collaboration float64 `json:"collaboration"`
	Action        float64 `json:"action"`
	Risk          float64 `json:"risk"`
	Empathy       float64 `json:"empathy"`
	Vision        float64 `json:"vision"`
	Adaptability  float64 `json:"adaptability"`
}

type Archetype struct {
	Name    string       `json:"name"`
	Metrics MetricScores `json:"metrics"`
}

type ArchetypeResult struct {
	Name    string  `json:"name"`
	Score   float64 `json:"score"`
	Percent float64 `json:"percent"`
}

func cosineSimilarity(a, b MetricScores) float64 {
	aVals := []float64{a.Creativity, a.Collaboration, a.Action, a.Risk, a.Empathy, a.Vision, a.Adaptability}
	bVals := []float64{b.Creativity, b.Collaboration, b.Action, b.Risk, b.Empathy, b.Vision, b.Adaptability}

	var dot, magA, magB float64
	for i := 0; i < len(aVals); i++ {
		dot += aVals[i] * bVals[i]
		magA += aVals[i] * aVals[i]
		magB += bVals[i] * bVals[i]
	}

	if magA == 0 || magB == 0 {
		return 0
	}

	return dot / (math.Sqrt(magA) * math.Sqrt(magB))
}

func CalculateArchetypeScores(user MetricScores, archetypes []Archetype) []ArchetypeResult {
	results := make([]ArchetypeResult, 0, len(archetypes))
	var total float64

	for _, a := range archetypes {
		sim := cosineSimilarity(user, a.Metrics)
		results = append(results, ArchetypeResult{
			Name:  a.Name,
			Score: sim,
		})
		total += sim
	}

	for i := range results {
		if total > 0 {
			results[i].Percent = (results[i].Score / total) * 100
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Percent > results[j].Percent
	})

	return results
}
