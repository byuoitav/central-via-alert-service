package comms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	//"go.uber.org/zap/zapcore"
)

type Systems struct {
	Displays []struct {
		Name    string `json: "name"`
		Power   string `json: "power"`
		Input   string `json: "input,omitempty"`
		Blanked bool   `json: "blanked"`
	} `json: "displays"`
	AudioDevices []struct {
		Name   string `json: "name"`
		Power  string `json: "power,omitempty"`
		Input  string `json: "input,omitempty"`
		Muted  bool   `json: "muted,omitempty"`
		Volume int    `json: "volume"`
	} `json: "audioDevices"`
}

type AlertMessage struct {
	Message string `json: "message"`
}

func worker(wg *sync.WaitGroup, m []string, status_url string, alert_url string, reset_url string, contenttype string, alert_time int, orig []uint8) {
	// After a determined about of time - Stop the function and exit
	timing := time.Duration(alert_time)
	timer := time.After(timing * time.Minute)
	fmt.Printf("Alert URL: %v\n", alert_url)
	defer wg.Done()
	var alertMessage AlertMessage

	for range time.Tick(time.Second * 10) {
		for {
			select {
			case <-timer:
				fmt.Printf("Worker has finished based on Timer")
				reqType := "PUT"

				// Reset the VIA to clear the alert
				_, err := http.Get(reset_url)
				if err != nil {
					fmt.Printf("Error sending reset command: %v\n", err.Error())
				}

				time.Sleep(10 * time.Second)

				// Return system back to original state
				// Reset the room back to the original room status
				final, err := SendRequest(reqType, status_url, orig)
				if err != nil {
					fmt.Printf("Error Getting Status: %v\n", err.Error())
				}
				f := string([]byte(final))
				fmt.Printf("Finishing Output: %v\n", f)

				return

			default:
				for _, part := range m {
					// build the alert message to send
					fmt.Printf("Text: %v\n", part)
					alertMessage.Message = part
					req, err := json.Marshal(alertMessage)
					fmt.Println(string(req))
					if err != nil {
						fmt.Printf("JSON Marshal did not work")
					}

					// Send Alert Message to the VIA
					reqType := "POST"
					resp, err := SendRequest(reqType, alert_url, req)
					if err != nil {
						fmt.Printf("Error sending alert to via: %v\n", err.Error())
					}
					s := string([]byte(resp))
					fmt.Printf("Worker Response: %v\n", s)
					time.Sleep(time.Second * 5)

				}
			}
		}
	}
}

func SendRequest(rtype string, url string, body []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(rtype, url, bytes.NewBuffer(body))
	if err != nil {
		errStr := fmt.Sprintf("Error creating request: %v\n", err.Error())
		return nil, errors.New(errStr)

	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error on Client: %v\n", err.Error())
		return nil, fmt.Errorf("Error Sending Request to System: %v\n", err.Error())
	}

	defer req.Body.Close()

	done, _ := io.ReadAll(resp.Body)
	s := string([]byte(done))
	fmt.Printf("Response: %v\n", s)

	return done, nil

}

func SendMessage(m []string, via string, alert_time int, L *zap.SugaredLogger) error {
	var systems Systems
	var wg sync.WaitGroup

	L.Infof("Sending Message Abroad...Ready.........")

	// break down where the device lives - Room
	split := strings.Split(via, "-")
	bldg := split[0]
	room_num := split[1]

	//room := bldg + room_num
	cp := bldg + "-" + room_num + "-" + "CP1"

	// fully qualify domain name for each VIA
	vn := via + ".byu.edu"

	// Build ye all of the urls
	// Status URL
	status_url := "http://" + cp + ":8000/buildings/" + bldg + "/rooms/" + room_num

	// Alert url
	alert_url := "http://" + cp + ":8058/" + vn + "/alert/message"

	// Reset url
	reset_url := "http://" + cp + ":8058/" + vn + "/reset"

	// get current status of the room
	resp, err := http.Get(status_url)
	if err != nil {
		L.Errorf("Error Getting Status: %v\n", err.Error())
		return err
	}

	defer resp.Body.Close()

	status, err := io.ReadAll(resp.Body)

	// Save that for later (It will become important)
	// We are going to use this to restore the existing state of the room
	orig := status

	// Get all of the displays in the room and reconfigure the json and reassert
	err = json.Unmarshal(status, &systems)
	if err != nil {
		errStr := fmt.Sprintf("Error in unmarshalling json: %v\n", err.Error())
		return errors.New(errStr)
	}

	// Find all the Displays in the room
	// Change them all to on and set them all to the VIA1
	L.Debugf("Displays: %v\n", systems.Displays)
	fmt.Println("")
	for i, _ := range systems.Displays {
		systems.Displays[i].Input = "VIA1"
		systems.Displays[i].Power = "on"
	}

	for i, _ := range systems.AudioDevices {
		re := regexp.MustCompile(`D+[0-9]+`)
		test := re.MatchString(systems.AudioDevices[i].Name)
		if test == true {
			systems.AudioDevices[i].Input = "VIA1"
			systems.AudioDevices[i].Power = "on"
		}
	}

	// build a new body to send that will turn on displays and set them to the VIA.
	body, err := json.Marshal(systems)
	if err != nil {
		errStr := fmt.Sprintf("Error in Marshal: %v\n", err.Error())
		return errors.New(errStr)
	}

	contenttype := "application/json"
	reqType := "PUT"

	// Send command to power devices and switch to VIA
	sr, err := SendRequest(reqType, status_url, body)
	if err != nil {
		errStr := fmt.Sprintf("Error in Posting Content")
		return errors.New(errStr)
	}

	s := string([]byte(sr))
	L.Debugf("Main body Response: %v\n", s)

	// Send the alert messages to the VIA1
	// Loop over and over for the specified time
	for i := 0; i < 1; i++ {
		L.Debugf("Starting worker")
		wg.Add(1)
		go worker(&wg, m, status_url, alert_url, reset_url, contenttype, alert_time, orig)
	}

	return nil

}
