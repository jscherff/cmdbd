// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import `fmt`

// SerialNumService is an interface that creates serial numbers from seed values.
type SerialNumService interface {
	Create(key string, seed int64) (serialNum string, err error)
	Format(key string) (serialFmt string, err error)
}

// serialNumService is a service that implements the SerialNumService interface.
type serialNumService struct {
	SerialFormat map[string]string
}

// NewSerialNumService returns an object that implements the SerialNumService interface.
func NewSerialNumService(conf *Config) (SerialNumService, error) {

	if conf.SerialFormat == nil || len(conf.SerialFormat) == 0 {
		return nil, fmt.Errorf(`empty serial format map`)
	}

	return &serialNumService{conf.SerialFormat}, nil
}

// Format returns the format string of the provided format key.
func (this *serialNumService) Format(key string) (string, error) {

	if format, ok := this.SerialFormat[key]; ok {
		return format, nil
	} else if format, ok := this.SerialFormat[`Default`]; ok {
		return format, nil
	} else {
		return ``, fmt.Errorf(`format key %q not found`, key)
	}
}

// Create generates a new serial number using the provided format key and seed.
func (this *serialNumService) Create(key string, seed int64) (string, error) {

	if format, err := this.Format(key); err != nil {
		return ``, err
	} else {
		return fmt.Sprintf(format, seed), nil
	}
}
