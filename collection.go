package main

import (
  "os"
  "encoding/json"
  "time"
  "errors"
  "strconv"
  "strings"
)

type Collection struct {
  Todos []*Todo
}

func (c *Collection) RemoveAtIndex(item  int) {
  s := *c
  s.Todos = append(s.Todos[:item], s.Todos[item+1:]...)
  *c = s
}

func (c *Collection) RetrieveTodos() (error) {
  file, err := os.OpenFile(os.Getenv("TODO_DB_PATH"), os.O_RDONLY, 0600)
  if err != nil {
    return err
  }

  defer file.Close()

  err = json.NewDecoder(file).Decode(&c.Todos)
  return err
}

func (c *Collection) WriteTodos() (error) {
  file, err := os.OpenFile(os.Getenv("TODO_DB_PATH"), os.O_RDWR|os.O_TRUNC, 0600)
  if err != nil {
    return err
  }

  defer file.Close()

  err = json.NewEncoder(file).Encode(&c.Todos)
  return err
}

func (c *Collection) ListPendingTodos() () {
  for i := len(c.Todos)-1; i >=0; i-- {
    if c.Todos[i].Status != "pending" {
      c.RemoveAtIndex(i)
    }
  }
}

func (c *Collection) ListDoneTodos() () {
  for i := len(c.Todos)-1; i >=0; i-- {
    if c.Todos[i].Status != "done" {
      c.RemoveAtIndex(i)
    }
  }
}

func (c *Collection) CreateTodo(newTodo *Todo) (error) {
  var highestId int64 = 0
  for _, todo := range c.Todos {
    if todo.Id > highestId {
      highestId = todo.Id
    }
  }

  newTodo.Id = (highestId + 1)
  c.Todos = append(c.Todos, newTodo)

  err := c.WriteTodos()
  return err
}

func (c *Collection) Find(id int64) (foundedTodo *Todo, err error) {
  founded := false
  for _, todo := range c.Todos {
    if id == todo.Id {
      foundedTodo = todo
      founded = true
    }
  }
  if !founded {
    err = errors.New("The todo with the id "+strconv.FormatInt(id, 10)+" was not found.")
  }
  return
}

func (c *Collection) Toggle(id int64) (*Todo, error) {
  todo, err := c.Find(id)

  if err != nil {
    return todo, err
  }

  if todo.Status == "done" {
    todo.Status = "pending"
  } else {
    todo.Status = "done"
  }
  todo.Modified = time.Now().Local().String()

  err = c.WriteTodos()
  if err != nil {
    err = errors.New("Todos couldn't be saved")
    return todo, err
  }

  return todo, err
}

func (c * Collection) RemoveFinishedTodos() (error) {
  c.ListPendingTodos()
  err := c.WriteTodos()
  return err
}

func (c *Collection) Reorder() (error) {
  for i, todo := range c.Todos {
    todo.Id = int64(i+1)
  }
  err := c.WriteTodos()
  return err
}

func (c *Collection) Search(sentence string) {
  for i := len(c.Todos)-1; i >=0; i-- {
    if !strings.Contains(c.Todos[i].Desc, sentence) {
      c.RemoveAtIndex(i)
    }
  }
}