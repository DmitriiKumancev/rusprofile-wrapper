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
		Items []struct {
			INN  string `json:"inn"`
			Name string `json:"raw_name"`
			CEO  string `json:"ceo_name"`
			URL  string `json:"url"`
		} `json:"ul"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal to json: %w", err)
	}

	if len(data.Items) == 0 {
		return nil, fmt.Errorf("no company with provided INN")
	}

	company := data.Items[0]
	req, err = http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s", serverAddress, company.URL), nil)
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

	kpp := selection.Nodes[0].Data

	return &Company{
		INN:  strings.Trim(company.INN, "!~"),
		KPP:  kpp,
		Name: company.Name,
		Ceo:  company.CEO,
	}, nil
}

func (c CompanyFinder) mustEmbedUnimplementedCompanyFinderServer() {
	panic("implement me")
}
