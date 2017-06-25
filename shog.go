package main

import (
  "bufio"
  "fmt"
  "os"
  "io"
  "./terminal"
)

func main() {
  t := terminal.Terminal{}
  reader := bufio.NewReader(os.Stdin)
  for true {
    fmt.Print(t.Invitation())
    command, err := reader.ReadString('\n')
    if command == "exit\n" || err == io.EOF {
      break
    } else {
      fmt.Print(command)
    }
    
  }
  fmt.Println("Bye!")
}