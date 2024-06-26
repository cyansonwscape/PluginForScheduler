// pkg/plugin/myplugin.go
package plugin

import (
	"context"
	"log"

	// 导入随机数生成包
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// Name is the name of the plugin used in the plugin registry and configurations.

const Name = "greedy"

// Sort is a plugin that implements QoS class based sorting.

type sample struct{}

var _ framework.QueueSortPlugin = &sample{}
var _ framework.FilterPlugin = &sample{}
var _ framework.PreScorePlugin = &sample{}

// 添加 framework.ScorePlugin 接口的实现
var _ framework.ScorePlugin = &sample{}

// New initializes a new plugin and returns it.
func New(_ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
	return &sample{}, nil
}

// Name returns name of the plugin.
func (pl *sample) Name() string {
	return Name
}
func (s *Sample) Less(pod1, pod2 *framework.PodInfo) bool {
	// 获得pod1和pod2对CPU资源的请求量
	req1, _ := pod1.Pod.Spec.Containers[0].Resources.Requests.Cpu().AsInt64()
	req2, _ := pod2.Pod.Spec.Containers[0].Resources.Requests.Cpu().AsInt64()

	// 按CPU请求量从小到大排序
	return req1 < req2
}

func (pl *sample) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	log.Printf("filter pod: %v, node: %v", pod.Name, nodeInfo)
	log.Println(state)

	// 排除没有cpu=true标签的节点
	if nodeInfo.Node().Labels["cpu"] != "true" {
		return framework.NewStatus(framework.Unschedulable, "Node: "+nodeInfo.Node().Name)
	}
	return framework.NewStatus(framework.Success, "Node: "+nodeInfo.Node().Name)
}

func (pl *sample) PreScore(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodes []*v1.Node) *framework.Status {
	log.Println(nodes)
	return framework.NewStatus(framework.Success, "Node: "+pod.Name)
}

// ...

// Score is invoked after Filter and presumably all Filter and PreScore plugins have succeeded. To ensure Pod's schedulability, it is recommended to not make this function side effect.

// 假设的节点资源信息
type NodeResources struct {
	AllocatableCpu int64 // 节点可分配的CPU总量，单位是millicores
	AllocatedCpu   int64 // 节点已分配的CPU总量，单位是millicores
}

// 假设我们有一个全局变量，包含各节点的资源信息
var nodeResourcesMap = make(map[string]*NodeResources)

func (pl *sample) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	log.Printf("Scoring pod: %v on node: %v", pod.Name, nodeName)

	nodeResources := nodeResourcesMap[nodeName]

	// 计算节点剩余的CPU资源
	remainingCpu := nodeResources.AllocatableCpu - nodeResources.AllocatedCpu

	// 剩余资源较少的节点打高分（这里假设节点最多可能有1e4 millicores CPU未分配）
	score := int64(100) - (remainingCpu * int64(100) / int64(1e4))

	// 确保分数在[0, 100]范围内
	if score < 0 {
		score = 0
	} else if score > framework.MaxNodeScore {
		score = framework.MaxNodeScore
	}

	return score, framework.NewStatus(framework.Success)
}

// ScoreExtensions returns a ScoreExtensions interface if it should be invoked when 'score' is invoked.
func (pl *sample) ScoreExtensions() framework.ScoreExtensions {
	return nil
}
