/******
** @创建时间 : 2020/6/10 11:00
** @作者 : SongZhiBin
******/
package astar

import (
	"fmt"
	"math"
	"os"
	"sync"
)

// core
// Applied to two dimensional matrices
/*
a、将开始点记录为当前点P
b、将当前点P放入封闭列表
c、搜寻点P所有邻近点，假如某邻近点既没有在开放列表或封闭列表里面，则计算出该邻近点的F值，并设父节点为P，然后将其放入开放列表
d、判断开放列表是否已经空了，如果没有说明在达到结束点前已经找完了所有可能的路径点，寻路失败，算法结束；否则继续。
e、从开放列表拿出一个F值最小的点，作为寻路路径的下一步。
f、判断该点是否为结束点，如果是，则寻路成功，算法结束；否则继续。
g、将该点设为当前点P，跳回步骤c。
*/

// Matrix为全属性地图
type matrix [][]int

var Matrix matrix

const (
	// X轴
	XAxis = 4
	// Y轴
	YAxis = 4
)

// 初始化地图
func init() {
	Matrix = make([][]int, YAxis)
	for i, _ := range Matrix {
		(Matrix)[i] = make([]int, XAxis)
	}
}

// 主函数
func (m *matrix) Run(startX, startY, endX, endY int) ([]*point, bool) {
	// 参数为坐标 自己生产*point对象
	start := &point{
		X: startX,
		Y: startY,
	}
	ok := m.isExist(start)
	if !ok {
		fmt.Fprintln(os.Stdout, "start no exist")
		return nil, false
	}
	end := &point{
		X: endX,
		Y: endY,
	}
	ok = m.isExist(end)
	if !ok {
		fmt.Fprintln(os.Stdout, "end no exist")
		return nil, false
	}
	endPoint, ok := m.aStar(start, end)
	if !ok {
		return nil, false
	}
	// 迭代
	ret := make([]*point, 0)
	for {
		if endPoint == nil {
			m.reverse(&ret)
			return ret, true
		}
		//fmt.Printf("(%d,%d)\n", endPoint.X, endPoint.Y)
		ret = append(ret, endPoint)
		endPoint = endPoint.fatherPoint
	}
}

// 反转
func (m *matrix) reverse(slices *[]*point) {
	start, end := 0, len(*slices)-1
	for start < end {
		(*slices)[start], (*slices)[end] = (*slices)[end], (*slices)[start]
		start++
		end--
	}
}

// 核心 返回路径和是否找到路径
func (m *matrix) aStar(start, end *point) (*point, bool) {
	// 传入起始位置和预计要去的终点
	// 校验start,end中x,y是否有效
	ok := m.isValid(start)
	if !ok {
		fmt.Fprintln(os.Stdout, "start is not valid")
		return nil, false
	}
	ok = m.isValid(end)
	if !ok {
		fmt.Fprintln(os.Stdout, "end is not valid")
		return nil, false
	}

	// 这里进行一个小判断 优化程序 判断如果statPoint == endPoint 直接跳出
	if start.X == end.X && start.Y == end.Y {
		return start, true
	}

	// 为了避免变量污染让此函数可以并发的执行 所以这里用了局部变量

	// 创建两个slice分别用于存放 `可以访问的node` 和 `已经访问过and不能访问的路径`
	opens := make([]*point, 0)
	closes := make([]*point, 0)

	// a、将开始点记录为当前点P
	// 直接将起点放到已走过的位置
	// b、将当前点P放入封闭列表
	opens = append(opens, start)
	now := opens[0]
	for {
		// c、搜寻点P所有邻近点，假如某邻近点既没有在开放列表或封闭列表里面，则计算出该邻近点的F值，并设父节点为P，然后将其放入开放列表
		m.scoutPath(now, end, &opens, &closes)
		// 搜索完后 删除取出的首节点 加入到closes对列中
		index, ok := m.isInSlice(now, &opens)
		if !ok {
			return nil, false
		}
		m.deleteSlice(index, &opens)
		closes = append(closes, now)
		// d.判断开放列表是否已经空了，如果没有说明在达到结束点前已经找完了所有可能的路径点，寻路失败，算法结束；否则继续。
		if len(opens) <= 0 {
			// 如果为空则跳出
			fmt.Fprintln(os.Stdout, "no find path")
			return nil, false
		}
		// 这里进行排序
		m.quickSort(0, len(opens)-1, &opens)

		// e.从开放列表拿出一个F值最小的点，作为寻路路径的下一步。
		next := opens[0]
		// f.判断该点是否为结束点，如果是，则寻路成功，算法结束；否则继续。
		if next.X == end.X && next.Y == end.Y {
			// 寻路成功
			return next, true
		}

		// g、将该点设为当前点P，跳回步骤c。
		now = next
	}
}

// 删除元素
func (m *matrix) deleteSlice(index int, opens *[]*point) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	if index < 0 {
		return
	}
	newOpens := make([]*point, len(*opens)-1)
	if index == 0 {
		copy(newOpens, (*opens)[1:])
		*opens = newOpens
		return
	}
	if index == len(*opens) {
		copy(newOpens, (*opens)[:len(*opens)-1])
		*opens = newOpens
		return
	}
	copy(newOpens, (*opens)[:index])
	copy(newOpens[index:], (*opens)[index+1:])
	*opens = newOpens
}

// 探路
func (m *matrix) scoutPath(now, end *point, opens, closes *[]*point) {
	// 顺序:上->右->下->左
	// 上 x=now.x y=now.y-1
	upNode := m.nextNode(now.X, now.Y-1, now, end)
	rightNode := m.nextNode(now.X+1, now.Y, now, end)
	downNode := m.nextNode(now.X, now.Y+1, now, end)
	leftNode := m.nextNode(now.X-1, now.Y, now, end)
	m.choiceSlice(now, end, opens, closes, &sync.Mutex{}, upNode, rightNode, downNode, leftNode)
}

// 用于探路增加新节点
func (m *matrix) nextNode(x, y int, now, end *point) *point {
	res := &point{
		X:           x,
		Y:           y,
		g:           now.g + 1,
		h:           int(math.Abs(float64(end.X-now.X)) + math.Abs(float64(end.Y-now.Y))), // abs(dx-x)+abs(dy-y)
		fatherPoint: now,
	}
	res.f = res.g + res.h
	return res
}

// 用于探路创建新节点判断放入opens和close两个通道之一
func (m *matrix) choiceSlice(now, end *point, opens, closes *[]*point, lock *sync.Mutex, points ...*point) {
	for _, point := range points {
		// 判断地图是否有该节点
		ok := m.isExist(point)
		// 存在则进行下一步判断
		if !ok {
			continue
		}
		// 这里判断该节点是否已经在closes中
		_, ok = m.isInSlice(point, closes)

		if ok {
			continue
		}
		lock.Lock()
		if (*m)[point.Y][point.X] == 0 {
			// 将其放到 opens
			index, ok := m.isInSlice(point, opens)
			// 判断是否存在于opens中
			if ok {
				// 检查f值是否更小符合更新条件
				if (*opens)[index].f > now.g+1+int(math.Abs(float64(end.X-now.X))+math.Abs(float64(end.Y-now.Y))) {
					continue
				}
				// 已经存在进行更新操作 更新g,h,f
				(*opens)[index].g = now.g + 1
				(*opens)[index].h = int(math.Abs(float64(end.X-now.X)) + math.Abs(float64(end.Y-now.Y)))
				(*opens)[index].f = (*opens)[index].g + (*opens)[index].h
				(*opens)[index].fatherPoint = now
			} else {
				*opens = append(*opens, point)
			}
		} else {
			*closes = append(*closes, point)
		}
		lock.Unlock()
	}
}

// 判断节点是否在某一个slice中
func (m *matrix) isInSlice(node *point, slices *[]*point) (int, bool) {
	for index, slice := range *slices {
		if node.X == slice.X && node.Y == slice.Y {
			return index, true
		}
	}
	return 0, false
}

// 判断point节点是否有效
func (m *matrix) isValid(point *point) bool {
	if point.X < 0 || point.X >= XAxis || point.Y < 0 || point.Y >= YAxis || (*m)[point.Y][point.X] != 0 {
		return false
	}
	return true
}

// 判断point节点是否有
func (m *matrix) isExist(point *point) bool {
	if point.X < 0 || point.X >= XAxis || point.Y < 0 || point.Y >= YAxis {
		return false
	}
	return true
}

// node设置
func (m *matrix) AddObstacles(x, y int, tag int) bool {
	// 判断x,y是否合法
	if x < 0 || y < 0 || x >= XAxis || y >= YAxis {
		return false
	}
	Matrix[y][x] = tag
	return true
}

// 快排核心
func (m *matrix) sortCore(start, end int, opens *[]*point) int {
	mid := (*opens)[start]
	for start < end {
		for start < end && (*opens)[end].f > mid.f {
			end--
		}
		if start < end {
			(*opens)[start] = (*opens)[end]
			start++
		}
		for start < end && (*opens)[start].f < mid.f {
			start++
		}
		if start < end {
			(*opens)[end] = (*opens)[start]
			end--
		}
	}
	(*opens)[start] = mid
	return start
}

// 快排
func (m *matrix) quickSort(start, end int, opens *[]*point) {
	if start < end {
		mid := m.sortCore(start, end, opens)
		m.quickSort(start, mid-1, opens)
		m.quickSort(mid+1, end, opens)
	}

}

// 地图demo
/*
0:是可通过的路径
1:是不可通过的障碍物
5:是起点
9:是终点
2:已走过的
3:临时占位-> 预留位

 0 - 1 - 2 - 3 - 4 - 5 - 6
	--------------------------
 1	2 | 1 | 1 | 0 | 0 | 0 |
	--------------------------
 2	2 | p | 1 | 0 | 2 | 9 |
	--------------------------
 3	2 | 1 | 1 | 0 | 2 | 0 |
	--------------------------
 4	2 | 1 | 1 | 1 | 2 | 0 |
	--------------------------
 5	2 | 2 | 2 | 2 | 2 | 1 |
	--------------------------
(1,1)->(1,2)->(1,3)->(1,4)->(1,5)->(2,5)->(3,5)->(4,5)->(5,5)->(5,4)->(5,3)->(5,2)->(6,2)
*/

// 点属性
type point struct {
	// X,Y 为坐标点
	// G值  是从起点到某一点积累的移动值。一般我们将从起点按非斜向方向移动的G值定为10，
	// 斜向为20，走一步必须将此节点的值增加，也就是说移动到的节点的G值等于移动前节点的G值加上按方向走到该店的G值的增加量。
	// 打个比方，移动前(0,0)的G值为0，则(0,1)的G值就是10，(1,2)的G值就是30

	// H值
	// 该值和终点坐标相关，一般常用的曼哈顿算法为: abs(dx-x)+abs(dy-y)，也就是横坐标的差值加上纵坐标的差值。

	// F值
	// F值为G值和H值之和
	X, Y        int
	g, h, f     int
	fatherPoint *point
}
