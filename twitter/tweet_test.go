package twitter

import (
	"testing"
)

func TestGetCreationDate(t *testing.T) {
	tweet := &Tweet{CreatedAt: "Sun Jun 09 11:50:49 +0000 2019"}

	time, err := tweet.GetCreationDate()
	if err != nil {
		t.Errorf(err.Error())
	}
	if 49 != time.Second() {
		t.Errorf("Second must be 49. Got %d", time.Second())
	}
}
