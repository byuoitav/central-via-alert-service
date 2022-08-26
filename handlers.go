package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	comms "github.com/byuoitav/central-via-alert-service/comms"
	couch "github.com/byuoitav/central-via-alert-service/couch"
	message "github.com/byuoitav/central-via-alert-service/message"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

const CouchDB string = "https://couchdb-prd.avs.byu.edu"

type Handlers struct {
	CreateServer func(string) *AlertServer
}

func test(via string, message []string) {
	fmt.Printf("VIA to Execute: %v\n", via)
	fmt.Printf("Message: %v\n", message)
}

func AlertAsync(alertMessage string, deviceList []string, alertTime int, L *zap.SugaredLogger) error {
	for _, dev := range devicesList {
		err = comms.SendMessage(me, dev, alert_time, L)
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
		}
	}
}

func (h *Handlers) RegisterRoutes(e *echo.Group) {

	// Production Endpoint for sending messages to all devices
	e.POST("/emessage/timer/:timing/all", func(c echo.Context) error {
		var alertmess string

		t := c.Param("timing")
		build := h.CreateServer("all")
		u := build.Username
		p := build.Password
		L := build.Logger
		fmt.Printf("Username: %v\n", u)

		query := map[string]interface{}{
			"fields": []string{"_id"},
			"limit":  2000,
			"selector": map[string]interface{}{
				"_id": map[string]interface{}{
					"$regex": "VIA1",
				},
			},
		}

		database := "devices"

		// Set Deadline for pulling data from the couch database
		shortDuration := 30 * time.Second
		d := time.Now().Add(shortDuration)
		ctx, cancel := context.WithDeadline(context.Background(), d)
		defer cancel()

		// Create a new client for connecting to the production couch database
		client, err := couch.NewClient(ctx, u, p, CouchDB)
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// Get all of the first VIAs in each room
		devices, err := couch.CouchQuery(ctx, client, query, database)
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}

		fmt.Printf("Devices: %v\n", devices)

		alert_time, err := strconv.Atoi(t)

		// pull the message from the request
		alert := echo.Map{}

		err = c.Bind(&alert)
		if err != nil {
			fmt.Printf("No message received: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		L.Debug("Alert Message: ", alert)

		// Checking the message to make sure the body is formated right so it doesn't cause issues down the line
		if _, ok := alert["Message"]; ok {
			L.Debug("Message has been received properly.....")
			alertmess = alert["Message"].(string)
		} else {
			L.Debug("Message key not formated properly or missing, Please use Message as the key in the body")
			mes := fmt.Sprintf("Message not formated properly or missing")
			return c.String(http.StatusInternalServerError, mes)
		}

		if alertmess == "" {
			L.Debug("Message is blank....")
			mes := fmt.Sprintf("Message is blank or missing")
			return c.String(http.StatusInternalServerError, mes)
		}

		// Transform the text into an array of text strings and prep for sending to VIAs
		me := message.Transform(alertmess)

		// Send the message to the specified VIA
		// Go Routine when sending to more than one device
		go func() {
			err = AlertAsync(me, dev, alert_time, L)
			if err != nil {
				fmt.Printf("Error: %v\n", err.Error())
			}
		}()

		return c.JSON(http.StatusOK, fmt.Sprintf("Successfully Sent Message to Campus"))

	})

	// Test Endpoint against larger test group
	e.POST("/emessage/timer/:timing/test", func(c echo.Context) error {
		// Test group? ITB? JKB? TLRB? JFSB? Ye olden
		t := c.Param("timing")
		var alertmess string
		build := h.CreateServer("all")
		u := build.Username
		p := build.Password
		L := build.Logger
		fmt.Printf("Username: %v\n", u)

		query := map[string]interface{}{
			"fields": []string{"_id"},
			"limit":  2000,
			"selector": map[string]interface{}{
				"_id": map[string]interface{}{
					"$regex": "ITB-1.*-VIA1",
				},
			},
		}

		database := "devices"

		shortDuration := 30 * time.Second
		d := time.Now().Add(shortDuration)
		ctx, cancel := context.WithDeadline(context.Background(), d)
		defer cancel()

		// Create a new client for connecting to the production couch database
		client, err := couch.NewClient(ctx, u, p, CouchDB)
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// Get all of the first VIAs in each room
		devices, err := couch.CouchQuery(ctx, client, query, database)
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}

		fmt.Printf("Devices: %v\n", devices)

		alert_time, err := strconv.Atoi(t)

		// pull the message from the request
		alert := echo.Map{}

		err = c.Bind(&alert)
		if err != nil {
			fmt.Printf("No message received: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		L.Debug("Alert Message: ", alert)

		// Checking the message to make sure the body is formated right so it doesn't cause issues down the line
		if _, ok := alert["Message"]; ok {
			L.Debug("Message has been received properly.....")
			alertmess = alert["Message"].(string)
		} else {
			L.Debug("Message key not formated properly or missing, Please use Message as the key in the body")
			mes := fmt.Sprintf("Message not formated properly or missing")
			return c.String(http.StatusInternalServerError, mes)
		}

		if alertmess == "" {
			L.Debug("Message is blank....")
			mes := fmt.Sprintf("Message is blank or missing")
			return c.String(http.StatusInternalServerError, mes)
		}

		// Transform the text into an array of text strings and prep for sending to VIAs
		me := message.Transform(alertmess)

		// Send the message to the specified VIA
		// Go Routine when sending to more than one device
		/*
			for _, dev := range devices {
				err = comms.SendMessage(me, dev, alert_time, L)
				if err != nil {
					fmt.Printf("Error: %v\n", err.Error())
				}
			}
		*/
		go func() {
			err = AlertAsync(me, dev, alert_time, L)
			if err != nil {
				fmt.Printf("Error: %v\n", err.Error())
			}
		}()

		response := fmt.Sprintf("Successfully sent message to test group in the ITB")
		return c.JSON(http.StatusOK, response)

	})

	// Endpoint for testing just against a single via
	e.POST("/emessage/timer/:timing/via/:vianame", func(c echo.Context) error {
		var alertmess string

		build := h.CreateServer("All")
		L := build.Logger
		t := c.Param("timing")
		via := c.Param("vianame")

		L.Info("Sending message to %v", via)

		alert_time, err := strconv.Atoi(t)
		if err != nil {
			L.Errorf("Error Converting string to int")
		}

		alert := make(map[string]interface{})

		err = c.Bind(&alert)
		if err != nil {
			L.Debugf("No message received: %s\n", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		L.Debug("Alert Message: ", alert)

		// Checking the message to make sure the body is formated right so it doesn't cause issues down the line
		if _, ok := alert["Message"]; ok {
			L.Debug("Message has been received properly.....")
			alertmess = alert["Message"].(string)
		} else {
			L.Debug("Message key not formated properly or missing, Please use Message as the key in the body")
			mes := fmt.Sprintf("Message not formated properly or missing")
			return c.String(http.StatusInternalServerError, mes)
		}

		if alertmess == "" {
			L.Debug("Message is blank....")
			mes := fmt.Sprintf("Message is blank or missing")
			return c.String(http.StatusInternalServerError, mes)
		}

		// Transform the text into an array of text strings and prep for sending to VIAs
		me := message.Transform(alertmess)

		// Send the message to the specified VIA
		// Go Routine when sending to more than one device
		err = comms.SendMessage(me, via, alert_time, L)
		if err != nil {
			L.Debug("Error: %v\n", err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}

		L.Debug("Single Endpoint Used: %v", via)
		response := fmt.Sprintf("Sending message to single endpoint: %v", via)
		return c.JSON(http.StatusOK, response)

	})

	// Get all the buildings in the database if needed
	e.GET("/emessage/buildings", func(c echo.Context) error {
		fmt.Println("Getting a list of all the buildings on campus")
		build := h.CreateServer("all")
		u := build.Username
		p := build.Password
		//L := build.Logger

		// couchdb query
		query := map[string]interface{}{
			"fields": []string{"_id"},
			"limit":  2000,
			"selector": map[string]interface{}{
				"_id": map[string]interface{}{
					"$regex": "",
				},
			},
		}

		// Which Database will you pull from
		database := "buildings"

		shortDuration := 10 * time.Second
		d := time.Now().Add(shortDuration)
		ctx, cancel := context.WithDeadline(context.Background(), d)
		defer cancel()

		client, err := couch.NewClient(ctx, u, p, CouchDB)
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// Get all the buildings in the couch database
		buildings, err := couch.CouchQuery(ctx, client, query, database)
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}

		for _, bldg := range buildings {
			fmt.Printf("Building: %v\n", bldg)
		}

		return c.JSON(http.StatusOK, buildings)
	})

	// Endpoint for executing against a single building
	e.POST("/emessage/timer/:timing/building/:bldg", func(c echo.Context) error {
		building := c.Param("bldg")
		var alertmess string
		t := c.Param("timing")
		build := h.CreateServer("all")
		u := build.Username
		p := build.Password
		L := build.Logger
		fmt.Printf("Username: %v\n", u)

		regParam := building + "-.*-VIA1"

		query := map[string]interface{}{
			"fields": []string{"_id"},
			"limit":  2000,
			"selector": map[string]interface{}{
				"_id": map[string]interface{}{
					"$regex": regParam,
				},
			},
		}

		database := "devices"

		shortDuration := 30 * time.Second
		d := time.Now().Add(shortDuration)
		ctx, cancel := context.WithDeadline(context.Background(), d)
		defer cancel()

		// Create a new client for connecting to the production couch database
		client, err := couch.NewClient(ctx, u, p, CouchDB)
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}

		// Get all of the first VIAs in each room
		devices, err := couch.CouchQuery(ctx, client, query, database)
		if err != nil {
			fmt.Printf("Error: %v\n", err.Error())
			return c.String(http.StatusInternalServerError, err.Error())
		}

		fmt.Printf("Devices: %v\n", devices)

		alert_time, err := strconv.Atoi(t)

		// pull the message from the request
		alert := echo.Map{}

		err = c.Bind(&alert)
		if err != nil {
			fmt.Printf("No message received: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		L.Debug("Alert Message: ", alert)

		// Checking the message to make sure the body is formated right so it doesn't cause issues down the line
		if _, ok := alert["Message"]; ok {
			L.Debug("Message has been received properly.....")
			alertmess = alert["Message"].(string)
		} else {
			L.Debug("Message key not formated properly or missing, Please use Message as the key in the body")
			mes := fmt.Sprintf("Message not formated properly or missing")
			return c.String(http.StatusInternalServerError, mes)
		}

		if alertmess == "" {
			L.Debug("Message is blank....")
			mes := fmt.Sprintf("Message is blank or missing")
			return c.String(http.StatusInternalServerError, mes)
		}

		// Transform the text into an array of text strings and prep for sending to VIAs
		me := message.Transform(alertmess)

		// Send the message to the specified VIA
		// Go Routine when sending to more than one device
		/*
			for _, dev := range devices {
				err = comms.SendMessage(me, dev, alert_time, L)
				if err != nil {
					fmt.Printf("Error: %v\n", err.Error())
				}
			}
		*/
		go func() {
			err = AlertAsync(me, dev, alert_time, L)
			if err != nil {
				fmt.Printf("Error: %v\n", err.Error())
			}
		}()

		response := fmt.Sprintf("Sending message to building: %v", building)
		return c.JSON(http.StatusOK, response)
	})
}
