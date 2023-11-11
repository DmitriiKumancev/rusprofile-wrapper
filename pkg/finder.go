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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	comp, err := fetchGeneralInfo(ctx, inn)
	if err != nil {
		return nil, err
	}

	kpp, err := fetchKPP(ctx, comp)
	if err != nil {
		return nil, err
	}

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

func fetchGeneralInfo(ctx context.Context, inn *INN) (*company, error) {
	searchURL := fmt.Sprintf(fmt.Sprintf("%s/%s", serverAddress, searchQuery), inn.INN)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "request failed: %v", err)
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}(resp.Body)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't read request body: %v", err)
	}

	var data struct {
		Items []company `json:"ul"`
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't unmarshal to json: %v", err)
	}

	var comp *company
	for _, elem := range data.Items {
		elem.INN = strings.Trim(elem.INN, "!~")
		if elem.INN == inn.INN {
			comp = &elem
			break
		}
	}

	if comp == nil {
		return nil, status.Error(codes.NotFound, "no company with provided INN")
	}

	return comp, nil
}

func fetchKPP(ctx context.Context, comp *company) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s", serverAddress, comp.URL), nil)
	if err != nil {
		return "", status.Errorf(codes.Internal, "couldn't create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", status.Errorf(codes.Internal, "request failed: %v", err)
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}(resp.Body)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", status.Errorf(codes.Internal, "couldn't create html document: %v", err)
	}

	selection := doc.Find("#clip_kpp")
	if len(selection.Nodes) == 0 {
		return "", status.Errorf(codes.Internal, "no kpp value was found")
	}

	kpp := selection.Nodes[0].FirstChild.Data
	return kpp, nil
}
