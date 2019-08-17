package tgbot

import "github.com/joomcode/errorx"

var Errors = errorx.NewNamespace("tgbot")
var ReqErr = Errors.NewType("request")
var BuilderErr = Errors.NewType("builder")

func newReqError(err error, method string, req request) error {
	return ReqErr.Wrap(err, "request failed. method: %s, params: %s", method, req.String())
}

func newApiError(err error, method string, req request, resp response) error {
	return ReqErr.Wrap(err, "request api failed. code: %d, description: %s, method: %s, req: %s",
		resp.ErrorCode, resp.Description, method, req.String())
}

