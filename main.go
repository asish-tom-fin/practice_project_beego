package main

import (
	_ "database/sql"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/buger/goterm"
	_ "github.com/go-sql-driver/mysql"
)

type Blog struct {
	Id        int64  `orm:"column(id);auto" json:"id"`
	Post      string `orm:"column(post);size(2000);null" json:"post"`
	CreatedBy string `orm:"column(created_by);size(200);null" json:"created_by"`
}

const (
	CREATE_BLOG = 1
	LIST_BLOG   = 2
	SEARCH_BLOG = 3
	DELETE_BLOG = 4
	EDIT_BLOG   = 5
	EXIT_APP    = 6
)

func init() {
	orm.RegisterModel(new(Blog))
	name := "default"
	orm.RegisterDataBase("default", "mysql", "root:password@/practice_db?charset=utf8", 30)
	force := false
	verbose := true
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Innit Completed..!")
}

func main() {
	o := orm.NewOrm()
	o.Using("default")
	for {
		if readChoiceAndSwitch(o) {
			break
		}
	}
	return
}

func readChoiceAndSwitch(o orm.Ormer) bool {
	var choice int
	printChoices()
	fmt.Scanf("%d", &choice)
	switch choice {
	case CREATE_BLOG:
		newBlog(o)
		break
	case LIST_BLOG:
		listAllBlogs(o)
		break
	case SEARCH_BLOG:
		searchAndListByID(o)
		break
	case DELETE_BLOG:
		searchAndDeleteByID(o)
		break
	case EDIT_BLOG:
		searchAndEditByID(o)
	case EXIT_APP:
		return true

	}
	return false
}

func printChoices() {
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	fmt.Println("Enter Your choice:")
	fmt.Println("1: Create new blog")
	fmt.Println("2: List all blog")
	fmt.Println(`3: Search and list blog by "ID"`)
	fmt.Println(`4: Search and Delete blog by "ID"`)
	fmt.Println(`5: Search and Edit by blog "ID"`)
	fmt.Println("9: Exit")
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	return
}

func newBlog(o orm.Ormer) {
	var post string
	var createdBy string
	fmt.Println("Enter your name")
	fmt.Scanf("%s", &createdBy)
	fmt.Println(`Enter your post in " ":`)
	fmt.Scanf("%q", &post)
	new_blog := Blog{Post: post, CreatedBy: createdBy}
	fmt.Println("Saving to db")
	fmt.Println(o.Insert(&new_blog))
	fmt.Scanf("%s", &post)
	return
}

func listAllBlogs(o orm.Ormer) {
	qs := o.QueryTable("blog")
	var blogs []orm.Params
	qs.Values(&blogs)
	renderBlogList(blogs)
	return
}

func renderBlogList(blogs []orm.Params) {
	println(`Result:///////////////////////////////////////////////////////`)
	for _, bl := range blogs {
		fmt.Println(bl["Id"], bl["Post"], bl["CreatedBy"])
	}
	println(`/////////////////////////////////////////////////////////////`)
	return

}

func searchAndListByID(o orm.Ormer) {
	fmt.Println("Enter id to search")
	var id int
	fmt.Scanf("%d", &id)
	qs := o.QueryTable("blog").Filter("Id", id)
	var blogs []orm.Params
	count, _ := qs.Values(&blogs)
	if count != 0 {
		renderBlogList(blogs)
	} else {
		fmt.Println("Unable to find the search")
	}
	return
}

func searchAndDeleteByID(o orm.Ormer) {
	fmt.Println("Enter id to search")
	var id int
	fmt.Scanf("%d", &id)
	count, _ := o.QueryTable("blog").Filter("Id", id).Delete()
	if count != 0 {
		qs := o.QueryTable("blog")
		var blogs []orm.Params
		qs.Values(&blogs)
		renderBlogList(blogs)
	} else {
		fmt.Println("Unable to find the search")
	}
	return

}

func searchAndEditByID(o orm.Ormer) {
	fmt.Println("Enter id to search")
	var id int
	var post string
	fmt.Scanf("%d", &id)
	fmt.Println(`Enter updated blog in " ":`)
	fmt.Scanf("%q", &post)
	count, _ := o.QueryTable("blog").Filter("Id", id).Update(orm.Params{
		"Post": post,
	})
	if count != 0 {
		qs := o.QueryTable("blog").Filter("Id", id)
		var blogs []orm.Params
		qs.Values(&blogs)
		renderBlogList(blogs)
	} else {
		fmt.Println("Unable to find the search")
	}
	fmt.Scanf("%s", &post)
	return
}
