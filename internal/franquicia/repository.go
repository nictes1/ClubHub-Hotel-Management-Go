package franquicia

import (
	"clubhub-hotel-management/internal/domain"
	"database/sql"
)

const (
	SAVE_USER = ""
	GET_ONE   = ""
	GET_ALL   = ""
	EXIST     = ""
	DELETE    = ""
)

type Repository interface {
	Create(franquicia *domain.Franquicia) error
	Update(f domain.Franquicia) error
	GetOne(f domain.Franquicia) error
	GetAll() ([]domain.Franquicia, error)
}

type repo struct {
	db *sql.DB
}

// NewRepository crea un nuevo repositorio de franquicia.
func NewRepository(mysql *sql.DB) Repository {
	return &repo{
		db: mysql,
	}
}

func (r *repo) Create(franquicia *domain.Franquicia) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var locationID int64
	err = r.db.QueryRow("SELECT id FROM locations WHERE city = ? AND country = ? AND address = ? AND zip_code = ?", franquicia.Location.City, franquicia.Location.Country, franquicia.Location.Address, franquicia.Location.ZipCode).Scan(&locationID)
	if err == sql.ErrNoRows {
		res, err := tx.Exec("INSERT INTO locations (city, country, address, zip_code) VALUES (?, ?, ?, ?)",
			franquicia.Location.City, franquicia.Location.Country, franquicia.Location.Address, franquicia.Location.ZipCode)
		if err != nil {
			tx.Rollback()
			return err
		}
		locationID, err = res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return err
		}
	} else if err != nil {
		tx.Rollback()
		return err
	}

	res, err := tx.Exec("INSERT INTO franchises (name, url, location_id, logo_url, is_website_live, created_date, expiry_date, registrar_name, contact_email, protocol, is_protocol_secure, ssl_grade) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		franquicia.Name, franquicia.URL, locationID, franquicia.LogoURL, franquicia.IsWebsiteLive, franquicia.DomainInfo.CreatedDate, franquicia.DomainInfo.ExpiryDate, franquicia.DomainInfo.RegistrarName, franquicia.DomainInfo.ContactEmail, franquicia.DomainInfo.Protocol, franquicia.DomainInfo.IsProtocolSecure, franquicia.DomainInfo.SSLGrade)
	if err != nil {
		tx.Rollback()
		return err
	}

	franquiciaID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, record := range franquicia.DomainInfo.DNSRecords {
		_, err = tx.Exec("INSERT INTO dns_records (franchise_id, type, value, ttl, priority) VALUES (?, ?, ?, ?, ?)",
			franquiciaID, record.Type, record.Value, record.TTL, record.Priority)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *repo) Update(f domain.Franquicia) error {
	return nil
}

func (r *repo) GetOne(f domain.Franquicia) error {
	return nil
}

func (r *repo) GetAll() ([]domain.Franquicia, error) {
	var franquicias []domain.Franquicia

	rows, err := r.db.Query("SELECT f.id, f.name, f.url, l.city, l.country, l.address, l.zip_code, f.logo_url, f.is_website_live, f.created_date, f.expiry_date, f.registrar_name, f.contact_email, f.protocol, f.is_protocol_secure, f.ssl_grade FROM franchises f LEFT JOIN locations l ON f.location_id = l.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f domain.Franquicia
		if err := rows.Scan(&f.ID, &f.Name, &f.URL, &f.Location.City, &f.Location.Country, &f.Location.Address, &f.Location.ZipCode, &f.LogoURL, &f.IsWebsiteLive, &f.DomainInfo.CreatedDate, &f.DomainInfo.ExpiryDate, &f.DomainInfo.RegistrarName, &f.DomainInfo.ContactEmail, &f.DomainInfo.Protocol, &f.DomainInfo.IsProtocolSecure, &f.DomainInfo.SSLGrade); err != nil {
			return nil, err
		}
		franquicias = append(franquicias, f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return franquicias, nil
}
