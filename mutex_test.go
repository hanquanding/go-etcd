package getcd

import (
	"context"
	"github.com/hanquanding/getcd"
	"testing"
)

func TestMutext(t *testing.T) {
	etcd, err := getcd.New("http://localhost:2379")
	if err != nil {
		t.Fatal(err)
	}

	mutext, err := etcd.NewMutex("test")
	if err != nil {
		t.Fatal(err)
	}
	err = mutext.Lock(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	t.Log("success")
}
