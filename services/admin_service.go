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

type AdminService struct {
	pb.UnimplementedAdminServiceServer
}

func (s *AdminService) CreateAdmin(ctx context.Context, req *pb.Admin) (*pb.AdminResponse, error) {
	if req.Name == "" || req.Contact == "" || req.Role == "" || req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "All fields are required")
	}

	admin := models.Admin{
		Name:    req.Name,
		Contact: req.Contact,
		Role:    req.Role,
		Address: req.Address,
	}
	database.DB.Create(&admin)

	return &pb.AdminResponse{Message: "Admin created successfully"}, nil
}

func (s *AdminService) GetAdmin(ctx context.Context, req *pb.AdminRequest) (*pb.Admin, error) {
	var admin models.Admin
	if err := database.DB.First(&admin, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Admin not found")
	}

	return &pb.Admin{
		Id:      int32(admin.ID),
		Name:    admin.Name,
		Contact: admin.Contact,
		Role:    admin.Role,
		Address: admin.Address,
	}, nil
}

func (s *AdminService) ListAdmins(ctx context.Context, req *pb.Empty) (*pb.AdminList, error) {
	var admins []models.Admin
	database.DB.Find(&admins)

	var pbAdmins []*pb.Admin
	for _, admin := range admins {
		pbAdmins = append(pbAdmins, &pb.Admin{
			Id:      int32(admin.ID),
			Name:    admin.Name,
			Contact: admin.Contact,
			Role:    admin.Role,
			Address: admin.Address,
		})
	}

	return &pb.AdminList{Admins: pbAdmins}, nil
}

func (s *AdminService) UpdateAdmin(ctx context.Context, req *pb.Admin) (*pb.AdminResponse, error) {
	var admin models.Admin
	if err := database.DB.First(&admin, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Admin not found")
	}

	if req.Name == "" || req.Contact == "" || req.Role == "" || req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "All fields are required")
	}

	admin.Name = req.Name
	admin.Contact = req.Contact
	admin.Role = req.Role
	admin.Address = req.Address
	database.DB.Save(&admin)

	return &pb.AdminResponse{Message: "Admin updated successfully"}, nil
}

func (s *AdminService) DeleteAdmin(ctx context.Context, req *pb.AdminRequest) (*pb.AdminResponse, error) {
	var admin models.Admin
	if err := database.DB.First(&admin, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Admin not found")
	}
	database.DB.Delete(&admin)

	return &pb.AdminResponse{Message: "Admin deleted successfully"}, nil
}
