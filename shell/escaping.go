package shell
import "fmt"

const (
  KEY_UP=65
  KEY_DOWN=66
  KEY_RIGHT=67
  KEY_LEFT=68
)

func (s Shell) MoveCursor(x int, y int) {
  fmt.Print("\033[%d;%dH", x, y)
}

// Move cursor up relative the current position
func (s Shell) MoveCursorUp(num int) {
  fmt.Print("\033[%dA", num);
}

// Move cursor down relative the current position
func (s Shell) MoveCursorDown(num int) {
  fmt.Print("\033[%dB", num);
}

// Move cursor forward relative the current position
func (s Shell) MoveCursorRight(num int) {
  fmt.Print("\033[%dC", num);
}

// Move cursor backward relative the current position
func (s Shell) MoveCursorLeft(num int) {
  fmt.Print("\033[%dD", num);
}

func (s Shell) ClearLine(){
  fmt.Print("\r\033[K")
}