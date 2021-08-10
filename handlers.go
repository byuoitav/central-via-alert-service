package main

import (
	"context"
	"fmt"
	//"log"
	"net/http"
	"strconv"
	//"os"
	//"strconv"
	"time"

	comms "github.com/byuoitav/central-via-alert-service/comms"
	couch "github.com/byuoitav/central-via-alert-service/couch"
	message "github.com/byuoitav/central-via-alert-service/message"
	//viadriver "github.com/byuoitav/kramer-driver"
	"github.com/labstack/echo"
)

func test(via string, message []string) {
	fmt.Printf("VIA to Execute: %v\n", via)
	fmt.Printf("Message: %v\n", message)
}

type Handlers struct {
	CreateServer func(string) *AlertServer
}

func (h *Handlers) RegisterRoutes(e *echo.Group) {

	// Production Endpoint for sending messages to all devices
	e.POST("/emessage/timer/:timing/all", func(c echo.Context) error {
		//t := c.Param("timing")
		build := h.CreateServer("all")
		u := build.Username
		p := build.Password
		fmt.Printf("Username: %v\n", u)

		shortDuration := 10 * time.Second
		d := time.Now().Add(shortDuration)
		ctx, cancel := context.WithDeadline(context.Background(), d)
		defer cancel()

		// Create a new client for connecting to the production couch database
		client, err := couch.NewClient(ctx, u, p, "https://couchdb-prd.avs.byu.edu")
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// Get all of the first VIAs in each room
		devices, err := couch.Devices(ctx, client)
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}

		fmt.Printf("Devices: %v\n", devices)

		//alert_time, err := strconv.Atoi(t)

		// pull the message from the request
		messages := echo.Map{}

		err = c.Bind(&messages)
		if err != nil {
			fmt.Printf("No message received: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		alertmess := messages["Message"].(string)

		// Transform the text into an array of text strings and prep for sending to VIAs
		me := message.Transform(alertmess)

		for _, dev := range devices {
			go func(dev string, me []string) {
				test(dev, me)
			}(dev, me)
		}

		// Send the message to the specified VIA
		// Go Routine when sending to more than one device
		/*
			err = comms.SendMessage(me, via, alert_time)
			if err != nil {
				fmt.Printf("Error: %v\n", err.Error())
			}
		*/

		// Get all the VIA's in the database and dump to array

		// Interate oer list of VIAs and executing against each one

		// TODO Go routine for executing against a large list of sadness

		// TODO Return status

		return c.JSON(http.StatusOK, fmt.Sprintf("Work in progress"))

	})

	// Test Endpoint against larger test group
	e.POST("/emessage/test", func(c echo.Context) error {
		// Test group? ITB? JKB? TLRB? JFSB? Ye olden

		return c.JSON(http.StatusOK, fmt.Sprintf("Still implementing endpoint"))

	})

	// Endpoint for testing just against a single
	e.POST("/emessage/timer/:timing/via/:vianame", func(c echo.Context) error {
		t := c.Param("timing")
		via := c.Param("vianame")

		alert_time, err := strconv.Atoi(t)
		if err != nil {
			fmt.Errorf("Error Converting string to int")
		}

		alert := make(map[string]interface{})

		err = c.Bind(&alert)
		if err != nil {
			fmt.Printf("No message received: %s\n", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		alertmess := alert["Message"].(string)

		// Transform the text into an array of text strings and prep for sending to VIAs
		me := message.Transform(alertmess)

		// Send the message to the specified VIA
		// Go Routine when sending to more than one device
		err = comms.SendMessage(me, via, alert_time)
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
		}

		fmt.Printf("Single Endpoint Used: %v\n", via)
		return c.JSON(http.StatusOK, fmt.Sprintf("Message: %v\n", me))

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
