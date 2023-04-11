package server

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"sync"

	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"time"

	"github.com/PuerkitoBio/goquery"
	"whs.su/rusprofile/src/rpc"
)

var IsDigit = regexp.MustCompile(`^[0-9]+$`).MatchString

type HttpFetcher interface {
	Get(ctx context.Context, url string, headers map[string]string) ([]byte,error)
}

type defaultHttpFetcher struct {

}

func DefaultHttpFetcher() HttpFetcher { return &defaultHttpFetcher{} }

func (this *defaultHttpFetcher) Get(ctx context.Context, url string, headers map[string]string) ([]byte,error) {
	log.Printf("fetch url '%s'", url)
	client := http.Client{Timeout: 25 * time.Second}
	if req, err := http.NewRequest(http.MethodGet, url, nil); err != nil {
		return nil, err
	} else {

		for k, v := range headers {
			req.Header.Set(k, v)
		}

		if resp, err := client.Do(req.WithContext(ctx)); err != nil {
			return nil, err
		} else {
			if resp.Body == nil {
				return nil, fmt.Errorf("empty response")
			}
			defer resp.Body.Close()

			if body, err := ioutil.ReadAll(resp.Body); err != nil {
				return nil, err
			} else {
				return body, nil
			}
		}
	}
}

type Server struct {
	Fetcher HttpFetcher
	guard sync.Mutex
	rpc.UnimplementedRusprofileServer
}

type RusprofileDetailsResponse struct {
	Kpp string `json:"kpp"`
	Inn string `json:"inn"`
}

func NewRusprofileClientResponseFromHtml(body []byte) (*RusprofileDetailsResponse, error) {
	if doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body)); err != nil {
		return nil, err
	} else {
		kpps := doc.Selection.Find("span#clip_kpp").First()
		if kpps.Length() < 1 {
			return nil, fmt.Errorf("kpp node not found in body")
		}
		var result = &RusprofileDetailsResponse{Kpp: kpps.Text()}
		return result, nil
	}
}

type RusprofileResponseEntry struct {
	Name    string `json:"name"`
	RawName string `json:"raw_name"`
	Ogrn    string `json:"ogrn"`
	Ceo     string `json:"ceo_name"`
	Url     string `json:"url"`
}

type RusprofileResponse struct {
	Count int                       `json:"ul_count"`
	Items []RusprofileResponseEntry `json:"ul"`
}

func (this *Server) httpGet(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	return this.Fetcher.Get(ctx,url,headers)
}

func (this *Server) Get(ctx context.Context, inn *rpc.InnRequest) (*rpc.ProfileResponse, error) {
	this.guard.Lock()
	defer func() {
		go func() { // release guard after 5 seconds - prevent too often requests to rusprofle.ru
			time.Sleep(5 * time.Second)
			this.guard.Unlock()
		}()
	}()
	// call ajax to rusprofile and return results
	if !IsDigit(inn.GetINN()) {
		return nil, fmt.Errorf("invalid inn format")
	}
	cacheKey := fmt.Sprintf("0.%d", time.Now().UnixMilli())
	url := fmt.Sprintf("https://www.rusprofile.ru/ajax.php?query=%s&action=search&cacheKey=%s", inn.GetINN(), cacheKey)
	if data, err := this.httpGet(ctx, url, makeHeaders()); err != nil {
		return nil, err
	} else {
		var shortInfo RusprofileResponse
		if err := json.Unmarshal(data, &shortInfo); err != nil {
			return nil, err
		} else {
			if shortInfo.Count < 1 || len(shortInfo.Items) < 1 {
				log.Printf("empty response for %s", inn.GetINN())
				return nil, fmt.Errorf("empty response")
			}
			entry := shortInfo.Items[0]
			url = fmt.Sprintf("https://www.rusprofile.ru%s", entry.Url)

			// add random delay to omit ban from antibot system
			delay := 5 + int(math.Round(rand.Float64()*3.0))
			log.Printf("sleep %d seconds before next request", delay)
			time.Sleep(time.Duration(delay) * time.Second)
			if data, err := this.httpGet(ctx, url, makeHeaders()); err != nil {
				return nil, fmt.Errorf("empty response for details: %s", err.Error())
			} else {
				if details, err := NewRusprofileClientResponseFromHtml(data); err != nil {
					log.Printf("ERROR for '%s' : %s ", inn.INN, err.Error())
					return nil, fmt.Errorf("error reading kpp: %s", err.Error())
				} else {
					response := &rpc.ProfileResponse{}
					response.KPP = details.Kpp
					response.INN = inn.INN
					response.Director = entry.Ceo
					response.Company = entry.Name
					return response, nil
				}
			}
		}
	}
}

func NewServer() rpc.RusprofileServer {
	return &Server{ Fetcher: DefaultHttpFetcher() }
}
