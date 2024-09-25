package audio
import (
    "testing"

    "github.com/golang/mock/gomock"
)


func TestPlay(t *testing.T){

	instance := NewAudioPlayer()

	if instance.Done != nil {
		t.Error("done channel should be nil on initialization")
	}



}
