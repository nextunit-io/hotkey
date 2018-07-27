package hotkey

import (
	"bytes"
	"errors"
	"fmt"
	"syscall"
	"time"
	"unsafe"

	log "github.com/sirupsen/logrus"
)

// The WindowsHotkey struct provides the methods
// for the Windows OS to register hotkeys and
// consume them.
type WindowsHotkey struct {
	hotkeyObject
	runBackgroundProcess bool
	callback             func(id int)
}

type windowsMsg struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}

var (
	running = map[int]*WindowsHotkey{}
)

// NewWindowsHotkey creates a new WindowsHotkey object.
func NewWindowsHotkey(hotkey hotkeyObject) *WindowsHotkey {
	return &WindowsHotkey{
		hotkeyObject:         hotkey,
		runBackgroundProcess: false,
		callback:             func(id int) {},
	}
}

// String returns a human-friendly display name of the hotkey
// such as "Hotkey[Id: 1, Alt+Ctrl+O]"
func (hotkey *WindowsHotkey) String() string {
	mod := &bytes.Buffer{}
	if hotkey.Modifiers&ModAlt != 0 {
		mod.WriteString("Alt+")
	}
	if hotkey.Modifiers&ModCtrl != 0 {
		mod.WriteString("Ctrl+")
	}
	if hotkey.Modifiers&ModShift != 0 {
		mod.WriteString("Shift+")
	}
	if hotkey.Modifiers&ModWin != 0 {
		mod.WriteString("Win+")
	}
	return fmt.Sprintf("Hotkey[ID: %d, %s%c]", hotkey.ID, mod, hotkey.KeyCode)
}

// Register is a function to register a Windows Hotkey.
func (hotkey *WindowsHotkey) Register(fn func(id int)) error {
	out := make(chan error)
	go func() {
		log.Debugf("[WindowsHotkey] - Register Windows Hotkey: %s", hotkey)

		user32 := syscall.MustLoadDLL("user32")
		defer user32.Release()

		reghotkey := user32.MustFindProc("RegisterHotKey")

		r1, _, err := reghotkey.Call(
			0, uintptr(hotkey.ID), uintptr(hotkey.Modifiers), uintptr(hotkey.KeyCode))
		if r1 == 1 {
			running[hotkey.ID] = hotkey
			out <- nil
			close(out)

			hotkey.runBackgroundProcess = true
			hotkey.callback = fn
			log.Debugf("[WindowsHotkey] - Start background process.")

			hotkey.windowsBackgroundProcess()
		} else {
			out <- fmt.Errorf("Failed to register %s, error: %s", hotkey, err)
			close(out)
		}
	}()

	return <-out
}

// Deactivate provides the finishing task to deactivate all registered
// hotkeys.
func (hotkey *WindowsHotkey) Deactivate() error {
	if _, ok := running[hotkey.ID]; !ok {
		return errors.New("No hotkey with that ID is running")
	}

	hotkey.stopHotkey()

	delete(running, hotkey.ID)
	return nil
}

// WindowsClose is a method which has to be called at the end, when the application
// is terminated.
func WindowsClose() {
	for _, hotkey := range running {
		hotkey.Deactivate()
	}
}

func (hotkey *WindowsHotkey) stopHotkey() {
	hotkey.runBackgroundProcess = false
}

func (hotkey *WindowsHotkey) windowsBackgroundProcess() {
	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

	peekmsg := user32.MustFindProc("PeekMessageW")

	for hotkey.runBackgroundProcess {
		var msg = &windowsMsg{}
		peekmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

		id := msg.WPARAM

		// Registered id is in the WPARAM field:
		if id != 0 {
			if _, ok := running[int(id)]; ok {
				log.Debugf("Starting callback for hotkey ID %d", int(id))
				running[int(id)].callback(int(id))
			} else {
				log.Debugf("Can find hotkey with ID %d", int(id))
			}
		}

		time.Sleep(time.Millisecond * 50)
	}
}
