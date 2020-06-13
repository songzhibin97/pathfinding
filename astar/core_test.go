/******
** @创建时间 : 2020/6/13 10:57
** @作者 : SongZhiBin
******/
package astar

import (
	"reflect"
	"testing"
)

func Test_matrix_Run(t *testing.T) {
	// todo 因为结果是根据 AddObstacles动态生成障碍点 所以不能使用测试组
	//Matrix.Init()
	// 添加新的障碍点
	Matrix.AddObstacles(1, 0, 2)
	Matrix.AddObstacles(1, 1, 2)
	Matrix.AddObstacles(1, 2, 2)
	got, ok := Matrix.Run(0, 0, 0, 0)
	if !ok {
		return
	}
	want := []*point{&point{X: 0, Y: 0}}
	if !reflect.DeepEqual(want,got) {
		 t.Errorf("excepted:%v, got:%v", want, got) // 测试失败输出错误提示
	}
}
