package handlers

import (
	"building-distributed-app-in-gin-chapter06/api/common"
	"building-distributed-app-in-gin-chapter06/api/models"
	"building-distributed-app-in-gin-chapter06/api/response"
	"building-distributed-app-in-gin-chapter06/api/vo"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostController struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewPostController(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *PostController {
	return &PostController{
		collection:  collection,
		ctx:         ctx,
		redisClient: redisClient,
	}
}

func (p PostController) Create(c *gin.Context) {
	var requestPost vo.CreatePostRequest
	if err := c.ShouldBindJSON(&requestPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post = models.Post{}
	post.Title = requestPost.Title
	post.Content = requestPost.Content
	post.CategoryId, _ = strconv.Atoi(requestPost.CategoryId)
	post.HeadImg = requestPost.HeadImg
	post.ID = primitive.NewObjectID()
	post.PublishedAt = time.Now()
	_, err := p.collection.InsertOne(p.ctx, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new post"})
		return
	}

	log.Println("Remove data from Redis")
	p.redisClient.Del("posts")

	c.JSON(http.StatusOK, post)
}

func (p PostController) Update(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := p.collection.UpdateOne(p.ctx, bson.M{
		"_id": objectId,
	}, bson.D{{"$set", bson.D{
		{"title", post.Title},
		{"content", post.Content},
	}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response.Success(c, gin.H{"post": post}, "更新成功")
}

func (p PostController) Show(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cur := p.collection.FindOne(p.ctx, bson.M{
		"_id": objectId,
	})
	var post models.Post
	err := cur.Decode(&post)
	if err != nil {
		response.Fail(c, nil, "文章不存在")
		return
	}

	response.Success(c, gin.H{"data": post}, "成功")
}

func (p PostController) Edit(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cur := p.collection.FindOne(p.ctx, bson.M{
		"_id": objectId,
	})
	var post models.Post
	err := cur.Decode(&post)
	if err != nil {
		response.Fail(c, nil, "文章不存在")
		return
	}

	response.Success(c, gin.H{"data": post}, "成功")
}

func (p PostController) Delete(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := p.collection.DeleteOne(p.ctx, bson.M{
		"_id": objectId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response.Success(c, gin.H{"post": ""}, "删除成功")
}

func (p PostController) NewBlog(c *gin.Context) {
	var requestPost vo.CreateBlogRequest
	if err := c.ShouldBindJSON(&requestPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tokenString := c.GetHeader("Authorization")
	_, claims, _ := common.ParseToken(tokenString)
	userName := claims.UserName

	var post = models.Post{}
	post.Title = requestPost.Title
	post.Content = requestPost.Content
	post.Tags = requestPost.TagStr
	post.ID = primitive.NewObjectID()
	post.PublishedAt = time.Now()
	post.UserName = userName

	_, err := p.collection.InsertOne(p.ctx, post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new post"})
		return
	}

	log.Println("Remove data from Redis")
	p.redisClient.Del("posts")

	c.JSON(http.StatusOK, post)
}

func (p PostController) PageList(c *gin.Context) {
	// 获取分页参数
	pageNum, _ := strconv.Atoi(c.DefaultQuery("currentPage", ""))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", ""))

	//val, err := p.redisClient.Get("posts").Result()
	//if err == redis.Nil {
	log.Printf("Request to MongoDB")

	findOptions := options.Find()
	// Sort by `price` field descending
	findOptions.SetSort(bson.D{{"publishedAt", -1}})
	cur, err := p.collection.Find(p.ctx, bson.M{}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(p.ctx)

	posts := make([]models.Post, 0)
	for cur.Next(p.ctx) {
		var post models.Post
		cur.Decode(&post)
		posts = append(posts, post)
	}
	total := len(posts)

	//data, _ := json.Marshal(posts)
	//p.redisClient.Set("posts", string(data), 0)
	a := (pageNum - 1) * pageSize
	b := pageNum * pageSize
	if b > total {
		b = total
	}
	response.Success(c, gin.H{"rows": posts[a:b], "total": total}, "成功")
	//} else if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//} else {
	//	log.Printf("Request to Redis")
	//	posts := make([]models.Post, 0)
	//	total := len(posts)
	//	json.Unmarshal([]byte(val), &posts)
	//	a := (pageNum - 1) * pageSize
	//	b := pageNum * pageSize
	//	response.Success(c, gin.H{"data": posts[a:b], "total": total}, "成功")
	//}
}
