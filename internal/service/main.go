package service

import (
	"github.com/riazahmedshah/go-booking/internal/booking"
	"github.com/riazahmedshah/go-booking/internal/property"
	"github.com/riazahmedshah/go-booking/internal/repository"
	"github.com/riazahmedshah/go-booking/internal/server"
	"github.com/riazahmedshah/go-booking/internal/user"
)

type Service struct {
	UserService     *user.UserService
	PropertyService *property.PropertyService
	BookingService  *booking.BookingService
}

func NewService(server *server.Server, repository *repository.Repositories) (*Service, error) {
	userService := user.NewUserService(server, repository.UserRepository)
	propertyService := property.NewPropertyService(server, repository.PropertyRepository)
	bookingService := booking.NewBookingService(server, repository.BookingRepository)
	return &Service{
		UserService:     userService,
		PropertyService: propertyService,
		BookingService:  bookingService,
	}, nil
}
