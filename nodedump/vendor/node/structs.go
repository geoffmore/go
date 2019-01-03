package node

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	k8s_core "k8s.io/api/core/v1"
	"os"
)

type k8sNode struct {
	name string
	// Resources
	capacity    k8s_core.ResourceList
	allocatable k8s_core.ResourceList
	// Labels
	nodeLabels map[string]string // .ObjectMeta.Labels
	// Annotations
	nodeAnnotations map[string]string // .ObjectMeta.Annotations
	// NodeInfo
	runtimeVer string //  .Status.NodeInfo.ContainerRuntimeVersion
	kernelVer  string //  .Status.NodeInfo.KernelVersion
	kubeletVer string //  .Status.NodeInfo.KubeletVersion
	osType     string //  .Status.NodeInfo.OperatingSystem
	osImage    string //  .Status.NodeInfo.OSImage
	arch       string //  .Status.NodeInfo.Architecture
	// Conditions
	isOutOfDisk     string
	hasDiskPressure string
	hasMemPressure  string
	hasPIDPressure  string
	isReady         string
}

type k8sInterface interface {
	labels() [][]string
	annotations() [][]string
	conditions() [][]string
	resources() [][]string
	info() [][]string
	Dump()
	Diff()
}

func (k k8sNode) info() [][]string {
	var nodeInfo [][]string

	nodeInfo = append(nodeInfo,
		[]string{headerPadding, "NODE INFO"},
		[]string{"arch", k.arch},
		[]string{"os", k.osType},
		[]string{"os version", k.osImage},
		[]string{"runtime", k.runtimeVer},
		[]string{"kernel", k.kernelVer},
		[]string{"kubelet", k.kubeletVer},
	)
	return nodeInfo
}

func (k k8sNode) labels() [][]string {
	var nodeLabels [][]string
	nodeLabels = append(nodeLabels, []string{headerPadding, "LABELS"})
	for label, val := range k.nodeLabels {
		someFunc(&val, true)
		nodeLabels = append(nodeLabels, []string{label, val})
	}
	return nodeLabels
}

func (k k8sNode) annotations() [][]string {
	var nodeAnnotations [][]string

	nodeAnnotations = append(nodeAnnotations, []string{headerPadding, "ANNOTATIONS"})
	for annotation, val := range k.nodeAnnotations {
		someFunc(&val, true)
		nodeAnnotations = append(nodeAnnotations, []string{annotation, val})
	}
	return nodeAnnotations
}

func (k k8sNode) conditions() [][]string {
	var nodeConditions [][]string
	nodeConditions = append(nodeConditions,
		[]string{headerPadding, "CONDITIONS"},
		[]string{"OutOfDisk", k.isOutOfDisk},
		[]string{"DiskPressure", k.hasDiskPressure},
		[]string{"MemoryPressure", k.hasMemPressure},
		[]string{"PIDPresssure", k.hasPIDPressure},
		[]string{"Ready", k.isReady},
	)
	return nodeConditions
}

func (k k8sNode) resources() [][]string {
	var nodeResources [][]string
	nodeResources = append(nodeResources,
		[]string{headerPadding, "RESOURCES"},
		[]string{"Capacity", ""},
		[]string{"cpu", fmt.Sprintf("%s", k.capacity.Cpu())},
		[]string{"ephemeral-storage", fmt.Sprintf("%s", k.capacity.StorageEphemeral())},
		[]string{"memory", fmt.Sprintf("%s", k.capacity.Memory())},
		[]string{"pods", fmt.Sprintf("%s", k.capacity.Pods())},
		[]string{"Allocated", ""},
		[]string{"cpu", fmt.Sprintf("%s", k.allocatable.Cpu())},
		[]string{"ephemeral-storage", fmt.Sprintf("%s", k.allocatable.StorageEphemeral())},
		[]string{"memory", fmt.Sprintf("%s", k.allocatable.Memory())},
		[]string{"pods", fmt.Sprintf("%s", k.allocatable.Pods())},
	)
	return nodeResources
}

func (k k8sNode) Dump() {
	var data [][]string

	if allFlag {
		data = append(data, k.info()...)
		data = append(data, k.resources()...)
		data = append(data, k.conditions()...)
		data = append(data, k.labels()...)
		data = append(data, k.conditions()...)
	} else {
		if infoFlag {
			data = append(data, k.info()...)
		}
		if specFlag {
			data = append(data, k.resources()...)

		}
		if condFlag {
			data = append(data, k.conditions()...)
		}
		if labelFlag {
			data = append(data, k.labels()...)
		}
		if annotationFlag {
			data = append(data, k.annotations()...)
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", k.name})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output

}

func (k k8sNode) infoDiff(k2 k8sNode) [][]string {
	var infoDiff [][]string
	infoDiff = append(infoDiff,
		[]string{headerPadding, "NODE INFO", headerPadding},
		[]string{"arch", k.arch, k2.arch},
		[]string{"os", k.osType, k2.osType},
		[]string{"os version", k.osImage, k2.osImage},
		[]string{"runtime", k.runtimeVer, k2.runtimeVer},
		[]string{"kernel", k.kernelVer, k2.kernelVer},
		[]string{"kubelet", k.kubeletVer, k2.kubeletVer},
	)
	return infoDiff
}

func (k k8sNode) resourceDiff(k2 k8sNode) [][]string {
	var resourceDiff [][]string

	resourceDiff = append(resourceDiff,
		[]string{headerPadding, "RESOURCES", headerPadding},
		[]string{"Capacity", "", ""},
		[]string{"cpu", fmt.Sprintf("%s", k.capacity.Cpu()), fmt.Sprintf("%s", k2.capacity.Cpu())},
		[]string{"ephemeral-storage", fmt.Sprintf("%s", k.capacity.StorageEphemeral()), fmt.Sprintf("%s", k2.capacity.StorageEphemeral())},
		[]string{"memory", fmt.Sprintf("%s", k.capacity.Memory()), fmt.Sprintf("%s", k2.capacity.Memory())},
		[]string{"pods", fmt.Sprintf("%s", k.capacity.Pods()), fmt.Sprintf("%s", k2.capacity.Pods())},
		[]string{"Allocated", "", ""},
		[]string{"cpu", fmt.Sprintf("%s", k.allocatable.Cpu()), fmt.Sprintf("%s", k2.allocatable.Cpu())},
		[]string{"ephemeral-storage", fmt.Sprintf("%s", k.allocatable.StorageEphemeral()), fmt.Sprintf("%s", k2.allocatable.StorageEphemeral())},
		[]string{"memory", fmt.Sprintf("%s", k.allocatable.Memory()), fmt.Sprintf("%s", k2.allocatable.Memory())},
		[]string{"pods", fmt.Sprintf("%s", k.allocatable.Pods()), fmt.Sprintf("%s", k2.allocatable.Pods())},
	)
	return resourceDiff
}

func (k k8sNode) conditionDiff(k2 k8sNode) [][]string {
	var conditionDiff [][]string
	conditionDiff = append(conditionDiff,
		[]string{headerPadding, "CONDITIONS", headerPadding},
		[]string{"OutOfDisk", k.isOutOfDisk, k2.isOutOfDisk},
		[]string{"DiskPressure", k.hasDiskPressure, k2.hasDiskPressure},
		[]string{"MemoryPressure", k.hasMemPressure, k2.hasMemPressure},
		[]string{"PIDPresssure", k.hasPIDPressure, k2.hasPIDPressure},
		[]string{"Ready", k.isReady, k2.isReady},
	)
	return conditionDiff
}

func (k k8sNode) labelDiff(k2 k8sNode) [][]string {
	var labelDiff [][]string

	labelDiff = append(labelDiff, []string{headerPadding, "LABELS", headerPadding})

	var usedLabels []string
	for label, val := range k.nodeLabels {
		val2, exists := k2.nodeLabels[label]
		someFunc(&val, true)
		someFunc(&val2, exists)
		labelDiff = append(labelDiff, []string{label, val, val2})
		usedLabels = append(usedLabels, label)
	}
	for label2, val2 := range k2.nodeLabels {
		if !stringInSlice(label2, usedLabels) {
			val, exists := k.nodeLabels[label2]
			someFunc(&val2, true)
			someFunc(&val, exists)
			labelDiff = append(labelDiff, []string{label2, val, val2})
		}
	}
	return labelDiff
}

func (k k8sNode) annotationDiff(k2 k8sNode) [][]string {
	var annotationDiff [][]string
	var usedAnnotations []string
	annotationDiff = append(annotationDiff, []string{headerPadding, "ANNOTATIONS", headerPadding})

	for annotation, val := range k.nodeAnnotations {
		val2, exists := k.nodeAnnotations[annotation]
		someFunc(&val, true)
		someFunc(&val2, exists)
		annotationDiff = append(annotationDiff, []string{annotation, val, val2})
		usedAnnotations = append(usedAnnotations, annotation)
	}
	for annotation2, val2 := range k.nodeAnnotations {
		if !stringInSlice(annotation2, usedAnnotations) {
			val, exists := k.nodeAnnotations[annotation2]
			someFunc(&val2, true)
			someFunc(&val, exists)
			annotationDiff = append(annotationDiff, []string{annotation2, val, val2})
		}
	}
	return annotationDiff
}

func (k k8sNode) Diff(k2 k8sNode) {
	var data [][]string

	if allFlag {
		data = append(data, k.infoDiff(k2)...)
		data = append(data, k.resourceDiff(k2)...)
		data = append(data, k.conditionDiff(k2)...)
		data = append(data, k.labelDiff(k2)...)
		data = append(data, k.annotationDiff(k2)...)
	} else {
		if infoFlag {
			data = append(data, k.infoDiff(k2)...)
		}
		if specFlag {
			data = append(data, k.resourceDiff(k2)...)
		}
		if condFlag {
			data = append(data, k.conditionDiff(k2)...)
		}
		if labelFlag {
			data = append(data, k.labelDiff(k2)...)
		}
		if annotationFlag {
			data = append(data, k.annotationDiff(k2)...)
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", k.name, k2.name})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}
