package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"shuffler/app"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var h Handler

func main() {
	h = Handler{
		sessionMap: make(map[string]Session),
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/start", startGame)
	e.POST("/correct", correct)
	e.POST("/wrong", wrong)
	e.POST("/end", end)

	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}

type Session struct {
	Name string
	App app.App
	Input chan string
	Output chan app.Card
	Close chan bool
}

type Handler struct {
	sessionMap map[string]Session
}

func startGame(c echo.Context) error {
	session := c.Request().Header.Get("session")
	if session == "" {
		session = uuid.NewString()
	}

	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = session
	cookie.Expires = time.Now().Add(5 * 24 * time.Hour)
	c.SetCookie(cookie)

	newApp := app.NewApp()
	input, output, close := newApp.Run(true)
	h.sessionMap[session] = Session{
		Name: session,
		App: newApp,
		Input: input,
		Output: output,
		Close: close,
	}

	card := <- output

	return c.JSON(http.StatusOK, &card)
}

func getSession(c echo.Context) (Session, error) {
	session, err := c.Cookie("session")
	if err != nil {
		return Session{}, errors.New("unable to read session cookie")
	}

	if _, ok := h.sessionMap[session.Value]; !ok {
		return Session{}, errors.New("session not found")
	}

	return h.sessionMap[session.Value], nil
}

func correct(c echo.Context) error {
	session, err := getSession(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	session.Input <- "t"
	card := <- session.Output

	return c.JSON(http.StatusOK, &card)
}

func wrong(c echo.Context) error {
	session, err := getSession(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	session.Input <- "f"
	card := <- session.Output

	return c.JSON(http.StatusOK, &card)
}

func end(c echo.Context) error {
	session, err := getSession(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	session.Input <- "e"
	wrongCards := h.sessionMap[session.Name].App.GetWrongCards()
	
	delete(h.sessionMap, session.Name)
	
	return c.JSON(http.StatusOK, wrongCards)
}