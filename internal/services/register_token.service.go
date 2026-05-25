package services

import (
	"context"
	"fmt"
	"time"

	firestore "cloud.google.com/go/firestore"
	"github.com/krisnaganesha1609/IoTDrainage-BE/internal/requests"
)

func (s *Service) RegisterToken(ctx context.Context, request *requests.RegisterTokenRequest) error {

	docRef := s.Firebase.Firestore.Collection("devices").Doc(request.DeviceID)

	_, err := docRef.Set(ctx, map[string]interface{}{
		"mobile_devices": map[string]interface{}{
			request.AppInstanceID: map[string]interface{}{
				"fcm_token":  request.FCMToken,
				"updated_at": time.Now(),
			},
		},
	}, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("Gagal menyimpan token ke Firestore: %v", err)
	}
	return nil
}
