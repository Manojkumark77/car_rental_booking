package main

import (
	"carrental/database"
	pb "carrental/pb"
	"carrental/services"
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func main() {
	database.ConnectDatabase()

	startHTTPServer()
}

func startHTTPServer() {
	mux := runtime.NewServeMux()
	ctx := context.Background()

	err := pb.RegisterCustomerServiceHandlerServer(ctx, mux, &services.CustomerService{})
	if err != nil {
		log.Fatalf("Failed to register gRPC Gateway: %v", err)
	}

	err = pb.RegisterAdminServiceHandlerServer(ctx, mux, &services.AdminService{})
	if err != nil {
		log.Fatalf("Failed to register AdminService gRPC Gateway: %v", err)
	}

	err = pb.RegisterVehicleServiceHandlerServer(ctx, mux, &services.VehicleService{})
	if err != nil {
		log.Fatalf("Failed to register AdminService gRPC Gateway: %v", err)
	}

	err = pb.RegisterBookingServiceHandlerServer(ctx, mux, &services.BookingService{})
	if err != nil {
		log.Fatalf("Failed to register AdminService gRPC Gateway: %v", err)
	}

	err = pb.RegisterPaymentServiceHandlerServer(ctx, mux, &services.PaymentService{})
	if err != nil {
		log.Fatalf("Failed to register AdminService gRPC Gateway: %v", err)
	}

	err = pb.RegisterReviewServiceHandlerServer(ctx, mux, &services.ReviewService{})
	if err != nil {
		log.Fatalf("Failed to register AdminService gRPC Gateway: %v", err)
	}

	log.Println("REST API running on port 8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
