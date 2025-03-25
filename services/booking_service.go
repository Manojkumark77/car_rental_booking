package services

import (
	"carrental/database"
	"carrental/models"
	pb "carrental/pb"
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var allowedBookingStatuses = map[string]bool{
	"Pending":   true,
	"Confirmed": true,
	"Cancelled": true,
	"Completed": true,
}

type BookingService struct {
	pb.UnimplementedBookingServiceServer
}

func (s *BookingService) CreateBooking(ctx context.Context, req *pb.Booking) (*pb.BookingResponse, error) {

	if req.CustomerId == 0 || req.VehicleId == 0 || req.StartDate == "" || req.EndDate == "" || req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "All fields are required")
	}

	if !allowedBookingStatuses[req.Status] {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid booking status: %s. Allowed statuses: Pending, Confirmed, Cancelled, Completed", req.Status)
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid start date format. Use YYYY-MM-DD")
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid end date format. Use YYYY-MM-DD")
	}
	if startDate.After(endDate) {
		return nil, status.Error(codes.InvalidArgument, "Start date cannot be after end date")
	}

	var customer models.Customer
	if err := database.DB.First(&customer, req.CustomerId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Customer not found")
	}

	var vehicle models.Vehicle
	if err := database.DB.First(&vehicle, req.VehicleId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Vehicle not found")
	}

	if vehicle.Availability != "Available" {
		return nil, status.Error(codes.NotFound, "Vehicle Already Booked")
	}

	booking := models.Booking{
		CustomerID: uint(req.CustomerId),
		VehicleID:  uint(req.VehicleId),
		StartDate:  startDate,
		EndDate:    endDate,
		Status:     req.Status,
	}
	database.DB.Create(&booking)

	vehicle.Availability = "Rented"
	database.DB.Save(&vehicle)

	return &pb.BookingResponse{Message: "Booking created successfully"}, nil
}

func (s *BookingService) GetBooking(ctx context.Context, req *pb.BookingRequest) (*pb.Booking, error) {
	var booking models.Booking

	if err := database.DB.First(&booking, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Booking not found")
	}

	return &pb.Booking{
		Id:         int32(booking.ID),
		CustomerId: int32(booking.CustomerID),
		VehicleId:  int32(booking.VehicleID),
		StartDate:  booking.StartDate.Format("2006-01-02"),
		EndDate:    booking.EndDate.Format("2006-01-02"),
		Status:     booking.Status,
	}, nil
}

func (s *BookingService) ListBookings(ctx context.Context, req *pb.Empty) (*pb.BookingList, error) {
	var bookings []models.Booking
	database.DB.Find(&bookings)

	var pbBookings []*pb.Booking
	for _, booking := range bookings {
		pbBookings = append(pbBookings, &pb.Booking{
			Id:         int32(booking.ID),
			CustomerId: int32(booking.CustomerID),
			VehicleId:  int32(booking.VehicleID),
			StartDate:  booking.StartDate.Format("2006-01-02"),
			EndDate:    booking.EndDate.Format("2006-01-02"),
			Status:     booking.Status,
		})
	}

	return &pb.BookingList{Bookings: pbBookings}, nil
}

func (s *BookingService) UpdateBooking(ctx context.Context, req *pb.Booking) (*pb.BookingResponse, error) {
	var booking models.Booking

	if err := database.DB.First(&booking, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Booking not found")
	}

	if req.CustomerId == 0 || req.VehicleId == 0 || req.StartDate == "" || req.EndDate == "" || req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "All fields are required")
	}

	if !allowedBookingStatuses[req.Status] {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid booking status: %s. Allowed statuses: Pending, Confirmed, Cancelled, Completed", req.Status)
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid start date format. Use YYYY-MM-DD")
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid end date format. Use YYYY-MM-DD")
	}
	if startDate.After(endDate) {
		return nil, status.Error(codes.InvalidArgument, "Start date cannot be after end date")
	}

	var customer models.Customer
	if err := database.DB.First(&customer, req.CustomerId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Customer not found")
	}

	var vehicle models.Vehicle
	if err := database.DB.First(&vehicle, req.VehicleId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Vehicle not found")
	}

	booking.CustomerID = uint(req.CustomerId)
	booking.VehicleID = uint(req.VehicleId)
	booking.StartDate = startDate
	booking.EndDate = endDate
	booking.Status = req.Status

	database.DB.Save(&booking)

	return &pb.BookingResponse{Message: "Booking updated successfully"}, nil
}

func (s *BookingService) DeleteBooking(ctx context.Context, req *pb.BookingRequest) (*pb.BookingResponse, error) {
	var booking models.Booking

	if err := database.DB.First(&booking, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Booking not found")
	}

	database.DB.Delete(&booking)

	return &pb.BookingResponse{Message: "Booking deleted successfully"}, nil
}
