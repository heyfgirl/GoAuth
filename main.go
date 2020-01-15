package main

import (
	middleware "GoAuth/middlewares"
	router "GoAuth/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	// 要在路由组之前全局使用「跨域中间件」, 否则OPTIONS会返回404
	app.Use(middleware.Cors())
	//获取参数中间件
	app.Use(middleware.Recover())

	router.InitRouter(app)
	app.Run(":8080")

}
func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
	var totnums []int = append(nums1, nums2...)
	var num float64
	totnums2 := []int{}
	for i := 0; i < len(nums1); i++ {
		for t := 0; t < len(nums2); t++ {
			if i {

			}
		}
	}

	return num
}
