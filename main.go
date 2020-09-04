package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
)

type Todo struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Status bool `json:"status"`
}

func initMysql()(err error)  {
	dsn :="root:mysql@tcp(127.0.0.1:3306)/gin_play"
	DB,err = gorm.Open("mysql",dsn)
	if err != nil{
		return
	}
	return DB.DB().Ping()
}

func main() {
	//连接数据库
	err := initMysql()
	if err != nil{
		panic(err)
	}
	defer DB.Close()

	//模型绑定
	DB.AutoMigrate(&Todo{})

	r:= gin.Default()
	//告诉gin静态文件
	r.Static("/static","static")
	//告诉gin找模板文件
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK,"index.html",nil)
	})

	v1Group := r.Group("v1")
	{
		//添加
		v1Group.POST("/todo",func(c *gin.Context) {
			var todo Todo
			c.BindJSON(&todo)
			if err = DB.Create(&todo).Error;err!=nil{
				c.JSON(http.StatusOK,gin.H{"error":err.Error()})
			}else{
				c.JSON(http.StatusOK,gin.H{
					"code":200,
					"msg":"success",
					"data":todo,
				})
			}
		})
		//查看所有
		v1Group.GET("/todo", func(c *gin.Context) {
			var todoList []Todo
			if err =DB.Find(&todoList).Error;err!=nil{
				c.JSON(http.StatusOK,gin.H{"error":err.Error()})
			}else{
				c.JSON(http.StatusOK,todoList)
			}
		})
		//查看一个
		v1Group.GET("/todo/:id", func(c *gin.Context) {

		})
		//修改
		v1Group.PUT("/todo/:id", func(c *gin.Context) {
			id,ok := c.Params.Get("id")
			if !ok{
				c.JSON(http.StatusOK,gin.H{"error":err.Error()})
				return
			}
			var todo Todo
			if err = DB.Where("id=?",id).First(&todo).Error;err != nil{
				c.JSON(http.StatusOK,gin.H{"error":err.Error()})
				return
			}
			c.BindJSON(&todo)
			if err = DB.Save(&todo).Error;err!=nil{
				c.JSON(http.StatusOK,gin.H{"error":err.Error()})
			}else{
				c.JSON(http.StatusOK,todo)
			}
		})
		//删除
		v1Group.DELETE("/todo/:id", func(c *gin.Context) {
			id,ok := c.Params.Get("id")
			if !ok{
				c.JSON(http.StatusOK,gin.H{"error":err.Error()})
				return
			}
			if err = DB.Where("id=?",id).Delete(Todo{}).Error;err!=nil{
				c.JSON(http.StatusOK,gin.H{"error":err.Error()})
			}else{
				c.JSON(http.StatusOK,gin.H{
					"msg":"delete success",
				})
			}
		})
	}
	r.Run()
}

