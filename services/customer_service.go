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

type CustomerService struct {
	pb.UnimplementedCustomerServiceServer
}

func (s *CustomerService) CreateCustomer(ctx context.Context, req *pb.Customer) (*pb.CustomerResponse, error) {
	if req.Name == "" || req.Contact == "" || req.Address == "" || req.LicenseNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "All fields are required")
	}

	var existingCustomer models.Customer
	if err := database.DB.Where("license_number = ?", req.LicenseNumber).First(&existingCustomer).Error; err == nil {
		return nil, status.Error(codes.AlreadyExists, "Customer with this License Number already exists")
	}

	customer := models.Customer{
		Name:          req.Name,
		Contact:       req.Contact,
		Address:       req.Address,
		LicenseNumber: req.LicenseNumber,
	}
	database.DB.Create(&customer)

	return &pb.CustomerResponse{Message: "Customer created successfully"}, nil
}

func (s *CustomerService) GetCustomer(ctx context.Context, req *pb.CustomerRequest) (*pb.Customer, error) {
	var customer models.Customer
	if err := database.DB.First(&customer, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Customer not found")
	}

	return &pb.Customer{
		Id:            int32(customer.ID),
		Name:          customer.Name,
		Contact:       customer.Contact,
		Address:       customer.Address,
		LicenseNumber: customer.LicenseNumber,
	}, nil
}

func (s *CustomerService) ListCustomers(ctx context.Context, req *pb.Empty) (*pb.CustomerList, error) {
	var customers []models.Customer
	database.DB.Find(&customers)

	var pbCustomers []*pb.Customer
	for _, customer := range customers {
		pbCustomers = append(pbCustomers, &pb.Customer{
			Id:            int32(customer.ID),
			Name:          customer.Name,
			Contact:       customer.Contact,
			Address:       customer.Address,
			LicenseNumber: customer.LicenseNumber,
		})
	}

	return &pb.CustomerList{Customers: pbCustomers}, nil
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, req *pb.Customer) (*pb.CustomerResponse, error) {
	var customer models.Customer
	if err := database.DB.First(&customer, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Customer not found")
	}

	if req.Name == "" || req.Contact == "" || req.Address == "" || req.LicenseNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "All fields are required")
	}

	customer.Name = req.Name
	customer.Contact = req.Contact
	customer.Address = req.Address
	customer.LicenseNumber = req.LicenseNumber
	database.DB.Save(&customer)

	return &pb.CustomerResponse{Message: "Customer updated successfully"}, nil
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, req *pb.CustomerRequest) (*pb.CustomerResponse, error) {
	var customer models.Customer
	if err := database.DB.First(&customer, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Customer not found")
	}
	database.DB.Delete(&customer)

	return &pb.CustomerResponse{Message: "Customer deleted successfully"}, nil
}
