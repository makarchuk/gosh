package shell

import (
  "os"
  "fmt"
  "bufio"
  "syscall"
  "strings"
  "os/user"
  "os/exec"
  "strconv"
  "os/signal"
)

type Shell struct {
  Dir string
}

func InitShell() Shell {
  s := Shell{ "" }
  s.Dir = s.Pwd()
  return s
} 

func (s Shell) Invitation() string {
  // TODO: Build from format string. Makre it customizable
  return s.Username() + "@" + s.Hostname() + ":" + s.Dir + "$ "
}

func (s Shell) Hostname() string {
  host, err := os.Hostname()
  if err != nil {
    return "localhost"
  } else {
    return host
  }
}

func (s Shell) Pwd() string {
  dir, err := os.Getwd()
  if err != nil {
    return "/"
  }
  return dir
}

func (s Shell) user() (user.User, error) {
  user_id := strconv.Itoa(os.Getuid())
  current_user, err := user.LookupId(user_id) 
  if err != nil{
    return user.User{"", "", "", "", ""}, err
  } else {
    return *current_user, nil
  }
}

func (s Shell) Username() string {
  current_user, err := s.user()
  if err != nil {
    return "nobody"
  } else {
    return current_user.Username
  }
}

// TODO: Change string to array of strings and call it like a decent person
func (s Shell) RunBinary(binary string, args ...string) {
  cmd := exec.Command(binary, args...)
  stdout, _ := cmd.StdoutPipe()
  err := cmd.Start()
  if err != nil {
    fmt.Println(err)
  }
  output := bufio.NewReader(stdout)
  for true {
    chunk, err := output.ReadString('\n')
    fmt.Print(chunk)
    if err != nil {
      break
    }
  }
}

func (s Shell) HandleInput(inp string) {
  chunks := strings.Split(inp, " ")
  command := chunks[0]
  args := chunks[1:]
  s.RunBinary(command, args...)
}

func fixTerminal(){
  exec.Command("stty", "-F", "/dev/tty", "echo").Run()
}

func setRawTerminal(){
  //Change Shell to some strange state
  //https://stackoverflow.com/a/17278730/4091324
  exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
  exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
  c := make(chan os.Signal, 2)
  signal.Notify(c, os.Interrupt, syscall.SIGTERM)
  go func() {
    <-c
    fixTerminal()
    os.Exit(0)
  }()
}

func (s Shell) Run() {
  setRawTerminal()
  defer fixTerminal()
  var b []byte = make([]byte, 1) 
  buffer := ""
  fmt.Print(s.Invitation())  
  for {
    os.Stdin.Read(b)
    chr := b[0]
    if chr == 4 && buffer=="" { //Ctrl+D 
      break
    } else if chr == 127 { //backspace
      buffer = buffer[:len(buffer)-1]
    } else if chr == 10 { //enter
      buffer = strings.Trim(buffer, "\n\t ")
      s.HandleInput(buffer)
      fmt.Print(s.Invitation())
      buffer = ""
    } else {
      buffer = buffer + string(b)
    }
    fmt.Print(string(chr))
  }
  fmt.Println("Bye!")
}