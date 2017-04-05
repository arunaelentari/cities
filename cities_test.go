package main

import (
	"testing"
)

func TestCities_Equal(t *testing.T) {
	type testCase struct {
		c1, c2 cities
		want   bool
	}
	cases := []testCase{
		{
			c1:   cities{},
			c2:   cities{},
			want: true,
		},
		{
			c1: cities{},
			c2: cities{
				city{name: "Barcelona", population: 1.6e6, cost: ReasonableCost, climate: GreatClimate},
			},
			want: false,
		},
	}
	for _, tc := range cases {
		//	c1 := cities{}
		//	c2 := cities{}
		//		want := true
		if tc.c1.Equal(tc.c2) != tc.want {
			t.Errorf("cities.Equal() should be %v for cities\n%q\nand\n%q\n", tc.want, tc.c1, tc.c2)
		}
	}
}

func TestCities_sortBy(t *testing.T) {
	c := cities{
		city{name: "Barcelona", population: 1.6e6, cost: ReasonableCost, climate: GreatClimate},
		city{name: "Seattle", population: 652405, cost: ExpensiveCost, climate: GoodClimate},
		city{name: "New York", population: 8.406e6, cost: ExpensiveCost, climate: GoodClimate},
		city{name: "Copenhagen", population: 562379, cost: ExpensiveCost, climate: PoorClimate},
		city{name: "Stockholm", population: 789024, cost: ExpensiveCost, climate: PoorClimate},
		city{name: "Deviltown", population: 1233567890, cost: VeryExpensiveCost, climate: NastyClimate},
		city{name: "Paradisio", population: 1e6, cost: CheapCost, climate: PerfectClimate},
	}
	c.sortBy("name")
	want := cities{
		city{name: "Barcelona", population: 1.6e6, cost: ReasonableCost, climate: GreatClimate},
		city{name: "Copenhagen", population: 562379, cost: ExpensiveCost, climate: PoorClimate},
		city{name: "Deviltown", population: 1233567890, cost: VeryExpensiveCost, climate: NastyClimate},
		city{name: "New York", population: 8.406e6, cost: ExpensiveCost, climate: GoodClimate},
		city{name: "Paradisio", population: 1e6, cost: CheapCost, climate: PerfectClimate},
		city{name: "Seattle", population: 652405, cost: ExpensiveCost, climate: GoodClimate},
		city{name: "Stockholm", population: 789024, cost: ExpensiveCost, climate: PoorClimate},
	}
	if !c.Equal(want) {
		t.Errorf("Not in the same order, cities sortBy(\"name\"):\n%v\nWant\n%v\n", c, want)
	}
}
