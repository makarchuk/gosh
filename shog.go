package main

import (
  "./shell"
)

func main() {
  s := shell.InitShell()
  s.Read()
}