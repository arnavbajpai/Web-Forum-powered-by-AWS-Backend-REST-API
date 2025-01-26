package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/arnavbajpai/web-forum-project/internal/dataaccess"
	"github.com/arnavbajpai/web-forum-project/internal/database"
	"github.com/arnavbajpai/web-forum-project/internal/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Payload struct {
	Meta json.RawMessage `json:"meta,omitempty"`
	Data json.RawMessage `json:"data,omitempty"`
}

type Response struct {
	Payload   Payload  `json:"payload"`
	Messages  []string `json:"messages"`
	ErrorCode int      `json:"errorCode"`
}

func WrapResponse(data interface{}, messages []string) (Response, error) {
	payloadData, err := json.Marshal(data)
	if err != nil {
		return Response{}, err
	}
	return Response{
		Payload:  Payload{Data: payloadData},
		Messages: messages,
	}, nil
}

func ValidateToken(c *gin.Context) int {
	email, exists := c.Get("userEmail")
	if !exists {
		fmt.Println("Error 1 email = ", email)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email from context"})
		return 0
	}
	emailStr, ok := email.(string)
	if !ok {
		fmt.Println("Error 2 email str = ", emailStr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid email format in context"})
		return 0
	}
	username := strings.Split(emailStr, "@")[0]
	fmt.Println(username)
	client, err := dataaccess.FindUser(database.DBCon, username, 0)
	if err == nil {
		return client.UserID
	} else {
		fmt.Println(err.Error())
		firstName, _ := c.Get("userGivenName")
		firstNameStr, ok := firstName.(string)
		if !ok {
			fmt.Println("Error 3 firstName =", firstName)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid format in context"})
			return 0
		}
		surname, _ := c.Get("userFamilyName")
		surnameStr, ok := surname.(string)
		if !ok {
			fmt.Println("Error 4 surname =", surname)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid format in context"})
			return 0
		}
		query := "INSERT INTO Users (UserAlias, Email, FirstName , Surname) VALUES (?, ?, ?, ?)"
		res, err := database.DBCon.Exec(query, username, email, firstNameStr, surnameStr)
		if err != nil {
			fmt.Println("Error 5 : ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid format in context"})
			return 0
		}
		newUserID, err := res.LastInsertId()
		if err != nil {
			fmt.Println("Error 6 : ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid format in context"})
			return 0
		}
		return int(newUserID)
	}
}

func AddCategory(c *gin.Context) {
	var newCategory models.Category
	if err := c.BindJSON(&newCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := dataaccess.InsertCategory(database.DBCon, newCategory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user"})
		return
	}
	c.IndentedJSON(http.StatusOK, nil)
}
func GetCategories(c *gin.Context) {
	categoryName := c.Query("name")
	category, err := dataaccess.FindCategories(database.DBCon, categoryName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No categories found")})
		return
	}
	resp, err := WrapResponse(category, []string{"Category fetched successfully"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create response"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func AddUser(c *gin.Context) {
	var newUser models.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := dataaccess.FindUser(database.DBCon, newUser.UserAlias, 0); err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": fmt.Sprintf("User with username %s already exists", newUser.UserAlias)})
		return
	}
	if err := dataaccess.InsertUser(database.DBCon, newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user"})
		return
	}
	c.IndentedJSON(http.StatusOK, nil)
}
func GetUserByID(c *gin.Context) {

	userID, _ := strconv.Atoi(c.Param("userID"))

	user, err := dataaccess.FindUser(database.DBCon, "", userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("User with userID %d not found", userID)})
		return
	}
	resp, err := WrapResponse(user, []string{"User fetched successfully"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create response"})
		return
	}
	c.IndentedJSON(http.StatusOK, resp)
}

func GetUserByAlias(c *gin.Context) {

	userAlias := c.Query("userAlias")
	if userAlias == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid username provided."})
		return
	}
	user, err := dataaccess.FindUser(database.DBCon, userAlias, 0)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("User with username %s not found", userAlias)})
		return
	}

	resp, err := WrapResponse(user, []string{"User fetched successfully"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create response"})
		return
	}
	fmt.Println(user)
	c.IndentedJSON(http.StatusOK, resp)
}

func UpdateUser(c *gin.Context) {
	clientID := ValidateToken(c)
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId path parameter is required"})
		return
	}
	if clientID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Anauthorised access. Client with id %d tried to update UserID %d", clientID, userID)})
		return
	}
	user, err := dataaccess.UpdateUserRecord(database.DBCon, userID, c.Query("userAlias"), c.Query("firstName"), c.Query("surname"), c.Query("phone"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	resp, err := WrapResponse(user, []string{"User fetched successfully"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create response"})
		return
	}
	c.IndentedJSON(http.StatusOK, resp)
}

func DeleteUser(c *gin.Context) {
	clientID := ValidateToken(c)
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId path parameter is required"})
		return
	}
	if clientID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Anauthorised access. Client with id %d tried to delete UserID %d", clientID, userID)})
		return
	}

	if err := dataaccess.RemoveUser(database.DBCon, userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func AddPost(c *gin.Context) {
	var request models.AddPostRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tags := request.Tags
	newPost := request.Post
	fmt.Println(newPost.UserID)
	if poster, err := dataaccess.FindUser(database.DBCon, "", newPost.UserID); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": fmt.Sprintf("User %s does not exist.", poster.UserAlias)})
		return
	}
	if err := dataaccess.InsertPost(database.DBCon, newPost, tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add Post"})
		return
	}
	c.IndentedJSON(http.StatusOK, nil)
}

func GetPosts(c *gin.Context) {
	clientID := ValidateToken(c)
	categoryID, _ := strconv.Atoi(c.Query("categoryID"))
	tagID, _ := strconv.Atoi(c.Query("tagID"))
	parentId, _ := strconv.Atoi(c.Query("parentID"))
	posts, err := dataaccess.FindPost(database.DBCon, categoryID, tagID, parentId, clientID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No posts found")})
		return
	}
	resp, err := WrapResponse(posts, []string{"Posts fetched successfully"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create response"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
func GetPostByID(c *gin.Context) {
	clientID := ValidateToken(c)
	postID, _ := strconv.Atoi(c.Param("postID"))
	post, tags, err := dataaccess.FindPostByID(database.DBCon, postID, clientID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No posts found")})
		return
	}
	responseData := models.PostWithTagsResponse{
		Post: post,
		Tags: tags,
	}
	resp, err := WrapResponse(responseData, []string{"Post and associated tags fetched successfully"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create response"})
		return
	}
	c.JSON(http.StatusOK, resp)
}
func UpdatePost(c *gin.Context) {
	clientID := ValidateToken(c)
	postID, err := strconv.Atoi(c.Param("postID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "postId path parameter is required"})
		return
	}
	targetPost, _, err := dataaccess.FindPostByID(database.DBCon, postID, clientID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No posts found")})
		return
	}
	if targetPost.UserID != clientID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Anauthorised access. Client with id %d tried to update post with author %d", clientID, targetPost.UserID)})
		return
	}
	post, tags, err := dataaccess.UpdatePostRecord(database.DBCon, postID, c.Query("content"), c.Query("statusName"), c.Query("title"), clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		fmt.Println(err.Error())
		return
	}
	responseData := models.PostWithTagsResponse{
		Post: post,
		Tags: tags,
	}
	resp, err := WrapResponse(responseData, []string{"Post updated successfully"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create response"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func DeletePost(c *gin.Context) {
	clientID := ValidateToken(c)
	postID, err := strconv.Atoi(c.Param("postID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "postId path parameter is required"})
		return
	}
	targetPost, _, err := dataaccess.FindPostByID(database.DBCon, postID, clientID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No posts found")})
		return
	}
	if targetPost.UserID != clientID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Anauthorised access. Client with id %d tried to delete post with author %d", clientID, targetPost.UserID)})
		return
	}
	if err := dataaccess.RemovePost(database.DBCon, postID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete post"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func GetTag(c *gin.Context) {
	tagID, err := strconv.Atoi(c.Param("tagID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "postId path parameter is required"})
		return
	}
	newTag, err := dataaccess.FindTagByID(database.DBCon, tagID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No tag found.")})
		return
	}
	resp, err := WrapResponse(newTag, []string{"Tag fetched successfully"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create response"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func GetTopTags(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Param("limit"))
	tags, err := dataaccess.FindTags(database.DBCon, limit)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No tags found")})
		return
	}
	resp, err := WrapResponse(tags, []string{"tags fetched successfully"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create response"})
		return
	}
	c.JSON(http.StatusOK, resp)

}

func DeleteTag(c *gin.Context) { // CRON JOB
	tagID, err := strconv.Atoi(c.Param("tagID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tagId path parameter is required"})
		return
	}
	if err := dataaccess.RemoveTag(database.DBCon, tagID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to delete tag"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
