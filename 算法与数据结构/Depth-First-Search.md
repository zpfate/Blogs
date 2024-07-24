# Depth First Search

## 前言

> **深度优先搜索算法**（英语：Depth-First-Search，DFS）是一种用于遍历或搜索[树](https://zh.wikipedia.org/wiki/树_(数据结构))或[图](https://zh.wikipedia.org/wiki/图_(数学))的[算法](https://zh.wikipedia.org/wiki/算法)。这个算法会尽可能深地搜索树的分支。当节点v的所在边都己被探寻过，搜索将回溯到发现节点v的那条边的起始节点。这一过程一直进行到已发现从源节点可达的所有节点为止

<font color=red>通常隐含了**栈**的实现</font>

## 问题分类

1. 树的DFS
2. 一般 DFS 问题：**一维问题(回溯法)**和**二维问题(图)**
3. 大多数一般递归问题也是利用 DFS 求解

## 二叉树遍历法

```java

public class DepthFirstSearchInBinaryTree {
    public static class TreeNode<T> {
        T val;
        TreeNode<T> left;
        TreeNode<T> right;

        TreeNode(T rootData) {
            val = rootData;
        }
    }

    public <T> void dfs(TreeNode<T> node) {
        // 需要将具体问题转化 => 具体问题需要做哪些事情
        doSomething(node);
        dfs(node.left);
        dfs(node.right);
    }
}
```

## 一般DFS问题模板

```java
public class DepthFirstSearch {
    class Point {
        int num;
        int value;
    }

    private boolean[] marked;
    private int count;

    public DepthFirstSearch(Graph graph, Point start) {
        marked = new boolean[graph.length()];
        dfs(graph, start);
    }

    public void dfs(Graph graph, Point point) {
        marked[point.num] = true;
        count++;

        // graph.adj(point) 表示和 point 相邻的所有节点 周围
        for (Point aroundPoint : graph.adj(point)) {
            if (!marked[aroundPoint.num]) {
                dfs(graph, aroundPoint);
            }
        }
    }

    public int count() {
        return count;
    }
}
```

