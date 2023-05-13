package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sd "github.com/toysmars/distlib-go/pkg/servicediscovery"
)

type operatorTestDriver struct {
	ctx context.Context
	id  sd.Identity
	op  sd.Operator
}

type operatorTestAction struct {
	Action            string
	Instance          string
	ExpectedError     error
	ExpectedInstances []string
}

func createOperatorTestDriver(id sd.Identity, op sd.Operator) operatorTestDriver {
	return operatorTestDriver{
		ctx: context.Background(),
		id:  id,
		op:  op,
	}
}

func (ts *operatorTestDriver) RunTestActions(t *testing.T, actions []operatorTestAction) {
	t.Helper()

	availableDuration := 100 * time.Millisecond
	newIdWithInstance := func(instance string) sd.Identity {
		return sd.Identity{
			Namespace: ts.id.Namespace,
			Service:   ts.id.Service,
			Instance:  instance,
		}
	}
	newHealthStatus := func(instance string) sd.HealthStatus {
		now := time.Now()
		return sd.HealthStatus{
			Identity:       newIdWithInstance(instance),
			LastHeartbeat:  now,
			AvailableUntil: now.Add(availableDuration),
		}
	}
	containsInstance := func(healthStatusList []sd.HealthStatus, instance string) bool {
		for _, hs := range healthStatusList {
			if hs.Identity.Instance == instance {
				return true
			}
		}
		return false
	}

	for _, action := range actions {
		var healthStatusList []sd.HealthStatus
		var err error
		switch action.Action {
		case "list":
			healthStatusList, err = ts.op.ListHealthStatus(ts.ctx, ts.id)
		case "register":
			err = ts.op.RegisterHealthStatus(ts.ctx, newHealthStatus(action.Instance))
		case "deregister":
			err = ts.op.DeregisterHealthStatus(ts.ctx, newIdWithInstance(action.Instance))
		}
		if action.ExpectedError == nil {
			assert.NoError(t, err, action)
		} else {
			assert.True(t, errors.Is(err, action.ExpectedError), action)
		}
		if action.ExpectedInstances != nil {
			assert.Len(t, healthStatusList, len(action.ExpectedInstances))
			for _, instance := range action.ExpectedInstances {
				assert.True(t, containsInstance(healthStatusList, instance))
			}
		}
	}
}

func RunOperatorUnitTests(t *testing.T, op sd.Operator) {
	t.Run("service discovery operator test", func(t *testing.T) {
		id := sd.Identity{
			Namespace: "unit-test",
			Service:   "my-unit-test-service",
		}
		ts := createOperatorTestDriver(id, op)
		ts.RunTestActions(t, []operatorTestAction{
			{Action: "list", ExpectedInstances: []string{}},
			{Action: "register", Instance: "x"},
			{Action: "list", ExpectedInstances: []string{"x"}},
			{Action: "register", Instance: "x"},
			{Action: "list", ExpectedInstances: []string{"x"}},
			{Action: "register", Instance: "y"},
			{Action: "list", ExpectedInstances: []string{"x", "y"}},
			{Action: "register", Instance: "z"},
			{Action: "list", ExpectedInstances: []string{"x", "y", "z"}},
			{Action: "deregister", Instance: "y"},
			{Action: "list", ExpectedInstances: []string{"x", "z"}},
		})
	})
}
