package shell

import (
  "os"
  "fmt"
  "bufio"
  "strings"
  "os/user"
  "os/exec"
  "strconv"
)

type Shell struct {
  Dir string
  term *Term
}

func (s Shell) Read() {
  for {
    ch := make([]byte, 3)
    s.term.Read(ch)
    fmt.Println(ch)
    s.bytesToKey(ch)

  }
}


func InitShell() Shell {
  term, _ := open("/dev/tty")
  s := Shell{ "", term}
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