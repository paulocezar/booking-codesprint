package search

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/paulocezar/booking-codesprint/passions"
)

// SimpleSearchService implements PassionServicesServer interface
type SimpleSearchService struct {
	cityToInt         map[string]int
	passionToInt      map[string]int
	intTocity         map[int]string
	citiesWithPassion map[int][]int
	maxEndorsement    map[int]int
	maxEndorsementP   map[int]int
	endorsements      map[int]map[int]int
}

type candidate struct {
	id    int
	score float64
}

type ByScore []candidate

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].score > a[j].score }

func (s *SimpleSearchService) Search(c context.Context, req *passions.SearchRequest) (*passions.SearchResponse, error) {

	passionRelevance := make(map[int]float64)
	cityPassionRelevance := make(map[int]float64)

	for _, sp := range req.Passions {
		sp = strings.ToLower(sp)
		p, exists := s.passionToInt[sp]
		if !exists {
			continue
		}

		for _, c := range s.citiesWithPassion[p] {
			passionRelevance[c] += float64(s.endorsements[c][p]) / float64(s.maxEndorsementP[p])
			cityPassionRelevance[c] += float64(s.endorsements[c][p]) / float64(s.maxEndorsement[c])
		}
	}

	dc := make([]candidate, 0, len(passionRelevance))
	ds := make([]string, 0)

	for k, v := range passionRelevance {
		score := v
		score += 0.6 * cityPassionRelevance[k]
		dc = append(dc, candidate{id: k, score: score})
	}

	sort.Sort(ByScore(dc))

	for _, ax := range dc {
		ds = append(ds, s.intTocity[ax.id])
	}

	return &passions.SearchResponse{
		Destinations: ds,
	}, nil
}

func NewSimpleSearchServer(dbPath string) (*SimpleSearchService, error) {

	f, err := os.Open(os.ExpandEnv(dbPath))

	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)

	s := SimpleSearchService{}

	s.cityToInt = make(map[string]int)
	s.passionToInt = make(map[string]int)
	s.intTocity = make(map[int]string)
	s.citiesWithPassion = make(map[int][]int)
	s.maxEndorsement = make(map[int]int)
	s.maxEndorsementP = make(map[int]int)
	s.endorsements = make(map[int]map[int]int)

	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
		}

		if _, exists := s.cityToInt[rec[0]]; !exists {
			s.intTocity[len(s.cityToInt)] = rec[0]
			s.cityToInt[rec[0]] = len(s.cityToInt)
			s.endorsements[s.cityToInt[rec[0]]] = make(map[int]int)
		}

		rec[1] = strings.ToLower(rec[1])
		if _, exists := s.passionToInt[rec[1]]; !exists {
			s.passionToInt[rec[1]] = len(s.passionToInt)
			s.citiesWithPassion[s.passionToInt[rec[1]]] = make([]int, 0)
		}

		passion := s.passionToInt[rec[1]]
		city := s.cityToInt[rec[0]]
		endorsement, _ := strconv.Atoi(rec[2])

		s.maxEndorsement[city] = max(endorsement, s.maxEndorsement[city])
		s.maxEndorsementP[passion] = max(endorsement, s.maxEndorsementP[passion])

		if _, seen := s.endorsements[city][passion]; !seen {
			s.citiesWithPassion[passion] = append(s.citiesWithPassion[passion], city)
		}
		s.endorsements[city][passion] += endorsement
	}

	return &s, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
