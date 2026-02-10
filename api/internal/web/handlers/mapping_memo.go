package handlers

import (
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/isutare412/web-memo/api/internal/core/ent"
	"github.com/isutare412/web-memo/api/internal/core/model"
	"github.com/isutare412/web-memo/api/internal/web/gen"
)

// MemoToWeb converts an ent.Memo and its associated tags to the generated Memo type.
func MemoToWeb(memo *ent.Memo, tags []string) gen.Memo {
	return gen.Memo{
		ID:          memo.ID,
		OwnerID:     memo.OwnerID,
		CreateTime:  memo.CreateTime,
		UpdateTime:  memo.UpdateTime,
		Title:       memo.Title,
		Content:     memo.Content,
		IsPublished: memo.IsPublished,
		Version:     memo.Version,
		Tags:        tags,
	}
}

// MemoSearchResultToWeb converts a MemoSearchResult to the generated Memo type,
// including search relevance scores.
func MemoSearchResultToWeb(result *model.MemoSearchResult) gen.Memo {
	m := MemoToWeb(result.Memo, memoTagNames(result.Memo))
	m.Scores = &gen.MemoScores{
		Rrf:      result.RRFScore,
		Semantic: result.SemanticScore,
		Bm25:     result.BM25Score,
	}
	return m
}

// MemosToListResponse converts a list of ent.Memo with per-memo tags and pagination
// metadata to a generated ListMemosResponse.
func MemosToListResponse(memos []*ent.Memo, tagsByMemo [][]string, page, pageSize, lastPage, totalCount int) gen.ListMemosResponse {
	webMemos := make([]gen.Memo, 0, len(memos))
	for i, m := range memos {
		var tags []string
		if i < len(tagsByMemo) {
			tags = tagsByMemo[i]
		}
		webMemos = append(webMemos, MemoToWeb(m, tags))
	}

	return gen.ListMemosResponse{
		Page:           lo.ToPtr(page),
		PageSize:       lo.ToPtr(pageSize),
		LastPage:       lo.ToPtr(lastPage),
		TotalMemoCount: lo.ToPtr(totalCount),
		Memos:          webMemos,
	}
}

// SearchResultsToListResponse converts search results to a generated
// ListMemosResponse without pagination metadata (fields left nil).
func SearchResultsToListResponse(results []*model.MemoSearchResult) gen.ListMemosResponse {
	return gen.ListMemosResponse{
		Memos: lo.Map(results, func(r *model.MemoSearchResult, _ int) gen.Memo {
			return MemoSearchResultToWeb(r)
		}),
	}
}

// SubscribersToWeb converts a model.ListSubscribersResponse to the generated
// ListSubscribersResponse.
func SubscribersToWeb(resp *model.ListSubscribersResponse) gen.ListSubscribersResponse {
	return gen.ListSubscribersResponse{
		MemoOwnerID: resp.MemoOwnerID,
		Subscribers: lo.Map(resp.Subscribers, func(u *ent.User, _ int) gen.Subscriber {
			return gen.Subscriber{ID: u.ID}
		}),
	}
}

// CollaboratorsToWeb converts a model.ListCollaboratorsResponse to the generated
// ListCollaboratorsResponse. Each user's Collaborations edge must be loaded so that
// the approval status can be determined.
func CollaboratorsToWeb(resp *model.ListCollaboratorsResponse, memoID uuid.UUID) gen.ListCollaboratorsResponse {
	collaborators := make([]gen.Collaborator, 0, len(resp.Collaborators))
	for _, user := range resp.Collaborators {
		for _, collabo := range user.Edges.Collaborations {
			if collabo.MemoID != memoID {
				continue
			}
			collaborators = append(collaborators, gen.Collaborator{
				ID:         user.ID,
				UserName:   user.UserName,
				PhotoURL:   user.PhotoURL,
				IsApproved: collabo.Approved,
			})
		}
	}

	return gen.ListCollaboratorsResponse{
		MemoOwnerID:   resp.MemoOwnerID,
		Collaborators: collaborators,
	}
}

// memoTagNames extracts tag names from a memo's eager-loaded Tags edge.
func memoTagNames(memo *ent.Memo) []string {
	return lo.Map(memo.Edges.Tags, func(t *ent.Tag, _ int) string {
		return t.Name
	})
}
