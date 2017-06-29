package shell
import "fmt"

const (
  KEY_UP=65
  KEY_DOWN=66
  KEY_RIGHT=67
  KEY_LEFT=68
  KEY_BACKSPACE=8
  KEY_DELETE=127
)

func (s Shell) MoveCursor(x int, y int) {
  fmt.Print("\033[%d;%dH", x, y)
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

func (s Shell) bytesToKey(b []byte){
  if b[0] == 27 && b[1] == 91 {
    if b[2] == KEY_LEFT {
      s.MoveCursorLeft(1)
    } else if b[2] == KEY_RIGHT {
      s.MoveCursorRight(1)
    }
  } else if b[0] == 127 {
    //TODO: Acutally remove chars here
    s.MoveCursorLeft(1)
  } else {
    //fmt.Print(b[0])
  }
}
