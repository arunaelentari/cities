package main

import "testing"

func TestAddition(t *testing.T) {
	got := add(1, 1)
	want := 2
	if got != want {
		t.Errorf("add(%v, %v)=%v, want %v\n", 1, 1, got, want)
	}
}

func TestMultiplication(t *testing.T) {
	type testCase struct {
		a, b, want int
	}
	cases := []testCase{
		{a: 2, b: 2, want: 4},
		{a: 3, b: 3, want: 9},
		{a: 4, b: 4, want: 16},
		{a: 5, b: 5, want: 25},
	}
	for i, tc := range cases {
		got := multiply(tc.a, tc.b)
		if got != cases[i].want {
			t.Errorf("multiply(%v, %v)=%v, want %v\n", tc.a, tc.b, got, tc.want)
		}
	}
}

func multiply(a, b int) int {
	return a * b
}

func add(a, b int) int {
	return a + b
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
		t.Errorf("Not in the same order: %v, want %v\n", c, want)
	}
}
