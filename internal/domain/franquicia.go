package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type FranquiciaRequest struct {
	ID       string   `json:"id,omitempty" bson:"_id,omitempty"`
	URL      string   `json:"url" bson:"url"`
	Name     string   `json:"name" bson:"name"`
	Location Location `json:"location" bson:"location"`
}

type Franquicia struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	Name          string             `json:"name" bson:"name"`
	URL           string             `json:"url" bson:"url"`
	Location      Location           `json:"location" bson:"location"`
	LogoURL       string             `json:"logo_url,omitempty" bson:"logo_url,omitempty"`
	IsWebsiteLive bool               `json:"is_website_live" bson:"is_website_live"`
	DomainInfo    DomainInfo         `json:"domain_info,omitempty" bson:"domain_info,omitempty"`
}

type DomainInfo struct {
	CreatedDate      string        `json:"created_date,omitempty" bson:"created_date,omitempty"`
	ExpiryDate       string        `json:"expiry_date,omitempty" bson:"expiry_date,omitempty"`
	RegistrarName    string        `json:"registrar_name,omitempty" bson:"registrar_name,omitempty"`
	ContactEmail     string        `json:"contact_email,omitempty" bson:"contact_email,omitempty"`
	Protocol         string        `json:"protocol,omitempty" bson:"protocol,omitempty"`
	IsProtocolSecure bool          `json:"is_protocol_secure,omitempty" bson:"is_protocol_secure,omitempty"`
	ServerHops       []string      `json:"server_hops,omitempty" bson:"server_hops,omitempty"`
	SSLGrade         string        `json:"ssl_grade,omitempty" bson:"ssl_grade,omitempty"`
	DNSRecords       []DNSRecord   `json:"dns_records,omitempty" bson:"dns_records,omitempty"`
	RegistrarInfo    RegistrarInfo `json:"registrar_info,omitempty" bson:"registrar_info,omitempty"`
	TechnicalInfo    TechnicalInfo `json:"technical_info,omitempty" bson:"technical_info,omitempty"`
}
type DNSRecord struct {
	Type     string `json:"type" bson:"type"`
	Value    string `json:"value" bson:"value"`
	TTL      int    `json:"ttl" bson:"ttl"`
	Priority int    `json:"priority" bson:"priority"`
}

type SSLInfo struct {
	Host            string        `json:"host" bson:"host"`
	Port            int           `json:"port" bson:"port"`
	Protocol        string        `json:"protocol" bson:"protocol"`
	IsPublic        bool          `json:"isPublic" bson:"isPublic"`
	Status          string        `json:"status" bson:"status"`
	StartTime       int64         `json:"startTime" bson:"startTime"`
	TestTime        int64         `json:"testTime" bson:"testTime"`
	EngineVersion   string        `json:"engineVersion" bson:"engineVersion"`
	CriteriaVersion string        `json:"criteriaVersion" bson:"criteriaVersion"`
	Endpoints       []SSLEndpoint `json:"endpoints" bson:"endpoints"`
}

type SSLEndpoint struct {
	IPAddress         string `json:"ipAddress" bson:"ipAddress"`
	ServerName        string `json:"serverName" bson:"serverName"`
	StatusMessage     string `json:"statusMessage" bson:"statusMessage"`
	Grade             string `json:"grade" bson:"grade"`
	GradeTrustIgnored string `json:"gradeTrustIgnored" bson:"gradeTrustIgnored"`
	HasWarnings       bool   `json:"hasWarnings" bson:"hasWarnings"`
	IsExceptional     bool   `json:"isExceptional" bson:"isExceptional"`
	Progress          int    `json:"progress" bson:"progress"`
	Duration          int    `json:"duration" bson:"duration"`
	Delegation        int    `json:"delegation" bson:"delegation"`
}

type RegistrarInfo struct {
	Organization string `json:"organization,omitempty" bson:"organization,omitempty"`
	Address      string `json:"address,omitempty" bson:"address,omitempty"`
	City         string `json:"city,omitempty" bson:"city,omitempty"`
	State        string `json:"state,omitempty" bson:"state,omitempty"`
	PostalCode   string `json:"postal_code,omitempty" bson:"postal_code,omitempty"`
	Country      string `json:"country,omitempty" bson:"country,omitempty"`
	Phone        string `json:"phone,omitempty" bson:"phone,omitempty"`
	Fax          string `json:"fax,omitempty" bson:"fax,omitempty"`
	Email        string `json:"email,omitempty" bson:"email,omitempty"`
}

type TechnicalInfo struct {
	Organization string `json:"organization,omitempty" bson:"organization,omitempty"`
	Address      string `json:"address,omitempty" bson:"address,omitempty"`
	City         string `json:"city,omitempty" bson:"city,omitempty"`
	State        string `json:"state,omitempty" bson:"state,omitempty"`
	PostalCode   string `json:"postal_code,omitempty" bson:"postal_code,omitempty"`
	Country      string `json:"country,omitempty" bson:"country,omitempty"`
	Phone        string `json:"phone,omitempty" bson:"phone,omitempty"`
	Fax          string `json:"fax,omitempty" bson:"fax,omitempty"`
	Email        string `json:"email,omitempty" bson:"email,omitempty"`
}

type Location struct {
	City      string  `json:"city" bson:"city"`
	Country   string  `json:"country" bson:"country"`
	Address   string  `json:"address" bson:"address"`
	ZipCode   string  `json:"zip_code" bson:"zip_code"`
	Latitude  float64 `json:"latitude,omitempty" bson:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty" bson:"longitude,omitempty"`
}
