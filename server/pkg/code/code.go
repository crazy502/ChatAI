package code

type Code int64

const (
	CodeSuccess Code = 1000

	CodeInvalidParams    Code = 2001
	CodeUserExist        Code = 2002
	CodeUserNotExist     Code = 2003
	CodeInvalidPassword  Code = 2004
	CodeNotMatchPassword Code = 2005
	CodeInvalidToken     Code = 2006
	CodeNotLogin         Code = 2007
	CodeInvalidCaptcha   Code = 2008
	CodeRecordNotFound   Code = 2009
	CodeIllegalPassword  Code = 2010
	CodeEmailExist       Code = 2011
	CodeTooManyRequests  Code = 2012

	CodeForbidden Code = 3001

	CodeServerBusy Code = 4001

	AIModelNotFind    Code = 5001
	AIModelCannotOpen Code = 5002
	AIModelFail       Code = 5003
)

var msg = map[Code]string{
	CodeSuccess:          "success",
	CodeInvalidParams:    "\u8bf7\u6c42\u53c2\u6570\u9519\u8bef",
	CodeUserExist:        "\u7528\u6237\u5df2\u5b58\u5728",
	CodeUserNotExist:     "\u7528\u6237\u4e0d\u5b58\u5728",
	CodeInvalidPassword:  "\u7528\u6237\u540d\u6216\u5bc6\u7801\u9519\u8bef",
	CodeNotMatchPassword: "\u4e24\u6b21\u5bc6\u7801\u4e0d\u4e00\u81f4",
	CodeInvalidToken:     "\u65e0\u6548\u7684 Token",
	CodeNotLogin:         "\u7528\u6237\u672a\u767b\u5f55",
	CodeInvalidCaptcha:   "\u9a8c\u8bc1\u7801\u9519\u8bef",
	CodeRecordNotFound:   "\u8bb0\u5f55\u4e0d\u5b58\u5728",
	CodeIllegalPassword:  "\u5bc6\u7801\u4e0d\u5408\u6cd5",
	CodeEmailExist:       "\u90ae\u7bb1\u5df2\u5b58\u5728",
	CodeTooManyRequests:  "\u8bf7\u6c42\u8fc7\u4e8e\u9891\u7e41\uff0c\u8bf7\u7a0d\u540e\u91cd\u8bd5",
	CodeForbidden:        "\u6743\u9650\u4e0d\u8db3",
	CodeServerBusy:       "\u670d\u52a1\u7e41\u5fd9\uff0c\u8bf7\u7a0d\u540e\u91cd\u8bd5",
	AIModelNotFind:       "\u6a21\u578b\u4e0d\u5b58\u5728",
	AIModelCannotOpen:    "\u6a21\u578b\u521d\u59cb\u5316\u5931\u8d25",
	AIModelFail:          "\u6a21\u578b\u8c03\u7528\u5931\u8d25",
}

func (c Code) Code() int64 {
	return int64(c)
}

func (c Code) Msg() string {
	if text, ok := msg[c]; ok {
		return text
	}
	return msg[CodeServerBusy]
}
