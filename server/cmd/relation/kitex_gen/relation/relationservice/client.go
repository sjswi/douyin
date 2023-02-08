// Code generated by Kitex v0.4.4. DO NOT EDIT.

package relationservice

import (
	"context"
	relation "douyin_rpc/server/cmd/relation/kitex_gen/relation"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	Action(ctx context.Context, req *relation.RelationActionRequest, callOptions ...callopt.Option) (r *relation.RelationActionResponse, err error)
	FollowList(ctx context.Context, req *relation.RelationFollowListRequest, callOptions ...callopt.Option) (r *relation.RelationFollowListResponse, err error)
	FollowerList(ctx context.Context, req *relation.RelationFollowerListRequest, callOptions ...callopt.Option) (r *relation.RelationFollowerListResponse, err error)
	FriendList(ctx context.Context, req *relation.RelationFriendListRequest, callOptions ...callopt.Option) (r *relation.RelationFriendListResponse, err error)
	GetRelation(ctx context.Context, req *relation.GetRelationRequest, callOptions ...callopt.Option) (r *relation.GetRelationResponse, err error)
	GetCount(ctx context.Context, req *relation.GetCountRequest, callOptions ...callopt.Option) (r *relation.GetCountResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kRelationServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kRelationServiceClient struct {
	*kClient
}

func (p *kRelationServiceClient) Action(ctx context.Context, req *relation.RelationActionRequest, callOptions ...callopt.Option) (r *relation.RelationActionResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.Action(ctx, req)
}

func (p *kRelationServiceClient) FollowList(ctx context.Context, req *relation.RelationFollowListRequest, callOptions ...callopt.Option) (r *relation.RelationFollowListResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.FollowList(ctx, req)
}

func (p *kRelationServiceClient) FollowerList(ctx context.Context, req *relation.RelationFollowerListRequest, callOptions ...callopt.Option) (r *relation.RelationFollowerListResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.FollowerList(ctx, req)
}

func (p *kRelationServiceClient) FriendList(ctx context.Context, req *relation.RelationFriendListRequest, callOptions ...callopt.Option) (r *relation.RelationFriendListResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.FriendList(ctx, req)
}

func (p *kRelationServiceClient) GetRelation(ctx context.Context, req *relation.GetRelationRequest, callOptions ...callopt.Option) (r *relation.GetRelationResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetRelation(ctx, req)
}

func (p *kRelationServiceClient) GetCount(ctx context.Context, req *relation.GetCountRequest, callOptions ...callopt.Option) (r *relation.GetCountResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetCount(ctx, req)
}
