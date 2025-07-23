package idpay

import (
	"context"
	"errors"
	"github.com/amiranmanesh/paygap/client"
	"github.com/amiranmanesh/paygap/status"
	"google.golang.org/grpc/codes"
	"net/http"
	"reflect"
	"strconv"
)

const API_VERSION = "1.1"

const (
	IDPAY_HOST = "https://api.idpay.ir"

	PAYMENT_ENDPOINT      = "/v1.1/payment"
	VERIFY_ENDPOINT       = "/v1.1/payment/verify"
	INQUIRY_ENDPOINT      = "/v1.1/payment/inquiry"
	TRANSACTIONS_ENDPOINT = "/v1.1/payment/transactions"
)

type (
	CallBackPostFunc func(ctx context.Context, status int, trackId int, id string, orderId string, amount int, cardNo string, hashedCardNo string, date uint) error
	CallBackGetFunc  func(ctx context.Context, status int, trackId int, id string, orderId string) error
)

// New create idpay object for create new request
func New(client client.Transporter, apiKey string, sandbox bool) (*IdPay, error) {
	if client == nil {
		return nil, status.ERR_CLIENT_IS_NIL
	}

	idpay := new(IdPay)

	idpay.client = client
	idpay.apiKey = apiKey
	idpay.sandbox = sandbox
	idpay.baseUrl = IDPAY_HOST
	idpay.paymentEndpoint = PAYMENT_ENDPOINT
	idpay.verifyEndpoint = VERIFY_ENDPOINT
	idpay.inquiryEndpoint = INQUIRY_ENDPOINT
	idpay.transactionsEndpoint = TRANSACTIONS_ENDPOINT

	if err := client.GetValidator().Struct(idpay); err != nil {
		return nil, status.New(0, http.StatusBadRequest, codes.InvalidArgument, err.Error())
	}

	return idpay, nil
}

// CreateTransaction create a new transaction and return id and link
func (i *IdPay) CreateTransaction(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error) {
	if err := i.client.GetValidator().Struct(req); err != nil {
		return nil, status.New(0, http.StatusBadRequest, codes.InvalidArgument, err.Error())
	}
	return request[*PaymentRequest, *PaymentResponse](ctx, i, req, i.baseUrl, i.paymentEndpoint, nil)
}

// VerifyTransaction verify an transaction base on id and order id
func (i *IdPay) VerifyTransaction(ctx context.Context, req *VerifyRequest) (*VerifyResponse, error) {
	if err := i.client.GetValidator().Struct(req); err != nil {
		return nil, status.New(0, http.StatusBadRequest, codes.InvalidArgument, err.Error())
	}

	return request[*VerifyRequest, *VerifyResponse](ctx, i, req, i.baseUrl, i.verifyEndpoint, nil)
}

// TransactionStatus check transaction status and return transaction details
func (i *IdPay) TransactionStatus(ctx context.Context, req *TransactionStatusRequest) (*TransactionStatusResponse, error) {
	if err := i.client.GetValidator().Struct(req); err != nil {
		return nil, status.New(0, http.StatusBadRequest, codes.InvalidArgument, err.Error())
	}

	return request[*TransactionStatusRequest, *TransactionStatusResponse](ctx, i, req, i.baseUrl, i.inquiryEndpoint, nil)
}

// TransactionList get list of transaction with set page size and page
func (i *IdPay) TransactionList(ctx context.Context, req *TransactionListRequest) (*TransactionListResponse, error) {
	if err := i.client.GetValidator().Struct(req); err != nil {
		return nil, status.New(0, http.StatusBadRequest, codes.InvalidArgument, err.Error())
	}

	query := make(map[string]string)

	if req.PageSize != 0 {
		query["page"] = strconv.Itoa(req.Page)
		query["page_size"] = strconv.Itoa(req.PageSize)
	}

	return request[*TransactionListRequest, *TransactionListResponse](ctx, i, req, i.baseUrl, i.transactionsEndpoint, query)
}

func request[RQ any, RS any](ctx context.Context, i *IdPay, req RQ, baseUrl string, endpoint string, queryParams map[string]string) (response RS, err error) {
	r, ok := reflect.New(reflect.TypeOf(response).Elem()).Interface().(RS)
	if !ok {
		return response, errors.New("response type is invalid")
	}

	headers := make(map[string]string)
	headers["X-API-KEY"] = i.apiKey
	headers["Content-Type"] = "application/json"

	if i.sandbox {
		headers["X-SANDBOX"] = "1"
	}

	errResp := &ErrorResponse{}
	resp, err := i.client.Post(ctx, &client.APIConfig{Host: baseUrl, Path: endpoint, Headers: headers, Query: queryParams}, req)
	if err != nil {
		return response, status.New(0, http.StatusInternalServerError, codes.Internal, err.Error())
	}

	if resp.GetHttpResponse().StatusCode != http.StatusOK|http.StatusCreated {
		if err := resp.GetJSON(errResp); err != nil {
			return response, status.New(0, http.StatusInternalServerError, codes.Internal, err.Error())
		}
		return response, status.New(errResp.ErrorCode, http.StatusFailedDependency, codes.OK, errResp.ErrorMessage)
	}

	if err := resp.GetJSON(r); err != nil {
		return response, status.New(0, http.StatusInternalServerError, codes.Internal, err.Error())
	}

	return r, nil
}
