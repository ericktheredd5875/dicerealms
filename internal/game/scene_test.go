package game

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ericktheredd5875/dicerealms/config"
)

func TestSceneLogEntryAndSummary(t *testing.T) {
	scene := &Scene{
		Title:     "The Arrival",
		Mood:      "Tense",
		StartedBy: "DM",
		StartedAt: time.Now(),
	}

	scene.LogEntry("A stranger appears at the edge of town.")
	scene.LogEntry("The wind picks up suddenly.")

	if len(scene.Log) != 2 {
		t.Errorf("Expected 2 log entries, got %d", len(scene.Log))
	}

	summary := scene.Summary()
	if !strings.Contains(summary, "The Arrival") || !strings.Contains(summary, "A stranger appears") {
		t.Errorf("Summary missing expected content:\n%s", summary)
	}
}

func TestSanitizeFileName(t *testing.T) {
	raw := "The/Scene:Name With Spaces"
	safe := sanitizeFileName(raw)
	expected := "The-Scene_Name_With_Spaces"
	if safe != expected {
		t.Errorf("Expected %q, got %q", expected, safe)
	}
}

func TestStartAndEndScene(t *testing.T) {

	config.SceneDir = "../../logs/scenes"

	r := &Room{
		Name: "TestRoom",
	}
	r.StartScene("Test Scene", "Mysterious", "DM")
	if r.ActiveScene == nil {
		t.Fatal("Expected scene to be initialized")
	}
	if r.ActiveScene.Title != "Test Scene" {
		t.Errorf("Unexpected scene title: %s", r.ActiveScene.Title)
	}
	r.ActiveScene.LogEntry("It begins.")

	// Simulate EndScene and log writing
	summary := r.EndScene("bob")
	fmt.Println(summary)
	// Check if file was created
	files, err := os.ReadDir(config.SceneDir)
	if err != nil {
		t.Fatalf("Could not read logs directory: %v", err)
	}
	found := false
	for _, f := range files {
		if strings.Contains(f.Name(), "scene_") && strings.HasSuffix(f.Name(), ".txt") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected scene log file not found in logs/")
	}
}
