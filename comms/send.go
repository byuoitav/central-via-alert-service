package comms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
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

/*
type displays struct {
	name    string `json: "name"`
	power   string `json: "power"`
	input   string `json: "input"`
	blanked bool   `json: "blanked"`
}
*/
type AlertMessage struct {
	Message string `json: "message"`
}

func worker(wg *sync.WaitGroup, m []string, alert_url string, contenttype string) {
	// After 10 minutes - Stop the function and exit
	timer := time.After(10 * time.Minute)
	fmt.Printf("Alert URL: %v\n", alert_url)
	defer wg.Done()
	var alertMessage AlertMessage

	for range time.Tick(time.Second * 5) {
		for {
			select {
			case <-timer:
				fmt.Printf("Worker has finished based on Timer")
				return
			default:
				for _, part := range m {
					fmt.Printf("Text: %v\n", part)
					alertMessage.Message = part
					req, err := json.Marshal(alertMessage)
					fmt.Println(string(req))
					if err != nil {
						fmt.Printf("JSON Marshal did not work")
					}
					//resp, err := http.Post(alert_url, contenttype, bytes.NewBuffer(req))
					reqType := "POST"
					resp, err := SendRequest(reqType, alert_url, req)
					if err != nil {
						fmt.Printf("Error: %v\n", err.Error())
					}
					s := string([]byte(resp))
					fmt.Printf("Worker Response: %v\n", s)
					time.Sleep(time.Second * 5)

				}
			}
			/*
				for range time.Tick(time.Second * 5) {
					for _, part := range m {
						fmt.Printf("Text: %v\n", part)
						alertMessage.Message = part
						req, _ := json.Marshal(alertMessage)
						//resp, err := http.Post(alert_url, contenttype, bytes.NewBuffer(req))
						reqType := "POST"
						resp, err := SendRequest(reqType, alert_url, req)
						if err != nil {
							fmt.Printf("Error: %v\n", err.Error())
						}
						s := string([]byte(resp))
						fmt.Printf("Worker Response: %v\n", s)
						time.Sleep(time.Second * 5)

			*/
		}
	}
}

func SendRequest(rtype string, url string, body []byte) ([]byte, error) {
	client := &http.Client{}

	fmt.Printf("I AM HERE!")

	req, err := http.NewRequest(rtype, url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("I am going to lose my job: %v\n", err.Error())

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

/*
func SendGet () {

}
*/
func SendMessage(m []string, via string) error {
	var systems Systems
	var wg sync.WaitGroup

	fmt.Println("Sending Message Abroad...Ready.........")

	// break down where the device lives - Room
	split := strings.Split(via, "-")
	bldg := split[0]
	room_num := split[1]
	//room := bldg + room_num
	cp := bldg + "-" + room_num + "-" + "CP1"

	// Build ye the status url
	status_url := "http://" + cp + ":8000/buildings/" + bldg + "/rooms/" + room_num

	// Build ye the via url
	alert_url := "http://" + cp + ":8058/" + via + "/alert/message"

	// get current status of the room
	resp, err := http.Get(status_url)
	if err != nil {
		fmt.Sprintf("Error Getting Status: %v\n", err.Error())
		return err
	}

	defer resp.Body.Close()

	status, err := io.ReadAll(resp.Body)

	// Save that for later (It will become important)
	orig := status

	// Get all of the displays in the room and reconfigure the json and reassert
	err = json.Unmarshal(status, &systems)
	if err != nil {
		fmt.Printf("Error in unmarshalling json: %v\n", err.Error())
		return err
	}

	// Find all the Displays in the room
	// Change them all to on and set them all to the VIA1
	fmt.Printf("Displays: %v\n", systems.Displays)
	fmt.Println("")
	for i, _ := range systems.Displays {
		systems.Displays[i].Input = "VIA1"
		systems.Displays[i].Power = "on"
		fmt.Printf("Display: %v\n", systems.Displays[i])
	}

	for i, _ := range systems.AudioDevices {
		fmt.Printf("Which one is running: %v and its at index: %v\n", systems.AudioDevices[i], i)
		re := regexp.MustCompile(`D+[0-9]+`)
		test := re.MatchString(systems.AudioDevices[i].Name)
		if test == true {
			systems.AudioDevices[i].Input = "VIA1"
			systems.AudioDevices[i].Power = "on"
			fmt.Printf("Display: %v\n", systems.AudioDevices[i])
		}
	}

	fmt.Printf("%v\n", systems)
	fmt.Println("")

	// build a new body to send that will turn on displays and set them to the VIA.
	body, err := json.Marshal(systems)
	if err != nil {
		fmt.Printf("Error in Marshal: %v\n", err.Error())
	}

	contenttype := "application/json"
	reqType := "PUT"

	// Send Body to system
	//req, err := http.Post(status_url, contenttype, bytes.NewBuffer(body))
	//req, err := SendPut(status_url, bytes.NewBuffer(body))
	sr, err := SendRequest(reqType, status_url, body)
	if err != nil {
		fmt.Printf("Error in Posting Content")
	}

	//defer req.Body.Close()
	//done, _ := io.ReadAll(req.Body)
	s := string([]byte(sr))
	fmt.Printf("Main body Response: %v\n", s)

	// Send the alert messages to the VIA1
	// Loop over and over for the next 5 minutes - logic in comms
	for i := 0; i < 1; i++ {
		fmt.Println("Starting worker")
		wg.Add(1)
		go worker(&wg, m, alert_url, contenttype)
	}

	// Reset the room back to the original room status
	final, err := http.Post(status_url, contenttype, bytes.NewBuffer(orig))
	if err != nil {
		fmt.Sprintf("Error Getting Status: %v\n", err.Error())
		return err
	}

	defer final.Body.Close()

	status, err = io.ReadAll(final.Body)

	fmt.Printf("Finishing Output: %v\n", status)

	return nil

}
