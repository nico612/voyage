package post

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/nico612/go-project/examples/miniblog/internal/miniblog/store"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/errno"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/log"
	"github.com/nico612/go-project/examples/miniblog/internal/pkg/model"
	v1 "github.com/nico612/go-project/examples/miniblog/pkg/api/miniblog/v1"
	"gorm.io/gorm"
)

//go:generate mockgen -destination mock_post.go -package post github.com/nico612/go-project/examples/miniblog/internal/miniblog/biz/post PostBiz
type PostBiz interface {
	Create(ctx context.Context, username string, r *v1.CreatePostRequest) (*v1.CreatePostResponse, error)
	Update(ctx context.Context, username, postID string, r *v1.UpdatePostRequest) error
	Delete(ctx context.Context, username, postID string) error
	DeleteCollection(ctx context.Context, username string, postIDs []string) error
	Get(ctx context.Context, username, postID string) (*v1.GetPostResponse, error)
	List(ctx context.Context, username string, offset, limit int) (*v1.ListPostResponse, error)
}

type postBiz struct {
	ds store.IStore
}

var _ PostBiz = (*postBiz)(nil)

func New(ds store.IStore) *postBiz {
	return &postBiz{ds: ds}
}

func (p *postBiz) Create(ctx context.Context, username string, r *v1.CreatePostRequest) (*v1.CreatePostResponse, error) {
	var post model.PostM
	_ = copier.Copy(&post, r)
	post.Username = username

	// 创建后，grom 框架会把自动创建的数据回写到post中
	if err := p.ds.Posts().Create(ctx, &post); err != nil {
		return nil, err
	}
	return &v1.CreatePostResponse{
		PostID: post.PostID,
	}, nil
}

func (p *postBiz) Update(ctx context.Context, username, postID string, r *v1.UpdatePostRequest) error {
	post, err := p.ds.Posts().Get(ctx, username, postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.ErrPostNotFound
		}
		return err
	}

	if r.Title != nil {
		post.Title = *r.Title
	}

	if r.Content != nil {
		post.Content = *r.Title
	}

	return p.ds.Posts().Update(ctx, post)
}

// Delete is the implementation of the `Delete` method in PostBiz interface.
func (p *postBiz) Delete(ctx context.Context, username, postID string) error {
	if err := p.ds.Posts().Delete(ctx, username, []string{postID}); err != nil {
		return err
	}
	return nil
}

// DeleteCollection is the implementation of the `DeleteCollection` method in PostBiz interface.
func (p *postBiz) DeleteCollection(ctx context.Context, username string, postIDs []string) error {
	if err := p.ds.Posts().Delete(ctx, username, postIDs); err != nil {
		return err
	}
	return nil
}

func (p *postBiz) Get(ctx context.Context, username, postID string) (*v1.GetPostResponse, error) {
	post, err := p.ds.Posts().Get(ctx, username, postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrPostNotFound
		}
		return nil, err
	}
	var resp v1.GetPostResponse
	_ = copier.Copy(&resp, post)
	resp.CreatedAt = post.CreatedAt.Format("2006-01-02 15:04:05")
	resp.UpdatedAt = post.UpdatedAt.Format("2006-01-02 15:04:05")
	return &resp, nil
}

func (p *postBiz) List(ctx context.Context, username string, offset, limit int) (*v1.ListPostResponse, error) {
	total, list, err := p.ds.Posts().List(ctx, username, offset, limit)
	if err != nil {
		log.C(ctx).Errorw("Failed to list posts from storage", "err", err)
		return nil, err
	}

	posts := make([]*v1.PostInfo, 0, len(list))
	for _, item := range list {

		posts = append(posts, &v1.PostInfo{
			Username:  item.Username,
			PostID:    item.PostID,
			Title:     item.Title,
			Content:   item.Content,
			CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: item.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &v1.ListPostResponse{TotalCount: total, Posts: posts}, nil
}
