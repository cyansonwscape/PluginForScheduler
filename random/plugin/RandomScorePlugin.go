// pkg/plugin/myplugin.go
package plugin

// import (
// 	"context"
// 	"log"
// 	"math/rand" // 导入随机数生成包

// 	"k8s.io/apimachinery/pkg/runtime"
// 	"k8s.io/kubernetes/pkg/scheduler/framework"
// )

import (
	"context"
	"log"
	"math/rand"
)

// Name is the name of the plugin used in the plugin registry and configurations.

const Name = "RandomScorePlugin"

// Sort is a plugin that implements QoS class based sorting.

type sample struct{}

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
// 返回随机分数
func (pl *sample) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	log.Printf("Scoring pod: %v on node: %v", pod.Name, nodeName)

	// 生成一个随机分数
	score := rand.Int63n(framework.MaxNodeScore)
	return score, framework.NewStatus(framework.Success)
}

// ScoreExtensions returns a ScoreExtensions interface if it should be invoked when 'score' is invoked.
func (pl *sample) ScoreExtensions() framework.ScoreExtensions {
	return nil
}
