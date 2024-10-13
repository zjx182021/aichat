package server

import (
	"context"
	"keywords_filter/pkg/filter"
	"keywords_filter/proto/proto"
)

type filterService struct {
	proto.UnimplementedFilterServer
	filer filter.IFilter
}

func NewFilterService(filer filter.IFilter) proto.FilterServer {
	return &filterService{
		filer: filer,
	}

}

func (f *filterService) Validate(_ context.Context, in *proto.FilterReq) (*proto.ValidateRes, error) {
	ok, words := f.filer.Validate(in.Text)
	return &proto.ValidateRes{
		Ok:      ok,
		Keyword: words,
	}, nil
}
func (f *filterService) FindAll(_ context.Context, in *proto.FilterReq) (*proto.FindAllRes, error) {
	words := f.filer.FindAll(in.Text)
	return &proto.FindAllRes{
		Keywords: words,
	}, nil
}
