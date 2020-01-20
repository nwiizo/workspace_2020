package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

//--------------------------------------------------------------------------------------------
//	Main
//--------------------------------------------------------------------------------------------

func main() {

	home := homeDir()

	logWriter, err := os.OpenFile(filepath.Join(home, "k8watcher.log"),
		os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err.Error())
	}
	log.SetOutput(logWriter)
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	log.Println("start...")

	//
	//   KubeUi <----------- KubeWatcher
	//            kubeEvent

	eventChan := make(chan *KubeEvent, 10)
	kubeWatcher := NewKubeWatcher(eventChan)
	kubeUI := NewKubeUI(eventChan)

	go kubeWatcher.Run()
	kubeUI.Run()
}

//--------------------------------------------------------------------------------------------
//	Event object
//--------------------------------------------------------------------------------------------

// KubeEvent .
type KubeEvent struct {
	eventType EventType
	newObj    interface{}
}

// NewKubeEvent .
func NewKubeEvent(eventType EventType, newObj interface{}) *KubeEvent {
	return &KubeEvent{
		eventType: eventType,
		newObj:    newObj,
	}
}

// EventType .
type EventType int

const (
	// NodeAdd .
	NodeAdd EventType = iota
	// NodeUpdate .
	NodeUpdate
	// NodeDelete .
	NodeDelete
	// PodAdd .
	PodAdd
	// PodUpdate .
	PodUpdate
	// PodDelete .
	PodDelete
)

//--------------------------------------------------------------------------------------------
//	Kubenetes Node/Pod Watcher
//--------------------------------------------------------------------------------------------

// KubeWatcher .
type KubeWatcher struct {
	Sender chan<- *KubeEvent
}

// NewKubeWatcher .
func NewKubeWatcher(sender chan *KubeEvent) *KubeWatcher {
	return &KubeWatcher{
		Sender: sender,
	}
}

// Run .
func (kw *KubeWatcher) Run() {

	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	watchNodes := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(), "nodes", v1.NamespaceAll, fields.Everything())

	_, nodesController := cache.NewInformer(
		watchNodes, &v1.Node{}, 0,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    func(obj interface{}) { kw.Sender <- NewKubeEvent(NodeAdd, obj) },
			DeleteFunc: func(obj interface{}) { kw.Sender <- NewKubeEvent(NodeDelete, obj) },
			UpdateFunc: func(old, new interface{}) { kw.Sender <- NewKubeEvent(NodeUpdate, new) },
		},
	)
	watchPods := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(), string(v1.ResourcePods), v1.NamespaceAll, fields.Everything())
	_, podsController := cache.NewInformer(
		watchPods, &v1.Pod{}, 0,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    func(obj interface{}) { kw.Sender <- NewKubeEvent(PodAdd, obj) },
			DeleteFunc: func(obj interface{}) { kw.Sender <- NewKubeEvent(PodDelete, obj) },
			UpdateFunc: func(old, new interface{}) { kw.Sender <- NewKubeEvent(PodUpdate, new) },
		},
	)

	stop := make(chan struct{})
	defer close(stop)
	go nodesController.Run(stop)
	go podsController.Run(stop)

	for {
		time.Sleep(time.Second)
	}
}

//--------------------------------------------------------------------------------------------
//	Text User Interface
//--------------------------------------------------------------------------------------------

// KubeUI .
type KubeUI struct {
	Receiver        <-chan *KubeEvent
	Nodes           []*v1.Node
	Pods            []*v1.Pod
	app             *tview.Application
	nodeGrid        *tview.Grid
	nodeViews       []*tview.Table
	nodeColumnCount int
}

// NewKubeUI .
func NewKubeUI(rec chan *KubeEvent) *KubeUI {

	runewidth.DefaultCondition.EastAsianWidth = false
	ui := &KubeUI{
		Receiver:        rec,
		Nodes:           []*v1.Node{},
		Pods:            []*v1.Pod{},
		nodeViews:       []*tview.Table{},
		app:             tview.NewApplication(),
		nodeColumnCount: 1,
	}

	rootView := tview.NewGrid().SetRows(3, -1, 2)
	rootView.SetBackgroundColor(tcell.NewHexColor(0xe0e0e0))

	headView := tview.NewTextView()
	headView.SetDynamicColors(true).SetBackgroundColor(tcell.NewHexColor(0x303030))
	headView.SetBorderPadding(1, 1, 2, 0)
	headView.SetText("[#306ee3]⎈ [white]Kubernetes Watcher")

	ui.nodeGrid = tview.NewGrid()
	ui.nodeGrid.SetBorderPadding(1, 1, 2, 2)
	ui.nodeGrid.SetBackgroundColor(tcell.NewHexColor(0xf6f6f4))
	ui.nodeGrid.SetGap(1, 2)

	bottomView := tview.NewTextView()
	bottomView.SetDynamicColors(true).SetBackgroundColor(tcell.NewHexColor(0xe0e0e0))
	bottomView.SetText("[#306ee3] watching...     CTRL+C -> Exit")

	ui.app.SetAfterDrawFunc(func(screen tcell.Screen) {
		ui.drowNodeGrid(false)
	})

	ui.app.SetRoot(
		rootView.
			AddItem(headView, 0, 0, 1, 1, 0, 0, false).
			AddItem(ui.nodeGrid, 1, 0, 1, 1, 0, 0, false).
			AddItem(bottomView, 2, 0, 1, 1, 0, 0, false),
		true,
	)

	return ui
}

// Run .
func (ui *KubeUI) Run() {

	go ui.EventReciever()

	if err := ui.app.Run(); err != nil {
		panic(err)
	}
}

// EventReciever .
func (ui *KubeUI) EventReciever() {
	for {
		kubeEvent := <-ui.Receiver

		log.Println(kubeEvent)
		switch kubeEvent.eventType {
		case NodeAdd:
			ui.AddNode(kubeEvent.newObj.(*v1.Node))
		case NodeUpdate:
			ui.UpdateNode(kubeEvent.newObj.(*v1.Node))
		case NodeDelete:
			ui.RemoveNode(kubeEvent.newObj.(*v1.Node))
		case PodAdd:
			ui.AddPod(kubeEvent.newObj.(*v1.Pod))
		case PodUpdate:
			ui.UpdatePod(kubeEvent.newObj.(*v1.Pod))
		case PodDelete:
			ui.RemovePod(kubeEvent.newObj.(*v1.Pod))
		}
	}
}

// AddNode .
func (ui *KubeUI) AddNode(v1Node *v1.Node) {
	nodeView := tview.NewTable()
	nodeView.SetBackgroundColor(tcell.NewHexColor(0x454545))
	nodeView.Select(0, 0).SetFixed(1, 1).SetSelectable(true, false)
	nodeView.SetBorder(true).SetBorderPadding(1, 1, 1, 1)
	for _, nodeAddress := range v1Node.Status.Addresses {
		if nodeAddress.Type == v1.NodeInternalIP {
			nodeView.SetTitleAlign(tview.AlignLeft).SetTitle(nodeAddress.Address)
			break
		}
	}
	ui.Nodes = append(ui.Nodes, v1Node)
	ui.nodeViews = append(ui.nodeViews, nodeView)
	ui.drawPods(v1Node.Name)
	ui.drowNodeGrid(true)
}

// UpdateNode .
func (ui *KubeUI) UpdateNode(v1Node *v1.Node) {
	if i, _, _, exist := ui.getNode(v1Node.Name); exist {
		ui.Nodes[i] = v1Node
		ui.drawPods(v1Node.Name)
	}
}

// RemoveNode .
func (ui *KubeUI) RemoveNode(v1Node *v1.Node) {
	if i, _, _, exist := ui.getNode(v1Node.Name); exist {
		ui.Nodes = append(ui.Nodes[:i], ui.Nodes[i+1:]...)
		ui.nodeViews = append(ui.nodeViews[:i], ui.nodeViews[i+1:]...)
		ui.drowNodeGrid(true)
	}
}

// AddPod .
func (ui *KubeUI) AddPod(v1Pod *v1.Pod) {
	log.Println("addPod", v1Pod.Spec.NodeName, v1Pod.Name)
	ui.Pods = append(ui.Pods, v1Pod)
	ui.drawPods(v1Pod.Spec.NodeName)
	ui.appDrow()
}

// UpdatePod .
func (ui *KubeUI) UpdatePod(v1Pod *v1.Pod) {
	for i, p := range ui.Pods {
		if p.Name == v1Pod.Name {
			log.Println("UpdatePod", v1Pod.Spec.NodeName, v1Pod.Name)
			ui.Pods[i] = v1Pod
			ui.drawPods(v1Pod.Spec.NodeName)
			ui.appDrow()
			break
		}
	}
}

// RemovePod .
func (ui *KubeUI) RemovePod(v1Pod *v1.Pod) {
	for i, p := range ui.Pods {
		if p.Name == v1Pod.Name {
			log.Println("RemovePod", v1Pod.Spec.NodeName, v1Pod.Name)
			ui.Pods = append(ui.Pods[:i], ui.Pods[i+1:]...)
			ui.drawPods(v1Pod.Spec.NodeName)
			ui.appDrow()
			break
		}
	}
}

func (ui *KubeUI) appDrow() {
	ui.app.QueueUpdateDraw(func() {
	})
}

func (ui *KubeUI) getNode(nodeName string) (int, *v1.Node, *tview.Table, bool) {
	for i, n := range ui.Nodes {
		if n.Name == nodeName {
			return i, n, ui.nodeViews[i], true
		}
	}
	return -1, nil, nil, false
}

func (ui *KubeUI) drowNodeGrid(force bool) {
	if len(ui.nodeViews) > 0 {
		// border<2> -----48------ gap<2> -----48------ <border<2>
		_, _, w, _ := ui.nodeGrid.GetRect()
		c := (w - 4 + 2) / (48 + 2)
		if c == 0 {
			c = 1
		}
		log.Println("adjustNodeGrid", c, w)
		if c != ui.nodeColumnCount || force {
			ui.nodeColumnCount = c
			ui.nodeGrid.Clear()
			for i, nodeView := range ui.nodeViews {
				ui.nodeGrid.AddItem(nodeView, i/c, i%c, 1, 1, 0, 0, true)
			}
			ui.appDrow()
		}
	}
}

func (ui *KubeUI) drawPods(nodeName string) {
	if _, _, nodeView, exist := ui.getNode(nodeName); exist {
		nodeView.Clear()
		nodeView.SetCell(0, 0, tview.NewTableCell("[#306ee3]NAMESPACE").SetAlign(tview.AlignCenter))
		nodeView.SetCell(0, 1, tview.NewTableCell("[#306ee3]POD NAME").SetAlign(tview.AlignCenter))
		nodeView.SetCell(0, 2, tview.NewTableCell("[#306ee3]STATUS").SetAlign(tview.AlignCenter))

		row := 0
		for _, v1Pod := range ui.Pods {
			if v1Pod.Spec.NodeName == nodeName {
				row++
				nodeView.SetCell(row, 0, tview.NewTableCell(v1Pod.Namespace).SetMaxWidth(11))
				nodeView.SetCell(row, 1, tview.NewTableCell(v1Pod.Name).SetMaxWidth(20))
				nodeView.SetCell(row, 2, tview.NewTableCell(podStats(v1Pod)).SetMaxWidth(17))
			}
		}
	}
}

// podStats .
// 参考）
// https://github.com/kubernetes/kubernetes/blob/master/pkg/printers/internalversion/printers.go#printPod
func podStats(pod *v1.Pod) string {
	reason := string(pod.Status.Phase)
	if pod.Status.Reason != "" {
		reason = pod.Status.Reason
	}

	initializing := false
	for i := range pod.Status.InitContainerStatuses {
		container := pod.Status.InitContainerStatuses[i]
		switch {
		case container.State.Terminated != nil && container.State.Terminated.ExitCode == 0:
			continue
		case container.State.Terminated != nil:
			if len(container.State.Terminated.Reason) == 0 {
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Init:Signal:%d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("Init:ExitCode:%d", container.State.Terminated.ExitCode)
				}
			} else {
				reason = "Init:" + container.State.Terminated.Reason
			}
			initializing = true
		case container.State.Waiting != nil && len(container.State.Waiting.Reason) > 0 && container.State.Waiting.Reason != "PodInitializing":
			reason = "Init:" + container.State.Waiting.Reason
			initializing = true
		default:
			reason = fmt.Sprintf("Init:%d/%d", i, len(pod.Spec.InitContainers))
			initializing = true
		}
		break
	}
	if !initializing {
		hasRunning := false
		for i := len(pod.Status.ContainerStatuses) - 1; i >= 0; i-- {
			container := pod.Status.ContainerStatuses[i]
			if container.State.Waiting != nil && container.State.Waiting.Reason != "" {
				reason = container.State.Waiting.Reason
			} else if container.State.Terminated != nil && container.State.Terminated.Reason != "" {
				reason = container.State.Terminated.Reason
			} else if container.State.Terminated != nil && container.State.Terminated.Reason == "" {
				if container.State.Terminated.Signal != 0 {
					reason = fmt.Sprintf("Signal:%d", container.State.Terminated.Signal)
				} else {
					reason = fmt.Sprintf("ExitCode:%d", container.State.Terminated.ExitCode)
				}
			} else if container.Ready && container.State.Running != nil {
				hasRunning = true
			}
		}

		if reason == "Completed" && hasRunning {
			reason = "Running"
		}
	}

	// if pod.DeletionTimestamp != nil && pod.Status.Reason == node.NodeUnreachablePodReason {
	// 	reason = "Unknown"
	// } else
	if pod.DeletionTimestamp != nil {
		reason = "Terminating"
	}
	return reason
}

//--------------------------------------------------------------------------------------------
//	Utilities
//--------------------------------------------------------------------------------------------
func jsonDump(a ...interface{}) {
	jsonBytes, err := json.Marshal(a)
	if err != nil {
		panic(err.Error())
	}
	out := new(bytes.Buffer)
	json.Indent(out, jsonBytes, "", "    ")
	fmt.Println(out.String())
}

// homeDir .
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
