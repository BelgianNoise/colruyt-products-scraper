package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {

	url := "https://www.ticketswap.com/api/graphql/public"
	method := "POST"
	count := 0

	for {
		payload := strings.NewReader(`
		{
			"operationName": "GetListingsForNewEventPage",
			"query": "query GetListingsForNewEventPage($eventTypeId: ID!, $first: Int!, $currency: CurrencyCode) {\n  node(id: $eventTypeId) {\n __typename\n ... on EventType {\n id\n availableListings: listings(\n first: $first\n filter: {listingStatus: AVAILABLE}\n ) {\n __typename\n edges {\n __typename\n node {\n __typename\n id\n hash\n price {\n __typename\n totalPriceWithTransactionFee(toCurrency: $currency) {\n __typename\n amount\n currency\n }\n }\n comment: description\n numberOfTicketsStillForSale\n numberOfTicketsInListing\n }\n }\n }\n }\n }\n}\n",
			"variables": {
				"currency": "EUR",
				"eventTypeId": "RXZlbnRUeXBlOjI1ODczMjQ=",
				"first": 10
			}
		}`)

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)
		if err != nil {
			panic(err)
		}
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[%d] %s\n", count, string(body[:50]))
		res.Body.Close()
		count++
		time.Sleep(500 * time.Millisecond)
	}
}
