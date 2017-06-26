package terminal

import (
  "io"
  "os"
  "fmt"
  "bufio"
  "strings"
  "os/user"
  "os/exec"
  "strconv"
)

type Terminal struct {
  Dir string
}

func InitTerminal() Terminal {
  t := Terminal{ "" }
  t.Dir = t.Pwd()
  return t
} 

func (t Terminal) Invitation() string {
  // TODO: Build from format string. Makre it customizable
  return t.Username() + "@" + t.Hostname() + ":" + t.Dir + "$ "
}

func (t Terminal) Hostname() string {
  host, err := os.Hostname()
  if err != nil {
    return "localhost"
  } else {
    return host
  }
}

func (t Terminal) Pwd() string {
  dir, err := os.Getwd()
  if err != nil {
    return "/"
  }
  return dir
}

func (t Terminal) user() (user.User, error) {
  user_id := strconv.Itoa(os.Getuid())
  current_user, err := user.LookupId(user_id) 
  if err != nil{
    return user.User{"", "", "", "", ""}, err
  } else {
    return *current_user, nil
  }
}

func (t Terminal) Username() string {
  current_user, err := t.user()
  if err != nil {
    return "nobody"
  } else {
    return current_user.Username
  }
}

// TODO: Change string to array of strings and call it like a decent person
func (t Terminal) RunBinary(binary string, args ...string) {
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

func (t Terminal) HandleInput(inp string) {
  chunks := strings.Split(inp, " ")
  command := chunks[0]
  args := chunks[1:]
  t.RunBinary(command, args...)
}


func (t Terminal) Run() {
  reader := bufio.NewReader(os.Stdin)
  for true {
    fmt.Print(t.Invitation())
    input, err := reader.ReadString('\n')
    if input == "exit\n" || err == io.EOF {
      break
    } else {
      input = strings.Trim(input, "\n ")
      t.HandleInput(input)
    }   
  }
  fmt.Println("Bye!")
}