// Copyright 2023 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package moarray

import (
	"github.com/arjunsk/kmeans/utils/assertx"
	"reflect"
	"testing"
)

func TestNormalizeL2(t *testing.T) {
	type args struct {
		argF32 []float32
		argF64 []float64
	}
	type testCase struct {
		name string
		args args

		wantF32 []float32
		wantF64 []float64
		wantErr bool
	}
	tests := []testCase{
		{
			name:    "Test1 - float32 - zero vector",
			args:    args{argF32: []float32{0, 0, 0}},
			wantF32: []float32{0, 0, 0},
		},
		{
			name:    "Test1.b - float32",
			args:    args{argF32: []float32{1, 2, 3}},
			wantF32: []float32{0.26726124, 0.5345225, 0.80178374},
		},
		{
			name:    "Test1.c - float32",
			args:    args{argF32: []float32{10, 3.333333333333333, 4, 5}},
			wantF32: []float32{0.8108108, 0.27027026, 0.32432434, 0.4054054},
		},
		{
			name:    "Test2 - float64 - zero vector",
			args:    args{argF64: []float64{0, 0, 0}},
			wantF64: []float64{0, 0, 0},
		},
		{
			name:    "Test3 - float64",
			args:    args{argF64: []float64{1, 2, 3}},
			wantF64: []float64{0.2672612419124244, 0.5345224838248488, 0.8017837257372732},
		},
		{
			name:    "Test4 - float64",
			args:    args{argF64: []float64{-1, 2, 3}},
			wantF64: []float64{-0.2672612419124244, 0.5345224838248488, 0.8017837257372732},
		},
		{
			name:    "Test5 - float64",
			args:    args{argF64: []float64{10, 3.333333333333333, 4, 5}},
			wantF64: []float64{0.8108108108108107, 0.27027027027027023, 0.3243243243243243, 0.4054054054054054},
		},
		{
			name:    "Test6 - float64",
			args:    args{argF64: []float64{1, 2, 3.6666666666666665, 4.666666666666666}},
			wantF64: []float64{0.15767649936829103, 0.31535299873658207, 0.5781471643504004, 0.7358236637186913},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.args.argF32 != nil {
				if tt.wantErr {
					if _, err := NormalizeL2[float32](tt.args.argF32); err == nil {
						t.Errorf("NormalizeL2() should throw error")
					}
				} else if gotRes, err := NormalizeL2[float32](tt.args.argF32); err != nil || !reflect.DeepEqual(tt.wantF32, gotRes) {
					t.Errorf("NormalizeL2() = %v, want %v", gotRes, tt.wantF32)
				}
			}
			if tt.args.argF64 != nil {
				if tt.wantErr {
					if _, err := NormalizeL2[float64](tt.args.argF64); err == nil {
						t.Errorf("NormalizeL2() should throw error")
					}
				} else if gotRes, err := NormalizeL2[float64](tt.args.argF64); err != nil || !assertx.InEpsilonF64Slice(tt.wantF64, gotRes) {
					t.Errorf("NormalizeL2() = %v, want %v", gotRes, tt.wantF64)
				}
			}
		})
	}
}
