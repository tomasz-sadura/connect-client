package connect

import (
	"context"
	"fmt"
	"net/url"
)

func (c *Client) ListConnectors(ctx context.Context) ([]string, error) {
	var connectorNames []string
	response, err := c.client.NewRequest().
		SetContext(ctx).
		SetResult(&connectorNames).
		SetError(ApiError{}).
		Get("/connectors")
	if err != nil {
		return nil, err
	}

	err = getErrorFromResponse(response)
	if err != nil {
		return nil, err
	}

	return connectorNames, nil
}

// ListConnectorsOptions describe the available options to list connectors. Either Status or Info must be set to true.
type ListConnectorsOptions struct {
	ExpandStatus bool
	ExpandInfo   bool
}

func (l *ListConnectorsOptions) Validate() error {
	if !l.ExpandStatus && !l.ExpandInfo {
		return fmt.Errorf("either info or status must be set to true")
	}

	return nil
}

// ListConnectorsResponseExpanded is the response to /connectors if the expand query parameters are set.
type ListConnectorsResponseExpanded struct {
	Info   ConnectorInfo      `json:"info"`
	Status ConnectorStateInfo `json:"status"`
}

// ListConnectorsResponseExpandedInfo represents the Info object for described connectors.
// Deprecated: Use ConnectorInfo instead, the same information is available there
type ListConnectorsResponseExpandedInfo struct {
	Name   string            `json:"name"`
	Config map[string]string `json:"config"`
	Tasks  []struct {
		Connector string `json:"connector"`
		Task      int    `json:"task"`
	} `json:"tasks"`
	Type string `json:"type"`
}

// ListConnectorsResponseExpandedStatus represents the Status object for described connectors.
// Deprecated: Use ConnectorStateInfo instead, the same information is available there
type ListConnectorsResponseExpandedStatus struct {
	Name      string `json:"name"`
	Connector struct {
		State    string `json:"state"`
		WorkerID string `json:"worker_id"`
		Trace    string `json:"trace,omitempty"`
	}
	Tasks []struct {
		ID       int    `json:"id"`
		State    string `json:"state"`
		WorkerID string `json:"worker_id"`
		Trace    string `json:"trace,omitempty"`
	} `json:"tasks"`
	Type string `json:"type"`
}

func (c *Client) ListConnectorsExpanded(ctx context.Context) (map[string]ListConnectorsResponseExpanded, error) {
	// Adds additional options that show us more information about the connectors list
	expands := []string{"info", "status"}

	connectors := map[string]ListConnectorsResponseExpanded{}
	response, err := c.client.NewRequest().
		SetContext(ctx).
		SetResult(&connectors).
		SetError(ApiError{}).
		SetQueryParamsFromValues(url.Values{"expand": expands}).
		Get("/connectors")
	if err != nil {
		return nil, err
	}

	err = getErrorFromResponse(response)
	if err != nil {
		return nil, err
	}

	return connectors, nil
}
