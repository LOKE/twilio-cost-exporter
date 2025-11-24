package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Current Invoice Metrics
	herokuInvoiceCurrentTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_invoice_current_total_usd",
			Help: "Current Heroku invoice total in USD",
		},
		[]string{"payment_status", "state", "period_end"},
	)
	herokuInvoiceCurrentCharges = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_invoice_current_charges_total_usd",
			Help: "Current Heroku invoice charges total in USD",
		},
		[]string{"period_end"},
	)
	herokuInvoiceCurrentCredits = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_invoice_current_credits_total_usd",
			Help: "Current Heroku invoice credits total in USD",
		},
		[]string{"period_end"},
	)
	herokuInvoiceCurrentAddons = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_invoice_current_addons_total_usd",
			Help: "Current Heroku invoice addons total in USD",
		},
		[]string{"period_end"},
	)
	herokuInvoiceCurrentDatabase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_invoice_current_database_total_usd",
			Help: "Current Heroku invoice database total in USD",
		},
		[]string{"period_end"},
	)
	herokuInvoiceCurrentPlatform = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_invoice_current_platform_total_usd",
			Help: "Current Heroku invoice platform total in USD",
		},
		[]string{"period_end"},
	)
	herokuInvoiceCurrentDynoUnits = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_invoice_current_dyno_units",
			Help: "Current Heroku invoice dyno units",
		},
		[]string{"period_end"},
	)
	herokuInvoiceCurrentWeightedDynoHours = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_invoice_current_weighted_dyno_hours",
			Help: "Current Heroku invoice weighted dyno hours",
		},
		[]string{"period_end"},
	)

	// Previous Invoice Metrics
	herokuInvoicePreviousTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_invoice_previous_total_usd",
			Help: "Previous Heroku invoice total in USD",
		},
		[]string{"payment_status", "state", "period_end"},
	)
	herokuInvoicePreviousCharges = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_invoice_previous_charges_total_usd",
			Help: "Previous Heroku invoice charges total in USD",
		},
		[]string{"period_end"},
	)
	herokuInvoicePreviousCredits = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_invoice_previous_credits_total_usd",
			Help: "Previous Heroku invoice credits total in USD",
		},
		[]string{"period_end"},
	)
	herokuInvoicePreviousAddons = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_invoice_previous_addons_total_usd",
			Help: "Previous Heroku invoice addons total in USD",
		},
		[]string{"period_end"},
	)
	herokuInvoicePreviousDatabase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_invoice_previous_database_total_usd",
			Help: "Previous Heroku invoice database total in USD",
		},
		[]string{"period_end"},
	)
	herokuInvoicePreviousPlatform = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_invoice_previous_platform_total_usd",
			Help: "Previous Heroku invoice platform total in USD",
		},
		[]string{"period_end"},
	)
	herokuInvoicePreviousDynoUnits = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_invoice_previous_dyno_units",
			Help: "Previous Heroku invoice dyno units",
		},
		[]string{"period_end"},
	)
	herokuInvoicePreviousWeightedDynoHours = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_invoice_previous_weighted_dyno_hours",
			Help: "Previous Heroku invoice weighted dyno hours",
		},
		[]string{"period_end"},
	)

	herokuTeamUsageCurrentUSD = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_team_usage_current_usd",
			Help: "Current month Heroku team usage cost in USD",
		},
		[]string{"team_name", "app_name", "type", "month"},
	)

	herokuTeamUsagePreviousUSD = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_team_usage_previous_usd",
			Help: "Previous month Heroku team usage cost in USD",
		},
		[]string{"team_name", "app_name", "type", "month"},
	)

	herokuTeamUsageCurrentDynoUnits = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_team_usage_current_dyno_units",
			Help: "Current month Heroku team usage dyno units",
		},
		[]string{"team_name", "app_name", "month"},
	)

	herokuTeamUsagePreviousDynoUnits = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "heroku_team_usage_previous_dyno_units",
			Help: "Previous month Heroku team usage dyno units",
		},
		[]string{"team_name", "app_name", "month"},
	)
)

func init() {
	prometheus.MustRegister(herokuInvoiceCurrentTotal)
	prometheus.MustRegister(herokuInvoiceCurrentCharges)
	prometheus.MustRegister(herokuInvoiceCurrentCredits)
	prometheus.MustRegister(herokuInvoiceCurrentAddons)
	prometheus.MustRegister(herokuInvoiceCurrentDatabase)
	prometheus.MustRegister(herokuInvoiceCurrentPlatform)
	prometheus.MustRegister(herokuInvoiceCurrentDynoUnits)
	prometheus.MustRegister(herokuInvoiceCurrentWeightedDynoHours)

	prometheus.MustRegister(herokuInvoicePreviousTotal)
	prometheus.MustRegister(herokuInvoicePreviousCharges)
	prometheus.MustRegister(herokuInvoicePreviousCredits)
	prometheus.MustRegister(herokuInvoicePreviousAddons)
	prometheus.MustRegister(herokuInvoicePreviousDatabase)
	prometheus.MustRegister(herokuInvoicePreviousPlatform)
	prometheus.MustRegister(herokuInvoicePreviousDynoUnits)
	prometheus.MustRegister(herokuInvoicePreviousWeightedDynoHours)

	prometheus.MustRegister(herokuTeamUsageCurrentUSD)
	prometheus.MustRegister(herokuTeamUsagePreviousUSD)
	prometheus.MustRegister(herokuTeamUsageCurrentDynoUnits)
	prometheus.MustRegister(herokuTeamUsagePreviousDynoUnits)
}

type CostExporter struct {
}

func NewCostExporter() *CostExporter {
	return &CostExporter{}
}

type HerokuInvoice struct {
	ChargesTotal      int64   `json:"charges_total"`
	CreatedAt         string  `json:"created_at"`
	CreditsTotal      int64   `json:"credits_total"`
	ID                string  `json:"id"`
	Number            int64   `json:"number"`
	PeriodEnd         string  `json:"period_end"`
	PeriodStart       string  `json:"period_start"`
	State             int     `json:"state"`
	Total             int64   `json:"total"`
	UpdatedAt         string  `json:"updated_at"`
	AddonsTotal       int64   `json:"addons_total"`
	DatabaseTotal     int64   `json:"database_total"`
	DynoUnits         float64 `json:"dyno_units"`
	PlatformTotal     int64   `json:"platform_total"`
	PaymentStatus     string  `json:"payment_status"`
	WeightedDynoHours float64 `json:"weighted_dyno_hours"`
}

type HerokuTeamUsage struct {
	ID                  string  `json:"id"`
	Name                string  `json:"name"`
	Month               string  `json:"month"`
	Dynos               float64 `json:"dynos"`
	Data                float64 `json:"data"`
	Addons              float64 `json:"addons"`
	Connect             float64 `json:"connect"`
	Partner             float64 `json:"partner"`
	PrivateSpace        float64 `json:"private_space"`
	PrivateSpaceCredits float64 `json:"private_space_credits"`
	ShieldSpace         float64 `json:"shield_space"`
	ShieldSpaceCredits  float64 `json:"shield_space_credits"`
	Space               float64 `json:"space"`
	Apps                []struct {
		AppName string  `json:"app_name"`
		Dynos   float64 `json:"dynos"`
		Addons  float64 `json:"addons"`
		Data    float64 `json:"data"`
		Connect float64 `json:"connect"`
		Partner float64 `json:"partner"`
	} `json:"apps"`
}

func (e *CostExporter) updateHerokuMetrics() error {
	apiKey := os.Getenv("HEROKU_API_KEY")
	if apiKey == "" {
		log.Println("HEROKU_API_KEY not set, skipping Heroku invoice metrics")
		return nil
	}

	url := "https://api.heroku.com/account/invoices"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch heroku invoices: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("heroku api returned status %d: %s", resp.StatusCode, string(body))
	}

	var invoices []HerokuInvoice
	if err := json.NewDecoder(resp.Body).Decode(&invoices); err != nil {
		return fmt.Errorf("failed to decode heroku invoices: %w", err)
	}

	if len(invoices) == 0 {
		log.Println("No Heroku invoices found")
		return nil
	}

	// Sort by PeriodEnd
	sort.Slice(invoices, func(i, j int) bool {
		return invoices[i].PeriodEnd < invoices[j].PeriodEnd
	})

	if len(invoices) == 0 {
		log.Println("No Heroku invoices found")
		return nil
	}

	// Reset gauges
	herokuInvoiceCurrentTotal.Reset()
	herokuInvoiceCurrentCharges.Reset()
	herokuInvoiceCurrentCredits.Reset()
	herokuInvoiceCurrentAddons.Reset()
	herokuInvoiceCurrentDatabase.Reset()
	herokuInvoiceCurrentPlatform.Reset()
	herokuInvoiceCurrentDynoUnits.Reset()
	herokuInvoiceCurrentWeightedDynoHours.Reset()

	herokuInvoicePreviousTotal.Reset()
	herokuInvoicePreviousCharges.Reset()
	herokuInvoicePreviousCredits.Reset()
	herokuInvoicePreviousAddons.Reset()
	herokuInvoicePreviousDatabase.Reset()
	herokuInvoicePreviousPlatform.Reset()
	herokuInvoicePreviousDynoUnits.Reset()
	herokuInvoicePreviousWeightedDynoHours.Reset()

	// Current Invoice (Last one)
	current := invoices[len(invoices)-1]
	log.Printf("Current Heroku invoice: %s (Period End: %s, Total: %d, State: %d)", current.ID, current.PeriodEnd, current.Total, current.State)

	herokuInvoiceCurrentTotal.WithLabelValues(current.PaymentStatus, fmt.Sprintf("%d", current.State), current.PeriodEnd).Set(float64(current.Total) / 100.0)
	herokuInvoiceCurrentCharges.WithLabelValues(current.PeriodEnd).Set(float64(current.ChargesTotal) / 100.0)
	herokuInvoiceCurrentCredits.WithLabelValues(current.PeriodEnd).Set(float64(current.CreditsTotal) / 100.0)
	herokuInvoiceCurrentAddons.WithLabelValues(current.PeriodEnd).Set(float64(current.AddonsTotal) / 100.0)
	herokuInvoiceCurrentDatabase.WithLabelValues(current.PeriodEnd).Set(float64(current.DatabaseTotal) / 100.0)
	herokuInvoiceCurrentPlatform.WithLabelValues(current.PeriodEnd).Set(float64(current.PlatformTotal) / 100.0)
	herokuInvoiceCurrentDynoUnits.WithLabelValues(current.PeriodEnd).Set(float64(current.DynoUnits))
	herokuInvoiceCurrentWeightedDynoHours.WithLabelValues(current.PeriodEnd).Set(float64(current.WeightedDynoHours))

	// Previous Invoice (Second to last, if exists)
	if len(invoices) >= 2 {
		previous := invoices[len(invoices)-2]
		log.Printf("Previous Heroku invoice: %s (Period End: %s, Total: %d, State: %d)", previous.ID, previous.PeriodEnd, previous.Total, previous.State)

		herokuInvoicePreviousTotal.WithLabelValues(previous.PaymentStatus, fmt.Sprintf("%d", previous.State), previous.PeriodEnd).Set(float64(previous.Total) / 100.0)
		herokuInvoicePreviousCharges.WithLabelValues(previous.PeriodEnd).Set(float64(previous.ChargesTotal) / 100.0)
		herokuInvoicePreviousCredits.WithLabelValues(previous.PeriodEnd).Set(float64(previous.CreditsTotal) / 100.0)
		herokuInvoicePreviousAddons.WithLabelValues(previous.PeriodEnd).Set(float64(previous.AddonsTotal) / 100.0)
		herokuInvoicePreviousDatabase.WithLabelValues(previous.PeriodEnd).Set(float64(previous.DatabaseTotal) / 100.0)
		herokuInvoicePreviousPlatform.WithLabelValues(previous.PeriodEnd).Set(float64(previous.PlatformTotal) / 100.0)
		herokuInvoicePreviousDynoUnits.WithLabelValues(previous.PeriodEnd).Set(float64(previous.DynoUnits))
		herokuInvoicePreviousWeightedDynoHours.WithLabelValues(previous.PeriodEnd).Set(float64(previous.WeightedDynoHours))
	}

	return nil
}

func (e *CostExporter) updateHerokuTeamUsageMetrics() error {
	apiKey := os.Getenv("HEROKU_API_KEY")
	teamID := os.Getenv("HEROKU_TEAM_ID")
	if apiKey == "" || teamID == "" {
		log.Println("HEROKU_API_KEY or HEROKU_TEAM_ID not set, skipping Heroku team usage metrics")
		return nil
	}

	// Fetch for previous and current month
	start := time.Now().AddDate(0, -1, 0).Format("2006-01")
	end := time.Now().Format("2006-01")
	url := fmt.Sprintf("https://api.heroku.com/teams/%s/usage/monthly?start=%s&end=%s", teamID, start, end)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch heroku team usage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("heroku api returned status %d: %s", resp.StatusCode, string(body))
	}

	var usages []HerokuTeamUsage
	if err := json.NewDecoder(resp.Body).Decode(&usages); err != nil {
		return fmt.Errorf("failed to decode heroku team usage: %w", err)
	}

	if len(usages) == 0 {
		log.Println("No Heroku team usage found")
		return nil
	}

	// Assuming we want the usage for the requested team/month.
	// The API returns an array, likely containing the requested team's usage.
	// We'll process all returned items just in case.

	herokuTeamUsageCurrentUSD.Reset()
	herokuTeamUsagePreviousUSD.Reset()
	herokuTeamUsageCurrentDynoUnits.Reset()
	herokuTeamUsagePreviousDynoUnits.Reset()

	currentMonth := time.Now().Format("2006-01")
	previousMonth := time.Now().AddDate(0, -1, 0).Format("2006-01")

	for _, usage := range usages {
		log.Printf("Updating metrics for team %s (Month: %s)", usage.Name, usage.Month)

		if usage.Month == currentMonth {
			// Team Totals
			herokuTeamUsageCurrentDynoUnits.WithLabelValues(usage.Name, "_team_total", usage.Month).Set(usage.Dynos)

			herokuTeamUsageCurrentUSD.WithLabelValues(usage.Name, "_team_total", "data", usage.Month).Set(usage.Data)
			herokuTeamUsageCurrentUSD.WithLabelValues(usage.Name, "_team_total", "addons", usage.Month).Set(usage.Addons)
			herokuTeamUsageCurrentUSD.WithLabelValues(usage.Name, "_team_total", "connect", usage.Month).Set(usage.Connect)
			herokuTeamUsageCurrentUSD.WithLabelValues(usage.Name, "_team_total", "partner", usage.Month).Set(usage.Partner)

			// App Breakdowns
			for _, app := range usage.Apps {
				herokuTeamUsageCurrentDynoUnits.WithLabelValues(usage.Name, app.AppName, usage.Month).Set(app.Dynos)

				herokuTeamUsageCurrentUSD.WithLabelValues(usage.Name, app.AppName, "addons", usage.Month).Set(app.Addons)
				herokuTeamUsageCurrentUSD.WithLabelValues(usage.Name, app.AppName, "data", usage.Month).Set(app.Data)
				herokuTeamUsageCurrentUSD.WithLabelValues(usage.Name, app.AppName, "connect", usage.Month).Set(app.Connect)
				herokuTeamUsageCurrentUSD.WithLabelValues(usage.Name, app.AppName, "partner", usage.Month).Set(app.Partner)
			}
		} else if usage.Month == previousMonth {
			// Team Totals
			herokuTeamUsagePreviousDynoUnits.WithLabelValues(usage.Name, "_team_total", usage.Month).Set(usage.Dynos)

			herokuTeamUsagePreviousUSD.WithLabelValues(usage.Name, "_team_total", "data", usage.Month).Set(usage.Data)
			herokuTeamUsagePreviousUSD.WithLabelValues(usage.Name, "_team_total", "addons", usage.Month).Set(usage.Addons)
			herokuTeamUsagePreviousUSD.WithLabelValues(usage.Name, "_team_total", "connect", usage.Month).Set(usage.Connect)
			herokuTeamUsagePreviousUSD.WithLabelValues(usage.Name, "_team_total", "partner", usage.Month).Set(usage.Partner)

			// App Breakdowns
			for _, app := range usage.Apps {
				herokuTeamUsagePreviousDynoUnits.WithLabelValues(usage.Name, app.AppName, usage.Month).Set(app.Dynos)

				herokuTeamUsagePreviousUSD.WithLabelValues(usage.Name, app.AppName, "addons", usage.Month).Set(app.Addons)
				herokuTeamUsagePreviousUSD.WithLabelValues(usage.Name, app.AppName, "data", usage.Month).Set(app.Data)
				herokuTeamUsagePreviousUSD.WithLabelValues(usage.Name, app.AppName, "connect", usage.Month).Set(app.Connect)
				herokuTeamUsagePreviousUSD.WithLabelValues(usage.Name, app.AppName, "partner", usage.Month).Set(app.Partner)
			}
		}
	}

	return nil
}

func (e *CostExporter) updateMetrics(ctx context.Context) error {
	if err := e.updateHerokuMetrics(); err != nil {
		log.Printf("Error updating Heroku invoice metrics: %v", err)
	}

	if err := e.updateHerokuTeamUsageMetrics(); err != nil {
		log.Printf("Error updating Heroku team usage metrics: %v", err)
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

	log.Printf("Starting Heroku Cost Exporter on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
