package hechodoc

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestOpenApiPathFromEchoPath(t *testing.T) {
	data:= []struct {
		Tag string
		EchoPath string
		OAIPath string
	}{
		{"t1","/","/"},
		{"t2","a","a"},
		{"t3","a/:id","a/{id}"},
		{"t4","a/:id/:code","a/{id}/{code}"},
		{"t5","a/:id/:code/","a/{id}/{code}/"},
		{"t6","a/b/:c/:d/:e/f/:g/h","a/b/{c}/{d}/{e}/f/{g}/h"},
	}

	for _,d:=range data{
		assert.Equal(t,OpenApiPathFromEchoPath(d.EchoPath),d.OAIPath,d.Tag)
	}
}
