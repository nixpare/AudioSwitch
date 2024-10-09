package main

/*
#cgo LDFLAGS: -lole32 -loleaut32 -luuid -lmmdevapi
*/
import "C"
import (
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v3/pkg/application"
	"golang.org/x/sys/windows"
)

//go:embed frontend/dist
var assets embed.FS

var (
	app           *application.App
	audioService  *AudioService
	windowService *WindowService
)

var (
	saveDir            string
	windowSaveFilePath string
	audioSaveFilePath  string
)

var (
	ErrUserDataDir = errors.New("user data directory")
	ErrLogFile     = errors.New("log file creation")
)

func init() {
	err := initSaveDirAndPaths()
	if err != nil {
		log.Fatalln(err)
	}

	err = initLogs()
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	var exitCode int
	defer os.Exit(exitCode)

	log.Println("Starting AudioSwitch")
	defer log.Println("Stopping AudioSwitch")

	audioService = newAudioService()

	var err error
	windowService, err = newWindowService()
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		err := windowService.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	app = application.New(application.Options{
		Name:        "Audio Switch",
		Description: "A switch for toggling audio devices",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Services: []application.Service{
			application.NewService(audioService),
			application.NewService(windowService),
		},
		ErrorHandler: func(err error) {
			log.Println(err)
		},
	})

	err = audioService.Start()
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		err := audioService.Stop()
		if err != nil {
			log.Println(err)
		}
	}()

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Fake main thread after app.Run()
		
		windowService.CreateWindow()
		windowService.CreateOverlay()
	}()

	if err = app.Run(); err != nil {
		log.Printf("app fatal error: %v\n", err)
		exitCode = 1
	}
}

func initSaveDirAndPaths() error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	dir, err = buildSaveDirectory(dir, "Nixpare")
	if err != nil {
		return err
	}

	dir, err = buildSaveDirectory(dir, "AudioSwitch")
	if err != nil {
		return err
	}
	
	saveDir = dir
	audioSaveFilePath = filepath.Join(dir, "audio_save.json")
	windowSaveFilePath = filepath.Join(dir, "window_save.json")

	return nil
}

func buildSaveDirectory(base string, dir string) (string, error) {
	dir = filepath.Join(base, dir)
	info, err := os.Stat(dir)
	
	if err == nil && !info.IsDir() {
		return "", fmt.Errorf("%w <%s>: is not a directory", ErrUserDataDir, dir)
	}

	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("%w <%s>: %w", ErrUserDataDir, dir, err)
		}

		err = os.Mkdir(dir, 0660)
		if err != nil {
			return "", fmt.Errorf("%w <%s>: %w", ErrUserDataDir, dir, err)
		}
	}

	return dir, nil
}

func initLogs() error {
	if !ProductionBuild {
		return nil
	}

	f, err := os.OpenFile(filepath.Join(saveDir, "app.log"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrLogFile, err)
	}
	// avoid defer f.Close()

	os.Stdout, os.Stderr = f, f
	log.SetOutput(f)

	err = windows.SetStdHandle(windows.STD_OUTPUT_HANDLE, windows.Handle(f.Fd()))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrLogFile, err)
	}

	err = windows.SetStdHandle(windows.STD_ERROR_HANDLE, windows.Handle(f.Fd()))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrLogFile, err)
	}

	_, err = f.WriteString("\n------------------------------------------------\n\n")
	if err != nil {
		return fmt.Errorf("%w: %w", ErrLogFile, err)
	}

	return nil
}
