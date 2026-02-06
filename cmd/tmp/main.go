package main

import (
	"fmt"

	jmespath "github.com/jmespath/go-jmespath"
)

func main() {

	parseAndOutput("[].{name:name, foo:foo}")

	parseAndOutput("locations[?state == 'WA'].name | sort(@) | {WashingtonCities: join(', ', @)}")
	parseAndOutput("locations[?state == 'WA']")
}

func parseAndOutput(query string) {
	parser := jmespath.NewParser()
	ast, err := parser.Parse(query)
	if err != nil {
		fmt.Printf("Error parsing: %s\n", err)
	}

	fmt.Printf("%s\n\n", query)
	fmt.Printf("%s\n\n", ast.PrettyPrint(2))
	fields, err := ast.GetResultFields()
	if err != nil {
		fmt.Printf("Error getting fields: %s\n", err)
	}
	fmt.Printf("Fields: %v\n\n", fields)

}
