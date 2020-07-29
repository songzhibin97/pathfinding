/******
** @创建时间 : 2020/6/13 10:57
** @作者 : SongZhiBin
******/
package astar

import (
	"fmt"
	"testing"
)

func Test_matrix_Run(t *testing.T) {
	// todo 因为结果是根据 AddObstacles动态生成障碍点 所以不能使用测试组
	//Matrix.Init()
	// 添加新的障碍点
	Matrix.AddObstacles(1, 0, 2)
	Matrix.AddObstacles(1, 1, 2)
	Matrix.AddObstacles(1, 3, 2)

	got, ok := Matrix.Run(0, 0, 2, 0)
	fmt.Println(ok)
	for _, g := range got {
		fmt.Println(g)
	}

}
