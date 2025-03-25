package services

import (
	"carrental/database"
	"carrental/models"
	pb "carrental/pb"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type VehicleService struct {
	pb.UnimplementedVehicleServiceServer
}

var allowedVehicleTypes = map[string]bool{
	"Suzuki": true, "Toyota": true, "Honda": true, "Tata": true,
}
var allowedAvailability = map[string]bool{
	"Available": true, "Rented": true, "Maintenance": true,
}

func (s *VehicleService) CreateVehicle(ctx context.Context, req *pb.Vehicle) (*pb.VehicleResponse, error) {
	if req.Model == "" || req.Year == 0 || req.RentalRate == 0 || req.Availability == "" || req.Type == "" || req.Mileage == 0 {
		return nil, status.Error(codes.InvalidArgument, "All fields are required")
	}

	if !allowedVehicleTypes[req.Type] {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid vehicle type: %s. Allowed types: Suzuki, Toyota, Honda, Tata", req.Type)
	}

	if !allowedAvailability[req.Availability] {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid availability: %s. Allowed values: Available, Rented, Maintenance", req.Availability)
	}

	vehicle := models.Vehicle{
		Model:        req.Model,
		Year:         int(req.Year),
		RentalRate:   req.RentalRate,
		Availability: req.Availability,
		Type:         req.Type,
		Mileage:      int(req.Mileage),
	}

	database.DB.Create(&vehicle)

	return &pb.VehicleResponse{Message: "Vehicle created successfully"}, nil
}

func (s *VehicleService) GetVehicle(ctx context.Context, req *pb.VehicleRequest) (*pb.Vehicle, error) {
	var vehicle models.Vehicle

	if err := database.DB.First(&vehicle, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Vehicle not found")
	}

	return &pb.Vehicle{
		Id:           int32(vehicle.ID),
		Model:        vehicle.Model,
		Year:         int32(vehicle.Year),
		RentalRate:   vehicle.RentalRate,
		Availability: vehicle.Availability,
		Type:         vehicle.Type,
		Mileage:      int32(vehicle.Mileage),
	}, nil
}

func (s *VehicleService) ListVehicles(ctx context.Context, req *pb.Empty) (*pb.VehicleList, error) {
	var vehicles []models.Vehicle
	database.DB.Find(&vehicles)

	var pbVehicles []*pb.Vehicle
	for _, vehicle := range vehicles {
		pbVehicles = append(pbVehicles, &pb.Vehicle{
			Id:           int32(vehicle.ID),
			Model:        vehicle.Model,
			Year:         int32(vehicle.Year),
			RentalRate:   vehicle.RentalRate,
			Availability: vehicle.Availability,
			Type:         vehicle.Type,
			Mileage:      int32(vehicle.Mileage),
		})
	}

	return &pb.VehicleList{Vehicles: pbVehicles}, nil
}

func (s *VehicleService) UpdateVehicle(ctx context.Context, req *pb.Vehicle) (*pb.VehicleResponse, error) {
	var vehicle models.Vehicle

	if err := database.DB.First(&vehicle, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Vehicle not found")
	}

	if req.Model == "" || req.Year == 0 || req.RentalRate == 0 || req.Availability == "" || req.Type == "" || req.Mileage == 0 {
		return nil, status.Error(codes.InvalidArgument, "All fields are required")
	}

	if !allowedVehicleTypes[req.Type] {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid vehicle type: %s. Allowed types: Suzuki, Toyota, Honda, Tata", req.Type)
	}

	if !allowedAvailability[req.Availability] {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid availability: %s. Allowed values: Available, Rented, Maintenance", req.Availability)
	}

	vehicle.Model = req.Model
	vehicle.Year = int(req.Year)
	vehicle.RentalRate = req.RentalRate
	vehicle.Availability = req.Availability
	vehicle.Type = req.Type
	vehicle.Mileage = int(req.Mileage)
	database.DB.Save(&vehicle)

	return &pb.VehicleResponse{Message: "Vehicle updated successfully"}, nil
}

func (s *VehicleService) DeleteVehicle(ctx context.Context, req *pb.VehicleRequest) (*pb.VehicleResponse, error) {
	var vehicle models.Vehicle

	if err := database.DB.First(&vehicle, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Vehicle not found")
	}
	database.DB.Delete(&vehicle)

	return &pb.VehicleResponse{Message: "Vehicle deleted successfully"}, nil
}
