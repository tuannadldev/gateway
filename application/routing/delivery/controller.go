package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	r9 "github.com/redis/go-redis/v9"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"

	"gateway/application/model"
	"gateway/application/routing"
	"gateway/config"
	"gateway/pkg/grpc_errors"
	"gateway/pkg/utils"
)

type routingHandler struct {
	routingUC routing.RoutingUseCase
}

func RegisterRoutes(r *gin.Engine, c *config.Config, rdb *r9.Client, routingUC routing.RoutingUseCase) {
	handler := &routingHandler{
		routingUC: routingUC,
	}

	registerRoutes(r, handler.handle(c, rdb))
}

func (handler *routingHandler) handle(config *config.Config, rdb *r9.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		route := ctx.Request.Method + ":" + ctx.FullPath()
		routing_config, found := routingRegistry[route]
		if !found {

			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "An unexpected error has occurred. Please retry your request later"})
			return
		}
		metadata := make(map[string]string)
		// Begin: check auth ------------------------------------------------------------------------
		// Begin: check auth
		if !routing.ApiPublic[ctx.FullPath()] {
			auth_header := strings.Split(ctx.Request.Header.Get("Authorization"), "Bearer ")
			if len(auth_header) != 2 {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "API token required"})
				return
			} else {
				jwtToken := strings.TrimSpace(auth_header[1])
				if len(jwtToken) == 0 {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid API token"})
					return
				}
				if !handler.checkHardcodeApi(ctx.FullPath(), jwtToken, config) {
					claims, err := utils.ValidateToken(jwtToken, config)
					if err != nil {
						ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "401 Unauthorized"})
						return
					}

					// Check token active in redis
					err = rdb.Get(ctx, fmt.Sprintf("user::token::%d", claims.UserID)).Err()
					if err != nil {
						ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Your token has been revoked"})
						return
					}

					ctx.Set("user_login_id", claims.UserID)
					metadata["uid"] = strconv.FormatInt(claims.UserID, 10)
					if routing_config.remoteServicePermission != "" {
						// permission_codes := []interface{}{routing_config.remoteServicePermission}
						// arr_code := strings.Split(routing_config.remoteServicePermission, "::")
						// if len(arr_code) != 2 {
						// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Permission code invalid"})
						// 	return
						// }
						// if arr_code[1] != "all" {
						// 	permission_codes = append(permission_codes, arr_code[0]+"::all")
						// }
						ok, err := rdb.SIsMember(ctx, fmt.Sprintf("auth:user:privilege:userid:%d", claims.UserID), routing_config.remoteServicePermission).Result()
						if err != nil || !ok {
							ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Permission denied"})
							return
						}
					}
				}
			}
		}

		// End: check auth ------------------------------------------------------------------------

		data, err := routing_config.requestHandler.Handle(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		routing_data := &model.RoutingData{
			ServiceName:   routing_config.remoteServiceName,
			ServiceMethod: routing_config.remoteServiceMethod,
			Payload:       data,
			Metadata:      metadata,
		}
		fmt.Println(routing_data)
		res, err := handler.routingUC.Forward(routing_data)
		if err != nil {
			err_status, _ := status.FromError(err)
			map_code := grpc_errors.MapGRPCErrCodeToHttpStatus(err_status.Code())
			ctx.AbortWithStatusJSON(map_code, gin.H{"message": err_status.Message()})
			return
		}

		ctx.JSON(http.StatusOK, res)
	}
}

func (handler *routingHandler) checkHardcodeApi(url string, jwtToken string, config *config.Config) bool {
	if routing.ApiHardCodeToken[url] {
		if jwtToken == string(config.SecretToken) {
			return true
		}
	}
	return false
}
