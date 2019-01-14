package main

import "fmt"
import "time"
import "math"
import "math/rand"

const Male = true
const Female = false
const BabyPercentage = 30.0

const Speed = 100 // This is actually more like 'days per year'

var day int
var year int
var decimalPlaces float64

func init() {
	day = 0
	year = 100
	decimalPlaces = 1.0 / float64(Speed)
}

func chance(chance float32) bool {
	result := rand.Float32()
	if result <= chance {
		return true
	} else {
		return false
	}
}

func randInt(max int) int {
	return rand.Intn(max)
}

func randFloat(max int) float32 {
	max32 := float32(max)
	return rand.Float32() * max32
}

func getDay() float32 {
	day32 := float32(day) / float32(Speed)
	yearAndDay := float32(year) + day32
	return round(yearAndDay, decimalPlaces)
}

func getAge(day, birthday float32) int {
	diff := day - birthday
	return int(diff)
}

func round(x float32, unit float64) float32 {
	return float32(math.Round(float64(x)/unit) * unit)
}

func createRandomCitizen() Citizen {
	name := randInt(999999)
	sex := chance(0.5)
	birthday := randFloat(80)
	return Citizen{name, nil, nil, nil, 0.0, sex, round(birthday, decimalPlaces)}
}

func die(age int) bool {
	var c float32
	switch {
	case age <= 3:
		c = 0.00001
	case age >3 && age <= 10:
		c = 0.000002
	case age >10 && age <=50:
		c = 0.000005
	case age >50 && age <=60:
		c = 0.00003
	case age >60 && age <=70:
		c = 0.0002
	case age >70 && age <=80:
		c = 0.0004
	case age >80 && age <=90:
		c = 0.0009
	default:
		c = 0.009
	}
	return chance(c)
}

type Citizen struct {
	Name int
	Mother *Citizen
	Father *Citizen
	Partner *Citizen
	Fertile float32 // The day we can have a baby again
	Sex bool
	Birthday float32
}

type City struct {
	Name string
	Citizens []*Citizen
}

func createMates(males, females map[*Citizen]int) int {
	// lenMales := len(males)
	// lenFemales := len(females)
	mates := 0
	for male, maleAge := range males {
		for female, femaleAge := range females {
			ageDiff := maleAge - femaleAge
			if ageDiff >= -20 && ageDiff <= 20 {
				// Add in a semblance of chance
				// if chance(0.7) {
					male.Partner = female
					female.Partner = male
					mates++
					delete(females, female)
					delete(males, male)
					break
				// }
			}
		}
	}
	// fmt.Printf("Of the %d males and %d females I had, I made %d mates.\n", lenMales, lenFemales, mates)
	return mates
}

func makeBabies(babyMakers map[*Citizen]int) []*Citizen {
	var mother, father *Citizen
	var babies []*Citizen
	for babyMaker, _ := range babyMakers {
		if chance(0.002) {
			mother = babyMaker
			father = babyMaker.Partner
			today := getDay()
			baby := Citizen{randInt(999999), mother, father, nil, 0.0, chance(0.5), today}
			babies = append(babies, &baby)
			// Mother can't have another baby for a year
			mother.Fertile = today + 1.0
		}
	}
	// sex := "boy"
	// if baby.Sex == Female { sex = "girl" }
	// fmt.Printf("%d and %d had a baby %s and named it %d. It was born on %f.\n", mother.Name, father.Name, sex, baby.Name, baby.Birthday)

	return babies
}

func main() {
	rand.Seed(time.Now().UnixNano())
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

	babiesThisYear := 0
	matesThisYear := 0
	deathsThisYear := 0
	deathsBaby := 0
	deathsTeen := 0
	deathsAdult := 0
	deaths60Plus := 0

	accumulativeBabies := 0
	accumulativeMates := 0
	accumulativeDeaths := 0
	deathsBabyAccumulative := 0
	deathsTeenAccumulative := 0
	deathsAdultAccumulative := 0
	deaths60PlusAccumulative := 0

	var dead []int

	for {
		babyMakers := make(map[*Citizen]int)
		singleMales := make(map[*Citizen]int)
		singleFemales := make(map[*Citizen]int)
		dead = nil
		today := getDay()
		for index, citizen := range citizens {
			age := getAge(getDay(), citizen.Birthday)
			// Add the index of this citizen to another array if they die
			if die(age) {
				dead = append(dead, index)
				// fmt.Printf("%d died aged %d\n", citizen.Name, age)
				if citizen.Partner != nil {
					citizen.Partner.Partner = nil // Make this citizen's partner single again
				}
				// If this dead citizen's partner has already been added to the baby makers
				// that citizen can no longer make babies.
				delete(babyMakers, citizen.Partner)
				switch {
				case age <= 12:
					deathsBaby++
				case age > 12 && age <= 18:
					deathsTeen++
				case age > 18 && age <= 60:
					deathsAdult++
				case age > 60:
					deaths60Plus++
				}
				continue
			}
			// If you don't have a partner, you might get one.
			if citizen.Partner == nil {
				if citizen.Sex == Male {
					if age > 18 && age < 60 {
						singleMales[citizen] = age
					}
				} else {
					if age > 18 && age < 60 {
						singleFemales[citizen] = age
					}
				}
			} else {
				if citizen.Fertile < today && citizen.Sex == Female && age > 18 && age < 50 {
					// We don't actually use the map values, but using a map makes it easier to remove elements.
					// It also randomises the order.
					babyMakers[citizen] = 0
				}
			}
			// fmt.Printf("Boys: %d     Girls: %d\n", boys, girls)
		}

		// Kill the dead.
		deathsThisYear += len(dead)
		for i, index := range dead {
    	index = index - i // The size of the citizens array will decrease by one each time
    	copy(citizens[index:], citizens[index+1:])
    	citizens[len(citizens)-1] = nil
    	citizens = citizens[:len(citizens)-1]
  	}

		// Pass all the singles to the mate maker
		matesThisYear += createMates(singleMales, singleFemales)

		// Pass the baby makers to the baby maker
		babies := makeBabies(babyMakers)
		for _, baby := range babies {
			citizens = append(citizens, baby)
		}
		babiesThisYear += len(babies)

		time.Sleep(10 * time.Nanosecond)

		if day == Speed { // If this is a new year
			day = 0
			year++
			// Reset counters
			babiesThisYear = 0
			matesThisYear = 0
			deathsThisYear = 0
			deathsBaby = 0
			deathsTeen = 0
			deathsAdult = 0
			deaths60Plus = 0
			// fmt.Printf("New year: %d\n", year)
			// fmt.Printf("Population: %d\n", len(citizens))
			// fmt.Printf("Babies this year: %d\n", babiesThisYear)
			// fmt.Printf("Mates this year: %d\n", matesThisYear)
			// fmt.Printf("Failed mates this year: %d\n", failedMatesThisYear)
			// fmt.Printf("Deaths this year: %d\n", deathsThisYear)
			// Give us a report every 100 years
			if math.Mod(float64(year), 100.0) == 0 {
				fmt.Printf("Year: %d\n", year)
				fmt.Printf("Population: %d\n", len(citizens))
				fmt.Printf("Potential baby makers right now: %d\n", len(babyMakers) * 2)
				fmt.Printf("Average babies: %d\n", accumulativeBabies / Speed)
				fmt.Printf("Average mates: %d\n", accumulativeMates / Speed)
				fmt.Printf("Average deaths: %d\n", accumulativeDeaths / Speed)
				fmt.Printf("Baby deaths: %d\n", deathsBabyAccumulative / Speed)
				fmt.Printf("Teen deaths: %d\n", deathsTeenAccumulative / Speed)
				fmt.Printf("Adult deaths: %d\n", deathsAdultAccumulative / Speed)
				fmt.Printf("60+ deaths: %d\n", deaths60PlusAccumulative / Speed)
				// Reset counters
				accumulativeBabies = 0
				accumulativeMates = 0
				accumulativeDeaths = 0
				deathsBabyAccumulative = 0
				deathsTeenAccumulative = 0
				deathsAdultAccumulative = 0
				deaths60PlusAccumulative = 0
			}
		} else {
			day++
		}
		// Accumulate the accumulators
		accumulativeBabies += babiesThisYear
		accumulativeMates += matesThisYear
		accumulativeDeaths += deathsThisYear
		deathsBabyAccumulative += deathsBaby
		deathsTeenAccumulative += deathsTeen
		deathsAdultAccumulative += deathsAdult
		deaths60PlusAccumulative += deaths60Plus
	}
}
