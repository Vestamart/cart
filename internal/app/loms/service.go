package loms

import (
	"context"
	"errors"
	"github.com/vestamart/homework/internal/domain"
	"github.com/vestamart/homework/internal/repository"
	desc "github.com/vestamart/homework/pkg/api/loms/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OrdersRepository и StocksStorage интерфейсы для взаимодействия с репозиториями
//
//go:generate minimock -i github.com/vestamart/homework/internal/app/loms.OrdersRepository -o ./mock/orders_repository_mock.go -n OrdersRepositoryMock -p mock
type OrdersRepository interface {
	Create(_ context.Context, userID int64, items *[]domain.Item) (int64, error)
	SetStatus(_ context.Context, orderID int64, status domain.OrderStatus) error
	GetByID(_ context.Context, orderID int64) (*domain.Order, error)
}

//go:generate minimock -i github.com/vestamart/homework/internal/app/loms.StocksStorage -o ./mock/stock_repository_mock.go -n StocksStorageMock -p mock
type StocksStorage interface {
	Reserve(_ context.Context, sku uint32, count uint32) error
	ReserveRemove(_ context.Context, sku uint32, count uint32) error
	ReserveCancel(_ context.Context, skus map[uint32]uint32) error
	GetBySKU(_ context.Context, sku uint32) (uint32, error)
	RollbackReserve(_ context.Context, skus map[uint32]uint32) error
}

type Service struct {
	desc.UnimplementedLomsServer
	ordersRepository OrdersRepository
	stocksRepository StocksStorage
}

func NewService(ordersRepository OrdersRepository, stocksRepository StocksStorage) *Service {
	return &Service{ordersRepository: ordersRepository, stocksRepository: stocksRepository}
}

func (s Service) OrderCreate(ctx context.Context, request *desc.OrderCreateRequest) (*desc.OrderCreateResponse, error) {
	if request == nil || (request.User == 0 || len(request.Items) == 0) {
		return nil, status.Errorf(codes.InvalidArgument, "empty order create request")
	}

	items := make([]domain.Item, 0, len(request.Items))
	for _, v := range request.Items {
		items = append(items, domain.Item{
			Sku:   v.Sku,
			Count: v.Count,
		})
	}

	orderId, err := s.ordersRepository.Create(ctx, request.User, &items)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var reservedSKUs = make(map[uint32]uint32)
	for _, v := range items {
		err = s.stocksRepository.Reserve(ctx, v.Sku, v.Count)
		if err != nil {
			_ = s.stocksRepository.RollbackReserve(ctx, reservedSKUs)
			if errStatus := s.ordersRepository.SetStatus(ctx, orderId, 2); errStatus != nil {
				return nil, status.Error(codes.Internal, errStatus.Error())
			}
			if errors.Is(err, repository.SKUNotExistErr) {
				return nil, status.Errorf(codes.NotFound, "SKU %d not found", v.Sku)
			}
			if errors.Is(err, domain.ItemNotEnoughErr) {
				return nil, status.Errorf(codes.ResourceExhausted, "SKU %d is not enough", v.Sku)
			}
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		reservedSKUs[v.Sku] = v.Count
	}
	if err = s.ordersRepository.SetStatus(ctx, orderId, 1); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.OrderCreateResponse{OrderId: orderId}, status.Errorf(codes.OK, "")
}

func (s Service) OrderInfo(ctx context.Context, request *desc.OrderInfoRequest) (*desc.OrderInfoResponse, error) {
	rawResponse, err := s.ordersRepository.GetByID(ctx, request.OrderId)
	if err != nil {
		if errors.Is(err, repository.OrderNotExistErr) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, err
	}
	items := make([]*desc.Item, 0, len(rawResponse.Items))
	for _, v := range rawResponse.Items {
		items = append(items, &desc.Item{
			Sku:   v.Sku,
			Count: v.Count,
		})
	}

	response := &desc.OrderInfoResponse{
		Status: desc.OrderStatus(rawResponse.Status),
		User:   rawResponse.UserID,
		Items:  items,
	}

	return response, status.Error(codes.OK, "")
}

func (s Service) OrderPay(ctx context.Context, request *desc.OrderPayRequest) (*desc.OrderPayResponse, error) {
	getByID, err := s.ordersRepository.GetByID(ctx, request.OrderID)
	if err != nil {
		if errors.Is(err, repository.OrderNotExistErr) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, v := range getByID.Items {
		if err = s.stocksRepository.ReserveRemove(ctx, v.Sku, v.Count); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	err = s.ordersRepository.SetStatus(ctx, request.OrderID, 3)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &desc.OrderPayResponse{}, status.Error(codes.OK, "")
}

func (s Service) OrderCancel(ctx context.Context, request *desc.OrderCancelRequest) (*desc.OrderCancelResponse, error) {
	rawResponse, err := s.ordersRepository.GetByID(ctx, request.OrderID)
	if err != nil {
		if errors.Is(err, repository.OrderNotExistErr) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	items := make(map[uint32]uint32)
	for _, v := range rawResponse.Items {
		items[v.Sku] = v.Count
	}

	if err = s.stocksRepository.ReserveCancel(ctx, items); err != nil {
		if errors.Is(err, repository.OrderNotExistErr) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err = s.ordersRepository.SetStatus(ctx, request.OrderID, 4); err != nil {
	}
	return &desc.OrderCancelResponse{}, status.Error(codes.OK, "")
}
func (s Service) StocksInfo(ctx context.Context, request *desc.StocksInfoRequest) (*desc.StocksInfoResponse, error) {
	v, err := s.stocksRepository.GetBySKU(ctx, request.Sku)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &desc.StocksInfoResponse{Count: uint64(v)}, status.Error(codes.OK, "")
}
