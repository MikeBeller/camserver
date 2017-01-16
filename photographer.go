package main

import (
	"io/ioutil"
	"os/exec"
	"time"
)

var execCommand = exec.Command

type Photographer struct {
	requestChannel       chan int
	responseChannel      chan PhotoResponse
	minimumPhotoInterval time.Duration
	lastPhotoTime        time.Time
	lastPhotoData        []byte
}

type PhotoResponse struct {
	Error       error
	IsDuplicate bool
	TimeTaken   time.Time
	Data        []byte
}

const minSeconds = 10

func NewPhotographer() *Photographer {
	p := &Photographer{
		requestChannel:       make(chan int),
		responseChannel:      make(chan PhotoResponse),
		minimumPhotoInterval: minSeconds * time.Second,
		lastPhotoTime:        time.Now().Add(-minSeconds * time.Second),
		lastPhotoData:        nil,
	}
	go p.run()
	return p
}

func (p *Photographer) takePicture() ([]byte, error) {
	cmd := execCommand("raspistill", "-o", "-")
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(pipe)
	if err != nil {
		return nil, err
	}
	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (p *Photographer) run() {
	for {
		<-p.requestChannel
		if time.Since(p.lastPhotoTime) < p.minimumPhotoInterval {
			p.responseChannel <- PhotoResponse{
				Error:       nil,
				IsDuplicate: true,
				TimeTaken:   p.lastPhotoTime,
				Data:        p.lastPhotoData,
			}
			continue
		}

		// Else take a new one
		p.lastPhotoTime = time.Now()
		data, err := p.takePicture()
		if err != nil {
			p.responseChannel <- PhotoResponse{
				Error: err,
			}
			continue
		}

		p.lastPhotoData = data
		p.responseChannel <- PhotoResponse{
			Error:       nil,
			IsDuplicate: false,
			TimeTaken:   p.lastPhotoTime,
			Data:        p.lastPhotoData,
		}
	}
}

// Returns a "photo response" struct containing either the new
// photo, a duplicate of an older photo (if taken too soon after
// another) or an error indicating what went wrong.
func (p *Photographer) GetPic() PhotoResponse {
	p.requestChannel <- 1
	resp := <-p.responseChannel
	return resp
}
