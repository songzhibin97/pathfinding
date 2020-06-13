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
## aStar定制性很强,你可以根据不同的业务去设置一些你自己的逻辑 例如
```
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
 1	5 | 1 | 1 | 0 | 0 | 0 |
	--------------------------
 2	2 | p | 1 | 0 | 2 | 9 |
	--------------------------
 3	2 | 1 | 1 | 0 | 2 | 0 |
	--------------------------
 4	2 | 1 | 1 | 1 | 2 | 0 |
	--------------------------
 5	2 | 2 | 2 | 2 | 2 | 1 |
	--------------------------
*/
```
## 你可以在 .choiceSlice方法下更改你所需要的数字的判断逻辑去实现你的业务需求 比如:"所谓的智能寻路助手,让你在所有路径下寻求最优解,当然你还需要在.aStar下实现你的选择逻辑"
