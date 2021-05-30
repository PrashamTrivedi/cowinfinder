package cowinfinder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type States struct {
	StatesData []statesData `json:"states,omitempty"`
}

type statesData struct {
	StateId    int    `json:"state_id,omitempty"`
	StateName  string `json:"state_name,omitempty"`
	StateName1 string `json:"state_name_1,omitempty"`
}

type Districts struct {
	Districts []districtsData `json:"districts,omitempty"`
}

type districtsData struct {
	StateId       int    `json:"state_id,omitempty"`
	DistrictId    int    `json:"district_id,omitempty"`
	DistrictName  string `json:"district_name,omitempty"`
	DistrictName1 string `json:"district_name_1,omitempty"`
}

type Centers struct {
	Centers []centersData `json:"centers,omitempty"`
}

type centersData struct {
	CenterId      int           `json:"center_id,omitempty"`
	Name          string        `json:"name,omitempty"`
	Name1         string        `json:"name_1,omitempty"`
	Address       string        `json:"address,omitempty"`
	Address1      string        `json:"address_1,omitempty"`
	StateName     string        `json:"state_name,omitempty"`
	StateName1    string        `json:"state_name_1,omitempty"`
	DistrictName  string        `json:"district_name,omitempty"`
	DistrictName1 string        `json:"district_name_1,omitempty"`
	BlockName     string        `json:"block_name,omitempty"`
	BlockName1    string        `json:"block_name_1,omitempty"`
	Pincode       int           `json:"pincode,omitempty"`
	Lat           float32       `json:"lat,omitempty"`
	Lon           float32       `json:"lon,omitempty"`
	From          string        `json:"from,omitempty"`
	To            string        `json:"to,omitempty"`
	FeeType       string        `json:"fee_type,omitempty"`
	VaccineFees   []vaccineFees `json:"vaccine_fees,omitempty"`
	Sessions      []sessions    `json:"sessions,omitempty"`
}

type vaccineFees struct {
	Vaccine string `json:"vaccine,omitempty"`
	Fee     string `json:"fee,omitempty"`
}

type sessions struct {
	SessionId              string   `json:"session_id,omitempty"`
	Date                   string   `json:"date,omitempty"`
	AvailableCapacity      int      `json:"available_capacity,omitempty"`
	AvailableCapacityDose1 int      `json:"available_capacity_dose1,omitempty"`
	AvailableCapacityDose2 int      `json:"available_capacity_dose2,omitempty"`
	MinAgeLimit            int      `json:"min_age_limit,omitempty"`
	Vaccine                string   `json:"vaccine,omitempty"`
	Slots                  []string `json:"slots,omitempty"`
}

func GetCalendarSlots(districtIds []string, age int, firstDose bool) []centersData {
	centers := make([]centersData, 0)
	for _, districtId := range districtIds {
		centerChannels := make(chan centersData)
		go getCalendarForDistrict(districtId, centerChannels, age, firstDose)

		for center := range centerChannels {
			if len(center.Sessions) > 0 {
				centers = append(centers, center)
			}
		}
	}
	return centers

}

func getCalendarForDistrict(districtId string, centerChannels chan centersData, age int, firstDose bool) {
	var centers Centers
	today := time.Now().Format("02-01-2006")
	url := fmt.Sprintf("https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/calendarByDistrict?district_id=%s&date=%s", districtId, today)

	if err := getData(url, &centers); err != nil {
		panic(err)
	}

	for _, center := range centers.Centers {

		sessions := make([]sessions, 0)
		//Find All sessions where
		//Min Age = given age, and dose1 or dose2 count > 0 based on given flag
		//Add all sessions
		for _, session := range center.Sessions {
			availableDose := session.AvailableCapacityDose2
			if firstDose {
				availableDose = session.AvailableCapacityDose1
			}
			if session.MinAgeLimit == age && session.AvailableCapacity > 10 && availableDose > 10 {

				sessions = append(sessions, session)
			}
		}

		center.Sessions = sessions
		centerChannels <- center

	}

	close(centerChannels)
}
func GetStates() ([]statesData, error) {
	var states States

	url := "https://cdn-api.co-vin.in/api/v2/admin/location/states"
	// fmt.Println("Getting States")
	if err := getData(url, &states); err != nil {
		return nil, err
	}
	return states.StatesData, nil

	// for _, storyId := range storyIds[:10] {
	// 	storiesChannel <- storyId
	// }
	// close(storiesChannel)
}

func getData(url string, typeData interface{}) error {
	// fmt.Println(url)

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if err := json.NewDecoder(response.Body).Decode(&typeData); err != nil {
		return err
	}

	return nil

}
