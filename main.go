package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	twilioUsageRecordValue = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "twilio_usage",
			Help: "Twilio usage record value",
		},
		[]string{"category", "subresource", "unit"},
	)
)

func init() {
	prometheus.MustRegister(twilioUsageRecordValue)
}

type CostExporter struct {
}

func NewCostExporter() *CostExporter {
	return &CostExporter{}
}

type TwilioUsageRecord struct {
	Category    string `json:"category"`
	Count       string `json:"count"`
	CountUnit   string `json:"count_unit"`
	Description string `json:"description"`
	Price       string `json:"price"`
	PriceUnit   string `json:"price_unit"`
	Usage       string `json:"usage"`
	UsageUnit   string `json:"usage_unit"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type TwilioUsageResponse struct {
	UsageRecords []TwilioUsageRecord `json:"usage_records"`
}

func (e *CostExporter) fetchUsageRecords(subresource string) ([]TwilioUsageRecord, error) {
	accountID := os.Getenv("TWILIO_ACCOUNT_ID")
	sid := os.Getenv("TWILIO_SID")
	secret := os.Getenv("TWILIO_SECRET")

	if accountID == "" || sid == "" || secret == "" {
		return nil, fmt.Errorf("TWILIO_ACCOUNT_ID, TWILIO_SID, or TWILIO_SECRET not set")
	}

	url := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Usage/Records/%s.json", accountID, subresource)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", sid, secret)))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch twilio usage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("twilio api returned status %d: %s", resp.StatusCode, string(body))
	}

	var usageResponse TwilioUsageResponse
	if err := json.NewDecoder(resp.Body).Decode(&usageResponse); err != nil {
		return nil, fmt.Errorf("failed to decode twilio usage: %w", err)
	}

	return usageResponse.UsageRecords, nil
}

func (e *CostExporter) updateTwilioMetrics() error {
	subresources := []string{"Today", "Yesterday", "ThisMonth", "LastMonth"}

	for _, subresource := range subresources {
		records, err := e.fetchUsageRecords(subresource)
		if err != nil {
			log.Printf("Error fetching usage records for %s: %v", subresource, err)
			continue
		}

		for _, record := range records {
			// Price
			var price float64
			fmt.Sscanf(record.Price, "%f", &price)
			if price > 0 {
				twilioUsageRecordValue.WithLabelValues(record.Category, subresource, "usd").Set(price)
			}

			// Usage
			var usage float64
			fmt.Sscanf(record.Usage, "%f", &usage)
			if usage > 0 {
				twilioUsageRecordValue.WithLabelValues(record.Category, subresource, record.UsageUnit).Set(usage)
			}

			// Count
			var count float64
			fmt.Sscanf(record.Count, "%f", &count)
			if count > 0 {
				twilioUsageRecordValue.WithLabelValues(record.Category, subresource, record.CountUnit).Set(count)
			}
		}
	}

	return nil
}

func (e *CostExporter) updateMetrics(ctx context.Context) error {
	if err := e.updateTwilioMetrics(); err != nil {
		log.Printf("Error updating Twilio metrics: %v", err)
	}
	return nil
}

func (e *CostExporter) startMetricsUpdater(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour) // Update every hour
	defer ticker.Stop()

	e.updateMetrics(ctx)

	for {
		select {
		case <-ticker.C:
			if err := e.updateMetrics(ctx); err != nil {
				log.Printf("Error updating metrics: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	ctx := context.Background()

	exporter := NewCostExporter()

	// Update metrics immediately on startup
	log.Printf("Updating metrics on startup...")
	if err := exporter.updateMetrics(ctx); err != nil {
		log.Printf("Warning: Failed to update metrics on startup: %v", err)
	}

	go exporter.startMetricsUpdater(ctx)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Twilio Cost Exporter on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
