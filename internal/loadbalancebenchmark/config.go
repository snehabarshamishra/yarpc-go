package loadbalancingbenchmark

import (
	"github.com/pkg/errors"
	"strconv"
)

type TestConfig struct {
	m map[string]string
}

func NewTestConfig(m map[string]string) (*TestConfig, error) {
	return &TestConfig{
		m: m,
	}, nil
}

func (config *TestConfig) GetChooserType() (ChooserType, error) {
	ct, ok := config.m["ChooserType"]
	if !ok {
		return Unknown, errors.New("ChooserType field is not set")
	}
	return ChooserType(ct), nil
}

func (config *TestConfig) GetServerCount() (int, error) {
	scnt, ok := config.m["ServerCount"]
	if !ok {
		return -1, errors.New("ServerCount field is not set")
	}
	cnt, err := strconv.Atoi(scnt)
	if err != nil {
		return -1, err
	}
	return cnt, nil
}

func (config *TestConfig) GetClientCount() (int, error) {
	scnt, ok := config.m["ClientCount"]
	if !ok {
		return -1, errors.New("ClientCount field is not set")
	}
	cnt, err := strconv.Atoi(scnt)
	if err != nil {
		return -1, err
	}
	return cnt, nil
}
