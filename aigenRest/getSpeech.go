package aigenRest

import (
	"aigen/aigenAudioAutoPlay"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	audioPath = "voicenotes/"
	format    = ".mp3"
)

// SpeakOut Get speech from Azure
func SpeakOut(innerVoice string) (string, error) {
	// Set the input parameters
	innerVoiceLang := "en-US"
	innerVoiceName := "en-US-DavisNeural"
	format := ".mp3"

	// TODO: Support msspeech voices other than en-US-DavisNeural (e.g. en-US-JessaNeural) to match the inner voice sentinment (e.g. cheerful, sad, etc.)
	//nDefault\nChat\nAngry\nCheerful\nExcited\nFriendly\nHopeful\nSad\nShouting\nTerrified\nUnfriendly\nWhispering\
	//mood, _ := emotionalAI(innerVoice)
	//fmt.Println("mood: ", mood)
	//
	//speechStyle := "Default"
	//if mood == "Default" {
	//	speechStyle = "Default"
	//} else if mood == "Chat" {
	//	speechStyle = "Chat"
	//} else if mood == "Angry" {
	//	speechStyle = "Angry"
	//} else if mood == "Cheerful" {
	//	speechStyle = "Cheerful"
	//} else if mood == "Excited" {
	//	speechStyle = "Excited"
	//} else if mood == "Friendly" {
	//	speechStyle = "Friendly"
	//} else if mood == "Hopeful" {
	//	speechStyle = "Hopeful"
	//} else if mood == "Sad" {
	//	speechStyle = "Sad"
	//} else if mood == "Shouting" {
	//	speechStyle = "Shouting"
	//} else if mood == "Terrified" {
	//	speechStyle = "Terrified"
	//} else if mood == "Unfriendly" {
	//	speechStyle = "Unfriendly"
	//} else if mood == "Whispering" {
	//	speechStyle = "Whispering"
	//}

	// Send a request to get an authentication token
	tokenUrl := "https://eastus.api.cognitive.microsoft.com/sts/v1.0/issuetoken"
	tokenReq, issueTokenSuccess := http.NewRequest("POST", tokenUrl, nil)
	if issueTokenSuccess != nil {
		return "", issueTokenSuccess
	}
	tokenReq.Header.Set("Ocp-Apim-Subscription-Key", os.Getenv("SPEECH_KEY"))
	tokenResp, err := http.DefaultClient.Do(tokenReq)

	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		receiveSpeech := Body.Close()
		if receiveSpeech != nil {
			log.Printf("Could not use Subscription Token Or Something with error: %s", receiveSpeech)
		}
	}(tokenResp.Body)

	if tokenResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("SpeakOut: unexpected status code %d", tokenResp.StatusCode)
	}
	tokenBody, tokenSuccess := ioutil.ReadAll(tokenResp.Body)

	if tokenSuccess != nil {
		return "", tokenSuccess
	}
	token := string(tokenBody)

	// Send a request to generate the audio file
	url := "https://westus.tts.speech.microsoft.com/cognitiveservices/v1"
	xml := fmt.Sprintf("<speak version='1.0' xml:lang='%s'><voice xml:lang='%s' xml:gender='Male' name='%s'>%s</voice></speak>", innerVoiceLang, innerVoiceLang, innerVoiceName, innerVoice)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(xml))

	if err != nil {
		return "", err
	}

	req.Header.Set("X-Microsoft-OutputFormat", audio16khz128kbitratemonomp3)
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("User-Agent", "AiGen")
	resp, restSuccess := http.DefaultClient.Do(req)

	if restSuccess != nil {
		return "", restSuccess
	}
	defer func(Body io.ReadCloser) {
		speechOut := Body.Close()
		if speechOut != nil {
			log.Printf("Could not execute speech functionality %s", speechOut)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("SpeakOut: unexpected status code %d", resp.StatusCode)
	}
	body, audioSuccess := ioutil.ReadAll(resp.Body)
	if audioSuccess != nil {
		return "", audioSuccess
	}

	// Save the audio file to disk
	generateLetters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randomString := ""
	for i := 0; i < 10; i++ {
		randomString += string(generateLetters[rand.Intn(len(generateLetters))])
	}
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	timestamp = strings.ReplaceAll(timestamp, "-", "")
	randomString += timestamp
	err = ioutil.WriteFile(audioPath+randomString+format, body, 0644)
	if err != nil {
		return "", err
	}

	joinedFileName := joinFileName(audioPath, randomString, format)
	log.Printf("File saved to %s", joinedFileName)
	_, out := aigenAudioAutoPlay.UpdateBotChatAudioPath(joinedFileName)
	if out != nil {
		log.Printf("Error updating bot chat audio path: %v", err)
	}

	return joinedFileName, nil
}

// joinFileName joins the audio path, random string, and format
// to create the file name
func joinFileName(audioPath string, randomString string, format string) string {
	return audioPath + randomString + format
}

// GptSpeakOut Get Speech from OpenAI
func GptSpeakOut(innerVoice string) (string, error) {

	url := "https://api.openai.com/v1/audio/speech"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
    "model": "tts-1",
    "input": "%s",
    "voice": "shimmer"
  }`, innerVoice))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("OPENAI"))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(res.Body)

	body, audioSuccess := ioutil.ReadAll(res.Body)
	if audioSuccess != nil {
		return "", audioSuccess
	}

	// Save the audio file to disk
	generateLetters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randomString := ""
	for i := 0; i < 10; i++ {
		randomString += string(generateLetters[rand.Intn(len(generateLetters))])
	}
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	timestamp = strings.ReplaceAll(timestamp, "-", "")
	randomString += timestamp
	err = ioutil.WriteFile(audioPath+randomString+format, body, 0644)
	if err != nil {
		return "", err
	}

	joinedFileName := joinFileName(audioPath, randomString, format)
	log.Printf("File saved to %s", joinedFileName)
	_, out := aigenAudioAutoPlay.UpdateBotChatAudioPath(joinedFileName)
	if out != nil {
		log.Printf("Error updating bot chat audio path: %v", err)
	}

	return joinedFileName, nil
}
