package shell

import (
  "io"
  "os"
  "unsafe"
  "syscall"
 )

//Dark magic from pkg/term here
type Term struct {
  name string
  fd   int
  orig syscall.Termios // original state of the terminal, see Open and Restore
}

const (
  TCIFLUSH  = 0
  TCOFLUSH  = 1
  TCIOFLUSH = 2

  TCSANOW   = 0
  TCSADRAIN = 1
  TCSAFLUSH = 2
  
  TCSETS  = 0x5402
  TCSETSW = 0x5403
  TCSETSF = 0x5404
  TCFLSH  = 0x540B
  TCSBRK  = 0x5409
  TCSBRKP = 0x5425  
  IXON    = 0x00000400
  IXANY   = 0x00000800
  IXOFF   = 0x00001000
  CRTSCTS = 0x80000000
)


func ioctl(fd, request, argp uintptr) error {
  if _, _, e := syscall.Syscall6(syscall.SYS_IOCTL, fd, request, argp, 0, 0, 0); e != 0 {
    return e
  }
  return nil
}

func tcgetattr(fd uintptr, argp *syscall.Termios) error {
  return ioctl(fd, syscall.TCGETS, uintptr(unsafe.Pointer(argp)))
}

func tcsetattr(fd, action uintptr, argp *syscall.Termios) error {
  var request uintptr
  switch action {
  case TCSANOW:
    request = TCSETS
  case TCSADRAIN:
    request = TCSETSW
  case TCSAFLUSH:
    request = TCSETSF
  default:
    return syscall.EINVAL
  }
  return ioctl(fd, request, uintptr(unsafe.Pointer(argp)))
}

func cfmakeraw(attr *syscall.Termios) {
  attr.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
  attr.Oflag &^= syscall.OPOST
  attr.Cflag &^= syscall.CSIZE | syscall.PARENB
  attr.Cflag |= syscall.CS8
  attr.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
  attr.Cc[syscall.VMIN] = 1
  attr.Cc[syscall.VTIME] = 0
}

func (t Term) makeRaw() error {
  var a syscall.Termios
  if err := tcgetattr(uintptr(t.fd), &a); err != nil {
    return err
  }
  cfmakeraw((*syscall.Termios)(&a))
  return tcsetattr(uintptr(t.fd), TCSANOW, &a)
}


func open(name string) (*Term, error) {
  fd, e := syscall.Open(name, syscall.O_NOCTTY|syscall.O_CLOEXEC|syscall.O_NDELAY|syscall.O_RDWR, 0666)
  if e != nil {
    return nil, &os.PathError{"open", name, e}
  }

  t := Term{name: name, fd: fd}
  if err := tcgetattr(uintptr(t.fd), &t.orig); err != nil {
    return nil, err
  }
  return &t, syscall.SetNonblock(t.fd, false)
}

func (t Term) Restore() error {
  return tcsetattr(uintptr(t.fd), TCIOFLUSH, &t.orig)
}

func (t *Term) Write(b []byte) (int, error) {
  n, e := syscall.Write(t.fd, b)
  if n < 0 {
    n = 0
  }
  if n != len(b) {
    return n, io.ErrShortWrite
  }
  if e != nil {
    return n, &os.PathError{"write", t.name, e}
  }
  return n, nil
}


func (t *Term) Read(b []byte) (int, error) {
  n, e := syscall.Read(t.fd, b)
  if n < 0 {
    n = 0
  }
  if n == 0 && len(b) > 0 && e == nil {
    return 0, io.EOF
  }
  if e != nil {
    return n, &os.PathError{"read", t.name, e}
  }
  return n, nil
}
