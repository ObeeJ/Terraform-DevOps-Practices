package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CarbonService struct {
	db    *sql.DB
	cache *redis.Client
}

type CalculateRequest struct {
	Activity  string                 `json:"activity"`
	Weight    float64                `json:"weight,omitempty"`
	Distance  float64                `json:"distance,omitempty"`
	From      string                 `json:"from,omitempty"`
	To        string                 `json:"to,omitempty"`
	Transport string                 `json:"transport,omitempty"`
	Amount    float64                `json:"amount,omitempty"`
	Unit      string                 `json:"unit,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type CalculateResponse struct {
	CarbonFootprint float64                `json:"carbon_footprint"`
	Unit            string                 `json:"unit"`
	Breakdown       map[string]interface{} `json:"breakdown"`
	Suggestions     []string               `json:"suggestions"`
	Calculation     map[string]interface{} `json:"calculation"`
	Timestamp       time.Time              `json:"timestamp"`
}

type EmissionFactor struct {
	ID            int     `json:"id"`
	Activity      string  `json:"activity"`
	TransportMode string  `json:"transport_mode"`
	Factor        float64 `json:"factor"`
	Unit          string  `json:"unit"`
	Source        string  `json:"source"`
}

func NewCarbonService(db *sql.DB, cache *redis.Client) *CarbonService {
	return &CarbonService{
		db:    db,
		cache: cache,
	}
}

func (cs *CarbonService) CalculateCarbon(c *fiber.Ctx) error {
	start := time.Now()

	var req CalculateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request format",
		})
	}

	// Validate required fields
	if req.Activity == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Activity is required",
		})
	}

	// Calculate carbon footprint
	result, err := cs.calculateCarbonFootprint(req)
	if err != nil {
		log.Printf("Calculation error: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to calculate carbon footprint",
		})
	}

	// Store calculation in database
	go cs.storeCalculation(req, result)

	// Track API usage
	go cs.trackAPIUsage("calculate", c.Get("User-ID", "anonymous"), time.Since(start))

	return c.JSON(result)
}

func (cs *CarbonService) calculateCarbonFootprint(req CalculateRequest) (*CalculateResponse, error) {
	// Get emission factor from database
	factor, err := cs.getEmissionFactor(req.Activity, req.Transport)
	if err != nil {
		return nil, err
	}

	var carbonFootprint float64
	var breakdown map[string]interface{}
	var calculation map[string]interface{}

	switch req.Activity {
	case "shipping":
		carbonFootprint, breakdown, calculation = cs.calculateShipping(req, factor)
	case "electricity":
		carbonFootprint, breakdown, calculation = cs.calculateElectricity(req, factor)
	case "fuel":
		carbonFootprint, breakdown, calculation = cs.calculateFuel(req, factor)
	default:
		// Generic calculation
		carbonFootprint = req.Amount * factor.Factor
		breakdown = map[string]interface{}{
			"activity": req.Activity,
			"amount":   req.Amount,
			"factor":   factor.Factor,
			"unit":     factor.Unit,
		}
		calculation = map[string]interface{}{
			"formula": "amount × emission_factor",
			"values":  breakdown,
		}
	}

	// Generate suggestions
	suggestions := cs.generateSuggestions(req, carbonFootprint)

	return &CalculateResponse{
		CarbonFootprint: math.Round(carbonFootprint*1000) / 1000, // Round to 3 decimal places
		Unit:            "kg_co2e",
		Breakdown:       breakdown,
		Suggestions:     suggestions,
		Calculation:     calculation,
		Timestamp:       time.Now(),
	}, nil
}

func (cs *CarbonService) calculateShipping(req CalculateRequest, factor EmissionFactor) (float64, map[string]interface{}, map[string]interface{}) {
	// If distance not provided, estimate based on from/to
	distance := req.Distance
	if distance == 0 && req.From != "" && req.To != "" {
		distance = cs.estimateDistance(req.From, req.To)
	}

	// Calculate: weight (tonnes) × distance (km) × emission factor
	weightTonnes := req.Weight / 1000.0
	carbonFootprint := weightTonnes * distance * factor.Factor

	breakdown := map[string]interface{}{
		"weight_kg":       req.Weight,
		"weight_tonnes":   weightTonnes,
		"distance_km":     distance,
		"transport_mode":  req.Transport,
		"emission_factor": factor.Factor,
		"from":            req.From,
		"to":              req.To,
	}

	calculation := map[string]interface{}{
		"formula": "weight_tonnes × distance_km × emission_factor",
		"values":  breakdown,
		"result":  carbonFootprint,
	}

	return carbonFootprint, breakdown, calculation
}

func (cs *CarbonService) calculateElectricity(req CalculateRequest, factor EmissionFactor) (float64, map[string]interface{}, map[string]interface{}) {
	// Calculate: amount (kWh) × emission factor
	carbonFootprint := req.Amount * factor.Factor

	breakdown := map[string]interface{}{
		"energy_kwh":      req.Amount,
		"energy_source":   req.Transport,
		"emission_factor": factor.Factor,
		"grid_mix":        "regional_average", // Could be enhanced with location-specific data
	}

	calculation := map[string]interface{}{
		"formula": "energy_kwh × emission_factor",
		"values":  breakdown,
		"result":  carbonFootprint,
	}

	return carbonFootprint, breakdown, calculation
}

func (cs *CarbonService) calculateFuel(req CalculateRequest, factor EmissionFactor) (float64, map[string]interface{}, map[string]interface{}) {
	// Calculate: amount (liters) × emission factor
	carbonFootprint := req.Amount * factor.Factor

	breakdown := map[string]interface{}{
		"fuel_liters":     req.Amount,
		"fuel_type":       req.Transport,
		"emission_factor": factor.Factor,
		"combustion":      "direct_emissions",
	}

	calculation := map[string]interface{}{
		"formula": "fuel_liters × emission_factor",
		"values":  breakdown,
		"result":  carbonFootprint,
	}

	return carbonFootprint, breakdown, calculation
}

func (cs *CarbonService) getEmissionFactor(activity, transport string) (EmissionFactor, error) {
	var factor EmissionFactor

	query := `
		SELECT id, activity, transport_mode, factor, unit, source 
		FROM emission_factors 
		WHERE activity = $1
	`

	args := []interface{}{activity}

	if transport != "" {
		query += " AND transport_mode = $2"
		args = append(args, transport)
	}

	query += " ORDER BY id LIMIT 1"

	err := cs.db.QueryRow(query, args...).Scan(
		&factor.ID,
		&factor.Activity,
		&factor.TransportMode,
		&factor.Factor,
		&factor.Unit,
		&factor.Source,
	)

	if err != nil {
		// Fallback to default factors if database is unavailable
		return cs.getDefaultEmissionFactor(activity, transport), nil
	}

	return factor, nil
}

func (cs *CarbonService) getDefaultEmissionFactor(activity, transport string) EmissionFactor {
	// Fallback emission factors for when database is unavailable
	defaultFactors := map[string]map[string]EmissionFactor{
		"shipping": {
			"air":  {Activity: "shipping", TransportMode: "air", Factor: 0.996, Unit: "kg_co2e_per_tonne_km", Source: "IPCC 2023"},
			"sea":  {Activity: "shipping", TransportMode: "sea", Factor: 0.015, Unit: "kg_co2e_per_tonne_km", Source: "IPCC 2023"},
			"road": {Activity: "shipping", TransportMode: "road", Factor: 0.209, Unit: "kg_co2e_per_tonne_km", Source: "IPCC 2023"},
			"rail": {Activity: "shipping", TransportMode: "rail", Factor: 0.028, Unit: "kg_co2e_per_tonne_km", Source: "IPCC 2023"},
		},
		"electricity": {
			"grid":  {Activity: "electricity", TransportMode: "grid", Factor: 0.525, Unit: "kg_co2e_per_kwh", Source: "IEA 2023"},
			"solar": {Activity: "electricity", TransportMode: "solar", Factor: 0.041, Unit: "kg_co2e_per_kwh", Source: "IPCC 2023"},
			"wind":  {Activity: "electricity", TransportMode: "wind", Factor: 0.011, Unit: "kg_co2e_per_kwh", Source: "IPCC 2023"},
		},
		"fuel": {
			"gasoline":    {Activity: "fuel", TransportMode: "gasoline", Factor: 2.31, Unit: "kg_co2e_per_liter", Source: "EPA 2023"},
			"diesel":      {Activity: "fuel", TransportMode: "diesel", Factor: 2.68, Unit: "kg_co2e_per_liter", Source: "EPA 2023"},
			"natural_gas": {Activity: "fuel", TransportMode: "natural_gas", Factor: 0.202, Unit: "kg_co2e_per_kwh", Source: "EPA 2023"},
		},
	}

	if activityFactors, exists := defaultFactors[activity]; exists {
		if factor, exists := activityFactors[transport]; exists {
			return factor
		}
		// Return first available factor for the activity
		for _, factor := range activityFactors {
			return factor
		}
	}

	// Ultimate fallback
	return EmissionFactor{
		Activity:      activity,
		TransportMode: transport,
		Factor:        1.0,
		Unit:          "kg_co2e_per_unit",
		Source:        "Default",
	}
}

func (cs *CarbonService) estimateDistance(from, to string) float64 {
	// Simplified distance estimation - in production, use a proper geocoding service
	distances := map[string]map[string]float64{
		"NYC": {
			"London":     5585, // km
			"Paris":      5837,
			"Tokyo":      10847,
			"Sydney":     15993,
			"LosAngeles": 3944,
		},
		"London": {
			"NYC":   5585,
			"Paris": 344,
			"Tokyo": 9561,
		},
		"Paris": {
			"NYC":    5837,
			"London": 344,
			"Tokyo":  9714,
		},
	}

	if fromDistances, exists := distances[from]; exists {
		if distance, exists := fromDistances[to]; exists {
			return distance
		}
	}

	// Default distance for unknown routes
	return 1000.0
}

func (cs *CarbonService) generateSuggestions(req CalculateRequest, carbonFootprint float64) []string {
	var suggestions []string

	switch req.Activity {
	case "shipping":
		if req.Transport == "air" {
			suggestions = append(suggestions, "Consider sea freight to reduce emissions by 98%")
			suggestions = append(suggestions, "Use rail transport when possible for 97% reduction")
		}
		if req.Transport == "road" {
			suggestions = append(suggestions, "Switch to rail transport for 87% emissions reduction")
			suggestions = append(suggestions, "Optimize routes to reduce distance")
		}
	case "electricity":
		if req.Transport == "grid" {
			suggestions = append(suggestions, "Switch to renewable energy for 92% reduction")
			suggestions = append(suggestions, "Install solar panels for clean energy")
		}
	case "fuel":
		suggestions = append(suggestions, "Consider electric vehicles for zero direct emissions")
		suggestions = append(suggestions, "Use biofuels to reduce carbon intensity")
	}

	// General suggestions based on footprint size
	if carbonFootprint > 1000 {
		suggestions = append(suggestions, "This is a high-impact activity - consider carbon offsetting")
	}

	return suggestions
}

func (cs *CarbonService) storeCalculation(req CalculateRequest, result *CalculateResponse) {
	inputJSON, _ := json.Marshal(req)

	_, err := cs.db.Exec(`
		INSERT INTO calculations (activity, input_data, carbon_footprint, unit, user_id)
		VALUES ($1, $2, $3, $4, $5)
	`, req.Activity, inputJSON, result.CarbonFootprint, result.Unit, uuid.New().String())

	if err != nil {
		log.Printf("Failed to store calculation: %v", err)
	}
}

func (cs *CarbonService) trackAPIUsage(endpoint, userID string, responseTime time.Duration) {
	_, err := cs.db.Exec(`
		INSERT INTO api_usage (endpoint, user_id, response_time_ms)
		VALUES ($1, $2, $3)
	`, endpoint, userID, responseTime.Milliseconds())

	if err != nil {
		log.Printf("Failed to track API usage: %v", err)
	}
}

func (cs *CarbonService) GetActivities(c *fiber.Ctx) error {
	activities := map[string]interface{}{
		"shipping": map[string]interface{}{
			"description":     "Calculate carbon footprint for freight transport",
			"transport_modes": []string{"air", "sea", "road", "rail"},
			"required_fields": []string{"activity", "weight", "distance_or_locations", "transport"},
			"example": map[string]interface{}{
				"activity":  "shipping",
				"weight":    500,
				"from":      "NYC",
				"to":        "London",
				"transport": "air",
			},
		},
		"electricity": map[string]interface{}{
			"description":     "Calculate carbon footprint for electricity consumption",
			"energy_sources":  []string{"grid", "solar", "wind", "coal", "gas"},
			"required_fields": []string{"activity", "amount", "transport"},
			"example": map[string]interface{}{
				"activity":  "electricity",
				"amount":    100,
				"unit":      "kwh",
				"transport": "grid",
			},
		},
		"fuel": map[string]interface{}{
			"description":     "Calculate carbon footprint for fuel consumption",
			"fuel_types":      []string{"gasoline", "diesel", "natural_gas"},
			"required_fields": []string{"activity", "amount", "transport"},
			"example": map[string]interface{}{
				"activity":  "fuel",
				"amount":    50,
				"unit":      "liters",
				"transport": "gasoline",
			},
		},
	}

	return c.JSON(fiber.Map{
		"activities": activities,
		"total":      len(activities),
	})
}

func (cs *CarbonService) GetEmissionFactors(c *fiber.Ctx) error {
	rows, err := cs.db.Query(`
		SELECT activity, transport_mode, factor, unit, source
		FROM emission_factors
		ORDER BY activity, transport_mode
	`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch emission factors",
		})
	}
	defer rows.Close()

	var factors []EmissionFactor
	for rows.Next() {
		var factor EmissionFactor
		err := rows.Scan(
			&factor.Activity,
			&factor.TransportMode,
			&factor.Factor,
			&factor.Unit,
			&factor.Source,
		)
		if err != nil {
			continue
		}
		factors = append(factors, factor)
	}

	return c.JSON(fiber.Map{
		"emission_factors": factors,
		"total":            len(factors),
		"source":           "IPCC 2023, IEA 2023, EPA 2023",
	})
}

func (cs *CarbonService) GetAnalytics(c *fiber.Ctx) error {
	var totalCalculations int
	var avgResponseTime float64
	var totalCarbonCalculated float64

	// Get total calculations
	cs.db.QueryRow("SELECT COUNT(*) FROM calculations").Scan(&totalCalculations)

	// Get average response time
	cs.db.QueryRow("SELECT AVG(response_time_ms) FROM api_usage").Scan(&avgResponseTime)

	// Get total carbon calculated
	cs.db.QueryRow("SELECT SUM(carbon_footprint) FROM calculations").Scan(&totalCarbonCalculated)

	// Get top activities
	rows, _ := cs.db.Query(`
		SELECT activity, COUNT(*) as count
		FROM calculations
		GROUP BY activity
		ORDER BY count DESC
		LIMIT 5
	`)
	defer rows.Close()

	topActivities := make(map[string]int)
	for rows.Next() {
		var activity string
		var count int
		rows.Scan(&activity, &count)
		topActivities[activity] = count
	}

	return c.JSON(fiber.Map{
		"analytics": map[string]interface{}{
			"total_calculations":      totalCalculations,
			"avg_response_time_ms":    math.Round(avgResponseTime*100) / 100,
			"total_carbon_calculated": math.Round(totalCarbonCalculated*100) / 100,
			"top_activities":          topActivities,
		},
		"timestamp": time.Now(),
	})
}
