package main

import "fmt"
import "time"
import "math/rand"

const Male = true
const Female = false
const MateCC = 0.000005
const BabyCC = 0.000005

func chance(chance float32) bool {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	result := r.Float32()
	if result <= chance {
		return true
	} else {
		return false
	}
}

func createRandomCitizen() Citizen {
	name := randInt(999999)
	age := randInt(80)
	sex := chance(0.5)
	birthday := randInt(999999)
	return Citizen{name, nil, nil, nil,	age, sex, birthday}
}

func randInt(max int) int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return r.Intn(max)
}

type Citizen struct {
	Name int
	Mother *Citizen
	Father *Citizen
	Partner *Citizen
	Age int
	Sex bool
	Birthday int
}

type City struct {
	Name string
	Citizens []*Citizen
}

func createMate(citizens []*Citizen) {
	var citizen1, citizen2 *Citizen
	for _, citizen := range citizens {
		if citizen1 == nil {
			if citizen.Partner == nil && citizen.Age > 18 {
				citizen1 = citizen
			}
		} else {
			if citizen.Sex != citizen1.Sex && citizen.Partner == nil && citizen.Age > 18 {
				ageDiff := citizen1.Age - citizen.Age
				if (ageDiff >= -8) && (ageDiff <= 8) {
					citizen2 = citizen
					break
				}
			}
		}
	}
	if citizen1 != nil && citizen2 != nil {
		citizen1.Partner = citizen2
		citizen2.Partner = citizen1
		fmt.Printf("%d and %d have mated!\n", citizen1.Name, citizen2.Name)
	} else {
		fmt.Println("Tried to create mating pair but no suitable mates available :-(")
	}
}

func makeBaby(citizens []*Citizen, birthday int) Citizen {
	// Collect all partnered citizens
	babyMakers := []*Citizen{}
	for _, citizen := range citizens {
		if citizen.Partner != nil {
			babyMakers = append(babyMakers, citizen)
		}
	}
	// Pick a baby maker at random
	var mother, father *Citizen
	l := len(babyMakers)
	if l == 0 { return Citizen{} }
	i := randInt(l)
	if babyMakers[i].Sex == Male {
		father = babyMakers[i]
		mother = babyMakers[i].Partner
	} else {
		mother = babyMakers[i]
		father = babyMakers[i].Partner
	}

	baby := Citizen{randInt(999999), mother, father, nil,	0, chance(0.5), birthday}
	sex := "boy"
	if baby.Sex == Female { sex = "girl" }
	fmt.Printf("%d and %d had a baby %s and named it %d. It was born on %d.\n", mother.Name, father.Name, sex, baby.Name, baby.Birthday)
	return baby
}

func age(citizens []*Citizen, day int) {
	for _, citizen := range citizens {
		if citizen.Birthday == day {
			citizen.Age ++
			fmt.Printf("%d aged and is now %d years old. (Birthday %d).\n", citizen.Name, citizen.Age, citizen.Birthday)
		}
	}
}

func main() {
	cities := []*City{}
	citizens := []*Citizen{}
	// Initialise a city with some citizens
	cities = append(cities, &City{ "city1", nil } )
	// Create 100 citizens in city1
	for i := 0; i < 100; i++ {
		c := createRandomCitizen()
		citizens = append(citizens, &c)
	}
	fmt.Println(cities[0].Name)
	l := len(citizens)
	for i := 0; i < l; i++ {
		fmt.Println(citizens[i].Age)
	}
	d := 0
	for {
		if chance(MateCC) {
			createMate(citizens)
		}
		if chance(BabyCC) {
			baby := makeBaby(citizens, d)
			p := &baby
			citizens = append(citizens, p)
		}
		// fmt.Println(d)
		age(citizens, d)
		d ++
		if d == 999999 {
			fmt.Println("New year!")
			d = 0
		}
	}
}

