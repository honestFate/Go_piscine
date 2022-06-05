package randServer

import (
	"math/rand"
	api "randsig/pkg/api"
	"time"

	"github.com/google/uuid"
)

type GRPCServer struct {
	api.UnimplementedRandomSignalerServer
}

func (s *GRPCServer) RandSignal(req *api.RandSignalRequest, stream api.RandomSignaler_RandSignalServer) error {
	sessionId := uuid.New()
	rand.Seed(time.Now().UnixNano())
	mean := rand.Intn(21) - 10
	std := 0.3 + rand.Float64()*(1.5-0.3)
	for i := 0; i < 10; i++ {
		frequency := rand.NormFloat64()*std + float64(mean)
		if err := stream.Send(&api.RandSignalResponse{sessionId.String(), frequency, CurrentTimestamp: time.Now().String()}); err != nil {
			return err
		}
	}
	return nil
}
