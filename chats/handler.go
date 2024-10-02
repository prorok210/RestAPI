package chats

import "RestAPI/core"

/*
docs(

	name: GetChat;
	tag: chats;
	path: /chats;
	method: GET;
	summary: Verify user;
	description: Verify user with the given data;
	isAuth: false;
	req_content_types: application/json, multipart/form-data;
	requestbody: {
		"mobile": "string",
		"otp": "string"
	};
	resp_content_type: application/json;
	responsebody: {
		"Message": "User verified"
	};

)docs
*/
func GetChat(request core.HttpRequest) core.HttpResponse {
	return core.HTTP200
}

/*
docs(

	name: VerifyUserHandler;
	tag: chats;
	path: /user/verify;
	method: GET;
	summary: Verify user;
	description: Verify user with the given data;
	isAuth: false;
	req_content_types: multipart/form-data;
	requestbody: {
		"mobile": "string",
		"otp": "string"
	};
	resp_content_type: application/json;
	responsebody: {
		"Message": "User verified"
	};

)docs
*/
