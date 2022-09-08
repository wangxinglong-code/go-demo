package models

import "math"

//分页
type PageableList struct {
	Page      int         `json:"page"`      //当前页
	PageCount int         `json:"pageCount"` //总页数
	List      interface{} `json:"list"`      //列表内容
}

type Limit struct {
	Offset   uint
	LimitNum uint
}

type CommonPageable struct {
	PageCount uint        `json:"page_count"` //总页数
	CurPage   int         `json:"cur_page"`   //当前页
	PageLimit int         `json:"page_limit"` //页条数
	Total     int64       `json:"total"`      //总条数
	List      interface{} `json:"list"`       //列表内容
}

func PageData(curPage, limit int, total int64, data interface{}) (pageable CommonPageable) {
	if 0 != limit {
		pageable.PageCount = uint(math.Ceil(float64(total) / float64(limit)))
	}
	pageable.CurPage = curPage
	pageable.PageLimit = limit
	pageable.Total = total
	pageable.List = data
	return
}
