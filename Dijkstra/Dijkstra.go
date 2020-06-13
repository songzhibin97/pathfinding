/******
** @创建时间 : 2020/6/13 11:09
** @作者 : SongZhiBin
******/
package Dijkstra

import (
	"fmt"
	"github.com/tinygo-org/tinygo/src/os"
	"math"
)

// 迪克斯特拉算法
/*
1、找到“最便宜”的节点（可在最短时间内到达的节点）。
2、更新该节点的邻居节点的开销。
3、重复这个过程，直到对图中每个节点都做了。
4、计算最终路径。
*/

// drawing 使用二维数组表示加权有向图
type drawings [][]int

// 点属性
type point struct {
	// num表示node编号 比如1,2,3代表a,b,c 确切的数为x轴坐标
	num        int
	fatherNode *point
	minPower   int
	// 与上一个节点相聚的位置
	distance int
}

var (
	Drawings drawings
	MinPath  map[int]*point
)

const (
	// node表示顶点有4个
	// 初始化图会生成 node * node的二维矩阵
	node = 4
	// 无穷大 用于表示start->end未知的距离
	Infinity = math.MaxInt64
)

// 初始化有向图
func init() {
	Drawings = make(drawings, node)
	MinPath = make(map[int]*point)
	for i := 0; i < node; i++ {
		Drawings[i] = make([]int, node)
		MinPath[i] = &point{
			num:      i,
			minPower: Infinity,
		}
	}
	for i := 0; i < node; i++ {
		for j := 0; j < node; j++ {
			// 将所有矩阵内的数字设置为 最大值 1<<63 - 1
			Drawings[i][j] = Infinity
		}
	}
}

// 添加有向图关系

func (d *drawings) AddConnect(src, des, power int) bool {
	// src -> des 有向指向
	// power 权位
	// 判断src以及des是否合法
	if src < 0 || src >= node || des < 0 || des >= node {
		fmt.Fprintln(os.Stdout, "src or des not value")
		return false
	}
	(*d)[src][des] = power
	return true
}

func (d *drawings) Run(startI int) {
	start := &point{
		num:      startI,
		minPower: 0,
	}
	d.core(start)
	for key, value := range MinPath {
		fmt.Fprintf(os.Stdout, "key:%v father:%#v\n", key, value)
	}
}

// 核心
func (d *drawings) core(start *point) {
	// 创建两个slice分别用于存放 `可以访问的node` 和 `已经访问过and不能访问的路径`
	opens := []*point{start}
	closes := make([]*point, 0)
	// 此时opens已经包含start节点
	for len(opens) > 0 {
		now := opens[0]
		// 广度优先搜索
		d.scoutPath(now, &opens, &closes)
		// 删除 now
		d.deleteSlice(0, &opens)
		// 添加到close中
		closes = append(closes, now)
		// 排序 按与now节点 *** 一定要排序
		d.quickSort(0, len(opens)-1, &opens)
	}
}

// 探路
func (d *drawings) scoutPath(now *point, opens, closes *[]*point) {
	// 搜索
	for i := 0; i < node; i++ {
		if Drawings[now.num][i] == Infinity {
			// 代表无路
			continue
		}
		// 判断now节点是否已经在close中
		_, ok := d.isInSlice(now, closes)
		if ok {
			continue
		}
		// 判断更新子节点最小的路径以及父节点
		// now.num 是 y对应的 i是x
		if Drawings[now.num][i]+now.minPower < MinPath[i].minPower {
			// 如果小于 更新minPower and 父节点
			MinPath[i].minPower = Drawings[now.num][i] + now.minPower
			MinPath[i].fatherNode = now
		}
		MinPath[i].distance = Drawings[now.num][i]
		// 将该子节点加入到opens中
		*opens = append(*opens, MinPath[i])
	}
}

// 判断节点是否在某一个slice中
func (d *drawings) isInSlice(node *point, slices *[]*point) (int, bool) {
	for index, slice := range *slices {
		if node.num == slice.num {
			return index, true
		}
	}
	return 0, false
}

// 删除元素
func (d *drawings) deleteSlice(index int, opens *[]*point) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	if index < 0 {
		return
	}
	if index == 0 {
		newOpens := make([]*point, len(*opens)-1)
		copy(newOpens, (*opens)[1:])
		*opens = newOpens
		return
	}
	if index == len(*opens) {
		newOpens := make([]*point, len(*opens)-1)
		copy(newOpens, (*opens)[:len(*opens)-1])
		*opens = newOpens
		return
	}
	*opens = (*opens)[:copy(*opens, (*opens)[index:])]
}

// 快排核心
func (d *drawings) sortCore(start, end int, opens *[]*point) int {
	mid := (*opens)[start]
	for start < end {
		for start < end && (*opens)[end].distance > mid.distance {
			end--
		}
		if start < end {
			(*opens)[start] = (*opens)[end]
			start++
		}
		for start < end && (*opens)[start].distance < mid.distance {
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
func (d *drawings) quickSort(start, end int, opens *[]*point) {
	if start < end {
		mid := d.sortCore(start, end, opens)
		d.quickSort(start, mid-1, opens)
		d.quickSort(mid+1, end, opens)
	}

}
