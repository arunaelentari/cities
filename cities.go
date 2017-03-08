// This program ranks cities based on a number of criteria.
//
// Criteria have a positive weight. This allows us to compare two cities and say which one is greater.

package main

import (
	"fmt"
	"sort"
	"strings"
	"net/http"
	"crypto/tls"

	"golang.org/x/crypto/acme/autocert"
)

type (
	// cost is the cost of living in a place.
	cost int
	// climate is the quality of typical weather.
	climate int

	// city is a place where we might want to live.
	city struct {
		name       string
		population int
		cost       cost
		climate    climate
	}

	// cities is a collection of city.
	cities []city

	criteria struct {
		weight float64
		name   string      // e.g. "population"
		value  interface{} // this is an int or a cost or a climate
	}
)

const (
	CheapCost cost = 1 + iota
	VeryReasonableCost
	ReasonableCost
	ExpensiveCost
	VeryExpensiveCost
)

const (
	NastyClimate climate = 1 + iota
	PoorClimate
	GoodClimate
	GreatClimate
	PerfectClimate
)

var (
	ClimateDesc = map[climate]string{
		NastyClimate:   "nasty",
		PoorClimate:    "poor",
		GoodClimate:    "good",
		GreatClimate:   "great",
		PerfectClimate: "perfect",
	}
	CostDesc = map[cost]string{
		CheapCost:          "cheap",
		VeryReasonableCost: "very reasonable",
		ReasonableCost:     "reasonable",
		ExpensiveCost:      "expensive",
		VeryExpensiveCost:  "very expensive",
	}
	// TODO: this should be eventually read from a user.
	Cities = cities{
		city{name: "Barcelona", population: 1.6e6, cost: ReasonableCost, climate: GreatClimate},
		city{name: "Seattle", population: 652405, cost: ExpensiveCost, climate: GoodClimate},
		city{name: "New York", population: 8.406e6, cost: ExpensiveCost, climate: GoodClimate},
		city{name: "Copenhagen", population: 562379, cost: ExpensiveCost, climate: PoorClimate},
		city{name: "Stockholm", population: 789024, cost: ExpensiveCost, climate: PoorClimate},
		city{name: "Deviltown", population: 1233567890, cost: VeryExpensiveCost, climate: NastyClimate},
		city{name: "Paradisio", population: 1e6, cost: CheapCost, climate: PerfectClimate},
	}
)

// String returns a description of the climate.
func (c climate) String() string {
	return ClimateDesc[c]
}

// String returns a description of the cost.
func (c cost) String() string {
	return CostDesc[c]
}

// String returns a description of the city.
func (c city) String() string {
	if c.name == "" {
		return "city with empty name, you dummy!"
	}
	p := ""
	if c.population < 1e6 {
		p = fmt.Sprintf("%v", c.population)
		if len(p) > 3 {
			from := len(p) - 3
			p = p[:from] + " " + p[from:]
		}
	} else {
		p = fmt.Sprintf("%.1fM", float64(c.population)/1e6)
	}
	return fmt.Sprintf(
		"%v: %v, cost: %v, climate: %v",
		c.name,
		p,
		c.cost,
		c.climate,
	)
}

// String returns a description of the cities.
func (cs cities) String() string {
	desc := make([]string, len(cs), len(cs))
	for i := range cs {
		desc[i] = fmt.Sprintf("  * %v", cs[i])
	}
	return strings.Join(desc, "\n")
}

func (cs cities) getNames() string {
	names := make([]string, len(cs), len(cs))
	for i := range cs {
		names[i] = cs[i].name
	}
	return strings.Join(names, ", ")

}

// sortBy sorts cities by given criteria.
func (cs cities) sortBy(criteria string) {
	if criteria == "name" {
		sort.Slice(cs, func(i, j int) bool { return cs[i].name < cs[j].name })
	}
	if criteria == "population" {
		sort.Slice(cs, func(i, j int) bool { return cs[i].population < cs[j].population })
	}
	if criteria == "cost" {
		sort.Slice(cs, func(i, j int) bool { return cs[i].cost < cs[j].cost })
	}
	if criteria == "climate" {
		sort.Slice(cs, func(i, j int) bool { return cs[i].climate < cs[j].climate })
	}

}

// getInfo returns information on cities.
func (cs cities) getInfo() string{
	// We want a way to sort by a weighted set of criteria, e.g:
	// Cities.sortByCriteria(criteria{"climate", 2}, criteria{"cost", 1})
	// Or:
	// Cities.sortByCriteria(map[string]int{"climate": 2, "cost": 1})
	//
	// Regardless of which function signature we pick, we then want results
	// where the cities are sorted in ascending order (worst to best), like:
	//
	// The sorted cities by climate (67%) and cost (33%) are:
	// * Deviltown: 1234.6M, cost: very expensive, climate: nasty
  	// * Copenhagen: 562 379, cost: expensive, climate: poor
  	// * Stockholm: 789 024, cost: expensive, climate: poor
	// * New York: 8.4M, cost: expensive, climate: good
	// * Barcelona: 1.6M, cost: reasonable, climate: great
  	// * Paradisio: 1.0M, cost: cheap, climate: perfect

	n := "The sorted cities by name are:\n"
	cs.sortBy("name")
	n += cs.String()
	n += "\nThe sorted cities by population are:\n"
	cs.sortBy("population")
	n += cs.String()
	n += "\nThe sorted cities by cost are:\n"
	cs.sortBy("cost")
	n += cs.String()
	n += "\nThe sorted cities by climate are:\n"
	cs.sortBy("climate")
	n += cs.String()
	return n
}

// indexHandler writes the http reply to the request for the index page.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("You are all my minions, beware %v !\n", r.RemoteAddr)
  	fmt.Fprintf(w, Cities.getInfo())
}

func main() {
	// TODO: we need to undestand how the private key and certificate is cached.
	// https://godoc.org/golang.org/x/crypto/acme/autocert#Manager.
	m := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("cities.hkjn.me"),
	}
	s := &http.Server{
		Addr: ":https",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
	}

	fmt.Println("Dobroe utro, Larsik!! Where shall we live?")
	fmt.Printf("We have %v cities: %v\n", len(Cities), Cities.getNames())
	fmt.Printf("I will now be a webe server forever, you puny minions, hahahaha!\n")
	http.HandleFunc("/", indexHandler)
	err := s.ListenAndServeTLS("", "")
	panic(err)
}
