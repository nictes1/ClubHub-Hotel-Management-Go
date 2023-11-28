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
	getDomainInfo(ctx *gin.Context, domainReq string) (*domain.DomainInfo, *domain.Location, error)

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
	log.Println("Iniciando la creación de franquicia")

	var info *domain.DomainInfo
	var location *domain.Location
	var sslInfo *domain.SSLInfo
	var wg sync.WaitGroup
	errs := make(chan error, 3)

	var err error

	wg.Add(2)

	// Goroutine para obtener información SSL
	go func() {
		defer wg.Done()
		sslInfo, err = s.getSSLInfo(ctx, req.URL)
		if err != nil {
			errs <- fmt.Errorf("error obteniendo información SSL: %w", err)
			return
		}
		s.assignSSLInfo(req, sslInfo)
	}()

	go func() {
		defer wg.Done()
		info, location, err = s.getDomainInfo(ctx, req.URL)
		if err != nil {
			errs <- err
			return
		}
		req.DomainInfo = *info
		req.Location = *location

	}()

	// go func() {
	// 	defer wg.Done()
	// 	logoURL, err := s.scrapeLogoURL(ctx, ensureURLScheme(req.URL))
	// 	if err != nil {
	// 		errs <- err
	// 		return
	// 	}
	// 	req.LogoURL = logoURL
	// }()

	wg.Wait()
	close(errs)

	log.Println("Finalizando la espera de goroutines en CreateFranquicia")

	for err := range errs {
		if err != nil {
			return err
		}
	}

	req.ID = primitive.NewObjectID()
	req.Name = info.RegistrarName
	log.Println("Finalizando la espera de goroutines en CreateFranquicia")
	err = s.repo.Create(ctx, req)
	if err != nil {
		log.Printf("Error al crear franquicia: %v", err)
		return err
	}

	log.Println("Franquicia creada con éxito: ", req)
	return nil
}

func (s *service) assignSSLInfo(req *domain.Franquicia, sslInfo *domain.SSLInfo) {
	if sslInfo != nil && len(sslInfo.Endpoints) > 0 {
		req.DomainInfo.SSLGrade = sslInfo.Endpoints[0].Grade
		req.DomainInfo.Protocol = sslInfo.Protocol
		req.DomainInfo.IsProtocolSecure = sslInfo.Protocol == "HTTPS"
		for _, endpoint := range sslInfo.Endpoints {
			req.DomainInfo.ServerHops = append(req.DomainInfo.ServerHops, endpoint.ServerName)
		}
	}
}

func (s *service) getSSLInfo(ctx *gin.Context, url string) (*domain.SSLInfo, error) {
	log.Printf("Obteniendo información SSL para URL: %s", url)
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
	log.Printf("Buscando logo en URL: %s", url)
	c := colly.NewCollector(
		colly.Async(true),
		// colly timeout
		colly.MaxDepth(1),
	)
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	var logoURL string
	var found bool

	c.OnHTML("img", func(e *colly.HTMLElement) {
		if !found {
			logoURL = e.Request.AbsoluteURL(e.Attr("src"))
			log.Println("Found image:", logoURL)
			found = true
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := c.Visit(url)
	if err != nil {
		return "", err
	}

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

func (s *service) getDomainInfo(ctx *gin.Context, domainReq string) (*domain.DomainInfo, *domain.Location, error) {
	log.Printf("Obteniendo información de dominio para: %s", domainReq)
	domainName, err := extractDomainName(domainReq)
	if err != nil {
		return nil, nil, err
	}

	whoisResult, err := whois.Whois(domainName)
	if err != nil {
		return nil, nil, err
	}

	parsedResult, err := whoisparser.Parse(whoisResult)
	if err != nil {
		return nil, nil, err
	}

	createdDate, err := time.Parse(time.RFC3339, parsedResult.Domain.CreatedDate)
	if err != nil {
		return nil, nil, err
	}
	expiryDate, err := time.Parse(time.RFC3339, parsedResult.Domain.ExpirationDate)
	if err != nil {
		return nil, nil, err
	}

	registrarInfo := domain.RegistrarInfo{
		Organization: parsedResult.Registrar.Name,
		Address:      parsedResult.Registrar.Street,
		City:         parsedResult.Registrar.City,
		State:        parsedResult.Registrar.Province,
		PostalCode:   parsedResult.Registrar.PostalCode,
		Country:      parsedResult.Registrar.Country,
		Phone:        parsedResult.Registrar.Phone,
		Fax:          parsedResult.Registrar.Fax,
		Email:        parsedResult.Registrar.Email,
	}

	technicalInfo := domain.TechnicalInfo{
		Organization: parsedResult.Technical.Organization,
		Address:      parsedResult.Technical.Street,
		City:         parsedResult.Technical.City,
		State:        parsedResult.Technical.Province,
		PostalCode:   parsedResult.Technical.PostalCode,
		Country:      parsedResult.Technical.Country,
		Phone:        parsedResult.Technical.Phone,
		Fax:          parsedResult.Technical.Fax,
		Email:        parsedResult.Technical.Email,
	}

	location := &domain.Location{
		City:    parsedResult.Administrative.City,
		Country: parsedResult.Administrative.Country,
		Address: parsedResult.Administrative.Street,
		ZipCode: parsedResult.Administrative.PostalCode,
	}

	domainInfo := &domain.DomainInfo{
		CreatedDate:   createdDate.Format("2006-01-02 15:04:05"),
		ExpiryDate:    expiryDate.Format("2006-01-02 15:04:05"),
		RegistrarName: parsedResult.Administrative.Name,
		ContactEmail:  parsedResult.Administrative.Email,
		RegistrarInfo: registrarInfo,
		TechnicalInfo: technicalInfo,
	}

	fmt.Println("ParsedAdministrative result: ", parsedResult.Administrative)

	return domainInfo, location, nil
}

// Funcion para obtener la ubicacion de una pagina. (params: etiquetas css.)
func (s *service) ScrapeLocationInfo(urlReq, city, country, adrress, info string) (*domain.Location, error) {

	c := colly.NewCollector(colly.Async(true))
	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	location := &domain.Location{}

	c.OnHTML(city, func(e *colly.HTMLElement) {
		location.City = strings.TrimSpace(e.Text)
	})
	c.OnHTML(country, func(e *colly.HTMLElement) {
		location.Country = strings.TrimSpace(e.Text)
	})
	c.OnHTML(adrress, func(e *colly.HTMLElement) {
		location.Address = strings.TrimSpace(e.Text)
	})
	c.OnHTML(info, func(e *colly.HTMLElement) {
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
