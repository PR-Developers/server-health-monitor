package repository

import (
	"fmt"

	"github.com/PR-Developers/server-health-monitor/internal/consts"
	"github.com/PR-Developers/server-health-monitor/internal/database"
	"github.com/PR-Developers/server-health-monitor/internal/logger"
	"github.com/PR-Developers/server-health-monitor/internal/types"
	"github.com/PR-Developers/server-health-monitor/internal/utils"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type healthRepository struct {
	*baseRepository
}

var (
	_ IHealthRepository = (*healthRepository)(nil)
)

// NewHealthRepository returns an instanced health repository
func NewHealthRepository() IHealthRepository {
	db, _ := database.Instance()

	return &healthRepository{
		baseRepository: &baseRepository{
			db:             db,
			collection:     db.Client().Database(utils.GetVariable(consts.DB_NAME)).Collection(consts.COLLECTION_HEALTH),
			collectionName: consts.COLLECTION_HEALTH,
			log:            logger.Instance(),
		},
	}
}

// Find all health data given a certain query
func (r *healthRepository) Find(query interface{}) ([]types.Health, error) {
	return r.FindWithFilter(query, nil)
}

// FindWithFilter returns all health data given a certain query and options
func (r *healthRepository) FindWithFilter(query interface{}, options *options.FindOptions) ([]types.Health, error) {
	cursor, err := r.collection.Find(r.db.Context(), query, options)
	if err != nil {
		msg := fmt.Sprintf("failed to read data from collection: %s with query: %s (%s)", r.collectionName, query, err.Error())
		r.log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	var data []types.Health
	defer cursor.Close(r.db.Context())
	for cursor.Next(r.db.Context()) {
		var record types.Health
		if err = cursor.Decode(&record); err != nil {
			r.log.Warningf("failed to read record on %s with query: %s", r.collectionName, query)
		}
		data = append(data, record)
	}

	return data, nil
}

// Insert a single health record into the database
func (r *healthRepository) Insert(data *types.Health) (string, error) {
	res, err := r.collection.InsertOne(r.db.Context(), data)
	if err != nil {
		msg := fmt.Sprintf("failed to insert data into collection: %s", r.collectionName)
		r.log.Error(msg)
		return "", fmt.Errorf(msg)
	}
	return fmt.Sprintf("%x", res.InsertedID), nil
}
