package main_test

import (
	"context"
	"log"
	"testing"

	"whs.su/rusprofile/src/rpc"
	rusprofile "whs.su/rusprofile/src/server"
)

func TestCall(t *testing.T) {
	srv := rusprofile.NewServer()
	if response, err := srv.Get(context.TODO(), &rpc.InnRequest{INN: "7736207543"}); err != nil {
		t.Fatalf("error: %v",err)
	} else {
		log.Printf("success response: %v", response)
	}
}

