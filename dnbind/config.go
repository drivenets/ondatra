// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dnbind

import (
	"fmt"

	"github.com/openconfig/ondatra/dnbind/creds"
)

// Config contains parameters to configure the Drivenets binding.
type Config struct {
	Credentials *creds.Credentials `yaml:"credentials"`
}

func (c *Config) String() string {
	return fmt.Sprintf("%+v", *c)
}

// ValidateConfig checks if the provided config is valid.
func ValidateConfig(cfg *Config) error {
	return nil
}
