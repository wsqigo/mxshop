package handler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func TestGetAllCategoryList(t *testing.T) {
	rsp, err := client.GetAllCategoryList(context.Background(), &emptypb.Empty{})
	assert.Nil(t, err)
	t.Log(rsp.Total)
	for _, category := range rsp.Data {
		t.Log(category.Name)
	}
}
