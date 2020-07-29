# Dijkstra
```go
package main

import (
	"Songzhibin/test/pathfinding/Dijkstra"
)

func main() {
  	// 添加一些向量
  	// AddConnect(src,des,pow)
  	// src -> des
	Dijkstra.Drawings.AddConnect(0, 1, 6)
	Dijkstra.Drawings.AddConnect(0, 2, 2)
	Dijkstra.Drawings.AddConnect(2, 1, 3)
	Dijkstra.Drawings.AddConnect(1, 3, 1)
	Dijkstra.Drawings.AddConnect(2, 3, 5)
  	// Run(start) 启动 以start为点寻找到各个点权重最小的路径and权位 
	Dijkstra.Drawings.Run(0)
}
```

# aStar
```go
 func main() {
	// 添加新的障碍点
  	// 这里是二维矩阵
  	// AddObstacles(x,y,tag)
  	// 分别对应x,y坐标以及设置的tag
	Matrix.AddObstacles(1, 0, 2)
	Matrix.AddObstacles(1, 1, 2)
	Matrix.AddObstacles(1, 2, 2)
  	// Run(startX,startY,endX,endY)
  	// 分别对应起点终点坐标
	Matrix.Run(0, 0, 0, 0)
}
```
#### aStar定制性很强,你可以根据不同的业务去设置一些你自己的逻辑 例如
```
// 地图demo
/*
0:是可通过的路径
1:是不可通过的障碍物
5:是起点
9:是终点
2:已走过的
3:临时占位-> 预留位

    0 	- 1 - 2 - 3 - 4 - 5 - 6
	--------------------------
 1	5 | 1 | 1 | 0 | 0 | 0 |
	--------------------------
 2	2 | 1 | 1 | 0 | 2 | 9 |
	--------------------------
 3	2 | 1 | 1 | 0 | 2 | 0 |
	--------------------------
 4	2 | 1 | 1 | 1 | 2 | 0 |
	--------------------------
 5	2 | 2 | 2 | 2 | 2 | 1 |
	--------------------------
*/
```
#### 你可以在 .choiceSlice方法下更改你所需要的数字的判断逻辑去实现你的业务需求 比如:"所谓的智能寻路助手,让你在所有路径下寻求最优解,当然你还需要在.aStar下实现你的选择逻辑"
```go
// .choiceSlice
		// 实现不同数字逻辑...例如一些channel 告诉你前面已经堵了?或是根据压载去选择较轻松的一边?
		
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
			...
```
```go
// .aStar
	// 不同与 .choiceSlice 在这里可以实现一些其他的功能 例如没有最优解会指派到一些临时处理口 例如高速公路上的一些应急口?
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
```

