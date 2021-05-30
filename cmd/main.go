package main

import (
	"flag"
	"fmt"
	"prashamhtrivedi/cowinfinder"
	"strings"
	"time"
)

func main() {
	fmt.Println("Welcome to Cowin-finder.")
	districtIds := flag.String("districtIds", "154,770,153,772", "District Ids to Search, Defaults to Ahmedabad and Gandhinagar Area")
	minAge := flag.Int("minAge", 18, "Min age for Vaccine, Accepted values are 18 or 45, If you enter any other value, it will change to 18")
	firstDose := flag.Bool("firstDose", true, "If you're taking first dose or second?")
	flag.Parse()

	districtIdArr := strings.Split(*districtIds, ",")
	if *minAge != 45 {
		*minAge = 18
	}
	getDistrictData(districtIdArr, minAge, firstDose)
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()
	done := make(chan bool)

	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case <-ticker.C:
			getDistrictData(districtIdArr, minAge, firstDose)
		}
	}

}

func getDistrictData(districtIdArr []string, minAge *int, firstDose *bool) {
	availableSlots := cowinfinder.GetCalendarSlots(districtIdArr, *minAge, *firstDose)
	// fmt.Println(availableSlots)
	if len(availableSlots) > 0 {

		for _, availableSlot := range availableSlots {
			fmt.Println("==================================")
			fmt.Printf("Slots Available at %s (%s)\n", availableSlot.Name, availableSlot.Address)
			fmt.Printf("Time From:%s To: %s\n", availableSlot.From, availableSlot.To)

			fmt.Println("Slots:")
			vaccineFeeData := make(map[string]string)
			for _, vaccine := range availableSlot.VaccineFees {
				vaccineFeeData[vaccine.Vaccine] = vaccine.Fee
			}
			for _, session := range availableSlot.Sessions {
				fmt.Printf("Vaccine Available: %s\n", session.Vaccine)
				fmt.Printf("Vaccine Price: %s\n", vaccineFeeData[session.Vaccine])
				fmt.Printf("Total Doses Available: %d, For First Dose: %d, For Seond Dose: %d\n", session.AvailableCapacity, session.AvailableCapacityDose1, session.AvailableCapacityDose2)
			}
			fmt.Println("==================================")
		}
	} else {
		fmt.Println("No slots available now")
	}
}
