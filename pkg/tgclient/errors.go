package tgclient

import "github.com/joomcode/errorx"

var Errors = errorx.NewNamespace("tg_errors")
var ParseErr = Errors.NewType("parse")
var TimeoutErr = Errors.NewType("timeout")
var AuthErr = Errors.NewType("auth")
var RequestErr = Errors.NewType("request")
