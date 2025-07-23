package zarinpal

import (
	"context"
	"github.com/amiranmanesh/paygap/client"
	"github.com/amiranmanesh/paygap/status"
	"google.golang.org/grpc/codes"
	"net/http"
)

const API_VERSION = "4"

const (
	ZARINPAL_HOST         = "https://api.zarinpal.com"
	ZARINPAL_SANDBOX_HOST = "https://sandbox.zarinpal.com"
)

const (
	ZARINPAL_REQUEST_API_ENDPOINT                = "/pg/v4/payment/request.json"
	ZARINPAL_VERIFY_API_ENDPOINT                 = "/pg/v4/payment/verify.json"
	ZARINPAL_UNVERIFIED_TRANSACTION_API_ENDPOINT = "/pg/v4/payment/unVerified.json"
)

// New create zarinpal provider object for user factory request methods
func New(client client.Transporter, merchantID string, sandbox bool) (*Zarinpal, error) {
	if client == nil {
		return nil, status.ERR_CLIENT_IS_NIL
	}

	zarinpal := &Zarinpal{
		client:             client,
		merchantID:         merchantID,
		baseUrl:            ZARINPAL_HOST,
		requestEndpoint:    ZARINPAL_REQUEST_API_ENDPOINT,
		verifyEndpoint:     ZARINPAL_VERIFY_API_ENDPOINT,
		unverifiedEndpoint: ZARINPAL_UNVERIFIED_TRANSACTION_API_ENDPOINT,
	}

	if sandbox {
		zarinpal.baseUrl = ZARINPAL_SANDBOX_HOST
	}

	if err := client.GetValidator().Struct(zarinpal); err != nil {
		return nil, status.New(0, http.StatusBadRequest, codes.InvalidArgument, err.Error())
	}

	return zarinpal, nil
}

// RequestPayment create payment request and return status code and authority
func (z *Zarinpal) RequestPayment(ctx context.Context, amount uint, callBackUrl, currency, description string, metaData map[string]interface{}) (*PaymentResponse, error) {
	req := &paymentRequest{
		merchantID:  z.merchantID,
		Amount:      amount,
		Currency:    currency,
		CallBackURL: callBackUrl,
		Description: description,
		MetaData:    metaData,
	}

	if err := z.client.GetValidator().Struct(req); err != nil {
		return nil, status.New(0, http.StatusBadRequest, codes.InvalidArgument, err.Error())
	}

	response := &PaymentResponse{}
	resp, err := z.client.Post(ctx, &client.APIConfig{Host: z.baseUrl, Path: z.requestEndpoint}, req)
	if err != nil {
		return nil, status.New(0, http.StatusInternalServerError, codes.Internal, err.Error())
	}

	if err := resp.GetJSON(response); err != nil {
		return nil, status.New(0, http.StatusInternalServerError, codes.Internal, err.Error())
	}

	return response, nil
}

// VerifyPayment transaction by merchant id, amount and authority to payment provider
func (z *Zarinpal) VerifyPayment(ctx context.Context, amount uint, authority string) (*VerifyResponse, error) {
	req := &verifyRequest{
		merchantID: z.merchantID,
		Amount:     amount,
		Authority:  authority,
	}

	if err := z.client.GetValidator().Struct(req); err != nil {
		return nil, status.New(0, http.StatusBadRequest, codes.InvalidArgument, err.Error())
	}

	response := &VerifyResponse{}
	resp, err := z.client.Post(ctx, &client.APIConfig{Host: z.baseUrl, Path: z.verifyEndpoint}, req)
	if err != nil {
		return nil, status.New(0, http.StatusInternalServerError, codes.Internal, err.Error())
	}

	if err := resp.GetJSON(response); err != nil {
		return nil, status.New(0, http.StatusInternalServerError, codes.Internal, err.Error())
	}

	return response, nil
}

// UnverifiedTransactions get unverified transactions from provider
func (z *Zarinpal) UnverifiedTransactions(ctx context.Context) (*UnverifiedTransactionsResponse, error) {
	req := &unverifiedTransactionsRequest{
		merchantID: z.merchantID,
	}

	if err := z.client.GetValidator().Struct(req); err != nil {
		return nil, status.New(0, http.StatusBadRequest, codes.InvalidArgument, err.Error())
	}

	response := &UnverifiedTransactionsResponse{}
	resp, err := z.client.Post(ctx, &client.APIConfig{Host: z.baseUrl, Path: z.unverifiedEndpoint}, req)
	if err != nil {
		return nil, status.New(0, http.StatusInternalServerError, codes.Internal, err.Error())
	}

	if err := resp.GetJSON(response); err != nil {
		return nil, status.New(0, http.StatusInternalServerError, codes.Internal, err.Error())
	}

	return response, nil
}

// FloatingShareSettlement a special method is used for sellers who benefit from an incoming amount,
// more information in https://docs.zarinpal.com/paymentGateway/setshare.html#%D8%AA%D8%B3%D9%88%DB%8C%D9%87-%D8%A7%D8%B4%D8%AA%D8%B1%D8%A7%DA%A9%DB%8C-%D8%B4%D9%86%D8%A7%D9%88%D8%B1
func (z *Zarinpal) FloatingShareSettlement(ctx context.Context, amount uint, description, callbackUrl string, wages []*Wages, metaData map[string]interface{}) (*FloatingShareSettlementResponse, error) {
	req := &floatingShareSettlementRequest{
		merchantID:  z.merchantID,
		Amount:      amount,
		CallBackURL: callbackUrl,
		Description: description,
		MetaData:    metaData,
		Wages:       wages,
	}

	if err := z.client.GetValidator().Struct(req); err != nil {
		return nil, status.New(0, http.StatusBadRequest, codes.InvalidArgument, err.Error())
	}

	response := &FloatingShareSettlementResponse{}
	resp, err := z.client.Post(ctx, &client.APIConfig{Host: z.baseUrl, Path: z.requestEndpoint}, req)
	if err != nil {
		return nil, status.New(0, http.StatusInternalServerError, codes.Internal, err.Error())
	}

	if err := resp.GetJSON(response); err != nil {
		return nil, status.New(0, http.StatusInternalServerError, codes.Internal, err.Error())
	}

	return response, nil
}

// VerifyFloatingShareSettlement verify floating share settlement
func (z *Zarinpal) VerifyFloatingShareSettlement(ctx context.Context, amount uint, authority string) (*VerifyFloatingShareSettlementResponse, error) {
	req := &verifyRequest{
		merchantID: z.merchantID,
		Amount:     amount,
		Authority:  authority,
	}

	if err := z.client.GetValidator().Struct(req); err != nil {
		return nil, status.New(0, http.StatusBadRequest, codes.InvalidArgument, err.Error())
	}

	response := &VerifyFloatingShareSettlementResponse{}
	resp, err := z.client.Post(ctx, &client.APIConfig{Host: z.baseUrl, Path: z.requestEndpoint}, req)
	if err != nil {
		return nil, status.New(0, http.StatusInternalServerError, codes.Internal, err.Error())
	}

	if err := resp.GetJSON(response); err != nil {
		return nil, status.New(0, http.StatusInternalServerError, codes.Internal, err.Error())
	}

	return response, nil
}
