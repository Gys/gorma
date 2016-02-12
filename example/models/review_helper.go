//************************************************************************//
// API "congo": Model Helpers
//
// Generated with goagen v0.0.1, command line:
// $ goagen
// --out=$(GOPATH)/src/github.com/goadesign/gorma/example
// --design=github.com/goadesign/gorma/example/design
//
// The content of this file is auto-generated, DO NOT MODIFY
//************************************************************************//

package models

import (
	"github.com/goadesign/goa"
	"github.com/goadesign/gorma/example/app"
	"github.com/jinzhu/gorm"
	"time"
)

// MediaType Retrieval Functions
// ListReview returns an array of view: default
func (m *ReviewDB) ListAppReview(ctx *goa.Context, proposalid int, userid int) []*app.Review {
	now := time.Now()
	defer ctx.Info("ListReview", "duration", time.Since(now))
	var objs []*app.Review
	err := m.Db.Scopes(ReviewFilterByProposal(proposalid, &m.Db), ReviewFilterByUser(userid, &m.Db)).Table(m.TableName()).Find(&objs).Error

	//	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil {
		ctx.Error("error listing Review", "error", err.Error())
		return objs
	}

	return objs
}

func (m *Review) ReviewToAppReview() *app.Review {
	review := &app.Review{}
	review.ID = &m.ID
	review.Comment = m.Comment
	review.Rating = &m.Rating

	return review
}

// OneAppReview returns an array of view: default
func (m *ReviewDB) OneReview(ctx *goa.Context, id int, proposalid int, userid int) (*app.Review, error) {
	now := time.Now()
	var native Review
	defer ctx.Info("OneReview", "duration", time.Since(now))
	err := m.Db.Scopes(ReviewFilterByProposal(proposalid, &m.Db), ReviewFilterByUser(userid, &m.Db)).Table(m.TableName()).Preload("Proposal").Preload("User").Where("id = ?", id).Find(&native).Error

	if err != nil && err != gorm.RecordNotFound {
		ctx.Error("error getting Review", "error", err.Error())
		return nil, err
	}

	view := *native.ReviewToAppReview()
	return &view, err

}

// MediaType Retrieval Functions
// ListReviewLink returns an array of view: link
func (m *ReviewDB) ListAppReviewLink(ctx *goa.Context, proposalid int, userid int) []*app.ReviewLink {
	now := time.Now()
	defer ctx.Info("ListReviewLink", "duration", time.Since(now))
	var objs []*app.ReviewLink
	err := m.Db.Scopes(ReviewFilterByProposal(proposalid, &m.Db), ReviewFilterByUser(userid, &m.Db)).Table(m.TableName()).Find(&objs).Error

	//	err := m.Db.Table(m.TableName()).Find(&objs).Error
	if err != nil {
		ctx.Error("error listing Review", "error", err.Error())
		return objs
	}

	return objs
}

func (m *Review) ReviewToAppReviewLink() *app.ReviewLink {
	review := &app.ReviewLink{}
	review.ID = &m.ID

	return review
}

// OneAppReviewLink returns an array of view: link
func (m *ReviewDB) OneReviewLink(ctx *goa.Context, id int, proposalid int, userid int) (*app.ReviewLink, error) {
	now := time.Now()
	var native Review
	defer ctx.Info("OneReviewLink", "duration", time.Since(now))
	err := m.Db.Scopes(ReviewFilterByProposal(proposalid, &m.Db), ReviewFilterByUser(userid, &m.Db)).Table(m.TableName()).Preload("Proposal").Preload("User").Where("id = ?", id).Find(&native).Error

	if err != nil && err != gorm.RecordNotFound {
		ctx.Error("error getting Review", "error", err.Error())
		return nil, err
	}

	view := *native.ReviewToAppReviewLink()
	return &view, err

}
