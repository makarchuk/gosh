package terminal

import (
  "os"
  "os/user"
  "strconv"
)

type Terminal struct {
}

func InitTerminal() Terminal {
  return Terminal{}
} 

func (t Terminal) Invitation() string {
  return t.Username() + ":" + t.Pwd() + "$ "
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
