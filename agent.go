package thousandeyes

import (
	"encoding/json"
	"fmt"
)

// Agents - list of agents
type Agents []Agent

// Agent - Agent struct
type Agent struct {
	AgentID               *int64               `json:"agentId,omitempty"`
	AgentName             *string              `json:"agentName,omitempty"`
	AgentType             *string              `json:"agentType,omitempty"`
	CountryID             *string              `json:"countryId,omitempty"`
	ClusterMembers        *[]ClusterMember     `json:"clusterMembers,omitempty"`
	IPAddresses           *[]string            `json:"ipAddresses,omitempty"`
	Groups                *GroupLabels         `json:"groups,omitempty"`
	Location              *string              `json:"location,omitempty"`
	ErrorDetails          *[]AgentErrorDetails `json:"errorDetails,omitempty"`
	Hostname              *string              `json:"hostname,omitempty"`
	Prefix                *string              `json:"prefix,omitempty"`
	Enabled               *bool                `json:"enabled,omitempty" te:"int-bool"`
	Network               *string              `json:"network,omitempty"`
	CreatedDate           *string              `json:"createdDate,omitempty"`
	LastSeen              *string              `json:"lastSeen,omitempty"`
	AgentState            *string              `json:"agentState,omitempty"`
	VerifySslCertificates *bool                `json:"verifySslCertificate,omitempty" te:"int-bool"`
	KeepBrowserCache      *bool                `json:"keepBrowserCache,omitempty" te:"int-bool"`
	Utilization           *int                 `json:"utilization,omitempty"`
	Ipv6Policy            *string              `json:"IPV6Policy,omitempty"`
	TargetForTests        *string              `json:"targetForTests,omitempty"`
}

//ClusterMember - ClusterMember struct
type ClusterMember struct {
	MemberID          *int64    `json:"memberId,omitempty"`
	Name              *string   `json:"name,omitempty"`
	IPAddresses       *[]string `json:"IPAddresses,omitempty"`
	PublicIPAddresses *[]string `json:"PublicIPAddresses,omitempty"`
	Prefix            *string   `json:"Prefix,omitempty"`
	Network           *string   `json:"network,omitempty"`
	LastSeen          *string   `json:"lastSeen,omitempty"`
	AgentState        *string   `json:"agentState,omitempty"`
	Utilization       *int      `json:"utilization,omitempty"`
	TargetForTests    *string   `json:"targetForTests,omitempty"`
}

// AgentErrorDetails - Agent error details
type AgentErrorDetails struct {
	Code        *string `json:"code,omitempty"`
	Description *string `json:"description,omitempty"`
}

// MarshalJSON implements the json.Marshaler interface. It ensures
// that ThousandEyes int fields that only use the values 0 or 1 are
// treated as booleans.
func (t Agent) MarshalJSON() ([]byte, error) {
	type alias Agent

	data, err := json.Marshal((alias)(t))
	if err != nil {
		return nil, err
	}

	return jsonBoolToInt(&t, data)
}

// UnmarshalJSON implements the json.Unmarshaler interface. It ensures
// that ThousandEyes int fields that only use the values 0 or 1 are
// treated as booleans.
func (t *Agent) UnmarshalJSON(data []byte) error {
	type alias Agent
	test := (*alias)(t)

	data, err := jsonIntToBool(t, data)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &test)
}

// GetAgents - Get agents
func (c *Client) GetAgents() (*Agents, error) {
	resp, err := c.get("/agents")
	if err != nil {
		return &Agents{}, err
	}
	var target map[string]Agents
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	agents := target["agents"]
	return &agents, nil
}

// GetAgent - Get agent
func (c *Client) GetAgent(id int) (*Agent, error) {
	resp, err := c.get(fmt.Sprintf("/agents/%d", id))
	if err != nil {
		return nil, err
	}
	var target map[string][]Agent
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	agent := target["agents"][0]
	return &agent, nil
}

// AddAgentsToCluster - add agent to cluster
func (c *Client) AddAgentsToCluster(cluster int, ids []int) (*[]Agent, error) {
	resp, err := c.post(fmt.Sprintf("/agents/%d/add-to-cluster", cluster), ids, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to add agents to cluster, response code %d", resp.StatusCode)
	}
	var target map[string][]Agent
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	agent := target["agents"]
	return &agent, nil
}

// RemoveAgentsFromCluster - remove agent from cluster
func (c *Client) RemoveAgentsFromCluster(cluster int, ids []int) (*[]Agent, error) {
	resp, err := c.post(fmt.Sprintf("/agents/%d/remove-from-cluster", cluster), ids, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to remove agents from cluster, response code %d", resp.StatusCode)
	}
	var target map[string][]Agent
	if dErr := c.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	agent := target["agents"]
	return &agent, nil
}
