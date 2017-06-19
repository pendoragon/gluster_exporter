package info

type GlusterInfo struct {
	VolumeDetail VolumeDetail
	PeerStatus   PeerStatus
}

// VolumeStatusXML XML type of "gluster volume status"
type VolumeDetail struct {
	Volumes []Volume
}

type Volume struct {
	Status    string
	VolName   string
	NodeCount int
	Nodes      []Node
}

type Node struct {
	Hostname string
	Path     string
	PeerID   string
	Status   int
	Port     int
	Ports    []Port
	Pid        int
	SizeTotal  uint64
	SizeFree   uint64
	Device     string
	BlockSize  int
	MntOptions string
	FsName     string
}

type Port struct {
	TCP  int
	RDMA string
}
