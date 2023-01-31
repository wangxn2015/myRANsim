package ue_location

import (
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	model "github.com/wangxn2015/myRANsim/api/ue_location"
	"github.com/wangxn2015/myRANsim/pkg/store/ue_location_store"
	"github.com/wangxn2015/onos-lib-go/pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var log = logging.GetLogger()

func NewService(ueStore ue_location_store.UeStore) service.Service {
	return &Service{
		ueStore: ueStore,
	}
}

type Service struct {
	ueStore ue_location_store.UeStore
}

func (s Service) Register(r *grpc.Server) {
	server := &Server{
		ueStore: s.ueStore,
	}
	model.RegisterUeLocationServiceServer(r, server)
}

type Server struct {
	model.UnimplementedUeLocationServiceServer
	ueStore ue_location_store.UeStore
}

func (s Server) GetUes(request *model.UesLocationRequest, server model.UeLocationService_GetUesServer) error {
	//TODO implement me

	panic("implement me")
}

func (s Server) GetUe(request *model.UeLocationRequest, server model.UeLocationService_GetUeServer) error {
	log.Info("GetUe request: %+v", request)

	err := s.ueStore.Search(
		server.Context(),
		func(ue *model.UeInfo) error {
			//res := &track_msg.UEInfo{Laptop: ue}
			res := ue
			err := server.Send(res)
			if err != nil {
				return err
			}

			log.Info("sent ue with id: %d\n", ue.GetImsi())
			return nil
		},
	)
	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}
	return nil
}
