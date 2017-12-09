package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	monitoring "google.golang.org/api/monitoring/v3"
)

const projectID = "projects/default-1296"

type mon struct {
	s *monitoring.Service
}

// getCustomMetric reads the custom metric created.
func (m mon) getCustomMetric(metricType string) (*monitoring.ListMetricDescriptorsResponse, error) {
	resp, err := m.s.Projects.MetricDescriptors.List(projectID).
		Filter(fmt.Sprintf("metric.type=\"%s\"", metricType)).Do()
	if err != nil {
		return nil, fmt.Errorf("Could not get custom metric: %v", err)
	}

	//log.Printf("getCustomMetric: %s\n", formatResource(resp))
	return resp, nil
}

// writeTimeSeriesValue writes a value for the custom metric created
func (m mon) writeTimeSeriesValue(datapoint int64, metricType string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	timeseries := monitoring.TimeSeries{
		Metric: &monitoring.Metric{
			Type: metricType,
			Labels: map[string]string{
				"environment": "test",
			},
		},
		Resource: &monitoring.MonitoredResource{
			Type: "global",
			Labels: map[string]string{
				"project_id": "default-1296",
			},
		},
		Points: []*monitoring.Point{
			{
				Interval: &monitoring.TimeInterval{
					StartTime: now,
					EndTime:   now,
				},
				Value: &monitoring.TypedValue{
					Int64Value: &datapoint,
				},
			},
		},
	}

	createTimeseriesRequest := monitoring.CreateTimeSeriesRequest{
		TimeSeries: []*monitoring.TimeSeries{&timeseries},
	}

	//log.Printf("writeTimeseriesRequest: %s\n", formatResource(createTimeseriesRequest))
	_, err := m.s.Projects.TimeSeries.Create(projectID, &createTimeseriesRequest).Do()
	if err != nil {
		return fmt.Errorf("Could not write time series value, %v ", err)
	}
	return nil
}

func createService(ctx context.Context) (*monitoring.Service, error) {
	hc, err := google.DefaultClient(ctx, monitoring.MonitoringScope)
	if err != nil {
		return nil, err
	}
	s, err := monitoring.New(hc)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// here we create our custom metrics.
func (m mon) createCustomMetric() error {

	// first metric: ticketsSold
	ld := monitoring.LabelDescriptor{Key: "environment", ValueType: "STRING", Description: "Tickets currently sold"}
	md := monitoring.MetricDescriptor{
		Type:        "custom.googleapis.com/tickets_sold",
		Labels:      []*monitoring.LabelDescriptor{&ld},
		MetricKind:  "GAUGE",
		ValueType:   "INT64",
		Unit:        "tickets",
		Description: "Tickets currently sold",
		DisplayName: "Tickets sold",
	}
	_, err := m.s.Projects.MetricDescriptors.Create(projectID, &md).Do()
	if err != nil {
		return fmt.Errorf("Could not create custom metric: %v", err)
	}
	//log.Printf("createCustomMetric: %s\n", formatResource(resp))

	// second metric: BuyTime
	ld = monitoring.LabelDescriptor{Key: "environment", ValueType: "STRING", Description: "Time to buy a ticket"}
	md = monitoring.MetricDescriptor{
		Type:        "custom.googleapis.com/timetobuy",
		Labels:      []*monitoring.LabelDescriptor{&ld},
		MetricKind:  "GAUGE",
		ValueType:   "INT64",
		Unit:        "ms",
		Description: "Time to buy a ticket",
		DisplayName: "Time to buy",
	}

	//_, err = m.s.Projects.MetricDescriptors.Delete("custom.googleapis.com/timetobuy").Do()

	_, err = m.s.Projects.MetricDescriptors.Create(projectID, &md).Do()
	if err != nil {
		return fmt.Errorf("Could not create custom metric: %v", err)
	}

	// second metric: BuyTimeClient
	ld = monitoring.LabelDescriptor{Key: "environment", ValueType: "STRING", Description: "Time to buy a ticket (client)"}
	md = monitoring.MetricDescriptor{
		Type:        "custom.googleapis.com/timetobuyClient",
		Labels:      []*monitoring.LabelDescriptor{&ld},
		MetricKind:  "GAUGE",
		ValueType:   "INT64",
		Unit:        "ms",
		Description: "Time to buy a ticket (client)",
		DisplayName: "Time to buy on client",
	}

	_, err = m.s.Projects.MetricDescriptors.Create(projectID, &md).Do()
	if err != nil {
		return fmt.Errorf("Could not create custom metric: %v", err)
	}

	//log.Printf("createCustomMetric: %s\n", formatResource(resp))
	return nil
}

func createMon() *mon {

	myMon := new(mon)

	ctx := context.Background()
	s, err := createService(ctx)

	if err != nil {
		log.Fatal(err)
	}

	myMon.s = s
	return myMon
}

// create the metrics and check if the are available.
func (m mon) init() {

	metricType1 := "custom.googleapis.com/tickets_sold"
	metricType2 := "custom.googleapis.com/timetobuy"

	// Create the metrics.
	if err := m.createCustomMetric(); err != nil {
		log.Fatal(err)
	}

	// Wait until the new metric can be read back.
	for {
		resp, err := m.getCustomMetric(metricType1)
		if err != nil {
			log.Fatal(err)
		}
		if len(resp.MetricDescriptors) != 0 {
			break
		}
		time.Sleep(2 * time.Second)
	}
	for {
		resp, err := m.getCustomMetric(metricType2)
		if err != nil {
			log.Fatal(err)
		}
		if len(resp.MetricDescriptors) != 0 {
			break
		}
		time.Sleep(2 * time.Second)
	}

	// Read the TimeSeries for the last 5 minutes for that metric.
	// if err := readTimeSeriesValue(s, projectID, metricType); err != nil {
	// 	log.Fatal(err)
	// }
}

// formatResource marshals a response object as JSON.
func formatResource(resource interface{}) []byte {
	b, err := json.MarshalIndent(resource, "", "    ")
	if err != nil {
		panic(err)
	}
	return b
}

// readTimeSeriesValue reads the TimeSeries for the value specified by metric type in a time window from the last 5 minutes.
// func readTimeSeriesValue() error {
// 	startTime := time.Now().UTC().Add(time.Minute * -5)
// 	endTime := time.Now().UTC()
// 	resp, err := s.Projects.TimeSeries.List(projectResource(projectID)).
// 		Filter(fmt.Sprintf("metric.type=\"%s\"", metricType)).
// 		IntervalStartTime(startTime.Format(time.RFC3339Nano)).
// 		IntervalEndTime(endTime.Format(time.RFC3339Nano)).
// 		Do()
// 	if err != nil {
// 		return fmt.Errorf("Could not read time series value, %v ", err)
// 	}
// 	log.Printf("readTimeseriesValue: %s\n", formatResource(resp))
// 	return nil
// }
