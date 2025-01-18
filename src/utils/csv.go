package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type ApolloData struct {
	Email         string
	FirstName     string
	LastName      string
	Title         string
	State         string
	City          string
	Country       string
	CompanyName   string
	Industry      string
	LinkedInLink  string
	EmployeeCount int
}

func generatePossibleHeadersForOne(header string) []string {
	return []string{
		header,
		strings.ToLower(header),
		strings.ToUpper(header),
		fmt.Sprintf("%s%s", strings.ToUpper(header[0:1]), strings.ToLower(header[1:])),
	}
}

func generatePossibleHeadersForTwo(h1, h2 string) []string {
	h1Lower := strings.ToLower(h1)
	h2Lower := strings.ToLower(h2)
	h1Upper := strings.ToUpper(h1)
	h2Upper := strings.ToUpper(h2)
	h1Title := fmt.Sprintf("%s%s", strings.ToUpper(h1[0:1]), strings.ToLower(h1[1:]))
	h2Title := fmt.Sprintf("%s%s", strings.ToUpper(h2[0:1]), strings.ToLower(h2[1:]))

	return []string{
		h1 + h2,
		h1 + " " + h2,
		h1 + "_" + h2,
		h1Lower + h2Lower,
		h1Lower + " " + h2Lower,
		h1Lower + "_" + h2Lower,
		h1Title + h2Title,
		h1Title + " " + h2Title,
		h1Title + "_" + h2Title,
		h1Upper + h2Upper,
		h1Upper + " " + h2Upper,
		h1Upper + "_" + h2Upper,
	}
}

func ApolloDataAsBestAsPossible(headers []string, row []string) ApolloData {
	var data ApolloData

	headerIndex := make(map[string]int)
	for i, header := range headers {
		possibilities := generatePossibleHeadersForOne(header)
		for _, possibility := range possibilities {
			headerIndex[possibility] = i
		}
	}

	state := generatePossibleHeadersForOne("state")
	state = append(state, generatePossibleHeadersForTwo("company", "state")...)
	country := generatePossibleHeadersForOne("country")
	country = append(country, generatePossibleHeadersForTwo("company", "country")...)
	city := generatePossibleHeadersForOne("city")
	city = append(city, generatePossibleHeadersForTwo("company", "city")...)

	possibleFields := map[string][]string{
		"Email":         generatePossibleHeadersForOne("email"),
		"FirstName":     generatePossibleHeadersForTwo("first", "name"),
		"LastName":      generatePossibleHeadersForTwo("last", "name"),
		"Title":         generatePossibleHeadersForOne("title"),
		"State":         state,
		"City":          city,
		"Country":       country,
		"CompanyName":   generatePossibleHeadersForTwo("company", "name"),
		"Industry":      generatePossibleHeadersForOne("industry"),
		"LinkedInLink":  generatePossibleHeadersForTwo("LinkedIn", "Link"),
		"EmployeeCount": generatePossibleHeadersForTwo("Employee", "Count"),
	}

	for field, possibilities := range possibleFields {
		for _, possibility := range possibilities {
			if index, ok := headerIndex[possibility]; ok && index < len(row) {
				switch field {
				case "Email":
					data.Email = row[index]
				case "FirstName":
					data.FirstName = row[index]
				case "LastName":
					data.LastName = row[index]
				case "Title":
					data.Title = row[index]
				case "State":
					data.State = row[index]
				case "City":
					data.City = row[index]
				case "Country":
					data.Country = row[index]
				case "CompanyName":
					data.CompanyName = row[index]
				case "Industry":
					data.Industry = row[index]
				case "LinkedInLink":
					data.LinkedInLink = row[index]
				case "EmployeeCount":
					{
						v, err := strconv.Atoi(row[index])
						if err == nil {
							data.EmployeeCount = v
						}
					}
				}
				break
			}
		}
	}

	return data
}
