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

type PaymentService struct {
	pb.UnimplementedPaymentServiceServer
}

var allowedPaymentMethods = map[string]bool{
	"Credit Card": true, "Debit Card": true, "Cash": true,
}

var allowedPaymentStatuses = map[string]bool{
	"Pending": true, "Completed": true, "Failed": true,
}

func (s *PaymentService) CreatePayment(ctx context.Context, req *pb.Payment) (*pb.PaymentResponse, error) {
	if req.BookingId == 0 || req.Amount <= 0 || req.PaymentMethod == "" || req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "All fields are required")
	}

	if !allowedPaymentMethods[req.PaymentMethod] {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid payment method: %s. Allowed: Credit Card, Debit Card, Cash", req.PaymentMethod)
	}

	if !allowedPaymentStatuses[req.Status] {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid payment status: %s. Allowed: Pending, Completed, Failed", req.Status)
	}

	var booking models.Booking
	if err := database.DB.First(&booking, req.BookingId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Booking not found")
	}

	payment := models.Payment{
		BookingID:     uint(req.BookingId),
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Status:        req.Status,
	}

	database.DB.Create(&payment)

	return &pb.PaymentResponse{Message: "Payment created successfully"}, nil
}

func (s *PaymentService) GetPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.Payment, error) {
	var payment models.Payment

	if err := database.DB.First(&payment, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Payment not found")
	}

	return &pb.Payment{
		Id:            int32(payment.ID),
		BookingId:     int32(payment.BookingID),
		Amount:        payment.Amount,
		PaymentMethod: payment.PaymentMethod,
		Status:        payment.Status,
	}, nil
}

func (s *PaymentService) ListPayment(ctx context.Context, req *pb.Empty) (*pb.PaymentList, error) {
	var payments []models.Payment
	database.DB.Find(&payments)

	var pbPayments []*pb.Payment
	for _, payment := range payments {
		pbPayments = append(pbPayments, &pb.Payment{
			Id:            int32(payment.ID),
			BookingId:     int32(payment.BookingID),
			Amount:        payment.Amount,
			PaymentMethod: payment.PaymentMethod,
			Status:        payment.Status,
		})
	}

	return &pb.PaymentList{Payments: pbPayments}, nil
}

func (s *PaymentService) UpdatePayment(ctx context.Context, req *pb.Payment) (*pb.PaymentResponse, error) {
	var payment models.Payment

	if err := database.DB.First(&payment, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Payment not found")
	}

	if req.BookingId == 0 || req.Amount <= 0 || req.PaymentMethod == "" || req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "All fields are required")
	}

	if !allowedPaymentMethods[req.PaymentMethod] {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid payment method: %s. Allowed: Credit Card, Debit Card, Cash", req.PaymentMethod)
	}

	if !allowedPaymentStatuses[req.Status] {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid payment status: %s. Allowed: Pending, Completed, Failed", req.Status)
	}

	var booking models.Booking
	if err := database.DB.First(&booking, req.BookingId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Booking not found")
	}

	payment.BookingID = uint(req.BookingId)
	payment.Amount = req.Amount
	payment.PaymentMethod = req.PaymentMethod
	payment.Status = req.Status
	database.DB.Save(&payment)

	return &pb.PaymentResponse{Message: "Payment updated successfully"}, nil
}

func (s *PaymentService) DeletePayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	var payment models.Payment

	if err := database.DB.First(&payment, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Payment not found")
	}

	database.DB.Delete(&payment)

	return &pb.PaymentResponse{Message: "Payment deleted successfully"}, nil
}
