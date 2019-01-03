package node

import (
	"encoding/json"
	k8s_core "k8s.io/api/core/v1"
)

const commandPrefix string = "kubectl get node"
const commandSuffix string = "-o json"
const headerPadding string = "---"

func conditionStatus2String(some k8s_core.ConditionStatus) string {
	var status string
	switch some {
	case k8s_core.ConditionTrue:
		status = "true"
	case k8s_core.ConditionFalse:
		status = "false"
	case k8s_core.ConditionUnknown:
		status = "unknown"
	default:
		status = "logical error"
	}
	return status
}

func genNodeInfo(node string) k8sNode {
	var nodeStructRaw k8s_core.Node
	var nodeStruct k8sNode

	cmdString := genCommand(commandPrefix, node, commandSuffix)
	rawJSON := runCommand(cmdString)

	// Marshal -> JSON
	// Unmarshal <- JSON
	json.Unmarshal(rawJSON, &nodeStructRaw)

	nodeStruct.arch = nodeStructRaw.Status.NodeInfo.Architecture
	nodeStruct.runtimeVer = nodeStructRaw.Status.NodeInfo.ContainerRuntimeVersion
	nodeStruct.kernelVer = nodeStructRaw.Status.NodeInfo.KernelVersion
	nodeStruct.kubeletVer = nodeStructRaw.Status.NodeInfo.KubeletVersion
	nodeStruct.osType = nodeStructRaw.Status.NodeInfo.OperatingSystem
	nodeStruct.osImage = nodeStructRaw.Status.NodeInfo.OSImage
	nodeStruct.nodeLabels = nodeStructRaw.ObjectMeta.Labels
	nodeStruct.capacity = nodeStructRaw.Status.Capacity
	nodeStruct.allocatable = nodeStructRaw.Status.Allocatable
	nodeStruct.nodeAnnotations = nodeStructRaw.ObjectMeta.Annotations
	nodeStruct.name = nodeStructRaw.Name

	nodeConditions := nodeStructRaw.Status.Conditions
	for _, condition := range nodeConditions {
		switch condition.Type {
		case "OutOfDisk":
			nodeStruct.isOutOfDisk = conditionStatus2String(condition.Status)
		case "MemoryPressure":
			nodeStruct.hasMemPressure = conditionStatus2String(condition.Status)
		case "DiskPressure":
			nodeStruct.hasDiskPressure = conditionStatus2String(condition.Status)
		case "PIDPressure":
			nodeStruct.hasPIDPressure = conditionStatus2String(condition.Status)
		case "Ready":
			nodeStruct.isReady = conditionStatus2String(condition.Status)
		}
	}
	return nodeStruct

}

func InitNodes(node1string, node2string string) (node1, node2 k8sNode) {

	// Generate the struct we need to compare nodes
	node1 = genNodeInfo(node1string)
	if node2string != "" {
		node2 = genNodeInfo(node2string)
	}
	return node1, node2

}
