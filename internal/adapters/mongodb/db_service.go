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
