package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

type responseStruct struct {
	Connections []struct {
		From struct {
			DepartureTimestamp int    `json:"departureTimestamp"`
			Platform           string `json:"platform"`
		} `json:"from"`
		To struct {
			ArrivalTimestamp int `json:"arrivalTimestamp"`
		} `json:"to"`
	} `json:"connections"`
	From struct {
		Name string `json:"name"`
	} `json:"from"`
	To struct {
		Name string `json:"name"`
	} `json:"to"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <FROM> <TO>\n", os.Args[0])
		os.Exit(1)
	}

	response := makeRequest(os.Args[1], os.Args[2]) // TODO validation

	fmt.Printf("\n%v to %v\n\n", strings.ToUpper(response.From.Name), strings.ToUpper(response.To.Name))

	tableData := buildTableData(response)
	tableHeader := []string{"Nr", "Platform", "Departure", "Arrival"}

	writeTable(tableHeader, tableData)
	println()
}

func buildTableData(response responseStruct) [][]string {
	var tableData [][]string
	for i, con := range response.Connections {
		tableData = append(tableData, []string{
			strconv.Itoa(i),
			con.From.Platform,
			time.Unix(int64(con.From.DepartureTimestamp), 0).Format("15:04"),
			time.Unix(int64(con.To.ArrivalTimestamp), 0).Format("15:04")})
	}

	return tableData
}

func makeRequest(from string, to string) responseStruct {
	getResponse, err := http.Get("http://transport.opendata.ch/v1/connections?from=" + from + "&to=" + to + "&fields[]=from/name&fields[]=to/name&fields[]=connections/from/platform&fields[]=connections/from/departureTimestamp&fields[]=connections/to/arrivalTimestamp")
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(getResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer getResponse.Body.Close()

	var response responseStruct
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}

	return response
}

func writeTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}
