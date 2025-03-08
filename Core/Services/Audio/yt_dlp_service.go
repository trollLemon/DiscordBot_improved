package audio



import "os/exec"


type YtDLP struct {}




func (ytdlp *YtDLP) GetAudioStream(url string) (string, error) {
	
	
	cmd := exec.Command("yt-dlp", "--get-url", "-x", url)
	
	output, err := cmd.Output()
	
	if err != nil {
		return "",err
	}

	streamUrl := string(output)

	return streamUrl, nil

}
