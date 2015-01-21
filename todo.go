package main

import (
  "github.com/fatih/color"
  "fmt"
)

type Todo struct {
  Id int64 `json:"id"`
  Desc string `json:"desc"`
  Status string `json:"status"`
  Modified string `json:"modified"`
}

func (t *Todo) MakeOutput() {
  var symbole string
  var colorFunction func(...interface {}) string

  if(t.Status == "done"){
    colorFunction = color.New(color.FgGreen).SprintFunc()
    symbole = "✓"
  } else {
    colorFunction = color.New(color.FgRed).SprintFunc()
    symbole = "✕"
  }

  fmt.Println(t.Id, colorFunction(symbole), t.Desc)
}