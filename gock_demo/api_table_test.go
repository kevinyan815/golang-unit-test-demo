package gock_demo

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"testing"
)

func TestGetResultByAPI2(t *testing.T) {
	t.Parallel()     // 将 TLog 标记为能够与其他测试并行运行
	defer gock.Off() // 测试执行后刷新挂起的mock

	type args struct {
		mockURI        string
		path           string
		responseStatus int
		reply          int
		paramX         int
		paramY         int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "case1",
			args: args{
				mockURI:        "http://your-api.com",
				path:           "/post",
				responseStatus: 200,
				reply:          100,
				paramX:         1,
				paramY:         1,
			},
			want: 101,
		},
		{
			name: "case2",
			args: args{
				mockURI:        "http://your-api.com",
				path:           "/post",
				responseStatus: 200,
				reply:          200,
				paramX:         2,
				paramY:         2,
			},
			want: 202,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gock.New(tt.args.mockURI).
				Path(tt.args.path).
				MatchType("json").
				JSON(map[string]int{"x": tt.args.paramX}).
				Reply(200).
				JSON(map[string]int{"value": tt.args.reply})

			// 调用我们的业务函数
			res := GetResultByAPI(tt.args.paramX, tt.args.paramY)
			// 校验返回结果是否符合预期
			assert.Equal(t, res, tt.want)
			assert.True(t, gock.IsDone())

		})
	}
}
