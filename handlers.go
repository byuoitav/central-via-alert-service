package main

import (
	//"context"
	"fmt"
	//"log"
	"net/http"
	//"os"
	//"strconv"
	//"time"

	message "github.com/byuoitav/central-via-alert-service/message"
	viadriver "github.com/byuoitav/kramer-driver"
	"github.com/labstack/echo"
)

type Handlers struct {
	CreateServer func(string) *viadriver.Via
}

func (h *Handlers) RegisterRoutes(e *echo.Group) {

	// Production Endpoint for sending messages to all devices
	e.POST("/emessage/all", func(c echo.Context) error {
		// pull the message from the request
		messages := echo.Map{}

		err := c.Bind(&messages)
		if err != nil {
			fmt.Printf("No message received: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		// TODO Transform message into an array of messages if needs be

		// TODO Get all the VIA's in the database and dump to array

		// TODO Interate oer list of VIAs and executing against each one

		// TODO Go routine for executing against a large list of sadness

		// TODO Return status

		return c.JSON(http.StatusOK, fmt.Sprintf("Still implementing endpoint"))

	})

	// Test Endpoint against larger test group
	e.POST("/emessage/test", func(c echo.Context) error {
		return c.JSON(http.StatusOK, fmt.Sprintf("Still implementing endpoint"))

	})

	// Endpoint for testing just against ITB-1106
	e.POST("/emessage/1106", func(c echo.Context) error {
		alert := make(map[string]interface{})

		// Largest size of a word that can be displayed before being broken into multiple words
		maxlength := 23

		// Largest size a message can be before being broken into multiple messages
		maxSize := 140

		err := c.Bind(&alert)
		if err != nil {
			fmt.Printf("No message received: %s\n", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		alertmess := alert["message"].(string)

		fmt.Printf("Received Message: %s\n", alertmess)

		// shorten any string down to below a character threshood.
		wordshorten := message.LongWords(alertmess, maxlength)

		// break longer messages down into smaller groups
		alerts := message.WordChunks(wordshorten, maxSize)
		fmt.Printf("Message: %v\n", alerts)
		// Send the message to ITB-1106-GO1
		// Go Routine which every VIA will end up using

		fmt.Printf("1106 Endpoint Used")
		return c.JSON(http.StatusOK, fmt.Sprintf("Message: %v\n", alerts))

	})

	// Get all the Buildings in the database
	e.GET("/emessage/buildings", func(c echo.Context) error {
		fmt.Printf("Getting a list of all the buildings on campus")
		return c.JSON(http.StatusOK, fmt.Sprintf("Still implementing endpoint"))

	})

	// Endpoint for executing against a single building
	e.POST("/emessage/bldg/:bldg", func(c echo.Context) error {
		fmt.Printf("Blinded by the torch light!")
		return c.JSON(http.StatusOK, fmt.Sprintf("Blinded by the torch light!!!!!!"))
	})
}
