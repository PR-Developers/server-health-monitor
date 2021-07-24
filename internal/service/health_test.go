package service

import (
	"fmt"
	"testing"

	"github.com/PR-Developers/server-health-monitor/internal/repository"
	"github.com/PR-Developers/server-health-monitor/internal/service/mocks"
	"github.com/PR-Developers/server-health-monitor/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

//go:generate mockery --dir=../../ -r --name IHealthRepository

type testServiceHelper struct {
	service IHealthService
	repo    repository.IHealthRepository
	mock    *mock.Mock
}

var (
	healthData []types.Health = []types.Health{
		{
			AgentID: "1",
			Host: types.Host{
				Hostname: "test",
				Platform: "ubuntu",
				OS:       "linux",
			},
			CreateTime: 1,
		},
		{
			AgentID: "2",
			Host: types.Host{
				Hostname: "test",
				Platform: "ubuntu",
				OS:       "linux",
			},
			CreateTime: 2,
		},
	}
)

func getInitializedHealthService() testServiceHelper {
	repo := new(mocks.IHealthRepository)
	// repo.On
	service := NewHealthService(repo)

	return testServiceHelper{
		service: service,
		repo:    repo,
		mock:    &repo.Mock,
	}
}

func TestHealth_GetHealth_ReturnsExpectedHealthData(t *testing.T) {
	helper := getInitializedHealthService()
	helper.mock.On("Find", bson.M{}).Return(healthData, nil)

	response := helper.service.GetHealth("1")

	data := response.Data.([]types.Health)

	assert.Equal(t, 2, len(data))
	assert.Equal(t, int64(1), data[0].CreateTime)
	assert.Equal(t, "test", data[0].Host.Hostname)
	assert.Equal(t, "ubuntu", data[0].Host.Platform)
	assert.Equal(t, "linux", data[0].Host.OS)
	assert.Equal(t, int64(2), data[1].CreateTime)

	helper.mock.AssertExpectations(t)
}

func TestHealth_GetHealth_HandlesError(t *testing.T) {
	helper := getInitializedHealthService()
	helper.mock.On("Find", bson.M{}).Return(nil, fmt.Errorf("failed to get data from DB"))

	response := helper.service.GetHealth("1")

	assert.Nil(t, response.Data)
	assert.Equal(t, 500, response.StatusCode)
	assert.Equal(t, "failed to get all health data - Request ID: 1", response.Error)
	assert.False(t, response.Success)

	helper.mock.AssertExpectations(t)
}

func TestHealth_GetHealthByAgentId_ReturnsExpectedHealthData(t *testing.T) {
	helper := getInitializedHealthService()
	helper.mock.On("Find", bson.M{"agentID": "1"}).Return([]types.Health{healthData[0]}, nil)

	response := helper.service.GetHealthByAgentID("1", "1")

	data := response.Data.([]types.Health)

	assert.Equal(t, 1, len(data))

	helper.mock.AssertExpectations(t)
}

func TestHealth_GetHealthByAgentId_HandlesError(t *testing.T) {
	helper := getInitializedHealthService()
	helper.mock.On("Find", bson.M{"agentID": "4"}).Return(nil, fmt.Errorf("failed to get data from DB"))

	response := helper.service.GetHealthByAgentID("1", "4")

	assert.Nil(t, response.Data)
	assert.Equal(t, 500, response.StatusCode)
	assert.Equal(t, "failed to get data for agent: 4 - Request ID: 1", response.Error)
	assert.False(t, response.Success)

	helper.mock.AssertExpectations(t)
}

func TestHealth_AddHealth_AddsExpectedHealthData(t *testing.T) {
	helper := getInitializedHealthService()
	helper.mock.On("Insert", &healthData[0]).Return("1234567", nil)

	response := helper.service.AddHealth("1", "1", &healthData[0])

	data := response.Data.(*types.Health)

	assert.True(t, response.Success)
	assert.NotEmpty(t, data.ID)

	helper.mock.AssertExpectations(t)
}

func TestHealth_AddHealth_HandlesError(t *testing.T) {
	helper := getInitializedHealthService()
	helper.mock.On("Insert", &healthData[1]).Return("", fmt.Errorf("failed to insert data into DB"))

	response := helper.service.AddHealth("1", "2", &healthData[1])

	assert.Nil(t, response.Data)
	assert.Equal(t, 500, response.StatusCode)
	assert.Equal(t, "failed to insert data for agent: 2 - Request ID 1", response.Error)
	assert.False(t, response.Success)

	helper.mock.AssertExpectations(t)
}