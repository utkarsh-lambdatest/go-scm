// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"

	"github.com/drone/go-scm/scm"
)

type organizationService struct {
	client *wrapper
}

func (s *organizationService) Find(ctx context.Context, name string) (*scm.Organization, *scm.Response, error) {
	path := fmt.Sprintf("orgs/%s", name)
	out := new(organization)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertOrganization(out), res, err
}

func (s *organizationService) FindMembership(ctx context.Context, name, username string) (*scm.Membership, *scm.Response, error) {
	path := fmt.Sprintf("orgs/%s/memberships/%s", name, username)
	out := new(membership)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertMembership(out), res, err
}

func (s *organizationService) ListMemberships(ctx context.Context, orgNameList []string, username string, opts scm.ListOptions) ([]*scm.Membership, *scm.Response, error) {
	path := fmt.Sprintf("/user/memberships/orgs?%s", encodeListOptions(opts))
	out := []*membership{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	orgMemberships := convertMemberships(out)

	usernameMembership := &scm.Membership{
		Role:   scm.RoleAdmin,
		Active: true,
		Organization: scm.Organization{
			Name:   username,
			Avatar: "",
		},
	}
	orgMemberships = append(orgMemberships, usernameMembership)

	return orgMemberships, res, err
}

func (s *organizationService) List(ctx context.Context, opts scm.ListOptions) ([]*scm.Organization, *scm.Response, error) {
	path := fmt.Sprintf("user/orgs?%s", encodeListOptions(opts))
	out := []*organization{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertOrganizationList(out), res, err
}

func convertOrganizationList(from []*organization) []*scm.Organization {
	to := []*scm.Organization{}
	for _, v := range from {
		to = append(to, convertOrganization(v))
	}
	return to
}

type organization struct {
	Login  string `json:"login"`
	Avatar string `json:"avatar_url"`
}

type membership struct {
	State        string       `json:"state"`
	Role         string       `json:"role"`
	Organization organization `json:"organization,omitempty"`
}

func convertOrganization(from *organization) *scm.Organization {
	return &scm.Organization{
		Name:   from.Login,
		Avatar: from.Avatar,
	}
}

func convertMemberships(from []*membership) []*scm.Membership {
	to := []*scm.Membership{}
	for _, v := range from {
		to = append(to, convertMembership(v))
	}
	return to
}

func convertMembership(from *membership) *scm.Membership {
	to := new(scm.Membership)
	if from.State == "active" {
		to.Active = true
	}
	switch from.Role {
	case "admin":
		to.Role = scm.RoleAdmin
	case "member":
		to.Role = scm.RoleMember
	default:
		to.Role = scm.RoleViewer
	}

	to.Organization.Name = from.Organization.Login
	to.Organization.Avatar = from.Organization.Avatar
	return to
}
