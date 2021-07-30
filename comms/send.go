package comms

import (
	"encoding/json"
	"fmt"
	"http"
	"strings"
	"sync"
	"time"
)

type displays struct {
	name    string `json: "name"`
	power   string `json: "power"`
	input   string `json: "input"`
	blanked bool   `json: "blanked"`
}

type AlertMessage struct {
	Message string `json: "message"`
}

func worker(wg *sync.WaitGroup, m []string, alert_via string, contenttype string) {
	defer wg.Done()
	var alertMessage AlertMessage

	for range time.Tick(time.Second * 10) {
		for _, part := range m {
			alertMessage.Message = m
			req := json.Marshal(alertMessage)
			resp, err := http.Post(alert_url, contenttype, req)
		}
	}
}

func SendMessage(m []string, via string) {
	var systems Systems
	var wg sync.WaitGroup

	fmt.Println("Sending Message Abroad...Ready.........")

	// break down where the device lives - Room
	split := strings.Split(via, "-")
	bldg := split[0]
	room_num := split[1]
	room := bldg + room_num
	cp := bldg + room_num + "CP1"

	// Build ye the status url
	status_url := "http://" + cp + ":8000/buildings/" + bldg + "/rooms/" + room_num

	// Build ye the via url
	alert_url := "http://" + cp + ":8058/" + via + "alert/message"

	// get current status of the room
	resp, err := http.Get(status_url)
	if err != nil {
		fmt.Printf("Error Getting Status: %v\n", err.Error())
	}

	defer resp.Body.Close()

	status, err := io.ReadAll(resp.Body)

	// Save that for later (It will become important)
	orig := status

	// Get all of the displays in the room and reconfigure the json and reassert
	err = json.Unmarshal(status, &systems)
	if err != nil {
		fmt.Printf("Error in unmarshalling json: %v\n", err.Error())
		return
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

	// Send Body to system
	req, err := http.Put(status_url, contenttype, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error in Posting Content")
	}

	// Send the alert messages to the VIA1
	for i := 0; i < 1; i++ {
		fmt.Println("Starting worker")
		wg.Add(1)
		go worker(&wg, m, alert_url, contenttype)
	}
	// Loop over and over for the next 5 minutes
	// Reset the room back to the original room status

}
