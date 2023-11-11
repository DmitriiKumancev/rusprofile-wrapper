package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	serverAddress = "https://www.rusprofile.ru"
	searchQuery   = "ajax.php?query=%s&action=search"
)

type company struct {
	INN  string `json:"inn"`
	Name string `json:"raw_name"`
	CEO  string `json:"ceo_name"`
	URL  string `json:"url"`
}

type CompanyFinder struct {
	UnimplementedCompanyFinderServer
}

func (c CompanyFinder) ByINN(ctx context.Context, inn *INN) (*Company, error) {
	searchURL := fmt.Sprintf(fmt.Sprintf("%s/%s", serverAddress, searchQuery), inn.INN)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}(resp.Body)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read request body: %w", err)
	}

	var data struct {
		Items []company `json:"ul"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal to json: %w", err)
	}

	var comp *company
	for i := range data.Items {
		elem := &data.Items[i]
		elem.INN = strings.Trim(elem.INN, "!~")

		if elem.INN == inn.INN {
			comp = elem
			break
		}
	}

	if comp == nil {
		return nil, fmt.Errorf("no company with provided INN")
	}

	req, err = http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s", serverAddress, comp.URL), nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't create request: %w", err)
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}(resp.Body)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't create html document: %w", err)
	}

	selection := doc.Find("#clip_kpp")
	if len(selection.Nodes) == 0 {
		return nil, fmt.Errorf("no kpp value was found")
	}

	kpp := selection.Nodes[0].FirstChild.Data

	return &Company{
		INN:  comp.INN,
		KPP:  kpp,
		Name: comp.Name,
		CEO:  comp.CEO,
	}, nil
}

func (c CompanyFinder) mustEmbedUnimplementedCompanyFinderServer() {
	panic("implement me")
}
