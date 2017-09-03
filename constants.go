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

package main

import (
	"github.com/RackSec/srslog"
	"os"
)

const (
	AccessPriority = srslog.LOG_DAEMON|srslog.LOG_INFO
	ErrorPriority = srslog.LOG_DAEMON|srslog.LOG_ERR
	LogFileFlags = os.O_APPEND|os.O_CREATE|os.O_WRONLY
	LogFileMode = 0644
	LogDirMode = 0755
)

