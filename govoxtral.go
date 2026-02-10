package voxtral

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	CppLoadModel  func(modelDir string) int
	CppTranscribe func(wavPath string) string
	CppFreeResult func()
)

type Voxtral struct{}

func (v *Voxtral) Load(modelFile string) error {
	if ret := CppLoadModel(modelFile); ret != 0 {
		return fmt.Errorf("failed to load Voxtral model from %s", modelFile)
	}
	return nil
}

func (v *Voxtral) Transcribe(audioPath string) (string, error) {
	// Temporary directory for the converted wav
	dir, err := os.MkdirTemp("", "voxtral")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)

	convertedPath := dir + "/converted.wav"

	// Convert audio to 16kHz mono wav using ffmpeg
	// -ar 16000: set audio sample rate to 16000Hz
	// -ac 1: set number of audio channels to 1 (mono)
	// -y: overwrite output files without asking
	cmd := exec.Command("ffmpeg", "-i", audioPath, "-ar", "16000", "-ac", "1", "-y", convertedPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("ffmpeg conversion failed: %w, output: %s", err, string(out))
	}

	result := strings.Clone(CppTranscribe(convertedPath))
	CppFreeResult()

	return strings.TrimSpace(result), nil
}
