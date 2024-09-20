package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const endpointModifyGroupedLight = "/clip/v2/resource/grouped_light/893be5a6-4170-45a1-9871-66096c07ca8f"
const hueApplicationKey = "3PlEPx8u3Jr5pcVn8gNXocMaZo7p9HFDbDDTD40G"
const hueBrigeURL = "https://192.168.1.171"

type UpdateRequest struct {
	Button int `json:"button"`
}

type Scene struct {
	Brightness float32
	X          float32
	Y          float32
}

const buttonValueOff = 0

type PowerPayload struct {
	On struct {
		On bool `json:"on"`
	} `json:"on"`
}

type DimmingColorPayload struct {
	On struct {
		On bool `json:"on"`
	} `json:"on"`
	Dimming struct {
		Brightness float32 `json:"brightness"`
	} `json:"dimming"`
	Color struct {
		XY struct {
			X float32 `json:"x"`
			Y float32 `json:"y"`
		} `json:"xy"`
	} `json:"color"`
}

var scenes = [3]Scene{
	// bright
	{
		X:          0.4572,
		Y:          0.4099,
		Brightness: 100.0,
	},
	// relaxed
	{
		X:          0.4998,
		Y:          0.4152,
		Brightness: 56.0,
	},
	// red
	{
		Brightness: 20.0,
		X:          0.6687,
		Y:          0.3189,
	},
}

func main() {

	http.HandleFunc("/update", handleUpdate)

	fmt.Println("Server is running")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit")

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var updateReq UpdateRequest
	if err := json.Unmarshal(reqBody, &updateReq); err != nil {
		http.Error(w, "Error parsing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if updateReq.Button >= 4 || updateReq.Button < 0 {
		http.Error(w, "Error. Received button index out of bounds", http.StatusBadRequest)
		return
	}

	var payload interface{}
	if updateReq.Button == buttonValueOff {
		payload = createPowerPayload(false)
	} else {
		scene := scenes[updateReq.Button-1]
		payload = createDimmingColorPayload(scene.X, scene.Y, scene.Brightness)
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Error building JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(payload)
	fmt.Printf("%s\n", jsonData)

	fullUrl, err := url.JoinPath(hueBrigeURL, endpointModifyGroupedLight)
	if err != nil {
		http.Error(w, "Error creating request: "+err.Error(), http.StatusInternalServerError)
	}
	req, err := http.NewRequest(http.MethodPut, fullUrl, bytes.NewBuffer(jsonData))
	req.Header.Add("hue-application-key", hueApplicationKey)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		http.Error(w, "Error creating request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := buildClient()
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error making request: "+err.Error(), http.StatusInternalServerError)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(string(body))
	w.WriteHeader(resp.StatusCode)
	w.Write(nil)

}

func buildClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return client
}

func createPowerPayload(isOn bool) PowerPayload {
	return PowerPayload{
		On: struct {
			On bool `json:"on"`
		}{On: isOn},
	}
}

func createDimmingColorPayload(x, y, brightness float32) DimmingColorPayload {
	return DimmingColorPayload{
		On: struct {
			On bool `json:"on"`
		}{On: true},
		Dimming: struct {
			Brightness float32 `json:"brightness"`
		}{
			Brightness: brightness,
		},
		Color: struct {
			XY struct {
				X float32 `json:"x"`
				Y float32 `json:"y"`
			} `json:"xy"`
		}{
			XY: struct {
				X float32 `json:"x"`
				Y float32 `json:"y"`
			}{
				X: x,
				Y: y,
			},
		},
	}
}
