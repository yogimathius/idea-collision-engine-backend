package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// TestRunner provides comprehensive test execution for the Idea Collision Engine API
func main() {
	fmt.Println("🧪 Idea Collision Engine - Test Suite Runner")
	fmt.Println("============================================")
	
	startTime := time.Now()
	
	// Test suites to run
	testSuites := []TestSuite{
		{
			Name:        "Unit Tests",
			Command:     "go test -v ./internal/...",
			Description: "Core business logic and algorithms",
		},
		{
			Name:        "Integration Tests", 
			Command:     "go test -v -tags=integration ./...",
			Description: "API endpoints and database integration",
		},
		{
			Name:        "Benchmark Tests",
			Command:     "go test -bench=. -benchmem ./internal/collision/...",
			Description: "Performance benchmarks for collision generation",
		},
		{
			Name:        "Race Condition Tests",
			Command:     "go test -race ./internal/...",
			Description: "Concurrent access safety",
		},
		{
			Name:        "Coverage Report",
			Command:     "go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html",
			Description: "Test coverage analysis",
		},
	}
	
	var results []TestResult
	
	for _, suite := range testSuites {
		fmt.Printf("\n📋 Running %s\n", suite.Name)
		fmt.Printf("   %s\n", suite.Description)
		fmt.Println("   " + strings.Repeat("-", 50))
		
		result := runTestSuite(suite)
		results = append(results, result)
		
		if result.Success {
			fmt.Printf("   ✅ %s completed successfully (%.2fs)\n", suite.Name, result.Duration.Seconds())
		} else {
			fmt.Printf("   ❌ %s failed (%.2fs)\n", suite.Name, result.Duration.Seconds())
			if result.Output != "" {
				fmt.Printf("   Error: %s\n", result.Output)
			}
		}
	}
	
	// Print summary
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 TEST SUMMARY")
	fmt.Println(strings.Repeat("=", 60))
	
	passed := 0
	failed := 0
	
	for _, result := range results {
		status := "❌ FAILED"
		if result.Success {
			status = "✅ PASSED"
			passed++
		} else {
			failed++
		}
		
		fmt.Printf("%-25s %s (%.2fs)\n", result.SuiteName, status, result.Duration.Seconds())
	}
	
	totalTime := time.Since(startTime)
	
	fmt.Printf("\n📈 Results: %d passed, %d failed\n", passed, failed)
	fmt.Printf("⏱️  Total execution time: %.2fs\n", totalTime.Seconds())
	
	if failed > 0 {
		fmt.Printf("\n⚠️  %d test suite(s) failed. Review the output above.\n", failed)
		os.Exit(1)
	} else {
		fmt.Println("\n🎉 All test suites passed!")
		
		// Additional success information
		fmt.Println("\n📋 Next Steps:")
		fmt.Println("   • Review coverage.html for test coverage details")
		fmt.Println("   • Run 'make benchmark' for detailed performance analysis")
		fmt.Println("   • Deploy with confidence! 🚀")
	}
}

type TestSuite struct {
	Name        string
	Command     string
	Description string
}

type TestResult struct {
	SuiteName string
	Success   bool
	Duration  time.Duration
	Output    string
}

func runTestSuite(suite TestSuite) TestResult {
	start := time.Now()
	
	// Split command into parts
	parts := strings.Fields(suite.Command)
	if len(parts) == 0 {
		return TestResult{
			SuiteName: suite.Name,
			Success:   false,
			Duration:  time.Since(start),
			Output:    "Empty command",
		}
	}
	
	// Handle compound commands (with &&)
	if strings.Contains(suite.Command, "&&") {
		// For compound commands, execute as shell command
		cmd := exec.Command("sh", "-c", suite.Command)
		output, err := cmd.CombinedOutput()
		
		return TestResult{
			SuiteName: suite.Name,
			Success:   err == nil,
			Duration:  time.Since(start),
			Output:    string(output),
		}
	}
	
	// Single command
	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	
	return TestResult{
		SuiteName: suite.Name,
		Success:   err == nil,
		Duration:  time.Since(start),
		Output:    string(output),
	}
}