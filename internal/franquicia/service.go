package franquicia

import (
	"clubhub-hotel-management/internal/domain"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

type service struct {
	repo Repository
}

type Service interface {
	CreateFranquicia(*gin.Context, *domain.Franquicia) error
	getSSLInfo(*gin.Context, string) (*domain.SSLInfo, error)
	scrapeLogoURL(*gin.Context, string) (string, error)
	getDomainInfo(*gin.Context, string) (*domain.DomainInfo, error)

	GetFranquiciaByID(*gin.Context, domain.Franquicia) error
	GetAllFranquicias(*gin.Context) ([]domain.Franquicia, error)
	UpdateFranquicia(*gin.Context, domain.Franquicia) error
}

// NewService crea un nuevo servicio de franquicia.
func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) CreateFranquicia(ctx *gin.Context, req *domain.Franquicia) error {
	var wg sync.WaitGroup
	var errWhois, errScrape, errSSL error
	var sslInfo *domain.SSLInfo
	var domainInfo *domain.DomainInfo
	var logoURL string

	wg.Add(3)

	go func() {
		defer wg.Done()
		var err error
		sslInfo, err = s.getSSLInfo(ctx, req.URL)
		if err != nil {
			log.Printf("error obteniendo SSL Info: %v", err)
			errSSL = err
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		domainInfo, err = s.getDomainInfo(ctx, req.URL)
		if err != nil {
			log.Printf("error obteniendo Domain Info: %v", err)
			errWhois = err
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		logoURL, err = s.scrapeLogoURL(ctx, ensureURLScheme(req.URL))
		if err != nil {
			log.Printf("error haciendo scraping del logo: %v", err)
			errScrape = err
		}
	}()

	wg.Wait()

	if errWhois != nil {
		log.Printf("error: %v", errWhois)
		return errWhois
	}
	if errScrape != nil {
		log.Printf("error: %v", errScrape)
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

	err := s.repo.Create(req)
	if err != nil {
		return err
	}
	log.Println("Franquicia created: \n", req)
	return nil
}

func (s *service) getSSLInfo(ctx *gin.Context, url string) (*domain.SSLInfo, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	apiURL := fmt.Sprintf("https://api.ssllabs.com/api/v3/analyze?host=%s", url)

	resp, err := client.Get(apiURL)
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

func (s *service) scrapeLogoURL(ctx *gin.Context, url string) (string, error) {
	c := colly.NewCollector(colly.Async(true))
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	var logoURL string
	c.OnHTML("img", func(e *colly.HTMLElement) {
		logoURL = e.Request.AbsoluteURL(e.Attr("src"))
		log.Println("Found image:", logoURL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(url)
	c.Wait()

	if logoURL == "" {
		return "", fmt.Errorf("no se pudo encontrar el logo en la URL proporcionada")
	}

	return logoURL, nil
}

func ensureURLScheme(urlStr string) string {
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "http://" + urlStr
	}
	return urlStr
}

func extractDomainName(urlStr string) (string, error) {
	parsedURL, err := url.Parse(ensureURLScheme(urlStr))
	if err != nil {
		return "", err
	}
	return parsedURL.Hostname(), nil
}

func (s *service) getDomainInfo(ctx *gin.Context, domainReq string) (*domain.DomainInfo, error) {
	domainName, err := extractDomainName(domainReq)
	if err != nil {
		return nil, err
	}

	whoisResult, err := whois.Whois(domainName)
	if err != nil {
		return nil, err
	}

	parsedResult, err := whoisparser.Parse(whoisResult)
	if err != nil {
		return nil, err
	}

	createdDate, err := time.Parse(time.RFC3339, parsedResult.Domain.CreatedDate)
	if err != nil {
		return nil, err
	}
	expiryDate, err := time.Parse(time.RFC3339, parsedResult.Domain.ExpirationDate)
	if err != nil {
		return nil, err
	}

	domainInfo := &domain.DomainInfo{
		CreatedDate:   createdDate.Format("2006-01-02 15:04:05"),
		ExpiryDate:    expiryDate.Format("2006-01-02 15:04:05"),
		RegistrarName: parsedResult.Registrar.Name,
		ContactEmail:  parsedResult.Administrative.Email,
	}

	return domainInfo, nil
}

func (s *service) GetFranquiciaByID(ctx *gin.Context, f domain.Franquicia) error {
	return s.repo.GetOne(f)
}

func (s *service) GetAllFranquicias(ctx *gin.Context) ([]domain.Franquicia, error) {
	fs, err := s.repo.GetAll()
	if err != nil {
		return []domain.Franquicia{}, err
	}
	return fs, nil
}

func (s *service) UpdateFranquicia(ctx *gin.Context, f domain.Franquicia) error {
	return s.repo.Update(f)
}
