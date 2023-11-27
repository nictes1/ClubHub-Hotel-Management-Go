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
	"github.com/joho/godotenv"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	repo Repository
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

type Service interface {
	CreateFranquicia(*gin.Context, *domain.Franquicia) error
	getSSLInfo(*gin.Context, string) (*domain.SSLInfo, error)
	scrapeLogoURL(*gin.Context, string) (string, error)
	getDomainInfo(*gin.Context, string) (*domain.DomainInfo, error)

	GetFranquiciaByID(ctx *gin.Context, id string) (domain.Franquicia, error)
	GetByLocation(ctx *gin.Context, city, country string) ([]domain.Franquicia, error)
	GetByDateRange(ctx *gin.Context, startDate, endDate string) ([]domain.Franquicia, error)
	GetByFranchiseName(ctx *gin.Context, name string) ([]domain.Franquicia, error)
	GetAllFranquicias(*gin.Context) ([]domain.Franquicia, error)
	UpdateFranquicia(*gin.Context, domain.Franquicia) error
}

// NewService crea un nuevo servicio de franquicia.
func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

// CreateFranquicia maneja la creación de una nueva franquicia.
func (s *service) CreateFranquicia(ctx *gin.Context, req *domain.Franquicia) error {
	var wg sync.WaitGroup
	errs := make(chan error, 4)

	wg.Add(4)

	// Goroutine para obtener información SSL
	go func() {
		defer wg.Done()
		sslInfo, err := s.getSSLInfo(ctx, req.URL)
		if err != nil {
			errs <- err
			return
		}

		if sslInfo != nil && len(sslInfo.Endpoints) > 0 {
			req.DomainInfo.SSLGrade = sslInfo.Endpoints[0].Grade
			req.DomainInfo.Protocol = sslInfo.Protocol
			req.DomainInfo.IsProtocolSecure = sslInfo.Protocol == "HTTPS"

			var serverHops []string
			for _, endpoint := range sslInfo.Endpoints {
				serverHops = append(serverHops, endpoint.ServerName)
			}
			req.DomainInfo.ServerHops = serverHops
		}
	}()

	go func() {
		defer wg.Done()
		info, err := s.getDomainInfo(ctx, req.URL)
		if err != nil {
			errs <- err
			return
		}
		req.DomainInfo.CreatedDate = info.CreatedDate
		req.DomainInfo.ExpiryDate = info.ExpiryDate
		req.DomainInfo.RegistrarName = info.RegistrarName
		req.DomainInfo.ContactEmail = info.ContactEmail
		req.DomainInfo.Protocol = info.Protocol
		req.DomainInfo.IsProtocolSecure = info.IsProtocolSecure
		req.DomainInfo.DNSRecords = info.DNSRecords
	}()

	go func() {
		defer wg.Done()
		logoURL, err := s.scrapeLogoURL(ctx, ensureURLScheme(req.URL))
		if err != nil {
			errs <- err
			return
		}
		req.LogoURL = logoURL
	}()

	// Goroutine para obtener la ubicación
	go func() {
		defer wg.Done()
		location, err := s.scrapeLocationInfo(ensureURLScheme(req.URL))
		if err != nil {
			errs <- err
			return
		}
		req.Location = *location
	}()

	wg.Wait()
	close(errs)

	for err := range errs {
		if err != nil {
			return err
		}
	}

	req.ID = primitive.NewObjectID()
	// Guardar la información de la franquicia
	err := s.repo.Create(ctx, req)
	if err != nil {
		log.Printf("Error al crear franquicia: %v", err)
		return err
	}

	log.Println("Franquicia creada con éxito: ", req)
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
		urlStr = "https://" + urlStr
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

	fmt.Println("administrative whois \n", parsedResult.Administrative)
	fmt.Println("billing whois\n", parsedResult.Billing)
	fmt.Println("billing whois\n", parsedResult.Billing)
	fmt.Println("Domain whois\n", parsedResult.Domain)
	fmt.Println("Registrant whois\n", parsedResult.Registrant)
	fmt.Println("Registrar whois\n", parsedResult.Registrar)
	fmt.Println("Technical whois\n", parsedResult.Technical)

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

func (s *service) scrapeLocationInfo(urlReq string) (*domain.Location, error) {

	c := colly.NewCollector(colly.Async(true))
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	location := &domain.Location{}

	c.OnHTML(".location-info .city", func(e *colly.HTMLElement) {
		location.City = strings.TrimSpace(e.Text)
	})
	c.OnHTML(".location-info .country", func(e *colly.HTMLElement) {
		location.Country = strings.TrimSpace(e.Text)
	})
	c.OnHTML(".location-info .address", func(e *colly.HTMLElement) {
		location.Address = strings.TrimSpace(e.Text)
	})
	c.OnHTML(".location-info .zip", func(e *colly.HTMLElement) {
		location.ZipCode = strings.TrimSpace(e.Text)
	})

	c.OnScraped(func(r *colly.Response) {
		log.Printf("Finished scraping: %s", r.Request.URL)
		if location.City == "" || location.Country == "" || location.Address == "" || location.ZipCode == "" {
			log.Println("No se pudo obtener la información completa de ubicación")
		}
	})

	err := c.Visit(urlReq)
	if err != nil {
		return nil, err
	}

	c.Wait()

	return location, nil
}

func (s *service) GetFranquiciaByID(ctx *gin.Context, id string) (domain.Franquicia, error) {
	result, err := s.repo.GetOne(ctx, id)
	if err != nil {
		return domain.Franquicia{}, err
	}

	return result, nil
}

func (s *service) GetByLocation(ctx *gin.Context, city, country string) ([]domain.Franquicia, error) {
	result, err := s.repo.GetByLocation(ctx, city, country)
	if err != nil {
		return []domain.Franquicia{}, err
	}
	return result, nil
}

func (s *service) GetByDateRange(ctx *gin.Context, startDate, endDate string) ([]domain.Franquicia, error) {
	result, err := s.repo.GetByDateRange(ctx, startDate, endDate)
	if err != nil {
		return []domain.Franquicia{}, err
	}
	return result, nil
}

func (s *service) GetByFranchiseName(ctx *gin.Context, name string) ([]domain.Franquicia, error) {
	result, err := s.repo.GetByFranchiseName(ctx, name)
	if err != nil {
		return []domain.Franquicia{}, err
	}
	return result, nil
}

func (s *service) GetAllFranquicias(ctx *gin.Context) ([]domain.Franquicia, error) {
	fs, err := s.repo.GetAll(ctx)
	if err != nil {
		return []domain.Franquicia{}, err
	}
	return fs, nil
}

func (s *service) UpdateFranquicia(ctx *gin.Context, f domain.Franquicia) error {
	return s.repo.Update(ctx, f)
}
