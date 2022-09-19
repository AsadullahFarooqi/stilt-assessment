package main

import (
	"testing"
)

func Test_findCourier(t *testing.T) {
	type args struct {
		couriers []*Courier
		prepTime int64
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success test",
			args: args{
				couriers: []*Courier{
					&Courier{
						Name:       "12",
						ArriveTime: 12,
					},
				},
				prepTime: 12,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "failed test",
			args: args{
				couriers: []*Courier{
					&Courier{
						Name:       "12",
						ArriveTime: 12,
					},
				},
				prepTime: 13,
			},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findCourier(tt.args.couriers, tt.args.prepTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("findCourier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("findCourier() = %v, want %v", got, tt.want)
			}
		})
	}
}
