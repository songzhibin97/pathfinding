/******
** @创建时间 : 2020/6/10 11:01
** @作者 : SongZhiBin
******/
package astar

// interface

type Grid interface {

	// 用于初始化地图之后添加障碍点
	AddObstacles(x, y int, tag int) bool

	// 分别传入 startX,starY,endX,endY 返回具体的路径列表 和 是否成功
	Run(startX, starY, endX, endY int) ([]*point, bool)
}
