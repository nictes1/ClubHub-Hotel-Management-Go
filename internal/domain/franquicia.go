package domain

// Franquicia representa los datos de una franquicia hotelera.
type Franquicia struct {
	Name          string     `json:"name"`
	URL           string     `json:"url"`
	Location      Location   `json:"location"`
	LogoURL       string     `json:"logo_url,omitempty"`    // URL del logo obtenida por scraping
	IsWebsiteLive bool       `json:"is_website_live"`       // Estado de disponibilidad del sitio web
	DomainInfo    DomainInfo `json:"domain_info,omitempty"` // Información del dominio obtenida mediante la librería Whois o API de SSL Labs
}

// DomainInfo contiene información detallada sobre el dominio de una franquicia.
type DomainInfo struct {
	CreatedDate      string      `json:"created_date,omitempty"`       // Fecha de creación del dominio
	ExpiryDate       string      `json:"expiry_date,omitempty"`        // Fecha de expiración del dominio
	RegistrarName    string      `json:"registrar_name,omitempty"`     // Nombre del registrador del dominio
	ContactEmail     string      `json:"contact_email,omitempty"`      // Email de contacto del titular del dominio
	Protocol         string      `json:"protocol,omitempty"`           // Protocolo de comunicación (ej. HTTP, HTTPS)
	IsProtocolSecure bool        `json:"is_protocol_secure,omitempty"` // Indica si el protocolo de comunicación es seguro (HTTPS)
	ServerHops       []string    `json:"server_hops,omitempty"`        // Nombres de los servidores por los que pasa la solicitud antes de llegar al host
	SSLGrade         string      `json:"ssl_grade,omitempty"`          // Calificación SSL del sitio web (obtenido de SSL Labs)
	DNSRecords       []DNSRecord `json:"dns_records,omitempty"`        // Registros DNS del dominio
}

type DNSRecord struct {
	Type     string `json:"type"`     // Tipo de registro DNS (ej. A, AAAA, CNAME, MX, etc.)
	Value    string `json:"value"`    // Valor del registro DNS (ej. dirección IP para registros tipo A)
	TTL      int    `json:"ttl"`      // Tiempo de vida del registro (Time To Live)
	Priority int    `json:"priority"` // Prioridad del registro (usado principalmente para registros MX)
}

type SSLInfo struct {
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	Protocol        string        `json:"protocol"`
	IsPublic        bool          `json:"isPublic"`
	Status          string        `json:"status"`
	StartTime       int64         `json:"startTime"`
	TestTime        int64         `json:"testTime"`
	EngineVersion   string        `json:"engineVersion"`
	CriteriaVersion string        `json:"criteriaVersion"`
	Endpoints       []SSLEndpoint `json:"endpoints"`
}

type SSLEndpoint struct {
	IPAddress         string `json:"ipAddress"`
	ServerName        string `json:"serverName"`
	StatusMessage     string `json:"statusMessage"`
	Grade             string `json:"grade"`
	GradeTrustIgnored string `json:"gradeTrustIgnored"`
	HasWarnings       bool   `json:"hasWarnings"`
	IsExceptional     bool   `json:"isExceptional"`
	Progress          int    `json:"progress"`
	Duration          int    `json:"duration"`
	Delegation        int    `json:"delegation"`
}
