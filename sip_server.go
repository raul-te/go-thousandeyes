package thousandeyes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// SIPServer - SIPServer trace test
type SIPServer struct {
	// Common test fields
	AlertsEnabled      *bool                `json:"alertsEnabled,omitempty" te:"int-bool"`
	AlertRules         *[]AlertRule         `json:"alertRules,omitempty"`
	APILinks           *[]APILink           `json:"apiLinks,omitempty"`
	CreatedBy          *string              `json:"createdBy,omitempty"`
	CreatedDate        *string              `json:"createdDate,omitempty"`
	Description        *string              `json:"description,omitempty"`
	Enabled            *bool                `json:"enabled,omitempty" te:"int-bool"`
	Groups             *[]GroupLabel        `json:"groups,omitempty"`
	ModifiedBy         *string              `json:"modifiedBy,omitempty"`
	ModifiedDate       *string              `json:"modifiedDate,omitempty"`
	SavedEvent         *bool                `json:"savedEvent,omitempty" te:"int-bool"`
	SharedWithAccounts *[]SharedWithAccount `json:"sharedWithAccounts,omitempty"`
	TestID             *int64               `json:"testId,omitempty"`
	TestName           *string              `json:"testName,omitempty"`
	Type               *string              `json:"type,omitempty"`
	LiveShare          *bool                `json:"liveShare,omitempty" te:"int-bool"`

	// Fields unique to this test
	Agents                *[]Agent     `json:"agents,omitempty"`
	BandwidthMeasurements *bool        `json:"bandwidthMeasurements,omitempty" te:"int-bool"`
	BGPMeasurements       *bool        `json:"bgpMeasurements,omitempty" te:"int-bool"`
	Interval              *int         `json:"interval,omitempty"`
	MTUMeasurements       *bool        `json:"mtuMeasurements,omitempty" te:"int-bool"`
	NetworkMeasurements   *bool        `json:"networkMeasurements,omitempty" te:"int-bool"`
	NumPathTraces         *int         `json:"numPathTraces,omitempty"`
	OptionsRegex          *string      `json:"options_regex,omitempty"`
	PathTraceMode         *string      `json:"pathTraceMode,omitempty"`
	ProbeMode             *string      `json:"probeMode,omitempty"`
	RegisterEnabled       *bool        `json:"registerEnabled,omitempty" te:"int-bool"`
	SIPTargetTime         *int         `json:"sipTargetTime,omitempty"`
	SIPTimeLimit          *int         `json:"sipTimeLimit,omitempty"`
	TargetSIPCredentials  *SIPAuthData `json:"targetSipCredentials,omitempty"`
	UsePublicBGP          *bool        `json:"usePublicBgp,omitempty" te:"int-bool"`
}

// MarshalJSON implements the json.Marshaler interface. It ensures
// that ThousandEyes int fields that only use the values 0 or 1 are
// treated as booleans.
func (t SIPServer) MarshalJSON() ([]byte, error) {
	type aliasTest SIPServer

	data, err := json.Marshal((aliasTest)(t))
	if err != nil {
		return nil, err
	}

	return jsonBoolToInt(&t, data)
}

// UnmarshalJSON implements the json.Unmarshaler interface. It ensures
// that ThousandEyes int fields that only use the values 0 or 1 are
// treated as booleans.
func (t *SIPServer) UnmarshalJSON(data []byte) error {
	type aliasTest SIPServer
	test := (*aliasTest)(t)

	data, err := jsonIntToBool(t, data)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &test)
}

// AddAgent - Add agemt to sip server  test
func (t *SIPServer) AddAgent(id int) {
	agent := Agent{AgentID: Int(id)}
	*t.Agents = append(*t.Agents, agent)
}

// AddAlertRule - Adds an alert to agent test
func (t *SIPServer) AddAlertRule(id int) {
	alertRule := AlertRule{RuleID: Int(id)}
	*t.AlertRules = append(*t.AlertRules, alertRule)
}

// GetSIPServer  - get sip server test
func (c *Client) GetSIPServer(id int) (*SIPServer, error) {
	resp, err := c.get(fmt.Sprintf("/tests/%d", id))
	if err != nil {
		return &SIPServer{}, err
	}

	// Duplicate http response so we can read JSON directly
	// and still use the normal client interaction to process
	// the http response.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not decode HTTP response: %v", err)
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBody))

	var target map[string][]SIPServer
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}

	// A design flaw in ThousandEyes V6 API results in field on sip-server tests which
	// should be part of a targetSipCredentials object (matching the behavior of voice-call
	// tests) are instead part of the sip-server test object itself.
	// As this is not intended to be fixed until V7, The solution will be to have a
	// separate struct for reads, which will be converted before being passed.
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBody))
	var sipTarget map[string][]SIPAuthData
	if dErr := c.decodeJSON(resp, &sipTarget); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	for i := range target["test"] {
		if sipTarget["test"][i].AuthUser != nil {
			sipAuth := &SIPAuthData{
				AuthUser:     sipTarget["test"][i].AuthUser,
				Password:     sipTarget["test"][i].Password,
				Port:         sipTarget["test"][i].Port,
				Protocol:     sipTarget["test"][i].Protocol,
				SIPProxy:     sipTarget["test"][i].SIPProxy,
				SIPRegistrar: sipTarget["test"][i].SIPRegistrar,
				User:         sipTarget["test"][i].User,
			}
			target["test"][i].TargetSIPCredentials = sipAuth
		}
	}
	return &target["test"][0], nil
}

//CreateSIPServer - Create sip server test
func (c Client) CreateSIPServer(t SIPServer) (*SIPServer, error) {
	resp, err := c.post("/tests/sip-server/new", t, nil)
	if err != nil {
		return &t, err
	}
	if resp.StatusCode != 201 {
		return &t, fmt.Errorf("failed to create sip-server test, response code %d", resp.StatusCode)
	}
	var target map[string][]SIPServer
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	return &target["test"][0], nil
}

//DeleteSIPServer - delete sip server test
func (c *Client) DeleteSIPServer(id int) error {
	resp, err := c.post(fmt.Sprintf("/tests/sip-server/%d/delete", id), nil, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		return fmt.Errorf("failed to delete sip test, response code %d", resp.StatusCode)
	}
	return nil
}

//UpdateSIPServer - - update sip server test
func (c *Client) UpdateSIPServer(id int, t SIPServer) (*SIPServer, error) {
	resp, err := c.post(fmt.Sprintf("/tests/sip-server/%d/update", id), t, nil)
	if err != nil {
		return &t, err
	}
	if resp.StatusCode != 200 {
		return &t, fmt.Errorf("failed to update test, response code %d", resp.StatusCode)
	}
	var target map[string][]SIPServer
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	return &target["test"][0], nil
}
