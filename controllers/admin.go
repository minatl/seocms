package controllers

import (
    //"github.com/astaxie/beedb"
    "github.com/astaxie/beego"
    "time"
)

type AdminController struct {
    beego.Controller
}

var err error

func (this *AdminController) Get() {
    this.Layout = "layout_admin.tpl"
    object := this.Ctx.Params[":object"]
    action := this.Ctx.Params[":action"]
    if object == "article" {
        switch action {
        case "list":
            this.Data["PageTitle"] = "文章列表_文章管理_SEOCMS"
            articles := []Article{}
            orm := InitDb()
            err = orm.OrderBy("-pubdate").FindAll(&articles)
            Check(err)
            this.Data["Articles"] = articles
            this.TplNames = "admin/article_list.tpl"
        case "add":
            this.Data["PageTitle"] = "添加文章_文章管理_SEOCMS"

            //this.Data["Categories"] = []string{"博客", "笔记"}    // 测试数据
            orm := InitDb()
            categories := []Category{}
            err = orm.OrderBy("name").FindAll(&categories)
            if err != nil {
                this.Data["Categories"] = []string{}
            } else {
                categoryList := []string{}
                for _, category := range(categories) {
                    categoryList = append(categoryList, category.Name)
                }
                this.Data["Categories"] = categoryList
            }

            this.TplNames = "admin/add_article.tpl"
        case "edit":
            //this.Data["Id"] = this.Ctx.Params[":id"]
            id := this.Ctx.Params[":id"]

            orm := InitDb()
            article := Article{}
            err = orm.Where("id=?", id).Find(&article)
            Check(err)
            this.Data["Article"] = article

            orm = InitDb()
            categories := []Category{}
            err = orm.OrderBy("name").FindAll(&categories)
            if err != nil {
                this.Data["Categories"] = []string{}
            } else {
                categoryList := []string{}
                for _, category := range(categories) {
                    categoryList = append(categoryList, category.Name)
                }
                this.Data["Categories"] = categoryList
            }

            this.Data["Pubdate"] = article.Pubdate.Format("2006-01-02")

            this.Data["PageTitle"] = "编辑文章_文章管理_SEOCMS"
            this.TplNames = "admin/edit_article.tpl"
        case "delete":
            this.Data["Id"] = this.Ctx.Params[":id"]
            this.Data["PageTitle"] = "删除文章_文章管理_SEOCMS"
            this.TplNames = "admin/del_article.tpl"
        }
    } else if object == "category" {
        switch action {
        case "list":    // 分类列表
            categories := []Category{}
            orm := InitDb()
            err = orm.OrderBy("name").FindAll(&categories)
            Check(err)
            this.Data["PageTitle"] = "分类列表_文章管理_SEOCMS"
            this.Data["Categories"] = categories
            this.TplNames = "admin/category_list.tpl"
        case "add":    // 添加分类
            this.Data["PageTitle"] = "添加分类_文章管理_SEOCMS"
            this.TplNames = "admin/add_category.tpl"
        case "edit":    // 修改分类
            id := this.Ctx.Params[":id"]
            //this.Data["Id"] = id

            orm := InitDb()
            category := Category{}
            err = orm.Where("id=?", id).Find(&category)
            Check(err)
            this.Data["Category"] = category

            this.Data["PageTitle"] = "修改分类_文章管理_SEOCMS"
            this.TplNames = "admin/edit_category.tpl"
        case "delete":    // 删除分类
            //this.Data["Id"] = this.Ctx.Params[":id"]
            //this.Data["PageTitle"] = "删除分类_文章管理_SEOCMS"
            //this.TplNames = "admin/delete_category.tpl"
            id := this.Ctx.Params[":id"]

            orm := InitDb()
            category := Category{}
            err = orm.Where("id=?", id).Find(&category)
            Check(err)
            orm.Delete(&category)

            this.Ctx.Redirect(301, "/category/list")
        }
    }
}

func (this *AdminController) Post() {
    object := this.Ctx.Params[":object"]
    action := this.Ctx.Params[":action"]
    this.Layout = "layout_admin.tpl"
    if object == "article" {
        switch action {
        case "add":    // 处理添加文章
            title := this.Input().Get("title")
            pubdate := this.Input().Get("pubdate")
            abstract := this.Input().Get("abstract")
            content := this.Input().Get("content")
            category := this.Input().Get("category")
            tags := this.Input().Get("tags")
            Debug("%s, %s, %s, %s, %s, %s", title, pubdate, abstract, content, category, tags)

            orm := InitDb()
            categoryObj := Category{}
            err = orm.Where("name=?", category).Find(&categoryObj)
            Check(err)
            categoryId := categoryObj.Id
            Debug("Category id is `%d`.", categoryId)

            article := Article{}
            article.Title = title
            article.Pubdate, _ = Str2date(pubdate)
            article.Updated = time.Now()
            article.Abstract = abstract
            article.AbstractHtml = Markdown2html(abstract)
            article.Content = content
            article.ContentHtml = Markdown2html(content)
            article.Category = categoryId
            Debug("Abstract is `%s`.", article.AbstractHtml)
            orm = InitDb()
            err = orm.Save(&article)
            Check(err)

        case "edit":    // 处理修改文章
            id := this.Ctx.Params[":id"]    // 文章ID

            orm := InitDb()
            article := Article{}
            err = orm.Where("id=?", id).Find(&article)
            Check(err)

            article.Title = this.Input().Get("title")
            pubdate := this.Input().Get("pubdate")
            article.Abstract = this.Input().Get("abstract")
            article.Content = this.Input().Get("content")
            category := this.Input().Get("category")
            tags := this.Input().Get("tags")

            this.Data["Article"] = article
            this.Data["Pubdate"] = pubdate
            this.Data["Category"] = category
            this.Data["Tags"] = tags

            // 检查文章标题或发布日期是否为空
            if article.Title == "" || pubdate == "" {
                this.Data["Message"] = "文章标题或发布日期不能为空。"
                this.TplNames = "admin/edit_article.tpl"
                return
            }

            // 检查发布日期格式是否合法
            article.Pubdate, err = Str2date(pubdate)
            if err != nil {
                this.Data["Message"] = "发布日期格式不正确，请采用`2006-01-02`这样的格式。"
                this.TplNames = "admin/edit_article.tpl"
                return
            }

            orm = InitDb()
            categoryObj := Category{}
            err = orm.Where("name=?", category).Find(&categoryObj)
            Check(err)
            categoryId := categoryObj.Id

            article.Updated = time.Now()
            article.AbstractHtml = Markdown2html(article.Abstract)
            article.ContentHtml = Markdown2html(article.Content)
            article.Category = categoryId
            orm = InitDb()
            err = orm.Save(&article)
            Check(err)
        }
        this.Ctx.Redirect(302, "/article/list")
    } else if object == "category" {
        switch action {
        case "add":    // 处理添加分类
            name := this.Input().Get("name")
            nameEn := this.Input().Get("name_en")
            description := this.Input().Get("description")
            alias := this.Input().Get("alias")

            this.Data["Name"] = name
            this.Data["NameEn"] = nameEn
            this.Data["Description"] = description
            this.Data["Alias"] = alias

            // 检查分类名称或分类英文名称是否为空
            if name == "" || nameEn == "" {
                this.Data["Message"] = "分类名称或分类英文名称不能为空。"
                this.TplNames = "admin/add_category.tpl"
                return
            }

            // 检查分类名称或分类英文名称是否已存在
            orm := InitDb()
            category := Category{}
            err = orm.Where("name=? or name_en=?", name, nameEn).Find(&category)
            if err != nil {
                orm = InitDb()
            } else {
                this.Data["Message"] = "分类名称或分类英文名称已存在。"
                this.TplNames = "admin/add_category.tpl"
                return
            }

            category.Name = name
            category.NameEn = nameEn
            category.Description = description
            category.Alias = alias
            //category.Alias = Str2slice(alias)
            //Debug("%s, %s, %s, %v", name, nameEn, description, alias)
            Debug("%s, %s, %s, %v", category.Name, category.NameEn, category.Description, category.Alias)
            err = orm.Save(&category)
            Check(err)

            this.Ctx.Redirect(302, "/category/list")
        case "edit":    // 处理修改分类
            id := this.Ctx.Params[":id"]    // 分类ID

            orm := InitDb()
            category := Category{}
            err = orm.Where("id=?", id).Find(&category)
            Check(err)

            category.Name = this.Input().Get("name")
            category.NameEn = this.Input().Get("name_en")
            category.Description = this.Input().Get("description")
            category.Alias = this.Input().Get("alias")

            this.Data["Category"] = category

            // 检查分类名称或分类英文名称是否为空
            if category.Name == "" || category.NameEn == "" {
                this.Data["Message"] = "分类名称或分类英文名称已存在。"
                this.TplNames = "admin/edit_category.tpl"
                return
            }

            // 检查分类名称或分类英文名称是否已存在
            err = orm.Where("id!=? and (name=? or name_en=?)", id, category.Name, category.NameEn).Find(&category)
            if err != nil {
                orm = InitDb()
            } else {
                this.Data["Message"] = "分类名称或分类英文名称已存在"
                this.TplNames = "admin/edit_category.tpl"
                return
            }

            err = orm.Save(&category)
            Check(err)

            this.Ctx.Redirect(302, "/category/list")
        }
    }
}
