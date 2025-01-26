package dataaccess

import (
	"database/sql"
	"fmt"

	"github.com/arnavbajpai/web-forum-project/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

func InsertUser(db *sql.DB, newUser models.User) error {
	query := `INSERT INTO Users (UserAlias, Email, FirstName, Surname, Phone, PasswordHash) 
	VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, newUser.UserAlias, newUser.Email, newUser.FirstName,
		newUser.Surname, newUser.Phone, newUser.PasswordHash)
	if err != nil {
		return err
	}
	return nil
}

func InsertCategory(db *sql.DB, newCategory models.Category) error {
	query := `INSERT INTO Categories (Name, Description) 
	VALUES (?, ?)`
	_, err := db.Exec(query, newCategory.Name, newCategory.Description)
	if err != nil {
		return err
	}
	return nil
}
func FindCategories(db *sql.DB, categoryName string) (models.Category, error) {
	query := "SELECT * FROM Categories WHERE Name = ?"
	row := db.QueryRow(query, categoryName)
	var cat models.Category
	if err := row.Scan(&cat.CategoryID, &cat.Name, &cat.Description, &cat.CreatedAt); err != nil {
		return models.Category{}, err
	}
	return cat, nil
}

func FindUser(db *sql.DB, userAlias string, userID int) (models.User, error) {
	var query string
	var args []interface{}
	if userAlias == "" {
		query = "SELECT * FROM Users WHERE UserID = ?"
		args = append(args, userID)
	} else {
		query = "SELECT * FROM Users WHERE UserAlias = ?"
		args = append(args, userAlias)
	}
	row := db.QueryRow(query, args...)
	var user models.User
	if err := row.Scan(&user.UserID, &user.UserAlias, &user.Email, &user.FirstName,
		&user.Surname, &user.Phone, &user.PasswordHash, &user.StatusId, &user.CreatedAt,
		&user.UpdatedAt, &user.RoleID); err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}

func UpdateUserRecord(db *sql.DB, userID int, userAlias string, firstName string, surname string, phone string) (models.User, error) {
	if userAlias != "" {
		query := "UPDATE Users SET UserAlias = ? WHERE UserID = ?"
		_, err := db.Exec(query, userAlias, userID)
		if err != nil {
			return models.User{}, err
		}
	}
	if firstName != "" {
		query := "UPDATE Users SET FirstName = ? WHERE UserID = ?"
		_, err := db.Exec(query, firstName, userID)
		if err != nil {
			return models.User{}, err
		}
	}
	if surname != "" {
		query := "UPDATE Users SET Surname = ? WHERE UserID = ?"
		_, err := db.Exec(query, surname, userID)
		if err != nil {
			return models.User{}, err
		}
	}
	if phone != "" {
		query := "UPDATE Users SET Phone = ? WHERE UserID = ?"
		_, err := db.Exec(query, phone, userID)
		if err != nil {
			return models.User{}, err
		}
	}
	return FindUser(db, "", userID)
}

func RemoveUser(db *sql.DB, userID int) error {
	query := "UPDATE Users SET StatusId = ? WHERE UserID = ?"
	result, err := db.Exec(query, 2, userID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("Not Found.")
	}
	return nil
}

func InsertPost(db *sql.DB, newPost models.Post, tags []string) error {
	if !newPost.IsTopic { // comment/reply
		postQuery := `INSERT INTO Posts (ParentPostID, UserID, CategoryID, Content)
		VALUES (?, ?, ?, ?)`
		_, err := db.Exec(postQuery, newPost.ParentPostID, newPost.UserID, newPost.CategoryID,
			newPost.Content)
		if err != nil {
			return err
		}
		return nil
	} else { // topic
		postQuery := `INSERT INTO Posts (UserID, CategoryID, Title, Content, IsTopic)
		VALUES (?, ?, ?, ?, ?)`
		result, err := db.Exec(postQuery, newPost.UserID, newPost.CategoryID,
			newPost.Title, newPost.Content, true)
		if err != nil {
			return err
		}
		postID, err := result.LastInsertId()
		if err != nil {
			return err
		}
		for _, tagName := range tags {
			var tagID int
			tagQuery := `SELECT TagID FROM Tags WHERE TagName = ?`
			err := db.QueryRow(tagQuery, tagName).Scan(&tagID)
			if err != nil {
				if err == sql.ErrNoRows {
					insertTagQuery := `INSERT INTO Tags (TagName) VALUES (?)`
					tagResult, err := db.Exec(insertTagQuery, tagName)
					if err != nil {
						return err
					}
					tagID64, err := tagResult.LastInsertId()
					if err != nil {
						return err
					}
					tagID = int(tagID64)
				} else {
					return err
				}
			}
			postTagQuery := `INSERT INTO Post_Tags (PostID, TagID) VALUES (?, ?)`
			_, err = db.Exec(postTagQuery, postID, tagID)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func FindPost(db *sql.DB, categoryId int, tagId int, parentId int, clientID int) ([]models.Post, error) {
	var query string
	var args []interface{}
	if tagId == 0 && parentId == 0 {
		query = `SELECT * FROM Posts 
                 WHERE CategoryID = ? AND IsTopic = TRUE`
		args = append(args, categoryId)
	} else if parentId == 0 {
		query = `SELECT p.* FROM Posts p 
                 JOIN Post_Tags pt ON p.PostID = pt.PostID 
                 WHERE p.CategoryID = ? AND pt.TagID = ? AND p.IsTopic = TRUE`
		args = append(args, categoryId, tagId)
	} else {
		query = `SELECT * FROM Posts 
                 WHERE ParentPostID = ?`
		args = append(args, parentId)
	}
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.PostID, &post.ParentPostID, &post.UserID, &post.CategoryID,
			&post.Title, &post.Content, &post.StatusId, &post.CreatedAt, &post.UpdatedAt, &post.IsTopic)
		if err != nil {
			return nil, err
		}
		if post.UserID == clientID {
			post.Owner = true
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func FindPostByID(db *sql.DB, postID int, clientID int) (models.Post, []models.Tag, error) {
	query := "SELECT * FROM Posts WHERE PostID = ?"
	row := db.QueryRow(query, postID)
	var newPost models.Post
	if err := row.Scan(&newPost.PostID, &newPost.ParentPostID, &newPost.UserID, &newPost.CategoryID,
		&newPost.Title, &newPost.Content, &newPost.StatusId, &newPost.CreatedAt,
		&newPost.UpdatedAt, &newPost.IsTopic); err != nil {
		if err == sql.ErrNoRows {
			return models.Post{}, []models.Tag{}, fmt.Errorf("Post not found")
		}
		return models.Post{}, []models.Tag{}, err
	}
	if newPost.UserID == clientID {
		newPost.Owner = true
	}
	query2 := "SELECT t.* FROM Tags t JOIN Post_Tags pt ON t.TagID = pt.TagID WHERE pt.PostID = ?"
	rows, err := db.Query(query2, postID)
	if err != nil {
		return models.Post{}, []models.Tag{}, err
	}
	defer rows.Close()
	var tags []models.Tag
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.TagID, &tag.TagName, &tag.CreatedAt)
		if err != nil {
			return models.Post{}, []models.Tag{}, err
		}
		tags = append(tags, tag)
	}
	if err = rows.Err(); err != nil {
		return models.Post{}, []models.Tag{}, err
	}
	return newPost, tags, nil

}

func UpdatePostRecord(db *sql.DB, postID int, content string, status string, title string, clientID int) (models.Post, []models.Tag, error) {
	if content != "" {
		query := "UPDATE Posts SET Content = ? WHERE PostID = ?"
		_, err := db.Exec(query, content, postID)
		if err != nil {
			return models.Post{}, []models.Tag{}, err
		}
	}
	if status != "" {
		if status == "Resolved" || status == "Closed" || status == "Complete" {
			query := "UPDATE Posts SET StatusId = ? WHERE PostID = ?"
			_, err := db.Exec(query, 2, postID)
			if err != nil {
				return models.Post{}, []models.Tag{}, err
			}
		}
	}
	if title != "" {
		query := "UPDATE Posts SET Title = ? WHERE PostID = ?"
		_, err := db.Exec(query, title, postID)
		if err != nil {
			return models.Post{}, []models.Tag{}, err
		}
	}
	return FindPostByID(db, postID, clientID)
}

func RemovePost(db *sql.DB, postID int) error {
	query := `DELETE FROM Posts WHERE PostID = ?`
	_, err := db.Exec(query, postID)
	if err != nil {
		return err
	}
	return nil
}
func FindTagByID(db *sql.DB, tagID int) (models.Tag, error) {

	var newTag models.Tag
	query := "SELECT * FROM Tags WHERE TagID = ?"
	row := db.QueryRow(query, tagID)
	if err := row.Scan(&newTag.TagID, &newTag.TagName, &newTag.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return models.Tag{}, fmt.Errorf("Tag not found")
		}
		return models.Tag{}, err
	}
	return newTag, nil
}

func FindTags(db *sql.DB, limit int) ([]models.Tag, error) {

	query := `
		SELECT t.TagID, t.TagName, t.CreatedAt
		FROM Tags t
		JOIN (
			SELECT TagID, COUNT(*) AS usage_count
			FROM Post_Tags
			GROUP BY TagID
			ORDER BY usage_count DESC
			LIMIT ?
		) top_tags
		ON t.TagID = top_tags.TagID
	`

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var tags []models.Tag

	for rows.Next() {
		var tag models.Tag
		if err := rows.Scan(&tag.TagID, &tag.TagName, &tag.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return tags, nil
}

func RemoveTag(db *sql.DB, tagID int) error {
	query := `DELETE FROM Tags WHERE TagID = ?`
	_, err := db.Exec(query, tagID)
	if err != nil {
		return err
	}
	return nil
}
