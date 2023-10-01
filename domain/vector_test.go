package domain

import "testing"

func TestVector_Add(t *testing.T) {
	type args struct {
		vec Vector
	}
	tests := []struct {
		name string
		v    Vector
		args args
		want Vector
	}{
		{
			name: "test1",
			v:    Vector{1, 2, 3},
			args: args{Vector{-1, 4, 6}},
			want: Vector{0, 6, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.v.Add(tt.args.vec)
			if tt.v.Compare(tt.want) != 0 {
				t.Errorf("Vector.Add() = %v, want %v", tt.v, tt.want)
			}
		})
	}
}

func TestVector_Mul(t *testing.T) {
	type args struct {
		scalar float64
	}
	tests := []struct {
		name string
		v    Vector
		args args
		want Vector
	}{
		{
			name: "test1",
			v:    Vector{1, 2, 3},
			args: args{4},
			want: Vector{4, 8, 12},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.v.Mul(tt.args.scalar)
			if tt.v.Compare(tt.want) != 0 {
				t.Errorf("Vector.Add() = %v, want %v", tt.v, tt.want)
			}
		})
	}
}

func TestVector_Compare(t *testing.T) {
	type args struct {
		v2 Vector
	}
	tests := []struct {
		name string
		v    Vector
		args args
		want int
	}{

		{
			name: "test1",
			v:    Vector{1, 2, 3},
			args: args{Vector{1, 2, 3}},
			want: 0,
		},
		{
			name: "test2",
			v:    Vector{1, 2, 3},
			args: args{Vector{1, 2, 4}},
			want: -1,
		},
		{
			name: "test3",
			v:    Vector{1, 2, 3},
			args: args{Vector{1, 2, 2}},
			want: 1,
		},
		{
			name: "test4",
			v:    Vector{1, 2, 3},
			args: args{Vector{1, 2, 3, 4}},
			want: -1,
		},
		{
			name: "test5",
			v:    Vector{1, 2, 3, 4},
			args: args{Vector{1, 2, 3}},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.Compare(tt.args.v2); got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
