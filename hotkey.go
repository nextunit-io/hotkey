package hotkey

import (
	"errors"
	"fmt"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// This constances containing the possible keys for building a combination of them.
const (
	ModAlt   = 1 << iota // ALT KEY
	ModCtrl              // CTRL KEY
	ModShift             // SHIFT KEY
	ModWin               // WINDOWS KEY
)

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel LogLevel = iota
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
)

// The LogLevel is a struct to provide you the possibility
// to change the loglevel which fits to your needs.
type LogLevel log.Level

type hotkeyObject struct {
	ID        int // Unique id
	Modifiers int // Mask of modifiers
	KeyCode   int // Key code, e.g. 'A'
}

// Hotkey is the interface definition. There are some needed functions, that the
// function need to provide.
type Hotkey interface {
	Register(fn func(id int)) error
	Deactivate() error
	String() string
}

// Create is a method which creates (whatever OS is in use) a hotkey to register.
func Create(id, modifieres, keyCode int) (Hotkey, error) {
	hotkey := hotkeyObject{
		ID:        id,
		Modifiers: modifieres,
		KeyCode:   keyCode,
	}

	log.Debugf("[Hotkey] - Create new Hotkey object. Scanning for System...")

	switch os := runtime.GOOS; os {
	case "darwin":
		return nil, errors.New("Unsupported OS")
	case "linux":
		return nil, errors.New("Unsupported OS")
	case "windows":
		log.Debugf("[Hotkey] - using Windows Hotkey.")

		return NewWindowsHotkey(hotkey), nil
	default:
		return nil, errors.New("Unsupported OS")
	}
}

// Close is the final function to tear down the hotkeys.
// This method should be called in defer of the main method.
func Close() {
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("Unsupported OS")
		break
	case "linux":
		fmt.Println("Unsupported OS")
		break
	case "windows":
		WindowsClose()
	default:
		fmt.Println("Unsupported OS")
		break
	}
}

// SetLogLevel provides a option to change the loglevel of
// the hotkey module.
func SetLogLevel(level LogLevel) {
	log.SetLevel(log.Level(level))
}
