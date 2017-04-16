// This program is a website that helps you choose a dream city based on various rankings.
//
// The home page includes a list of cities.
//
// A user can sort cities based on cost, climate, population and other criteria.
//
// These are shown on separate pages.
// - GET /: the index page, gives links to the other pages.
// - GET /by-cost: ranks cities by cost.
// - GET /by-climate: ranks cities by climate.
// - GET /by-population: ranks cities by population.
// - POST /city: allows users to enter a city

package main

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

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
	pageData struct {
		Title    string
		Criteria string
		Cities   cities
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

func (c1 cities) Equal(c2 cities) bool {
	if len(c1) != len(c2) {
		return false
	}
	for i := range c1 {
		if !c1[i].Equal(c2[i]) {
			return false
		}
	}
	return true
}

func (c1 city) Equal(c2 city) bool {
	if c1.name != c2.name {
		return false
	}
	if c1.population != c2.population {
		return false
	}
	if c1.cost != c2.cost {
		return false
	}
	if c1.climate != c2.climate {
		return false
	}
	return true
}

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
//
// TODO: We want a way to sort by a weighted set of criteria, e.g:
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

// indexHandler writes the http reply to the request for the index page.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("You are all my minions, %v, beware  %v, %v!\n", r.RemoteAddr, r.Method, r.URL)
	if r.Method != "GET" {
		log.Printf("This ain't right: %v!\n", r.Method)
		http.Error(w, "This is a bad request. Try again!", http.StatusBadRequest)
		return
	}
	if r.URL.Path != "/" {
		log.Printf("This ain't right: %v!\n", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		htmlo, err := ioutil.ReadFile("404.html")
		if err != nil {
			log.Panicf("Oioioi, there is a problem reading the file: %v\n", err)
		}
		fmt.Fprintf(w, string(htmlo))
		return
	}
	htmlo, err := ioutil.ReadFile("index.html.tmpl") // TODO: we should read this once on startup, not on every request
	if err != nil {
		log.Panicf("Oibai, there is a problem reading the file: %v\n", err)
	}
	t, err := template.New("webpage").Parse(string(htmlo))
	if err != nil {
		log.Panicf("Help, I couldn't parse the %v\n", err)
	}
	Cities.sortBy("cost")
	data := pageData{
		Title:  "Welcome",
		Cities: Cities,
	}

	err = t.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

// TODO: create by-climate and by-population handlers, add links back to index page, change html
func getCriteriaHandler(c string) func(http.ResponseWriter, *http.Request) {
	log.Printf("Yo, I am supposed to get back the handler to criterion %s, maan\n", c)

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("You are all my minions, %v, beware  %v, %v!\n", r.RemoteAddr, r.Method, r.URL)
		if r.Method != "GET" {
			log.Printf("This ain't right: %v!\n", r.Method) // TODO: make it html, add a link
			http.Error(w, "This is a bad request. Try again!", http.StatusBadRequest)
			return
		}
		if r.URL.Path != fmt.Sprintf("/by-%s", c) {
			log.Printf("This ain't right: %v!\n", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			htmlo, err := ioutil.ReadFile("404.html")
			if err != nil {
				log.Panicf("Oioioi, there is a problem reading the file: %v\n", err)
			}
			fmt.Fprintf(w, string(htmlo))
			return
		}
		htmlo, err := ioutil.ReadFile("cities.html.tmpl")
		if err != nil {
			log.Panicf("Oivey, there is a problem reading the file: %v\n", err)
		}
		t, err := template.New("webpage").Parse(string(htmlo))
		if err != nil {
			log.Panicf("Help, I couldn't parse the %v\n", err)
		}
		Cities.sortBy(c)
		data := pageData{
			Title:    fmt.Sprintf("By %s", c),
			Criteria: c,
			Cities:   Cities,
		}

		err = t.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	version := "1.0.0"
	log.Printf("Salem, all is good.  I am the version %v\n", version)
	prod := os.Getenv("CITIES_ISPROD") == "true"
	addr := ":1025"
	if prod {
		addr = ":https"
	}
	s := &http.Server{
		Addr: addr,
	}
	if prod {
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache("cache"),
			HostPolicy: autocert.HostWhitelist("cities.hkjn.me"),
		}
		s.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
	}
	log.Printf("We have %v cities: %v\n", len(Cities), Cities.getNames())
	log.Printf("I will now be a webe server forever at %v, you puny minions, hahahaha!\n", addr)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/by-cost", getCriteriaHandler("cost"))
	http.HandleFunc("/by-population", getCriteriaHandler("population"))
	http.HandleFunc("/by-climate", getCriteriaHandler("climate"))
	if prod {
		panic(s.ListenAndServeTLS("", ""))
	} else {
		panic(s.ListenAndServe())
	}
}
