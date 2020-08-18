package model

import (
	"sort"
	"strings"

	"github.com/convox/console/pkg/license"
	"github.com/pkg/errors"
)

type Organization struct {
	ID string `dynamo:"id"`

	Administrators []string `dynamo:"administrator-ids"`
	Operators      []string `dynamo:"operator-ids"`
	Users          []string `dynamo:"user-ids"`

	Creator             string            `dynamo:"creator"`
	Locked              bool              `dynamo:"locked"`
	MaxUsers            int               `dynamo:"max-users"`
	Name                string            `dynamo:"name"`
	OverrideConcurrency int               `dynamo:"override-concurrency"`
	Plan                string            `dynamo:"plan"`
	Restrictions        map[string]string `dynamo:"restrictions"`
	StripeCustomer      string            `dynamo:"stripe-id"`
	StripeSubscription  string            `dynamo:"plan-subscription-id"`
}

type Organizations []Organization

func (m *Model) OrganizationGet(id string) (*Organization, error) {
	o := &Organization{}

	if err := m.storage.Get("organizations", id, o); err != nil {
		return nil, errors.WithStack(err)
	}

	return o, nil
}

func (m *Model) OrganizationIntegrations(oid string) (Integrations, error) {
	var is Integrations

	if err := m.storage.GetIndex("integrations", "organization-id-index", map[string]string{"organization-id": oid}, &is); err != nil {
		return nil, errors.WithStack(err)
	}

	return is, nil
}

func (m *Model) OrganizationRacks(oid string) (Racks, error) {
	var rs Racks

	if err := m.storage.GetIndex("racks", "organization-id-index", map[string]string{"organization-id": oid}, &rs); err != nil {
		return nil, errors.WithStack(err)
	}

	sort.Slice(rs, rs.Less)

	// if len(rs) > 0 {
	// 	rs = append(rs, rs[0], rs[0], rs[0], rs[0], rs[0])
	// }

	return rs, nil
}

func (m *Model) OrganizationSave(o *Organization) error {
	if err := m.storage.Put("organizations", o); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (o *Organization) JobConcurrency() int {
	if o.OverrideConcurrency > 0 {
		return o.OverrideConcurrency
	}

	if !license.Current.Public {
		return 1000
	}

	switch {
	case o.PlanPro():
		return 3
	case o.PlanEnterprise():
		return 5
	default:
		return 1
	}
}

func (o *Organization) PlanEnterprise() bool {
	return o.PlanEnterpriseDynamic() || o.PlanEnterpriseStatic()
}

func (o *Organization) PlanEnterpriseDynamic() bool {
	switch {
	case strings.HasPrefix(o.Plan, "enterprise-dynamic-"):
		return true
	default:
		return false
	}
}

func (o *Organization) PlanEnterpriseStatic() bool {
	switch o.Plan {
	case "enterprise":
		return true
	default:
		return false
	}
}

func (o *Organization) PlanPro() bool {
	switch o.Plan {
	case "pro", "pro-user-month", "pro-v3", "pro-v3-annual", "pro-v3-annual-full", "pro-v3-static":
		return true
	default:
		return false
	}
}

func (os Organizations) Less(i, j int) bool {
	return os[i].Name < os[j].Name
}
