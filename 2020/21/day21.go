package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

var inputFile = flag.String("f", "testInput", "Relative file path to use as input.")

type Food struct {
	ID          int
	Ingredients map[string]*Ingredient
	Allergens   map[string]bool
}
type Allergen struct {
	Name       string
	Ingredient *Ingredient
	Foods      []*Food
}
type Ingredient struct {
	name              string
	possibleAllergens map[string]bool
	allergen          *Allergen
	foods             []*Food
}

func (i Ingredient) String() string {
	var sb strings.Builder
	sb.WriteString("{\n")
	sb.WriteString(fmt.Sprintf("\tName %s\n", i.name))
	sb.WriteString("\tFoods:")
	for _, f := range i.foods {
		sb.WriteString(fmt.Sprintf("%d,", f.ID))
	}
	sb.WriteString("\n")
	sb.WriteString("\tAllergens:")
	for a := range i.possibleAllergens {
		sb.WriteString(fmt.Sprintf("%s,", a))
	}
	sb.WriteString("\n")

	sb.WriteString("}")
	return sb.String()
}

func uniqueStrings(input []string) (result []string) {
	tmpMap := make(map[string]bool)
	for _, i := range input {
		if _, ok := tmpMap[i]; !ok {
			result = append(result, i)
			tmpMap[i] = true
		}
	}
	return
}
func ingredientOcurrences(input []string) map[string]int {
	result := make(map[string]int)

	for _, i := range input {
		if _, ok := result[i]; !ok {
			result[i] = 0
		}
		result[i] = result[i] + 1
	}
	return result
}
func main() {
	flag.Parse()
	input, _ := ioutil.ReadFile(*inputFile)
	trimmedInput := strings.Split(strings.TrimSpace(string(input)), "\n")
	ingredients, allergens := parse(trimmedInput)

	//For each allergen:
	for _, allergen := range allergens {
		//- Get list of foods that have this allergen
		globalIngredients := make(map[string]int)
		for _, f := range allergen.Foods {
			for _, i := range f.Ingredients {
				globalIngredients[i.name] = globalIngredients[i.name] + 1
			}
		}
		//- Generate global list of ingredients
		//- keep all the ingredients that are available on every food (ingredient map[ingredient]count==number of foods)
		for ingredient, count := range globalIngredients {
			//- If count ==  => Keep candidate
			//- This are the possible ingredients for this allergen
			if count == len(allergen.Foods) {
				ingredients[ingredient].possibleAllergens[allergen.Name] = true
			}
		}
	}

	ingredientsInFoodWithoutAllergens := 0
	matchedAllergens := make(map[string]bool)

	for _, ing := range ingredients {
		var allergensToClear []string
		for possibleAllergen := range ing.possibleAllergens {
			_, ok := matchedAllergens[possibleAllergen]
			if ok {
				allergensToClear = append(allergensToClear, possibleAllergen)
			}
		}
		for _, allergen := range allergensToClear {
			delete(ing.possibleAllergens, allergen)
		}
		if len(ing.possibleAllergens) == 0 {
			ingredientsInFoodWithoutAllergens += len(ing.foods)
			ing.allergen = nil
		} else if len(ing.possibleAllergens) == 1 {
			for allergen := range ing.possibleAllergens {
				allergens[allergen].Ingredient = ing
				ing.allergen = allergens[allergen]
				matchedAllergens[allergen] = true
				ing.possibleAllergens = nil
			}
		}
	}

	for len(matchedAllergens) != len(allergens) {
		for _, ing := range ingredients {
			if len(ing.possibleAllergens) == 0 || ing.allergen != nil {
				continue
			}
			var allergensToClear []string
			for possibleAllergen := range ing.possibleAllergens {
				_, ok := matchedAllergens[possibleAllergen]
				if ok {
					allergensToClear = append(allergensToClear, possibleAllergen)
				}
			}
			for _, allergen := range allergensToClear {
				delete(ing.possibleAllergens, allergen)
			}

			if len(ing.possibleAllergens) == 1 {
				for allergen := range ing.possibleAllergens {
					allergens[allergen].Ingredient = ing
					ing.allergen = allergens[allergen]
					matchedAllergens[allergen] = true
					ing.possibleAllergens = nil
				}
			}
		}
		fmt.Printf("Missing %d\n", len(allergens)-len(matchedAllergens))
	}
	allergenNames := make([]string, 0, len(allergens))
	for allergenName := range allergens {
		allergenNames = append(allergenNames, allergenName)
	}
	sort.Strings(allergenNames)
	ingredientsByAllergen := make([]string, 0, len(allergens))
	for _, allergenName := range allergenNames {
		allergen := allergens[allergenName]
		ingredientsByAllergen = append(ingredientsByAllergen, allergen.Ingredient.name)
	}
	fmt.Printf("Part1: %d\n", ingredientsInFoodWithoutAllergens)
	// fmt.Printf("Sorted allergens: %s", strings.Join(allergenNames, ","))
	fmt.Printf("Part2: %s", strings.Join(ingredientsByAllergen, ","))
}

func parse(s []string) (map[string]*Ingredient, map[string]*Allergen) {

	ingredients := make(map[string]*Ingredient)
	allergens := make(map[string]*Allergen)
	for k, l := range s {

		currentFood := &Food{ID: k, Allergens: make(map[string]bool), Ingredients: make(map[string]*Ingredient)}
		contents := strings.Split(strings.TrimRight(l, ")"), "(contains")

		foodAlergens := make([]string, 0, 10)

		for _, allergenName := range strings.Split(contents[1], ", ") {
			allergenName = strings.TrimSpace(allergenName)
			var allergen *Allergen
			var ok bool
			if allergen, ok = allergens[allergenName]; !ok {
				allergen = &Allergen{
					Name: allergenName,
				}
			}
			allergen.Foods = append(allergen.Foods, currentFood)
			allergens[allergenName] = allergen
			currentFood.Allergens[allergenName] = true
			foodAlergens = append(foodAlergens, strings.TrimSpace(allergenName))
		}
		for _, ing := range strings.Split(strings.TrimSpace(contents[0]), " ") {

			var ingredient *Ingredient
			var ok bool
			if ingredient, ok = ingredients[strings.TrimSpace(ing)]; !ok {
				ingredient = &Ingredient{
					name:              ing,
					possibleAllergens: make(map[string]bool),
				}
			}
			ingredient.foods = append(ingredient.foods, currentFood)
			currentFood.Ingredients[ingredient.name] = ingredient
			ingredients[ingredient.name] = ingredient
		}
	}
	return ingredients, allergens

}
