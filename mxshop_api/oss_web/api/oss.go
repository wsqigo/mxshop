package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mxshop_api/oss_web/utils"
	"net/http"
)

func GetToken(ctx *gin.Context) {
	response := utils.GetPolicyToken()
	ctx.Header("Access-Control-Allow-Methods", "POST")
	ctx.Header("Access-Control-Allow-Origin", "*")

	ctx.String(http.StatusOK, response)
}

func HandlerRequest(ctx *gin.Context) {
	fmt.Println("\nHandle Post Request ... ")

	// Get PublicKey bytes
	bytePublicKey, err := utils.GetPublicKey(ctx)
	if err != nil {
		utils.ResponseFailed(ctx)
		return
	}

	// Get Authorization bytes : decode from Base64String
	byteAuthorization, err := utils.GetAuthorization(ctx)
	if err != nil {
		utils.ResponseFailed(ctx)
		return
	}

	// Get MD5 bytes from Newly Constructed Authrization String.
	byteMD5, err := utils.GetMD5FromNewAuthString(ctx)
	if err != nil {
		utils.ResponseFailed(ctx)
		return
	}

	// verifySignature and response to client
	if utils.VerifySignature(bytePublicKey, byteMD5, byteAuthorization) {
		// do something you want accoding to callback_body ...

		utils.ResponseSuccess(ctx) // response OK : 200
	} else {
		utils.ResponseFailed(ctx) // response FAILED : 400
	}
}
