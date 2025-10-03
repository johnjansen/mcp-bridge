package bdd

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name: "mcp-bridge-bdd",
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			InitializeScenario(ctx)
			InitializeDebugScenario(ctx)
		},
		Options: &godog.Options{
			Format:    "pretty",
			Paths:     []string{"../features"},
			Randomize: 0,
		},
	}
	if suite.Run() != 0 {
		t.Fail()
	}
}
