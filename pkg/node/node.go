package node

// Node represents a node in the cluster
// It is a physical representation of a Worker, Manager, or Scheduler.
// It is ANY machine in the cluster

type Node struct {
	Name            string
	IP              string
	Cores           int
	Memory          int
	MemoryAllocated int
	Disk            int
	DiskAllocated   int
	Role            string
	TaskCount       int
}
