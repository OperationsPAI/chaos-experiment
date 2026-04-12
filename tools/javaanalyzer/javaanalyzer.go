package javaanalyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ClassMethodEntry represents a class-method pair from the Java analysis
type ClassMethodEntry struct {
	AppName    string         `json:"app,omitempty"`
	ClassName  string         `json:"className"`
	MethodName string         `json:"methodName"`
	Mutations  []MutationSpec `json:"mutations,omitempty"`
}

// UnmarshalJSON supports both legacy extractor keys (className/methodName)
// and new keys (class/method) that include mutation details.
func (e *ClassMethodEntry) UnmarshalJSON(data []byte) error {
	type rawEntry struct {
		AppName          string         `json:"app"`
		AppNameLegacy    string         `json:"appName"`
		ClassName        string         `json:"className"`
		ClassNameLegacy  string         `json:"class"`
		MethodName       string         `json:"methodName"`
		MethodNameLegacy string         `json:"method"`
		Mutations        []MutationSpec `json:"mutations"`
	}

	var raw rawEntry
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	e.AppName = raw.AppName
	if e.AppName == "" {
		e.AppName = raw.AppNameLegacy
	}

	e.ClassName = raw.ClassName
	if e.ClassName == "" {
		e.ClassName = raw.ClassNameLegacy
	}

	e.MethodName = raw.MethodName
	if e.MethodName == "" {
		e.MethodName = raw.MethodNameLegacy
	}

	e.Mutations = raw.Mutations
	return nil
}

// PathResult represents the results for a specific path
type PathResult struct {
	PathName string             `json:"pathName"`
	Methods  []ClassMethodEntry `json:"methods"`
}

// AnalyzeJavaPath analyzes a single Java source path and returns the method entries
func AnalyzeJavaPath(sourcePath string, jarPath string) ([]ClassMethodEntry, error) {
	// Create a temporary file for output
	tempFile, err := os.CreateTemp("", "java-analysis-*.json")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	tempFileName := tempFile.Name()
	if err := tempFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temp file: %w", err)
	}
	defer func() {
		if err := os.Remove(tempFileName); err != nil {
			// Log the error but don't fail the function
			fmt.Fprintf(os.Stderr, "warning: failed to remove temp file %s: %v\n", tempFileName, err)
		}
	}()

	// Run the Java analyzer
	cmd := exec.Command("java", "-jar", jarPath, sourcePath, tempFileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error running Java analyzer: %w", err)
	}

	// Read the results
	data, err := os.ReadFile(tempFileName)
	if err != nil {
		return nil, fmt.Errorf("error reading output file: %w", err)
	}

	var entries []ClassMethodEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return entries, nil
}

// AnalyzeJavaPaths analyzes multiple Java source paths and returns path results
func AnalyzeJavaPaths(sourcePaths []string) ([]PathResult, error) {
	// Absolute path to the analyzer JAR
	jarPath := "tools/javaanalyzer/method-extractor.jar"

	// Check if the JAR exists
	if _, err := os.Stat(jarPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("analyzer JAR not found at %s, please build with 'mvn package'", jarPath)
	}

	var results []PathResult

	for _, path := range sourcePaths {
		if path == "" {
			continue
		}

		// Get the last directory name from the path
		pathName := filepath.Base(path)

		// Analyze the path
		entries, err := AnalyzeJavaPath(path, jarPath)
		if err != nil {
			return nil, fmt.Errorf("error analyzing path %s: %w", path, err)
		}

		// Add to results
		result := PathResult{
			PathName: pathName,
			Methods:  entries,
		}
		results = append(results, result)
	}

	return results, nil
}

// SaveResultsToFile saves the analysis results to the specified JSON file
func SaveResultsToFile(results []PathResult, outputFile string) error {
	// Create parent directory if it doesn't exist
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Marshal to JSON
	resultJSON, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("error creating JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(outputFile, resultJSON, 0644); err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}

	return nil
}
