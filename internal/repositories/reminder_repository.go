package repositories

import (
	"context"
	"time"

	"github.com/ronak4195/personal-assistant/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReminderFilter struct {
	UserID   string
	IsActive *bool
	From     *time.Time
	To       *time.Time
}

type ReminderRepository interface {
	Create(ctx context.Context, r *models.Reminder) error
	FindByID(ctx context.Context, id, userID string) (*models.Reminder, error)
	List(ctx context.Context, f ReminderFilter) ([]models.Reminder, error)
	Update(ctx context.Context, r *models.Reminder) error
	Delete(ctx context.Context, id, userID string) error
}

type reminderRepository struct {
	col *mongo.Collection
}

func NewReminderRepository(db *mongo.Database) ReminderRepository {
	return &reminderRepository{
		col: db.Collection("reminders"),
	}
}

func (r *reminderRepository) Create(ctx context.Context, rem *models.Reminder) error {
	now := time.Now().UTC()
	rem.CreatedAt = now
	rem.UpdatedAt = now

	res, err := r.col.InsertOne(ctx, rem)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(interface{ Hex() string }); ok {
		rem.ID = oid.Hex()
	}
	return nil
}

func (r *reminderRepository) FindByID(ctx context.Context, id, userID string) (*models.Reminder, error) {
	var rem models.Reminder
	err := r.col.FindOne(ctx, bson.M{
		"_id":    id,
		"userId": userID,
	}).Decode(&rem)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rem, nil
}

func (r *reminderRepository) List(ctx context.Context, f ReminderFilter) ([]models.Reminder, error) {
	filter := bson.M{"userId": f.UserID}
	if f.IsActive != nil {
		filter["isActive"] = *f.IsActive
	}
	if f.From != nil || f.To != nil {
		due := bson.M{}
		if f.From != nil {
			due["$gte"] = *f.From
		}
		if f.To != nil {
			due["$lte"] = *f.To
		}
		filter["dueAt"] = due
	}

	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var res []models.Reminder
	if err := cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *reminderRepository) Update(ctx context.Context, rem *models.Reminder) error {
	rem.UpdatedAt = time.Now().UTC()
	_, err := r.col.UpdateOne(ctx, bson.M{
		"_id":    rem.ID,
		"userId": rem.UserID,
	}, bson.M{
		"$set": rem,
	})
	return err
}

func (r *reminderRepository) Delete(ctx context.Context, id, userID string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{
		"_id":    id,
		"userId": userID,
	})
	return err
}
