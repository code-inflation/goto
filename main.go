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

	response, err := http.Get("http://transport.opendata.ch/v1/connections?from=" + os.Args[1] + "&to=" + os.Args[2] + "&fields[]=from/name&fields[]=to/name&fields[]=connections/from/platform&fields[]=connections/from/departureTimestamp&fields[]=connections/to/arrivalTimestamp")
	if err != nil {
		log.Fatal(err)
	} else {
		body, err := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		var response responseStruct
		err = json.Unmarshal(body, &response)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\n%v to %v\n\n", strings.ToUpper(response.From.Name), strings.ToUpper(response.To.Name))

		var data [][]string

		for i, con := range response.Connections {
			data = append(data, []string{
				strconv.Itoa(i),
				con.From.Platform,
				time.Unix(int64(con.From.DepartureTimestamp), 0).Format("15:04:05"),
				time.Unix(int64(con.To.ArrivalTimestamp), 0).Format("15:04:05")})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Nr", "Platform", "Departure Time", "Arrival Time"})

		for _, v := range data {
			table.Append(v)
		}
		table.Render() // Send output
	}
}
