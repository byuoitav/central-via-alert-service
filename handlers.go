package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/byuoitav/common/status"
	viadriver "github.com/byuoitav/kramer-driver"
	"github.com/labstack/echo"
)

type Handlers struct {
	CreateServer func(string) *viakramer.Via
}

func (h *Handlers) RegisterRoutes(e *echo.Group) {

	// Endpoint for sending messages to all devices
	e.POST("/emessage/all", func(c echo.Context) error {
		// pull the message from the request
		messages := echo.Map{}

		err := c.Bind(&messages)
		if err != nil {
			l.Printf("No message received: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}
		// Get all the VIA's in the database

		vs := h.CreateVideoSwitcher(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("Getting inputs")

		inputs, err := vs.AudioVideoInputs(c.Request().Context())
		if err != nil {
			l.Printf("unable to get inputs: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		out := c.Param("output")
		in, ok := inputs[out]
		if !ok {
			l.Printf("invalid output %q requested", out)
			return c.String(http.StatusBadRequest, "invalid output")
		}

		l.Printf("Got inputs: %+v", inputs)
		return c.JSON(http.StatusOK, status.Input{
			Input: fmt.Sprintf("%v:%v", in, out),
		})
	})
	// Endpoint for testing
	e.POST("/emessage/test", func(c echo.Context) error {
		addr := c.Param("address")
		vs := h.CreateVideoSwitcher(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("Getting volumes")

		vols, err := vs.Volumes(c.Request().Context(), []string{})
		if err != nil {
			l.Printf("unable to get volumes: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		block := c.Param("block")
		vol, ok := vols[block]
		if !ok {
			l.Printf("invalid block %q requested", block)
			return c.String(http.StatusBadRequest, "invalid block")
		}

		l.Printf("Got volumes: %+v", vols)
		return c.JSON(http.StatusOK, status.Volume{
			Volume: vol,
		})
	})
	// Endpoint for testing just against ITB-1106
	e.POST("/emessage/1106", func(c echo.Context) error {
		addr := c.Param("address")
		vs := h.CreateVideoSwitcher(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

		l.Printf("Getting mutes")

		mutes, err := vs.Mutes(c.Request().Context(), []string{})
		if err != nil {
			l.Printf("unable to get mutes: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		block := c.Param("block")
		mute, ok := mutes[block]
		if !ok {
			l.Printf("invalid block %q requested", block)
			return c.String(http.StatusBadRequest, "invalid block")
		}

		l.Printf("Got mutes: %+v", mutes)
		return c.JSON(http.StatusOK, status.Mute{
			Muted: mute,
		})
	})

	// Get all the Buildings in the database
	e.GET("/emessage/buildings", func(c echo.Context) error {
		addr := c.Param("address")
		vs := h.CreateVideoSwitcher(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)
		out := c.Param("output")
		in := c.Param("input")

		l.Printf("Setting AV input on %q to %q", out, in)

		err := vs.SetAudioVideoInput(c.Request().Context(), out, in)
		if err != nil {
			l.Printf("unable to set AV input: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Set AV input")
		return c.JSON(http.StatusOK, status.Input{
			Input: fmt.Sprintf("%v:%v", in, out),
		})
	})

	// Endpoint for executing against a single building
	e.POST("/emessage/bldg/:bldg", func(c echo.Context) error {
		addr := c.Param("address")
		vs := h.CreateVideoSwitcher(addr)
		l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)
		block := c.Param("block")

		vol, err := strconv.Atoi(c.Param("volume"))
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		l.Printf("Setting volume on %q to %d", block, vol)

		err = vs.SetVolume(c.Request().Context(), block, vol)
		if err != nil {
			l.Printf("unable to set volume: %s", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		l.Printf("Set volume")
		return c.JSON(http.StatusOK, status.Volume{
			Volume: vol,
		})
	})
	/*
		ps62.GET("/block/:block/muted/:mute", func(c echo.Context) error {
			addr := c.Param("address")
			vs := h.CreateVideoSwitcher(addr)
			l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)
			block := c.Param("block")

			mute, err := strconv.ParseBool(c.Param("mute"))
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			l.Printf("Setting mute on %q to %t", block, mute)

			err = vs.SetMute(c.Request().Context(), block, mute)
			if err != nil {
				l.Printf("unable to set mute: %s", err)
				return c.String(http.StatusInternalServerError, err.Error())
			}

			l.Printf("Set mute")
			return c.JSON(http.StatusOK, status.Mute{
				Muted: mute,
			})
		})

		gain60 := group.Group("/AT-GAIN-60/:address")

		// get state
		gain60.GET("/block/:block/volume", func(c echo.Context) error {
			addr := c.Param("address")
			amp := h.CreateAmp(addr)
			l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

			l.Printf("Getting volumes")

			ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
			defer cancel()

			vols, err := amp.Volumes(ctx, []string{})
			if err != nil {
				l.Printf("unable to get volumes: %s", err)
				return c.String(http.StatusInternalServerError, err.Error())
			}

			l.Printf("Got volumes: %+v", vols)

			block := c.Param("block")
			vol, ok := vols[block]
			if !ok {
				l.Printf("invalid block %q requested", block)
				return c.String(http.StatusBadRequest, "invalid block")
			}

			return c.JSON(http.StatusOK, status.Volume{
				Volume: vol,
			})
		})

		gain60.GET("/block/:block/muted", func(c echo.Context) error {
			addr := c.Param("address")
			amp := h.CreateAmp(addr)
			l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)

			l.Printf("Getting mutes")

			ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
			defer cancel()

			mutes, err := amp.Mutes(ctx, []string{})
			if err != nil {
				l.Printf("unable to get mutes: %s", err)
				return c.String(http.StatusInternalServerError, err.Error())
			}

			l.Printf("Got mutes: %+v", mutes)

			block := c.Param("block")
			mute, ok := mutes[block]
			if !ok {
				l.Printf("invalid block %q requested", block)
				return c.String(http.StatusBadRequest, "invalid block")
			}

			return c.JSON(http.StatusOK, status.Mute{
				Muted: mute,
			})
		})

		// set state
		gain60.GET("/block/:block/volume/:volume", func(c echo.Context) error {
			addr := c.Param("address")
			amp := h.CreateAmp(addr)
			l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)
			block := c.Param("block")

			vol, err := strconv.Atoi(c.Param("volume"))
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			l.Printf("Setting volume on %q to %d", block, vol)

			ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
			defer cancel()

			err = amp.SetVolume(ctx, block, vol)
			if err != nil {
				l.Printf("unable to set volume: %s", err)
				return c.String(http.StatusInternalServerError, err.Error())
			}

			l.Printf("Set volume")
			return c.JSON(http.StatusOK, status.Volume{
				Volume: vol,
			})
		})

		gain60.GET("/block/:block/muted/:mute", func(c echo.Context) error {
			addr := c.Param("address")
			amp := h.CreateAmp(addr)
			l := log.New(os.Stderr, fmt.Sprintf("[%v] ", addr), log.Ldate|log.Ltime|log.Lmicroseconds)
			block := c.Param("block")

			mute, err := strconv.ParseBool(c.Param("mute"))
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}

			l.Printf("Setting mute on %q to %t", block, mute)

			ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
			defer cancel()

			err = amp.SetMute(ctx, block, mute)
			if err != nil {
				l.Printf("unable to set mute: %s", err)
				return c.String(http.StatusInternalServerError, err.Error())
			}

			l.Printf("Set mute")
			return c.JSON(http.StatusOK, status.Mute{
				Muted: mute,
			})
		})
	*/
}
