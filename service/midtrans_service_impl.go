package service

import (
	"be-go-fiber-ecommerce/entity"
	"be-go-fiber-ecommerce/model"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransServiceImpl struct {
	DB *gorm.DB
}

func NewMidTransServiceImpl(db *gorm.DB) *MidtransServiceImpl {
	return &MidtransServiceImpl{
		DB: db,
	}
}

func (service *MidtransServiceImpl) CartTransaction(c *fiber.Ctx) (model.MidtransResponse, error) {
	userId := c.Locals("userID").(uint)

	var user entity.User
	if err := service.DB.Preload("Cart.Items.Product.Category").First(&user, userId).Error; err != nil {
		return model.MidtransResponse{}, fmt.Errorf("user not found or Cart haven't initialized")
	}

	tx := service.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var totalGross int64
	itemsDetails := make([]midtrans.ItemDetails, 0)
	for _, item := range user.Cart.Items {
		totalGross += int64(item.Quantity) * int64(item.Product.Price)

		itemsDetails = append(itemsDetails, midtrans.ItemDetails{
			ID:       fmt.Sprintf("Property-%d", item.ProductID),
			Name:     item.Product.Name,
			Price:    int64(item.Product.Price),
			Qty:      int32(item.Quantity),
			Category: item.Product.Category.Name,
		})

		newStock := item.Product.Stock - item.Quantity
		if newStock < 0 {
			tx.Rollback()
			return model.MidtransResponse{}, fmt.Errorf("not enough stock for product: %s", item.Product.Name)
		}

		if err := tx.Model(&entity.Product{}).Where("id = ?", item.ProductID).Update("stock", newStock).Error; err != nil {
			tx.Rollback()
			return model.MidtransResponse{}, fmt.Errorf("update error (product: %s): %w", item.Product.Name, err)
		}
	}

	orderSuffix := fmt.Sprintf("%d", time.Now().UnixNano())

	var snapClient = snap.Client{}
	snapClient.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)

	custAddress := &midtrans.CustomerAddress{
		FName:       extractEmailUsername(user.Email),
		LName:       "Synapsis",
		Phone:       "085311111010",
		Address:     "St. Kerto No. 4",
		City:        "Yogyakarta",
		Postcode:    "55165",
		CountryCode: "IDN",
	}

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  fmt.Sprintf("MID-User%d-Order-%s", userId, orderSuffix),
			GrossAmt: totalGross,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName:    extractEmailUsername(user.Email),
			LName:    "Synapsis",
			Email:    "info@synapsis.id",
			Phone:    "085311111010",
			BillAddr: custAddress,
			ShipAddr: custAddress,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		Items:           &itemsDetails,
	}

	response, errSnap := snapClient.CreateTransaction(req)
	if errSnap != nil {
		tx.Rollback()
		return model.MidtransResponse{}, fmt.Errorf(errSnap.Message)
	}

	if err := tx.Where("cart_id = ?", user.Cart.ID).Delete(&entity.CartItem{}).Error; err != nil {
		tx.Rollback()
		return model.MidtransResponse{}, fmt.Errorf("error clearing cart items")
	}

	if err := tx.Commit().Error; err != nil {
		return model.MidtransResponse{}, fmt.Errorf("transaction commit error: %w", err)
	}

	midtransResponse := model.MidtransResponse{
		Token:       response.Token,
		RedirectUrl: response.RedirectURL,
	}

	return midtransResponse, nil
}
