package repositories

import (
	"context"
	"time"

	"github.com/ronak4195/personal-assistant/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransactionFilter struct {
	UserID        string
	Type          *models.TransactionType
	From          *time.Time
	To            *time.Time
	CategoryID    *string
	SubcategoryID *string
	Limit         int64
	Offset        int64
	SortDateAsc   bool
}

type TransactionRepository interface {
	Create(ctx context.Context, t *models.Transaction) error
	FindByID(ctx context.Context, id, userID string) (*models.Transaction, error)
	List(ctx context.Context, f TransactionFilter) ([]models.Transaction, int64, error)
	Update(ctx context.Context, t *models.Transaction) error
	Delete(ctx context.Context, id, userID string) error
	ListByDateRange(ctx context.Context, userID string, from, to time.Time) ([]models.Transaction, error)
}

type transactionRepository struct {
	col *mongo.Collection
}

func NewTransactionRepository(db *mongo.Database) TransactionRepository {
	return &transactionRepository{
		col: db.Collection("transactions"),
	}
}

func (r *transactionRepository) Create(ctx context.Context, t *models.Transaction) error {
	now := time.Now().UTC()
	t.CreatedAt = now
	t.UpdatedAt = now

	res, err := r.col.InsertOne(ctx, t)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(interface{ Hex() string }); ok {
		t.ID = oid.Hex()
	}
	return nil
}

func (r *transactionRepository) FindByID(ctx context.Context, id, userID string) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.col.FindOne(ctx, bson.M{
		"_id":    id,
		"userId": userID,
	}).Decode(&tx)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) List(ctx context.Context, f TransactionFilter) ([]models.Transaction, int64, error) {
	filter := bson.M{"userId": f.UserID}
	if f.Type != nil {
		filter["type"] = *f.Type
	}
	if f.From != nil || f.To != nil {
		dateRange := bson.M{}
		if f.From != nil {
			dateRange["$gte"] = *f.From
		}
		if f.To != nil {
			dateRange["$lte"] = *f.To
		}
		filter["date"] = dateRange
	}
	if f.CategoryID != nil {
		filter["categoryId"] = *f.CategoryID
	}
	if f.SubcategoryID != nil {
		filter["subcategoryId"] = *f.SubcategoryID
	}

	count, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOpts := options.Find().
		SetSkip(f.Offset).
		SetLimit(f.Limit)

	sort := bson.D{{Key: "date", Value: -1}}
	if f.SortDateAsc {
		sort = bson.D{{Key: "date", Value: 1}}
	}
	findOpts.SetSort(sort)

	cursor, err := r.col.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var res []models.Transaction
	if err := cursor.All(ctx, &res); err != nil {
		return nil, 0, err
	}
	return res, count, nil
}

func (r *transactionRepository) Update(ctx context.Context, t *models.Transaction) error {
	t.UpdatedAt = time.Now().UTC()
	_, err := r.col.UpdateOne(ctx, bson.M{
		"_id":    t.ID,
		"userId": t.UserID,
	}, bson.M{
		"$set": t,
	})
	return err
}

func (r *transactionRepository) Delete(ctx context.Context, id, userID string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{
		"_id":    id,
		"userId": userID,
	})
	return err
}

func (r *transactionRepository) ListByDateRange(ctx context.Context, userID string, from, to time.Time) ([]models.Transaction, error) {
	filter := bson.M{
		"userId": userID,
		"date": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}
	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var res []models.Transaction
	if err := cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}
