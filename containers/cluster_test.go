package containers

import (
	"reflect"
	"testing"
)

func TestCluster_Recenter(t *testing.T) {
	type fields struct {
		center  Vector
		members []Vector
	}
	tests := []struct {
		name       string
		fields     fields
		wantCenter Vector
	}{
		{
			name: "test1",
			fields: fields{
				center:  Vector{1, 1},
				members: []Vector{{1, 1}, {2, 2}},
			},
			wantCenter: Vector{1.5, 1.5},
		},
		{
			name: "test2",
			fields: fields{
				center:  Vector{1, 1},
				members: []Vector{{1, 1}, {2, 2}, {3, 3}},
			},
			wantCenter: Vector{2, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cluster{
				center:  tt.fields.center,
				members: tt.fields.members,
			}
			c.Recenter()
			if c.center.Compare(tt.wantCenter) != 0 {
				t.Errorf("Recenter() gotCenter = %v, want %v", c.center, tt.wantCenter)
			}
		})
	}
}

func TestCluster_RecenterReturningMovedDistance(t *testing.T) {
	type fields struct {
		center  Vector
		members []Vector
	}
	type args struct {
		distFn DistanceFunction
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantMoveDistances float64
		wantCenter        Vector
		wantErr           bool
	}{
		{
			name: "empty cluster",
			fields: fields{
				center:  Vector{1, 1},
				members: []Vector{},
			},
			args:       args{distFn: EuclideanDistance},
			wantCenter: Vector{1, 1}, // unchanged
		},
		{
			name: "non-empty cluster",
			fields: fields{
				center:  Vector{1, 1},
				members: []Vector{{1, 1}, {2, 2}},
			},
			args:              args{distFn: EuclideanDistance},
			wantMoveDistances: 0.7071067811865476,
			wantCenter:        Vector{1.5, 1.5}, // changed
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cluster{
				center:  tt.fields.center,
				members: tt.fields.members,
			}
			gotMoveDistance, err := c.RecenterWithMovedDistance(tt.args.distFn)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecenterReturningMovedDistance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMoveDistance != tt.wantMoveDistances {
				t.Errorf("RecenterReturningMovedDistance() gotMoveDistance = %v, want %v", gotMoveDistance, tt.wantMoveDistances)
			}
			if c.center.Compare(tt.wantCenter) != 0 {
				t.Errorf("Recenter() gotCenter = %v, want %v", c.center, tt.wantCenter)
			}
		})
	}
}

func TestCluster_Reset(t *testing.T) {
	type fields struct {
		center  Vector
		members []Vector
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Test1",
			fields: fields{
				center:  Vector{1, 1},
				members: []Vector{{1, 1}, {2, 2}, {3, 3}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cluster{
				center:  tt.fields.center,
				members: tt.fields.members,
			}
			c.Reset()
			if len(c.members) != 0 {
				t.Errorf("Reset() = %v, want %v", c.members, []Vector{})
			}
			if c.center.Compare(tt.fields.center) != 0 {
				t.Errorf("Reset() = %v, want %v", c.center, tt.fields.center)
			}
		})
	}
}

func TestCluster_SSE(t *testing.T) {
	type fields struct {
		center  Vector
		members []Vector
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Test1",
			fields: fields{
				center:  Vector{1, 1},
				members: []Vector{{1, 1}, {3, 3}, {3, 3}},
			},
			want: 16.000000000000004,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cluster{
				center:  tt.fields.center,
				members: tt.fields.members,
			}
			if got := c.SSE(); got != tt.want {
				t.Errorf("SSE() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCluster_String(t *testing.T) {
	type fields struct {
		center  Vector
		members []Vector
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test1",
			fields: fields{
				center:  Vector{1, 1},
				members: []Vector{{1, 1}, {3, 3}, {3, 3}},
			},
			want: "center: [1 1], members: [[1 1] [3 3] [3 3]]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cluster{
				center:  tt.fields.center,
				members: tt.fields.members,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCluster_AddMember_Members(t *testing.T) {

	tests := []struct {
		name        string
		cluster     *Cluster
		argMembers  []Vector
		wantMembers []Vector
	}{
		{
			name: "Test1",
			cluster: &Cluster{
				center:  Vector{1, 1},
				members: []Vector{{1, 2}, {2, 3}},
			},
			argMembers:  []Vector{{4, 5}, {6, 7}},
			wantMembers: []Vector{{1, 2}, {2, 3}, {4, 5}, {6, 7}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, member := range tt.argMembers {
				tt.cluster.AddMember(member)
			}

			if !reflect.DeepEqual(tt.cluster.Members(), tt.wantMembers) {
				t.Errorf("Members() = %v, want %v", tt.cluster.Members(), tt.wantMembers)
			}

		})
	}
}
