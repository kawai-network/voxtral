package voxtral

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	sampleAudio = "https://qianwen-res.oss-cn-beijing.aliyuncs.com/Qwen3-ASR-Repo/asr_en.wav"
)

func skipIfNoModel(t *testing.T) string {
	t.Helper()
	modelDir := os.Getenv("VOXTRAL_MODEL_DIR")
	if modelDir == "" {
		t.Skip("VOXTRAL_MODEL_DIR not set, skipping test (set to voxtral model directory)")
	}
	// Check for a file that should exist in the model directory.
	// Adjust this filename if necessary based on what files really exist.
	if _, err := os.Stat(filepath.Join(modelDir, "consolidated.safetensors")); os.IsNotExist(err) {
		t.Skipf("Model file not found in %s, skipping", modelDir)
	}
	return modelDir
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP GET failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func TestLoadModel(t *testing.T) {
	modelDir := skipIfNoModel(t)

	// Initialize the library
	// Assumes the library (libgovoxtral.so/.dylib) is in the current directory or path
	if err := Init(""); err != nil {
		t.Fatalf("Failed to init library: %v", err)
	}

	v := &Voxtral{}
	err := v.Load(modelDir)
	if err != nil {
		t.Fatalf("LoadModel failed: %v", err)
	}
}

func TestAudioTranscription(t *testing.T) {
	modelDir := skipIfNoModel(t)

	tmpDir, err := os.MkdirTemp("", "voxtral-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Download sample audio
	audioFile := filepath.Join(tmpDir, "sample.wav")
	t.Log("Downloading sample audio...")
	if err := downloadFile(sampleAudio, audioFile); err != nil {
		t.Fatalf("Failed to download sample audio: %v", err)
	}

	// Initialize the library
	if err := Init(""); err != nil {
		t.Fatalf("Failed to init library: %v", err)
	}

	v := &Voxtral{}
	// Load model
	err = v.Load(modelDir)
	if err != nil {
		t.Fatalf("LoadModel failed: %v", err)
	}

	// Transcribe
	text, err := v.Transcribe(audioFile)
	if err != nil {
		t.Fatalf("Transcribe failed: %v", err)
	}

	t.Logf("Transcribed text: %s", text)

	if text == "" {
		t.Fatal("Transcription returned empty text")
	}

	allText := strings.ToLower(text)
	t.Logf("All text: %s", allText)

	if !strings.Contains(allText, "big") {
		t.Errorf("Expected 'big' in transcription, got: %s", allText)
	}
}
