package franquicia

import (
	"clubhub-hotel-management/internal/domain"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

type service struct {
	repo Repository
}

type Service interface {
	CreateFranquicia(req *domain.Franquicia) error

	getSSLInfo(url string) (*domain.SSLInfo, error)
	scrapeLogoURL(url string) (string, error)
	getDomainInfo(domainReq string) (*domain.DomainInfo, error)

	GetFranquiciaByID(f domain.Franquicia) error
	GetAllFranquicias(f domain.Franquicia) error
	UpdateFranquicia(f domain.Franquicia) error
}

// NewService crea un nuevo servicio de franquicia.
func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) CreateFranquicia(req *domain.Franquicia) error {
	var wg sync.WaitGroup
	var errWhois, errScrape, errSSL error
	var sslInfo *domain.SSLInfo
	var domainInfo *domain.DomainInfo
	var logoURL string

	wg.Add(3)

	go func() {
		defer wg.Done()
		var err error
		sslInfo, err = s.getSSLInfo(req.URL)
		if err != nil {
			log.Printf("error obteniendo SSL Info: %v", err)
			errSSL = err
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		domainInfo, err = s.getDomainInfo(req.URL)
		if err != nil {
			log.Printf("error obteniendo Domain Info: %v", err)
			errWhois = err
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		logoURL, err = s.scrapeLogoURL(req.URL)
		if err != nil {
			log.Printf("error haciendo scraping del logo: %v", err)
			errScrape = err
		}
	}()

	wg.Wait()

	if errWhois != nil {
		return errWhois
	}
	if errScrape != nil {
		return errScrape
	}
	if errSSL != nil {
		return errSSL
	}

	if sslInfo != nil && len(sslInfo.Endpoints) > 0 {
		endpoint := sslInfo.Endpoints[0]
		req.DomainInfo.SSLGrade = endpoint.Grade
		req.DomainInfo.Protocol = sslInfo.Protocol
		req.DomainInfo.IsProtocolSecure = sslInfo.Protocol == "HTTPS"

		var serverHops []string
		for _, ep := range sslInfo.Endpoints {
			serverHops = append(serverHops, ep.ServerName)
		}
		req.DomainInfo.ServerHops = serverHops
	}

	if domainInfo != nil {
		req.DomainInfo.CreatedDate = domainInfo.CreatedDate
		req.DomainInfo.ExpiryDate = domainInfo.ExpiryDate
		req.DomainInfo.RegistrarName = domainInfo.RegistrarName
		req.DomainInfo.ContactEmail = domainInfo.ContactEmail
		req.DomainInfo.Protocol = domainInfo.Protocol
		req.DomainInfo.IsProtocolSecure = domainInfo.IsProtocolSecure

		req.DomainInfo.DNSRecords = append(req.DomainInfo.DNSRecords, domainInfo.DNSRecords...)
	}
	req.LogoURL = logoURL
	req.IsWebsiteLive = (sslInfo != nil)

	err := s.repo.Create(*req)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) getSSLInfo(url string) (*domain.SSLInfo, error) {
	apiURL := fmt.Sprintf("https://api.ssllabs.com/api/v3/analyze?host=%s", url)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sslInfo domain.SSLInfo
	if err := json.NewDecoder(resp.Body).Decode(&sslInfo); err != nil {
		return nil, err
	}

	return &sslInfo, nil
}

func (s *service) scrapeLogoURL(url string) (string, error) {
	c := colly.NewCollector()

	var logoURL string
	c.OnHTML("img[class*='logo']", func(e *colly.HTMLElement) {
		logoURL = e.Attr("src")
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(url)

	if logoURL == "" {
		return "", fmt.Errorf("no se pudo encontrar el logo en la URL proporcionada")
	}

	return logoURL, nil
}

func (s *service) getDomainInfo(domainReq string) (*domain.DomainInfo, error) {
	whoisResult, err := whois.Whois(domainReq)
	if err != nil {
		return nil, err
	}

	parsedResult, err := whoisparser.Parse(whoisResult)
	if err != nil {
		return nil, err
	}

	domainInfo := &domain.DomainInfo{
		CreatedDate:   parsedResult.Domain.CreatedDate,
		ExpiryDate:    parsedResult.Domain.ExpirationDate,
		RegistrarName: parsedResult.Registrar.Name,
		ContactEmail:  parsedResult.Administrative.Email,
	}

	return domainInfo, nil
}

func (s *service) GetFranquiciaByID(f domain.Franquicia) error {
	return s.repo.GetOne(f)
}

func (s *service) GetAllFranquicias(f domain.Franquicia) error {
	return s.repo.Create(f)
}

func (s *service) UpdateFranquicia(f domain.Franquicia) error {
	return s.repo.Create(f)
}
