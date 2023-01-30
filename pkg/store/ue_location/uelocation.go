package ue_location

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/wangxn2015/myRANsim/api/ue_location"
	"math/rand"
	"time"
)

type UeStore interface {
	Search(ctx context.Context, found func(ue *ue_location.UeInfo) error) error
}

type InMemoryUeStore struct {
	data map[uint64]*ue_location.UeInfo //指针 or 实际值？

}

func NewInMemoryUeStore() *InMemoryUeStore {
	return &InMemoryUeStore{
		data: make(map[uint64]*ue_location.UeInfo),
	}
}

func (store InMemoryUeStore) generateData() {
	for i := 0; i < 3; i++ {
		store.data[uint64(2023+i)] = &ue_location.UeInfo{
			Imsi: uint64(2023 + i),
			Location: &ue_location.Location{
				Lat: float64(i),
				Lng: float64(i + 10),
			},
		}
	}
}

func (store *InMemoryUeStore) moveUes() {
	for i := 0; i < 3; i++ {
		Location := &ue_location.Location{
			Lat: 40.075 + float64(i+1)*0.0005*rand.Float64(),
			Lng: 116.24 + float64(i+1)*0.0005*rand.Float64(),
		}

		store.data[uint64(2000+i)].Location = Location
	}

}

func (store *InMemoryUeStore) Search(ctx context.Context, found func(ue *ue_location.UeInfo) error) error {

	if len(store.data) == 0 {
		store.generateData()
	}

	for j := 0; j < 100; j++ {
		store.moveUes()
		for _, ue := range store.data {
			if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
				fmt.Println("context is cancelled")
				return nil
			}
			other, err := deepCopy(ue)
			if err != nil {
				return err
			}
			err = found(other)
			if err != nil {
				return err
			}
		}
		time.Sleep(time.Second)
	}

	return nil
}

func deepCopy(ue *ue_location.UeInfo) (*ue_location.UeInfo, error) {
	other := &ue_location.UeInfo{}

	err := copier.Copy(other, ue)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %w", err)
	}

	return other, nil
}
