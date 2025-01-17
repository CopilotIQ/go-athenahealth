package athenahealth

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type Provider struct {
	ANSINameCode                string       `json:"ansinamecode"`
	ANSISpecialtyCode           string       `json:"ansispecialtycode"`
	Billable                    bool         `json:"billable"`
	CreateEncounterOnCheckIn    bool         `json:"createencounteroncheckin"`
	DisplayName                 string       `json:"displayname"`
	EntityType                  string       `json:"entitytype"`
	FirstName                   string       `json:"firstname"`
	HideInPortal                bool         `json:"hideinportal"`
	LastName                    string       `json:"lastname"`
	NPI                         int          `json:"npi"`
	OtherProviderIDList         []string     `json:"otherprovideridlist"`
	ProviderID                  int          `json:"providerid"`
	ProviderType                string       `json:"providertype"`
	ProviderTypeID              string       `json:"providertypeid"`
	ProviderUsername            string       `json:"providerusername"`
	SchedulingName              string       `json:"schedulingname"`
	Sex                         string       `json:"sex"`
	Specialty                   string       `json:"specialty"`
	SpecialtyID                 NumberString `json:"specialtyid"`
	SupervisingProviderID       int          `json:"supervisingproviderid"`
	SupervisingProviderUsername string       `json:"supervisingproviderusername"`
}

// GetProvider - Get information for the specific provider
//
// GET /v1/{practiceid}/providers/{providerid}
//
// https://docs.athenahealth.com/api/api-ref/provider#Get-information-of-given-provider
func (h *HTTPClient) GetProvider(ctx context.Context, id string) (*Provider, error) {
	out := []*Provider{}

	_, err := h.Get(ctx, fmt.Sprintf("/providers/%s", id), nil, &out)
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, errors.New("Unexpected length returned")
	}

	return out[0], nil
}

type ListChangedProviderOptions struct {
	LeaveUnprocessed           bool
	ShowProcessedEndDatetime   time.Time
	ShowProcessedStartDatetime time.Time
}

type listChangedProvidersResponse struct {
	ChangedProviders []*Provider `json:"providers"`
}

// ListChangedProviders - A list of all changes to providers
//
// GET /v1/{practiceid}/providers/changed
//
// https://docs.athenahealth.com/api/api-ref/provider#Get-list-of-changes-in-providers
func (h *HTTPClient) ListChangedProviders(ctx context.Context, opts *ListChangedProviderOptions) ([]*Provider, error) {
	out := &listChangedProvidersResponse{}

	q := url.Values{}

	if opts != nil {
		if opts.LeaveUnprocessed {
			q.Add("leaveunprocessed", strconv.FormatBool(opts.LeaveUnprocessed))
		}

		if !opts.ShowProcessedEndDatetime.IsZero() {
			q.Add("showprocessedenddatetime", opts.ShowProcessedEndDatetime.Format("01/02/2006 15:04:05"))
		}

		if !opts.ShowProcessedStartDatetime.IsZero() {
			q.Add("showprocessedstartdatetime", opts.ShowProcessedStartDatetime.Format("01/02/2006 15:04:05"))
		}
	}

	_, err := h.Get(ctx, "/providers/changed", q, out)
	if err != nil {
		return nil, err
	}

	return out.ChangedProviders, nil
}

type ListProvidersOptions struct {
	ShowAllProviderIDs bool

	Pagination *PaginationOptions
}

type ListProvidersResult struct {
	Providers []*Provider

	Pagination *PaginationResult
}

type ListProvidersResponse struct {
	Providers []*Provider `json:"providers"`

	PaginationResponse
}

// ListProviders - List all providers available in a specific practice
//
// GET /v1/{practiceid}/providers
//
// https://docs.athenahealth.com/api/api-ref/provider-reference#Get-list-of-all-providers
func (h *HTTPClient) ListProviders(ctx context.Context, opts *ListProvidersOptions) (*ListProvidersResult, error) {
	out := &ListProvidersResponse{}

	q := url.Values{}

	if opts != nil {
		if opts.ShowAllProviderIDs {
			q.Add("showallproviderids", "true")
		}

		if opts.Pagination != nil {
			if opts.Pagination.Limit > 0 {
				q.Add("limit", strconv.Itoa(opts.Pagination.Limit))
			}

			if opts.Pagination.Offset > 0 {
				q.Add("offset", strconv.Itoa(opts.Pagination.Offset))
			}
		}
	}

	_, err := h.Get(ctx, "/providers", q, out)
	if err != nil {
		return nil, err
	}

	return &ListProvidersResult{
		Providers:  out.Providers,
		Pagination: makePaginationResult(out.Next, out.Previous, out.TotalCount),
	}, nil
}
