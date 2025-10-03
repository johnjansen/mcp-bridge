package bdd

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name:                "mcp-bridge-bdd",
		ScenarioInitializer: InitializeScenario,
	}
	if suite.Run() != 0 {
		t.Fail()
	}
}
