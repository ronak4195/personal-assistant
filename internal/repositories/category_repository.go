package repositories

import (
	"context"
	"log"
	"time"

	"github.com/ronak4195/personal-assistant/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryRepository interface {
	Create(ctx context.Context, c *models.Category) error
	FindByID(ctx context.Context, id, userID string) (*models.Category, error)
	List(ctx context.Context, userID string, parentID *string) ([]models.Category, error)
	Update(ctx context.Context, c *models.Category) error
	Delete(ctx context.Context, id, userID string) error
}

type categoryRepository struct {
	col *mongo.Collection
}

func NewCategoryRepository(db *mongo.Database) CategoryRepository {
	return &categoryRepository{
		col: db.Collection("categories"),
	}
}

func (r *categoryRepository) Create(ctx context.Context, c *models.Category) error {
	now := time.Now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now

	res, err := r.col.InsertOne(ctx, c)
	if err != nil {
		return err
	}

	if oid, ok := res.InsertedID.(interface{ Hex() string }); ok {
		c.ID = oid.Hex()
	}
	return nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id, userID string) (*models.Category, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err // invalid ID format
	}

	var cat models.Category
	log.Default().Println("Finding category by ID:", id, "for user:", userID)

	err = r.col.FindOne(ctx, bson.M{
		"_id":    objectID,
		"userId": userID,
	}).Decode(&cat)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &cat, nil
}

func (r *categoryRepository) List(ctx context.Context, userID string, parentID *string) ([]models.Category, error) {
	filter := bson.M{"userId": userID}
	if parentID != nil {
		filter["parentId"] = *parentID
	}

	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var res []models.Category
	if err := cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *categoryRepository) Update(ctx context.Context, c *models.Category) error {
	c.UpdatedAt = time.Now().UTC()
	_, err := r.col.UpdateOne(ctx, bson.M{
		"_id":    c.ID,
		"userId": c.UserID,
	}, bson.M{
		"$set": bson.M{
			"name":      c.Name,
			"parentId":  c.ParentID,
			"updatedAt": c.UpdatedAt,
		},
	})
	return err
}

func (r *categoryRepository) Delete(ctx context.Context, id, userID string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{
		"_id":    id,
		"userId": userID,
	})
	return err
}
