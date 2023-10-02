package containers

import (
	"fmt"
	"math"
)

// TODO: don't expose center and members

type Cluster struct {
	Center  Vector
	Members []Vector
}

func (c *Cluster) Add(member Vector) {
	c.Members = append(c.Members, member)
}

func (c *Cluster) Recenter() error {
	memberCnt := len(c.Members)
	if memberCnt == 0 {
		return nil
		//return fmt.Errorf("there is no mean for an empty cluster")
	}

	// newCenter = "Mean" of the Members
	newCenter := make(Vector, len(c.Center))
	for _, member := range c.Members {
		newCenter.Add(member)
	}
	newCenter.Mul(1 / float64(memberCnt))

	c.Center = newCenter
	return nil
}

func (c *Cluster) RecenterReturningMovedDistance(distFn DistanceFunction) (moveDistances float64, err error) {
	memberCnt := len(c.Members)
	if memberCnt == 0 {
		//return 0, errors.New("kmeans: there is no mean for an empty cluster")
		return 0, nil
	}

	// newCenter is the "Mean" of the Members
	newCenter := make(Vector, len(c.Center))
	for _, member := range c.Members {
		newCenter.Add(member)
	}
	newCenter.Mul(1 / float64(memberCnt))

	moveDistances, err = distFn(c.Center, newCenter)
	if err != nil {
		return 0, err
	}

	c.Center = newCenter

	return moveDistances, nil
}

func (c *Cluster) SSE() float64 {
	sse := 0.0
	for i := 0; i < len(c.Members); i++ {
		dist, _ := EuclideanDistance(c.Center, c.Members[i])
		sse += math.Pow(dist, 2)
	}
	return sse
}

// Reset only resets the members of the cluster. The center is not reset.
func (c *Cluster) Reset() {
	c.Members = []Vector{}
}

func (c *Cluster) String() string {
	return fmt.Sprintf("Center: %v, Members: %v", c.Center, c.Members)
}
