package handlers

import (
	"building-distributed-app-in-gin-chapter06/api/common"
	"building-distributed-app-in-gin-chapter06/api/dto"
	"building-distributed-app-in-gin-chapter06/api/models"
	"building-distributed-app-in-gin-chapter06/api/response"
	"building-distributed-app-in-gin-chapter06/api/util"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var jwtKey = []byte("JWT_SECRET")

type AuthHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewAuthHandler(ctx context.Context, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		collection: collection,
		ctx:        ctx,
	}
}

func (handler *AuthHandler) RegisterHandler(c *gin.Context) {
	var form models.Form
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name := form.Name
	telephone := form.Telephone
	password := form.Password
	inviteCode := form.InviteCode
	mail := form.Mail

	//数据验证
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	if inviteCode != "123321" {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "邀请码错误")
		return
	}
	//如果没有输入名称，则给一个10位的随机字符串
	log.Println(name)
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	log.Println(name)

	//查询是否存在
	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"telephone": telephone,
	})
	err := cur.Decode(&form)
	if err == nil {
		response.Response(c, http.StatusUnprocessableEntity, 423, nil, "用户已注册")
		return
	}

	id := primitive.NewObjectID()
	registeredAt := time.Now()
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "加密错误")
		return
	}

	//创建用户
	newForm := models.Form{
		ID:           id,
		Name:         name,
		Telephone:    telephone,
		Password:     string(hasedPassword),
		InviteCode:   inviteCode,
		Mail:         mail,
		RegisteredAt: registeredAt,
	}

	//写入数据库
	_, error := handler.collection.InsertOne(handler.ctx, newForm)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new user"})
		return
	}

	//发放token
	token, err := common.ReleaseToken(newForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统异常"})
		log.Printf("token generate error : %v", err)
		return
	}

	//返回结果
	response.Success(c, gin.H{"token": token, "name": telephone, "roles": "admin"}, "注册成功")
}

func (handler *AuthHandler) SignInHandler(c *gin.Context) {
	//var user models.User
	var form models.Form

	log.Println("0")

	//if err := c.ShouldBindJSON(&user); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//}
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//获取参数
	//telephone := user.Telephone
	//password := user.Password
	telephone := form.Telephone
	password := form.Password

	////数据验证
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}

	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}

	//查询是否存在
	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"telephone": telephone,
	})
	var dbuser models.Form
	err := cur.Decode(&dbuser)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "该手机号未注册")
		return
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(dbuser.Password), []byte(password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
		return
	}

	//发放token
	token, err := common.ReleaseToken(dbuser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "系统异常"})
		log.Printf("token generate error : %v", err)
		return
	}

	//返回结果
	response.Success(c, gin.H{"token": token, "name": dbuser.Name, "roles": "admin"}, "登陆成功")
}

func (handler *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取authorization header
		tokenString := ctx.GetHeader("Authorization")
		log.Println(tokenString)

		//validate token format
		//if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		//	ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足1"})
		//	ctx.Abort()
		//	return
		//}
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足1"})
			ctx.Abort()
			return
		}
		log.Println("0")

		//tokenString = tokenString[7:]

		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足2"})
			ctx.Abort()
			return
		}
		log.Println("1")

		//验证通过后获取claims中的userid
		userName := claims.UserName
		var user models.Form

		log.Println(userName)
		cur := handler.collection.FindOne(handler.ctx, bson.M{
			"name": userName,
		})
		err1 := cur.Decode(&user)
		if err1 != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足3"})
			ctx.Abort()
			return
		}
		log.Println("3")

		//如果用户存在，将user的信息写入上下文
		ctx.Set("user", user)
		ctx.Next()
	}
}

func (handler *AuthHandler) Info(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(models.User))}})
}

func (handler *AuthHandler) LogOutHandler(c *gin.Context) {
	response.Success(c, gin.H{}, "成功退出登录")
}

func (handler *AuthHandler) IntroductionHandler(c *gin.Context) {
	response.Success(c, gin.H{"introduction": "这里是DXW的博客"}, "介绍")
}

type tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (handler *AuthHandler) TagHandler(c *gin.Context) {
	golang := tag{1, "golang"}
	java := tag{2, "java"}
	mysql := tag{3, "mysql"}
	tags := make([]tag, 0)
	tags = append(tags, golang, java, mysql)
	response.Success(c, gin.H{"data": tags}, "tags")
}

//
//func (handler *AuthHandler) RefreshHandler(c *gin.Context) {
//	tokenValue := c.GetHeader("Authorization")
//
//	tokenValue = tokenValue[7:]
//	claims := &Claims{}
//	tkn, err := jwt.ParseWithClaims(tokenValue, claims, func(token *jwt.Token) (interface{}, error) {
//		return jwtKey, nil
//	})
//	if err != nil {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
//		return
//	}
//	if tkn == nil || !tkn.Valid {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
//		return
//	}
//	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Token is not expired yet"})
//		return
//	}
//	expirationTime := time.Now().Add(5 * time.Minute)
//	claims.ExpiresAt = expirationTime.Unix()
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	tokenString, err := token.SignedString(os.Getenv("JWT_SECRET"))
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	jwtOutput := JWTOutput{
//		Token:   tokenString,
//		Expires: expirationTime,
//	}
//	c.JSON(http.StatusOK, jwtOutput)
//}
