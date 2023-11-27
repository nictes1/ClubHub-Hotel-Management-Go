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
	ID               string      `json:"id" bson:"_id"`
	CreatedDate      string      `json:"created_date,omitempty" bson:"created_date,omitempty"`
	ExpiryDate       string      `json:"expiry_date,omitempty" bson:"expiry_date,omitempty"`
	RegistrarName    string      `json:"registrar_name,omitempty" bson:"registrar_name,omitempty"`
	ContactEmail     string      `json:"contact_email,omitempty" bson:"contact_email,omitempty"`
	Protocol         string      `json:"protocol,omitempty" bson:"protocol,omitempty"`
	IsProtocolSecure bool        `json:"is_protocol_secure,omitempty" bson:"is_protocol_secure,omitempty"`
	ServerHops       []string    `json:"server_hops,omitempty" bson:"server_hops,omitempty"`
	SSLGrade         string      `json:"ssl_grade,omitempty" bson:"ssl_grade,omitempty"`
	DNSRecords       []DNSRecord `json:"dns_records,omitempty" bson:"dns_records,omitempty"`
}

type DNSRecord struct {
	ID       string `json:"id" bson:"_id"`
	Type     string `json:"type" bson:"type"`
	Value    string `json:"value" bson:"value"`
	TTL      int    `json:"ttl" bson:"ttl"`
	Priority int    `json:"priority" bson:"priority"`
}

type SSLInfo struct {
	ID              string        `json:"id" bson:"_id"`
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
	ID                string `json:"id" bson:"_id"`
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
