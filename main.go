package main

import (
	_ "database/sql"
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/buger/goterm"
	_ "github.com/go-sql-driver/mysql"
)

// This struct Blog connects blog to db
type Blog struct {
	ID            int64     `orm:"column(id);auto" json:"id"`
	Post          string    `orm:"column(post);size(2000);null" json:"post"`
	CreatedAt     time.Time `orm:"auto_now_add;column(created_at);type(datetime);null" json:"createdAt"`
	LastUpdatedAt time.Time `orm:"auto_now;column(last_updated_at);type(datetime);null" json:"lastUpdatedAt"`
	CreatedBy     string    `orm:"column(created_by);size(200);null" json:"createdBy"`
}

const (
	// Const to give meaning for switch case choices
	CreateBlog = 1
	ListBlog   = 2
	SearchBlog = 3
	DeleteBlog = 4
	EditBlog   = 5
	ExitApp    = 9
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
	case CreateBlog:
		newBlog(o)
		break
	case ListBlog:
		listAllBlogs(o)
		break
	case SearchBlog:
		searchAndListByID(o)
		break
	case DeleteBlog:
		searchAndDeleteByID(o)
		break
	case EditBlog:
		searchAndEditByID(o)
	case ExitApp:
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
	fmt.Println("Enter your first name")
	fmt.Scanf("%s", &createdBy)
	fmt.Println(`Enter your post in " " (double quotes):`)
	fmt.Scanf("%q", &post)
	newBlog := Blog{Post: post, CreatedBy: createdBy}
	fmt.Println("Saving to db")
	fmt.Println(o.Insert(&newBlog))
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
	println(`Result:////// in the order ID | Post | Name | Created At| Last Updated At`)
	for _, bl := range blogs {
		fmt.Println(bl["ID"], " | ", bl["Post"], " | ", bl["CreatedBy"], " | ",
			bl["CreatedAt"], " | ", bl["LastUpdatedAt"])
	}
	println(`/////////////////////////////////////////////////////////////`)
	return

}

func searchAndListByID(o orm.Ormer) {
	fmt.Println("Enter id to search")
	var id int
	fmt.Scanf("%d", &id)
	qs := o.QueryTable("blog").Filter("ID", id)
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
	count, _ := o.QueryTable("blog").Filter("ID", id).Delete()
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
	fmt.Println(`Enter updated blog in " "(double quotes):`)
	fmt.Scanf("%q", &post)
	count, _ := o.QueryTable("blog").Filter("ID", id).Update(orm.Params{
		"Post": post,
	})
	if count != 0 {
		qs := o.QueryTable("blog").Filter("ID", id)
		var blogs []orm.Params
		qs.Values(&blogs)
		renderBlogList(blogs)
	} else {
		fmt.Println("Unable to find the search")
	}
	fmt.Scanf("%s", &post)
	return
}
