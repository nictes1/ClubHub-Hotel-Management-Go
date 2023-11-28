package franquicia

import (
	"clubhub-hotel-management/internal/domain"
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(ctx context.Context, franquicia *domain.Franquicia) error
	Update(ctx context.Context, f domain.Franquicia) error
	GetOne(ctx context.Context, id string) (domain.Franquicia, error)
	GetAll(ctx context.Context) ([]domain.Franquicia, error)
	GetByDateRange(ctx context.Context, startDate, endDate string) ([]domain.Franquicia, error)
	GetByLocation(ctx context.Context, city, country string) ([]domain.Franquicia, error)
	GetByFranchiseName(ctx context.Context, name string) ([]domain.Franquicia, error)
}

type repository struct {
	db *mongo.Collection
}

func NewRepository(db *mongo.Collection) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, franquicia *domain.Franquicia) error {
	_, err := r.db.InsertOne(ctx, franquicia)
	return err
}

func (r *repository) Update(ctx context.Context, f domain.Franquicia) error {
	filter := bson.M{"_id": f.ID}
	update := bson.M{"$set": bson.M{}}

	val := reflect.ValueOf(f)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag.Get("json")

		if tag == "_id" || (!field.IsValid() || field.IsZero()) {
			continue
		}

		update["$set"].(bson.M)[tag] = field.Interface()
	}
	_, err := r.db.UpdateOne(ctx, filter, update)
	return err
}

func (r *repository) GetOne(ctx context.Context, id string) (domain.Franquicia, error) {
	var franquicia domain.Franquicia
	filter := bson.M{"_id": id}
	err := r.db.FindOne(ctx, filter).Decode(&franquicia)
	return franquicia, err
}

func (r *repository) GetAll(ctx context.Context) ([]domain.Franquicia, error) {
	var franquicias []domain.Franquicia
	cursor, err := r.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var franquicia domain.Franquicia
		if err := cursor.Decode(&franquicia); err != nil {
			return nil, err
		}
		franquicias = append(franquicias, franquicia)
	}

	return franquicias, nil
}

func (r *repository) GetByFranchiseName(ctx context.Context, name string) ([]domain.Franquicia, error) {
	var franquicias []domain.Franquicia
	filter := bson.M{"name": bson.M{"$regex": name, "$options": "i"}}
	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var franquicia domain.Franquicia
		if err := cursor.Decode(&franquicia); err != nil {
			return nil, err
		}
		franquicias = append(franquicias, franquicia)
	}

	return franquicias, nil
}

func (r *repository) GetByLocation(ctx context.Context, city, country string) ([]domain.Franquicia, error) {
	var franquicias []domain.Franquicia
	filter := bson.M{
		"location.city":    city,
		"location.country": country,
	}
	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var franquicia domain.Franquicia
		if err := cursor.Decode(&franquicia); err != nil {
			return nil, err
		}
		franquicias = append(franquicias, franquicia)
	}

	return franquicias, nil
}

func (r *repository) GetByDateRange(ctx context.Context, startDate, endDate string) ([]domain.Franquicia, error) {
	var franquicias []domain.Franquicia
	filter := bson.M{
		"domain_info.created_date": bson.M{"$gte": startDate},
		"domain_info.expiry_date":  bson.M{"$lte": endDate},
	}
	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var franquicia domain.Franquicia
		if err := cursor.Decode(&franquicia); err != nil {
			return nil, err
		}
		franquicias = append(franquicias, franquicia)
	}

	return franquicias, nil
}
