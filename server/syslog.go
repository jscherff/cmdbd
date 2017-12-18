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

package server

import (
	`fmt`
	`github.com/RackSec/srslog`
	`github.com/jscherff/cmdbd/common`
)

var (
	FacilityMap = map[string]srslog.Priority {

		`LOG_KERN`:	srslog.LOG_KERN,
		`LOG_USER`:	srslog.LOG_USER,
		`LOG_MAIL`:	srslog.LOG_MAIL,
		`LOG_DAEMON`:	srslog.LOG_DAEMON,
		`LOG_AUTH`:	srslog.LOG_AUTH,
		`LOG_SYSLOG`:	srslog.LOG_SYSLOG,
		`LOG_LPR`:	srslog.LOG_LPR,
		`LOG_NEWS`:	srslog.LOG_NEWS,
		`LOG_UUCP`:	srslog.LOG_UUCP,
		`LOG_CRON`:	srslog.LOG_CRON,
		`LOG_AUTHPRIV`:	srslog.LOG_AUTHPRIV,
		`LOG_FTP`:	srslog.LOG_FTP,

		`LOG_LOCAL0`:	srslog.LOG_LOCAL0,
		`LOG_LOCAL1`:	srslog.LOG_LOCAL1,
		`LOG_LOCAL2`:	srslog.LOG_LOCAL2,
		`LOG_LOCAL3`:	srslog.LOG_LOCAL3,
		`LOG_LOCAL4`:	srslog.LOG_LOCAL4,
		`LOG_LOCAL5`:	srslog.LOG_LOCAL5,
		`LOG_LOCAL6`:	srslog.LOG_LOCAL6,
		`LOG_LOCAL7`:	srslog.LOG_LOCAL7,
	}

	SeverityMap = map[string]srslog.Priority {

		`LOG_ALERT`:	srslog.LOG_ALERT,
		`LOG_CRIT`:	srslog.LOG_CRIT,
		`LOG_ERR`:	srslog.LOG_ERR,
		`LOG_WARNING`:	srslog.LOG_WARNING,
		`LOG_NOTICE`:	srslog.LOG_NOTICE,
		`LOG_INFO`:	srslog.LOG_INFO,
		`LOG_DEBUG`:	srslog.LOG_DEBUG,
	}
)

// Syslog contains the Syslog configuration.
type Syslog struct {
	*srslog.Writer
	Protocol	string
	Port		string
	Host		string
	Tag		string
	Facility	string
	Severity	string
}

// NewSyslog creates and initializes a new Syslog instance.
func NewSyslog(cf string) (*Syslog, error) {

	this := &Syslog{}

	if err := common.LoadConfig(this, cf); err != nil {
		return nil, err
	}

	var facility, severity srslog.Priority

	if fac, ok := FacilityMap[this.Facility]; !ok {
		facility = srslog.LOG_LOCAL7
	} else {
		facility = fac
	}

	if sev, ok := SeverityMap[this.Severity]; !ok {
		severity = srslog.LOG_INFO
	} else {
		severity = sev
	}

	raddr := fmt.Sprintf(`%s:%s`, this.Host, this.Port)

	if writer, err := srslog.Dial(this.Protocol, raddr, facility|severity, this.Tag); err != nil {
		return nil, err
	} else {
		this.Writer = writer
	}

	return this, nil
}
