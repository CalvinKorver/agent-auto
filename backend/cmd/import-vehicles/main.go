package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"carbuyer/internal/db"
	"carbuyer/internal/db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VehicleRow struct {
	ID                  string
	Make                string
	Model               string
	Year                string
	Trim                string
	TrimDescription     string
	BaseMSRP            string
	ImageURL            string
	BodyType            string
	LengthIn            string
	WidthIn             string
	HeightIn            string
	WheelbaseIn         string
	CurbWeightLbs       string
	Cylinders           string
	EngineSizeL         string
	HorsepowerHP        string
	HorsepowerRPM       string
	TorqueFtLbs         string
	TorqueRPM           string
	DriveType           string
	Transmission        string
	EngineType          string
	FuelType            string
	FuelTankCapacityGal  string
	EPACombinedMPG      string
	EPACityHighwayMPG   string
	RangeMiles          string
	CountryOfOrigin      string
	CarClassification    string
	PlatformCode        string
	DateAdded           string
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <path-to-sql-file> [database-url]")
	}

	sqlFilePath := os.Args[1]
	databaseURL := os.Getenv("DATABASE_URL")
	if len(os.Args) > 2 {
		databaseURL = os.Args[2]
	}

	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable or argument is required")
	}

	log.Println("Connecting to database...")
	database, err := db.NewDatabase(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("Running migrations...")
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Printf("Reading SQL file: %s", sqlFilePath)
	file, err := os.Open(sqlFilePath)
	if err != nil {
		log.Fatalf("Failed to open SQL file: %v", err)
	}
	defer file.Close()

	rows, err := parseSQLFile(file)
	if err != nil {
		log.Fatalf("Failed to parse SQL file: %v", err)
	}

	log.Printf("Parsed %d rows from SQL file", len(rows))

	// Build unique makes and models
	makeMap := make(map[string]bool)
	modelMap := make(map[string]models.Model) // key: "make|model"
	trimRows := []VehicleRow{}

	for _, row := range rows {
		// Skip header row
		if row.ID == "ID" || row.Make == "Make" {
			continue
		}

		// Skip empty rows
		if row.Make == "" || row.Model == "" {
			continue
		}

		makeMap[row.Make] = true
		modelKey := fmt.Sprintf("%s|%s", row.Make, row.Model)
		if _, exists := modelMap[modelKey]; !exists {
			modelMap[modelKey] = models.Model{
				Name: row.Model,
			}
		}
		trimRows = append(trimRows, row)
	}

	log.Printf("Found %d unique makes", len(makeMap))
	log.Printf("Found %d unique models", len(modelMap))
	log.Printf("Found %d trim rows", len(trimRows))

	// Import makes
	log.Println("Importing makes...")
	makeIDMap, err := importMakes(database.DB, makeMap)
	if err != nil {
		log.Fatalf("Failed to import makes: %v", err)
	}
	log.Printf("Imported %d makes", len(makeIDMap))

	// Import models
	log.Println("Importing models...")
	modelIDMap, err := importModels(database.DB, modelMap, makeIDMap)
	if err != nil {
		log.Fatalf("Failed to import models: %v", err)
	}
	log.Printf("Imported %d models", len(modelIDMap))

	// Import vehicle trims
	log.Println("Importing vehicle trims...")
	imported, err := importVehicleTrims(database.DB, trimRows, modelIDMap)
	if err != nil {
		log.Fatalf("Failed to import vehicle trims: %v", err)
	}
	log.Printf("Imported %d vehicle trims", imported)

	log.Println("Import completed successfully!")
}

func parseSQLFile(file *os.File) ([]VehicleRow, error) {
	var rows []VehicleRow
	scanner := bufio.NewScanner(file)

	// Pattern to match INSERT statements
	insertPattern := regexp.MustCompile(`INSERT INTO.*VALUES`)
	inInsert := false
	var currentValues strings.Builder

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if insertPattern.MatchString(line) {
			inInsert = true
			// Extract values part
			valuesIdx := strings.Index(line, "VALUES")
			if valuesIdx != -1 {
				currentValues.WriteString(line[valuesIdx+6:])
			}
			continue
		}

		if inInsert {
			currentValues.WriteString(" " + line)
			// Check if this line ends the INSERT statement
			if strings.HasSuffix(line, ";") || strings.HasSuffix(line, ");") {
				valuesStr := currentValues.String()
				parsedRows := parseValues(valuesStr)
				rows = append(rows, parsedRows...)
				currentValues.Reset()
				inInsert = false
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return rows, nil
}

func parseValues(valuesStr string) []VehicleRow {
	var rows []VehicleRow

	// Remove leading/trailing whitespace and parentheses
	valuesStr = strings.TrimSpace(valuesStr)
	if strings.HasPrefix(valuesStr, "(") {
		valuesStr = valuesStr[1:]
	}
	if strings.HasSuffix(valuesStr, ");") {
		valuesStr = valuesStr[:len(valuesStr)-2]
	} else if strings.HasSuffix(valuesStr, ";") {
		valuesStr = valuesStr[:len(valuesStr)-1]
	}

	// Split by ),( to get individual rows
	rowPattern := regexp.MustCompile(`\),\s*\(`)
	rowStrings := rowPattern.Split(valuesStr, -1)

	for _, rowStr := range rowStrings {
		rowStr = strings.TrimSpace(rowStr)
		if strings.HasPrefix(rowStr, "(") {
			rowStr = rowStr[1:]
		}
		if strings.HasSuffix(rowStr, ")") {
			rowStr = rowStr[:len(rowStr)-1]
		}

		row := parseRow(rowStr)
		rows = append(rows, row)
	}

	return rows
}

func parseRow(rowStr string) VehicleRow {
	// Parse comma-separated values, handling quoted strings
	values := parseCSVLine(rowStr)

	row := VehicleRow{}
	if len(values) > 0 {
		row.ID = unquote(values[0])
	}
	if len(values) > 1 {
		row.Make = unquote(values[1])
	}
	if len(values) > 2 {
		row.Model = unquote(values[2])
	}
	if len(values) > 3 {
		row.Year = unquote(values[3])
	}
	if len(values) > 4 {
		row.Trim = unquote(values[4])
	}
	if len(values) > 5 {
		row.TrimDescription = unquote(values[5])
	}
	if len(values) > 6 {
		row.BaseMSRP = unquote(values[6])
	}
	if len(values) > 7 {
		row.ImageURL = unquote(values[7])
	}
	if len(values) > 8 {
		row.BodyType = unquote(values[8])
	}
	if len(values) > 9 {
		row.LengthIn = unquote(values[9])
	}
	if len(values) > 10 {
		row.WidthIn = unquote(values[10])
	}
	if len(values) > 11 {
		row.HeightIn = unquote(values[11])
	}
	if len(values) > 12 {
		row.WheelbaseIn = unquote(values[12])
	}
	if len(values) > 13 {
		row.CurbWeightLbs = unquote(values[13])
	}
	if len(values) > 14 {
		row.Cylinders = unquote(values[14])
	}
	if len(values) > 15 {
		row.EngineSizeL = unquote(values[15])
	}
	if len(values) > 16 {
		row.HorsepowerHP = unquote(values[16])
	}
	if len(values) > 17 {
		row.HorsepowerRPM = unquote(values[17])
	}
	if len(values) > 18 {
		row.TorqueFtLbs = unquote(values[18])
	}
	if len(values) > 19 {
		row.TorqueRPM = unquote(values[19])
	}
	if len(values) > 20 {
		row.DriveType = unquote(values[20])
	}
	if len(values) > 21 {
		row.Transmission = unquote(values[21])
	}
	if len(values) > 22 {
		row.EngineType = unquote(values[22])
	}
	if len(values) > 23 {
		row.FuelType = unquote(values[23])
	}
	if len(values) > 24 {
		row.FuelTankCapacityGal = unquote(values[24])
	}
	if len(values) > 25 {
		row.EPACombinedMPG = unquote(values[25])
	}
	if len(values) > 26 {
		row.EPACityHighwayMPG = unquote(values[26])
	}
	if len(values) > 27 {
		row.RangeMiles = unquote(values[27])
	}
	if len(values) > 28 {
		row.CountryOfOrigin = unquote(values[28])
	}
	if len(values) > 29 {
		row.CarClassification = unquote(values[29])
	}
	if len(values) > 30 {
		row.PlatformCode = unquote(values[30])
	}
	if len(values) > 31 {
		row.DateAdded = unquote(values[31])
	}

	return row
}

func parseCSVLine(line string) []string {
	var values []string
	var current strings.Builder
	inQuotes := false
	escapeNext := false

	for i, char := range line {
		if escapeNext {
			current.WriteRune(char)
			escapeNext = false
			continue
		}

		if char == '\\' {
			escapeNext = true
			continue
		}

		if char == '\'' {
			if inQuotes && i+1 < len(line) && line[i+1] == '\'' {
				// Escaped quote
				current.WriteRune('\'')
				i++ // Skip next quote
				continue
			}
			inQuotes = !inQuotes
			continue
		}

		if char == ',' && !inQuotes {
			values = append(values, current.String())
			current.Reset()
			continue
		}

		current.WriteRune(char)
	}

	// Add last value
	if current.Len() > 0 || !inQuotes {
		values = append(values, current.String())
	}

	return values
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		return s[1 : len(s)-1]
	}
	return s
}

func importMakes(db *gorm.DB, makeMap map[string]bool) (map[string]uuid.UUID, error) {
	makeIDMap := make(map[string]uuid.UUID)

	for makeName := range makeMap {
		var make models.Make
		err := db.Where("name = ?", makeName).First(&make).Error
		if err == nil {
			// Make already exists
			makeIDMap[makeName] = make.ID
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to check make: %w", err)
		}

		// Create new make
		make = models.Make{Name: makeName}
		if err := db.Create(&make).Error; err != nil {
			return nil, fmt.Errorf("failed to create make %s: %w", makeName, err)
		}
		makeIDMap[makeName] = make.ID
	}

	return makeIDMap, nil
}

func importModels(db *gorm.DB, modelMap map[string]models.Model, makeIDMap map[string]uuid.UUID) (map[string]uuid.UUID, error) {
	modelIDMap := make(map[string]uuid.UUID)

	for modelKey, model := range modelMap {
		parts := strings.Split(modelKey, "|")
		if len(parts) != 2 {
			continue
		}
		makeName := parts[0]
		makeID, exists := makeIDMap[makeName]
		if !exists {
			log.Printf("Warning: Make %s not found for model %s", makeName, model.Name)
			continue
		}

		var existingModel models.Model
		err := db.Where("make_id = ? AND name = ?", makeID, model.Name).First(&existingModel).Error
		if err == nil {
			// Model already exists
			modelIDMap[modelKey] = existingModel.ID
			continue
		}
		if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to check model: %w", err)
		}

		// Create new model
		model.MakeID = makeID
		if err := db.Create(&model).Error; err != nil {
			return nil, fmt.Errorf("failed to create model %s: %w", model.Name, err)
		}
		modelIDMap[modelKey] = model.ID
	}

	return modelIDMap, nil
}

func importVehicleTrims(db *gorm.DB, trimRows []VehicleRow, modelIDMap map[string]uuid.UUID) (int, error) {
	imported := 0
	batchSize := 100

	for i := 0; i < len(trimRows); i += batchSize {
		end := i + batchSize
		if end > len(trimRows) {
			end = len(trimRows)
		}

		batch := trimRows[i:end]
		for _, row := range batch {
			modelKey := fmt.Sprintf("%s|%s", row.Make, row.Model)
			modelID, exists := modelIDMap[modelKey]
			if !exists {
				log.Printf("Warning: Model %s|%s not found, skipping trim", row.Make, row.Model)
				continue
			}

			year, err := strconv.Atoi(row.Year)
			if err != nil {
				log.Printf("Warning: Invalid year %s, skipping row", row.Year)
				continue
			}

			trim := models.VehicleTrim{
				ModelID:   modelID,
				Year:      year,
				TrimName:  row.Trim,
				SourceID:  sql.NullString{String: row.ID, Valid: row.ID != ""},
			}

			// Set nullable fields
			if row.TrimDescription != "" {
				trim.TrimDescription = sql.NullString{String: row.TrimDescription, Valid: true}
			}
			if row.BaseMSRP != "" {
				msrp := parseMSRP(row.BaseMSRP)
				if msrp > 0 {
					trim.BaseMSRP = sql.NullFloat64{Float64: msrp, Valid: true}
				}
			}
			if row.ImageURL != "" {
				trim.ImageURLs = sql.NullString{String: row.ImageURL, Valid: true}
			}
			if row.BodyType != "" {
				trim.BodyType = sql.NullString{String: row.BodyType, Valid: true}
			}

			// Parse numeric fields
			if row.LengthIn != "" {
				if val, err := strconv.ParseFloat(row.LengthIn, 64); err == nil {
					trim.LengthIn = sql.NullFloat64{Float64: val, Valid: true}
				}
			}
			if row.WidthIn != "" {
				if val, err := strconv.ParseFloat(row.WidthIn, 64); err == nil {
					trim.WidthIn = sql.NullFloat64{Float64: val, Valid: true}
				}
			}
			if row.HeightIn != "" {
				if val, err := strconv.ParseFloat(row.HeightIn, 64); err == nil {
					trim.HeightIn = sql.NullFloat64{Float64: val, Valid: true}
				}
			}
			if row.WheelbaseIn != "" {
				if val, err := strconv.ParseFloat(row.WheelbaseIn, 64); err == nil {
					trim.WheelbaseIn = sql.NullFloat64{Float64: val, Valid: true}
				}
			}
			if row.CurbWeightLbs != "" {
				if val, err := strconv.Atoi(row.CurbWeightLbs); err == nil {
					trim.CurbWeightLbs = sql.NullInt32{Int32: int32(val), Valid: true}
				}
			}
			if row.EngineSizeL != "" {
				if val, err := strconv.ParseFloat(row.EngineSizeL, 64); err == nil {
					trim.EngineSizeL = sql.NullFloat64{Float64: val, Valid: true}
				}
			}
			if row.HorsepowerHP != "" {
				if val, err := strconv.Atoi(row.HorsepowerHP); err == nil {
					trim.HorsepowerHP = sql.NullInt32{Int32: int32(val), Valid: true}
				}
			}
			if row.HorsepowerRPM != "" {
				if val, err := strconv.Atoi(row.HorsepowerRPM); err == nil {
					trim.HorsepowerRPM = sql.NullInt32{Int32: int32(val), Valid: true}
				}
			}
			if row.TorqueFtLbs != "" {
				if val, err := strconv.Atoi(row.TorqueFtLbs); err == nil {
					trim.TorqueFtLbs = sql.NullInt32{Int32: int32(val), Valid: true}
				}
			}
			if row.TorqueRPM != "" {
				if val, err := strconv.Atoi(row.TorqueRPM); err == nil {
					trim.TorqueRPM = sql.NullInt32{Int32: int32(val), Valid: true}
				}
			}
			if row.FuelTankCapacityGal != "" {
				if val, err := strconv.ParseFloat(row.FuelTankCapacityGal, 64); err == nil {
					trim.FuelTankCapacityGal = sql.NullFloat64{Float64: val, Valid: true}
				}
			}
			if row.EPACombinedMPG != "" {
				if val, err := strconv.Atoi(row.EPACombinedMPG); err == nil {
					trim.EPACombinedMPG = sql.NullInt32{Int32: int32(val), Valid: true}
				}
			}

			// Set string fields
			if row.Cylinders != "" {
				trim.Cylinders = sql.NullString{String: row.Cylinders, Valid: true}
			}
			if row.DriveType != "" {
				trim.DriveType = sql.NullString{String: row.DriveType, Valid: true}
			}
			if row.Transmission != "" {
				trim.Transmission = sql.NullString{String: row.Transmission, Valid: true}
			}
			if row.EngineType != "" {
				trim.EngineType = sql.NullString{String: row.EngineType, Valid: true}
			}
			if row.FuelType != "" {
				trim.FuelType = sql.NullString{String: row.FuelType, Valid: true}
			}
			if row.EPACityHighwayMPG != "" {
				trim.EPACityHighwayMPG = sql.NullString{String: row.EPACityHighwayMPG, Valid: true}
			}
			if row.RangeMiles != "" {
				trim.RangeMiles = sql.NullString{String: row.RangeMiles, Valid: true}
			}
			if row.CountryOfOrigin != "" {
				trim.CountryOfOrigin = sql.NullString{String: row.CountryOfOrigin, Valid: true}
			}
			if row.CarClassification != "" {
				trim.CarClassification = sql.NullString{String: row.CarClassification, Valid: true}
			}
			if row.PlatformCode != "" {
				trim.PlatformCode = sql.NullString{String: row.PlatformCode, Valid: true}
			}

			if err := db.Create(&trim).Error; err != nil {
				log.Printf("Warning: Failed to create trim: %v", err)
				continue
			}
			imported++
		}

		if (i+batchSize)%1000 == 0 {
			log.Printf("Processed %d/%d rows...", i+batchSize, len(trimRows))
		}
	}

	return imported, nil
}

func parseMSRP(msrpStr string) float64 {
	// Remove $ and commas, then parse
	msrpStr = strings.ReplaceAll(msrpStr, "$", "")
	msrpStr = strings.ReplaceAll(msrpStr, ",", "")
	msrpStr = strings.TrimSpace(msrpStr)
	val, err := strconv.ParseFloat(msrpStr, 64)
	if err != nil {
		return 0
	}
	return val
}

