package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Clip struct {
	ID string `json:"ID"`
	Title string `json:"Title"`
	AudioURL string `json:"AudioURL"`
	Slug string `json:"Slug"`
}

type APIResponse struct {
	Clips []Clip `json:"Clips"`
}

func main() {
	baseURL := "http://api.omny.fm/orgs/8f7208d2-db6e-4bfa-85b5-ad3d00776f1f/programs/5da87d46-6804-4f65-97b6-ad4b0113f567/clips"
	pageSize := 10
	cursor := 393

	for {
		apiURL := fmt.Sprintf("%s?cursor=%d&pageSize=%d", baseURL, cursor, pageSize)
		response, err := http.Get(apiURL)
		if err != nil {
			fmt.Println("Error fetching API: ", err)
			return
		}
		defer response.Body.Close()

		var apiResponse APIResponse
		err = json.NewDecoder(response.Body).Decode(&apiResponse)
		if err != nil {
			fmt.Println("Error decoding API response:", err)
			return
		}

		fmt.Println("Total clips in API response:", len(apiResponse.Clips))
		fmt.Println("Cursor:", cursor)
		if len(apiResponse.Clips) == 0 {
			fmt.Println("No more clips found. Exiting loop.")
			break
		}
		for _, clip := range apiResponse.Clips {
			if containsFamouseFreshFridays(clip.Title) {
				err = downloadFile(clip.AudioURL, clip.Slug+".mp3")
				if err != nil {
					fmt.Println("Error downloading audio:", err)
				} else {
					fmt.Println("Audio downloaded:", clip.Title)
				}
			}
		}
		cursor += 1
	}
}


func downloadFile(url string, outputPath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}

func containsFamouseFreshFridays(title string) bool {
	return strings.Contains(title, "#FAMOUSFRESHFRIDAYS") || strings.Contains(title, "#FamousFreshFridays")
}