package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/analyticsreporting/v4"
	"google.golang.org/api/option"
)

const (
	scopes          = "https://www.googleapis.com/auth/analytics.readonly"
	keyFileLocation = "pablo-testing-382714-307479c4bad3.json"
	viewID          = "<>" // Replace with the desired Analytics Property ID View
	dateRangeStart  = "30daysAgo"
	dateRangeEnd    = "yesterday"
	sessions        = "ga:sessions"
	users           = "ga:users"
	requestsCr      = "ga:goal8ConversionRate"
	requestValue    = "ga:goal8Value"
	requests        = "ga:goal8Completions"

	dimensionDate    = "ga:date"
	dimensionChannel = "ga:channelGrouping"
)

func initializeAnalyticsReporting() (*analyticsreporting.Service, error) {
	b, err := os.ReadFile(keyFileLocation)
	if err != nil {
		return nil, fmt.Errorf("error reading key file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(b, scopes)
	if err != nil {
		return nil, fmt.Errorf("error creating JWT config: %v", err)
	}

	client := config.Client(context.Background())

	srv, err := analyticsreporting.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("error creating Analytics Reporting service: %v", err)
	}

	return srv, nil
}

func getReport(service *analyticsreporting.Service) (*analyticsreporting.GetReportsResponse, error) {
	request := &analyticsreporting.GetReportsRequest{
		ReportRequests: []*analyticsreporting.ReportRequest{
			{
				ViewId: viewID,
				DateRanges: []*analyticsreporting.DateRange{
					{
						StartDate: dateRangeStart,
						EndDate:   dateRangeEnd,
					},
				},
				Metrics: []*analyticsreporting.Metric{
					{
						Expression: sessions,
					},
					{
						Expression: users,
					},
					{
						Expression: requestsCr,
					},
					{
						Expression: requestValue,
					},
					{
						Expression: requests,
					},
				},
				Dimensions: []*analyticsreporting.Dimension{
					{
						Name: dimensionDate,
					},
					{
						Name: dimensionChannel,
					},
				},
			},
		},
	}

	response, err := service.Reports.BatchGet(request).Do()
	if err != nil {
		return nil, fmt.Errorf("error getting Analytics report: %v", err)
	}

	return response, nil
}

func printResponse(response *analyticsreporting.GetReportsResponse) {
	for _, report := range response.Reports {
		columnHeader := report.ColumnHeader
		dimensionHeaders := columnHeader.Dimensions
		metricHeaders := columnHeader.MetricHeader.MetricHeaderEntries

		for _, row := range report.Data.Rows {
			dimensions := row.Dimensions
			dateRangeValues := row.Metrics

			for i, dimension := range dimensions {
				fmt.Printf("%s: %s\n", dimensionHeaders[i], dimension)
			}

			for i, values := range dateRangeValues {
				fmt.Printf("Date range: %d\n", i)
				for j, value := range values.Values {
					fmt.Printf("%s: %s\n", metricHeaders[j].Name, value)
				}
			}
		}
	}
}

func main() {
	srv, err := initializeAnalyticsReporting()
	if err != nil {
		log.Fatalf("Error initializing Analytics Reporting service: %v", err)
	}

	response, err := getReport(srv)
	if err != nil {
		log.Fatalf("Error getting Analytics report: %v", err)
	}

	printResponse(response)
}
