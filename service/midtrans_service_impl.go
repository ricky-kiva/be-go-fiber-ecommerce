package service

import (
	"be-go-fiber-ecommerce/models"
	"be-go-fiber-ecommerce/models/web"
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

func (service *MidtransServiceImpl) CartTransaction(c *fiber.Ctx) (web.MidtransResponse, error) {
	var totalGross int64
	userId := c.Locals("userID").(uint)

	var user models.User
	if err := service.DB.Preload("Cart.Items.Product.Category").First(&user, userId).Error; err != nil {
		return web.MidtransResponse{}, fmt.Errorf("user not found or Cart haven't initialized")
	}

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
		return web.MidtransResponse{}, fmt.Errorf(errSnap.Message)
	}

	if err := service.DB.Where("cart_id = ?", user.Cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		return web.MidtransResponse{}, fmt.Errorf("error clearing cart items")
	}

	midtransResponse := web.MidtransResponse{
		Token:       response.Token,
		RedirectUrl: response.RedirectURL,
	}

	return midtransResponse, nil
}
