package http

import (
	"html/template"
	"net/http"
	"path"
	"strconv"

	"github.com/Toringol/group-golang-1/tree/master/s_shepelev/blog/app/blog"
	"github.com/Toringol/group-golang-1/tree/master/s_shepelev/blog/app/model"
	"github.com/labstack/echo/v4"
	"github.com/russross/blackfriday"
)

var (
	layoutTemplatePosts = template.Must(template.ParseFiles(path.Join("../views", "layout.html"),
		path.Join("../views", "posts.html")))
	layoutTemplateActualPost = template.Must(template.ParseFiles(path.Join("../views", "layout.html"),
		path.Join("../views", "post.html")))
	layoutTemplateAddPost = template.Must(template.ParseFiles(path.Join("../views", "layout.html"),
		path.Join("../views", "newPost.html")))
)

// blogHandlers - http handlers structure
type blogHandlers struct {
	usecase blog.Usecase
}

// NewBlogHandler - deliver our handlers in http
func NewBlogHandler(e *echo.Echo, us blog.Usecase) {
	handlers := blogHandlers{usecase: us}

	// Blog handlers
	e.GET("/", handlers.getAllPostsHandlers)
	e.GET("/posts/:postID", handlers.getPostInfoHandler)
	e.PUT("/posts/:postID", handlers.changePostHandler)
	e.DELETE("/posts/:postID", handlers.deletePostHandler)
	e.GET("/addPost", handlers.getNewPostHandler)
	e.POST("/addPost", handlers.createNewPostHandler)
}

func (bh *blogHandlers) getAllPostsHandlers(ctx echo.Context) error {
	posts, err := bh.usecase.ListPosts()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal DB Error")
	}

	if err := layoutTemplatePosts.ExecuteTemplate(ctx.Response(), "layout", posts); err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, "Internal DB Error")
	}

	return nil
}

func (bh *blogHandlers) getPostInfoHandler(ctx echo.Context) error {
	postID, err := strconv.ParseInt(ctx.Param("postID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal DB Error")
	}

	post, err := bh.usecase.SelectPostByID(postID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal DB Error")
	}

	if err := layoutTemplateActualPost.ExecuteTemplate(ctx.Response(), "layout", post); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal DB Error")
	}

	return nil
}

func (bh *blogHandlers) changePostHandler(ctx echo.Context) error {
	postID, err := strconv.ParseInt(ctx.Param("postID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal DB Error")
	}

	oldPostData, err := bh.usecase.SelectPostByID(postID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal DB Error")
	}

	postData := new(model.Post)
	postData.ID = postID

	if err := ctx.Bind(postData); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal DB Error")
	}

	if postData.Author == "" {
		postData.Author = oldPostData.Author
	}
	if postData.Description == "" {
		postData.Description = oldPostData.Description
	}
	if postData.Title == "" {
		postData.Title = oldPostData.Title
	}

	postData.Description = template.HTML(blackfriday.Run([]byte(postData.Description)))

	_, err = bh.usecase.UpdatePost(postData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal DB Error")
	}

	return ctx.JSON(http.StatusCreated, nil)
}

func (bh *blogHandlers) deletePostHandler(ctx echo.Context) error {
	postID, err := strconv.ParseInt(ctx.Param("postID"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal DB Error")
	}

	_, err = bh.usecase.DeletePost(postID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal DB Error")
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (bh *blogHandlers) getNewPostHandler(ctx echo.Context) error {
	if err := layoutTemplateAddPost.ExecuteTemplate(ctx.Response(), "layout", ""); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal DB Error")
	}

	return nil
}

func (bh *blogHandlers) createNewPostHandler(ctx echo.Context) error {
	postData := new(model.Post)

	if err := ctx.Bind(postData); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal DB Error")
	}

	postData.Description = template.HTML(blackfriday.Run([]byte(postData.Description)))

	_, err := bh.usecase.CreatePost(postData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal DB Error")
	}

	return ctx.JSON(http.StatusCreated, nil)
}
