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

type ReviewService struct {
	pb.UnimplementedReviewServiceServer
}

func (s *ReviewService) CreateReview(ctx context.Context, req *pb.Review) (*pb.ReviewResponse, error) {
	if req.CustomerId == 0 || req.VehicleId == 0 || req.Rating < 1 || req.Rating > 5 {
		return nil, status.Error(codes.InvalidArgument, "CustomerID, VehicleID, and Rating (1-5) are required")
	}

	var customer models.Customer
	if err := database.DB.First(&customer, req.CustomerId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Customer not found")
	}

	var vehicle models.Vehicle
	if err := database.DB.First(&vehicle, req.VehicleId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Vehicle not found")
	}

	review := models.Review{
		CustomerID: uint(req.CustomerId),
		VehicleID:  uint(req.VehicleId),
		Rating:     int(req.Rating),
		Comments:   req.Comments,
	}

	database.DB.Create(&review)

	return &pb.ReviewResponse{Message: "Review created successfully"}, nil
}

func (s *ReviewService) GetReview(ctx context.Context, req *pb.ReviewRequest) (*pb.Review, error) {
	var review models.Review

	if err := database.DB.First(&review, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Review not found")
	}

	return &pb.Review{
		Id:         int32(review.ID),
		CustomerId: int32(review.CustomerID),
		VehicleId:  int32(review.VehicleID),
		Rating:     int32(review.Rating),
		Comments:   review.Comments,
	}, nil
}

func (s *ReviewService) ListReviews(ctx context.Context, req *pb.Empty) (*pb.ReviewList, error) {
	var reviews []models.Review
	database.DB.Find(&reviews)

	var pbReviews []*pb.Review
	for _, review := range reviews {
		pbReviews = append(pbReviews, &pb.Review{
			Id:         int32(review.ID),
			CustomerId: int32(review.CustomerID),
			VehicleId:  int32(review.VehicleID),
			Rating:     int32(review.Rating),
			Comments:   review.Comments,
		})
	}

	return &pb.ReviewList{Reviews: pbReviews}, nil
}

func (s *ReviewService) UpdateReview(ctx context.Context, req *pb.Review) (*pb.ReviewResponse, error) {
	var review models.Review

	if err := database.DB.First(&review, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Review not found")
	}

	if req.CustomerId == 0 || req.VehicleId == 0 || req.Rating < 1 || req.Rating > 5 {
		return nil, status.Error(codes.InvalidArgument, "CustomerID, VehicleID, and Rating (1-5) are required")
	}

	review.CustomerID = uint(req.CustomerId)
	review.VehicleID = uint(req.VehicleId)
	review.Rating = int(req.Rating)
	review.Comments = req.Comments
	database.DB.Save(&review)

	return &pb.ReviewResponse{Message: "Review updated successfully"}, nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, req *pb.ReviewRequest) (*pb.ReviewResponse, error) {
	var review models.Review

	if err := database.DB.First(&review, req.Id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.NotFound, "Review not found")
	}

	database.DB.Delete(&review)

	return &pb.ReviewResponse{Message: "Review deleted successfully"}, nil
}
