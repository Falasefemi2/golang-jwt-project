// package helpers

// import (
// 	"errors"

// 	"github.com/gin-gonic/gin"
// )

// func CheckUserType(c *gin.Context, role string) (err error) {
// 	userType := c.GetString("user_type")
// 	err = nil
// 	if userType != role {
// 		err = errors.New("You are not authorized to perform this action")
// 		return err
// 	}
// 	return err
// }

// func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
// 	userType := c.GetStringMap("user_type")
// 	uid := c.GetString("uid")
// 	err = nil

// 	if userType == "USER" && uid != userId {
// 		err = errors.New("You are not authorized to perform this action")
// 		return err
// 	}

// 	err = CheckUserType(c, userType)
// 	return err
// }

package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userType, exists := c.Get("user_type")
	if !exists {
		err = errors.New("user_type not found in context")
		return
	}

	userTypeStr, ok := userType.(string)
	if !ok {
		err = errors.New("user_type is not a string")
		return
	}

	if userTypeStr != role {
		err = errors.New("you are not authorized to perform this action")
		return
	}

	return
}

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userTypeMap, exists := c.Get("user_type")
	if !exists {
		err = errors.New("user_type not found in context")
		return
	}

	userTypeMapValue, ok := userTypeMap.(map[string]interface{})
	if !ok {
		err = errors.New("user_type is not a map")
		return
	}

	userTypeStr, ok := userTypeMapValue["user_type"].(string)
	if !ok {
		err = errors.New("user_type value is not a string")
		return
	}

	uid, exists := c.Get("uid")
	if !exists {
		err = errors.New("uid not found in context")
		return
	}

	uidStr, ok := uid.(string)
	if !ok {
		err = errors.New("uid is not a string")
		return
	}

	if userTypeStr == "USER" && uidStr != userId {
		err = errors.New("you are not authorized to perform this action")
		return
	}

	err = CheckUserType(c, userTypeStr)
	return
}
