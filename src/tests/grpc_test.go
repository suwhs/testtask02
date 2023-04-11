package main_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"whs.su/rusprofile/src/rest"
	"whs.su/rusprofile/src/rpc"
	"whs.su/rusprofile/src/server"
)

type MockHttpServer struct {

}

func (this *MockHttpServer) Get(ctx context.Context, url string, headers map[string]string) ([]byte,error) {
	log.Printf("mock url: %s", url)	
	if strings.Contains(url,"ajax.php") {
		return []byte(`{
			"ul_count": 1,
			"ul": [
			  {
				"name": "ООО \"Яндекс\"",
				"raw_name": "ООО \"Яндекс\"",
				"many_ceo": 1,
				"link": "/id/189505",
				"ogrn": "1027700229193",
				"raw_ogrn": "1027700229193",
				"inn": "!~~7736207543~~!",
				"region": "Москва",
				"address": "119021, город Москва, ул. Льва Толстого, д.16",
				"inactive": 0,
				"status_extended": 0,
				"ceo_name": "Савиновский Артем Геннадьевич",
				"ceo_type": "Генеральный директор",
				"snippet_string": "Савиновский Артем Геннадьевич",
				"snippet_type": "Генеральный директор",
				"status_code": null,
				"svprekrul_date": null,
				"main_okved_id": "62.01",
				"okved_descr": "Разработка компьютерного программного обеспечения",
				"authorized_capital": "16605000.0000",
				"reg_date": "2000-09-14",
				"okpo": null,
				"url": "/id/189505",
				"aci_id": "189505"
			  }
			],
			"ip_count": 0,
			"success": true,
			"code": 0,
			"message": "OK"
		  }
		  `), nil

	} else if strings.Contains(url,"/id/") {
		return []byte(`<html><head></head><body><span id="clip_kpp">770901001</span></body>`),nil
	} else {
		return nil, fmt.Errorf("unknown mock url scheme ")
	}
}

func TestGrpcConnection(t *testing.T) {

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	rpc.RegisterRusprofileServer(grpcServer, &server.Server{Fetcher: &MockHttpServer{}})

	if lis, err := net.Listen("tcp", "localhost:9999"); err != nil {
		t.Fatalf("could not listen port 9999: %s", err.Error())
	} else {
		go func() {
			grpcServer.Serve(lis)
		}()
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second)

		if conn, err := grpc.Dial("localhost:9999", grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
			t.Fatalf("could not connect to grpc server: %s", err.Error())
		} else {
			defer conn.Close()
			client := rpc.NewRusprofileClient(conn)
			if response, err := client.Get(context.TODO(), &rpc.InnRequest{INN: "7736207543"}); err != nil {
				t.Fatalf("fail grpc call: %s", err.Error())
			} else {
				log.Printf("grpc call result: %v", response)
			}
		}
	}()

	// _testRest(t)
	wg.Wait()
	grpcServer.GracefulStop()
}

func TestRestConnection(t *testing.T) {
	ctx := context.TODO()
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	rpc.RegisterRusprofileServer(grpcServer, &server.Server{Fetcher: &MockHttpServer{}})

	if lis, err := net.Listen("tcp", "localhost:9999"); err != nil {
		t.Fatalf("could not listen port 9999: %s", err.Error())
	} else {
		go func() {
			grpcServer.Serve(lis)
		}()
		go func() {
			rest.RunRestServer(ctx, "localhost:9999")
		}()
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(3 * time.Second)
		client := http.Client{Timeout: 25 * time.Second}
		if req, err := http.NewRequest(http.MethodGet, "http://localhost:8081/search/7736207543", nil); err != nil {
			t.Fatalf("error connecting to rest service: %s", err.Error())
		} else {
			if resp, err := client.Do(req.WithContext(ctx)); err != nil {
				t.Fatalf("error executing request to rest server: %s",err.Error())
			} else {
				if resp.Body == nil {
					t.Fatalf("empty response")
				}
				defer resp.Body.Close()

				if body, err := ioutil.ReadAll(resp.Body); err != nil {
					t.Fatalf("error read response: %s", err.Error())
				} else {
					log.Printf("body: %s", string(body))
				}
			}
		}

	}()

	wg.Wait()
	grpcServer.GracefulStop()
}
