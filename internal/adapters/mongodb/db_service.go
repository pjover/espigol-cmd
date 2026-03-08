package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/pjover/espigol/internal/adapters/mongodb/dbo"
	"github.com/pjover/espigol/internal/domain/model"
	"github.com/pjover/espigol/internal/domain/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbService struct {
	configService ports.ConfigService
	uri           string
	database      string
}

func NewDbService(configService ports.ConfigService) ports.DbService {
	uri := configService.GetString("db.server")
	database := configService.GetString("db.name")

	return &dbService{
		configService: configService,
		uri:           uri,
		database:      database,
	}
}

func (d *dbService) UpsertPartner(partner *model.Partner) error {
	dboPartner := dbo.ConvertPartnerToDbo(partner)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.uri))
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer client.Disconnect(context.Background())

	coll := client.Database(d.database).Collection("partner")
	opts := options.Replace().SetUpsert(true)
	_, err = coll.ReplaceOne(ctx, bson.D{{Key: "_id", Value: dboPartner.Id}}, dboPartner, opts)
	if err != nil {
		return fmt.Errorf("upserting partner: %w", err)
	}

	return nil
}

func (d *dbService) FindPartnerByEmail(email string) (*model.Partner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.uri))
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	defer client.Disconnect(context.Background())

	coll := client.Database(d.database).Collection("partner")
	var result dbo.Partner
	err = coll.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("partner with email %s not found", email)
		}
		return nil, fmt.Errorf("finding partner by email: %w", err)
	}

	return dbo.ConvertPartnerToModel(result), nil
}

func (d *dbService) UpsertExpenseForecast(forecast *model.ExpenseForecast) error {
	dboForecast := dbo.ConvertExpenseForecastToDbo(forecast)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.uri))
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer client.Disconnect(context.Background())

	coll := client.Database(d.database).Collection("expense_forecast")
	opts := options.Replace().SetUpsert(true)
	_, err = coll.ReplaceOne(ctx, bson.D{{Key: "_id", Value: dboForecast.Id}}, dboForecast, opts)
	if err != nil {
		return fmt.Errorf("upserting expense forecast: %w", err)
	}

	return nil
}

func (d *dbService) GetPartnerByID(id int) (*model.Partner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.uri))
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	defer client.Disconnect(context.Background())

	coll := client.Database(d.database).Collection("partner")
	var result dbo.Partner
	err = coll.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("partner with ID %d not found", id)
		}
		return nil, fmt.Errorf("finding partner by ID: %w", err)
	}

	return dbo.ConvertPartnerToModel(result), nil
}

func (d *dbService) GetAllPartners() ([]*model.Partner, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.uri))
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	defer client.Disconnect(context.Background())

	coll := client.Database(d.database).Collection("partner")
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("finding all partners: %w", err)
	}
	defer cursor.Close(ctx)

	var partners []*model.Partner
	for cursor.Next(ctx) {
		var result dbo.Partner
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("decoding partner: %w", err)
		}
		partners = append(partners, dbo.ConvertPartnerToModel(result))
	}

	return partners, nil
}

func (d *dbService) DeletePartner(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.uri))
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer client.Disconnect(context.Background())

	coll := client.Database(d.database).Collection("partner")
	res, err := coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("deleting partner: %w", err)
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("partner with ID %d not found", id)
	}

	return nil
}

func (d *dbService) GetExpenseForecastByID(id int) (*model.ExpenseForecast, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.uri))
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	defer client.Disconnect(context.Background())

	coll := client.Database(d.database).Collection("expense_forecast")
	var result dbo.ExpenseForecast
	err = coll.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("expense forecast with ID %d not found", id)
		}
		return nil, fmt.Errorf("finding expense forecast by ID: %w", err)
	}

	partner, err := d.GetPartnerByID(result.PartnerId)
	if err != nil {
		return nil, fmt.Errorf("finding partner for expense forecast: %w", err)
	}

	return dbo.ConvertExpenseForecastToModel(result, partner), nil
}

func (d *dbService) GetAllExpenseForecasts() ([]*model.ExpenseForecast, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.uri))
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}
	defer client.Disconnect(context.Background())

	coll := client.Database(d.database).Collection("expense_forecast")
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("finding all expense forecasts: %w", err)
	}
	defer cursor.Close(ctx)

	var forecasts []*model.ExpenseForecast
	for cursor.Next(ctx) {
		var result dbo.ExpenseForecast
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("decoding expense forecast: %w", err)
		}
		partner, err := d.GetPartnerByID(result.PartnerId)
		if err != nil {
			return nil, fmt.Errorf("finding partner for expense forecast %d: %w", result.Id, err)
		}
		forecasts = append(forecasts, dbo.ConvertExpenseForecastToModel(result, partner))
	}

	return forecasts, nil
}

func (d *dbService) DeleteExpenseForecast(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(d.uri))
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer client.Disconnect(context.Background())

	coll := client.Database(d.database).Collection("expense_forecast")
	res, err := coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("deleting expense forecast: %w", err)
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("expense forecast with ID %d not found", id)
	}

	return nil
}
