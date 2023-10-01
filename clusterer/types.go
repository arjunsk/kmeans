package clusterer

import (
	"errors"
	"github.com/arjunsk/go-kmeans/domain"
)

type Clusterer interface {
	Cluster() (domain.Clusters, error)
}

func StdInputCheck(clusterCnt, vectorCnt int) error {
	if clusterCnt > vectorCnt {
		return errors.New("kmeans: The count of the data set must at least equal k")
	}
	return nil
}
